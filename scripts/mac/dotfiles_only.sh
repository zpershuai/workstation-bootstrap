#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# shellcheck source=../lib.sh
source "${ROOT_DIR}/scripts/lib.sh"

run_module "${ROOT_DIR}/scripts/mac/check_env.sh"
run_module "${ROOT_DIR}/scripts/mac/25_repos.sh"
run_module "${ROOT_DIR}/scripts/mac/30_dotfiles.sh"
