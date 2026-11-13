-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAlterTableDropFormsTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_alter_drop
< ok 0 0
> exec
| CREATE DATABASE omni_alter_drop
< ok 0 0
> exec
| USE omni_alter_drop
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, n INT, m INT, k INT, pid INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))
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
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN IF EXISTS m RESTRICT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP IF EXISTS m CASCADE'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN IF EXISTS m'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP IF EXISTS no_such'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT CASCADE'
< error
| 1064 42000 …line 1 column 44 near "CASCADE" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT RESTRICT'
< error
| 1064 42000 …line 1 column 45 near "RESTRICT" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m RESTRICT FIRST'
< error
| 1064 42000 …line 1 column 42 near "FIRST" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN m, RESTRICT'
< error
| 1064 42000 …line 1 column 37 near "RESTRICT" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN restrict'
< error
| 1064 42000 …line 1 column 34 near "restrict" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP cascade'
< error
| 1064 42000 …line 1 column 26 near "cascade" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP INDEX i RESTRICT'
< error
| 1064 42000 …line 1 column 35 near "RESTRICT" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP KEY i CASCADE'
< error
| 1064 42000 …line 1 column 32 near "CASCADE" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP KEY IF EXISTS i CASCADE'
< error
| 1064 42000 …line 1 column 42 near "CASCADE" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY f RESTRICT'
< error
| 1064 42000 …line 1 column 41 near "RESTRICT" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT ck_n CASCADE'
< error
| 1064 42000 …line 1 column 42 near "CASCADE" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CHECK ck_n RESTRICT'
< error
| 1064 42000 …line 1 column 38 near "RESTRICT" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY KEY CASCADE'
< error
| 1064 42000 …line 1 column 38 near "CASCADE" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pt DROP PARTITION p0 RESTRICT'
< error
| 1064 42000 …line 1 column 41 near "RESTRICT" 
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
| 1064 42000 …line 1 column 26 near "" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)'
< error
| 1064 42000 …line 1 column 27 near ", ADD PRIMARY KEY (id)" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PRIMARY INDEX'
< error
| 1064 42000 …line 1 column 32 near "INDEX" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN f'
< error
| 1064 42000 …line 1 column 28 near "f" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN INDEX f'
< error
| 1064 42000 …line 1 column 32 near "INDEX f" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN IF EXISTS f'
< error
| 1064 42000 …line 1 column 29 near "IF EXISTS f" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP FOREIGN KEY IF EXISTS f'
< error
| 1064 42000 …line 1 column 33 near "IF EXISTS f" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n'
< error
| 1064 42000 …line 1 column 32 near "IF EXISTS ck_n" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP CHECK IF EXISTS ck_n'
< error
| 1064 42000 …line 1 column 27 near "IF EXISTS ck_n" 
> exec
| DROP DATABASE IF EXISTS omni_alter_drop
< ok 0 0
