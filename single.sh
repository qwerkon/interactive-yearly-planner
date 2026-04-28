#!/usr/bin/env bash

set -eo pipefail

if [ -z "${CFG:-}" ]; then
  echo "CFG is required. Pass a comma-separated list of YAML config files." >&2
  exit 1
fi

if [ -z "${PLANNERGEN_BINARY:-}" ]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "go is required. Run ./install.sh or install Go manually." >&2
    exit 1
  fi

  GO_CMD=(go run cmd/plannergen/plannergen.go)
else
  GO_CMD=("$PLANNERGEN_BINARY")
  echo "Building using plannergen binary at \"${PLANNERGEN_BINARY}\""
fi

if ! command -v xelatex >/dev/null 2>&1; then
  echo "xelatex is required. Run ./install.sh or install TeX Live manually." >&2
  exit 1
fi

if [ -z "${PREVIEW:-}" ]; then
  "${GO_CMD[@]}" --config "${CFG}"
else
  "${GO_CMD[@]}" --preview --config "${CFG}"
fi

last_cfg="${CFG##*,}"
nakedname="${last_cfg##*/}"
nakedname="${nakedname%.yaml}"

if [ -n "${TRANSLATION:-}" ]; then
  if ! command -v python3 >/dev/null 2>&1; then
    echo "python3 is required for translations. Run ./install.sh or install Python manually." >&2
    exit 1
  fi

  python3 translate.py "${TRANSLATION}"
fi

_passes=(1)

if [[ -n "${PASSES:-}" ]]; then
  # shellcheck disable=SC2207
  _passes=($(seq 1 "${PASSES}"))
fi

for _ in "${_passes[@]}"; do
  xelatex \
    -file-line-error \
    -interaction=nonstopmode \
    -synctex=1 \
    -output-directory=./out \
    "out/${nakedname}.tex"
done

if [ -n "${NAME:-}" ]; then
  cp "out/${nakedname}.pdf" "${NAME}.pdf"
  echo "created ${NAME}.pdf"
else
  cp "out/${nakedname}.pdf" "${nakedname}.pdf"
  echo "created ${nakedname}.pdf"
fi
