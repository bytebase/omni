-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_CTENameIsWrittenUnquoted default#1
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
| CREATE VIEW v_cte AS WITH `select` AS (SELECT a FROM t) SELECT a FROM `select`
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'v_cte'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| with select as (select `testdb`.`t`.`a` AS `a` from `testdb`.`t`)select `select`.`a` AS `a` from `select`
> exec
| SELECT 1 FROM v_cte
< error
| 1064 42000 …near 'select as (select `testdb`.`t`.`a` AS `a` from `testdb`.`t`)select `select`.`...' at line 1
