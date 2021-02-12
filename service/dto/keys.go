package dto

import (
	"fmt"
	"runtime"
	"time"

	"github.com/micmonay/keybd_event"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

// KeysCommand is a command to switch to another page.
// Using "page" for the page name
type KeysCommand struct {
	Parameters map[string]interface{}
}

// Execute the command
func (p *KeysCommand) Execute() (bool, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		clog.Logger.Errorf("error: %v", err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	value, found := p.Parameters["keys"]
	if found {
		keyValue, ok := value.(string)
		if ok {
			clog.Logger.Infof("Key pressed: %s", keyValue)
			kb.SetKeys(keybd_event.VK_A, keybd_event.VK_B)
		} else {
			keys := make([]string, 0)
			for _, keyValue := range value.([]interface{}) {
				keys = append(keys, keyValue.(string))
			}
			for _, keyValue := range keys {
				clog.Logger.Infof("Key pressed: %s", keyValue)
			}
		}
		// Press the selected keys
		err = kb.Launching()
		if err != nil {
			panic(err)
		}
	} else {
		return false, fmt.Errorf("The command parameter is missing")
	}
	return true, nil
}
