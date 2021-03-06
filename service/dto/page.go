package dto

import (
	"fmt"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PageCommandTypeInfo switch to another page
var PageCommandTypeInfo = models.CommandTypeInfo{"PAGE", "Page", "switching to another page", []models.CommandParameterInfo{
	{"page", "string", "name of the page to switch to", make([]string, 0)},
}}

// PageCommand is a command to switch to another page.
// Using "page" for the page name
type PageCommand struct {
	Parameters map[string]interface{}
}

// Init the command
func (p *PageCommand) Init(a *Action) (bool, error) {
	return true, nil
}

// Stop the command
func (p *PageCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (p *PageCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	value, found := p.Parameters["page"]
	if found {
		cmdValue, ok := value.(string)
		if ok {
			message := models.Message{
				Profile:  a.Profile,
				Page:     cmdValue,
				ImageURL: "check_mark.png",
				State:    0,
			}
			api.SendMessage(message)
		} else {
			return false, fmt.Errorf("The command parameter is in wrong format. Please use string as format")
		}
	} else {
		return false, fmt.Errorf("The command parameter is missing")
	}
	return true, nil
}
