package pac

// homematic setting a parameter with a float value
import (
	"errors"
	"fmt"

	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/smarthome"
)

// HMDimmerCommandTypeInfo showing hardware sensor data
var HMDimmerCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Smarthome",
	Type:             "HMDIMMER",
	Name:             "HomematicDimmer",
	Description:      "activating a homematic dimmer",
	Icon:             "bulb.png",
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
	},
}

type HMDimmerCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	name        string
	prgID       string
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
		programList, err := hm.ProgramList()
		if !ok {
			return HMDimmerCommandTypeInfo, err
		}
		programs := programList.Programs
		HMDimmerCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, program := range programs {
			HMDimmerCommandTypeInfo.Parameters[index].List = append(HMDimmerCommandTypeInfo.Parameters[index].List, program.Name)
		}
	}

	return HMDimmerCommandTypeInfo, nil
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

	hm, ok := smarthome.GetHomematic()
	if !ok {
		return false, errors.New("homematic not configured")
	}
	programList, err := hm.ProgramList()
	if err != nil {
		return false, err
	}
	programs := programList.Programs
	for _, program := range programs {
		if program.Name == h.name {
			h.prgID = program.ID
		}
	}
	if h.prgID == "" {
		return true, fmt.Errorf("can't find program with name %s", h.name)
	}
	return true, nil
}

// Stop nothing
func (h *HMDimmerCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (h *HMDimmerCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if h.prgID == "" {
		return true, fmt.Errorf("can't find program with name %s", h.name)
	}
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return true, errors.New("homematic not configured")
	}
	_, err := hm.RunProgram(h.prgID)
	if err != nil {
		return true, err
	}
	return true, nil
}
