package templates

import (
	"strings"

	"gopkg.in/yaml.v3"
	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/web"
)

func Templates() []models.Template {
	templates := make([]models.Template, 0)
	files, err := web.Templates.ReadDir("templates")
	if err != nil {
		clog.Logger.Debugf("Error reading mapper file: %v", err)
	}
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".yaml") {
			bytes, err := web.Templates.ReadFile("templates/" + file.Name())
			if err != nil {
				clog.Logger.Errorf("Error reading mapper file: %v", err)
				continue
			}
			profile := models.Profile{}
			err = yaml.Unmarshal(bytes, &profile)
			if err != nil {
				clog.Logger.Errorf("Error unmarshalling mapper file: %v", err)
				continue
			}
			templates = append(templates, models.Template{Group: profile.Group, Label: profile.Label, Name: profile.Name, Description: profile.Description})
		}
	}
	return templates
}

func GetTemplate(templatename string) models.Profile {
	files, err := web.Templates.ReadDir("templates")
	if err != nil {
		clog.Logger.Debugf("Error reading mapper file: %v", err)
	}
	for _, file := range files {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".yaml") {
			bytes, err := web.Templates.ReadFile("templates/" + file.Name())
			if err != nil {
				clog.Logger.Errorf("Error reading mapper file: %v", err)
				continue
			}
			profile := models.Profile{}
			err = yaml.Unmarshal(bytes, &profile)
			if err != nil {
				clog.Logger.Errorf("Error unmarshalling mapper file: %v", err)
				continue
			}
			if profile.Name == templatename {
				return profile
			}
		}
	}
	return models.Profile{}
}
