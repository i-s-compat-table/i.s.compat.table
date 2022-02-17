CREATE VIEW cols AS
  SELECT
      printf('%016x', id) AS id
    , db_name
    , table_name
    , column_name
    , column_type
    , note
    , group_concat(distinct version) AS versions
    , (
        SELECT url
        FROM urls AS link
        WHERE link.id = col.url_id AND col.is_current
        LIMIT 1
      ) AS url
    -- , col.license
    -- , col.link AS license_link
    -- , col.attribution
  FROM
    (
      SELECT
          cv.id
        , db.name "db_name"
        , v.version "version"
        , t.name table_name
        , col.name column_name
        , column_type.name column_type
        , cv.nullable
        , note.note
        , link.id url_id
        , v.is_current
        , license.license
        , license_url.url
        , license.attribution
      FROM column_versions AS cv
      JOIN versions AS v ON v.id = cv.version_id
      JOIN dbs AS db ON db.id = v.db_id
      JOIN columns AS col ON col.id = cv.column_id
      JOIN tables AS t ON t.id = col.table_id
      JOIN types AS column_type ON column_type.id = cv.type_id
      LEFT OUTER JOIN urls AS link ON cv.url_id = link.id
      LEFT OUTER JOIN notes AS note ON cv.note_id = note.id
      LEFT OUTER JOIN licenses AS license ON cv.note_license_id = license.id
      LEFT OUTER JOIN urls AS license_url ON license.link_id = license_url.url
      ORDER BY 2, 4, 5, cast(v.version AS REAL), v.version
    ) AS col
  GROUP BY col.id
  ORDER BY table_name, column_name;
