package pac

import (
	"errors"
	"fmt"
	"image/color"
	"strings"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/lighting"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// PHueLightsCommandTypeInfo showing hardware sensor data
var PHueLightsCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Lighting",
	Type:             "PHUELIGHTS",
	Name:             "PhilipsHueLights",
	Description:      "control a hue light and get a feedback",
	Icon:             "light_bulb.svg",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "name",
			Type:           "string",
			Description:    "the philips hue light to control",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
			GroupedList:    true,
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

type PHueLightsCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	ticker      *time.Ticker
	done        chan bool
	name        string
	bright      int
	saturation  int
	hue         int
	colortemp   int
	color       color.Color
}

var (
	hueColorModeHS      string = "hs"
	hueColorModeXY      string = "xy"
	hueColorModeCT      string = "ct"
	hueDefaultColorTemp int    = 4000
	hueLightPrefix      string = "Light"
	hueGroupPrefix      string = "Group"
)

// EnrichType enrich the type info with the informations from the profile
func (p *PHueLightsCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return PHueLightsCommandTypeInfo, errors.New("philips Hue not configured")
	}

	index := -1
	for x, parameter := range PHueLightsCommandTypeInfo.Parameters {
		if parameter.Name == "name" {
			index = x
		}
	}
	if index >= 0 {
		PHueLightsCommandTypeInfo.Parameters[index].List = make([]string, 0)
		lights := hue.Lights
		for _, light := range lights {
			PHueLightsCommandTypeInfo.Parameters[index].List = append(PHueLightsCommandTypeInfo.Parameters[index].List, p.buildLightName(light.Name))
		}
		groups := hue.Groups
		for _, group := range groups {
			PHueLightsCommandTypeInfo.Parameters[index].List = append(PHueLightsCommandTypeInfo.Parameters[index].List, p.buildGroupName(group.Name))
		}
	}

	return PHueLightsCommandTypeInfo, nil
}

// Init nothing
func (p *PHueLightsCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	p.action = a
	p.commandName = commandName
	p.name, err = ConvertParameter2String(p.Parameters, "name", "")
	if err != nil {
		return false, fmt.Errorf("the light parameter is in wrong format. Please use string as format")
	}
	p.bright, err = ConvertParameter2Int(p.Parameters, "brightness", 254)
	if err != nil {
		return false, fmt.Errorf("the brightness parameter is in wrong format. Please use int as format")
	}
	p.saturation, err = ConvertParameter2Int(p.Parameters, "saturation", 0)
	if err != nil {
		return false, fmt.Errorf("the saturation parameter is in wrong format. Please use int as format")
	}
	p.hue, err = ConvertParameter2Int(p.Parameters, "hue", 0)
	if err != nil {
		return false, fmt.Errorf("the hue parameter is in wrong format. Please use int as format")
	}
	p.colortemp, err = ConvertParameter2Int(p.Parameters, "colortemp", hueDefaultColorTemp)
	if err != nil {
		return false, fmt.Errorf("the colortemp parameter is in wrong format. Please use int as format")
	}
	p.color, err = ConvertParameter2Color(p.Parameters, "color", nil)
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
				var name string
				if p.isGroupName(p.name) {
					name = p.getName(p.name)
					on, err := hue.GroupIsOn(name)
					if err != nil {
						text = fmt.Sprintf("error getting light with name: %s", p.name)
					} else {
						if on {
							text = "group \"%s\" is on"
						} else {
							text = "group \"%s\" is off"
						}
					}
				} else {
					name = p.getName(p.name)
					on, err := hue.LightIsOn(name)
					if err != nil {
						text = fmt.Sprintf("error getting light with name: %s", p.name)
					} else {
						if on {
							text = "light \"%s\" on"
						} else {
							text = "light \"%s\" off"
						}
					}
				}
				text = fmt.Sprintf(text, name)
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

	name := p.getName(p.name)

	if p.isGroupName(p.name) {
		return p.doGroupWork(name)
	} else {
		return p.doLightWork(name)
	}
}

func (p *PHueLightsCommand) doGroupWork(name string) (bool, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return true, errors.New("philips hue not configured")
	}
	group, ok := hue.Group(name)
	if !ok {
		return true, fmt.Errorf("can't find group with name. %s", name)
	}
	clog.Logger.Infof("light colormode: %s", group.State.ColorMode)
	colorMode := p.getMyColormode()

	if (colorMode == group.State.ColorMode) && group.IsOn() {
		group.Off()
	} else {
		if colorMode == hueColorModeXY {
			if p.color != nil {
				group.Col(p.color)
				return true, nil
			}
		}
		if p.bright > 0 {
			group.Bri(uint8(p.bright))
		}
		if colorMode == hueColorModeHS {
			if p.saturation > 0 {
				group.Sat(uint8(p.saturation))
			}
			if p.hue > 0 {
				group.Hue(uint16(p.hue))
			}
		}
		if colorMode == hueColorModeCT {
			ct := uint16(250)
			if p.colortemp > 0 {
				ct = uint16(kelvin2Mired(p.colortemp))
			}
			group.Ct(ct)
		}
	}
	return true, nil
}

func (p *PHueLightsCommand) doLightWork(name string) (bool, error) {
	hue, ok := lighting.GetPhilipsHue()
	if !ok {
		return true, errors.New("philips hue not configured")
	}
	light, ok := hue.Light(name)
	if !ok {
		return true, fmt.Errorf("can't find light with name. %s", name)
	}
	clog.Logger.Infof("light colormode: %s", light.State.ColorMode)
	colorMode := p.getMyColormode()

	if (colorMode == light.State.ColorMode) && light.IsOn() {
		light.Off()
	} else {
		if colorMode == hueColorModeXY {
			if p.color != nil {
				light.Col(p.color)
				return true, nil
			}
		}
		if p.bright > 0 {
			light.Bri(uint8(p.bright))
		}
		if colorMode == hueColorModeHS {
			if p.saturation > 0 {
				light.Sat(uint8(p.saturation))
			}
			if p.hue > 0 {
				light.Hue(uint16(p.hue))
			}
		}
		if colorMode == hueColorModeCT {
			ct := uint16(250)
			if p.colortemp > 0 {
				ct = uint16(kelvin2Mired(p.colortemp))
			}
			light.Ct(ct)
		}
	}
	return true, nil
}

func kelvin2Mired(kelvin int) int {
	return 1000000 / kelvin
}

func (p *PHueLightsCommand) getMyColormode() string {
	if p.saturation > 0 {
		return hueColorModeHS
	}
	if p.hue > 0 {
		return hueColorModeHS
	}
	if p.color != nil {
		return hueColorModeXY
	}
	// this is the default, thats why it's the latest to test.
	if p.colortemp > 0 {
		return hueColorModeCT
	}
	return hueColorModeCT
}

func (p *PHueLightsCommand) buildLightName(name string) string {
	return fmt.Sprintf("%s: %s", hueLightPrefix, name)
}

func (p *PHueLightsCommand) buildGroupName(name string) string {
	return fmt.Sprintf("%s: %s", hueGroupPrefix, name)
}

func (p *PHueLightsCommand) isGroupName(name string) bool {
	return strings.HasPrefix(name, hueGroupPrefix)
}

func (p *PHueLightsCommand) getName(name string) string {
	return strings.TrimSpace(name[strings.Index(name, ":")+1:])
}
