-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestCastUnknownTypesOracle default#1
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
| PREPARE parse_probe FROM 'SELECT CAST(a AS INTERVAL DAY)'
< error
| 4161 HY000 Unknown data type: 'INTERVAL'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS DAY)'
< error
| 4161 HY000 Unknown data type: 'DAY'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS GEOMCOLLECTION)'
< error
| 4161 HY000 Unknown data type: 'GEOMCOLLECTION'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS `CHAR`)'
< error
| 4161 HY000 Unknown data type: 'CHAR'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS `SIGNED`)'
< error
| 4161 HY000 Unknown data type: 'SIGNED'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS `DATE`)'
< error
| 4161 HY000 Unknown data type: 'DATE'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS NULL)'
< error
| 4161 HY000 Unknown data type: 'NULL'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS ROARINGBITMAP)'
< error
| 4161 HY000 Unknown data type: 'ROARINGBITMAP'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS ARRAY(INT))'
< error
| 4161 HY000 Unknown data type: 'ARRAY'
> exec
| PREPARE parse_probe FROM 'SELECT CAST(a AS ARRAY<INT>)'
< error
| 4161 HY000 Unknown data type: 'ARRAY'
