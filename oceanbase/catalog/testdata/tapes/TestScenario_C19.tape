-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "EXPRESSION" "VARBINARY"
< row
: bytes
| SYS_NC18$
: bytes
| LOWER(`name`)
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(64) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_lower` ((LOWER(`name`))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
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
< column "EXPRESSION" "VARBINARY"
< row
: bytes
| k_cast
: bytes
| cast(JSON_EXTRACT(`payload`,'$.age') as unsigned)
< row
: bytes
| k_low
: bytes
| LOWER(`name`)
< row
: bytes
| k_sum
: bytes
| (`a` + `b`)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 6 0
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
< column "EXPRESSION" "VARBINARY"
< row
: bytes
| LOWER(`name`)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t4a (a INT, INDEX ((a + RAND())))
< error
| 3758 HY000 Expression of functional index contains a disallowed function.
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `doc` json DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_name` ((cast(JSON_UNQUOTE(JSON_EXTRACT(`doc`,'$.name')) as char(64) charset utf8mb4))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE TABLE t_bad (doc JSON, INDEX ((doc->>'$.name')))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 4 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(64) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| CREATE INDEX idx_lower ON t ((LOWER(name)))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`
< error
| 1091 42000 Can't DROP '!hidden!idx_lower!0!0'; check that column/key exists
> exec
| ALTER TABLE t RENAME INDEX idx_lower TO idx_lc
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(64) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_lc` ((LOWER(`name`))) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
