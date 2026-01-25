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
3. Core dev trio: **zsh + tmux + neovim**
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

1. **Prerequisites** (Xcode CLT, Homebrew, PATH sanity)
2. **Install software**  
   - `brew bundle` from `brew/Brewfile` (includes `cask` for GUI apps)
   - `npm install -g` from `npm/packages.txt`
3. **Restore configs (dotfiles)**  
   - symlink `config/` files into `$HOME` and `$XDG_CONFIG_HOME`
4. **External config repos (important)**  
   Some dotfiles are managed as separate Git repos.  
   We do NOT embed `.git` inside the dotfiles repo.  
   Instead, we keep a repo manifest file and clone/pull them during setup.

---

## Repository layout

Recommended structure:

dotfiles/
setup.sh                    # entrypoint
README.md

scripts/
lib.sh                    # logging, backup, safe_link helpers
mac/
00_prereq.sh
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
zsh/.zshrc
tmux/.tmux.conf
nvim/                    # for configs not managed externally
git/.gitconfig
…etc

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

## Step-by-step implementation plan (macOS Phase 1)

### Step 1 — Create the repo skeleton
1. Create directories:
   - `scripts/`, `scripts/mac/`, `brew/`, `npm/`, `config/`, `repos/`
2. Add `setup.sh` (entrypoint) and `scripts/lib.sh` (shared helpers)

### Step 2 — Capture current machine state into lock files
- Export Brewfile:
  ```bash
  brew bundle dump --describe --file ./brew/Brewfile --force

	•	Export global npm packages:

npm ls -g --depth=0 --parseable | tail -n +2 | awk -F/ '{print $NF}' | sort -u > ./npm/packages.txt



Step 3 — Copy critical dotfiles into config/

Start with the essentials:
	•	zsh: .zshrc, .zprofile (and plugin manager config if any)
	•	tmux: .tmux.conf
	•	nvim: either config/nvim/ OR external repo via repos.lock

Step 4 — Add external repos manifest if needed
	•	Put standalone config repos into repos/repos.lock
	•	Implement scripts/mac/25_repos.sh to:
	•	clone if missing
	•	fetch + checkout ref (if provided)
	•	pull updates
	•	fail loudly with clear errors (ssh key, permission, network)

Step 5 — Implement linking strategy

Implement scripts/mac/30_dotfiles.sh:
	•	safe_link to:
	•	backup existing files into ~/.dotfiles_backup/<timestamp>/
	•	create symlinks to repo-managed paths

Step 6 — Run and verify
	•	Run:

./setup.sh


	•	Verify:
	•	brew bundle check
	•	command -v nvim && nvim --version
	•	tmux and zsh behave as expected
	•	symlinks are correct and backups exist

⸻

Usage

Fresh machine bootstrap

git clone <this-repo> ~/dotfiles
cd ~/dotfiles
./setup.sh

After editing Brewfile / packages list

./setup.sh

Update external repos only

bash scripts/mac/25_repos.sh


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

Run dotfiles steps
	•	Sync external repos:
	•	bash scripts/mac/25_repos.sh
	•	Link configs:
	•	bash scripts/mac/30_dotfiles.sh

Add a new config
	•	IN_REPO: put file/dir under config/ and add a safe_link in scripts/mac/30_dotfiles.sh
	•	EXTERNAL_REPO: add entry to repos/repos.lock and link from ~/.dotfiles.d/repos/<name>

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
