-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparse_Section_7_3_ExpressionViews default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_arithmetic_expr
< ok 0 0
> exec
| CREATE VIEW v_arithmetic_expr AS SELECT a + b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_arithmetic_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_arithmetic_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_arithmetic_expr` AS select (`t`.`a` + `t`.`b`) AS `a + b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_function_call
< ok 0 0
> exec
| CREATE VIEW v_function_call AS SELECT CONCAT(a, b) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_function_call
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_function_call
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_function_call` AS select concat(`t`.`a`,`t`.`b`) AS `CONCAT(a, b)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_case_expr
< ok 0 0
> exec
| CREATE VIEW v_case_expr AS SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_case_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_case_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_case_expr` AS select (case when (`t`.`a` > 0) then 'pos' else 'neg' end) AS `CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_cast_expr
< ok 0 0
> exec
| CREATE VIEW v_cast_expr AS SELECT CAST(a AS CHAR) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_cast_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_cast_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_cast_expr` AS select cast(`t`.`a` as char charset utf8mb4) AS `CAST(a AS CHAR)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_aggregate_expr
< ok 0 0
> exec
| CREATE VIEW v_aggregate_expr AS SELECT COUNT(*), SUM(a) FROM t GROUP BY a HAVING SUM(a) > 10
< ok 0 0
> query
| SHOW CREATE VIEW v_aggregate_expr
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_aggregate_expr
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_aggregate_expr` AS select count(0) AS `COUNT(*)`,sum(`t`.`a`) AS `SUM(a)` from `t` group by `t`.`a` having (sum(`t`.`a`) > 10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
