package clocks

import (
	"bytes"
	"image/color"
	"image/png"
	"time"

	"github.com/fogleman/gg"
)

var (
	clockBerlinCase   color.Color = color.RGBA{R: 0xB0, G: 0xB0, B: 0xB0, A: 0xFF}
	clockBerlinYellow color.Color = color.RGBA{R: 0xDC, G: 0x88, B: 0x00, A: 0xFF}
	clockBerlinRed    color.Color = color.RGBA{R: 0xFF, G: 0x2A, B: 0x03, A: 0xFF}
)

// generateBerlin generates a nice clock png
func GenerateBerlin(timeToRender time.Time, width int, height int) []byte {
	fheight := float64(height)
	fwidth := float64(width)
	if fwidth > (0.8 * fheight) {
		fwidth = 0.8 * fheight
	} else {
		fheight = 1.25 * fwidth
	}
	dc := gg.NewContext(int(fwidth), int(fheight))
	dc.SetColor(clockBerlinCase)
	dc.SetLineCapButt()

	halfWidth := fwidth / 2.0
	caseWidth := fheight / 40.0
	halfCaseWidth := caseWidth / 2.0
	bWidth := fwidth/4.0 - halfCaseWidth
	//	dc.DrawRectangle(0.0, 0.0, fwidth, fheight)
	//	dc.Stroke()
	dc.SetLineWidth(caseWidth)
	secondRadius := 0.1 * fheight
	yPos := secondRadius + caseWidth
	dc.DrawCircle(halfWidth, yPos, secondRadius)

	dc.Stroke()
	dc.SetLineWidth(2 * caseWidth)
	dc.DrawLine(halfWidth-(secondRadius/2.0), yPos+secondRadius+caseWidth, halfWidth+(secondRadius/2.0), yPos+secondRadius+caseWidth)

	dc.Stroke()
	dc.SetLineWidth(caseWidth)
	quadHeight := 0.15*fheight - caseWidth
	yPos = 0.15*fheight + yPos
	hour := timeToRender.Hour()
	tHour := hour / 5
	for i := 0; i < 4; i++ {
		if i >= tHour {
			dc.SetColor(colorDarkSegment)
		} else {
			dc.SetColor(clockBerlinRed)
		}
		xPos := float64(i) * (fwidth / 4.0)
		dc.DrawRoundedRectangle(xPos, yPos, bWidth, quadHeight, 0.03*fheight)
		dc.Fill()
	}
	dc.Stroke()
	dc.SetColor(clockBerlinCase)
	dc.DrawRoundedRectangle(halfCaseWidth, yPos, fwidth-caseWidth, quadHeight, 0.03*fheight)
	dc.DrawLine(0.25*fwidth, yPos, 0.25*fwidth, yPos+quadHeight)
	dc.DrawLine(0.5*fwidth, yPos, 0.5*fwidth, yPos+quadHeight)
	dc.DrawLine(0.75*fwidth, yPos, 0.75*fwidth, yPos+quadHeight)
	dc.Stroke()

	yPos = 0.15*fheight + yPos
	barWidth := 0.23 * fwidth
	barHeight := 2.0 * caseWidth
	//barRadius := 0.03 * fheight
	dc.SetLineWidth(barHeight)
	dc.DrawLine(0.15*fwidth, yPos, 0.15*fwidth+barWidth, yPos)
	dc.DrawLine(fwidth-(0.15*fwidth)-barWidth, yPos, fwidth-(0.15*fwidth), yPos)

	dc.Stroke()
	dc.SetLineWidth(caseWidth)
	yPos = yPos + 1.5*caseWidth

	eHour := hour % 5
	for i := 0; i < 4; i++ {
		if i >= eHour {
			dc.SetColor(colorDarkSegment)
		} else {
			dc.SetColor(clockBerlinRed)
		}
		xPos := float64(i) * (fwidth / 4.0)
		dc.DrawRoundedRectangle(xPos, yPos, bWidth, quadHeight, 0.03*fheight)
		dc.Fill()
	}
	dc.Stroke()
	dc.SetColor(clockBerlinCase)

	dc.DrawRoundedRectangle(halfCaseWidth, yPos, fwidth-caseWidth, quadHeight, 0.03*fheight)
	dc.DrawLine(0.25*fwidth, yPos, 0.25*fwidth, yPos+quadHeight)
	dc.DrawLine(0.5*fwidth, yPos, 0.5*fwidth, yPos+quadHeight)
	dc.DrawLine(0.75*fwidth, yPos, 0.75*fwidth, yPos+quadHeight)
	dc.Stroke()

	yPos = 0.15*fheight + yPos
	dc.SetLineWidth(barHeight)
	dc.SetLineCapButt()
	dc.DrawLine(0.15*fwidth, yPos, 0.15*fwidth+barWidth, yPos)
	dc.DrawLine(fwidth-(0.15*fwidth)-barWidth, yPos, fwidth-(0.15*fwidth), yPos)
	dc.Stroke()

	yPos = yPos + 1.5*caseWidth
	minutes := timeToRender.Minute()
	tMinute := minutes / 5
	for i := 0; i < 11; i++ {
		if i >= tMinute {
			dc.SetColor(colorDarkSegment)
		} else {
			if (i % 3) == 2 {
				dc.SetColor(clockBerlinRed)
			} else {
				dc.SetColor(clockBerlinYellow)
			}
		}
		xPos := float64(i) * (fwidth / 11.0)
		dc.DrawRoundedRectangle(xPos, yPos, fwidth/11.0, quadHeight, 0.03*fheight)
		dc.Fill()
	}
	dc.Stroke()

	dc.SetColor(clockBerlinCase)
	dc.SetLineWidth(caseWidth)
	dc.DrawRoundedRectangle(halfCaseWidth, yPos, fwidth-caseWidth, quadHeight, 0.03*fheight)
	for i := 1; i < 11; i++ {
		xPos := (0.09 * float64(i)) * fwidth
		dc.DrawLine(xPos, yPos, xPos, yPos+quadHeight)
	}
	dc.Stroke()

	yPos = 0.15*fheight + yPos
	dc.SetLineWidth(barHeight)
	dc.SetLineCapButt()
	dc.DrawLine(0.15*fwidth, yPos, 0.15*fwidth+barWidth, yPos)
	dc.DrawLine(fwidth-(0.15*fwidth)-barWidth, yPos, fwidth-(0.15*fwidth), yPos)
	dc.Stroke()

	yPos = yPos + 1.5*caseWidth
	eMinute := minutes % 5
	for i := 0; i < 4; i++ {
		if i >= eMinute {
			dc.SetColor(colorDarkSegment)
		} else {
			dc.SetColor(clockBerlinYellow)
		}
		xPos := float64(i) * (fwidth / 4.0)
		dc.DrawRoundedRectangle(xPos, yPos, bWidth, quadHeight, 0.03*fheight)
		dc.Fill()
	}
	dc.Stroke()
	dc.SetColor(clockBerlinCase)
	dc.SetLineWidth(caseWidth)
	dc.DrawRoundedRectangle(halfCaseWidth, yPos, fwidth-caseWidth, quadHeight, 0.03*fheight)
	dc.DrawLine(0.25*fwidth, yPos, 0.25*fwidth, yPos+quadHeight)
	dc.DrawLine(0.5*fwidth, yPos, 0.5*fwidth, yPos+quadHeight)
	dc.DrawLine(0.75*fwidth, yPos, 0.75*fwidth, yPos+quadHeight)
	//dc.Fill()
	dc.Stroke()

	if (timeToRender.Second() % 2) == 1 {
		dc.SetColor(colorDarkSegment)
	} else {
		dc.SetColor(clockBerlinYellow)
	}
	dc.DrawCircle(halfWidth, secondRadius+caseWidth, secondRadius-0.5*caseWidth)
	dc.Fill()
	dc.Stroke()

	//	minutes := timeToRender.Minute()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)

	return buff.Bytes()
}
