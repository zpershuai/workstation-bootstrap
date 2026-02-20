#!/bin/bash
# Backup Manico configuration to XML format for version control
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log() { echo "[manico-backup] $*"; }

backup_plist() {
  local name="$1"
  local src="${HOME}/Library/Preferences/${name}.plist"
  local dest="${SCRIPT_DIR}/${name}.plist.xml"

  if [[ ! -f "$src" ]]; then
    log "Warning: $src not found, skipping"
    return 0
  fi

  # Convert binary plist to XML format
  plutil -convert xml1 "$src" -o "$dest"
  log "Backed up: $name"
}

log "Backing up Manico configuration..."

backup_plist "com.lintie.manico"
backup_plist "com.lintie.manico.helper"

log "Backup complete. Files saved to:"
log "  ${SCRIPT_DIR}/"
