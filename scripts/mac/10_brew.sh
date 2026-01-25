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

  get_bundle_ids() {
    local cask_name="$1"
    if ! command -v python3 >/dev/null 2>&1; then
      return 0
    fi
    local json
    json="$(brew info --cask --json=v2 "${cask_name}" 2>/dev/null || true)"
    if [[ -z "${json}" ]]; then
      return 0
    fi
    JSON_INPUT="${json}" python3 - <<'PY'
import json
import os
raw = os.environ.get("JSON_INPUT", "").strip()
if not raw:
    sys.exit(0)
data = json.loads(raw)
casks = data.get("casks", [])
if not casks:
    sys.exit(0)
entry = casks[0]
bundle_id = entry.get("bundle_id", [])
if isinstance(bundle_id, str):
    bundle_id = [bundle_id]
for bid in bundle_id:
    if bid:
        print(bid)
PY
  }

  get_app_names_from_json() {
    local cask_name="$1"
    if ! command -v python3 >/dev/null 2>&1; then
      return 0
    fi
    local json
    json="$(brew info --cask --json=v2 "${cask_name}" 2>/dev/null || true)"
    if [[ -z "${json}" ]]; then
      return 0
    fi
    JSON_INPUT="${json}" python3 - <<'PY'
import json
import os
raw = os.environ.get("JSON_INPUT", "").strip()
if not raw:
    raise SystemExit(0)
data = json.loads(raw)
casks = data.get("casks", [])
if not casks:
    raise SystemExit(0)
entry = casks[0]
artifacts = entry.get("artifacts", [])
for artifact in artifacts:
    if "app" in artifact:
        for app in artifact["app"]:
            print(app)
PY
  }

  has_bundle_id_installed() {
    local cask_name="$1"
    if ! command -v mdfind >/dev/null 2>&1; then
      return 1
    fi
    local found=1
    while IFS= read -r bid; do
      [[ -z "${bid}" ]] && continue
      if mdfind "kMDItemCFBundleIdentifier == '${bid}'" | grep -q .; then
        found=0
        break
      fi
    done < <(get_bundle_ids "${cask_name}")
    return "${found}"
  }

  cask_skip_list=()
  while read -r cask; do
    [[ -z "${cask}" ]] && continue

    if brew list --cask "${cask}" >/dev/null 2>&1; then
      cask_skip_list+=("${cask}")
      continue
    fi

    app_names="$(brew info --cask "${cask}" 2>/dev/null | awk '/\\(App\\)$/ {print $1}')"
    if [[ -z "${app_names}" ]]; then
      app_names="$(get_app_names_from_json "${cask}")"
    fi
    if [[ -n "${app_names}" ]]; then
      while IFS= read -r app; do
        if [[ -d "/Applications/${app}" || -d "${HOME}/Applications/${app}" ]]; then
          log "Found ${app} in /Applications; skipping cask ${cask}"
          cask_skip_list+=("${cask}")
          break
        fi
      done <<< "${app_names}"
    fi

    if has_bundle_id_installed "${cask}"; then
      log "Found bundle id installed; skipping cask ${cask}"
      cask_skip_list+=("${cask}")
      continue
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
