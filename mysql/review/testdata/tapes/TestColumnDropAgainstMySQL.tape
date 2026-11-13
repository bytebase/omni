-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestColumnDropAgainstMySQL/a_RENAME_TABLE_of_a_temporary_table_onto_its_name_after,_which_MySQL_does_not_rename default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE tmp (a INT)
< ok 0 0
> exec
| RENAME TABLE tmp TO x
< error
| 1050 42S01 Table 'x' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_USE_of_a_database_the_server_lacks_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| USE nosuch_db
< error
| 1049 42000 Unknown database 'nosuch_db'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_USE_of_its_database_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| USE shop
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_check_of_another_column_added_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)
< ok 1 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x ADD CHECK (v > 0)
< ok 1 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 1 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_check_of_both_columns_of_the_key,_both_dropped default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), CONSTRAINT ck CHECK (a < b))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| x
: bytes
| ck
: bytes
| (`a` < `b`)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a, DROP COLUMN b
< error
| 3959 HY000 Check constraint 'ck' uses column 'a', hence column cannot be dropped or renamed.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| a,b

== TestColumnDropAgainstMySQL/a_check_of_it_alone default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY CHECK (id > 0), v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| x
: bytes
| x_chk_1
: bytes
| (`id` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_check_reads_it_and_another_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| x
: bytes
| ck
: bytes
| (`id` < `v`)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 3959 HY000 Check constraint 'ck' uses column 'id', hence column cannot be dropped or renamed.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_column_of_only_some_unique_keys_dropped_before_from_a_table_without_a_primary_key_partitioned_by_KEY_() default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT NOT NULL, w INT, UNIQUE KEY uu (u), UNIQUE KEY uv (u, v)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| kp
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| kp
: bytes
| u
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| kp
: bytes
| uu
: int64 0
: bytes
| u
: int64 -1
< row
: bytes
| kp
: bytes
| uv
: int64 0
: bytes
| u
: int64 -1
< row
: bytes
| kp
: bytes
| uv
: int64 0
: bytes
| v
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| kp
: bytes
| p0
: bytes
|
< row
: bytes
| kp
: bytes
| p1
: bytes
|
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE kp DROP COLUMN v
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_column_of_the_key_held_by_a_prefix_dropped_before_from_a_table_partitioned_by_KEY_(),_which_leaves_the_prefix_out default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (a VARCHAR(20), b INT, v INT, PRIMARY KEY (a(10), b), UNIQUE KEY ub (b)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| kp
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| kp
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| kp
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 10
< row
: bytes
| kp
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
< row
: bytes
| kp
: bytes
| ub
: int64 0
: bytes
| b
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| kp
: bytes
| p0
: bytes
|
< row
: bytes
| kp
: bytes
| p1
: bytes
|
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE kp DROP COLUMN a
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_column_out_of_any_key_dropped_before_from_a_table_without_a_primary_key_partitioned_by_KEY_() default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| kp
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| kp
: bytes
| u
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| kp
: bytes
| uu
: int64 0
: bytes
| u
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| kp
: bytes
| p0
: bytes
|
< row
: bytes
| kp
: bytes
| p1
: bytes
|
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE kp DROP COLUMN v
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_column_out_of_the_key_dropped_before_from_a_table_partitioned_by_KEY_() default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| kp
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| kp
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| kp
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| kp
: bytes
| p0
: bytes
|
< row
: bytes
| kp
: bytes
| p1
: bytes
|
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE kp DROP COLUMN v
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_column_the_table_lacks_inserted_into_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| INSERT INTO x (missing) VALUES (1)
< error
| 1054 42S22 Unknown column 'missing' in 'field list'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_database_the_server_lacks_dropped_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| DROP DATABASE missing_db
< error
| 1008 HY000 Can't drop database 'missing_db'; database doesn't exist
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_drop_after_of_it,_which_a_foreign_key_naming_a_table_of_its_name_in_another_database_does_not_refuse default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (id INT PRIMARY KEY); CREATE TABLE c (xid INT, FOREIGN KEY (xid) REFERENCES other.x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| c
: bytes
| xid
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| xid
: int64 1
: bytes
| xid
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| c_ibfk_1
: bytes
| x
: bytes
| xid
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| DROP TABLE x
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_drop_after_of_it,_which_a_foreign_key_of_the_database_refuses default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE c (xv INT, FOREIGN KEY (xv) REFERENCES x (v))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| c
: bytes
| xv
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| xv
: int64 1
: bytes
| xv
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| c_ibfk_1
: bytes
| x
: bytes
| xv
: bytes
| v
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| DROP TABLE x
< error
| 3730 HY000 Cannot drop table 'x' referenced by a foreign key constraint 'c_ibfk_1' on table 'c'.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_foreign_key_of_another_database_references_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.ch (x_id INT, CONSTRAINT fk_och FOREIGN KEY (x_id) REFERENCES shop.x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_och' of table 'ch'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_foreign_key_of_another_table_references_another_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE ch (x_v INT, CONSTRAINT fk_ch FOREIGN KEY (x_v) REFERENCES x (v))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ch
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| ch
: bytes
| x_v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| ch
: bytes
| fk_ch
: int64 1
: bytes
| x_v
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| ch
: bytes
| fk_ch
: bytes
| x
: bytes
| x_v
: bytes
| v
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_foreign_key_of_another_table_references_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ch
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| ch
: bytes
| x_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| ch
: bytes
| fk_ch
: int64 1
: bytes
| x_id
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| ch
: bytes
| fk_ch
: bytes
| x
: bytes
| x_id
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_ch' of table 'ch'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_foreign_key_of_another_table_references_it,_with_foreign_key_checks_off default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ch
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| ch
: bytes
| x_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| ch
: bytes
| fk_ch
: int64 1
: bytes
| x_id
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| ch
: bytes
| fk_ch
: bytes
| x
: bytes
| x_id
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_ch' of table 'ch'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_foreign_key_of_the_table_holds_another_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (v) REFERENCES p (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| fk_x
: int64 1
: bytes
| v
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
< row
: bytes
| x
: bytes
| fk_x
: bytes
| p
: bytes
| v
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_foreign_key_of_the_table_holds_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
< row
: bytes
| x
: bytes
| fk_x
: bytes
| p
: bytes
| id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1828 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_x'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_foreign_key_to_it_added_and_dropped_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (k INT PRIMARY KEY, x_id INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| ch
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| ch
: bytes
| k
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| ch
: bytes
| x_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| ch
: bytes
| PRIMARY
: int64 0
: bytes
| k
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE ch ADD CONSTRAINT f FOREIGN KEY (x_id) REFERENCES x (id)
< ok 0 0
> exec
| ALTER TABLE ch DROP FOREIGN KEY f
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_functional_index_reads_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY fx ((id + v)))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| fx
: int64 1
: bytes
| ((`id` + `v`))
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 3837 HY000 Column 'id' has a functional index dependency and cannot be dropped or renamed.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_generated_column_naming_it_only_in_a_literal default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g VARCHAR(20) AS (CONCAT('id', v)))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| g
: bytes
| concat(_utf8mb4\'id\',`v`)
: bytes
| VIRTUAL GENERATED
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_generated_column_reading_it,_dropped_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| g
: bytes
| (`id` + 1)
: bytes
| VIRTUAL GENERATED
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN g
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_generated_column_reading_it,_dropped_too default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| g
: bytes
| (`id` + 1)
: bytes
| STORED GENERATED
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN g
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_generated_column_reads_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| g
: bytes
| (`id` + 1)
: bytes
| VIRTUAL GENERATED
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 3108 HY000 Column 'id' has a generated column dependency.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_key_added_after_a_USE_of_another_database,_to_a_table_of_its_name_there default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (id INT, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| USE other
< ok 0 0
> exec
| ALTER TABLE x ADD PRIMARY KEY (v)
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_key_added_back_and_a_rename_onto_a_table_the_database_holds_after,_which_fail default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| ALTER TABLE x ADD PRIMARY KEY (v), RENAME TO y
< error
| 1050 42S01 Table 'y' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_key_is_added_back default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ADD PRIMARY KEY (v)
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| v

== TestColumnDropAgainstMySQL/a_key_of_a_column_prefix default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (s VARCHAR(50), v INT, PRIMARY KEY (s(10)))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| s
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| s
: int64 10
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN s
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_key_of_one_AUTO_INCREMENT_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT AUTO_INCREMENT PRIMARY KEY, v INT); INSERT INTO x (v) VALUES (1), (2)
< ok 2 1
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
| auto_increment
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP id
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_key_of_one_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_key_of_one_column,_with_a_unique_key_on_a_NOT_NULL_column_left default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, UNIQUE KEY uv (v))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_later_statement_adds_a_key_back default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| ALTER TABLE x ADD PRIMARY KEY (v)
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| v

== TestColumnDropAgainstMySQL/a_later_statement_that_adds_a_key_back,_but_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN missing, ADD PRIMARY KEY (v)
< error
| 1091 42000 Can't DROP 'missing'; check that column/key exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_literal_default_names_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(10) DEFAULT 'id')
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| 'id'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_read-only_transaction_started_before,_which_the_drop_commits default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| START TRANSACTION READ ONLY
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rename_after_onto_a_table_an_earlier_pair_moves_away default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| RENAME TABLE y TO y2, x TO y
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rename_after_with_a_pair_onto_a_name_an_earlier_pair_made,_which_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
< row
: bytes
| z
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| z
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| z
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| RENAME TABLE x TO spare, y TO spare
< error
| 1050 42S01 Table 'spare' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rename_after_with_a_pair_onto_a_name_an_earlier_pair_made_and_a_later_one_moved_away default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
< row
: bytes
| z
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| z
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| z
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| RENAME TABLE x TO spare, spare TO s2, y TO spare
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rename_after_with_a_pair_onto_a_name_an_earlier_pair_moved_a_table_back_onto,_which_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
< row
: bytes
| z
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| z
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| z
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| RENAME TABLE y TO spare, x TO y, z TO y
< error
| 1050 42S01 Table 'y' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rename_of_its_name_onto_a_table_the_database_holds_after,_which_MySQL_runs_on_it_though_a_temporary_table_of_its_name_exists,_and_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE x (a INT)
< ok 0 0
> exec
| RENAME TABLE x TO y
< error
| 1050 42S01 Table 'y' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rename_onto_a_table_the_database_holds_after,_which_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| RENAME TABLE x TO y
< error
| 1050 42S01 Table 'y' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_rollback_to_a_savepoint_the_transaction_lacks_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| BEGIN
< ok 0 0
> exec
| ROLLBACK TO SAVEPOINT nosuch
< error
| 1305 42000 SAVEPOINT nosuch does not exist
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_row_inserted_before_into_a_table_the_database_lacks default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| INSERT INTO missing VALUES (1)
< error
| 1146 42S02 Table 'shop.missing' doesn't exist
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_scope_for_a_system_variable_before_a_user_variable default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET SESSION foreign_key_checks = 1, @x = 1
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_statement_on_another_table_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE y ADD COLUMN b INT
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_statement_prepared_and_deallocated_after,_which_runs_nothing default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| PREPARE s FROM 'ALTER TABLE x ADD PRIMARY KEY (v)'
< ok 0 0
> exec
| DEALLOCATE PREPARE s
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_statement_prepared_before_of_a_table_the_database_lacks default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| PREPARE s FROM 'SELECT * FROM missing'
< error
| 1146 42S02 Table 'shop.missing' doesn't exist
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_table_CREATE_TABLE_IF_NOT_EXISTS_leaves,_whose_foreign_key_would_reference_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE c (k INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| c
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| c
: bytes
| k
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| c
: bytes
| PRIMARY
: int64 0
: bytes
| k
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| CREATE TABLE IF NOT EXISTS c (k INT PRIMARY KEY, FOREIGN KEY (k) REFERENCES x (id))
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_table_altered_to_be_renamed_onto_it_after,_which_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| ALTER TABLE y RENAME TO x
< error
| 1050 42S01 Table 'x' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_table_made_before_whose_column's_REFERENCES,_which_MySQL_ignores,_names_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| CREATE TABLE c (k INT REFERENCES x (id), PRIMARY KEY (k))
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_table_made_of_its_name_after,_which_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY)
< error
| 1050 42S01 Table 'x' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_table_of_another_database_renamed_to_its_name_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.t (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE other.t RENAME TO x
< error
| 1050 42S01 Table 'x' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_table_renamed_onto_it_after,_which_fails default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| RENAME TABLE y TO x
< error
| 1050 42S01 Table 'x' already exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_table_the_database_lacks_dropped_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| DROP TABLE missing
< error
| 1051 42S02 Unknown table 'shop.missing'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_temporary_table_of_its_name_dropped_after,_and_then_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE x (a INT)
< ok 0 0
> exec
| DROP TABLE x
< ok 0 0
> exec ^
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_temporary_table_of_its_name_dropped_twice default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| DROP TEMPORARY TABLE IF EXISTS x
< ok 0 0
> exec
| DROP TEMPORARY TABLE x
< error
| 1051 42S02 Unknown table 'shop.x'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_temporary_table_of_its_name_made_after,_given_a_key,_and_dropped default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE x (a INT)
< ok 0 0
> exec
| ALTER TABLE x ADD PRIMARY KEY (a)
< ok 0 0
> exec
| DROP TABLE x
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_temporary_table_renamed_onto_its_name_after,_and_dropped default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE tmp (a INT)
< ok 0 0
> exec
| ALTER TABLE tmp RENAME TO x
< ok 0 0
> exec
| DROP TABLE x
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_truncate_after default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| TRUNCATE TABLE x
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_truncate_of_a_table_a_foreign_key_references_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE tp (a INT PRIMARY KEY); CREATE TABLE tc (a INT, FOREIGN KEY (a) REFERENCES tp (a))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| tc
: bytes
| BASE TABLE
< row
: bytes
| tp
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| tc
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| tp
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| tc
: bytes
| a
: int64 1
: bytes
| a
: int64 -1
< row
: bytes
| tp
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
< row
: bytes
| tc
: bytes
| tc_ibfk_1
: bytes
| tp
: bytes
| a
: bytes
| a
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| TRUNCATE TABLE tp
< error
| 1701 42000 Cannot truncate a table referenced in a foreign key constraint (`shop`.`tc`, CONSTRAINT `tc_ibfk_1`)
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_unique_key_of_it_alone default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id)); INSERT INTO x VALUES (1, 1), (2, 1)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_unique_key_of_it_and_a_prefix_of_another_column,_beside_a_unique_key_of_the_other default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(20) NOT NULL, UNIQUE KEY u (v), UNIQUE KEY uv (id, v(2))); INSERT INTO x VALUES (1, 'xx1'), (2, 'xx2')
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
: int64 2
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1062 23000 Duplicate entry 'xx' for key 'x.uv'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_unique_key_of_it_and_another_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_unique_key_of_it_and_another_column,_beside_a_unique_key_of_the_other default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, w INT, UNIQUE KEY u (v), UNIQUE KEY uv (id, v)); INSERT INTO x VALUES (1, 1, 1), (2, 2, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_unique_key_of_it_and_another_column,_both_dropped default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN v
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_unique_key_of_it_and_another_column,_the_other_repeating default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1), (2, 1)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1062 23000 Duplicate entry '1' for key 'x.u'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_user_variable_named_sql_require_primary_key_set_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET @sql_require_primary_key = 1
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/a_user_variable_set_to_a_double_out_of_range_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET @x = 1e400
< error
| 1367 22007 Illegal double '1e400' value found during parsing
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_user_variable_set_to_a_hexadecimal_value_of_an_uppercase_prefix_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET @x = 0X41
< error
| 1054 42S22 Unknown column '0X41' in 'field list'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_user_variable_set_to_a_malformed_hexadecimal_value_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET @x = X'4'
< error
| 1064 42000 …near 'X'4'' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/a_view_or_a_trigger_reads_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE VIEW xv AS SELECT id FROM x; CREATE TRIGGER xt BEFORE INSERT ON x FOR EACH ROW SET NEW.v = NEW.id
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| xv
: bytes
| VIEW
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| xv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
< row
: bytes
| x
: bytes
| xt
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/an_expression_default_reads_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT DEFAULT (id + 1))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
| DEFAULT_GENERATED
: bytes
| ((`id` + 1))
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 3766 HY000 Column 'id' of table 'x' has a default value expression dependency and cannot be dropped or renamed.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/an_invisible_column_and_a_visible_one_left default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE, w INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
| INVISIBLE
: bytes
| NULL
< row
: bytes
| x
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/another_column_dropped_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN v
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/another_index_holds_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY iv (v, id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| iv
: int64 1
: bytes
| v
: int64 -1
< row
: bytes
| x
: bytes
| iv
: int64 1
: bytes
| id
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/another_table_locked_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| LOCK TABLES y WRITE
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1100 HY000 Table 'x' was not locked with LOCK TABLES
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/both_columns_of_a_key_of_two default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN b, DROP COLUMN a
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/both_columns_of_a_key_of_two,_in_two_statements default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN b
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/both_columns_of_a_key_of_two,_in_two_statements,_a_unique_key_of_a_prefix_of_the_other_keeping_its_values_apart default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b VARCHAR(20), v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b(2))); INSERT INTO x VALUES (1, 'aa1', 1), (1, 'bb2', 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
< row
: bytes
| x
: bytes
| ub
: int64 0
: bytes
| b
: int64 2
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN b
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/both_columns_of_a_key_of_two,_in_two_statements,_a_unique_key_of_the_other_keeping_its_values_apart default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
< row
: bytes
| x
: bytes
| ub
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN b
< ok 2 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/foreign_key_checks_set_with_:=_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET foreign_key_checks := OFF
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/it_is_added_back_under_its_name default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ADD COLUMN id INT
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/it_is_the_only_column default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1090 42000 You can't delete all columns with ALTER TABLE; use DROP TABLE instead
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/one_column_of_a_key_of_two default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| b

== TestColumnDropAgainstMySQL/one_column_of_a_key_of_two,_the_other_repeating default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a
< error
| 1062 23000 Duplicate entry '1' for key 'x.PRIMARY'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| a,b

== TestColumnDropAgainstMySQL/one_column_of_a_key_of_two_whose_other_is_a_prefix,_beside_a_unique_key_of_the_whole_other default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (a INT, b VARCHAR(20), v INT, PRIMARY KEY (a, b(2)), UNIQUE KEY ub (b)); INSERT INTO x VALUES (1, 'xx1', 1), (2, 'xx2', 2)
< ok 2 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| b
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
: int64 2
< row
: bytes
| x
: bytes
| ub
: int64 0
: bytes
| b
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN a
< error
| 1062 23000 Duplicate entry 'xx' for key 'x.PRIMARY'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| a,b

== TestColumnDropAgainstMySQL/only_an_invisible_column_left default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
| INVISIBLE
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/rows_of_the_table_updated_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)
< ok 1 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| UPDATE x SET v = id + 1 WHERE id IN (SELECT id FROM (SELECT id FROM x) d)
< ok 1 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 1 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/sql_generate_invisible_primary_key_on,_then_off,_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/sql_generate_invisible_primary_key_on_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1235 42000 This version of MySQL doesn't yet support 'existing primary key drop without adding a new primary key. In @@sql_generate_invisible_primary_key=ON mode table should have a primary key. Please add a new primary key to be able to drop existing primary key.'
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/sql_require_primary_key_off_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET SESSION sql_require_primary_key = OFF
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/sql_require_primary_key_on_after default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| SET SESSION sql_require_primary_key = ON
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/sql_require_primary_key_on_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET SESSION sql_require_primary_key = ON
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 3750 HY000 Unable to create or change a table without a primary key, when the system variable 'sql_require_primary_key' is set. Add a primary key to the table or unset this variable to avoid this message. Note that tables without a primary key can cause performance problems in row-based replication, so please consult your DBA before changing this setting.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/sql_require_primary_key_set_with_a_variable_the_server_reads_only_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SET SESSION sql_require_primary_key = OFF, version = 'x'
< error
| 1238 HY000 Variable 'version' is a read only variable
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/statements_on_another_table,_on_it_that_give_it_no_key,_or_on_none_after default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
< row
: bytes
| y
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| y
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| INSERT INTO y VALUES (1)
< ok 1 0
> exec
| ALTER TABLE y ADD COLUMN b INT
< ok 0 0
> exec
| CREATE INDEX iv ON x (v)
< ok 0 0
> exec
| ALTER TABLE x ADD COLUMN w INT
< ok 0 0
> exec
| INSERT INTO x (v) VALUES (1)
< ok 1 0
> exec
| CREATE VIEW xv AS SELECT v FROM x
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/statements_that_cannot_fail_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| BEGIN
< ok 0 0
> exec
| SAVEPOINT s
< ok 0 0
> exec
| ROLLBACK
< ok 0 0
> exec
| COMMIT
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| ANALYZE TABLE missing
< ok 0 0
> exec
| ANALYZE TABLE x UPDATE HISTOGRAM ON missing
< ok 0 0
> exec
| SET foreign_key_checks = OFF, @@session.unique_checks = 0, @x = 'a', @h = X'41', @b = b'101', @n = -1.5e3, @h2 = X'aB', @h3 = 0xaB, @b2 = B'01', foreign_key_checks = DEFAULT, unique_checks = DEFAULT
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/tables_flushed_after default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| FLUSH TABLES
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/tables_flushed_with_a_read_lock_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| FLUSH TABLES WITH READ LOCK
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1223 HY000 Can't execute the query because you have a conflicting read lock
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/tables_locked_after default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| LOCK TABLES x READ
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec ^
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/tables_locked_for_reading_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| LOCK TABLES x READ
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1099 HY000 Table 'x' was locked with a READ lock and can't be updated
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/the_database's_tables_shown_before default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| SHOW TABLES FROM shop
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/the_table_is_partitioned_by_KEY_(),_its_primary_key default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| x
: bytes
| p0
: bytes
|
< row
: bytes
| x
: bytes
| p1
: bytes
|
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/the_table_is_partitioned_by_it default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY HASH (id) PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| x
: bytes
| p0
: bytes
| `id`
< row
: bytes
| x
: bytes
| p1
: bytes
| `id`
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 3855 HY000 Column 'id' has a partitioning function dependency and cannot be dropped or renamed.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/the_unique_key's_column_dropped_before_from_a_table_without_a_primary_key_partitioned_by_KEY_() default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| kp
: bytes
| BASE TABLE
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| kp
: bytes
| u
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| kp
: bytes
| w
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| kp
: bytes
| uu
: int64 0
: bytes
| u
: int64 -1
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| kp
: bytes
| p0
: bytes
|
< row
: bytes
| kp
: bytes
| p1
: bytes
|
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE kp DROP COLUMN u
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/twice default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN id
< error
| 1091 42000 Can't DROP 'id'; check that column/key exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_ALGORITHM=INPLACE default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INPLACE
< error
| 1846 0A000 ALGORITHM=INPLACE is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try ALGORITHM=COPY.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_ALGORITHM_and_LOCK_left_at_DEFAULT default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, LOCK DEFAULT
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/with_ALGORITHM_set_to_DEFAULT,_then_INSTANT default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, ALGORITHM=INSTANT
< error
| 1846 0A000 ALGORITHM=INSTANT is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try ALGORITHM=COPY/INPLACE.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_ALGORITHM_set_to_INSTANT,_then_DEFAULT default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, ALGORITHM=DEFAULT
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: null

== TestColumnDropAgainstMySQL/with_ALGORITHM_set_to_INSTANT_before_LOCK_left_at_DEFAULT default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, LOCK=DEFAULT
< error
| 1846 0A000 ALGORITHM=INSTANT is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try ALGORITHM=COPY/INPLACE.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_LOCK=NONE default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, LOCK=NONE
< error
| 1846 0A000 LOCK=NONE is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try LOCK=SHARED.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_a_check_the_table_lacks_altered default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x ALTER CHECK missing ENFORCED, DROP COLUMN id
< error
| 3821 HY000 Check constraint 'missing' is not found in the table.
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_a_column_the_table_lacks default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN missing
< error
| 1091 42000 Can't DROP 'missing'; check that column/key exists
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id

== TestColumnDropAgainstMySQL/with_a_partition_dropped default#1
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| SET sql_require_primary_key = OFF
< ok 0 0
> exec
| SET sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 0 0
> exec
| DROP DATABASE IF EXISTS bbdataarchive
< ok 0 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLE x (id INT, v INT, PRIMARY KEY (id)) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "ENUM"
< row
: bytes
| nokey
: bytes
| BASE TABLE
< row
: bytes
| p
: bytes
| BASE TABLE
< row
: bytes
| plain
: bytes
| BASE TABLE
< row
: bytes
| pv
: bytes
| VIEW
< row
: bytes
| t
: bytes
| BASE TABLE
< row
: bytes
| x
: bytes
| BASE TABLE
> query
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA,
| 		COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "TEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(CASE WHEN EXTRA LIKE '%DEFAULT_GENERATED%' THEN CONCAT('(', COLUMN_DEFAULT, ')') ELSE QUOTE(COLUMN_DEFAULT) END, '')" "TEXT"
< row
: bytes
| nokey
: bytes
| a
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| p
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| plain
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| pv
: bytes
| note
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| code
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| p_id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| t
: bytes
| n
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| id
: bytes
|
: bytes
|
: bytes
| NULL
< row
: bytes
| x
: bytes
| v
: bytes
|
: bytes
|
: bytes
| NULL
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')')), COALESCE(SUB_PART, -1) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, CONCAT('(', EXPRESSION, ')'))" "TEXT"
< column "COALESCE(SUB_PART, -1)" "BIGINT"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
: int64 -1
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
: int64 -1
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
: int64 -1
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
: int64 -1
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY TABLE_NAME, CONSTRAINT_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
: bytes
| p_id
: bytes
| id
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_NAME, k.CHECK_CLAUSE FROM information_schema.TABLE_CONSTRAINTS c JOIN information_schema.CHECK_CONSTRAINTS k
| 		ON k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME WHERE c.CONSTRAINT_SCHEMA = 'shop' AND c.CONSTRAINT_TYPE = 'CHECK' ORDER BY c.TABLE_NAME, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "CHECK_CLAUSE" "TEXT"
< row
: bytes
| t
: bytes
| ck_n
: bytes
| (`n` > 0)
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "TEXT"
< row
: bytes
| x
: bytes
| p0
: bytes
| `id`
< row
: bytes
| x
: bytes
| p1
: bytes
| `id`
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> exec
| ALTER TABLE x DROP COLUMN id, DROP PARTITION p0
< error
| 1064 42000 …near 'PARTITION p0' at line 1
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS shop.x, shop.y, shop.tmp
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "LONGTEXT"
< row
: bytes
| id
