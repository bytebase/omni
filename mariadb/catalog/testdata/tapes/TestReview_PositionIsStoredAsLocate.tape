-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_PositionIsStoredAsLocate default#1
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
| CREATE VIEW pv1 AS SELECT POSITION('1' IN '1') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW pv1
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| pv1
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `pv1` AS select locate('1','1') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW re_pv1 AS select locate('1','1') AS `x` from `t`
< ok 0 0
> exec
| CREATE VIEW pv2 AS SELECT POSITION((1 = 1) IN '1') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW pv2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| pv2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `pv2` AS select locate(1 = 1,'1') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW re_pv2 AS select locate(1 = 1,'1') AS `x` from `t`
< ok 0 0
> exec
| CREATE VIEW pv3 AS SELECT POSITION((a + 1) IN '1') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW pv3
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| pv3
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `pv3` AS select locate(`t`.`a` + 1,'1') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW re_pv3 AS select locate(`t`.`a` + 1,'1') AS `x` from `t`
< ok 0 0
> exec
| CREATE VIEW pv4 AS SELECT LOCATE((1 = 1), '1') AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW pv4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| pv4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `pv4` AS select locate(1 = 1,'1') AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW re_pv4 AS select locate(1 = 1,'1') AS `x` from `t`
< ok 0 0
