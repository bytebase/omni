-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestReservedWordFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select a full from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| set maxvalue = 1
< error
| 1064 42000 …near 'maxvalue = 1' at line 1
> exec ^
: string
| set range = 1
< error
| 1064 42000 …near 'range = 1' at line 1
> exec ^
: string
| set global maxvalue = 1
< ok 0 0
> exec ^
: string
| set global range = 1
< ok 0 0
> exec ^
: string
| select count from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count(*) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select sum from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select sum(a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select now from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select now() from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select now.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a now from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a as now from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select cast from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select cast(a as char) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select extract from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select extract(year from a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select trim(a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select substring from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select substring(a, 1) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select position from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select position('a' in a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select group_concat from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select group_concat(a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select max from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select max(a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select variance from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select variance(a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count + 1 from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 from t where count = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count from t order by count
< error
| 1046 3D000 No database selected
> exec ^
: string
| select t.count from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from count
< error
| 1046 3D000 No database selected
> exec ^
: string
| select if from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select if(1, 2, 3)
< ok 0 0
> exec ^
: string
| select left from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select left(a, 1) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select rank from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select rank() over () from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select lag from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select lag(a) over () from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select row_number from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select row_number() over () from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select first_value from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select first_value(a) over () from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select grouping from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a, grouping(a) from t group by a with rollup
< error
| 1046 3D000 No database selected
> exec ^
: string
| select lateral from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t, lateral (select 1) x
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from lateral
< error
| 1046 3D000 No database selected
> exec ^
: string
| select json_table from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$')) j
< ok 0 0
> exec ^
: string
| select * from json_table
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t window w as ()
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a, rank() over w from t window w as ()
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a window from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select * from t window
< error
| 1064 42000 …near 'window' at line 1
> exec ^
: string
| select 1 from t except select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 from t intersect select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t except select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t intersect select 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a except from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select * from t except
< error
| 1064 42000 …near 'except' at line 1
> exec ^
: string
| select a over from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count(*) over () from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| create database d
< ok 0 0
> exec ^
: string
| drop database d
< ok 0 0
> exec ^
: string
| create schema d
< ok 0 0
> exec ^
: string
| drop schema d
< ok 0 0
> exec ^
: string
| show databases
< ok 0 0
> exec ^
: string
| show schemas
< ok 0 0
> exec ^
: string
| select database()
< ok 0 0
> exec ^
: string
| select database from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select schema from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| create function f() returns int return 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| drop function f
< ok 0 0
> exec ^
: string
| create table function (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (function int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| get diagnostics @a = number
< ok 0 0
> exec ^
: string
| select get from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| create table get (a int)
< error
| 1064 42000 …near 'get (a int)' at line 1
> exec ^
: string
| grant select on t to a with grant option
< error
| 1046 3D000 No database selected
> exec ^
: string
| select option from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| create table option (a int)
< error
| 1064 42000 …near 'option (a int)' at line 1
> exec ^
: string
| create view v as select 1 with check option
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int generated always as (1) stored)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int generated always as (1) virtual)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int as (1) stored)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int as (1) virtual)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (generated int)
< error
| 1064 42000 …near 'generated int)' at line 1
> exec ^
: string
| create table x (stored int)
< error
| 1064 42000 …near 'stored int)' at line 1
> exec ^
: string
| create table x (virtual int)
< error
| 1064 42000 …near 'virtual int)' at line 1
> exec ^
: string
| select system from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a from t for system_time as of now()
< error
| 1064 42000 …near 'system_time as of now()' at line 1
> exec ^
: string
| select * from t for update of t
< error
| 1064 42000 …near 'of t' at line 1
> exec ^
: string
| select of from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table of (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a dec)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a dec(3))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (dec int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int) partition by range (a) (partition p values less than (1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int) partition by range (a) (partition p values less than maxvalue)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table range (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table maxvalue (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select cube from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a from t group by a with cube
< error
| 1064 42000 …near 'cube' at line 1
> exec ^
: string
| select recursive from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| with recursive c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| select empty from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$' default '0' on empty)) j
< ok 0 0
> exec ^
: string
| create table x (a timestamp default now)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| create table x (a timestamp default now())
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a timestamp default current_timestamp)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a timestamp default current_timestamp(6))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int default b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| create table x (a int default 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int default (1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int default (b))
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a timestamp on update now)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| create table x (a timestamp on update now())
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a timestamp on update current_timestamp)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table x (a int on update b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| create table x (a timestamp default count)
< error
| 1064 42000 …near 'count)' at line 1
> exec ^
: string
| select a into from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select * from t into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| select a into @b from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a from t into @b
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 into @a, @b
< ok 0 0
> exec ^
: string
| select a into outfile 'x' from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a into dumpfile 'x' from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select into from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| create role event
< error
| 1064 42000 …near 'event' at line 1
> exec ^
: string
| create role `event`
< ok 0 0
> exec ^
: string
| create role 'event'
< ok 0 0
> exec ^
: string
| create role a, event
< error
| 1064 42000 …near 'event' at line 1
> exec ^
: string
| create role process
< error
| 1064 42000 …near 'process' at line 1
> exec ^
: string
| create role none
< error
| 1064 42000 …near 'none' at line 1
> exec ^
: string
| create role shutdown
< error
| 1064 42000 …near 'shutdown' at line 1
> exec ^
: string
| create role super
< error
| 1064 42000 …near 'super' at line 1
> exec ^
: string
| create role a
< ok 0 0
> exec ^
: string
| create role a, b
< ok 0 0
> exec ^
: string
| create role if not exists event
< error
| 1064 42000 …near 'event' at line 1
> exec ^
: string
| create role if not exists a
< ok 0 0
> exec ^
: string
| create procedure p() close: loop leave close; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() open: loop leave open; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin: loop leave begin; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() add: loop leave add; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() select: loop leave select; end loop
< error
| 1064 42000 …near 'select: loop leave select' at line 1
> exec ^
: string
| create procedure p() l: loop leave l; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin close: loop leave close; end loop; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() l: begin select 1; end l
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare c cursor for select 1; open c; close c; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin begin select 1; end; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; begin end; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| set event = 1
< error
| 1064 42000 …near 'event = 1' at line 1
> exec ^
: string
| set global event = 1
< ok 0 0
> exec ^
: string
| set @a = 1, event = 2
< error
| 1064 42000 …near 'event = 2' at line 1
> exec ^
: string
| set process = 1
< error
| 1064 42000 …near 'process = 1' at line 1
> exec ^
: string
| set none = 1
< error
| 1064 42000 …near 'none = 1' at line 1
> exec ^
: string
| set skip = 1
< error
| 1064 42000 …near 'skip = 1' at line 1
> exec ^
: string
| select skip from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| create table skip (a int)
< error
| 1064 42000 …near 'skip (a int)' at line 1
