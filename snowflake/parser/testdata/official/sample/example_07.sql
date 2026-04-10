SELECT *
  FROM (
       SELECT *
         FROM t1 JOIN t2
           ON t1.a = t2.c
       ) SAMPLE (1);
