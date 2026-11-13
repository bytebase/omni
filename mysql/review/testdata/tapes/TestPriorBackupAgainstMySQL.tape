-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestPriorBackupAgainstMySQL/a_column_USING_names default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN p USING (note) WHERE t.id = 1
< error
| 1054 42S22 Unknown column 'note' in 'from clause'

== TestPriorBackupAgainstMySQL/a_column_USING_names,_a_derived_table_on_its_side_may_hold default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN (SELECT 1 AS note) d ON d.note = t.id JOIN p USING (note) WHERE t.id = 3
< ok 0 0

== TestPriorBackupAgainstMySQL/a_column_in_a_subquery,_not_judged default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`plain`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `plain`.* FROM plain WHERE id IN (SELECT gone FROM t)
< error
| 1054 42S22 Unknown column 'gone' in 'field list'

== TestPriorBackupAgainstMySQL/a_column_of_a_joined_table default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN p ON p.gone = t.n WHERE t.id = 3
< error
| 1054 42S22 Unknown column 'p.gone' in 'on clause'

== TestPriorBackupAgainstMySQL/a_column_of_a_table_the_change_recreates default#1
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
| CREATE TABLE `bbdataarchive`.`bk_2_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_2_0` SELECT `nokey`.* FROM nokey WHERE b = 1
< error
| 1054 42S22 Unknown column 'b' in 'where clause'

== TestPriorBackupAgainstMySQL/a_column_the_SET_clause_names,_not_copied default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE id = 1
< ok 1 0

== TestPriorBackupAgainstMySQL/a_column_the_change_adds,_in_ORDER_BY default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE a = 1 ORDER BY gone LIMIT 1
< error
| 1054 42S22 Unknown column 'gone' in 'order clause'

== TestPriorBackupAgainstMySQL/a_column_the_change_adds,_qualified default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE nokey.gone = 1
< error
| 1054 42S22 Unknown column 'nokey.gone' in 'where clause'

== TestPriorBackupAgainstMySQL/a_column_the_change_adds,_unqualified default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE gone = 1
< error
| 1054 42S22 Unknown column 'gone' in 'where clause'

== TestPriorBackupAgainstMySQL/a_column_the_change_changes_the_name_of default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`p`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `p`.* FROM p WHERE memo = 'zz'
< error
| 1054 42S22 Unknown column 'memo' in 'where clause'

== TestPriorBackupAgainstMySQL/a_column_the_change_renames_to,_through_an_alias default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`p`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `x`.* FROM p AS x WHERE x.memo = 'a'
< error
| 1054 42S22 Unknown column 'x.memo' in 'where clause'

== TestPriorBackupAgainstMySQL/a_plain_backup default#1
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
| CREATE TABLE `bbdataarchive`.`bk_0_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_0_0` SELECT `t`.* FROM t WHERE id = 1
< ok 1 0
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE a = 1
< ok 1 0

== TestPriorBackupAgainstMySQL/a_table_the_change_creates default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`n2`
< error
| 1146 42S02 Table 'shop.n2' doesn't exist

== TestPriorBackupAgainstMySQL/a_table_the_change_creates,_in_a_subquery default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE id IN (SELECT id FROM n2)
< error
| 1146 42S02 Table 'shop.n2' doesn't exist

== TestPriorBackupAgainstMySQL/a_table_the_change_creates,_joined default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1
< error
| 1146 42S02 Table 'shop.n2' doesn't exist

== TestPriorBackupAgainstMySQL/a_table_the_change_renames_to default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nk`
< error
| 1146 42S02 Table 'shop.nk' doesn't exist

== TestPriorBackupAgainstMySQL/a_view_as_the_target default#1
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
| CREATE TABLE `bbdataarchive`.`bk_0_0` LIKE `shop`.`pv`
< error
| 1347 HY000 'shop.pv' is not BASE TABLE

== TestPriorBackupAgainstMySQL/statements_before_a_USE default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE gone = 1
< error
| 1054 42S22 Unknown column 'gone' in 'where clause'
> exec
| CREATE TABLE `bbdataarchive`.`bk_3_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_3_0` SELECT `t`.* FROM shop.t WHERE id = 1
< ok 1 0
