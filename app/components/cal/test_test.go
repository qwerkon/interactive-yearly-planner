package cal_test

import (
	"strings"
	"testing"
	"time"

	"github.com/qwerkon/interactive-yearly-planner/app/components/cal"
	"github.com/qwerkon/interactive-yearly-planner/app/config"
)

func TestWeekRefsAreStableAcrossYearBoundaries(t *testing.T) {
	year := cal.NewYear(time.Monday, 2026)
	weeks := cal.NewWeeksForYear(time.Monday, year)

	for _, week := range weeks {
		if week.Ref() == "" {
			t.Fatal("week ref must not be empty")
		}
		if !strings.Contains(week.Target(), week.Ref()) {
			t.Fatalf("target %q must contain ref %q", week.Target(), week.Ref())
		}
	}
}

func TestMonthRowsUseWeekRefForDayLinks(t *testing.T) {
	year := cal.NewYear(time.Monday, 2026)
	month := cal.NewMonth(time.Monday, year, cal.NewQuarter(time.Monday, year, 1), time.January)

	for _, week := range month.Weeks {
		ref := week.Ref()
		for _, day := range week.Days {
			if day.Time.IsZero() {
				continue
			}
			link := day.Day(nil, false, ref)
			if !strings.Contains(link, ref) {
				t.Fatalf("day link %q must point to row ref %q", link, ref)
			}
		}
	}
}

func TestMonthRowsLatexLinkSnapshot(t *testing.T) {
	cfg := config.Config{
		Year:      2026,
		WeekStart: time.Monday,
		Pages: config.Pages{
			{RenderBlocks: config.RenderBlocks{{FuncName: "monthly"}, {FuncName: "weekly"}}},
		},
		Layout: config.Layout{Colors: config.Colors{
			Saturday:      "gray!60!black",
			Sunday:        "red!70!black",
			PublicHoliday: "red!70!black",
		}},
	}
	year := cal.NewYear(time.Monday, 2026)
	month := cal.NewMonth(time.Monday, year, cal.NewQuarter(time.Monday, year, 1), time.January)

	got := renderMonthRowsLatexSnapshot(cfg, month)
	want := strings.Join([]string{
		`\hyperlink{Week 1}{\rotatebox[origin=tr]{90}{\makebox[\myLenMonthlyCellHeight][c]{Week 1}}} &  &  &  & \hyperlink{Week 1}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}1\\ \hline\end{tabular}} & \hyperlink{Week 1}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}2\\ \hline\end{tabular}} & \hyperlink{Week 1}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}\textcolor{gray!60!black}{3}\\ \hline\end{tabular}} & \hyperlink{Week 1}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}\textcolor{red!70!black}{4}\\ \hline\end{tabular}}`,
		`\hyperlink{Week 2}{\rotatebox[origin=tr]{90}{\makebox[\myLenMonthlyCellHeight][c]{Week 2}}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}5\\ \hline\end{tabular}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}6\\ \hline\end{tabular}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}7\\ \hline\end{tabular}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}8\\ \hline\end{tabular}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}9\\ \hline\end{tabular}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}\textcolor{gray!60!black}{10}\\ \hline\end{tabular}} & \hyperlink{Week 2}{\begin{tabular}{@{}p{5mm}@{}|}\hfil{}\textcolor{red!70!black}{11}\\ \hline\end{tabular}}`,
	}, "\n")

	if got != want {
		t.Fatalf("latex link snapshot mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestAnnualWithoutWeeklyLinksWeekNumberToFirstVisibleDay(t *testing.T) {
	cfg := config.Config{
		Year:      2026,
		WeekStart: time.Monday,
		Pages: config.Pages{
			{RenderBlocks: config.RenderBlocks{{FuncName: "annual"}, {FuncName: "daily"}}},
		},
	}
	year := cal.NewYear(time.Monday, 2026)
	month := cal.NewMonth(time.Monday, year, cal.NewQuarter(time.Monday, year, 1), time.January)
	week := month.Weeks[0]

	weekNumber := week.WeekNumberDay(false)
	if !strings.Contains(weekNumber, "2026-01-01T00:00:00") {
		t.Fatalf("week number should link to first visible day, got %q", weekNumber)
	}

	dayLink := week.Days[3].Day(nil, false, cfg)
	if !strings.Contains(dayLink, "2026-01-01T00:00:00") {
		t.Fatalf("day should link to daily page, got %q", dayLink)
	}
	if strings.Contains(dayLink, "Week 1") {
		t.Fatalf("day should not link to week when daily pages exist, got %q", dayLink)
	}
}

func renderMonthRowsLatexSnapshot(cfg config.Config, month *cal.Month) string {
	rows := make([]string, 0, 2)
	for _, week := range month.Weeks[:2] {
		cells := []string{week.WeekNumber(true)}
		for _, day := range week.Days {
			cells = append(cells, day.Day(nil, true, cfg, week.Ref()))
		}
		rows = append(rows, strings.Join(cells, " & "))
	}

	return strings.Join(rows, "\n")
}
