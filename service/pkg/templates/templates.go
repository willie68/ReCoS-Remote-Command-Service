package templates

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

func Templates() []models.Template {
	templates := make([]models.Template, 0)
	templates = append(templates, models.Template{Group: "Elgato", Name: "Streamdeck"})
	templates = append(templates, models.Template{Group: "Elgato", Name: "Streamdeck_Mini"})
	templates = append(templates, models.Template{Group: "Elgato", Name: "Streamdeck_XL"})
	templates = append(templates, models.Template{Group: "Phones", Name: "phone_5x3_sp"})
	templates = append(templates, models.Template{Group: "Phones", Name: "phone_5x3_mp"})
	templates = append(templates, models.Template{Group: "Tablet", Name: "tab_8x5_sp"})
	templates = append(templates, models.Template{Group: "Tablet", Name: "phone_8x5_mp"})
	return templates
}
