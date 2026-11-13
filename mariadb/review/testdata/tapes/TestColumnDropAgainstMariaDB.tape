-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestColumnDropAgainstMariaDB default#1
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1), (2, 2)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT AUTO_INCREMENT PRIMARY KEY, v INT); INSERT INTO x (v) VALUES (1), (2)
< ok 2 1
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (s VARCHAR(50), v INT, PRIMARY KEY (s(10)))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| s
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, UNIQUE KEY uv (v))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
> exec
| ALTER TABLE x DROP COLUMN a
< error
| 1072 42000 Key column 'a' doesn't exist in table
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| a,b
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
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
> exec
| ALTER TABLE x DROP COLUMN a
< error
| 1072 42000 Key column 'a' doesn't exist in table
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| a,b
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
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT)
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
< row
: bytes
| x
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))
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
< row
: bytes
| x
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
| 'NULL'
< row
: bytes
| x
: bytes
| g
: bytes
| `id` + 1
: bytes
| VIRTUAL GENERATED
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED)
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
< row
: bytes
| x
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
| 'NULL'
< row
: bytes
| x
: bytes
| g
: bytes
| `id` + 1
: bytes
| STORED GENERATED
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY CHECK (id > 0), v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< row
: bytes
| x
: bytes
| id
: bytes
| `id` > 0
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), CONSTRAINT ck CHECK (a < b))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
< row
: bytes
| x
: bytes
| ck
: bytes
| `a` < `b`
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
> exec
| ALTER TABLE x DROP COLUMN a, DROP COLUMN b
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN IF EXISTS id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN IF EXISTS id, DROP COLUMN IF EXISTS id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN IF EXISTS missing
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY iv (v, id))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| iv
: int64 1
: bytes
| v
< row
: bytes
| x
: bytes
| iv
: int64 1
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE VIEW xv AS SELECT id FROM x; CREATE TRIGGER xt BEFORE INSERT ON x FOR EACH ROW SET NEW.v = NEW.id
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< row
: bytes
| x
: bytes
| xt
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE ch (x_v INT, CONSTRAINT fk_ch FOREIGN KEY (x_v) REFERENCES x (v))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| ch
: bytes
| x_v
: bytes
|
: bytes
|
: bytes
| 'NULL'
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
| ch
: bytes
| fk_ch
: int64 1
: bytes
| x_v
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (v) REFERENCES p (id))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| fk_x
: int64 1
: bytes
| v
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 8 0
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
< row
: bytes
| x
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
| 'NULL'
< row
: bytes
| x
: bytes
| g
: bytes
| concat('id',`v`)
: bytes
| VIRTUAL GENERATED
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| ch
: bytes
| x_id
: bytes
|
: bytes
|
: bytes
| 'NULL'
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
| ch
: bytes
| fk_ch
: int64 1
: bytes
| x_id
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_ch' of table `shop`.`ch`
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| ch
: bytes
| x_id
: bytes
|
: bytes
|
: bytes
| 'NULL'
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
| ch
: bytes
| fk_ch
: int64 1
: bytes
| x_id
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_ch' of table `shop`.`ch`
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.ch (x_id INT, CONSTRAINT fk_och FOREIGN KEY (x_id) REFERENCES shop.x (id))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'fk_och' of table `other`.`ch`
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 1 0
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1))
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
< row
: bytes
| x
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
| 'NULL'
< row
: bytes
| x
: bytes
| g
: bytes
| `id` + 1
: bytes
| VIRTUAL GENERATED
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1054 42S22 Unknown column 'id' in 'GENERATED ALWAYS AS'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v))
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< row
: bytes
| x
: bytes
| ck
: bytes
| `id` < `v`
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1054 42S22 Unknown column 'id' in 'CHECK'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY HASH (id) PARTITIONS 2
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
| 1054 42S22 Unknown column 'id' in 'PARTITION BY'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN missing
< error
| 1091 42000 Can't DROP COLUMN `missing`; check that it exists
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN id
< error
| 1091 42000 Can't DROP COLUMN `id`; check that it exists
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| v
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| v
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN missing, ADD PRIMARY KEY (v)
< error
| 1091 42000 Can't DROP COLUMN `missing`; check that it exists
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)
< ok 1 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE ch (k INT PRIMARY KEY, x_id INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| 'NULL'
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
| ch
: bytes
| PRIMARY
: int64 0
: bytes
| k
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE c (k INT PRIMARY KEY)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| c
: bytes
| PRIMARY
: int64 0
: bytes
| k
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| CREATE TABLE c (k INT REFERENCES x (id), PRIMARY KEY (k))
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1829 HY000 Cannot drop column 'id': needed in a foreign key constraint 'c_ibfk_1' of table `shop`.`c`
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT, v INT, PRIMARY KEY (id)) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1113 42000 A table must have at least 1 column
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE, w INT)
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
< row
: bytes
| x
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT DEFAULT (id + 1))
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
< row
: bytes
| x
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
| '(`id` + 1)'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1054 42S22 Unknown column 'id' in 'DEFAULT'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(10) DEFAULT 'id')
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
< row
: bytes
| x
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
| '\'id\''
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.t (a INT PRIMARY KEY)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 1 0
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE tp (a INT PRIMARY KEY); CREATE TABLE tc (a INT, FOREIGN KEY (a) REFERENCES tp (a))
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
| 'NULL'
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
< row
: bytes
| tc
: bytes
| a
: int64 1
: bytes
| a
< row
: bytes
| tp
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| TRUNCATE TABLE tp
< error
| 1701 42000 Cannot truncate a table referenced in a foreign key constraint (`shop`.`tc`, CONSTRAINT `tc_ibfk_1` FOREIGN KEY (`a`) REFERENCES `shop`.`tp` (`a`))
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 8 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); INSERT INTO x VALUES (1, 1)
< ok 1 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| INSERT INTO x (missing) VALUES (1)
< error
| 1054 42S22 Unknown column 'missing' in 'INSERT INTO'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| SET foreign_key_checks = 0, version = 'x'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ANALYZE TABLE x UPDATE HISTOGRAM ON v
< error
| 1064 42000 …near 'UPDATE HISTOGRAM ON v' at line 1
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1), (2, 2)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1072 42000 Key column 'id' doesn't exist in table
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, UNIQUE KEY u (id, v)); INSERT INTO x VALUES (1, 1, 1), (2, 1, 2)
< ok 2 0
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
< row
: bytes
| x
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id)); INSERT INTO x VALUES (1, 1), (2, 1)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| 'NULL'
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
| kp
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| 'NULL'
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
| 'NULL'
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
| kp
: bytes
| uu
: int64 0
: bytes
| u
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| 'NULL'
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
| 'NULL'
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
| kp
: bytes
| uu
: int64 0
: bytes
| u
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b)); INSERT INTO x VALUES (1, 1, 1), (1, 2, 2)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| b
< row
: bytes
| x
: bytes
| ub
: int64 0
: bytes
| b
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
> exec
| ALTER TABLE x DROP COLUMN a
< error
| 1072 42000 Key column 'a' doesn't exist in table
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| a,b
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, w INT, UNIQUE KEY u (v), UNIQUE KEY uv (id, v)); INSERT INTO x VALUES (1, 1, 1), (2, 2, 2)
< ok 2 0
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| u
: int64 0
: bytes
| v
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
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
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 1072 42000 Key column 'id' doesn't exist in table
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| SET @x = 0X41
< error
| 1054 42S22 Unknown column '0X41' in 'SET'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE OR REPLACE TABLE x (id INT PRIMARY KEY)
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE OR REPLACE DATABASE shop
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE kp (u INT NOT NULL, v INT NOT NULL, w INT, UNIQUE KEY uu (u), UNIQUE KEY uv (u, v)) PARTITION BY KEY () PARTITIONS 2
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| 'NULL'
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
| kp
: bytes
| uu
: int64 0
: bytes
| u
< row
: bytes
| kp
: bytes
| uv
: int64 0
: bytes
| u
< row
: bytes
| kp
: bytes
| uv
: int64 0
: bytes
| v
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< error
| 1072 42000 Key column 'v' doesn't exist in table
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY); CREATE TABLE z (b INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| z
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 8 0
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| z
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 8 0
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
< row
: bytes
| z
: bytes
| PRIMARY
: int64 0
: bytes
| b
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 0 0
> exec
| DROP DATABASE IF EXISTS shop
< ok 8 0
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v)); CREATE TABLE c (xv INT, FOREIGN KEY (xv) REFERENCES x (v))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| c
: bytes
| xv
: bytes
|
: bytes
|
: bytes
| 'NULL'
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
| c
: bytes
| xv
: int64 1
: bytes
| xv
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| x
: bytes
| uv
: int64 0
: bytes
| v
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
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| DROP TABLE x
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (id INT PRIMARY KEY); CREATE TABLE c (xid INT, FOREIGN KEY (xid) REFERENCES other.x (id))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY BINARY TABLE_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "TABLE_TYPE" "VARCHAR"
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
| c
: bytes
| xid
: bytes
|
: bytes
|
: bytes
| 'NULL'
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
| c
: bytes
| xid
: int64 1
: bytes
| xid
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| DROP DATABASE IF EXISTS other
< ok 1 0
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE TABLE y (a INT PRIMARY KEY)
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
| 'NULL'
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| y
: bytes
| PRIMARY
: int64 0
: bytes
| a
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
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE x (a INT)
< ok 0 0
> exec
| RENAME TABLE x TO y
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id
< ok 0 0
> exec
| CREATE TEMPORARY TABLE tmp (a INT)
< ok 0 0
> exec
| RENAME TABLE tmp TO x
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (s DATE, e DATE, v INT, PRIMARY KEY (s), PERIOD FOR p (s, e))
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
< row
: bytes
| x
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
| e
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| s
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
< row
: bytes
| x
: bytes
| p
: bytes
| `s` < `e`
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
> exec
| ALTER TABLE x DROP COLUMN s, DROP COLUMN e
< error
| 1054 42S22 Unknown column 's' in 'p'
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| s
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, LOCK=DEFAULT
< error
| 1846 0A000 ALGORITHM=INSTANT is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try ALGORITHM=COPY
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, ALGORITHM=INSTANT
< error
| 1846 0A000 ALGORITHM=INSTANT is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try ALGORITHM=COPY
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INPLACE
< error
| 1846 0A000 ALGORITHM=INPLACE is not supported. Reason: Dropping a primary key is not allowed without also adding a new primary key. Try ALGORITHM=COPY
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, LOCK=NONE
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=`DEFAULT`
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: bytes
| id
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
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT); CREATE DATABASE other; CREATE TABLE other.x (id INT, v INT)
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
< row
: bytes
| x
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
< row
: bytes
| x
: bytes
| PRIMARY
: int64 0
: bytes
| id
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
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "TEXT"
< row
: null
