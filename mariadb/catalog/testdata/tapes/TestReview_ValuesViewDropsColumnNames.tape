-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ValuesViewDropsColumnNames default#1
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
| CREATE VIEW v_values(x, y) AS VALUES (1, 2)
< ok 0 0
> query
| SHOW CREATE VIEW v_values
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_values
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_values` AS values (1,2)
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'v_values' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1
< row
: bytes
| 2
> exec
| SELECT x FROM v_values
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec
| CREATE VIEW v_select(x, y) AS SELECT 1, 2
< ok 0 0
> query
| SHOW CREATE VIEW v_select
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| v_select
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `v_select` AS select 1 AS `x`,2 AS `y`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
