-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C14 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT c_pos CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t ALTER CHECK c_pos NOT ENFORCED
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| ALTER TABLE t ALTER CHECK c_pos ENFORCED
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
|     a INT,
|     g INT AS (a * 2) STORED,
|     CONSTRAINT c_g_nonneg CHECK (g >= 0)
| )
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb' AND CONSTRAINT_NAME='c_g_nonneg'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| INSERT INTO t (a) VALUES (5)
< ok 1 0
> exec
| INSERT INTO t (a) VALUES (-1)
< ok 1 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE other (id INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CHECK (a IN (SELECT id FROM other)))
< ok 0 0
> exec
| CREATE TABLE t3 (a INT, CHECK (a < NOW()))
< ok 0 0
> exec
| CREATE TABLE t4 (a INT, CHECK (a < @max))
< ok 0 0
> exec
| CREATE TABLE t5 (a INT, CHECK (a < RAND()))
< ok 0 0
