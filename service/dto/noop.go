package dto

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// NoopCommandTypeInfo is a command with no operation, but the possibility to change text and icon
var NoopCommandTypeInfo = models.CommandTypeInfo{"NOOP", "Noop", "do nothing", false, []models.CommandParameterInfo{}}

// NoopCommand is a command to do nothing.
type NoopCommand struct {
	Parameters map[string]interface{}
	action     *Action
}

// Init nothing
func (d *NoopCommand) Init(a *Action) (bool, error) {
	return true, nil
}

// Stop nothing
func (d *NoopCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (d *NoopCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if IsSingleClick(requestMessage) {
	}
	return true, nil
}
