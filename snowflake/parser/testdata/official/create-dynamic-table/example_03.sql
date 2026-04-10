CREATE DYNAMIC TABLE my_dynamic_table (date TIMESTAMP_NTZ, id NUMBER, content VARIANT)
  TARGET_LAG = '20 minutes'
  WAREHOUSE = mywh
  CLUSTER BY (date, id)
  AS
    SELECT product_id, product_name FROM staging_table;
