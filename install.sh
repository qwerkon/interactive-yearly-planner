#!/usr/bin/env bash

set -euo pipefail

missing_commands=()

for cmd in go xelatex; do
  if ! command -v "${cmd}" >/dev/null 2>&1; then
    missing_commands+=("${cmd}")
  fi
done

if [ "${#missing_commands[@]}" -eq 0 ]; then
  echo "All required dependencies are installed."
  exit 0
fi

if ! command -v apt-get >/dev/null 2>&1; then
  echo "Missing commands: ${missing_commands[*]}" >&2
  echo "Automatic installation is only supported on systems with apt-get." >&2
  exit 1
fi

packages=()

if ! command -v go >/dev/null 2>&1; then
  packages+=(golang-go)
fi

if ! command -v xelatex >/dev/null 2>&1; then
  packages+=(texlive-xetex texlive-latex-extra texlive-fonts-recommended texlive-lang-polish)
fi

if [ "${#packages[@]}" -eq 0 ]; then
  echo "No apt packages need to be installed."
  exit 0
fi

echo "Installing missing dependencies: ${packages[*]}"
sudo apt-get update
sudo apt-get install -y "${packages[@]}"

echo "Installed dependencies successfully."
