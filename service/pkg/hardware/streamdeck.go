// Special implementation for openhardwaremonitor app
// +build windows
package hardware

import (
	"fmt"
	"os/exec"
	"strings"

	"wkla.no-ip.biz/remote-desk-service/config"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var StreamdeckIntegInfo = models.IntegInfo{
	Category:    "System",
	Name:        "streamdeck",
	Description: "Elagto Streamdeck (C) Integeration integrates an elgato streamdeck as client for the ReCoS.",
	Image:       "monitor.svg",
	Parameters: []models.ParamInfo{
		{
			Name:           "active",
			Type:           "bool",
			Description:    "activate the streamdeck",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		/*		{
					Name:           "program",
					Type:           "string",
					Description:    "path to the StreamDeckService [Optional]",
					WizardPossible: false,
					List:           make([]string, 0),
				},
		*/
		{
			Name:           "profile",
			Type:           "string",
			Unit:           "",
			Description:    "which profile to use with the streamdeck",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

type StreamDeckInteg struct {
	Program    string
	Parameters []string
	Active     bool
	Profile    string
	cmd        *exec.Cmd
}

var StreamDeckIntegInstance StreamDeckInteg

// InitOpenHardwareMonitor initialise the open hardware monitor connection
func InitStreamDeckInteg(extconfig map[string]interface{}, serviceConfig config.Config) error {
	var err error
	value, ok := extconfig["streamdeck"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Debug("hardware:streamdeck: found config")
			active, ok := config["active"].(bool)
			if !ok {
				active = false
			}
			var program string
			var profile string
			var args []string
			if active {
				clog.Logger.Debug("hardware:streamdeck: active")
				profile, ok = config["profile"].(string)
				if !ok {
					profile = ""
				}
				program, ok = config["program"].(string)
				if !ok {
					program = "./streamdeck"
				}
				if !strings.HasSuffix(program, "/") {
					program = program + "/"
				}
				program = program + "StreamDeckService.exe"
				args = make([]string, 0)
				args = append(args, "-u", fmt.Sprintf("http://127.0.0.1:%d", serviceConfig.Port))

				if profile != "" {
					args = append(args, "-p", profile)
				}
			}
			StreamDeckIntegInstance = StreamDeckInteg{
				Active:     active,
				Profile:    profile,
				Program:    program,
				Parameters: args,
			}
			StreamDeckIntegInstance.Connect()
		}
	}
	return err
}

func DisposeStreamDeck() {
	StreamDeckIntegInstance.Dispose()
}

func (s *StreamDeckInteg) Connect() error {
	if s.Active {
		s.cmd = exec.Command(s.Program, s.Parameters...)

		clog.Logger.Debugf("hardware:streamdeck:execute command line: %v", s.cmd.String())

		err := s.cmd.Start()
		if err != nil {
			clog.Logger.Errorf("error: %v\r\n", err)
		}
	}
	return nil
}

func (s *StreamDeckInteg) Dispose() {
	if s.Active {
		s.cmd.Process.Kill()
		s.cmd.Process.Release()
	}
}
