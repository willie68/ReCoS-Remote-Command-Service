package pac

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
	"wkla.no-ip.biz/remote-desk-service/internal/utils"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

// KeysCommandTypeInfo sending key strokes to the active program
var KeysCommandTypeInfo = models.CommandTypeInfo{
	Category:         "System",
	Type:             "KEYS",
	Name:             "Keys",
	Description:      "Typing on a virtual keyboard",
	Icon:             "keystroke_monitoring.png",
	WizardPossible:   true,
	WizardActionType: models.Single,
	Parameters: []models.CommandParameterInfo{
		{
			Name:           "keylayout",
			Type:           "string",
			Description:    "defining the layout of the keyboard used to send the data",
			Unit:           "",
			WizardPossible: false,
			List:           []string{"en", "de"},
		},
		{
			Name:           "keystrokes",
			Type:           "string",
			Description:    "keys are the string with the keys used to send, example: 'akteon{enter}'",
			Unit:           "",
			WizardPossible: true,
			List:           make([]string, 0),
		},
	},
}

// KeysCommand is a command to send keystroke to the active application.
// parameters: layout for the keyboard layout: en or de
// keys: string with the keys to press. Using macros for special keys, see sendMacro methode.
type KeysCommand struct {
	Parameters map[string]interface{}
}

// EnrichType enrich the type info with the informations from the profile
func (p *KeysCommand) EnrichType(profile models.Profile) (models.CommandTypeInfo, error) {
	return KeysCommandTypeInfo, nil
}

// Init the command
func (p *KeysCommand) Init(a *Action, commandName string) (bool, error) {
	return true, nil
}

// Stop the command
func (p *KeysCommand) Stop(a *Action) (bool, error) {
	return true, nil
}

// Execute the command
func (p *KeysCommand) Execute(a *Action, requestMessage models.Message) (bool, error) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		clog.Logger.Errorf("error: %v", err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	layout := "en"
	lvalue, found := p.Parameters["keylayout"]
	if found {
		layoutValue, ok := lvalue.(string)
		if ok {
			layout = layoutValue
		}
	}
	keyboardLayout := utils.KeyEnglish
	switch layout {
	case "de":
		keyboardLayout = utils.KeyGerman
	}

	value, found := p.Parameters["keystrokes"]
	if found {
		keyValue, ok := value.(string)
		if ok {
			clog.Logger.Infof("Key pressed: %s", keyValue)
			sendString(keyValue, kb, keyboardLayout)
		}
	} else {
		return false, fmt.Errorf("the command parameter is missing")
	}
	return true, nil
}

func sendString(value string, kb keybd_event.KeyBonding, layout []utils.KeyItem) error {
	runeValue := []rune(value)
	inMacro := false
	macro := make([]rune, 0)
	for index := 0; index < len(runeValue); index++ {
		char := runeValue[index]
		if char == rune('~') {
			if runeValue[index+1] == rune('~') {
				index++
			} else {
				fmt.Println("\r\nwaiting 1 second")
				time.Sleep(1 * time.Second)
				continue
			}
		}
		if char == rune('{') {
			if runeValue[index+1] == rune('{') {
				index++
			} else {
				fmt.Println("\r\nstart of macro")
				inMacro = true
				macro = make([]rune, 0)
				continue
			}
		}
		if inMacro && (char == rune('}')) {
			fmt.Println("end of macro: " + string(macro))
			sendMacro(string(macro), kb)
			inMacro = false
			continue
		}
		if inMacro {
			macro = append(macro, char)
			continue
		} else {
			fmt.Print(string(char))
		}
		err := sendRune(char, kb, layout)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendMacro(macro string, kb keybd_event.KeyBonding) error {
	kb.Clear()
	switch strings.ToLower(macro) {
	case "backspace", "bs", "bksp":
		kb.AddKey(keybd_event.VK_BACKSPACE)
	case "break":
		kb.AddKey(keybd_event.VK_PAUSE)
	case "capslock":
		kb.AddKey(keybd_event.VK_CAPSLOCK)
	case "del", "delete":
		kb.AddKey(keybd_event.VK_DELETE)
	case "down":
		kb.AddKey(keybd_event.VK_DOWN)
	case "end":
		kb.AddKey(keybd_event.VK_END)
	case "enter":
		kb.AddKey(keybd_event.VK_ENTER)
	case "esc":
		kb.AddKey(keybd_event.VK_ESC)
	case "help":
		kb.AddKey(keybd_event.VK_HELP)
	case "home":
		kb.AddKey(keybd_event.VK_HOME)
	case "ins", "insert":
		kb.AddKey(keybd_event.VK_INSERT)
	case "left":
		kb.AddKey(keybd_event.VK_LEFT)
	case "num":
		kb.AddKey(keybd_event.VK_NUMLOCK)
	case "pgdn":
		kb.AddKey(keybd_event.VK_PAGEDOWN)
	case "pgup":
		kb.AddKey(keybd_event.VK_PAGEUP)
	case "prtsc":
		kb.AddKey(keybd_event.VK_PRINT)
	case "right":
		kb.AddKey(keybd_event.VK_RIGHT)
	case "scrolllock":
		kb.AddKey(keybd_event.VK_SCROLLLOCK)
	case "tab":
		kb.AddKey(keybd_event.VK_TAB)
	case "up":
		kb.AddKey(keybd_event.VK_UP)
	case "f1":
		kb.AddKey(keybd_event.VK_F1)
	case "f2":
		kb.AddKey(keybd_event.VK_F2)
	case "f3":
		kb.AddKey(keybd_event.VK_F3)
	case "f4":
		kb.AddKey(keybd_event.VK_F4)
	case "f5":
		kb.AddKey(keybd_event.VK_F5)
	case "f6":
		kb.AddKey(keybd_event.VK_F6)
	case "f7":
		kb.AddKey(keybd_event.VK_F7)
	case "f8":
		kb.AddKey(keybd_event.VK_F8)
	case "f9":
		kb.AddKey(keybd_event.VK_F9)
	case "f10":
		kb.AddKey(keybd_event.VK_F10)
	case "f11":
		kb.AddKey(keybd_event.VK_F11)
	case "f12":
		kb.AddKey(keybd_event.VK_F12)
	}
	time.Sleep(100 * time.Millisecond)
	err := kb.Launching()
	time.Sleep(100 * time.Millisecond)
	return err
}

func sendRune(char rune, kb keybd_event.KeyBonding, keyboardLayout []utils.KeyItem) error {
	kb.Clear()
	if int(char) < len(keyboardLayout) {

		keyItem := keyboardLayout[char]
		if keyItem.Shift {
			kb.HasSHIFT(true)
		}
		if keyItem.Altgr {
			kb.HasALTGR(true)
		}
		kb.AddKey(keyItem.KeyCode)
		err := kb.Launching()
		return err
	}
	return fmt.Errorf("can't translate this key \"%s\" to keystrokes", string(char))
}
