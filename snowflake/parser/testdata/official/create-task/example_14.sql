CREATE OR ALTER TASK my_task
  WAREHOUSE = regress
  AFTER my_other_task
  AS
    SELECT 2 * PI();
