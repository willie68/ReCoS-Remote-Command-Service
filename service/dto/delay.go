package dto

import (
	"fmt"
	"time"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var DelayCommandTypeInfo = models.CommandTypeInfo{
	Type:           "DELAY",
	Name:           "Delay",
	Description:    "Setting up a short delay",
	Icon:           "",
	WizardPossible: false,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "time",
			Type:           "int",
			Description:    "delay time in seconds",
			Unit:           " Seconds",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

// DelayCommand is a command to execute a delay. Using time for getting the ttime in seconds to delay the execution.
type DelayCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (d *DelayCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return DelayCommandTypeInfo, nil
}

// Init a delay in the actual context
func (d *DelayCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop a delay in the actual context
func (d *DelayCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a delay in the actual context
func (d *DelayCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	value, found := d.Parameters["time"]
	if found {
		delayValue, ok := value.(int)
		if ok {
			clog.Logger.Infof("delay with %v seconds", delayValue)
			time.Sleep(time.Duration(delayValue) * time.Second)
		} else {
			return false, fmt.Errorf("Time is in wrong format. Please use integer as format")
		}
	} else {
		return false, fmt.Errorf("Time is missing")
	}
	return true, nil
}
