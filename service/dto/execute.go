package dto

import (
	"fmt"
	"os/exec"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

// ExecuteCommand is a command to execute a program or batch file.
// Using "command" for getting the command line to execute.
// Using "args" for optional parameters
type ExecuteCommand struct {
	Parameters map[string]interface{}
}

// Execute the command
func (e *ExecuteCommand) Execute() (bool, error) {
	value, found := e.Parameters["command"]
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

			if output, err := cmd.Output(); err != nil {
				clog.Logger.Errorf("error: %v\r\n", err)
			} else {
				clog.Logger.Debug(string(output))
			}
		} else {
			return false, fmt.Errorf("The command parameter is in wrong format. Please use string as format")
		}
	} else {
		return false, fmt.Errorf("The command parameter is missing")
	}
	return true, nil
}
