-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScanRulesAgainstMariaDB default#1
> exec
| SET foreign_key_checks = 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY fk_p, DROP CONSTRAINT ck_n
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| ALTER TABLE t ADD INDEX IF NOT EXISTS uk_code (code)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| ALTER TABLE t ADD (INDEX IF NOT EXISTS uk_code (code))
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| ALTER TABLE plain DROP CONSTRAINT PRIMARY KEY
< ok 2 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n, DROP FOREIGN KEY IF EXISTS fk_p
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| ALTER TABLE t DROP INDEX IF EXISTS uk_code
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| PRIMARY KEY
: bytes
| PRIMARY
: bytes
| id
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 7 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 7 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 6 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| CREATE OR REPLACE TABLE plain (id INT)
< ok 0 0
> query
| SELECT c.TABLE_NAME, c.CONSTRAINT_TYPE, c.CONSTRAINT_NAME,
| 		COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k
| 			WHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),
| 			(SELECT ck.CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS ck
| 			WHERE ck.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND ck.CONSTRAINT_NAME = c.CONSTRAINT_NAME), '')
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 6 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 5 0
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
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, COALESCE(GENERATION_EXPRESSION, ''), EXTRA, COALESCE(QUOTE(COLUMN_DEFAULT), '')
| 		FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "COALESCE(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "EXTRA" "VARCHAR"
< column "COALESCE(QUOTE(COLUMN_DEFAULT), '')" "MEDIUMTEXT"
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
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
| 'NULL'
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY INDEX_NAME, SEQ_IN_INDEX
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| plain
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| fk_p
: int64 1
: bytes
| p_id
< row
: bytes
| t
: bytes
| idx_n
: int64 1
: bytes
| n
< row
: bytes
| t
: bytes
| uk_code
: int64 0
: bytes
| code
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME, COLUMN_NAME, REFERENCED_COLUMN_NAME FROM information_schema.KEY_COLUMN_USAGE
| 		WHERE TABLE_SCHEMA = 'shop' AND POSITION_IN_UNIQUE_CONSTRAINT IS NOT NULL ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME, ORDINAL_POSITION
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
| SELECT TABLE_NAME, CONSTRAINT_NAME, CHECK_CLAUSE FROM information_schema.CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME, BINARY CONSTRAINT_NAME
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
| `n` > 0
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop' ORDER BY BINARY EVENT_OBJECT_TABLE, BINARY TRIGGER_NAME
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT TABLE_NAME, PARTITION_NAME, COALESCE(PARTITION_EXPRESSION, '') FROM information_schema.PARTITIONS
| 		WHERE TABLE_SCHEMA = 'shop' AND PARTITION_NAME IS NOT NULL ORDER BY BINARY TABLE_NAME, PARTITION_ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "PARTITION_NAME" "VARCHAR"
< column "COALESCE(PARTITION_EXPRESSION, '')" "LONGTEXT"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "UNSIGNED BIGINT"
< row
: uint64 0
> query
| SELECT t.TABLE_NAME, c.CONSTRAINT_NAME IS NOT NULL FROM information_schema.TABLES t
| 		LEFT JOIN information_schema.TABLE_CONSTRAINTS c ON c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME AND c.CONSTRAINT_TYPE = 'PRIMARY KEY'
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| `n` > 0
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
| 		FROM information_schema.TABLE_CONSTRAINTS c WHERE c.CONSTRAINT_SCHEMA = 'shop'
| 		ORDER BY BINARY c.TABLE_NAME, c.CONSTRAINT_TYPE, BINARY c.CONSTRAINT_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_TYPE" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COALESCE((SELECT GROUP_CONCAT(k.COLUMN_NAME ORDER BY k.ORDINAL_POSITION) FROM information_schema.KEY_COLUMN_USAGE k\n\t\t\tWHERE k.CONSTRAINT_SCHEMA = c.CONSTRAINT_SCHEMA AND k.TABLE_NAME = c.TABLE_NAME AND k.CONSTRAINT_NAME = c.CONSTRAINT_NAME),\n\t\t\t(SELECT c" "LONGTEXT"
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
| 		WHERE t.TABLE_SCHEMA = 'shop' AND t.TABLE_TYPE = 'BASE TABLE' ORDER BY BINARY t.TABLE_NAME
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
