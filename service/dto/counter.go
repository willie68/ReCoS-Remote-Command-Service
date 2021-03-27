package dto

import (
	"fmt"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/session"
)

// CounterCommandTypeInfo switch to another page
var CounterCommandTypeInfo = models.CommandTypeInfo{
	Type:             "COUNTER",
	Name:             "Counter",
	Description:      "Counting button clicks",
	Icon:             "slot_machine.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "persist",
			Type:           "bool",
			Description:    "persist the value between restarts",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// CounterCommand is a command to switch to another page.
// Using "page" for the page name
type CounterCommand struct {
	Parameters  map[string]interface{}
	a           *Action
	commandName string
	countValue  float64
	persist     bool
}

// EnrichType enrich the type info with the informations from the profile
func (c *CounterCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return CounterCommandTypeInfo, nil
}

// Init the command
func (c *CounterCommand) Init(a *Action, commandName string) (bool, error) {
	c.a = a
	c.commandName = commandName
	c.persist = false
	value, found := c.Parameters["persist"]
	if found {
		var ok bool
		c.persist, ok = value.(bool)
		if !ok {
			return false, fmt.Errorf("Persist is in wrong format. Please use boolean as format")
		}
	}

	if c.persist {
		value, ok := session.SessionCache.RetrieveCommandData(a.Profile, a.Name, c.commandName)
		if ok {
			c.countValue = value.(float64)
		}
	}

	return true, nil
}

// Stop the command
func (c *CounterCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (c *CounterCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if c.persist {
		value, ok := session.SessionCache.RetrieveCommandData(a.Profile, a.Name, c.commandName)
		if ok {
			c.countValue = value.(float64)
		}
	}

	if IsDblClick(requestMessage) {
		c.countValue = 0
	}
	if IsSingleClick(requestMessage) {

		c.countValue += 1
	}
	if c.persist {
		session.SessionCache.StoreCommandData(a.Profile, a.Name, c.commandName, c.countValue)
	}
	c.UpdateClients(a, c.commandName)
	return false, nil
}

func (c *CounterCommand) UpdateClients(a *Action, commandName string) (bool, error) {
	message := models.Message{
		Profile: a.Profile,
		Action:  a.Name,
		Text:    fmt.Sprintf("%d", int(c.countValue)),
		State:   0,
	}
	api.SendMessage(message)
	return true, nil
}
