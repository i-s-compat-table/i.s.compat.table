#!/bin/bash
here="${BASH_SOURCE[0]%/*}"
repo_root="$(cd "$here/.." && pwd)"
db_file="$here/db.sqlite"
db_init_script="$repo_root/scripts/db.sql"
# for dataset in datasets:
#  sqlite3 "ATTACH DATABASE '$input_db' AS other; $bulk_sql; DETACH other;"
bulk_sql="
INSERT INTO 
"
main() {
  sqlite3 -init "$db_init_script" "$db_file"
}
