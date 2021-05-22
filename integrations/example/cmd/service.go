package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"wkla.no-ip.biz/pl-example/api"
	"wkla.no-ip.biz/pl-example/api/handler"
	"wkla.no-ip.biz/pl-example/api/routes"
	"wkla.no-ip.biz/pl-example/error/serror"
	"wkla.no-ip.biz/pl-example/health"
	"wkla.no-ip.biz/pl-example/logging"

	config "wkla.no-ip.biz/pl-example/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/pl-example/crypt"
	clog "wkla.no-ip.biz/pl-example/logging"

	flag "github.com/spf13/pflag"
)

/*
apVersion implementing api version for this service
*/
const apiVersion = "1"
const servicename = "pl_example"

var port int
var sslport int
var serviceURL string
var ssl bool
var configFile string
var serviceConfig config.Config
var sslsrv *http.Server
var srv *http.Server

func init() {
	// variables for parameter override
	ssl = false
	clog.Logger.Info("init service")
	flag.IntVarP(&port, "port", "p", 0, "port of the http server.")
	flag.IntVarP(&sslport, "sslport", "t", 0, "port of the https server.")
	flag.StringVarP(&configFile, "config", "c", "", "this is the path and filename to the config file")
	flag.StringVarP(&serviceURL, "serviceURL", "u", "", "service url from outside")
}

func apiRoutes() *chi.Mux {
	baseURL := fmt.Sprintf("/api/v%s", apiVersion)
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		//middleware.DefaultCompress,
		middleware.Recoverer,
		cors.Handler(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-mcs-username", "X-mcs-password", "X-mcs-profile"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)

	router.Route("/", func(r chi.Router) {
		r.Mount(baseURL+"/action", routes.ActionRoutes())
		r.Mount("/health", health.Routes())
	})

	return router
}

func healthRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		//middleware.DefaultCompress,
		middleware.Recoverer,
	)

	router.Route("/", func(r chi.Router) {
		r.Mount("/", health.Routes())
	})
	return router
}

func main() {
	flag.Parse()

	clog.Logger.Info("starting server")
	defer clog.Logger.Close()

	serror.Service = servicename
	if configFile == "" {
		configFolder, err := config.GetDefaultConfigFolder()
		if err != nil {
			clog.Logger.Alertf("can't load config file: %s", err.Error())
			os.Exit(1)
		}
		configFolder = fmt.Sprintf("%s/service/", configFolder)
		err = os.MkdirAll(configFolder, os.ModePerm)
		if err != nil {
			clog.Logger.Alertf("can't load config file: %s", err.Error())
			os.Exit(1)
		}
		configFile = configFolder + "/service.yaml"

		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			config.SaveConfig(configFile, config.DefaulConfig, false)
		}
	}

	config.File = configFile

	if err := config.Load(); err != nil {
		clog.Logger.Alertf("can't load config file: %s", err.Error())
		os.Exit(1)
	}

	serviceConfig = config.Get()
	initConfig()

	clog.Logger.Info("service is starting")

	healthCheckConfig := health.CheckConfig(serviceConfig.HealthCheck)
	health.InitHealthSystem(healthCheckConfig)

	if serviceConfig.Sslport > 0 {
		ssl = true
		clog.Logger.Info("ssl active")
	}

	clog.Logger.Infof("ssl: %t", ssl)
	clog.Logger.Infof("serviceURL: %s", serviceConfig.ServiceURL)
	clog.Logger.Infof("%s api routes", servicename)
	router := apiRoutes()
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		clog.Logger.Infof("%s %s", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		clog.Logger.Alertf("could not walk api routes. %s", err.Error())
	}
	clog.Logger.Info("health api routes")
	healthRouter := healthRoutes()
	if err := chi.Walk(healthRouter, walkFunc); err != nil {
		clog.Logger.Alertf("could not walk health routes. %s", err.Error())
	}

	router.HandleFunc("/ws", api.ServeWs)
	if ssl {
		gc := crypt.GenerateCertificate{
			Organization: "MCS",
			Host:         "127.0.0.1",
			ValidFor:     10 * 365 * 24 * time.Hour,
			IsCA:         false,
			EcdsaCurve:   "P384",
			Ed25519Key:   false,
		}
		tlsConfig, err := gc.GenerateTLSConfig()
		if err != nil {
			clog.Logger.Alertf("could not create tls config. %s", err.Error())
		}

		sslsrv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Sslport),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
			TLSConfig:    tlsConfig,
		}
		go func() {
			clog.Logger.Infof("starting https server on address: %s", sslsrv.Addr)
			if err := sslsrv.ListenAndServeTLS("", ""); err != nil {
				clog.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      healthRouter,
		}
		go func() {
			clog.Logger.Infof("starting http server on address: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				clog.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
	} else {
		// own http server for the healthchecks
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
		}
		go func() {
			clog.Logger.Infof("starting http server on address: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				clog.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
	}

	clog.Logger.Info("waiting for clients")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	if ssl {
		sslsrv.Shutdown(ctx)
	}

	clog.Logger.Info("finished")

	os.Exit(0)
}

func initConfig() {
	if port > 0 {
		serviceConfig.Port = port
	}
	if sslport > 0 {
		serviceConfig.Sslport = sslport
	}
	if serviceURL != "" {
		serviceConfig.ServiceURL = serviceURL
	}

	handler.AuthenticationConfig.Password = serviceConfig.Password

	var err error
	logging.Logger.SetLevel(serviceConfig.Logging.Level)
	serviceConfig.Logging.Filename, err = config.ReplaceConfigdir(serviceConfig.Logging.Filename)
	if err != nil {
		clog.Logger.Alertf("error wrong logging folder: %s", err.Error())
		os.Exit(1)
	}

	logging.Logger.Filename = serviceConfig.Logging.Filename
	logging.Logger.InitGelf()
}

func getApikey() string {
	value := fmt.Sprintf("%s_%s", servicename, "default")
	apikey := fmt.Sprintf("%x", md5.Sum([]byte(value)))
	return strings.ToLower(apikey)
}
