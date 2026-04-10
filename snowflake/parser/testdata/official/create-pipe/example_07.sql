CREATE PIPE mypipe_aws
  AUTO_INGEST = TRUE
  AS
  COPY INTO snowpipe_db.public.mytable
  FROM @snowpipe_db.public.mystage
  FILE_FORMAT = (TYPE = 'JSON');
