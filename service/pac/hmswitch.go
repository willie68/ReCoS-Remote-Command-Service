package pac

// homematic programs is the command to start a dedicated program on the homematic server
import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/smarthome"
)

// HMSwitchCommandTypeInfo showing hardware sensor data
var HMSwitchCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Smarthome",
	Type:             "HMSWITCH",
	Name:             "HomematicSwitch",
	Description:      "activating a homematic switch",
	Icon:             "light_bulb.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "name",
			Type:           "string",
			Description:    "the homematic switches to use",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
			GroupedList:    true,
		},
		{
			Name:           "onicon",
			Type:           "icon",
			Description:    "the icon used to show the on state",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "officon",
			Type:           "icon",
			Description:    "the icon used to show the off state",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

type HMSwitchCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	ticker      *time.Ticker
	done        chan bool

	name    string
	onIcon  string
	offIcon string
	iseID   string
}

// EnrichType enrich the type info with the informations from the profile
func (h *HMSwitchCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return HMSwitchCommandTypeInfo, errors.New("homematic not configured")
	}

	index := -1
	for x, parameter := range HMSwitchCommandTypeInfo.Parameters {
		if parameter.Name == "name" {
			index = x
		}
	}
	if index >= 0 {
		deviceList, err := hm.DeviceList()
		if !ok {
			return HMSwitchCommandTypeInfo, err
		}
		devices := deviceList.Devices
		HMSwitchCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, device := range devices {
			channels := device.Channels
			for _, channel := range channels {
				if strings.ToLower(channel.Direction) == "receiver" {
					HMSwitchCommandTypeInfo.Parameters[index].List = append(HMSwitchCommandTypeInfo.Parameters[index].List, h.buildSwitchName(device.Name, channel.Name))
				}
			}
		}
	}

	return HMSwitchCommandTypeInfo, nil
}

func (p *HMSwitchCommand) buildSwitchName(device, name string) string {
	return fmt.Sprintf("%s: %s", device, name)
}

// Init nothing
func (h *HMSwitchCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	h.action = a
	h.commandName = commandName
	h.name, err = ConvertParameter2String(h.Parameters, "name", "")
	if err != nil {
		return false, err
	}

	h.onIcon, err = ConvertParameter2String(h.Parameters, "onicon", "")
	if err != nil {
		return false, err
	}
	h.offIcon, err = ConvertParameter2String(h.Parameters, "officon", "")
	if err != nil {
		return false, err
	}

	hm, ok := smarthome.GetHomematic()
	if !ok {
		return false, errors.New("homematic not configured")
	}
	deviceList, err := hm.DeviceList()
	if err != nil {
		return false, err
	}
	devices := deviceList.Devices
	for _, device := range devices {
		channels := device.Channels
		for _, channel := range channels {
			deviceChannel := h.buildSwitchName(device.Name, channel.Name)
			if h.name == deviceChannel {
				h.iseID = channel.Ise_id
			}
		}
	}
	if h.iseID == "" {
		return true, fmt.Errorf("can't find device/channel with name %s", h.name)
	}
	h.ticker = time.NewTicker(time.Duration(hm.UpdatePeriod()) * time.Second)
	h.done = make(chan bool)
	go func() {
		for {
			select {
			case <-h.done:
				return
			case <-h.ticker.C:
				value, err := h.getState()
				if err != nil {
					clog.Logger.Errorf("can't get state of device/channel with name %s", h.name)
				}
				text := fmt.Sprintf("the %s is %+v", h.action.Config.Title, value)
				icon := ""
				if value && h.onIcon != "" {
					icon = h.onIcon
				}
				if !value && h.offIcon != "" {
					icon = h.offIcon
				}

				message := models.Message{
					Profile: h.action.Profile,
					Action:  h.action.Name,
					Text:    text,
					State:   0,
				}
				if (h.onIcon != "") || (h.offIcon != "") {
					message.ImageURL = icon
				}
				api.SendMessage(message)
			}
		}
	}()

	return true, nil
}

// Stop nothing
func (h *HMSwitchCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (h *HMSwitchCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if h.iseID == "" {
		return true, fmt.Errorf("can't find device/channel with name %s", h.name)
	}
	value, err := h.getState()
	if err != nil {
		return true, fmt.Errorf("can't get state of device/channel with name %s", h.name)
	}
	value = !value
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return true, errors.New("homematic not configured")
	}
	if value {
		hm.ChangeState(h.iseID, 1.0)
	} else {
		hm.ChangeState(h.iseID, 0.0)
	}
	return true, nil
}

func (h *HMSwitchCommand) getState() (bool, error) {
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return true, errors.New("homematic not configured")
	}
	datapoints, err := hm.State(h.iseID)
	if err != nil {
		return false, err
	}
	for _, datapoint := range datapoints {
		if strings.ToUpper(datapoint.Type) == "STATE" {
			switch datapoint.ValueType {
			case 2:
				value, _ := strconv.ParseBool(datapoint.Value)
				return value, nil
			case 4:
				value, _ := strconv.ParseFloat(datapoint.Value, 64)
				return value >= 0.5, nil
			case 16:
				value, _ := strconv.ParseInt(datapoint.Value, 10, 64)
				return value != 0, nil
			}
		}
	}
	jsonStr, _ := json.Marshal(datapoints)
	clog.Logger.Infof("found datapoints: %s", jsonStr)
	return false, fmt.Errorf("can't find device/channel value with name %s", h.name)
}
