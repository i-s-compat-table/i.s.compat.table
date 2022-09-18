-- fixup trivial note differences

UPDATE column_versions
SET note_id = (
    SELECT id
    FROM notes
    WHERE note = 'Name of the role to which this role membership was granted (can be the current user, or a different role in case of nested role memberships)'
  )                                                                           --
WHERE note_id = (
  SELECT id
  FROM notes
  WHERE note = 'Name of the role to which this role membership was granted (may be the current user, or a different role in case of nested role memberships)'
);                                                                          --


UPDATE column_versions
SET note_id = (
  SELECT id FROM notes
  WHERE note = 'If data_type identifies a character type, the maximum possible length in octets (bytes) of a datum; null for all other data types. The maximum octet length depends on the declared character maximum length (see above) and the server encoding.'
)
WHERE note_id = (
  SELECT id FROM notes
  WHERE note = 'If data_type identifies a character type, the maximum possible length in octets (bytes) of a datum (this should not be of concern to PostgreSQL users); null for all other data types.'
);                                                                                                                 --- unneccessary editorial ------------------------

-- delete the newly-unreferenced unreferenced notes
WITH unreferenced AS (
  SELECT * FROM notes AS note
  LEFT OUTER JOIN column_versions cv ON note.id = cv.note_id
  WHERE cv.id IS NULL
)
DELETE FROM notes AS note WHERE note.id IN (SELECT id FROM unreferenced);
