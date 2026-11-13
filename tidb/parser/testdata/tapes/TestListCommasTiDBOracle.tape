-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestListCommasTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_list_comma
< ok 0 0
> exec
| CREATE DATABASE omni_list_comma
< ok 0 0
> exec
| USE omni_list_comma
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
| 1064 42000 …line 1 column 24 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (id INT, PRIMARY KEY (id),)'
< error
| 1064 42000 …line 1 column 42 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t ()'
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10),)'
< error
| 1064 42000 …line 1 column 83 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) ()'
< error
| 1064 42000 …line 1 column 48 near ")" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) (SUBPARTITION s0,))'
< error
| 1064 42000 …line 1 column 125 near "))" 
> exec
| PREPARE parse_probe FROM 'CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) (PARTITION p0 VALUES LESS THAN (10) ())'
< error
| 1064 42000 …line 1 column 109 near "))" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr ADD PARTITION (PARTITION p2 VALUES LESS THAN (30))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr ADD PARTITION (PARTITION p2 VALUES LESS THAN (30),)'
< error
| 1064 42000 …line 1 column 66 near ")" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr ADD PARTITION ()'
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20))'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (15), PARTITION p1b VALUES LESS THAN (20),)'
< error
| 1064 42000 …line 1 column 119 near ")" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE pr REORGANIZE PARTITION p1 INTO ()'
< error
| 1064 42000 …line 1 column 46 near ")" 
> exec
| DROP DATABASE IF EXISTS omni_list_comma
< ok 0 0
