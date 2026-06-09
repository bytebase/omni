DELETE FROM my_table PARTITION p1 WHERE k1 = 3;
DELETE FROM my_table PARTITION p1 WHERE k1 >= 3 AND k2 = "abc";
DELETE FROM my_table PARTITIONS (p1, p2) WHERE k1 >= 3 AND k2 = "abc";
DELETE FROM t1 USING t2 INNER JOIN t3 ON t2.id = t3.id WHERE t1.id = t2.id
