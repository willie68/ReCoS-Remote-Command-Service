package models

import (
	"fmt"
	"strconv"
	"strings"
)

// SensorCategorie is the general type of the sensor
type SensorCategorie string

const (
	// MainBoard sensors are  the category holding all sensors of the mainboard
	MainBoard SensorCategorie = "MainBoard"
	// CPU summarizes all sensors of the CPU
	CPU = "CPU"
	// Memory summarizes all sensors of the installed memory
	Memory = "Memory"
	// GPU summarizes all sensors of the installed grafic cards
	GPU = "GPU"
	// Storage summarizes all sensors of the installed storage
	Storage = "Storage"
	// Unknown summarizes all unknown sensors
	Unknown = "Unknown"
)

// SensorType is the general type of the sensor
type SensorType string

const (
	// Load sensor type holds the load of the specific hardware, value in %
	Load SensorType = "Load"
	// Temperature sensor type holds the temperature of the specific hardware, value in °C
	Temperature = "Temperature"
	// Clocks sensor type holds the Clocks of the specific hardware, value in MHz
	Clocks = "Clocks"
	// Powers sensor type holds the Powers of the specific hardware, value in W
	Powers = "Powers"
	// Data sensor type holds different datas of the specific hardware
	Data = "Data"
	// Voltages sensor type holds the Voltages of the specific hardware, value in V
	Voltages = "Voltages"
	// Fans sensor type holds the fans (e.g. the rpms) of the specific hardware, value in RPM
	Fans = "Fans"
)

// Sensor this is the sturcture of one sensor
type Sensor struct {
	Categorie    SensorCategorie
	Hardwarename string
	Type         SensorType
	Name         string
	Value        float64
	Min          float64
	Max          float64
	ValueStr     string
	MinStr       string
	MaxStr       string
}

// GetFullSensorName getting an unique sensor name
func (s *Sensor) GetFullSensorName() string {
	return fmt.Sprintf("%s/%s/%s/%s", s.Categorie, s.Hardwarename, s.Type, s.Name)
}

// ParseValues parses all values dependuing on thier type
func (s *Sensor) ParseValues() {
	valueStr := s.ValueStr
	if valueStr == "" {
		return
	}
	if strings.Index(valueStr, " ") >= 0 {
		valueStr = strings.Split(valueStr, " ")[0]
	}
	if strings.Index(valueStr, ",") >= 0 {
		valueStr = strings.ReplaceAll(valueStr, ",", ".")
	}
	value, err := strconv.ParseFloat(valueStr, 32)
	if err != nil {
		value = 0.0
	}
	s.Value = value
}
