package dto

import (
	"errors"
	"fmt"
	"strconv"

	"wkla.no-ip.biz/remote-desk-service/pkg/lighting"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PHueScenesCommandTypeInfo showing hardware sensor data
var PHueScenesCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Lighting",
	Type:             "PHUESCENES",
	Name:             "PhilipsHueScenes",
	Description:      "applying a scene to a hue group",
	Icon:             "audio_equalizer.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "name",
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
			FilteredList:   "name",
		},
		{
			Name:           "brightness",
			Type:           "int",
			Description:    "brightness of the light",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

type PHueScenesCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	group       string
	scene       string
	bright      int
}

var hueColorModeScene string = "scene"

// EnrichType enrich the type info with the informations from the profile
func (d *PHueScenesCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return PHueScenesCommandTypeInfo, errors.New("philips Hue not configured")
	}

	groupMap := make(map[string]string, 0)
	index := -1
	for x, parameter := range PHueScenesCommandTypeInfo.Parameters {
		if parameter.Name == "name" {
			index = x
		}
	}
	if index >= 0 {
		groups := hue.Groups
		PHueScenesCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, group := range groups {
			PHueScenesCommandTypeInfo.Parameters[index].List = append(PHueScenesCommandTypeInfo.Parameters[index].List, group.Name)
			groupMap[strconv.Itoa(group.ID)] = group.Name
		}
	}

	index = -1
	for x, parameter := range PHueScenesCommandTypeInfo.Parameters {
		if parameter.Name == "scene" {
			index = x
		}
	}
	if index >= 0 {
		scenes := hue.Scenes
		PHueScenesCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, scene := range scenes {
			groupName, ok := groupMap[scene.Group]
			if !ok {
				continue
			}
			name := fmt.Sprintf("%s: %s", groupName, scene.Name)
			PHueScenesCommandTypeInfo.Parameters[index].List = append(PHueScenesCommandTypeInfo.Parameters[index].List, name)
		}
	}

	return PHueScenesCommandTypeInfo, nil
}

// Init nothing
func (p *PHueScenesCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	p.action = a
	p.commandName = commandName
	p.group, err = ConvertParameter2String(p.Parameters, "name", "")
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
	return true, nil
}

// Stop nothing
func (p *PHueScenesCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (p *PHueScenesCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return true, errors.New("philips hue not configured")
	}
	group, ok := hue.Group(p.group)
	if !ok {
		return true, fmt.Errorf("can't find group with name. %s", p.group)
	}

	if p.bright > 0 {
		group.Bri(uint8(p.bright))
	}
	if p.scene != "" {
		scene, ok := hue.SceneForGroup(*group, p.scene)
		if !ok {
			return true, fmt.Errorf("scene \"%s\" not found", p.scene)
		}
		group.Scene(scene.ID)
	}
	return true, nil
}
