package lighting

import (
	"fmt"
	"sync"

	"github.com/amimof/huego"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

var (
	username  string
	device    string
	ipaddress string
	onceHue   sync.Once
)

func InitPhilipsHue(extconfig map[string]interface{}) error {
	value, ok := extconfig["philipshue"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Debug("lighting:philipshue: found config")
			username, ok = config["username"].(string)
			if !ok {
				return fmt.Errorf("philipshue: no username given")
			}
			device, ok = config["device"].(string)
			if !ok {
				return fmt.Errorf("philipshue: no device given")
			}
			ipaddress, ok = config["ipaddress"].(string)
			if !ok {
				return fmt.Errorf("philipshue: no ipaddress given")
			}
		}
	}
	onceHue.Do(func() {
		bridge := huego.New(ipaddress, username)
		l, err := bridge.GetLights()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Found %d lights", len(l))
	})
	return nil
}
