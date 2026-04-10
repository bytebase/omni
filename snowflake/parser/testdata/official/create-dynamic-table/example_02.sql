CREATE DYNAMIC ICEBERG TABLE my_dynamic_table (date TIMESTAMP_NTZ, id NUMBER, content STRING)
  TARGET_LAG = '20 minutes'
  WAREHOUSE = mywh
  EXTERNAL_VOLUME = 'my_external_volume'
  CATALOG = 'SNOWFLAKE'
  BASE_LOCATION = 'my_iceberg_table'
  AS
    SELECT product_id, product_name FROM staging_table;
