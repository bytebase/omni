-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_AX default#1
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
| CREATE TABLE t (a INT, b INT, c INT)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN d INT
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN e INT FIRST
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN f INT AFTER a
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
|          WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME="t" ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| e
< row
: bytes
| a
< row
: bytes
| f
< row
: bytes
| b
< row
: bytes
| c
< row
: bytes
| d
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
| CREATE TABLE t (a INT, b INT, c INT, INDEX idx_a(a), INDEX idx_ab(a, b), INDEX idx_bc(b, c))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
|          WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME="t" AND INDEX_NAME <> 'PRIMARY'
|          ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| idx_ab
< row
: bytes
| idx_bc
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 1090 42000 You can't delete all columns with ALTER TABLE; use DROP TABLE instead
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
| CREATE TABLE t1 (a INT, b INT, CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a + 1))
< ok 0 0
> exec
| ALTER TABLE t1 DROP COLUMN a
< ok 0 0
> exec
| ALTER TABLE t2 DROP COLUMN a
< error
| 1054 42S22 Unknown column 'a' in 'GENERATED ALWAYS AS'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT PRIMARY KEY, b INT NOT NULL DEFAULT 5 COMMENT 'x')
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN b BIGINT
< ok 0 0
> query
| SELECT IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_DEFAULT" "TEXT"
< column "COLUMN_COMMENT" "VARCHAR"
< row
: bytes
| YES
: bytes
| NULL
: bytes
|
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| ALTER TABLE t CHANGE COLUMN a b BIGINT NOT NULL
< ok 0 0
> query
| SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "COLUMN_NAME" "VARCHAR"
< column "DATA_TYPE" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| b
: bytes
| bigint
: bytes
| NO
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
| CREATE TABLE t (a INT, b INT, INDEX (a))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX (b)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX (a, b)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
|          WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME="t" AND INDEX_NAME <> 'PRIMARY'
|          ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| a_2
< row
: bytes
| b
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
| CREATE TABLE t (a INT, b VARCHAR(100)) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE (a)
< ok 0 0
> exec
| ALTER TABLE t ADD KEY (b)
< ok 0 0
> exec
| ALTER TABLE t ADD FULLTEXT (b)
< ok 0 0
> query
| SELECT INDEX_NAME, NON_UNIQUE, INDEX_TYPE FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| a
: int64 0
: bytes
| BTREE
< row
: bytes
| b
: int64 1
: bytes
| BTREE
< row
: bytes
| b_2
: int64 1
: bytes
| FULLTEXT
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
| CREATE TABLE parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t1 (a INT)
< ok 0 0
> exec
| ALTER TABLE t1 ADD FOREIGN KEY (a) REFERENCES parent(id)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT REFERENCES parent(id))
< ok 0 0
> query
| SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1'
|               AND CONSTRAINT_TYPE='FOREIGN KEY'
|         ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t1_ibfk_1
> query
| SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2'
|               AND CONSTRAINT_TYPE='FOREIGN KEY'
|         ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t2_ibfk_1
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
| CREATE TABLE t (a INT NOT NULL DEFAULT 5 COMMENT 'hi', INDEX (a))
< ok 0 0
> exec
| ALTER TABLE t RENAME COLUMN a TO aa
< ok 0 0
> query
| SELECT COLUMN_NAME, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT
|             FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "COLUMN_DEFAULT" "TEXT"
< column "COLUMN_COMMENT" "VARCHAR"
< row
: bytes
| aa
: bytes
| NO
: bytes
| 5
: bytes
| hi
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
| CREATE TABLE t (a INT, b INT, INDEX old_name (a, b))
< ok 0 0
> exec
| ALTER TABLE t RENAME INDEX old_name TO new_name
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
|          WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME="t" AND INDEX_NAME <> 'PRIMARY'
|          ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| new_name
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
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| ALTER TABLE t RENAME TO t2
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.TABLES
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 0
> query
| SELECT COUNT(*) FROM information_schema.TABLES
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
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
| CREATE TABLE t (a INT NOT NULL, b INT)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 5
< ok 0 0
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: bytes
| 5
> exec
| ALTER TABLE t ALTER COLUMN a DROP DEFAULT
< ok 0 0
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> exec
| ALTER TABLE t MODIFY b INT INVISIBLE
< ok 0 0
> query
| SELECT EXTRA FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "EXTRA" "VARCHAR"
< row
: bytes
| INVISIBLE
> exec
| ALTER TABLE t MODIFY b INT
< ok 0 0
> query
| SELECT EXTRA FROM information_schema.COLUMNS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "EXTRA" "VARCHAR"
< row
: bytes
|
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
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE t
|   ADD COLUMN c INT,
|   DROP COLUMN b,
|   ADD INDEX (c),
|   RENAME COLUMN a TO aa
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
|          WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME="t" ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| aa
< row
: bytes
| c
> query
| SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
|          WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME="t" AND INDEX_NAME <> 'PRIMARY'
|          ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| c
