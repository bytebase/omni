-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestDeparseContainer_7_1_WindowFunctions default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_wf_rownum
< ok 0 0
> exec
| CREATE VIEW v_wf_rownum AS SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_rownum
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_rownum
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_rownum` (`a`, `ROW_NUMBER() OVER (ORDER BY a)`) AS SELECT `a` AS `a`,ROW_NUMBER() OVER (ORDER BY `a`) AS `ROW_NUMBER() OVER (ORDER BY a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_wf_sum_part
< ok 0 0
> exec
| CREATE VIEW v_wf_sum_part AS SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_sum_part
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_sum_part
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_sum_part` (`a`, `SUM(b) OVER (PARTITION BY a ORDER BY b)`) AS SELECT `a` AS `a`,SUM(`b`) OVER (PARTITION BY `a` ORDER BY `b`) AS `SUM(b) OVER (PARTITION BY a ORDER BY b)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_wf_rows
< ok 0 0
> exec
| CREATE VIEW v_wf_rows AS SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_rows
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_rows
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_rows` (`a`, `name_exp_2`) AS SELECT `a` AS `a`,SUM(`b`) OVER (PARTITION BY `a` ORDER BY `b` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS `SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_wf_range
< ok 0 0
> exec
| CREATE VIEW v_wf_range AS SELECT a, AVG(b) OVER (ORDER BY a RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_range
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_range
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_range` (`a`, `name_exp_2`) AS SELECT `a` AS `a`,AVG(`b`) OVER (ORDER BY `a` RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS `AVG(b) OVER (ORDER BY a RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_wf_named
< ok 0 0
> exec
| CREATE VIEW v_wf_named AS SELECT a, RANK() OVER w, DENSE_RANK() OVER w FROM t WINDOW w AS (ORDER BY a)
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_named
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_named
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_named` (`a`, `RANK() OVER w`, `DENSE_RANK() OVER w`) AS SELECT `a` AS `a`,RANK() OVER `w` AS `RANK() OVER w`,DENSE_RANK() OVER `w` AS `DENSE_RANK() OVER w` FROM `test`.`t` WINDOW `w` AS (ORDER BY `a`)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_wf_multi
< ok 0 0
> exec
| CREATE VIEW v_wf_multi AS SELECT a, ROW_NUMBER() OVER (ORDER BY a), SUM(b) OVER (ORDER BY a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_multi
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_multi
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_multi` (`a`, `ROW_NUMBER() OVER (ORDER BY a)`, `SUM(b) OVER (ORDER BY a)`) AS SELECT `a` AS `a`,ROW_NUMBER() OVER (ORDER BY `a`) AS `ROW_NUMBER() OVER (ORDER BY a)`,SUM(`b`) OVER (ORDER BY `a`) AS `SUM(b) OVER (ORDER BY a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_wf_lagld
< ok 0 0
> exec
| CREATE VIEW v_wf_lagld AS SELECT a, LAG(a, 1) OVER (ORDER BY a), LEAD(a, 1) OVER (ORDER BY a) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_wf_lagld
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_wf_lagld
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_lagld` (`a`, `LAG(a, 1) OVER (ORDER BY a)`, `LEAD(a, 1) OVER (ORDER BY a)`) AS SELECT `a` AS `a`,LAG(`a`, 1) OVER (ORDER BY `a`) AS `LAG(a, 1) OVER (ORDER BY a)`,LEAD(`a`, 1) OVER (ORDER BY `a`) AS `LEAD(a, 1) OVER (ORDER BY a)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
