#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

log() {
  printf '[backup] %s\n' "$*"
}

BACKUP_DIR="${HOME}/.dotfiles_backup/$(date +%Y%m%d-%H%M%S)"
log "Backup dir: ${BACKUP_DIR}"
mkdir -p "${BACKUP_DIR}"

paths=(
  "${HOME}/.zshrc"
  "${HOME}/.zprofile"
  "${HOME}/.profile"
  "${HOME}/.tmux.conf"
  "${HOME}/.gitconfig"
  "${HOME}/.config/git/.gitconfig.base"
  "${HOME}/.config/fish"
  "${HOME}/.config/iterm2"
  "${HOME}/.config/ghostty"
  "${HOME}/.config/karabiner"
  "${HOME}/.config/uv"
  "${HOME}/.config/nvim"
  "${HOME}/.tmux"
  "${HOME}/.cc-switch"
  "${HOME}/.dotfiles"
)

for p in "${paths[@]}"; do
  if [[ -e "${p}" || -L "${p}" ]]; then
    log "Moving ${p} -> ${BACKUP_DIR}"
    mv "${p}" "${BACKUP_DIR}/"
  else
    log "Skip missing ${p}"
  fi
done

log "Done."
