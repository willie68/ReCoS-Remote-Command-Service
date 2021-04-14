package dto

import (
	"reflect"

	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var typeRegistry = make(map[string]reflect.Type)

func RegisterType(typedNil interface{}) {
	t := reflect.TypeOf(typedNil).Elem()
	typeRegistry[t.PkgPath()+"."+t.Name()] = t
}

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
	AudioVolumeCommandTypeInfo,
	MediaPlayCommandTypeInfo,
	DaysRemainCommandTypeInfo,
	BrowseCommandTypeInfo,
	PingCommandTypeInfo,
	CounterCommandTypeInfo,
	DiceCommandTypeInfo,
	RndWordsCommandTypeInfo,
	PlayAudioCommandTypeInfo,
	PHueLightsCommandTypeInfo,
	PHueScenesCommandTypeInfo,
}

func InitCommand() {
	InitOSCommand()
}

func EnrichTypes(types []models.CommandTypeInfo, profile models.Profile) ([]models.CommandTypeInfo, error) {
	localTypes := make([]models.CommandTypeInfo, 0)
	for _, commandType := range types {
		command := models.Command{
			Type: commandType.Type,
		}
		commandExecutor := GetCommand(command)
		newType, err := commandExecutor.EnrichType(profile)
		if err != nil {
			return nil, err
		}
		localTypes = append(localTypes, newType)
	}
	return localTypes, nil
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
	case AudioVolumeCommandTypeInfo.Type:
		{
			cmdExecutor = &AudioVolumeCommand{
				Parameters: command.Parameters,
			}
		}
	case MediaPlayCommandTypeInfo.Type:
		{
			cmdExecutor = &MediaPlayCommand{
				Parameters: command.Parameters,
			}
		}
	case DaysRemainCommandTypeInfo.Type:
		{
			cmdExecutor = &DaysRemainCommand{
				Parameters: command.Parameters,
			}
		}
	case BrowseCommandTypeInfo.Type:
		{
			cmdExecutor = &BrowseCommand{
				Parameters: command.Parameters,
			}
		}
	case PingCommandTypeInfo.Type:
		{
			cmdExecutor = &PingCommand{
				Parameters: command.Parameters,
			}
		}
	case CounterCommandTypeInfo.Type:
		{
			cmdExecutor = &CounterCommand{
				Parameters: command.Parameters,
			}
		}
	case DiceCommandTypeInfo.Type:
		{
			cmdExecutor = &DiceCommand{
				Parameters: command.Parameters,
			}
		}
	case RndWordsCommandTypeInfo.Type:
		{
			cmdExecutor = &RndWordsCommand{
				Parameters: command.Parameters,
			}
		}
	case SendMessageCommandTypeInfo.Type:
		{
			cmdExecutor = &SendMessageCommand{
				Parameters: command.Parameters,
			}
		}
	case PlayAudioCommandTypeInfo.Type:
		{
			cmdExecutor = &PlayAudioCommand{
				Parameters: command.Parameters,
			}
		}
	case PHueLightsCommandTypeInfo.Type:
		{
			cmdExecutor = &PHueLightsCommand{
				Parameters: command.Parameters,
			}
		}
	case PHueScenesCommandTypeInfo.Type:
		{
			cmdExecutor = &PHueScenesCommand{
				Parameters: command.Parameters,
			}
		}
	}
	if cmdExecutor == nil {
		cmdExecutor = GetOSCommand(command)
	}
	return cmdExecutor
}
