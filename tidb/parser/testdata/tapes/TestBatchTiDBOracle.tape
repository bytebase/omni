-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBatchTiDBOracle default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_batch_test
< ok 0 0
> exec
| USE omni_batch_test
< ok 0 0
> exec
| BATCH ON id LIMIT 5000 DELETE FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t' doesn't exist
> exec
| BATCH LIMIT 5000 DELETE FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t' doesn't exist
> exec
| BATCH ON id LIMIT 1000 DRY RUN DELETE FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t' doesn't exist
> exec
| BATCH ON id LIMIT 1000 DRY RUN QUERY DELETE FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t' doesn't exist
> exec
| BATCH ON id LIMIT 100 UPDATE t SET x=1 WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t' doesn't exist
> exec
| BATCH ON id LIMIT 100 INSERT INTO t2 SELECT * FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t2' doesn't exist
> exec
| BATCH ON id LIMIT 100 INSERT IGNORE INTO t2 SELECT * FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t2' doesn't exist
> exec
| BATCH ON id LIMIT 100 REPLACE INTO t2 SELECT * FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t2' doesn't exist
> exec
| BATCH ON t.id LIMIT 100 DELETE FROM t WHERE 1=1
< error
| 1146 42S02 Table 'omni_batch_test.t' doesn't exist
> exec
| BATCH DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 12 near "DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON LIMIT 100 DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 14 near "LIMIT 100 DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON id LIMIT DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 24 near "DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON id LIMIT (1+1) DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 19 near "(1+1) DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON id LIMIT -100 DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 19 near "-100 DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON id LIMIT 100 SELECT * FROM t
< error
| 1064 42000 …line 1 column 28 near "SELECT * FROM t" 
> exec
| BATCH ON id LIMIT 100 CREATE TABLE x (id INT)
< error
| 1064 42000 …line 1 column 28 near "CREATE TABLE x (id INT)" 
> exec
| BATCH ON id LIMIT 100 DRY DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 32 near "DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON id LIMIT 100 DRY RUN FOO DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 33 near "FOO DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON id LIMIT 100
< error
| 1064 42000 …line 1 column 21 near "" 
> exec
| BATCH
< error
| 1064 42000 …line 1 column 5 near "" 
> exec
| BATCH ON * LIMIT 100 DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 10 near "* LIMIT 100 DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON t.* LIMIT 100 DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 12 near "* LIMIT 100 DELETE FROM t WHERE 1=1" 
> exec
| BATCH ON db.t.* LIMIT 100 DELETE FROM t WHERE 1=1
< error
| 1064 42000 …line 1 column 15 near "* LIMIT 100 DELETE FROM t WHERE 1=1" 
> exec
| DROP DATABASE IF EXISTS omni_batch_test
< ok 0 0
