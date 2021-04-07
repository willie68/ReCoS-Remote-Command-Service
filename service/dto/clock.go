package dto

import (
	"image/color"
	"strconv"
	"time"

	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/dto/clocks"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ClockCommandTypeInfo is a clock
var ClockCommandTypeInfo = models.CommandTypeInfo{
	Type:             "CLOCK",
	Name:             "Clock",
	Description:      "Displaying a nice clock",
	Icon:             "clock.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "format",
			Type:           "string",
			Description:    "Format string for formatting the clock",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "timezone",
			Type:           "string",
			Description:    "time zone of the clock",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "dateformat",
			Type:           "string",
			Description:    "Format string for formatting the date",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "analog",
			Type:           "bool",
			Description:    "Showing a nice analog clock",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "showseconds",
			Type:           "bool",
			Description:    "Showing seconds on a analog clock",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "showdate",
			Type:           "bool",
			Description:    "Showing the date on a analog clock",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "design",
			Type:           "string",
			Description:    "design pattern for the clock",
			Unit:           "",
			WizardPossible: true,
			List:           []string{"analog", "digital", "berlin", "roman"},
		},
		{
			Name:           "color",
			Type:           "color",
			Description:    "color of the display",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// ClockCommand is a command to execute a delay. Using time for getting the ttime in seconds to delay the execution.
type ClockCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	stop        bool
	ticker      *time.Ticker
	done        chan bool
	format      string
	dateformat  string
	analog      bool
	showseconds bool
	showdate    bool
	commandName string
	design      string
	color       color.Color
	timezone    string
}

// EnrichType enrich the type info with the informations from the profile
func (c *ClockCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	index := -1
	for x, parameter := range ClockCommandTypeInfo.Parameters {
		if parameter.Name == "timezone" {
			index = x
		}
	}
	ClockCommandTypeInfo.Parameters[index].List = make([]string, 0)
	ClockCommandTypeInfo.Parameters[index].List = append(ClockCommandTypeInfo.Parameters[index].List, GetIANANames()...)

	return ClockCommandTypeInfo, nil
}

// Init a delay in the actual context
func (c *ClockCommand) Init(a *Action, commandName string) (bool, error) {
	c.action = a
	c.stop = false
	c.ticker = time.NewTicker(1 * time.Second)
	c.format = "15:04:05"
	c.analog = false
	c.commandName = commandName
	c.done = make(chan bool)

	//	GetIANANames()

	value, err := ConvertParameter2Bool(c.Parameters, "analog", false)
	if err != nil {
		clog.Logger.Errorf("error in getting analog: %v", err)
		return false, err
	}
	c.analog = value

	value, err = ConvertParameter2Bool(c.Parameters, "showseconds", false)
	if err != nil {
		clog.Logger.Errorf("error in getting showseconds: %v", err)
		return false, err
	}
	c.showseconds = value

	value, err = ConvertParameter2Bool(c.Parameters, "showdate", false)
	if err != nil {
		clog.Logger.Errorf("error in getting showdate: %v", err)
		return false, err
	}
	c.showdate = value

	svalue, err := ConvertParameter2String(c.Parameters, "timezone", "")
	if err != nil {
		clog.Logger.Errorf("error in getting timezone: %v", err)
		return false, err
	}
	c.timezone = svalue

	svalue, err = ConvertParameter2String(c.Parameters, "format", "15:04:05")
	if err != nil {
		clog.Logger.Errorf("error in getting format: %v", err)
		return false, err
	}
	c.format = svalue

	svalue, err = ConvertParameter2String(c.Parameters, "dateformat", "02.01")
	if err != nil {
		clog.Logger.Errorf("error in getting format: %v", err)
		return false, err
	}
	c.dateformat = svalue

	c.design = "analog"
	svalue, err = ConvertParameter2String(c.Parameters, "design", "analog")
	if err != nil {
		clog.Logger.Errorf("error in getting format: %v", err)
		return false, err
	}
	c.design = svalue

	cvalue, err := ConvertParameter2Color(c.Parameters, "color", colorSegments)
	if err != nil {
		clog.Logger.Errorf("error in getting color: %v", err)
		return false, err
	}
	c.color = cvalue

	go func() {
		for {
			select {
			case <-c.done:
				return
			case t := <-c.ticker.C:
				if api.HasConnectionWithProfile(a.Profile) {
					title := t.Format(c.format)
					text := ""
					if c.timezone != "" {
						text = c.timezone
					}
					if c.analog {
						c.SendGraphics(title, text)
					} else {
						message := models.Message{
							Profile: a.Profile,
							Action:  a.Name,
							State:   1,
							Title:   title,
							Text:    text,
						}
						api.SendMessage(message)
					}
				}
			}
		}
	}()
	return true, nil
}

// Stop stops the actual command
func (c *ClockCommand) Stop(a *Action) (bool, error) {
	c.ticker.Stop()
	c.done <- true
	return true, nil
}

// Execute a delay in the actual context
func (c *ClockCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	return true, nil
}

// GetGraphics creates a clock graphics from the id
func (c *ClockCommand) GetGraphics(id string, width int, height int) (models.GraphicsInfo, error) {
	timeToRender := idToTime(id)
	if width <= 0 {
		width = clocks.ClockImageWidth
	}
	if height <= 0 {
		height = clocks.ClockImageHeight
	}
	if c.timezone != "" {
		location, err := time.LoadLocation(c.timezone)
		if err != nil {
			clog.Logger.Errorf("can't load location of %s", c.timezone)
		} else {
			timeToRender = timeToRender.In(location)
		}
	}

	var model models.GraphicsInfo
	switch c.design {
	case "analog":
		buff := clocks.GenerateAnalog(timeToRender, width, height, c.showseconds, c.showdate, c.action.Config.Fontsize, c.dateformat)
		model = models.GraphicsInfo{
			Mimetype: "image/bmp",
			Data:     buff,
		}
	case "berlin":
		buff := clocks.GenerateBerlin(timeToRender, width, height)
		model = models.GraphicsInfo{
			Mimetype: "image/png",
			Data:     buff,
		}
	case "digital":
		buff := clocks.GenerateDigital(timeToRender, width, height, c.color, c.showseconds)
		model = models.GraphicsInfo{
			Mimetype: "image/png",
			Data:     buff,
		}
	case "roman":
		buff := clocks.GenerateRoman(timeToRender, width, height)
		model = models.GraphicsInfo{
			Mimetype: "image/png",
			Data:     buff,
		}

	}
	return model, nil
}

// SendPNG sending this array to the client
func (c *ClockCommand) SendGraphics(value, text string) {
	now := time.Now()
	if c.timezone != "" {
		location, err := time.LoadLocation(c.timezone)
		if err != nil {
			clog.Logger.Errorf("can't load location of %s", c.timezone)
			return
		}
		now = now.In(location)
	}

	id := timeToID(now)
	image := GetImageURL(c.action, c.commandName, id)
	message := models.Message{
		Profile:  c.action.Profile,
		Action:   c.action.Name,
		ImageURL: image,
		Title:    "",
		Text:     text,
		State:    0,
	}
	api.SendMessage(message)
}

func timeToID(time time.Time) string {
	return strconv.FormatInt(time.Unix(), 10)
}

func idToTime(id string) time.Time {
	value, _ := strconv.ParseInt(id, 10, 64)
	return time.Unix(value, 0)
}
