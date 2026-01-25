#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
MANIFEST="${ROOT_DIR}/repos/repos.lock"

say() {
  local prefix="$1"
  shift
  printf '%s %b\n' "${prefix}" "$*"
}

check() { say "[check]" "$*"; }
warn() { say "[warn]" "$*"; }
error() { say "[error]" "$*"; exit 1; }
hint() { say "[hint]" "$*"; }

check "platform"
if [[ "$(uname -s)" != "Darwin" ]]; then
  error "Unsupported platform. This script only supports macOS (Darwin)."
fi

check "required commands"
command -v bash >/dev/null 2>&1 || error "Missing bash. Fix: xcode-select --install"
command -v git >/dev/null 2>&1 || error "Missing git. Fix: xcode-select --install"
command -v brew >/dev/null 2>&1 || error "Missing Homebrew. Fix: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
command -v node >/dev/null 2>&1 || error "Missing Node.js. Fix: brew install node"
command -v npm >/dev/null 2>&1 || error "Missing npm. Fix: brew install node"

check "repos manifest"
if [[ ! -f "${MANIFEST}" ]]; then
  error "Missing repos manifest: ${MANIFEST}"
fi

ssh_urls=0
https_urls=0
declare -a hosts
hosts=()
line_no=0

while IFS= read -r line || [[ -n "${line}" ]]; do
  line_no=$((line_no + 1))
  line="${line%%#*}"
  line="$(printf '%s' "${line}" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')"
  [[ -z "${line}" ]] && continue

  read -r name url dest ref <<<"${line}"
  if [[ -z "${name}" || -z "${url}" || -z "${dest}" ]]; then
    error "Invalid repos.lock line ${line_no}: expected 'name url dest [ref]'"
  fi

  if [[ "${url}" == git@* ]]; then
    ssh_urls=1
    host="${url#git@}"
    host="${host%%:*}"
    hosts+=("${host}")
  elif [[ "${url}" == ssh://* ]]; then
    ssh_urls=1
    host="${url#ssh://}"
    host="${host#*@}"
    host="${host%%/*}"
    hosts+=("${host}")
  elif [[ "${url}" == https://* ]]; then
    https_urls=1
    host="${url#https://}"
    host="${host%%/*}"
    hosts+=("${host}")
  fi

done < "${MANIFEST}"

if [[ ${ssh_urls} -eq 1 ]]; then
  check "SSH keys"
  if ! ls "${HOME}/.ssh/id_"* >/dev/null 2>&1; then
    error $'No SSH keys found in ~/.ssh. Fix:\n  ssh-keygen -t ed25519 -C "your_email@example.com"\n  eval "$(ssh-agent -s)"\n  ssh-add --apple-use-keychain ~/.ssh/id_ed25519\n  pbcopy < ~/.ssh/id_ed25519.pub\nThen add the public key to your Git host.'
  fi

  if ! ssh-add -l >/dev/null 2>&1; then
    warn $'ssh-agent has no keys loaded. Fix:\n  eval "$(ssh-agent -s)"\n  ssh-add --apple-use-keychain ~/.ssh/id_ed25519'
  fi

  check "SSH connectivity"
  declare -a uniq_hosts=()
  if [[ ${#hosts[@]} -eq 0 ]]; then
    warn "No SSH hosts parsed from repos.lock despite SSH URLs."
  fi
  for h in "${hosts[@]}"; do
    if [[ " ${uniq_hosts[*]-} " != *" ${h} "* ]]; then
      uniq_hosts+=("${h}")
    fi
  done

  for h in "${uniq_hosts[@]}"; do
    if ! ssh -o BatchMode=yes -o ConnectTimeout=5 "git@${h}" -T >/dev/null 2>&1; then
      warn "SSH handshake failed for ${h}. Possible causes: missing key, host not trusted, or no access."
      hint "Try: ssh -T git@${h}"
      hint "If Permission denied(publickey): add your key to ${h} and re-run."
    fi
  done
fi

if [[ ${https_urls} -eq 1 ]]; then
  check "HTTPS auth"
  if command -v gh >/dev/null 2>&1; then
    if ! gh auth status >/dev/null 2>&1; then
      warn "GitHub CLI not authenticated. Fix: gh auth login"
    fi
  else
    warn "gh not installed; HTTPS private repos may require authentication. Fix: brew install gh"
  fi
fi

check "network"
if ! ping -c 1 -t 2 github.com >/dev/null 2>&1; then
  warn "Network check failed (github.com unreachable). Clone may fail until connectivity is restored."
fi

say "[check]" "Environment check passed ✅"
