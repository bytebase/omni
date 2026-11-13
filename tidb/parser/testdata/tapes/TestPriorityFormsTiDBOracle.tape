-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestPriorityFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select low_priority * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select low_priority * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select delayed * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select delayed * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select high_priority * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select high_priority * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select low_priority high_priority delayed a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select low_priority high_priority delayed a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select distinct low_priority straight_join a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select distinct low_priority straight_join a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| (select low_priority 1) union (select delayed 2)
< ok 0 0
> exec ^
: string
| select low_priority from t
< error
| 1064 42000 …line 1 column 24 near "from t" 
> exec ^
: string
| select a from t where low_priority = 1
< error
| 1064 42000 …line 1 column 34 near "low_priority = 1" 
> exec ^
: string
| select SQL_CACHE * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select SQL_CACHE * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_cache sql_no_cache a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_cache sql_no_cache a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_no_cache sql_cache a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_no_cache sql_cache a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_cache.a from sql_cache
< error
| 1146 42S02 Table 'test.sql_cache' doesn't exist
> exec ^
: string
| select sql_cache.a from sql_cache
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_cache from t
< error
| 1064 42000 …line 1 column 21 near "from t" 
> exec ^
: string
| select `sql_cache` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `sql_cache` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update low_priority t set a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update low_priority t set a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update high_priority t set a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update high_priority t set a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update delayed t set a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update delayed t set a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update high_priority ignore t set a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update high_priority ignore t set a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update high_priority t1, t2 set t1.a = 2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| update high_priority t1, t2 set t1.a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update ignore high_priority t set a = 2
< error
| 1064 42000 …line 1 column 27 near "high_priority t set a = 2" 
> exec ^
: string
| update low_priority high_priority t set a = 2
< error
| 1064 42000 …line 1 column 33 near "high_priority t set a = 2" 
> exec ^
: string
| delete high_priority from t where a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete high_priority from t where a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete delayed from t where a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete delayed from t where a = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete low_priority quick ignore from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete low_priority quick ignore from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete delayed quick ignore t1 from t1, t2 where t1.a = t2.a
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| delete delayed quick ignore t1 from t1, t2 where t1.a = t2.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete high_priority from t1 using t1, t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| delete high_priority from t1 using t1, t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| with c as (select 1) delete high_priority from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| with c as (select 1) delete high_priority from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete quick high_priority from t
< error
| 1064 42000 …line 1 column 26 near "high_priority from t" 
> exec ^
: string
| delete high_priority low_priority from t
< error
| 1064 42000 …line 1 column 33 near "low_priority from t" 
> exec ^
: string
| insert delayed ignore into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert delayed ignore into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace high_priority into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace high_priority into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace delayed into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace delayed into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
