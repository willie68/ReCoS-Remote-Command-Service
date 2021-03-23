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
