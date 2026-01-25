#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

MANIFEST="${ROOT_DIR}/repos/repos.lock"
if [[ ! -f "${MANIFEST}" ]]; then
  log "No repo manifest at ${MANIFEST}"
  exit 0
fi

if ! command -v git >/dev/null 2>&1; then
  die "git is required to sync external repos"
fi

log "Syncing external config repos"
while read -r name url dest ref; do
  [[ -z "${name}" || "${name}" == \#* ]] && continue

  if [[ -z "${url}" || -z "${dest}" ]]; then
    log "Invalid entry for ${name}; expected: name url dest [ref]"
    continue
  fi

  dest="${dest/#\~/$HOME}"
  ensure_dir "$(dirname "${dest}")"

  if [[ ! -d "${dest}/.git" ]]; then
    log "Cloning ${name} -> ${dest}"
    git clone "${url}" "${dest}"
  fi

  git -C "${dest}" fetch --all --tags
  if [[ -n "${ref:-}" ]]; then
    git -C "${dest}" checkout "${ref}"
  else
    git -C "${dest}" checkout "$(git -C "${dest}" symbolic-ref --short HEAD)"
    git -C "${dest}" pull --ff-only
  fi
done < "${MANIFEST}"
