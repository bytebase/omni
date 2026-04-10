CREATE DYNAMIC TABLE my_dynamic_table
  TARGET_LAG = '1 hour'
  WAREHOUSE = my_warehouse
  IMMUTABLE WHERE (ts < CURRENT_TIMESTAMP() - INTERVAL '1 day')
AS
  SELECT * FROM source_table;
