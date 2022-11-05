-- see https://www.sqlite.org/pragma.html#pragma_user_version
-- Increment this integer whenever you add/remove/rename a table name or
-- add/remove/move/rename/alter the type of a column
pragma user_version = 2; 

CREATE TABLE dbs(
    id INTEGER PRIMARY KEY -- xxhash3_64(name)
  , name TEXT NOT NULL
);
CREATE TABLE versions (
    id INTEGER PRIMARY KEY     -- xxhash3_64 of db_id, version
  , db_id INTEGER NOT NULL REFERENCES dbs(id)
  , version      TEXT NOT NULL
  -- mutable metadata:
  , release_date TEXT  NULL -- iso-8601 date, manually supplied
  , is_current BOOLEAN NULL
  , version_order INT8 NOT NULL
);

CREATE TABLE tables(
    id INTEGER PRIMARY KEY --xxhash3_64(name)
  , name TEXT NOT NULL -- always lowercase
);
CREATE TABLE columns(
    id INTEGER PRIMARY KEY -- xxhash3_64 of table_id, name.
  , table_id INTEGER NOT NULL REFERENCES tables(id)
  , name     TEXT    NOT NULL -- always lowercase
);
CREATE TABLE types(
    id INTEGER PRIMARY KEY -- xxhash3_64(name)
  , name TEXT NOT NULL     -- always uppercase
); 
CREATE TABLE urls(
    id INTEGER PRIMARY KEY
  , url TEXT NOT NULL
);

CREATE TABLE notes(
    id INTEGER PRIMARY KEY -- xxhash3_64(note)
  , note TEXT NOT NULL -- whitespace trimmed, normalized => single spaces.
);
CREATE TABLE licenses(
    id INTEGER PRIMARY KEY -- xxhash3_64 of the license text + attribution text
  , license TEXT NOT NULL -- ideally a SPDX expression
  , attribution TEXT NOT NULL -- should always start with a copyright symbol
  , url_id INTEGER NULL REFERENCES urls(id) -- mutable metadata
);

CREATE TABLE column_versions (
    id INTEGER PRIMARY KEY -- xxhash3_64 of column_id, version_id
  , column_id INTEGER NOT NULL REFERENCES columns(id)
  , version_id INTEGER NOT NULL REFERENCES versions(id)

  -- mutable metadata:
  , type_id INTEGER NULL REFERENCES types(id)
  , nullable BOOLEAN NULL
    -- true if the column is nullable, false if the column is not nullable,
    -- and null if we don't know whether the column is nullable.
  , url_id INTEGER NULL REFERENCES urls(id)
  , note_id INTEGER NULL REFERENCES notes(id)
  , note_license_id INTEGER NULL REFERENCES licenses(id) -- mutable metadata
);
