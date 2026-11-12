-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainer_Section_1_9_GeneratedColumns default#1
> exec
| DROP TABLE IF EXISTS t_gen1
< ok 0 0
> exec
| CREATE TABLE t_gen1 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL)
< ok 0 0
> query
| SHOW CREATE TABLE t_gen1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen1
: bytes
| CREATE TABLE `t_gen1` (
|   `col1` int DEFAULT NULL,
|   `col2` int DEFAULT NULL,
|   `col3` int GENERATED ALWAYS AS (`col1` + `col2`) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_gen2
< ok 0 0
> exec
| CREATE TABLE t_gen2 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 * col2) STORED)
< ok 0 0
> query
| SHOW CREATE TABLE t_gen2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen2
: bytes
| CREATE TABLE `t_gen2` (
|   `col1` int DEFAULT NULL,
|   `col2` int DEFAULT NULL,
|   `col3` int GENERATED ALWAYS AS (`col1` * `col2`) STORED
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_gen3
< ok 0 0
> exec
| CREATE TABLE t_gen3 (first_name VARCHAR(50), last_name VARCHAR(50), full_name VARCHAR(255) AS (CONCAT(first_name, ' ', last_name)) VIRTUAL)
< ok 0 0
> query
| SHOW CREATE TABLE t_gen3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen3
: bytes
| CREATE TABLE `t_gen3` (
|   `first_name` varchar(50) DEFAULT NULL,
|   `last_name` varchar(50) DEFAULT NULL,
|   `full_name` varchar(255) GENERATED ALWAYS AS (concat(`first_name`, _utf8mb4' ', `last_name`)) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_gen4
< ok 0 0
> exec
| CREATE TABLE t_gen4 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) STORED NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_gen4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen4
: bytes
| CREATE TABLE `t_gen4` (
|   `col1` int DEFAULT NULL,
|   `col2` int DEFAULT NULL,
|   `col3` int GENERATED ALWAYS AS (`col1` + `col2`) STORED NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_gen5
< ok 0 0
> exec
| CREATE TABLE t_gen5 (col1 INT, col2 INT, col3 INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL COMMENT 'sum of cols')
< ok 0 0
> query
| SHOW CREATE TABLE t_gen5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen5
: bytes
| CREATE TABLE `t_gen5` (
|   `col1` int DEFAULT NULL,
|   `col2` int DEFAULT NULL,
|   `col3` int GENERATED ALWAYS AS (`col1` + `col2`) VIRTUAL COMMENT 'sum of cols'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t_gen7
< ok 0 0
> exec
| CREATE TABLE t_gen7 (data JSON, name VARCHAR(255) GENERATED ALWAYS AS (JSON_EXTRACT(data, '$.name')) VIRTUAL)
< ok 0 0
> query
| SHOW CREATE TABLE t_gen7
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_gen7
: bytes
| CREATE TABLE `t_gen7` (
|   `data` json DEFAULT NULL,
|   `name` varchar(255) GENERATED ALWAYS AS (json_extract(`data`, _utf8mb4'$.name')) VIRTUAL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
