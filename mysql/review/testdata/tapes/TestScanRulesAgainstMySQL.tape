-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScanRulesAgainstMySQL/CREATE_TABLE_AS_SELECT,_which_copies_no_key default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE c AS SELECT id, v FROM plain
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| c
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| c
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE d (PRIMARY KEY (id)) AS SELECT id FROM plain
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| d
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| c
: int64 0
< row
: bytes
| d
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| c
: int64 0
< row
: bytes
| d
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/CREATE_TABLE_IF_NOT_EXISTS,_not_judged default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE IF NOT EXISTS a (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/CREATE_TABLE_IF_NOT_EXISTS_of_a_table_with_a_key default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE IF NOT EXISTS plain (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/DROP_INDEX_of_the_primary_key default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| DROP INDEX `PRIMARY` ON plain
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 0
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 0
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_key_added_later default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE a ADD PRIMARY KEY (id)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| a
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_primary_key_dropped,_then_added_back default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 0
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 0
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE plain ADD PRIMARY KEY (v)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| v
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_primary_key_replaced default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE plain DROP PRIMARY KEY, ADD PRIMARY KEY (id, v)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id,v
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_table_created_without_a_key default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE b (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| b
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| b
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| a
: int64 0
< row
: bytes
| b
: int64 1
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_table_renamed_after_losing_its_key default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE plain DROP PRIMARY KEY
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 0
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 0
< row
: bytes
| t
: int64 1
> exec
| RENAME TABLE plain TO plain2
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain2
: int64 0
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain2
: int64 0
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_table_whose_database_the_change_drops default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE DATABASE d
< ok 1 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE d.e (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| DROP DATABASE d
< ok 1 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_unique_index_dropped_by_name default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE t DROP INDEX uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| DROP INDEX idx_n ON t
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/a_unique_key_the_change_makes,_then_drops default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE nokey ADD UNIQUE KEY uk_a (a)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| nokey
: bytes
| UNIQUE
: bytes
| uk_a
: bytes
| a
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| nokey
: bytes
| UNIQUE
: bytes
| uk_a
: bytes
| a
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE nokey DROP INDEX uk_a
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/an_index_of_a_table_CREATE_TABLE_IF_NOT_EXISTS_leaves default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| CREATE TABLE IF NOT EXISTS t (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| DROP INDEX uk_code ON t
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1

== TestScanRulesAgainstMySQL/an_index_of_a_table_the_change_recreates default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| DROP TABLE t
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
> exec
| CREATE TABLE t (code INT, KEY uk_code (code))
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
> exec
| ALTER TABLE t DROP INDEX uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0

== TestScanRulesAgainstMySQL/constraint_drops_the_SQL_names default#1
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
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| CHECK
: bytes
| ck_n
: bytes
| (`n` > 0)
< row
: bytes
| t
: bytes
| FOREIGN KEY
: bytes
| fk_p
: bytes
| p_id
< row
: bytes
| t
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 1
> exec
| ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CHECK ck_n
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| t
: bytes
| UNIQUE
: bytes
| uk_code
: bytes
| code
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
> exec
| ALTER TABLE t DROP CONSTRAINT uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop' ORDER BY c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "TEXT"
< row
: bytes
| p
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY t.TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
> query ^
< result
< column "TABLE_NAME" "VARCHAR"
< column "c.CONSTRAINT_NAME IS NOT NULL" "INT"
< row
: bytes
| nokey
: int64 0
< row
: bytes
| p
: int64 1
< row
: bytes
| plain
: int64 1
< row
: bytes
| t
: int64 0
