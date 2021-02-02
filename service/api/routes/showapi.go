package routes

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ProfilesRoutes getting all routes for the profile endpoint
func ShowRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetUIProfilesEndpoint)
	router.Get("/{profileName}", GetUIProfileEndpoint)
	return router
}

// GetUIProfilesEndpoint getting all profile names
func GetUIProfilesEndpoint(response http.ResponseWriter, request *http.Request) {
	// user := getUsername(request)
	// if user == "" {
	// 	msg := fmt.Sprintf("user header %s missing", api.UserHeader)
	// 	api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
	// 	return
	// }
	var profileNames []string
	profileNames = make([]string, 0)
	for _, profile := range config.Profiles {
		profileNames = append(profileNames, profile.Name)
	}
	render.JSON(response, request, profileNames)
}

// GetUIProfileEndpoint getting all profile names
func GetUIProfileEndpoint(response http.ResponseWriter, request *http.Request) {
	// user := getUsername(request)
	// if user == "" {
	// 	msg := fmt.Sprintf("user header %s missing", api.UserHeader)
	// 	api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
	// 	return
	// }

	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	for _, profile := range config.Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			// making a deep copy of this profile
			uiProfile := profile.Copy()
			// for every action deleting the commands
			for index, _ := range uiProfile.Actions {
				uiProfile.Actions[index].Commands = make([]models.Command, 0)
			}
			render.JSON(response, request, uiProfile)
			return
		}
	}
	return
}
