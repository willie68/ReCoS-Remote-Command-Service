package dto

import (
	"fmt"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ScreenshotCommand is a command to do a sceen shot and store this into the filesystem
type ScreenshotCommand struct {
	Parameters map[string]interface{}
}

// Init a delay in the actual context
func (s *ScreenshotCommand) Init(a *Action) (bool, error) {
	return true, nil
}

// Stop a delay in the actual context
func (s *ScreenshotCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a delay in the actual context
func (s *ScreenshotCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	value, found := s.Parameters["saveto"]
	if !found {
		return false, fmt.Errorf("Folder is missing")
	}
	folder, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("saveto is in wrong format. Please use string as format")
	}

	display := -1
	value, found = s.Parameters["display"]
	if found {
		valueInt, ok := value.(int)
		if !ok {
			return false, fmt.Errorf("saveto is in wrong format. Please use int as format")
		}
		display = valueInt
	}
	clog.Logger.Infof("folder: %s, display: %d ", folder, display)

	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		if (display == i) || (display == -1) {
			bounds := screenshot.GetDisplayBounds(i)

			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				return false, fmt.Errorf("capturing sceenshot missed.%v", err)
			}
			found = true
			filename := ""
			for count := 0; found; count++ {
				filename = fmt.Sprintf("%s/screen_%d_%d.png", folder, count, i)
				_, err := os.Stat(filename)
				if os.IsNotExist(err) {
					found = false
				}
			}
			clog.Logger.Infof("save to: %s ", filename)
			file, _ := os.Create(filename)
			defer file.Close()
			png.Encode(file, img)
		}
	}

	return true, nil
}
