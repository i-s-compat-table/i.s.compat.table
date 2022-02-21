CREATE VIEW cols AS
  SELECT
      printf('%016x', id) AS id
    , db_name
    , table_name
    , column_name
    , group_concat(distinct column_type) AS column_type
    , (
        CASE group_concat(distinct nullable)
          WHEN '1' THEN 'true'
          WHEN '0' THEN 'false'
          ELSE group_concat(distinct nullable)
        END
      ) AS nullable
    , (
        SELECT urls.url
        FROM urls
        WHERE urls.id = url_id
        LIMIT 1
      ) AS url
    , note
    , group_concat(distinct version) AS versions
    , col.license
    , (
        SELECT urls.url
        FROM urls
        WHERE urls.id = col.license_url_id
        LIMIT 1
      ) AS license_url
    , col.attribution
  FROM
    (
      SELECT
          cv.id
        , db.name "db_name"
        , v.version "version"
        , t.name table_name
        , col.name column_name
        , col.id column_id
        , column_type.name column_type
        , cv.nullable
        , note.note
        , link.id url_id
        , v.is_current
        , license.license
        , license.url_id AS license_url_id
        , license.attribution
      FROM column_versions AS cv
      JOIN versions AS v ON v.id = cv.version_id
      JOIN dbs AS db ON db.id = v.db_id
      JOIN columns AS col ON col.id = cv.column_id
      JOIN tables AS t ON t.id = col.table_id
      LEFT OUTER JOIN types AS column_type ON column_type.id = cv.type_id
      -- ^ sometimes scraped docs don't include types
      LEFT OUTER JOIN urls AS link ON cv.url_id = link.id
      -- ^ observed columns don't come with reference urls
      LEFT OUTER JOIN notes AS note ON cv.note_id = note.id
      LEFT OUTER JOIN licenses AS license ON cv.note_license_id = license.id
      ORDER BY 2, 4, 5, cast(v.version AS REAL) DESC, v.version
    ) AS col
  GROUP BY col.column_id, col.note, col.column_type
  ORDER BY table_name, column_name, db_name, 7; -- 7 = versions

CREATE VIEW IF NOT EXISTS column_support AS
  SELECT
      table_name
    , column_name
    , group_concat(distinct db_name)
  FROM  (
      SELECT
          cv.id
        , db.name "db_name"
        , v.version "version"
        , t.name table_name
        , col.name column_name
        , col.id column_id
        , column_type.name column_type
        , cv.nullable
        , note.note
        , link.id url_id
        , v.is_current
        , license.license
        , license.url_id AS license_url_id
        , license.attribution
      FROM column_versions AS cv
      JOIN versions AS v ON v.id = cv.version_id
      JOIN dbs AS db ON db.id = v.db_id
      JOIN columns AS col ON col.id = cv.column_id
      JOIN tables AS t ON t.id = col.table_id
      LEFT OUTER JOIN types AS column_type ON column_type.id = cv.type_id
      -- ^ sometimes scraped docs don't include types
      LEFT OUTER JOIN urls AS link ON cv.url_id = link.id
      -- ^ observed columns don't come with reference urls
      LEFT OUTER JOIN notes AS note ON cv.note_id = note.id
      LEFT OUTER JOIN licenses AS license ON cv.note_license_id = license.id
      ORDER BY 4, 5, 2, cast(v.version AS REAL) DESC, v.version
    )
  GROUP BY table_name, column_name
  ORDER BY table_name, column_name;
  
CREATE VIEW IF NOT EXISTS relation_support AS
  SELECT
      table_name
    , group_concat(distinct db_name)
  FROM  (
      SELECT
          cv.id
        , db.name "db_name"
        , v.version "version"
        , t.name table_name
        , col.name column_name
        , col.id column_id
        , column_type.name column_type
        , cv.nullable
        , note.note
        , link.id url_id
        , v.is_current
        , license.license
        , license.url_id AS license_url_id
        , license.attribution
      FROM column_versions AS cv
      JOIN versions AS v ON v.id = cv.version_id
      JOIN dbs AS db ON db.id = v.db_id
      JOIN columns AS col ON col.id = cv.column_id
      JOIN tables AS t ON t.id = col.table_id
      LEFT OUTER JOIN types AS column_type ON column_type.id = cv.type_id
      -- ^ sometimes scraped docs don't include types
      LEFT OUTER JOIN urls AS link ON cv.url_id = link.id
      -- ^ observed columns don't come with reference urls
      LEFT OUTER JOIN notes AS note ON cv.note_id = note.id
      LEFT OUTER JOIN licenses AS license ON cv.note_license_id = license.id
      ORDER BY 4, 5, 2, cast(v.version AS REAL) DESC, v.version
    )
  GROUP BY table_name
  ORDER BY table_name;