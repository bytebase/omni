-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBDropPeriodOracle default#1
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
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR app'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR `app`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR period'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR SYSTEM_TIME_X'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR `if`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD /* c */ IF /* c */ EXISTS /* c */ FOR app'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS /*!100000 FOR app */'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR app, DROP PERIOD IF EXISTS FOR app'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR SYSTEM_TIME'
< error
| 1064 42000 …near 'SYSTEM_TIME' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR system_time'
< error
| 1064 42000 …near 'system_time' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR /* c */ SYSTEM_TIME'
< error
| 1064 42000 …near 'SYSTEM_TIME' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS /*!FOR SYSTEM_TIME */'
< error
| 1064 42000 …near 'SYSTEM_TIME */' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR "SYSTEM_TIME"'
< error
| 1064 42000 …near '"SYSTEM_TIME"' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD FOR app'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD FOR SYSTEM_TIME'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD FOR `SYSTEM_TIME`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS app'
< error
| 1064 42000 …near 'app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF FOR app'
< error
| 1064 42000 …near 'FOR app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD EXISTS FOR app'
< error
| 1064 42000 …near 'EXISTS FOR app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF NOT EXISTS FOR app'
< error
| 1064 42000 …near 'NOT EXISTS FOR app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR if'
< error
| 1064 42000 …near 'if' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR FOR'
< error
| 1064 42000 …near 'FOR' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR t.app'
< error
| 1064 42000 …near '.app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR app RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD IF EXISTS FOR app CASCADE'
< error
| 1064 42000 …near 'CASCADE' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD FOR app RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN PERIOD IF EXISTS FOR app'
< error
| 1064 42000 …near 'IF EXISTS FOR app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP period'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP PERIOD'
< error
| 1064 42000 …near '' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP period, DROP id'
< error
| 1064 42000 …near ' DROP id' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP period RESTRICT'
< error
| 1064 42000 …near 'RESTRICT' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP `period`'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP COLUMN period'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP IF EXISTS period'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP `PERIOD` FOR app'
< error
| 1064 42000 …near 'FOR app' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t DROP `PERIOD` IF EXISTS FOR app'
< error
| 1064 42000 …near 'IF EXISTS FOR app' at line 1
