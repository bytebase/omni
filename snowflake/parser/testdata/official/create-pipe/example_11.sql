CREATE PIPE mypipe
  AUTO_INGEST = TRUE
  INTEGRATION = 'MYINT'
  AS
  COPY INTO snowpipe_db.public.my_v3_iceberg_table
  FROM @snowpipe_db.public.mystage
  FILE_FORMAT = (TYPE = 'JSON');
