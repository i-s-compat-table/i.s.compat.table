CREATE TABLE information_schema_columns(
    db_name         TEXT NOT NULL
  , table_name      TEXT NOT NULL
  , column_name     TEXT NOT NULL
  , column_type     TEXT     NULL
  , column_nullable TEXT     NULL
  , url             TEXT     NULL
  , notes           TEXT     NULL
  , versions        TEXT     NULL
  , CONSTRAINT info_schema_columns_pk PRIMARY KEY (
        db_name
      , table_name
      , column_name
      , column_type
      , column_nullable
    )
  -- should depend on column_type and column_nullable, too :/
);

-- CREATE VIEW information_schema_table_support AS
--   SELECT db_name, table_name, url, notes
--   FROM information_schema_columns
--   ORDER BY db_name ASC, table_name ASC;

-- CREATE VIEW latest_information_schema_support AS
--   SELECT db_name, '[' ||  table_name || '](' || url || ')', notes
--   FROM information_schema_table_support
--   ORDER BY db_name ASC, table_name ASC;
 
-- -- TODO: copy into individual dbs, attach those dbs, then
-- -- insert into main.information_schema_columns values
-- -- (select * from $other.information_schema_columns)
-- -- on conflict do update set versions = ? || ',' || excluded.versions 
-- .mode tabs
-- .import ./data/cockroachdb.tsv information_schema_columns
-- .import ./data/mariadb.tsv information_schema_columns
-- .import ./data/mssql.tsv information_schema_columns
-- .import ./data/mysql.tsv information_schema_columns
-- .import ./data/postgres.tsv information_schema_columns
-- .import ./data/snowflakedb.tsv information_schema_columns
