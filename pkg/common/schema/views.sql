CREATE VIEW cols AS
  SELECT printf('%016x', id) AS id
    -- , db_name
    , table_name
    , column_name
    , column_type
    , notes
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
          col.id
        , col.table_name
        , col.column_name
        , col.column_type
        , col.notes
        , link.id AS url_id
        , v.version
        , v.is_current
        , license.license
        , license.link
        , license.attribution
      FROM columns AS col
      JOIN version_columns vc ON col.id = vc.column_id
      JOIN urls AS link ON vc.url_id = link.id
      LEFT OUTER JOIN licenses AS license ON link.license_id = license.id
      JOIN versions v ON vc.version_id = v.id
      ORDER BY col.table_name, col.column_name, cast(v.version AS REAL), v.version
    ) AS col
  GROUP BY col.id
  ORDER BY table_name, column_name;

--
-- SELECT
--     table_name
--   , column_name
--   , url
--   , is_current
-- FROM columns AS col
-- JOIN version_columns vc ON col.id = vc.column_id
-- JOIN urls AS link ON vc.url_id = link.id
-- JOIN versions v ON vc.version_id = v.id
-- WHERE is_current
-- ORDER BY table_name, column_name;
