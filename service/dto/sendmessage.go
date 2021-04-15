package dto

import (
	"fmt"

	"github.com/sqweek/dialog"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// SendMessageCommandTypeInfo saving to the file system
var SendMessageCommandTypeInfo = models.CommandTypeInfo{
	Category:         "ReCoS",
	Type:             "SENDMESSAGE",
	Name:             "Sendmessage",
	Description:      "Sending a Message",
	Icon:             "send.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "receiver",
			Type:           "string",
			Description:    "who is the receiver of the message",
			Unit:           "",
			WizardPossible: true,
			List:           []string{"client", "service"},
		},
		{
			Name:           "message",
			Type:           "string",
			Description:    "the message to send",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

// SendMessageCommand is a command to do a sceen shot and store this into the filesystem
type SendMessageCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (s *SendMessageCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return SendMessageCommandTypeInfo, nil
}

// Init a delay in the actual context
func (s *SendMessageCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop a delay in the actual context
func (s *SendMessageCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a delay in the actual context
func (s *SendMessageCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	msg, err := ConvertParameter2String(s.Parameters, "message", "")
	if err != nil {
		return false, fmt.Errorf("error getting parameter: %v", err)
	}

	receiver, err := ConvertParameter2String(s.Parameters, "receiver", "client")
	if err != nil {
		return false, fmt.Errorf("error getting parameter: %v", err)
	}
	clog.Logger.Infof("receiver: %s, message: %s ", receiver, msg)

	if receiver == "client" {
		message := models.Message{
			Profile:  a.Profile,
			Action:   a.Name,
			ImageURL: "check_mark.png",
			Command:  "sendmessage",
			Text:     msg,
			State:    0,
		}
		api.SendMessage(message)
	}
	if receiver == "service" {
		go func() {
			dialog.Message("Message:\r\n%s", msg).Title("ReCoS Service Message").Info()
		}()
		clog.Logger.Infof("message dialog")
	}

	return true, nil
}
