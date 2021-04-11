package dto

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"math/rand"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/dto/clocks"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// DiceCommandTypeInfo start an browser with directly with a url or filepath
var DiceCommandTypeInfo = models.CommandTypeInfo{
	Category:         "Games",
	Type:             "DICE",
	Name:             "Dice",
	Description:      "rolling a Dice",
	Icon:             "games_dice.png.png",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "sides",
			Type:           "int",
			Description:    "number of sides of the dice",
			Unit:           " Sides",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	}}

// DiceCommand is a command to check the ping times.
// Using "sides" for the number of sides of the dice.
type DiceCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	commandName string
	sides       int
}

var (
	diceColorDot   color.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xFF}
	diceColorSides color.Color = color.RGBA{R: 0xff, G: 0xe6, B: 0x99, A: 0xFF}
)

// EnrichType enrich the type info with the informations from the profile
func (d *DiceCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return DiceCommandTypeInfo, nil
}

// Init the command
func (d *DiceCommand) Init(a *Action, commandName string) (bool, error) {
	d.action = a
	d.commandName = commandName

	valueInt, err := ConvertParameter2Int(d.Parameters, "sides", 6)
	if err != nil {
		return false, fmt.Errorf("The sides parameter is in wrong format. Please use int as format")
	}
	d.sides = valueInt

	rand.Seed(time.Now().UnixNano())
	return true, nil
}

// Stop the command
func (d *DiceCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (d *DiceCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	text := fmt.Sprintf("%d", rand.Intn(d.sides)+1)
	message := models.Message{
		Profile: d.action.Profile,
		Action:  d.action.Name,
		Text:    text,
		State:   0,
	}
	api.SendMessage(message)
	if d.sides < 10 {
		d.SendGraphics()
	}
	return false, nil
}

// GetGraphics creates a clock graphics from the id
func (d *DiceCommand) GetGraphics(id string, width int, height int) (models.GraphicsInfo, error) {
	number, err := strconv.Atoi(id)
	if err != nil {
		clog.Logger.Errorf("Error : %v", err)
		return models.GraphicsInfo{}, err
	}
	if width <= 0 {
		width = clocks.ClockImageWidth
	}
	if height <= 0 {
		height = clocks.ClockImageHeight
	}
	var model models.GraphicsInfo
	buff := d.generateDice(number, width, height)
	model = models.GraphicsInfo{
		Mimetype: "image/bmp",
		Data:     buff,
	}
	return model, nil
}

// SendPNG sending this array to the client
func (d *DiceCommand) SendGraphics() {
	id := strconv.Itoa(rand.Intn(d.sides) + 1)
	image := GetImageURL(d.action, d.commandName, id)
	message := models.Message{
		Profile:  d.action.Profile,
		Action:   d.action.Name,
		ImageURL: image,
		Title:    " ",
		Text:     id,
		State:    0,
	}
	api.SendMessage(message)
}

type num2Dice struct {
	a bool
	b bool
	c bool
	d bool
	e bool
	f bool
	g bool
	h bool
	i bool
}

var num2DiceArray = []num2Dice{
	{false, false, false, false, false, false, false, false, false}, //0
	{false, false, false, false, true, false, false, false, false},  //1
	{true, false, false, false, false, false, false, false, true},   //2
	{true, false, false, false, true, false, false, false, true},    //3
	{true, false, true, false, false, false, true, false, true},     //4
	{true, false, true, false, true, false, true, false, true},      //5
	{true, true, true, false, false, false, true, true, true},       //6
	{true, true, true, false, true, false, true, true, true},        //7
	{true, true, true, true, false, true, true, true, true},         //8
	{true, true, true, true, true, true, true, true, true},          //9
}

// generateBMP generates a nice clock bmp
func (d *DiceCommand) generateDice(number int, width int, height int) []byte {
	edgeSize := height
	if width < height {
		edgeSize = width
	}
	dc := gg.NewContext(edgeSize, edgeSize)
	halfWidth := float64(edgeSize / 2)
	halfHeight := float64(edgeSize / 2)
	floatEdgeSize := float64(edgeSize)

	padding := floatEdgeSize / 10.0

	floatEdgeSize = floatEdgeSize - (2 * padding)
	dotRadius := floatEdgeSize / 12.0

	dc.InvertY()
	dc.SetColor(color.Black)

	dc.SetLineWidth(1.0)
	dc.DrawRoundedRectangle(padding, padding, floatEdgeSize, floatEdgeSize, floatEdgeSize/10.0)
	dc.SetColor(diceColorSides)
	dc.Fill()
	dc.Stroke()

	dc.SetColor(diceColorDot)
	myNum2Dice := num2DiceArray[number]

	if myNum2Dice.a {
		dc.DrawCircle(3*padding, 3*padding, dotRadius)
	}
	if myNum2Dice.b {
		dc.DrawCircle(halfWidth, 3*padding, dotRadius)
	}
	if myNum2Dice.c {
		dc.DrawCircle(floatEdgeSize-padding, 3*padding, dotRadius)
	}
	if myNum2Dice.d {
		dc.DrawCircle(3*padding, halfHeight, dotRadius)
	}
	if myNum2Dice.e {
		dc.DrawCircle(halfWidth, halfHeight, dotRadius)
	}
	if myNum2Dice.f {
		dc.DrawCircle(floatEdgeSize-padding, halfHeight, dotRadius)
	}
	if myNum2Dice.g {
		dc.DrawCircle(3*padding, floatEdgeSize-padding, dotRadius)
	}
	if myNum2Dice.h {
		dc.DrawCircle(halfWidth, floatEdgeSize-padding, dotRadius)
	}
	if myNum2Dice.i {
		dc.DrawCircle(floatEdgeSize-padding, floatEdgeSize-padding, dotRadius)
	}
	dc.Fill()
	dc.Stroke()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)
	return buff.Bytes()
}
