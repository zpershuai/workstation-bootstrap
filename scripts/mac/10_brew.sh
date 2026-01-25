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

  cask_skip_list=()
  while read -r cask; do
    [[ -z "${cask}" ]] && continue

    if brew list --cask "${cask}" >/dev/null 2>&1; then
      cask_skip_list+=("${cask}")
      continue
    fi

    app_names="$(brew info --cask "${cask}" 2>/dev/null | awk '/\\(App\\)$/ {print $1}')"
    if [[ -n "${app_names}" ]]; then
      while IFS= read -r app; do
        if [[ -d "/Applications/${app}" || -d "${HOME}/Applications/${app}" ]]; then
          log "Found ${app} in /Applications; skipping cask ${cask}"
          cask_skip_list+=("${cask}")
          break
        fi
      done <<< "${app_names}"
    fi
  done < <(awk -F'"' '/^cask /{print $2}' "${BREWFILE}")

  if [[ "${DRY_RUN:-0}" == "1" ]]; then
    log "DRY_RUN: would run brew bundle --file ${BREWFILE}"
    if [[ ${#cask_skip_list[@]} -gt 0 ]]; then
      log "DRY_RUN: would skip casks: ${cask_skip_list[*]}"
    fi
    exit 0
  fi

  if [[ ${#cask_skip_list[@]} -gt 0 ]]; then
    HOMEBREW_BUNDLE_CASK_SKIP="$(printf '%s ' "${cask_skip_list[@]}" | sed 's/ *$//')" \
      brew bundle --file "${BREWFILE}"
  else
    brew bundle --file "${BREWFILE}"
  fi
else
  log "No Brewfile found at ${BREWFILE}"
fi
