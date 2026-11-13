-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C19 default#1
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
< column "EXPRESSION" "TEXT"
< row
: null
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
|   PRIMARY KEY (`id`),
|   KEY `idx_lower` ((lower(`name`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< column "EXPRESSION" "TEXT"
< row
: bytes
| k_cast
: bytes
| cast(json_extract(`payload`,_utf8mb4\'$.age\') as unsigned)
< row
: bytes
| k_low
: bytes
| lower(`name`)
< row
: bytes
| k_sum
: bytes
| (`a` + `b`)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< column "EXPRESSION" "TEXT"
< row
: bytes
| lower(`name`)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t4a (a INT, INDEX ((a + RAND())))
< error
| 3758 HY000 Expression of functional index 'functional_index' contains a disallowed function.
> exec
| CREATE TABLE t4b (a INT, INDEX ((a)))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| CREATE TABLE t4c (payload JSON, INDEX ((payload->'$.name')))
< error
| 3753 42000 Cannot create a functional index on a function that returns a JSON or GEOMETRY value.
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
|   PRIMARY KEY (`id`),
|   KEY `idx_name` ((cast(json_unquote(json_extract(`doc`,_utf8mb4'$.name')) as char(64) charset utf8mb4)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE TABLE t_bad (doc JSON, INDEX ((doc->>'$.name')))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE INDEX idx_lower ON t ((LOWER(name)))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`
< error
| 3755 HY000 Cannot drop column '!hidden!idx_lower!0!0' because it is used by a functional index. In order to drop the column, you must remove the functional index.
> exec
| ALTER TABLE t RENAME INDEX idx_lower TO idx_lc
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
|   PRIMARY KEY (`id`),
|   KEY `idx_lc` ((lower(`name`)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
