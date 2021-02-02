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
	Actions []Action `json:"actions"`
}

// Page is the most visible part of this. Every Page is organised in Rows and Columns. And with this every cell is a place for holding an action
type Page struct {
	Name string `json:"name"`
	// Columns of this page
	Columns int `json:"columns"`
	// Rows of this page
	Rows int `json:"rows"`
	// Cells are the ui container for the actions, only the action id will be used here
	Cells []string `json:"cells"`
}

//ActionType the type of the action enum
type ActionType string

const (
	//Single is a simple single shot. Normally displayed as a button. Simply one action
	Single ActionType = "SINGLE"
	//Toggle has a state, true or false and on every change it does an individual action
	Toggle = "TOGGLE"
	//MultiStage is an action which has more than 2 states and step to the next stage on every execution. After the last stage it return to the first one
	MultiStage = "MULTISTAGE"
)

// Action type
type Action struct {
	// Type is the type of an action
	Type ActionType `json:"type"`
	// Name is the name/id for this action
	Name string `json:"name"`
	// Title of this action for display
	Title string `json:"title"`
	// Description of this action for information
	Description string `json:"description"`
	// Commands are the magic behind this
	Commands []Command `json:"commands"`
}

// CommandType the command type
type CommandType string

const (
	// Delay dalay further execution a defined time
	Delay CommandType = "DELAY"
	// Execute start an application or shell script and optionally waits for it's finishing
	Execute = "EXECUTE"
)

// Command type
type Command struct {
	// Type is the type of an command
	Type CommandType `json:"type"`
	// Name is the command
	Name string `json:"name"`
	// Parameters of this command
	Parameters map[string]interface{} `json:"parameters"`
}

// Copy make a deep copy of this profile
func (p *Profile) Copy() Profile {
	profile := Profile{
		Name:        p.Name,
		Description: p.Description,
	}
	profile.Actions = make([]Action, 0)
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
func (a *Action) Copy() Action {
	action := Action{
		Name:        a.Name,
		Description: a.Description,
		Title:       a.Title,
		Type:        a.Type,
	}
	action.Commands = make([]Command, 0)
	for _, command := range a.Commands {
		action.Commands = append(action.Commands, command.Copy())
	}
	return action
}

// Copy make a deep copy of this action
func (p *Page) Copy() Page {
	page := Page{
		Name:    p.Name,
		Columns: p.Columns,
		Rows:    p.Rows,
	}
	page.Cells = make([]string, 0)
	for _, cell := range p.Cells {
		page.Cells = append(page.Cells, cell)
	}
	return page
}

// Copy make a deep copy of this action
func (c *Command) Copy() Command {
	command := Command{
		Name: c.Name,
		Type: c.Type,
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	enc.Encode(c.Parameters)
	command.Parameters = make(map[string]interface{})
	dec.Decode(&command.Parameters)
	return command
}
