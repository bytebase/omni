INSERT OVERWRITE ALL
  INTO t1
  INTO t1 (c1, c2, c3) VALUES (n2, n1, DEFAULT)
  INTO t2 (c1, c2, c3)
  INTO t2 VALUES (n3, n2, n1)
SELECT n1, n2, n3 from src;
