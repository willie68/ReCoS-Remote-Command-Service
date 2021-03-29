package dto

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/bmp"
	"wkla.no-ip.biz/remote-desk-service/api"
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
			Name:           "design",
			Type:           "string",
			Description:    "design pattern for the clock",
			Unit:           "",
			WizardPossible: false,
			List:           []string{"analog", "digital"},
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
	analog      bool
	showseconds bool
	commandName string
	design      string
	color       color.Color
}

const clockImageWidth = 200
const clockImageHeight = 200

var (
	colorTicks    color.Color = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	colorArHour   color.Color = color.RGBA{R: 0x40, G: 0x40, B: 0x80, A: 0xFF}
	colorArMinute color.Color = color.RGBA{R: 0x40, G: 0x40, B: 0x80, A: 0xFF}
	colorArSecond color.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xFF}
	tickLength    float64     = 15
)

// EnrichType enrich the type info with the informations from the profile
func (c *ClockCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
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

	value, found := c.Parameters["analog"]
	if found {
		var ok bool
		c.analog, ok = value.(bool)
		if !ok {
			return false, fmt.Errorf("Analog is in wrong format. Please use boolean as format")
		}
	}

	c.showseconds = false
	value, found = c.Parameters["showseconds"]
	if found {
		var ok bool
		c.showseconds, ok = value.(bool)
		if !ok {
			return false, fmt.Errorf("Showseconds is in wrong format. Please use boolean as format")
		}
	}

	c.done = make(chan bool)
	value, found = c.Parameters["format"]
	if found {
		var ok bool
		c.format, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("Format is in wrong format. Please use string as format")
		}
	}

	c.design = "analog"
	value, found = c.Parameters["design"]
	if found {
		var ok bool
		c.design, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("Design is in wrong format. Please use string as format")
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

	go func() {
		for {
			select {
			case <-c.done:
				return
			case t := <-c.ticker.C:
				if api.HasConnectionWithProfile(a.Profile) {
					title := t.Format(c.format)
					if c.analog {
						c.SendGraphics(title)
					} else {
						message := models.Message{
							Profile: a.Profile,
							Action:  a.Name,
							State:   1,
							Title:   title,
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
		width = clockImageWidth
	}
	if height <= 0 {
		height = clockImageHeight
	}
	var model models.GraphicsInfo
	switch c.design {
	case "analog":
		buff := c.generateAnalog(timeToRender, width, height)
		model = models.GraphicsInfo{
			Mimetype: "image/bmp",
			Data:     buff,
		}
	case "digital":
		buff := c.generateDigital(timeToRender, width, height)
		model = models.GraphicsInfo{
			Mimetype: "image/png",
			Data:     buff,
		}

	}
	return model, nil
}

// SendPNG sending this array to the client
func (c *ClockCommand) SendGraphics(value string) {
	now := time.Now()

	//	buff := c.generateBMP(now)

	// Encode the bytes in the buffer to a base64 string
	//	encodedString := base64.StdEncoding.EncodeToString(buff)

	// You can embed it in an html doc with this string
	//	image := "data:image/bmp;base64," + encodedString
	id := timeToID(now)
	image := GetImageURL(c.action, c.commandName, id)
	message := models.Message{
		Profile:  c.action.Profile,
		Action:   c.action.Name,
		ImageURL: image,
		Title:    "",
		Text:     "",
		State:    0,
	}
	api.SendMessage(message)
}

// generateBMP generates a nice clock bmp
func (c *ClockCommand) generateAnalog(timeToRender time.Time, width int, height int) []byte {
	edgeSize := height
	if width < height {
		edgeSize = width
	}
	dc := gg.NewContext(edgeSize, edgeSize)
	halfWidth := float64(edgeSize / 2)
	halfHeight := float64(edgeSize / 2)
	floatEdgeSize := float64(edgeSize)
	myTicklength := tickLength * floatEdgeSize / float64(clockImageHeight)
	dc.SetColor(colorTicks)
	dc.InvertY()

	dc.SetLineWidth(1.0)
	dc.DrawCircle(halfWidth, halfHeight, halfHeight-1)
	dc.MoveTo(halfWidth-1, floatEdgeSize)
	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)

	for i := 0; i < 12; i++ {
		deg := float64(30.0 * i)
		dc.MoveTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-myTicklength)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-myTicklength)))
		dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-1)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-1)))
	}
	dc.Stroke()

	dc.SetLineWidth(4.0)

	dc.MoveTo(halfWidth-1, floatEdgeSize)
	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)
	dc.MoveTo(halfWidth-1, 0)
	dc.LineTo(halfWidth-1, myTicklength)

	dc.MoveTo(0, halfHeight-1)
	dc.LineTo(myTicklength, halfHeight-1)
	dc.MoveTo(floatEdgeSize-myTicklength, halfHeight-1)
	dc.LineTo(floatEdgeSize, halfHeight-1)

	dc.Stroke()

	if c.showseconds {
		dc.SetColor(colorArSecond)
		dc.SetLineWidth(1.0)
		seconds := timeToRender.Second()

		deg := float64(6.0 * seconds)
		dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*10), halfHeight-(math.Cos(deg2Rad(deg))*10))
		dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-2)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-2)))

		dc.Stroke()
	}

	dc.SetColor(colorArMinute)
	dc.SetLineWidth(3.0)
	minute := timeToRender.Minute()

	deg := float64(6.0 * minute)
	dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*2), halfHeight-(math.Cos(deg2Rad(deg))*2))
	dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-10)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-10)))

	dc.Stroke()

	dc.SetColor(colorArHour)
	dc.SetLineWidth(6.0)
	hour := timeToRender.Hour()

	deg = float64(30.0*hour + (minute / 2))
	dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*2), halfHeight-(math.Cos(deg2Rad(deg))*2))
	dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth*1/2)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight*1/2)))

	dc.Stroke()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	bmp.Encode(&buff, myImage)
	return buff.Bytes()
}

// generateBMP generates a nice clock bmp
func (c *ClockCommand) generateDigital(timeToRender time.Time, width int, height int) []byte {
	dc := gg.NewContext(width, height)
	dc.SetColor(c.color)
	dc.InvertY()
	/*	digits := float64(4)
		if c.showseconds {
			digits = float64(6)
		}
	*/
	xWidthDigit := float64(width) / 4.5
	if c.showseconds {
		xWidthDigit = float64(width) / 7.0
	}
	yHeightDigit := xWidthDigit * 1.5
	segmentThickness := xWidthDigit / 15.0
	yStartDigit := (float64(height) - yHeightDigit) / 2.0
	xStartDigit := 0.0

	hour := timeToRender.Hour()

	myValue := hour
	x := float64(1)
	for x >= 0 {
		xPos := xStartDigit + (xWidthDigit * x)
		yPos := yStartDigit
		digit := myValue % 10
		myValue = myValue / 10
		writeSegment(digit, xPos, yPos, xWidthDigit, yHeightDigit, segmentThickness, dc, c.color, true)
		x--
	}

	xPos := xStartDigit + (xWidthDigit * 2.25)
	yPos := yStartDigit + yHeightDigit/3.0
	clog.Logger.Debugf("Clock circle: %.0f %.0f %.0f", xPos, yPos, segmentThickness)
	dc.SetColor(c.color)
	if (timeToRender.Second()%2) == 1 && !c.showseconds {
		dc.SetColor(colorDarkSegment)
	}
	dc.DrawCircle(xPos, yPos, segmentThickness)
	yPos = yStartDigit + yHeightDigit*2.0/3.0
	dc.DrawCircle(xPos, yPos, segmentThickness)
	dc.Fill()
	dc.Stroke()

	xDelta := xWidthDigit / 2
	minutes := timeToRender.Minute()

	myValue = minutes
	x = float64(3)
	for x >= 2 {
		digit := myValue % 10
		writeSegment(digit, xDelta+xStartDigit+(xWidthDigit*x), yStartDigit, xWidthDigit, yHeightDigit, segmentThickness, dc, c.color, true)
		myValue = myValue / 10
		x--
	}

	if c.showseconds {
		xPos := xStartDigit + (xWidthDigit * 4.75)
		yPos := yStartDigit + yHeightDigit/3.0
		clog.Logger.Debugf("Clock circle: %.0f %.0f %.0f", xPos, yPos, segmentThickness)
		dc.SetColor(c.color)
		if (timeToRender.Second()%2) == 1 && !c.showseconds {
			dc.SetColor(colorDarkSegment)
		}
		dc.DrawCircle(xPos, yPos, segmentThickness)
		yPos = yStartDigit + yHeightDigit*2.0/3.0
		dc.DrawCircle(xPos, yPos, segmentThickness)
		dc.Fill()
		dc.Stroke()

		xDelta := xWidthDigit
		seconds := timeToRender.Second()

		myValue = seconds
		x = float64(5)
		for x >= 4 {
			digit := myValue % 10
			writeSegment(digit, xDelta+xStartDigit+(xWidthDigit*x), yStartDigit, xWidthDigit, yHeightDigit, segmentThickness, dc, c.color, true)
			myValue = myValue / 10
			x--
		}
	}

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)

	return buff.Bytes()
}

func timeToID(time time.Time) string {
	return strconv.FormatInt(time.Unix(), 10)
}

func idToTime(id string) time.Time {
	value, _ := strconv.ParseInt(id, 10, 64)
	return time.Unix(value, 0)
}
