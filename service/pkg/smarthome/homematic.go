package smarthome

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html/charset"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var HomematicIntegInfo = models.IntegInfo{
	Category:    "Smarthome",
	Name:        "homematic",
	Description: "homematic is a smart home system.",
	Image:       "home.svg",
	Parameters: []models.ParamInfo{
		{
			Name:           "active",
			Type:           "bool",
			Description:    "activate the homematic integration",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "url",
			Type:           "string",
			Description:    "url to the homematic server",
			WizardPossible: false,
			List:           make([]string, 0),
		},
		{
			Name:           "updateperiod",
			Type:           "int",
			Unit:           " sec",
			Description:    "update period in seconds to update the used devices",
			WizardPossible: false,
			List:           make([]string, 0),
		},
	},
}

type Homematic struct {
	url     string
	periode int
	reload  sync.Mutex
	ticker  *time.Ticker
	done    chan bool
}

type Devicelist struct {
	XMLName xml.Name `xml:"deviceList"`
	Devices []Device `xml:"device"`
}

type Statelist struct {
	XMLName xml.Name `xml:"stateList"`
	Devices []Device `xml:"device"`
}

type State struct {
	XMLName xml.Name `xml:"state"`
	Devices []Device `xml:"device"`
}

type Device struct {
	XMLName    xml.Name  `xml:"device"`
	Name       string    `xml:"name,attr"`
	Address    string    `xml:"address,attr"`
	Ise_id     string    `xml:"ise_id,attr"`
	Interface  string    `xml:"interface,attr"`
	DeviceType string    `xml:"device_type,attr"`
	ReadConfig bool      `xml:"read_config,attr"`
	Channels   []Channel `xml:"channel"`
}

type Channel struct {
	XMLName          xml.Name    `xml:"channel"`
	Name             string      `xml:"name,attr"`
	ChannelType      string      `xml:"type,attr"`
	Address          string      `xml:"address,attr"`
	Ise_id           string      `xml:"ise_id,attr"`
	Direction        string      `xml:"direction,attr"`
	ParentDevice     string      `xml:"parent_device,attr"`
	Index            int         `xml:"index,attr"`
	GroupPartner     string      `xml:"group_partner,attr"`
	AesAvailable     bool        `xml:"aes_available,attr"`
	TransmissionMode string      `xml:"transmission_mode,attr"`
	Visible          bool        `xml:"visible,attr"`
	ReadConfig       bool        `xml:"read_config,attr"`
	Operate          bool        `xml:"operate,attr"`
	Datapoints       []Datapoint `xml:"datapoint"`
}

type Datapoint struct {
	XMLName xml.Name `xml:"datapoint"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Ise_id  string   `xml:"ise_id,attr"`
	Value   string   `xml:"value,attr"`
	//ValueType 2 is bool, 4 is float32, 8 ?, 16 int, 20 string
	ValueType  int    `xml:"valuetype,attr"`
	ValueUnit  string `xml:"valueunit,attr"`
	TimeStamp  int    `xml:"timestamp,attr"`
	Operations int    `xml:"operations,attr"`
}

type ProgramList struct {
	XMLName  xml.Name  `xml:"programList"`
	Programs []Program `xml:"program"`
}

type Program struct {
	XMLName     xml.Name `xml:"program"`
	ID          string   `xml:"id,attr"`
	Active      bool     `xml:"active,attr"`
	TimeStamp   int      `xml:"timestamp,attr"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description,attr"`
	Visible     bool     `xml:"visible,attr"`
	Operate     bool     `xml:"operate,attr"`
}

type ProgramResult struct {
	XMLName xml.Name `xml:"result"`
	Started []struct {
		ProgramID string `xml:"program_id,attr"`
	} `xml:"started"`
}

type StateResult struct {
	XMLName xml.Name `xml:"result"`
	Changed []struct {
		ID       string  `xml:"id,attr"`
		NewValue float64 `xml:"new_value,attr"`
	} `xml:"changed"`
}

var homematic *Homematic

func GetHomematic() (*Homematic, bool) {
	if homematic != nil {
		return homematic, true
	}
	return nil, false
}

func InitHomematic(extconfig map[string]interface{}) error {
	value, ok := extconfig["homematic"]
	if ok {
		config := value.(map[string]interface{})
		if config != nil {
			clog.Logger.Debug("smarthome:homematic: found config")
			active, ok := config["active"].(bool)
			if !ok {
				active = false
			}
			if active {
				clog.Logger.Debug("smarthome:homematic: active")
				url, ok := config["url"].(string)
				if !ok {
					return fmt.Errorf("homematic: no ipaddress given")
				}
				updatePeriod, ok := config["updateperiod"].(int)
				if !ok {
					updatePeriod = 10
				}
				homematic = &Homematic{
					url:     url,
					periode: updatePeriod,
					ticker:  time.NewTicker(time.Duration(updatePeriod) * time.Second),
					done:    make(chan bool),
				}
				err := homematic.init()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (h *Homematic) init() error {
	go func() {
		h.reloadAll()
	}()
	go func() {
		for {
			select {
			case <-h.done:
				return
			case <-h.ticker.C:
				h.reloadAll()
			}
		}
	}()

	return nil
}

func (h *Homematic) UpdatePeriod() int {
	return h.periode
}

func (h *Homematic) reloadAll() error {
	h.reload.Lock()
	defer h.reload.Unlock()
	//clog.Logger.Info("homematic: realod all")

	/*
		result, err := h.ChangeState("2134", 0)
		if err != nil {
			return fmt.Errorf("GET error: %v", err)
		}
		jsonStr, _ := json.Marshal(result)
		clog.Logger.Infof("Result: %s", string(jsonStr))
			result, err = h.RunProgram("2514")
			if err != nil {
				return fmt.Errorf("GET error: %v", err)
			}
			jsonStr, _ = json.Marshal(result)
			clog.Logger.Infof("Result: %s", string(jsonStr))
	*/
	return nil
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("read body: %v", err)
	}

	return data, nil
}

func (h *Homematic) DeviceList() (*Devicelist, error) {
	deviceUrl := h.url + "/addons/xmlapi/devicelist.cgi"

	if xmlBytes, err := getXML(deviceUrl); err != nil {
		clog.Logger.Errorf("error evaluating devices: %v", err)
		return nil, err
	} else {
		reader := bytes.NewReader(xmlBytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		//		var devicelist devicelist
		var devicelist Devicelist
		err = decoder.Decode(&devicelist)
		if err != nil {
			clog.Logger.Errorf("error evaluating devices: %v", err)
			return nil, err
		}
		return &devicelist, nil
	}
}

func (h *Homematic) StateList() (*Statelist, error) {
	stateUrl := h.url + "/addons/xmlapi/statelist.cgi"

	if xmlBytes, err := getXML(stateUrl); err != nil {
		clog.Logger.Errorf("error evaluating states: %v", err)
		return nil, err
	} else {
		reader := bytes.NewReader(xmlBytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		//		var devicelist devicelist
		var statelist Statelist
		err = decoder.Decode(&statelist)
		if err != nil {
			clog.Logger.Errorf("error evaluating states: %v", err)
			return nil, err
		}
		return &statelist, nil
	}
}

func (h *Homematic) ProgramList() (*ProgramList, error) {
	programsUrl := h.url + "/addons/xmlapi/programlist.cgi"

	if xmlBytes, err := getXML(programsUrl); err != nil {
		clog.Logger.Errorf("error evaluating programs: %v", err)
		return nil, err
	} else {
		reader := bytes.NewReader(xmlBytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		//		var devicelist devicelist
		var prglist ProgramList
		err = decoder.Decode(&prglist)
		if err != nil {
			clog.Logger.Errorf("error evaluating states: %v", err)
			return nil, err
		}
		return &prglist, nil
	}
}

func (h *Homematic) RunProgram(programID string) (*ProgramResult, error) {
	programUrl := h.url + "/addons/xmlapi/runprogram.cgi?program_id=" + programID

	if xmlBytes, err := getXML(programUrl); err != nil {
		clog.Logger.Errorf("error runnong program: %v", err)
		return nil, err
	} else {
		reader := bytes.NewReader(xmlBytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		//		var devicelist devicelist
		var result ProgramResult
		err = decoder.Decode(&result)
		if err != nil {
			clog.Logger.Errorf("error evaluating states: %v", err)
			return nil, err
		}
		return &result, nil
	}
}

func (h *Homematic) ChangeState(IseID string, value float64) (*StateResult, error) {
	stateUrl := fmt.Sprintf("%s/addons/xmlapi/statechange.cgi?ise_id=%s&new_value=%.2f", h.url, IseID, value)

	if xmlBytes, err := getXML(stateUrl); err != nil {
		clog.Logger.Errorf("error runnong program: %v", err)
		return nil, err
	} else {
		reader := bytes.NewReader(xmlBytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		//		var devicelist devicelist
		var result StateResult
		err = decoder.Decode(&result)
		if err != nil {
			clog.Logger.Errorf("error evaluating states: %v", err)
			return nil, err
		}
		return &result, nil
	}
}

func (h *Homematic) State(IseID string) ([]Datapoint, error) {
	stateUrl := fmt.Sprintf("%s/addons/xmlapi/state.cgi?channel_id=%s", h.url, IseID)

	if xmlBytes, err := getXML(stateUrl); err != nil {
		clog.Logger.Errorf("error runnong program: %v", err)
		return nil, err
	} else {
		reader := bytes.NewReader(xmlBytes)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		//		var devicelist devicelist
		var result State
		err = decoder.Decode(&result)
		if err != nil {
			clog.Logger.Errorf("error evaluating states: %v", err)
			return nil, err
		}
		for _, device := range result.Devices {
			for _, channel := range device.Channels {
				if channel.Ise_id == IseID {
					return channel.Datapoints, nil
				}
			}
		}
		return nil, fmt.Errorf("channnel with id %s not found", IseID)
	}
}
