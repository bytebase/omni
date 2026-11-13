-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   PRIMARY KEY (`id`),
|   KEY `idx_name` (`name`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `idx_email` (`email`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_ci3
< ok 0 0
> exec
| CREATE TABLE t_ci3 (id INT PRIMARY KEY, content TEXT)
< ok 0 0
> exec
| CREATE FULLTEXT INDEX idx_ft ON t_ci3 (content)
< ok 0 0
> query
| SHOW CREATE TABLE t_ci3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ci3
: bytes
| CREATE TABLE `t_ci3` (
|   `id` int NOT NULL,
|   `content` text,
|   PRIMARY KEY (`id`),
|   FULLTEXT KEY `idx_ft` (`content`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_ci4
< ok 0 0
> exec
| CREATE TABLE t_ci4 (id INT PRIMARY KEY, geo GEOMETRY NOT NULL)
< ok 0 0
> exec
| CREATE SPATIAL INDEX idx_sp ON t_ci4 (geo)
< ok 0 0
> query
| SHOW CREATE TABLE t_ci4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ci4
: bytes
| CREATE TABLE `t_ci4` (
|   `id` int NOT NULL,
|   `geo` geometry NOT NULL,
|   PRIMARY KEY (`id`),
|   SPATIAL KEY `idx_sp` (`geo`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< error
| 1064 42000 …near 'IF NOT EXISTS idx_val ON t_ci5 (val)' at line 1
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
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_di2
< ok 0 0
> exec
| CREATE TABLE t_di2 (id INT PRIMARY KEY)
< ok 0 0
> exec
| DROP INDEX idx_noexist ON t_di2
< error
| 1091 42000 Can't DROP 'idx_noexist'; check that column/key exists
