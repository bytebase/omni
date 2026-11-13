-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBDropClauseOracle default#1
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
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS no_such'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX IF EXISTS no_such'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP KEY IF EXISTS i'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY IF EXISTS f'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN IF EXISTS no_such'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP IF EXISTS no_such'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX IF EXISTS `if`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS c1, DROP INDEX IF EXISTS i1, DROP FOREIGN KEY IF EXISTS f1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pt DROP PARTITION IF EXISTS no_such'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pt DROP PARTITION IF EXISTS p0, p1'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT ck_n'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT PRIMARY KEY'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CHECK ck_n'
< error
| 1064 42000 …near 'CHECK ck_n' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CHECK IF EXISTS ck_n'
< error
| 1064 42000 …near 'CHECK IF EXISTS ck_n' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX IF no_such'
< error
| 1064 42000 …near 'no_such' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX EXISTS no_such'
< error
| 1064 42000 …near 'EXISTS no_such' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT PRIMARY'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS PRIMARY KEY'
< error
| 1064 42000 …near 'PRIMARY KEY' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY KEY IF EXISTS'
< error
| 1064 42000 …near 'IF EXISTS' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pt DROP PARTITION IF EXISTS'
< error
| 1064 42000 …near '' at line 1
