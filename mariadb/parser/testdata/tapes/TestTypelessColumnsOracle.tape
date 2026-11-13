-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestTypelessColumnsOracle default#1
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
| PREPARE parse_probe FROM 'CREATE TABLE tl (a key)'
< error
| 4161 HY000 Unknown data type: 'key'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a primary key)'
< error
| 4161 HY000 Unknown data type: 'primary'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a unique)'
< error
| 4161 HY000 Unknown data type: 'unique'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a unique key)'
< error
| 4161 HY000 Unknown data type: 'unique'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a comment ''x'')'
< error
| 4161 HY000 Unknown data type: 'comment'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a not null)'
< error
| 4161 HY000 Unknown data type: 'not'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a null)'
< error
| 4161 HY000 Unknown data type: 'null'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a default 1)'
< error
| 4161 HY000 Unknown data type: 'default'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a auto_increment)'
< error
| 4161 HY000 Unknown data type: 'auto_increment'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a collate utf8mb4_bin)'
< error
| 4161 HY000 Unknown data type: 'collate'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a references r (id))'
< error
| 4161 HY000 Unknown data type: 'references'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a generated always as (1))'
< error
| 4161 HY000 Unknown data type: 'generated'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a generated always as (1) stored)'
< error
| 4161 HY000 Unknown data type: 'generated'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a generated always as (1) virtual)'
< error
| 4161 HY000 Unknown data type: 'generated'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a as (1))'
< error
| 4161 HY000 Unknown data type: 'as'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a as (1) stored)'
< error
| 4161 HY000 Unknown data type: 'as'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a on update now())'
< error
| 4161 HY000 Unknown data type: 'on'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a column_format fixed)'
< error
| 4161 HY000 Unknown data type: 'column_format'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a storage disk)'
< error
| 4161 HY000 Unknown data type: 'storage'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a srid 0)'
< error
| 4161 HY000 Unknown data type: 'srid'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a invisible)'
< error
| 4161 HY000 Unknown data type: 'invisible'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a visible)'
< error
| 4161 HY000 Unknown data type: 'visible'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a auto_random)'
< error
| 4161 HY000 Unknown data type: 'auto_random'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a key) SELECT 1 AS a'
< error
| 4161 HY000 Unknown data type: 'key'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a default 1) SELECT 1 AS a'
< error
| 4161 HY000 Unknown data type: 'default'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a not null) AS SELECT 1 AS a'
< error
| 4161 HY000 Unknown data type: 'not'
> exec
| PREPARE parse_probe FROM 'CREATE TABLE tl (a comment ''x'') AS SELECT 1 AS a'
< error
| 4161 HY000 Unknown data type: 'comment'
> exec
| PREPARE parse_probe FROM 'CREATE TEMPORARY TABLE tl (a key)'
< error
| 4161 HY000 Unknown data type: 'key'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t ADD a key'
< error
| 4161 HY000 Unknown data type: 'key'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t ADD COLUMN a comment ''x'''
< error
| 4161 HY000 Unknown data type: 'comment'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t ADD COLUMN a generated always as (1)'
< error
| 4161 HY000 Unknown data type: 'generated'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t ADD (a key)'
< error
| 4161 HY000 Unknown data type: 'key'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t ADD COLUMN (a not null, b int)'
< error
| 4161 HY000 Unknown data type: 'not'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t MODIFY a comment ''x'''
< error
| 4161 HY000 Unknown data type: 'comment'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t MODIFY COLUMN a not null'
< error
| 4161 HY000 Unknown data type: 'not'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t MODIFY a generated always as (1)'
< error
| 4161 HY000 Unknown data type: 'generated'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t CHANGE a b key'
< error
| 4161 HY000 Unknown data type: 'key'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t CHANGE COLUMN a b comment ''x'''
< error
| 4161 HY000 Unknown data type: 'comment'
> exec
| PREPARE parse_probe FROM 'ALTER TABLE t CHANGE a b generated always as (1)'
< error
| 4161 HY000 Unknown data type: 'generated'
