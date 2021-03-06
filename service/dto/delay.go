package dto

import (
	"fmt"
	"time"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var DelayCommandTypeInfo = models.CommandTypeInfo{"DELAY", "Delay", "Setting up a short delay", []models.CommandParameterInfo{}}

// DelayCommand is a command to execute a delay. Using time for getting the ttime in seconds to delay the execution.
type DelayCommand struct {
	Parameters map[string]interface{}
}

// Init a delay in the actual context
func (d *DelayCommand) Init(a *Action) (bool, error) {
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
