// +build windows
package pac

import (
	"wkla.no-ip.biz/remote-desk-service/api"
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

// OBSStartStopCommand This command connects to the openhardwaremonitor application on windows.
// With this you can get different sensors of your computer. For using the webserver of the openhardwaremonitor
// app, you have to add another external configuration into the main service configuration.
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
