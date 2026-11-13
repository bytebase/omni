-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestPlanReplayerFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| plan replayer dump explain select 1
< ok 0 0
> exec ^
: string
| plan replayer dump explain analyze select 1
< ok 0 0
> exec ^
: string
| plan replayer dump explain select a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain select a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain select * from t where a > 10
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain select * from t where a > 10
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain values row(1)
< ok 0 0
> exec ^
: string
| plan replayer dump explain table t1 union table t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| plan replayer dump explain table t1 union table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain values row(1) union values row(2)
< ok 0 0
> exec ^
: string
| plan replayer dump explain table t1 union select 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| plan replayer dump explain table t1 union select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain analyze table t1 union table t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| plan replayer dump explain analyze table t1 union table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| plan replayer dump explain (select 1)
< ok 0 0
> exec ^
: string
| plan replayer dump explain select 1 union select 2
< ok 0 0
> exec ^
: string
| plan replayer dump explain insert into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain insert into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain replace into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain replace into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain update t set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain update t set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain delete from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain delete from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain alter table t add column c int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain alter table t add column c int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain analyze alter table t add column c int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain analyze alter table t add column c int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain import into t from '/tmp/a.csv'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain import into t from '/tmp/a.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain analyze import into t from '/tmp/a.csv'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain analyze import into t from '/tmp/a.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain import into t from table t1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain import into t from table t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain begin
< error
| 1064 42000 …line 1 column 32 near "begin" 
> exec ^
: string
| plan replayer dump explain set a = 1
< error
| 1064 42000 …line 1 column 30 near "set a = 1" 
> exec ^
: string
| plan replayer dump explain load data infile '/tmp/a.csv' into table t
< error
| 1064 42000 …line 1 column 31 near "load data infile '/tmp/a.csv' into table t" 
> exec ^
: string
| plan replayer dump explain create table t2 (a int)
< error
| 1064 42000 …line 1 column 33 near "create table t2 (a int)" 
> exec ^
: string
| plan replayer dump explain truncate table t
< error
| 1064 42000 …line 1 column 35 near "truncate table t" 
> exec ^
: string
| plan replayer dump explain alter database d charset utf8mb4
< error
| 1064 42000 …line 1 column 41 near "database d charset utf8mb4" 
> exec ^
: string
| plan replayer dump explain drop table t
< error
| 1064 42000 …line 1 column 31 near "drop table t" 
> exec ^
: string
| plan replayer dump explain show tables
< error
| 1064 42000 …line 1 column 31 near "show tables" 
> exec ^
: string
| plan replayer dump explain explain select 1
< error
| 1064 42000 …line 1 column 34 near "explain select 1" 
> exec ^
: string
| plan replayer dump explain t
< error
| 1064 42000 …line 1 column 28 near "t" 
> exec ^
: string
| plan replayer dump explain analyze table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| plan replayer dump explain analyze table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump explain slow query
< ok 0 0
> exec ^
: string
| plan replayer dump explain analyze slow query
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query limit 10
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query limit 1, 10
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query where a > 10 limit 10 offset 5
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query where a > 10 and t < 1 limit 10
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query where a > 10 order by t
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query order by t
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow query where a > 10 and t < 1 order by t limit 10
< ok 0 0
> exec ^
: string
| plan replayer dump explain analyze slow query where a > 10 and t < 1 order by t limit 10
< ok 0 0
> exec ^
: string
| plan replayer dump explain slow
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| plan replayer dump explain slow queries
< error
| 1064 42000 …line 1 column 39 near "queries" 
> exec ^
: string
| plan replayer dump explain slow query limit 10 where a > 10
< error
| 1064 42000 …line 1 column 52 near "where a > 10" 
> exec ^
: string
| plan replayer dump explain slow query order by t where a > 10
< error
| 1064 42000 …line 1 column 54 near "where a > 10" 
> exec ^
: string
| plan replayer dump explain slow query where
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| plan replayer dump explain 'sql.txt'
< ok 0 0
> exec ^
: string
| plan replayer dump explain analyze 'sql.txt'
< ok 0 0
> exec ^
: string
| plan replayer dump explain `sql.txt`
< error
| 1064 42000 …line 1 column 36 near "`sql.txt`" 
> exec ^
: string
| plan replayer dump explain analyze 'a' 'b'
< error
| 1064 42000 …line 1 column 42 near "'b'" 
> exec ^
: string
| plan replayer dump with stats as of timestamp '2024-01-01 00:00:00' explain select 1
< ok 0 0
> exec ^
: string
| plan replayer dump with stats as of timestamp now() explain analyze select 1
< ok 0 0
> exec ^
: string
| plan replayer dump with stats as of timestamp '2024-01-01' explain slow query limit 10
< ok 0 0
> exec ^
: string
| plan replayer dump with stats as of timestamp '2024-01-01' explain 'sql.txt'
< ok 0 0
> exec ^
: string
| plan replayer dump with stats explain select 1
< error
| 1064 42000 …line 1 column 37 near "explain select 1" 
> exec ^
: string
| plan replayer dump with stats as of explain select 1
< error
| 1064 42000 …line 1 column 43 near "explain select 1" 
> exec ^
: string
| plan replayer dump with stats as of timestamp explain select 1
< error
| 1064 42000 …line 1 column 53 near "explain select 1" 
> exec ^
: string
| plan replayer dump with explain select 1
< error
| 1064 42000 …line 1 column 31 near "explain select 1" 
> exec ^
: string
| plan replayer dump as of timestamp '2024-01-01' explain select 1
< error
| 1064 42000 …line 1 column 24 near "of timestamp '2024-01-01' explain select 1" 
> exec ^
: string
| plan replayer with stats as of timestamp '2024-01-01' dump explain select 1
< error
| 1064 42000 …line 1 column 18 near "with stats as of timestamp '2024-01-01' dump explain select 1" 
> exec ^
: string
| plan replayer load '/tmp/sdfaalskdjf.zip'
< ok 0 0
> exec ^
: string
| plan replayer load
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| plan replayer load 'a.zip' 'b.zip'
< error
| 1064 42000 …line 1 column 34 near "'b.zip'" 
> exec ^
: string
| plan replayer load `a.zip`
< error
| 1064 42000 …line 1 column 26 near "`a.zip`" 
> exec ^
: string
| plan replayer capture '123' '123'
< ok 0 0
> exec ^
: string
| plan replayer capture remove '123' '123'
< ok 0 0
> exec ^
: string
| plan replayer capture '123'
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| plan replayer capture
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| plan replayer capture '123' '123' '123'
< error
| 1064 42000 …line 1 column 39 near "'123'" 
> exec ^
: string
| plan replayer capture remove '123'
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| plan replayer capture remove
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| plan replayer capture 123 123
< error
| 1064 42000 …line 1 column 25 near "123 123" 
> exec ^
: string
| plan replayer
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| plan replayer dump
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| plan replayer dump explain
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| plan replayer dump select 1
< error
| 1064 42000 …line 1 column 25 near "select 1" 
> exec ^
: string
| plan replayer explain select 1
< error
| 1064 42000 …line 1 column 21 near "explain select 1" 
> exec ^
: string
| plan replayer dumping explain select 1
< error
| 1064 42000 …line 1 column 21 near "dumping explain select 1" 
> exec ^
: string
| plan replayers dump explain select 1
< error
| 1064 42000 …line 1 column 14 near "replayers dump explain select 1" 
> exec ^
: string
| plan dump explain select 1
< error
| 1064 42000 …line 1 column 9 near "dump explain select 1" 
> exec ^
: string
| plan
< error
| 1064 42000 …line 1 column 4 near "" 
> exec ^
: string
| replayer dump explain select 1
< error
| 1064 42000 …line 1 column 8 near "replayer dump explain select 1" 
> exec ^
: string
| plan `replayer` dump explain select 1
< error
| 1064 42000 …line 1 column 15 near "`replayer` dump explain select 1" 
> exec ^
: string
| `plan` replayer dump explain select 1
< error
| 1064 42000 …line 1 column 6 near "`plan` replayer dump explain select 1" 
> exec ^
: string
| plan replayer `dump` explain select 1
< error
| 1064 42000 …line 1 column 20 near "`dump` explain select 1" 
> exec ^
: string
| plan replayer `capture` '123' '123'
< error
| 1064 42000 …line 1 column 23 near "`capture` '123' '123'" 
