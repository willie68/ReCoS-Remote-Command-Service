package dto

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"time"

	"github.com/fogleman/gg"
	"golang.org/x/image/bmp"
	"wkla.no-ip.biz/remote-desk-service/api"
	"wkla.no-ip.biz/remote-desk-service/internal"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// StopwatchCommand is a command to simulate a stopwatch.
type StopwatchCommand struct {
	Parameters map[string]interface{}
	action     *Action
	stop       bool
	ticker     *time.Ticker
	done       chan bool
	format     string
	analog     bool

	startTime time.Time
	stopTime  time.Time
	running   bool
}

// Init initialise the stopwatch parmeters
func (s *StopwatchCommand) Init(a *Action) (bool, error) {
	s.action = a
	s.stop = false
	s.format = "%0h:%0m:%0s"
	s.analog = false
	s.done = make(chan bool)

	/*	value, found := c.Parameters["analog"]
		if found {
			var ok bool
			c.analog, ok = value.(bool)
			if !ok {
				return false, fmt.Errorf("Analog is in wrong format. Please use boolean as format")
			}
		}
	*/
	value, found := s.Parameters["format"]
	if found {
		var ok bool
		s.format, ok = value.(string)
		if !ok {
			return false, fmt.Errorf("Format is in wrong format. Please use string as format")
		}
	}
	return true, nil
}

// Stop stops the actual command
func (s *StopwatchCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

func (s *StopwatchCommand) startStopwatch() {
	s.running = true
	s.ticker = time.NewTicker(1 * time.Millisecond)
	go func() {
		for {
			select {
			case <-s.done:
				return
			case t := <-s.ticker.C:
				if api.HasConnectionWithProfile(s.action.Profile) {
					tdelta := t.Sub(s.startTime)
					title, _ := durationfmt.Format(tdelta, s.format)
					if s.analog {
						s.SendPNG(title)
					} else {
						message := models.Message{
							Profile: s.action.Profile,
							Action:  s.action.Name,
							State:   1,
							Title:   title,
						}
						api.SendMessage(message)
					}
				}
			}
		}
	}()
}

func (s *StopwatchCommand) stopStopwatch() {
	s.ticker.Stop()
	s.done <- true
	s.running = false
}

// Execute a delay in the actual context
func (s *StopwatchCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	timeNow := time.Now()
	if s.running {
		// stop the running clock
		s.stopTime = timeNow
		s.stopStopwatch()

		go func() {
			tdelta := s.stopTime.Sub(s.startTime)
			title, _ := durationfmt.Format(tdelta, s.format)
			if s.analog {
				s.SendPNG(title)
			} else {
				message := models.Message{
					Profile: s.action.Profile,
					Action:  s.action.Name,
					State:   1,
					Title:   s.action.Config.Title,
					Text:    title,
				}
				api.SendMessage(message)
			}
			time.Sleep(3 * time.Second)
		}()
	} else {
		// start the stopwatch
		s.startTime = timeNow
		s.startStopwatch()
	}
	return false, nil
}

// SendPNG sending this array to the client
func (s *StopwatchCommand) SendPNG(value string) {
	now := time.Now()
	dc := gg.NewContext(clockImageWidth, clockImageHeight)
	halfWidth := float64(clockImageHeight / 2)
	halfHeight := float64(clockImageHeight / 2)
	dc.SetColor(colorTicks)
	dc.InvertY()

	dc.SetLineWidth(1.0)
	dc.DrawCircle(halfWidth, halfHeight, halfHeight-1)
	dc.MoveTo(halfWidth-1, clockImageHeight)
	dc.LineTo(halfWidth-1, clockImageHeight-tickLength)

	for i := 0; i < 12; i++ {
		deg := float64(30.0 * i)
		dc.MoveTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-tickLength)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-tickLength)))
		dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-1)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-1)))
	}
	dc.Stroke()

	dc.SetLineWidth(4.0)

	dc.MoveTo(halfWidth-1, clockImageHeight)
	dc.LineTo(halfWidth-1, clockImageHeight-tickLength)
	dc.MoveTo(halfWidth-1, 0)
	dc.LineTo(halfWidth-1, tickLength)

	dc.MoveTo(0, halfHeight-1)
	dc.LineTo(tickLength, halfHeight-1)
	dc.MoveTo(clockImageWidth-tickLength, halfHeight-1)
	dc.LineTo(clockImageWidth, halfHeight-1)

	dc.Stroke()

	dc.SetColor(colorArSecond)
	dc.SetLineWidth(1.0)
	seconds := now.Second()

	deg := float64(6.0 * seconds)
	dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*10), halfHeight-(math.Cos(deg2Rad(deg))*10))
	dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-2)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-2)))

	dc.Stroke()

	dc.SetColor(colorArMinute)
	dc.SetLineWidth(3.0)
	minute := now.Minute()

	deg = float64(6.0 * minute)
	dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*2), halfHeight-(math.Cos(deg2Rad(deg))*2))
	dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth-10)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight-10)))

	dc.Stroke()

	dc.SetColor(colorArHour)
	dc.SetLineWidth(6.0)
	hour := now.Hour()

	deg = float64(30.0*hour + (minute / 2))
	dc.MoveTo(halfWidth-(math.Sin(deg2Rad(deg))*2), halfHeight-(math.Cos(deg2Rad(deg))*2))
	dc.LineTo(halfWidth+(math.Sin(deg2Rad(deg))*(halfWidth*1/2)), halfHeight+(math.Cos(deg2Rad(deg))*(halfHeight*1/2)))

	dc.Stroke()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	bmp.Encode(&buff, myImage)

	// Encode the bytes in the buffer to a base64 string
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	// You can embed it in an html doc with this string
	image := "data:image/bmp;base64," + encodedString

	//image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEgAAABICAYAAABV7bNHAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAABTBJREFUeNrsmz1s01AQx5+Tpl+0JUNZkIAMSDQLFLGwIMLEhJSyIYaWgQ4sbYWYAQmBkJCAhYUBOrAw0CAmGFAQDCygIIYGJKTQBSQ6hEJTmk/unPfSF5O2sX12/ML7S69xktZ2f+/u3p19ZkxLS0tLS0tVGZ08+KmJ04ktvs4/W3iS+W8AAYwYvCRhHIcxDiPW5p9m+HgFIwXQ8l0DCKBE4WUKxiSHQqEUjKcA6qGygDiYWRgzMKLyd4OGwfb19LCxSISNhkNsVyjUch9fyxW2XK2wJXhdLJVa/UoOxjyMO15aleEBnCtWMAjlWH8fO9Lba4JxonfFInsP4/Wf9Vag5gBSKtCAAAy60APZldBCJgYH2bG+PrITLtRq7PnaHxhr5rakNIwJamsyiOCgO932EkyboPIcUjowgADOAx6ITZ0c6DfhoFv5oeVqld3/9dsap85RBfEwFRwEcmFkGAANsIjhX3ol4hvOdXYDUvJAPB79nF183jFAMhx0qUsjIyzuMABTCI+9F1bHj8US45iOAiQDIKV9ByTDwZO6HN0JkMKs09oN53CwN8LerhcFpARA+gqQMr7FIIAzxVcr03KuRaO+xZt2tVQus+s/V+Tgfdhp2RJyuJSbUGaHhwMHR1j12aEd8kcLPHn1FpCAgzo/PGSeSFCFKQauplwxOQ3xBBDPkMfFUo5ZcdA1MTggT+LUNlcPnAPi5jkjJ4GqaLrZ1S57ZUGzorbyMwmkikf1PKm+qsFkJ0kBWa3Hy/LBK53d0WRFM9QWlJStR0VtZNsNKxqnBDTTOIiC1rMRsJsmd5IEEL9MatKWZkBJjYZC8oqWpLKgxrKowrLeTm4k8qJ23Ww7QMfFxlgHC1G6gran5eS7ATQuKuVukCXzP0QGaKxLAFkmO+YKEA/QTOQ/3aK9PWEyF2sA2uzWjJo5kb3/xVU5bkBuEd6zJ5Agyp8+kezHFSCEM3TxYiAB5aenW+dDNsNF9/hOmypUa/5ZEJrxZjMVWEC1mrYgSm0FqHGRGxsJukXYECHClCtA8j1uaafKS5rsDIWLpfHHYqncNYDwlhAloIzYqd3gFkRlm+/ff6AA9EpsbNLEpJSwx8jqHSQuhnrzb+OSyoAyEGNzrgHxQJ0SO1fZzfD8lytV8XaeMg9q7AwbllTVi+ZzT5EB4r1/uTqgNSWtCIOzFEMftutedjLpqyJNV9GKFgpr8tu75KUGb2fL1Q9WMNveVIo9FuvJkAPimhMb2BOoSmEqnWteeIIdtd0W9jm7mD0Qj+M16rG6BRmBv5h/D+AsVRqlxU0nvdR2q/lzoshDV3u9Htzc6NHqqpz3pAHOFepqfrO86ETjJH6vyrVNYIQTJy0mZu+0033Z7rwEV/uOjZGwmcTQhw2T2Di5MyAX9hGOJe6csLOsuwbEIWVkSC+hDMEu130dbslDt3q8WpA/OuO2695x764MCd+/56XIwQ7cw8fj3lr5ZVqzZDlnKB5wcdXcbIX0BeLRu2KJ7QdL8svlMBDf+LnCvlUqVrdKU+yf6mEWXP4XmHSzEdtlsCdn1CNQWD5ghmy5DINQSJ/4oXwcCrvQsElyVv4cQZ3s7ydrGUaLwcLTAgaB3HW6lPsCSAKV4KAS8ucICHuMjvRGbMNCKGgxlksWcmU+52al8hXQdqCE4vyRzNFQ6zCI+dUPyNi3yLOwPpynfDbMV0CW+DTJA3nM5e6w0MTrUymvLMZ3QC1goUUd4rASW/x6nm08Ev6Blws5pqWlpaWlpaWlpRUE/RVgAAD3F9MyT2oUAAAAAElFTkSuQmCC"
	message := models.Message{
		Profile:  s.action.Profile,
		Action:   s.action.Name,
		ImageURL: image,
		Text:     value,
		State:    0,
	}
	api.SendMessage(message)
}