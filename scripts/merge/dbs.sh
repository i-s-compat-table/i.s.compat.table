#!/usr/bin/env bash
### USAGE assemble.sh [-h|--help] OUTPUT INPUTS...
### merge each input-database in INPUTS into OUTPUT
### ARGS:
###   -h|--help  print this message and exit
###   OUTPUT     a path where you'd like to write the output of the merge
###   INPUTS...  paths to sqlite3 databases

here="${BASH_SOURCE[0]%/*}"
bulk_merge_sql_file="$here/merge.sql"

# shellcheck source=../common.sh
. "$here/../common.sh"

db_init_script="$repo_root/pkg/common/schema/db.sql"
current_db_version="$(
  grep -ie 'pragma user_version' "$db_init_script" |
    awk -F'( |=|;)' '{ print $5 }'
)"

check_user_version() {
  local input_path="$1"
  this_user_version="$(sqlite3 "$input_path" 'pragma user_version;')"
  if [ "$this_user_version" != "$current_db_version" ]; then
    fail "expected \`pragma user_version = $current_db_version\`; got $this_user_version from '$input_path'"
  fi
}

main() {
  set -euo pipefail
  local bulk_merge_sql=""
  bulk_merge_sql="$(<"$bulk_merge_sql_file")"
  local output_path=""
  local input_paths=()
  while [ -n "${1:-}" ]; do
    case "$1" in
    -h | --help) usage && exit 0 ;;
    -*) fail "invalid argument '$1'" ;;
    *)
      output_path="$1"
      shift
      while [ -n "${1:-}" ]; do
        case "$1" in
        -*) fail "invalid argument '$1'" ;;
        *)
          input_paths+=("$1")
          shift
          ;;
        esac
      done
      ;;
    esac
  done

  for input_path in "${input_paths[@]}"; do
    if ! test -e "$input_path"; then
      fail "input '$input_path' does not exist"
    fi

    if ! test -f "$input_path"; then
      fail "input '$input_path' is not a file"
    fi
    echo 'select 1;' | sqlite3 "$input_path" &>/dev/null || fail "sqlite3 can't connect to $input_path"
    check_user_version "$input_path"
  done
  if ! test -e "$output_path"; then
     sqlite3 "$output_path" <"$db_init_script"
  elif test -f "$output_path"; then
    check_user_version "$output_path"
  else
    fail "output '$output_path' is not a file"
  fi

  for input_db in "${input_paths[@]}"; do
    sql="ATTACH DATABASE '$input_db' AS other; $bulk_merge_sql; DETACH other;"
    printf "starting %s..." "$input_db"
    sqlite3 "$output_path" "$sql"
    echo "done"
  done
}

if [ "${BASH_SOURCE[0]}" = "$0" ]; then main "$@"; fi
