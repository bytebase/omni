CREATE OR ALTER DYNAMIC TABLE my_dynamic_table
 TARGET_LAG = '5 minutes'
 WAREHOUSE = my_other_wh
 CLUSTER BY (a)
 AS
   SELECT a, b FROM t;
