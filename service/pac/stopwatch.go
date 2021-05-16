package pac

import (
	"fmt"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	durationfmt "wkla.no-ip.biz/remote-desk-service/internal"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// StopwatchCommandTypeInfo is a stopwatch with seconds
var StopwatchCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Time",
	Type:             "STOPWATCH",
	Name:             "Stopwatch",
	Description:      "Measure time with a stopwatch",
	Icon:             "rate.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "format",
			Type:           "string",
			Description:    "the format of the time",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

// StopwatchCommand is a command to simulate a stopwatch.
type StopwatchCommand struct {
	Parameters map[string]interface{}
	action     *Action
	stop       bool
	ticker     *time.Ticker
	done       chan bool
	format     string
	startTime  time.Time
	stopTime   time.Time
	running    bool
}

// EnrichType enrich the type info with the informations from the profile
func (s *StopwatchCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return StopwatchCommandTypeInfo, nil
}

// Init initialise the stopwatch parmeters
func (s *StopwatchCommand) Init(a *Action, commandName string) (bool, error) {
	s.action = a
	s.stop = false
	s.format = "%0h:%0m:%0s"
	s.done = make(chan bool)

	/*	value, found := c.Parameters["analog"]
		if found {
			var ok bool
			c.analog, ok = value.(bool)
			if !ok {
				return false, fmt.Errorf("Analog is in wrong format. Please use boolean as format")
			}
		}
	*/
	value, found := s.Parameters["format"]
	if found {
		var ok bool
		s.format, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("format is in wrong format. Please use string as format")
		}
	}
	return true, nil
}

// Stop stops the actual command
func (s *StopwatchCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

func (s *StopwatchCommand) startStopwatch() {
	s.running = true
	s.ticker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-s.done:
				return
			case t := <-s.ticker.C:
				if api.HasConnectionWithProfile(s.action.Profile) {
					tdelta := t.Sub(s.startTime)
					title, _ := durationfmt.Format(tdelta, s.format)
					message := models.Message{
						Profile: s.action.Profile,
						Action:  s.action.Name,
						State:   1,
						Title:   title,
					}
					api.SendMessage(message)
				}
			}
		}
	}()
}

func (s *StopwatchCommand) stopStopwatch() {
	s.ticker.Stop()
	s.done <- true
	s.running = false
}

// Execute a delay in the actual context
func (s *StopwatchCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	timeNow := time.Now()
	if IsDblClick(requestMessage) {
		if s.running {
			s.stopTime = timeNow
			s.stopStopwatch()
		}
		title, _ := durationfmt.Format(0, s.format)
		message := models.Message{
			Profile: s.action.Profile,
			Action:  s.action.Name,
			State:   1,
			Title:   s.action.Config.Title,
			Text:    title,
		}
		api.SendMessage(message)
	} else {
		if s.running {
			// stop the running clock
			s.stopTime = timeNow
			s.stopStopwatch()

			go func() {
				tdelta := s.stopTime.Sub(s.startTime)
				title, _ := durationfmt.Format(tdelta, s.format)
				message := models.Message{
					Profile: s.action.Profile,
					Action:  s.action.Name,
					State:   1,
					Title:   s.action.Config.Title,
					Text:    title,
				}
				api.SendMessage(message)
				time.Sleep(3 * time.Second)
			}()
		} else {
			// start the stopwatch
			s.startTime = timeNow
			s.startStopwatch()
		}
	}
	return false, nil
}
