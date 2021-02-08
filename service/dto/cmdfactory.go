package dto

import "wkla.no-ip.biz/remote-desk-service/pkg/models"

// GetCommand return the command worker class responsible for executing the command definition
func GetCommand(command models.Command) CommandExecutor {
	var cmdExecutor CommandExecutor
	switch command.Type {
	case models.Delay:
		{
			cmdExecutor = &DelayCommand{
				Parameters: command.Parameters,
			}
		}
	case models.Execute:
		{
			cmdExecutor = &ExecuteCommand{
				Parameters: command.Parameters,
			}
		}
	}
	return cmdExecutor
}
