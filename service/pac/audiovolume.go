package pac

import (
	"fmt"
	"math"
	"time"

	"wkla.no-ip.biz/remote-desk-service/internal/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/audio"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var parameterDeviceName = "device"
var parameterCommandName = "command"

// AudioVolumeCommandTypeInfo switch to another page
var AudioVolumeCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio-Video",
	Type:             "AUDIOVOLUME",
	Name:             "AudioVolume",
	Description:      "setting the volume of an audio device",
	Icon:             "speaker.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:        parameterDeviceName,
			Type:        "string",
			Description: "name of the device to control",
			Unit:        "", WizardPossible: true,
			List: make([]string, 0),
		},
		{
			Name:           parameterCommandName,
			Type:           "string",
			Description:    "the command to send",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

var audioCommandArray = []string{"mute", "volume up", "volume down"}

// AudioVolumeCommand is a command to switch to another page.
// Using "page" for the page name
type AudioVolumeCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	devicename  string
	command     string
	ticker      *time.Ticker
	done        chan bool
}

// EnrichType enrich the type info with the informations from the profile
func (v *AudioVolumeCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	commandInfo := AudioVolumeCommandTypeInfo
	index := -1
	for x, param := range commandInfo.Parameters {
		if param.Name == parameterCommandName {
			index = x
		}
	}
	if index >= 0 {
		commandInfo.Parameters[index].List = make([]string, 0)
		commandInfo.Parameters[index].List = append(commandInfo.Parameters[index].List, audioCommandArray...)
	}

	index = -1
	for x, param := range commandInfo.Parameters {
		if param.Name == parameterDeviceName {
			index = x
		}
	}
	if index >= 0 {
		commandInfo.Parameters[index].List = make([]string, 0)
		commandInfo.Parameters[index].List = append(commandInfo.Parameters[index].List, audio.GetSessionNames()...)
	}
	return commandInfo, nil
}

// Init the command
func (v *AudioVolumeCommand) Init(a *Action, commandName string) (bool, error) {
	v.commandName = commandName
	v.action = a

	object, ok := v.Parameters[parameterDeviceName]
	if !ok {
		return false, fmt.Errorf("the device parameter is empty")
	}

	v.devicename, ok = object.(string)
	if !ok {
		return false, fmt.Errorf("the device parameter should be a string")
	}

	object, ok = v.Parameters[parameterCommandName]
	if !ok {
		return false, fmt.Errorf("the command parameter is empty")
	}

	v.command, ok = object.(string)
	if !ok {
		return false, fmt.Errorf("the command parameter should be a string")
	}

	v.ticker = time.NewTicker(1 * time.Second)
	v.done = make(chan bool)
	go func() {
		for {
			select {
			case <-v.done:
				return
			case <-v.ticker.C:
				if api.HasConnectionWithProfile(a.Profile) {
					session, ok := audio.GetSession(v.devicename)
					if ok {
						v.SendSessionMute(session.GetMute(), session.GetVolume(), session.IsInput())
					}
					continue
				}
			}
		}
	}()

	return true, nil
}

// Stop the command
func (v *AudioVolumeCommand) Stop(a *Action) (bool, error) {
	v.done <- true
	return true, nil
}

// Execute the command
func (v *AudioVolumeCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	session, ok := audio.GetSession(v.devicename)
	if !ok {
		return false, fmt.Errorf("device %s not found", v.devicename)
	}
	switch v.command {
	case audioCommandArray[0]:
		{
			muted := session.GetMute()
			session.SetMute(!muted)
		}
	case audioCommandArray[1]:
		{
			volume := session.GetVolume()
			volume += 0.05
			if volume > 1 {
				volume = 1.0
			}
			session.SetVolume(volume)
		}
	case audioCommandArray[2]:
		{
			volume := session.GetVolume()
			volume -= 0.05
			if volume < 0 {
				volume = 0.0
			}
			session.SetVolume(volume)
		}
	}

	v.SendSessionMute(session.GetMute(), session.GetVolume(), session.IsInput())

	return false, nil
}

func (v *AudioVolumeCommand) SendSessionMute(muted bool, volume float32, isInput bool) {
	var image string
	var text string
	if isInput {
		image = "radio_microphone.svg"
		text = " "
		if muted {
			image = "radio_microphone_off.svg"
			text = "muted"
		}
	} else {
		image = "audio_volume_high.svg"
		percent := int(math.Round(float64(volume) * 100.0))
		text = fmt.Sprintf("Volume: %d%%", percent)
		if volume < 0.66 {
			image = "audio_volume_medium.svg"
		}
		if volume < 0.33 {
			image = "audio_volume_low.svg"
		}
		if muted {
			image = "audio_volume_mute.svg"
			text = "muted"
		}
	}
	message := models.Message{
		Profile:  v.action.Profile,
		Action:   v.action.Name,
		ImageURL: image,
		Text:     text,
		State:    0,
	}
	api.SendMessage(message)
}
