package dto

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// Profiles contains all profiles defined for this server
var Profiles []Profile

// CommandExecutor is an interface for executing a command. Every command implementation has to implement this.
type CommandExecutor interface {
	Init(a *Action) (bool, error)
	Execute(a *Action, requestMessage models.Message) (bool, error)
	Stop(a *Action) (bool, error)
}

// Profile holding state informations about one profile
type Profile struct {
	Name    string
	Config  models.Profile
	Actions []Action
}

// Action holding status of one action and can execute this action
type Action struct {
	Profile  string
	RunOne   bool
	Name     string
	Title    string
	Config   models.Action
	m        sync.Mutex
	counter  int
	Commands map[string]CommandExecutor
}

// InitProfiles initialse the dto profiles for saving/retrieving status of and executing every action
func InitProfiles(configProfiles []models.Profile) error {
	count := 0
	Profiles = make([]Profile, 0)
	for _, configProfile := range configProfiles {
		dtoProfile := Profile{
			Name:    configProfile.Name,
			Config:  configProfile,
			Actions: make([]Action, 0),
		}

		for _, configAction := range configProfile.Actions {
			count++
			title := fmt.Sprintf("%s_%d", configAction.Name, count)
			action := Action{
				Profile:  configProfile.Name,
				RunOne:   configAction.RunOne,
				Name:     configAction.Name,
				Config:   configAction,
				Title:    title,
				counter:  0,
				Commands: make(map[string]CommandExecutor),
			}
			for _, command := range configAction.Commands {
				commandExecutor := GetCommand(command)
				if commandExecutor != nil {
					commandExecutor.Init(&action)
					action.Commands[command.Name] = commandExecutor
				}
			}
			dtoProfile.Actions = append(dtoProfile.Actions, action)
		}

		Profiles = append(Profiles, dtoProfile)
	}
	return nil
}

// Execute an action from a profile
func Execute(profileName string, actionName string, message models.Message) (bool, error) {
	profile, err := GetProfile(profileName)
	if err != nil {
		return false, err
	}
	action, err := profile.GetAction(actionName)
	if err != nil {
		return false, err
	}
	go doExecute(action, message)
	return true, nil
}

func doExecute(action *Action, message models.Message) {
	_, err := action.Execute(message)
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
func (p *Profile) GetAction(actionName string) (*Action, error) {
	for index := range p.Actions {
		action := &p.Actions[index]
		if strings.EqualFold(action.Name, actionName) {
			return action, nil
		}
	}
	return &Action{}, fmt.Errorf("Action %s not found", actionName)
}

// Execute an action
func (a *Action) Execute(requestMessage models.Message) (bool, error) {
	a.m.Lock()
	a.counter++
	counter := a.counter
	a.m.Unlock()
	if a.RunOne {
		a.m.Lock()
		defer a.m.Unlock()
	}
	clog.Logger.Debugf("execution action %s_%d", a.Title, counter)
	switch a.Config.Type {
	case models.Single:
		lastTitle := ""
		sendPostMessage := true
		for index, command := range a.Config.Commands {
			imageName := fmt.Sprintf("hourglass%d.png", index%4)
			if command.Icon != "" {
				imageName = command.Icon
			}
			title := ""
			if command.Title != "" {
				title = command.Title
				lastTitle = title
			}
			message := models.Message{
				Profile:  a.Profile,
				Action:   a.Name,
				ImageURL: imageName,
				Title:    title,
				State:    index + 1,
			}
			api.SendMessage(message)
			//cmdExecutor := GetCommand(command)
			cmdExecutor := a.Commands[command.Name]
			if cmdExecutor == nil {
				clog.Logger.Errorf("can't find command with type: %s", command.Type)
			}
			ok, err := cmdExecutor.Execute(a, requestMessage)
			if err != nil {
				clog.Logger.Errorf("error executing command: %v", err)
			}
			clog.Logger.Debugf("executing command result: %v", ok)
			sendPostMessage = sendPostMessage && ok
		}
		if sendPostMessage {

			message := models.Message{
				Profile:  a.Profile,
				Action:   a.Name,
				ImageURL: "check_mark.png",
				Title:    lastTitle,
				State:    0,
			}
			api.SendMessage(message)
			go func() {
				time.Sleep(3 * time.Second)
				message := models.Message{
					Profile:  a.Profile,
					Action:   a.Name,
					ImageURL: "",
					Title:    "",
					State:    0,
				}
				api.SendMessage(message)
			}()
		}
	}
	return true, nil
}

func IsSingleClick(message models.Message) bool {
	if message.Command == "" {
		return false
	}
	return message.Command == "click"
}

func IsDblClick(message models.Message) bool {
	if message.Command == "" {
		return false
	}
	return message.Command == "dblclick"
}
