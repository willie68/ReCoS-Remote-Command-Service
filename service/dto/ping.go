package dto

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-ping/ping"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PingCommandTypeInfo start an browser with directly with a url or filepath
var PingCommandTypeInfo = models.CommandTypeInfo{
	Type:             "PING",
	Name:             "Ping",
	Description:      "Execute a ping to a url and show the answer in ms",
	Icon:             "world_shipping.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "url",
			Type:           "string",
			Description:    "the url to ping to",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "period",
			Type:           "int",
			Description:    "period in secods",
			Unit:           "Seconds",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	}}

// PingCommand is a command to check the ping times.
// Using "url" for the server.
type PingCommand struct {
	Parameters map[string]interface{}
	action     *Action
	ticker     *time.Ticker
	done       chan bool
	url        string
	period     int
}

// EnrichType enrich the type info with the informations from the profile
func (p *PingCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return PingCommandTypeInfo, nil
}

// Init the command
func (p *PingCommand) Init(a *Action, commandName string) (bool, error) {
	p.action = a
	value, found := p.Parameters["url"]
	if found {
		url, ok := value.(string)
		if ok {
			p.url = url
		} else {
			return false, fmt.Errorf("The url parameter is in wrong format. Please use string as format")
		}
	} else {
		return false, fmt.Errorf("The url parameter is missing")
	}
	value, found = p.Parameters["period"]
	if found {
		period, ok := value.(int)
		if ok {
			p.period = period
		} else {
			return false, fmt.Errorf("The period parameter is in wrong format. Please use int as format")
		}
	} else {
		return false, fmt.Errorf("The period parameter is missing")
	}

	p.ticker = time.NewTicker(time.Duration(p.period) * time.Second)
	p.done = make(chan bool)
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				text := ""
				pingtime, err := p.getPingTime(p.url)
				if err != nil {
					text = fmt.Sprintf("error %v", err)
				}
				text = fmt.Sprintf("%.2fms", pingtime)
				message := models.Message{
					Profile: p.action.Profile,
					Action:  p.action.Name,
					Text:    text,
					State:   0,
				}
				api.SendMessage(message)
				continue
			}
		}
	}()
	return true, nil
}

// Stop the command
func (p *PingCommand) Stop(a *Action) (bool, error) {
	p.done <- true
	return true, nil
}

// Execute the command
func (p *PingCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	value, found := p.Parameters["url"]
	if found {
		url, ok := value.(string)
		if ok {
			pingtime, err := p.getPingTime(url)
			if err != nil {
				clog.Logger.Errorf("error: %v\r\n", err)
				return false, fmt.Errorf("Error executing the url. %v", err)
			}
			clog.Logger.Infof("Ping time: %.2fms", pingtime)
		} else {
			return false, fmt.Errorf("The url parameter is in wrong format. Please use string as format")
		}
	} else {
		return false, fmt.Errorf("The url parameter is missing")
	}
	return true, nil
}

func (p *PingCommand) getPingTime2(url string) (float64, error) {
	start := time.Now()
	resp, err := http.Get(url)
	stop := time.Now()
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return float64(stop.Sub(start).Nanoseconds()) / 1000000.0, nil
}

func (p *PingCommand) getPingTime(url string) (float64, error) {
	pinger, err := ping.NewPinger(url)
	if err != nil {
		panic(err)
	}
	pinger.SetPrivileged(true)

	pinger.Count = 1
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	clog.Logger.Debugf("ping %v", stats)
	return float64(stats.AvgRtt.Nanoseconds()) / 1000000.0, nil
}
