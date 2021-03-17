package routes

import (
	"net/http"
	"os"
	"path/filepath"
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
)

var icons []string
var mu sync.Mutex

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
	if icons == nil || len(icons) == 0 {

		icons = make([]string, 0)

		err := filepath.Walk(config.Get().Icons, func(path string, info os.FileInfo, err error) error {
			pathNames := strings.Split(path, "/")
			if len(pathNames) == 1 {
				pathNames = strings.Split(path, "\\")
			}
			icon := pathNames[len(pathNames)-1]
			if strings.HasSuffix(icon, ".png") {
				icons = append(icons, icon)
			}
			return nil
		})
		if err != nil {
			clog.Logger.Debug("Error reading icon files:" + err.Error())
			api.Err(response, request, err)
			return
		}
		sort.Slice(icons, func(i, j int) bool { return strings.ToLower(icons[i]) < strings.ToLower(icons[j]) })
	}
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
