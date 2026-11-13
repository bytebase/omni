-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestReservedWordFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| select a full from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a full from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select count(*) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select sum from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select sum(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select now from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now() from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select now() from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now.a from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select now.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a now from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a now from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a as now from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a as now from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select cast from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select cast(a as char) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select extract from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select extract from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select extract(year from a) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select extract(year from a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select trim from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select trim(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substring from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select substring from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select substring(a, 1) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select substring(a, 1) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select position from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select position from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select position('a' in a) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select position('a' in a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select group_concat from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select group_concat(a) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select group_concat(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select max from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select max(a) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select max(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select variance from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select variance from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select variance(a) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select variance(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count + 1 from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select count + 1 from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 from t where count = 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select 1 from t where count = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count from t order by count
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select count from t order by count
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.count from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t.count from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from count
< error
| 1146 42S02 Table 'omni_t_mysql_parser.count' doesn't exist
> exec ^
: string
| select * from count
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select left(a, 1) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select rank from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select rank() over () from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select rank() over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select lag from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select lag(a) over () from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select lag(a) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select row_number from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select row_number() over () from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select row_number() over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select first_value from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select first_value(a) over () from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select first_value(a) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select grouping from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select a, grouping(a) from t group by a with rollup
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a, grouping(a) from t group by a with rollup
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select lateral from t
< error
| 1064 42000 …near 'lateral from t' at line 1
> exec ^
: string
| select * from t, lateral (select 1) x
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t, lateral (select 1) x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from lateral
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select json_table from t
< error
| 1064 42000 …near 'json_table from t' at line 1
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$')) j
< ok 0 0
> exec ^
: string
| select * from json_table
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select * from t window w as ()
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t window w as ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a, rank() over w from t window w as ()
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a, rank() over w from t window w as ()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a window from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select * from t window
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select 1 from t except select 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select 1 from t except select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 from t intersect select 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select 1 from t intersect select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t except select 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t except select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t intersect select 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t intersect select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a except from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select * from t except
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select a over from t
< error
| 1064 42000 …near 'over from t' at line 1
> exec ^
: string
| select count(*) over () from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select count(*) over () from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create function f() returns int return 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| drop function f
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| drop function f
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table function (a int)
< error
| 1064 42000 …near 'function (a int)' at line 1
> exec ^
: string
| create table x (function int)
< error
| 1064 42000 …near 'function int)' at line 1
> exec ^
: string
| get diagnostics @a = number
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| get diagnostics @a = number
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select get from t
< error
| 1064 42000 …near 'get from t' at line 1
> exec ^
: string
| create table get (a int)
< error
| 1064 42000 …near 'get (a int)' at line 1
> exec ^
: string
| grant select on t to a with grant option
< ok 0 0
> exec ^
: string
| select option from t
< error
| 1064 42000 …near 'option from t' at line 1
> exec ^
: string
| create table option (a int)
< error
| 1064 42000 …near 'option (a int)' at line 1
> exec ^
: string
| create view v as select 1 with check option
< ok 0 0
> exec ^
: string
| create table x (a int generated always as (1) stored)
< ok 0 0
> exec ^
: string
| create table x (a int generated always as (1) virtual)
< ok 0 0
> exec ^
: string
| create table x (a int as (1) stored)
< ok 0 0
> exec ^
: string
| create table x (a int as (1) virtual)
< ok 0 0
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
| 1064 42000 …near 'system from t' at line 1
> exec ^
: string
| select a from t for system_time as of now()
< error
| 1064 42000 …near 'system_time as of now()' at line 1
> exec ^
: string
| select * from t for update of t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t for update of t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select of from t
< error
| 1064 42000 …near 'of from t' at line 1
> exec ^
: string
| create table of (a int)
< error
| 1064 42000 …near 'of (a int)' at line 1
> exec ^
: string
| create table x (a dec)
< ok 0 0
> exec ^
: string
| create table x (a dec(3))
< ok 0 0
> exec ^
: string
| create table x (dec int)
< error
| 1064 42000 …near 'dec int)' at line 1
> exec ^
: string
| create table x (a int) partition by range (a) (partition p values less than (1))
< ok 0 0
> exec ^
: string
| create table x (a int) partition by range (a) (partition p values less than maxvalue)
< ok 0 0
> exec ^
: string
| create table range (a int)
< error
| 1064 42000 …near 'range (a int)' at line 1
> exec ^
: string
| create table maxvalue (a int)
< error
| 1064 42000 …near 'maxvalue (a int)' at line 1
> exec ^
: string
| select cube from t
< error
| 1064 42000 …near 'cube from t' at line 1
> exec ^
: string
| select a from t group by a with cube
< error
| 1064 42000 …near 'with cube' at line 1
> exec ^
: string
| select recursive from t
< error
| 1064 42000 …near 'recursive from t' at line 1
> exec ^
: string
| with recursive c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| select empty from t
< error
| 1064 42000 …near 'empty from t' at line 1
> exec ^
: string
| select empty(a) from t
< error
| 1064 42000 …near 'empty(a) from t' at line 1
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$' default '0' on empty)) j
< ok 0 0
> exec ^
: string
| create table x (a timestamp default now)
< error
| 1064 42000 …near 'now)' at line 1
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
| create table x (a timestamp default current_timestamp(6))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| create table x (a timestamp default current_timestamp(6))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table x (a int default b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| create table x (a int default 1)
< ok 0 0
> exec ^
: string
| create table x (a int default (1))
< ok 0 0
> exec ^
: string
| create table x (a int default (b))
< ok 0 0
> exec ^
: string
| create table x (a timestamp on update now)
< error
| 1064 42000 …near 'now)' at line 1
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
| 1064 42000 …near '' at line 1
> exec ^
: string
| select a into @b from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a into @b from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t into @b
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a from t into @b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 into @a, @b
< error
| 1222 21000 The used SELECT statements have a different number of columns
> exec ^
: string
| select 1 into @a, @b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a into outfile 'x' from t
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| select a into outfile 'x' from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a into dumpfile 'x' from t
< error
| 1227 42000 Access denied; you need (at least one of) the FILE privilege(s) for this operation
> exec ^
: string
| select a into dumpfile 'x' from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select into from t
< error
| 1064 42000 …near 'into from t' at line 1
> exec ^
: string
| create role event
< error
| 1064 42000 …near 'event' at line 1
> exec ^
: string
| create role `event`
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create role `event`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create role 'event'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create role 'event'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create role a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create role a, b
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create role a, b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create role if not exists event
< error
| 1064 42000 …near 'event' at line 1
> exec ^
: string
| create role if not exists a
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create role if not exists a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() close: loop leave close; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() close: loop leave close; end loop
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() open: loop leave open; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() open: loop leave open; end loop
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin: loop leave begin; end loop
< error
| 1064 42000 …near ': loop leave begin; end loop' at line 1
> exec ^
: string
| create procedure p() add: loop leave add; end loop
< error
| 1064 42000 …near 'add: loop leave add; end loop' at line 1
> exec ^
: string
| create procedure p() select: loop leave select; end loop
< error
| 1064 42000 …near ': loop leave select; end loop' at line 1
> exec ^
: string
| create procedure p() l: loop leave l; end loop
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() l: loop leave l; end loop
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin close: loop leave close; end loop; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin close: loop leave close; end loop; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() l: begin select 1; end l
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() l: begin select 1; end l
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare c cursor for select 1; open c; close c; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare c cursor for select 1; open c; close c; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin begin select 1; end; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin begin select 1; end; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; begin end; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; begin end; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() end: loop leave end; end loop
< error
| 1064 42000 …near 'end: loop leave end; end loop' at line 1
> exec ^
: string
| select skip from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select skip from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table skip (a int)
< ok 0 0
