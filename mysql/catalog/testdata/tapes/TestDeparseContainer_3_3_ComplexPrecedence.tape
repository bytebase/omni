-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestDeparseContainer_3_3_ComplexPrecedence default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_left_add
< ok 0 0
> exec
| CREATE VIEW v_left_add AS SELECT a + b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_left_add
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_left_add
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_left_add` AS select ((`t`.`a` + `t`.`b`) + `t`.`c`) AS `a + b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_mul_add
< ok 0 0
> exec
| CREATE VIEW v_mul_add AS SELECT a * b + c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_mul_add
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mul_add
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mul_add` AS select ((`t`.`a` * `t`.`b`) + `t`.`c`) AS `a * b + c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_add_mul
< ok 0 0
> exec
| CREATE VIEW v_add_mul AS SELECT a + b * c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_add_mul
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_add_mul
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_add_mul` AS select (`t`.`a` + (`t`.`b` * `t`.`c`)) AS `a + b * c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_paren_mul
< ok 0 0
> exec
| CREATE VIEW v_paren_mul AS SELECT (a + b) * c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_paren_mul
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_paren_mul
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_paren_mul` AS select ((`t`.`a` + `t`.`b`) * `t`.`c`) AS `(a + b) * c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_or_and
< ok 0 0
> exec
| CREATE VIEW v_or_and AS SELECT a OR b AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_or_and
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_or_and
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_or_and` AS select ((0 <> `t`.`a`) or ((0 <> `t`.`b`) and (0 <> `t`.`c`))) AS `a OR b AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_paren_or
< ok 0 0
> exec
| CREATE VIEW v_paren_or AS SELECT (a OR b) AND c FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_paren_or
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_paren_or
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_paren_or` AS select (((0 <> `t`.`a`) or (0 <> `t`.`b`)) and (0 <> `t`.`c`)) AS `(a OR b) AND c` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| DROP VIEW IF EXISTS v_mixed
< ok 0 0
> exec
| CREATE VIEW v_mixed AS SELECT a > 0 AND b < 10 OR c = 5 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_mixed
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_mixed
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mixed` AS select (((`t`.`a` > 0) and (`t`.`b` < 10)) or (`t`.`c` = 5)) AS `a > 0 AND b < 10 OR c = 5` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
