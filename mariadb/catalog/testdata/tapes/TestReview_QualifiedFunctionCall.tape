-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_QualifiedFunctionCall default#1
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
| CREATE FUNCTION `upper`(x INT) RETURNS INT DETERMINISTIC RETURN x+1
< ok 0 0
> exec
| CREATE FUNCTION plainfn(x INT) RETURNS INT DETERMINISTIC RETURN x+2
< ok 0 0
> exec
| CREATE FUNCTION `odd name`(x INT) RETURNS INT DETERMINISTIC RETURN x+3
< ok 0 0
> exec
| CREATE FUNCTION MixedFn(x INT) RETURNS INT DETERMINISTIC RETURN x+4
< ok 0 0
> exec
| CREATE VIEW qv_upper AS SELECT testdb.`upper`(1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qv_upper'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `testdb`.`upper`(1) AS `x`
> exec
| CREATE VIEW qv_plain AS SELECT testdb.plainfn(1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qv_plain'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `testdb`.`plainfn`(1) AS `x`
> exec
| CREATE VIEW qv_odd AS SELECT testdb.`odd name`(1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qv_odd'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `testdb`.`odd name`(1) AS `x`
> exec
| CREATE VIEW qv_mixed AS SELECT testdb.MixedFn(1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qv_mixed'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `testdb`.`MixedFn`(1) AS `x`
> exec
| CREATE VIEW qv_bare AS SELECT `upper`(1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qv_bare'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select ucase(1) AS `x`
> exec
| CREATE VIEW qv_builtin AS SELECT UPPER('a') AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qv_builtin'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select ucase('a') AS `x`
