SELECT
    col.table_name
  , col.column_name
  , col.ordinal_position
  , col.is_nullable
  , col.data_type
FROM information_schema.columns AS col
WHERE lower(table_schema) = 'information_schema';
