package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

var config = Config{
	Port:       0,
	Sslport:    0,
	ServiceURL: "http://127.0.0.1",
	HealthCheck: HealthCheck{
		Period: 30,
	},
}

// File the config file
var File = "config/pl_example.yaml"

// Get returns loaded config
func Get() Config {
	return config
}

// Load loads the config
func Load() error {
	_, err := os.Stat(File)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(File)
	if err != nil {
		return fmt.Errorf("can't load config file: %s", err.Error())
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("can't unmarshal config file: %s", err.Error())
	}
	return nil
}

func Save() error {
	return SaveConfig(File, config, true)
}

// SaveConfig saving the config
func SaveConfig(filename string, config Config, overwrite bool) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) || overwrite {
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
