package pac

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/kbinani/screenshot"
	"wkla.no-ip.biz/remote-desk-service/api"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// ScreenshotCommandTypeInfo saving to the file system
var ScreenshotCommandTypeInfo = models.CommandTypeInfo{
	Category:         "useful",
	Type:             "SCREENSHOT",
	Name:             "Screenshot",
	Description:      "Taking a Screenshot",
	Icon:             "camera.svg",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.ParamInfo{
		{
			Name:           "saveto",
			Type:           "string",
			Description:    "the folder where the screenshot should be saved",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
		{
			Name:           "display",
			Type:           "int",
			Description:    "the display number",
			Unit:           "",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

// ScreenshotCommand is a command to do a sceen shot and store this into the filesystem
type ScreenshotCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (s *ScreenshotCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return ScreenshotCommandTypeInfo, nil
}

// Init a delay in the actual context
func (s *ScreenshotCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop a delay in the actual context
func (s *ScreenshotCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute a delay in the actual context
func (s *ScreenshotCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {

	folder, err := ConvertParameter2String(s.Parameters, "saveto", "")
	if err != nil {
		clog.Logger.Errorf("screenshot: error in getting saveto folder: %v", err)
		return false, err
	}

	if folder == "" {
		return false, fmt.Errorf("saveto folder should not be empty")
	}

	display, err := ConvertParameter2Int(s.Parameters, "display", -1)
	if err != nil {
		clog.Logger.Errorf("screenshot: error in getting display: %v", err)
		return false, err
	}
	clog.Logger.Infof("folder: %s, display: %d ", folder, display)

	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		names := ""
		if (display == i) || (display == -1) {
			bounds := screenshot.GetDisplayBounds(i)

			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				return false, fmt.Errorf("capturing sceenshot missed.%v", err)
			}
			found := true
			filename := ""
			for count := 0; found; count++ {
				filename = fmt.Sprintf("%s/screen_%d_%d.svg", folder, count, i)
				_, err := os.Stat(filename)
				if os.IsNotExist(err) {
					found = false
				}
			}
			clog.Logger.Infof("save to: %s ", filename)
			file, _ := os.Create(filename)
			defer file.Close()
			png.Encode(file, img)
			names = names + filename + "\r\n"
		}
		message := models.Message{
			Profile:  a.Profile,
			Action:   a.Name,
			ImageURL: "check_mark.svg",
			Title:    "done",
			Text:     names,
			State:    0,
		}
		api.SendMessage(message)
		go func() {
			time.Sleep(3 * time.Second)
			message := models.Message{
				Profile:  a.Profile,
				Action:   a.Name,
				ImageURL: a.Config.Icon,
				Title:    a.Config.Title,
				Text:     "",
				State:    0,
			}
			api.SendMessage(message)
		}()
	}

	return false, nil
}
