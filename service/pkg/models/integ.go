package models

type IntegInfo struct {
	// Category of this plugin
	Category string `json:"category"`
	// Name is the command
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
	// Image is a simple image to show on the page
	Image string `json:"image"`
	// Parameters describes the needed parameters
	Parameters []ParamInfo `json:"parameter"`
}
