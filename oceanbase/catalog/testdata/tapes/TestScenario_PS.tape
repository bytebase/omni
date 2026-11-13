-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| t_chk_5
< row
: bytes
| t_OBCHECK_0000000000000001
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
| t_OBCHECK_0000000000000002
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
| child_ibfk_5
< row
: bytes
| child_OBFK_0000000000000003
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
| child_OBFK_0000000000000004
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
< ok 2 0
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
