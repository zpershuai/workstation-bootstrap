# Repository Guidelines

## Project Overview

This is a macOS-first dotfiles and machine bootstrap repository. It automates the setup of a fresh macOS (and future Ubuntu) machine using shell scripts, Homebrew, npm packages, and external git repositories.

## Project Structure & Module Organization

- `setup.sh` - Main entrypoint at repo root
- `scripts/` - Platform-specific modules
  - `mac/` - macOS scripts (numbered: `NN_topic.sh`)
  - `ubuntu/` - Ubuntu scripts (numbered: `NN_topic.sh`)
  - `lib.sh` - Shared library with logging, backup, and linking utilities
- `brew/Brewfile` - Homebrew packages source of truth
- `npm/packages.txt` - Global npm packages
- `repos/repos.lock` - External git repositories manifest
- `config/` - Dotfiles stored without leading dots
- `misc/` - Platform-agnostic tools and resources

## Build, Test, and Development Commands

### Main Commands
- `./setup.sh` - Run the full bootstrap flow
- `DRY_RUN=1 ./setup.sh` - Preview changes without applying

### Module Commands
- `bash scripts/mac/check_env.sh` - Verify environment and auth
- `bash scripts/mac/00_prereq.sh` - Install prerequisites
- `bash scripts/mac/05_fonts.sh` - Install fonts
- `bash scripts/mac/10_brew.sh` - Sync Homebrew packages
- `bash scripts/mac/20_npm.sh` - Install global npm packages
- `bash scripts/mac/25_repos.sh` - Clone/pull external repos
- `bash scripts/mac/30_dotfiles.sh` - Link dotfiles
- `bash scripts/mac/checkhealth.sh` - Verify setup state

### Running Single Tests
- `bash scripts/mac/checkhealth.sh` - Health check for links and state
- `bash scripts/lib.sh` - Test library loading
- Individual script testing: `bash -n scripts/mac/10_brew.sh` (syntax check)

## Coding Style & Naming Conventions

### Shell Scripts
- **Shebang**: `#!/usr/bin/env bash`
- **Strict mode**: `set -euo pipefail` at start of every script
- **Indentation**: 2 spaces, no tabs
- **Script naming**: `NN_topic.sh` with numeric prefix for ordering

### Variables
- **Locals**: Use `local var_name` for function variables
- **Constants**: UPPER_SNAKE_CASE (e.g., `ROOT_DIR`, `BACKUP_DIR`)
- **Environment**: Optional with defaults (e.g., `DRY_RUN="${DRY_RUN:-0}"`)
- **Quoting**: Always quote variables: `"${var}"`

### Functions
- **Naming**: lowercase_snake_case
- **Definitions**: `func_name() {`
- **Returns**: Use `return` for status codes, `echo` for output
- **Logging**: Use `log()`, `warn()`, `die()` from lib.sh

### Imports
Always source lib.sh with shellcheck directive:
```bash
# shellcheck source=../lib.sh
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/lib.sh"
```

### Error Handling
- Use `set -euo pipefail` in all scripts
- Use `die()` from lib.sh for fatal errors
- Check command existence: `command -v cmd >/dev/null 2>&1`
- Validate file existence: `[[ -f "${file}" ]]` before operations

### Path Handling
- Compute ROOT_DIR: `ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"`
- Use absolute paths with `"${ROOT_DIR}/..."`
- Create directories: `ensure_dir "${path}"` from lib.sh

### Linking
- Use `safe_link "${src}" "${dest}"` from lib.sh
- Links target non-hidden files in repo to hidden paths in $HOME
- Backup existing files automatically via `backup_path()`

### Logging
- Use provided log functions from lib.sh:
  - `log "message"` - Info messages
  - `warn "message"` - Warnings
  - `die "message"` - Fatal errors (exits 1)
- Prefix: `[dotfiles] message`

## Testing Guidelines

### Testing Approach
- No formal test framework; use simple shell-based checks
- Health checks in `scripts/mac/checkhealth.sh` verify:
  - Symlink correctness
  - Binary presence (fish, starship)
  - External repo state
  - Font installation

### Adding Tests
- Add health check functions to `checkhealth.sh`
- Use `expect_link "${src}" "${dest}"` pattern
- Test individual scripts with `bash -n script.sh`
- Use `DRY_RUN=1` for safe testing

### Test Commands
```bash
# Syntax check
bash -n scripts/mac/10_brew.sh

# Dry run
DRY_RUN=1 bash scripts/mac/30_dotfiles.sh

# Health check
bash scripts/mac/checkhealth.sh
```

## Security Guidelines

- **Never commit secrets**: API keys, tokens, private keys, passwords
- **Prefer local files**: Use `~/.config/secrets/env` (git-ignored)
- **Source secrets conditionally**: Check file exists before sourcing
- **Git identity**: Store in `~/.gitconfig.local`, not in repo

## Commit & Pull Request Guidelines

### Commit Messages
- Use concise, imperative subjects (50 chars or less)
- Examples: "Add brew bundle script", "Fix font installation check"

### PR Requirements
- Brief summary of changes and intent
- Commands run and their results
- Risk notes (especially for destructive changes)
- Test output from checkhealth.sh

## Idempotency & Safety

- Scripts must be safe to run multiple times
- Backup existing files before overwriting
- Check links are already correct before recreating
- Use `DRY_RUN` mode for non-destructive previews

## Platform Support

### macOS (Primary)
- Scripts in `scripts/mac/`
- Homebrew as package manager
- Xcode CLT required

### Ubuntu (Planned)
- Scripts in `scripts/ubuntu/`
- apt-based package management
- Same external repo and dotfile structure

## External Repositories

Manifest format in `repos/repos.lock`:
```
name    url                              dest                              ref(optional)
nvim    git@github.com:xxx/nvim.git      ~/.dotfiles.d/repos/nvim          main
tmux    git@github.com:xxx/tmux.git      ~/.dotfiles.d/repos/tmux          v1.2.0
```

- Clone to `~/.dotfiles.d/repos/<name>/`
- Symlink to appropriate config locations
- Optional `ref` for branch/tag/commit pinning
