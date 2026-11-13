-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestTableValueConstructorOracle default#1
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE u (a INT, b INT)
< ok 0 0
> exec
| INSERT INTO t VALUES (1, 2), (3, 4)
< ok 2 0
> exec
| TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| TABLE t ORDER BY a LIMIT 1
< error
| 1064 42000 …near 'TABLE t ORDER BY a LIMIT 1' at line 1
> exec
| TABLE t INTO @x
< error
| 1064 42000 …near 'TABLE t INTO @x' at line 1
> exec
| (TABLE t)
< error
| 1064 42000 …near 'TABLE t)' at line 1
> exec
| SELECT * FROM t UNION TABLE u
< error
| 1064 42000 …near 'TABLE u' at line 1
> exec
| SELECT * FROM (TABLE t) x
< error
| 1064 42000 …near 'TABLE t) x' at line 1
> exec
| SELECT * FROM t WHERE a IN (TABLE u)
< error
| 1064 42000 …near 'TABLE u)' at line 1
> exec
| SELECT (TABLE t)
< error
| 1064 42000 …near 'TABLE t)' at line 1
> exec
| WITH c AS (TABLE t) SELECT * FROM c
< error
| 1064 42000 …near 'TABLE t) SELECT * FROM c' at line 1
> exec
| INSERT INTO t TABLE u
< error
| 1064 42000 …near 'TABLE u' at line 1
> exec
| REPLACE INTO t TABLE u
< error
| 1064 42000 …near 'TABLE u' at line 1
> exec
| CREATE TABLE t2 AS TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| CREATE VIEW v AS TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| EXPLAIN TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| CREATE PROCEDURE tvc_p1() TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| VALUES ROW(1, 2)
< error
| 1064 42000 …near 'ROW(1, 2)' at line 1
> exec
| VALUES ROW(1, 2), ROW(3, 4)
< error
| 1064 42000 …near 'ROW(1, 2), ROW(3, 4)' at line 1
> exec
| VALUES (1), ROW(2)
< error
| 1064 42000 …near 'ROW(2)' at line 1
> exec
| SELECT * FROM (VALUES ROW(1, 2)) x
< error
| 1064 42000 …near 'ROW(1, 2)) x' at line 1
> exec
| SELECT * FROM t UNION VALUES ROW(1, 2)
< error
| 1064 42000 …near 'ROW(1, 2)' at line 1
> exec
| INSERT INTO t VALUES ROW(1, 2)
< error
| 1064 42000 …near 'ROW(1, 2)' at line 1
> exec
| INSERT INTO t VALUES (1, 2), ROW(3, 4)
< error
| 1064 42000 …near 'ROW(3, 4)' at line 1
> exec
| INSERT INTO t VALUE ROW(1, 2)
< error
| 1064 42000 …near 'ROW(1, 2)' at line 1
> exec
| REPLACE INTO t VALUES ROW(1, 2)
< error
| 1064 42000 …near 'ROW(1, 2)' at line 1
> exec
| EXPLAIN VALUES ROW(1, 2)
< error
| 1064 42000 …near 'ROW(1, 2)' at line 1
> exec
| SELECT ROW(1, 2) = ROW(1, 2)
< ok 0 0
> exec
| SELECT * FROM t WHERE (a, b) IN (ROW(1, 2))
< ok 0 0
> exec
| VALUES (1, 2)
< ok 0 0
> exec
| VALUES (1, 2), (3, 4)
< ok 0 0
> exec
| VALUES ()
< error
| 4141 HY000 Row with no elements is not allowed in table value constructor in this context
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| VALUES ()
< error
| 4141 HY000 Row with no elements is not allowed in table value constructor in this context
> exec ^
: string
| VALUES ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| VALUES (1, 2), (3)
< error
| 4099 HY000 The used table value constructor has a different number of values
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| VALUES (1, 2), (3)
< error
| 4099 HY000 The used table value constructor has a different number of values
> exec ^
: string
| VALUES (1, 2), (3)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| VALUES (DEFAULT, 1)
< error
| 4172 HY000 'default' is not allowed in this context
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| VALUES (DEFAULT, 1)
< error
| 4172 HY000 'default' is not allowed in this context
> exec ^
: string
| VALUES (DEFAULT, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| VALUES ((1), (2))
< ok 0 0
> exec
| VALUES ((SELECT 1), (VALUES (2)))
< ok 0 0
> exec
| VALUES (@a := 1)
< ok 0 0
> exec
| VALUES (a)
< error
| 4100 HY000 Field reference 'a' can't be used in table value constructor
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| VALUES (a)
< error
| 4100 HY000 Field reference 'a' can't be used in table value constructor
> exec ^
: string
| VALUES (a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| VALUES
< error
| 1064 42000 …near '' at line 1
> exec
| VALUES 1
< error
| 1064 42000 …near '1' at line 1
> exec
| VALUES (1,)
< error
| 1064 42000 …near ')' at line 1
> exec
| VALUES (1) (2)
< error
| 1064 42000 …near '(2)' at line 1
> exec
| VALUES (1 AS x)
< error
| 1064 42000 …near 'AS x)' at line 1
> exec
| VALUES (*)
< error
| 1064 42000 …near '*)' at line 1
> exec
| VALUES (t.*)
< error
| 1064 42000 …near '*)' at line 1
> exec
| VALUES (1), (2) AS x
< error
| 1064 42000 …near 'AS x' at line 1
> exec
| VALUE (1)
< error
| 1064 42000 …near 'VALUE (1)' at line 1
> exec
| VALUES (1, 2) ORDER BY 1
< ok 0 0
> exec
| VALUES (1, 2) ORDER BY 1 DESC LIMIT 1
< ok 0 0
> exec
| VALUES (1, 2) LIMIT 1 OFFSET 1
< ok 0 0
> exec
| VALUES (1, 2) LIMIT 1, 1
< ok 0 0
> exec
| VALUES (1, 2) FOR UPDATE
< ok 0 0
> exec
| VALUES (1, 2) LOCK IN SHARE MODE
< ok 0 0
> exec
| VALUES (1) ORDER BY 1 LIMIT 1 LOCK IN SHARE MODE
< ok 0 0
> exec
| VALUES (1, 2) INTO @x, @y
< ok 1 0
> exec
| VALUES (1) INTO OUTFILE '/tmp/tvc_oracle_out'
< ok 1 0
> exec
| VALUES (1) INTO @x FOR UPDATE
< ok 1 0
> exec
| VALUES (1) FOR UPDATE INTO @x
< ok 1 0
> exec
| VALUES (1) ORDER BY 1 LIMIT 1 FOR UPDATE INTO @x
< ok 1 0
> exec
| (VALUES (1)) INTO @x
< ok 1 0
> exec
| (VALUES (1)) FOR UPDATE
< ok 0 0
> exec
| VALUES (1) UNION VALUES (2) INTO @x
< error
| 1172 42000 Result consisted of more than one row
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| VALUES (1) UNION VALUES (2) INTO @x
< ok 0 0
> exec
| VALUES (1) UNION VALUES (2) FOR UPDATE
< ok 0 0
> exec
| WITH c AS (SELECT 1) VALUES (1) INTO @x
< ok 1 0
> exec
| VALUES (1) ORDER BY (VALUES (1) ORDER BY (SELECT 1))
< ok 0 0
> exec
| VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (SELECT 1))))))))))))))))))))))))))))))
< ok 0 0
> exec
| SELECT (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (VALUES (1) ORDER BY (SELECT 1)))))))))))))))))))))))))))))))
< ok 0 0
> exec
| VALUES (1) ORDER BY 1 ORDER BY 1
< error
| 1064 42000 …near 'ORDER BY 1' at line 1
> exec
| VALUES (1) LIMIT 1 LIMIT 1
< error
| 1064 42000 …near 'LIMIT 1' at line 1
> exec
| VALUES (1) LIMIT 1 ORDER BY 1
< error
| 1064 42000 …near 'ORDER BY 1' at line 1
> exec
| VALUES (1) INTO @x INTO @y
< error
| 1064 42000 …near 'INTO @y' at line 1
> exec
| VALUES (1) INTO @x ORDER BY 1
< error
| 1064 42000 …near 'ORDER BY 1' at line 1
> exec
| VALUES (1) INTO @x UNION VALUES (2)
< error
| 1064 42000 …near 'UNION VALUES (2)' at line 1
> exec
| VALUES (1) ORDER BY 1 WITH ROLLUP
< error
| 1064 42000 …near 'ROLLUP' at line 1
> exec
| VALUES (1) WHERE 1
< error
| 1064 42000 …near 'WHERE 1' at line 1
> exec
| VALUES (1) GROUP BY 1
< error
| 1064 42000 …near 'GROUP BY 1' at line 1
> exec
| VALUES (1) RETURNING a
< error
| 1064 42000 …near 'RETURNING a' at line 1
> exec
| SELECT * FROM (VALUES (1) INTO @x) t
< error
| 1064 42000 …near 'INTO @x) t' at line 1
> exec
| SELECT * FROM t WHERE a IN (VALUES (1) INTO @x)
< error
| 1064 42000 …near 'INTO @x)' at line 1
> exec
| (VALUES (1) INTO @x)
< error
| 1064 42000 …near 'INTO @x)' at line 1
> exec
| CREATE VIEW v AS VALUES (1) INTO @x
< error
| 1064 42000 …near 'INTO @x' at line 1
> exec
| CREATE TABLE t2 VALUES (1) INTO @x
< error
| 1064 42000 …near 'INTO @x' at line 1
> exec
| INSERT INTO t VALUES (1, 2) INTO @x
< error
| 1064 42000 …near 'INTO @x' at line 1
> exec
| (VALUES (1, 2))
< ok 0 0
> exec
| ((VALUES (1, 2)))
< ok 0 0
> exec
| (VALUES (1, 2)) ORDER BY 1
< ok 0 0
> exec
| (VALUES (1, 2) ORDER BY 1) LIMIT 1
< ok 0 0
> exec
| (VALUES (1, 2)) UNION (VALUES (3, 4))
< ok 0 0
> exec
| VALUES (1, 2) UNION VALUES (3, 4)
< ok 0 0
> exec
| VALUES (1, 2) UNION ALL VALUES (3, 4)
< ok 0 0
> exec
| VALUES (1, 2) UNION SELECT 3, 4
< ok 0 0
> exec
| SELECT 1, 2 UNION VALUES (3, 4)
< ok 0 0
> exec
| SELECT 1, 2 UNION VALUES (3, 4) ORDER BY 1
< ok 0 0
> exec
| SELECT 1, 2 UNION ALL VALUES (3, 4) LIMIT 1
< ok 0 0
> exec
| VALUES (1, 2) INTERSECT VALUES (1, 2)
< ok 0 0
> exec
| VALUES (1, 2) EXCEPT VALUES (3, 4)
< ok 0 0
> exec
| VALUES (1, 2) EXCEPT ALL VALUES (3, 4)
< ok 0 0
> exec
| SELECT 1 UNION (VALUES (2))
< ok 0 0
> exec
| VALUES (1) UNION VALUES (2) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| VALUES (1) UNION (VALUES (2) ORDER BY 1 LIMIT 1) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| ((VALUES (1)) UNION VALUES (2)) ORDER BY 1
< ok 0 0
> exec
| VALUES (1) UNION VALUES (2) INTERSECT VALUES (3)
< ok 0 0
> exec
| (VALUES (1) LIMIT 1) UNION (VALUES (2) LIMIT 1)
< ok 0 0
> exec
| SELECT 1 UNION VALUES (2) UNION SELECT 3
< ok 0 0
> exec
| VALUES (1) LIMIT 1 UNION VALUES (2)
< error
| 1064 42000 …near 'UNION VALUES (2)' at line 1
> exec
| VALUES (1) ORDER BY 1 UNION VALUES (2)
< error
| 1064 42000 …near 'UNION VALUES (2)' at line 1
> exec
| VALUES (1) FOR UPDATE UNION VALUES (2)
< error
| 1064 42000 …near 'UNION VALUES (2)' at line 1
> exec
| SELECT 1 UNION VALUES (2) LIMIT 1 UNION SELECT 3
< error
| 1064 42000 …near 'UNION SELECT 3' at line 1
> exec
| SELECT (VALUES (1) LIMIT 1 UNION VALUES (VALUES (1) LIMIT 1 UNION VALUES (VALUES (1) LIMIT 1 UNION VALUES (1))))
< error
| 1064 42000 …near 'UNION VALUES (VALUES (1) LIMIT 1 UNION VALUES (VALUES (1) LIMIT 1 UNION VALUE...' at line 1
> exec
| SELECT * FROM (VALUES (1, 2)) x
< ok 0 0
> exec
| SELECT * FROM (VALUES (1, 2)) AS x
< ok 0 0
> exec
| SELECT * FROM (VALUES (1, 2)) AS `x`
< ok 0 0
> exec
| SELECT * FROM (VALUES (1, 2), (3, 4)) x (c1, c2)
< ok 0 0
> exec
| SELECT * FROM (VALUES (1, 2) ORDER BY 1 LIMIT 1) x
< ok 0 0
> exec
| SELECT * FROM (VALUES (1) FOR UPDATE) t
< ok 0 0
> exec
| SELECT * FROM ((VALUES (1, 2))) x
< ok 0 0
> exec
| SELECT * FROM t JOIN (VALUES (1, 2)) x ON t.a = x.`1`
< ok 0 0
> exec
| SELECT * FROM t NATURAL JOIN (VALUES (1, 2)) AS x (a, b)
< ok 0 0
> exec
| SELECT * FROM (VALUES (1, 2) UNION VALUES (3, 4)) x
< ok 0 0
> exec
| SELECT * FROM (WITH c AS (SELECT 1) VALUES (1)) x
< ok 0 0
> exec
| DELETE t FROM t JOIN (VALUES (1)) x ON t.a = x.`1`
< ok 1 0
> exec
| UPDATE t JOIN (VALUES (1)) x ON t.a = x.`1` SET t.b = 0
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN (VALUES (1), (3))
< ok 0 0
> exec
| SELECT * FROM t WHERE a NOT IN (VALUES (1), (3))
< ok 0 0
> exec
| SELECT * FROM t WHERE (a, b) IN (VALUES (1, 2))
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN ((VALUES (1)))
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN (VALUES (1) UNION VALUES (2))
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN (SELECT 1 UNION VALUES (2))
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN (WITH c AS (SELECT 1) VALUES (1))
< ok 0 0
> exec
| SELECT * FROM t WHERE a IN (VALUES ())
< error
| 4141 HY000 Row with no elements is not allowed in table value constructor in this context
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t WHERE a IN (VALUES ())
< error
| 4141 HY000 Row with no elements is not allowed in table value constructor in this context
> exec ^
: string
| SELECT * FROM t WHERE a IN (VALUES ())
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t WHERE a = (VALUES (1))
< ok 0 0
> exec
| SELECT * FROM t WHERE a = ANY (VALUES (1))
< ok 0 0
> exec
| SELECT * FROM t WHERE a = ALL (VALUES (1))
< ok 0 0
> exec
| SELECT a FROM t WHERE a > SOME (VALUES (1), (2))
< ok 0 0
> exec
| SELECT * FROM t WHERE (a, b) = (VALUES (1, 2))
< ok 0 0
> exec
| SELECT * FROM t WHERE (VALUES (1)) = a
< ok 0 0
> exec
| SELECT * FROM t WHERE EXISTS (VALUES (1))
< ok 0 0
> exec
| SELECT NOT EXISTS (VALUES (1))
< ok 0 0
> exec
| SELECT (VALUES (1))
< ok 0 0
> exec
| SELECT (VALUES (1)) + 1
< ok 0 0
> exec
| SELECT (VALUES (1), (2) LIMIT 1)
< ok 0 0
> exec
| SELECT (WITH c AS (SELECT 1) VALUES (1))
< ok 0 0
> exec
| SELECT 1 IN (VALUES (1))
< ok 0 0
> exec
| DO (VALUES (1))
< ok 0 0
> exec
| SET @x = (VALUES (1))
< ok 0 0
> exec
| UPDATE t SET a = (VALUES (1))
< ok 1 0
> exec
| DELETE FROM t WHERE a IN (VALUES (1))
< ok 1 0
> exec
| WITH c AS (VALUES (1, 2)) SELECT * FROM c
< ok 0 0
> exec
| WITH c (x, y) AS (VALUES (1, 2)) SELECT * FROM c
< ok 0 0
> exec
| WITH c AS (VALUES (1, 2)) VALUES (3, 4)
< ok 0 0
> exec
| WITH c AS (SELECT 1) VALUES (3, 4)
< ok 0 0
> exec
| WITH c AS (SELECT 1) (VALUES (3, 4))
< ok 0 0
> exec
| WITH RECURSIVE c (n) AS (VALUES (1) UNION ALL SELECT n + 1 FROM c WHERE n < 3) SELECT * FROM c
< ok 0 0
> exec
| SELECT * FROM LATERAL (VALUES (1)) x
< error
| 1064 42000 …near '(VALUES (1)) x' at line 1
> exec
| SELECT * FROM (VALUES (1)) x USE INDEX (i)
< error
| 1064 42000 …near 'USE INDEX (i)' at line 1
> exec
| SELECT * FROM t WHERE a IN (VALUES 1)
< error
| 1064 42000 …near '1)' at line 1
> exec
| INSERT INTO t VALUES (1, 2)
< ok 1 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> exec
| INSERT INTO t VALUES (), ()
< ok 2 0
> exec
| INSERT INTO t VALUE (1, 2)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ORDER BY 1
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) LIMIT 1
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) FOR UPDATE
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 2), (3, 4) UNION VALUES (5, 6)
< ok 3 0
> exec
| INSERT INTO t VALUES (1, 2) UNION ALL VALUES (3, 4) UNION DISTINCT SELECT 5, 6
< ok 3 0
> exec
| INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4) ORDER BY 1 LIMIT 1
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) INTERSECT VALUES (1, 2)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) EXCEPT VALUES (3, 4)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) UNION (VALUES (3, 4))
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 2) UNION (SELECT 3, 4)
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4) ON DUPLICATE KEY UPDATE a = 1
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4) RETURNING a
< ok 0 0
> exec
| INSERT INTO t VALUE (1, 2) UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t (a, b) VALUES (1, 2) UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t (a, b) VALUE (1, 2) UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t (a, b) VALUES (1, 2) UNION SELECT 3, 4
< ok 2 0
> exec
| INSERT INTO t PARTITION (p0) VALUES (1, 2) UNION VALUES (3, 4)
< error
| 1747 HY000 PARTITION () clause on non partitioned table
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO t PARTITION (p0) VALUES (1, 2) UNION VALUES (3, 4)
< error
| 1747 HY000 PARTITION () clause on non partitioned table
> exec ^
: string
| INSERT INTO t PARTITION (p0) VALUES (1, 2) UNION VALUES (3, 4)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| INSERT IGNORE INTO t VALUES (1, 2) UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t SELECT 1, 2 UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t SELECT * FROM (VALUES (1, 2)) x
< ok 1 0
> exec
| INSERT INTO t (VALUES (1, 2))
< ok 1 0
> exec
| INSERT INTO t (a, b) (VALUES (1, 2))
< ok 1 0
> exec
| INSERT INTO t ((VALUES (1, 2)))
< ok 1 0
> exec
| INSERT INTO t ((VALUES (1, 2)) UNION (VALUES (3, 4)))
< ok 2 0
> exec
| INSERT INTO t (VALUES (1, 2)) UNION VALUES (3, 4)
< ok 2 0
> exec
| INSERT INTO t (VALUES (1, 2)) ON DUPLICATE KEY UPDATE a = 1
< ok 1 0
> exec
| INSERT INTO t (VALUES (1, 2)) ON DUPLICATE KEY UPDATE a = VALUES(a)
< ok 1 0
> exec
| INSERT INTO t (VALUES (1, 2)) RETURNING a
< ok 0 0
> exec
| INSERT INTO t (SELECT 1, 2) ON DUPLICATE KEY UPDATE a = 1
< ok 1 0
> exec
| INSERT INTO t (WITH c AS (SELECT 1, 2) VALUES (5, 6))
< ok 1 0
> exec
| INSERT INTO t WITH c AS (SELECT 1, 2) VALUES (5, 6)
< ok 1 0
> exec
| INSERT INTO t WITH c AS (SELECT 1, 2) VALUES (5, 6) UNION VALUES (7, 8)
< ok 2 0
> exec
| INSERT INTO t WITH c AS (VALUES (1, 2)) SELECT * FROM c
< ok 1 0
> exec
| REPLACE INTO t VALUES (1, 2) UNION VALUES (3, 4)
< ok 2 0
> exec
| REPLACE INTO t VALUE (1, 2) ORDER BY 1
< ok 1 0
> exec
| REPLACE INTO t (VALUES (1, 2))
< ok 1 0
> exec
| REPLACE INTO t (a, b) (VALUES (1, 2))
< ok 1 0
> exec
| REPLACE INTO t SELECT * FROM (VALUES (1, 2)) x
< ok 1 0
> exec
| CREATE EVENT tvc_ev ON SCHEDULE EVERY 1 DAY DO INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4)
< ok 0 0
> exec
| INSERT INTO t VALUE (1, 2) UNION VALUE (3, 4)
< error
| 1064 42000 …near 'VALUE (3, 4)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) UNION VALUE (3, 4)
< error
| 1064 42000 …near 'VALUE (3, 4)' at line 1
> exec
| INSERT INTO t SELECT 1, 2 UNION VALUE (3, 4)
< error
| 1064 42000 …near 'VALUE (3, 4)' at line 1
> exec
| INSERT INTO t (VALUE (1, 2))
< error
| 1064 42000 …near '(1, 2))' at line 1
> exec
| INSERT INTO t SET a = 1 UNION VALUES (3, 4)
< error
| 1064 42000 …near 'UNION VALUES (3, 4)' at line 1
> exec
| INSERT INTO t VALUES (*)
< error
| 1064 42000 …near '*)' at line 1
> exec
| INSERT INTO t VALUES (t.*)
< error
| 1064 42000 …near '*)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) AS new
< error
| 1064 42000 …near 'AS new' at line 1
> exec
| INSERT INTO t (VALUES (1, 2)) AS x
< error
| 1064 42000 …near 'AS x' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a) + VALUES(b)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(t.a)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(test.t.a)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(`a`)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (VALUES(a))
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (SELECT 1 FROM (VALUES (1)) x)
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUE(a)
< ok 1 0
> exec
| INSERT INTO t SELECT 1, 2 ON DUPLICATE KEY UPDATE a = VALUES(a)
< ok 1 0
> exec
| INSERT INTO t SET a = 1 ON DUPLICATE KEY UPDATE a = VALUES(a)
< ok 1 0
> exec
| SELECT VALUE(a) FROM t
< ok 0 0
> exec
| UPDATE t SET a = VALUE(a)
< ok 69 0
> exec
| CREATE TABLE t3 (a INT) PARTITION BY LIST (VALUES(a)) (PARTITION p0 VALUES IN (1))
< error
| 1564 HY000 This partition function is not allowed
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t3 (a INT) PARTITION BY LIST (VALUES(a)) (PARTITION p0 VALUES IN (1))
< ok 0 0
> exec
| CREATE TABLE t3 (a INT) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (1))
< ok 0 0
> exec
| CREATE TABLE t3 (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< error
| 1050 42S01 Table 't3' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t3 (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| SELECT VALUES(a) FROM t
< error
| 1064 42000 …near 'VALUES(a) FROM t' at line 1
> exec
| SELECT * FROM t WHERE VALUES(a) = 1
< error
| 1064 42000 …near 'VALUES(a) = 1' at line 1
> exec
| UPDATE t SET a = VALUES(a)
< error
| 1064 42000 …near 'VALUES(a)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (SELECT VALUES(a))
< error
| 1064 42000 …near 'VALUES(a))' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (VALUES (1))
< error
| 1064 42000 …near '1))' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = a IN (VALUES (1))
< error
| 1064 42000 …near '1))' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(1)
< error
| 1064 42000 …near '1)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a + 1)
< error
| 1064 42000 …near '+ 1)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(@x)
< error
| 1064 42000 …near '@x)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(t.*)
< error
| 1064 42000 …near '*)' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES()
< error
| 1064 42000 …near ')' at line 1
> exec
| INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a) UNION VALUES (1)
< error
| 1064 42000 …near 'UNION VALUES (1)' at line 1
> exec
| CREATE TABLE t3 (a INT) PARTITION BY HASH (VALUES(a))
< error
| 1064 42000 …near 'VALUES(a))' at line 1
> exec
| CREATE TABLE t3 (a INT DEFAULT VALUES(a))
< error
| 1064 42000 …near 'VALUES(a))' at line 1
> exec
| CREATE TABLE t3 (a INT, CHECK (VALUES(a) > 0))
< error
| 1064 42000 …near 'VALUES(a) > 0))' at line 1
> exec
| CREATE TABLE t2 AS VALUES (1, 2)
< ok 1 0
> exec
| CREATE TABLE t2 VALUES (1, 2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 VALUES (1, 2)
< ok 0 0
> exec
| CREATE TABLE t2 (VALUES (1, 2))
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (VALUES (1, 2))
< ok 0 0
> exec
| CREATE TABLE t2 AS (VALUES (1, 2))
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 AS (VALUES (1, 2))
< ok 0 0
> exec
| CREATE TABLE t2 ((VALUES (1)))
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 ((VALUES (1)))
< ok 0 0
> exec
| CREATE TABLE t2 (WITH c AS (SELECT 1) VALUES (1))
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (WITH c AS (SELECT 1) VALUES (1))
< ok 0 0
> exec
| CREATE TABLE t2 (VALUES (1)) UNION VALUES (2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (VALUES (1)) UNION VALUES (2)
< ok 0 0
> exec
| CREATE TABLE t2 AS SELECT * FROM (VALUES (1, 2)) x
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 AS SELECT * FROM (VALUES (1, 2)) x
< ok 0 0
> exec
| CREATE TABLE t2 AS VALUES (1, 2) UNION VALUES (3, 4)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 AS VALUES (1, 2) UNION VALUES (3, 4)
< ok 0 0
> exec
| CREATE TABLE t2 VALUES (1) UNION VALUES (2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 VALUES (1) UNION VALUES (2)
< ok 0 0
> exec
| CREATE TABLE t2 VALUES (1) ORDER BY 1 LIMIT 1
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 VALUES (1) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| CREATE TABLE t2 VALUES (1) FOR UPDATE
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 VALUES (1) FOR UPDATE
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT) AS VALUES (1, 2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT, b INT) AS VALUES (1, 2)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT) VALUES (1, 2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT, b INT) VALUES (1, 2)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT) (VALUES (1))
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT) (VALUES (1))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT) (SELECT 1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT) (SELECT 1)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT) IGNORE VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT) IGNORE VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT) ENGINE=InnoDB VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT) ENGINE=InnoDB VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT) PARTITION BY HASH(a) VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 (a INT) PARTITION BY HASH(a) VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 AS WITH c AS (SELECT 1) VALUES (1, 2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 AS WITH c AS (SELECT 1) VALUES (1, 2)
< ok 0 0
> exec
| CREATE TABLE t2 ENGINE=InnoDB VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 ENGINE=InnoDB VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 ENGINE=InnoDB AS VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 ENGINE=InnoDB AS VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 ENGINE=InnoDB VALUE (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 ENGINE=InnoDB VALUE (1)
< ok 0 0
> exec
| CREATE TABLE t2 IGNORE VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 IGNORE VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 REPLACE VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 REPLACE VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 IGNORE AS VALUES (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 IGNORE AS VALUES (1)
< ok 0 0
> exec
| CREATE TABLE t2 VALUE (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 VALUE (1)
< ok 0 0
> exec
| CREATE TABLE t2 AS VALUE (1)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 AS VALUE (1)
< ok 0 0
> exec
| CREATE TABLE t2 VALUE (1) UNION VALUES (2)
< error
| 1050 42S01 Table 't2' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t2 VALUE (1) UNION VALUES (2)
< ok 0 0
> exec
| CREATE OR REPLACE TABLE t2 AS VALUES (1)
< ok 1 0
> exec
| CREATE TEMPORARY TABLE t2 VALUES (1)
< ok 1 0
> exec
| CREATE TABLE t2 (a INT) VALUE (1)
< error
| 1064 42000 …near 'VALUE (1)' at line 1
> exec
| CREATE TABLE t2 (a INT) AS VALUE (1)
< error
| 1064 42000 …near 'VALUE (1)' at line 1
> exec
| CREATE TABLE t2 (VALUE (1))
< error
| 1064 42000 …near '(1))' at line 1
> exec
| CREATE TABLE t2 VALUE (1) UNION VALUE (2)
< error
| 1064 42000 …near 'VALUE (2)' at line 1
> exec
| CREATE VIEW v AS VALUES (1, 2)
< ok 0 0
> exec
| CREATE VIEW v AS (VALUES (1, 2))
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS (VALUES (1, 2))
< ok 0 0
> exec
| CREATE VIEW v AS ((VALUES (1)))
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS ((VALUES (1)))
< ok 0 0
> exec
| CREATE VIEW v AS VALUES (1, 2) UNION VALUES (3, 4)
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS VALUES (1, 2) UNION VALUES (3, 4)
< ok 0 0
> exec
| CREATE VIEW v AS (VALUES (1)) UNION (VALUES (2)) ORDER BY 1 LIMIT 1
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS (VALUES (1)) UNION (VALUES (2)) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| CREATE VIEW v (x, y) AS VALUES (1, 2)
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v (x, y) AS VALUES (1, 2)
< ok 0 0
> exec
| CREATE VIEW v AS VALUES (1, 2) WITH CHECK OPTION
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS VALUES (1, 2) WITH CHECK OPTION
< ok 0 0
> exec
| CREATE VIEW v AS VALUES (1) FOR UPDATE
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS VALUES (1) FOR UPDATE
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (VALUES (1, 2)) SELECT * FROM c
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS WITH c AS (VALUES (1, 2)) SELECT * FROM c
< ok 0 0
> exec
| CREATE OR REPLACE VIEW v AS VALUES (1, 2) ORDER BY 1
< ok 0 0
> exec
| CREATE ALGORITHM=MERGE VIEW v AS VALUES (1)
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE ALGORITHM=MERGE VIEW v AS VALUES (1)
< ok 0 0
> exec
| ALTER VIEW v AS VALUES (1, 2)
< ok 0 0
> exec
| EXPLAIN VALUES (1, 2)
< ok 0 0
> exec
| EXPLAIN (VALUES (1, 2))
< ok 0 0
> exec
| EXPLAIN VALUES (1, 2) UNION VALUES (3, 4)
< ok 0 0
> exec
| EXPLAIN EXTENDED (VALUES (1)) UNION VALUES (2)
< ok 0 0
> exec
| EXPLAIN FORMAT=JSON VALUES (1, 2)
< ok 0 0
> exec
| EXPLAIN EXTENDED VALUES (1, 2)
< ok 0 0
> exec
| EXPLAIN PARTITIONS VALUES (1)
< ok 0 0
> exec
| EXPLAIN WITH c AS (SELECT 1) VALUES (1)
< ok 0 0
> exec
| EXPLAIN VALUES (1) INTO @x
< ok 0 0
> exec
| EXPLAIN VALUES (1) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| EXPLAIN INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4)
< ok 0 0
> exec
| CREATE PROCEDURE tvc_p2() BEGIN VALUES (1); END
< ok 0 0
> exec
| CREATE PROCEDURE tvc_p3() VALUES (1)
< ok 0 0
> exec
| CREATE PROCEDURE tvc_p4() BEGIN DECLARE c CURSOR FOR VALUES (1); END
< ok 0 0
> exec
| CREATE PROCEDURE tvc_p5() BEGIN DECLARE c CURSOR FOR (VALUES (1)) UNION VALUES (2); END
< ok 0 0
> exec
| CREATE FUNCTION tvc_f1() RETURNS INT RETURN (VALUES (1))
< ok 0 0
> exec
| CREATE TRIGGER tvc_tg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = (VALUES (1))
< ok 0 0
> exec
| PREPARE tvc_s FROM 'VALUES (1, 2)'
< ok 0 0
> exec
| CREATE VIEW v AS VALUE (1)
< error
| 1064 42000 …near 'VALUE (1)' at line 1
> exec
| CREATE PROCEDURE tvc_p6() BEGIN DECLARE c CURSOR FOR TABLE t; END
< error
| 1064 42000 …near 'TABLE t; END' at line 1
> exec
| VALUES (1), (2) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| VALUES (1) LIMIT 1 INTO @x
< ok 1 0
> exec
| INSERT INTO t VALUES (1, 2) ORDER BY 1 LIMIT 1
< ok 1 0
> exec
| EXPLAIN VALUES (1) ORDER BY 1
< ok 0 0
> exec
| SELECT (VALUES (1) ORDER BY 1 LIMIT 1)
< ok 0 0
> exec
| (VALUES (1) LIMIT 1) INTO @x
< ok 1 0
> exec
| SELECT * FROM (VALUES (1), (2) ORDER BY 1 LIMIT 1) x
< ok 0 0
> exec
| WITH c AS (VALUES (1) ORDER BY 1 LIMIT 1) SELECT * FROM c
< ok 0 0
> exec
| CREATE TABLE n AS VALUES (1) ORDER BY 1 LIMIT 1
< ok 1 0
> exec
| CREATE VIEW v AS VALUES (1) UNION VALUES (2) ORDER BY 1 LIMIT 2
< error
| 1050 42S01 Table 'v' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE VIEW v AS VALUES (1) UNION VALUES (2) ORDER BY 1 LIMIT 2
< ok 0 0
> exec
| SELECT 1 UNION VALUES (2) ORDER BY 1 LIMIT 1
< ok 0 0
> exec
| VALUES (1) UNION VALUES (2) LIMIT 1 INTO @x
< ok 1 0
> exec
| SELECT 1 UNION (VALUES (1) ORDER BY 1 LIMIT 1) ORDER BY 1 LIMIT 2
< ok 0 0
> exec
| (VALUES (1) ORDER BY 1) ORDER BY 1
< ok 0 0
> exec
| VALUES (1) ORDER BY (VALUES (2) ORDER BY (VALUES (3) LIMIT 1))
< ok 0 0
