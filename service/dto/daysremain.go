package dto

import (
	"fmt"
	"strings"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// DaysRemainCommandTypeInfo is a count down timer but just for days for a defined end date
var DaysRemainCommandTypeInfo = models.CommandTypeInfo{
	Type:             "DAYSREMAIN",
	Name:             "Days remain",
	Description:      "Displaying the days remain to a date",
	Icon:             "calendar_year.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "date",
			Type:           "date",
			Description:    "date to count to",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "formattitle",
			Type:           "string",
			Description:    "the title message for the response, defaults %d",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "formattext",
			Type:           "string",
			Description:    "the text message for the response, defaults days remain",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "finished",
			Type:           "string",
			Description:    "the message at the end of the Counter, defaults: finished",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// DaysRemainCommand is a command to start a days count down.
// there are several options, formattitle, formattext, finished and of course the date for the enddate
// Use %d for inserting the actual days to wait.
type DaysRemainCommand struct {
	Parameters  map[string]interface{}
	date        time.Time
	formatTitle string
	formatText  string
	finished    string
	ticker      *time.Ticker
	done        chan bool
}

// EnrichType enrich the type info with the informations from the profile
func (d *DaysRemainCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return DaysRemainCommandTypeInfo, nil
}

// Init a days remain in the actual context
func (d *DaysRemainCommand) Init(a *Action, commandName string) (bool, error) {
	value, found := d.Parameters["date"]
	if found {
		var ok bool
		dateValue, ok := value.(string)
		if !ok {
			return false, fmt.Errorf("Date is in wrong format. Please use string as format")
		}
		layout := "2006-01-02"
		ddate, err := time.Parse(layout, dateValue)
		if err != nil {
			layout = "2006-01-02T15:04:05.000Z"
			ddate, err = time.Parse(layout, dateValue)
			ddate = time.Date(ddate.Year(), ddate.Month(), ddate.Day(), 0, 0, 0, 0, ddate.Location())
			if err != nil {
				return false, fmt.Errorf("Date format is not correct. %v", err)
			}
		}
		d.date = ddate
	}

	format, err := ConvertParameter2String(d.Parameters, "formattitle", "%d")
	if err != nil {
		return false, err
	}
	d.formatTitle = format

	format, err = ConvertParameter2String(d.Parameters, "formattext", "days remain")
	if err != nil {
		return false, err
	}
	d.formatText = format

	format, err = ConvertParameter2String(d.Parameters, "finished", "finnished")
	if err != nil {
		return false, err
	}
	d.finished = format

	d.ticker = time.NewTicker(10 * time.Second)
	d.done = make(chan bool)
	go func() {
		for {
			select {
			case <-d.done:
				return
			case <-d.ticker.C:
				d.sendMessage(a)
			}
		}
	}()
	return true, nil
}

func (d *DaysRemainCommand) sendMessage(a *Action) {
	duration := d.date.Sub(time.Now())
	remainDays := int((duration.Hours() + 24) / 24.0)
	title := d.formatTitle
	if strings.Count(d.formatTitle, "%") > 0 {
		title = fmt.Sprintf(d.formatTitle, remainDays)
	}
	text := d.formatText
	if strings.Count(d.formatText, "%") > 0 {
		text = fmt.Sprintf(d.formatText, remainDays)
	}
	if remainDays <= 0 {
		title = d.finished
		layout := "2006-01-02"
		text = d.date.Format(layout)
	}
	message := models.Message{
		Profile: a.Profile,
		Action:  a.Name,
		State:   1,
		Title:   title,
		Text:    text,
	}
	api.SendMessage(message)
}

// Stop a timer in the actual context
func (d *DaysRemainCommand) Stop(a *Action) (bool, error) {
	d.ticker.Stop()
	d.done <- true
	return true, nil
}

// Execute a timer in the actual context
func (d *DaysRemainCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	d.sendMessage(a)
	return false, nil
}
