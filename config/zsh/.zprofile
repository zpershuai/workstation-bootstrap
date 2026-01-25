# Login shell config (primary settings).

# Load local secrets if present.
if [[ -f "${HOME}/.config/secrets/env" ]]; then
  source "${HOME}/.config/secrets/env"
fi

# Oh My Zsh config (managed by .zprofile).
export ZSH="${HOME}/.oh-my-zsh"
ZSH_THEME="ys"
plugins=(autojump macos zsh-autosuggestions zsh-syntax-highlighting)
if [[ -f "${ZSH}/oh-my-zsh.sh" ]]; then
  source "${ZSH}/oh-my-zsh.sh"
fi

# Homebrew shell env (Apple Silicon default path).
if [[ -x "/opt/homebrew/bin/brew" ]]; then
  eval "$(/opt/homebrew/bin/brew shellenv)"
fi

# Add user-local bin paths.
export PATH="${HOME}/.local/bin:${PATH}"

# bun completions
if [[ -s "${HOME}/.bun/_bun" ]]; then
  source "${HOME}/.bun/_bun"
fi

# bun
export BUN_INSTALL="${HOME}/.bun"
export PATH="${BUN_INSTALL}/bin:${PATH}"

# Local env bootstrap if present.
if [[ -f "${HOME}/.local/bin/env" ]]; then
  . "${HOME}/.local/bin/env"
fi

alias claude-mem='${HOME}/.bun/bin/bun "${HOME}/.claude/plugins/marketplaces/thedotmack/plugin/scripts/worker-service.cjs"'
