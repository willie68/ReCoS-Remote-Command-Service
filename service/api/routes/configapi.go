package routes

import (
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/api/handler"
	"wkla.no-ip.biz/remote-desk-service/config"
	"wkla.no-ip.biz/remote-desk-service/dto"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/web"
)

var icons []string
var mu sync.Mutex
var getIconsDo sync.Once

/*
ConfigRoutes getting all routes for the config endpoint
*/
func ConfigRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/icons", GetIcons)
	router.Get("/commands", GetCommands)
	router.With(handler.AuthCheck()).Get("/check", GetCheck)
	return router
}

/*
GetIcons list of all possible icon names
*/
func GetIcons(response http.ResponseWriter, request *http.Request) {
	mu.Lock()
	getIconsDo.Do(func() {
		icons = make([]string, 0)
		files, err := web.WebClientAssets.ReadDir("webclient/assets")
		if err != nil {
			clog.Logger.Debug("Error reading icon files:" + err.Error())
			api.Err(response, request, err)
			return
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".png") {
				icons = append(icons, file.Name())
			}
		}

		sort.Slice(icons, func(i, j int) bool { return strings.ToLower(icons[i]) < strings.ToLower(icons[j]) })
	})
	mu.Unlock()
	render.JSON(response, request, icons)
}

/*
GetCommands list of all possible command types
*/
func GetCommands(response http.ResponseWriter, request *http.Request) {
	profileName := request.Header.Get("X-mcs-profile")
	var profile models.Profile
	if profileName != "" {
		profile, _ = config.GetProfile(profileName)
	}
	types := dto.CommandTypes
	types, err := dto.EnrichTypes(types, profile)
	if err != nil {
		clog.Logger.Debug("Error reading icon files:" + err.Error())
		api.Err(response, request, err)
		return
	}
	render.JSON(response, request, types)
}

/*
GetCheck simply checks the authentication
*/
func GetCheck(response http.ResponseWriter, request *http.Request) {
	render.JSON(response, request, "ok")
}
