package pac

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pac/clocks"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// TimerCommandTypeInfo is a count down timer, just showing the count down time in the title
var TimerCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Time",
	Type:             "TIMER",
	Name:             "Timer",
	Description:      "Starting a count down timer",
	Icon:             "timer.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "time",
			Type:           "int",
			Description:    "time to delay in Seconds",
			Unit:           " Seconds",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "format",
			Type:           "string",
			Description:    "the message for the response, defaults %d seconds",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "finished",
			Type:           "string",
			Description:    "the message at the end of the timer, defaults: finished",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

// TimerCommand is a command to start a timer. Using time for getting the time in seconds.
// For formatting the response the parameters fomat and finished are responsible.
// Use %d for inserting the actual time to wait.
type TimerCommand struct {
	action      *Action
	commandName string
	format      string
	finished    string
	time        int
	Parameters  map[string]interface{}
}

var (
	timerColorTicks  color.Color = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	timerColorGreen  color.Color = color.RGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xFF}
	timerColorYellow color.Color = color.RGBA{R: 0xff, G: 0xff, B: 0x00, A: 0xFF}
	timerColorRed    color.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xFF}
	timerTickLength  float64     = 15
)

// EnrichType enrich the type info with the informations from the profile
func (t *TimerCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return TimerCommandTypeInfo, nil
}

// Init a timer in the actual context
func (t *TimerCommand) Init(a *Action, commandName string) (bool, error) {
	var err error
	t.action = a
	t.commandName = commandName
	t.format, err = ConvertParameter2String(t.Parameters, "format", "%d seconds")
	if err != nil {
		return false, err
	}
	t.finished, err = ConvertParameter2String(t.Parameters, "finished", "finished")
	if err != nil {
		return false, err
	}
	t.time, err = ConvertParameter2Int(t.Parameters, "time", 10)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Stop a timer in the actual context
func (t *TimerCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a timer in the actual context
func (t *TimerCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	delayValue := t.time
	clog.Logger.Infof("count down with %v seconds", delayValue)
	for ; delayValue > 0; delayValue-- {
		// TODO get this from the config
		title := fmt.Sprintf(t.format, delayValue)
		t.SendGraphics(delayValue, title)
		time.Sleep(1 * time.Second)
	}
	message := models.Message{
		Profile:  a.Profile,
		Action:   a.Name,
		State:    0,
		Title:    t.finished,
		ImageURL: "check_mark.png",
	}
	api.SendMessage(message)
	time.Sleep(3 * time.Second)
	return true, nil
}

// SendPNG sending this array to the client
func (t *TimerCommand) SendGraphics(value int, text string) {
	image := GetImageURL(t.action, t.commandName, strconv.Itoa(value))
	message := models.Message{
		Profile:  t.action.Profile,
		Action:   t.action.Name,
		ImageURL: image,
		State:    value,
		Title:    text,
	}
	api.SendMessage(message)
}

// GetGraphics creates a clock graphics from the id
func (t *TimerCommand) GetGraphics(id string, width int, height int) (models.GraphicsInfo, error) {
	var model models.GraphicsInfo
	value, err := strconv.Atoi(id)
	if err != nil {
		return model, err
	}
	if width <= 0 {
		width = clocks.ClockImageWidth
	}
	if height <= 0 {
		height = clocks.ClockImageHeight
	}
	buff := t.generateTimerWatch(value, width, height)
	model = models.GraphicsInfo{
		Mimetype: "image/png",
		Data:     buff,
	}
	return model, nil
}

// generateTimerWatch generates a nice clock bmp
func (t *TimerCommand) generateTimerWatch(value int, width int, height int) []byte {
	edgeSize := height
	if width < height {
		edgeSize = width
	}
	dc := gg.NewContext(edgeSize, edgeSize)
	halfWidth := float64(edgeSize / 2)
	halfHeight := float64(edgeSize / 2)
	floatEdgeSize := float64(edgeSize)
	myTicklength := timerTickLength * floatEdgeSize / float64(clocks.ClockImageHeight)

	dc.InvertY()

	// Draw the case around the dsísplay
	dc.SetColor(timerColorTicks)
	dc.SetLineWidth(1.0)
	dc.DrawCircle(halfWidth, halfHeight, halfHeight-1)
	dc.MoveTo(halfWidth-1, floatEdgeSize)
	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)
	dc.Stroke()

	// draw the cake
	dc.SetColor(timerColorGreen)
	if value <= (t.time / 4) {
		dc.SetColor(timerColorYellow)
	}
	if value <= 1 {
		dc.SetColor(timerColorRed)
	}

	np := math.Pi / 2.0
	pos := float64(value) / float64(t.time) * -2.0 * math.Pi
	if value < t.time {

		clog.Logger.Infof("pos: %f", pos)

		dc.MoveTo(halfWidth, halfHeight)
		dc.DrawArc(halfWidth, halfHeight, halfWidth-1, np, pos+np)
		dc.LineTo(halfWidth, halfHeight)

	} else {
		dc.DrawCircle(halfWidth, halfHeight, halfWidth-1)
	}
	dc.Fill()
	dc.Stroke()

	// Draw some ticks
	dc.SetColor(timerColorTicks)
	dc.SetLineWidth(4.0)
	for i := 0; i < t.time; i++ {
		deg := 360.0 / float64(t.time) * float64(i)
		dc.MoveTo(halfWidth+(math.Sin(Deg2Rad(deg))*(halfWidth-myTicklength)), halfHeight+(math.Cos(Deg2Rad(deg))*(halfHeight-myTicklength)))
		dc.LineTo(halfWidth+(math.Sin(Deg2Rad(deg))*(halfWidth-1)), halfHeight+(math.Cos(Deg2Rad(deg))*(halfHeight-1)))
	}
	dc.Stroke()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)
	return buff.Bytes()
}
