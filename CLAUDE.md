# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a workstation bootstrap repository for automating macOS and Ubuntu machine setup. The design prioritizes idempotency, safety (backups before overwrites), and separation of concerns between software installation and configuration restoration.

## Running the Setup

```bash
./setup.sh                    # Full bootstrap flow (auto-detects OS)
bash scripts/mac/10_brew.sh   # Run individual macOS module
bash scripts/mac/25_repos.sh  # Sync external config repos only
```

To verify after running:
```bash
brew bundle check             # Verify Homebrew packages
command -v nvim && nvim --version
```

## Architecture

### Entry Point & OS Detection

`setup.sh` is the single entrypoint. It detects the OS via `uname -s` and dispatches to platform-specific module chains in `scripts/mac/` or `scripts/ubuntu/`. Each module is a standalone bash script that can be run independently.

### Module Execution Order

Modules are named with numeric prefixes (`NN_topic.sh`) to enforce execution order:

- `00_prereq.sh` - System prerequisites (Xcode CLT, git, curl)
- `10_brew.sh` / `10_apt.sh` - Package manager installation and bundle
- `20_npm.sh` - Global npm packages
- `25_repos.sh` - External config repo synchronization
- `30_dotfiles.sh` - Symlink dotfiles from `config/` to home
- `40_macos_defaults.sh` - macOS system preferences (macOS only)

### Core Library Functions

`scripts/lib.sh` provides shared utilities used across all modules:

- `log "message"` - Standardized logging output
- `die "error"` - Fatal error with exit
- `run_module "path"` - Execute a module if it exists, skip if missing
- `ensure_dir "path"` - Create directory if absent
- `backup_path "path"` - Move existing file/dir to timestamped backup in `~/.dotfiles_backup/`
- `safe_link "src" "dest"` - Idempotent symlink: backups existing, skips correct links, creates new

Critical pattern: All modules must source lib.sh at the top:
```bash
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"
```

### Source-of-Truth Files

- `brew/Brewfile` - Homebrew packages and casks (generate via `brew bundle dump`)
- `npm/packages.txt` - Global npm packages (one per line, `#` for comments)
- `repos/repos.lock` - External Git repo manifest (see below)
- `config/` - Dotfiles stored in this repo, linked to home via `30_dotfiles.sh`

### External Config Repos Pattern

Some dotfiles are managed in separate Git repos (e.g., Neovim config). Rather than nesting Git repos or using submodules, this repo uses a manifest-driven approach:

`repos/repos.lock` format (space-separated):
```
name    url                    dest                          ref
nvim    git@github.com:user/nvim.git   ~/.dotfiles.d/repos/nvim    main
tmux    git@github.com:user/tmux.git   ~/.dotfiles.d/repos/tmux    v1.2.0
```

The `25_repos.sh` module:
1. Parses the manifest line-by-line
2. Clones if missing, fetches if exists
3. Checks out `ref` if provided, otherwise pulls current branch
4. Expands `~` to `$HOME` in dest paths

After cloning, symlinks are created separately in `30_dotfiles.sh` to point from target locations (e.g., `~/.config/nvim -> ~/.dotfiles.d/repos/nvim`).

### Safety & Idempotency

- All modules use `set -euo pipefail` for strict error handling
- `safe_link` never destroys data: existing files are backed up with timestamp
- Modules can be re-run safely; checks prevent duplicate operations
- Missing modules are skipped gracefully via `run_module`
