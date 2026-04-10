SELECT table_id,
  table_name,
  table_schema,
  table_catalog,
  created,
  deleted,
  comment
FROM SNOWFLAKE.ACCOUNT_USAGE.TABLES
WHERE table_catalog = 'DB1'
AND table_schema = 'S1'
AND table_name = 'MY_TABLE'
AND deleted IS NOT NULL
ORDER BY deleted;
