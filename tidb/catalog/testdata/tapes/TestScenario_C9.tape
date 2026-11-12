-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C9 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT GENERATED ALWAYS AS (a+1))
< ok 0 0
> query
| SELECT IFNULL(EXTRA,'') FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(EXTRA,'')" "VARCHAR"
< row
: bytes
| VIRTUAL GENERATED
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT,
|             b INT GENERATED ALWAYS AS (a+1) STORED,
|             FOREIGN KEY (b) REFERENCES p(id)
|         )
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE cv (
|             a INT,
|             b INT GENERATED ALWAYS AS (a+1) VIRTUAL,
|             FOREIGN KEY (b) REFERENCES p(id)
|         )
< error
| 3733 HY000 Foreign key 'fk_1' uses virtual column 'b' which is not supported.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, KEY (b))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL PRIMARY KEY)
< error
| 3106 HY000 'Defining a virtual generated column as primary key' is not supported for generated columns.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b TIMESTAMP GENERATED ALWAYS AS (NOW()) VIRTUAL)
< error
| 3102 HY000 Expression of generated column 'b' contains a disallowed function.
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, b INT GENERATED ALWAYS AS (a+1) STORED UNIQUE)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL UNIQUE)
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.STATISTICS
|                  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t1'
|                  AND COLUMN_NAME='b' AND NON_UNIQUE=0
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> query
| SELECT COUNT(*) FROM information_schema.STATISTICS
|                  WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2'
|                  AND COLUMN_NAME='b' AND NON_UNIQUE=0
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
|             a INT NULL,
|             b INT GENERATED ALWAYS AS (a+1) VIRTUAL NOT NULL
|         )
< ok 0 0
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| NO
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (
|             a INT,
|             b INT GENERATED ALWAYS AS (a+1) STORED,
|             UNIQUE KEY (b)
|         )
< ok 0 0
> exec
| CREATE TABLE c (x INT, FOREIGN KEY (x) REFERENCES p(b))
< ok 0 0
> query
| SELECT COUNT(*) FROM information_schema.REFERENTIAL_CONSTRAINTS
|              WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='c'
|              AND REFERENCED_TABLE_NAME='p'
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (
|             a VARCHAR(10) CHARACTER SET latin1,
|             b VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a, 'x')) VIRTUAL
|         )
< ok 0 0
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4
