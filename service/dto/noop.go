package dto

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// NoopCommandTypeInfo is a command with no operation, but the possibility to change text and icon
var NoopCommandTypeInfo = models.CommandTypeInfo{
	Category:         "useful",
	Type:             "NOOP",
	Name:             "Noop",
	Description:      "do nothing",
	Icon:             "",
	WizardPossible:   false,
	WizardActionType: models.Single,
	Parameters:       []models.CommandParameterInfo{},
	/*
		{
			Name:           "first",
			Type:           "string",
			Description:    "name of the page to switch to",
			Unit:           "",
			WizardPossible: true,
			List:           []string{"one", "two", "three"},
		},
		{
			Name:           "second",
			Type:           "string",
			Description:    "name of the page to switch to",
			Unit:           "",
			WizardPossible: true,
			List:           []string{"one: 1", "two: 1", "two: 2", "three: 1", "three: 2", "three: 3", "four: 1"},
			FilteredList:   "first",
		},
	},*/
}

// NoopCommand is a command to do nothing.
type NoopCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (d *NoopCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return NoopCommandTypeInfo, nil
}

// Init nothing
func (d *NoopCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop nothing
func (d *NoopCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (d *NoopCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	return true, nil
}
