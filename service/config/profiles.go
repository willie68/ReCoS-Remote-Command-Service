package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var Profiles []models.Profile

// InitProfiles read all profile files from the filesystem
func InitProfiles(folder string) error {
	Profiles = make([]models.Profile, 0)
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), ".yaml") {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		var profile models.Profile
		err = yaml.Unmarshal(fileContent, &profile)
		if err != nil {
			return err
		}
		Profiles = append(Profiles, profile)
	}
	return nil
}
