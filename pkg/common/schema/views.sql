CREATE VIEW cols AS
  SELECT printf('%016x', id)
    , table_name
    , column_name
    , column_type
    , notes
    , group_concat(distinct version)
  FROM
    (
      SELECT col.id, col.table_name, col.column_name, col.column_type, col.notes, v.version
      FROM columns AS col
      JOIN version_columns vc ON col.id = vc.column_id
      JOIN versions v ON vc.version_id = v.id
      ORDER BY col.table_name, col.column_name, v.version
    )
  group by id
  order by table_name, column_name;

-- CREATE VIEW latest_information_schema_support AS
--   SELECT db_name, '[' ||  table_name || '](' || url || ')', notes
--   FROM information_schema_table_support
--   ORDER BY db_name ASC, table_name ASC;
