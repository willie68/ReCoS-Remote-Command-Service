package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"wkla.no-ip.biz/remote-desk-service/error/serror"

	"wkla.no-ip.biz/remote-desk-service/api"
)

// TenantHeader in this header thr right tenant should be inserted
const timeout = 1 * time.Minute

//APIKey the apikey of this service
var APIKey string

/*
ConfigDescription describres all metadata of a config
*/
type ConfigDescription struct {
	StoreID string `json:"storeid"`
	UserID  string `json:"userID"`
	Size    int    `json:"size"`
}

/*
ConfigRoutes getting all routes for the config endpoint
*/
func ConfigRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", PostConfigEndpoint)
	router.Get("/", GetConfigEndpoint)
	router.Delete("/", DeleteConfigEndpoint)
	router.Get("/size", GetConfigSizeEndpoint)
	return router
}

/*
GetConfigEndpoint getting if a store for a tenant is initialised
because of the automatic store creation, the value is more likely that data is stored for this tenant
*/
func GetConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	user := getUsername(request)
	if user == "" {
		msg := fmt.Sprintf("user header %s missing", api.UserHeader)
		api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
		return
	}

	c := ConfigDescription{
		StoreID: "myNewStore",
		UserID:  user,
		Size:    1234567,
	}
	render.JSON(response, request, c)
}

/*
PostConfigEndpoint create a new store for a tenant
because of the automatic store creation, this method will always return 201
*/
func PostConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	user := getUsername(request)
	if user == "" {
		msg := fmt.Sprintf("user header %s missing", api.UserHeader)
		api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
		return
	}

	render.Status(request, http.StatusCreated)
	render.JSON(response, request, user)
}

/*
DeleteConfigEndpoint deleting store for a tenant, this will automatically delete all data in the store
*/
func DeleteConfigEndpoint(response http.ResponseWriter, request *http.Request) {
	user := getUsername(request)
	if user == "" {
		msg := fmt.Sprintf("user header %s missing", api.UserHeader)
		api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
		return
	}

	render.JSON(response, request, user)
}

/*
GetConfigSizeEndpoint size of the store for a tenant
*/
func GetConfigSizeEndpoint(response http.ResponseWriter, request *http.Request) {
	user := getUsername(request)
	if user == "" {
		msg := fmt.Sprintf("user header %s missing", api.UserHeader)
		api.Err(response, request, serror.BadRequest(nil, "missing-user", msg))
		return
	}

	render.JSON(response, request, user)
}

/*
getUsername getting the tenant from the request
*/
func getUsername(req *http.Request) string {
	return req.Header.Get(api.UserHeader)
}
