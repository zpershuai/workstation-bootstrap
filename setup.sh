#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# shellcheck source=./scripts/lib.sh
source "${ROOT_DIR}/scripts/lib.sh"

OS_NAME="$(uname -s)"

case "${OS_NAME}" in
  Darwin)
    run_module "${ROOT_DIR}/scripts/mac/00_prereq.sh"
    run_module "${ROOT_DIR}/scripts/mac/10_brew.sh"
    run_module "${ROOT_DIR}/scripts/mac/20_npm.sh"
    run_module "${ROOT_DIR}/scripts/mac/25_repos.sh"
    run_module "${ROOT_DIR}/scripts/mac/30_dotfiles.sh"
    ;;
  Linux)
    run_module "${ROOT_DIR}/scripts/ubuntu/00_prereq.sh"
    run_module "${ROOT_DIR}/scripts/ubuntu/10_apt.sh"
    run_module "${ROOT_DIR}/scripts/ubuntu/20_npm.sh"
    run_module "${ROOT_DIR}/scripts/ubuntu/25_repos.sh"
    run_module "${ROOT_DIR}/scripts/ubuntu/30_dotfiles.sh"
    ;;
  *)
    die "Unsupported OS: ${OS_NAME}"
    ;;
esac

log "Done."
