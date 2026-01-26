# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a workstation bootstrap repository for automating macOS (Phase 1) and Ubuntu (Phase 2) machine setup. The design prioritizes idempotency, safety (backups before overwrites), and separation of concerns between software installation and configuration restoration.

**Scope:** macOS first, Ubuntu later. Windows is out of scope.

## Running the Setup

```bash
./setup.sh                      # Full bootstrap flow (auto-detects OS)
bash scripts/mac/10_brew.sh     # Run individual macOS module
bash scripts/mac/25_repos.sh    # Sync external config repos only
bash scripts/mac/dotfiles_only.sh  # Run env check + repos + dotfiles only
```

Utility scripts:
```bash
bash scripts/mac/check_env.sh   # Verify git/ssh/gh auth before cloning
bash scripts/mac/05_fonts.sh    # Install fonts from misc/fonts/
bash scripts/mac/checkhealth.sh # Verify symlinks and setup state
bash scripts/mac/backup_test.sh # Move managed dotfiles to backup for testing
```

To verify after running:
```bash
brew bundle check               # Verify Homebrew packages
command -v nvim && nvim --version
```

## Architecture

### Entry Point & OS Detection

`setup.sh` is the single entrypoint. It detects the OS via `uname -s` and dispatches to platform-specific module chains in `scripts/mac/` or `scripts/ubuntu/`. Each module is a standalone bash script that can be run independently.

### Environment Gate

Before running the full bootstrap, `check_env.sh` verifies tools and authentication (git/ssh/gh) to catch issues early. This prevents failures during external repo cloning.

**Common failures:**
- **Missing SSH key**: Run `ssh-keygen -t ed25519`, add to ssh-agent, and upload public key to Git host
- **Permission denied (publickey)**: Verify SSH key is added to Git host with `ssh -T git@github.com`
- **repos.lock format error**: Each line must be `name url dest [ref]` (space-separated)
- **gh not logged in**: Run `gh auth login` for HTTPS repos

### Module Execution Order

Modules are named with numeric prefixes (`NN_topic.sh`) to enforce execution order:

| Module | Purpose |
|--------|---------|
| `check_env.sh` | Environment gate (verify git/ssh/gh) |
| `00_prereq.sh` | System prerequisites (Xcode CLT, oh-my-zsh) |
| `05_fonts.sh` | Copy fonts from `misc/fonts/` to `~/Library/Fonts` |
| `10_brew.sh` / `10_apt.sh` | Package manager installation and bundle |
| `20_npm.sh` | Global npm packages |
| `25_repos.sh` | External config repo synchronization |
| `30_dotfiles.sh` | Symlink dotfiles from `config/` to home |
| `40_macos_defaults.sh` | macOS system preferences (optional, macOS only) |

### Core Library Functions

`scripts/lib.sh` provides shared utilities used across all modules:

- `log "message"` - Standardized logging output
- `die "error"` - Fatal error with exit
- `run_module "path"` - Execute a module if it exists, skip if missing
- `ensure_dir "path"` - Create directory if absent
- `backup_path "path"` - Move existing file/dir to timestamped backup in `~/.dotfiles_backup/`
- `safe_link "src" "dest"` - Idempotent symlink: backups existing, skips correct links, creates new

**Critical pattern:** All modules must source lib.sh at the top:
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

**`repos/repos.lock` format** (space-separated):
```
name    url                    dest                          ref
nvim    git@github.com:user/nvim.git   ~/.dotfiles.d/repos/nvim    main
tmux    git@github.com:user/tmux.git   ~/.dotfiles.d/repos/tmux    v1.2.0
```

Rules:
- `name`: Unique identifier
- `url`: Git clone URL (SSH recommended)
- `dest`: Absolute path where repo should live
- `ref`: Optional branch/tag/commit. If omitted, uses default branch

The `25_repos.sh` module:
1. Parses the manifest line-by-line
2. Clones if missing, fetches if exists
3. Checks out `ref` if provided, otherwise pulls current branch
4. Expands `~` to `$HOME` in dest paths

After cloning, symlinks are created separately in `30_dotfiles.sh` to point from target locations (e.g., `~/.config/nvim -> ~/.dotfiles.d/repos/nvim`).

### Oh My Zsh

Oh My Zsh is installed via the official script during `00_prereq.sh`:
```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

It is **not** tracked in `repos/repos.lock` since the installer manages its own updates.

### Safety & Idempotency

- All modules use `set -euo pipefail` for strict error handling
- `safe_link` never destroys data: existing files are backed up with timestamp
- Modules can be re-run safely; checks prevent duplicate operations
- Missing modules are skipped gracefully via `run_module`

## Coding Style & Naming Conventions

- **Shell interpreter**: Use `bash` or `zsh` with a shebang (`#!/bin/bash` or `#!/bin/zsh`)
- **Indentation**: 2 spaces; avoid tabs
- **Script naming**: `NN_topic.sh` to preserve execution order
- **Dotfile naming**: Repo dotfiles should be non-hidden (no leading dot) and linked to hidden targets in `$HOME`
  - Example: `config/zsh/zshrc` → `~/.zshrc`
  - Example: `config/git/gitconfig` → `~/.gitconfig`
- **No formatter**: Keep scripts readable and defensive
- **Error handling**: Always use `set -euo pipefail`

## Repository Structure

```
workstation-bootstrap/
├── setup.sh                    # Entrypoint
├── README.md                   # User documentation
├── AGENTS.md                   # AI agent guidelines
├── CLAUDE.md                   # This file
├── scripts/
│   ├── lib.sh                  # Core library functions
│   ├── mac/
│   │   ├── check_env.sh        # Environment gate
│   │   ├── 00_prereq.sh
│   │   ├── 05_fonts.sh
│   │   ├── 10_brew.sh
│   │   ├── 20_npm.sh
│   │   ├── 25_repos.sh
│   │   ├── 30_dotfiles.sh
│   │   ├── 40_macos_defaults.sh
│   │   ├── backup_test.sh      # Testing utilities
│   │   ├── checkhealth.sh
│   │   └── dotfiles_only.sh
│   └── ubuntu/                 # Future support
├── brew/
│   └── Brewfile                # Homebrew packages lock
├── npm/
│   └── packages.txt            # Global npm packages
├── repos/
│   └── repos.lock              # External repos manifest
├── config/                     # In-repo dotfiles (non-hidden)
│   ├── zsh/
│   │   └── zprofile
│   ├── tmux/
│   │   └── tmux.conf
│   └── git/
│       ├── gitconfig
│       └── gitignore_global
└── misc/                       # Bundled assets and snapshots
    ├── fonts/                  # Fonts to install locally
    ├── cc-switch/              # App-specific config snapshots
    └── dotfiles/               # Platform-agnostic scripts
```

## Adding New Configurations

### In-Repo Dotfiles
1. Place file/dir under `config/` (non-hidden, no leading dot)
2. Add a `safe_link` entry in `scripts/mac/30_dotfiles.sh`

### External Repo Configs
1. Add entry to `repos/repos.lock`
2. Add symlink in `scripts/mac/30_dotfiles.sh` from `~/.dotfiles.d/repos/<name>`

### Misc App Snapshots
Place app-specific data under `misc/` and link from `scripts/mac/30_dotfiles.sh`

## Secrets & Security

**Do NOT commit:**
- API keys, tokens, private SSH keys
- Personal certificates
- Passwords

**Preferred approach:**
- Local secrets: `~/.config/secrets/env` (gitignored)
- Source from shell config if exists: `[[ -f ~/.config/secrets/env ]] && source ~/.config/secrets/env`
- Git identity: Keep in `~/.gitconfig.local` (not in repo)

**Git config structure:**
- `~/.gitconfig` - Main config (symlinked from repo)
- `~/.config/git/.gitconfig.base` - Base settings (symlinked from repo)
- `~/.gitconfig.local` - User identity and local overrides (NOT in repo)

## Design Principles

1. **Idempotent**: Safe to run multiple times; no duplicated installs or broken links
2. **Minimal surprises**: Separate "install software" from "restore configs"
3. **Safe linking**: Existing files are backed up before being replaced by symlinks
4. **Extensible**: macOS first; Ubuntu later with the same structure
5. **Boring and reliable**: Prefer simple, battle-tested solutions over complex tooling
