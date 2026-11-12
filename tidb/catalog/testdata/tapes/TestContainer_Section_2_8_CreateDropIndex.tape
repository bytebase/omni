-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_2_8_CreateDropIndex default#1
> exec
| DROP TABLE IF EXISTS t_ci1
< ok 0 0
> exec
| CREATE TABLE t_ci1 (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| CREATE INDEX idx_name ON t_ci1 (name)
< ok 0 0
> query
| SHOW CREATE TABLE t_ci1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ci1
: bytes
| CREATE TABLE `t_ci1` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_ci2
< ok 0 0
> exec
| CREATE TABLE t_ci2 (id INT PRIMARY KEY, email VARCHAR(255))
< ok 0 0
> exec
| CREATE UNIQUE INDEX idx_email ON t_ci2 (email)
< ok 0 0
> query
| SHOW CREATE TABLE t_ci2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ci2
: bytes
| CREATE TABLE `t_ci2` (
|   `id` int NOT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_ci3
< ok 0 0
> exec
| CREATE TABLE t_ci3 (id INT PRIMARY KEY, content TEXT)
< ok 0 0
> exec
| CREATE FULLTEXT INDEX idx_ft ON t_ci3 (content)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP TABLE IF EXISTS t_ci4
< ok 0 0
> exec
| CREATE TABLE t_ci4 (id INT PRIMARY KEY, geo GEOMETRY NOT NULL)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 52 near "GEOMETRY NOT NULL)" 
> exec
| DROP TABLE IF EXISTS t_ci5
< ok 0 0
> exec
| CREATE TABLE t_ci5 (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| CREATE INDEX idx_val ON t_ci5 (val)
< ok 0 0
> exec
| CREATE INDEX IF NOT EXISTS idx_val ON t_ci5 (val)
< ok 0 0
> query
| SHOW CREATE TABLE t_ci5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ci5
: bytes
| CREATE TABLE `t_ci5` (
|   `id` int NOT NULL,
|   `val` int DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_val` (`val`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_ci6
< ok 0 0
> exec
| CREATE TABLE t_ci6 (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| CREATE INDEX idx_a ON t_ci6 (a)
< ok 0 0
> exec
| CREATE INDEX idx_a ON t_ci6 (b)
< error
| 1061 42000 Duplicate key name 'idx_a'
> exec
| DROP TABLE IF EXISTS t_di1
< ok 0 0
> exec
| CREATE TABLE t_di1 (id INT PRIMARY KEY, name VARCHAR(100), KEY idx_name (name))
< ok 0 0
> exec
| DROP INDEX idx_name ON t_di1
< ok 0 0
> query
| SHOW CREATE TABLE t_di1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_di1
: bytes
| CREATE TABLE `t_di1` (
|   `id` int NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_di2
< ok 0 0
> exec
| CREATE TABLE t_di2 (id INT PRIMARY KEY)
< ok 0 0
> exec
| DROP INDEX idx_noexist ON t_di2
< error
| 1091 42000 index idx_noexist doesn't exist
