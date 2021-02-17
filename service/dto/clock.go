package dto

import (
	"fmt"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ClockCommand is a command to execute a delay. Using time for getting the ttime in seconds to delay the execution.
type ClockCommand struct {
	Parameters map[string]interface{}
	action     *Action
	stop       bool
	ticker     *time.Ticker
	done       chan bool
	format     string
}

// Init a delay in the actual context
func (c *ClockCommand) Init(a *Action) (bool, error) {
	c.action = a
	c.stop = false
	c.ticker = time.NewTicker(1 * time.Second)
	c.done = make(chan bool)
	value, found := c.Parameters["format"]
	c.format = "15:04:05"
	if found {
		var ok bool
		c.format, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("Format is in wrong format. Please use string as format")
		}
	}
	go func() {
		for {
			select {
			case <-c.done:
				return
			case t := <-c.ticker.C:
				title := t.Format(c.format)
				message := models.Message{
					Profile: a.Profile,
					Action:  a.Name,
					State:   1,
					Title:   title,
				}
				api.SendMessage(message)
			}
		}
	}()
	return true, nil
}

// Stop stops the actual command
func (c *ClockCommand) Stop(a *Action) (bool, error) {
	c.ticker.Stop()
	c.done <- true
	return true, nil
}

// Execute a delay in the actual context
func (c *ClockCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	return true, nil
}
