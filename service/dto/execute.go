package dto

import (
	"fmt"
	"os/exec"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ExecuteCommandTypeInfo start an application or shell script and optionally waits for it's finishing
var ExecuteCommandTypeInfo = models.CommandTypeInfo{
	Type:             "EXECUTE",
	Name:             "Execute",
	Description:      "Execute an application",
	Icon:             "flash.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "command",
			Type:           "string",
			Description:    "the command to execute",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "waitOnClose",
			Type:           "bool",
			Description:    "wait's till the command is finnished",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "args",
			Type:           "[]string",
			Description:    "a list of arguments",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	}}

// ExecuteCommand is a command to execute a program or batch file.
// Using "command" for getting the command line to execute.
// Using "args" for optional parameters
type ExecuteCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (e *ExecuteCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return ExecuteCommandTypeInfo, nil
}

// Init the command
func (e *ExecuteCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop the command
func (e *ExecuteCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (e *ExecuteCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	waitOnExit := true
	value, found := e.Parameters["waitOnClose"]
	if found {
		waitValue, ok := value.(bool)
		if ok {
			waitOnExit = waitValue
		}
	}
	value, found = e.Parameters["command"]
	if found {
		cmdValue, ok := value.(string)
		if ok {
			cmd := exec.Command(cmdValue)

			argsValue, found := e.Parameters["args"]
			var args []string
			if found {
				args = make([]string, 0)
				args = append(args, cmd.Args[0])
				for _, argValue := range argsValue.([]interface{}) {
					args = append(args, argValue.(string))
				}
			}
			cmd.Args = args

			clog.Logger.Debugf("execute command line: %v", cmd.String())

			if waitOnExit {
				err := cmd.Run()
				if err != nil {
					clog.Logger.Errorf("error: %v\r\n", err)
				}
			} else {
				err := cmd.Start()
				if err != nil {
					clog.Logger.Errorf("error: %v\r\n", err)
				}
				go func() {
					err = cmd.Wait()
					if err != nil {
						clog.Logger.Errorf("error: %v\r\n", err)
					}
				}()
			}
		} else {
			return false, fmt.Errorf("The command parameter is in wrong format. Please use string as format")
		}
	} else {
		return false, fmt.Errorf("The command parameter is missing")
	}
	return true, nil
}
