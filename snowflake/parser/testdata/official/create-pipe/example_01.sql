CREATE PIPE mypipe
  AS
  COPY INTO mytable
  FROM @mystage
  FILE_FORMAT = (TYPE = 'JSON');
