package pac

import (
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ShowTextCommandTypeInfo saving to the file system
var ShowTextCommandTypeInfo = models.CommandTypeInfo{
	Category:         "ReCoS",
	Type:             "SHOWTEXT",
	Name:             "ShowText",
	Description:      "showing a simple text",
	Icon:             "send.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "text",
			Type:           "string",
			Description:    "the text that should be sendet",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

// SendMessageCommand is a command to do a sceen shot and store this into the filesystem
type ShowTextCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (s *ShowTextCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return ShowTextCommandTypeInfo, nil
}

// Init a delay in the actual context
func (s *ShowTextCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop a delay in the actual context
func (s *ShowTextCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a delay in the actual context
func (s *ShowTextCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	msg, err := ConvertParameter2String(s.Parameters, "text", "")
	if err != nil {
		return false, err
	}
	message := models.Message{
		Profile: a.Profile,
		Action:  a.Name,
		Text:    msg,
		State:   0,
	}
	api.SendMessage(message)
	return false, nil
}
