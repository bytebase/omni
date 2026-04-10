CREATE PIPE mypipe_azure
  AUTO_INGEST = TRUE
  INTEGRATION = 'MYINT'
  AS
  COPY INTO snowpipe_db.public.mytable
  FROM @snowpipe_db.public.mystage
  FILE_FORMAT = (TYPE = 'JSON');
