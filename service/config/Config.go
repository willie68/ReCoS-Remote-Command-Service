package config

import (
	"errors"
	"fmt"
	"os"

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

	HealthCheck HealthCheck `yaml:"healthcheck"`

	Profiles string `yaml:"profiles"`

	ExternalConfig map[string]interface{} `yaml:"extconfig"`
}

// HealthCheck configuration for the health check system
type HealthCheck struct {
	Period int `yaml:"period"`
}

var DefaulConfig = Config{
	Port:       9280,
	Sslport:    0,
	ServiceURL: "http://127.0.0.1:9280",
	SecretFile: "",
	HealthCheck: HealthCheck{
		Period: 30,
	},
	Profiles: "${configdir}/profiles",
	ExternalConfig: map[string]interface{}{
		"openhardwaremonitor": map[string]interface{}{
			"url":          "http://127.0.0.1:12999/data.json",
			"updateperiod": "5",
		},
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
