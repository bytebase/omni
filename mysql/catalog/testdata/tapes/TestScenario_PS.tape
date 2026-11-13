-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
| t_chk_1
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
| t_chk_20
< row
: bytes
| t_chk_21
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
< error
| 1067 42000 Invalid default value for 'a'
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
| 1826 HY000 Duplicate foreign key constraint name 'c_ibfk_1'
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
< error
| 3822 HY000 Duplicate check constraint name 'my_rule'.
