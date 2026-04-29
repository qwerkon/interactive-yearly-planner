package compose

import (
	"strconv"

	"github.com/qwerkon/interactive-yearly-planner/app/components/cal"
	"github.com/qwerkon/interactive-yearly-planner/app/components/header"
	"github.com/qwerkon/interactive-yearly-planner/app/components/page"
	"github.com/qwerkon/interactive-yearly-planner/app/config"
)

func Annual(cfg config.Config, tpls []string) (page.Modules, error) {
	year := cal.NewYear(cfg.WeekStart, cfg.Year)
	extra := header.Items{}
	if notesEnabled(cfg) {
		extra = header.Items{header.NewTextItem("Notes").RefText("Notes Index")}
	}

	return page.Modules{{
		Cfg: cfg,
		Tpl: tpls[0],
		Body: map[string]interface{}{
			"Year":         year,
			"Breadcrumb":   year.Breadcrumb(),
			"HeadingMOS":   year.HeadingMOS(),
			"SideQuarters": year.SideQuarters(0),
			"SideMonths":   year.SideMonths(0),
			"Extra":        extra.WithTopRightCorner(cfg.ClearTopRightCorner),
			"Extra2":       extra2(cfg, true, false, nil, 0),
		},
	}}, nil
}

func notesEnabled(cfg config.Config) bool {
	return cfg.Layout.Numbers.NotesIndexPages > 0 && cfg.Layout.Numbers.NotesOnPage > 0
}

func extra2(cfg config.Config, sel1, sel2 bool, week *cal.Week, idxPage int) header.Items {
	items := make(header.Items, 0, 3)

	if cfg.Pages.WeeklyEnabled() && week != nil {
		items = append(items, header.NewCellItem(week.Name()))
	}

	items = append(items, header.NewCellItem("Calendar").Selected(sel1))

	if notesEnabled(cfg) {
		if idxPage > 0 {
			suffix := ""
			if idxPage > 1 {
				suffix = " " + strconv.Itoa(idxPage)
			}

			items = append(items, header.NewCellItem("Notes").Refer("Notes Index"+suffix).Selected(sel2))
		} else {
			items = append(items, header.NewCellItem("Notes").Refer("Notes Index").Selected(sel2))
		}
	}

	if len(items) == 1 {
		return items
	}

	return items.WithTopRightCorner(cfg.ClearTopRightCorner)
}
