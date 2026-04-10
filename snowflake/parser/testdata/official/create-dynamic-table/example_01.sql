CREATE OR REPLACE DYNAMIC TABLE my_dynamic_table
  TARGET_LAG = '20 minutes'
  WAREHOUSE = mywh
  AS
    SELECT product_id, product_name FROM staging_table;
