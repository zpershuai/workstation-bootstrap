# Oh My Zsh config.

export ZSH="${HOME}/.oh-my-zsh"
ZSH_THEME="ys"
plugins=(autojump macos zsh-autosuggestions zsh-syntax-highlighting)

if [[ -f "${ZSH}/oh-my-zsh.sh" ]]; then
  source "${ZSH}/oh-my-zsh.sh"
fi

# Load main config from .zprofile.
if [[ -f "${HOME}/.zprofile" ]]; then
  source "${HOME}/.zprofile"
fi
