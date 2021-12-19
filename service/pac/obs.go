// +build windows
package pac

import (
	"wkla.no-ip.biz/remote-desk-service/internal/api"
	"wkla.no-ip.biz/remote-desk-service/internal/utils"
	"wkla.no-ip.biz/remote-desk-service/pkg/video"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

const (
	MODE_RECORDING = "recording"
	MODE_STREAMING = "streaming"
)

// OBSStartStopCommandTypeInfo showing hardware sensor data
var OBSStartStopCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio-Video",
	Type:             "OBSSTARTSTOP",
	Name:             "OBSStartStop",
	Description:      "Start/Stop recording/streaming on OBS",
	Icon:             "play_pause.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "mode",
			Type:           "string",
			Description:    "the mode recording or streaming",
			Unit:           "",
			WizardPossible: true,
			List:           []string{MODE_RECORDING, MODE_STREAMING},
		},
	},
}

// OBSProfileCommandTypeInfo showing hardware sensor data
var OBSProfileCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio-Video",
	Type:             "OBSPROFILE",
	Name:             "OBSProfile",
	Description:      "switch profile on OBS",
	Icon:             "folder.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "profile",
			Type:           "string",
			Description:    "the profile to switch to",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// OBSSceneCollectionCommandTypeInfo showing hardware sensor data
var OBSSceneCollectionCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio-Video",
	Type:             "OBSSCENECOLLECTION",
	Name:             "OBSSceneCollection",
	Description:      "switch the scene collection on OBS",
	Icon:             "backgrounds.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "scenecollection",
			Type:           "string",
			Description:    "the scenecollection to switch to",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

const (
	SCENECMD_NEXT   = "next"
	SCENECMD_PREV   = "previous"
	SCENECMD_FIRST  = "first"
	SCENECMD_LAST   = "last"
	SCENECMD_SWITCH = "switch"
)

// OBSSceneCommandTypeInfo showing hardware sensor data
var OBSSceneCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio-Video",
	Type:             "OBSSCENE",
	Name:             "OBSScene",
	Description:      "up/down and switch scenes on obs",
	Icon:             "scene2.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "scenecommand",
			Type:           "string",
			Description:    "the scene command",
			Unit:           "",
			WizardPossible: true,
			List:           []string{SCENECMD_SWITCH, SCENECMD_FIRST, SCENECMD_NEXT, SCENECMD_PREV, SCENECMD_LAST},
		},
		{
			Name:           "scenename",
			Type:           "string",
			Description:    "the scene name to switch to, only used in switch command ",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// OBSStartStopCommand
// This command has the following parameters:
// mode: the mode recording or streaming
type OBSStartStopCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	mode        string
}

// EnrichType enrich the type info with the informations from the profile
func (o *OBSStartStopCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return OBSStartStopCommandTypeInfo, nil
}

// Init nothing
func (o *OBSStartStopCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	o.action = a
	o.commandName = commandName
	o.mode, err = ConvertParameter2String(o.Parameters, "mode", MODE_RECORDING)
	if err != nil {
		clog.Logger.Errorf("error in getting mode: %v", err)
		return false, err
	}
	return true, nil
}

// Stop nothing
func (o *OBSStartStopCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (o *OBSStartStopCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if IsSingleClick(requestMessage) {
		if video.OBSInstance != nil {
			var ok bool
			if o.mode == MODE_RECORDING {
				ok = video.OBSInstance.SwitchRecording()
			}
			if o.mode == MODE_STREAMING {
				video.OBSInstance.SwitchStreaming()
			}
			message := models.Message{
				Profile: o.action.Profile,
				Action:  o.action.Name,
				State:   0,
			}
			if ok {
				message.ImageURL = "pause.svg"
			} else {
				message.ImageURL = "play_pause.svg"
			}
			api.SendMessage(message)
		}
	}
	return false, nil
}

// OBSProfileCommand
// This command has the following parameters:
// profile: the to use
type OBSProfileCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	profile     string
}

// EnrichType enrich the type info with the informations from the profile
func (o *OBSProfileCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	profiles, err := video.OBSInstance.GetProfiles()
	if err != nil {
		clog.Logger.Errorf("error in getting profiles: %v", err)
		return OBSProfileCommandTypeInfo, nil
	}
	index := GetIndexOfParameter(OBSProfileCommandTypeInfo.Parameters, "profile")
	if index >= 0 {
		OBSProfileCommandTypeInfo.Parameters[index].List = make([]string, 0)
		OBSProfileCommandTypeInfo.Parameters[index].List = append(OBSProfileCommandTypeInfo.Parameters[index].List, profiles...)
	}
	return OBSProfileCommandTypeInfo, nil
}

// Init nothing
func (o *OBSProfileCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	o.action = a
	o.commandName = commandName
	o.profile, err = ConvertParameter2String(o.Parameters, "profile", "")
	if err != nil {
		clog.Logger.Errorf("error in getting profile: %v", err)
		return false, err
	}
	return true, nil
}

// Stop nothing
func (o *OBSProfileCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (o *OBSProfileCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if IsSingleClick(requestMessage) {
		if video.OBSInstance != nil {
			err := video.OBSInstance.SetProfile(o.profile)
			if err != nil {
				message := models.Message{
					Profile:  o.action.Profile,
					Action:   o.action.Name,
					Text:     err.Error(),
					State:    0,
					ImageURL: "close.svg",
				}
				api.SendMessage(message)
			}
			return true, err
		}
	}
	return true, nil
}

// OBSSceneCollectionCommand
// This command has the following parameters:
// profile: the to use
type OBSSceneCollectionCommand struct {
	Parameters      map[string]interface{}
	action          *Action
	commandName     string
	scenecollection string
}

// EnrichType enrich the type info with the informations from the profile
func (o *OBSSceneCollectionCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	sceneCollections, err := video.OBSInstance.GetSceneCollections()
	if err != nil {
		clog.Logger.Errorf("error in getting scene collections: %v", err)
		return OBSSceneCollectionCommandTypeInfo, nil
	}
	index := GetIndexOfParameter(OBSSceneCollectionCommandTypeInfo.Parameters, "scenecollection")
	if index >= 0 {
		OBSSceneCollectionCommandTypeInfo.Parameters[index].List = make([]string, 0)
		OBSSceneCollectionCommandTypeInfo.Parameters[index].List = append(OBSSceneCollectionCommandTypeInfo.Parameters[index].List, sceneCollections...)
	}
	return OBSSceneCollectionCommandTypeInfo, nil
}

// Init nothing
func (o *OBSSceneCollectionCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	o.action = a
	o.commandName = commandName
	o.scenecollection, err = ConvertParameter2String(o.Parameters, "scenecollection", "")
	if err != nil {
		clog.Logger.Errorf("error in getting scenecollection: %v", err)
		return false, err
	}
	return true, nil
}

// Stop nothing
func (o *OBSSceneCollectionCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (o *OBSSceneCollectionCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if IsSingleClick(requestMessage) {
		if video.OBSInstance != nil {
			err := video.OBSInstance.SetSceneCollection(o.scenecollection)
			if err != nil {
				message := models.Message{
					Profile:  o.action.Profile,
					Action:   o.action.Name,
					Text:     err.Error(),
					State:    0,
					ImageURL: "close.svg",
				}
				api.SendMessage(message)
			}
			return true, err
		}
	}
	return true, nil
}

// OBSSceneCommand
// This command has the following parameters:
// profile: the to use
type OBSSceneCommand struct {
	Parameters   map[string]interface{}
	action       *Action
	commandName  string
	scenecommand string
	scenename    string
}

// EnrichType enrich the type info with the informations from the profile
func (o *OBSSceneCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	scenes, err := video.OBSInstance.GetScenes()
	if err != nil {
		clog.Logger.Errorf("error in getting scenes: %v", err)
		return OBSSceneCommandTypeInfo, nil
	}
	index := GetIndexOfParameter(OBSSceneCommandTypeInfo.Parameters, "scenename")
	if index >= 0 {
		OBSSceneCommandTypeInfo.Parameters[index].List = make([]string, 0)
		OBSSceneCommandTypeInfo.Parameters[index].List = append(OBSSceneCommandTypeInfo.Parameters[index].List, scenes...)
	}
	return OBSSceneCommandTypeInfo, nil
}

// Init nothing
func (o *OBSSceneCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	o.action = a
	o.commandName = commandName
	o.scenecommand, err = ConvertParameter2String(o.Parameters, "scenecommand", SCENECMD_SWITCH)
	if err != nil {
		clog.Logger.Errorf("error in getting scene command: %v", err)
		return false, err
	}
	o.scenename, err = ConvertParameter2String(o.Parameters, "scenename", "")
	if err != nil {
		clog.Logger.Errorf("error in getting scene name: %v", err)
		return false, err
	}
	return true, nil
}

// Stop nothing
func (o *OBSSceneCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (o *OBSSceneCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if IsSingleClick(requestMessage) {
		if video.OBSInstance != nil {
			var err error
			var scenes []string
			var scene string
			switch o.scenecommand {
			case SCENECMD_SWITCH:
				err = video.OBSInstance.SetScene(o.scenename)
			case SCENECMD_FIRST:
				scenes, err = video.OBSInstance.GetScenes()
				if err == nil {
					err = video.OBSInstance.SetScene(scenes[0])
				}
			case SCENECMD_LAST:
				scenes, err = video.OBSInstance.GetScenes()
				if err == nil {
					err = video.OBSInstance.SetScene(scenes[len(scenes)-1])
				}
			case SCENECMD_NEXT:
				scenes, err = video.OBSInstance.GetScenes()
				if err == nil {
					scene, err = video.OBSInstance.GetCurrentScene()
					if err == nil {
						index := utils.SliceIndex(len(scenes), func(i int) bool { return scenes[i] == scene })
						index = (index + 1) % len(scenes)
						err = video.OBSInstance.SetScene(scenes[index])
					}
				}
			case SCENECMD_PREV:
				scenes, err = video.OBSInstance.GetScenes()
				if err == nil {
					scene, err = video.OBSInstance.GetCurrentScene()
					if err == nil {
						index := utils.SliceIndex(len(scenes), func(i int) bool { return scenes[i] == scene })
						index = (index - 1)
						if index < 0 {
							index = len(scenes) - 1
						}
						err = video.OBSInstance.SetScene(scenes[index])
					}
				}
			}
			if err != nil {
				message := models.Message{
					Profile:  o.action.Profile,
					Action:   o.action.Name,
					Text:     err.Error(),
					State:    0,
					ImageURL: "close.svg",
				}
				api.SendMessage(message)
			}
			return true, err
		}
	}
	return true, nil
}
