-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_DefaultOutsideTheConvertibleBand default#1
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
| CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT (X'80'))
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 80
> exec
| DROP TABLE t
< ok 0 0
> exec
| CREATE TABLE `t` (
|   `a` varchar(20) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT X'80'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
< ok 0 0
> exec
| INSERT INTO t VALUES ()
< ok 1 0
> query
| SELECT HEX(a) FROM t
< result
< column "HEX(a)" "VARCHAR"
< row
: bytes
| 80
