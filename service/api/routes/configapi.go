package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/config"
	"wkla.no-ip.biz/remote-desk-service/dto"
)

var icons []string

/*
ConfigRoutes getting all routes for the config endpoint
*/
func ConfigRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/icons", GetIcons)
	router.Get("/commands", GetCommands)
	return router
}

/*
GetIcons list of all possible icon names
*/
func GetIcons(response http.ResponseWriter, request *http.Request) {
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
		panic(err)
	}

	render.JSON(response, request, icons)
}

/*
GetCommands list of all possible command types
*/
func GetCommands(response http.ResponseWriter, request *http.Request) {
	render.JSON(response, request, dto.CommandTypes)
}
