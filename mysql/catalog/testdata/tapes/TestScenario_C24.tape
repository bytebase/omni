-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C24 default#1
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
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = OFF
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
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
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = ON
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
|   `my_row_id` bigint unsigned NOT NULL AUTO_INCREMENT /*!80023 INVISIBLE */,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   PRIMARY KEY (`my_row_id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = OFF
< ok 0 0
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
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = ON
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, EXTRA, COLUMN_KEY
| 			 FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='my_row_id'
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< column "COLUMN_KEY" "ENUM"
< row
: bytes
| my_row_id
: bytes
| bigint unsigned
: bytes
| NO
: bytes
| auto_increment INVISIBLE
: bytes
| PRI
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
| 			 ORDER BY ORDINAL_POSITION LIMIT 1
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| my_row_id
> query
| SHOW INDEX FROM t
< result
< column "Table" "VARCHAR"
< column "Non_unique" "INT"
< column "Key_name" "VARCHAR"
< column "Seq_in_index" "UNSIGNED INT"
< column "Column_name" "VARCHAR"
< column "Collation" "VARCHAR"
< column "Cardinality" "BIGINT"
< column "Sub_part" "BIGINT"
< column "Packed" "NULL"
< column "Null" "VARCHAR"
< column "Index_type" "VARCHAR"
< column "Comment" "VARCHAR"
< column "Index_comment" "VARCHAR"
< column "Visible" "VARCHAR"
< column "Expression" "TEXT"
< row
: bytes
| t
: int64 0
: bytes
| PRIMARY
: int64 1
: bytes
| my_row_id
: bytes
| A
: int64 0
: null
: null
: bytes
|
: bytes
| BTREE
: bytes
|
: bytes
|
: bytes
| YES
: null
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = OFF
< ok 0 0
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
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = ON
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
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t3' AND COLUMN_NAME='my_row_id'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| CREATE TABLE t (my_row_id INT, a INT)
< error
| 4108 HY000 Failed to generate invisible primary key. Column 'my_row_id' already exists.
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
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
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = OFF
< ok 0 0
