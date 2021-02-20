// +build windows
package dto

import "wkla.no-ip.biz/remote-desk-service/pkg/models"

// GetOSCommand return the command worker class responsible for executing the command definition
// THis one is only for windows specific commands
func GetOSCommand(command models.Command) CommandExecutor {
	var cmdExecutor CommandExecutor
	switch command.Type {
	case models.WindowCtrlCommand:
		{
			cmdExecutor = &WindowCtrlCommand{
				Parameters: command.Parameters,
			}
		}
	case models.HardwareMonitorCommand:
		{
			cmdExecutor = &HardwareMonitorCommand{
				Parameters: command.Parameters,
			}
		}
	}
	return cmdExecutor
}
