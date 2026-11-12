-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestQueryFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| table t order by a
< error
| 1064 42000 …near 'table t order by a' at line 1
> exec ^
: string
| table t limit 1
< error
| 1064 42000 …near 'table t limit 1' at line 1
> exec ^
: string
| table t into @x
< error
| 1064 42000 …near 'table t into @x' at line 1
> exec ^
: string
| values row(1)
< ok 0 0
> exec ^
: string
| (table t)
< error
| 1064 42000 …near 'table t)' at line 1
> exec ^
: string
| ((table t))
< error
| 1064 42000 …near 'table t))' at line 1
> exec ^
: string
| (values row(1))
< ok 0 0
> exec ^
: string
| select (table t)
< error
| 1064 42000 …near 't)' at line 1
> exec ^
: string
| select ((table t) = 1)
< error
| 1064 42000 …near 't) = 1)' at line 1
> exec ^
: string
| select exists (table t)
< error
| 1064 42000 …near 'table t)' at line 1
> exec ^
: string
| select * from t where a in (table t)
< error
| 1064 42000 …near 't)' at line 1
> exec ^
: string
| select * from t where a = any (table t)
< error
| 1064 42000 …near 't)' at line 1
> exec ^
: string
| select * from t where a in (select a from t)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from (table t) d
< error
| 1064 42000 …near 't) d' at line 1
> exec ^
: string
| select * from t join (table u) d on 1
< error
| 1064 42000 …near 'u) d on 1' at line 1
> exec ^
: string
| select * from (values row(1)) d
< ok 0 0
> exec ^
: string
| with c as (table t) select * from c
< error
| 1064 42000 …near 'table t) select * from c' at line 1
> exec ^
: string
| with c as (select 1) table c
< error
| 1064 42000 …near 'table c' at line 1
> exec ^
: string
| with c as (values row(1)) select * from c
< ok 0 0
> exec ^
: string
| select 1 union table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| table t union select 1
< error
| 1064 42000 …near 'table t union select 1' at line 1
> exec ^
: string
| (select 1) union (table t)
< error
| 1064 42000 …near 'table t)' at line 1
> exec ^
: string
| values row(1) union table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| select 1 union values row(1)
< ok 0 0
> exec ^
: string
| (values row(1)) union (select 1)
< ok 0 0
> exec ^
: string
| insert into t table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t (table u)
< error
| 1064 42000 …near 'u)' at line 1
> exec ^
: string
| insert into t select 1 union table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t select 1 union select 2
< error
| 1046 3D000 No database selected
> exec ^
: string
| replace into t table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| explain table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| explain (table t)
< error
| 1064 42000 …near 'table t)' at line 1
> exec ^
: string
| explain select 1 union table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| explain select 1 union select 2
< ok 0 0
> exec ^
: string
| create view v as table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| create view v as select 1 union table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| create view v as select 1 union select 2
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table n as table t
< error
| 1064 42000 …near 'table t' at line 1
> exec ^
: string
| create table n (table t)
< error
| 1064 42000 …near 'table t)' at line 1
> exec ^
: string
| create table n as values row(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create procedure p() begin table t; end
< error
| 1064 42000 …near 'table' at line 1
> exec ^
: string
| create procedure p() begin declare c cursor for table t; end
< error
| 1064 42000 …near 'table' at line 1
> exec ^
: string
| create procedure p() begin declare c cursor for select 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| (select 1) lock in share mode
< error
| 1064 42000 …near 'lock in share mode' at line 1
> exec ^
: string
| (select 1) for update nowait
< error
| 1064 42000 …near 'for update nowait' at line 1
> exec ^
: string
| (select 1) for update wait 1
< error
| 1064 42000 …near 'for update wait 1' at line 1
> exec ^
: string
| (select 1) for update skip locked
< error
| 1064 42000 …near 'for update skip locked' at line 1
> exec ^
: string
| ((select 1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| ((select 1) order by 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| ((select 1) for update)
< error
| 1064 42000 …near 'for update)' at line 1
> exec ^
: string
| (select 1 for update) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| ((select 1) union (select 2)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| (select 1 union select 2) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| with c as (select 1) (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| select 1 for update
< ok 0 0
> exec ^
: string
| (select 1 for update)
< ok 0 0
> exec ^
: string
| (select 1) order by 1 for update
< ok 0 0
> exec ^
: string
| (select 1) limit 1 for update
< ok 0 0
> exec ^
: string
| (select 1) order by 1 limit 1 for update
< ok 0 0
> exec ^
: string
| (select 1) order by 1 lock in share mode
< ok 0 0
> exec ^
: string
| (select 1) order by 1 for update nowait
< ok 0 0
> exec ^
: string
| (select 1) limit 1 for update wait 1
< ok 0 0
> exec ^
: string
| (select 1) order by 1 for update skip locked
< ok 0 0
> exec ^
: string
| ((select 1) order by 1 for update)
< ok 0 0
> exec ^
: string
| (select 1 for update) order by 1 for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) for update
< ok 0 0
> exec ^
: string
| (select 1) union all (select 2) for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) lock in share mode
< ok 0 0
> exec ^
: string
| select 1 union (select 2) for update
< ok 0 0
> exec ^
: string
| (select 1) union select 2 for update
< ok 0 0
> exec ^
: string
| with c as (select 1) (select 1) order by 1 for update
< ok 0 0
> exec ^
: string
| select * from ((select 1) for update) d
< error
| 1064 42000 …near 'for update) d' at line 1
> exec ^
: string
| select * from ((select 1) order by 1 for update) d
< ok 0 0
> exec ^
: string
| select * from ((select 1) union (select 2) for update) d
< ok 0 0
> exec ^
: string
| select ((select 1) for update)
< error
| 1064 42000 …near 'for update)' at line 1
> exec ^
: string
| select (select 1 for update)
< ok 0 0
> exec ^
: string
| select 1 in ((select 1) for update)
< error
| 1064 42000 …near 'for update)' at line 1
> exec ^
: string
| select exists ((select 1) for update)
< error
| 1064 42000 …near 'for update)' at line 1
> exec ^
: string
| select exists ((select 1) order by 1 for update)
< ok 0 0
> exec ^
: string
| insert into t (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t (select 1) order by 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t select 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| explain (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| explain (select 1) order by 1 for update
< ok 0 0
> exec ^
: string
| create view v as (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create view v as (select 1) order by 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create view v as select 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table n (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table n (select 1) order by 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table n as (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table n as (select 1) limit 1 for update
< error
| 1046 3D000 No database selected
> exec ^
: string
| create procedure p() begin (select 1) for update; end
< error
| 1064 42000 …near '(select 1) for update' at line 1
> exec ^
: string
| create procedure p() begin (select 1) order by 1 for update; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare c cursor for (select 1) for update; end
< error
| 1064 42000 …near '(select 1) for update' at line 1
> exec ^
: string
| create procedure p() begin declare c cursor for (select 1) limit 1 for update; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
