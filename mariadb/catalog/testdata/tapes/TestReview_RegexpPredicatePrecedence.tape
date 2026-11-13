-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_RegexpPredicatePrecedence default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v_rx1 AS SELECT (a = 1) REGEXP '0' AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rx1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rx1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rx1` AS select (`t`.`a` = 1) regexp '0' AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v_rx2 AS SELECT a REGEXP 'x' = 1 AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW v_rx2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_rx2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_rx2` AS select `t`.`a` regexp 'x' = 1 AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
