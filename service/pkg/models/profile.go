package models

import (
	"bytes"
	"encoding/gob"
)

// Profile is the container for different pages. In UI you can switch between Profiles. Every Profile consist of a name and different pages to navigate between
type Profile struct {
	// Name of this profile
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
	// Pages are the UI structure for the different pages
	Pages []Page `json:"pages"`
	// Actions contains the action definitions
	Actions []*Action `json:"actions"`
}

type ToolbarType string

const (
	ToolbarShow ToolbarType = "show"
	ToolbarHide             = "hide"
)

// Page is the most visible part of this. Every Page is organised in Rows and Columns. And with this every cell is a place for holding an action
type Page struct {
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
	// Title of this action for display
	Icon string `json:"icon"`
	// Columns of this page
	Columns int `json:"columns"`
	// Rows of this page
	Rows int `json:"rows"`
	// Rows of this page
	Toolbar ToolbarType `json:"toolbar"`
	// Cells are the ui container for the actions, only the action id will be used here
	Cells []string `json:"cells"`
}

//ActionType the type of the action enum
type ActionType string

const (
	//Single is a simple single shot. Normally displayed as a button. Simply one action
	Single ActionType = "SINGLE"
	//Display is only showing title and icon, you can't click on it
	Display = "DISPLAY"
	//Toggle has a state, true or false and on every change it does an individual action
	Toggle = "TOGGLE"
	//Multi is an action which has more than 2 states and step to the next stage on every execution. After the last stage it return to the first one
	Multi = "MULTI"
)

// Action type
type Action struct {
	// Type is the type of an action
	Type ActionType `json:"type"`
	// Name is the name/id for this action
	Name string `json:"name"`
	// Title of this action for display
	Title string `json:"title"`
	// Title of this action for display
	Icon string `json:"icon"`
	// Description of this action for information
	Description string `json:"description"`
	// Fontsize size of the title and text font
	Fontsize int `json:"fontsize"`
	// Fontcolor color of the title and text font
	Fontcolor string `json:"fontcolor"`
	// Outlined the font of the title and text font
	Outlined bool `json:"outlined"`
	// RunOne means only run one instance of this action at one time.
	// Scheduling more than one execution will lead into a sequentiell execution
	RunOne bool `json:"runone"`
	// Commands are the magic behind this
	Commands []*Command `json:"commands"`
	// Actions are the actions to execute behind a multi action button
	Actions []string `json:"actions"`
}

// CommandType the command type
type CommandType string

type CommandTypeInfo struct {
	// Type is the type of an command
	Type CommandType `json:"type"`
	// Name is the command
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
	// WizardPossible this command can be used in the wizard
	WizardPossible bool `json:"wizard"`
	// Parameters describes the needed parameters
	Parameters []CommandParameterInfo `json:"parameter"`
}

type GraphicsInfo struct {
	Mimetype string
	Data     []byte
}

type CommandParameterInfo struct {
	// Name is the command
	Name string `json:"name"`
	// Type is the type of an command
	Type string `json:"type"`
	// Description of this action for information
	Description string `json:"description"`
	// Name is the command
	Unit string `json:"unit"`
	// WizardPossible this command can be used in the wizard
	WizardPossible bool `json:"wizard"`
	// List is a enumeration of possible values
	List []string `json:"list"`
}

// Command type
type Command struct {
	// ID is for internal use only
	ID string
	// Type is the type of an command
	Type CommandType `json:"type"`
	// Name is the command
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
	// Icon is the icon to show when this command is executing
	Icon string `json:"icon"`
	// Title is the title to show when this command is executing
	Title string `json:"title"`
	// Parameters of this command
	Parameters map[string]interface{} `json:"parameters"`
}

// Copy make a deep copy of this profile
func (p *Profile) Copy() Profile {
	profile := Profile{
		Name:        p.Name,
		Description: p.Description,
	}
	profile.Actions = make([]*Action, 0)
	for _, action := range p.Actions {
		profile.Actions = append(profile.Actions, action.Copy())
	}
	profile.Pages = make([]Page, 0)
	for _, page := range p.Pages {
		profile.Pages = append(profile.Pages, page.Copy())
	}
	return profile
}

// Copy make a deep copy of this action
func (a *Action) Copy() *Action {
	action := Action{
		Name:        a.Name,
		Description: a.Description,
		Title:       a.Title,
		Type:        a.Type,
		RunOne:      a.RunOne,
		Icon:        a.Icon,
		Fontsize:    a.Fontsize,
		Fontcolor:   a.Fontcolor,
		Outlined:    a.Outlined,
	}
	action.Commands = make([]*Command, 0)
	for _, command := range a.Commands {
		action.Commands = append(action.Commands, command.Copy())
	}
	action.Actions = make([]string, 0)
	for _, actionName := range a.Actions {
		action.Actions = append(action.Actions, actionName)
	}
	return &action
}

// Copy make a deep copy of this action
func (p *Page) Copy() Page {
	page := Page{
		Name:        p.Name,
		Columns:     p.Columns,
		Description: p.Description,
		Rows:        p.Rows,
		Toolbar:     p.Toolbar,
	}
	page.Cells = make([]string, 0)
	for _, cell := range p.Cells {
		page.Cells = append(page.Cells, cell)
	}
	return page
}

// Copy make a deep copy of this action
func (c *Command) Copy() *Command {
	command := Command{
		Name:        c.Name,
		Description: c.Description,
		Type:        c.Type,
		Icon:        c.Icon,
		Title:       c.Title,
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	enc.Encode(c.Parameters)
	command.Parameters = make(map[string]interface{})
	dec.Decode(&command.Parameters)
	return &command
}
