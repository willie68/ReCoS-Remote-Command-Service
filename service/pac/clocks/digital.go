package clocks

import (
	"bytes"
	"image/color"
	"image/png"
	"time"

	"github.com/fogleman/gg"
)

// GenerateDigital generates a nice clock png
func GenerateDigital(timeToRender time.Time, width int, height int, color color.Color, showseconds bool) []byte {
	dc := gg.NewContext(width, height)
	dc.SetColor(color)
	dc.InvertY()
	xWidthDigit := float64(width) / 4.5
	if showseconds {
		xWidthDigit = float64(width) / 7.0
	}
	yHeightDigit := xWidthDigit * 1.5
	segmentThickness := xWidthDigit / 15.0
	yStartDigit := (float64(height) - yHeightDigit) / 2.0
	xStartDigit := 0.0

	hour := timeToRender.Hour()

	myValue := hour
	x := float64(1)
	for x >= 0 {
		xPos := xStartDigit + (xWidthDigit * x)
		yPos := yStartDigit
		digit := myValue % 10
		myValue = myValue / 10
		WriteSegment(digit, xPos, yPos, xWidthDigit, yHeightDigit, segmentThickness, dc, color, true)
		x--
	}

	xPos := xStartDigit + (xWidthDigit * 2.25)
	yPos := yStartDigit + yHeightDigit/3.0
	dc.SetColor(color)
	if (timeToRender.Second()%2) == 1 && !showseconds {
		dc.SetColor(colorDarkSegment)
	}
	dc.DrawCircle(xPos, yPos, segmentThickness)
	yPos = yStartDigit + yHeightDigit*2.0/3.0
	dc.DrawCircle(xPos, yPos, segmentThickness)
	dc.Fill()
	dc.Stroke()

	xDelta := xWidthDigit / 2
	minutes := timeToRender.Minute()

	myValue = minutes
	x = float64(3)
	for x >= 2 {
		digit := myValue % 10
		WriteSegment(digit, xDelta+xStartDigit+(xWidthDigit*x), yStartDigit, xWidthDigit, yHeightDigit, segmentThickness, dc, color, true)
		myValue = myValue / 10
		x--
	}

	if showseconds {
		xPos := xStartDigit + (xWidthDigit * 4.75)
		yPos := yStartDigit + yHeightDigit/3.0
		dc.SetColor(color)
		if (timeToRender.Second()%2) == 1 && !showseconds {
			dc.SetColor(colorDarkSegment)
		}
		dc.DrawCircle(xPos, yPos, segmentThickness)
		yPos = yStartDigit + yHeightDigit*2.0/3.0
		dc.DrawCircle(xPos, yPos, segmentThickness)
		dc.Fill()
		dc.Stroke()

		xDelta := xWidthDigit
		seconds := timeToRender.Second()

		myValue = seconds
		x = float64(5)
		for x >= 4 {
			digit := myValue % 10
			WriteSegment(digit, xDelta+xStartDigit+(xWidthDigit*x), yStartDigit, xWidthDigit, yHeightDigit, segmentThickness, dc, color, true)
			myValue = myValue / 10
			x--
		}
	}

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)

	return buff.Bytes()
}
