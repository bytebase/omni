CREATE EXTERNAL TABLE my_ext_table (
  date_part date as to_date(substr(metadata$filename, 1, 10), 'YYYY/MM/DD'),
  ts timestamp AS (value:time::timestamp),
  user_id varchar AS (value:userId::varchar),
  color varchar AS (value:color::varchar)
) PARTITION BY (date_part)
  LOCATION=@my_ext_stage
  AUTO_REFRESH = false
  FILE_FORMAT=(TYPE=JSON);
