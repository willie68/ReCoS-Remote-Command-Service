// +build windows
package osdependent

import (
	"wkla.no-ip.biz/remote-desk-service/config"
	"wkla.no-ip.biz/remote-desk-service/pkg/hardware"
)

// InitOSDependend initialise windows depending components
func InitOSDependend(config config.Config) error {
	extConfig := config.ExternalConfig
	// initialise OpenHardwareMonitor component
	err := hardware.InitOpenHardwareMonitor(extConfig)
	if err != nil {
		return err
	}
	return nil
}
