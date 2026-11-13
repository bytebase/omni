-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestWildcardItemFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| select * from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select *, a from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select *, a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select s.t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select s.t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a, t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a, t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select *, t.*, a from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select *, t.*, a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select *from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select *from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t . * from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t . * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t. * from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t. * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t .* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t .* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select distinct * from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select distinct * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t union all select t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t union all select t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select (select * from t limit 1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select (select * from t limit 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select exists (select * from t)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select exists (select * from t)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into t select * from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t select * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| with c as (select * from t) select c.* from c
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| with c as (select * from t) select c.* from c
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select off.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select off.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select off .* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select off .* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select ``.t.* from t
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select ``.t.* from t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select ``.``.* from t
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select ``.``.* from t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select `a `.t.* from t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select `a `.t.* from t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select ``.t.a from t
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select ``.t.a from t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select `a/b`.t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select `a/b`.t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select ``.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select ``.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select `*` from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select `*` from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a * b from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a * b from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.a * 2 from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t.a * 2 from t
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
| select count(all *) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select count(all *) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count(*) = 1 from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select count(*) = 1 from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| select *, 1 from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select *, 1 from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select *, t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select *, t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.*, t.* from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select t.*, t.* from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select high_priority *, a from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select high_priority *, a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a, (select * from t limit 1) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a, (select * from t limit 1) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a, exists (select * from t) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a, exists (select * from t) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a, count(*) from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select a, count(*) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select t.*, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select *, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select *, a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select count(*), * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select 1, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select a as x, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select distinct a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select all a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select sql_no_cache a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * from t union select a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select (select a, * from t limit 1)
< error
| 1064 42000 …near '* from t limit 1)' at line 1
> exec ^
: string
| select a from t where exists (select a, * from t)
< error
| 1064 42000 …near '* from t)' at line 1
> exec ^
: string
| select * from (select a, * from t) x
< error
| 1064 42000 …near '* from t) x' at line 1
> exec ^
: string
| with c as (select a, * from t) select * from c
< error
| 1064 42000 …near '* from t) select * from c' at line 1
> exec ^
: string
| insert into t select a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| create view v9 as select a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| create table t9 as select a, * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * as x from t
< error
| 1064 42000 …near 'as x from t' at line 1
> exec ^
: string
| select * x from t
< error
| 1064 42000 …near 'x from t' at line 1
> exec ^
: string
| select t.* x from t
< error
| 1064 42000 …near 'x from t' at line 1
> exec ^
: string
| select * = t
< error
| 1064 42000 …near '= t' at line 1
> exec ^
: string
| select * = 1 from t
< error
| 1064 42000 …near '= 1 from t' at line 1
> exec ^
: string
| select *=1 from t
< error
| 1064 42000 …near '=1 from t' at line 1
> exec ^
: string
| select t.* = 1 from t
< error
| 1064 42000 …near '= 1 from t' at line 1
> exec ^
: string
| select s.t.* = 1 from t
< error
| 1064 42000 …near '= 1 from t' at line 1
> exec ^
: string
| select 1 = * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select 1 = t.* from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * <=> 1 from t
< error
| 1064 42000 …near '<=> 1 from t' at line 1
> exec ^
: string
| select * + 1 from t
< error
| 1064 42000 …near '+ 1 from t' at line 1
> exec ^
: string
| select 1 + * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select 1 + t.* from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * * 2 from t
< error
| 1064 42000 …near '* 2 from t' at line 1
> exec ^
: string
| select 2 * * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * div 2 from t
< error
| 1064 42000 …near 'div 2 from t' at line 1
> exec ^
: string
| select - * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select -t.* from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select not * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * and 1 from t
< error
| 1064 42000 …near 'and 1 from t' at line 1
> exec ^
: string
| select 1 and t.* from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select * or 1 from t
< error
| 1064 42000 …near 'or 1 from t' at line 1
> exec ^
: string
| select * is null from t
< error
| 1064 42000 …near 'is null from t' at line 1
> exec ^
: string
| select t.* is null from t
< error
| 1064 42000 …near 'is null from t' at line 1
> exec ^
: string
| select * in (1) from t
< error
| 1064 42000 …near 'in (1) from t' at line 1
> exec ^
: string
| select * between 1 and 2 from t
< error
| 1064 42000 …near 'between 1 and 2 from t' at line 1
> exec ^
: string
| select * like 'a' from t
< error
| 1064 42000 …near 'like 'a' from t' at line 1
> exec ^
: string
| select * regexp 'a' from t
< error
| 1064 42000 …near 'regexp 'a' from t' at line 1
> exec ^
: string
| select * collate utf8mb4_bin from t
< error
| 1064 42000 …near 'collate utf8mb4_bin from t' at line 1
> exec ^
: string
| select @a := * from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select (*) from t
< error
| 1064 42000 …near '*) from t' at line 1
> exec ^
: string
| select (t.*) from t
< error
| 1064 42000 …near '*) from t' at line 1
> exec ^
: string
| select case * when 1 then 2 end from t
< error
| 1064 42000 …near '* when 1 then 2 end from t' at line 1
> exec ^
: string
| select case when * then 2 end from t
< error
| 1064 42000 …near '* then 2 end from t' at line 1
> exec ^
: string
| select case when 1 then * end from t
< error
| 1064 42000 …near '* end from t' at line 1
> exec ^
: string
| select case when 1 then 2 else t.* end from t
< error
| 1064 42000 …near '* end from t' at line 1
> exec ^
: string
| select interval * day + d from t
< error
| 1064 42000 …near '* day + d from t' at line 1
> exec ^
: string
| select d + interval * day from t
< error
| 1064 42000 …near '* day from t' at line 1
> exec ^
: string
| select cast(* as char) from t
< error
| 1064 42000 …near '* as char) from t' at line 1
> exec ^
: string
| select cast(t.* as char) from t
< error
| 1064 42000 …near '* as char) from t' at line 1
> exec ^
: string
| select convert(*, char) from t
< error
| 1064 42000 …near '*, char) from t' at line 1
> exec ^
: string
| select convert(* using utf8mb4) from t
< error
| 1064 42000 …near '* using utf8mb4) from t' at line 1
> exec ^
: string
| select count(t.*) from t
< error
| 1064 42000 …near '*) from t' at line 1
> exec ^
: string
| select count(distinct *) from t
< error
| 1064 42000 …near '*) from t' at line 1
> exec ^
: string
| select 1 from t where *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t where t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t where a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t where a = t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t group by *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t group by t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t having *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t order by *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t order by t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| set @a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| set @a = t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| update t set a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| update t set a = t.*
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| insert into t (a) values (*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| insert into t (a) values (t.*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| insert into t (a) values (1) on duplicate key update a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| do *
< error
| 1096 HY000 No tables used
> exec ^
: string
| do *
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| do t.*
< error
| 1051 42S02 Unknown table 't'
> exec ^
: string
| do t.*
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| do *, 1
< error
| 1096 HY000 No tables used
> exec ^
: string
| do *, 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| do 1, t.*
< error
| 1051 42S02 Unknown table 't'
> exec ^
: string
| do 1, t.*
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| do t.*, 1
< error
| 1051 42S02 Unknown table 't'
> exec ^
: string
| do t.*, 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| do 1, *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| do *, *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| do * = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| do 1 = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| do t.* + 1
< error
| 1064 42000 …near '+ 1' at line 1
