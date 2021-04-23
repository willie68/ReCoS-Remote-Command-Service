package pac

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// Profiles contains all profiles defined for this server
var Profiles []*Profile
var count int

// CommandExecutor is an interface for executing a command. Every command implementation has to implement this.
type CommandExecutor interface {
	EnrichType(profile models.Profile) (models.CommandTypeInfo, error)
	Init(a *Action, commandName string) (bool, error)
	Execute(a *Action, requestMessage models.Message) (bool, error)
	Stop(a *Action) (bool, error)
}

// CommandExecutor is an interface for executing a command. Every command implementation has to implement this.
type GraphicsCommandExecutor interface {
	GetGraphics(id string, width int, height int) (models.GraphicsInfo, error)
}

// ClientUpdatesCommandExecutor is an intrefce for the command executor, which methods will be called to update the clients with the alst state of the underlying command
type ClientUpdatesCommandExecutor interface {
	UpdateClients(a *Action, commandName string) (bool, error)
}

// Profile holding state informations about one profile
type Profile struct {
	Name    string
	Config  models.Profile
	Actions []*Action
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
	State    int
	Commands map[string]CommandExecutor
	Actions  []string
}

// InitProfiles initialse the dto profiles for saving/retrieving status of and executing every action
func InitProfiles(configProfiles []models.Profile) error {
	api.ClientUpdateCallback = func(profileName string) {
		UpdateClient4Profile(profileName)
	}
	Profiles = make([]*Profile, 0)
	for _, configProfile := range configProfiles {
		dtoProfile, err := InitProfile(configProfile.Name)
		if err == nil {
			Profiles = append(Profiles, dtoProfile)
		}
	}
	return nil
}

// InitProfile initialse the dto profiles for saving/retrieving status of and executing every action
func InitProfile(profileName string) (*Profile, error) {
	for _, configProfile := range config.Profiles {
		if profileName != configProfile.Name {
			continue
		}
		dtoProfile := Profile{
			Name:    configProfile.Name,
			Config:  configProfile,
			Actions: make([]*Action, 0),
		}

		for _, configAction := range configProfile.Actions {
			count++
			title := fmt.Sprintf("%s_%d", configAction.Name, count)
			action := Action{
				Profile:  configProfile.Name,
				RunOne:   configAction.RunOne,
				Name:     configAction.Name,
				Config:   *configAction,
				Title:    title,
				counter:  0,
				Commands: make(map[string]CommandExecutor),
				Actions:  configAction.Actions,
				State:    -1,
			}
			if len(action.Actions) > 0 {
				action.State = 0
			}
			for _, command := range configAction.Commands {
				commandExecutor := GetCommand(*command)
				if commandExecutor != nil {
					ok, err := commandExecutor.Init(&action, command.Name)
					if ok {
						action.Commands[command.ID] = commandExecutor
					}
					if err != nil {
						clog.Logger.Errorf("error initialising command %s#%s: %v", action.Name, command.Name, err)
					}
				}
			}
			dtoProfile.Actions = append(dtoProfile.Actions, &action)
		}
		return &dtoProfile, nil
	}
	return nil, errors.New("profile not found")
}

// ReinitProfile reinitialse the dto profiles
func ReinitProfile(profileName string) error {
	dtoProfile, err := InitProfile(profileName)

	if err == nil {
		Profiles = append(Profiles, dtoProfile)
	}

	return err
}

// ReinitProfiles reinitialse the dto profiles
func ReinitProfiles(configProfiles []models.Profile) error {
	for _, profile := range Profiles {
		CloseProfile(profile.Name)
	}
	InitProfiles(configProfiles)
	return nil
}

func UpdateClient4Profile(profileName string) {
	for _, profile := range Profiles {
		if profile.Name == profileName {
			profile.UpdateClients()
		}
	}
}

// CloseProfile closes a profiles
func CloseProfile(profileName string) error {
	actualProfile, err := GetProfile(profileName)
	if err == nil {
		for _, action := range actualProfile.Actions {
			actualAction, err := actualProfile.GetAction(action.Name)
			if err == nil {
				actualAction.Close()
			}
		}
	}
	return nil
}

// RemoveProfile removing a profile from the active profiles array
func RemoveProfile(profileName string) error {
	pos := -1
	for x, profile := range Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			pos = x
		}
	}

	if pos >= 0 {
		Profiles[pos] = Profiles[len(Profiles)-1]
		// We do not need to put s[i] at the end, as it will be discarded anyway
		Profiles = Profiles[:len(Profiles)-1]
		return nil
	}
	return fmt.Errorf("Profile %s not found", profileName)
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

func Graphics(profileName string, actionName string, commandName string, id string, width int, height int) (models.GraphicsInfo, error) {
	empty := models.GraphicsInfo{}
	profile, err := GetProfile(profileName)
	if err != nil {
		return empty, err
	}
	action, err := profile.GetAction(actionName)
	if err != nil {
		return empty, err
	}
	return action.GetGraphics(commandName, id, width, height)
}

// GetProfile return the action with the name actionName if present otherwise an error
func GetProfile(profileName string) (*Profile, error) {
	for _, profile := range Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			return profile, nil
		}
	}
	return nil, fmt.Errorf("Profile %s not found", profileName)
}

// GetAction return the action with the name actionName if present otherwise an error
func (p *Profile) GetAction(actionName string) (*Action, error) {
	for index := range p.Actions {
		action := p.Actions[index]
		if strings.EqualFold(action.Name, actionName) {
			return action, nil
		}
	}
	return nil, fmt.Errorf("Action %s not found", actionName)
}

// GetAction return the action with the name actionName if present otherwise an error
func (p *Profile) UpdateClients() {
	for index := range p.Actions {
		action := p.Actions[index]
		action.UpdateClients()
	}
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
		doWorkSingle(a, a, requestMessage)
	case models.Multi:
		doWorkMulti(a, requestMessage)
	}
	return true, nil
}

func doWorkMulti(a *Action, requestMessage models.Message) (bool, error) {
	profile, err := GetProfile(a.Profile)
	if err != nil {
		return false, err
	}
	workingAction, err := profile.GetAction(a.Actions[a.State])
	if err != nil {
		return false, err
	}
	doWorkSingle(workingAction, a, requestMessage)
	a.State++
	if a.State >= len(a.Actions) {
		a.State = 0
	}
	return true, nil
}

func doWorkSingle(a *Action, sendingAction *Action, requestMessage models.Message) {
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
			Profile:  sendingAction.Profile,
			Action:   sendingAction.Name,
			ImageURL: imageName,
			Title:    title,
			State:    index + 1,
		}
		api.SendMessage(message)
		//cmdExecutor := GetCommand(command)
		cmdExecutor := a.Commands[command.ID]
		if cmdExecutor == nil {
			clog.Logger.Errorf("can't find command with type: %s", command.Type)
			return
		}
		ok, err := cmdExecutor.Execute(sendingAction, requestMessage)
		if err != nil {
			clog.Logger.Errorf("error executing command: %v", err)
		}
		//clog.Logger.Debugf("executing command result: %v", ok)
		sendPostMessage = sendPostMessage && ok
	}
	if sendPostMessage {
		message := models.Message{
			Profile:  sendingAction.Profile,
			Action:   sendingAction.Name,
			ImageURL: "check_mark.png",
			Title:    lastTitle,
			State:    0,
		}
		api.SendMessage(message)
		go func() {
			time.Sleep(3 * time.Second)
			message := models.Message{
				Profile:  sendingAction.Profile,
				Action:   sendingAction.Name,
				ImageURL: a.Config.Icon,
				Title:    a.Config.Title,
				Text:     "",
				State:    sendingAction.State,
			}
			api.SendMessage(message)
		}()
	}
}

func (a *Action) GetGraphics(commandName string, id string, width int, height int) (models.GraphicsInfo, error) {
	empty := models.GraphicsInfo{}
	for _, command := range a.Config.Commands {
		if command.Name == commandName {
			cmdExecutor := a.Commands[command.ID]
			if cmdExecutor == nil {
				clog.Logger.Errorf("can't find command with type: %s", command.Type)
				return empty, fmt.Errorf("can't find command with type: %s", command.Type)
			}
			v, ok := interface{}(cmdExecutor).(GraphicsCommandExecutor)
			if !ok {
				clog.Logger.Errorf("command can't create graphics: %s", command.Type)
				return empty, fmt.Errorf("command can't create graphics: %s", command.Type)
			}
			graphicsInfo, err := v.GetGraphics(id, width, height)
			if err != nil {
				clog.Logger.Errorf("error executing command: %v", err)
				return empty, err
			}
			//clog.Logger.Debugf("executing command result: true", ok)
			return graphicsInfo, nil
		}
	}
	return empty, fmt.Errorf("can't find command with name: %s", commandName)
}

// UpdateClients updating all clients with the atual state
func (a *Action) UpdateClients() {
	for _, command := range a.Config.Commands {
		cmdExecutor := a.Commands[command.ID]
		if cmdExecutor == nil {
			continue
		}
		clientUpdates, ok := interface{}(cmdExecutor).(ClientUpdatesCommandExecutor)
		if !ok {
			continue
		}
		clientUpdates.UpdateClients(a, command.Name)
	}
}

// Close an action will close/stop all dedicated commands
func (a *Action) Close() error {
	if (a.Config.Type == models.Single) || (a.Config.Type == models.Display) {
		a.m.Lock()
		clog.Logger.Debugf("close action %s", a.Title)

		for _, command := range a.Config.Commands {
			//cmdExecutor := GetCommand(command)
			cmdExecutor := a.Commands[command.ID]
			if cmdExecutor == nil {
				clog.Logger.Errorf("can't find command: %s, (%s#%s)", command.ID, command.Type, command.Name)
				continue
			}
			_, err := cmdExecutor.Stop(a)
			if err != nil {
				clog.Logger.Errorf("error stopping command: %v", err)
			}
		}

		a.Commands = make(map[string]CommandExecutor)
		a.m.Unlock()
	}

	return nil
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
