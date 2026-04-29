package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPublicHolidaysIncludeHolidayFiltersTypes(t *testing.T) {
	holidays := PublicHolidays{
		IncludeTypes: []string{"Public", "Bank"},
		ExcludeTypes: []string{"Bank"},
	}

	if holidays.includeHoliday(PublicHoliday{Date: "2026-01-01", Types: []string{"Bank"}}) {
		t.Fatal("excluded type must win over included type")
	}
	if !holidays.includeHoliday(PublicHoliday{Date: "2026-01-02", Types: []string{"Public"}}) {
		t.Fatal("public holiday should be included")
	}
}

func TestPublicHolidaysIncludeHolidayFiltersSubdivisions(t *testing.T) {
	holidays := PublicHolidays{Subdivisions: []string{"US-CA"}}

	if !holidays.includeHoliday(PublicHoliday{Date: "2026-01-01", Global: true, Types: []string{"Public"}}) {
		t.Fatal("global holiday should be included when subdivisions are selected")
	}
	if !holidays.includeHoliday(PublicHoliday{Date: "2026-03-31", Counties: []string{"US-CA"}, Types: []string{"Public"}}) {
		t.Fatal("matching subdivision holiday should be included")
	}
	if holidays.includeHoliday(PublicHoliday{Date: "2026-04-01", Counties: []string{"US-NY"}, Types: []string{"Public"}}) {
		t.Fatal("non-matching subdivision holiday should be excluded")
	}
}

func TestPublicHolidaysIncludeHolidayKeepsAllSubdivisionsByDefault(t *testing.T) {
	holidays := PublicHolidays{}

	if !holidays.includeHoliday(PublicHoliday{Date: "2026-03-31", Counties: []string{"US-CA"}, Types: []string{"Public"}}) {
		t.Fatal("regional public holidays should stay included unless subdivisions are configured")
	}
}

func TestConfigValidationRejectsInvalidRenderBlock(t *testing.T) {
	cfg := Config{
		Year:      2026,
		WeekStart: time.Monday,
		Pages:     Pages{{Name: "bad", RenderBlocks: RenderBlocks{{FuncName: "missing", Tpls: []string{"x.tpl"}}}}},
	}

	if err := cfg.defaultAndValidate(); err == nil {
		t.Fatal("invalid render block should fail validation")
	}
}

func TestConfigValidationRejectsInvalidCustomDate(t *testing.T) {
	cfg := Config{
		Year:      2026,
		WeekStart: time.Monday,
		Pages:     Pages{{Name: "weekly", RenderBlocks: RenderBlocks{{FuncName: "weekly", Tpls: []string{"weekly.tpl"}}}}},
		Events:    Events{Custom: []PublicHoliday{{Date: "2026/01/01", Name: "Bad", Types: []string{"Event"}}}},
	}

	if err := cfg.defaultAndValidate(); err == nil {
		t.Fatal("invalid event date should fail validation")
	}
}

func TestEventsLoadYAMLFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.yaml")
	if err := os.WriteFile(path, []byte(`events:
  - date: "2026-05-04"
    name: "Time off"
    shortName: "Off"
    types: [Event]
  - date: "2027-05-04"
    name: "Other year"
    types: [Event]
`), 0600); err != nil {
		t.Fatal(err)
	}

	holidays := PublicHolidays{}
	if err := (Events{Files: []string{path}}).Load(2026, &holidays); err != nil {
		t.Fatal(err)
	}

	event, ok := holidays.Holiday(time.Date(2026, time.May, 4, 0, 0, 0, 0, time.Local))
	if !ok {
		t.Fatal("expected 2026 event to be loaded")
	}
	if !event.IsEvent() || event.ShortName != "Off" {
		t.Fatalf("unexpected event: %#v", event)
	}
	if _, ok := holidays.Holiday(time.Date(2027, time.May, 4, 0, 0, 0, 0, time.Local)); ok {
		t.Fatal("event for another year should not be loaded")
	}
}

func TestReadICSHolidaysSupportsMultiDayAndCategories(t *testing.T) {
	path := writeTempICS(t, `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART;VALUE=DATE:20260504
DTEND;VALUE=DATE:20260507
SUMMARY:Time\, off
CATEGORIES:Vacation,Personal
END:VEVENT
END:VCALENDAR
`)

	events, err := readICSHolidays(path, 2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 all-day events, got %d: %#v", len(events), events)
	}
	if events[0].Date != "2026-05-04" || events[2].Date != "2026-05-06" {
		t.Fatalf("unexpected dates: %#v", events)
	}
	if events[0].Name != "Time, off" || !contains(events[0].Types, "Vacation") || !contains(events[0].Types, "Personal") {
		t.Fatalf("unexpected parsed event: %#v", events[0])
	}
}

func TestReadICSHolidaysSupportsTimezoneAndWeeklyRecurrence(t *testing.T) {
	path := writeTempICS(t, `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART;TZID=Europe/Warsaw:20260105T090000
DTEND;TZID=Europe/Warsaw:20260105T100000
SUMMARY:Standup
RRULE:FREQ=WEEKLY;COUNT=3
END:VEVENT
END:VCALENDAR
`)

	events, err := readICSHolidays(path, 2026)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"2026-01-05", "2026-01-12", "2026-01-19"}
	if len(events) != len(want) {
		t.Fatalf("expected %d events, got %d: %#v", len(want), len(events), events)
	}
	for i, date := range want {
		if events[i].Date != date {
			t.Fatalf("event %d date mismatch: want %s got %s", i, date, events[i].Date)
		}
	}
}

func TestReadICSHolidaysSupportsFoldedLinesAndUntil(t *testing.T) {
	path := writeTempICS(t, "BEGIN:VCALENDAR\nBEGIN:VEVENT\nDTSTART;VALUE=DATE:20261230\nSUMMARY:Long\n summary\nRRULE:FREQ=DAILY;UNTIL=20270102\nEND:VEVENT\nEND:VCALENDAR\n")

	events, err := readICSHolidays(path, 2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events in 2026, got %d: %#v", len(events), events)
	}
	if events[0].Name != "Longsummary" || events[1].Date != "2026-12-31" {
		t.Fatalf("unexpected folded/until events: %#v", events)
	}
}

func writeTempICS(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "events.ics")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	return path
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}

	return false
}
