-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_left_add` (`a + b + c`) AS SELECT `a`+`b`+`c` AS `a + b + c` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mul_add` (`a * b + c`) AS SELECT `a`*`b`+`c` AS `a * b + c` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_add_mul` (`a + b * c`) AS SELECT `a`+`b`*`c` AS `a + b * c` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_paren_mul` (`(a + b) * c`) AS SELECT (`a`+`b`)*`c` AS `(a + b) * c` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_or_and` (`a OR b AND c`) AS SELECT `a` OR `b` AND `c` AS `a OR b AND c` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_paren_or` (`(a OR b) AND c`) AS SELECT (`a` OR `b`) AND `c` AS `(a OR b) AND c` FROM `test`.`t`
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_mixed` (`a > 0 AND b < 10 OR c = 5`) AS SELECT `a`>0 AND `b`<10 OR `c`=5 AS `a > 0 AND b < 10 OR c = 5` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
