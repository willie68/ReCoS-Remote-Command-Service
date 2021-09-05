package transformer

import (
	"fmt"

	clog "wkla.no-ip.biz/remote-desk-service/logging"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
)

func ExportPage(p models.Profile, pageName string) (*models.ProfileExchange, error) {
	page, err := p.GetPage(pageName)
	if err != nil {
		return nil, fmt.Errorf("can't find a page named \"%s\" in this profile \"%s\"", pageName, p.Name)
	}

	actions := make([]*models.Action, 0)
	for _, cell := range page.Cells {
		action, err := p.GetAction(cell)
		if err != nil {
			clog.Logger.Errorf("can't find action %s", cell)
			continue
		}
		actions = append(actions, action)
	}
	exchange := models.ProfileExchange{
		Type:    models.ExchangePage,
		Pages:   []models.Page{*page},
		Actions: actions,
	}
	return &exchange, nil
}

func CombineProfile(profile models.Profile, e models.ProfileExchange) (*models.Profile, error) {
	// adding actions to action and rename them is needed
	for _, eaction := range e.Actions {
		name := eaction.Name
		index := 0
		for {
			if profile.HasAction(eaction.Name) {
				eaction.Name = fmt.Sprintf("%s_%d", name, index)
				index++
				continue
			}
			// change cell name for that action
			for _, page := range e.Pages {
				for x, cell := range page.Cells {
					if cell == name {
						page.Cells[x] = eaction.Name
					}
				}
			}

			profile.Actions = append(profile.Actions, eaction)
			break
		}
	}

	// addding page as new page
	for _, epage := range e.Pages {
		name := epage.Name
		index := 0
		for {
			if profile.HasPage(epage.Name) {
				epage.Name = fmt.Sprintf("%s_%d", name, index)
				index++
				continue
			}

			profile.Pages = append(profile.Pages, epage)
			break
		}
	}

	return &profile, nil
}
