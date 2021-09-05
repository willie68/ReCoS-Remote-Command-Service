package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"wkla.no-ip.biz/remote-desk-service/pkg/transformer"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/api/handler"
	"wkla.no-ip.biz/remote-desk-service/config"
	"wkla.no-ip.biz/remote-desk-service/error/serror"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pac"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ProfilesRoutes getting all routes for the profile endpoint
func ProfilesRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetProfiles)
	router.With(handler.AuthCheck()).Post("/", PostProfile)
	router.With(handler.AuthCheck()).Put("/{profileName}", PutProfile)
	router.Get("/{profileName}", GetProfile)
	router.With(handler.AuthCheck()).Delete("/{profileName}", DeleteProfile)
	router.Get("/{profileName}/export", GetExportProfile)
	router.Get("/{profileName}/actions/{actionName}/export", GetExportAction)
	router.Get("/{profileName}/pages/{pageName}/export", GetExportPage)
	router.With(handler.AuthCheck()).Post("/{profileName}/combine/", PostCombine)
	return router
}

// GetProfiles getting all profile names
func GetProfiles(response http.ResponseWriter, request *http.Request) {
	var profileNames []string
	profileNames = make([]string, 0)
	for _, profile := range config.Profiles {
		profileNames = append(profileNames, profile.Name)
	}

	sort.Slice(profileNames, func(i, j int) bool { return strings.ToLower(profileNames[i]) < strings.ToLower(profileNames[j]) })

	render.JSON(response, request, profileNames)
}

// GetProfile getting a profile
func GetProfile(response http.ResponseWriter, request *http.Request) {
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

	err = config.SaveProfileFile(profile)
	if err != nil {
		clog.Logger.Debug("Error saving profile:" + err.Error())
		api.Err(response, request, err)
		return
	}

	go func() {
		config.SaveProfileFile(profile)
		config.AddProfile(profile)
		if _, err := pac.InitProfile(profile.Name); err != nil {
			clog.Logger.Alertf("can't create profiles: %s", err.Error())
		}
	}()

	render.Status(request, http.StatusCreated)
	render.JSON(response, request, profile)
}

// PutProfile create a new profile
func PutProfile(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var profile models.Profile
	err = decoder.Decode(&profile)
	if err != nil {
		clog.Logger.Debug("Error reading json body:" + err.Error())
		api.Err(response, request, err)
		return
	}

	if profileName != profile.Name {
		msg := fmt.Sprintf("profile names not equal: %s != %s", profileName, profile.Name)
		clog.Logger.Debug(msg)
		api.Err(response, request, serror.BadRequest(nil, "profile-name", msg))
		return
	}

	err = config.UpdateProfileFile(profile)
	if err != nil {
		clog.Logger.Debug("Error saving profile:" + err.Error())
		api.Err(response, request, err)
		return
	}

	go func() {
		pac.CloseProfile(profile.Name)
		pac.RemoveProfile(profile.Name)
		config.UpdateProfile(profile)
		if err := pac.ReinitProfile(profile.Name); err != nil {
			clog.Logger.Alertf("can't create profiles: %s", err.Error())
		}
	}()

	render.JSON(response, request, profile)
}

// DeleteProfile getting a profile
func DeleteProfile(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	for _, profile := range config.Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			pac.CloseProfile(profileName)
			profile, err := config.RemoveProfile(profileName)
			if err != nil {
				clog.Logger.Debug("Error deleting profile: \n" + err.Error())
				api.Err(response, request, err)
				return
			}
			err = config.DeleteProfileFile(profileName)
			if err != nil {
				clog.Logger.Debug("Error deleting profile file: \n" + err.Error())
				api.Err(response, request, err)
				return
			}
			render.JSON(response, request, profile)
			return
		}
	}

	api.NotFound(response, request, "profile", profileName)
	return
}

// GetExportProfile getting a profile
func GetExportProfile(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	profile, ok := getProfile(profileName)
	if !ok {
		clog.Logger.Debugf("Profile %s not found", profileName)
		api.NotFound(response, request, "profile", profileName)
		return
	}

	body, err := json.Marshal(profile)
	if err != nil {
		clog.Logger.Debug("Error reading profile: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	response.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.profile\"", profileName))
	render.Data(response, request, body)
}

func getProfile(profileName string) (models.Profile, bool) {

	for _, profile := range config.Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			return profile, true
		}
	}
	return models.Profile{}, false
}

// GetExportAction exporting a action from a profile as a file
func GetExportAction(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	profile, ok := getProfile(profileName)
	if !ok {
		clog.Logger.Debugf("Profile %s not found", profileName)
		api.NotFound(response, request, "profile", profileName)
		return
	}

	actionName, err := api.Param(request, "actionName")
	if err != nil {
		clog.Logger.Debug("Error reading action name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	action, err := profile.GetAction(actionName)
	if err != nil {
		clog.Logger.Debug("Error getting action: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	body, err := json.Marshal(action)
	if err != nil {
		clog.Logger.Debug("Error serialising action: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	response.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.action\"", actionName))
	render.Data(response, request, body)
}

// GetExportPage exporting a page from a profile as a file
func GetExportPage(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	profile, ok := getProfile(profileName)
	if !ok {
		clog.Logger.Debugf("Profile %s not found", profileName)
		api.NotFound(response, request, "profile", profileName)
		return
	}

	pageName, err := api.Param(request, "pageName")
	if err != nil {
		clog.Logger.Debug("Error reading page name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	exchange, err := transformer.ExportPage(profile, pageName)
	if err != nil {
		clog.Logger.Debug("Error getting page: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	body, err := json.Marshal(exchange)
	if err != nil {
		clog.Logger.Debug("Error serialising action: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	response.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.page\"", pageName))
	render.Data(response, request, body)
}

// PostCombine combines an export file with a profile
func PostCombine(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	profile, ok := getProfile(profileName)
	if !ok {
		clog.Logger.Debugf("Profile %s not found", profileName)
		api.NotFound(response, request, "profile", profileName)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var profileExchange models.ProfileExchange
	err = decoder.Decode(&profileExchange)
	if err != nil {
		clog.Logger.Debug("Error reading json body:" + err.Error())
		api.Err(response, request, err)
		return
	}

	newProfile, err := transformer.CombineProfile(profile, profileExchange)
	if err != nil {
		clog.Logger.Debug("Error reading json body:" + err.Error())
		api.Err(response, request, err)
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(response, request, newProfile)
}
