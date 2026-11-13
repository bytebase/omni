-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAlterIgnoreTableTiDBOracle default#1
> exec
| DROP DATABASE IF EXISTS omni_alter_ignore
< ok 0 0
> exec
| CREATE DATABASE omni_alter_ignore
< ok 0 0
> exec
| USE omni_alter_ignore
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE TABLE t'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE TABLE t ADD COLUMN z INT'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE TABLE t RENAME TO t2'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE TABLE t ANALYZE PARTITION p0'
< ok 0 0
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE DATABASE d'
< error
| 1064 42000 …line 1 column 21 near "DATABASE d" 
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE SEQUENCE s'
< error
| 1064 42000 …line 1 column 21 near "SEQUENCE s" 
> exec
| PREPARE parse_probe FROM 'ALTER TABLE IGNORE t'
< error
| 1064 42000 …line 1 column 18 near "IGNORE t" 
> exec
| PREPARE parse_probe FROM 'ALTER IGNORE'
< error
| 1064 42000 …line 1 column 12 near "" 
> exec
| DROP DATABASE IF EXISTS omni_alter_ignore
< ok 0 0
