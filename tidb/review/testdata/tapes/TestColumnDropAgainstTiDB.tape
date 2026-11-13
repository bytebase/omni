-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestColumnDropAgainstTiDB/a_key_added_back default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| id
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN id, ADD PRIMARY KEY (v)
< error
| 8200 HY000 Unsupported drop integer primary key
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstTiDB/a_key_of_one_column,_under_IF_EXISTS default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| id
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN IF EXISTS id
< error
| 8200 HY000 Unsupported drop integer primary key
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstTiDB/a_key_of_one_integer_column default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| id
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 8200 HY000 Unsupported drop integer primary key
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstTiDB/a_key_of_one_integer_column,_nonclustered default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| id
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN id
< error
| 8200 HY000 can't drop column id with composite index covered or Primary Key covered now
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstTiDB/a_key_of_one_string_column default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (s VARCHAR(20) PRIMARY KEY, v INT)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| s
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN s
< error
| 8200 HY000 can't drop column s with composite index covered or Primary Key covered now
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| s

== TestColumnDropAgainstTiDB/both_columns_of_a_key_of_two default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b) NONCLUSTERED)
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| a
< row
: bytes
| x
: bytes
| b
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN a, DROP COLUMN b
< error
| 8200 HY000 can't drop column a with composite index covered or Primary Key covered now
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstTiDB/one_column_of_a_key_of_two default#1
> exec
| CREATE DATABASE shop
< ok 0 0
> exec
| CREATE DATABASE bbdataarchive
< ok 0 0
> exec
| USE shop
< ok 0 0
> exec
|
| CREATE TABLE p (id INT, note VARCHAR(20), PRIMARY KEY (id) NONCLUSTERED);
| CREATE TABLE t (id INT, code VARCHAR(20), p_id INT, n INT,
|   PRIMARY KEY (id) NONCLUSTERED, UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
| CREATE TABLE nokey (a INT);
| CREATE TABLE plain (id INT, v INT, PRIMARY KEY (id) NONCLUSTERED);
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
| INSERT INTO nokey VALUES (1), (2);
| INSERT INTO plain VALUES (1, 1), (2, 2);
| CREATE VIEW pv AS SELECT id, note FROM p;
|
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b))
< ok 0 0
> query
| SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME
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
| SELECT TABLE_NAME, COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'shop' ORDER BY TABLE_NAME, ORDINAL_POSITION
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| nokey
: bytes
| a
< row
: bytes
| p
: bytes
| id
< row
: bytes
| p
: bytes
| note
< row
: bytes
| plain
: bytes
| id
< row
: bytes
| plain
: bytes
| v
< row
: bytes
| pv
: bytes
| id
< row
: bytes
| pv
: bytes
| note
< row
: bytes
| t
: bytes
| id
< row
: bytes
| t
: bytes
| code
< row
: bytes
| t
: bytes
| p_id
< row
: bytes
| t
: bytes
| n
< row
: bytes
| x
: bytes
| a
< row
: bytes
| x
: bytes
| b
< row
: bytes
| x
: bytes
| v
> query
| SELECT TABLE_NAME, INDEX_NAME, MIN(NON_UNIQUE) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'shop' GROUP BY TABLE_NAME, INDEX_NAME ORDER BY TABLE_NAME, INDEX_NAME
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "MIN(NON_UNIQUE)" "VARCHAR"
< row
: bytes
| p
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| plain
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| PRIMARY
: bytes
| 0
< row
: bytes
| t
: bytes
| fk_p
: bytes
| 1
< row
: bytes
| t
: bytes
| idx_n
: bytes
| 1
< row
: bytes
| t
: bytes
| uk_code
: bytes
| 0
< row
: bytes
| x
: bytes
| PRIMARY
: bytes
| 0
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME, REFERENCED_TABLE_NAME FROM information_schema.REFERENTIAL_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| fk_p
: bytes
| p
> query
| SELECT TABLE_NAME, CONSTRAINT_NAME FROM information_schema.TIDB_CHECK_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = 'shop'
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| ck_n
> query
| SELECT EVENT_OBJECT_TABLE, TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = 'shop'
< result
< column "EVENT_OBJECT_TABLE" "VARCHAR"
< column "TRIGGER_NAME" "VARCHAR"
> query
| SELECT @@lower_case_table_names
< result
< column "@@lower_case_table_names" "VARCHAR"
< row
: bytes
| 2
> exec
| ALTER TABLE x DROP COLUMN a
< error
| 8200 HY000 can't drop column a with composite index covered or Primary Key covered now
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = 'x' AND INDEX_NAME = 'PRIMARY'
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
