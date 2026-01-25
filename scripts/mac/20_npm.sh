#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

if ! command -v npm >/dev/null 2>&1; then
  if command -v brew >/dev/null 2>&1; then
    log "npm not found; installing node via Homebrew"
    brew install node
  else
    log "npm not found and Homebrew missing; install Node.js or Homebrew first"
    exit 0
  fi
fi

PKG_FILE="${ROOT_DIR}/npm/packages.txt"
if [[ -f "${PKG_FILE}" ]]; then
  log "Installing global npm packages"
  grep -vE '^[[:space:]]*#' "${PKG_FILE}" | sed '/^[[:space:]]*$/d' | xargs -n 1 npm install -g
else
  log "No npm package list found at ${PKG_FILE}"
fi
