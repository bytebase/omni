-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C19 default#1
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
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64))
< ok 0 0
> exec
| CREATE INDEX idx_lower ON t ((LOWER(name)))
< ok 0 0
> query
|
| 			SELECT COLUMN_NAME, EXPRESSION
| 			  FROM information_schema.STATISTICS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='idx_lower'
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| NULL
: bytes
| lower(`name`)
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
|   `id` int NOT NULL,
|   `name` varchar(64) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_lower` ((lower(`name`)))
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
| CREATE TABLE t (
| 			a INT, b INT,
| 			name VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
| 			payload JSON,
| 			INDEX k_sum ((a + b)),
| 			INDEX k_low ((LOWER(name))),
| 			INDEX k_cast ((CAST(payload->'$.age' AS UNSIGNED)))
| 		)
< ok 0 0
> query
|
| 			SELECT INDEX_NAME, EXPRESSION
| 			  FROM information_schema.STATISTICS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| k_cast
: bytes
| cast(json_extract(`payload`, _utf8mb4'$.age') as unsigned)
< row
: bytes
| k_low
: bytes
| lower(`name`)
< row
: bytes
| k_sum
: bytes
| `a` + `b`
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
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64))
< ok 0 0
> exec
| CREATE INDEX idx_lower ON t ((LOWER(name)))
< ok 0 0
> query
|
| 			SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| name
> query
|
| 			SELECT EXPRESSION FROM information_schema.STATISTICS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='idx_lower'
< result
< column "EXPRESSION" "VARCHAR"
< row
: bytes
| lower(`name`)
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
| CREATE TABLE t4a (a INT, INDEX ((a + RAND())))
< error
| 3758 HY000 Expression of expression index 'expression_index' contains a disallowed function
> exec
| CREATE TABLE t4b (a INT, INDEX ((a)))
< error
| 3762 HY000 Expression index on a column is not supported. Consider using a regular index instead
> exec
| CREATE TABLE t4c (payload JSON, INDEX ((payload->'$.name')))
< error
| 3753 HY000 Cannot create an expression index on a function that returns a JSON or GEOMETRY value
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
| 			id INT PRIMARY KEY,
| 			doc JSON,
| 			INDEX idx_name ((CAST(doc->>'$.name' AS CHAR(64))))
| 		)
< ok 0 0
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
|   `id` int NOT NULL,
|   `doc` json DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_name` ((cast(json_unquote(json_extract(`doc`, _utf8mb4'$.name')) as char(64))))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| CREATE TABLE t_bad (doc JSON, INDEX ((doc->>'$.name')))
< error
| 3757 HY000 Cannot create an expression index on an expression that returns a BLOB or TEXT. Please consider using CAST
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
| CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64))
< ok 0 0
> exec
| CREATE INDEX idx_lower ON t ((LOWER(name)))
< ok 0 0
> exec
| DROP INDEX idx_lower ON t
< ok 0 0
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
|   `id` int NOT NULL,
|   `name` varchar(64) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| CREATE INDEX idx_lower ON t ((LOWER(name)))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`
< error
| 1091 42000 Can't DROP '!hidden!idx_lower!0!0'; check that column/key exists
