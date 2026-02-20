#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

log() { printf '[check] %s\n' "$*"; }
warn() { printf '[warn] %s\n' "$*"; }
err() { printf '[error] %s\n' "$*"; }

errors=0

expect_link() {
  local src="$1"
  local dest="$2"
  if [[ ! -L "${dest}" ]]; then
    err "Missing symlink: ${dest}"
    errors=$((errors + 1))
    return 0
  fi
  local target
  target="$(readlink "${dest}")"
  if [[ "${target}" != "${src}" ]]; then
    err "Link mismatch: ${dest} -> ${target} (expected ${src})"
    errors=$((errors + 1))
  else
    log "OK: ${dest} -> ${src}"
  fi
}

log "Checking dotfile links"
expect_link "${HOME}/.dotfiles.d/repos/nvim" "${HOME}/.config/nvim"
expect_link "${HOME}/.dotfiles.d/repos/tmux" "${HOME}/.tmux"
expect_link "${ROOT_DIR}/misc/dotfiles" "${HOME}/.dotfiles"
expect_link "${ROOT_DIR}/misc/cc-switch" "${HOME}/.cc-switch"

expect_link "${ROOT_DIR}/config/zsh/zshrc" "${HOME}/.zshrc"
expect_link "${ROOT_DIR}/config/zsh/zprofile" "${HOME}/.zprofile"
expect_link "${ROOT_DIR}/config/shell/profile" "${HOME}/.profile"
expect_link "${HOME}/.tmux/.tmux.conf" "${HOME}/.tmux.conf"
expect_link "${ROOT_DIR}/config/git/gitconfig" "${HOME}/.gitconfig"
expect_link "${ROOT_DIR}/config/git/gitconfig.base" "${HOME}/.config/git/.gitconfig.base"

log "Checking oh-my-zsh"
if [[ ! -d "${HOME}/.oh-my-zsh" ]]; then
  warn "Missing ~/.oh-my-zsh"
fi

log "Checking fonts"
if [[ -f "${HOME}/Library/Fonts/SauceCodeProNerdFontCompleteNerdFontNerdFont-Regular.ttf" ]]; then
  log "OK: SauceCodePro Nerd Font installed"
else
  warn "Missing font in ~/Library/Fonts"
fi

log "Checking external repos"
for repo in nvim tmux; do
  if [[ -d "${HOME}/.dotfiles.d/repos/${repo}/.git" ]]; then
    log "OK: repo ${repo} present"
  else
    warn "Missing repo: ${HOME}/.dotfiles.d/repos/${repo}"
  fi
done

if [[ ${errors} -gt 0 ]]; then
  err "Health check failed with ${errors} error(s)"
  exit 1
fi

log "Environment looks good"
