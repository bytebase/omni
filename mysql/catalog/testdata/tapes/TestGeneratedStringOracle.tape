-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestGeneratedStringOracle default#1
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS ('it''s'))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS (_utf8mb4'it\'s') VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS (_latin1 X'27'))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS (_latin1'\'') VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS (_latin1 X'5c'))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS (_latin1'\\') VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS (_latin1 X'5c6e'))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS (_latin1'\\n') VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS ('a\\b'))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS (_utf8mb4'a\\b') VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a VARCHAR(10) AS (_latin1'x'))
< ok 0 0
> query
| SHOW CREATE TABLE g
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| g
: bytes
| CREATE TABLE `g` (
|   `a` varchar(10) GENERATED ALWAYS AS (_latin1'x') VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
