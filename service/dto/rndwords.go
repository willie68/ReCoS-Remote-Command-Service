package dto

import (
	"fmt"
	"math/rand"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// RndWordsCommandTypeInfo start an browser with directly with a url or filepath
var RndWordsCommandTypeInfo = models.CommandTypeInfo{
	Type:             "RNDWORDS",
	Name:             "RandomWords",
	Description:      "randomly select word",
	Icon:             "games_dice.png.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "words",
			Type:           "[]string",
			Description:    "the words to select one from",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	}}

// RndWordsCommand is a command to check the ping times.
// Using "sides" for the number of sides of the dice.
type RndWordsCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	words       []string
}

// EnrichType enrich the type info with the informations from the profile
func (r *RndWordsCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return RndWordsCommandTypeInfo, nil
}

// Init the command
func (r *RndWordsCommand) Init(a *Action, commandName string) (bool, error) {
	r.action = a
	r.commandName = commandName

	values, err := ConvertParameter2StringArray(r.Parameters, "words")
	if err != nil {
		return false, fmt.Errorf("The sides parameter is in wrong format. Please use int as format")
	}
	r.words = values

	rand.Seed(time.Now().UnixNano())
	return true, nil
}

// Stop the command
func (r *RndWordsCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (r *RndWordsCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if len(r.words) > 0 {
		index := rand.Intn(len(r.words))
		text := fmt.Sprintf("%d: %s", index+1, r.words[index])
		message := models.Message{
			Profile: r.action.Profile,
			Action:  r.action.Name,
			Text:    text,
			State:   0,
		}
		api.SendMessage(message)
	}
	return false, nil
}
