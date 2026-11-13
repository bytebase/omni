-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_UnionOrderByKeepsWrittenAlias default#1
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
| CREATE VIEW v_un AS SELECT a AS x FROM t UNION SELECT b FROM t ORDER BY x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'v_un'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `testdb`.`t`.`a` AS `x` from `testdb`.`t` union select `testdb`.`t`.`b` AS `b` from `testdb`.`t` order by `x`
