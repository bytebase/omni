-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_QuantifiedComparisonDropsParens default#1
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
| INSERT INTO t VALUES (1)
< ok 1 0
> exec
| CREATE VIEW q AS SELECT (1 OR 0) > ALL (SELECT 2) AS x FROM t
< ok 0 0
> query
| SHOW CREATE VIEW q
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| q
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `q` AS select 1 or 0 > all (select 2) AS `x` from `t`
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
> query
| SELECT (1 OR 0) > ALL (SELECT 2)
< result
< column "(1 OR 0) > ALL (SELECT 2)" "INT"
< row
: int64 0
> query
| SELECT x FROM q
< result
< column "x" "INT"
< row
: int64 1
