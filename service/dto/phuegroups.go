package dto

import (
	"errors"
	"fmt"
	"image/color"
	"time"

	"wkla.no-ip.biz/remote-desk-service/pkg/lighting"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PHueGroupsCommandTypeInfo showing hardware sensor data
var PHueGroupsCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Lighting",
	Type:             "PHUEGROUPS",
	Name:             "PhilipsHueGroups",
	Description:      "control a hue group, like room or zone and get a feedback",
	Icon:             "light_bulb.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "group",
			Type:           "string",
			Description:    "the philips hue group to control",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "scene",
			Type:           "string",
			Description:    "the philips hue scene to âpply",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "brightness",
			Type:           "int",
			Description:    "brightness of the light",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "saturation",
			Type:           "int",
			Description:    "saturation of the light",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "hue",
			Type:           "int",
			Description:    "hue of the light, this is a color value ranging from 1..65535",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "colortemp",
			Type:           "int",
			Description:    "color temperatur of the light, this is a value ranging from 2000..6500",
			Unit:           " Kelvin",
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

type PHueGroupsCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	ticker      *time.Ticker
	done        chan bool
	group       string
	scene       string
	bright      int
	saturation  int
	hue         int
	colortemp   int
	color       color.Color
}

// EnrichType enrich the type info with the informations from the profile
func (d *PHueGroupsCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return PHueGroupsCommandTypeInfo, errors.New("philips Hue not configured")
	}

	index := -1
	for x, parameter := range PHueGroupsCommandTypeInfo.Parameters {
		if parameter.Name == "group" {
			index = x
		}
	}
	if index >= 0 {
		lights := hue.Groups
		PHueGroupsCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, light := range lights {
			PHueGroupsCommandTypeInfo.Parameters[index].List = append(PHueGroupsCommandTypeInfo.Parameters[index].List, light.Name)
		}
	}

	index = -1
	for x, parameter := range PHueGroupsCommandTypeInfo.Parameters {
		if parameter.Name == "scene" {
			index = x
		}
	}
	if index >= 0 {
		lights := hue.Scenes
		PHueGroupsCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, light := range lights {
			PHueGroupsCommandTypeInfo.Parameters[index].List = append(PHueGroupsCommandTypeInfo.Parameters[index].List, light.Name)
		}
	}

	return PHueGroupsCommandTypeInfo, nil
}

// Init nothing
func (p *PHueGroupsCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	p.action = a
	p.commandName = commandName
	p.group, err = ConvertParameter2String(p.Parameters, "group", "")
	if err != nil {
		return false, fmt.Errorf("the light parameter is in wrong format. Please use string as format")
	}
	p.scene, err = ConvertParameter2String(p.Parameters, "scene", "")
	if err != nil {
		return false, fmt.Errorf("the scene parameter is in wrong format. Please use string as format")
	}
	p.bright, err = ConvertParameter2Int(p.Parameters, "brightness", 254)
	if err != nil {
		return false, fmt.Errorf("the brightness parameter is in wrong format. Please use int as format")
	}
	p.saturation, err = ConvertParameter2Int(p.Parameters, "saturation", 254)
	if err != nil {
		return false, fmt.Errorf("the saturation parameter is in wrong format. Please use int as format")
	}
	p.hue, err = ConvertParameter2Int(p.Parameters, "hue", 0)
	if err != nil {
		return false, fmt.Errorf("the hue parameter is in wrong format. Please use int as format")
	}
	p.colortemp, err = ConvertParameter2Int(p.Parameters, "colortemp", 0)
	if err != nil {
		return false, fmt.Errorf("the colortemp parameter is in wrong format. Please use int as format")
	}
	p.color, err = ConvertParameter2Color(p.Parameters, "color", color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})
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
			}
		}
	}()
	return true, nil
}

// Stop nothing
func (p *PHueGroupsCommand) Stop(a *Action) (bool, error) {
	p.done <- true
	return true, nil
}

// Execute nothing
func (p *PHueGroupsCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return true, errors.New("philips hue not configured")
	}
	group, ok := hue.Group(p.group)
	if !ok {
		return true, fmt.Errorf("can't find light with name. %s", p.group)
	}
	if group.IsOn() {
		group.Off()
	} else {
		if p.scene != "" {
			group.Scene(p.scene)
		}
		if p.bright > 0 {
			group.Bri(uint8(p.bright))
		}
		if p.saturation > 0 {
			group.Sat(uint8(p.saturation))
		}
		if p.hue > 0 {
			group.Hue(uint16(p.hue))
		}
		if p.colortemp > 0 {
			ct := uint16(1000000 / p.colortemp)
			group.Ct(ct)
		}
		if p.color != nil {
			group.Col(p.color)
		}
		group.On()
	}
	return true, nil
}
