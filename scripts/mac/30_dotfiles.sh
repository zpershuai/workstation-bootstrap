#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

log "Linking dotfiles"

safe_link "${ROOT_DIR}/config/zsh/.zshrc" "${HOME}/.zshrc"
safe_link "${ROOT_DIR}/config/tmux/.tmux.conf" "${HOME}/.tmux.conf"
safe_link "${ROOT_DIR}/config/git/.gitconfig" "${HOME}/.gitconfig"

ensure_dir "${HOME}/.config"
if [[ -d "${ROOT_DIR}/config/nvim" ]]; then
  safe_link "${ROOT_DIR}/config/nvim" "${HOME}/.config/nvim"
fi
