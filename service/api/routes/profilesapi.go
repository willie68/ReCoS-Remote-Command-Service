package routes

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

// ProfilesRoutes getting all routes for the profile endpoint
func ProfilesRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetProfilesEndpoint)
	router.Get("/{profileName}", GetProfileEndpoint)
	return router
}

// GetProfilesEndpoint getting all profile names
func GetProfilesEndpoint(response http.ResponseWriter, request *http.Request) {
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

// GetProfileEndpoint getting all profile names
func GetProfileEndpoint(response http.ResponseWriter, request *http.Request) {
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
			render.JSON(response, request, profile)
			return
		}
	}
}