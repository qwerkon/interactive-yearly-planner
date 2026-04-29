package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug Debug

	Year                int `env:"PLANNER_YEAR"`
	WeekStart           time.Weekday
	Dotted              bool
	CalAfterSchedule    bool
	ClearTopRightCorner bool
	AMPMTime            bool
	AddLastHalfHour     bool
	PublicHolidays      PublicHolidays
	Events              Events

	Pages Pages

	Layout Layout
}

type Debug struct {
	ShowFrame bool
	ShowLinks bool
}

type PublicHolidays struct {
	CountryCode  string `env:"PLANNER_PUBLIC_HOLIDAYS_COUNTRY_CODE"`
	CountryCodes []string
	ShowNames    bool `env:"PLANNER_PUBLIC_HOLIDAYS_SHOW_NAMES"`
	BaseURL      string
	CacheDir     string `env:"PLANNER_PUBLIC_HOLIDAYS_CACHE_DIR"`
	RefreshCache bool   `env:"PLANNER_PUBLIC_HOLIDAYS_REFRESH_CACHE"`
	IncludeTypes []string
	ExcludeTypes []string
	Subdivisions []string
	ExcludeDates []string
	Custom       []PublicHoliday
	ICSFiles     []string

	holidays map[string]PublicHoliday
}

type Events struct {
	Files  []string
	Custom []PublicHoliday
}

type PublicHoliday struct {
	Date        string   `json:"date" yaml:"date"`
	LocalName   string   `json:"localName" yaml:"localName"`
	Name        string   `json:"name" yaml:"name"`
	ShortName   string   `json:"shortName" yaml:"shortName"`
	CountryCode string   `json:"countryCode" yaml:"countryCode"`
	Global      bool     `json:"global" yaml:"global"`
	Counties    []string `json:"counties" yaml:"counties"`
	Types       []string `json:"types" yaml:"types"`
}

type Pages []Page
type Page struct {
	Name         string
	RenderBlocks RenderBlocks
}

type RenderBlocks []RenderBlock

func (r Pages) WeeklyEnabled() bool {
	for _, s := range r {
		for _, block := range s.RenderBlocks {
			if block.FuncName == "weekly" {
				return true
			}
		}
	}

	return false
}

func (r Pages) DailyEnabled() bool {
	for _, s := range r {
		for _, block := range s.RenderBlocks {
			if block.FuncName == "daily" {
				return true
			}
		}
	}

	return false
}

func (r Pages) MonthlyEnabled() bool {
	for _, s := range r {
		for _, block := range s.RenderBlocks {
			if block.FuncName == "monthly" {
				return true
			}
		}
	}

	return false
}

func (r Pages) QuarterlyEnabled() bool {
	for _, s := range r {
		for _, block := range s.RenderBlocks {
			if block.FuncName == "quarterly" {
				return true
			}
		}
	}

	return false
}

func (r Pages) DailyReflectEnabled() bool {
	for _, s := range r {
		for _, block := range s.RenderBlocks {
			if block.FuncName == "daily_reflect" {
				return true
			}
		}
	}

	return false
}

func (r Pages) DailyNotesEnabled() bool {
	for _, s := range r {
		for _, block := range s.RenderBlocks {
			if block.FuncName == "daily_notes" {
				return true
			}
		}
	}

	return false
}

type RenderBlock struct {
	FuncName string
	Tpls     []string
}

type Colors struct {
	Gray          string
	LightGray     string
	Saturday      string
	Sunday        string
	PublicHoliday string
	Event         string
}

type Layout struct {
	Paper Paper

	Numbers     Numbers
	Lengths     Lengths
	Colors      Colors
	DeskWeekly  DeskWeekly
	DeskMonthly DeskMonthly
}

type DeskMonthly struct {
	HeaderFont     string
	SubHeaderFont  string
	SidePanelWidth string
	LegendFont     string
	ShowLegend     bool
	ShowSidePanel  bool
}

type DeskWeekly struct {
	StartHour         int
	EndHour           int
	HourLineHeight    string
	ColumnPadding     string
	HeaderPadding     string
	HeaderHeight      string
	MiniCalendarWidth string
	DayNumberFont     string
	WeekdayFont       string
	HourFont          string
	HolidayNameFont   string
	HolidayMarker     string
	EventMarker       string
	ShowMiniCalendar  bool
	ShowNotes         bool
	ShowHolidayMarker bool
	ShowHolidayLegend bool
}

type Numbers struct {
	ArrayStretch        float64
	QuarterlyLines      int
	WeeklyLines         int
	DailyTodos          int
	DailyNotes          int
	DailyPersonal       int
	DailyBottomHour     int
	DailyTopHour        int
	DailyDiaryGoals     int
	DailyDiaryGrateful  int
	DailyDiaryBest      int
	DailyDiaryLog       int
	TodoLinesInTodoPage int
	IndexMeetingNotes   int
	NotesIndexPages     int
	NotesOnPage         int
	DotHeightFull       int
	DotWidthFull        int
	DotWidthTwoThirds   int
}

type Paper struct {
	Width  string `env:"PLANNER_LAYOUT_PAPER_WIDTH"`
	Height string `env:"PLANNER_LAYOUT_PAPER_HEIGHT"`

	Margin Margin

	ReverseMargins bool
	MarginParWidth string
	MarginParSep   string
}

type Margin struct {
	Top    string `env:"PLANNER_LAYOUT_PAPER_MARGIN_TOP"`
	Bottom string `env:"PLANNER_LAYOUT_PAPER_MARGIN_BOTTOM"`
	Left   string `env:"PLANNER_LAYOUT_PAPER_MARGIN_LEFT"`
	Right  string `env:"PLANNER_LAYOUT_PAPER_MARGIN_RIGHT"`
}

func New(pathConfigs ...string) (Config, error) {
	var (
		bts []byte
		err error
		cfg Config
	)

	for _, filepath := range pathConfigs {
		if bts, err = ioutil.ReadFile(strings.ToLower(filepath)); err != nil {
			return cfg, fmt.Errorf("ioutil read file: %w", err)
		}

		if err = yaml.Unmarshal(bts, &cfg); err != nil {
			return cfg, fmt.Errorf("yaml unmarshal: %w", err)
		}
	}

	if err = env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("env parse: %w", err)
	}

	if cfg.Year == 0 {
		cfg.Year = time.Now().Year()
	}
	if err = cfg.defaultAndValidate(); err != nil {
		return cfg, err
	}

	if err = cfg.PublicHolidays.Load(cfg.Year); err != nil {
		return cfg, fmt.Errorf("public holidays load: %w", err)
	}
	if err = cfg.Events.Load(cfg.Year, &cfg.PublicHolidays); err != nil {
		return cfg, fmt.Errorf("events load: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) defaultAndValidate() error {
	if cfg.Year < 1900 || cfg.Year > 2200 {
		return fmt.Errorf("invalid year %d", cfg.Year)
	}
	if cfg.WeekStart < time.Sunday || cfg.WeekStart > time.Saturday {
		return fmt.Errorf("invalid weekstart %d", cfg.WeekStart)
	}
	if err := cfg.Pages.validate(); err != nil {
		return err
	}
	if cfg.Layout.Colors.Saturday == "" {
		cfg.Layout.Colors.Saturday = "gray!60!black"
	}
	if cfg.Layout.Colors.Sunday == "" {
		cfg.Layout.Colors.Sunday = "red!70!black"
	}
	if cfg.Layout.Colors.PublicHoliday == "" {
		cfg.Layout.Colors.PublicHoliday = cfg.Layout.Colors.Sunday
	}
	if cfg.Layout.Colors.Event == "" {
		cfg.Layout.Colors.Event = cfg.Layout.Colors.PublicHoliday
	}
	if cfg.Layout.DeskWeekly.StartHour == 0 {
		cfg.Layout.DeskWeekly.StartHour = 6
	}
	if cfg.Layout.DeskWeekly.EndHour == 0 {
		cfg.Layout.DeskWeekly.EndHour = 22
	}
	if cfg.Layout.DeskWeekly.StartHour < 0 || cfg.Layout.DeskWeekly.EndHour > 23 || cfg.Layout.DeskWeekly.StartHour > cfg.Layout.DeskWeekly.EndHour {
		return fmt.Errorf("invalid layout.deskweekly hour range %d-%d", cfg.Layout.DeskWeekly.StartHour, cfg.Layout.DeskWeekly.EndHour)
	}
	if cfg.Layout.DeskWeekly.HourLineHeight == "" {
		cfg.Layout.DeskWeekly.HourLineHeight = "4.25mm"
	}
	if cfg.Layout.DeskWeekly.ColumnPadding == "" {
		cfg.Layout.DeskWeekly.ColumnPadding = "1mm"
	}
	if cfg.Layout.DeskWeekly.HeaderPadding == "" {
		cfg.Layout.DeskWeekly.HeaderPadding = "1mm"
	}
	if cfg.Layout.DeskWeekly.HeaderHeight == "" {
		cfg.Layout.DeskWeekly.HeaderHeight = "13mm"
	}
	if cfg.Layout.DeskWeekly.MiniCalendarWidth == "" {
		cfg.Layout.DeskWeekly.MiniCalendarWidth = "4.25cm"
	}
	if cfg.Layout.DeskWeekly.DayNumberFont == "" {
		cfg.Layout.DeskWeekly.DayNumberFont = "18"
	}
	if cfg.Layout.DeskWeekly.WeekdayFont == "" {
		cfg.Layout.DeskWeekly.WeekdayFont = "6.5"
	}
	if cfg.Layout.DeskWeekly.HourFont == "" {
		cfg.Layout.DeskWeekly.HourFont = "5.5"
	}
	if cfg.Layout.DeskWeekly.HolidayNameFont == "" {
		cfg.Layout.DeskWeekly.HolidayNameFont = "4.3"
	}
	if cfg.Layout.DeskWeekly.HolidayMarker == "" {
		cfg.Layout.DeskWeekly.HolidayMarker = "\\textbullet"
	}
	if cfg.Layout.DeskWeekly.EventMarker == "" {
		cfg.Layout.DeskWeekly.EventMarker = "*"
	}
	if cfg.Layout.DeskMonthly.HeaderFont == "" {
		cfg.Layout.DeskMonthly.HeaderFont = "15"
	}
	if cfg.Layout.DeskMonthly.SubHeaderFont == "" {
		cfg.Layout.DeskMonthly.SubHeaderFont = "7"
	}
	if cfg.Layout.DeskMonthly.SidePanelWidth == "" {
		cfg.Layout.DeskMonthly.SidePanelWidth = "4.2cm"
	}
	if cfg.Layout.DeskMonthly.LegendFont == "" {
		cfg.Layout.DeskMonthly.LegendFont = "5.5"
	}
	if err := cfg.PublicHolidays.validate(); err != nil {
		return err
	}
	if err := cfg.Events.validate(); err != nil {
		return err
	}

	return nil
}

func (r Pages) validate() error {
	if len(r) == 0 {
		return fmt.Errorf("pages must not be empty")
	}
	allowed := map[string]bool{
		"title": true, "annual": true, "quarterly": true, "monthly": true, "weekly": true,
		"daily": true, "daily_reflect": true, "daily_notes": true, "notes_indexed": true,
	}
	for i, page := range r {
		if strings.TrimSpace(page.Name) == "" {
			return fmt.Errorf("pages[%d].name is required", i)
		}
		if len(page.RenderBlocks) == 0 {
			return fmt.Errorf("pages[%d].renderblocks must not be empty", i)
		}
		for j, block := range page.RenderBlocks {
			if !allowed[block.FuncName] {
				return fmt.Errorf("pages[%d].renderblocks[%d].funcname %q is invalid", i, j, block.FuncName)
			}
			if len(block.Tpls) == 0 {
				return fmt.Errorf("pages[%d].renderblocks[%d].tpls must not be empty", i, j)
			}
		}
	}

	return nil
}

func (p PublicHolidays) validate() error {
	for _, code := range p.countryCodes() {
		if len(code) != 2 {
			return fmt.Errorf("invalid publicholidays country code %q", code)
		}
	}
	for _, date := range p.ExcludeDates {
		if err := validateDate(date); err != nil {
			return fmt.Errorf("publicholidays.excludedates %q: %w", date, err)
		}
	}
	for i, holiday := range p.Custom {
		if err := holiday.validate(fmt.Sprintf("publicholidays.custom[%d]", i)); err != nil {
			return err
		}
	}
	for _, path := range p.ICSFiles {
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("publicholidays.icsfiles contains an empty path")
		}
	}

	return nil
}

func (e Events) validate() error {
	for _, path := range e.Files {
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("events.files contains an empty path")
		}
	}
	for i, event := range e.Custom {
		if err := event.validate(fmt.Sprintf("events.custom[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

func (p PublicHoliday) validate(prefix string) error {
	if err := validateDate(p.Date); err != nil {
		return fmt.Errorf("%s.date %q: %w", prefix, p.Date, err)
	}
	if strings.TrimSpace(p.Name) == "" && strings.TrimSpace(p.LocalName) == "" && strings.TrimSpace(p.ShortName) == "" {
		return fmt.Errorf("%s requires name, localName, or shortName", prefix)
	}
	if len(p.Types) == 0 {
		return fmt.Errorf("%s.types must not be empty", prefix)
	}

	return nil
}

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf("expected YYYY-MM-DD")
	}

	return nil
}

func (d DeskWeekly) Hours() []string {
	hours := make([]string, 0, d.EndHour-d.StartHour+1)
	for hour := d.StartHour; hour <= d.EndHour; hour++ {
		hours = append(hours, fmt.Sprintf("%02d:00", hour))
	}

	return hours
}

func (p *PublicHolidays) Load(year int) error {
	countryCodes := p.countryCodes()
	if len(countryCodes) == 0 {
		return p.loadLocal(year)
	}

	baseURL := strings.TrimRight(p.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://date.nager.at/api/v3"
	}

	p.holidays = make(map[string]PublicHoliday)
	client := http.Client{Timeout: 10 * time.Second}
	for _, countryCode := range countryCodes {
		holidays, err := p.loadCountry(client, baseURL, year, countryCode)
		if err != nil {
			return err
		}
		p.addHolidays(holidays)
	}

	return p.loadLocal(year)
}

func (p *PublicHolidays) loadCountry(client http.Client, baseURL string, year int, countryCode string) ([]PublicHoliday, error) {
	cachePath := p.cachePath(year, countryCode)
	if cachePath != "" && !p.RefreshCache {
		if bts, err := ioutil.ReadFile(cachePath); err == nil {
			var holidays []PublicHoliday
			if err = json.Unmarshal(bts, &holidays); err != nil {
				return nil, fmt.Errorf("%s cache: %w", countryCode, err)
			}

			return holidays, nil
		}
	}

	resp, err := client.Get(fmt.Sprintf("%s/PublicHolidays/%d/%s", baseURL, year, countryCode))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", countryCode, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s: unexpected status %s", countryCode, resp.Status)
	}

	var holidays []PublicHoliday
	if err = json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, fmt.Errorf("%s: %w", countryCode, err)
	}

	if cachePath != "" {
		if err = os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
			return nil, fmt.Errorf("%s cache mkdir: %w", countryCode, err)
		}
		bts, _ := json.MarshalIndent(holidays, "", "  ")
		if err = ioutil.WriteFile(cachePath, bts, 0600); err != nil {
			return nil, fmt.Errorf("%s cache write: %w", countryCode, err)
		}
	}

	return holidays, nil
}

func (p *PublicHolidays) loadLocal(year int) error {
	if p.holidays == nil {
		p.holidays = map[string]PublicHoliday{}
	}
	p.addHolidays(p.Custom)

	for _, path := range p.ICSFiles {
		holidays, err := readICSHolidays(path, year)
		if err != nil {
			return err
		}
		p.addHolidays(holidays)
	}

	for _, date := range p.ExcludeDates {
		delete(p.holidays, date)
	}

	return nil
}

func (e Events) Load(year int, p *PublicHolidays) error {
	if p.holidays == nil {
		p.holidays = map[string]PublicHoliday{}
	}
	p.addHolidays(normalizeEvents(e.Custom))

	for _, path := range e.Files {
		events, err := readYAMLEvents(path, year)
		if err != nil {
			return err
		}
		p.addHolidays(events)
	}

	return nil
}

func normalizeEvents(events []PublicHoliday) []PublicHoliday {
	normalized := make([]PublicHoliday, 0, len(events))
	for _, event := range events {
		if len(event.Types) == 0 {
			event.Types = []string{"Event"}
		}
		normalized = append(normalized, event)
	}

	return normalized
}

func (p PublicHolidays) cachePath(year int, countryCode string) string {
	if p.CacheDir == "" {
		return ""
	}

	return filepath.Join(p.CacheDir, fmt.Sprintf("%d-%s.json", year, countryCode))
}

func (p *PublicHolidays) addHolidays(holidays []PublicHoliday) {
	for _, holiday := range holidays {
		if p.includeHoliday(holiday) {
			p.holidays[holiday.Date] = holiday
		}
	}
}

func (p PublicHolidays) includeHoliday(holiday PublicHoliday) bool {
	if holiday.Date == "" {
		return false
	}
	if p.excludedByType(holiday) {
		return false
	}
	if !p.includedBySubdivision(holiday) {
		return false
	}
	if len(p.IncludeTypes) == 0 {
		return holiday.IsPublic() || holiday.IsEvent()
	}

	types := normalizedSet(p.IncludeTypes)
	for _, typ := range holiday.Types {
		if types[strings.ToUpper(strings.TrimSpace(typ))] {
			return true
		}
	}

	return false
}

func (p PublicHolidays) excludedByType(holiday PublicHoliday) bool {
	if len(p.ExcludeTypes) == 0 {
		return false
	}

	types := normalizedSet(p.ExcludeTypes)
	for _, typ := range holiday.Types {
		if types[strings.ToUpper(strings.TrimSpace(typ))] {
			return true
		}
	}

	return false
}

func (p PublicHolidays) includedBySubdivision(holiday PublicHoliday) bool {
	if len(p.Subdivisions) == 0 {
		return true
	}
	if holiday.Global {
		return true
	}

	subdivisions := normalizedSet(p.Subdivisions)
	for _, county := range holiday.Counties {
		if subdivisions[strings.ToUpper(strings.TrimSpace(county))] {
			return true
		}
	}

	return false
}

func normalizedSet(values []string) map[string]bool {
	set := map[string]bool{}
	for _, value := range values {
		value = strings.ToUpper(strings.TrimSpace(value))
		if value != "" {
			set[value] = true
		}
	}

	return set
}

func (p PublicHolidays) countryCodes() []string {
	codes := make([]string, 0, len(p.CountryCodes)+1)
	seen := map[string]bool{}

	for _, code := range p.CountryCodes {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code != "" && !seen[code] {
			codes = append(codes, code)
			seen[code] = true
		}
	}

	for _, code := range strings.Split(p.CountryCode, ",") {
		code = strings.ToUpper(strings.TrimSpace(code))
		if code != "" && !seen[code] {
			codes = append(codes, code)
			seen[code] = true
		}
	}

	return codes
}

func (p PublicHolidays) IsPublicHoliday(day time.Time) bool {
	_, ok := p.Holiday(day)
	return ok
}

func (p PublicHolidays) Holiday(day time.Time) (PublicHoliday, bool) {
	if len(p.holidays) == 0 {
		return PublicHoliday{}, false
	}

	holiday, ok := p.holidays[day.Format("2006-01-02")]
	return holiday, ok
}

func (p PublicHoliday) IsPublic() bool {
	for _, typ := range p.Types {
		if typ == "Public" {
			return true
		}
	}

	return false
}

func (p PublicHoliday) IsEvent() bool {
	for _, typ := range p.Types {
		if strings.EqualFold(typ, "Event") {
			return true
		}
	}

	return false
}

func readICSHolidays(path string, year int) ([]PublicHoliday, error) {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ics read %s: %w", path, err)
	}

	var holidays []PublicHoliday
	var event icsEvent
	inEvent := false
	for _, line := range unfoldICSLines(string(bts)) {
		switch {
		case line == "BEGIN:VEVENT":
			inEvent = true
			event = icsEvent{}
		case line == "END:VEVENT":
			if inEvent {
				expanded, err := event.Holidays(year)
				if err != nil {
					return nil, fmt.Errorf("ics event %q: %w", event.Summary, err)
				}
				holidays = append(holidays, expanded...)
			}
			inEvent = false
		case inEvent:
			event.ApplyLine(line)
		}
	}

	sort.Slice(holidays, func(i, j int) bool { return holidays[i].Date < holidays[j].Date })
	return holidays, nil
}

type icsEvent struct {
	Start      icsDate
	End        icsDate
	Summary    string
	Categories []string
	RRule      map[string]string
}

type icsDate struct {
	Time   time.Time
	AllDay bool
	Valid  bool
}

func (e *icsEvent) ApplyLine(line string) {
	name, params, value, ok := splitICSLine(line)
	if !ok {
		return
	}

	switch name {
	case "DTSTART":
		e.Start = parseICSDate(value, params)
	case "DTEND":
		e.End = parseICSDate(value, params)
	case "SUMMARY":
		e.Summary = unescapeICSValue(value)
	case "CATEGORIES":
		e.Categories = splitICSCategories(value)
	case "RRULE":
		e.RRule = parseICSRRule(value)
	}
}

func (e icsEvent) Holidays(year int) ([]PublicHoliday, error) {
	if !e.Start.Valid || strings.TrimSpace(e.Summary) == "" {
		return nil, nil
	}

	starts, err := e.occurrences(year)
	if err != nil {
		return nil, err
	}

	types := append([]string{"Event"}, e.Categories...)
	durationDays := e.durationDays()
	seen := map[string]bool{}
	holidays := []PublicHoliday{}
	for _, start := range starts {
		for i := 0; i < durationDays; i++ {
			day := start.AddDate(0, 0, i)
			if day.Year() != year {
				continue
			}
			date := day.Format("2006-01-02")
			if seen[date] {
				continue
			}
			seen[date] = true
			holidays = append(holidays, PublicHoliday{Date: date, LocalName: e.Summary, Name: e.Summary, ShortName: e.Summary, Types: types})
		}
	}

	sort.Slice(holidays, func(i, j int) bool { return holidays[i].Date < holidays[j].Date })
	return holidays, nil
}

func (e icsEvent) durationDays() int {
	if !e.End.Valid || !e.End.Time.After(e.Start.Time) {
		return 1
	}
	end := e.End.Time
	if e.End.AllDay {
		end = end.AddDate(0, 0, -1)
	}
	days := int(end.Sub(startOfDay(e.Start.Time)).Hours()/24) + 1
	if days < 1 {
		return 1
	}

	return days
}

func (e icsEvent) occurrences(year int) ([]time.Time, error) {
	if len(e.RRule) == 0 {
		return []time.Time{startOfDay(e.Start.Time)}, nil
	}

	freq := strings.ToUpper(e.RRule["FREQ"])
	if freq == "" {
		return []time.Time{startOfDay(e.Start.Time)}, nil
	}
	count := 370
	if raw := e.RRule["COUNT"]; raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid RRULE COUNT %q", raw)
		}
		count = parsed
	}
	interval := 1
	if raw := e.RRule["INTERVAL"]; raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 1 {
			return nil, fmt.Errorf("invalid RRULE INTERVAL %q", raw)
		}
		interval = parsed
	}
	until := time.Date(year, time.December, 31, 23, 59, 59, 0, e.Start.Time.Location())
	if raw := e.RRule["UNTIL"]; raw != "" {
		parsed := parseICSDate(raw, nil)
		if !parsed.Valid {
			return nil, fmt.Errorf("invalid RRULE UNTIL %q", raw)
		}
		until = parsed.Time
	}

	starts := []time.Time{}
	current := startOfDay(e.Start.Time)
	for i := 0; i < count && !current.After(until); i++ {
		if current.Year() == year {
			starts = append(starts, current)
		}
		if current.Year() > year && current.After(until) {
			break
		}
		switch freq {
		case "DAILY":
			current = current.AddDate(0, 0, interval)
		case "WEEKLY":
			current = current.AddDate(0, 0, 7*interval)
		case "MONTHLY":
			current = current.AddDate(0, interval, 0)
		case "YEARLY":
			current = current.AddDate(interval, 0, 0)
		default:
			return nil, fmt.Errorf("unsupported RRULE FREQ %q", freq)
		}
	}

	return starts, nil
}

func unfoldICSLines(content string) []string {
	lines := []string{}
	for _, raw := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		raw = strings.TrimRight(raw, "\r")
		if raw == "" {
			continue
		}
		if (strings.HasPrefix(raw, " ") || strings.HasPrefix(raw, "\t")) && len(lines) > 0 {
			lines[len(lines)-1] += strings.TrimLeft(raw, " \t")
			continue
		}
		lines = append(lines, strings.TrimSpace(raw))
	}

	return lines
}

func splitICSLine(line string) (string, map[string]string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", nil, "", false
	}
	left := strings.Split(parts[0], ";")
	name := strings.ToUpper(strings.TrimSpace(left[0]))
	params := map[string]string{}
	for _, param := range left[1:] {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			params[strings.ToUpper(strings.TrimSpace(kv[0]))] = strings.Trim(strings.TrimSpace(kv[1]), `"`)
		}
	}

	return name, params, strings.TrimSpace(parts[1]), true
}

func parseICSDate(value string, params map[string]string) icsDate {
	value = strings.TrimSpace(strings.TrimSuffix(value, "Z"))
	if len(value) >= len("20060102T150405") {
		loc := time.Local
		if tzid := params["TZID"]; tzid != "" {
			if loaded, err := time.LoadLocation(tzid); err == nil {
				loc = loaded
			}
		}
		if parsed, err := time.ParseInLocation("20060102T150405", value[:len("20060102T150405")], loc); err == nil {
			return icsDate{Time: parsed, Valid: true}
		}
	}
	if len(value) >= len("20060102") {
		if parsed, err := time.ParseInLocation("20060102", value[:len("20060102")], time.Local); err == nil {
			return icsDate{Time: parsed, AllDay: true, Valid: true}
		}
	}

	return icsDate{}
}

func parseICSRRule(value string) map[string]string {
	rule := map[string]string{}
	for _, part := range strings.Split(value, ";") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			rule[strings.ToUpper(strings.TrimSpace(kv[0]))] = strings.ToUpper(strings.TrimSpace(kv[1]))
		}
	}

	return rule
}

func splitICSCategories(value string) []string {
	categories := []string{}
	seen := map[string]bool{"EVENT": true}
	for _, category := range strings.Split(value, ",") {
		category = strings.TrimSpace(unescapeICSValue(category))
		if category == "" || seen[strings.ToUpper(category)] {
			continue
		}
		seen[strings.ToUpper(category)] = true
		categories = append(categories, category)
	}

	return categories
}

func unescapeICSValue(value string) string {
	value = strings.ReplaceAll(value, `\n`, " ")
	value = strings.ReplaceAll(value, `\N`, " ")
	value = strings.ReplaceAll(value, `\,`, ",")
	value = strings.ReplaceAll(value, `\;`, ";")
	value = strings.ReplaceAll(value, `\\`, `\`)
	return strings.TrimSpace(value)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func readYAMLEvents(path string, year int) ([]PublicHoliday, error) {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("events read %s: %w", path, err)
	}

	var file struct {
		Events []PublicHoliday
	}
	if err = yaml.Unmarshal(bts, &file); err != nil {
		return nil, fmt.Errorf("events yaml %s: %w", path, err)
	}

	events := make([]PublicHoliday, 0, len(file.Events))
	for i, event := range file.Events {
		if err := event.validate(fmt.Sprintf("%s.events[%d]", path, i)); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(event.Date, strconv.Itoa(year)) {
			continue
		}
		events = append(events, normalizeEvents([]PublicHoliday{event})[0])
	}

	sort.Slice(events, func(i, j int) bool { return events[i].Date < events[j].Date })
	return events, nil
}

func formatICSDate(date string) string {
	date = strings.TrimSuffix(date, "Z")
	if len(date) >= len("20060102") {
		date = date[:len("20060102")]
	}
	if len(date) != len("20060102") {
		return date
	}

	return date[:4] + "-" + date[4:6] + "-" + date[6:8]
}
