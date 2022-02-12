-- see https://www.sqlite.org/pragma.html#pragma_user_version
-- Increment this integer whenever you add/remove/rename a table name or
-- add/remove/move/rename/alter the type of a column
pragma user_version = 1; 

CREATE TABLE columns(
    id INTEGER PRIMARY KEY
    -- Always the xxhash3_64 of at least db_name, table_name, and column_name.
    -- If more information is scrapable, `id` might depend on column_type,
    -- column_nullable, and maybe even notes.
  , db_name         TEXT    NOT NULL -- always lowercase
  , table_name      TEXT    NOT NULL -- always lowercase
  , column_name     TEXT    NOT NULL -- always lowercase
  , column_type     TEXT        NULL -- always lowercase
  , column_nullable BOOLEAN     NULL
    -- true if the column is nullable, false if the column is not nullable,
    -- and null if we don't know whether the column is nullable.
  , notes           TEXT        NULL -- no leading/trailing whitespace.
);

CREATE TABLE versions (
    id INTEGER PRIMARY KEY     -- xxhash3_64 of db_name, version
  , db_name TEXT               -- implicitly references columns.db_name
  , version TEXT      NOT NULL
  , release_date TEXT     NULL -- iso-8601 date, manually supplied
);

CREATE TABLE urls(id INTEGER PRIMARY KEY, url TEXT);

CREATE TABLE version_columns(
    version_id    INTEGER REFERENCES versions(id)
  , column_id     INTEGER REFERENCES columns(id)
  , column_number INTEGER -- order of the column in this version of the table
  , CONSTRAINT version_col_pk PRIMARY KEY (version_id, column_id)
);

CREATE TABLE column_reference_urls(
    column_id INTEGER REFERENCES columns(id)
  , url_id    INTEGER REFERENCES urls(id)
  , CONSTRAINT col_ref_url_pk PRIMARY KEY (column_id, url_id)
);

CREATE TABLE licenses(
  id INTEGER PRIMARY KEY -- xxhash3_64 of the license text + attribution text
  , license TEXT
  , link TEXT
  , attribution TEXT
);
-- TODO: create collation for semver-ish version-number strings
