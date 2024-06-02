SELECT
  col.table_name
  , col.column_name
  , col.ordinal_position
  , col.is_nullable
  , coalesce(col.domain_name, col.data_type) -- not all dbs have domains
FROM information_schema.columns AS col
WHERE lower(col.table_schema) = 'information_schema'
