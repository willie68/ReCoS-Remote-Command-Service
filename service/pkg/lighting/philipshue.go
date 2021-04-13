package lighting

import (
	"fmt"
	"sync"

	"github.com/amimof/huego"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

type PhilipsHue struct {
	username  string
	device    string
	ipaddress string
	bridge    huego.Bridge
	Lights    []huego.Light
	Groups    []huego.Group
	Scenes    []huego.Scene
	reload    sync.Mutex
}

var philipsHue *PhilipsHue

func GetPhilipsHue() (*PhilipsHue, bool) {
	if philipsHue != nil {
		return philipsHue, true
	}
	return nil, false
}

func InitPhilipsHue(extconfig map[string]interface{}) error {
	value, ok := extconfig["philipshue"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Debug("lighting:philipshue: found config")
			username, ok := config["username"].(string)
			if !ok {
				return fmt.Errorf("philipshue: no username given")
			}
			device, ok := config["device"].(string)
			if !ok {
				return fmt.Errorf("philipshue: no device given")
			}
			ipaddress, ok := config["ipaddress"].(string)
			if !ok {
				return fmt.Errorf("philipshue: no ipaddress given")
			}
			philipsHue = &PhilipsHue{
				username:  username,
				device:    device,
				ipaddress: ipaddress,
			}
			err := philipsHue.init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PhilipsHue) init() error {
	p.bridge = *huego.New(p.ipaddress, p.username)
	return p.reloadAll()
}

func (p *PhilipsHue) reloadAll() error {
	var err error
	p.Lights, err = p.getLights()
	if err != nil {
		clog.Logger.Errorf("error evaluating lights: %v", err)
		return err
	}

	p.Groups, err = p.getGroups()
	if err != nil {
		clog.Logger.Errorf("error evaluating groups: %v", err)
		return err
	}

	p.Scenes, err = p.getScenes()
	if err != nil {
		clog.Logger.Errorf("error evaluating scenes: %v", err)
		return err
	}
	return nil
}

func (p *PhilipsHue) getLights() ([]huego.Light, error) {
	l, err := p.bridge.GetLights()
	if err != nil {
		clog.Logger.Errorf("error evaluating lights: %v", err)
		return nil, err
	}
	clog.Logger.Infof("Found %d lights", len(l))
	for x, light := range l {
		clog.Logger.Debugf("light: %d = %s", x, light.Name)
	}
	return l, nil
}

func (p *PhilipsHue) getGroups() ([]huego.Group, error) {
	g, err := p.bridge.GetGroups()
	if err != nil {
		clog.Logger.Errorf("error evaluating groups: %v", err)
		return nil, err
	}
	for x, group := range g {
		clog.Logger.Debugf("group: %d = %s", x, group.Name)
	}
	return g, nil
}

func (p *PhilipsHue) getScenes() ([]huego.Scene, error) {
	scenes, err := p.bridge.GetScenes()
	if err != nil {
		clog.Logger.Errorf("error evaluating scenes: %v", err)
		return nil, err
	}
	for x, scene := range scenes {
		clog.Logger.Debugf("scene: %d = %s", x, scene.Name)
	}
	return scenes, nil
}

func (p *PhilipsHue) ToggleLight(lightname string) error {
	for _, light := range p.Lights {
		if light.Name == lightname {
			if light.IsOn() {
				light.Off()
			} else {
				light.On()
			}
		}
	}
	return nil
}

func (p *PhilipsHue) LightIsOn(lightname string) (bool, error) {
	light, ok := p.getLight(lightname)
	if !ok {
		p.reload.Lock()
		defer p.reload.Unlock()
		lights, err := p.getLights()
		if err != nil {
			return false, err
		}
		p.Lights = lights
		light, ok = p.getLight(lightname)
		if !ok {
			return false, fmt.Errorf("light with name %s not found", lightname)
		}
	}
	return light.IsOn(), nil
}

func (p *PhilipsHue) getLight(lightname string) (*huego.Light, bool) {
	for _, light := range p.Lights {
		if light.Name == lightname {
			return &light, true
		}
	}
	return nil, false
}
