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

type ActionExecutor interface {
	Execute() (bool, error)
}

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

		for _, configAction := range configProfile.Actions {
			action := Action{
				Name:   configAction.Name,
				Config: configAction,
			}
			dtoProfile.Actions = append(dtoProfile.Actions, action)
		}

		Profiles = append(Profiles, dtoProfile)
	}
	return nil
}

// Execute an action from a profile
func Execute(profileName string, actionName string) (bool, error) {
	profile, err := GetProfile(profileName)
	if err != nil {
		return false, err
	}
	action, err := profile.GetAction(actionName)
	if err != nil {
		return false, err
	}
	go doExecute(action)
	return true, nil
}

func doExecute(action Action) {
	_, err := action.Execute()
	if err != nil {
		clog.Logger.Errorf("Error executing action: %v", err)
	}
}

// GetProfile return the action with the name actionName if present otherwise an error
func GetProfile(profileName string) (Profile, error) {
	for _, profile := range Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			return profile, nil
		}
	}
	return Profile{}, fmt.Errorf("Profile %s not found", profileName)
}

// GetAction return the action with the name actionName if present otherwise an error
func (p *Profile) GetAction(actionName string) (Action, error) {
	for _, action := range p.Actions {
		if strings.EqualFold(action.Name, actionName) {
			return action, nil
		}
	}
	return Action{}, fmt.Errorf("Action %s not found", actionName)
}

// Execute an action
func (a *Action) Execute() (bool, error) {
	switch a.Config.Type {
	case models.Single:
		for _, command := range a.Config.Commands {
			switch command.Type {
			case models.Delay:
				{
					value, found := command.Parameters["time"]
					if found {
						delayValue, ok := value.(int)
						if ok {
							clog.Logger.Infof("delay with %v seconds", delayValue)
							time.Sleep(time.Duration(delayValue) * time.Second)
						} else {
							clog.Logger.Errorf("time is in wrong format")
						}
					} else {
						clog.Logger.Errorf("time is missing")
					}
				}
			}
		}
	}
	return true, nil
}
