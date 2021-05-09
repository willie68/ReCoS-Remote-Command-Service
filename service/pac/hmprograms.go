package pac

// homematic programs is the command to start a dedicated program on the homematic server
import (
	"errors"
	"fmt"

	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/smarthome"
)

// HMPrgCommandTypeInfo showing hardware sensor data
var HMPrgCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Smarthome",
	Type:             "HMPROGRAMS",
	Name:             "HomematicPrograms",
	Description:      "executoing a homematic program",
	Icon:             "history2.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "name",
			Type:           "string",
			Description:    "the homematic program to execute",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

type HMPrgCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	name        string
	prgID       string
}

// EnrichType enrich the type info with the informations from the profile
func (h *HMPrgCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	hm, ok := smarthome.GetHomematic()
	if !ok {
		return HMPrgCommandTypeInfo, errors.New("homematic not configured")
	}

	index := -1
	for x, parameter := range HMPrgCommandTypeInfo.Parameters {
		if parameter.Name == "name" {
			index = x
		}
	}
	if index >= 0 {
		programList, err := hm.ProgramList()
		if !ok {
			return HMPrgCommandTypeInfo, err
		}
		programs := programList.Programs
		HMPrgCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, program := range programs {
			HMPrgCommandTypeInfo.Parameters[index].List = append(HMPrgCommandTypeInfo.Parameters[index].List, program.Name)
		}
	}

	return HMPrgCommandTypeInfo, nil
}

// Init nothing
func (h *HMPrgCommand) Init(a *Action, commandName string) (bool, error) {
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
func (h *HMPrgCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (h *HMPrgCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
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
