package dto

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/audio"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var deviceName = "device name"
var commandName = "command name"

// AudioVolumeCommandTypeInfo switch to another page
var AudioVolumeCommandTypeInfo = models.CommandTypeInfo{"AUDIOVOLUME", "AudioVolume", "setting the volume of an audio device", true, []models.CommandParameterInfo{
	{deviceName, "string", "name of the device to control", "", true, make([]string, 0)},
	{commandName, "string", "the command to send", "", true, make([]string, 0)},
}}

var audioCommandArray = []string{"mute", "volume up", "volume down"}

// AudioVolumeCommand is a command to switch to another page.
// Using "page" for the page name
type AudioVolumeCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (v *AudioVolumeCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	commandInfo := AudioVolumeCommandTypeInfo
	index := -1
	for x, param := range commandInfo.Parameters {
		if param.Name == commandName {
			index = x
		}
	}
	if index >= 0 {
		commandInfo.Parameters[index].List = make([]string, 0)
		for _, audioCommand := range audioCommandArray {
			commandInfo.Parameters[index].List = append(commandInfo.Parameters[index].List, audioCommand)
		}
	}

	index = -1
	for x, param := range commandInfo.Parameters {
		if param.Name == deviceName {
			index = x
		}
	}
	if index >= 0 {
		commandInfo.Parameters[index].List = make([]string, 0)
		for _, sessionName := range audio.GetSessionNames() {
			commandInfo.Parameters[index].List = append(commandInfo.Parameters[index].List, sessionName)
		}
	}
	return commandInfo, nil
}

// Init the command
func (v *AudioVolumeCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop the command
func (v *AudioVolumeCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (v *AudioVolumeCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	return true, nil
}
