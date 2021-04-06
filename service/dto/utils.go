package dto

import (
	"archive/zip"
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"
	"sync"
	"syscall"
)

func parseHexColor(s string) (c color.RGBA, err error) {
	s = strings.TrimPrefix(s, "#")
	c.A = 0xff
	switch len(s) {
	case 8:
		_, err = fmt.Sscanf(s, "%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	return
}

func deg2Rad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func rad2Deg(rad float64) float64 {
	return rad * (180.0 / math.Pi)
}

func GetImageURL(action *Action, commandName string, id string) string {
	return fmt.Sprintf("/api/v1/show/%s/%s/%s/%s", action.Profile, action.Name, commandName, id)
}

func ConvertParameter2Bool(parameters map[string]interface{}, parameterName string, defaultValue bool) (bool, error) {
	myBool := defaultValue
	value, found := parameters[parameterName]
	if found {
		var ok bool
		myBool, ok = value.(bool)
		if !ok {
			return false, fmt.Errorf("%s is in wrong format. Please use string as format", parameterName)
		}
	}
	return myBool, nil
}

func ConvertParameter2Int(parameters map[string]interface{}, parameterName string, defaultValue int) (int, error) {
	valueInt := defaultValue
	value, found := parameters[parameterName]
	if found {
		var ok bool
		valueInt, ok = value.(int)
		if !ok {
			fValue, ok := value.(float64)
			if ok {
				valueInt = int(fValue)
			} else {
				return 0, fmt.Errorf("%s is in wrong format. Please use string as format", parameterName)
			}
		}
	}
	return valueInt, nil
}

func ConvertParameter2String(parameters map[string]interface{}, parameterName string, defaultValue string) (string, error) {
	valueStr := defaultValue
	value, found := parameters[parameterName]
	if found {
		var ok bool
		valueStr, ok = value.(string)
		if !ok {
			return "", fmt.Errorf("%s is in wrong format. Please use string as format", parameterName)
		}
	}
	return valueStr, nil
}

func ConvertParameter2StringArray(parameters map[string]interface{}, parameterName string) ([]string, error) {
	values := make([]string, 0)
	value, found := parameters[parameterName]
	if found {
		for _, iValue := range value.([]interface{}) {
			myValue, ok := iValue.(string)
			if ok {
				values = append(values, myValue)
			}
		}
	}
	return values, nil
}

func ConvertParameter2Color(parameters map[string]interface{}, parameterName string, defaultValue color.Color) (color.Color, error) {
	valueColor := defaultValue
	value, found := parameters[parameterName]
	if found && value != "" {
		var err error
		valueColor, err = parseHexColor(value.(string))
		if err != nil {
			return color.Black, fmt.Errorf("%s is in wrong format. Please use string as format", parameterName)
		}
	}
	return valueColor, nil
}

var zones []string
var loadIANAOnce sync.Once

func GetIANANames() []string {
	loadIANAOnce.Do(func() {
		env, _ := syscall.Getenv("ZONEINFO")
		zones = make([]string, 0)
		r, err := zip.OpenReader(env)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		// Iterate through the files in the archive,
		// printing some of their contents.
		for _, f := range r.File {
			if !strings.HasSuffix(f.Name, "/") {
				zones = append(zones, f.Name)
			}
		}
	})
	return zones
}
