package dto

import (
	"fmt"

	"github.com/skratchdot/open-golang/open"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// BrowseCommandTypeInfo start an browser with directly with a url or filepath
var BrowseCommandTypeInfo = models.CommandTypeInfo{
	Type:             "BROWSE",
	Name:             "Browse",
	Description:      "Execute the default browser with a url or a filepath",
	Icon:             "world.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "url",
			Type:           "string",
			Description:    "the url/filepath to show",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	}}

// BrowseCommand is a command to execute the internet browser.
// Using "webpage" for getting the webpage to show.
type BrowseCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (e *BrowseCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return BrowseCommandTypeInfo, nil
}

// Init the command
func (e *BrowseCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop the command
func (e *BrowseCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (e *BrowseCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	value, found := e.Parameters["url"]
	if found {
		url, ok := value.(string)
		if ok {
			err := open.Run(url)
			clog.Logger.Debugf("start browser with: %s", url)
			if err != nil {
				clog.Logger.Errorf("error: %v\r\n", err)
				return false, fmt.Errorf("Error executing the url. %v", err)
			}
		} else {
			return false, fmt.Errorf("The url parameter is in wrong format. Please use string as format")
		}
	} else {
		return false, fmt.Errorf("The url parameter is missing")
	}
	return true, nil
}
