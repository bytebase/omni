-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_PS default#1
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
| child_ibfk_5
< row
: bytes
| fk_2
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
| fk_2
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
| CREATE TABLE t (a DATETIME(6) DEFAULT NOW())
< error
| 1067 42000 Invalid default value for 'a'
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
|             CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id),
|             b INT,
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
| CREATE TABLE t1 (a INT, CONSTRAINT my_rule CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (b INT, CONSTRAINT my_rule CHECK (b > 0))
< ok 0 0
