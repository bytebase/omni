-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestScenario_C22 default#1
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 ADD COLUMN b INT
< ok 0 0
> exec
| ALTER TABLE t1 ADD COLUMN c INT, ALGORITHM=DEFAULT, LOCK=DEFAULT
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1'
| 			 ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| c
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 ADD INDEX ix_a (a)
< ok 0 0
> exec
| ALTER TABLE t1 ADD INDEX ix_a2 (a), LOCK=DEFAULT
< ok 0 0
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB ROW_FORMAT=DYNAMIC
< ok 0 0
> exec
| ALTER TABLE t1 ADD COLUMN b VARCHAR(32) NULL
< ok 0 0
> exec
| ALTER TABLE t1 ADD COLUMN c INT NULL, ALGORITHM=INSTANT
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='c'
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT, b INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 DROP COLUMN b
< ok 0 0
> exec
| CREATE TABLE t2 (id INT PRIMARY KEY, a INT, b INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t2 DROP COLUMN b, ALGORITHM=INSTANT
< ok 0 0
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
| CREATE TABLE t1 (id INT PRIMARY KEY, old_name INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 RENAME COLUMN old_name TO new_name
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='new_name'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> exec
| ALTER TABLE t1 CHANGE COLUMN new_name newer_name INT
< ok 0 0
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 CHANGE COLUMN a a BIGINT
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "TEXT"
< row
: bytes
| bigint
> exec
| ALTER TABLE t1 MODIFY COLUMN a INT, ALGORITHM=INPLACE
< ok 0 0
> exec
| ALTER TABLE t1 MODIFY COLUMN a INT UNSIGNED, LOCK=NONE
< ok 0 0
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 DROP PRIMARY KEY, ADD PRIMARY KEY (a), ALGORITHM=INSTANT
< ok 0 0
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
| CREATE TABLE t1 (id INT PRIMARY KEY, a INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE t1 CHANGE COLUMN a a BIGINT UNSIGNED, LOCK=NONE
< ok 0 0
> exec
| ALTER TABLE t1 ADD COLUMN b INT, ALGORITHM=INSTANT, LOCK=NONE
< ok 0 0
> exec
| ALTER TABLE t1 ADD COLUMN c INT, ALGORITHM=INSTANT
< ok 0 0
