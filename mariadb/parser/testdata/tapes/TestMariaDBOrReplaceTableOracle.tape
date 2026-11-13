-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBOrReplaceTableOracle default#1
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
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY TABLE r (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r LIKE t'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r AS SELECT 1 AS a'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE r (id INT) SELECT 1 AS id'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE IF NOT EXISTS r (id INT)'
< error
| 1221 HY000 Incorrect usage of OR REPLACE and IF NOT EXISTS
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TEMPORARY TABLE IF NOT EXISTS r (id INT)'
< error
| 1221 HY000 Incorrect usage of OR REPLACE and IF NOT EXISTS
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE IF NOT EXISTS r LIKE t'
< error
| 1221 HY000 Incorrect usage of OR REPLACE and IF NOT EXISTS
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE VECTOR INDEX IF NOT EXISTS i ON t (v)'
< error
| 1221 HY000 Incorrect usage of OR REPLACE and IF NOT EXISTS
> exec
| PREPARE parse_probe FROM 'CREATE OR TABLE r (id INT)'
< error
| 1064 42000 …near 'TABLE r (id INT)' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE OR REPLACE TABLE IF EXISTS r (id INT)'
< error
| 1064 42000 …near 'EXISTS r (id INT)' at line 1
