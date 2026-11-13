-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestResourceControlMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| create resource group x ru_per_sec=1000
< ok 0 0
> exec ^
: string
| create resource group if not exists x ru_per_sec=1000
< ok 0 0
> exec ^
: string
| create resource group `default` ru_per_sec=100
< ok 0 0
> exec ^
: string
| create resource group `@rg` ru_per_sec=100
< ok 0 0
> exec ^
: string
| alter resource group `@rg` ru_per_sec=100
< ok 0 0
> exec ^
: string
| drop resource group `@rg`
< ok 0 0
> exec ^
: string
| drop resource group `@`
< ok 0 0
> exec ^
: string
| set resource group `@rg`
< ok 0 0
> exec ^
: string
| set resource group `@@x`
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1 QUERY_LIMIT=(ACTION SWITCH_GROUP(`@rg`))
< ok 0 0
> exec ^
: string
| create resource group x RU_PER_SEC=1000
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec 1000
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=0
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=18446744073709551615
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=UNLIMITED
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=unlimited
< ok 0 0
> exec ^
: string
| create resource group x priority=low
< ok 0 0
> exec ^
: string
| create resource group x priority=medium
< ok 0 0
> exec ^
: string
| create resource group x priority=high
< ok 0 0
> exec ^
: string
| create resource group x priority high
< ok 0 0
> exec ^
: string
| create resource group x burstable
< ok 0 0
> exec ^
: string
| create resource group x burstable=true
< ok 0 0
> exec ^
: string
| create resource group x burstable=false
< ok 0 0
> exec ^
: string
| create resource group x burstable=1
< ok 0 0
> exec ^
: string
| create resource group x burstable=0
< ok 0 0
> exec ^
: string
| create resource group x burstable=2
< ok 0 0
> exec ^
: string
| create resource group x burstable=9223372036854775807
< ok 0 0
> exec ^
: string
| create resource group x burstable true
< ok 0 0
> exec ^
: string
| create resource group x burstable 1
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=20, priority=LOW, burstable
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=20 priority=LOW burstable
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=2000, BURSTABLE
< ok 0 0
> exec ^
: string
| create resource group x burstable ru_per_sec=4000
< ok 0 0
> exec ^
: string
| create resource group x priority=low burstable ru_per_sec=1 query_limit=() background=()
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=()
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT ()
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=NULL
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT NULL
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED "10s")
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED='10s' PROCESSED_KEYS=100)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED='10s', PROCESSED_KEYS=100, RU=100)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU=100)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS=100)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1h30m')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1.5s')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '-1s')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1s' EXEC_ELAPSED '2s')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU 1, RU 2)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(PROCESSED_KEYS 1 PROCESSED_KEYS 2)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION KILL EXEC_ELAPSED='10m')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION DRYRUN)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10m' ACTION COOLDOWN)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(rg1))
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION=SWITCH_GROUP(rg1))
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(default))
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(`default`))
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH=EXACT)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH SIMILAR DURATION '10m')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH PLAN DURATION='10m')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION UNLIMITED)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION = UNLIMITED)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 'UNLIMITED')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1H')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1M')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1MS')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1H30M')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 'UnLiMiTeD')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION COOLDOWN WATCH EXACT DURATION '10m')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background (task_types 'br,lightning')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background (task_types='br,lightning')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background = (task_types='')
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT=50)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT 0)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background (task_types='br', utilization_limit=50)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background=()
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background ()
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background=NULL
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 background NULL
< ok 0 0
> exec ^
: string
| alter resource group x ru_per_sec=1000
< ok 0 0
> exec ^
: string
| alter resource group if exists x ru_per_sec=1000
< ok 0 0
> exec ^
: string
| alter resource group default priority = high
< ok 0 0
> exec ^
: string
| alter resource group x BURSTABLE
< ok 0 0
> exec ^
: string
| alter resource group x BURSTABLE, ru_per_sec=3000
< ok 0 0
> exec ^
: string
| alter resource group x ru_per_sec=UNLIMITED, BURSTABLE
< ok 0 0
> exec ^
: string
| alter resource group x QUERY_LIMIT=()
< ok 0 0
> exec ^
: string
| alter resource group x QUERY_LIMIT=NULL
< ok 0 0
> exec ^
: string
| alter resource group x background NULL
< ok 0 0
> exec ^
: string
| alter resource group default priority=low background = ( task_types "ttl" )
< ok 0 0
> exec ^
: string
| drop resource group x
< ok 0 0
> exec ^
: string
| drop resource group DEFAULT
< ok 0 0
> exec ^
: string
| drop resource group `default`
< ok 0 0
> exec ^
: string
| drop resource group if exists x
< ok 0 0
> exec ^
: string
| set resource group x
< ok 0 0
> exec ^
: string
| set resource group default
< ok 0 0
> exec ^
: string
| set resource group `default`
< ok 0 0
> exec ^
: string
| set resource group ``
< ok 0 0
> exec ^
: string
| calibrate resource
< ok 0 0
> exec ^
: string
| CALIBRATE RESOURCE
< ok 0 0
> exec ^
: string
| calibrate resource workload tpcc
< ok 0 0
> exec ^
: string
| calibrate resource workload oltp_read_write
< ok 0 0
> exec ^
: string
| calibrate resource workload oltp_read_only
< ok 0 0
> exec ^
: string
| calibrate resource workload oltp_write_only
< ok 0 0
> exec ^
: string
| calibrate resource workload tpch_10
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME '2021-04-15 00:00:00'
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME = '2021-04-15 00:00:00'
< ok 0 0
> exec ^
: string
| calibrate resource END_TIME now()
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME '2023-04-01 13:00:00' END_TIME '2023-04-01 16:00:00' DURATION '20m'
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME '2023-04-01 13:00:00',END_TIME='2023-04-01 16:00:00'
< ok 0 0
> exec ^
: string
| calibrate resource DURATION='20m' START_TIME '2023-04-01 13:00:00'
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 20 minute
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME now() - 1000 END_TIME current_timestamp()
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME a
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME t.a
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME (select 1)
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME @v
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME -1
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME null
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME a is null
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME a between 1 and 2
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME exists (select 1)
< ok 0 0
> exec ^
: string
| calibrate resource DURATION '20d'
< ok 0 0
> exec ^
: string
| calibrate resource DURATION '1d12h'
< ok 0 0
> exec ^
: string
| calibrate resource DURATION '1.5h'
< ok 0 0
> exec ^
: string
| calibrate resource DURATION ''
< ok 0 0
> exec ^
: string
| calibrate resource DURATION interval 15 minute
< ok 0 0
> exec ^
: string
| calibrate resource DURATION = interval 15 minute
< ok 0 0
> exec ^
: string
| calibrate resource DURATION interval '15' minute
< ok 0 0
> exec ^
: string
| calibrate resource DURATION interval 15 MINUTE
< ok 0 0
> exec ^
: string
| calibrate resource DURATION interval 15 day_hour
< ok 0 0
> exec ^
: string
| calibrate resource DURATION interval 15 sql_tsi_minute
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 20 minute DURATION interval 15 minute
< ok 0 0
> exec ^
: string
| calibrate resource START_TIME '2020-01-01' END_TIME '2020-01-02' DURATION '1m' DURATION '2m'
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT SIMILAR to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT EXACT to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT PLAN to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST 'b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377'
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377
< ok 0 0
> exec ^
: string
| query watch add PLAN DIGEST `5e3ddd388f6012e328233dbcdda5d48f404e0536c6c54d9618233210f3d5762a`
< ok 0 0
> exec ^
: string
| query watch add PLAN DIGEST @digest1
< ok 0 0
> exec ^
: string
| query watch add resource group rg1 SQL TEXT SIMILAR to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add resource group `default` SQL TEXT SIMILAR to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add resource group default SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| query watch add resource group @rg SQL TEXT SIMILAR to @sql1
< ok 0 0
> exec ^
: string
| query watch add resource group rg1
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT SIMILAR to 'select 1' resource group rg1
< ok 0 0
> exec ^
: string
| query watch add ACTION = KILL SQL TEXT SIMILAR to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add ACTION COOLDOWN resource group rg1 SQL TEXT SIMILAR to 'select 1'
< ok 0 0
> exec ^
: string
| query watch add ACTION DRYRUN SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| query watch add ACTION SWITCH_GROUP(rg1) SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| query watch add ACTION KILL
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST 'a', resource group rg1
< ok 0 0
> exec ^
: string
| query watch add resource group rg1, action kill, sql digest 'a'
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT SIMILAR to 'select 1' resource group rg1 action kill
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST concat('a','b')
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST -1
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT EXACT to 'select 1' COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST 'a' COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| query watch add PLAN DIGEST 'a' COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT EXACT to 'select 1' COLLATE 'utf8mb4_bin'
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT EXACT to 'select 1' COLLATE binary
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT EXACT to 'select 1' COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST 'a' COLLATE utf8mb4_bin, resource group rg
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST cfg->>'$.digest'
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST cfg->'$.digest'
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST d.cfg->'$.a'
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST a.b.cfg->'$.a'
< ok 0 0
> exec ^
: string
| query watch add PLAN DIGEST cfg->>'$.digest'
< ok 0 0
> exec ^
: string
| query watch add SQL TEXT EXACT to cfg->>'$.t'
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST cfg->>'$.a' COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST cfg->>'$.a' resource group rg
< ok 0 0
> exec ^
: string
| query watch add SQL DIGEST cfg->>'$.a' action kill
< ok 0 0
> exec ^
: string
| query watch add resource group `@rg` SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| query watch add resource group `@@x` SQL DIGEST 'a'
< ok 0 0
> exec ^
: string
| query watch remove 1
< ok 0 0
> exec ^
: string
| query watch remove 0
< ok 0 0
> exec ^
: string
| query watch remove 9223372036854775807
< ok 0 0
> exec ^
: string
| QUERY WATCH REMOVE 1
< ok 0 0
> exec ^
: string
| CREATE USER 'root'@'127.0.0.1' IDENTIFIED BY 'hashstring' RESOURCE GROUP rg1
< ok 0 0
> exec ^
: string
| ALTER USER 'root'@'localhost' RESOURCE GROUP rg2
< ok 0 0
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP default
< ok 0 0
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP `rg`
< ok 0 0
> exec ^
: string
| create user 'u'@'h' comment 'c' RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| create user 'u'@'h' account lock RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| create user 'u'@'h' require ssl RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| create user 'u'@'h' with max_queries_per_hour 1 RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| alter user if exists 'u'@'h' RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| alter user 'u'@'h' identified by 'p' RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| create table priority (duration int, plan int, watch int, background int, ru int, exact int, similar int, digest int, unlimited int, low int, high int, start_time int, end_time int, burstable int, calibrate int, workload int, task_types int, dryrun int, cooldown int, switch_group int, tpcc int, query_limit int, exec_elapsed int, processed_keys int, utilization_limit int, oltp_read_write int, tpch_10 int)
< ok 0 0
> exec ^
: string
| select priority, duration, plan, watch from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select priority, duration, plan, watch from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.plan, t.duration from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.plan, t.duration from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 as priority, 2 as duration, 3 as plan, 4 as watch, 5 as calibrate, 6 as unlimited
< ok 0 0
> exec ^
: string
| select * from t priority
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t priority
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as plan
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as plan
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t duration join t2 workload on duration.a = workload.a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t duration join t2 workload on duration.a = workload.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t plan where plan.id=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t plan where plan.id=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t priority where priority.id=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t priority where priority.id=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t duration where duration.id=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t duration where duration.id=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t text where text.id=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t text where text.id=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t action where action.id=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t action where action.id=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t config where config.id=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t config where config.id=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t priority set priority.a=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t priority set priority.a=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create index priority on t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index priority on t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create index duration on t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index duration on t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add index duration (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add index duration (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t rename column a to priority
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t rename column a to priority
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| savepoint priority
< ok 0 0
> exec ^
: string
| savepoint duration
< ok 0 0
> exec ^
: string
| release savepoint plan
< ok 0 0
> exec ^
: string
| create database priority
< ok 0 0
> exec ^
: string
| create database duration
< ok 0 0
> exec ^
: string
| with priority as (select 1) select * from priority
< ok 0 0
> exec ^
: string
| create user duration
< ok 0 0
> exec ^
: string
| create user unlimited
< ok 0 0
> exec ^
: string
| create role `plan`
< ok 0 0
> exec ^
: string
| select current_resource_group()
< ok 0 0
> exec ^
: string
| create resource group x burstable=18446744073709551615
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=18446744073709551615
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=9223372036854775808
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=9223372036854775808
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 9223372036854775808)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 9223372036854775808)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 9223372036854775808
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 9223372036854775808
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 18446744073709551615
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 18446744073709551615
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| alter resource group x
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create resource group if exists x ru_per_sec=1000
< error
| 1064 42000 …line 1 column 31 near "exists x ru_per_sec=1000" 
> exec ^
: string
| create resource group x ru_per_sec=1000 ru_per_sec=2000
< error
| 1064 42000 …line 1 column 55 near ""Dupliated options specified 
> exec ^
: string
| create resource group x priority=LOW ru_per_sec=1000 priority=HIGH
< error
| 1064 42000 …line 1 column 66 near ""Dupliated options specified 
> exec ^
: string
| create resource group x ru_per_sec=1000, ru_per_sec=2000
< error
| 1064 42000 …line 1 column 56 near ""Dupliated options specified 
> exec ^
: string
| create resource group x foo=1
< error
| 1064 42000 …line 1 column 27 near "foo=1" 
> exec ^
: string
| create resource group x cpu ='8c'
< error
| 1064 42000 …line 1 column 27 near "cpu ='8c'" 
> exec ^
: string
| create resource group x region ='us, 3'
< error
| 1064 42000 …line 1 column 30 near "region ='us, 3'" 
> exec ^
: string
| create resource group x followers=0
< error
| 1064 42000 …line 1 column 33 near "followers=0" 
> exec ^
: string
| create resource group x cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'
< error
| 1064 42000 …line 1 column 27 near "cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'" 
> exec ^
: string
| create resource group x ru_per_sec='check'
< error
| 1064 42000 …line 1 column 42 near "'check'" 
> exec ^
: string
| create resource group x ru_per_sec=
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| create resource group x ru_per_sec=-1
< error
| 1064 42000 …line 1 column 36 near "-1" 
> exec ^
: string
| create resource group x ru_per_sec=+1000
< error
| 1064 42000 …line 1 column 36 near "+1000" 
> exec ^
: string
| create resource group x ru_per_sec=1.5
< error
| 1064 42000 …line 1 column 38 near "1.5" 
> exec ^
: string
| create resource group x ru_per_sec=18446744073709551616
< error
| 1064 42000 …line 1 column 55 near "18446744073709551616" 
> exec ^
: string
| create resource group x priority=mid
< error
| 1064 42000 …line 1 column 36 near "mid" 
> exec ^
: string
| create resource group x priority=1
< error
| 1064 42000 …line 1 column 34 near "1" 
> exec ^
: string
| create resource group x priority='low'
< error
| 1064 42000 …line 1 column 38 near "'low'" 
> exec ^
: string
| create resource group x ru_per_sec=20 priority=MID BURSTABLE
< error
| 1064 42000 …line 1 column 50 near "MID BURSTABLE" 
> exec ^
: string
| create resource group x `ru_per_sec`=1000
< error
| 1064 42000 …line 1 column 36 near "`ru_per_sec`=1000" 
> exec ^
: string
| create resource group x ru_per_sec=1000,
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| create resource group x burstable=disable
< error
| 1064 42000 …line 1 column 41 near "disable" 
> exec ^
: string
| create resource group x burstable=unlimited
< error
| 1064 42000 …line 1 column 43 near "unlimited" 
> exec ^
: string
| create resource group x burstable=moderated
< error
| 1064 42000 …line 1 column 43 near "moderated" 
> exec ^
: string
| create resource group x burstable=off
< error
| 1064 42000 …line 1 column 37 near "off" 
> exec ^
: string
| create resource group x burstable=on
< error
| 1064 42000 …line 1 column 36 near "on" 
> exec ^
: string
| create resource group x burstable=null
< error
| 1064 42000 …line 1 column 38 near "null" 
> exec ^
: string
| create resource group x burstable=default
< error
| 1064 42000 …line 1 column 41 near "default" 
> exec ^
: string
| create resource group x burstable='true'
< error
| 1064 42000 …line 1 column 40 near "'true'" 
> exec ^
: string
| create resource group x burstable=-1
< error
| 1064 42000 …line 1 column 35 near "-1" 
> exec ^
: string
| create resource group x burstable=1.5
< error
| 1064 42000 …line 1 column 37 near "1.5" 
> exec ^
: string
| create resource group x burstable=18446744073709551616
< error
| 1064 42000 …line 1 column 54 near "18446744073709551616" 
> exec ^
: string
| alter resource group x BURSTABLE = UNLIMITED ru_per_sec=4000
< error
| 1064 42000 …line 1 column 44 near "UNLIMITED ru_per_sec=4000" 
> exec ^
: string
| alter resource group x ru_per_sec=2000, BURSTABLE=yes
< error
| 1064 42000 …line 1 column 53 near "yes" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY=(EXEC_ELAPSED '10s')
< error
| 1064 42000 …line 1 column 45 near "QUERY=(EXEC_ELAPSED '10s')" 
> exec ^
: string
| create resource group x ru_per_sec=1000 LIMIT=(EXEC_ELAPSED '10s')
< error
| 1064 42000 …line 1 column 45 near "LIMIT=(EXEC_ELAPSED '10s')" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=EXEC_ELAPSED '10s'
< error
| 1064 42000 …line 1 column 64 near "EXEC_ELAPSED '10s'" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT EXEC_ELAPSED '10s'
< error
| 1064 42000 …line 1 column 64 near "EXEC_ELAPSED '10s'" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s'
< error
| 1064 42000 …line 1 column 73 near "" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(DURATION '1m')
< error
| 1064 42000 …line 1 column 61 near "DURATION '1m')" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED 'xyz')
< error
| 1064 42000 …line 1 column 72 near ")"The EXEC_ELAPSED option is not a valid duration: time: invalid duration "xyz" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10')
< error
| 1064 42000 …line 1 column 71 near ")"The EXEC_ELAPSED option is not a valid duration: time: missing unit in duration "10" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '')
< error
| 1064 42000 …line 1 column 69 near ")"The EXEC_ELAPSED option is not a valid duration: time: invalid duration "" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1d')
< error
| 1064 42000 …line 1 column 71 near ")"The EXEC_ELAPSED option is not a valid duration: time: unknown unit "d" in duration "1d" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED 100)
< error
| 1064 42000 …line 1 column 69 near "100)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1H')
< error
| 1064 42000 …line 1 column 71 near ")"The EXEC_ELAPSED option is not a valid duration: time: unknown unit "H" in duration "1H" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1M')
< error
| 1064 42000 …line 1 column 71 near ")"The EXEC_ELAPSED option is not a valid duration: time: unknown unit "M" in duration "1M" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '1MS')
< error
| 1064 42000 …line 1 column 72 near ")"The EXEC_ELAPSED option is not a valid duration: time: unknown unit "MS" in duration "1MS" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(PROCESSED_KEYS='100')
< error
| 1064 42000 …line 1 column 73 near "'100')" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(PROCESSED_KEYS -1)
< error
| 1064 42000 …line 1 column 69 near "-1)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU='100')
< error
| 1064 42000 …line 1 column 61 near "'100')" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU 1.5)
< error
| 1064 42000 …line 1 column 59 near "1.5)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU -1)
< error
| 1064 42000 …line 1 column 57 near "-1)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU +1)
< error
| 1064 42000 …line 1 column 57 near "+1)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(RU 18446744073709551616)
< error
| 1064 42000 …line 1 column 76 near "18446744073709551616)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN ACTION KILL)
< error
| 1064 42000 …line 1 column 100 near ")"Dupliated runaway options specified 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 100 ACTION DRYRUN ACTION KILL)
< error
| 1064 42000 …line 1 column 100 near ")"Dupliated runaway options specified 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT WATCH SIMILAR)
< error
| 1064 42000 …line 1 column 79 near ")"Dupliated runaway options specified 
> exec ^
: string
| alter resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN WATCH SIMILAR DURATION '10m' ACTION COOLDOWN)
< error
| 1064 42000 …line 1 column 132 near ")"Dupliated runaway options specified 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION STOP)
< error
| 1064 42000 …line 1 column 64 near "STOP)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP())
< error
| 1064 42000 …line 1 column 74 near "))" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION SWITCH_GROUP(@rg))
< error
| 1064 42000 …line 1 column 76 near "@rg))" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH ROUGH)
< error
| 1064 42000 …line 1 column 64 near "ROUGH)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 'xyz')
< error
| 1064 42000 …line 1 column 80 near ")"The WATCH DURATION option is not a valid duration: time: invalid duration "xyz" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1d')
< error
| 1064 42000 …line 1 column 79 near ")"The WATCH DURATION option is not a valid duration: time: unknown unit "d" in duration "1d" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION 10)
< error
| 1064 42000 …line 1 column 76 near "10)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT=(WATCH EXACT DURATION '1D')
< error
| 1064 42000 …line 1 column 79 near ")"The WATCH DURATION option is not a valid duration: time: unknown unit "d" in duration "1d" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(task_types br)
< error
| 1064 42000 …line 1 column 65 near "br)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(task_types 1)
< error
| 1064 42000 …line 1 column 64 near "1)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(utilization_limit '50')
< error
| 1064 42000 …line 1 column 74 near "'50')" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT="NAN")
< error
| 1064 42000 …line 1 column 77 near ""NAN")" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(utilization_limit=-1)
< error
| 1064 42000 …line 1 column 71 near "-1)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(utilization_limit 18446744073709551616)
< error
| 1064 42000 …line 1 column 90 near "18446744073709551616)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(task_types 'br' task_types 'lightning')
< error
| 1064 42000 …line 1 column 91 near ")"Dupliated background options specified 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(foo 1)
< error
| 1064 42000 …line 1 column 55 near "foo 1)" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=task_types 'br'
< error
| 1064 42000 …line 1 column 61 near "task_types 'br'" 
> exec ^
: string
| create resource group x ru_per_sec=1000 query_limit=(exec_elapsed '1s',)
< error
| 1064 42000 …line 1 column 72 near ")" 
> exec ^
: string
| create resource group x ru_per_sec=1000 background=(task_types 'br',)
< error
| 1064 42000 …line 1 column 69 near ")" 
> exec ^
: string
| create resource group x ru_per_sec=1000 query_limit=(,)
< error
| 1064 42000 …line 1 column 54 near ",)" 
> exec ^
: string
| drop resource group x,y
< error
| 1064 42000 …line 1 column 22 near ",y" 
> exec ^
: string
| drop resource group if exists x,y
< error
| 1064 42000 …line 1 column 32 near ",y" 
> exec ^
: string
| drop resource group
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set resource group x y
< error
| 1064 42000 …line 1 column 22 near "y" 
> exec ^
: string
| set resource group
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create resource group x type=user
< error
| 1064 42000 …line 1 column 28 near "type=user" 
> exec ^
: string
| create resource group x type=user vcpu=0-63 thread_priority=0 enable
< error
| 1064 42000 …line 1 column 28 near "type=user vcpu=0-63 thread_priority=0 enable" 
> exec ^
: string
| create resource group x vcpu=0-63
< error
| 1064 42000 …line 1 column 28 near "vcpu=0-63" 
> exec ^
: string
| create resource group x thread_priority=0
< error
| 1064 42000 …line 1 column 39 near "thread_priority=0" 
> exec ^
: string
| create resource group x ru_per_sec=1000 enable
< error
| 1064 42000 …line 1 column 46 near "enable" 
> exec ^
: string
| alter resource group x force
< error
| 1064 42000 …line 1 column 28 near "force" 
> exec ^
: string
| alter resource group x enable
< error
| 1064 42000 …line 1 column 29 near "enable" 
> exec ^
: string
| drop resource group x force
< error
| 1064 42000 …line 1 column 27 near "force" 
> exec ^
: string
| set resource group x for 1
< error
| 1064 42000 …line 1 column 24 near "for 1" 
> exec ^
: string
| calibrate resource foo
< error
| 1064 42000 …line 1 column 22 near "foo" 
> exec ^
: string
| calibrate resource workload
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| calibrate resource workload foo
< error
| 1064 42000 …line 1 column 31 near "foo" 
> exec ^
: string
| calibrate resource workload = tpcc
< error
| 1064 42000 …line 1 column 29 near "= tpcc" 
> exec ^
: string
| calibrate resource workload = oltp_read_write START_TIME '2023-04-01 13:00:00'
< error
| 1064 42000 …line 1 column 29 near "= oltp_read_write START_TIME '2023-04-01 13:00:00'" 
> exec ^
: string
| calibrate resource workload tpcc START_TIME '2023-04-01 13:00:00'
< error
| 1064 42000 …line 1 column 43 near "START_TIME '2023-04-01 13:00:00'" 
> exec ^
: string
| calibrate resource START_TIME '2023-04-01 13:00:00' workload tpcc
< error
| 1064 42000 …line 1 column 60 near "workload tpcc" 
> exec ^
: string
| calibrate resource workload tpcc, workload tpcc
< error
| 1064 42000 …line 1 column 33 near ", workload tpcc" 
> exec ^
: string
| calibrate resource START_TIME
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| calibrate resource DURATION
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| calibrate resource DURATION 20
< error
| 1064 42000 …line 1 column 30 near "20" 
> exec ^
: string
| calibrate resource DURATION 'xyz'
< error
| 1064 42000 …line 1 column 33 near ""The DURATION option is not a valid duration: fail to read an integer 
> exec ^
: string
| calibrate resource DURATION '20'
< error
| 1064 42000 …line 1 column 32 near ""The DURATION option is not a valid duration: fail to read an integer 
> exec ^
: string
| calibrate resource DURATION '1w'
< error
| 1064 42000 …line 1 column 32 near ""The DURATION option is not a valid duration: unknown unit w 
> exec ^
: string
| calibrate resource DURATION '-1m'
< error
| 1064 42000 …line 1 column 33 near ""The DURATION option is not a valid duration: fail to read an integer 
> exec ^
: string
| calibrate resource DURATION 'unlimited'
< error
| 1064 42000 …line 1 column 39 near ""The DURATION option is not a valid duration: fail to read an integer 
> exec ^
: string
| calibrate resource DURATION interval 15
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| calibrate resource DURATION interval 15 `MINUTE`
< error
| 1064 42000 …line 1 column 48 near "`MINUTE`" 
> exec ^
: string
| calibrate resource DURATION interval 15 `day_hour`
< error
| 1064 42000 …line 1 column 50 near "`day_hour`" 
> exec ^
: string
| calibrate resource DURATION interval 15 'minute'
< error
| 1064 42000 …line 1 column 48 near "'minute'" 
> exec ^
: string
| calibrate resource START_TIME '2020-01-01',
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| calibrate resource START_TIME '2020-01-01' START_TIME '2020-01-02'
< error
| 1064 42000 …line 1 column 66 near ""Dupliated options specified 
> exec ^
: string
| calibrate resource START_TIME '2020-01-01' END_TIME '2020-01-02' END_TIME '2020-01-03'
< error
| 1064 42000 …line 1 column 86 near ""Dupliated options specified 
> exec ^
: string
| calibrate resource DURATION '1m' DURATION '2m'
< error
| 1064 42000 …line 1 column 46 near ""Dupliated options specified 
> exec ^
: string
| calibrate resource DURATION '1m' START_TIME '2020-01-01' DURATION '2m'
< error
| 1064 42000 …line 1 column 70 near ""Dupliated options specified 
> exec ^
: string
| query watch add
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| query watch add SQL SIMILAR to 'select 1'
< error
| 1064 42000 …line 1 column 27 near "SIMILAR to 'select 1'" 
> exec ^
: string
| query watch add SQL TEXT SIMILAR 'select 1'
< error
| 1064 42000 …line 1 column 43 near "'select 1'" 
> exec ^
: string
| query watch add SQL TEXT SIMILAR to
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| query watch add SQL TEXT FOO to 'select 1'
< error
| 1064 42000 …line 1 column 28 near "FOO to 'select 1'" 
> exec ^
: string
| query watch add PLAN DIGEST
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| query watch add SQL DIGEST 'a' SQL DIGEST 'b'
< error
| 1064 42000 …line 1 column 45 near ""Dupliated options specified 
> exec ^
: string
| query watch add SQL DIGEST 'a' PLAN DIGEST 'b'
< error
| 1064 42000 …line 1 column 46 near ""Dupliated options specified 
> exec ^
: string
| query watch add resource group `default` resource group `rg1` SQL TEXT SIMILAR to 'select 1'
< error
| 1064 42000 …line 1 column 65 near "SQL TEXT SIMILAR to 'select 1'"Dupliated options specified 
> exec ^
: string
| query watch add ACTION KILL ACTION DRYRUN SQL DIGEST 'a'
< error
| 1064 42000 …line 1 column 45 near "SQL DIGEST 'a'"Dupliated options specified 
> exec ^
: string
| query watch add ACTION SWITCH_GROUP(@rg) SQL DIGEST 'a'
< error
| 1064 42000 …line 1 column 39 near "@rg) SQL DIGEST 'a'" 
> exec ^
: string
| query watch add RESOURCE GROUP @@rg SQL DIGEST 'a'
< error
| 1064 42000 …line 1 column 35 near "@@rg SQL DIGEST 'a'" 
> exec ^
: string
| query watch add RESOURCE GROUP @@global.rg SQL DIGEST 'a'
< error
| 1064 42000 …line 1 column 42 near "@@global.rg SQL DIGEST 'a'" 
> exec ^
: string
| query watch add SQL DIGEST 'a',
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| query watch add SQL DIGEST 'a' extra
< error
| 1064 42000 …line 1 column 36 near "extra" 
> exec ^
: string
| query watch add SQL TEXT EXACT to 'select 1' COLLATE
< error
| 1064 42000 …line 1 column 52 near "" 
> exec ^
: string
| query watch add SQL DIGEST cfg->>1
< error
| 1064 42000 …line 1 column 34 near "1" 
> exec ^
: string
| query watch add SQL DIGEST cfg->1
< error
| 1064 42000 …line 1 column 33 near "1" 
> exec ^
: string
| query watch add SQL DIGEST 'a'->>'$.x'
< error
| 1064 42000 …line 1 column 33 near "->>'$.x'" 
> exec ^
: string
| query watch add SQL DIGEST (cfg)->>'$.x'
< error
| 1064 42000 …line 1 column 35 near "->>'$.x'" 
> exec ^
: string
| query watch add SQL DIGEST now()->>'$.a'
< error
| 1064 42000 …line 1 column 35 near "->>'$.a'" 
> exec ^
: string
| query watch add SQL DIGEST @v->>'$.a'
< error
| 1064 42000 …line 1 column 32 near "->>'$.a'" 
> exec ^
: string
| query watch add SQL DIGEST cfg->'$.a'->'$.b'
< error
| 1064 42000 …line 1 column 39 near "->'$.b'" 
> exec ^
: string
| query watch add SQL DIGEST cfg->>'$.d'->>'$.e'
< error
| 1064 42000 …line 1 column 41 near "->>'$.e'" 
> exec ^
: string
| query watch add SQL DIGEST cfg COLLATE utf8mb4_bin ->> '$.x'
< error
| 1064 42000 …line 1 column 54 near "->> '$.x'" 
> exec ^
: string
| query watch add SQL DIGEST cfg->>
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| query watch remove
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| query watch remove -1
< error
| 1064 42000 …line 1 column 20 near "-1" 
> exec ^
: string
| query watch remove 'a'
< error
| 1064 42000 …line 1 column 22 near "'a'" 
> exec ^
: string
| query watch remove @v
< error
| 1064 42000 …line 1 column 21 near "@v" 
> exec ^
: string
| query watch remove 1 2
< error
| 1064 42000 …line 1 column 22 near "2" 
> exec ^
: string
| query watch remove 18446744073709551616
< error
| 1064 42000 …line 1 column 39 near "18446744073709551616" 
> exec ^
: string
| query watch remove resource group rg1
< error
| 1064 42000 …line 1 column 27 near "resource group rg1" 
> exec ^
: string
| query watch remove resource group @rg
< error
| 1064 42000 …line 1 column 27 near "resource group @rg" 
> exec ^
: string
| query watch 1
< error
| 1064 42000 …line 1 column 13 near "1" 
> exec ^
: string
| drop query watch 1
< error
| 1064 42000 …line 1 column 10 near "query watch 1" 
> exec ^
: string
| drop query watch resource group rg1
< error
| 1064 42000 …line 1 column 10 near "query watch resource group rg1" 
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP rg comment 'c'
< error
| 1064 42000 …line 1 column 45 near "comment 'c'" 
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP rg account lock
< error
| 1064 42000 …line 1 column 45 near "account lock" 
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP rg RESOURCE GROUP rg2
< error
| 1064 42000 …line 1 column 46 near "RESOURCE GROUP rg2" 
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP 'rg'
< error
| 1064 42000 …line 1 column 39 near "'rg'" 
> exec ^
: string
| create user 'u'@'h' RESOURCE GROUP
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| create user 'u'@'h' RESOURCE rg
< error
| 1064 42000 …line 1 column 31 near "rg" 
> exec ^
: string
| CREATE ROLE RESOURCE
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role plan
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role duration
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role priority
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| create role action
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| create role text
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| create role r1, resource
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| drop role resource
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set role resource
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set role all except resource
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| set default role resource to u1
< error
| 1064 42000 …line 1 column 28 near "to u1" 
