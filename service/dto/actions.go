package dto

import (
	"fmt"
	"strings"
	"time"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// Profiles contains all profiles defined for this server
var Profiles []Profile

// Profile holding state informations about one profile
type Profile struct {
	Name    string
	Config  models.Profile
	Actions []Action
}

// Action holding status of one action and can execute this action
type Action struct {
	Name   string
	Config models.Action
}

// InitProfiles initialse the dto profiles for saving/retrieving status of and executing every action
func InitProfiles(configProfiles []models.Profile) error {
	Profiles = make([]Profile, 0)
	for _, configProfile := range configProfiles {
		dtoProfile := Profile{
			Name:    configProfile.Name,
			Config:  configProfile,
			Actions: make([]Action, 0),
		}
		Profiles = append(Profiles, dtoProfile)

		for _, configAction := range configProfile.Actions {
			action := Action{
				Name:   configAction.Name,
				Config: configAction,
			}
			dtoProfile.Actions = append(dtoProfile.Actions, action)
		}
	}
	return nil
}

// Execute an action from a profile
func Execute(profileName string, action string) (bool, error) {
	found := false
	for _, profile := range Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			found = true
		}
	}
	if !found {
		return false, fmt.Errorf("no profile with the name %s found", profileName)
	}
	return true, nil
}

// Execute an action
func (a *Action) Execute() {
	switch a.Config.Type {
	case models.Single:
		for _, command := range a.Config.Commands {
			switch command.Type {
			case models.Delay:
				{
					value, found := command.Parameters["delayTime"]
					if found {
						delayValue, ok := value.(int)
						if ok {
							clog.Logger.Infof("delay with %v seconds", delayValue)
							time.Sleep(time.Duration(delayValue) * time.Second)
						}
					}
				}
			}
		}
	}
}
