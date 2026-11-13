-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_rownum` AS select `t`.`a` AS `a`,row_number() OVER (ORDER BY `t`.`a` )  AS `ROW_NUMBER() OVER (ORDER BY a)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_sum_part` AS select `t`.`a` AS `a`,sum(`t`.`b`) OVER (PARTITION BY `t`.`a` ORDER BY `t`.`b` )  AS `SUM(b) OVER (PARTITION BY a ORDER BY b)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_rows` AS select `t`.`a` AS `a`,sum(`t`.`b`) OVER (PARTITION BY `t`.`a` ORDER BY `t`.`b` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)  AS `Name_exp_2` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_range` AS select `t`.`a` AS `a`,avg(`t`.`b`) OVER (ORDER BY `t`.`a` RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING)  AS `Name_exp_2` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_named` AS select `t`.`a` AS `a`,rank() OVER `w` AS `RANK() OVER w`,dense_rank() OVER `w` AS `DENSE_RANK() OVER w` from `t` window `w` AS (ORDER BY `t`.`a` ) 
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_multi` AS select `t`.`a` AS `a`,row_number() OVER (ORDER BY `t`.`a` )  AS `ROW_NUMBER() OVER (ORDER BY a)`,sum(`t`.`b`) OVER (ORDER BY `t`.`a` )  AS `SUM(b) OVER (ORDER BY a)` from `t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_wf_lagld` AS select `t`.`a` AS `a`,lag(`t`.`a`,1) OVER (ORDER BY `t`.`a` )  AS `LAG(a, 1) OVER (ORDER BY a)`,lead(`t`.`a`,1) OVER (ORDER BY `t`.`a` )  AS `LEAD(a, 1) OVER (ORDER BY a)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
