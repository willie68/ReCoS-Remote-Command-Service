package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/config"
	"wkla.no-ip.biz/remote-desk-service/dto"
	"wkla.no-ip.biz/remote-desk-service/error/serror"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ProfilesRoutes getting all routes for the profile endpoint
func ProfilesRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetProfiles)
	router.Post("/", PostProfile)
	router.Get("/{profileName}", GetProfile)
	return router
}

// GetProfiles getting all profile names
func GetProfiles(response http.ResponseWriter, request *http.Request) {
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

// GetProfile getting a profile
func GetProfile(response http.ResponseWriter, request *http.Request) {
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

// PostProfile create a new profile
func PostProfile(response http.ResponseWriter, request *http.Request) {
	user := getUsername(request)
	if user == "" {
		msg := fmt.Sprintf("user header %s missing", api.UserHeader)
		api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
		return
	}

	decoder := json.NewDecoder(request.Body)
	var profile models.Profile
	err := decoder.Decode(&profile)
	if err != nil {
		clog.Logger.Debug("Error reading json body:" + err.Error())
		api.Err(response, request, err)
		return
	}

	if config.HasProfile(profile.Name) {
		msg := fmt.Sprintf("profile already exists: %s", profile.Name)
		api.Err(response, request, serror.BadRequest(nil, "profile-exists", msg))
		return
	}

	err = config.SaveProfile(profile)
	if err != nil {
		clog.Logger.Debug("Error saving profile:" + err.Error())
		api.Err(response, request, err)
		return
	}

	go func() {
		if err := dto.ReinitProfiles(config.Profiles); err != nil {
			clog.Logger.Alertf("can't create profiles: %s", err.Error())
		}
	}()

	render.Status(request, http.StatusCreated)
	render.JSON(response, request, profile)
}
