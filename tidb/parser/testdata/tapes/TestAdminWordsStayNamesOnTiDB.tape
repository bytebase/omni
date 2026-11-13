-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAdminWordsStayNamesOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE jobs (queries int, bindings int, thread int, bdr int)
< ok 0 0
> exec ^
: string
| CREATE TABLE ddl (next_row_id int, plan_cache int, statistics int)
< ok 0 0
> exec ^
: string
| SELECT next_row_id FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT next_row_id FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS admin
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t AS ddl
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t AS ddl
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT ddl.id FROM t ddl
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT ddl.id FROM t ddl
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t ddl WHERE ddl.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t ddl WHERE ddl.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t queries WHERE queries.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t queries WHERE queries.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t jobs SET jobs.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t jobs SET jobs.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX queries ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX queries ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX next_row_id ON jobs (bdr)
< error
| 1146 42S02 Table 'test.jobs' doesn't exist
> exec ^
: string
| CREATE INDEX next_row_id ON jobs (bdr)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO bindings (capture, evolve) VALUES (1, 2)
< error
| 1146 42S02 Table 'test.bindings' doesn't exist
> exec ^
: string
| INSERT INTO bindings (capture, evolve) VALUES (1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t WHERE cleanup = 1 AND top = 2 AND recent = 3
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t WHERE cleanup = 1 AND top = 2 AND recent = 3
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT thread, batch_size, max_write_speed FROM jobs
< error
| 1146 42S02 Table 'test.jobs' doesn't exist
> exec ^
: string
| SELECT thread, batch_size, max_write_speed FROM jobs
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
