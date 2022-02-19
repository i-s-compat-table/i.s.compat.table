#!/usr/bin/env bash
### USAGE patch_docs.sh [-h|--help] PATH/TO/DIR
### Given a directory containing files `db.sqlite` and `patch.sql`, run
### `patch.sql` on `db.sqlite`, updating the db in-place
# shellcheck disable=2162

# shellcheck source=./common.sh
. "${BASH_SOURCE[0]%/*}/common.sh"

main() {
  set -euo pipefail
  local input_dir=""
  local patch_script=""
  local doc_db=""
  while test -n "${1:-}"; do
    case "$1" in
    -h | --help) usage && exit 0 ;;
    -*) fail "invalid option: $1" ;;
    *)
      if test -z "$input_dir"; then
        input_dir="$1"
        shift
      else
        fail "duplicate input-dir argument"
      fi
      ;;
    esac
  done
  doc_db="$input_dir/db.sqlite"
  patch_script="$input_dir/patch.sql"
  if ! test -d "$input_dir"; then fail "input-dir '$input_dir' is not a directory"; fi
  if ! test -f "$doc_db"; then fail "input path '$doc_db' is not a file"; fi
  if ! test -f "$patch_script"; then fail "no such file '$patch_script'"; fi
  if ! (read <"$doc_db" &>/dev/null); then
    fail "input path '$doc_db' can't be read"
  fi
  if ! (read <"$patch_script" &>/dev/null); then
    fail "input path '$patch_script' can't be read"
  fi
  sqlite3 "$doc_db" <"$patch_script"
}

if test "${BASH_SOURCE[0]}" = "$0"; then main "$@"; fi
