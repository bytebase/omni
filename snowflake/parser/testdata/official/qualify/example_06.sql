SELECT p, SUM(o) OVER (PARTITION BY p) AS r
  FROM qt
  WHERE o < 4
  GROUP BY p, o
  HAVING SUM(i) > 3
  QUALIFY r IN (
    SELECT MIN(i)
      FROM qt
      GROUP BY p
      HAVING MIN(i) > 3);
