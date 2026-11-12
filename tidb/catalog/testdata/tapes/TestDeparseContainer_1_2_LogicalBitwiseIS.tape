-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_and_or` (`a AND b`, `a OR b`) AS SELECT `a` AND `b` AS `a AND b`,`a` OR `b` AS `a OR b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_xor` (`a XOR b`) AS SELECT `a` XOR `b` AS `a XOR b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_not` (`NOT a`) AS SELECT NOT `a` AS `NOT a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bitwise` (`a | b`, `a & b`, `a ^ b`) AS SELECT `a`|`b` AS `a | b`,`a`&`b` AS `a & b`,`a`^`b` AS `a ^ b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_shifts` (`a << b`, `a >> b`) AS SELECT `a`<<`b` AS `a << b`,`a`>>`b` AS `a >> b` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_bitnot` (`~a`) AS SELECT ~`a` AS `~a` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_isnull` (`a IS NULL`, `a IS NOT NULL`) AS SELECT `a` IS NULL AS `a IS NULL`,`a` IS NOT NULL AS `a IS NOT NULL` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_istf` (`a IS TRUE`, `a IS FALSE`) AS SELECT `a` IS TRUE AS `a IS TRUE`,`a` IS FALSE AS `a IS FALSE` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_in` (`a IN (1,2,3)`, `a NOT IN (1,2,3)`) AS SELECT `a` IN (1,2,3) AS `a IN (1,2,3)`,`a` NOT IN (1,2,3) AS `a NOT IN (1,2,3)` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_between` (`a BETWEEN 1 AND 10`, `a NOT BETWEEN 1 AND 10`) AS SELECT `a` BETWEEN 1 AND 10 AS `a BETWEEN 1 AND 10`,`a` NOT BETWEEN 1 AND 10 AS `a NOT BETWEEN 1 AND 10` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_like` (`a LIKE 'foo%'`, `a NOT LIKE 'bar%'`) AS SELECT `a` LIKE _UTF8MB4'foo%' AS `a LIKE 'foo%'`,`a` NOT LIKE _UTF8MB4'bar%' AS `a NOT LIKE 'bar%'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
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
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_like_esc` (`a LIKE 'x' ESCAPE '\\'`) AS SELECT `a` LIKE _UTF8MB4'x' AS `a LIKE 'x' ESCAPE '\\'` FROM `test`.`t`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
