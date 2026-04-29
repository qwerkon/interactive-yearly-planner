package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kudrykv/latex-yearly-planner/app"
	"github.com/urfave/cli/v2"
)

var months = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
var monthsShort = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
var weekdays = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
var weekdaysShort = []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}

type profile struct {
	Config string
	Name   string
}

var profiles = map[string]profile{
	"rm2-minimal": {
		Config: "cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.minimal.yaml,cfg/template_months_on_side_minimal.yaml,cfg/rm2.mos.default.yaml",
		Name:   "rm2.minimal",
	},
	"rm2-desk-weekly": {
		Config: "cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.desk.weekly.yaml,cfg/template_desk_weekly_rm2.yaml,cfg/rm2.mos.default.yaml",
		Name:   "rm2.desk.weekly",
	},
	"rm2-desk-monthly-weekly": {
		Config: "cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.desk.weekly.yaml,cfg/rm2.desk.monthly.yaml,cfg/template_desk_monthly_rm2.yaml,cfg/rm2.mos.default.yaml",
		Name:   "rm2.desk.monthly-weekly",
	},
}

func main() {
	app := &cli.App{
		Name:  "planner",
		Usage: "build configured PDF planners",
		Commands: []*cli.Command{
			buildCommand(false),
			buildCommand(true),
			{
				Name:  "profiles",
				Usage: "list available planner profiles",
				Action: func(*cli.Context) error {
					for name := range profiles {
						fmt.Println(name)
					}
					return nil
				},
			},
			{
				Name:  "build-range",
				Usage: "build a profile for an inclusive year range",
				Flags: commonFlags(),
				Action: func(c *cli.Context) error {
					from := c.Int("from")
					to := c.Int("to")
					if from == 0 || to == 0 || from > to {
						return fmt.Errorf("pass a valid --from and --to range")
					}

					for year := from; year <= to; year++ {
						if err := runBuild(c, year, false); err != nil {
							return err
						}
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func buildCommand(preview bool) *cli.Command {
	name := "build"
	usage := "build a full planner PDF"
	if preview {
		name = "preview"
		usage = "build a preview planner PDF"
	}

	return &cli.Command{
		Name:  name,
		Usage: usage,
		Flags: commonFlags(),
		Action: func(c *cli.Context) error {
			year := c.Int("year")
			if year == 0 {
				year = time.Now().Year() + 1
			}

			return runBuild(c, year, preview)
		},
	}
}

func commonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{Name: "profile", Value: "rm2-desk-weekly", Usage: "profile name"},
		&cli.StringFlag{Name: "config", Usage: "override comma-separated config files"},
		&cli.IntFlag{Name: "year", Usage: "planner year"},
		&cli.IntFlag{Name: "from", Usage: "first year for build-range"},
		&cli.IntFlag{Name: "to", Usage: "last year for build-range"},
		&cli.StringFlag{Name: "lang", Aliases: []string{"translation"}, Usage: "translation name"},
		&cli.StringFlag{Name: "country", Usage: "holiday country codes, comma-separated"},
		&cli.StringFlag{Name: "name", Usage: "output PDF name without .pdf"},
		&cli.IntFlag{Name: "passes", Value: 2, Usage: "XeLaTeX passes"},
		&cli.BoolFlag{Name: "refresh-holidays", Usage: "refresh Nager.Date cache"},
	}
}

func runBuild(c *cli.Context, year int, preview bool) error {
	prof, ok := profiles[c.String("profile")]
	if !ok {
		return fmt.Errorf("unknown profile %q", c.String("profile"))
	}

	config := prof.Config
	if c.String("config") != "" {
		config = c.String("config")
	}

	name := c.String("name")
	if name == "" {
		parts := []string{prof.Name}
		if c.String("lang") != "" {
			parts = append(parts, c.String("lang"))
		}
		parts = append(parts, strconv.Itoa(year))
		if preview {
			parts = append([]string{parts[0], "preview"}, parts[1:]...)
		}
		name = strings.Join(parts, ".")
	}

	if err := os.Setenv("PLANNER_YEAR", strconv.Itoa(year)); err != nil {
		return err
	}
	if c.String("lang") != "" {
		if err := os.Setenv("TRANSLATION", c.String("lang")); err != nil {
			return err
		}
	}
	if c.String("country") != "" {
		if err := os.Setenv("PLANNER_PUBLIC_HOLIDAYS_COUNTRY_CODE", c.String("country")); err != nil {
			return err
		}
	}
	if c.Bool("refresh-holidays") {
		if err := os.Setenv("PLANNER_PUBLIC_HOLIDAYS_REFRESH_CACHE", "true"); err != nil {
			return err
		}
	}

	if err := os.MkdirAll("out", 0755); err != nil {
		return fmt.Errorf("create out dir: %w", err)
	}
	if err := ensurePDFDir(); err != nil {
		return err
	}

	args := []string{"plannergen", "--config", config}
	if preview {
		args = append(args, "--preview")
	}

	if err := app.New().RunContext(context.Background(), args); err != nil {
		return fmt.Errorf("generate latex: %w", err)
	}

	if lang := c.String("lang"); lang != "" {
		if err := translate(lang); err != nil {
			return fmt.Errorf("translate latex: %w", err)
		}
	}

	rootTex := filepath.Join("out", app.RootFilename(lastConfig(config)))
	for pass := 1; pass <= c.Int("passes"); pass++ {
		fmt.Printf("XeLaTeX pass %d/%d\n", pass, c.Int("passes"))
		if err := runXeLaTeX(rootTex); err != nil {
			return err
		}
	}

	rootPDF := strings.TrimSuffix(rootTex, ".tex") + ".pdf"
	outputPDF := filepath.Join("pdf", name+".pdf")
	if err := copyFile(rootPDF, outputPDF); err != nil {
		return fmt.Errorf("copy pdf: %w", err)
	}

	fmt.Printf("created %s\n", outputPDF)
	return nil
}

func ensurePDFDir() error {
	if err := os.MkdirAll("pdf", 0755); err != nil {
		return fmt.Errorf("create pdf dir: %w", err)
	}

	gitkeep := filepath.Join("pdf", ".gitkeep")
	file, err := os.OpenFile(gitkeep, os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("create %s: %w", gitkeep, err)
	}

	return file.Close()
}

func runXeLaTeX(rootTex string) error {
	if _, err := exec.LookPath("xelatex"); err != nil {
		return fmt.Errorf("xelatex is required. Run ./install.sh or install TeX Live manually: %w", err)
	}

	cmd := exec.Command("xelatex", "-file-line-error", "-interaction=nonstopmode", "-synctex=1", "-output-directory=./out", rootTex)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("xelatex %s: %w", rootTex, err)
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}

func lastConfig(config string) string {
	parts := strings.Split(config, ",")
	return parts[len(parts)-1]
}

func translate(language string) error {
	language = strings.ToLower(language)
	path := filepath.Join("translations", language+".json")

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("requested translation is not currently supported: %w", err)
	}

	translation := map[string]string{}
	if err := json.Unmarshal(content, &translation); err != nil {
		return err
	}

	fmt.Printf("Translating pdf to %s\n", language)

	if err := translateIfExists("out/annual.tex", translation, annualReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/quarterly.tex", translation, quarterlyReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/monthly.tex", translation, monthlyReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/weekly.tex", translation, weeklyReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/daily.tex", translation, dailyReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/daily_reflect.tex", translation, dailyReflectReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/daily_notes.tex", translation, dailyNotesReplacements); err != nil {
		return err
	}
	if err := translateIfExists("out/notes_indexed.tex", translation, notesIndexedReplacements); err != nil {
		return err
	}

	paths, err := filepath.Glob("out/*.tex")
	if err != nil {
		return err
	}
	for _, path := range paths {
		if err := applyTranslation(path, commonGeneratedReplacements(translation)); err != nil {
			return err
		}
	}

	return nil
}

func translateIfExists(path string, translation map[string]string, replacements func(map[string]string) map[string]string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	return applyTranslation(path, replacements(translation))
}

func applyTranslation(path string, replacements map[string]string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	for from, to := range replacements {
		text = strings.ReplaceAll(text, from, to)
	}

	return os.WriteFile(path, []byte(text), 0600)
}

func addIdentifiers(replacements map[string]string, translation map[string]string, keys []string, fn func(string) string) {
	for _, key := range keys {
		translated, ok := translation[key]
		if !ok {
			continue
		}
		replacements[fn(key)] = fn(translated)
	}
}

func annualReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "{" + s + "}}" })
	addIdentifiers(replacements, translation, months, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "{" + s + "}" })
	return replacements
}

func quarterlyReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "{" + s + "}}" })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "{" + s + "}" })
	return replacements
}

func monthlyReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return s })
	addIdentifiers(replacements, translation, []string{"Week"}, func(s string) string { return "[c]{" + s })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "{" + s + "}" })
	return replacements
}

func weeklyReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, months, func(s string) string { return `\textbf{` + s })
	addIdentifiers(replacements, translation, months, func(s string) string { return " / " + s + "}" })
	addIdentifiers(replacements, translation, []string{"Week"}, func(s string) string { return "}{" + s })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return ", " + s + "}" })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return "{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "{" + s })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return `{\scriptsize ` + s + "}" })
	return replacements
}

func dailyReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"Week"}, func(s string) string { return "}{" + s })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return "}{" + s + "," })
	addIdentifiers(replacements, translation, weekdaysShort, func(s string) string { return "}{" + s + "," })
	addIdentifiers(replacements, translation, []string{"Schedule", "Top priorities"}, func(s string) string { return "{" + s + `\` })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "{" + s + " $" })
	addIdentifiers(replacements, translation, []string{"More", "Reflect"}, func(s string) string { return "{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"All notes"}, func(s string) string { return s })
	return replacements
}

func dailyReflectReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"Week"}, func(s string) string { return "}{" + s })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return "}{" + s + "," })
	addIdentifiers(replacements, translation, weekdaysShort, func(s string) string { return "}{" + s + "," })
	addIdentifiers(replacements, translation, []string{"Things I'm grateful for", "The best thing that happened today", "Daily log"}, func(s string) string { return s })
	addIdentifiers(replacements, translation, []string{"Reflect"}, func(s string) string { return "{" + s + "}" })
	return replacements
}

func dailyNotesReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, months, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"Week"}, func(s string) string { return "}{" + s })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return "}{" + s + "," })
	addIdentifiers(replacements, translation, weekdaysShort, func(s string) string { return "}{" + s + "," })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "{" + s + "}" })
	return replacements
}

func notesIndexedReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, []string{"Notes Index", "Note"}, func(s string) string { return "}{" + s })
	return replacements
}

func commonGeneratedReplacements(translation map[string]string) map[string]string {
	replacements := map[string]string{}
	addIdentifiers(replacements, translation, []string{"Calendar"}, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, monthsShort, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, []string{"Notes"}, func(s string) string { return "}{" + s + "}" })
	addIdentifiers(replacements, translation, months, func(s string) string { return " & " + s + " &" })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return `\textbf{` + s + "}" })
	addIdentifiers(replacements, translation, weekdays, func(s string) string { return " & " + s + " &" })
	addIdentifiers(replacements, translation, weekdaysShort, func(s string) string { return "{" + s + "," })
	addIdentifiers(replacements, translation, weekdaysShort, func(s string) string { return "}{" + s + "," })
	replacements["W & M & T & W & T & F & S & S"] = "T & P & W & Ś & C & P & S & N"
	notes := translation["Notes"]
	if notes == "" {
		notes = "Notes"
	}
	replacements["Notatkas"] = notes
	replacements["{Notatki Index}"] = "{Notes Index}"
	replacements["{Indeks notatek}"] = "{Notes Index}"
	return replacements
}
