// +build windows
package pac

import "wkla.no-ip.biz/remote-desk-service/pkg/models"

func InitOSCommand() {
	CommandTypes = append(CommandTypes, WindowCtrlCommandTypeInfo)
	CommandTypes = append(CommandTypes, HardwareMonitorCommandTypeInfo)
	CommandTypes = append(CommandTypes, OBSStartStopCommandTypeInfo)
	CommandTypes = append(CommandTypes, OBSProfileCommandTypeInfo)
	CommandTypes = append(CommandTypes, OBSSceneCollectionCommandTypeInfo)
	CommandTypes = append(CommandTypes, OBSSceneCommandTypeInfo)
}

// GetOSCommand return the command worker class responsible for executing the command definition
// THis one is only for windows specific commands
func GetOSCommand(command models.Command) CommandExecutor {
	var cmdExecutor CommandExecutor
	switch command.Type {
	case WindowCtrlCommandTypeInfo.Type:
		{
			cmdExecutor = &WindowCtrlCommand{
				Parameters: command.Parameters,
			}
		}
	case HardwareMonitorCommandTypeInfo.Type:
		{
			cmdExecutor = &HardwareMonitorCommand{
				Parameters: command.Parameters,
			}
		}
	case OBSStartStopCommandTypeInfo.Type:
		{
			cmdExecutor = &OBSStartStopCommand{
				Parameters: command.Parameters,
			}
		}
	case OBSProfileCommandTypeInfo.Type:
		{
			cmdExecutor = &OBSProfileCommand{
				Parameters: command.Parameters,
			}
		}
	case OBSSceneCollectionCommandTypeInfo.Type:
		{
			cmdExecutor = &OBSSceneCollectionCommand{
				Parameters: command.Parameters,
			}
		}
	case OBSSceneCommandTypeInfo.Type:
		{
			cmdExecutor = &OBSSceneCommand{
				Parameters: command.Parameters,
			}
		}
	}
	return cmdExecutor
}
