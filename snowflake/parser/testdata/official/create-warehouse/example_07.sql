-- Create a new warehouse with GENERATION = '1'
CREATE OR ALTER WAREHOUSE test_gen_warehouse
  WITH WAREHOUSE_SIZE = XSMALL
    GENERATION = '1'
    AUTO_SUSPEND = 60
    INITIALLY_SUSPENDED = TRUE;

-- Verify that it was created
SHOW WAREHOUSES LIKE 'test_gen_warehouse'
  ->> SELECT "name", "resource_constraint" FROM $1;

-- Alter it to GENERATION = '2'
CREATE OR ALTER WAREHOUSE test_gen_warehouse
  WITH WAREHOUSE_SIZE = SMALL
    GENERATION = '2'
    AUTO_SUSPEND = 120;

-- Verify that it was altered
SHOW WAREHOUSES LIKE 'test_gen_warehouse'
  ->> SELECT "name", "resource_constraint" FROM $1;

-- Clean up when done
DROP WAREHOUSE test_gen_warehouse;
