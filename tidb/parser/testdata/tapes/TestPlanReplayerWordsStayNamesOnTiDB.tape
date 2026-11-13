-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestPlanReplayerWordsStayNamesOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT * FROM replayer
< error
| 1146 42S02 Table 'test.replayer' doesn't exist
> exec ^
: string
| SELECT * FROM replayer
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM d.replayer
< error
| 1146 42S02 Table 'd.replayer' doesn't exist
> exec ^
: string
| SELECT * FROM d.replayer
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM dump
< error
| 1146 42S02 Table 'test.dump' doesn't exist
> exec ^
: string
| SELECT * FROM dump
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM capture
< error
| 1146 42S02 Table 'test.capture' doesn't exist
> exec ^
: string
| SELECT * FROM capture
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM plan
< error
| 1146 42S02 Table 'test.plan' doesn't exist
> exec ^
: string
| SELECT * FROM plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t replayer WHERE replayer.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t replayer WHERE replayer.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t dump WHERE dump.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t dump WHERE dump.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t capture WHERE capture.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t capture WHERE capture.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t plan WHERE plan.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t plan WHERE plan.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t replayer SET replayer.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t replayer SET replayer.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT dump.a FROM t dump
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT dump.a FROM t dump
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS replayer
< ok 0 0
> exec ^
: string
| SELECT 1 AS dump
< ok 0 0
> exec ^
: string
| SELECT 1 AS capture
< ok 0 0
> exec ^
: string
| SELECT replayer, dump, capture FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT replayer, dump, capture FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.plan FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.plan FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE replayer (dump int, capture int, plan int, stats int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (replayer int, INDEX dump (replayer))
< ok 0 0
> exec ^
: string
| CREATE INDEX replayer ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX replayer ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP TABLE replayer
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN replayer int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN replayer int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE plan RENAME TO replayer
< ok 0 0
> exec ^
: string
| PLAN REPLAYER DUMP EXPLAIN SELECT * FROM replayer
< error
| 1146 42S02 Table 'test.replayer' doesn't exist
> exec ^
: string
| PLAN REPLAYER DUMP EXPLAIN SELECT * FROM replayer
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| PLAN REPLAYER DUMP EXPLAIN SELECT dump.a FROM t dump
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| PLAN REPLAYER DUMP EXPLAIN SELECT dump.a FROM t dump
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| PLAN REPLAYER DUMP EXPLAIN DELETE FROM t plan WHERE plan.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| PLAN REPLAYER DUMP EXPLAIN DELETE FROM t plan WHERE plan.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
