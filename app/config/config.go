package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	holidays map[string]PublicHoliday
}

type PublicHoliday struct {
	Date        string   `json:"date"`
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Types       []string `json:"types"`
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
	Gray      string
	LightGray string
}

type Layout struct {
	Paper Paper

	Numbers Numbers
	Lengths Lengths
	Colors  Colors
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

	if err = cfg.PublicHolidays.Load(cfg.Year); err != nil {
		return cfg, fmt.Errorf("public holidays load: %w", err)
	}

	return cfg, nil
}

func (p *PublicHolidays) Load(year int) error {
	countryCodes := p.countryCodes()
	if len(countryCodes) == 0 {
		return nil
	}

	baseURL := strings.TrimRight(p.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://date.nager.at/api/v3"
	}

	p.holidays = make(map[string]PublicHoliday)
	client := http.Client{Timeout: 10 * time.Second}
	for _, countryCode := range countryCodes {
		resp, err := client.Get(fmt.Sprintf("%s/PublicHolidays/%d/%s", baseURL, year, countryCode))
		if err != nil {
			return fmt.Errorf("%s: %w", countryCode, err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return fmt.Errorf("%s: unexpected status %s", countryCode, resp.Status)
		}

		var holidays []PublicHoliday
		if err = json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
			resp.Body.Close()
			return fmt.Errorf("%s: %w", countryCode, err)
		}
		resp.Body.Close()

		for _, holiday := range holidays {
			if holiday.IsPublic() {
				p.holidays[holiday.Date] = holiday
			}
		}
	}

	return nil
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
