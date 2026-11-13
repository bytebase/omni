-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTraceFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| trace select 1
< ok 0 0
> exec ^
: string
| trace format = 'row' select 1
< ok 0 0
> exec ^
: string
| trace format = 'json' select 1
< ok 0 0
> exec ^
: string
| trace plan select 1
< ok 0 0
> exec ^
: string
| trace plan target = 'estimation' select 1
< ok 0 0
> exec ^
: string
| trace plan target = 'arandomstring' select 1
< error
| 1105 HY000 trace plan target should only be 'estimation'
> exec ^
: string
| trace plan target = 'arandomstring' select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace format = 'nosuchformat' select 1
< error
| 1105 HY000 trace format should be one of 'row', 'log' or 'json'
> exec ^
: string
| trace format = 'nosuchformat' select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace format = '' select 1
< error
| 1105 HY000 trace format should be one of 'row', 'log' or 'json'
> exec ^
: string
| trace format = '' select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace plan format = 'row' select 1
< error
| 1064 42000 …line 1 column 17 near "format = 'row' select 1" 
> exec ^
: string
| trace format = 'row' plan select 1
< error
| 1064 42000 …line 1 column 25 near "plan select 1" 
> exec ^
: string
| trace target = 'estimation' select 1
< error
| 1064 42000 …line 1 column 12 near "target = 'estimation' select 1" 
> exec ^
: string
| trace format = 'row' target = 'estimation' select 1
< error
| 1064 42000 …line 1 column 27 near "target = 'estimation' select 1" 
> exec ^
: string
| trace format = row select 1
< error
| 1064 42000 …line 1 column 18 near "row select 1" 
> exec ^
: string
| trace format = json select 1
< error
| 1064 42000 …line 1 column 19 near "json select 1" 
> exec ^
: string
| trace plan target = estimation select 1
< error
| 1064 42000 …line 1 column 30 near "estimation select 1" 
> exec ^
: string
| trace format = `row` select 1
< error
| 1064 42000 …line 1 column 20 near "`row` select 1" 
> exec ^
: string
| trace format = 1 select 1
< error
| 1064 42000 …line 1 column 16 near "1 select 1" 
> exec ^
: string
| trace format row select 1
< error
| 1064 42000 …line 1 column 16 near "row select 1" 
> exec ^
: string
| trace format = select 1
< error
| 1064 42000 …line 1 column 21 near "select 1" 
> exec ^
: string
| trace format = 'row'
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| trace plan target 'estimation' select 1
< error
| 1064 42000 …line 1 column 30 near "'estimation' select 1" 
> exec ^
: string
| trace plan target = select 1
< error
| 1064 42000 …line 1 column 26 near "select 1" 
> exec ^
: string
| trace
< error
| 1064 42000 …line 1 column 5 near "" 
> exec ^
: string
| trace plan
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| trace trace select 1
< error
| 1064 42000 …line 1 column 11 near "trace select 1" 
> exec ^
: string
| trace `select` 1
< error
| 1064 42000 …line 1 column 14 near "`select` 1" 
> exec ^
: string
| trace select c1 from t1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace select c1 from t1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace select 1 union select 2
< ok 0 0
> exec ^
: string
| trace select c1 from t1 union (select c2 from t2) limit 1, 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace select c1 from t1 union (select c2 from t2) limit 1, 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace (select 1)
< ok 0 0
> exec ^
: string
| trace (select 1) union (select 2)
< ok 0 0
> exec ^
: string
| trace with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| trace with recursive c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| trace table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace table t order by a limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace table t order by a limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace values row(1)
< ok 0 0
> exec ^
: string
| trace values row(1), row(2)
< ok 0 0
> exec ^
: string
| trace select 1 for update
< ok 0 0
> exec ^
: string
| trace select 1 for share
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| trace select 1 for share
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace table t1 for update
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace table t1 for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace table t1 union table t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace table t1 union table t2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace values row(1) union values row(2)
< ok 0 0
> exec ^
: string
| trace table t1 union select 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace table t1 union select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace (table t1) union (table t2)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace (table t1) union (table t2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace insert into t values (1), (2), (3)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace insert into t values (1), (2), (3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace insert into t select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace insert into t select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace replace into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace replace into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace update t set id = id + 1 order by id desc
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace update t set id = id + 1 order by id desc
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace update t, t2 set t.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace update t, t2 set t.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace delete from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace delete from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace delete from t where a = 1 limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace delete from t where a = 1 limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace delete t1, t2 from t1 inner join t2 where t1.id=t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| trace delete t1, t2 from t1 inner join t2 where t1.id=t2.id
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace begin
< ok 0 0
> exec ^
: string
| trace begin pessimistic
< ok 0 0
> exec ^
: string
| trace begin optimistic
< ok 0 0
> exec ^
: string
| trace begin work
< error
| 1064 42000 …line 1 column 16 near "work" 
> exec ^
: string
| trace start transaction
< ok 0 0
> exec ^
: string
| trace start transaction read only
< ok 0 0
> exec ^
: string
| trace start transaction with causal consistency only
< ok 0 0
> exec ^
: string
| trace start transaction read only as of timestamp now()
< ok 0 0
> exec ^
: string
| trace start transaction read only, with consistent snapshot
< error
| 1064 42000 …line 1 column 34 near ", with consistent snapshot" 
> exec ^
: string
| trace commit
< ok 0 0
> exec ^
: string
| trace commit and chain
< ok 0 0
> exec ^
: string
| trace commit work
< error
| 1064 42000 …line 1 column 17 near "work" 
> exec ^
: string
| trace commit and chain release
< error
| 1064 42000 …line 1 column 30 near "release" 
> exec ^
: string
| trace rollback
< ok 0 0
> exec ^
: string
| trace rollback to savepoint sp1
< ok 0 0
> exec ^
: string
| trace rollback to s
< ok 0 0
> exec ^
: string
| trace rollback to savepoint
< ok 0 0
> exec ^
: string
| trace rollback work
< error
| 1064 42000 …line 1 column 19 near "work" 
> exec ^
: string
| trace savepoint sp1
< ok 0 0
> exec ^
: string
| trace release savepoint sp1
< ok 0 0
> exec ^
: string
| trace release s
< error
| 1064 42000 …line 1 column 15 near "s" 
> exec ^
: string
| trace start replica
< error
| 1064 42000 …line 1 column 19 near "replica" 
> exec ^
: string
| trace set a = 1
< ok 0 0
> exec ^
: string
| trace set global tidb_mem_quota_query = 1
< ok 0 0
> exec ^
: string
| trace set session autocommit = 1
< ok 0 0
> exec ^
: string
| trace set @v = 1
< ok 0 0
> exec ^
: string
| trace set @@global.autocommit = 1
< ok 0 0
> exec ^
: string
| trace set names utf8mb4
< ok 0 0
> exec ^
: string
| trace set character set utf8mb4
< ok 0 0
> exec ^
: string
| trace set transaction isolation level read committed
< ok 0 0
> exec ^
: string
| trace set password = 'x'
< ok 0 0
> exec ^
: string
| trace set resource group rg1
< ok 0 0
> exec ^
: string
| trace set role r1
< error
| 1064 42000 …line 1 column 17 near "r1" 
> exec ^
: string
| trace set default role r1 to u1
< error
| 1064 42000 …line 1 column 17 near "default role r1 to u1" 
> exec ^
: string
| trace load data infile '/tmp/a.csv' into table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace load data infile '/tmp/a.csv' into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace load data local infile '/tmp/a.csv' into table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace load data local infile '/tmp/a.csv' into table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace load xml infile '/tmp/a.xml' into table t
< error
| 1064 42000 …line 1 column 14 near "xml infile '/tmp/a.xml' into table t" 
> exec ^
: string
| trace load stats 'x'
< error
| 1064 42000 …line 1 column 16 near "stats 'x'" 
> exec ^
: string
| trace analyze table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace analyze table t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace analyze table t index idx
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace analyze table t index idx
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace analyze table t all columns
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace analyze table t all columns
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace analyze incremental table t index idx
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace analyze incremental table t index idx
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace alter table t add column c int
< error
| 1064 42000 …line 1 column 11 near "alter table t add column c int" 
> exec ^
: string
| trace alter table t drop column c
< error
| 1064 42000 …line 1 column 11 near "alter table t drop column c" 
> exec ^
: string
| trace alter database d charset utf8mb4
< error
| 1064 42000 …line 1 column 11 near "alter database d charset utf8mb4" 
> exec ^
: string
| trace create table t2 (a int)
< error
| 1064 42000 …line 1 column 12 near "create table t2 (a int)" 
> exec ^
: string
| trace create index i on t (a)
< error
| 1064 42000 …line 1 column 12 near "create index i on t (a)" 
> exec ^
: string
| trace drop table t2
< error
| 1064 42000 …line 1 column 10 near "drop table t2" 
> exec ^
: string
| trace truncate table t
< error
| 1064 42000 …line 1 column 14 near "truncate table t" 
> exec ^
: string
| trace rename table t to t9
< error
| 1064 42000 …line 1 column 12 near "rename table t to t9" 
> exec ^
: string
| trace explain select 1
< error
| 1064 42000 …line 1 column 13 near "explain select 1" 
> exec ^
: string
| trace show tables
< error
| 1064 42000 …line 1 column 10 near "show tables" 
> exec ^
: string
| trace admin show ddl
< error
| 1064 42000 …line 1 column 11 near "admin show ddl" 
> exec ^
: string
| trace import into t from '/tmp/a.csv'
< error
| 1064 42000 …line 1 column 12 near "import into t from '/tmp/a.csv'" 
> exec ^
: string
| trace do 1
< error
| 1064 42000 …line 1 column 8 near "do 1" 
> exec ^
: string
| trace prepare s from 'select 1'
< error
| 1064 42000 …line 1 column 13 near "prepare s from 'select 1'" 
> exec ^
: string
| trace execute s
< error
| 1064 42000 …line 1 column 13 near "execute s" 
> exec ^
: string
| trace kill 1
< error
| 1064 42000 …line 1 column 10 near "kill 1" 
> exec ^
: string
| trace flush tables
< error
| 1064 42000 …line 1 column 11 near "flush tables" 
> exec ^
: string
| trace use d
< error
| 1064 42000 …line 1 column 9 near "use d" 
> exec ^
: string
| trace lock tables t read
< error
| 1064 42000 …line 1 column 10 near "lock tables t read" 
> exec ^
: string
| trace grant select on *.* to u
< error
| 1064 42000 …line 1 column 11 near "grant select on *.* to u" 
