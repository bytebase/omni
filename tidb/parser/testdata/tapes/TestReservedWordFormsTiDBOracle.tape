-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestReservedWordFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| alter table t add column a int after before
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int after before
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a int first
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int first
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column before int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column before int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create index before on t (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index before on t (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table function (a int)
< ok 0 0
> exec ^
: string
| create table loop (a int)
< ok 0 0
> exec ^
: string
| create table return (signal int, resignal int, get int, purge int, undo int, declare int)
< ok 0 0
> exec ^
: string
| create table t (a int, constraint signal unique (a))
< ok 0 0
> exec ^
: string
| create table t (a int, index return (a))
< ok 0 0
> exec ^
: string
| create user declare
< ok 0 0
> exec ^
: string
| create view loop as select 1
< ok 0 0
> exec ^
: string
| deallocate prepare signal
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| deallocate prepare signal
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from declare
< error
| 1146 42S02 Table 'test.declare' doesn't exist
> exec ^
: string
| delete from declare
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| execute signal
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| execute signal
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| prepare signal from 'select 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| prepare signal from 'select 1'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| rollback to loop
< ok 0 0
> exec ^
: string
| rollback to savepoint return
< ok 0 0
> exec ^
: string
| savepoint loop
< ok 0 0
> exec ^
: string
| select 'x' return
< ok 0 0
> exec ^
: string
| select (1) signal
< ok 0 0
> exec ^
: string
| select * from (select 1) return
< ok 0 0
> exec ^
: string
| select * from full
< error
| 1146 42S02 Table 'test.full' doesn't exist
> exec ^
: string
| select * from full
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from json_table a
< error
| 1146 42S02 Table 'test.json_table' doesn't exist
> exec ^
: string
| select * from json_table a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from json_table
< error
| 1146 42S02 Table 'test.json_table' doesn't exist
> exec ^
: string
| select * from json_table
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from json_table.a
< error
| 1146 42S02 Table 'json_table.a' doesn't exist
> exec ^
: string
| select * from json_table.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from lateral a
< error
| 1146 42S02 Table 'test.lateral' doesn't exist
> exec ^
: string
| select * from lateral a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from lateral
< error
| 1146 42S02 Table 'test.lateral' doesn't exist
> exec ^
: string
| select * from lateral
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from lateral.a
< error
| 1146 42S02 Table 'lateral.a' doesn't exist
> exec ^
: string
| select * from lateral.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t a join loop on a.b = loop.c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a join loop on a.b = loop.c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t a join return using (b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a join return using (b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t a straight_join signal on 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a straight_join signal on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t for update of empty
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t for update of empty
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t for update return
< error
| 1064 42000 …line 1 column 33 near "return" 
> exec ^
: string
| select * from t group by loop
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t group by loop
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t having get = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t having get = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t limit 1 loop
< error
| 1064 42000 …line 1 column 28 near "loop" 
> exec ^
: string
| select * from t order by return
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t order by return
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where a = 1 return
< error
| 1064 42000 …line 1 column 34 near "return" 
> exec ^
: string
| select * from t where each = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where each = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where function = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where function = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t, lateral
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t, lateral
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 + 1 get
< ok 0 0
> exec ^
: string
| select 1 from t tablesample system (1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select 1 from t tablesample system (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 from t tablesample system
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| select 1 loop
< ok 0 0
> exec ^
: string
| select 1 signal, 2 resignal
< ok 0 0
> exec ^
: string
| select @a := return
< error
| 1054 42S22 Unknown column 'return' in ''
> exec ^
: string
| select @a := return
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a as `return` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a as `return` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a before from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a before from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a from t a full join b on 1
< error
| 1064 42000 …line 1 column 22 near "full join b on 1" 
> exec ^
: string
| select a from t full join b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t full join b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a from t full outer join b
< error
| 1064 42000 …line 1 column 26 near "outer join b" 
> exec ^
: string
| select a from t full
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t full
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a from t group by a with cube
< error
| 1064 42000 …line 1 column 36 near "cube" 
> exec ^
: string
| select a from t group by cube(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t group by cube(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a full from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a full from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a full
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select a full
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a full, b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a full, b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a function from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a function from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a get from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a get from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a loop from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a loop from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a return from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a return from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a signal from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a signal from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a system from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a system from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.return from a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select a.return from a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select accessible from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select accessible from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select asensitive from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select asensitive from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select condition from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select condition from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cube(1)
< error
| 1305 42000 FUNCTION test.cube does not exist
> exec ^
: string
| select cube(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cube(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cube(a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select declare from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select declare from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select distinct sql_no_cache.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select distinct sql_no_cache.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select empty from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select empty from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select full from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select full from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select function from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select function from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select function.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select function.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select get from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select get from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select get
< error
| 1054 42S22 Unknown column 'get' in 'field list'
> exec ^
: string
| select get
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select get(1)
< error
| 1305 42000 FUNCTION test.get does not exist
> exec ^
: string
| select get(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select insensitive from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select insensitive from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select io_after_gtids from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select io_after_gtids from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select io_before_gtids from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select io_before_gtids from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select json_table from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select json_table from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select lateral from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select lateral from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select loop(1)
< error
| 1305 42000 FUNCTION test.loop does not exist
> exec ^
: string
| select loop(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select optimizer_costs from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select optimizer_costs from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select read_write from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select read_write from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select reads from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select reads from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select return(1)
< error
| 1305 42000 FUNCTION test.return does not exist
> exec ^
: string
| select return(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select return.a from return
< error
| 1146 42S02 Table 'test.return' doesn't exist
> exec ^
: string
| select return.a from return
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sensitive from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sensitive from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select signal(1)
< error
| 1305 42000 FUNCTION test.signal does not exist
> exec ^
: string
| select signal(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_buffer_result a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_buffer_result a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_buffer_result from t
< error
| 1064 42000 …line 1 column 29 near "from t" 
> exec ^
: string
| select sql_buffer_result.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_buffer_result.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_no_cache a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_no_cache a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_no_cache from t
< error
| 1064 42000 …line 1 column 24 near "from t" 
> exec ^
: string
| select sql_no_cache sql_no_cache.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_no_cache sql_no_cache.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sql_no_cache
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| select sql_no_cache.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sql_no_cache.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select system from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select system from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select system
< error
| 1054 42S22 Unknown column 'system' in 'field list'
> exec ^
: string
| select system
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update return set loop = 1
< error
| 1146 42S02 Table 'test.return' doesn't exist
> exec ^
: string
| update return set loop = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| use system
< ok 0 0
> exec ^
: string
| with get as (select 1) select * from get
< ok 0 0
> exec ^
: string
| alter table t modify a int after loop
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify a int after loop
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create database get
< ok 0 0
> exec ^
: string
| create function f() returns int not deterministic modifies sql data return 1
< error
| 1064 42000 …line 1 column 15 near "function f() returns int not deterministic modifies sql data return 1" 
> exec ^
: string
| create trigger before before insert on a for each row set @a = 1
< error
| 1064 42000 …line 1 column 14 near "trigger before before insert on a for each row set @a = 1" 
> exec ^
: string
| create trigger each after insert on a for each row set @a = 1
< error
| 1064 42000 …line 1 column 14 near "trigger each after insert on a for each row set @a = 1" 
> exec ^
: string
| create trigger t before insert on a for each row set @a = 1
< error
| 1064 42000 …line 1 column 14 near "trigger t before insert on a for each row set @a = 1" 
> exec ^
: string
| drop function if exists return
< error
| 1064 42000 …line 1 column 13 near "function if exists return" 
> exec ^
: string
| insert into return values (1)
< error
| 1146 42S02 Table 'test.return' doesn't exist
> exec ^
: string
| insert into return values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (return) values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (return) values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) return
< error
| 1064 42000 …line 1 column 31 near "return" 
> exec ^
: string
| select * from t group by count
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t group by count
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t order by count
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t order by count
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where count > 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where count > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 + count from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1 + count from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a count from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a count from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cast a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cast a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cast from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cast from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cast(a as char) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cast(a as char) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cast(a as signed) return from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cast(a as signed) return from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select cast.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cast.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select char from t
< error
| 1064 42000 …line 1 column 16 near "from t" 
> exec ^
: string
| select count (1) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count (1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count (a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count (a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count + 1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count + 1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count = 1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count = 1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count count from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count count from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(*) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(*) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(*) loop from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(*) loop from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(return) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(return) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select database from t
< error
| 1064 42000 …line 1 column 20 near "from t" 
> exec ^
: string
| select extract from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select extract from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select extract(day from a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select extract(day from a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select extract.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select extract.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select group_concat from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select group_concat from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select group_concat(a order by a separator ',') from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select group_concat(a order by a separator ',') from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select group_concat(a separator ',') from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select group_concat(a separator ',') from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select group_concat(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select group_concat(a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select grouping from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select grouping from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select grouping(a) from t group by a with rollup
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select grouping(a) from t group by a with rollup
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select grouping.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select grouping.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select if from t
< error
| 1064 42000 …line 1 column 14 near "from t" 
> exec ^
: string
| select if(1,2,3)
< ok 0 0
> exec ^
: string
| select if(return, 1, 2)
< error
| 1054 42S22 Unknown column 'return' in 'field list'
> exec ^
: string
| select if(return, 1, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select insert from t
< error
| 1064 42000 …line 1 column 18 near "from t" 
> exec ^
: string
| select left from t
< error
| 1064 42000 …line 1 column 16 near "from t" 
> exec ^
: string
| select max.max from max
< error
| 1146 42S02 Table 'test.max' doesn't exist
> exec ^
: string
| select max.max from max
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select mod from t
< error
| 1064 42000 …line 1 column 15 near "from t" 
> exec ^
: string
| select modifies from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select modifies from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select now
< error
| 1054 42S22 Unknown column 'now' in 'field list'
> exec ^
: string
| select now
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select now() function
< ok 0 0
> exec ^
: string
| select now()
< ok 0 0
> exec ^
: string
| select now(), system
< error
| 1054 42S22 Unknown column 'system' in 'field list'
> exec ^
: string
| select now(), system
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select now.a
< error
| 1054 42S22 Unknown column 'now.a' in 'field list'
> exec ^
: string
| select now.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select position from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select position from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select position(a in b) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select position(a in b) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select rank from t
< error
| 1064 42000 …line 1 column 16 near "from t" 
> exec ^
: string
| select rank() over ()
< ok 0 0
> exec ^
: string
| select rank.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select rank.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select repeat from t
< error
| 1064 42000 …line 1 column 18 near "from t" 
> exec ^
: string
| select replace from t
< error
| 1064 42000 …line 1 column 19 near "from t" 
> exec ^
: string
| select right from t
< error
| 1064 42000 …line 1 column 17 near "from t" 
> exec ^
: string
| select row_number from t
< error
| 1064 42000 …line 1 column 22 near "from t" 
> exec ^
: string
| select specific from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select specific from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select substring from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select substring from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select substring(a from 1) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select substring(a from 1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select substring(a, 1) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select substring(a, 1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sum from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sum(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sum(a) over (partition by loop) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) over (partition by loop) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sum(a) over return from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) over return from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select sum(a) over w from t window return as (order by a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) over w from t window return as (order by a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select trim(both 'a' from b) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select trim(both 'a' from b) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select variance from t group by variance
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select variance from t group by variance
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select variance from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select variance from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select variance(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select variance(a) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select variance.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select variance.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create function f() returns int deterministic reads sql data return 1
< error
| 1064 42000 …line 1 column 15 near "function f() returns int deterministic reads sql data return 1" 
> exec ^
: string
| create function f() returns int return 1
< error
| 1064 42000 …line 1 column 15 near "function f() returns int return 1" 
> exec ^
: string
| create function function() returns int return 1
< error
| 1064 42000 …line 1 column 15 near "function function() returns int return 1" 
> exec ^
: string
| create function return() returns int return 1
< error
| 1064 42000 …line 1 column 15 near "function return() returns int return 1" 
> exec ^
: string
| create procedure p() begin a: loop select 1; end loop a; end
< error
| 1064 42000 …line 1 column 34 near "loop select 1; end loop a; end" 
> exec ^
: string
| create procedure p() begin declare a condition for 1; end
< error
| 1064 42000 …line 1 column 46 near "condition for 1; end" 
> exec ^
: string
| create procedure p() begin declare a int; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin declare a int; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin declare c asensitive cursor for select 1; end
< error
| 1064 42000 …line 1 column 47 near "asensitive cursor for select 1; end" 
> exec ^
: string
| create procedure p() begin declare c cursor for select 1; declare c2 insensitive cursor for select 1; end
< error
| 1064 42000 …line 1 column 80 near "insensitive cursor for select 1; end" 
> exec ^
: string
| create procedure p() begin declare c sensitive cursor for select 1; end
< error
| 1064 42000 …line 1 column 46 near "sensitive cursor for select 1; end" 
> exec ^
: string
| create procedure p() begin declare exit handler for sqlexception set @a = 1; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin declare exit handler for sqlexception set @a = 1; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin declare loop int; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin declare loop int; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin loop = 1; end
< error
| 1064 42000 …line 1 column 33 near "= 1; end" 
> exec ^
: string
| create procedure p() begin loop select 1; end loop; end
< error
| 1064 42000 …line 1 column 38 near "select 1; end loop; end" 
> exec ^
: string
| create procedure p() begin loop: loop select 1; end loop; end
< error
| 1064 42000 …line 1 column 37 near "loop select 1; end loop; end" 
> exec ^
: string
| create procedure p() begin resignal; end
< error
| 1064 42000 …line 1 column 36 near "; end" 
> exec ^
: string
| create procedure p() begin select loop; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin select loop; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin set loop = 1; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin set loop = 1; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin signal sqlstate '45000'; end
< error
| 1064 42000 …line 1 column 42 near "sqlstate '45000'; end" 
> exec ^
: string
| declare a int
< error
| 1064 42000 …line 1 column 7 near "declare a int" 
> exec ^
: string
| declare c insensitive cursor for select 1
< error
| 1064 42000 …line 1 column 7 near "declare c insensitive cursor for select 1" 
> exec ^
: string
| declare
< error
| 1064 42000 …line 1 column 7 near "declare" 
> exec ^
: string
| drop function function
< error
| 1064 42000 …line 1 column 13 near "function function" 
> exec ^
: string
| drop function return
< error
| 1064 42000 …line 1 column 13 near "function return" 
> exec ^
: string
| grant execute on function f to u
< ok 0 0
> exec ^
: string
| grant execute on function function to u
< ok 0 0
> exec ^
: string
| grant select on function a to u
< ok 0 0
> exec ^
: string
| loop
< error
| 1064 42000 …line 1 column 4 near "loop" 
> exec ^
: string
| resignal sqlstate '45000'
< error
| 1064 42000 …line 1 column 8 near "resignal sqlstate '45000'" 
> exec ^
: string
| resignal
< error
| 1064 42000 …line 1 column 8 near "resignal" 
> exec ^
: string
| return 1
< error
| 1064 42000 …line 1 column 6 near "return 1" 
> exec ^
: string
| return
< error
| 1064 42000 …line 1 column 6 near "return" 
> exec ^
: string
| revoke execute on function f from u
< ok 0 0
> exec ^
: string
| set read_write = 1
< ok 0 0
> exec ^
: string
| set transaction read write
< ok 0 0
> exec ^
: string
| show create function f
< error
| 1064 42000 …line 1 column 20 near "function f" 
> exec ^
: string
| show create function function
< error
| 1064 42000 …line 1 column 20 near "function function" 
> exec ^
: string
| show function code f
< error
| 1064 42000 …line 1 column 18 near "code f" 
> exec ^
: string
| show function status
< ok 0 0
> exec ^
: string
| signal = 1
< error
| 1064 42000 …line 1 column 6 near "signal = 1" 
> exec ^
: string
| signal sqlstate '45000'
< error
| 1064 42000 …line 1 column 6 near "signal sqlstate '45000'" 
> exec ^
: string
| signal
< error
| 1064 42000 …line 1 column 6 near "signal" 
> exec ^
: string
| start transaction read write
< ok 0 0
> exec ^
: string
| alter table t add column a int default b
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| alter table t add column a timestamp default now
< error
| 1064 42000 …line 1 column 48 near "now" 
> exec ^
: string
| alter table t alter a set default (now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t alter a set default (now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change a b int default c
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| alter table t modify a int default b
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| create table t (a int default 'x')
< ok 0 0
> exec ^
: string
| create table t (a int default (now()))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t (a int default (now()))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int default -1)
< ok 0 0
> exec ^
: string
| create table t (a int default 1)
< ok 0 0
> exec ^
: string
| create table t (a int default `b`)
< error
| 1064 42000 …line 1 column 34 near ")" 
> exec ^
: string
| create table t (a int default b)
< error
| 1064 42000 …line 1 column 32 near ")" 
> exec ^
: string
| create table t (a int default b.c)
< error
| 1064 42000 …line 1 column 32 near ".c)" 
> exec ^
: string
| create table t (a int default count)
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec ^
: string
| create table t (a int default current_timestamp())
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t (a int default current_timestamp())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int default now())
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t (a int default now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int default now(1))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table t (a int default now(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int default now)
< error
| 1064 42000 …line 1 column 33 near "now)" 
> exec ^
: string
| create table t (a int default null)
< ok 0 0
> exec ^
: string
| create table t (a int default true)
< ok 0 0
> exec ^
: string
| create table t (a int references b (c) on update cascade)
< ok 0 0
> exec ^
: string
| create table t (a int, foreign key (a) references b (c) on update cascade on delete set null)
< ok 0 0
> exec ^
: string
| create table t (a timestamp default now() on update now())
< ok 0 0
> exec ^
: string
| create table t (a timestamp default now() on update now)
< error
| 1064 42000 …line 1 column 55 near "now)" 
> exec ^
: string
| create table t (a timestamp default now)
< error
| 1064 42000 …line 1 column 39 near "now)" 
> exec ^
: string
| create table t (a timestamp on update b)
< error
| 1064 42000 …line 1 column 39 near "b)" 
> exec ^
: string
| create table t (a timestamp on update current_timestamp())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now)
< error
| 1064 42000 …line 1 column 41 near "now)" 
