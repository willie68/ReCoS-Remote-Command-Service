package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
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
var File = "config/service.yaml"

// Get returns loaded config
func Get() *Config {
	return &config
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
	config.AppID = DefaultConfig.AppID

	for k, v := range DefaultConfig.ExternalConfig {
		_, ok := config.ExternalConfig[k]
		if !ok {
			config.ExternalConfig[k] = v
		}
	}
	return readSecret()
}

func Save() error {
	return SaveConfig(File, config, true)
}

// SaveConfig saving the config
func SaveConfig(filename string, config Config, overwrite bool) error {
	if config.AppUUID == "" {
		uuid := uuid.New()
		config.AppUUID = uuid.String()
	}
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

func readSecret() error {
	secretFile := config.SecretFile
	if secretFile != "" {
		data, err := ioutil.ReadFile(secretFile)
		if err != nil {
			return fmt.Errorf("can't load secret file: %s", err.Error())
		}
		var secretConfig Secret = Secret{}
		err = yaml.Unmarshal(data, &secretConfig)
		if err != nil {
			return fmt.Errorf("can't unmarshal secret file: %s", err.Error())
		}
		mergeSecret(secretConfig)
	}
	return nil
}

func mergeSecret(secret Secret) {
	//	config.MongoDB.Username = secret.MongoDB.Username
	//	config.MongoDB.Password = secret.MongoDB.Password
}
