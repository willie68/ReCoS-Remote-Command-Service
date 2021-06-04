// Special implementation for openhardwaremonitor app
// +build windows
package video

import (
	"fmt"
	"time"

	obsws "github.com/muesli/go-obs-websocket"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var OBSIntegInfo = models.IntegInfo{
	Category:    "System",
	Name:        "obs",
	Description: "obs is a software for streaming and recording.",
	Image:       "camcoder.svg",
	Parameters: []models.ParamInfo{
		{
			Name:           "active",
			Type:           "bool",
			Description:    "activate the obs integration",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "host",
			Type:           "string",
			Description:    "ip adress or name the pc where the obs is started. ",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name: "port",
			Type: "int",
			Unit: "",
			Description: "port of the obs websocket service. 	",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "password",
			Type:           "string",
			Description:    "password given in the obs websocket service. ",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

type OBS struct {
	host      string
	port      int
	password  string
	Active    bool
	Connected bool
	lastError error
	c         obsws.Client
}

var OBSInstance *OBS

// InitOBS initialise the obs websocket connection
func InitOBS(extconfig map[string]interface{}) error {
	var err error
	value, ok := extconfig["obs"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Debug("video:obs: found config")
			active, ok := config["active"].(bool)
			if !ok {
				active = false
			}
			var host string
			var password string
			var port int
			if active {
				clog.Logger.Debug("video:obs: active")
				host, ok = config["host"].(string)
				if !ok {
					err = fmt.Errorf("can't find host to connect to. %s", host)
				}
				port, ok = config["port"].(int)
				if !ok {
					port = 4444
				}
				password, ok = config["password"].(string)
				if !ok {
					password = ""
				}
			}
			OBSInstance = &OBS{
				Active:    active,
				Connected: false,
				host:      host,
				port:      port,
				password:  password,
			}
			err = OBSInstance.Connect()
		}
	}
	return err
}

func (o *OBS) Connect() error {
	o.c = obsws.Client{Host: o.host, Port: o.port, Password: o.password}
	if err := o.c.Connect(); err != nil {
		clog.Logger.Errorf("error connecting to obs. %v", err)
		return err
	}

	obsws.SetReceiveTimeout(time.Second * 2)
	go func() {
		profiles := o.GetSceneCollections()

		for _, profile := range profiles {
			clog.Logger.Info(profile)
		}

		o.SetProfile("GTA")
	}()
	return nil
}

func (o *OBS) Dispose() {
	o.c.Disconnect()
}

func DisposeOBS() {
	if OBSInstance != nil {
		OBSInstance.Dispose()
	}
}

func (o *OBS) IsRecording() bool {
	req := obsws.NewGetStreamingStatusRequest()
	// Send and receive a request asynchronously.
	resp, err := req.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	return resp.Recording
}

func (o *OBS) GetRecTimeCode() string {
	req := obsws.NewGetStreamingStatusRequest()
	// Send and receive a request asynchronously.
	resp, err := req.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	return resp.RecTimecode
}

// Start or stop recording, returns true if recording was startet
func (o *OBS) SwitchRecording() bool {
	srReq := obsws.NewStartStopRecordingRequest()
	srResp, err := srReq.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	clog.Logger.Debugf("start recording: %s", srResp.Status())
	time.Sleep(1 * time.Second)
	return o.IsRecording()
}

func (o *OBS) IsStreaming() bool {
	req := obsws.NewGetStreamingStatusRequest()
	// Send and receive a request asynchronously.
	resp, err := req.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	return resp.Streaming
}

func (o *OBS) GetStreamTimeCode() string {
	req := obsws.NewGetStreamingStatusRequest()
	// Send and receive a request asynchronously.
	resp, err := req.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	return resp.StreamTimecode
}

// SwitchStreaming Start or stop streaming, returns true if streaming was startet
func (o *OBS) SwitchStreaming() bool {
	srReq := obsws.NewStartStopStreamingRequest()
	srResp, err := srReq.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	clog.Logger.Debugf("start streaming: %s", srResp.Status())
	time.Sleep(1 * time.Second)
	return o.IsStreaming()
}

// GetScenes getting a list of scene names
func (o *OBS) GetSceneCollections() []string {
	srReq := obsws.NewListSceneCollectionsRequest()
	srResp, err := srReq.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	scenes := make([]string, 0)
	for _, profile := range srResp.SceneCollections {
		scenes = append(scenes, profile["sc-name"].(string))
	}
	clog.Logger.Debugf("get scenes: %v", srResp.Status())
	return scenes
}

// GetProfiles getting a list of profiles
func (o *OBS) GetProfiles() []string {
	srReq := obsws.NewListProfilesRequest()
	srResp, err := srReq.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	clog.Logger.Debugf("get profiles: %v", srResp.Status())
	profiles := make([]string, 0)
	for _, profile := range srResp.Profiles {
		profiles = append(profiles, profile["profile-name"].(string))
	}
	return profiles
}

func (o *OBS) SetProfile(name string) {
	srReq := obsws.NewSetCurrentProfileRequest(name)
	srResp, err := srReq.SendReceive(o.c)
	if err != nil {
		clog.Logger.Errorf("error in obs. %v", err)
	}
	clog.Logger.Debugf("set profiles: %v", srResp.Status())
}
