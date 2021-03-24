package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
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

	WebClient   string `yaml:"webclient"`
	AdminClient string `yaml:"adminclient"`

	Icons string `yaml:"icons"`

	HealthCheck HealthCheck `yaml:"healthcheck"`

	Profiles string `yaml:"profiles"`

	ExternalConfig map[string]interface{} `yaml:"extconfig"`

	Logging LoggingConfig `yaml:"logging"`
}

// HealthCheck configuration for the health check system
type HealthCheck struct {
	Period int `yaml:"period"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

var DefaulConfig = Config{
	Port:       9280,
	Sslport:    0,
	ServiceURL: "http://127.0.0.1:9280",
	SecretFile: "",
	Password:   "recosadmin",
	HealthCheck: HealthCheck{
		Period: 30,
	},
	Profiles:    "${configdir}/profiles",
	WebClient:   "${configdir}/webclient",
	AdminClient: "${configdir}/webadmin",
	Icons:       "${configdir}/webclient/assets",
	ExternalConfig: map[string]interface{}{
		"openhardwaremonitor": map[string]interface{}{
			"url":          "http://127.0.0.1:12999/data.json",
			"updateperiod": "5",
		},
	},
	Logging: LoggingConfig{
		Level: "INFO",
	},
}

// SaveConfig saving the config
func SaveConfig(folder string, config Config) error {
	filename := fmt.Sprintf("%s/service.yaml", folder)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// everything is ok, so please serialise the profile
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		err = yaml.NewEncoder(f).Encode(config)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Config already exists")
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
