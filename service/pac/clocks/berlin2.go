package clocks

import (
	"bytes"
	"fmt"
	"time"

	svg "github.com/ajstarks/svgo"
)

var (
	clockBerlin2Case        string = "rgb(176,176,176)"
	clockBerlin2Yellow      string = "rgb(220,136,00)"
	clockBerlin2Red         string = "rgb(255,42,03)"
	clockBerlin2DarkSegment string = "rgb(32,32,32)"
	clockBerlin2Segments    string = "rgb(255,0,0)"
	clockBerlin2Black       string = "rgb(0,0,0)"
)

// generateBerlin generates a nice clock png
func GenerateBerlin2(timeToRender time.Time, width int, height int) []byte {
	if width > height {
		width = height
	} else {
		height = width
	}
	var b bytes.Buffer

	halfWidth := width / 2

	dc := svg.New(&b)
	dc.Start(width, height)

	caseWidth := height / 40
	secondRadius := height / 12
	halfCaseWidth := caseWidth / 2
	bWidth := width/4 - halfCaseWidth

	yPos := secondRadius + caseWidth

	// drawing the case
	// case for second light
	var style string
	if (timeToRender.Second() % 2) == 1 {
		style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2DarkSegment, caseWidth, clockBerlin2Case)
	} else {
		style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2Yellow, caseWidth, clockBerlin2Case)
	}
	dc.Circle(halfWidth, yPos, secondRadius, style)

	style = fmt.Sprintf("fill:none;stroke-width:%d;stroke:%s", 2*caseWidth, clockBerlin2Case)
	dc.Line(halfWidth-(secondRadius/2), yPos+secondRadius+caseWidth, halfWidth+(secondRadius/2), yPos+secondRadius+caseWidth, style)

	// Hour case 10's
	quadHeight := height/6 - caseWidth
	yPos = secondRadius + yPos + caseWidth
	hour := timeToRender.Hour()
	tHour := hour / 5
	for i := 0; i < 4; i++ {
		if i >= tHour {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2DarkSegment, caseWidth, clockBerlin2Case)
		} else {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2Red, caseWidth, clockBerlin2Case)
		}
		xPos := (i * width / 4) + halfCaseWidth
		dc.Roundrect(xPos, yPos, bWidth, quadHeight, height/33, height/33, style)
	}

	style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", "none", caseWidth, clockBerlin2Case)
	radius := height / 33
	dc.Roundrect(halfCaseWidth, yPos, width-caseWidth, quadHeight, radius, radius, style)
	dc.Line(width/4+halfCaseWidth, yPos, width/4+halfCaseWidth, yPos+quadHeight, style)
	dc.Line(width/2+halfCaseWidth, yPos, width/2+halfCaseWidth, yPos+quadHeight, style)
	dc.Line(((3*width)/4)+halfCaseWidth, yPos, ((3*width)/4)+halfCaseWidth, yPos+quadHeight, style)

	yPos = height/6 + yPos
	barWidth := width / 4
	//barRadius := 0.03 * fheight
	style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", "none", 2*caseWidth, clockBerlin2Case)
	dc.Line(width/6, yPos, width/6+barWidth, yPos, style)
	dc.Line(width-(width/6)-barWidth, yPos, width-(width/6), yPos, style)

	// Hour case 1's
	yPos = yPos + caseWidth + caseWidth/2

	eHour := hour % 5
	for i := 0; i < 4; i++ {
		if i >= eHour {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2DarkSegment, caseWidth, clockBerlin2Case)
		} else {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2Red, caseWidth, clockBerlin2Case)
		}
		xPos := (i * width / 4) + halfCaseWidth
		dc.Roundrect(xPos, yPos, bWidth, quadHeight, height/33, height/33, style)
	}

	style = fmt.Sprintf("fill:none;stroke-width:%d;stroke:%s", caseWidth, clockBerlin2Case)
	dc.Roundrect(halfCaseWidth, yPos, width-caseWidth, quadHeight, radius, radius, style)
	dc.Line(width/4+halfCaseWidth, yPos, width/4+halfCaseWidth, yPos+quadHeight, style)
	dc.Line(width/2+halfCaseWidth, yPos, width/2+halfCaseWidth, yPos+quadHeight, style)
	dc.Line(((3*width)/4)+halfCaseWidth, yPos, ((3*width)/4)+halfCaseWidth, yPos+quadHeight, style)

	yPos = height/6 + yPos
	//barRadius := 0.03 * fheight
	style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", "none", 2*caseWidth, clockBerlin2Case)
	dc.Line(width/6, yPos, width/6+barWidth, yPos, style)
	dc.Line(width-(width/6)-barWidth, yPos, width-(width/6), yPos, style)

	//Minute case 10's
	yPos = yPos + caseWidth + halfCaseWidth
	minutes := timeToRender.Minute()
	tMinute := minutes / 5
	for i := 0; i < 11; i++ {
		if i >= tMinute {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2DarkSegment, caseWidth, clockBerlin2Case)
		} else {
			if (i % 3) == 2 {
				style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2Red, caseWidth, clockBerlin2Case)
			} else {
				style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2Yellow, caseWidth, clockBerlin2Case)
			}
		}
		xPos := ((i * width) / 11.0) + halfCaseWidth
		dc.Roundrect(xPos, yPos, width/11, quadHeight, radius, radius, style)
	}

	style = fmt.Sprintf("fill:none;stroke-width:%d;stroke:%s", caseWidth, clockBerlin2Case)
	dc.Roundrect(halfCaseWidth, yPos, width-caseWidth, quadHeight, radius, radius, style)
	for i := 1; i < 11; i++ {
		xPos := ((i * width) / 11) + halfCaseWidth
		dc.Line(xPos, yPos, xPos, yPos+quadHeight, style)
	}

	yPos = height/6 + yPos
	style = fmt.Sprintf("fill:none;stroke-width:%d;stroke:%s", 2*caseWidth, clockBerlin2Case)
	dc.Line(width/6, yPos, width/6+barWidth, yPos, style)
	dc.Line(width-(width/6)-barWidth, yPos, width-(width/6), yPos, style)

	yPos = yPos + caseWidth + halfCaseWidth
	eMinute := minutes % 5
	for i := 0; i < 4; i++ {
		if i >= eMinute {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2DarkSegment, caseWidth, clockBerlin2Case)
		} else {
			style = fmt.Sprintf("fill:%s;stroke-width:%d;stroke:%s", clockBerlin2Yellow, caseWidth, clockBerlin2Case)
		}
		xPos := ((i * width) / 4) + halfCaseWidth
		dc.Roundrect(xPos, yPos, bWidth, quadHeight, radius, radius, style)
	}

	style = fmt.Sprintf("fill:none;stroke-width:%d;stroke:%s", caseWidth, clockBerlin2Case)
	dc.Roundrect(halfCaseWidth, yPos, width-caseWidth, quadHeight, radius, radius, style)
	dc.Line(width/4+halfCaseWidth, yPos, width/4+halfCaseWidth, yPos+quadHeight, style)
	dc.Line(width/2+halfCaseWidth, yPos, width/2+halfCaseWidth, yPos+quadHeight, style)
	dc.Line(((3*width)/4)+halfCaseWidth, yPos, ((3*width)/4)+halfCaseWidth, yPos+quadHeight, style)

	dc.End()

	return b.Bytes()
}
