CREATE DYNAMIC TABLE my_dynamic_table
  TARGET_LAG = 'DOWNSTREAM'
  WAREHOUSE = mywh
  INITIALIZE = on_schedule
  REQUIRE USER
  AS
    SELECT product_id, product_name FROM staging_table;
