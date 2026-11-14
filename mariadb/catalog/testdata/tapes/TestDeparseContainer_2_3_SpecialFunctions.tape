-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseContainer_2_3_SpecialFunctions default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| CREATE TABLE IF NOT EXISTS tv (a VARCHAR(50), b VARCHAR(50), c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_trim_simple
< ok 0 0
> exec
| CREATE VIEW v_trim_simple AS SELECT TRIM(a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_simple
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_simple
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_simple` AS select trim(`tv`.`a`) AS `TRIM(a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_trim_leading
< ok 0 0
> exec
| CREATE VIEW v_trim_leading AS SELECT TRIM(LEADING 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_leading
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_leading
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_leading` AS select trim(leading 'x' from `tv`.`a`) AS `TRIM(LEADING 'x' FROM a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_trim_trailing
< ok 0 0
> exec
| CREATE VIEW v_trim_trailing AS SELECT TRIM(TRAILING 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_trailing
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_trailing
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_trailing` AS select trim(trailing 'x' from `tv`.`a`) AS `TRIM(TRAILING 'x' FROM a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_trim_both
< ok 0 0
> exec
| CREATE VIEW v_trim_both AS SELECT TRIM(BOTH 'x' FROM a) FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_both
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_both
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_both` AS select trim(both 'x' from `tv`.`a`) AS `TRIM(BOTH 'x' FROM a)` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_gc_basic
< ok 0 0
> exec
| CREATE VIEW v_gc_basic AS SELECT GROUP_CONCAT(a ORDER BY a SEPARATOR ',') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_gc_basic
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_gc_basic
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_gc_basic` AS select group_concat(`t`.`a` order by `t`.`a` ASC separator ',') AS `GROUP_CONCAT(a ORDER BY a SEPARATOR ',')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_gc_full
< ok 0 0
> exec
| CREATE VIEW v_gc_full AS SELECT GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';') FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_gc_full
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_gc_full
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_gc_full` AS select group_concat(distinct `t`.`a` order by `t`.`a` DESC separator ';') AS `GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_simple_case
< ok 0 0
> exec
| CREATE VIEW v_simple_case AS SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_simple_case
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_simple_case
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_simple_case` AS select case `t`.`a` when 1 then 'one' when 2 then 'two' else 'other' end AS `CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_trim_oracle
< ok 0 0
> exec
| CREATE VIEW v_trim_oracle AS SELECT TRIM_ORACLE(a) AS c1, TRIM_ORACLE(LEADING 'x' FROM a) AS c2, TRIM_ORACLE('x' FROM a) AS c3, TRIM_ORACLE(TRAILING FROM a) AS c4 FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_oracle
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_oracle
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_oracle` AS select oracle_schema.trim(`tv`.`a`) AS `c1`,oracle_schema.trim(leading 'x' from `tv`.`a`) AS `c2`,oracle_schema.trim(both 'x' from `tv`.`a`) AS `c3`,oracle_schema.rtrim(`tv`.`a`) AS `c4` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_trim_mariadb_schema
< ok 0 0
> exec
| CREATE VIEW v_trim_mariadb_schema AS SELECT mariadb_schema.trim(LEADING 'x' FROM a) AS c1, mariadb_schema.trim('x' FROM a) AS c2, mariadb_schema.trim(a) AS c3 FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_trim_mariadb_schema
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_trim_mariadb_schema
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_trim_mariadb_schema` AS select trim(leading 'x' from `tv`.`a`) AS `c1`,trim(both 'x' from `tv`.`a`) AS `c2`,trim(`tv`.`a`) AS `c3` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_value_func
< ok 0 0
> exec
| CREATE VIEW v_value_func AS SELECT VALUE(a) AS c1, VALUE(t.b) AS c2 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_value_func
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_value_func
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_value_func` AS select value(`t`.`a`) AS `c1`,value(`t`.`b`) AS `c2` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_dyncol_create
< ok 0 0
> exec
| CREATE VIEW v_dyncol_create AS SELECT COLUMN_CREATE('a', 1 AS INT, 'b', 'x' AS CHAR CHARACTER SET latin1, 'c', 1.5 AS DECIMAL(5,2), 'd', a AS UNSIGNED INT, 'e', 1 AS REAL, 'f', '10:00' AS TIME(2), 'g', 'y' AS NCHAR, 'h', 'z' AS CHAR BYTE, 'i', 2) AS c1 FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_dyncol_create
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_dyncol_create
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dyncol_create` AS select column_create('a',1 AS int,'b','x' AS char charset latin1 collate latin1_swedish_ci ,'c',1.5 AS decimal,'d',`tv`.`a` AS unsigned int,'e',1 AS double,'f','10:00' AS time,'g','y' AS char charset utf8mb3 collate utf8mb3_general_ci ,'h','z' AS char charset binary ,'i',2) AS `c1` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_dyncol_get
< ok 0 0
> exec
| CREATE VIEW v_dyncol_get AS SELECT COLUMN_ADD(a, 'x', 1 AS DOUBLE, 'y', 'q' AS CHAR CHARSET latin1 COLLATE latin1_bin) AS c1, COLUMN_DELETE(a, 'a', 'b') AS c2, COLUMN_GET(a, 'a' AS INT) AS c3, COLUMN_GET(a, 'a' AS CHAR(3)) AS c4, COLUMN_GET(a, 'a' AS DECIMAL(10,2)) AS c5 FROM tv
< ok 0 0
> query
| SHOW CREATE VIEW v_dyncol_get
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_dyncol_get
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_dyncol_get` AS select column_add(`tv`.`a`,'x',1 AS double,'y','q' AS char charset latin1 collate latin1_bin ) AS `c1`,column_add(`tv`.`a`,'a',NULL AS int,'b',NULL AS int) AS `c2`,column_get(`tv`.`a`,'a' as signed) AS `c3`,column_get(`tv`.`a`,'a' as char(3) charset utf8mb4) AS `c4`,column_get(`tv`.`a`,'a' as decimal(10,2)) AS `c5` from `tv`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
