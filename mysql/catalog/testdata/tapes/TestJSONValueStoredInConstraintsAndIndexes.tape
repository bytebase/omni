-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestJSONValueStoredInConstraintsAndIndexes default#1
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a JSON, CONSTRAINT c CHECK (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY) > 0))
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
|   `a` json DEFAULT NULL,
|   CONSTRAINT `c` CHECK ((json_value(`a`, _utf8mb4'$.b' returning signed default 7 on empty) > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b INT DEFAULT (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY)))
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
|   `a` json DEFAULT NULL,
|   `b` int DEFAULT (json_value(`a`, _utf8mb4'$.b' returning signed default 7 on empty))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a JSON, KEY k ((JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY))))
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
|   `a` json DEFAULT NULL,
|   KEY `k` ((json_value(`a`, _utf8mb4'$.b' returning signed default 7 on empty)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS g
< ok 0 0
> exec
| CREATE TABLE g (a JSON, b INT AS (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY)) STORED)
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
|   `a` json DEFAULT NULL,
|   `b` int GENERATED ALWAYS AS (json_value(`a`, _utf8mb4'$.b' returning signed default 7 on empty)) STORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
