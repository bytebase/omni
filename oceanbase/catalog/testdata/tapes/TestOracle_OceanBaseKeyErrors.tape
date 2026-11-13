-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseKeyErrors default#1
> exec
| DROP DATABASE IF EXISTS obkey_00
< ok 0 0
> exec
| CREATE DATABASE obkey_00
< ok 1 0
> exec
| USE obkey_00
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))
< error
| 1062 23000 Duplicated primary key
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT PRIMARY KEY)
< error
| 1068 42000 Multiple primary key defined
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT, PRIMARY KEY (b))
< error
| 1062 23000 Duplicated primary key
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (b)
< error
| 1068 42000 Multiple primary key defined
> exec
| ALTER TABLE t ADD COLUMN c INT PRIMARY KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS obkey_01
< ok 0 0
> exec
| CREATE DATABASE obkey_01
< ok 1 0
> exec
| USE obkey_01
< ok 0 0
> exec
| CREATE TABLE c (pid INT, FOREIGN KEY (pid) REFERENCES nosuch (id))
< error
| 1146 42S02 Table doesn't exist
> exec
| CREATE TABLE p (id INT PRIMARY KEY, v INT)
< ok 0 0
> exec
| CREATE TABLE c (pid INT, FOREIGN KEY (pid) REFERENCES p (nosuch))
< error
| 5031 HY000 Column not found
> exec
| CREATE TABLE c (pid INT, FOREIGN KEY (nosuch) REFERENCES p (id))
< error
| 5031 HY000 Column not found
> exec
| CREATE TABLE c (pid INT)
< ok 0 0
> exec
| ALTER TABLE c ADD FOREIGN KEY (pid) REFERENCES nosuch (id)
< error
| 1146 42S02 Table doesn't exist
> exec
| ALTER TABLE c ADD FOREIGN KEY (pid) REFERENCES p (nosuch)
< error
| 5031 HY000 Column not found
> exec
| ALTER TABLE c ADD FOREIGN KEY (nosuch) REFERENCES p (id)
< error
| 5031 HY000 Column not found
> exec
| DROP DATABASE IF EXISTS obkey_02
< ok 0 0
> exec
| CREATE DATABASE obkey_02
< ok 1 0
> exec
| USE obkey_02
< ok 0 0
> exec
| CREATE TABLE p (id INT, v INT, w INT, x INT, y INT, KEY kv (v), UNIQUE KEY uw (w), UNIQUE KEY uxy (x, y))
< ok 0 0
> exec
| CREATE TABLE c1 (pid INT, FOREIGN KEY (pid) REFERENCES p (id))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c2 (pid INT, FOREIGN KEY (pid) REFERENCES p (v))
< ok 0 0
> exec
| CREATE TABLE c3 (pid INT, FOREIGN KEY (pid) REFERENCES p (w))
< ok 0 0
> exec
| CREATE TABLE c4 (pid INT, FOREIGN KEY (pid) REFERENCES p (x))
< ok 0 0
> exec
| CREATE TABLE c5 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (x, y))
< ok 0 0
> exec
| CREATE TABLE c6 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (y, x))
< ok 0 0
> exec
| CREATE TABLE c7 (pid INT)
< ok 0 0
> exec
| ALTER TABLE c7 ADD FOREIGN KEY (pid) REFERENCES p (v)
< ok 0 0
> exec
| ALTER TABLE c7 ADD FOREIGN KEY (pid) REFERENCES p (id)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c8 (pid INT, FOREIGN KEY (pid) REFERENCES p (y))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c9 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (v, w))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE q (a INT, b INT, c INT, PRIMARY KEY (a, b), KEY kc (c, a))
< ok 0 0
> exec
| CREATE TABLE d1 (x INT, FOREIGN KEY (x) REFERENCES q (b))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE d2 (x INT, FOREIGN KEY (x) REFERENCES q (a))
< ok 0 0
> exec
| CREATE TABLE d3 (x INT, y INT, FOREIGN KEY (x, y) REFERENCES q (a, c))
< ok 0 0
> exec
| CREATE TABLE d4 (x INT, y INT, z INT, FOREIGN KEY (x, y, z) REFERENCES q (a, b, c))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE s (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES s (id))
< ok 0 0
> exec
| CREATE TABLE s2 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES s2 (id))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS obkey_03
< ok 0 0
> exec
| CREATE DATABASE obkey_03
< ok 1 0
> exec
| USE obkey_03
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c1 (pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES p (id))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c2 (pid BIGINT, FOREIGN KEY (pid) REFERENCES p (id))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c3 (pid INT UNSIGNED, FOREIGN KEY (pid) REFERENCES p (id))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c4 (pid INT)
< ok 0 0
> exec
| ALTER TABLE c4 ADD FOREIGN KEY (pid) REFERENCES p (id)
< ok 0 0
> exec
| CREATE TABLE c5 (pid BIGINT)
< ok 0 0
> exec
| ALTER TABLE c5 ADD FOREIGN KEY (pid) REFERENCES p (id)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS obkey_04
< ok 0 0
> exec
| CREATE DATABASE obkey_04
< ok 1 0
> exec
| USE obkey_04
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c1 (pid INT DEFAULT 0, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET DEFAULT)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c2 (pid INT DEFAULT 0, FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET DEFAULT)
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c3 (pid INT DEFAULT 0)
< ok 0 0
> exec
| ALTER TABLE c3 ADD FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET DEFAULT
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| DROP DATABASE IF EXISTS obkey_05
< ok 0 0
> exec
| CREATE DATABASE obkey_05
< ok 1 0
> exec
| USE obkey_05
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c1 (pid INT NOT NULL, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET NULL)
< error
| 1830 HY000 Column 'pid' cannot be NOT NULL: needed in a foreign key constraint 'c1_OBFK_0000000000000001' SET NULL
> exec
| CREATE TABLE c2 (pid INT NOT NULL, FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET NULL)
< error
| 1830 HY000 Column 'pid' cannot be NOT NULL: needed in a foreign key constraint 'c2_OBFK_0000000000000002' SET NULL
> exec
| DROP DATABASE IF EXISTS obkey_06
< ok 0 0
> exec
| CREATE DATABASE obkey_06
< ok 1 0
> exec
| USE obkey_06
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY, v INT GENERATED ALWAYS AS (id + 1) VIRTUAL, s INT GENERATED ALWAYS AS (id + 1) STORED, UNIQUE KEY uv (v), UNIQUE KEY us (s))
< ok 0 0
> exec
| CREATE TABLE c1 (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL, FOREIGN KEY (b) REFERENCES p (id))
< ok 0 0
> exec
| CREATE TABLE c2 (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED, FOREIGN KEY (b) REFERENCES p (id))
< ok 0 0
> exec
| CREATE TABLE c3 (a INT, FOREIGN KEY (a) REFERENCES p (v))
< ok 0 0
> exec
| CREATE TABLE c4 (a INT, FOREIGN KEY (a) REFERENCES p (s))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obkey_07
< ok 0 0
> exec
| CREATE DATABASE obkey_07
< ok 1 0
> exec
| USE obkey_07
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c1 (a INT, b INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id), CONSTRAINT fk FOREIGN KEY (b) REFERENCES p (id))
< error
| 1215 HY000 Cannot add foreign key constraint
> exec
| CREATE TABLE c2 (a INT, b INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id))
< ok 0 0
> exec
| CREATE TABLE c3 (a INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id))
< error
| 1022 23000 Can't write; duplicate key in table 'c3'
> exec
| ALTER TABLE c2 ADD CONSTRAINT fk FOREIGN KEY (b) REFERENCES p (id)
< error
| 1022 23000 Can't write; duplicate key in table 'c2'
