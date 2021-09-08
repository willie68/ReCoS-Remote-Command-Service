package transformer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

func TestCombine(t *testing.T) {

	p := models.Profile{
		Name: "default",
		Pages: []models.Page{
			{
				Name:  "default",
				Cells: []string{"ac_1", "ac_2"},
			},
			{
				Name:  "one",
				Cells: []string{"ac_1", "ac_3"},
			},
		},
		Actions: []*models.Action{
			{
				Name: "ac_1",
			},
			{
				Name: "ac_2",
			},
			{
				Name: "ac_3",
			},
		},
	}
	pe := models.ProfileExchange{
		Name: "default",
		Pages: []*models.Page{
			{
				Name:  "one",
				Cells: []string{"ac_1", "ac_4"},
			},
		},
		Actions: []*models.Action{
			{
				Name: "ac_1",
			},
			{
				Name: "ac_4",
			},
		},
	}

	profile, err := CombineProfile(p, pe)

	assert.Nil(t, err)
	assert.NotNil(t, profile)

	assert.True(t, profile.HasAction("ac_1_0"))
	assert.True(t, profile.HasAction("ac_4"))
	assert.True(t, profile.HasPage("one_0"))
	assert.Equal(t, 1, len(profile.Pages))
	assert.Equal(t, 2, len(profile.Actions))
}
