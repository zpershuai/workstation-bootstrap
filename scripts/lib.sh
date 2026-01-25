#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${HOME}/.dotfiles_backup/$(date +%Y%m%d-%H%M%S)"
DRY_RUN="${DRY_RUN:-0}"

log() {
  printf '[dotfiles] %s\n' "$*"
}

is_dry_run() {
  [[ "${DRY_RUN}" == "1" ]]
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
  if [[ -d "${dir}" ]]; then
    return 0
  fi
  if is_dry_run; then
    log "DRY_RUN: would create dir ${dir}"
    return 0
  fi
  mkdir -p "${dir}"
}

backup_path() {
  local path="$1"
  if [[ -e "${path}" || -L "${path}" ]]; then
    if is_dry_run; then
      log "DRY_RUN: would backup ${path} -> ${BACKUP_DIR}/"
      return 0
    fi
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

  if is_dry_run; then
    log "DRY_RUN: would link ${dest} -> ${src}"
    return 0
  fi

  backup_path "${dest}"
  ensure_dir "$(dirname "${dest}")"
  ln -s "${src}" "${dest}"
  log "Linked ${dest} -> ${src}"
}
