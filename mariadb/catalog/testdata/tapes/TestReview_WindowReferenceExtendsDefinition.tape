-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_WindowReferenceExtendsDefinition default#1
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
| CREATE VIEW w_order_added AS SELECT ROW_NUMBER() OVER (w ORDER BY b) AS x FROM t WINDOW w AS (PARTITION BY a)
< ok 0 0
> query
| SHOW CREATE VIEW w_order_added
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| w_order_added
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `w_order_added` AS select row_number() over ( partition by `t`.`a` order by `t`.`b`) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> exec
| CREATE VIEW w_frame_added AS SELECT SUM(b) OVER (w ROWS UNBOUNDED PRECEDING) AS x FROM t WINDOW w AS (ORDER BY a)
< ok 0 0
> query
| SHOW CREATE VIEW w_frame_added
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| w_frame_added
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `w_frame_added` AS select sum(`t`.`b`) over ( order by `t`.`a` rows between  unbounded  preceding  and  current row ) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
