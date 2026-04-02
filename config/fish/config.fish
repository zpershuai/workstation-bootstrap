# Interactive fish config (managed by workstation-bootstrap).

# 光标与按键模式配置
fish_vi_key_bindings
# Keep common readline-style navigation in vi mode.
bind -M insert \ca beginning-of-line
bind -M insert \ce end-of-line
bind -M default \ca beginning-of-line
bind -M default \ce end-of-line
set -g fish_vi_force_cursor 1
set -g fish_cursor_default block blink
set -g fish_cursor_insert block blink
set -g fish_cursor_replace_one block blink
set -g fish_cursor_replace block blink
set -g fish_cursor_external block blink
set -g fish_cursor_visual block blink

set -gx LANG en_US.UTF-8
set -gx LC_ALL en_US.UTF-8
set -gx VISUAL nvim
set -gx EDITOR nvim
set -gx LESSCHARSET utf-8
set -gx XDG_CONFIG_HOME $HOME/.config
set -gx CLAUDE_CODE_AUTO_UPDATER_DISABLED 1
set -gx DISABLE_AUTOUPDATER 1
set -gx BUN_INSTALL $HOME/.bun

function load_posix_env --argument-names file_path
  if not test -f $file_path
    return
  end

  bash -lc "source \"$file_path\" >/dev/null 2>&1; env" 2>/dev/null | while read -l line
    set -l kv (string split -m 1 = -- $line)
    if test (count $kv) -ne 2
      continue
    end

    set -l key $kv[1]
    set -l value $kv[2]

    switch $key
      case PWD SHLVL _ fish_pid hostname history
        continue
    end

    set -gx $key $value
  end
end

if test -f $HOME/.config/secrets/env
  load_posix_env $HOME/.config/secrets/env
end

if test -x /opt/homebrew/bin/brew
  set -gx HOMEBREW_PREFIX /opt/homebrew
  set -gx HOMEBREW_CELLAR /opt/homebrew/Cellar
  set -gx HOMEBREW_REPOSITORY /opt/homebrew
  fish_add_path --move --path /opt/homebrew/bin /opt/homebrew/sbin
end

fish_add_path --move --path \
  $HOME/.dotfiles/bin \
  $HOME/.local/bin \
  /opt/homebrew/opt/openjdk/bin \
  /opt/homebrew/opt/node/bin \
  $BUN_INSTALL/bin \
  /Users/perry/.antigravity/antigravity/bin \
  /Users/perry/.openclaw-wrapper/bin

if test -f $HOME/.mtl-cli/bin/mtl-activate
  load_posix_env $HOME/.mtl-cli/bin/mtl-activate
end

if command -q zoxide
  zoxide init fish | source
end

if command -q fzf
  if command -q fd
    set -gx FZF_DEFAULT_COMMAND 'fd --type f --strip-cwd-prefix --hidden --follow --exclude .git'
    set -gx FZF_CTRL_T_COMMAND $FZF_DEFAULT_COMMAND
    set -gx FZF_ALT_C_COMMAND 'fd --type d --strip-cwd-prefix --hidden --follow --exclude .git'
  end

  set -gx FZF_DEFAULT_OPTS "\
--height=70% \
--layout=reverse \
--border=rounded \
--preview 'bat --color=always --style=plain --line-range=:240 {} 2>/dev/null || ls -la {}' \
--preview-window=right,60%,wrap \
--bind=ctrl-u:half-page-up,ctrl-d:half-page-down"

  set -gx FZF_CTRL_T_OPTS "\
--preview 'bat --color=always --style=plain --line-range=:240 {} 2>/dev/null || ls -la {}' \
--preview-window=right,60%,wrap"

  set -gx FZF_ALT_C_OPTS "\
--preview 'eza -la --color=always {} 2>/dev/null || ls -la {}' \
--preview-window=right,60%,wrap"

  set -gx FZF_CTRL_R_OPTS "\
--sort \
--exact \
--preview 'echo {}' \
--preview-window=down,3,wrap"

  fzf --fish | source
end

if test -f $HOME/.local/bin/env
  load_posix_env $HOME/.local/bin/env
end

if command -q starship
  starship init fish | source
end

# Fish 语法高亮颜色配置
# 有效命令 -> 绿色 (与 starship success_symbol 一致)
set -g fish_color_command '#7bd389' --bold

# 无效/错误命令 -> 红色 (与 starship error_symbol 一致)
set -g fish_color_error '#ff6b6b' --bold

# 其他语法高亮配色
set -g fish_color_keyword '#ffd166' --bold    # 关键字 (if, for, while)
set -g fish_color_option '#8ecae6'            # 选项 (-flag)
set -g fish_color_param '#ffb454'             # 参数
set -g fish_color_redirection '#8ecae6'       # 重定向 (>, <, |)
set -g fish_color_quote '#7bd389'             # 引号字符串
set -g fish_color_comment 'brblack'           # 注释
set -g fish_color_autosuggestion 'brblack'    # 自动建议
set -g fish_color_valid_path --underline      # 存在的文件路径

alias cat bat
alias tmatt 'tmux att -t'
alias tmls 'tmux ls'
alias tmnew 'tmux new-session -s'
alias tmkill 'tmux kill-session -t'
alias vim nvim
alias ls $HOME/.dotfiles/bin/ls
function git
  env LANG=en_GB command git $argv
end

function rm
  set -l args
  set -l passthrough 0
  for arg in $argv
    if test $passthrough -eq 1
      set args $args $arg
      continue
    end

    switch $arg
      case '-*'
        continue
      case '--'
        set passthrough 1
      case '*'
        set args $args $arg
    end
  end
  command trash $args
end

function mkdir_date
  set -l script $HOME/.dotfiles/bin/mkdir_date
  if test -r $script
    set -l data_dir (sh $script | head -n 1)
    if test -n "$data_dir"
      cd $data_dir
      return
    end
  end

  set -l fallback_dir "$HOME/Log/"(date "+%Y_%m_%d")
  mkdir -p $fallback_dir
  cd $fallback_dir
end

function y
  set -l tmp (mktemp -t yazi-cwd.XXXXXX)
  command yazi $argv --cwd-file=$tmp
  if test -f $tmp
    set -l cwd (command cat $tmp)
    if test -n "$cwd" -a "$cwd" != "$PWD"
      cd $cwd
    end
    command rm -f $tmp
  end
end

function auto_delete_login --on-event fish_login
  set -l auto_delete_bin $HOME/.dotfiles/bin/auto_delete.sh
  if test -x $auto_delete_bin
    $auto_delete_bin --quiet
  end
end

function tmux_attach_on_ssh --on-event fish_login
  if test -z "$TMUX" -a -n "$SSH_CONNECTION"
    $HOME/.dotfiles/bin/tmwork
  end
end

# Ctrl-C: cancel current line and show new prompt on next line (like bash)
# Must clear commandline first to avoid ^C⏎ artifacts
bind \cc 'commandline ""; echo ""; commandline -f repaint'
