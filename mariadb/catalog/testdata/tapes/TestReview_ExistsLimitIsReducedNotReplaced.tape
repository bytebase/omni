-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ExistsLimitIsReducedNotReplaced default#1
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| CREATE VIEW v_ex1 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t)
< ok 0 0
> query
| SHOW CREATE VIEW v_ex1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ex1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ex1` AS select `t`.`a` AS `a` from `t` where exists(select 1 from `t` limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v_ex2 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t LIMIT 5)
< ok 0 0
> query
| SHOW CREATE VIEW v_ex2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ex2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ex2` AS select `t`.`a` AS `a` from `t` where exists(select 1 from `t` limit 1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v_ex3 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t LIMIT 0)
< ok 0 0
> query
| SHOW CREATE VIEW v_ex3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ex3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ex3` AS select `t`.`a` AS `a` from `t` where exists(select 1 from `t` limit 0)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW v_ex4 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t LIMIT 1 OFFSET 1)
< ok 0 0
> query
| SHOW CREATE VIEW v_ex4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_ex4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_ex4` AS select `t`.`a` AS `a` from `t` where exists(select 1 from `t` limit 1,1)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
