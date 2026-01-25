#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

log "Prerequisites: Xcode CLT, git, curl"

if ! xcode-select -p >/dev/null 2>&1; then
  log "Xcode CLT missing; triggering install"
  xcode-select --install || true
fi

if ! command -v git >/dev/null 2>&1; then
  log "git not found; install via Xcode CLT or Homebrew"
fi

if ! command -v curl >/dev/null 2>&1; then
  log "curl not found; install via Xcode CLT or Homebrew"
fi
