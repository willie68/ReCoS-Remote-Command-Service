package clocks

import (
	"bytes"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	colorTicks    color.Color = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	colorArHour   color.Color = color.RGBA{R: 0x40, G: 0x40, B: 0x80, A: 0xFF}
	colorArMinute color.Color = color.RGBA{R: 0x40, G: 0x40, B: 0x80, A: 0xFF}
	colorArSecond color.Color = color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xFF}
	tickLength    float64     = 15
)

// GenerateAnalog generates a nice clock bmp
func GenerateAnalog(timeToRender time.Time, width int, height int, showseconds bool, showDate bool, fontsize int, dateformat string) []byte {
	edgeSize := height
	if width < height {
		edgeSize = width
	}
	dc := gg.NewContext(edgeSize, edgeSize)
	halfWidth := float64(edgeSize / 2)
	halfHeight := float64(edgeSize / 2)
	floatEdgeSize := float64(edgeSize)
	myTicklength := tickLength * floatEdgeSize / float64(ClockImageHeight)

	if showDate {
		font, err := truetype.Parse(goregular.TTF)
		if err != nil {
			log.Fatal(err)
		}

		face := truetype.NewFace(font, &truetype.Options{Size: float64(fontsize)})
		dc.SetFontFace(face)
		dc.SetColor(colorArSecond)
		dc.SetLineWidth(2.0)
		dateStr := timeToRender.Format(dateformat)

		dc.DrawStringAnchored(dateStr, halfWidth, halfHeight/2.0, 0.5, 0.5)

		dc.Stroke()
	}
	dc.SetColor(colorTicks)
	dc.InvertY()

	dc.SetLineWidth(1.0)
	dc.DrawCircle(halfWidth, halfHeight, halfHeight-1)
	dc.MoveTo(halfWidth-1, floatEdgeSize)
	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)

	for i := 0; i < 12; i++ {
		deg := float64(30.0 * i)
		dc.MoveTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-myTicklength)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-myTicklength)))
		dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-1)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-1)))
	}
	dc.Stroke()

	dc.SetLineWidth(4.0)

	dc.MoveTo(halfWidth-1, floatEdgeSize)
	dc.LineTo(halfWidth-1, floatEdgeSize-myTicklength)
	dc.MoveTo(halfWidth-1, 0)
	dc.LineTo(halfWidth-1, myTicklength)

	dc.MoveTo(0, halfHeight-1)
	dc.LineTo(myTicklength, halfHeight-1)
	dc.MoveTo(floatEdgeSize-myTicklength, halfHeight-1)
	dc.LineTo(floatEdgeSize, halfHeight-1)
	dc.Stroke()

	if showseconds {
		dc.SetColor(colorArSecond)
		dc.SetLineWidth(1.0)
		seconds := timeToRender.Second()

		deg := float64(6.0 * seconds)
		dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*10), halfHeight-(math.Cos(deg2Rad(deg))*10))
		dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-2)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-2)))

		dc.Stroke()
	}

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

func deg2Rad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func rad2Deg(rad float64) float64 {
	return rad * (180.0 / math.Pi)
}
