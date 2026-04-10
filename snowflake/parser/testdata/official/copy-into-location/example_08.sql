CREATE or replace TABLE t1 (
  dt date,
  ts time
  )
AS
  SELECT TO_DATE($1)
        ,TO_TIME($2)
    FROM VALUES
           ('2020-01-28', '18:05')
          ,('2020-01-28', '22:57')
          ,('2020-01-28', NULL)
          ,('2020-01-29', '02:15')
;

SELECT * FROM t1;

COPY INTO @%t1
  FROM t1
  PARTITION BY ('date=' || to_varchar(dt, 'YYYY-MM-DD') || '/hour=' || to_varchar(date_part(hour, ts)))
  FILE_FORMAT = (TYPE=parquet)
  MAX_FILE_SIZE = 32000000
  HEADER=true;

LIST @%t1;
