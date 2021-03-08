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

// SaveProfileFile saving the profile
func SaveProfileFile(profile models.Profile) error {
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

// UpdateProfileFile saving the profile
func UpdateProfileFile(profile models.Profile) error {
	filename := fmt.Sprintf("%s/%s.yaml", profileFolder, profile.Name)
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

// HasProfile chacking if a profile is already defined
func HasProfile(profileName string) bool {
	for _, profile := range Profiles {
		if strings.EqualFold(profileName, profile.Name) {
			return true
		}
	}
	return false
}

// AddProfile adding a profile to the profile list
func AddProfile(profile models.Profile) error {
	if HasProfile(profile.Name) {
		return errors.New("Profile already exists")
	}
	Profiles = append(Profiles, profile)
	return nil
}

// UpdateProfile adding a profile to the profile list
func UpdateProfile(profile models.Profile) error {
	if HasProfile(profile.Name) {
		Profiles = append(Profiles, profile)
		return nil
	}
	DeleteProfile(profile.Name)
	AddProfile(profile)
	return nil
}

// DeleteProfile adding a profile to the profile list
func DeleteProfile(profileName string) (models.Profile, error) {
	if !HasProfile(profileName) {
		return models.Profile{}, errors.New("Profile not exists")
	}
	var myProfile models.Profile
	for x, profile := range Profiles {
		if strings.EqualFold(profile.Name, profileName) {
			myProfile = profile
			Profiles = remove(Profiles, x)
			break
		}
	}
	filename := fmt.Sprintf("%s/%s.yaml", profileFolder, profileName)
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		err := os.Remove(filename)
		return myProfile, err
	}
	return models.Profile{}, errors.New("Profile not exists")
}

func remove(s []models.Profile, i int) []models.Profile {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}
