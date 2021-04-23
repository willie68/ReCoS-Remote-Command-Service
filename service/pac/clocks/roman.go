package clocks

import (
	"bytes"
	"image/color"
	"math"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/bmp"
)

var (
	clockRomanColorTicks    color.Color = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	clockRomanColorArHour   color.Color = color.RGBA{R: 0x40, G: 0x40, B: 0x80, A: 0xFF}
	clockRomanColorArMinute color.Color = color.RGBA{R: 0x40, G: 0x40, B: 0x80, A: 0xFF}
	clockRomanColorArSecond color.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xFF}
	clockRomanTickLength    float64     = 10
	romanDigit                          = []string{"XII", "I", "II", "III", "IIII", "V", "VI", "VII", "VIII", "IX", "X", "XI"}
)

// GenerateRoman generates a nice clock bmp
func GenerateRoman(timeToRender time.Time, width int, height int) []byte {
	edgeSize := height
	if width < height {
		edgeSize = width
	}
	dc := gg.NewContext(edgeSize, edgeSize)
	halfWidth := float64(edgeSize / 2)
	halfHeight := float64(edgeSize / 2)
	floatEdgeSize := float64(edgeSize)
	myTicklength := tickLength * floatEdgeSize / float64(ClockImageHeight)

	dc.SetColor(colorTicks)
	dc.InvertY()

	dc.SetLineWidth(1.0)
	dc.DrawCircle(halfWidth, halfHeight, halfHeight-1)
	dc.MoveTo(halfWidth-1, floatEdgeSize)
	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)
	dc.DrawCircle(halfWidth, halfHeight, halfHeight-tickLength)

	dc.Stroke()
	dc.SetLineCapButt()
	for i := 0; i < 60; i++ {
		deg := float64(6.0 * i)
		if i%5 == 0 {
			dc.SetLineWidth(5.0)
			hour := i / 5
			xPos := halfWidth + (math.Sin(deg2Rad(deg)) * (halfWidth - 2*myTicklength))
			yPos := halfHeight + (math.Cos(deg2Rad(deg)) * (halfHeight - 2*myTicklength))
			s := romanDigit[hour]
			dc.DrawStringAnchored(s, xPos, yPos, 0.5, 0.5)
		} else {
			dc.SetLineWidth(1)
		}
		dc.MoveTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-myTicklength)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-myTicklength)))
		dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-1)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-1)))
		dc.Stroke()
	}

	dc.SetLineWidth(1.0)
	//	dc.SetLineWidth(4.0)

	//	dc.MoveTo(halfWidth-1, floatEdgeSize)
	//	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)
	//	dc.MoveTo(halfWidth-1, 0)
	//	dc.LineTo(halfWidth-1, myTicklength)

	//	dc.MoveTo(0, halfHeight-1)
	//	dc.LineTo(myTicklength, halfHeight-1)
	//	dc.MoveTo(floatEdgeSize-myTicklength, halfHeight-1)
	//	dc.LineTo(floatEdgeSize, halfHeight-1)
	//	dc.Stroke()

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
