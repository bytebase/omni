-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_1_15_ExpressionIndexes default#1
> exec
| DROP TABLE IF EXISTS t_expr1
< ok 0 0
> exec
| CREATE TABLE t_expr1 (name VARCHAR(100), KEY `idx` ((UPPER(name))))
< ok 0 0
> query
| SHOW CREATE TABLE t_expr1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_expr1
: bytes
| CREATE TABLE `t_expr1` (
|   `name` varchar(100) DEFAULT NULL,
|   KEY `idx` ((upper(`name`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_expr2
< ok 0 0
> exec
| CREATE TABLE t_expr2 (col1 INT, col2 INT, KEY `idx` ((col1 + col2)))
< ok 0 0
> query
| SHOW CREATE TABLE t_expr2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_expr2
: bytes
| CREATE TABLE `t_expr2` (
|   `col1` int DEFAULT NULL,
|   `col2` int DEFAULT NULL,
|   KEY `idx` (((`col1` + `col2`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_expr3
< ok 0 0
> exec
| CREATE TABLE t_expr3 (name VARCHAR(100), UNIQUE KEY `uidx` ((UPPER(name))))
< ok 0 0
> query
| SHOW CREATE TABLE t_expr3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_expr3
: bytes
| CREATE TABLE `t_expr3` (
|   `name` varchar(100) DEFAULT NULL,
|   UNIQUE KEY `uidx` ((upper(`name`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_expr4
< ok 0 0
> exec
| CREATE TABLE t_expr4 (a INT, b INT, KEY `idx_col` (a), KEY `idx_expr` ((a * b)))
< ok 0 0
> query
| SHOW CREATE TABLE t_expr4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_expr4
: bytes
| CREATE TABLE `t_expr4` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `idx_col` (`a`),
|   KEY `idx_expr` (((`a` * `b`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
