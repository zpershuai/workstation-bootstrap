#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${HOME}/.dotfiles_backup/$(date +%Y%m%d-%H%M%S)"

log() {
  printf '[dotfiles] %s\n' "$*"
}

need_cmd() {
  local cmd="$1"
  if ! command -v "${cmd}" >/dev/null 2>&1; then
    die "Missing required command: ${cmd}"
  fi
}

die() {
  printf '[dotfiles] ERROR: %s\n' "$*" >&2
  exit 1
}

run_module() {
  local module="$1"
  if [[ -f "${module}" ]]; then
    log "Running ${module}"
    bash "${module}"
  else
    log "Skipping missing module ${module}"
  fi
}

ensure_dir() {
  local dir="$1"
  [[ -d "${dir}" ]] || mkdir -p "${dir}"
}

backup_path() {
  local path="$1"
  if [[ -e "${path}" || -L "${path}" ]]; then
    ensure_dir "${BACKUP_DIR}"
    mv "${path}" "${BACKUP_DIR}/"
  fi
}

safe_link() {
  local src="$1"
  local dest="$2"

  if [[ ! -e "${src}" ]]; then
    log "Skip missing source ${src}"
    return 0
  fi

  if [[ -L "${dest}" && "$(readlink "${dest}")" == "${src}" ]]; then
    log "Link already correct: ${dest}"
    return 0
  fi

  backup_path "${dest}"
  ensure_dir "$(dirname "${dest}")"
  ln -s "${src}" "${dest}"
  log "Linked ${dest} -> ${src}"
}
