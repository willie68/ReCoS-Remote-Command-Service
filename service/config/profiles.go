package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

var Profiles []models.Profile
var profileFolder string

// GetProfileFolder returning the processed profile folder.
func GetProfileFolder(folder string) (string, error) {
	if strings.Contains(folder, "${configdir}") {
		configFolder, err := GetDefaultConfigFolder()
		if err != nil {
			return "", err
		}
		folder = fmt.Sprintf("%s/profiles", configFolder)
	}
	return folder, nil
}

// InitProfiles read all profile files from the filesystem
func InitProfiles(folder string) error {
	folder, err := GetProfileFolder(folder)
	if err != nil {
		return err
	}
	profileFolder = folder
	Profiles = make([]models.Profile, 0)
	var files []string
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return err
	}
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
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

// SaveProfile saving the profile
func SaveProfile(profile models.Profile) error {
	filename := fmt.Sprintf("%s/%s.yaml", profileFolder, profile.Name)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// everything is ok, so please serialise the profile
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		err = yaml.NewEncoder(f).Encode(profile)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Profile already exists")
}
