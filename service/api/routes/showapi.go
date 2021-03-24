package routes

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/config"
	"wkla.no-ip.biz/remote-desk-service/dto"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ShowRoutes getting all routes for the profile endpoint for showing in the client
func ShowRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", GetUIProfilesEndpoint)
	router.Get("/{profileName}", GetUIProfileEndpoint)
	router.Get("/{profileName}/{actionName}/{commandName}/{id}", GetGraphics)
	return router
}

// GetUIProfilesEndpoint getting all profile names
func GetUIProfilesEndpoint(response http.ResponseWriter, request *http.Request) {
	var profileInfos []models.ProfileShortInfo
	profileInfos = make([]models.ProfileShortInfo, 0)
	for _, profile := range config.Profiles {
		info := models.ProfileShortInfo{
			Name:        profile.Name,
			Description: profile.Description,
		}
		profileInfos = append(profileInfos, info)
	}

	sort.Slice(profileInfos, func(i, j int) bool {
		return strings.ToLower(profileInfos[i].Name) < strings.ToLower(profileInfos[j].Name)
	})

	result := models.ProfileInfos{
		Profiles: profileInfos,
	}
	render.JSON(response, request, result)
}

// GetUIProfileEndpoint getting all profile names
func GetUIProfileEndpoint(response http.ResponseWriter, request *http.Request) {
	profileName, err := api.Param(request, "profileName")
	if err != nil {
		clog.Logger.Debug("Error reading profile name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	for _, profile := range config.Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			// making a deep copy of this profile
			uiProfile := models.ProfileInfo{
				Name:        profile.Name,
				Description: profile.Description,
				Pages:       make([]models.PageInfo, 0),
				Actions:     make([]models.ActionInfo, 0),
			}
			// for every page creating the info object
			for _, page := range profile.Pages {
				pageInfo := models.PageInfo{
					Name:        page.Name,
					Description: page.Description,
					Icon:        page.Icon,
					Columns:     page.Columns,
					Rows:        page.Rows,
					Toolbar:     page.Toolbar,
					Cells:       make([]string, 0),
				}
				for _, cell := range page.Cells {
					pageInfo.Cells = append(pageInfo.Cells, cell)
				}
				uiProfile.Pages = append(uiProfile.Pages, pageInfo)
			}
			// for every action creating the info object
			for _, action := range profile.Actions {
				actionInfo := models.ActionInfo{
					Type:        action.Type,
					Name:        action.Name,
					Title:       action.Title,
					Description: action.Description,
					Icon:        action.Icon,
					Fontsize:    action.Fontsize,
					Fontcolor:   action.Fontcolor,
					Outlined:    action.Outlined,
				}
				uiProfile.Actions = append(uiProfile.Actions, actionInfo)
			}
			render.JSON(response, request, uiProfile)
			return
		}
	}
	return
}

// GetGraphics getting the graphics of the id
func GetGraphics(response http.ResponseWriter, request *http.Request) {
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
	commandName, err := api.Param(request, "commandName")
	if err != nil {
		clog.Logger.Debug("Error reading command name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	id, err := api.Param(request, "id")
	if err != nil {
		clog.Logger.Debug("Error reading id: \n" + err.Error())
		api.Err(response, request, err)
		return
	}
	width := 0
	widthStr, err := api.Query(request, "width")
	if err == nil {
		width, _ = strconv.Atoi(widthStr)
	}
	height := 0
	heightStr, err := api.Query(request, "height")
	if err == nil {
		height, _ = strconv.Atoi(heightStr)
	}

	graphicsInfo, err := dto.Graphics(profileName, actionName, commandName, id, width, height)
	if err != nil {
		clog.Logger.Debug("Error reading action name: \n" + err.Error())
		api.Err(response, request, err)
		return
	}

	response.Header().Set("Content-Type", graphicsInfo.Mimetype)
	response.Header().Set("Content-Length", strconv.Itoa(len(graphicsInfo.Data)))
	response.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	response.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		response.Write(graphicsInfo.Data)
	}
	return
}
