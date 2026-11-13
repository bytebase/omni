-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestPriorBackupAgainstMariaDB default#1
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
< ok 2 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`n2`
< error
| 1146 42S02 Table 'shop.n2' doesn't exist
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1
< error
| 1146 42S02 Table 'shop.n2' doesn't exist
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE id IN (SELECT id FROM n2)
< error
| 1146 42S02 Table 'shop.n2' doesn't exist
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nk`
< error
| 1146 42S02 Table 'shop.nk' doesn't exist
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE gone = 1
< error
| 1054 42S22 Unknown column 'gone' in 'WHERE'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE nokey.gone = 1
< error
| 1054 42S22 Unknown column 'nokey.gone' in 'WHERE'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE a = 1 ORDER BY gone LIMIT 1
< error
| 1054 42S22 Unknown column 'gone' in 'ORDER BY'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`p`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `x`.* FROM p AS x WHERE x.memo = 'a'
< error
| 1054 42S22 Unknown column 'x.memo' in 'WHERE'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`p`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `p`.* FROM p WHERE memo = 'zz'
< error
| 1054 42S22 Unknown column 'memo' in 'WHERE'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN p ON p.gone = t.n WHERE t.id = 3
< error
| 1054 42S22 Unknown column 'p.gone' in 'ON'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_2_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_2_0` SELECT `nokey`.* FROM nokey WHERE b = 1
< error
| 1054 42S22 Unknown column 'b' in 'WHERE'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`plain`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `plain`.* FROM plain WHERE id IN (SELECT gone FROM t)
< error
| 1054 42S22 Unknown column 'gone' in 'SELECT'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_0_0` LIKE `shop`.`pv`
< error
| 1347 HY000 'shop.pv' is not of type 'BASE TABLE'
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE id = 1
< ok 1 0
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN p USING (note) WHERE t.id = 1
< error
| 1054 42S22 Unknown column 'note' in 'FROM'
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN (SELECT 1 AS note) d ON d.note = t.id JOIN p USING (note) WHERE t.id = 3
< ok 0 0
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
< ok 1 0
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| CREATE DATABASE bbdataarchive
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
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
> exec
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE gone = 1
< error
| 1054 42S22 Unknown column 'gone' in 'WHERE'
> exec
| CREATE TABLE `bbdataarchive`.`bk_3_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_3_0` SELECT `t`.* FROM shop.t WHERE id = 1
< ok 1 0
