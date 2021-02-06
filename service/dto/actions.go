package dto

// Profiles contains all profiles defined for this server
var Profiles []Profile

// Profile holding informations about one profile
type Profile struct {
	Name    string
	Actions []Action
}

// Action holding status of one action
type Action struct {
}
