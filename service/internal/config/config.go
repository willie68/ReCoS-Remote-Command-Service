package config

import (
	"fmt"
	"os"
	"strings"
)

// Config our service configuration
type Config struct {
	//port of the http server
	Port int `yaml:"port"`
	//port of the https server
	Sslport int `yaml:"sslport"`
	//this is the url how to connect to this service from outside
	ServiceURL string `yaml:"serviceURL"`

	SecretFile string `yaml:"secretfile"`

	Password string `yaml:"password"`

	AppID   int    `yaml:"appid"`
	AppUUID string `yaml:"appuuid"`

	HealthCheck HealthCheck `yaml:"healthcheck"`

	Profiles     string `yaml:"profiles"`
	Sessions     string `yaml:"sessions"`
	TimezoneInfo string `yaml:"timezonezip"`

	ExternalConfig map[string]interface{} `yaml:"extconfig"`

	Logging LoggingConfig `yaml:"logging"`
}

// HealthCheck configuration for the health check system
type HealthCheck struct {
	Period int `yaml:"period"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`
}

var DefaultConfig = Config{
	Port:       9280,
	Sslport:    0,
	ServiceURL: "http://127.0.0.1:9280",
	SecretFile: "",
	Password:   "recosadmin",
	AppID:      73,
	AppUUID:    "",
	HealthCheck: HealthCheck{
		Period: 30,
	},
	Profiles:     "${configdir}/profiles",
	Sessions:     "${configdir}/sessions",
	TimezoneInfo: "${configdir}/zoneinfo.zip",
	ExternalConfig: map[string]interface{}{
		"openhardwaremonitor": map[string]interface{}{
			"active":       false,
			"url":          "http://127.0.0.1:12999/data.json",
			"updateperiod": 5,
		},
		"audioplayer": map[string]interface{}{
			"active":     false,
			"samplerate": 48000,
		},
		"philipshue": map[string]interface{}{
			"active":       false,
			"username":     "<the bridge generated username here>",
			"device":       "recos#hue_user",
			"ipaddress":    "127.0.0.1",
			"updateperiod": 5,
		},
		"homematic": map[string]interface{}{
			"active":       false,
			"url":          "http://192.168.172.10",
			"updateperiod": 5,
		},
		"streamdeck": map[string]interface{}{
			"active":  false,
			"program": "",
			"profile": "",
		},
		"obs": map[string]interface{}{
			"active":   false,
			"host":     "127.0.0.1",
			"port":     4444,
			"password": "",
		},
	},
	Logging: LoggingConfig{
		Level:    "INFO",
		Filename: "${configdir}/logging.log",
	},
}

// GetDefaultConfigFolder returning the default configuration folder of the system
func GetDefaultConfigFolder() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFolder := fmt.Sprintf("%s/ReCoS", home)
	err = os.MkdirAll(configFolder, os.ModePerm)
	if err != nil {
		return "", err
	}
	return configFolder, nil
}

func ReplaceConfigdir(s string) (string, error) {
	if strings.Contains(s, "${configdir}") {
		configFolder, err := GetDefaultConfigFolder()
		if err != nil {
			return "", err
		}
		return strings.Replace(s, "${configdir}", configFolder, -1), nil
	}
	return s, nil
}
