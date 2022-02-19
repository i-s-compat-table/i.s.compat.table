#!/bin/bash
### USAGE: dump_tsv.sh [-h|--help] [--dry-run] [--focus[=PROGRAM]]
###                    [-o|--ouptput[=]OUTPUT_PATH] INPUT_PATH
###
### dump a tsv of information_schema columns to OUTPUT_PATH
###
### ARGS:
###   -h|--help           print this message and exit
###   -o=OUTPUT|          where to write the TSV (default: /dev/stdout)
###   --output=OUTPUT
###   --focus[=PROGRAM]   run `sh -c "$PROGRAM '$OUTPUT_PATH'"`. (PROGRAM
###                       defaults to ${EDITOR} if set, else code/vi/ed)
###   --dry-run           validate and print input options

# shellcheck source=./common.sh
. "${BASH_SOURCE[0]%/*}/common.sh"

lookup_preferred_editor() {
  if test -n "${EDITOR:-}"; then echo "$EDITOR" && return 0; fi
  command -v editor || command -v vi || command -v ed
}

cmd="
-- configure output to be tsvs wth headers
.headers on
.mode csv
.separator '$(printf '\t')'

-- ensure the desired view exists
drop view if exists cols;
.read ./pkg/schema/views.sql

-- write to stdout
select * from cols;
"

main() {
  set -euo pipefail
  local input_path=""
  local output_path=""
  local should_focus=""
  local _editor=""
  _editor="$(lookup_preferred_editor)"
  local dry_run=false
  while test -n "${1:-}"; do
    case "$1" in
    --dry-run)
      shift
      dry_run=true
      ;;
    -h | --help) usage && exit 0 ;;
    -o=* | --output=*) output_path="" ;;
    -o | --output)
      shift
      if test -z "${1:-}"; then fail; fi
      case "$1" in
      -*) fail "invalid argument: $1" ;;
      *) output_path="$1" ;;
      esac
      shift
      ;;
    --focus=*)
      should_focus=true
      _editor="${1%*/}"
      echo "$_editor"
      shift
      ;;
    --focus)
      should_focus=true
      shift
      ;;
    -*) fail "invalid argument: $1" ;;
    *)
      if test -n "$input_path"; then
        fail "no input-path argument passed"
      else
        input_path="$1"
        shift
      fi
      ;;
    esac
  done

  # validate input
  if ! test -f "$input_path"; then
    fail "input path '$input_path' is not a file"
  elif ! (head -1 "$input_path" &>/dev/null); then
    fail "input path '$input_path' can't be read"
  fi

  # validate output
  if test -z "$output_path"; then
    if test "$should_focus" = true; then
      fail "can't focus stdout"
    else
      output_path="/dev/stdout"
    fi
  fi

  if test "$dry_run" = "true"; then
    echo "   output_path=$output_path"
    echo "    input_path=$input_path"
    echo "  should_focus=$should_focus"
    echo "       _editor=$_editor"
    exit 0
  fi

  echo "$cmd" | sqlite3 "$input_path" >"$output_path"
  if test "$should_focus" = true; then
    sh -c "$_editor $output_path"
  fi
}

if test "${BASH_SOURCE[0]}" = "$0"; then main "$@"; fi
