-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDB_DeparseRules default#1
> exec
| CREATE TABLE t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT 1
< ok 0 0
> query
| SHOW CREATE VIEW v1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` (`1`) AS SELECT 1 AS `1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v2 AS SELECT a FROM t WHERE a > 1 GROUP BY a HAVING COUNT(*) > 1 ORDER BY a LIMIT 10
< ok 0 0
> query
| SHOW CREATE VIEW v2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>1 GROUP BY `a` HAVING COUNT(1)>1 ORDER BY `a` LIMIT 10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v3 AS SELECT a, b AS alias1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v3` (`a`, `alias1`) AS SELECT `a` AS `a`,`b` AS `alias1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v4 AS SELECT t.a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v4` (`a`) AS SELECT `t`.`a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v5 AS SELECT `select` FROM t
< error
| 1054 42S22 Unknown column 'select' in 'field list'
> exec
| CREATE VIEW v6 AS SELECT a+b, a-b, a*b, a/b, a DIV b, a MOD b, a%b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v6
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v6
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v6` (`a+b`, `a-b`, `a*b`, `a/b`, `a DIV b`, `a MOD b`, `a%b`) AS SELECT `a`+`b` AS `a+b`,`a`-`b` AS `a-b`,`a`*`b` AS `a*b`,`a`/`b` AS `a/b`,`a` DIV `b` AS `a DIV b`,`a`%`b` AS `a MOD b`,`a`%`b` AS `a%b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v7 AS SELECT a=b, a<>b, a!=b, a>b, a<b, a>=b, a<=b, a<=>b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v7
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v7
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v7` (`a=b`, `a<>b`, `a!=b`, `a>b`, `a<b`, `a>=b`, `a<=b`, `a<=>b`) AS SELECT `a`=`b` AS `a=b`,`a`!=`b` AS `a<>b`,`a`!=`b` AS `a!=b`,`a`>`b` AS `a>b`,`a`<`b` AS `a<b`,`a`>=`b` AS `a>=b`,`a`<=`b` AS `a<=b`,`a`<=>`b` AS `a<=>b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v8 AS SELECT a AND b, a OR b, NOT a, a XOR b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v8` (`a AND b`, `a OR b`, `NOT a`, `a XOR b`) AS SELECT `a` AND `b` AS `a AND b`,`a` OR `b` AS `a OR b`,NOT `a` AS `NOT a`,`a` XOR `b` AS `a XOR b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v9 AS SELECT a|b, a&b, a^b, a<<b, a>>b, ~a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v9` (`a|b`, `a&b`, `a^b`, `a<<b`, `a>>b`, `~a`) AS SELECT `a`|`b` AS `a|b`,`a`&`b` AS `a&b`,`a`^`b` AS `a^b`,`a`<<`b` AS `a<<b`,`a`>>`b` AS `a>>b`,~`a` AS `~a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v10 AS SELECT a IS NULL, a IS NOT NULL, a IS TRUE, a IS FALSE FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v10` (`a IS NULL`, `a IS NOT NULL`, `a IS TRUE`, `a IS FALSE`) AS SELECT `a` IS NULL AS `a IS NULL`,`a` IS NOT NULL AS `a IS NOT NULL`,`a` IS TRUE AS `a IS TRUE`,`a` IS FALSE AS `a IS FALSE` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v11 AS SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v11
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v11
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v11` (`a IN (1,2,3)`, `a NOT IN (1,2,3)`) AS SELECT `a` IN (1,2,3) AS `a IN (1,2,3)`,`a` NOT IN (1,2,3) AS `a NOT IN (1,2,3)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v12 AS SELECT a BETWEEN 1 AND 10, a NOT BETWEEN 1 AND 10 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v12
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v12
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v12` (`a BETWEEN 1 AND 10`, `a NOT BETWEEN 1 AND 10`) AS SELECT `a` BETWEEN 1 AND 10 AS `a BETWEEN 1 AND 10`,`a` NOT BETWEEN 1 AND 10 AS `a NOT BETWEEN 1 AND 10` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v13 AS SELECT a LIKE 'foo%', a NOT LIKE 'bar%', a LIKE 'x' ESCAPE '\\' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v13
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v13
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v13` (`a LIKE 'foo%'`, `a NOT LIKE 'bar%'`, `a LIKE 'x' ESCAPE '\\'`) AS SELECT `a` LIKE _UTF8MB4'foo%' AS `a LIKE 'foo%'`,`a` NOT LIKE _UTF8MB4'bar%' AS `a NOT LIKE 'bar%'`,`a` LIKE _UTF8MB4'x' AS `a LIKE 'x' ESCAPE '\\'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v14 AS SELECT a REGEXP 'pattern', a NOT REGEXP 'pattern' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v14
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v14
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v14` (`a REGEXP 'pattern'`, `a NOT REGEXP 'pattern'`) AS SELECT `a` REGEXP _UTF8MB4'pattern' AS `a REGEXP 'pattern'`,`a` NOT REGEXP _UTF8MB4'pattern' AS `a NOT REGEXP 'pattern'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v15 AS SELECT COUNT(*), SUM(a), AVG(a), MAX(a), MIN(a), COUNT(DISTINCT a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v15
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v15
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v15` (`COUNT(*)`, `SUM(a)`, `AVG(a)`, `MAX(a)`, `MIN(a)`, `COUNT(DISTINCT a)`) AS SELECT COUNT(1) AS `COUNT(*)`,SUM(`a`) AS `SUM(a)`,AVG(`a`) AS `AVG(a)`,MAX(`a`) AS `MAX(a)`,MIN(`a`) AS `MIN(a)`,COUNT(DISTINCT `a`) AS `COUNT(DISTINCT a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v16 AS SELECT CONCAT(a, b), SUBSTRING(a, 1, 3), TRIM(a), UPPER(a), LOWER(a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v16
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v16
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v16` (`CONCAT(a, b)`, `SUBSTRING(a, 1, 3)`, `TRIM(a)`, `UPPER(a)`, `LOWER(a)`) AS SELECT CONCAT(`a`, `b`) AS `CONCAT(a, b)`,SUBSTRING(`a`, 1, 3) AS `SUBSTRING(a, 1, 3)`,TRIM(`a`) AS `TRIM(a)`,UPPER(`a`) AS `UPPER(a)`,LOWER(`a`) AS `LOWER(a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v17 AS SELECT NOW(), CURRENT_TIMESTAMP, CURRENT_DATE, CURRENT_TIME, CURRENT_USER FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v17
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v17
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v17` (`NOW()`, `CURRENT_TIMESTAMP`, `CURRENT_DATE`, `CURRENT_TIME`, `CURRENT_USER`) AS SELECT NOW() AS `NOW()`,CURRENT_TIMESTAMP() AS `CURRENT_TIMESTAMP`,CURRENT_DATE() AS `CURRENT_DATE`,CURRENT_TIME() AS `CURRENT_TIME`,CURRENT_USER() AS `CURRENT_USER` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v18 AS SELECT IFNULL(a, 0), COALESCE(a, b, 0), NULLIF(a, 0), IF(a > 0, 'yes', 'no') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v18
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v18
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v18` (`IFNULL(a, 0)`, `COALESCE(a, b, 0)`, `NULLIF(a, 0)`, `IF(a > 0, 'yes', 'no')`) AS SELECT IFNULL(`a`, 0) AS `IFNULL(a, 0)`,COALESCE(`a`, `b`, 0) AS `COALESCE(a, b, 0)`,NULLIF(`a`, 0) AS `NULLIF(a, 0)`,IF(`a`>0, _UTF8MB4'yes', _UTF8MB4'no') AS `IF(a > 0, 'yes', 'no')` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v19 AS SELECT CAST(a AS CHAR), CAST(a AS SIGNED), CONVERT(a, CHAR), CONVERT(a USING utf8mb4) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v19
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v19
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v19` (`CAST(a AS CHAR)`, `CAST(a AS SIGNED)`, `CONVERT(a, CHAR)`, `CONVERT(a USING utf8mb4)`) AS SELECT CAST(`a` AS CHAR(11)) AS `CAST(a AS CHAR)`,CAST(`a` AS SIGNED) AS `CAST(a AS SIGNED)`,CONVERT(`a`, CHAR(11)) AS `CONVERT(a, CHAR)`,CONVERT(`a` USING 'utf8mb4') AS `CONVERT(a USING utf8mb4)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v20 AS SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v20
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v20
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v20` (`CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END`) AS SELECT CASE WHEN `a`>0 THEN _UTF8MB4'pos' WHEN `a`<0 THEN _UTF8MB4'neg' ELSE _UTF8MB4'zero' END AS `CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v21 AS SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v21
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v21
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v21` (`CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END`) AS SELECT CASE `a` WHEN 1 THEN _UTF8MB4'one' WHEN 2 THEN _UTF8MB4'two' ELSE _UTF8MB4'other' END AS `CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v22 AS SELECT (a + b) * c, a + (b * c) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v22
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v22
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v22` (`(a + b) * c`, `a + (b * c)`) AS SELECT (`a`+`b`)*`c` AS `(a + b) * c`,`a`+(`b`*`c`) AS `a + (b * c)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v23 AS SELECT -a, +a, !a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v23
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v23
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v23` (`-a`, `+a`, `!a`) AS SELECT -`a` AS `-a`,+`a` AS `+a`,!`a` AS `!a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v24 AS SELECT 1, 1.5, 'hello', NULL, TRUE, FALSE FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v24
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v24
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v24` (`1`, `1.5`, `hello`, `NULL`, `TRUE`, `FALSE`) AS SELECT 1 AS `1`,1.5 AS `1.5`,_UTF8MB4'hello' AS `hello`,NULL AS `NULL`,TRUE AS `TRUE`,FALSE AS `FALSE` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v25 AS SELECT 0xFF, X'FF', 0b1010, b'1010' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v25
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v25
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v25` (`0xFF`, `X'FF'`, `0b1010`, `b'1010'`) AS SELECT x'ff' AS `0xFF`,x'ff' AS `X'FF'`,b'1010' AS `0b1010`,b'1010' AS `b'1010'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v26 AS SELECT _utf8mb4'hello', _latin1'world' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v26
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v26
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v26` (`hello`, `world`) AS SELECT _UTF8MB4'hello' AS `hello`,_LATIN1'world' AS `world` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v27 AS SELECT DATE '2024-01-01', TIME '12:00:00', TIMESTAMP '2024-01-01 12:00:00' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v27
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v27
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v27` (`DATE '2024-01-01'`, `TIME '12:00:00'`, `TIMESTAMP '2024-01-01 12:00:00'`) AS SELECT DATE '2024-01-01' AS `DATE '2024-01-01'`,TIME '12:00:00' AS `TIME '12:00:00'`,TIMESTAMP '2024-01-01 12:00:00' AS `TIMESTAMP '2024-01-01 12:00:00'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v28 AS SELECT (SELECT MAX(a) FROM t) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v28
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v28
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v28` (`(SELECT MAX(a) FROM t)`) AS SELECT (SELECT MAX(`a`) AS `MAX(a)` FROM `test`.`t`) AS `(SELECT MAX(a) FROM t)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v29 AS SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)
< ok 0 0
> query
| SHOW CREATE VIEW v29
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v29
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v29` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t` WHERE `a` IN (SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v30 AS SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t WHERE a > 0)
< ok 0 0
> query
| SHOW CREATE VIEW v30
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v30
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v30` (`a`, `b`, `c`) AS SELECT `test`.`t`.`a` AS `a`,`test`.`t`.`b` AS `b`,`test`.`t`.`c` AS `c` FROM `test`.`t` WHERE EXISTS (SELECT 1 AS `1` FROM `test`.`t` WHERE `a`>0)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v31 AS SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v31
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v31
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v31` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v32 AS SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v32
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v32
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v32` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` LEFT JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v33 AS SELECT t1.a, t2.b FROM t1 RIGHT JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v33
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v33
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v33` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` RIGHT JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v34 AS SELECT t1.a, t2.b FROM t1 CROSS JOIN t2
< ok 0 0
> query
| SHOW CREATE VIEW v34
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v34
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v34` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` JOIN `test`.`t2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v35 AS SELECT * FROM t1 NATURAL JOIN t2
< ok 0 0
> query
| SHOW CREATE VIEW v35
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v35
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v35` (`a`, `b`) AS SELECT `test`.`t1`.`a` AS `a`,`test`.`t1`.`b` AS `b` FROM `test`.`t1` NATURAL JOIN `test`.`t2`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v36 AS SELECT t1.a, t2.b FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v36
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v36
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v36` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` STRAIGHT_JOIN `test`.`t2` ON `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v37 AS SELECT t1.a, t2.b FROM t1 JOIN t2 USING (a)
< ok 0 0
> query
| SHOW CREATE VIEW v37
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v37
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v37` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM `test`.`t1` JOIN `test`.`t2` USING (`a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v38 AS SELECT t1.a, t2.b FROM t1, t2 WHERE t1.a = t2.a
< ok 0 0
> query
| SHOW CREATE VIEW v38
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v38
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v38` (`a`, `b`) AS SELECT `t1`.`a` AS `a`,`t2`.`b` AS `b` FROM (`test`.`t1`) JOIN `test`.`t2` WHERE `t1`.`a`=`t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v39 AS SELECT a FROM t UNION SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v39
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v39
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v39` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v40 AS SELECT a FROM t UNION ALL SELECT b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v40
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v40
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v40` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION ALL SELECT `b` AS `b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v41 AS SELECT a FROM t UNION SELECT b FROM t UNION SELECT c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v41
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v41
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v41` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` UNION SELECT `b` AS `b` FROM `test`.`t` UNION SELECT `c` AS `c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v42 AS SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v42
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v42
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v42` (`a`, `ROW_NUMBER() OVER (ORDER BY a)`) AS SELECT `a` AS `a`,ROW_NUMBER() OVER (ORDER BY `a`) AS `ROW_NUMBER() OVER (ORDER BY a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v43 AS SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v43
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v43
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v43` (`a`, `SUM(b) OVER (PARTITION BY a ORDER BY b)`) AS SELECT `a` AS `a`,SUM(`b`) OVER (PARTITION BY `a` ORDER BY `b`) AS `SUM(b) OVER (PARTITION BY a ORDER BY b)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v44 AS SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v44
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v44
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v44` (`a`, `name_exp_2`) AS SELECT `a` AS `a`,SUM(`b`) OVER (PARTITION BY `a` ORDER BY `b` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS `SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v45 AS SELECT a AS col1, b col2 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v45
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v45
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v45` (`col1`, `col2`) AS SELECT `a` AS `col1`,`b` AS `col2` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v46 AS SELECT a + b AS sum_col FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v46
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v46
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v46` (`sum_col`) AS SELECT `a`+`b` AS `sum_col` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v47 AS SELECT * FROM t AS t1
< ok 0 0
> query
| SHOW CREATE VIEW v47
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v47
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v47` (`a`, `b`, `c`) AS SELECT `test`.`t1`.`a` AS `a`,`test`.`t1`.`b` AS `b`,`test`.`t1`.`c` AS `c` FROM `test`.`t` AS `t1`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v48 AS SELECT a->>'$.key' FROM t
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_extract; a JSON string or JSON type is required.
> exec
| CREATE VIEW v49 AS SELECT INTERVAL 1 DAY + a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v49
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v49
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v49` (`INTERVAL 1 DAY + a`) AS SELECT DATE_ADD(`a`, INTERVAL 1 DAY) AS `INTERVAL 1 DAY + a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v50 AS SELECT a NOT LIKE 'x%' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v50
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v50
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v50` (`a NOT LIKE 'x%'`) AS SELECT `a` NOT LIKE _UTF8MB4'x%' AS `a NOT LIKE 'x%'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v51 AS SELECT a != b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v51
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v51
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v51` (`a != b`) AS SELECT `a`!=`b` AS `a != b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v52 AS SELECT COUNT(*) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v52
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v52
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v52` (`COUNT(*)`) AS SELECT COUNT(1) AS `COUNT(*)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v53 AS SELECT SUBSTRING('abc', 1, 2) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v53
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v53
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v53` (`SUBSTRING('abc', 1, 2)`) AS SELECT SUBSTRING(_UTF8MB4'abc', 1, 2) AS `SUBSTRING('abc', 1, 2)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v54 AS SELECT CURRENT_TIMESTAMP() FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v54
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v54
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v54` (`CURRENT_TIMESTAMP()`) AS SELECT CURRENT_TIMESTAMP() AS `CURRENT_TIMESTAMP()` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v55 AS SELECT a MOD b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v55
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v55
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v55` (`a MOD b`) AS SELECT `a`%`b` AS `a MOD b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v56 AS SELECT +a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v56
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v56
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v56` (`+a`) AS SELECT +`a` AS `+a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v57 AS SELECT !a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v57
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v57
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v57` (`!a`) AS SELECT !`a` AS `!a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v58 AS SELECT a AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v58
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v58
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v58` (`a AND b`) AS SELECT `a` AND `b` AS `a AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v59 AS SELECT a REGEXP 'x' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v59
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v59
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v59` (`a REGEXP 'x'`) AS SELECT `a` REGEXP _UTF8MB4'x' AS `a REGEXP 'x'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v60 AS SELECT a->'$.key' FROM t
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_extract; a JSON string or JSON type is required.
> exec
| CREATE VIEW v61 AS SELECT d.x FROM (SELECT a AS x FROM t) d
< ok 0 0
> query
| SHOW CREATE VIEW v61
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v61
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v61` (`x`) AS SELECT `d`.`x` AS `x` FROM (SELECT `a` AS `x` FROM `test`.`t`) AS `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v62 AS SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10
< ok 0 0
> query
| SHOW CREATE VIEW v62
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v62
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v62` (`x`) AS SELECT `d`.`x` AS `x` FROM (SELECT `a` AS `x` FROM `test`.`t` WHERE `a`>0) AS `d` WHERE `d`.`x`<10
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v63 AS SELECT GROUP_CONCAT(a ORDER BY a SEPARATOR ',') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v63
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v63
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v63` (`GROUP_CONCAT(a ORDER BY a SEPARATOR ',')`) AS SELECT GROUP_CONCAT(`a` ORDER BY `a` SEPARATOR ',') AS `GROUP_CONCAT(a ORDER BY a SEPARATOR ',')` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v64 AS SELECT GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v64
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v64
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v64` (`GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')`) AS SELECT GROUP_CONCAT(DISTINCT `a` ORDER BY `a` DESC SEPARATOR ';') AS `GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v65 AS SELECT CONCAT(a, b, c) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v65
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v65
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v65` (`CONCAT(a, b, c)`) AS SELECT CONCAT(`a`, `b`, `c`) AS `CONCAT(a, b, c)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v66 AS SELECT a IN (1, 2, 3) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v66
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v66
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v66` (`a IN (1, 2, 3)`) AS SELECT `a` IN (1,2,3) AS `a IN (1, 2, 3)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v67 AS SELECT a AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v67
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v67
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v67` (`x`) AS SELECT `a` AS `x` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v68 AS SELECT x.a FROM t AS x
< ok 0 0
> query
| SHOW CREATE VIEW v68
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v68
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v68` (`a`) AS SELECT `x`.`a` AS `a` FROM `test`.`t` AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v69 AS SELECT x.a FROM t x
< ok 0 0
> query
| SHOW CREATE VIEW v69
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v69
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v69` (`a`) AS SELECT `x`.`a` AS `a` FROM `test`.`t` AS `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v70 AS SELECT a + b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v70
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v70
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v70` (`a + b + c`) AS SELECT `a`+`b`+`c` AS `a + b + c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v71 AS SELECT a * b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v71
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v71
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v71` (`a * b + c`) AS SELECT `a`*`b`+`c` AS `a * b + c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v72 AS SELECT a + b * c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v72
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v72
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v72` (`a + b * c`) AS SELECT `a`+`b`*`c` AS `a + b * c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v73 AS SELECT a OR b AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v73
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v73
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v73` (`a OR b AND c`) AS SELECT `a` OR `b` AND `c` AS `a OR b AND c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v74 AS SELECT (a OR b) AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v74
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v74
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v74` (`(a OR b) AND c`) AS SELECT (`a` OR `b`) AND `c` AS `(a OR b) AND c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v75 AS SELECT a > 0 AND b < 10 OR c = 5 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v75
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v75
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v75` (`a > 0 AND b < 10 OR c = 5`) AS SELECT `a`>0 AND `b`<10 OR `c`=5 AS `a > 0 AND b < 10 OR c = 5` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v76 AS WITH cte AS (SELECT a FROM t) SELECT * FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v76
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v76
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v76` (`a`) AS WITH `cte` AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `cte`.`a` AS `a` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v77 AS WITH cte(x) AS (SELECT a FROM t) SELECT x FROM cte
< ok 0 0
> query
| SHOW CREATE VIEW v77
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v77
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v77` (`x`) AS WITH `cte` (`x`) AS (SELECT `a` AS `a` FROM `test`.`t`) SELECT `x` AS `x` FROM `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v78 AS SELECT (a + 1) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v78
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v78
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v78` (`(a + 1) AND b`) AS SELECT (`a`+1) AND `b` AS `(a + 1) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v79 AS SELECT (a > 0) AND (b + 1) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v79
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v79
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v79` (`(a > 0) AND (b + 1)`) AS SELECT (`a`>0) AND (`b`+1) AS `(a > 0) AND (b + 1)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v80 AS SELECT (a > 0) AND (b > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v80
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v80
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v80` (`(a > 0) AND (b > 0)`) AS SELECT (`a`>0) AND (`b`>0) AS `(a > 0) AND (b > 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v81 AS SELECT ABS(a) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v81
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v81
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v81` (`ABS(a) AND b`) AS SELECT ABS(`a`) AND `b` AS `ABS(a) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v82 AS SELECT COUNT(*) > 0 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v82
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v82
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v82` (`COUNT(*) > 0`) AS SELECT COUNT(1)>0 AS `COUNT(*) > 0` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v83 AS SELECT CASE WHEN a > 0 THEN 1 ELSE 0 END AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v83
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v83
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v83` (`CASE WHEN a > 0 THEN 1 ELSE 0 END AND b`) AS SELECT CASE WHEN `a`>0 THEN 1 ELSE 0 END AND `b` AS `CASE WHEN a > 0 THEN 1 ELSE 0 END AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v84 AS SELECT NOT (a + 1) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v84
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v84
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v84` (`NOT (a + 1)`) AS SELECT NOT (`a`+1) AS `NOT (a + 1)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v85 AS SELECT NOT (a > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v85
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v85
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v85` (`NOT (a > 0)`) AS SELECT NOT (`a`>0) AS `NOT (a > 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v86 AS SELECT IF(a > 0, 1, 0) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v86
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v86
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v86` (`IF(a > 0, 1, 0) AND b`) AS SELECT IF(`a`>0, 1, 0) AND `b` AS `IF(a > 0, 1, 0) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v87 AS SELECT (SELECT MAX(a) FROM t) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v87
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v87
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v87` (`(SELECT MAX(a) FROM t) AND b`) AS SELECT (SELECT MAX(`a`) AS `MAX(a)` FROM `test`.`t`) AND `b` AS `(SELECT MAX(a) FROM t) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v88 AS SELECT CAST(a AS CHAR(10)) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v88
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v88
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v88` (`CAST(a AS CHAR(10))`) AS SELECT CAST(`a` AS CHAR(10)) AS `CAST(a AS CHAR(10))` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v89 AS SELECT CAST(a AS BINARY) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v89
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v89
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v89` (`CAST(a AS BINARY)`) AS SELECT CAST(`a` AS BINARY(11)) AS `CAST(a AS BINARY)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v90 AS SELECT CAST(a AS DECIMAL(10,2)) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v90
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v90
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v90` (`CAST(a AS DECIMAL(10,2))`) AS SELECT CAST(`a` AS DECIMAL(10, 2)) AS `CAST(a AS DECIMAL(10,2))` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v91 AS SELECT CAST(a AS UNSIGNED) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v91
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v91
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v91` (`CAST(a AS UNSIGNED)`) AS SELECT CAST(`a` AS UNSIGNED) AS `CAST(a AS UNSIGNED)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v92 AS SELECT CAST(a AS DATE) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v92
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v92
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v92` (`CAST(a AS DATE)`) AS SELECT CAST(`a` AS DATE) AS `CAST(a AS DATE)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v93 AS SELECT CAST(a AS DATETIME) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v93
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v93
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v93` (`CAST(a AS DATETIME)`) AS SELECT CAST(`a` AS DATETIME) AS `CAST(a AS DATETIME)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v94 AS SELECT CAST(a AS JSON) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v94
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v94
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v94` (`CAST(a AS JSON)`) AS SELECT CAST(`a` AS JSON) AS `CAST(a AS JSON)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v95_setup AS SELECT 'hello' AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v95_setup
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v95_setup
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v95_setup` (`'hello' AND 1`) AS SELECT _UTF8MB4'hello' AND 1 AS `'hello' AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v96 AS SELECT a AND 'hello' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v96
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v96
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v96` (`a AND 'hello'`) AS SELECT `a` AND _UTF8MB4'hello' AS `a AND 'hello'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v97 AS SELECT CONCAT(a,b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v97
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v97
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v97` (`CONCAT(a,b) AND 1`) AS SELECT CONCAT(`a`, `b`) AND 1 AS `CONCAT(a,b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v98 AS SELECT IFNULL(a, 0) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v98
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v98
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v98` (`IFNULL(a, 0) AND b`) AS SELECT IFNULL(`a`, 0) AND `b` AS `IFNULL(a, 0) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v99 AS SELECT COALESCE(a, b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v99
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v99
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v99` (`COALESCE(a, b) AND 1`) AS SELECT COALESCE(`a`, `b`) AND 1 AS `COALESCE(a, b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v100 AS SELECT NULLIF(a, 0) AND b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v100
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v100
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v100` (`NULLIF(a, 0) AND b`) AS SELECT NULLIF(`a`, 0) AND `b` AS `NULLIF(a, 0) AND b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v101 AS SELECT GREATEST(a, b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v101
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v101
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v101` (`GREATEST(a, b) AND 1`) AS SELECT GREATEST(`a`, `b`) AND 1 AS `GREATEST(a, b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v102 AS SELECT LEAST(a, b) AND 1 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v102
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v102
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v102` (`LEAST(a, b) AND 1`) AS SELECT LEAST(`a`, `b`) AND 1 AS `LEAST(a, b) AND 1` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v103 AS SELECT (a = b) AND (a > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v103
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v103
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v103` (`(a = b) AND (a > 0)`) AS SELECT (`a`=`b`) AND (`a`>0) AS `(a = b) AND (a > 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v104 AS SELECT (a IN (1,2,3)) AND (b BETWEEN 1 AND 10) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v104
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v104
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v104` (`(a IN (1,2,3)) AND (b BETWEEN 1 AND 10)`) AS SELECT (`a` IN (1,2,3)) AND (`b` BETWEEN 1 AND 10) AS `(a IN (1,2,3)) AND (b BETWEEN 1 AND 10)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v105 AS SELECT (a IS NULL) AND (b LIKE 'x%') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v105
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v105
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v105` (`(a IS NULL) AND (b LIKE 'x%')`) AS SELECT (`a` IS NULL) AND (`b` LIKE _UTF8MB4'x%') AS `(a IS NULL) AND (b LIKE 'x%')` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v106 AS SELECT (a = 1) OR (b = 2) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v106
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v106
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v106` (`(a = 1) OR (b = 2)`) AS SELECT (`a`=1) OR (`b`=2) AS `(a = 1) OR (b = 2)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v107 AS SELECT EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v107
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v107
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v107` (`EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0)`) AS SELECT EXISTS (SELECT 1 AS `1` FROM `test`.`t` WHERE `a`>0) AND (`b`>0) AS `EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v108 AS SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE VIEW v109 AS SELECT a FROM t WHERE a > 0
< ok 0 0
> query
| SHOW CREATE VIEW v109
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v109
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v109` (`a`) AS SELECT `a` AS `a` FROM `test`.`t` WHERE `a`>0
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v110 AS SELECT t.a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v110
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v110
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v110` (`a`) AS SELECT `t`.`a` AS `a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v111 AS SELECT '' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v111
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v111
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v111` (``) AS SELECT _UTF8MB4'' FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v112 AS SELECT 'it''s' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v112
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v112
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v112` (`it's`) AS SELECT _UTF8MB4'it''s' AS `it's` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v113 AS SELECT 'back\\slash' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v113
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v113
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v113` (`back\slash`) AS SELECT _UTF8MB4'back\\slash' AS `back\slash` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v114 AS SELECT a + b + c + a + b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v114
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v114
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v114` (`a + b + c + a + b + c`) AS SELECT `a`+`b`+`c`+`a`+`b`+`c` AS `a + b + c + a + b + c` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v115 AS SELECT CONCAT(UPPER(TRIM(a)), LOWER(b)) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v115
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v115
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v115` (`CONCAT(UPPER(TRIM(a)), LOWER(b))`) AS SELECT CONCAT(UPPER(TRIM(`a`)), LOWER(`b`)) AS `CONCAT(UPPER(TRIM(a)), LOWER(b))` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE VIEW v116 AS SELECT COUNT(*), SUM(a), AVG(b), MAX(c) FROM t GROUP BY a
< ok 0 0
> query
| SHOW CREATE VIEW v116
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v116
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v116` (`COUNT(*)`, `SUM(a)`, `AVG(b)`, `MAX(c)`) AS SELECT COUNT(1) AS `COUNT(*)`,SUM(`a`) AS `SUM(a)`,AVG(`b`) AS `AVG(b)`,MAX(`c`) AS `MAX(c)` FROM `test`.`t` GROUP BY `a`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
