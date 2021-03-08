package dto

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var CommandTypes = []models.CommandTypeInfo{
	DelayCommandTypeInfo,
	ExecuteCommandTypeInfo,
	PageCommandTypeInfo,
	KeysCommandTypeInfo,
	NoopCommandTypeInfo,
	TimerCommandTypeInfo,
	ClockCommandTypeInfo,
	StopwatchCommandTypeInfo,
	ScreenshotCommandTypeInfo,
}

func InitCommand() {
	InitOSCommand()
}

// GetCommand return the command worker class responsible for executing the command definition
func GetCommand(command models.Command) CommandExecutor {
	var cmdExecutor CommandExecutor
	switch command.Type {
	case NoopCommandTypeInfo.Type:
		{
			cmdExecutor = &NoopCommand{
				Parameters: command.Parameters,
			}
		}
	case DelayCommandTypeInfo.Type:
		{
			cmdExecutor = &DelayCommand{
				Parameters: command.Parameters,
			}
		}
	case ExecuteCommandTypeInfo.Type:
		{
			cmdExecutor = &ExecuteCommand{
				Parameters: command.Parameters,
			}
		}
	case PageCommandTypeInfo.Type:
		{
			cmdExecutor = &PageCommand{
				Parameters: command.Parameters,
			}
		}
	case KeysCommandTypeInfo.Type:
		{
			cmdExecutor = &KeysCommand{
				Parameters: command.Parameters,
			}
		}
	case TimerCommandTypeInfo.Type:
		{
			cmdExecutor = &TimerCommand{
				Parameters: command.Parameters,
			}
		}
	case ClockCommandTypeInfo.Type:
		{
			cmdExecutor = &ClockCommand{
				Parameters: command.Parameters,
			}
		}
	case StopwatchCommandTypeInfo.Type:
		{
			cmdExecutor = &StopwatchCommand{
				Parameters: command.Parameters,
			}
		}
	case ScreenshotCommandTypeInfo.Type:
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
