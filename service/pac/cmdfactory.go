package pac

import (
	clog "wkla.no-ip.biz/remote-desk-service/logging"
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
	SendMessageCommandTypeInfo,
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
	HMPrgCommandTypeInfo,
	HMSwitchCommandTypeInfo,
	HMDimmerCommandTypeInfo,
	ShowTextCommandTypeInfo,
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
			clog.Logger.Errorf("Error enrich command type: %v", err)
			continue
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
	case HMPrgCommandTypeInfo.Type:
		{
			cmdExecutor = &HMPrgCommand{
				Parameters: command.Parameters,
			}
		}
	case HMSwitchCommandTypeInfo.Type:
		{
			cmdExecutor = &HMSwitchCommand{
				Parameters: command.Parameters,
			}
		}
	case HMDimmerCommandTypeInfo.Type:
		{
			cmdExecutor = &HMDimmerCommand{
				Parameters: command.Parameters,
			}
		}
	case ShowTextCommandTypeInfo.Type:
		{
			cmdExecutor = &ShowTextCommand{
				Parameters: command.Parameters,
			}
		}
	}
	if cmdExecutor == nil {
		cmdExecutor = GetOSCommand(command)
	}
	return cmdExecutor
}
