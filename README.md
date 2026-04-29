# interactive-yearly-planner

PDF planner designed for e-ink devices.

This fork adds a minimal reMarkable 2 planner variant and Polish translation support.

See [discussions](https://github.com/kudrykv/latex-yearly-planner/discussions) for available planners and their variations.

### Documentation work in progress
I am planning to write more documentation on how to use it and build it on your own.
Spoiler alert: it won't be easy.
Anyhow, more info on this will come.

### Tackling

I suggest looking at [rubify2](https://github.com/kudrykv/latex-yearly-planner/tree/rubify2) branch.
It is an ongoing refactor, and it can generate MOS template.

## Quick Start Guide
Here are the steps to quickly get the project up and running.

* Note: if you are here just for the planners you can find pre-generated
 planners in [2022-2032 Planners Discussions](https://github.com/kudrykv/latex-yearly-planner/discussions/57).

For the tinkerers, read on.

The following was tested with [POP_OS 22.04.1 LTS](https://pop.system76.com/) under [Virtualbox](https://www.virtualbox.org/) version 6.1

### Install Dependencies
1. [Go Language](https://go.dev/dl/)
2. [LaTex](https://miktex.org/download) & [PDFLaTeX](https://www.latex-project.org/get/)

On Debian/Ubuntu you can install the required packages with:

```bash
sudo apt update && sudo apt install golang-go texlive-xetex texlive-latex-extra texlive-lang-polish
```

3. From the project directory, run the following command after updating
 'PLANNER_YEAR' below. This should generate the PDF in the 'out' directory.
<code>PLANNER_YEAR=2022 \
PASSES=1 \
CFG="cfg/base.yaml,cfg/template_breadcrumb.yaml,cfg/sn_a5x.breadcrumb.default.yaml" \
NAME="sn_a5x.breadcrumb.default" \
./single.sh</code> 

[Source](https://github.com/kudrykv/latex-yearly-planner/discussions/34#discussioncomment-3128344)

4. Check the "out" directory for the 'pdf' planner. To move it to your device
, follow the manufacturer's instructions on how to load a PDF on your device.

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
PLANNER_YEAR=2026 \
PASSES=2 \
CFG="cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.minimal.yaml,cfg/template_months_on_side_minimal.yaml,cfg/rm2.mos.default.yaml" \
NAME="rm2.minimal.en.2026" \
./single.sh
```

Generate the Polish minimal RM2 planner:

```bash
PLANNER_YEAR=2026 \
PASSES=2 \
TRANSLATION=polish \
CFG="cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.minimal.yaml,cfg/template_months_on_side_minimal.yaml,cfg/rm2.mos.default.yaml" \
NAME="rm2.minimal.pl.2026" \
./single.sh
```

The minimal RM2 configuration uses a 24-hour clock and starts weeks on Monday:

```yaml
weekstart: 1
ampmtime: false
```

### Desk weekly reMarkable 2 planner

The desk weekly RM2 variant uses the reMarkable 2 page size in landscape orientation and adds a title page, annual page, and one compact weekly page per week. The weekly layout uses seven vertical day columns, hourly writing lines, and a small month calendar in the header.

Public holidays can be marked by country using Nager.Date. Set `publicholidays.countrycodes` in the planner config, or override it with `PLANNER_PUBLIC_HOLIDAYS_COUNTRY_CODE` using one code or a comma-separated list like `PL,DE`. The desk weekly RM2 config uses `PL` by default, so public holidays are downloaded from `https://date.nager.at/api/v3/PublicHolidays/<year>/<country>` and shown in red.

Holiday names from `localName` can be enabled in weekly day headers with `publicholidays.shownames: true` or `PLANNER_PUBLIC_HOLIDAYS_SHOW_NAMES=true`. This is disabled in the desk weekly RM2 config to avoid long holiday names wrapping in narrow columns.

Generate the Polish desk weekly RM2 planner:

```bash
PLANNER_YEAR=2026 \
PASSES=2 \
TRANSLATION=polish \
CFG="cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.desk.weekly.yaml,cfg/template_desk_weekly_rm2.yaml,cfg/rm2.mos.default.yaml" \
NAME="rm2.desk.weekly.pl.2026" \
./single.sh
```

### Polish translation

Polish translations are available through:

```text
translations/polish.json
```

Enable them with:

```bash
TRANSLATION=polish
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
TRANSLATION=polish ./preview.sh 2026
```

3. Generate the final PDF:

```bash
TRANSLATION=polish ./build.sh 2026
```

4. Use `single.sh` directly only when you need a custom config stack:

```bash
PLANNER_YEAR=2026 \
PASSES=2 \
TRANSLATION=polish \
CFG="cfg/base.yaml,cfg/rm2.base.yaml,cfg/rm2.minimal.yaml,cfg/template_months_on_side_minimal.yaml,cfg/rm2.mos.default.yaml" \
NAME="rm2.minimal.pl.2026" \
./single.sh
```

Generated PDF files are ignored by Git through `*.pdf` in `.gitignore`.

### Scripts

`install.sh` checks whether `go`, `python3`, and `xelatex` are available. On Debian/Ubuntu systems it installs missing dependencies with `apt-get`.

`single.sh` is the low-level build entrypoint. It runs the Go planner generator, optionally runs `translate.py`, compiles the generated LaTeX with XeLaTeX, and copies the final PDF to `NAME.pdf`.

`build.sh` is the default final-build wrapper. It defaults to the minimal RM2 config stack and writes progress output through `parser.py`. Override `CONFIG_FILES`, `NAME`, `PASSES`, or `TRANSLATION` when needed.

`preview.sh` is the same wrapper as `build.sh`, but passes `PREVIEW=1` to the generator. Use it for faster iteration while editing templates or config.

`parser.py` reads XeLaTeX output and prints compact page progress instead of streaming the full LaTeX log.

`release.sh` builds the release matrix for the current and next year. It now includes `rm2.minimal.default` and `rm2.minimal.pl.default` variants in addition to the upstream planner variants.

### Alternative install

Instead of installing the dependencies manually, this repository is defined as a Nix flake which specifies fixed versions of all the required dependencies. 

1. [Install Nix](https://nixos.org/download.html)
2. Build a planner pdf using `nix build`
3. Or, if you want to develop the code, enter a shell with all the dependencies present using `nix develop`
   
# Preview examples
<img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/01_annual.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/02_quarter.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/03_month.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/04_week.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/05_day.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/06_day_notes.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/07_day_reflect.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/08_todos_index.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/09_todos_page.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/10_notes_index.png" width="419"><img src="https://github.com/kudrykv/latex-yearly-planner/blob/main/examples/pictures/sn_a5x.planner/11_notes_page.png" width="419">
