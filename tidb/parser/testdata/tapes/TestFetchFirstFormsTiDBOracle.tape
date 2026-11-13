-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestFetchFirstFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select * from t fetch first 5 rows only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t fetch first 5 rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t fetch next 5 row only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t fetch next 5 row only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t fetch first row only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t fetch first row only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t fetch next rows only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t fetch next rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t order by a fetch first 2 rows only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t order by a fetch first 2 rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t fetch first 1 rows only for update
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t fetch first 1 rows only for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (select * from t fetch first 1 rows only) x
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from (select * from t fetch first 1 rows only) x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (select * from t fetch first 1 rows only) union (select * from t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (select * from t fetch first 1 rows only) union (select * from t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t union select * from t fetch first 1 rows only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t union select * from t fetch first 1 rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (select * from t) fetch first 1 rows only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (select * from t) fetch first 1 rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| table t fetch first 1 rows only
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| table t fetch first 1 rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| values row(1) fetch first 1 rows only
< error
| 1051 42S02 Unknown table ''
> exec ^
: string
| values row(1) fetch first 1 rows only
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 from dual fetch first 1 rows only
< ok 0 0
> exec ^
: string
| select * from ((select * from t) fetch first 1 rows only) x
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from ((select * from t) fetch first 1 rows only) x
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 in ((select 1) fetch first 1 rows only)
< ok 0 0
> exec ^
: string
| select (select 1 fetch first 1 rows only)
< ok 0 0
> exec ^
: string
| ((select 1) fetch first 1 rows only)
< ok 0 0
> exec ^
: string
| select * from t where a in (select a from t fetch first 1 rows only)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where a in (select a from t fetch first 1 rows only)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t limit 1 fetch first 1 rows only
< error
| 1064 42000 …line 1 column 29 near "fetch first 1 rows only" 
> exec ^
: string
| select * from t fetch first 1 rows only limit 1
< error
| 1064 42000 …line 1 column 45 near "limit 1" 
> exec ^
: string
| select * from t fetch first 5 rows only offset 1
< error
| 1064 42000 …line 1 column 46 near "offset 1" 
> exec ^
: string
| select * from t offset 1 rows fetch first 5 rows only
< error
| 1064 42000 …line 1 column 24 near "1 rows fetch first 5 rows only" 
> exec ^
: string
| select * from t fetch first 5 only
< error
| 1064 42000 …line 1 column 34 near "only" 
> exec ^
: string
| select * from t fetch first 5 rows
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| select * from t fetch 5 rows only
< error
| 1064 42000 …line 1 column 23 near "5 rows only" 
> exec ^
: string
| select * from t fetch first 1.5 rows only
< error
| 1064 42000 …line 1 column 31 near "1.5 rows only" 
> exec ^
: string
| select * from t fetch first -1 rows only
< error
| 1064 42000 …line 1 column 29 near "-1 rows only" 
> exec ^
: string
| select * from t fetch first (1) rows only
< error
| 1064 42000 …line 1 column 29 near "(1) rows only" 
> exec ^
: string
| select * from t fetch first a rows only
< error
| 1064 42000 …line 1 column 29 near "a rows only" 
> exec ^
: string
| update t set a = 1 fetch first 1 rows only
< error
| 1064 42000 …line 1 column 24 near "fetch first 1 rows only" 
> exec ^
: string
| delete from t fetch first 1 rows only
< error
| 1064 42000 …line 1 column 19 near "fetch first 1 rows only" 
