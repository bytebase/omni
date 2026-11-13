-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseForeignKeyIndex default#1
> exec
| DROP DATABASE IF EXISTS obfk_00
< ok 0 0
> exec
| CREATE DATABASE obfk_00
< ok 1 0
> exec
| USE obfk_00
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, x INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id), FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> exec
| ALTER TABLE c ADD CONSTRAINT fkx FOREIGN KEY (x) REFERENCES p (id)
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX fk
< error
| 1091 42000 Can't DROP 'fk'; check that column/key exists
> exec
| ALTER TABLE c DROP INDEX pid
< error
| 1091 42000 Can't DROP 'pid'; check that column/key exists
> exec
| ALTER TABLE c DROP INDEX fkx
< error
| 1091 42000 Can't DROP 'fkx'; check that column/key exists
> exec
| ALTER TABLE c DROP FOREIGN KEY fk, DROP INDEX fk
< error
| 1091 42000 Can't DROP 'fk'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS obfk_01
< ok 0 0
> exec
| CREATE DATABASE obfk_01
< ok 1 0
> exec
| USE obfk_01
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, x INT, KEY a1 (pid), KEY a2 (pid), KEY wide (pid, x), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX wide
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX a2
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX a1
< error
| 5320 HY000 'a1': used in a foreign key constraint
> exec
| ALTER TABLE c DROP FOREIGN KEY fk
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX a1
< ok 0 0
> exec
| ALTER TABLE c ADD CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obfk_02
< ok 0 0
> exec
| CREATE DATABASE obfk_02
< ok 1 0
> exec
| USE obfk_02
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> exec
| ALTER TABLE c ADD UNIQUE KEY late (pid)
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX late
< error
| 5320 HY000 'late': used in a foreign key constraint
> exec
| DROP DATABASE IF EXISTS obfk_03
< ok 0 0
> exec
| CREATE DATABASE obfk_03
< ok 1 0
> exec
| USE obfk_03
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, a INT, b INT, s VARCHAR(20), KEY ba (b, a), KEY pre (s(5)), CONSTRAINT fab FOREIGN KEY (a, b) REFERENCES p (u, v))
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX pre
< ok 0 0
> exec
| ALTER TABLE c DROP INDEX ba
< error
| 5320 HY000 'ba': used in a foreign key constraint
> exec
| DROP DATABASE IF EXISTS obfk_04
< ok 0 0
> exec
| CREATE DATABASE obfk_04
< ok 1 0
> exec
| USE obfk_04
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))
< ok 0 0
> exec
| CREATE TABLE c (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id), CONSTRAINT fk FOREIGN KEY (id) REFERENCES p (id))
< ok 0 0
> exec
| ALTER TABLE c DROP PRIMARY KEY
< error
| 4714 HY000 unique/primary keys in table referenced by enabled foreign keys
> exec
| ALTER TABLE c DROP PRIMARY KEY, ADD PRIMARY KEY (a)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obfk_05
< ok 0 0
> exec
| CREATE DATABASE obfk_05
< ok 1 0
> exec
| USE obfk_05
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pu INT, CONSTRAINT fk FOREIGN KEY (pu) REFERENCES p (u))
< ok 0 0
> exec
| ALTER TABLE p DROP PRIMARY KEY
< error
| 4714 HY000 unique/primary keys in table referenced by enabled foreign keys
> exec
| ALTER TABLE p DROP INDEX uu
< error
| 4714 HY000 unique/primary keys in table referenced by enabled foreign keys
> exec
| ALTER TABLE p DROP PRIMARY KEY, ADD PRIMARY KEY (id, v)
< ok 0 0
