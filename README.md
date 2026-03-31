# Dotfiles & Machine Bootstrap (macOS-first)

> Goal: Turn a freshly installed machine into my daily dev environment in a repeatable, auditable, and idempotent way.
>
> Scope (Phase 1): macOS  
> Future (Phase 2): Ubuntu  
> Out of scope: Windows

---

## Why this repo exists

I already have a working macOS setup. The problem is that my configuration is scattered across:

1. Homebrew installed packages (CLI tools + GUI apps via cask)
2. Global npm tools (e.g., codex / gemini cli and other developer utilities)
3. Core dev trio: **fish + tmux + neovim**
4. Dotfiles and configs distributed in different places
5. Some configs are **their own Git repos** (e.g., `~/.config/nvim` might be managed elsewhere)

This repository standardizes everything into a single source of truth and provides a one-command setup flow.

---

## Design principles (keep it boring and reliable)

- **Idempotent**: Safe to run multiple times; no duplicated installs or broken links.
- **Minimal surprises**: Separate “install software” from “restore configs”.
- **Source-of-truth files**:  
  - Brew packages live in `brew/Brewfile`  
  - Global npm packages live in `npm/packages.txt`  
  - Dotfiles live in `config/` (or via external repos)
- **Safe linking**: Existing files are backed up before being replaced by symlinks.
- **Extensible**: macOS first; Ubuntu later with the same structure.

---

## High-level workflow

Setup is split into modules:

1. **Environment gate** (`scripts/mac/check_env.sh`)  
   Verify tools and auth (git/ssh/gh) before any clone/pull, and auto-install Homebrew if missing.
2. **Prerequisites** (Xcode CLT, Homebrew, core tools)
3. **Fonts**  
   Copy fonts from `misc/fonts/` into `~/Library/Fonts`.
4. **Install software**  
   - `brew bundle` from `brew/Brewfile`
   - Includes `fish`, `starship`, `yazi`, and common preview/search dependencies (`chafa`, `sevenzip`, `fzf`, `zoxide`, `jq`)
   - `npm install -g` from `npm/packages.txt`
5. **External config repos**  
   Manifest-driven clone/pull from `repos/repos.lock`.
6. **Restore configs (dotfiles)**  
   Symlink `config/` files into `$HOME` and `$XDG_CONFIG_HOME`.

---

## Repository layout

Current structure:

dotfiles/
setup.sh                    # entrypoint
README.md

scripts/
lib.sh                    # logging, backup, safe_link helpers
mac/
check_env.sh             # environment gate (macOS)
00_prereq.sh
05_fonts.sh              # install fonts from misc/
backup_test.sh           # move managed dotfiles to backup for testing
checkhealth.sh           # verify links and setup state
dotfiles_only.sh         # run check_env + repos + dotfiles
10_brew.sh
20_npm.sh
25_repos.sh            # clone/pull external config repos (manifest-driven)
30_dotfiles.sh
40_macos_defaults.sh   # optional: system defaults
ubuntu/                   # reserved for future
00_prereq.sh
10_apt.sh
20_npm.sh
25_repos.sh
30_dotfiles.sh

brew/
Brewfile

npm/
packages.txt

repos/
repos.lock               # external git repos manifest (url + dest + ref)

config/
zsh/zshrc
zsh/zprofile
fish/config.fish
starship.toml
tmux/tmux.conf
git/gitconfig
git/gitignore_global
yazi/yazi.toml
yazi/keymap.toml
…etc

misc/
fonts/                   # bundled fonts to install locally
cc-switch/               # cc-switch config snapshot
dotfiles/                # platform-agnostic scripts from ~/.dotfiles

---

## Managing dotfiles that are separate Git repos

### Problem
Some configs are already standalone repos (e.g., my Neovim config).  
I want a clean bootstrap repo, not a nested Git mess.

### Solution (default): Manifest-driven external repos
- Keep a manifest file: `repos/repos.lock`
- During setup, clone/pull each repo to a controlled location:
  - default: `~/.dotfiles.d/repos/<name>`
- Then symlink from the correct target location:
  - example: `~/.config/nvim -> ~/.dotfiles.d/repos/nvim`

This keeps versioning clean and makes the setup reproducible.

### Shell stack
Primary interactive shell is `fish`, with `starship` as the prompt.

- `fish`, `starship`, `fzf`, and `zoxide` are installed from `brew/Brewfile`
- `ghostty` is configured to start `fish` explicitly
- after bootstrap, switch the login shell with:

```bash
chsh -s "$(command -v fish)"
```

Legacy `zsh` config remains in-repo as a fallback for tools that still expect `~/.zshrc` or `~/.zprofile`.

### Manifest format (example)
`repos/repos.lock`:

name   url                               dest                               ref(optional)

nvim     git@github.com:xxx/nvim.git       ~/.dotfiles.d/repos/nvim            main
tmux     git@github.com:xxx/tmux-conf.git  ~/.dotfiles.d/repos/tmux            v1.2.0

Rules:
- `name`: unique identifier
- `url`: git clone URL (SSH recommended)
- `dest`: absolute path where repo should live
- `ref`: optional (branch/tag/commit). If omitted, use default branch.

---
Usage

Fresh machine bootstrap

git clone <this-repo> ~/dotfiles
cd ~/dotfiles
./setup.sh

Run pieces

bash scripts/mac/check_env.sh
bash scripts/mac/05_fonts.sh
bash scripts/mac/10_brew.sh
bash scripts/mac/backup_test.sh
bash scripts/mac/checkhealth.sh
bash scripts/mac/dotfiles_only.sh
bash scripts/mac/25_repos.sh
bash scripts/mac/30_dotfiles.sh

`bash scripts/mac/10_brew.sh` syncs `brew/Brewfile`: it installs missing dependencies and upgrades outdated Homebrew-managed formulae/casks. To remove packages no longer present in `brew/Brewfile`, run `brew bundle cleanup --file brew/Brewfile --force`.

⸻

Environment Gate (check_env)

Why: catch missing tools and auth issues before cloning external repos.

Run:
	bash scripts/mac/check_env.sh

Common failures & fixes
	•	Missing SSH key:
	•	ssh-keygen -t ed25519 -C "your_email@example.com"
	•	eval "$(ssh-agent -s)"
	•	ssh-add --apple-use-keychain ~/.ssh/id_ed25519
	•	pbcopy < ~/.ssh/id_ed25519.pub
	•	Add the public key to your Git host
	•	Permission denied(publickey):
	•	ssh -T git@github.com
	•	Verify the key is added to the Git host and loaded in ssh-agent
	•	repos.lock format error:
	•	Each line must be: name url dest [ref]
	•	gh not logged in (HTTPS repos):
	•	gh auth login

⸻

Dotfiles (macOS)

Design principles
	•	Idempotent: safe to run repeatedly
	•	Backups: existing files are moved to ~/.dotfiles_backup/<timestamp>/
	•	External repos: git-backed configs are cloned into ~/.dotfiles.d/repos and linked

Shell layout
	•	~/.config/fish -> config/fish (primary interactive shell config)
	•	~/.config/starship.toml: symlink to config/starship.toml
	•	~/.zshrc: symlink to config/zsh/zshrc (legacy fallback)
	•	~/.zprofile: symlink to config/zsh/zprofile (legacy login-shell fallback)
	•	~/.tmux.conf: symlink to config/tmux/tmux.conf (repo-managed tmux entrypoint)
	•	config/tmux/tmux.conf sources ~/.tmux/.tmux.conf from the external tmux repo, then applies local overrides like the default shell
	•	Ghostty starts fish explicitly; run `chsh -s "$(command -v fish)"` to make other login-shell consumers default to fish too

Git config
	•	~/.gitconfig includes ~/.config/git/.gitconfig.base and ~/.gitconfig.local
	•	Keep identity in ~/.gitconfig.local (not in repo)
	•	Repo dotfiles are non-hidden (no leading dot); setup links them to hidden paths in $HOME

Yazi config
	•	~/.config/yazi -> config/yazi (managed by this repo)
	•	Default behavior is tuned for coding workflows: show hidden files, natural sort, larger preview pane
	•	Text/code files open in `$EDITOR` (default `nvim`)

Add a new config
	•	IN_REPO: put file/dir under config/ and add a safe_link in scripts/mac/30_dotfiles.sh
	•	EXTERNAL_REPO: add entry to repos/repos.lock and link from ~/.dotfiles.d/repos/<name>
	•	MISC snapshots: place app-specific data under misc/ and link from scripts/mac/30_dotfiles.sh

Secrets
	•	Local secrets live at ~/.config/secrets/env (ignored by git)
	•	Source it from shell config if present
	•	Git identity stays in ~/.gitconfig.local


⸻

Safety & secrets
	•	Do NOT commit secrets:
	•	API keys, tokens, private SSH keys, personal certificates
	•	Prefer:
	•	~/.config/secrets/env (ignored by git)
	•	source it from .zshrc if exists

⸻

Future: Ubuntu support

Planned mapping:
	•	brew bundle -> apt + snap/flatpak (to be decided)
	•	npm global tools and dotfiles logic remains similar
	•	keep the same manifest approach for external repos

⸻

Open questions (for discussion with Codex)
	1.	External repos: do we want to support “pin to commit hash” for reproducibility?
	2.	Where should repos live by default: ~/.dotfiles.d/repos or ~/Workspace?
	3.	Should we use stow/chezmoi later, or keep safe_link as the core?
	4.	Do we want optional modules (macOS defaults, fonts, login items, etc.)?

⸻

Non-goals (for now)
	•	Windows initialization
	•	Complex UI/interaction flows
	•	Full system preference automation (only if needed later)
