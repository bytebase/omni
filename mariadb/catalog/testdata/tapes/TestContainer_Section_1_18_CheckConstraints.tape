-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_18_CheckConstraints default#1
> exec
| DROP TABLE IF EXISTS t_chk1
< ok 0 0
> exec
| CREATE TABLE t_chk1 (a INT, CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE t_chk1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk1
: bytes
| CREATE TABLE `t_chk1` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chk2
< ok 0 0
> exec
| CREATE TABLE t_chk2 (a INT, CONSTRAINT `chk_name` CHECK (a > 0))
< ok 0 0
> query
| SHOW CREATE TABLE t_chk2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk2
: bytes
| CREATE TABLE `t_chk2` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `chk_name` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chk4
< ok 0 0
> exec
| CREATE TABLE t_chk4 (a INT, b INT, CHECK (a > 0), CHECK (b > 0))
< ok 0 0
> query
| SHOW CREATE TABLE t_chk4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk4
: bytes
| CREATE TABLE `t_chk4` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a` > 0),
|   CONSTRAINT `CONSTRAINT_2` CHECK (`b` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chk5
< ok 0 0
> exec
| CREATE TABLE t_chk5 (a INT, CHECK (a > 0 AND a < 100))
< ok 0 0
> query
| SHOW CREATE TABLE t_chk5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk5
: bytes
| CREATE TABLE `t_chk5` (
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a` > 0 and `a` < 100)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chk6
< ok 0 0
> exec
| CREATE TABLE t_chk6 (a INT, b INT, CHECK (a > b))
< ok 0 0
> query
| SHOW CREATE TABLE t_chk6
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk6
: bytes
| CREATE TABLE `t_chk6` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   CONSTRAINT `CONSTRAINT_1` CHECK (`a` > `b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
