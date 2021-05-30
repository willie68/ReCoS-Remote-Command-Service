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

	Password string `yaml:"password"`

	HealthCheck HealthCheck `yaml:"healthcheck"`

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

var DefaulConfig = Config{
	Port:       9281,
	Sslport:    0,
	ServiceURL: "http://127.0.0.1:9281",
	Password:   "recosadmin",
	HealthCheck: HealthCheck{
		Period: 30,
	},
	Logging: LoggingConfig{
		Level:    "INFO",
		Filename: "${configdir}/pl_example.log",
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
