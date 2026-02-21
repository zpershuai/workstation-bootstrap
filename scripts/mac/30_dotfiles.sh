#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

log "Linking dotfiles"

safe_link "${ROOT_DIR}/config/zsh/zshrc" "${HOME}/.zshrc"
safe_link "${ROOT_DIR}/config/zsh/zprofile" "${HOME}/.zprofile"
safe_link "${ROOT_DIR}/config/shell/profile" "${HOME}/.profile"
safe_link "${HOME}/.tmux/.tmux.conf" "${HOME}/.tmux.conf"
safe_link "${ROOT_DIR}/config/git/gitconfig" "${HOME}/.gitconfig"
safe_link "${ROOT_DIR}/config/git/gitignore_global" "${HOME}/.gitignore_global"

ensure_dir "${HOME}/.config/git"
safe_link "${ROOT_DIR}/config/git/gitconfig.base" "${HOME}/.config/git/.gitconfig.base"

ensure_dir "${HOME}/.config"
safe_link "${ROOT_DIR}/config/fish" "${HOME}/.config/fish"
safe_link "${ROOT_DIR}/config/ghostty" "${HOME}/.config/ghostty"
safe_link "${ROOT_DIR}/config/karabiner" "${HOME}/.config/karabiner"
safe_link "${ROOT_DIR}/config/uv" "${HOME}/.config/uv"

safe_link "${HOME}/.dotfiles.d/repos/nvim" "${HOME}/.config/nvim"
safe_link "${HOME}/.dotfiles.d/repos/tmux" "${HOME}/.tmux"

safe_link "${ROOT_DIR}/misc/cc-switch" "${HOME}/.cc-switch"
safe_link "${ROOT_DIR}/misc/dotfiles" "${HOME}/.dotfiles"

TMUX_INSTALL="${HOME}/.tmux/install.sh"
if [[ -x "${TMUX_INSTALL}" ]]; then
  if is_dry_run; then
    log "DRY_RUN: would run ${TMUX_INSTALL}"
  else
    log "Running tmux install script"
    bash "${TMUX_INSTALL}"
  fi
fi
