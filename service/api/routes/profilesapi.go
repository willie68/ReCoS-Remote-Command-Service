package routes

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gopkg.in/yaml.v3"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/api/handler"
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
	router.With(handler.AuthCheck()).Post("/", PostProfile)
	router.With(handler.AuthCheck()).Put("/{profileName}", PutProfile)
	router.Get("/{profileName}", GetProfile)
	router.With(handler.AuthCheck()).Delete("/{profileName}", DeleteProfile)
	router.Get("/{profileName}/export", GetExportProfile)
	return router
}

// GetProfiles getting all profile names
func GetProfiles(response http.ResponseWriter, request *http.Request) {
	var profileNames []string
	profileNames = make([]string, 0)
	for _, profile := range config.Profiles {
		profileNames = append(profileNames, profile.Name)
	}
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
		config.AddProfile(profile)
		if err := dto.ReinitProfiles(config.Profiles); err != nil {
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
		dto.CloseProfile(profile.Name)
		dto.RemoveProfile(profile.Name)
		config.UpdateProfile(profile)
		if err := dto.ReinitProfile(profile.Name); err != nil {
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
			dto.CloseProfile(profileName)
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

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	filename := fmt.Sprintf("%s.yaml", profileName)
	body, err := yaml.Marshal(profile)
	if err != nil {
		clog.Logger.Debug("Error reading profile: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	f, err := w.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		clog.Logger.Debug("Error writing profile: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	response.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", profileName))
	render.Data(response, request, buf.Bytes())
}

func getProfile(profileName string) (models.Profile, bool) {
	for _, profile := range config.Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			return profile, true
		}
	}
	return models.Profile{}, false
}
