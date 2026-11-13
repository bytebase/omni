-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestListCommasOceanBase default#1
> exec
| CREATE DATABASE IF NOT EXISTS oblistcomma
< ok 1 0
> exec
| USE oblistcomma
< ok 0 0
> exec
| DROP TABLE IF EXISTS pr
< ok 0 0
> exec
| CREATE TABLE pr (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (id INT)'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (id INT,)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id),)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t ()'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10),)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) ()'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0,))'
< error
| 1064 42000 …near '))' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) ())'
< error
| 1064 42000 …near '))' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITION TEMPLATE (SUBPARTITION s0) (PARTITION p0 VALUES LESS THAN (10))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITION TEMPLATE (SUBPARTITION s0,) (PARTITION p0 VALUES LESS THAN (10))'
< error
| 1064 42000 …near ') (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITION TEMPLATE () (PARTITION p0 VALUES LESS THAN (10))'
< error
| 1064 42000 …near ') (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr ADD PARTITION (PARTITION p2 VALUES LESS THAN (30),)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr ADD PARTITION ()'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20))'
< error
| 1235 0A000 reorganizing table without primary key(s) is not supported
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20),)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO ()'
< error
| 1064 42000 …near ')' at line 1
> exec
| DROP DATABASE IF EXISTS oblistcomma
< ok 1 0
