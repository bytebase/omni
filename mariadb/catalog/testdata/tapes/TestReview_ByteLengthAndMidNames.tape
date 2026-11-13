-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_ByteLengthAndMidNames default#1
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
| CREATE VIEW nv0 AS SELECT length('ab') AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'nv0'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select octet_length('ab') AS `x`
> exec
| CREATE VIEW nv1 AS SELECT mid('abc',1,1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'nv1'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select substr('abc',1,1) AS `x`
> exec
| CREATE VIEW nv2 AS SELECT octet_length('ab') AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'nv2'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select octet_length('ab') AS `x`
> exec
| CREATE VIEW nv3 AS SELECT char_length('ab') AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'nv3'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select char_length('ab') AS `x`
> exec
| CREATE VIEW nv4 AS SELECT substring('abc',1) AS x
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'nv4'
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select substr('abc',1) AS `x`
