-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestAlterTableDropFormsOracle default#1
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT, m INT, k INT, pid INT, `restrict` INT, `cascade` INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> exec
| CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP m RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP k CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT, DROP COLUMN k CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m CASCADE, ADD COLUMN z INT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m /* c */ RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN `restrict`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP `cascade`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT FIRST'
< error
| 1064 42000 …near 'FIRST' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m, RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN restrict'
< error
| 1064 42000 …near 'restrict' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP cascade'
< error
| 1064 42000 …near 'cascade' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX i RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP KEY i CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY f RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT ck_n CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CHECK ck_n RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY KEY CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pt DROP PARTITION p0 RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY KEY'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY f'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY f'
< error
| 1235 0A000 Multiple complex DDLs about primary in single stmt not supported
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY'
< error
| 1064 42000 …near 'PRIMARY' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)'
< error
| 1064 42000 …near ', ADD PRIMARY KEY (id)' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY INDEX'
< error
| 1064 42000 …near 'INDEX' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN f'
< error
| 1064 42000 …near 'f' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN INDEX f'
< error
| 1064 42000 …near 'INDEX f' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN IF EXISTS f'
< error
| 1064 42000 …near 'IF EXISTS f' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN IF EXISTS m'
< error
| 1064 42000 …near 'IF EXISTS m' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP IF EXISTS m'
< error
| 1064 42000 …near 'IF EXISTS m' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN IF EXISTS m RESTRICT'
< error
| 1064 42000 …near 'IF EXISTS m RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX IF EXISTS i'
< error
| 1064 42000 …near 'IF EXISTS i' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY IF EXISTS f'
< error
| 1064 42000 …near 'IF EXISTS f' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n'
< error
| 1064 42000 …near 'IF EXISTS ck_n' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CHECK IF EXISTS ck_n'
< error
| 1064 42000 …near 'IF EXISTS ck_n' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pt DROP PARTITION IF EXISTS p0'
< error
| 1064 42000 …near 'IF EXISTS p0' at line 1
