-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDeparseContainer_1_2_LogicalBitwiseIS default#1
> exec
| CREATE TABLE IF NOT EXISTS t (a INT, b INT, c INT)
< ok 0 0
> exec
| DROP VIEW IF EXISTS v_and_or
< ok 0 0
> exec
| CREATE VIEW v_and_or AS SELECT a AND b, a OR b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_and_or
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_and_or
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_and_or` AS select `t`.`a` <> 0 and `t`.`b` <> 0 AS `a AND b`,`t`.`a` <> 0 or `t`.`b` <> 0 AS `a OR b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_xor
< ok 0 0
> exec
| CREATE VIEW v_xor AS SELECT a XOR b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_xor
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_xor
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_xor` AS select `t`.`a` xor `t`.`b` AS `a XOR b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_not
< ok 0 0
> exec
| CREATE VIEW v_not AS SELECT NOT a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_not
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_not
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not` AS select `t`.`a` = 0 AS `NOT a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_bitwise
< ok 0 0
> exec
| CREATE VIEW v_bitwise AS SELECT a | b, a & b, a ^ b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bitwise
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bitwise
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bitwise` AS select `t`.`a` | `t`.`b` AS `a | b`,`t`.`a` & `t`.`b` AS `a & b`,`t`.`a` ^ `t`.`b` AS `a ^ b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_shifts
< ok 0 0
> exec
| CREATE VIEW v_shifts AS SELECT a << b, a >> b FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_shifts
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_shifts
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_shifts` AS select `t`.`a` << `t`.`b` AS `a << b`,`t`.`a` >> `t`.`b` AS `a >> b` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_bitnot
< ok 0 0
> exec
| CREATE VIEW v_bitnot AS SELECT ~a FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_bitnot
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_bitnot
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bitnot` AS select ~`t`.`a` AS `~a` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_isnull
< ok 0 0
> exec
| CREATE VIEW v_isnull AS SELECT a IS NULL, a IS NOT NULL FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_isnull
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_isnull
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_isnull` AS select `t`.`a` is null AS `a IS NULL`,`t`.`a` is not null AS `a IS NOT NULL` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_istf
< ok 0 0
> exec
| CREATE VIEW v_istf AS SELECT a IS TRUE, a IS FALSE FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_istf
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_istf
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_istf` AS select `t`.`a` is true AS `a IS TRUE`,`t`.`a` is false AS `a IS FALSE` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_in
< ok 0 0
> exec
| CREATE VIEW v_in AS SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_in
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_in
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_in` AS select `t`.`a` in (1,2,3) AS `a IN (1,2,3)`,`t`.`a` not in (1,2,3) AS `a NOT IN (1,2,3)` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_between
< ok 0 0
> exec
| CREATE VIEW v_between AS SELECT a BETWEEN 1 AND 10, a NOT BETWEEN 1 AND 10 FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_between
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_between
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_between` AS select `t`.`a` between 1 and 10 AS `a BETWEEN 1 AND 10`,`t`.`a` not between 1 and 10 AS `a NOT BETWEEN 1 AND 10` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_like
< ok 0 0
> exec
| CREATE VIEW v_like AS SELECT a LIKE 'foo%', a NOT LIKE 'bar%' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_like
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_like
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_like` AS select `t`.`a` like 'foo%' AS `a LIKE 'foo%'`,`t`.`a`  not like 'bar%' AS `a NOT LIKE 'bar%'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| DROP VIEW IF EXISTS v_like_esc
< ok 0 0
> exec
| CREATE VIEW v_like_esc AS SELECT a LIKE 'x' ESCAPE '\\' FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_like_esc
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_like_esc
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_like_esc` AS select `t`.`a` like 'x' escape '\\' AS `a LIKE 'x' ESCAPE '\\'` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
