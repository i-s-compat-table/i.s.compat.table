#!/bin/bash
here="${BASH_SOURCE[0]%/*}"
repo_root="$here/.."
target="$repo_root/data/information_schema/columns.tsv"
printf "db_name\ttable_name\tcolumn_name\tcolumn_type\tcolumn_nullable\turl\tnotes\tversions\n" > "$target"

dbs=(pg mariadb mssql)
for db in "${dbs[@]}"; do
  tail -n +2 "$repo_root/data/information_schema/$db.tsv" \
  | sed 's/"/\"/g' >> "$target"
done

tail -n +2 "$repo_root/data/information_schema/snowflakedb.tsv" \
  | awk -F "\t" '{
    print $1 "\t" $2 "\t" $3 "\t" $4 "\t"$5 "\t"$6 "\t" "\t"$8
  }' >> "$target"
