# interactive-yearly-planner

PDF planner generator for e-ink devices, focused on reMarkable 2 layouts. It builds hyperlinked yearly planners from YAML profiles, with RM2 minimal, desk weekly, and desk monthly-weekly variants.

The Go CLI can generate previews, final PDFs, all profiles for a selected year, or one profile across a range of years. The planner supports Polish translation, public holidays from Nager.Date with local cache, custom YAML events, ICS imports, configurable RM2 weekly/monthly layouts, and LaTeX snapshot tests for regression coverage.

## Quick Start Guide
### Install Dependencies
1. [Go Language](https://go.dev/dl/)
2. [LaTex](https://miktex.org/download) & [PDFLaTeX](https://www.latex-project.org/get/)

On Debian/Ubuntu you can install the required packages with:

```bash
sudo apt update && sudo apt install golang-go texlive-xetex texlive-latex-extra texlive-lang-polish
```

3. From the project directory, generate a PDF with the Go CLI:

```bash
go run cmd/planner/main.go build --profile rm2-desk-weekly --year 2026 --lang polish --country PL
```

4. Check the `pdf/` directory for the planner. To move it to your device, follow the manufacturer's instructions on how to load a PDF on your device.

If you encounter any problems related to '.sty' files you likely need to
 install some Latex related dependencies. Copying the error and search using
  your favorite search engine should get you on track to resolving the
   dependency issues. All the best!

### Minimal reMarkable 2 planner

The minimal reMarkable 2 variant is language-independent and is configured with:

```text
cfg/rm2.base.yaml
cfg/rm2.minimal.yaml
cfg/template_months_on_side_minimal.yaml
cfg/rm2.mos.default.yaml
```

It includes only:

```text
title
annual
daily
```

It excludes quarterly, monthly, weekly, daily reflection, daily notes, and indexed notes pages. Navigation links to disabled sections are hidden automatically, so the generated PDF does not contain dead links to removed views.

Generate the English minimal RM2 planner:

```bash
go run cmd/planner/main.go build --profile rm2-minimal --year 2026
```

Generate the Polish minimal RM2 planner:

```bash
go run cmd/planner/main.go build --profile rm2-minimal --year 2026 --lang polish --country PL
```

The minimal RM2 configuration uses a 24-hour clock and starts weeks on Monday:

```yaml
weekstart: 1
ampmtime: false
```

### Desk reMarkable 2 planner

The desk weekly RM2 variant uses the reMarkable 2 page size in landscape orientation and adds a title page, annual page, and one compact weekly page per week. The weekly layout uses seven vertical day columns, hourly writing lines, and a small month calendar in the header.

The monthly-weekly RM2 profile adds one landscape month page before the weekly pages. The month page uses a larger landscape grid, weekday/week links, a compact header, and a side panel for navigation and color legend:

```bash
go run cmd/planner/main.go preview --profile rm2-desk-monthly-weekly --year 2026 --lang polish --country PL
```

Monthly layout details can be tuned with `layout.deskmonthly`:

```yaml
layout:
  deskmonthly:
    headerfont: 15
    subheaderfont: 7
    sidepanelwidth: 4.2cm
    legendfont: 5.5
    showlegend: true
    showsidepanel: true
```

Public holidays can be marked by country using Nager.Date. Set `publicholidays.countrycodes` in the planner config, or override it with `PLANNER_PUBLIC_HOLIDAYS_COUNTRY_CODE` using one code or a comma-separated list like `PL,DE`. The desk weekly RM2 config uses `PL` by default, so public holidays are downloaded from `https://date.nager.at/api/v3/PublicHolidays/<year>/<country>` and shown in red.

Regional Nager.Date holidays can be filtered with `publicholidays.subdivisions`, using Nager.Date county/subdivision codes such as `US-CA` or `DE-BY`. When subdivisions are configured, global holidays are still included and regional holidays are included only when one of their `counties` matches. Holiday `types` can be narrowed with `includetypes` and removed with `excludetypes`; exclusions win over inclusions.

Holiday names from `localName` can be enabled in weekly day headers with `publicholidays.shownames: true` or `PLANNER_PUBLIC_HOLIDAYS_SHOW_NAMES=true`. This is disabled in the desk weekly RM2 config to avoid long holiday names wrapping in narrow columns.

Desk weekly layout details can be tuned without editing LaTeX:

```yaml
layout:
  deskweekly:
    starthour: 6
    endhour: 22
    hourlineheight: 4.25mm
    columnpadding: 1mm
    headerpadding: 1mm
    minicalendarwidth: 4.25cm
    daynumberfont: 18
    weekdayfont: 6.5
    hourfont: 5.5
    holidaynamefont: 4.3
    holidaymarker: \textbullet
    eventmarker: '*'
    showholidaymarker: true
    showholidaylegend: true
```

Holiday data is cached when `publicholidays.cachedir` is set. Use `publicholidays.refreshcache: true` or `PLANNER_PUBLIC_HOLIDAYS_REFRESH_CACHE=true` to refresh downloaded Nager.Date files. You can add local holidays, separate YAML event files, or all-day ICS events:

```yaml
publicholidays:
  countrycodes: [PL, DE]
  subdivisions: [DE-BY]
  includetypes: [Public]
  excludetypes: [Optional]
  cachedir: cache/holidays
  custom:
    - date: "2026-05-04"
      localName: Urlop
      name: Time off
      shortName: Off
      types: [Public]
    - date: "2026-06-02"
      localName: School event
      name: School event
      shortName: School
      types: [Event]
  icsfiles:
    - events.ics

events:
  files:
    - cfg/events.example.yaml
  custom:
    - date: "2026-06-02"
      localName: School event
      name: School event
      shortName: School
      types: [Event]
```

Public holidays use `layout.colors.publicholiday`; custom/ICS entries with `types: [Event]` use `layout.colors.event`. Weekly headers can show a marker and a short label using `shortName`, then `localName`, then `name` as fallback.

ICS import supports `VEVENT` entries with `DTSTART`, `DTEND`, `SUMMARY`, `CATEGORIES`, and basic `RRULE` recurrence. All-day and timezone-aware date-time starts are accepted, including `TZID` values such as `Europe/Warsaw`. Multi-day events are expanded into one planner marker per day. Supported recurrence frequencies are `DAILY`, `WEEKLY`, `MONTHLY`, and `YEARLY`, with `COUNT`, `UNTIL`, and `INTERVAL`.

An event file uses this shape:

```yaml
events:
  - date: "2026-05-04"
    localName: Urlop
    name: Time off
    shortName: Off
    types: [Event]
```

Config validation rejects empty page lists, unknown render block functions, render blocks without templates, invalid years/week starts, malformed `YYYY-MM-DD` dates, empty event file paths, and invalid two-letter country codes.

Generate the Polish desk weekly RM2 planner with the Go CLI:

```bash
go run cmd/planner/main.go build --profile rm2-desk-weekly --year 2026 --lang polish --country PL
```

The CLI runs the Go planner generator, applies translations, compiles with XeLaTeX, and copies the final PDF without using the Bash build wrappers.

Use a custom config stack when needed:

```bash
go run cmd/planner/main.go build \
  --config "cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.desk.weekly.yaml,cfg/template_desk_weekly_rm2.yaml,cfg/rm2.mos.default.yaml" \
  --year 2026 \
  --lang polish \
  --country PL \
  --name rm2.desk.weekly.pl.2026
```

Available profiles:

```bash
go run cmd/planner/main.go profiles
```

Build a range of years:

```bash
go run cmd/planner/main.go build-range --profile rm2-desk-weekly --from 2026 --to 2030 --lang polish --country PL
```

Build every available profile for one year:

```bash
go run cmd/planner/main.go build-all --year 2026 --lang polish --country PL
```

### Polish translation

Polish translations are available through:

```text
translations/polish.json
```

Enable them with:

```bash
go run cmd/planner/main.go build --profile rm2-desk-weekly --year 2026 --lang polish
```

The translation covers month names, weekday names, short labels, calendar labels, schedule labels, notes labels, and generated Go labels that appear after template rendering.

### Workflow

For everyday use, the intended flow is:

1. Install or verify local dependencies:

```bash
./install.sh
```

2. Generate a quick preview while changing layouts or translations:

```bash
go run cmd/planner/main.go preview --profile rm2-desk-weekly --year 2026 --lang polish --country PL
```

3. Generate the final PDF:

```bash
go run cmd/planner/main.go build --profile rm2-desk-weekly --year 2026 --lang polish --country PL
```

4. Build several years:

```bash
go run cmd/planner/main.go build-range --profile rm2-desk-weekly --from 2026 --to 2030 --lang polish --country PL
```

5. Build all profiles for one year:

```bash
go run cmd/planner/main.go build-all --year 2026 --lang polish --country PL
```

Generated PDF files from the Go CLI are written to `pdf/`. The directory is kept in Git with `pdf/.gitkeep`, while generated files inside it are ignored.

### Scripts

`install.sh` checks whether `go` and `xelatex` are available. On Debian/Ubuntu systems it installs missing dependencies with `apt-get`.

`cmd/planner` is the primary build entrypoint. It runs the Go planner generator, applies translations, compiles the generated LaTeX with XeLaTeX, and copies the final PDF to `--name` or the profile-derived output name.
