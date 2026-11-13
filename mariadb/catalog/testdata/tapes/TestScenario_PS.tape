-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_PS default#1
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
| CREATE TABLE t (
|             a INT,
|             CONSTRAINT t_chk_5 CHECK (a > 0),
|             b INT,
|             CHECK (b < 100)
|         )
< ok 0 0
> query
|
|             SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb'
|             ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| CONSTRAINT_1
< row
: bytes
| t_chk_5
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
| CREATE TABLE t (a INT, b INT, CONSTRAINT t_chk_20 CHECK (a>0))
< ok 0 0
> exec
| ALTER TABLE t ADD CHECK (b>0)
< ok 0 0
> query
|
|             SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb'
|             ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| CONSTRAINT_1
< row
: bytes
| t_chk_20
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
| CREATE TABLE child (
|             a INT,
|             CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES parent(id),
|             b INT,
|             FOREIGN KEY (b) REFERENCES parent(id)
|         )
< ok 0 0
> query
|
|             SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
|               AND CONSTRAINT_TYPE='FOREIGN KEY'
|             ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| child_ibfk_1
< row
: bytes
| child_ibfk_5
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
| CREATE TABLE parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE child (
|             a INT,
|             b INT,
|             CONSTRAINT child_ibfk_20 FOREIGN KEY (a) REFERENCES parent(id)
|         )
< ok 0 0
> exec
| ALTER TABLE child ADD FOREIGN KEY (b) REFERENCES parent(id)
< ok 0 0
> query
|
|             SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
|               AND CONSTRAINT_TYPE='FOREIGN KEY'
|             ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| child_ibfk_20
< row
: bytes
| child_ibfk_21
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
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW())
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` datetime(6) DEFAULT current_timestamp(6)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
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
| CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 3
< ok 0 0
> exec
| ALTER TABLE t ADD PARTITION PARTITIONS 2
< ok 0 0
> query
|
|             SELECT PARTITION_NAME FROM information_schema.PARTITIONS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
< row
: bytes
| p1
< row
: bytes
| p2
< row
: bytes
| p3
< row
: bytes
| p4
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
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (
|             a INT,
|             CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id),
|             b INT,
|             FOREIGN KEY (b) REFERENCES p(id)
|         )
< error
| 1005 HY000 Can't create table `testdb`.`c` (errno: 150 "Foreign key constraint is incorrectly formed")
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
| CREATE TABLE t1 (a INT, CONSTRAINT my_rule CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (b INT, CONSTRAINT my_rule CHECK (b > 0))
< ok 0 0
