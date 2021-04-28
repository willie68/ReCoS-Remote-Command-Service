package routes

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"gopkg.in/yaml.v3"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/api/handler"
	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pac"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/web"
)

var icons []string
var mu sync.Mutex
var getIconsDo sync.Once
var initIconsMapper sync.Once
var iconMapperMap map[string]map[string]string
var defaultIcon = "help.svg"

/*
ConfigRoutes getting all routes for the config endpoint
*/
func ConfigRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/icons", GetIcons)
	router.Get("/icons/{mapper}/{key}", GetIconMapperKey)
	router.Get("/commands", GetCommands)
	router.Get("/icons/{iconname}", GetIcon)
	router.With(handler.AuthCheck()).Get("/check", GetCheck)
	initIconMapper()
	return router
}

// initIconMapper initalise the maps for ths icon mapper
func initIconMapper() {
	initIconsMapper.Do(func() {

		iconMapperMap = make(map[string]map[string]string)

		files, err := web.IconMapperAssets.ReadDir("iconmapper")
		if err != nil {
			clog.Logger.Debugf("Error reading mapper file: %v", err)
		}
		for _, file := range files {
			clog.Logger.Infof("iconmapper file: %s", file.Name())
			mapperConfig := strings.ToLower(file.Name())
			if strings.HasSuffix(mapperConfig, ".yaml") {
				mapperName := mapperConfig[0:strings.LastIndex(mapperConfig, ".")]
				bytes, err := web.IconMapperAssets.ReadFile("iconmapper/" + file.Name())
				if err != nil {
					clog.Logger.Errorf("Error reading mapper file: %v", err)
					continue
				}
				mapper := make(map[string]string)
				err = yaml.Unmarshal(bytes, &mapper)
				if err != nil {
					clog.Logger.Errorf("Error unmarshalling mapper file: %v", err)
					continue
				}
				iconMapperMap[mapperName] = mapper
			}
		}
	})
}

/*
GetIcons list of all possible icon names
*/
func GetIcons(response http.ResponseWriter, request *http.Request) {
	mu.Lock()
	getIconsDo.Do(func() {
		icons = make([]string, 0)
		files, err := web.WebClientAssets.ReadDir("webclient/assets")
		if err != nil {
			clog.Logger.Debug("Error reading icon files:" + err.Error())
			api.Err(response, request, err)
			return
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".svg") {
				icons = append(icons, file.Name())
			}
		}

		sort.Slice(icons, func(i, j int) bool { return strings.ToLower(icons[i]) < strings.ToLower(icons[j]) })
	})
	mu.Unlock()
	render.JSON(response, request, icons)
}

/*
GetIcon returning a converted icon back to the client
*/
func GetIcon(response http.ResponseWriter, request *http.Request) {
	var err error
	w := 72
	h := 72
	heightStr := request.URL.Query().Get("height")
	if heightStr != "" {
		h, err = strconv.Atoi(heightStr)
		if err != nil {
			clog.Logger.Debugf("Error reading icon height: %v", err)
		}
	}
	widthStr := request.URL.Query().Get("width")
	if widthStr != "" {
		w, err = strconv.Atoi(widthStr)
		if err != nil {
			clog.Logger.Debugf("Error reading icon width: %v", err)
		}
	}
	iconName, err := api.Param(request, "iconname")
	if err != nil {
		clog.Logger.Debugf("Error reading icon name: %v", err)
		api.Err(response, request, err)
		return
	}
	index := strings.Index(iconName, ".")
	if index >= 0 {
		iconName = iconName[0:index]
	}

	clog.Logger.Infof("config: icon convert: %s", iconName)

	srcPath := fmt.Sprintf("webclient/assets/%s.svg", iconName)
	in, err := web.WebClientAssets.Open(srcPath)
	if err != nil {
		if os.IsNotExist(err) {
			value := fmt.Sprintf("%s.png", iconName)
			handleFile(response, request, value)
		} else {
			clog.Logger.Debugf("Error reading icon: %v", err)
			api.Err(response, request, err)
			return
		}
	}
	defer in.Close()

	var b bytes.Buffer
	icon, _ := oksvg.ReadIconStream(in)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	err = png.Encode(&b, rgba)
	if err != nil {
		clog.Logger.Debugf("Error encoding icon: %v", err)
		api.Err(response, request, err)
	}

	ctype := mime.TypeByExtension(".png")
	response.Header().Set("Content-Type", ctype)
	response.Header().Set("Content-Length", strconv.FormatInt(int64(b.Len()), 10))
	response.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		r := bytes.NewBuffer(b.Bytes())
		io.Copy(response, r)
	}

	//render.JSON(response, request, iconName)
}

/*
GetIconMapperKey get aicon mapped from a source
*/
func GetIconMapperKey(response http.ResponseWriter, request *http.Request) {
	valueOK := true
	value := defaultIcon
	var iconMapper map[string]string

	mapperName, err := api.Param(request, "mapper")
	if err != nil {
		clog.Logger.Debugf("Error reading mapper name: %v", err)
		valueOK = false
	}
	if valueOK {
		mapperName = strings.ToLower(mapperName)
		ok := true
		iconMapper, ok = iconMapperMap[mapperName]
		if !ok {
			err = fmt.Errorf("mapper name \"%s\" not found", mapperName)
			clog.Logger.Debugf("Error getting mapper: %v", err)
			valueOK = false
		}
	}

	keyName, err := api.Param(request, "key")
	if err != nil {
		clog.Logger.Debugf("Error reading mapper key: %v", err)
		api.Err(response, request, err)
		return
	}
	if valueOK {
		ok := true
		keyName = strings.ToLower(keyName)
		value, ok = iconMapper[keyName]
		if !ok {
			err = fmt.Errorf("key name \"%s\" not found", keyName)
			clog.Logger.Debugf("Error getting mapper: %v", err)
			valueOK = false
		}
	}
	if !valueOK {
		value = defaultIcon
	}
	clog.Logger.Debugf("Mapper: \"%s:%s\"=\"%s\"", mapperName, keyName, value)
	if strings.Index(value, "#") == 0 {
		render.JSON(response, request, value)
	}
	handleFile(response, request, value)
}

func handleFile(response http.ResponseWriter, request *http.Request, value string) {
	file, err := web.WebClientAssets.Open("webclient/assets/" + value)
	if err != nil {
		clog.Logger.Debugf("Error reading file: %v", err)
		api.Err(response, request, err)
		return
	}
	defer file.Close()

	filestat, err := file.Stat()
	if err != nil {
		clog.Logger.Debugf("Error reading file: %v", err)
		api.Err(response, request, err)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(value))
	response.Header().Set("Content-Type", ctype)
	response.Header().Set("Content-Length", strconv.FormatInt(filestat.Size(), 10))
	response.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		io.Copy(response, file)
	}
}

/*
GetCommands list of all possible command types
*/
func GetCommands(response http.ResponseWriter, request *http.Request) {
	profileName := request.Header.Get("X-mcs-profile")
	var profile models.Profile
	if profileName != "" {
		profile, _ = config.GetProfile(profileName)
	}
	types := pac.CommandTypes
	types, err := pac.EnrichTypes(types, profile)
	if err != nil {
		clog.Logger.Errorf("Error reading commands files: %v", err)
		api.Err(response, request, err)
		return
	}
	render.JSON(response, request, types)
}

/*
GetCheck simply checks the authentication
*/
func GetCheck(response http.ResponseWriter, request *http.Request) {
	render.JSON(response, request, "ok")
}
