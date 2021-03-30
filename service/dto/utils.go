package dto

import (
	"fmt"
	"image/color"
	"math"
	"strings"
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

func ConvertParameter2String(parameters map[string]interface{}, parameterName string) (string, error) {
	valueStr := ""
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

func ConvertParameter2Int(parameters map[string]interface{}, parameterName string) (int, error) {
	valueInt := 0
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
