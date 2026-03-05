#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

log "Prerequisites: Xcode CLT, Homebrew, git, curl, oh-my-zsh"

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

if ! command -v brew >/dev/null 2>&1; then
  log "Homebrew not found; installing"
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  # Add Homebrew to PATH for the current session if it was just installed
  if [[ -x /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
  elif [[ -x /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"
  fi
fi

if ! command -v node >/dev/null 2>&1; then
  log "Node.js not found; installing via Homebrew"
  brew install node
fi

if [[ ! -d "${HOME}/.oh-my-zsh" ]]; then
  if command -v curl >/dev/null 2>&1; then
    log "oh-my-zsh not found; installing"
    KEEP_ZSHRC=yes RUNZSH=no CHSH=no sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
  else
    log "oh-my-zsh missing; install via: sh -c \"$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)\""
  fi
fi
