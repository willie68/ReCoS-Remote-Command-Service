// +build windows
package osdependent

import (
	"os"
	"time"

	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/hardware"
)

// InitOSDependend initialise windows depending components
func InitOSDependend(config config.Config) error {
	extConfig := config.ExternalConfig
	// initialise streamdeck
	err := hardware.InitStreamDeckInteg(extConfig, config)
	if err != nil {
		clog.Logger.Errorf("error initialising streamdeck. %v", err)
		os.Exit(1)
	}
	// initialise OpenHardwareMonitor component
	err = hardware.InitOpenHardwareMonitor(extConfig)
	if err != nil {
		clog.Logger.Errorf("error initialising open hardware connection. %v", err)
		os.Exit(1)
	}
	if hardware.OpenHardwareMonitorInstance.Active {
		go func() {
			for !hardware.OpenHardwareMonitorInstance.Connected {
				err := hardware.OpenHardwareMonitorInstance.Connect()
				if err != nil {
					clog.Logger.Errorf("error connectiong to app. trying later again. %v", err)
					time.Sleep(10 * time.Second)
				}
			}
		}()
	}
	return nil
}

func DisposeOSDependend() {
	hardware.DisposeStreamDeck()
}
