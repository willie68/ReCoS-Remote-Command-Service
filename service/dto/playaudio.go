package dto

import (
	"fmt"

	"wkla.no-ip.biz/remote-desk-service/pkg/audio"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PlayAudioCommandTypeInfo switch to another page
var PlayAudioCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio-Video",
	Type:             "PLAYAUDIO",
	Name:             "Playaudio",
	Description:      "playing an audio file",
	Icon:             "music_beamed_note.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "file",
			Type:           "string",
			Description:    "name and path of the file to play",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// PageCommand is a command to switch to another page.
// Using "page" for the page name
type PlayAudioCommand struct {
	audiofile  string
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (p *PlayAudioCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return PlayAudioCommandTypeInfo, nil
}

// Init the command
func (p *PlayAudioCommand) Init(a *Action, commandName string) (bool, error) {
	audiofile, err := ConvertParameter2String(p.Parameters, "file", "")
	if err != nil {
		return false, fmt.Errorf("the file parameter is in wrong format. Please use string as format. \r\n%v", err)
	}
	p.audiofile = audiofile
	return true, nil
}

// Stop the command
func (p *PlayAudioCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (p *PlayAudioCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	go audio.PlayAudio(p.audiofile)
	return true, nil
}
