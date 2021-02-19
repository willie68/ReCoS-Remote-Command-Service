// Special implementation for openhardwaremonitor app
// +build windows
package hardware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

type OpenHardwareMonitor struct {
	baseURL    string
	Connected  bool
	periode    int
	lastError  error
	Sensorlist []models.Sensor
	ticker     *time.Ticker
	done       chan bool
	m          sync.Mutex
}

var OpenHardwareMonitorInstance OpenHardwareMonitor
var once sync.Once

// InitOpenHardwareMonitor initialise the open hardware monitor connection
func InitOpenHardwareMonitor(extconfig map[string]interface{}) error {
	var err error
	value, ok := extconfig["openhardwaremonitor"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Info("found config")
			url, ok := config["url"].(string)
			if !ok {
				err = fmt.Errorf("can't find url to connect to. %s", url)
			}
			updatePeriod, ok := config["updateperiod"].(int)
			if !ok {
				updatePeriod = 5
			}
			OpenHardwareMonitorInstance = OpenHardwareMonitor{
				periode:   updatePeriod,
				Connected: false,
				baseURL:   url,
				ticker:    time.NewTicker(time.Duration(updatePeriod) * time.Second),
				done:      make(chan bool),
			}
		}
	}
	return err
}

func (o *OpenHardwareMonitor) Connect() error {
	err := OpenHardwareMonitorInstance.updateSensorList()
	if err != nil {
		return err
	}
	for _, sensor := range OpenHardwareMonitorInstance.Sensorlist {
		clog.Logger.Infof("found sensor with name: %s", sensor.GetFullSensorName())
	}
	writingSensorList(OpenHardwareMonitorInstance.Sensorlist)
	go func() {
		for {
			select {
			case <-OpenHardwareMonitorInstance.done:
				return
			case <-OpenHardwareMonitorInstance.ticker.C:
				OpenHardwareMonitorInstance.lastError = OpenHardwareMonitorInstance.updateSensorList()
			}
		}
	}()
	return nil
}

// GetSensorList getting a sensor list from openhardwaremonitor
func (o *OpenHardwareMonitor) GetSensorList() ([]models.Sensor, error) {
	o.m.Lock()
	defer o.m.Unlock()
	if !o.Connected {
		clog.Logger.Debug("no connection")
	}
	return o.Sensorlist, o.lastError
}

// GetPeriod getting the update periode
func (o *OpenHardwareMonitor) GetPeriod() int {
	return o.periode
}

// GetSensorList getting a sensor list from openhardwaremonitor
func (o *OpenHardwareMonitor) updateSensorList() error {
	response, err := http.Get(o.baseURL)
	if err != nil {
		o.Connected = false
		return fmt.Errorf("error getting sensor values. %v", err)
	}
	defer response.Body.Close()
	o.Connected = true
	var target map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&target)
	if err != nil {
		return fmt.Errorf("error reading sensor values. %v", err)
	}
	values := target["Children"]
	if values == nil {
		return fmt.Errorf("error reading sensor values. %v", err)
	}
	computers := values.([]interface{})
	computer := computers[0].(map[string]interface{})

	values = computer["Children"]
	if values == nil {
		return fmt.Errorf("no sensors availble.")
	}

	sensors := make([]models.Sensor, 0)

	categories := values.([]interface{})
	for _, values := range categories {
		categorie := values.(map[string]interface{})
		name := categorie["Text"].(string)
		imageURL := categorie["ImageURL"].(string)
		sensorCat := string2SensorCategorie(imageURL)

		children := categorie["Children"].([]interface{})
		for _, values := range children {
			sensorType := values.(map[string]interface{})
			sensorTypeName := sensorType["Text"].(string)
			imageURL := sensorType["ImageURL"].(string)

			children := sensorType["Children"].([]interface{})
			for _, values := range children {
				sensor := values.(map[string]interface{})
				sensorName := sensor["Text"].(string)

				sensorObject := models.Sensor{
					Categorie:    sensorCat,
					Hardwarename: name,
					Type:         string2SensorType(sensorTypeName, imageURL),
					Name:         sensorName,
					ValueStr:     sensor["Value"].(string),
					MinStr:       sensor["Min"].(string),
					MaxStr:       sensor["Max"].(string),
				}
				sensorObject.ParseValues()
				sensors = append(sensors, sensorObject)
			}
		}
	}
	o.m.Lock()
	defer o.m.Unlock()
	o.Sensorlist = sensors
	return nil
}

func writingSensorList(sensorlist []models.Sensor) {
	f, err := os.Create("sensorlist.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString("Sensor full name;Category;Hardware;Type;Sensor name;Value (e.g.)\r\n")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	for _, sensor := range sensorlist {
		_, err := f.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s;%s\r\n", sensor.GetFullSensorName(), sensor.Categorie, sensor.Hardwarename, sensor.Type, sensor.Name, sensor.ValueStr))
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func string2SensorCategorie(value string) models.SensorCategorie {
	if strings.Contains(value, "mainboard") {
		return models.MainBoard
	} else if strings.Contains(value, "cpu") {
		return models.CPU
	} else if strings.Contains(value, "ram") {
		return models.Memory
	} else if strings.Contains(value, "ati") || strings.Contains(value, "nvidia") {
		return models.GPU
	} else if strings.Contains(value, "hdd") {
		return models.Storage
	}
	return models.Unknown
}

func string2SensorType(text, image string) models.SensorType {
	switch text {
	case "Temperatures":
		{
			return models.Temperature
		}
	case "Clocks":
		{
			return models.Clocks
		}
	case "Load":
		{
			return models.Load
		}
	case "Powers":
		{
			return models.Powers
		}
	case "Voltages":
		{
			return models.Voltages
		}
	case "Fans":
		{
			return models.Fans
		}
	}
	return models.Data
}
