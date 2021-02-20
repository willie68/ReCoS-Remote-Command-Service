package dto

import (
	"fmt"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// TimerCommand is a command to start a timer. Using time for getting the time in seconds.
// For formatting the response the parameters fomat and finished are responsible.
// Use %d for inserting the actual time to wait.
type TimerCommand struct {
	Parameters map[string]interface{}
}

// Init a timer in the actual context
func (d *TimerCommand) Init(a *Action) (bool, error) {
	return true, nil
}

// Stop a timer in the actual context
func (d *TimerCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a timer in the actual context
func (d *TimerCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	value, found := d.Parameters["format"]
	format := "%d seconds"
	if found {
		var ok bool
		format, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("Format is in wrong format. Please use string as format")
		}
	}
	value, found = d.Parameters["finished"]
	finnished := "finished"
	if found {
		var ok bool
		finnished, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("Format is in wrong format. Please use string as format")
		}
	}
	value, found = d.Parameters["time"]
	if found {
		delayValue, ok := value.(int)
		if ok {
			clog.Logger.Infof("count down with %v seconds", delayValue)
			for ; delayValue > 0; delayValue-- {
				// TODO get this from the config
				title := fmt.Sprintf(format, delayValue)
				icon := "point_green.png"
				if delayValue < 4 {
					icon = "point_yellow.png"
				}
				message := models.Message{
					Profile:  a.Profile,
					Action:   a.Name,
					State:    delayValue,
					Title:    title,
					ImageURL: icon,
				}
				api.SendMessage(message)
				time.Sleep(1 * time.Second)
			}
			message := models.Message{
				Profile:  a.Profile,
				Action:   a.Name,
				State:    0,
				ImageURL: "point_red.png",
				Title:    finnished,
			}
			api.SendMessage(message)
			time.Sleep(1 * time.Second)
		} else {
			return false, fmt.Errorf("Time is in wrong format. Please use integer as format")
		}
	} else {
		return false, fmt.Errorf("Time is missing")
	}
	return true, nil
}
