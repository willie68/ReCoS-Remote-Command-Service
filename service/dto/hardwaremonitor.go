// +build windows
package dto

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/color"
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
var HardwareMonitorCommandTypeInfo = models.CommandTypeInfo{"HARDWAREMONITOR", "HardwareMonitor", "getting data from the hardware sensors", []models.CommandParameterInfo{
	{"sensor", "string", "the sensor name like given above", "", make([]string, 0)},
	{"format", "string", "the format string for the textual representation", "", make([]string, 0)},
	{"display", "string", "text shows only the textual representation, graph shows both", "", []string{"text", "graph"}},
	{"ymin", "int", "the value for the floor of the graph", "", make([]string, 0)},
	{"ymax", "int", "the value for the top of the graph", "", make([]string, 0)},
	{"color", "color", "color of the graph", "", make([]string, 0)},
}}

const measurepoints = 100
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
	Parameters map[string]interface{}
	action     *Action
	stop       bool
	ticker     *time.Ticker
	done       chan bool
	temps      []float64
	sensors    []models.Sensor
	yMinValue  float64
	yMaxValue  float64
	yDelta     float64
	color      color.Color
	textonly   bool
}

// Init nothing
func (d *HardwareMonitorCommand) Init(a *Action) (bool, error) {
	d.temps = make([]float64, measurepoints)
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

	value, ok = d.Parameters["color"]
	if !ok {
		d.color = color.RGBA{
			R: 255, G: 0, B: 0, A: 0,
		}
	}
	color, err := parseHexColor(value.(string))
	if err != nil {
		clog.Logger.Errorf("error in getting sensors: %v", err)
		return false, err
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
				sensorname := d.Parameters["sensor"].(string)
				var temp float64
				var value string
				for _, sensor := range sensors {
					if strings.EqualFold(sensor.GetFullSensorName(), sensorname) {
						temp = sensor.Value
						value = sensor.ValueStr
					}
				}
				d.temps = append(d.temps, temp)
				if len(d.temps) > measurepoints {
					d.temps = d.temps[1:]
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
					d.SendPNG(value)
				}
			}
		}
	}()
	return true, nil
}

// Stop nothing
func (d *HardwareMonitorCommand) Stop(a *Action) (bool, error) {
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
		clog.Logger.Info(string(json))
	}

	return false, nil
}

// SendPNG sending this array to the client
func (d *HardwareMonitorCommand) SendPNG(value string) {
	dc := gg.NewContext(imageWidth, imageHeight)
	dc.SetColor(d.color)
	dc.InvertY()
	dc.MoveTo(0, 0)
	xLast := 0.0
	for index, temp := range d.temps {
		x, y := d.projectPoint(float64(index), temp)
		dc.LineTo(x, y)
		xLast = x
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
	for index, temp := range d.temps {
		x, y := d.projectPoint(float64(index), temp)
		dc.LineTo(x, y-1.0)
		xLast = x
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

	// Encode the bytes in the buffer to a base64 string
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	// You can embed it in an html doc with this string
	image := "data:image/bmp;base64," + encodedString

	//image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEgAAABICAYAAABV7bNHAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAABTBJREFUeNrsmz1s01AQx5+Tpl+0JUNZkIAMSDQLFLGwIMLEhJSyIYaWgQ4sbYWYAQmBkJCAhYUBOrAw0CAmGFAQDCygIIYGJKTQBSQ6hEJTmk/unPfSF5O2sX12/ML7S69xktZ2f+/u3p19ZkxLS0tLS0tVGZ08+KmJ04ktvs4/W3iS+W8AAYwYvCRhHIcxDiPW5p9m+HgFIwXQ8l0DCKBE4WUKxiSHQqEUjKcA6qGygDiYWRgzMKLyd4OGwfb19LCxSISNhkNsVyjUch9fyxW2XK2wJXhdLJVa/UoOxjyMO15aleEBnCtWMAjlWH8fO9Lba4JxonfFInsP4/Wf9Vag5gBSKtCAAAy60APZldBCJgYH2bG+PrITLtRq7PnaHxhr5rakNIwJamsyiOCgO932EkyboPIcUjowgADOAx6ITZ0c6DfhoFv5oeVqld3/9dsap85RBfEwFRwEcmFkGAANsIjhX3ol4hvOdXYDUvJAPB79nF183jFAMhx0qUsjIyzuMABTCI+9F1bHj8US45iOAiQDIKV9ByTDwZO6HN0JkMKs09oN53CwN8LerhcFpARA+gqQMr7FIIAzxVcr03KuRaO+xZt2tVQus+s/V+Tgfdhp2RJyuJSbUGaHhwMHR1j12aEd8kcLPHn1FpCAgzo/PGSeSFCFKQauplwxOQ3xBBDPkMfFUo5ZcdA1MTggT+LUNlcPnAPi5jkjJ4GqaLrZ1S57ZUGzorbyMwmkikf1PKm+qsFkJ0kBWa3Hy/LBK53d0WRFM9QWlJStR0VtZNsNKxqnBDTTOIiC1rMRsJsmd5IEEL9MatKWZkBJjYZC8oqWpLKgxrKowrLeTm4k8qJ23Ww7QMfFxlgHC1G6gran5eS7ATQuKuVukCXzP0QGaKxLAFkmO+YKEA/QTOQ/3aK9PWEyF2sA2uzWjJo5kb3/xVU5bkBuEd6zJ5Agyp8+kezHFSCEM3TxYiAB5aenW+dDNsNF9/hOmypUa/5ZEJrxZjMVWEC1mrYgSm0FqHGRGxsJukXYECHClCtA8j1uaafKS5rsDIWLpfHHYqncNYDwlhAloIzYqd3gFkRlm+/ff6AA9EpsbNLEpJSwx8jqHSQuhnrzb+OSyoAyEGNzrgHxQJ0SO1fZzfD8lytV8XaeMg9q7AwbllTVi+ZzT5EB4r1/uTqgNSWtCIOzFEMftutedjLpqyJNV9GKFgpr8tu75KUGb2fL1Q9WMNveVIo9FuvJkAPimhMb2BOoSmEqnWteeIIdtd0W9jm7mD0Qj+M16rG6BRmBv5h/D+AsVRqlxU0nvdR2q/lzoshDV3u9Htzc6NHqqpz3pAHOFepqfrO86ETjJH6vyrVNYIQTJy0mZu+0033Z7rwEV/uOjZGwmcTQhw2T2Di5MyAX9hGOJe6csLOsuwbEIWVkSC+hDMEu130dbslDt3q8WpA/OuO2695x764MCd+/56XIwQ7cw8fj3lr5ZVqzZDlnKB5wcdXcbIX0BeLRu2KJ7QdL8svlMBDf+LnCvlUqVrdKU+yf6mEWXP4XmHSzEdtlsCdn1CNQWD5ghmy5DINQSJ/4oXwcCrvQsElyVv4cQZ3s7ydrGUaLwcLTAgaB3HW6lPsCSAKV4KAS8ucICHuMjvRGbMNCKGgxlksWcmU+52al8hXQdqCE4vyRzNFQ6zCI+dUPyNi3yLOwPpynfDbMV0CW+DTJA3nM5e6w0MTrUymvLMZ3QC1goUUd4rASW/x6nm08Ev6Blws5pqWlpaWlpaWlpRUE/RVgAAD3F9MyT2oUAAAAAElFTkSuQmCC"
	message := models.Message{
		Profile:  d.action.Profile,
		Action:   d.action.Name,
		ImageURL: image,
		Text:     value,
		State:    0,
	}
	api.SendMessage(message)
}

func (d *HardwareMonitorCommand) projectPoint(x float64, y float64) (xDest float64, yDest float64) {
	xDest = x

	if y < d.yMinValue {
		yDest = 0
	} else if y > d.yMaxValue {
		yDest = d.yMaxValue
	} else {
		yDest = (float64(imageHeight) / d.yDelta) * (y - d.yMinValue)
	}
	return
}
