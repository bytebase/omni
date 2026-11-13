-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   `a` int DEFAULT NULL,
|   UNIQUE KEY `uk_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   UNIQUE KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   UNIQUE KEY `uk_ab` (`a`,`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   UNIQUE KEY `uk_a` (`a`),
|   UNIQUE KEY `uk_b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` int DEFAULT NULL,
|   UNIQUE KEY `a` (`a`),
|   UNIQUE KEY `a_2` (`a`,`b`),
|   UNIQUE KEY `a_3` (`a`,`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
