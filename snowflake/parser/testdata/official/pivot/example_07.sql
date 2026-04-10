CREATE OR REPLACE TABLE emp_manager(
    empid INT,
    managerid INT)
  AS SELECT * FROM VALUES
    (1, 7),
    (2, 8),
    (3, 9);
