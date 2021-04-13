package dto

import (
	"errors"
	"fmt"
	"image/color"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/lighting"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// HardwareMonitorCommandTypeInfo showing hardware sensor data
var PHueLightsCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Lighting",
	Type:             "PHUELIGHTS",
	Name:             "PhilipsHueLights",
	Description:      "control a hue light and get a feedback",
	Icon:             "light_bulb.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "light",
			Type:           "string",
			Description:    "the philips hue light to control",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "brightness",
			Type:           "integer",
			Description:    "brightness of the light",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "color",
			Type:           "color",
			Description:    "color of the light",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

type PHueLightsCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	ticker      *time.Ticker
	done        chan bool
	light       string
	bright      int
	color       color.Color
}

// EnrichType enrich the type info with the informations from the profile
func (d *PHueLightsCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return PHueLightsCommandTypeInfo, errors.New("philips Hue not configured")
	}

	index := -1
	for x, parameter := range PHueLightsCommandTypeInfo.Parameters {
		if parameter.Name == "light" {
			index = x
		}
	}
	if index >= 0 {
		lights := hue.Lights
		PHueLightsCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, light := range lights {
			PHueLightsCommandTypeInfo.Parameters[index].List = append(PHueLightsCommandTypeInfo.Parameters[index].List, light.Name)
		}
	}

	return PHueLightsCommandTypeInfo, nil
}

// Init nothing
func (p *PHueLightsCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	p.action = a
	p.commandName = commandName
	p.light, err = ConvertParameter2String(p.Parameters, "light", "")
	if err != nil {
		return false, fmt.Errorf("the light parameter is in wrong format. Please use string as format")
	}
	p.color, err = ConvertParameter2Color(p.Parameters, "color", color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF})
	if err != nil {
		return false, fmt.Errorf("the color parameter is in wrong format. Please use string as format")
	}

	p.ticker = time.NewTicker(1 * time.Second)
	p.done = make(chan bool)
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				text := ""
				hue, ok := lighting.GetPhilipsHue()
				if !ok {
					text = "philips hue not configured"
				}
				on, err := hue.LightIsOn(p.light)
				if err != nil {
					text = fmt.Sprintf("error getting light with name: %s", p.light)
				}
				if on {
					text = "light is on"
				} else {
					text = "light is off"
				}
				message := models.Message{
					Profile: a.Profile,
					Action:  a.Name,
					Command: p.commandName,
					Text:    text,
					State:   0,
				}
				api.SendMessage(message)
			}
		}
	}()
	return true, nil
}

// Stop nothing
func (p *PHueLightsCommand) Stop(a *Action) (bool, error) {
	p.done <- true
	return true, nil
}

// Execute nothing
func (p *PHueLightsCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return true, errors.New("philips hue not configured")
	}
	hue.ToggleLight(p.light)
	return true, nil
}
