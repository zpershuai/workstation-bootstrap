# Repository Guidelines

## Project Structure & Module Organization

This repo currently contains a single `README.md` that describes the intended bootstrap workflow and layout. The planned structure is:

- `setup.sh` entrypoint at repo root.
- `scripts/` with `mac/` modules named with numeric prefixes (e.g., `scripts/mac/10_brew.sh`).
- `brew/Brewfile`, `npm/packages.txt`, and `repos/repos.lock` as source-of-truth lock files.
- `config/` for dotfiles that live in this repo.

If you add these directories, keep to the naming patterns above and update the `README.md` to match.

## Build, Test, and Development Commands

There are no build or test commands in this repo yet. Expected commands once scripts land:

- `./setup.sh` to run the full bootstrap flow.
- `bash scripts/mac/10_brew.sh` to install Homebrew packages.
- `bash scripts/mac/25_repos.sh` to sync external config repos.

If you add new commands, document them in `README.md` and here.

## Coding Style & Naming Conventions

- Shell scripts should be `bash` or `zsh` and include a shebang.
- Indentation: 2 spaces; avoid tabs.
- Script names: `NN_topic.sh` to preserve execution order.
- Config files should match their target names (e.g., `config/zsh/.zshrc`).
- No formatter is enforced; keep scripts readable and defensive.

## Testing Guidelines

There is no test framework yet. If you add tests, prefer simple shell-based checks and document:

- How to run tests (e.g., `./scripts/test.sh`).
- Naming (e.g., `tests/test_*.sh`).
- Coverage expectations, if any.

## Commit & Pull Request Guidelines

Only a single commit exists, so there is no established commit message convention. Use concise, imperative subjects (e.g., "Add brew bundle script").

PRs should include:

- A brief summary of changes and intent.
- Any commands run and their results.
- Notes on risks (e.g., changes that may affect user machines).

## Security & Configuration Tips

- Do not commit secrets (API keys, tokens, private keys).
- Prefer local, ignored files like `~/.config/secrets/env` and source them from shell configs.
