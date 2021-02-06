package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

// ActionRoutes an action for a defined profile
func ActionRoutes() *chi.Mux {
	router := chi.NewRouter()
	//	router.Get("/{profile}", GetProfileActionsEndpoint)
	router.Post("/{profileName}/{actionName}", PostProfileActionEndpoint)
	return router
}

// PostProfileActionEndpoint getting all profile names
func PostProfileActionEndpoint(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	actionName, err := api.Param(request, "actionName")
	if err != nil {
		clog.Logger.Debug("Error reading action name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	//	var action models.Action
	for _, profile := range config.Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			// for every action creating the info object
			for _, action := range profile.Actions {
				log.Printf("%v\r\n", action)
			}
			render.JSON(response, request, actionName)
			return
		}
	}
	return
}
