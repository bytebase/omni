-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBAlterDropFormsOracle default#1
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT, pid INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))
< ok 0 0
> exec
| CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP n RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP pid CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN IF EXISTS n RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP IF EXISTS no_such CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n RESTRICT, DROP COLUMN pid CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n CASCADE, ADD COLUMN z INT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n /* c */ RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN `restrict`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP `cascade`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n RESTRICT CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n RESTRICT RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n RESTRICT FIRST'
< error
| 1064 42000 …near 'FIRST' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN n, RESTRICT'
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
| PREPARE parse_probe FROM 'ALTER TABLE t DROP KEY IF EXISTS i CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY f RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY IF EXISTS f RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT ck_n CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT PRIMARY KEY RESTRICT'
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
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)'
< error
| 1064 42000 …near ' ADD PRIMARY KEY (id)' at line 1
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
