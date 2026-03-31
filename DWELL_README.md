# dwell

Dwell is a development environment management tool. It synchronizes your dotfiles, external git repositories, and system packages.

## Installation

```bash
# Build from source
make build

# Or install to /usr/local/bin
make install
```

## Quick Start

```bash
# Check status of all modules
dwell status

# Sync all modules
dwell sync

# Sync a specific module
dwell sync nvim

# Run health checks
dwell doctor

# Preview changes (dry run)
dwell sync --dry-run
```

## Configuration

Dwell uses `dwell.yaml` as its primary configuration file. It is backward compatible with `repos/repos.lock`.

### Migrating from repos.lock

```bash
dwell init  # Creates dwell.yaml from repos.lock
```

### Configuration Format

```yaml
version: "1.0"

git:
  - name: nvim
    url: git@github.com:zpershuai/nvim.git
    path: ~/.dotfiles.d/repos/nvim
    ref: main
    links:
      - from: ~/.dotfiles.d/repos/nvim
        to: ~/.config/nvim
```

## Commands

| Command | Description |
|---------|-------------|
| `dwell sync [module]` | Sync all or specific module |
| `dwell status` | Show module status |
| `dwell doctor` | Run health checks |
| `dwell init` | Create dwell.yaml from repos.lock |
| `dwell --version` | Show version |

## Development

```bash
# Build
make build

# Run tests
make test

# Run locally
make dev -- status
```
