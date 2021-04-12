package models_test

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
	"wkla.no-ip.biz/remote-desk-service/dto"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

func TestNew(t *testing.T) {

	myCommand := models.Command{
		Type: dto.DelayCommandTypeInfo.Type,
		Name: "delay",
		Parameters: map[string]interface{}{
			"time":     10,
			"timeunit": time.Second,
		},
	}

	myAction := models.Action{
		Name:        "action1",
		Type:        models.Single,
		Description: "description for action",
		Title:       "Action Title",
		Commands:    []*models.Command{&myCommand},
	}

	profile := models.Profile{
		Name:        "Default",
		Description: "description for default",
		Pages: []models.Page{
			{
				Name:    "page1",
				Columns: 5,
				Rows:    5,
				Cells: []string{
					"action1",
				},
			},
		},
		Actions: []*models.Action{
			&myAction,
		},
	}

	//myString, _ := json.Marshal(profile)
	myString, _ := yaml.Marshal(profile)
	fmt.Println(string(myString))

}
