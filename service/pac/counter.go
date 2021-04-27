package pac

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strconv"

	"github.com/fogleman/gg"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pac/clocks"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/session"
)

// CounterCommandTypeInfo counting clicks
var CounterCommandTypeInfo = models.CommandTypeInfo{
	Category:         "useful",
	Type:             "COUNTER",
	Name:             "Counter",
	Description:      "Counting button clicks",
	Icon:             "slot_machine.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "persist",
			Type:           "bool",
			Description:    "persist the value between restarts",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "oldschool",
			Type:           "bool",
			Description:    "Showing the value as seven-segment display",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
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

// CounterCommand is a command to switch to another page.
// Using "page" for the page name
type CounterCommand struct {
	Parameters  map[string]interface{}
	a           *Action
	commandName string
	countValue  float64
	persist     bool
	oldschool   bool
	color       color.Color
}

var (
	colorSegments color.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
)

// EnrichType enrich the type info with the informations from the profile
func (c *CounterCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return CounterCommandTypeInfo, nil
}

// Init the command
func (c *CounterCommand) Init(a *Action, commandName string) (bool, error) {
	c.a = a
	c.commandName = commandName
	c.persist = false
	value, found := c.Parameters["persist"]
	if found {
		var ok bool
		c.persist, ok = value.(bool)
		if !ok {
			return false, fmt.Errorf("persist is in wrong format. Please use boolean as format")
		}
	}

	if c.persist {
		value, ok := session.SessionCache.RetrieveCommandData(a.Profile, a.Name, c.commandName)
		if ok {
			c.countValue = value.(float64)
		}
	}

	value, found = c.Parameters["oldschool"]
	if found {
		var ok bool
		c.oldschool, ok = value.(bool)
		if !ok {
			return false, fmt.Errorf("oldschool is in wrong format. Please use boolean as format")
		}
	}

	value, found = c.Parameters["color"]
	if !found {
		c.color = colorSegments
	} else {
		myColor, err := parseHexColor(value.(string))
		if err != nil {
			clog.Logger.Errorf("error in getting sensors: %v", err)
			return false, err
		}
		c.color = myColor
	}

	return true, nil
}

// Stop the command
func (c *CounterCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (c *CounterCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if c.persist {
		value, ok := session.SessionCache.RetrieveCommandData(a.Profile, a.Name, c.commandName)
		if ok {
			c.countValue = value.(float64)
		}
	}

	if IsDblClick(requestMessage) {
		c.countValue = 0
	}
	if IsSingleClick(requestMessage) {

		c.countValue += 1
	}
	if c.persist {
		session.SessionCache.StoreCommandData(a.Profile, a.Name, c.commandName, c.countValue)
	}
	c.UpdateClients(a, c.commandName)
	return false, nil
}

func (c *CounterCommand) UpdateClients(a *Action, commandName string) (bool, error) {
	if !c.oldschool {
		text := fmt.Sprintf("%d", int(c.countValue))
		message := models.Message{
			Profile: a.Profile,
			Action:  a.Name,
			Text:    text,
			State:   0,
		}
		api.SendMessage(message)
	} else {
		c.SendGraphics()
	}
	return true, nil
}

// SendGraphics sending a new image url to the client
func (c *CounterCommand) SendGraphics() {
	id := strconv.Itoa(int(c.countValue))
	image := GetImageURL(c.a, c.commandName, id)
	message := models.Message{
		Profile:  c.a.Profile,
		Action:   c.a.Name,
		ImageURL: image,
		State:    0,
	}
	api.SendMessage(message)
}

// GetGraphics creates a clock graphics from the id
func (c *CounterCommand) GetGraphics(id string, width int, height int) (models.GraphicsInfo, error) {
	if width <= 0 {
		width = imageWidth
	}
	if height <= 0 {
		height = imageHeight
	}
	value, err := strconv.Atoi(id)
	if err != nil {
		clog.Logger.Debugf("counter: get graphics value not correct: %s", id)
		return models.GraphicsInfo{}, err
	}
	clog.Logger.Debugf("counter: get graphics value: %d", value)
	buff := c.generatePNG(value, width, height)
	model := models.GraphicsInfo{
		Mimetype: "image/png",
		Data:     buff,
	}
	return model, nil
}

func (c *CounterCommand) generatePNG(value int, width int, height int) []byte {
	dc := gg.NewContext(width, height)
	dc.SetColor(c.color)
	dc.InvertY()
	digits := float64(0)
	myValue := value
	for ; myValue > 0; digits++ {
		myValue = myValue / 10
	}
	segmentThickness := float64(height / 15)
	yHeightDigit := float64(height * 2 / 3)
	xWidthDigit := yHeightDigit / 2
	yStartDigit := float64(height / 6)
	xStartDigit := (float64(width) - (xWidthDigit * digits)) / 2

	myValue = value
	x := digits - 1
	for myValue > 0 {
		digit := myValue % 10
		clocks.WriteSegment(digit, xStartDigit+(xWidthDigit*x), yStartDigit, xWidthDigit, yHeightDigit, segmentThickness, dc, c.color, false)
		myValue = myValue / 10
		x--
	}

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)

	return buff.Bytes()
}
