package models

// ProfileInfos information about all profiles
type ProfileInfos struct {
	// Profiles all profiles for this client
	Profiles []ProfileShortInfo `json:"profiles"`
}

// ProfileShortInfo is a short info about a profile
type ProfileShortInfo struct {
	// Name of this profile
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
}

// ProfileInfo is all needed info for the action client of a defined profile
type ProfileInfo struct {
	// Name of this profile
	Name string `json:"name"`
	// Description of this action for information
	Description string       `json:"description"`
	Pages       []PageInfo   `json:"pages"`
	Actions     []ActionInfo `json:"actions"`
}

// ActionInfo is all needed info for the action client of a defined profile action
type ActionInfo struct {
	Type        ActionType `json:"type"`
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

// PageInfo is all needed info for the action client of a defined profile page
type PageInfo struct {
	Name    string   `json:"name"`
	Columns int      `json:"columns"`
	Rows    int      `json:"rows"`
	Cells   []string `json:"cells"`
}
