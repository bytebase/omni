-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestPriorBackupAgainstTiDB/a_column_USING_names default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN p USING (note) WHERE t.id = 1
< error
| 1054 42S22 Unknown column 'note' in 'from clause'

== TestPriorBackupAgainstTiDB/a_column_USING_names,_a_derived_table_on_its_side_may_hold default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN (SELECT 'x' AS note) d ON d.note = t.code JOIN p USING (note) WHERE t.id = 3
< ok 0 0

== TestPriorBackupAgainstTiDB/a_column_in_a_subquery,_not_judged default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`plain`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `plain`.* FROM plain WHERE id IN (SELECT gone FROM t)
< error
| 1054 42S22 Unknown column 'gone' in 'field list'

== TestPriorBackupAgainstTiDB/a_column_of_a_joined_table default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN p ON p.gone = t.n WHERE t.id = 3
< error
| 1054 42S22 Unknown column 'p.gone' in 'on clause'

== TestPriorBackupAgainstTiDB/a_column_of_a_table_the_change_recreates default#1
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
| CREATE TABLE `bbdataarchive`.`bk_2_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_2_0` SELECT `nokey`.* FROM nokey WHERE b = 1
< error
| 1054 42S22 Unknown column 'b' in 'where clause'

== TestPriorBackupAgainstTiDB/a_column_the_SET_clause_names,_not_copied default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE id = 1
< ok 1 0

== TestPriorBackupAgainstTiDB/a_column_the_change_adds,_in_ORDER_BY default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE a = 1 ORDER BY gone LIMIT 1
< error
| 1054 42S22 Unknown column 'gone' in 'order clause'

== TestPriorBackupAgainstTiDB/a_column_the_change_adds,_qualified default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nokey`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `nokey`.* FROM nokey WHERE nokey.gone = 1
< error
| 1054 42S22 Unknown column 'nokey.gone' in 'where clause'

== TestPriorBackupAgainstTiDB/a_column_the_change_adds,_unqualified default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE gone = 1
< error
| 1054 42S22 Unknown column 'gone' in 'where clause'

== TestPriorBackupAgainstTiDB/a_column_the_change_changes_the_name_of default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`p`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `p`.* FROM p WHERE memo = 'zz'
< error
| 1054 42S22 Unknown column 'memo' in 'where clause'

== TestPriorBackupAgainstTiDB/a_column_the_change_renames_to,_through_an_alias default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`p`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `x`.* FROM p AS x WHERE x.memo = 'a'
< error
| 1054 42S22 Unknown column 'x.memo' in 'where clause'

== TestPriorBackupAgainstTiDB/a_plain_backup default#1
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

== TestPriorBackupAgainstTiDB/a_table_the_change_creates default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`n2`
< error
| 1146 42S02 Table 'shop.n2' doesn't exist

== TestPriorBackupAgainstTiDB/a_table_the_change_creates,_in_a_subquery default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t WHERE id IN (SELECT id FROM n2)
< error
| 1146 42S02 Table 'shop.n2' doesn't exist

== TestPriorBackupAgainstTiDB/a_table_the_change_creates,_joined default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`t`
< ok 0 0
> exec
| INSERT INTO `bbdataarchive`.`bk_1_0` SELECT `t`.* FROM t JOIN n2 ON t.id = n2.id WHERE t.n = 1
< error
| 1146 42S02 Table 'shop.n2' doesn't exist

== TestPriorBackupAgainstTiDB/a_table_the_change_renames_to default#1
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
| CREATE TABLE `bbdataarchive`.`bk_1_0` LIKE `shop`.`nk`
< error
| 1146 42S02 Table 'shop.nk' doesn't exist

== TestPriorBackupAgainstTiDB/a_view_as_the_target default#1
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
| CREATE TABLE `bbdataarchive`.`bk_0_0` LIKE `shop`.`pv`
< error
| 1347 HY000 'bbdataarchive.pv' is not BASE TABLE

== TestPriorBackupAgainstTiDB/statements_before_a_USE default#1
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
