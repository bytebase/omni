-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBindingFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)
< ok 0 0
> exec ^
: string
| CREATE GLOBAL BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)
< ok 0 0
> exec ^
: string
| CREATE SESSION BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)
< ok 0 0
> exec ^
: string
| create global binding for select * from t using select * from t use index(a)
< ok 0 0
> exec ^
: string
| CREATE BINDING GLOBAL FOR SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 21 near "GLOBAL FOR SELECT * FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE GLOBAL SESSION BINDING FOR SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 21 near "SESSION BINDING FOR SELECT * FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE SESSION GLOBAL BINDING FOR SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 21 near "GLOBAL BINDING FOR SELECT * FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE LOCAL BINDING FOR SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 12 near "LOCAL BINDING FOR SELECT * FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FOR GLOBAL SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 25 near "GLOBAL SELECT * FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t USING GLOBAL SELECT * FROM t
< error
| 1064 42000 …line 1 column 47 near "GLOBAL SELECT * FROM t" 
> exec ^
: string
| DROP BINDING GLOBAL FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 19 near "GLOBAL FOR SELECT * FROM t" 
> exec ^
: string
| SET GLOBAL BINDING ENABLED FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 26 near "ENABLED FOR SELECT * FROM t" 
> exec ^
: string
| SET SESSION BINDING ENABLED FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 27 near "ENABLED FOR SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t USING
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| CREATE BINDING SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 21 near "SELECT * FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING USING
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t USING SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 62 near "USING SELECT * FROM t" 
> exec ^
: string
| DROP BINDING FOR SELECT * FROM t USING SELECT * FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 60 near "USING SELECT * FROM t" 
> exec ^
: string
| DROP BINDING FOR
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| SET BINDING ENABLED
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| SET BINDING ENABLED SELECT * FROM t
< error
| 1064 42000 …line 1 column 26 near "SELECT * FROM t" 
> exec ^
: string
| SET BINDING ENABLED FOR USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 29 near "USING SELECT * FROM t" 
> exec ^
: string
| SET BINDING FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 15 near "FOR SELECT * FROM t" 
> exec ^
: string
| SET BINDING ENABLED FOR SELECT * FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SET BINDING ENABLED FOR SELECT * FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SET BINDING DISABLED FOR SELECT * FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SET BINDING DISABLED FOR SELECT * FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SET BINDING ENABLE FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 18 near "ENABLE FOR SELECT * FROM t" 
> exec ^
: string
| SET BINDING DISABLE FOR SELECT * FROM t
< error
| 1064 42000 …line 1 column 19 near "DISABLE FOR SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING USING SELECT * FROM t
< ok 0 0
> exec ^
: string
| CREATE GLOBAL BINDING USING SELECT * FROM t
< ok 0 0
> exec ^
: string
| CREATE SESSION BINDING USING SELECT * FROM t
< ok 0 0
> exec ^
: string
| CREATE BINDING USING SELECT /*+ USE_INDEX(t, a) */ * FROM t
< ok 0 0
> exec ^
: string
| DROP BINDING USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 18 near "USING SELECT * FROM t" 
> exec ^
: string
| SET BINDING ENABLED USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 25 near "USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FOR SELECT 1 USING SELECT 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR (SELECT * FROM t) USING (SELECT * FROM t)
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR SELECT 1 UNION SELECT 2 INTERSECT SELECT 3 USING SELECT 1 UNION SELECT 2 INTERSECT SELECT 3
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t FOR UPDATE USING SELECT * FROM t FOR UPDATE
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR TABLE t USING TABLE t
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR VALUES ROW(1, 2) USING VALUES ROW(1, 2)
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR UPDATE t SET a = 1 USING UPDATE t SET a = 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR INSERT INTO t1 SELECT * FROM t2 USING INSERT INTO t1 SELECT * FROM t2
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR REPLACE INTO t1 SELECT * FROM t2 USING REPLACE INTO t1 SELECT * FROM t2
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t WHERE a = 1 USING DELETE FROM t WHERE a = 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t USING DELETE FROM t
< ok 0 0
> exec ^
: string
| DROP BINDING FOR DELETE FROM t USING DELETE FROM t
< ok 0 0
> exec ^
: string
| SET BINDING ENABLED FOR DELETE FROM t USING DELETE FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SET BINDING ENABLED FOR DELETE FROM t USING DELETE FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FOR DELETE t1, t2 FROM t1, t2 USING DELETE t1, t2 FROM t1, t2
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t WHERE a = 1 ORDER BY a LIMIT 1 USING DELETE FROM t WHERE a = 1 ORDER BY a LIMIT 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t1 AS x WHERE x.a = 1 USING DELETE FROM t1 AS x WHERE x.a = 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t1 PARTITION (p0) WHERE a = 1 USING DELETE FROM t1 PARTITION (p0) WHERE a = 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR SHOW TABLES USING SHOW TABLES
< error
| 1064 42000 …line 1 column 23 near "SHOW TABLES USING SHOW TABLES" 
> exec ^
: string
| CREATE BINDING FOR EXPLAIN SELECT * FROM t USING EXPLAIN SELECT * FROM t
< error
| 1064 42000 …line 1 column 26 near "EXPLAIN SELECT * FROM t USING EXPLAIN SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FOR CREATE TABLE zz (a INT) USING CREATE TABLE zz (a INT)
< error
| 1064 42000 …line 1 column 25 near "CREATE TABLE zz (a INT) USING CREATE TABLE zz (a INT)" 
> exec ^
: string
| CREATE BINDING FOR SET @a = 1 USING SET @a = 1
< error
| 1064 42000 …line 1 column 22 near "SET @a = 1 USING SET @a = 1" 
> exec ^
: string
| CREATE BINDING FOR ANALYZE TABLE t USING ANALYZE TABLE t
< error
| 1064 42000 …line 1 column 26 near "ANALYZE TABLE t USING ANALYZE TABLE t" 
> exec ^
: string
| CREATE BINDING FOR LOAD DATA INFILE 'a' INTO TABLE t USING LOAD DATA INFILE 'a' INTO TABLE t
< error
| 1064 42000 …line 1 column 23 near "LOAD DATA INFILE 'a' INTO TABLE t USING LOAD DATA INFILE 'a' INTO TABLE t" 
> exec ^
: string
| CREATE BINDING FOR DO 1 USING DO 1
< error
| 1064 42000 …line 1 column 21 near "DO 1 USING DO 1" 
> exec ^
: string
| CREATE BINDING FOR USE test USING USE test
< error
| 1064 42000 …line 1 column 22 near "USE test USING USE test" 
> exec ^
: string
| DROP BINDING FOR SHOW TABLES
< error
| 1064 42000 …line 1 column 21 near "SHOW TABLES" 
> exec ^
: string
| SET BINDING ENABLED FOR SHOW TABLES
< error
| 1064 42000 …line 1 column 28 near "SHOW TABLES" 
> exec ^
: string
| CREATE BINDING FOR SELECT FROM t USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 30 near "FROM t USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t USING SELECT FROM t
< error
| 1064 42000 …line 1 column 52 near "FROM t" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 38 near "USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING USING UPDATE t SET USING
< error
| 1064 42000 …line 1 column 39 near "USING" 
> exec ^
: string
| DROP BINDING FOR UPDATE t SET
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| SET BINDING ENABLED FOR DELETE FROM WHERE a = 1
< error
| 1064 42000 …line 1 column 41 near "WHERE a = 1" 
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t1 USING t1, t2 USING DELETE FROM t1 USING t1, t2
< error
| 1064 42000 …line 1 column 42 near "t1, t2 USING DELETE FROM t1 USING t1, t2" 
> exec ^
: string
| DROP BINDING FOR DELETE FROM t1 USING t1, t2
< error
| 1064 42000 …line 1 column 40 near "t1, t2" 
> exec ^
: string
| CREATE BINDING FOR SELECT * FROM t USING DELETE FROM t1 USING t1, t2
< error
| 1064 42000 …line 1 column 61 near "USING t1, t2" 
> exec ^
: string
| SET BINDING ENABLED FOR DELETE FROM t1 USING t1, t2
< error
| 1064 42000 …line 1 column 47 near "t1, t2" 
> exec ^
: string
| CREATE BINDING FOR DELETE FROM t1, t2 USING DELETE FROM t1, t2
< error
| 1064 42000 …line 1 column 34 near ", t2 USING DELETE FROM t1, t2" 
> exec ^
: string
| DROP BINDING FOR DELETE FROM t1, t2
< error
| 1064 42000 …line 1 column 32 near ", t2" 
> exec ^
: string
| CREATE BINDING FOR WITH c AS (SELECT 1) SELECT * FROM c USING WITH c AS (SELECT 1) SELECT * FROM c
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR WITH RECURSIVE c(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM c WHERE n<3) SELECT * FROM c USING WITH RECURSIVE c(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM c WHERE n<3) SELECT * FROM c
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR WITH c AS (SELECT 1) UPDATE t SET a = 1 USING WITH c AS (SELECT 1) UPDATE t SET a = 1
< ok 0 0
> exec ^
: string
| DROP BINDING FOR WITH c AS (SELECT 1) UPDATE t SET a = 1
< ok 0 0
> exec ^
: string
| CREATE BINDING USING WITH c AS (SELECT 1) UPDATE t SET a = 1
< ok 0 0
> exec ^
: string
| CREATE BINDING FOR WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1 USING WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1
< error
| 1064 42000 …line 1 column 46 near "DELETE FROM t WHERE a = 1 USING WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1" 
> exec ^
: string
| CREATE BINDING USING WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1
< error
| 1064 42000 …line 1 column 48 near "DELETE FROM t WHERE a = 1" 
> exec ^
: string
| DROP BINDING FOR WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1
< error
| 1064 42000 …line 1 column 44 near "DELETE FROM t WHERE a = 1" 
> exec ^
: string
| SET BINDING ENABLED FOR WITH c AS (SELECT 1) DELETE FROM t WHERE a = 1
< error
| 1064 42000 …line 1 column 51 near "DELETE FROM t WHERE a = 1" 
> exec ^
: string
| CREATE BINDING FOR WITH c AS (SELECT 1) INSERT INTO t1 SELECT * FROM c USING WITH c AS (SELECT 1) INSERT INTO t1 SELECT * FROM c
< error
| 1064 42000 …line 1 column 46 near "INSERT INTO t1 SELECT * FROM c USING WITH c AS (SELECT 1) INSERT INTO t1 SELECT * FROM c" 
> exec ^
: string
| CREATE BINDING USING WITH c AS (SELECT 1) REPLACE INTO t1 SELECT * FROM c
< error
| 1064 42000 …line 1 column 49 near "REPLACE INTO t1 SELECT * FROM c" 
> exec ^
: string
| DROP BINDING FOR SELECT * FROM t
< ok 0 0
> exec ^
: string
| DROP GLOBAL BINDING FOR SELECT * FROM t
< ok 0 0
> exec ^
: string
| DROP SESSION BINDING FOR SELECT * FROM t
< ok 0 0
> exec ^
: string
| DROP BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)
< ok 0 0
> exec ^
: string
| SET BINDING ENABLED FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SET BINDING ENABLED FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 'sss'
< error
| 1105 HY000 can't find any plans for 'sss'
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 'sss'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE GLOBAL BINDING FROM HISTORY USING PLAN DIGEST 'sss'
< error
| 1105 HY000 can't find any plans for 'sss'
> exec ^
: string
| CREATE GLOBAL BINDING FROM HISTORY USING PLAN DIGEST 'sss'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE SESSION BINDING FROM HISTORY USING PLAN DIGEST 'sss'
< error
| 1105 HY000 can't find any plans for 'sss'
> exec ^
: string
| CREATE SESSION BINDING FROM HISTORY USING PLAN DIGEST 'sss'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a', 'b'
< error
| 1105 HY000 can't find any plans for 'a'
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a', 'b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @a
< error
| 1105 HY000 can't find specified user variable: a
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @a, @b, 'test1,test2', @c, 'test333'
< error
| 1105 HY000 can't find specified user variable: a
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @a, @b, 'test1,test2', @c, 'test333'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a',
< error
| 1064 42000 …line 1 column 50 near "" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 1
< error
| 1064 42000 …line 1 column 47 near "1" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 1+1
< error
| 1064 42000 …line 1 column 47 near "1+1" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST x
< error
| 1064 42000 …line 1 column 47 near "x" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @@global.x
< error
| 1064 42000 …line 1 column 56 near "@@global.x" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST CONCAT('a', 'b')
< error
| 1064 42000 …line 1 column 52 near "CONCAT('a', 'b')" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST 'a' USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 55 near "USING SELECT * FROM t" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING SQL DIGEST 'a'
< error
| 1064 42000 …line 1 column 37 near "SQL DIGEST 'a'" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING DIGEST 'a'
< error
| 1064 42000 …line 1 column 40 near "DIGEST 'a'" 
> exec ^
: string
| CREATE BINDING FROM HISTORY PLAN DIGEST 'a'
< error
| 1064 42000 …line 1 column 32 near "PLAN DIGEST 'a'" 
> exec ^
: string
| CREATE BINDING FOR HISTORY USING PLAN DIGEST 'a'
< error
| 1064 42000 …line 1 column 26 near "HISTORY USING PLAN DIGEST 'a'" 
> exec ^
: string
| CREATE BINDING FROM USING PLAN DIGEST 'a'
< error
| 1064 42000 …line 1 column 25 near "USING PLAN DIGEST 'a'" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| DROP GLOBAL BINDING FOR SQL DIGEST 's'
< ok 0 0
> exec ^
: string
| DROP SESSION BINDING FOR SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| DROP BINDING FOR SQL DIGEST 'a', 'b'
< ok 0 0
> exec ^
: string
| DROP BINDING FOR SQL DIGEST @a, @b, 'test1,test2', @c, 'test333'
< error
| 1105 HY000 can't find specified user variable: a
> exec ^
: string
| DROP BINDING FOR SQL DIGEST @a, @b, 'test1,test2', @c, 'test333'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST @`a`
< error
| 1105 HY000 can't find specified user variable: a
> exec ^
: string
| DROP BINDING FOR SQL DIGEST @`a`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @`a`
< error
| 1105 HY000 can't find specified user variable: a
> exec ^
: string
| CREATE BINDING FROM HISTORY USING PLAN DIGEST @`a`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST 1
< error
| 1064 42000 …line 1 column 29 near "1" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST a
< error
| 1064 42000 …line 1 column 29 near "a" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST @@global.x
< error
| 1064 42000 …line 1 column 38 near "@@global.x" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST CONCAT('a', 'b')
< error
| 1064 42000 …line 1 column 34 near "CONCAT('a', 'b')" 
> exec ^
: string
| DROP BINDING FOR SQL DIGEST 'a' USING SELECT * FROM t
< error
| 1064 42000 …line 1 column 37 near "USING SELECT * FROM t" 
> exec ^
: string
| DROP BINDING FOR SQL 'a'
< error
| 1064 42000 …line 1 column 24 near "'a'" 
> exec ^
: string
| DROP BINDING FOR DIGEST 'a'
< error
| 1064 42000 …line 1 column 23 near "DIGEST 'a'" 
> exec ^
: string
| DROP BINDING FOR PLAN DIGEST 'a'
< error
| 1064 42000 …line 1 column 21 near "PLAN DIGEST 'a'" 
> exec ^
: string
| SET BINDING ENABLED FOR SQL DIGEST '1'
< ok 0 0
> exec ^
: string
| SET BINDING DISABLED FOR SQL DIGEST '1'
< ok 0 0
> exec ^
: string
| SET BINDING ENABLED FOR SQL DIGEST 'a', 'b'
< error
| 1064 42000 …line 1 column 39 near ", 'b'" 
> exec ^
: string
| SET BINDING ENABLED FOR SQL DIGEST @a
< error
| 1064 42000 …line 1 column 37 near "@a" 
> exec ^
: string
| SET BINDING DISABLED FOR SQL DIGEST @a
< error
| 1064 42000 …line 1 column 38 near "@a" 
> exec ^
: string
| SET BINDING ENABLED FOR SQL DIGEST @a, @b
< error
| 1064 42000 …line 1 column 37 near "@a, @b" 
> exec ^
: string
| SET BINDING ENABLED FOR SQL DIGEST
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| SET BINDING ENABLED FOR PLAN DIGEST 'a'
< error
| 1064 42000 …line 1 column 28 near "PLAN DIGEST 'a'" 
> exec ^
: string
| CREATE BINDING USING SELECT * FROM *.t1
< ok 0 0
> exec ^
: string
| CREATE BINDING USING SELECT * FROM *.t1 WHERE t1.a > (SELECT MAX(a) FROM t2)
< ok 0 0
> exec ^
: string
| CREATE BINDING USING UPDATE *.t1 SET a = 1
< ok 0 0
> exec ^
: string
| SELECT * FROM *.t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| SELECT * FROM *.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM * . t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| SELECT * FROM * . t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM `*`.t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| SELECT * FROM `*`.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM *.t1 JOIN *.t2 ON t1.a = t2.a
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| SELECT a FROM *.t1 JOIN *.t2 ON t1.a = t2.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM *.t1 AS x
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| SELECT * FROM *.t1 AS x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE *.t1 SET a = 1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| UPDATE *.t1 SET a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM *.t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| DELETE FROM *.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO *.t1 VALUES (1, 2)
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| INSERT INTO *.t1 VALUES (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| REPLACE INTO *.t1 VALUES (1, 2)
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| REPLACE INTO *.t1 VALUES (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE *.zz (a INT)
< ok 0 0
> exec ^
: string
| DROP TABLE *.zz
< ok 0 0
> exec ^
: string
| ALTER TABLE *.t ADD COLUMN zz INT
< error
| 1146 42S02 Table '*.t' doesn't exist
> exec ^
: string
| ALTER TABLE *.t ADD COLUMN zz INT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRUNCATE TABLE *.t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| TRUNCATE TABLE *.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| RENAME TABLE *.t1 TO *.t2
< ok 0 0
> exec ^
: string
| ANALYZE TABLE *.t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| ANALYZE TABLE *.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX zz ON *.t1 (a)
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| CREATE INDEX zz ON *.t1 (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW CREATE TABLE *.t1
< error
| 1146 42S02 Table '*.t1' doesn't exist
> exec ^
: string
| SHOW CREATE TABLE *.t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM *.*
< error
| 1064 42000 …line 1 column 17 near "*" 
> exec ^
: string
| SELECT * FROM *.t1.a
< error
| 1064 42000 …line 1 column 19 near ".a" 
> exec ^
: string
| SELECT *.t1.a FROM *.t1
< error
| 1064 42000 …line 1 column 9 near ".t1.a FROM *.t1" 
> exec ^
: string
| SELECT * FROM test.*
< error
| 1064 42000 …line 1 column 20 near "*" 
> exec ^
: string
| SELECT .78*.21
< ok 0 0
> exec ^
: string
| SELECT 1 * .5
< ok 0 0
> exec ^
: string
| CREATE TABLE kw (binding INT, digest INT, plan INT, enabled INT, disabled INT)
< ok 0 0
> exec ^
: string
| SELECT binding, digest, plan, enabled, disabled FROM kw
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| SELECT binding, digest, plan, enabled, disabled FROM kw
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS binding, 2 AS digest, 3 AS plan, 4 AS enabled, 5 AS disabled
< ok 0 0
> exec ^
: string
| SELECT * FROM kw binding
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| SELECT * FROM kw binding
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM kw AS plan
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| SELECT * FROM kw AS plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM kw binding WHERE binding.a = 1
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| DELETE FROM kw binding WHERE binding.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM kw plan WHERE plan.a = 1
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| DELETE FROM kw plan WHERE plan.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE kw binding SET binding.binding = 1
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| UPDATE kw binding SET binding.binding = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SET binding = 1
< ok 0 0
> exec ^
: string
| SET plan = 1
< ok 0 0
> exec ^
: string
| SET enabled = 1
< ok 0 0
> exec ^
: string
| SET @binding = 1
< ok 0 0
> exec ^
: string
| SAVEPOINT binding
< ok 0 0
> exec ^
: string
| CREATE INDEX binding ON kw (binding)
< error
| 1146 42S02 Table 'test.kw' doesn't exist
> exec ^
: string
| CREATE INDEX binding ON kw (binding)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE DATABASE IF NOT EXISTS binding
< ok 0 0
