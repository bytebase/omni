CREATE PIPE mypipe2
  AS
  COPY INTO mytable(C1, C2)
  FROM (SELECT $5, $4 FROM @mystage)
  FILE_FORMAT = (TYPE = 'JSON');
