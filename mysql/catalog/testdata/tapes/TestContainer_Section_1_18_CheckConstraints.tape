-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   `a` int DEFAULT NULL,
|   CONSTRAINT `t_chk1_chk_1` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   CONSTRAINT `chk_name` CHECK ((`a` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_chk3
< ok 0 0
> exec
| CREATE TABLE t_chk3 (a INT, CHECK (a > 0) NOT ENFORCED)
< ok 0 0
> query
| SHOW CREATE TABLE t_chk3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chk3
: bytes
| CREATE TABLE `t_chk3` (
|   `a` int DEFAULT NULL,
|   CONSTRAINT `t_chk3_chk_1` CHECK ((`a` > 0)) /*!80016 NOT ENFORCED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   CONSTRAINT `t_chk4_chk_1` CHECK ((`a` > 0)),
|   CONSTRAINT `t_chk4_chk_2` CHECK ((`b` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   CONSTRAINT `t_chk5_chk_1` CHECK (((`a` > 0) and (`a` < 100)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   CONSTRAINT `t_chk6_chk_1` CHECK ((`a` > `b`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
