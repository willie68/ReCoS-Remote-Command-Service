// +build windows
package pac

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/bmp"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/hardware"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// HardwareMonitorCommandTypeInfo showing hardware sensor data
var HardwareMonitorCommandTypeInfo = models.CommandTypeInfo{
	Category:         "System",
	Type:             "HARDWAREMONITOR",
	Name:             "HardwareMonitor",
	Description:      "Displaying data from a sensors",
	Icon:             "tools.svg",
	WizardPossible:   true,
	WizardActionType: models.Display,
	Parameters: []models.ParamInfo{
		{
			Name:           "sensor",
			Type:           "string",
			Description:    "the sensor name like given above",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "format",
			Type:           "string",
			Description:    "the format string for the textual representation",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "display",
			Type:           "string",
			Description:    "text shows only the textual representation, graph shows both",
			Unit:           "",
			WizardPossible: true,
			List:           []string{"text", "graph"},
		},
		{
			Name:           "ymin",
			Type:           "int",
			Description:    "the value for the floor of the graph",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "ymax",
			Type:           "int",
			Description:    "the value for the top of the graph",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "color",
			Type:           "color",
			Description:    "color of the graph",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

const measurepoints = 500

const imageWidth = measurepoints
const imageHeight = measurepoints

// HardwareMonitorCommand This command connects to the openhardwaremonitor application on windows.
// With this you can get different sensors of your computer. For using the webserver of the openhardwaremonitor
// app, you have to add another external configuration into the main service configuration.
// This command has the following parameters:
// sensor: the sensor name like given above.
// format: the format string for the textual representation
// display: text, graph,  text shows only the textual representation, graph shows both
// ymin: the value for the floor of the graph
// ymax: the value for the bottom of the graph
// color: color of the graph
type HardwareMonitorCommand struct {
	Parameters  map[string]interface{}
	action      *Action
	stop        bool
	ticker      *time.Ticker
	done        chan bool
	datas       []float64
	sensors     []models.Sensor
	yMinValue   float64
	yMaxValue   float64
	yDelta      float64
	color       color.Color
	textonly    bool
	commandName string
}

// EnrichType enrich the type info with the informations from the profile
func (d *HardwareMonitorCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	sensors, err := hardware.OpenHardwareMonitorInstance.GetSensorList()
	if err != nil {
		clog.Logger.Errorf("error in getting sensors: %v", err)
		return HardwareMonitorCommandTypeInfo, nil
	}
	d.sensors = sensors
	index := -1
	for x, parameter := range HardwareMonitorCommandTypeInfo.Parameters {
		if parameter.Name == "sensor" {
			index = x
		}
	}
	if len(d.sensors) != len(HardwareMonitorCommandTypeInfo.Parameters[index].List) {
		HardwareMonitorCommandTypeInfo.Parameters[index].List = make([]string, 0)
		for _, sensor := range d.sensors {
			HardwareMonitorCommandTypeInfo.Parameters[index].List = append(HardwareMonitorCommandTypeInfo.Parameters[index].List, sensor.GetFullSensorName())
		}
	}

	return HardwareMonitorCommandTypeInfo, nil
}

// Init nothing
func (d *HardwareMonitorCommand) Init(a *Action, commandName string) (bool, error) {
	rand.Seed(42)
	d.datas = make([]float64, measurepoints)
	d.commandName = commandName
	object, ok := d.Parameters["sensor"]
	if !ok {
		return false, fmt.Errorf("the sensor parameter is empty")
	}
	sensorname, ok := object.(string)
	if !ok {
		return false, fmt.Errorf("the sensor parameter is in wrong format. Please use string as format")
	}
	value, ok := d.Parameters["ymin"]
	if !ok {
		d.yMinValue = 0
	} else {
		switch v := value.(type) {
		case int:
			d.yMinValue = float64(v)
		case float32:
			d.yMinValue = float64(v)
		case float64:
			d.yMinValue = v
		}
	}

	value, ok = d.Parameters["ymax"]
	if !ok {
		d.yMaxValue = 100
	} else {
		switch v := value.(type) {
		case int:
			d.yMaxValue = float64(v)
		case float32:
			d.yMaxValue = float64(v)
		case float64:
			d.yMaxValue = v
		}
	}
	d.yDelta = d.yMaxValue - d.yMinValue

	color, err := ConvertParameter2Color(d.Parameters, "color", color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF})
	if err != nil {
		return false, fmt.Errorf("the color parameter is in wrong format. Please use string as format")
	}
	d.color = color

	value, ok = d.Parameters["display"]
	if !ok {
		d.textonly = true
	} else {
		display := strings.ToLower(value.(string))
		switch display {
		case "text":
			d.textonly = true
		case "both":
			d.textonly = false
		}
	}
	d.action = a
	d.stop = false
	d.ticker = time.NewTicker(1 * time.Second)
	d.done = make(chan bool)
	go func() {
		for {
			select {
			case <-d.done:
				return
			case <-d.ticker.C:
				sensors, err := hardware.OpenHardwareMonitorInstance.GetSensorList()
				if err != nil {
					clog.Logger.Errorf("error in getting sensors: %v", err)
					continue
				}
				d.sensors = sensors
				index := -1
				for x, parameter := range HardwareMonitorCommandTypeInfo.Parameters {
					if parameter.Name == "sensor" {
						index = x
					}
				}
				if len(d.sensors) != len(HardwareMonitorCommandTypeInfo.Parameters[index].List) {
					HardwareMonitorCommandTypeInfo.Parameters[index].List = make([]string, 0)
					for _, sensor := range d.sensors {
						HardwareMonitorCommandTypeInfo.Parameters[index].List = append(HardwareMonitorCommandTypeInfo.Parameters[index].List, sensor.GetFullSensorName())
					}
				}
				var temp float64
				var value string
				found := false
				for _, sensor := range sensors {
					if strings.EqualFold(sensor.GetFullSensorName(), sensorname) {
						temp = sensor.Value
						value = sensor.ValueStr
						found = true
					}
				}
				if found {
					d.datas = append(d.datas, temp)
					if len(d.datas) > measurepoints {
						d.datas = d.datas[1:]
					}
					if api.HasConnectionWithProfile(a.Profile) {
						if d.textonly {
							message := models.Message{
								Profile: d.action.Profile,
								Action:  d.action.Name,
								Text:    value,
								State:   0,
							}
							api.SendMessage(message)
							continue
						}
						d.SendGraphics(value)
					}
				} else {
					message := models.Message{
						Profile: d.action.Profile,
						Action:  d.action.Name,
						Text:    fmt.Sprintf("Sensor %s not found", sensorname),
						State:   0,
					}
					api.SendMessage(message)
				}
			}
		}
	}()
	return true, nil
}

// Stop nothing
func (d *HardwareMonitorCommand) Stop(a *Action) (bool, error) {
	d.stop = true
	d.done <- true
	return true, nil
}

// Execute nothing
func (d *HardwareMonitorCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	if IsSingleClick(requestMessage) {
		sensors, err := hardware.OpenHardwareMonitorInstance.GetSensorList()
		if err != nil {
			return true, err
		}
		for _, sensor := range sensors {
			clog.Logger.Info(sensor.GetFullSensorName())
		}
		json, err := json.Marshal(sensors)
		if err != nil {
			return true, err
		}
		clog.Logger.Info(string(json))
	}

	return false, nil
}

// GetGraphics creates a clock graphics from the id
func (d *HardwareMonitorCommand) GetGraphics(id string, width int, height int) (models.GraphicsInfo, error) {
	if !hardware.OpenHardwareMonitorInstance.Active {
		return models.GraphicsInfo{}, errors.New("OpenHardwareMonitor is not active")
	}
	if width <= 0 {
		width = imageWidth
	}
	if width > measurepoints {
		width = measurepoints
	}

	if height <= 0 {
		height = imageHeight
	}
	buff := d.generateBMP(width, height)
	model := models.GraphicsInfo{
		Mimetype: "image/bmp",
		Data:     buff,
	}
	return model, nil
}

// SendGraphics sending a new image url to the client
func (d *HardwareMonitorCommand) SendGraphics(value string) {
	id := strconv.Itoa(rand.Intn(999999))
	image := GetImageURL(d.action, d.commandName, id)
	message := models.Message{
		Profile:  d.action.Profile,
		Action:   d.action.Name,
		ImageURL: image,
		Text:     value,
		State:    0,
	}
	api.SendMessage(message)
}

func (d *HardwareMonitorCommand) generateBMP(width int, height int) []byte {
	fHeight := float64(height)
	dc := gg.NewContext(width, height)
	dc.SetColor(d.color)
	dc.InvertY()
	dc.MoveTo(0, 0)
	xLast := 0.0
	for x := 0; x < width; x++ {
		index := measurepoints - width + x
		temp := d.datas[index]

		newX, y := d.projectPoint(float64(x), temp, fHeight)
		dc.LineTo(newX, y)
		xLast = newX
	}
	dc.LineTo(xLast, 0)
	dc.LineTo(0, 0)
	dc.Stroke()
	mycolor := d.color
	r, g, b, a := mycolor.RGBA()
	mycolor = color.NRGBA{
		R: uint8(r) / 2,
		G: uint8(g) / 2,
		B: uint8(b) / 2,
		A: uint8(a),
	}
	dc.SetColor(mycolor)
	dc.MoveTo(0, 0)
	xLast = 0.0
	for x := 0; x < width; x++ {
		index := measurepoints - width + x
		temp := d.datas[index]
		newX, y := d.projectPoint(float64(x), temp, fHeight)
		dc.LineTo(newX, y-1.0)
		xLast = newX
	}
	dc.LineTo(xLast, 0)
	dc.LineTo(0, 0)
	dc.Fill()
	dc.Stroke()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	bmp.Encode(&buff, myImage)

	return buff.Bytes()
}

func (d *HardwareMonitorCommand) projectPoint(x float64, y float64, height float64) (xDest float64, yDest float64) {
	xDest = x

	if y < d.yMinValue {
		yDest = 0
	} else if y > d.yMaxValue {
		yDest = d.yMaxValue
	} else {
		yDest = (float64(height) / d.yDelta) * (y - d.yMinValue)
	}
	return
}
