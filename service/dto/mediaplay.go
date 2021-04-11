package dto

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var mediaParameterCommandName = "command"

// MediaPlayCommandTypeInfo switch to another page
var MediaPlayCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Audio/Video",
	Type:             "MEDIAPLAY",
	Name:             "Mediaplay",
	Description:      "controlling a media player",
	Icon:             "music_beamed_note.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           mediaParameterCommandName,
			Type:           "string",
			Description:    "the command to send",
			Unit:           "",
			WizardPossible: true,
			List:           mediaCommandArray,
		},
	},
}

var mediaCommandArray = []string{"play", "stop", "next", "previous"}
var mediaKeysArray = []int{keybd_event.VK_MEDIA_PLAY_PAUSE, keybd_event.VK_MEDIA_STOP, keybd_event.VK_MEDIA_NEXT_TRACK, keybd_event.VK_MEDIA_PREV_TRACK}

// MediaPlayCommand is a command to switch to another page.
// Using "page" for the page name
type MediaPlayCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	command     string
	keys        int
}

// EnrichType enrich the type info with the informations from the profile
func (v *MediaPlayCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return MediaPlayCommandTypeInfo, nil
}

// Init the command
func (v *MediaPlayCommand) Init(a *Action, commandName string) (bool, error) {
	v.commandName = commandName
	v.action = a

	object, ok := v.Parameters[parameterCommandName]
	if !ok {
		return false, fmt.Errorf("The command parameter is empty.")
	}

	v.command, ok = object.(string)
	if !ok {
		return false, fmt.Errorf("The command parameter should be a string.")
	}
	found := false
	for x, cmd := range mediaCommandArray {
		if cmd == strings.ToLower(v.command) {
			v.keys = mediaKeysArray[x]
			found = true
		}
	}
	if !found {
		return false, errors.New(fmt.Sprintf("can't find command: %s", v.command))
	}
	return true, nil
}

// Stop the command
func (v *MediaPlayCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (v *MediaPlayCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		clog.Logger.Errorf("error: %v", err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	kb.Clear()

	kb.SetKeys(v.keys)
	err = kb.Launching()
	if err != nil {
		return true, err
	}
	return true, nil
}
