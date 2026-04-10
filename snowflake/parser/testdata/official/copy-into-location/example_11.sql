COPY INTO @%t1
  FROM t1
  FILE_FORMAT=(TYPE=parquet)
  INCLUDE_QUERY_ID=true;

SELECT last_query_id();

LS @%t1;
