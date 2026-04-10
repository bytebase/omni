CREATE OR REPLACE TABLE my_sort_example(a NUMBER, s VARCHAR, b BOOLEAN);
INSERT INTO my_sort_example VALUES
  (0, 'abc', TRUE),
  (0, 'abc', FALSE),
  (0, 'abc', NULL),
  (0, 'xyz', FALSE),
  (0, NULL, FALSE),
  (1, 'xyz', TRUE),
  (NULL, 'xyz', FALSE);
