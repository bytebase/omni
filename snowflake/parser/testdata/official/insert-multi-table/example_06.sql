INSERT ALL
  WHEN c > 10 THEN
    INTO t1 (col1, col2) VALUES (a, b)
SELECT a, b, c FROM src;
