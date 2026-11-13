-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMySQL_DeparseRules default#1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v1` AS select 1 AS `1`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v2` AS select `t`.`a` AS `a` from `t` where `t`.`a` > 1 group by `t`.`a` having count(0) > 1 order by `t`.`a` limit 10
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v3` AS select `t`.`a` AS `a`,`t`.`b` AS `alias1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v4` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v5 AS SELECT `select` FROM t
< error
| 1054 42S22 Unknown column 'select' in 'SELECT'
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v6` AS select `t`.`a` + `t`.`b` AS `a+b`,`t`.`a` - `t`.`b` AS `a-b`,`t`.`a` * `t`.`b` AS `a*b`,`t`.`a` / `t`.`b` AS `a/b`,`t`.`a` DIV `t`.`b` AS `a DIV b`,`t`.`a` MOD `t`.`b` AS `a MOD b`,`t`.`a` MOD `t`.`b` AS `a%b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v7` AS select `t`.`a` = `t`.`b` AS `a=b`,`t`.`a` <> `t`.`b` AS `a<>b`,`t`.`a` <> `t`.`b` AS `a!=b`,`t`.`a` > `t`.`b` AS `a>b`,`t`.`a` < `t`.`b` AS `a<b`,`t`.`a` >= `t`.`b` AS `a>=b`,`t`.`a` <= `t`.`b` AS `a<=b`,`t`.`a` <=> `t`.`b` AS `a<=>b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v8` AS select `t`.`a` <> 0 and `t`.`b` <> 0 AS `a AND b`,`t`.`a` <> 0 or `t`.`b` <> 0 AS `a OR b`,`t`.`a` = 0 AS `NOT a`,`t`.`a` xor `t`.`b` AS `a XOR b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v9` AS select `t`.`a` | `t`.`b` AS `a|b`,`t`.`a` & `t`.`b` AS `a&b`,`t`.`a` ^ `t`.`b` AS `a^b`,`t`.`a` << `t`.`b` AS `a<<b`,`t`.`a` >> `t`.`b` AS `a>>b`,~`t`.`a` AS `~a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v10` AS select `t`.`a` is null AS `a IS NULL`,`t`.`a` is not null AS `a IS NOT NULL`,`t`.`a` is true AS `a IS TRUE`,`t`.`a` is false AS `a IS FALSE` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v11` AS select `t`.`a` in (1,2,3) AS `a IN (1,2,3)`,`t`.`a` not in (1,2,3) AS `a NOT IN (1,2,3)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v12` AS select `t`.`a` between 1 and 10 AS `a BETWEEN 1 AND 10`,`t`.`a` not between 1 and 10 AS `a NOT BETWEEN 1 AND 10` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v13` AS select `t`.`a` like 'foo%' AS `a LIKE 'foo%'`,`t`.`a`  not like 'bar%' AS `a NOT LIKE 'bar%'`,`t`.`a` like 'x' escape '\\' AS `a LIKE 'x' ESCAPE '\\'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v14` AS select `t`.`a` regexp 'pattern' AS `a REGEXP 'pattern'`,!(`t`.`a` regexp 'pattern') AS `a NOT REGEXP 'pattern'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v15` AS select count(0) AS `COUNT(*)`,sum(`t`.`a`) AS `SUM(a)`,avg(`t`.`a`) AS `AVG(a)`,max(`t`.`a`) AS `MAX(a)`,min(`t`.`a`) AS `MIN(a)`,count(distinct `t`.`a`) AS `COUNT(DISTINCT a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v16` AS select concat(`t`.`a`,`t`.`b`) AS `CONCAT(a, b)`,substr(`t`.`a`,1,3) AS `SUBSTRING(a, 1, 3)`,trim(`t`.`a`) AS `TRIM(a)`,ucase(`t`.`a`) AS `UPPER(a)`,lcase(`t`.`a`) AS `LOWER(a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v17` AS select current_timestamp() AS `NOW()`,current_timestamp() AS `CURRENT_TIMESTAMP`,curdate() AS `CURRENT_DATE`,curtime() AS `CURRENT_TIME`,current_user() AS `CURRENT_USER` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v18` AS select ifnull(`t`.`a`,0) AS `IFNULL(a, 0)`,coalesce(`t`.`a`,`t`.`b`,0) AS `COALESCE(a, b, 0)`,nullif(`t`.`a`,0) AS `NULLIF(a, 0)`,if(`t`.`a` > 0,'yes','no') AS `IF(a > 0, 'yes', 'no')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v19` AS select cast(`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)`,cast(`t`.`a` as signed) AS `CAST(a AS SIGNED)`,cast(`t`.`a` as char charset utf8mb4) AS `CONVERT(a, CHAR)`,convert(`t`.`a` using utf8mb4) AS `CONVERT(a USING utf8mb4)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v20` AS select case when `t`.`a` > 0 then 'pos' when `t`.`a` < 0 then 'neg' else 'zero' end AS `CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v21` AS select case `t`.`a` when 1 then 'one' when 2 then 'two' else 'other' end AS `CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v22` AS select (`t`.`a` + `t`.`b`) * `t`.`c` AS `(a + b) * c`,`t`.`a` + `t`.`b` * `t`.`c` AS `a + (b * c)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v23` AS select -`t`.`a` AS `-a`,`t`.`a` AS `a`,`t`.`a` = 0 AS `!a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v24` AS select 1 AS `1`,1.5 AS `1.5`,'hello' AS `hello`,NULL AS `NULL`,1 AS `TRUE`,0 AS `FALSE` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v25` AS select 0xff AS `0xFF`,X'ff' AS `X'FF'`,0x0a AS `0b1010`,0x0a AS `b'1010'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v26` AS select _utf8mb4'hello' AS `hello`,_latin1'world' AS `world` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v27` AS select DATE'2024-01-01' AS `DATE '2024-01-01'`,TIME'12:00:00' AS `TIME '12:00:00'`,TIMESTAMP'2024-01-01 12:00:00' AS `TIMESTAMP '2024-01-01 12:00:00'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v28` AS select (select max(`t`.`a`) from `t`) AS `(SELECT MAX(a) FROM t)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v29` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where `t`.`a` in (select `t`.`a` from `t` where `t`.`a` > 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v30` AS select `t`.`a` AS `a`,`t`.`b` AS `b`,`t`.`c` AS `c` from `t` where exists(select 1 from `t` where `t`.`a` > 0 limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v31` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2` on(`t1`.`a` = `t2`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v32` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` left join `t2` on(`t1`.`a` = `t2`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v33` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t2` left join `t1` on(`t1`.`a` = `t2`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v34` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2`)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v35` AS select `t1`.`a` AS `a`,`t1`.`b` AS `b` from (`t1` join `t2` on(`t1`.`a` = `t2`.`a` and `t1`.`b` = `t2`.`b`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v36` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` straight_join `t2` on(`t1`.`a` = `t2`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v37` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2` on(`t1`.`a` = `t2`.`a`))
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v38` AS select `t1`.`a` AS `a`,`t2`.`b` AS `b` from (`t1` join `t2`) where `t1`.`a` = `t2`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v39` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v40` AS select `t`.`a` AS `a` from `t` union all select `t`.`b` AS `b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v41` AS select `t`.`a` AS `a` from `t` union select `t`.`b` AS `b` from `t` union select `t`.`c` AS `c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v42` AS select `t`.`a` AS `a`,row_number() over ( order by `t`.`a`) AS `ROW_NUMBER() OVER (ORDER BY a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v43` AS select `t`.`a` AS `a`,sum(`t`.`b`) over ( partition by `t`.`a` order by `t`.`b`) AS `SUM(b) OVER (PARTITION BY a ORDER BY b)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v44` AS select `t`.`a` AS `a`,sum(`t`.`b`) over ( partition by `t`.`a` order by `t`.`b` rows between  unbounded  preceding  and  current row ) AS `Name_exp_2` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v45` AS select `t`.`a` AS `col1`,`t`.`b` AS `col2` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v46` AS select `t`.`a` + `t`.`b` AS `sum_col` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v47` AS select `t1`.`a` AS `a`,`t1`.`b` AS `b`,`t1`.`c` AS `c` from `t` `t1`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v48 AS SELECT a->>'$.key' FROM t
< error
| 1064 42000 …near '>>'$.key' FROM t' at line 1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v49` AS select `t`.`a` + interval 1 day AS `INTERVAL 1 DAY + a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v50` AS select `t`.`a`  not like 'x%' AS `a NOT LIKE 'x%'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v51` AS select `t`.`a` <> `t`.`b` AS `a != b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v52` AS select count(0) AS `COUNT(*)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v53` AS select substr('abc',1,2) AS `SUBSTRING('abc', 1, 2)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v54` AS select current_timestamp() AS `CURRENT_TIMESTAMP()` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v55` AS select `t`.`a` MOD `t`.`b` AS `a MOD b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v56` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v57` AS select `t`.`a` = 0 AS `!a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v58` AS select `t`.`a` <> 0 and `t`.`b` <> 0 AS `a AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v59` AS select `t`.`a` regexp 'x' AS `a REGEXP 'x'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v60 AS SELECT a->'$.key' FROM t
< error
| 1064 42000 …near '>'$.key' FROM t' at line 1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v61` AS select `d`.`x` AS `x` from (select `t`.`a` AS `x` from `t`) `d`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v62` AS select `d`.`x` AS `x` from (select `t`.`a` AS `x` from `t` where `t`.`a` > 0) `d` where `d`.`x` < 10
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v63` AS select group_concat(`t`.`a` order by `t`.`a` ASC separator ',') AS `GROUP_CONCAT(a ORDER BY a SEPARATOR ',')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v64` AS select group_concat(distinct `t`.`a` order by `t`.`a` DESC separator ';') AS `GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v65` AS select concat(`t`.`a`,`t`.`b`,`t`.`c`) AS `CONCAT(a, b, c)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v66` AS select `t`.`a` in (1,2,3) AS `a IN (1, 2, 3)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v67` AS select `t`.`a` AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v68` AS select `x`.`a` AS `a` from `t` `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v69` AS select `x`.`a` AS `a` from `t` `x`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v70` AS select `t`.`a` + `t`.`b` + `t`.`c` AS `a + b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v71` AS select `t`.`a` * `t`.`b` + `t`.`c` AS `a * b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v72` AS select `t`.`a` + `t`.`b` * `t`.`c` AS `a + b * c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v73` AS select `t`.`a` <> 0 or `t`.`b` <> 0 and `t`.`c` <> 0 AS `a OR b AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v74` AS select (`t`.`a` <> 0 or `t`.`b` <> 0) and `t`.`c` <> 0 AS `(a OR b) AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v75` AS select `t`.`a` > 0 and `t`.`b` < 10 or `t`.`c` = 5 AS `a > 0 AND b < 10 OR c = 5` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v76` AS with cte as (select `t`.`a` AS `a` from `t`)select `cte`.`a` AS `a` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v77` AS with cte(`x`) as (select `t`.`a` AS `x` from `t`)select `cte`.`x` AS `x` from `cte`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v78` AS select `t`.`a` + 1 and `t`.`b` <> 0 AS `(a + 1) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v79` AS select `t`.`a` > 0 and `t`.`b` + 1 AS `(a > 0) AND (b + 1)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v80` AS select `t`.`a` > 0 and `t`.`b` > 0 AS `(a > 0) AND (b > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v81` AS select abs(`t`.`a`) and `t`.`b` <> 0 AS `ABS(a) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v82` AS select count(0) > 0 AS `COUNT(*) > 0` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v83` AS select case when `t`.`a` > 0 then 1 else 0 end and `t`.`b` <> 0 AS `CASE WHEN a > 0 THEN 1 ELSE 0 END AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v84` AS select !(`t`.`a` + 1) AS `NOT (a + 1)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v85` AS select `t`.`a` <= 0 AS `NOT (a > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v86` AS select if(`t`.`a` > 0,1,0) and `t`.`b` <> 0 AS `IF(a > 0, 1, 0) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v87` AS select (select max(`t`.`a`) from `t`) and `t`.`b` <> 0 AS `(SELECT MAX(a) FROM t) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v88` AS select cast(`t`.`a` as char(10) charset utf8mb4) AS `CAST(a AS CHAR(10))` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v89` AS select cast(`t`.`a` as char charset binary) AS `CAST(a AS BINARY)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v90` AS select cast(`t`.`a` as decimal(10,2)) AS `CAST(a AS DECIMAL(10,2))` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v91` AS select cast(`t`.`a` as unsigned) AS `CAST(a AS UNSIGNED)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v92` AS select cast(`t`.`a` as date) AS `CAST(a AS DATE)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v93` AS select cast(`t`.`a` as datetime) AS `CAST(a AS DATETIME)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v94 AS SELECT CAST(a AS JSON) FROM t
< error
| 1064 42000 …near 'JSON) FROM t' at line 1
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v95_setup` AS select 'hello' and 1 AS `'hello' AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v96` AS select `t`.`a` <> 0 and 'hello' AS `a AND 'hello'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v97` AS select concat(`t`.`a`,`t`.`b`) and 1 AS `CONCAT(a,b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v98` AS select ifnull(`t`.`a`,0) and `t`.`b` <> 0 AS `IFNULL(a, 0) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v99` AS select coalesce(`t`.`a`,`t`.`b`) and 1 AS `COALESCE(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v100` AS select nullif(`t`.`a`,0) and `t`.`b` <> 0 AS `NULLIF(a, 0) AND b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v101` AS select greatest(`t`.`a`,`t`.`b`) and 1 AS `GREATEST(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v102` AS select least(`t`.`a`,`t`.`b`) and 1 AS `LEAST(a, b) AND 1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v103` AS select `t`.`a` = `t`.`b` and `t`.`a` > 0 AS `(a = b) AND (a > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v104` AS select `t`.`a` in (1,2,3) and `t`.`b` between 1 and 10 AS `(a IN (1,2,3)) AND (b BETWEEN 1 AND 10)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v105` AS select `t`.`a` is null and `t`.`b` like 'x%' AS `(a IS NULL) AND (b LIKE 'x%')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v106` AS select `t`.`a` = 1 or `t`.`b` = 2 AS `(a = 1) OR (b = 2)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v107` AS select exists(select 1 from `t` where `t`.`a` > 0 limit 1) and `t`.`b` > 0 AS `EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v109` AS select `t`.`a` AS `a` from `t` where `t`.`a` > 0
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v110` AS select `t`.`a` AS `a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v111` AS select '' AS `Name_exp_1` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v112` AS select 'it\'s' AS `it's` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v113` AS select 'back\\slash' AS `back\slash` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v114` AS select `t`.`a` + `t`.`b` + `t`.`c` + `t`.`a` + `t`.`b` + `t`.`c` AS `a + b + c + a + b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v115` AS select concat(ucase(trim(`t`.`a`)),lcase(`t`.`b`)) AS `CONCAT(UPPER(TRIM(a)), LOWER(b))` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v116` AS select count(0) AS `COUNT(*)`,sum(`t`.`a`) AS `SUM(a)`,avg(`t`.`b`) AS `AVG(b)`,max(`t`.`c`) AS `MAX(c)` from `t` group by `t`.`a`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
