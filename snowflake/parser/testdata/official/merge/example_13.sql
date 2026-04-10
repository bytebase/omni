MERGE INTO target
  USING (SELECT k, MAX(v) AS v FROM src GROUP BY k) AS b
  ON target.k = b.k
  WHEN MATCHED THEN UPDATE SET target.v = b.v
  WHEN NOT MATCHED THEN INSERT (k, v) VALUES (b.k, b.v);
