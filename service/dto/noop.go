package dto

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// NoopCommand is a command to do nothing.
type NoopCommand struct {
	Parameters map[string]interface{}
	action     *Action
	stop       bool
	ticker     *time.Ticker
	done       chan bool
	temps      []float64
}

// Init nothing
func (d *NoopCommand) Init(a *Action) (bool, error) {
	clog.Logger.Info("initialising the clock")
	d.temps = make([]float64, 72)
	d.action = a
	d.stop = false
	d.ticker = time.NewTicker(1 * time.Second)
	d.done = make(chan bool)
	go func() {
		for {
			select {
			case <-d.done:
				return
			case <-d.ticker.C:
				temp := rand.Float64() * 100
				d.temps = append(d.temps, temp)
				if len(d.temps) > 72 {
					copy(d.temps, d.temps[1:])
				}
				d.SendPNG()
			}
		}
	}()
	return true, nil
}

// Stop nothing
func (d *NoopCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute nothing
func (d *NoopCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	return true, nil
}

// SendPNG sending this array to the client
func (d *NoopCommand) SendPNG() {
	dc := gg.NewContext(72, 72)
	dc.SetRGB(1, 0, 0)
	//dc.DrawRectangle(0, 0, 72, 72)
	dc.InvertY()
	dc.MoveTo(0, 0)
	xLast := 0.0
	for index, temp := range d.temps {
		var x float64
		var y float64
		x = float64(index)
		y = (72.0 / 70.0) * temp
		dc.LineTo(x, y)
		xLast = x
	}
	dc.LineTo(xLast, 0)
	dc.LineTo(0, 0)
	dc.Fill()
	dc.Stroke()

	myImage := dc.Image()
	var buff bytes.Buffer

	// The Buffer satisfies the Writer interface so we can use it with Encode
	// In previous example we encoded to a file, this time to a temp buffer
	png.Encode(&buff, myImage)

	// Encode the bytes in the buffer to a base64 string
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())

	// You can embed it in an html doc with this string
	image := "data:image/png;base64," + encodedString

	//image := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEgAAABICAYAAABV7bNHAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAABTBJREFUeNrsmz1s01AQx5+Tpl+0JUNZkIAMSDQLFLGwIMLEhJSyIYaWgQ4sbYWYAQmBkJCAhYUBOrAw0CAmGFAQDCygIIYGJKTQBSQ6hEJTmk/unPfSF5O2sX12/ML7S69xktZ2f+/u3p19ZkxLS0tLS0tVGZ08+KmJ04ktvs4/W3iS+W8AAYwYvCRhHIcxDiPW5p9m+HgFIwXQ8l0DCKBE4WUKxiSHQqEUjKcA6qGygDiYWRgzMKLyd4OGwfb19LCxSISNhkNsVyjUch9fyxW2XK2wJXhdLJVa/UoOxjyMO15aleEBnCtWMAjlWH8fO9Lba4JxonfFInsP4/Wf9Vag5gBSKtCAAAy60APZldBCJgYH2bG+PrITLtRq7PnaHxhr5rakNIwJamsyiOCgO932EkyboPIcUjowgADOAx6ITZ0c6DfhoFv5oeVqld3/9dsap85RBfEwFRwEcmFkGAANsIjhX3ol4hvOdXYDUvJAPB79nF183jFAMhx0qUsjIyzuMABTCI+9F1bHj8US45iOAiQDIKV9ByTDwZO6HN0JkMKs09oN53CwN8LerhcFpARA+gqQMr7FIIAzxVcr03KuRaO+xZt2tVQus+s/V+Tgfdhp2RJyuJSbUGaHhwMHR1j12aEd8kcLPHn1FpCAgzo/PGSeSFCFKQauplwxOQ3xBBDPkMfFUo5ZcdA1MTggT+LUNlcPnAPi5jkjJ4GqaLrZ1S57ZUGzorbyMwmkikf1PKm+qsFkJ0kBWa3Hy/LBK53d0WRFM9QWlJStR0VtZNsNKxqnBDTTOIiC1rMRsJsmd5IEEL9MatKWZkBJjYZC8oqWpLKgxrKowrLeTm4k8qJ23Ww7QMfFxlgHC1G6gran5eS7ATQuKuVukCXzP0QGaKxLAFkmO+YKEA/QTOQ/3aK9PWEyF2sA2uzWjJo5kb3/xVU5bkBuEd6zJ5Agyp8+kezHFSCEM3TxYiAB5aenW+dDNsNF9/hOmypUa/5ZEJrxZjMVWEC1mrYgSm0FqHGRGxsJukXYECHClCtA8j1uaafKS5rsDIWLpfHHYqncNYDwlhAloIzYqd3gFkRlm+/ff6AA9EpsbNLEpJSwx8jqHSQuhnrzb+OSyoAyEGNzrgHxQJ0SO1fZzfD8lytV8XaeMg9q7AwbllTVi+ZzT5EB4r1/uTqgNSWtCIOzFEMftutedjLpqyJNV9GKFgpr8tu75KUGb2fL1Q9WMNveVIo9FuvJkAPimhMb2BOoSmEqnWteeIIdtd0W9jm7mD0Qj+M16rG6BRmBv5h/D+AsVRqlxU0nvdR2q/lzoshDV3u9Htzc6NHqqpz3pAHOFepqfrO86ETjJH6vyrVNYIQTJy0mZu+0033Z7rwEV/uOjZGwmcTQhw2T2Di5MyAX9hGOJe6csLOsuwbEIWVkSC+hDMEu130dbslDt3q8WpA/OuO2695x764MCd+/56XIwQ7cw8fj3lr5ZVqzZDlnKB5wcdXcbIX0BeLRu2KJ7QdL8svlMBDf+LnCvlUqVrdKU+yf6mEWXP4XmHSzEdtlsCdn1CNQWD5ghmy5DINQSJ/4oXwcCrvQsElyVv4cQZ3s7ydrGUaLwcLTAgaB3HW6lPsCSAKV4KAS8ucICHuMjvRGbMNCKGgxlksWcmU+52al8hXQdqCE4vyRzNFQ6zCI+dUPyNi3yLOwPpynfDbMV0CW+DTJA3nM5e6w0MTrUymvLMZ3QC1goUUd4rASW/x6nm08Ev6Blws5pqWlpaWlpaWlpRUE/RVgAAD3F9MyT2oUAAAAAElFTkSuQmCC"
	message := models.Message{
		Profile:  d.action.Profile,
		Action:   d.action.Name,
		ImageURL: image,
		State:    0,
	}
	api.SendMessage(message)
}
