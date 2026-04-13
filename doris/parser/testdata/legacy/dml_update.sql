UPDATE test SET v1 = 1 WHERE k1 = 1 AND k2 = 2;
UPDATE test SET v1 = v1 + 1 WHERE k1 = 1;
UPDATE t1 SET t1.c1 = t2.c1, t1.c3 = t2.c3 * 100 FROM t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id
