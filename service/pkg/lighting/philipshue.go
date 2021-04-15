package lighting

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/amimof/huego"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
)

type PhilipsHue struct {
	username  string
	device    string
	ipaddress string
	periode   int
	bridge    huego.Bridge
	Lights    []huego.Light
	Groups    []huego.Group
	Scenes    []huego.Scene
	reload    sync.Mutex
	ticker    *time.Ticker
	done      chan bool
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
			updatePeriod, ok := config["updateperiod"].(int)
			if !ok {
				updatePeriod = 10
			}
			philipsHue = &PhilipsHue{
				username:  username,
				device:    device,
				ipaddress: ipaddress,
				periode:   updatePeriod,
				ticker:    time.NewTicker(time.Duration(updatePeriod) * time.Second),
				done:      make(chan bool),
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
	go func() {
		p.reloadAll()
	}()
	go func() {
		for {
			select {
			case <-p.done:
				return
			case <-p.ticker.C:
				p.reloadAll()
			}
		}
	}()

	return nil
}

func (p *PhilipsHue) reloadAll() error {
	var err error
	p.reload.Lock()
	p.Lights, err = p.getLights()
	p.reload.Unlock()
	if err != nil {
		clog.Logger.Errorf("error evaluating lights: %v", err)
		return err
	}

	p.reload.Lock()
	p.Groups, err = p.getGroups()
	p.reload.Unlock()
	if err != nil {
		clog.Logger.Errorf("error evaluating groups: %v", err)
		return err
	}

	p.reload.Lock()
	p.Scenes, err = p.getScenes()
	p.reload.Unlock()
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
	return l, nil
}

func (p *PhilipsHue) getGroups() ([]huego.Group, error) {
	g, err := p.bridge.GetGroups()
	if err != nil {
		clog.Logger.Errorf("error evaluating groups: %v", err)
		return nil, err
	}
	return g, nil
}

func (p *PhilipsHue) getScenes() ([]huego.Scene, error) {
	scenes, err := p.bridge.GetScenes()
	if err != nil {
		clog.Logger.Errorf("error evaluating scenes: %v", err)
		return nil, err
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

func (p *PhilipsHue) Light(lightname string) (*huego.Light, bool) {
	light, ok := p.getLight(lightname)
	if !ok {
		p.reload.Lock()
		defer p.reload.Unlock()
		lights, err := p.getLights()
		if err != nil {
			return nil, false
		}
		p.Lights = lights
		light, ok = p.getLight(lightname)
		if !ok {
			return nil, false
		}
	}
	return light, true
}

func (p *PhilipsHue) LightIsOn(lightname string) (bool, error) {
	light, ok := p.Light(lightname)
	if !ok {
		return false, fmt.Errorf("light with name %s not found", lightname)
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

func (p *PhilipsHue) Group(groupname string) (*huego.Group, bool) {
	group, ok := p.getGroup(groupname)
	if !ok {
		p.reload.Lock()
		defer p.reload.Unlock()
		groups, err := p.getGroups()
		if err != nil {
			return nil, false
		}
		p.Groups = groups
		group, ok = p.getGroup(groupname)
		if !ok {
			return nil, false
		}
	}
	return group, true
}

func (p *PhilipsHue) getGroup(groupname string) (*huego.Group, bool) {
	for _, group := range p.Groups {
		if group.Name == groupname {
			return &group, true
		}
	}
	return nil, false
}

func (p *PhilipsHue) GroupIsOn(lightname string) (bool, error) {
	group, ok := p.Group(lightname)
	if !ok {
		return false, fmt.Errorf("group with name %s not found", lightname)
	}
	return group.IsOn(), nil
}

func (p *PhilipsHue) Scene(scenename string) ([]*huego.Scene, bool) {
	scene, ok := p.getNamedScenes(scenename)
	if !ok {
		p.reload.Lock()
		defer p.reload.Unlock()
		scenes, err := p.getScenes()
		if err != nil {
			return nil, false
		}
		p.Scenes = scenes
		scene, ok = p.getNamedScenes(scenename)
		if !ok {
			return nil, false
		}
	}
	return scene, true
}

func (p *PhilipsHue) SceneForGroup(group huego.Group, scenename string) (*huego.Scene, bool) {
	scenes4Group, ok := p.getNamedScenes(scenename)
	if !ok {
		p.reload.Lock()
		defer p.reload.Unlock()
		scenes, err := p.getScenes()
		if err != nil {
			return nil, false
		}
		p.Scenes = scenes
		scenes4Group, ok = p.getNamedScenes(scenename)
		if !ok {
			return nil, false
		}
	}
	for _, scene := range scenes4Group {
		if scene.Group == strconv.Itoa(group.ID) {
			return scene, true
		}
	}
	return nil, false
}

func (p *PhilipsHue) getNamedScenes(scenename string) ([]*huego.Scene, bool) {
	found := false
	scenes := make([]*huego.Scene, 0)
	for x, scene := range p.Scenes {
		if scene.Name == scenename {
			scenes = append(scenes, &p.Scenes[x])
			found = true
		}
	}
	return scenes, found
}
