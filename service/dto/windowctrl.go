// +build windows
package dto

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/TheTitanrain/w32"
	"github.com/hnakamur/w32syscall"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// WindowCtrlCommandTypeInfo sending key strokes to the active program
var WindowCtrlCommandTypeInfo = models.CommandTypeInfo{"WINDOWCTRL", "WindowCtrl", "controlling windows on the desktop", []models.CommandParameterInfo{
	{"caption", "string", "the caption of the application window", "", make([]string, 0)},
	{"command", "string", "the command to execute on this window. Possible values are: minimize, activate, move  x y", "", make([]string, 0)},
}}

// WindowCtrlCommand is a command to execute a program or batch file.
// Using "command" for getting the command line to execute.
// Using "args" for optional parameters
type WindowCtrlCommand struct {
	Parameters map[string]interface{}
}

// Init the command
func (e *WindowCtrlCommand) Init(a *Action) (bool, error) {
	return true, nil
}

// Stop the command
func (e *WindowCtrlCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (e *WindowCtrlCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {

	value, found := e.Parameters["caption"]
	if !found {
		return false, fmt.Errorf("The caption parameter is missing")
	}
	caption, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("The caption parameter is in wrong format. Please use string as format")
	}
	if caption == "" {
		return false, fmt.Errorf("The caption parameter should not be empty")
	}

	value, found = e.Parameters["command"]
	if !found {
		return false, fmt.Errorf("The command parameter is missing")
	}

	cmdValue, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("The command parameter is in wrong format. Please use string as format")
	}

	if cmdValue == "" {
		return false, fmt.Errorf("The command parameter should not be empty")
	}

	return internalDoWork(caption, cmdValue)
}

func internalDoWork(caption string, command string) (bool, error) {
	command = strings.ToLower(strings.TrimSpace(command))
	err := w32syscall.EnumWindows(func(hwnd syscall.Handle, lparam uintptr) bool {
		h := w32.HWND(hwnd)
		text := w32.GetWindowText(h)
		if strings.Contains(text, caption) {
			clog.Logger.Infof("window \"%s\" found. executing \"%s\"", caption, command)
			switch command {
			case "minimize":
				return w32.ShowWindow(h, syscall.SW_MINIMIZE)
			case "activate":
				w32.ShowWindow(h, syscall.SW_RESTORE)
				return w32.ShowWindow(h, syscall.SW_SHOW)
			}
			if strings.HasPrefix(command, "move") {
				args := strings.Split(command, " ")
				if len(args) != 3 {
					clog.Logger.Error("illegal parameter count for move")
					return false
				}
				x, err := strconv.Atoi(strings.TrimSpace(args[1]))
				if err != nil {
					clog.Logger.Error("illegal x parameter for move")
					return false
				}
				y, err := strconv.Atoi(strings.TrimSpace(args[2]))
				if err != nil {
					clog.Logger.Error("illegal x parameter for move")
					return false
				}
				return w32.SetWindowPos(h, w32.HWND(0), x, y, -1, -1, 0)
			}
		}
		return true
	}, 0)
	if err != nil {
		return false, err
	}
	return true, nil
}
