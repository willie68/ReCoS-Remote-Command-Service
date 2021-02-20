package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/dto"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
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
	clog.Logger.Debugf("Action: %s:%s", profileName, actionName)
	decoder := json.NewDecoder(request.Body)
	var message models.Message
	err = decoder.Decode(&message)
	if err != nil {
		clog.Logger.Debug("Error reading json body:" + err.Error())
		api.Err(response, request, err)
		return
	}

	ok, err := dto.Execute(profileName, actionName, message)
	if err != nil {
		clog.Logger.Debug("Error reading action name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	render.JSON(response, request, ok)
	return
}
