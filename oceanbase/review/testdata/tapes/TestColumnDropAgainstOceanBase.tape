-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestColumnDropAgainstOceanBase/a_RENAME_TABLE_of_a_temporary_table_onto_its_name_after,_which_MySQL_does_not_rename default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_USE_of_a_database_the_server_lacks_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| USE nosuch_db;
< error
| 1049 42000 Unknown database 'nosuch_db'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_USE_of_its_database_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| USE shop;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_check_of_another_column_added_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1);
< ok 1 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x ADD CHECK (v > 0);
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_check_of_both_columns_of_the_key,_both_dropped default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), CONSTRAINT ck CHECK (a < b));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN a, DROP COLUMN b;
< error
| 3959 HY000 Check constraint 'ck' uses column 'a', hence column cannot be dropped or renamed.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/a_check_of_it_alone default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY CHECK (id > 0), v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_check_reads_it_and_another_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT ck CHECK (id < v));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 3959 HY000 Check constraint 'ck' uses column 'id', hence column cannot be dropped or renamed.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_column_of_only_some_unique_keys_dropped_before_from_a_table_without_a_primary_key_partitioned_by_KEY_() default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE kp (u INT NOT NULL, v INT NOT NULL, w INT, UNIQUE KEY uu (u), UNIQUE KEY uv (u, v)) PARTITION BY KEY () PARTITIONS 2;
< error
| 1488 HY000 Field in list of fields for partition function not found in table

== TestColumnDropAgainstOceanBase/a_column_of_the_key_held_by_a_prefix_dropped_before_from_a_table_partitioned_by_KEY_(),_which_leaves_the_prefix_out default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE kp (a VARCHAR(20), b INT, v INT, PRIMARY KEY (a(10), b), UNIQUE KEY ub (b)) PARTITION BY KEY () PARTITIONS 2;
< error
| 1064 42000 …near '(10), b), UNIQUE KEY ub (b)) PARTITION BY KEY () PARTITIONS 2;' at line 1

== TestColumnDropAgainstOceanBase/a_column_out_of_any_key_dropped_before_from_a_table_without_a_primary_key_partitioned_by_KEY_() default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2;
< error
| 1488 HY000 Field in list of fields for partition function not found in table

== TestColumnDropAgainstOceanBase/a_column_out_of_the_key_dropped_before_from_a_table_partitioned_by_KEY_() default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE kp (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2;
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE kp DROP COLUMN v;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_column_the_table_lacks_inserted_into_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| INSERT INTO x (missing) VALUES (1);
< error
| 1054 42S22 Unknown column 'missing' in 'field list'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_database_the_server_lacks_dropped_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| DROP DATABASE missing_db;
< error
| 1008 HY000 Can't drop database 'missing_db'; database doesn't exist
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_drop_after_of_it,_which_a_foreign_key_naming_a_table_of_its_name_in_another_database_does_not_refuse default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE DATABASE other;
< ok 1 0
> exec
| CREATE TABLE other.x (id INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE c (xid INT, FOREIGN KEY (xid) REFERENCES other.x (id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_drop_after_of_it,_which_a_foreign_key_of_the_database_refuses default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v));
< ok 0 0
> exec
| CREATE TABLE c (xv INT, FOREIGN KEY (xv) REFERENCES x (v));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_of_another_database_references_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE DATABASE other;
< ok 1 0
> exec
| CREATE TABLE other.ch (x_id INT, CONSTRAINT fk_och FOREIGN KEY (x_id) REFERENCES shop.x (id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 5320 HY000 'id': used in a foreign key constraint
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_of_another_table_references_another_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY uv (v));
< ok 0 0
> exec
| CREATE TABLE ch (x_v INT, CONSTRAINT fk_ch FOREIGN KEY (x_v) REFERENCES x (v));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_of_another_table_references_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 5320 HY000 'id': used in a foreign key constraint
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_of_another_table_references_it,_with_foreign_key_checks_off default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE ch (x_id INT, CONSTRAINT fk_ch FOREIGN KEY (x_id) REFERENCES x (id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET foreign_key_checks = 0;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 5320 HY000 'id': used in a foreign key constraint
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_of_the_table_holds_another_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (v) REFERENCES p (id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_of_the_table_holds_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, CONSTRAINT fk_x FOREIGN KEY (id) REFERENCES p (id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 5320 HY000 'id': used in a foreign key constraint
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_foreign_key_to_it_added_and_dropped_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE ch (k INT PRIMARY KEY, x_id INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE ch ADD CONSTRAINT f FOREIGN KEY (x_id) REFERENCES x (id);
< ok 0 0
> exec
| ALTER TABLE ch DROP FOREIGN KEY f;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_functional_index_reads_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY fx ((id + v)));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 3837 HY000 Column 'id' has a functional index dependency and cannot be dropped or renamed.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_generated_column_naming_it_only_in_a_literal default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g VARCHAR(20) AS (CONCAT('id', v)));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_generated_column_reading_it,_dropped_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN g;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_generated_column_reading_it,_dropped_too default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1) STORED);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN g;
< error
| 3108 HY000 Column 'id' has a generated column dependency
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_generated_column_reads_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, g INT AS (id + 1));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 3108 HY000 Column 'id' has a generated column dependency
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_key_added_after_a_USE_of_another_database,_to_a_table_of_its_name_there default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE DATABASE other;
< ok 1 0
> exec
| CREATE TABLE other.x (id INT, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_key_added_back_and_a_rename_onto_a_table_the_database_holds_after,_which_fail default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_key_is_added_back default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ADD PRIMARY KEY (v);
< error
| 1068 42000 Multiple primary key defined
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_key_of_a_column_prefix default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (s VARCHAR(50), v INT, PRIMARY KEY (s(10)));
< error
| 1064 42000 …near '(10)));' at line 1

== TestColumnDropAgainstOceanBase/a_key_of_one_AUTO_INCREMENT_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT AUTO_INCREMENT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO x (v) VALUES (1), (2);
< ok 2 1
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_key_of_one_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1), (2, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_key_of_one_column,_with_a_unique_key_on_a_NOT_NULL_column_left default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, UNIQUE KEY uv (v));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_later_statement_adds_a_key_back default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_later_statement_that_adds_a_key_back,_but_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_literal_default_names_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(10) DEFAULT 'id');
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_read-only_transaction_started_before,_which_the_drop_commits default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| START TRANSACTION READ ONLY;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1792 25006 Cannot execute statement in a READ ONLY transaction.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rename_after_onto_a_table_an_earlier_pair_moves_away default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rename_after_with_a_pair_onto_a_name_an_earlier_pair_made,_which_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE z (b INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rename_after_with_a_pair_onto_a_name_an_earlier_pair_made_and_a_later_one_moved_away default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE z (b INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rename_after_with_a_pair_onto_a_name_an_earlier_pair_moved_a_table_back_onto,_which_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE z (b INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rename_of_its_name_onto_a_table_the_database_holds_after,_which_MySQL_runs_on_it_though_a_temporary_table_of_its_name_exists,_and_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rename_onto_a_table_the_database_holds_after,_which_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_rollback_to_a_savepoint_the_transaction_lacks_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| BEGIN;
< ok 0 0
> exec
| ROLLBACK TO SAVEPOINT nosuch;
< error
| 1305 42000 savepoint does not exist
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_row_inserted_before_into_a_table_the_database_lacks default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| INSERT INTO missing VALUES (1);
< error
| 1146 42S02 Table 'shop.missing' doesn't exist
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_scope_for_a_system_variable_before_a_user_variable default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET SESSION foreign_key_checks = 1, @x = 1;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_statement_on_another_table_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE y ADD COLUMN b INT;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_statement_prepared_and_deallocated_after,_which_runs_nothing default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_statement_prepared_before_of_a_table_the_database_lacks default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| PREPARE s FROM 'SELECT * FROM missing';
< error
| 1146 42S02 Table 'shop.missing' doesn't exist
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_CREATE_TABLE_IF_NOT_EXISTS_leaves,_whose_foreign_key_would_reference_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE c (k INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| CREATE TABLE IF NOT EXISTS c (k INT PRIMARY KEY, FOREIGN KEY (k) REFERENCES x (id));
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_altered_to_be_renamed_onto_it_after,_which_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_made_before_whose_column's_REFERENCES,_which_MySQL_ignores,_names_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| CREATE TABLE c (k INT REFERENCES x (id), PRIMARY KEY (k));
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_made_of_its_name_after,_which_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_of_another_database_renamed_to_its_name_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE DATABASE other;
< ok 1 0
> exec
| CREATE TABLE other.t (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE other.t RENAME TO x;
< error
| 1050 42S01 Table 'x' already exists
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_renamed_onto_it_after,_which_fails default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_table_the_database_lacks_dropped_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| DROP TABLE missing;
< error
| 1051 42S02 Unknown table 'shop.missing'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_temporary_table_of_its_name_dropped_after,_and_then_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_temporary_table_of_its_name_dropped_twice default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| DROP TEMPORARY TABLE IF EXISTS x;
< ok 0 0
> exec
| DROP TEMPORARY TABLE x;
< error
| 1051 42S02 Unknown table 'shop.x'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_temporary_table_of_its_name_made_after,_given_a_key,_and_dropped default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_temporary_table_renamed_onto_its_name_after,_and_dropped default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_truncate_after default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_truncate_of_a_table_a_foreign_key_references_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE tp (a INT PRIMARY KEY);
< ok 0 0
> exec
| CREATE TABLE tc (a INT, FOREIGN KEY (a) REFERENCES tp (a));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| TRUNCATE TABLE tp;
< error
| 1701 42000 Cannot truncate a table referenced in a foreign key constraint (`shop`.`tc`, CONSTRAINT `tc_OBFK_0000000000000001` FOREIGN KEY (`a`) REFERENCES `shop`.`tp` (`a`))
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_unique_key_of_it_alone default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1), (2, 1);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_unique_key_of_it_and_a_prefix_of_another_column,_beside_a_unique_key_of_the_other default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v VARCHAR(20) NOT NULL, UNIQUE KEY u (v), UNIQUE KEY uv (id, v(2)));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 'xx1'), (2, 'xx2');
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_unique_key_of_it_and_another_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1), (2, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_unique_key_of_it_and_another_column,_beside_a_unique_key_of_the_other default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT NOT NULL, w INT, UNIQUE KEY u (v), UNIQUE KEY uv (id, v));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1, 1), (2, 2, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_unique_key_of_it_and_another_column,_both_dropped default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT, UNIQUE KEY u (id, v));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1, 1), (2, 1, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN v;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_unique_key_of_it_and_another_column,_the_other_repeating default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, UNIQUE KEY u (id, v));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1), (2, 1);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_user_variable_named_sql_require_primary_key_set_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET @sql_require_primary_key = 1;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_user_variable_set_to_a_double_out_of_range_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET @x = 1e400;
< error
| 1367 22007 Illegal double '1e400' value found during parsing
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_user_variable_set_to_a_hexadecimal_value_of_an_uppercase_prefix_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET @x = 0X41;
< error
| 1064 42000 …near '0X41;' at line 1
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_user_variable_set_to_a_malformed_hexadecimal_value_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET @x = X'4';
< error
| 1064 42000 …near 'X'4';' at line 1
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/a_view_or_a_trigger_reads_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE VIEW xv AS SELECT id FROM x;
< ok 0 0
> exec
| CREATE TRIGGER xt BEFORE INSERT ON x FOR EACH ROW SET NEW.v = NEW.id;
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/an_expression_default_reads_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT DEFAULT (id + 1));
< error
| 1235 0A000 Column reference in default value expression not supported

== TestColumnDropAgainstOceanBase/an_invisible_column_and_a_visible_one_left default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE, w INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/another_column_dropped_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, w INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN v;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/another_index_holds_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT, KEY iv (v, id));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/another_table_locked_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| LOCK TABLES y WRITE;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/both_columns_of_a_key_of_two default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1, 1), (2, 1, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN b, DROP COLUMN a;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/both_columns_of_a_key_of_two,_in_two_statements default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b));
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN a;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/both_columns_of_a_key_of_two,_in_two_statements,_a_unique_key_of_a_prefix_of_the_other_keeping_its_values_apart default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b VARCHAR(20), v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b(2)));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 'aa1', 1), (1, 'bb2', 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN a;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/both_columns_of_a_key_of_two,_in_two_statements,_a_unique_key_of_the_other_keeping_its_values_apart default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b), UNIQUE KEY ub (b));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1, 1), (1, 2, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN a;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/foreign_key_checks_set_with_:=_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET foreign_key_checks := OFF;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/it_is_added_back_under_its_name default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ADD COLUMN id INT;
< error
| 1060 42S21 Duplicate column name 'id'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/it_is_the_only_column default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1090 42000 You can't delete all columns with ALTER TABLE; use DROP TABLE instead
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/one_column_of_a_key_of_two default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1, 1), (1, 2, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN a;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/one_column_of_a_key_of_two,_the_other_repeating default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b INT, v INT, PRIMARY KEY (a, b));
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1, 1), (2, 1, 2);
< ok 2 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b
> exec
| ALTER TABLE x DROP COLUMN a;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| a,b

== TestColumnDropAgainstOceanBase/one_column_of_a_key_of_two_whose_other_is_a_prefix,_beside_a_unique_key_of_the_whole_other default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (a INT, b VARCHAR(20), v INT, PRIMARY KEY (a, b(2)), UNIQUE KEY ub (b));
< error
| 1064 42000 …near '(2)), UNIQUE KEY ub (b));' at line 1

== TestColumnDropAgainstOceanBase/only_an_invisible_column_left default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT INVISIBLE);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 4028 HY000 table must have at least one column that is not invisible
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/rows_of_the_table_updated_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO x VALUES (1, 1);
< ok 1 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| UPDATE x SET v = id + 1 WHERE id IN (SELECT id FROM (SELECT id FROM x) d);
< ok 1 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/sql_generate_invisible_primary_key_on,_then_off,_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET SESSION sql_generate_invisible_primary_key = ON;
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/sql_generate_invisible_primary_key_on_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET SESSION sql_generate_invisible_primary_key = ON;
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/sql_require_primary_key_off_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET SESSION sql_require_primary_key = OFF;
< error
| 1193 HY000 Unknown system variable 'sql_require_primary_key'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/sql_require_primary_key_on_after default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/sql_require_primary_key_on_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET SESSION sql_require_primary_key = ON;
< error
| 1193 HY000 Unknown system variable 'sql_require_primary_key'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/sql_require_primary_key_set_with_a_variable_the_server_reads_only_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SET SESSION sql_require_primary_key = OFF, version = 'x';
< error
| 1193 HY000 Unknown system variable 'sql_require_primary_key'
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/statements_on_another_table,_on_it_that_give_it_no_key,_or_on_none_after default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE y (a INT PRIMARY KEY);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/statements_that_cannot_fail_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| BEGIN;
< ok 0 0
> exec
| SAVEPOINT s;
< ok 0 0
> exec
| ROLLBACK;
< ok 0 0
> exec
| COMMIT;
< ok 0 0
> exec
| UNLOCK TABLES;
< ok 0 0
> exec
| ANALYZE TABLE missing;
< error
| 1146 42S02 Table doesn't exist
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/tables_flushed_after default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/tables_flushed_with_a_read_lock_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| FLUSH TABLES WITH READ LOCK;
< error
| 1235 0A000 This statement is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/tables_locked_after default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/tables_locked_for_reading_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| LOCK TABLES x READ;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/the_database's_tables_shown_before default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| SHOW TABLES FROM shop;
< ok 0 0
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/the_table_is_partitioned_by_KEY_(),_its_primary_key default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY KEY () PARTITIONS 2;
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 3855 HY000 Column 'id' has a partitioning function dependency and cannot be dropped or renamed.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/the_table_is_partitioned_by_it default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT) PARTITION BY HASH (id) PARTITIONS 2;
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id;
< error
| 3855 HY000 Column 'id' has a partitioning function dependency and cannot be dropped or renamed.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/the_unique_key's_column_dropped_before_from_a_table_without_a_primary_key_partitioned_by_KEY_() default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| CREATE TABLE kp (u INT NOT NULL, v INT, w INT, UNIQUE KEY uu (u)) PARTITION BY KEY () PARTITIONS 2;
< error
| 1488 HY000 Field in list of fields for partition function not found in table

== TestColumnDropAgainstOceanBase/twice default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN id;
< error
| 1091 42000 Can't DROP 'id'; check that column/key exists
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_ALGORITHM=INPLACE default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INPLACE;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_ALGORITHM_and_LOCK_left_at_DEFAULT default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, LOCK DEFAULT;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_ALGORITHM_set_to_DEFAULT,_then_INSTANT default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=DEFAULT, ALGORITHM=INSTANT;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_ALGORITHM_set_to_INSTANT,_then_DEFAULT default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, ALGORITHM=DEFAULT;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_ALGORITHM_set_to_INSTANT_before_LOCK_left_at_DEFAULT default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=INSTANT, LOCK=DEFAULT;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_LOCK=NONE default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, LOCK=NONE;
< error
| 1235 0A000 drop rowkey column is not supported
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_a_check_the_table_lacks_altered default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x ALTER CHECK missing ENFORCED, DROP COLUMN id;
< error
| 3821 HY000 Check constraint 'missing' is not found in the table.
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_a_column_the_table_lacks default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, DROP COLUMN missing;
< error
| 1091 42000 Can't DROP 'missing'; check that column/key exists
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_a_partition_dropped default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT, v INT, PRIMARY KEY (id)) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, DROP PARTITION p0;
< error
| 1064 42000 …near 'PARTITION p0;' at line 1
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id

== TestColumnDropAgainstOceanBase/with_a_quoted_DEFAULT default#1
> exec
| CREATE DATABASE shop
< ok 1 0
> exec
| USE shop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, note VARCHAR(20));
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, code VARCHAR(20), p_id INT, n INT,
|   UNIQUE KEY uk_code (code), KEY idx_n (n),
|   CONSTRAINT fk_p FOREIGN KEY (p_id) REFERENCES p (id),
|   CONSTRAINT ck_n CHECK (n > 0));
< ok 0 0
> exec
| CREATE TABLE nokey (a INT);
< ok 0 0
> exec
| CREATE TABLE plain (id INT PRIMARY KEY, v INT);
< ok 0 0
> exec
| INSERT INTO p VALUES (1, 'a'), (2, 'b');
< ok 2 0
> exec
| INSERT INTO t VALUES (1, 'x', 1, 1), (2, 'y', 2, 2);
< ok 2 0
> exec
| INSERT INTO nokey VALUES (1), (2);
< ok 2 0
> exec
| INSERT INTO plain VALUES (1, 1), (2, 2);
< ok 2 0
> exec
| CREATE VIEW pv AS SELECT id, note FROM p;
< ok 0 0
> exec
| CREATE TABLE x (id INT PRIMARY KEY, v INT);
< ok 0 0
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
> exec
| ALTER TABLE x DROP COLUMN id, ALGORITHM=`DEFAULT`;
< error
| 1064 42000 …near '`DEFAULT`;' at line 1
> query
| SELECT GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) FROM information_schema.STATISTICS
| 		WHERE TABLE_SCHEMA = 'shop' AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
: string
| x
< result
< column "GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX)" "VARCHAR"
< row
: bytes
| id
