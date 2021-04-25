package pac

// homematic setting a parameter with a float value
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

// HMDimmerCommandTypeInfo showing hardware sensor data
var HMDimmerCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Smarthome",
	Type:             "HMDIMMER",
	Name:             "HomematicDimmer",
	Description:      "activating a homematic dimmer",
	Icon:             "light_bulb.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "name",
			Type:           "string",
			Description:    "the homematic dimm value to use",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "action",
			Type:           "string",
			Description:    "the action to do with this device/channel",
			Unit:           "",
			WizardPossible: true,
			List:           []string{"set value", "up", "down"},
		},
		{
			Name:           "value",
			Type:           "int",
			Description:    "the value to set or the step interval (in percent)",
			Unit:           "%",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

type HMDimmerCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	ticker      *time.Ticker
	done        chan bool

	name  string
	cmd   string
	value int
	iseID string
}

// EnrichType enrich the type info with the informations from the profile
func (h *HMDimmerCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return HMDimmerCommandTypeInfo, errors.New("homematic not configured")
	}

	index := -1
	for x, parameter := range HMDimmerCommandTypeInfo.Parameters {
		if parameter.Name == "name" {
			index = x
		}
	}
	if index >= 0 {
		deviceList, err := hm.DeviceList()
		if !ok {
			return HMDimmerCommandTypeInfo, err
		}
		devices := deviceList.Devices
		HMDimmerCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, device := range devices {
			channels := device.Channels
			for _, channel := range channels {
				if strings.ToLower(channel.Direction) == "receiver" {
					HMDimmerCommandTypeInfo.Parameters[index].List = append(HMDimmerCommandTypeInfo.Parameters[index].List, h.buildName(device.Name, channel.Name))
				}
			}
		}
	}

	return HMDimmerCommandTypeInfo, nil
}

func (h *HMDimmerCommand) buildName(device, name string) string {
	return fmt.Sprintf("%s: %s", device, name)
}

// Init nothing
func (h *HMDimmerCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	h.action = a
	h.commandName = commandName
	h.name, err = ConvertParameter2String(h.Parameters, "name", "")
	if err != nil {
		return false, err
	}
	h.cmd, err = ConvertParameter2String(h.Parameters, "action", "")
	if err != nil {
		return false, err
	}
	h.value, err = ConvertParameter2Int(h.Parameters, "value", 50)
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
			deviceChannel := h.buildName(device.Name, channel.Name)
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
				perCent := int(value * 100)
				text := fmt.Sprintf("the %s is %d", h.action.Config.Title, perCent)

				message := models.Message{
					Profile: h.action.Profile,
					Action:  h.action.Name,
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
func (h *HMDimmerCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (h *HMDimmerCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if h.iseID == "" {
		return true, fmt.Errorf("can't find device/channel with name %s", h.name)
	}
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return true, errors.New("homematic not configured")
	}
	if h.cmd == "" {
		return true, nil
	}
	if h.cmd == "set value" {
		value := float64(h.value) / 100
		hm.ChangeState(h.iseID, value)
	} else {
		datas, err := hm.State(h.iseID)
		if err != nil {
			return true, err
		}
		value := -1.0
		for _, datapoint := range datas {
			if strings.ToUpper(datapoint.Type) == "LEVEL" {
				switch datapoint.ValueType {
				case 2:
					v, _ := strconv.ParseBool(datapoint.Value)
					if v {
						value = 1.0
					} else {
						value = 0.0
					}
				case 4:
					value, _ = strconv.ParseFloat(datapoint.Value, 64)
				case 16:
					v, _ := strconv.ParseInt(datapoint.Value, 10, 64)
					value = float64(v)
				}
			}
		}
		if value < 0 {
			return true, errors.New("datapoint not found")
		}
		diff := float64(h.value) / 100.0
		if h.cmd == "up" {
			value = value + diff
			if value > 1.0 {
				value = 1.0
			}
		} else {
			value = value - diff
			if value < 0 {
				value = 0.0
			}
		}
		hm.ChangeState(h.iseID, value)
	}
	return true, nil
}

func (h *HMDimmerCommand) getState() (float64, error) {
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return 0.0, errors.New("homematic not configured")
	}
	datapoints, err := hm.State(h.iseID)
	if err != nil {
		return 0.0, err
	}
	for _, datapoint := range datapoints {
		if strings.ToUpper(datapoint.Type) == "LEVEL" {
			switch datapoint.ValueType {
			case 2:
				value, _ := strconv.ParseBool(datapoint.Value)
				if value {
					return 1.0, nil
				} else {
					return 0.0, nil
				}
			case 4:
				value, _ := strconv.ParseFloat(datapoint.Value, 64)
				return value, nil
			case 16:
				value, _ := strconv.ParseInt(datapoint.Value, 10, 64)
				return float64(value), nil
			}
		}
	}
	jsonStr, _ := json.Marshal(datapoints)
	clog.Logger.Infof("found datapoints: %s", jsonStr)
	return 0.0, fmt.Errorf("can't find device/channel value with name %s", h.name)
}
