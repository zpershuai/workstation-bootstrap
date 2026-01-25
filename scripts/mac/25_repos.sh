#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

MANIFEST="${ROOT_DIR}/repos/repos.lock"
if [[ ! -f "${MANIFEST}" ]]; then
  log "No repo manifest at ${MANIFEST}"
  exit 0
fi

need_cmd git

log "Syncing external config repos"
while read -r name url dest ref; do
  [[ -z "${name}" || "${name}" == \#* ]] && continue

  if [[ -z "${url}" || -z "${dest}" ]]; then
    die "Invalid entry: ${name} (expected: name url dest [ref])"
  fi

  dest="${dest/#\~/$HOME}"
  if is_dry_run; then
    log "DRY_RUN: would ensure dir $(dirname "${dest}")"
  else
    ensure_dir "$(dirname "${dest}")"
  fi

  if is_dry_run; then
    log "DRY_RUN: would check access to ${url}"
    if ! git ls-remote --heads "${url}" >/dev/null 2>&1; then
      die "DRY_RUN access check failed for ${name} (${url}). Check SSH/HTTPS auth and network."
    fi
    if [[ -d "${dest}/.git" ]]; then
      log "DRY_RUN: would fetch/checkout/pull in ${dest}"
    else
      log "DRY_RUN: would clone ${name} -> ${dest}"
    fi
    continue
  fi

  if [[ ! -d "${dest}/.git" ]]; then
    log "Cloning ${name} -> ${dest}"
    if ! git clone "${url}" "${dest}"; then
      die "Clone failed for ${name}. Check SSH access (e.g., ssh -T git@github.com) or network."
    fi
  fi

  if ! git -C "${dest}" fetch --all --tags; then
    die "Fetch failed for ${name}. Check SSH access or network."
  fi

  if [[ -n "${ref:-}" ]]; then
    if ! git -C "${dest}" checkout "${ref}"; then
      die "Checkout failed for ${name} (${ref}). Verify ref exists."
    fi
  fi

  if [[ -n "${ref:-}" ]] && git -C "${dest}" show-ref --quiet "refs/heads/${ref}"; then
    git -C "${dest}" pull --ff-only
  elif [[ -z "${ref:-}" ]]; then
    git -C "${dest}" pull --ff-only
  fi
done < "${MANIFEST}"
