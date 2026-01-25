# Dotfiles Inventory

This document inventories dotfiles and config directories found in the current home directory. It is meant as a reference for deciding what to bring into this repo and what to leave system-managed.

## Top-level dotfiles (home)

- `~/.zshrc` / `~/.zprofile`: zsh runtime and login configuration.
- `~/.oh-my-zsh/`: Oh My Zsh framework and plugins.
- `~/.tmux.conf` / `~/.tmux/`: tmux configuration and plugin data.
- `~/.gitconfig`: Git user config (name/email, aliases, etc.).
- `~/.profile`: legacy shell profile; may be used by non-zsh shells.
- `~/.ssh/`: SSH keys and config (sensitive; do not commit).
- `~/.npm/`: npm cache and global state.
- `~/.bun/`: Bun runtime cache and globals.
- `~/.cache/`: general application caches (usually not tracked).
- `~/.local/`: user-local data and binaries (tool-specific).
- `~/.dotfiles/`: existing dotfiles repo or staging directory.
- `~/.codex/`: Codex CLI data/config.
- `~/.claude*`: Claude CLI config/memory (`.claude`, `.claude-mem`, `.claude.json`).
- `~/.supermaven/`: Supermaven config/cache.
- `~/.cc-switch/`: CC Switch app data.
- `~/.sogouinput/`: Sogou Input method data.
- `~/.viminfo`: Vim history (not usually tracked).
- `~/.zcompdump*`, `~/.zsh_history`, `~/.zsh_sessions/`: zsh cache/history (not tracked).
- `~/.DS_Store`, `~/.CFUserTextEncoding`, `~/.Trash/`: macOS system files (ignore).

## XDG config directory

- `~/.config/fish/`: fish shell configuration.
- `~/.config/iterm2/`: iTerm2 preferences and profiles.
- `~/.config/karabiner/`: Karabiner-Elements rules and profiles.
- `~/.config/nvim/`: Neovim configuration.
- `~/.config/uv/`: uv (Python tool) settings.

## Git-backed dotfiles/configs

These dotfile directories are Git repositories. If the remote repo name starts with `zpershuai`, it is marked as a personal project.

- `~/.tmux/` -> `git@github.com:zpershuai/tmux.git` (personal)
- `~/.config/nvim/` -> `git@github.com:zpershuai/nvim.git` (personal)
- `~/.oh-my-zsh/` -> `https://github.com/ohmyzsh/ohmyzsh.git`

## Notes for migration

- Prefer moving shell/editor configs into `config/` and symlinking via `scripts/mac/30_dotfiles.sh`.
- Keep secrets and private keys out of git (SSH keys, tokens, API keys).
- Cache/history directories should stay local and untracked.
