CREATE DYNAMIC ICEBERG TABLE my_dynamic_iceberg_v3_table (
    num_orders NUMBER(10,0),
    order_day
  )
  TARGET_LAG = '20 minutes'
  WAREHOUSE = my_warehouse
  EXTERNAL_VOLUME = 'my_external_volume'
  CATALOG = 'SNOWFLAKE'
  BASE_LOCATION = 'my_dynamic_iceberg_v3_table'
  ICEBERG_VERSION = 3
  AS
    SELECT
        COUNT(DISTINCT order_id)
        DATE_TRUNC('DAY', order_timestamp_ns) AS order_day
      FROM staging_v3_iceberg_table;
