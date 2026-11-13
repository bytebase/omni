-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestListCommasOracle default#1
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
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20),)'
< error
| 1064 42000 …near ')' at line 1
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO ()'
< error
| 1064 42000 …near ')' at line 1
