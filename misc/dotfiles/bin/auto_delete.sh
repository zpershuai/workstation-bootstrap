#!/bin/bash

# Delete old files from selected directories.
# Usage:
#   auto_delete.sh            # delete files older than 40 days
#   auto_delete.sh -n         # dry run
#   AUTO_DELETE_DAYS=30 auto_delete.sh

set -u

dry_run=0
quiet=0
case "${1:-}" in
  -n|--dry-run)
    dry_run=1
    shift
    ;;
  -q|--quiet)
    quiet=1
    shift
    ;;
esac

dayDiff="${AUTO_DELETE_DAYS:-40}"

dirList=(
  "${HOME}/Log"
  "${HOME}/Downloads"
)

file_find_opts=(
  -mindepth 1
  -type f
  -mtime "+${dayDiff}"
)

dir_find_opts=(
  -mindepth 1
  -type d
  -mtime "+${dayDiff}"
  -empty
)

for dir in "${dirList[@]}"; do
  if [[ ! -d "${dir}" ]]; then
    continue
  fi

  if [[ "${dry_run}" -eq 1 ]]; then
    if [[ "${quiet}" -ne 1 ]]; then
      echo "[dry-run] ${dir}"
      find "${dir}" "${file_find_opts[@]}" -print
      find "${dir}" "${dir_find_opts[@]}" -print
    else
      find "${dir}" "${file_find_opts[@]}" -print > /dev/null 2>&1
      find "${dir}" "${dir_find_opts[@]}" -print > /dev/null 2>&1
    fi
    continue
  fi

  if [[ "${quiet}" -eq 1 ]]; then
    find "${dir}" "${file_find_opts[@]}" -delete > /dev/null 2>&1
    # Remove now-empty old directories in a second pass.
    find "${dir}" -depth "${dir_find_opts[@]}" -delete > /dev/null 2>&1
  else
    echo "[delete] ${dir} (>${dayDiff} days)"
    find "${dir}" "${file_find_opts[@]}" -delete
    # Remove now-empty old directories in a second pass.
    find "${dir}" -depth "${dir_find_opts[@]}" -delete
  fi
done
