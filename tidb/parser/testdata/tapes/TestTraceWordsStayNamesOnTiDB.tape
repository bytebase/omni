-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTraceWordsStayNamesOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT * FROM trace
< error
| 1146 42S02 Table 'test.trace' doesn't exist
> exec ^
: string
| SELECT * FROM trace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM d.trace
< error
| 1146 42S02 Table 'd.trace' doesn't exist
> exec ^
: string
| SELECT * FROM d.trace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM target
< error
| 1146 42S02 Table 'test.target' doesn't exist
> exec ^
: string
| SELECT * FROM target
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t trace WHERE trace.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t trace WHERE trace.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t target WHERE target.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t target WHERE target.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t trace SET trace.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t trace SET trace.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT trace.a FROM t trace
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT trace.a FROM t trace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS trace
< ok 0 0
> exec ^
: string
| SELECT 1 AS target
< ok 0 0
> exec ^
: string
| SELECT trace, target FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT trace, target FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT t.trace FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT t.trace FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE trace (target int, trace int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (trace int, target int, INDEX trace (trace))
< ok 0 0
> exec ^
: string
| CREATE INDEX trace ON t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX trace ON t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP TABLE trace
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN trace int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN trace int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE trace RENAME TO target
< ok 0 0
> exec ^
: string
| TRACE SELECT * FROM trace
< error
| 1146 42S02 Table 'test.trace' doesn't exist
> exec ^
: string
| TRACE SELECT * FROM trace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRACE SELECT trace.a FROM t trace
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TRACE SELECT trace.a FROM t trace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRACE DELETE FROM t trace WHERE trace.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| TRACE DELETE FROM t trace WHERE trace.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| TRACE PLAN SELECT target FROM trace
< error
| 1146 42S02 Table 'test.trace' doesn't exist
> exec ^
: string
| TRACE PLAN SELECT target FROM trace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
