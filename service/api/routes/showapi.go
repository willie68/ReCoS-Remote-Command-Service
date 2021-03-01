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

// ShowRoutes getting all routes for the profile endpoint for showing in the client
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
	var profileInfos []models.ProfileShortInfo
	profileInfos = make([]models.ProfileShortInfo, 0)
	for _, profile := range config.Profiles {
		info := models.ProfileShortInfo{
			Name:        profile.Name,
			Description: profile.Description,
		}
		profileInfos = append(profileInfos, info)
	}
	result := models.ProfileInfos{
		Profiles: profileInfos,
	}
	render.JSON(response, request, result)
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
