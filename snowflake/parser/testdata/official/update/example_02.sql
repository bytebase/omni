UPDATE target
  SET v = src.v
  FROM src
  WHERE target.k = src.k;
