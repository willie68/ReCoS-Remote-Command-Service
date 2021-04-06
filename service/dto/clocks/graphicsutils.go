package clocks

import (
	"image/color"

	"github.com/fogleman/gg"
)

type num2Seg struct {
	a bool
	b bool
	c bool
	d bool
	e bool
	f bool
	g bool
}

var num2SegArray = []num2Seg{
	{true, true, true, true, true, true, false},     //0
	{false, true, true, false, false, false, false}, //1
	{true, true, false, true, true, false, true},    //2
	{true, true, true, true, false, false, true},    //3
	{false, true, true, false, false, true, true},   //4
	{true, false, true, true, false, true, true},    //5
	{true, false, true, true, true, true, true},     //6
	{true, true, true, false, false, false, false},  //7
	{true, true, true, true, true, true, true},      //8
	{true, true, true, true, false, true, true},     //9
}

const ClockImageWidth = 200
const ClockImageHeight = 200

var (
	colorDarkSegment color.Color = color.RGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
	colorSegments    color.Color = color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	colorBlack       color.Color = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

func WriteSegment(digit int, xStart, yStart, width, height, thickness float64, dc *gg.Context, drawColor color.Color, showDarkSegments bool) {
	//clog.Logger.Debugf("Counter: start segment at :%.0f %.0f %.0f %.0f %.0f", xStart, yStart, width, height, thickness)
	/*	dc.SetLineWidth(1)
		dc.SetColor(color.White)
		dc.MoveTo(xStart, yStart)
		dc.LineTo(xStart, yStart+height)
		dc.LineTo(xStart+width, yStart+height)
		dc.LineTo(xStart+width, yStart)
		dc.LineTo(xStart, yStart)
		dc.Fill()
		dc.Stroke()
		dc.SetColor(color.Black)
		dc.DrawLine(xStart, yStart+height/2, xStart+width, yStart+height/2)
		dc.Stroke()
	*/
	halfThickness := thickness / 2
	segmentHeight := height/2 - halfThickness
	segmentWidth := width - 2*thickness
	myNum2Seg := num2SegArray[digit]

	if myNum2Seg.a {
		// A
		writeHorizontalSegment(xStart+thickness*3/2, yStart+height-thickness, segmentWidth, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeHorizontalSegment(xStart+thickness*3/2, yStart+height-thickness, segmentWidth, thickness, dc, colorDarkSegment)
		}
	}
	if myNum2Seg.b {
		// B
		writeVerticalSegment(xStart+width-thickness*3/2, yStart+height/2, segmentHeight, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeVerticalSegment(xStart+width-thickness*3/2, yStart+height/2, segmentHeight, thickness, dc, colorDarkSegment)
		}
	}
	if myNum2Seg.c {
		// C
		writeVerticalSegment(xStart+width-2*thickness, yStart+halfThickness, segmentHeight, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeVerticalSegment(xStart+width-2*thickness, yStart+halfThickness, segmentHeight, thickness, dc, colorDarkSegment)
		}
	}
	if myNum2Seg.d {
		// D
		writeHorizontalSegment(xStart+halfThickness, yStart, segmentWidth, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeHorizontalSegment(xStart+halfThickness, yStart, segmentWidth, thickness, dc, colorDarkSegment)
		}
	}
	if myNum2Seg.e {
		// E
		writeVerticalSegment(xStart, yStart+halfThickness, segmentHeight, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeVerticalSegment(xStart, yStart+halfThickness, segmentHeight, thickness, dc, colorDarkSegment)
		}
	}
	if myNum2Seg.f {
		// F
		writeVerticalSegment(xStart+halfThickness, yStart+height/2, segmentHeight, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeVerticalSegment(xStart+halfThickness, yStart+height/2, segmentHeight, thickness, dc, colorDarkSegment)
		}
	}
	if myNum2Seg.g {
		// G
		writeHorizontalSegment(xStart+thickness, yStart+(height/2-halfThickness), segmentWidth, thickness, dc, drawColor)
	} else {
		if showDarkSegments {
			writeHorizontalSegment(xStart+thickness, yStart+(height/2-halfThickness), segmentWidth, thickness, dc, colorDarkSegment)
		}
	}
}

func writeHorizontalSegment(xStart float64, yStart float64, width float64, thickness float64, dc *gg.Context, color color.Color) {
	//clog.Logger.Debugf("Counter: start horizontal segment at :%.0f %.0f width %.0f  %.0f", xStart, yStart, width, thickness)
	xStartInt := xStart + 2
	halfThickness := thickness / 2
	yEnd := yStart + thickness
	xEnd := xStart + width - 2
	dc.SetColor(colorBlack)
	dc.MoveTo(xStartInt, yStart+halfThickness)
	dc.LineTo(xStartInt+halfThickness, yEnd)
	dc.LineTo(xEnd-halfThickness, yEnd)
	dc.LineTo(xEnd, yStart+halfThickness)
	dc.LineTo(xEnd-halfThickness, yStart)
	dc.LineTo(xStartInt+halfThickness, yStart)
	dc.LineTo(xStartInt, yStart+halfThickness)
	dc.SetColor(color)
	dc.Fill()
	dc.Stroke()
}

func writeVerticalSegment(xStart float64, yStart float64, height float64, thickness float64, dc *gg.Context, color color.Color) {
	//clog.Logger.Debugf("Counter: start vertical segment at : %.0f %.0f / height %.0f %.0f", xStart, yStart, height, thickness)
	halfThickness := thickness / 2
	yStartInt := yStart + 2
	yEnd := yStart + height - 2
	xEnd := xStart + thickness + halfThickness
	dc.SetColor(colorBlack)
	dc.MoveTo(xStart+halfThickness+1, yStartInt)
	dc.LineTo(xStart, yStartInt+halfThickness)
	dc.LineTo(xEnd-thickness, yEnd-halfThickness)
	dc.LineTo(xEnd-halfThickness-1, yEnd)
	dc.LineTo(xEnd, yEnd-halfThickness)
	dc.LineTo(xStart+thickness, yStartInt+halfThickness)
	dc.LineTo(xStart+halfThickness+1, yStartInt)
	dc.SetColor(color)
	dc.Fill()
	dc.Stroke()
}
