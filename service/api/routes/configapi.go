package routes

import (
	"bytes"
	"encoding/json"
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
	"wkla.no-ip.biz/remote-desk-service/pkg"
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
	router.Get("/integrations", GetInteg)
	router.With(handler.AuthCheck()).Post("/integrations/{integname}", PostInteg)
	router.Get("/credits", GetCredits)
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

/*
GetInteg get parameter config of the integrations
*/
func GetInteg(response http.ResponseWriter, request *http.Request) {
	var localConfig map[string]interface{}
	result := make(map[string]interface{})
	settings := make(map[string]interface{})
	for _, integration := range pkg.IntegInfos {
		localConfig = nil
		extconfig := config.Get().ExternalConfig
		value, ok := extconfig[integration.Name]
		if ok {
			localConfig = value.(map[string]interface{})
		}

		setting := make(map[string]interface{})
		for _, param := range integration.Parameters {
			if localConfig != nil {
				setting[param.Name] = localConfig[param.Name]
			} else {
				switch param.Type {
				case "string":
					setting[param.Name] = ""
				case "[]string":
					setting[param.Name] = []string{}
				case "int":
					setting[param.Name] = 0
				case "bool":
					setting[param.Name] = false
				case "color":
					setting[param.Name] = ""
				case "icon":
					setting[param.Name] = ""
				case "date":
					setting[param.Name] = "2006-01-02"
				}
			}
		}
		settings[integration.Name] = setting
	}
	result["infos"] = pkg.IntegInfos
	result["settings"] = settings
	render.JSON(response, request, result)
}

// PostInteg post a new config
func PostInteg(response http.ResponseWriter, request *http.Request) {
	integName, err := api.Param(request, "integname")
	if err != nil {
		clog.Logger.Errorf("Error reading integ name: %v", err)
		api.Err(response, request, err)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var params map[string]interface{}
	err = decoder.Decode(&params)
	if err != nil {
		clog.Logger.Errorf("Error reading json body: %v", err)
		api.Err(response, request, err)
		return
	}
	localConfig := config.Get()
	extConfig := localConfig.ExternalConfig
	integConfig, ok := extConfig[integName].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("error getting integ config for name: %s", integName)
		clog.Logger.Error(err.Error())
		api.Err(response, request, err)
		return
	}
	for k, v := range params {
		_, ok := integConfig[k]
		if !ok {
			err := fmt.Errorf("error parameter not found: %s", k)
			clog.Logger.Error(err.Error())
			api.Err(response, request, err)
			return
		}
		integConfig[k] = v
	}
	config.Save()
	render.JSON(response, request, integConfig)
}

/*
GetCredits returning a converted icon back to the client
*/
func GetCredits(response http.ResponseWriter, request *http.Request) {
	credits := web.CreditsAsset

	ctype := mime.TypeByExtension(".html")
	response.Header().Set("Content-Type", ctype)
	response.Header().Set("Content-Length", strconv.FormatInt(int64(len(credits)), 10))
	response.WriteHeader(http.StatusOK)

	if request.Method != "HEAD" {
		r := bytes.NewBuffer([]byte(credits))
		io.Copy(response, r)
	}
}
