#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

if ! command -v brew >/dev/null 2>&1; then
  log "Homebrew not found; installing"
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

BREWFILE="${ROOT_DIR}/brew/Brewfile"
if [[ -f "${BREWFILE}" ]]; then
  log "Installing Brewfile packages"
  brew bundle --file "${BREWFILE}"
else
  log "No Brewfile found at ${BREWFILE}"
fi
