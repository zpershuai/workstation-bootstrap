#!/bin/bash
# Restore Manico configuration from XML to binary plist
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log() { echo "[manico-restore] $*"; }

die() { echo "[manico-restore] Error: $*" >&2; exit 1; }

restore_plist() {
  local name="$1"
  local src="${SCRIPT_DIR}/${name}.plist.xml"
  local dest="${HOME}/Library/Preferences/${name}.plist"

  if [[ ! -f "$src" ]]; then
    log "Warning: $src not found, skipping"
    return 0
  fi

  # Backup existing if present
  if [[ -f "$dest" ]]; then
    cp "$dest" "${dest}.backup.$(date +%Y%m%d%H%M%S)"
    log "Existing config backed up"
  fi

  # Convert XML to binary plist
  plutil -convert binary1 "$src" -o "$dest"

  # Set correct permissions
  chmod 600 "$dest"

  log "Restored: $name"
}

log "Restoring Manico configuration..."

restore_plist "com.lintie.manico"
restore_plist "com.lintie.manico.helper"

# Restart Manico if running
if pgrep -q "Manico"; then
  log "Restarting Manico to apply changes..."
  killall "Manico" 2>/dev/null || true
  sleep 1
  open -a "Manico"
fi

log "Restore complete."
