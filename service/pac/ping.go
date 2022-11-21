package pac

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
	"wkla.no-ip.biz/remote-desk-service/internal/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PingCommandTypeInfo start an browser with directly with a url or filepath
var PingCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Network",
	Type:             "PING",
	Name:             "Ping",
	Description:      "Execute a ping to a url and show the answer in ms",
	Icon:             "world_shipping.svg",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.ParamInfo{
		{
			Name:           "name",
			Type:           "string",
			Description:    "the name to ping to",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "period",
			Type:           "int",
			Description:    "period in secods",
			Unit:           " Seconds",
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
	name       string
	period     int
}

// EnrichType enrich the type info with the informations from the profile
func (p *PingCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return PingCommandTypeInfo, nil
}

// Init the command
func (p *PingCommand) Init(a *Action, commandName string) (bool, error) {
	p.action = a

	name, err := ConvertParameter2String(p.Parameters, "name", "")
	if err != nil {
		return false, err
	}
	if name == "" {
		return false, fmt.Errorf("the name parameter is missing")
	}
	p.name = name

	valueInt, err := ConvertParameter2Int(p.Parameters, "period", 10)
	if err != nil {
		return false, fmt.Errorf("the period parameter is in wrong format. Please use int as format")
	}
	p.period = valueInt

	if p.period > 0 {
		p.ticker = time.NewTicker(time.Duration(p.period) * time.Second)
		p.done = make(chan bool)
		go func() {
			for {
				select {
				case <-p.done:
					return
				case <-p.ticker.C:
					text := ""
					pingtime, err := p.getPingTime(p.name)
					if err != nil {
						text = fmt.Sprintf("error %v", err)
						clog.Logger.Errorf("error: %v\r\n", err)
					} else {
						text = fmt.Sprintf("%.2fms", pingtime)
					}
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
	}
	return true, nil
}

// Stop the command
func (p *PingCommand) Stop(a *Action) (bool, error) {
	p.done <- true
	return true, nil
}

// Execute the command
func (p *PingCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	pingtime, err := p.getPingTime(p.name)
	if err != nil {
		clog.Logger.Errorf("error: %v\r\n", err)
		return false, fmt.Errorf("error executing the url. %v", err)
	}
	clog.Logger.Infof("Ping time: %.2fms", pingtime)
	return true, nil
}

func (p *PingCommand) getPingTime(url string) (float64, error) {
	pinger, err := ping.NewPinger(url)
	if err != nil {
		return 0, err
	}
	pinger.SetPrivileged(true)

	pinger.Count = 1
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	clog.Logger.Debugf("ping %v", stats)
	return float64(stats.AvgRtt.Nanoseconds()) / 1000000.0, nil
}
