#!/usr/bin/env bash
scripts_dir="${BASH_SOURCE[0]%/*}"
repo_root=""; repo_root="$(cd "$scripts_dir/.." && pwd)"; export repo_root;
usage() { grep '^###' "$0"  | sed 's/^### //g; s/^###//g'; }

fail() {
  local message="$*"
  local red=""
  local reset=""
  if test -n "$message"; then
    if test -t 1 && test -z "${NO_COLOR:-}"; then
      red="$(tput setaf 1)"
      reset="$(tput sgr0)"
    fi
    printf "%s%s%s\n" "$red" "$message" "$reset" >&2
  fi

  usage >&2
  exit 1
}
