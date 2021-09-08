package models

type ExchangeType string

const (
	ExchangePage   ExchangeType = "pages"
	ExchangeAction              = "actions"
)

// Profile is the container for different pages. In UI you can switch between Profiles. Every Profile consist of a name and different pages to navigate between
type ProfileExchange struct {
	// Name of this profile
	Name string `json:"name"`
	// Description of this action for information
	Description string `json:"description"`
	// Type of the exchange format
	Type ExchangeType `json:"type"`
	// Pages are the UI structure for the different pages
	Pages []*Page `json:"pages"`
	// Actions contains the action definitions
	Actions []*Action `json:"actions"`
}

// HasPage checking if the profile has a page with that name
func (p *ProfileExchange) HasPage(pageName string) bool {
	for _, page := range p.Pages {
		if page.Name == pageName {
			return true
		}
	}
	return false
}

// HasAction checking if the profile has an action with that name
func (p *ProfileExchange) HasAction(actionName string) bool {
	for _, action := range p.Actions {
		if action.Name == actionName {
			return true
		}
	}
	return false
}
