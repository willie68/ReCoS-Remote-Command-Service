package autostart

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

func init() {
}

func (a *App) IsEnabled() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		clog.Logger.Errorf("Error opening registry: %v", err)
		return false
	}
	defer k.Close()
	s, _, err := k.GetStringValue("ReCoS_Service")
	if err != nil {
		clog.Logger.Errorf("Error getting key: %v", err)
		return false
	}
	fmt.Printf("ReCoS Service is in %q\n", s)
	return s != ""
}

func (a *App) Enable() error {
	path := a.Exec[0]
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		clog.Logger.Errorf("Error opening registry: %v", err)
		return err
	}
	defer k.Close()

	err = k.SetStringValue("ReCoS_Service", path)
	if err != nil {
		clog.Logger.Errorf("Error setting key: %v", err)
		return err
	}
	return nil
}

func (a *App) Disable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		clog.Logger.Errorf("Error opening registry: %v", err)
		return err
	}
	defer k.Close()

	err = k.DeleteValue("ReCoS_Service")
	if err != nil {
		clog.Logger.Errorf("Error deleting key: %v", err)
		return err
	}
	return nil
}
