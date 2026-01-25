#!/usr/bin/env bash
set -euo pipefail

# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"

FONT_SRC="${ROOT_DIR}/misc/fonts/SauceCodeProNerdFontCompleteNerdFontNerdFont-Regular.ttf"
FONT_DEST_DIR="${HOME}/Library/Fonts"

if [[ ! -f "${FONT_SRC}" ]]; then
  die "Font not found in repo: ${FONT_SRC}"
fi

ensure_dir "${FONT_DEST_DIR}"
if [[ -f "${FONT_DEST_DIR}/$(basename "${FONT_SRC}")" ]]; then
  log "Font already installed: $(basename "${FONT_SRC}")"
  exit 0
fi

log "Installing font: $(basename "${FONT_SRC}")"
cp "${FONT_SRC}" "${FONT_DEST_DIR}/"
