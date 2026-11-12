-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C24 default#1
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
| SET SESSION sql_generate_invisible_primary_key = ON
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = ON
< error
| 1193 HY000 Unknown system variable 'show_gipk_in_create_table_and_information_schema'
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
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY ORDINAL_POSITION LIMIT 1
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE t2 (id INT NOT NULL UNIQUE, a INT)
< ok 0 0
> exec
| CREATE TABLE t3 (id INT AUTO_INCREMENT, a INT, PRIMARY KEY (id))
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 				 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 				 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 				 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t3' AND COLUMN_NAME='my_row_id'
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
| CREATE TABLE t (my_row_id INT, a INT)
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| my_row_id
