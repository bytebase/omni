CREATE PIPE mypipe_s3
  AUTO_INGEST = TRUE
  AWS_SNS_TOPIC = 'arn:aws:sns:us-west-2:001234567890:s3_mybucket'
  AS
  COPY INTO snowpipe_db.public.mytable
  FROM @snowpipe_db.public.mystage
  FILE_FORMAT = (TYPE = 'JSON');
