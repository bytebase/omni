-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReservedWordFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| create procedure p() after: loop leave after; end loop
< ok 0 0
> exec ^
: string
| create procedure p() absent: loop leave absent; end loop
< ok 0 0
> exec ^
: string
| create procedure p() begin after: loop leave after; end loop; end
< ok 0 0
> exec ^
: string
| create procedure p() database: loop leave database; end loop
< ok 0 0
> exec ^
: string
| create procedure p() clone: loop leave clone; end loop
< ok 0 0
> exec ^
: string
| create procedure p() format: loop leave format; end loop
< error
| 1064 42000 …near 'format: loop leave format; end loop' at line 1
> exec ^
: string
| create procedure p() close: loop leave close; end loop
< error
| 1064 42000 …near ': loop leave close; end loop' at line 1
> exec ^
: string
| create procedure p() get: loop leave get; end loop
< error
| 1064 42000 …near ': loop leave get; end loop' at line 1
> exec ^
: string
| create procedure p() window: loop leave window; end loop
< error
| 1064 42000 …near 'window: loop leave window; end loop' at line 1
> exec ^
: string
| create procedure p() select: loop leave select; end loop
< error
| 1064 42000 …near ': loop leave select; end loop' at line 1
> exec ^
: string
| create procedure p() begin declare c cursor for select 1; open c; close c; end
< ok 0 0
> exec ^
: string
| create procedure p() l: begin select 1; end l
< ok 0 0
> exec ^
: string
| create table x (a timestamp on update now)
< error
| 1064 42000 …near 'now)' at line 1
> exec ^
: string
| create table x (a int on update b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| create table x (a timestamp on update now())
< ok 0 0
> exec ^
: string
| create table x (a timestamp on update current_timestamp)
< ok 0 0
> exec ^
: string
| create table x (a timestamp default now())
< ok 0 0
> exec ^
: string
| create table x (a timestamp default current_timestamp)
< ok 0 0
> exec ^
: string
| select a into from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select * from t into
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select a into @b, from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select a into @b, @c from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a into @b, @c from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t into @b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t into @b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t for portion of p from '2020-01-01' to '2021-01-01' a set b = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t for portion of p from '2020-01-01' to '2021-01-01' a set b = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t for portion of p from '2020-01-01' to '2021-01-01' as a set b = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t for portion of p from '2020-01-01' to '2021-01-01' as a set b = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t for portion of p from '2020-01-01' to '2021-01-01' window set b = 1
< error
| 1064 42000 …near 'window set b = 1' at line 1
> exec ^
: string
| update t for portion of p from '2020-01-01' to '2021-01-01' as window set b = 1
< error
| 1064 42000 …near 'window set b = 1' at line 1
> exec ^
: string
| alter schema d character set utf8
< ok 0 0
> exec ^
: string
| alter table t add column a int after before
< error
| 1064 42000 …near 'before' at line 1
> exec ^
: string
| alter table t add column a int first
< ok 0 0
> exec ^
: string
| alter table t add column a int generated always as (1) stored
< ok 0 0
> exec ^
: string
| alter table t add column before int
< error
| 1064 42000 …near 'before int' at line 1
> exec ^
: string
| alter table t add system versioning
< ok 0 0
> exec ^
: string
| alter table t drop system versioning
< ok 0 0
> exec ^
: string
| alter table t modify a int after loop
< error
| 1064 42000 …near 'loop' at line 1
> exec ^
: string
| change master to master_ssl_verify_server_cert = 1
< ok 0 0
> exec ^
: string
| create database d
< ok 0 0
> exec ^
: string
| create database get
< ok 0 0
> exec ^
: string
| create function f() returns int deterministic reads sql data return 1
< ok 0 0
> exec ^
: string
| create function f() returns int not deterministic modifies sql data return 1
< ok 0 0
> exec ^
: string
| create function f() returns int return 1
< ok 0 0
> exec ^
: string
| create function function() returns int return 1
< ok 0 0
> exec ^
: string
| create function return() returns int return 1
< error
| 1064 42000 …near 'return() returns int return 1' at line 1
> exec ^
: string
| create index before on t (a)
< error
| 1064 42000 …near 'before on t (a)' at line 1
> exec ^
: string
| create procedure p() begin a: loop select 1; end loop a; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a condition for 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare exit handler for sqlexception set @a = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare loop int; end
< error
| 1064 42000 …near 'loop int; end' at line 1
> exec ^
: string
| create procedure p() begin get diagnostics @a = number; end
< ok 0 0
> exec ^
: string
| create procedure p() begin loop = 1; end
< error
| 1064 42000 …near '= 1; end' at line 1
> exec ^
: string
| create procedure p() begin loop select 1; end loop; end
< ok 0 0
> exec ^
: string
| create procedure p() begin loop: loop select 1; end loop; end
< error
| 1064 42000 …near ': loop select 1; end loop; end' at line 1
> exec ^
: string
| create procedure p() begin resignal; end
< ok 0 0
> exec ^
: string
| create function f() returns int begin return 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin return 1; end
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| create procedure p() begin return 1; end
| ; SELECT FROM
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| create procedure p() begin select loop; end
< error
| 1064 42000 …near 'loop; end' at line 1
> exec ^
: string
| create procedure p() begin set loop = 1; end
< error
| 1064 42000 …near 'loop = 1; end' at line 1
> exec ^
: string
| create procedure p() begin signal sqlstate '45000'; end
< ok 0 0
> exec ^
: string
| create schema d
< ok 0 0
> exec ^
: string
| create table function (a int)
< ok 0 0
> exec ^
: string
| create table generated (a int)
< ok 0 0
> exec ^
: string
| create table loop (a int)
< error
| 1064 42000 …near 'loop (a int)' at line 1
> exec ^
: string
| create table option (a int)
< ok 0 0
> exec ^
: string
| create table return (signal int, resignal int, get int, purge int, undo int, declare int)
< error
| 1064 42000 …near 'return (signal int, resignal int, get int, purge int, undo int, declare int)' at line 1
> exec ^
: string
| create table schema (a int)
< ok 0 0
> exec ^
: string
| create table system (a int)
< ok 0 0
> exec ^
: string
| create table t (a int as (1) stored)
< ok 0 0
> exec ^
: string
| create table t (a int as (1) virtual)
< ok 0 0
> exec ^
: string
| create table t (a int as (1))
< ok 0 0
> exec ^
: string
| create table t (a int generated always as (1) stored)
< ok 0 0
> exec ^
: string
| create table t (a int generated always as (1) virtual)
< ok 0 0
> exec ^
: string
| create table t (a int generated always as (1))
< ok 0 0
> exec ^
: string
| create table t (a int generated)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| create table t (a int without system versioning)
< ok 0 0
> exec ^
: string
| create table t (a int) with system versioning
< ok 0 0
> exec ^
: string
| create table t (a int) without system versioning
< error
| 1064 42000 …near 'without system versioning' at line 1
> exec ^
: string
| create table t (a int, constraint signal unique (a))
< error
| 1064 42000 …near 'signal unique (a))' at line 1
> exec ^
: string
| create table t (a int, index return (a))
< error
| 1064 42000 …near 'return (a))' at line 1
> exec ^
: string
| create table t (generated int, stored int, virtual int)
< ok 0 0
> exec ^
: string
| create table t (window int)
< ok 0 0
> exec ^
: string
| create trigger before before insert on a for each row set @a = 1
< error
| 1064 42000 …near 'before before insert on a for each row set @a = 1' at line 1
> exec ^
: string
| create trigger each after insert on a for each row set @a = 1
< error
| 1064 42000 …near 'each after insert on a for each row set @a = 1' at line 1
> exec ^
: string
| create trigger t before insert on a for each row set @a = 1
< ok 0 0
> exec ^
: string
| create user declare
< error
| 1064 42000 …near 'declare' at line 1
> exec ^
: string
| create view loop as select 1
< error
| 1064 42000 …near 'loop as select 1' at line 1
> exec ^
: string
| deallocate prepare signal
< error
| 1064 42000 …near 'signal' at line 1
> exec ^
: string
| delete from declare
< error
| 1064 42000 …near 'declare' at line 1
> exec ^
: string
| delete from t as window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| delete from t window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| delete from window
< error
| 1146 42S02 Table 'test.window' doesn't exist
> exec ^
: string
| delete from window
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| drop database d
< ok 0 0
> exec ^
: string
| drop function f
< ok 0 0
> exec ^
: string
| drop function function
< ok 0 0
> exec ^
: string
| drop function if exists f
< ok 0 0
> exec ^
: string
| drop function if exists return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| drop function return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| drop schema d
< ok 0 0
> exec ^
: string
| execute signal
< error
| 1064 42000 …near 'signal' at line 1
> exec ^
: string
| get diagnostics @a = number
< ok 0 0
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
| grant select on t to u with grant option
< ok 0 0
> exec ^
: string
| insert into return values (1)
< error
| 1064 42000 …near 'return values (1)' at line 1
> exec ^
: string
| insert into t (return) values (1)
< error
| 1064 42000 …near 'return) values (1)' at line 1
> exec ^
: string
| insert into t values (1) return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| insert into window values (1)
< error
| 1146 42S02 Table 'test.window' doesn't exist
> exec ^
: string
| insert into window values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| prepare signal from 'select 1'
< error
| 1064 42000 …near 'signal from 'select 1'' at line 1
> exec ^
: string
| revoke execute on function f from u
< ok 0 0
> exec ^
: string
| rollback to loop
< error
| 1064 42000 …near 'loop' at line 1
> exec ^
: string
| rollback to savepoint return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| savepoint loop
< error
| 1064 42000 …near 'loop' at line 1
> exec ^
: string
| select 'x' return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| select (1) signal
< error
| 1064 42000 …near 'signal' at line 1
> exec ^
: string
| select * from (select 1) a window w as ()
< ok 0 0
> exec ^
: string
| select * from (select 1) as window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| select * from (select 1) return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| select * from (select 1) window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| select * from a.window
< error
| 1146 42S02 Table 'a.window' doesn't exist
> exec ^
: string
| select * from a.window
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$')) t
< ok 0 0
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t a join loop on a.b = loop.c
< error
| 1064 42000 …near 'loop on a.b = loop.c' at line 1
> exec ^
: string
| select * from t a join return using (b)
< error
| 1064 42000 …near 'return using (b)' at line 1
> exec ^
: string
| select * from t a straight_join signal on 1
< error
| 1064 42000 …near 'signal on 1' at line 1
> exec ^
: string
| select * from t a window w as ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a window w as ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as rank
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as rank
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| select * from t for system_time all
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t for system_time all
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t for system_time as of 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t for system_time as of 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t for update of empty
< error
| 1064 42000 …near 'of empty' at line 1
> exec ^
: string
| select * from t for update of of
< error
| 1064 42000 …near 'of of' at line 1
> exec ^
: string
| select * from t for update of t
< error
| 1064 42000 …near 'of t' at line 1
> exec ^
: string
| select * from t for update return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| select * from t full
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t full
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t full join b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t full join b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t function
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t function
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t get
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t get
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t group by loop
< error
| 1064 42000 …near 'loop' at line 1
> exec ^
: string
| select * from t groups
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t groups
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t lateral
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t lateral
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t limit 1 loop
< error
| 1064 42000 …near 'loop' at line 1
> exec ^
: string
| select * from t of
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t of
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t order by return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| select * from t rank
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t rank
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t system
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t system
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t where a = 1 return
< error
| 1064 42000 …near 'return' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t where each = 1
< error
| 1064 42000 …near 'each = 1' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t window
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select * from t window w as ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t window w as ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from window
< error
| 1146 42S02 Table 'test.window' doesn't exist
> exec ^
: string
| select * from window
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from window w
< error
| 1146 42S02 Table 'test.window' doesn't exist
> exec ^
: string
| select * from window w
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from window, t
< error
| 1146 42S02 Table 'test.window' doesn't exist
> exec ^
: string
| select * from window, t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from window, window
< error
| 1066 42000 Not unique table/alias: 'window'
> exec ^
: string
| select * from window, window
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'window'
> exec ^
: string
| select 1 + 1 get
< ok 0 0
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 from t tablesample system
< error
| 1064 42000 …near 'system' at line 1
> exec ^
: string
| select 1 from t tablesample system (1)
< error
| 1064 42000 …near 'system (1)' at line 1
> exec ^
: string
| select 1 loop
< error
| 1064 42000 …near 'loop' at line 1
> exec ^
: string
| select 1 signal, 2 resignal
< error
| 1064 42000 …near 'signal, 2 resignal' at line 1
> exec ^
: string
| select 1 window w as ()
< error
| 1064 42000 …near 'w as ()' at line 1
> exec ^
: string
| select @a := return
< error
| 1064 42000 …near 'return' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a as window from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a as window from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a before from t
< error
| 1064 42000 …near 'before from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t a full join b on 1
< error
| 1064 42000 …near 'full join b on 1' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t full outer join b
< error
| 1064 42000 …near 'outer join b' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t into @b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t into @b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t join b as window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| select a from t join b window
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select a from t window w as ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t window w as ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a full
< error
| 1054 42S22 Unknown column 'a' in 'SELECT'
> exec ^
: string
| select a full
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a into @b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a into @b from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a loop from t
< error
| 1064 42000 …near 'loop from t' at line 1
> exec ^
: string
| select a rank from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a rank from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a return from t
< error
| 1064 42000 …near 'return from t' at line 1
> exec ^
: string
| select a signal from t
< error
| 1064 42000 …near 'signal from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a window from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a window from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a window w as ()
< error
| 1064 42000 …near 'w as ()' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select accessible from t
< error
| 1064 42000 …near 'accessible from t' at line 1
> exec ^
: string
| select asensitive from t
< error
| 1064 42000 …near 'asensitive from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select cast(a as signed) return from t
< error
| 1064 42000 …near 'return from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select char from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select condition from t
< error
| 1064 42000 …near 'condition from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count(*) loop from t
< error
| 1064 42000 …near 'loop from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count(return) from t
< error
| 1064 42000 …near 'return) from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select cube from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cube from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select cube(1)
< error
| 1630 42000 FUNCTION test.cube does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select cube(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select cume_dist() over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cume_dist() over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select database
< error
| 1054 42S22 Unknown column 'database' in 'SELECT'
> exec ^
: string
| select database
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select database from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select database from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select database()
< ok 0 0
> exec ^
: string
| select declare from t
< error
| 1064 42000 …near 'declare from t' at line 1
> exec ^
: string
| select dense_rank() over (order by a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select dense_rank() over (order by a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select first_value(a) over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select first_value(a) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select generated from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select generated from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select get
< error
| 1054 42S22 Unknown column 'get' in 'SELECT'
> exec ^
: string
| select get
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select groups from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select groups from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select if from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select if(1,2,3)
< ok 0 0
> exec ^
: string
| select if(return, 1, 2)
< error
| 1064 42000 …near 'return, 1, 2)' at line 1
> exec ^
: string
| select insensitive from t
< error
| 1064 42000 …near 'insensitive from t' at line 1
> exec ^
: string
| select insert from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select into from t
< error
| 1064 42000 …near 'into from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select lag(a) over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select lag(a) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select lead(a, 1) over (order by a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select lead(a, 1) over (order by a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select left from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select loop(1)
< error
| 1064 42000 …near 'loop(1)' at line 1
> exec ^
: string
| select master_bind from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select master_bind from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select master_ssl_verify_server_cert from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select master_ssl_verify_server_cert from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select max from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select min(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select min(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select mod from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select modifies from t
< error
| 1064 42000 …near 'modifies from t' at line 1
> exec ^
: string
| select now
< error
| 1054 42S22 Unknown column 'now' in 'SELECT'
> exec ^
: string
| select now
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now (1)
< error
| 1630 42000 FUNCTION test.now does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select now (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now (6) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select now (6) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now()
< ok 0 0
> exec ^
: string
| select now() function
< ok 0 0
> exec ^
: string
| select now(), system
< error
| 1054 42S22 Unknown column 'system' in 'SELECT'
> exec ^
: string
| select now(), system
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now.a
< error
| 1109 42S02 Unknown table 'now' in SELECT
> exec ^
: string
| select now.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select nth_value(a, 1) over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select nth_value(a, 1) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select ntile(2) over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select ntile(2) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select of from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select of from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select option from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select option from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select percent_rank() over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select percent_rank() over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select rank from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select rank from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select rank() over ()
< error
| 4018 HY000 No order list in window specification for 'rank'
> exec ^
: string
| select rank() over ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select rank() over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select rank() over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select read_write from t
< error
| 1064 42000 …near 'read_write from t' at line 1
> exec ^
: string
| select reads from t
< error
| 1064 42000 …near 'reads from t' at line 1
> exec ^
: string
| select repeat from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select replace from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select return(1)
< error
| 1064 42000 …near 'return(1)' at line 1
> exec ^
: string
| select return.a from return
< error
| 1064 42000 …near 'return' at line 1
> exec ^
: string
| select right from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select row_number from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select row_number() over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select row_number() over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select schema
< error
| 1054 42S22 Unknown column 'schema' in 'SELECT'
> exec ^
: string
| select schema
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select schema()
< ok 0 0
> exec ^
: string
| select sensitive from t
< error
| 1064 42000 …near 'sensitive from t' at line 1
> exec ^
: string
| select signal(1)
< error
| 1064 42000 …near 'signal(1)' at line 1
> exec ^
: string
| select specific from t
< error
| 1064 42000 …near 'specific from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sql_buffer_result from t
< error
| 1064 42000 …near 'from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sql_no_cache
< error
| 1064 42000 …near '' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sql_no_cache from t
< error
| 1064 42000 …near 'from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select stored from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select stored from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum(a) over (partition by loop) from t
< error
| 1064 42000 …near 'loop) from t' at line 1
> exec ^
: string
| select sum(a) over (rows between 1 preceding and current row) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) over (rows between 1 preceding and current row) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum(a) over return from t
< error
| 1064 42000 …near 'return from t' at line 1
> exec ^
: string
| select sum(a) over w from t window return as (order by a)
< error
| 1064 42000 …near 'return as (order by a)' at line 1
> exec ^
: string
| select sum(a) over w from t window w as (partition by b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) over w from t window w as (partition by b)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum(a) over window from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select sum(a) over window from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select system
< error
| 1054 42S22 Unknown column 'system' in 'SELECT'
> exec ^
: string
| select system
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select system(1)
< error
| 1630 42000 FUNCTION test.system does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| select system(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select variance (1) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select variance (1) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select virtual from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select virtual from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select window from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select window from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select window.a from window
< error
| 1146 42S02 Table 'test.window' doesn't exist
> exec ^
: string
| select window.a from window
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = system
< error
| 1054 42S22 Unknown column 'system' in 'SET'
> exec ^
: string
| set @a = system
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set autocommit = system
< ok 0 0
> exec ^
: string
| set autocommit = system + 1
< error
| 1054 42S22 Unknown column 'system' in 'SET'
> exec ^
: string
| set autocommit = system + 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set read_write = 1
< error
| 1064 42000 …near 'read_write = 1' at line 1
> exec ^
: string
| set transaction read write
< ok 0 0
> exec ^
: string
| show create database d
< ok 0 0
> exec ^
: string
| show create function f
< ok 0 0
> exec ^
: string
| show create function function
< ok 0 0
> exec ^
: string
| show create schema d
< ok 0 0
> exec ^
: string
| show databases
< ok 0 0
> exec ^
: string
| show function code f
< error
| 1289 HY000 The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working
> exec ^
: string
| show function status
< ok 0 0
> exec ^
: string
| show schemas
< ok 0 0
> exec ^
: string
| start transaction read write
< ok 0 0
> exec ^
: string
| update return set loop = 1
< error
| 1064 42000 …near 'return set loop = 1' at line 1
> exec ^
: string
| update t as window set a = 1
< error
| 1064 42000 …near 'window set a = 1' at line 1
> exec ^
: string
| update t window set a = 1
< error
| 1064 42000 …near 'window set a = 1' at line 1
> exec ^
: string
| use database
< ok 0 0
> exec ^
: string
| use schema
< ok 0 0
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
| with window as (select 1) select * from window
< ok 0 0
