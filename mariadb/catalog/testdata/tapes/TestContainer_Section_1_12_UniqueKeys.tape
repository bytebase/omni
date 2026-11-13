-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_12_UniqueKeys default#1
> exec
| DROP TABLE IF EXISTS t_uk1
< ok 0 0
> exec
| CREATE TABLE t_uk1 (a INT, UNIQUE KEY `uk_a` (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_uk1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uk1
: bytes
| CREATE TABLE `t_uk1` (
|   `a` int(11) DEFAULT NULL,
|   UNIQUE KEY `uk_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_uk2
< ok 0 0
> exec
| CREATE TABLE t_uk2 (a INT, UNIQUE KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_uk2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uk2
: bytes
| CREATE TABLE `t_uk2` (
|   `a` int(11) DEFAULT NULL,
|   UNIQUE KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_uk3
< ok 0 0
> exec
| CREATE TABLE t_uk3 (a INT, b INT, UNIQUE KEY `uk_ab` (a, b))
< ok 0 0
> query
| SHOW CREATE TABLE t_uk3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uk3
: bytes
| CREATE TABLE `t_uk3` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   UNIQUE KEY `uk_ab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_uk4
< ok 0 0
> exec
| CREATE TABLE t_uk4 (a INT, b INT, c INT, UNIQUE KEY `uk_a` (a), UNIQUE KEY `uk_b` (b))
< ok 0 0
> query
| SHOW CREATE TABLE t_uk4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uk4
: bytes
| CREATE TABLE `t_uk4` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_uk5
< ok 0 0
> exec
| CREATE TABLE t_uk5 (a INT, b INT, c INT, UNIQUE KEY (a), UNIQUE KEY (a, b), UNIQUE KEY (a, c))
< ok 0 0
> query
| SHOW CREATE TABLE t_uk5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_uk5
: bytes
| CREATE TABLE `t_uk5` (
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` int(11) DEFAULT NULL,
|   UNIQUE KEY `a` (`a`),
|   UNIQUE KEY `a_2` (`a`,`b`),
|   UNIQUE KEY `a_3` (`a`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
