INSERT OR IGNORE INTO main.dbs      SELECT * FROM other.dbs;
INSERT OR IGNORE INTO main.tables   SELECT * from other.tables;
INSERT OR IGNORE INTO main.columns  SELECT * FROM other.columns;
INSERT OR IGNORE INTO main.types    SELECT * FROM other.types;
INSERT OR IGNORE INTO main.notes    SELECT * FROM other.notes;
INSERT OR IGNORE INTO main.urls     SELECT * FROM other.urls;
INSERT OR IGNORE INTO main.licenses SELECT * FROM other.licenses;

INSERT INTO main.versions
  SELECT * FROM other.versions WHERE true
  -- "WHERE true" is required by the sqlite parser to resolve some parsing
  -- ambiguities. See https://sqlite.org/lang_upsert.html#parsing_ambiguity
  ON CONFLICT(id) DO UPDATE SET
    is_current = coalesce(excluded.is_current, is_current);

INSERT INTO main.column_versions
  SELECT * FROM other.column_versions WHERE true
  ON CONFLICT(id) DO UPDATE SET
      type_id = coalesce(excluded.type_id, type_id)
     , url_id = coalesce(excluded.url_id, url_id)
     , note_id = coalesce(excluded.note_id, note_id)
     , note_license_id = coalesce(excluded.note_license_id, note_license_id);

-- TODO: clean up dangling notes, licenses, types, urls
