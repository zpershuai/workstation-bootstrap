# fzf config (managed by workstation-bootstrap).

if ! command -v fzf >/dev/null 2>&1; then
  return 0
fi

# Use fd for fast and git-aware file/dir listing.
if command -v fd >/dev/null 2>&1; then
  export FZF_DEFAULT_COMMAND='fd --type f --strip-cwd-prefix --hidden --follow --exclude .git'
  export FZF_CTRL_T_COMMAND="${FZF_DEFAULT_COMMAND}"
  export FZF_ALT_C_COMMAND='fd --type d --strip-cwd-prefix --hidden --follow --exclude .git'
fi

# Preview with bat first, fall back to ls for non-regular files.
export FZF_DEFAULT_OPTS="
--height=70%
--layout=reverse
--border=rounded
--preview 'bat --color=always --style=plain --line-range=:240 {} 2>/dev/null || ls -la {}'
--preview-window=right,60%,wrap
--bind=ctrl-u:half-page-up,ctrl-d:half-page-down
"

export FZF_CTRL_T_OPTS="
--preview 'bat --color=always --style=plain --line-range=:240 {} 2>/dev/null || ls -la {}'
--preview-window=right,60%,wrap
"

export FZF_ALT_C_OPTS="
--preview 'eza -la --color=always {} 2>/dev/null || ls -la {}'
--preview-window=right,60%,wrap
"

export FZF_CTRL_R_OPTS="
--sort
--exact
--preview 'echo {}'
--preview-window=down,3,wrap
"

# Homebrew fzf shell completion + key bindings.
if [[ -f "/opt/homebrew/opt/fzf/shell/completion.zsh" ]]; then
  source "/opt/homebrew/opt/fzf/shell/completion.zsh"
fi
if [[ -f "/opt/homebrew/opt/fzf/shell/key-bindings.zsh" ]]; then
  source "/opt/homebrew/opt/fzf/shell/key-bindings.zsh"
fi
