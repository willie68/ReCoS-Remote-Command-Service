package dto

import "wkla.no-ip.biz/remote-desk-service/pkg/models"

// GetCommand return the command worker class responsible for executing the command definition
func GetCommand(command models.Command) CommandExecutor {
	var cmdExecutor CommandExecutor
	switch command.Type {
	case models.Noop:
		{
			cmdExecutor = &NoopCommand{
				Parameters: command.Parameters,
			}
		}
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
	case models.PageCommand:
		{
			cmdExecutor = &PageCommand{
				Parameters: command.Parameters,
			}
		}
	case models.KeysCommand:
		{
			cmdExecutor = &KeysCommand{
				Parameters: command.Parameters,
			}
		}
	case models.Timer:
		{
			cmdExecutor = &TimerCommand{
				Parameters: command.Parameters,
			}
		}
	case models.Clock:
		{
			cmdExecutor = &ClockCommand{
				Parameters: command.Parameters,
			}
		}
	case models.Screenshot:
		{
			cmdExecutor = &ScreenshotCommand{
				Parameters: command.Parameters,
			}
		}
	}
	if cmdExecutor == nil {
		cmdExecutor = GetOSCommand(command)
	}
	return cmdExecutor
}
