// +build !windows
package dto

import "wkla.no-ip.biz/remote-desk-service/pkg/models"

// GetOSCommand return the command worker class responsible for executing the command definition
func GetOSCommand(command models.Command) CommandExecutor {
}