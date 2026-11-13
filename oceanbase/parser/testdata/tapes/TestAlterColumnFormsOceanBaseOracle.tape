-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestAlterColumnFormsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obaltcol
< ok 0 0
> exec
| CREATE DATABASE obaltcol
< ok 1 0
> exec
| USE obaltcol
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| ALTER TABLE t ALTER b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b DROP NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| ALTER TABLE t ALTER b DROP NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| alter table t alter column b set not null
< error
| 1064 42000 …near 'not null' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b DROP NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b DROP NOT
< error
| 1064 42000 …near 'NOT' at line 1
> exec
| ALTER TABLE t ALTER COLUMN t.b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| ALTER TABLE t ALTER COLUMN t.b DROP NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET NOT NULL
< error
| 1064 42000 …near 'NOT NULL' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b DROP NOT NULL, ALTER COLUMN b DROP DEFAULT
< error
| 1064 42000 …near 'NOT NULL, ALTER COLUMN b DROP DEFAULT' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET
< error
| 1064 42000 …near 'SET' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b DROP
< error
| 1064 42000 …near 'DROP' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b
< error
| 1064 42000 …near 'b' at line 1
> exec
| ALTER TABLE t ALTER b
< error
| 1064 42000 …near 'b' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE t ALTER b DROP DEFAULT
< ok 0 0
> exec
| ALTER TABLE t ALTER b DROP DEFAULT, ALTER b SET DEFAULT 2
< error
| 1054 42S22 Unknown column 'b' in 't'
> exec
| ALTER TABLE t MODIFY b INT NOT NULL
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN b INT NULL
< ok 0 0
> exec
| ALTER TABLE t CHANGE b b INT NOT NULL
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN b SET VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec
| ALTER TABLE t ALTER b SET VISIBLE
< error
| 1064 42000 …near 'VISIBLE' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec
| ALTER TABLE t ALTER b SET INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec
| alter table t alter column b set invisible
< error
| 1064 42000 …near 'invisible' at line 1
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET INVISIBLE
< error
| 1064 42000 …near 'INVISIBLE' at line 1
> exec
| CREATE TABLE t2 (a INT, b INT INVISIBLE)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT VISIBLE)
< error
| 1050 42S01 Table 't2' already exists
> exec
| ALTER TABLE t ADD COLUMN c INT INVISIBLE
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN c INT VISIBLE
< error
| 1060 42S21 Duplicate column name 'c'
> exec
| ALTER TABLE t MODIFY b INT INVISIBLE
< ok 0 0
> exec
| ALTER TABLE t MODIFY b INT VISIBLE
< ok 0 0
> exec
| ALTER TABLE t CHANGE b b INT INVISIBLE
< ok 0 0
> exec
| ALTER TABLE t CHANGE b b INT VISIBLE
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obaltcol
< ok 2 0
