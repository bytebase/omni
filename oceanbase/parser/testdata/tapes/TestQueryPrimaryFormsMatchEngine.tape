-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestQueryPrimaryFormsMatchEngine default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_sweep
< ok 1 0
> exec
| USE omni_sweep
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM ?
: string
| table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| table u order by a
< error
| 1064 42000 …near 'table u order by a' at line 1
> exec ^
: string
| table u limit 1
< error
| 1064 42000 …near 'table u limit 1' at line 1
> exec ^
: string
| table u into @x
< error
| 1064 42000 …near 'table u into @x' at line 1
> exec ^
: string
| (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| ((table u))
< error
| 1064 42000 …near 'table u))' at line 1
> exec ^
: string
| explain table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| select * from (table u) v
< error
| 1064 42000 …near 'u) v' at line 1
> exec ^
: string
| select 1 union table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| table u union select 1
< error
| 1064 42000 …near 'table u union select 1' at line 1
> exec ^
: string
| table u union table v
< error
| 1064 42000 …near 'table u union table v' at line 1
> exec ^
: string
| select (table u)
< error
| 1064 42000 …near 'u)' at line 1
> exec ^
: string
| select * from t where a in (table u)
< error
| 1064 42000 …near 'u)' at line 1
> exec ^
: string
| select exists (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| with c as (table u) select * from c
< error
| 1064 42000 …near 'table u) select * from c' at line 1
> exec ^
: string
| create table n as table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create table n table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create view v as table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| create view v as select 1 union table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t (a) (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| insert into t () (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| insert into t ((table u))
< error
| 1064 42000 …near 'table u))' at line 1
> exec ^
: string
| insert into t select 1 union table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t select 1 union table u into @x
< error
| 1064 42000 …near 'table u into @x' at line 1
> exec ^
: string
| replace into t (a) (table u)
< error
| 1064 42000 …near 'table u)' at line 1
> exec ^
: string
| do (table u)
< error
| 1064 42000 …near '(table u)' at line 1
> exec ^
: string
| select * from t join (table u) v
< error
| 1064 42000 …near 'u) v' at line 1
> exec ^
: string
| select * from (select 1 union table u) v
< error
| 1064 42000 …near 'table u) v' at line 1
> exec ^
: string
| values row(1)
< ok 0 0
> exec ^
: string
| values row(1), row(2)
< ok 0 0
> exec ^
: string
| values row()
< error
| 3942 HY000 Each row of a VALUES clause must have at least one column, unless when used as source in an INSERT statement.
> exec ^
: string
| values row(1) order by 1
< ok 0 0
> exec ^
: string
| values row(1) limit 1
< ok 0 0
> exec ^
: string
| values row(1) order by 1 limit 1
< ok 0 0
> exec ^
: string
| values row(1) limit 1, 2
< ok 0 0
> exec ^
: string
| values row(1) limit 1 offset 2
< ok 0 0
> exec ^
: string
| values row(1) for update
< ok 0 0
> exec ^
: string
| values row(1) for share
< error
| 1064 42000 …near 'share' at line 1
> exec ^
: string
| values row(1) lock in share mode
< ok 0 0
> exec ^
: string
| values row(1) for update nowait
< ok 0 0
> exec ^
: string
| values row(1) for update wait 1
< ok 0 0
> exec ^
: string
| values row(1) order by 1 for update
< ok 0 0
> exec ^
: string
| values row(1) limit 1 for update
< ok 0 0
> exec ^
: string
| values row(1) order by 1 limit 1 for update
< ok 0 0
> exec ^
: string
| values row(1) into @x
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f'
< ok 0 0
> exec ^
: string
| values row(1) order by 1 into @x
< ok 0 0
> exec ^
: string
| values row(1) limit 1 into @x
< ok 0 0
> exec ^
: string
| values row(1) for update into @x
< ok 0 0
> exec ^
: string
| values row(1) into @x for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| values row(1) order by 1 order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| values row(1) limit 1 limit 1
< error
| 1064 42000 …near 'limit 1' at line 1
> exec ^
: string
| values row(1) limit 1 order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| explain values row(1)
< ok 0 0
> exec ^
: string
| explain values row(1) for update
< ok 0 0
> exec ^
: string
| explain values row(1) into @x
< ok 0 0
> exec ^
: string
| (values row(1))
< ok 0 0
> exec ^
: string
| ((values row(1)))
< ok 0 0
> exec ^
: string
| (values row(1)) order by 1
< ok 0 0
> exec ^
: string
| (values row(1)) limit 1
< ok 0 0
> exec ^
: string
| (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| (values row(1)) lock in share mode
< error
| 1064 42000 …near 'lock in share mode' at line 1
> exec ^
: string
| (values row(1)) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| (values row(1) for update)
< ok 0 0
> exec ^
: string
| (values row(1) into @x)
< error
| 1064 42000 …near 'into @x)' at line 1
> exec ^
: string
| (values row(1) order by 1)
< ok 0 0
> exec ^
: string
| (values row(1) limit 1)
< ok 0 0
> exec ^
: string
| (values row(1) order by 1) order by 1
< ok 0 0
> exec ^
: string
| (values row(1) limit 1) limit 1
< ok 0 0
> exec ^
: string
| (values row(1) for update) for update
< error
| 1064 42000 …near 'for update' at line 1
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
| (select 1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| (select 1 for update)
< ok 0 0
> exec ^
: string
| (select 1 into @x)
< ok 0 0
> exec ^
: string
| ((select 1)) for update
< error
| 1064 42000 …near 'for update' at line 1
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
| (select 1) order by 1 into @x
< ok 0 0
> exec ^
: string
| (select 1) limit 1 into @x
< ok 0 0
> exec ^
: string
| (select 1) limit 1 for update into @x
< ok 0 0
> exec ^
: string
| ((select 1)) order by 1 for update
< ok 0 0
> exec ^
: string
| (values row(1)) order by 1 for update
< ok 0 0
> exec ^
: string
| (values row(1)) limit 1 for update
< ok 0 0
> exec ^
: string
| (values row(1)) order by 1 into @x
< ok 0 0
> exec ^
: string
| (values row(1)) limit 1 for update into @x
< ok 0 0
> exec ^
: string
| ((values row(1))) order by 1 for update
< ok 0 0
> exec ^
: string
| select * from (select 1 into @x) v
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| (values row(1) for update) order by 1
< ok 0 0
> exec ^
: string
| (values row(1) order by 1) limit 1
< ok 0 0
> exec ^
: string
| (values row(1) order by 1) order by 1 limit 1
< ok 0 0
> exec ^
: string
| with c as (select 1) values row(1)
< ok 0 0
> exec ^
: string
| with c as (select 1) values row(1) order by 1
< ok 0 0
> exec ^
: string
| with c as (select 1) values row(1) for update
< ok 0 0
> exec ^
: string
| with c as (select 1) values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| with c as (select 1) (values row(1))
< ok 0 0
> exec ^
: string
| with c as (select 1) (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| with c as (values row(1) into @x) select * from c
< error
| 1064 42000 …near 'into @x) select * from c' at line 1
> exec ^
: string
| with c as (select 1) select 1 into @x
< ok 0 0
> exec ^
: string
| with c as (select 1) select 1 for update
< ok 0 0
> exec ^
: string
| explain (values row(1))
< ok 0 0
> exec ^
: string
| explain (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| explain values row(1) order by 1
< ok 0 0
> exec ^
: string
| explain values row(1) limit 1 for update
< ok 0 0
> exec ^
: string
| explain (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| do (values row(1) for update)
< ok 0 0
> exec ^
: string
| select (values row(1)) from t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t where exists (values row(1) for update)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t where a = (values row(1) limit 1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t join (values row(1)) v
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from (values row(1)) v join t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from (values row(1)) as v
< ok 0 0
> exec ^
: string
| select * from (values row(1)) as v (a)
< ok 0 0
> exec ^
: string
| select * from (values row(1)), t
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| values row(1) for update for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| values row(1) into @x into @y
< error
| 1064 42000 …near 'into @y' at line 1
> exec ^
: string
| values row(1) lock in share mode for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| values row(1) into @x, @y
< ok 0 0
> exec ^
: string
| values row(1) into dumpfile 'f'
< ok 0 0
> exec ^
: string
| declare c cursor for values row(1)
< error
| 1064 42000 …near 'declare' at line 1
> exec ^
: string
| explain with c as (select 1) values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| explain with c as (select 1) values row(1)
< ok 0 0
> exec ^
: string
| select * from (values row(1)) v
< ok 0 0
> exec ^
: string
| select * from (values row(1) for update) v
< ok 0 0
> exec ^
: string
| select * from (values row(1) order by 1) v
< ok 0 0
> exec ^
: string
| select * from (values row(1) limit 1) v
< ok 0 0
> exec ^
: string
| select * from (values row(1) into @x) v
< error
| 1064 42000 …near 'into @x) v' at line 1
> exec ^
: string
| select * from (values row(1) order by 1 for update) v
< ok 0 0
> exec ^
: string
| select * from ((values row(1))) v
< ok 0 0
> exec ^
: string
| select * from ((values row(1)) for update) v
< error
| 1064 42000 …near 'for update) v' at line 1
> exec ^
: string
| select (values row(1))
< ok 0 0
> exec ^
: string
| select (values row(1) for update)
< ok 0 0
> exec ^
: string
| select (values row(1) into @x)
< error
| 1064 42000 …near 'into @x)' at line 1
> exec ^
: string
| select * from t where a in (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t where a in (values row(1) for update)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select exists (values row(1))
< ok 0 0
> exec ^
: string
| with c as (values row(1)) select * from c
< ok 0 0
> exec ^
: string
| with c as (values row(1) for update) select * from c
< ok 0 0
> exec ^
: string
| with c as (values row(1) order by 1) select * from c
< ok 0 0
> exec ^
: string
| create view v as values row(1)
< ok 0 0
> exec ^
: string
| create view v as values row(1) for update
< ok 0 0
> exec ^
: string
| create view v as values row(1) order by 1
< ok 0 0
> exec ^
: string
| create view v as (values row(1))
< ok 0 0
> exec ^
: string
| create table n as values row(1)
< ok 0 0
> exec ^
: string
| create table n as values row(1) for update
< ok 0 0
> exec ^
: string
| create table n values row(1)
< ok 0 0
> exec ^
: string
| create table n (values row(1))
< ok 0 0
> exec ^
: string
| create table n as (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| do (values row(1))
< ok 0 0
> exec ^
: string
| create table n values row(1) for update
< ok 0 0
> exec ^
: string
| create table n values row(1) order by 1
< ok 0 0
> exec ^
: string
| create table n values row(1) into @x
< ok 0 0
> exec ^
: string
| create table n as values row(1) into @x
< ok 0 0
> exec ^
: string
| create table n as values row(1) order by 1
< ok 0 0
> exec ^
: string
| create table n as values row(1) union values row(2)
< ok 0 0
> exec ^
: string
| create table n as (values row(1))
< ok 0 0
> exec ^
: string
| create table n ((values row(1)))
< ok 0 0
> exec ^
: string
| create table n (values row(1)) union (values row(2))
< ok 0 0
> exec ^
: string
| create table n (select 1)
< ok 0 0
> exec ^
: string
| create table n (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table n as (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table n (a int) values row(1)
< ok 0 0
> exec ^
: string
| create table n (a int) (values row(1))
< ok 0 0
> exec ^
: string
| create table n (a int) as values row(1)
< ok 0 0
> exec ^
: string
| create table n (a int) select 1
< ok 0 0
> exec ^
: string
| create table n (select 1 union values row(1))
< ok 0 0
> exec ^
: string
| create table n select 1 union values row(1)
< ok 0 0
> exec ^
: string
| create table n values row(1) union select 1
< ok 0 0
> exec ^
: string
| create view v as values row(1) into @x
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause.
> exec ^
: string
| create view v as values row(1) union values row(2)
< ok 0 0
> exec ^
: string
| create view v as (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create view v as (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create view v (a) as values row(1)
< ok 0 0
> exec ^
: string
| alter view v as values row(1)
< ok 0 0
> exec ^
: string
| alter view v as values row(1) into @x
< error
| 1350 HY000 View's SELECT contains a 'INTO' clause.
> exec ^
: string
| alter view v as values row(1) for update
< ok 0 0
> exec ^
: string
| alter view v as (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| create table n (a int) (select 1)
< ok 0 0
> exec ^
: string
| create table n (a int) (select 1) union (select 2)
< ok 0 0
> exec ^
: string
| create table n (a int) with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table n (a int) engine = innodb values row(1)
< ok 0 0
> exec ^
: string
| create table n (a int) engine = innodb (select 1)
< ok 0 0
> exec ^
: string
| create table n engine = innodb values row(1)
< ok 0 0
> exec ^
: string
| create table n engine = innodb as values row(1)
< ok 0 0
> exec ^
: string
| create table n (a int) partition by hash(a) values row(1)
< ok 0 0
> exec ^
: string
| create table n with c as (select 1) values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| create table n as with c as (select 1) values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| create view v as with c as (select 1) values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| create table n (select 1) order by 1 into @x
< ok 0 0
> exec ^
: string
| create table n (values row(1)) limit 1 into @x
< ok 0 0
> exec ^
: string
| create table n (values row(1)) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| create table n (a int) as values row(1) into @x
< ok 0 0
> exec ^
: string
| create table n (a int) as (select 1) limit 1 into @x
< ok 0 0
> exec ^
: string
| create table n engine = innodb as values row(1) into @x
< ok 0 0
> exec ^
: string
| create table n engine = innodb as (select 1) limit 1 into @x
< ok 0 0
> exec ^
: string
| create table n (a int) engine = innodb as values row(1) into @x
< ok 0 0
> exec ^
: string
| create table n (a int) engine = innodb (values row(1)) limit 1 into @x
< ok 0 0
> exec ^
: string
| create table n with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| create table n as with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| values row(1) union values row(2)
< ok 0 0
> exec ^
: string
| values row(1) union all values row(2)
< ok 0 0
> exec ^
: string
| values row(1) union select 1
< ok 0 0
> exec ^
: string
| select 1 union values row(1)
< ok 0 0
> exec ^
: string
| select 1 union values row(1) for update
< ok 0 0
> exec ^
: string
| select 1 union values row(1) order by 1
< ok 0 0
> exec ^
: string
| select 1 union values row(1) limit 1
< ok 0 0
> exec ^
: string
| select 1 union values row(1) into @x
< ok 0 0
> exec ^
: string
| select 1 union values row(1) order by 1 for update
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) for update
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) order by 1
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) limit 1
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) into @x
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) order by 1 limit 1 for update
< ok 0 0
> exec ^
: string
| values row(1) for update union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| values row(1) order by 1 union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| values row(1) limit 1 union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| values row(1) into @x union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| (values row(1)) union (values row(2))
< ok 0 0
> exec ^
: string
| (values row(1)) union (values row(2)) for update
< ok 0 0
> exec ^
: string
| (values row(1)) union (values row(2)) into @x
< ok 0 0
> exec ^
: string
| (values row(1) for update) union (values row(2))
< ok 0 0
> exec ^
: string
| select 1 for update union values row(1)
< error
| 1064 42000 …near 'union values row(1)' at line 1
> exec ^
: string
| select 1 union select 2 into @x
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| select 1 union select 2 for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) into @x
< ok 0 0
> exec ^
: string
| select 1 into @x union select 2
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| select 1 for update union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 order by 1 union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 limit 1 union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 union select 2 for update union select 3
< error
| 1064 42000 …near 'union select 3' at line 1
> exec ^
: string
| select 1 union select 2 order by 1 union select 3
< error
| 1064 42000 …near 'union select 3' at line 1
> exec ^
: string
| select 1 union select 2 limit 1 union select 3
< error
| 1064 42000 …near 'union select 3' at line 1
> exec ^
: string
| select 1 union values row(1) for update union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 union values row(1) order by 1 union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 union values row(1) limit 1 union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 union values row(1) into @x union select 2
< error
| 1064 42000 …near 'union select 2' at line 1
> exec ^
: string
| select 1 union (select 2) for update union select 3
< error
| 1064 42000 …near 'union select 3' at line 1
> exec ^
: string
| select * from (values row(1) order by 1 union values row(2)) v
< error
| 1064 42000 …near 'union values row(2)) v' at line 1
> exec ^
: string
| select * from (values row(1) limit 1 union values row(2)) v
< error
| 1064 42000 …near 'union values row(2)) v' at line 1
> exec ^
: string
| select * from (values row(1) for update union values row(2)) v
< error
| 1064 42000 …near 'union values row(2)) v' at line 1
> exec ^
: string
| select * from (values row(1) union values row(2)) v
< ok 0 0
> exec ^
: string
| select * from (values row(1) union values row(2) for update) v
< ok 0 0
> exec ^
: string
| select * from (values row(1) union values row(2) into @x) v
< error
| 1064 42000 …near 'into @x) v' at line 1
> exec ^
: string
| select * from (select 1 union values row(1) into @x) v
< error
| 1064 42000 …near 'into @x) v' at line 1
> exec ^
: string
| (select 1 union values row(1) into @x)
< error
| 1064 42000 …near 'into @x)' at line 1
> exec ^
: string
| (select 1 union select 2 into @x)
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| (values row(1) union values row(2) into @x)
< error
| 1064 42000 …near 'into @x)' at line 1
> exec ^
: string
| (values row(1) union values row(2) for update)
< ok 0 0
> exec ^
: string
| (values row(1) union values row(2)) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| (values row(1) union values row(2)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| with c as (select 1) select 1 union values row(1)
< ok 0 0
> exec ^
: string
| explain values row(1) union values row(2)
< ok 0 0
> exec ^
: string
| select * from t where a in (values row(1) union values row(2))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select 1 union all (values row(1))
< ok 0 0
> exec ^
: string
| select 1 except values row(1)
< ok 0 0
> exec ^
: string
| select 1 intersect values row(1)
< ok 0 0
> exec ^
: string
| values row(1) except values row(2)
< ok 0 0
> exec ^
: string
| values row(1) intersect values row(2)
< ok 0 0
> exec ^
: string
| values row(1) union distinct values row(2)
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) union values row(3)
< ok 0 0
> exec ^
: string
| values row(1) union (values row(2))
< ok 0 0
> exec ^
: string
| (values row(1)) union values row(2)
< ok 0 0
> exec ^
: string
| (values row(1)) union values row(2) for update
< ok 0 0
> exec ^
: string
| (values row(1)) order by 1 union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| (values row(1)) limit 1 union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| (values row(1)) for update union values row(2)
< error
| 1064 42000 …near 'for update union values row(2)' at line 1
> exec ^
: string
| insert into t values row(1) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) order by 1 for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) lock in share mode
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) order by 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) limit 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t (values row(1)) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| insert into t (values row(1) for update)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1) into @x)
< error
| 1064 42000 …near 'into @x)' at line 1
> exec ^
: string
| insert into t (a) (values row(1) for update)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t ((values row(1)))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t ((values row(1))) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t (values row(1)) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) for update on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) into @x on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t select 1 union values row(1) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t select 1 union values row(1) into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union values row(2) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union values row(2) into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t (select 1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| insert into t (select 1 for update)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t ((select 1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t select 1 for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t select 1 into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values row(1) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t (select 1) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t values row(1) union values row(2) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) order by 1 on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) values row(1) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) values row(1) into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values row(1) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) select 1 for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) order by 1 for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) order by 1 for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) union (select 2) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) union (select 2) into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t with c as (select 1) values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t with c as (select 1) values row(1) for update
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values row(1) into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t (values row(1)) for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| insert into t with c as (select 1) values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| insert into t with c as (select 1) select 1 union values row(1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| select 1 for share
< error
| 1064 42000 …near 'share' at line 1
> exec ^
: string
| select * from t for share
< error
| 1064 42000 …near 'share' at line 1
> exec ^
: string
| select * from t for update of t
< error
| 1064 42000 …near 'of t' at line 1
> exec ^
: string
| select * from t for update nowait
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t for update wait 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t for update skip locked
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from t lock in share mode
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| (select 1) limit 1 for share
< error
| 1064 42000 …near 'share' at line 1
> exec ^
: string
| values row(1) for update of t
< error
| 1064 42000 …near 'of t' at line 1
> exec ^
: string
| values row(1) for update skip locked
< ok 0 0
> exec ^
: string
| values row(1) into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| values row(1) into outfile
< error
| 1064 42000 …near 'outfile' at line 1
> exec ^
: string
| values row(1) into dumpfile
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into @x,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| values row(1) into ,@x
< error
| 1064 42000 …near ',@x' at line 1
> exec ^
: string
| values row(1) into @x,,@y
< error
| 1064 42000 …near ',@y' at line 1
> exec ^
: string
| values row(1) into outfile 'f' 'g'
< ok 0 0
> exec ^
: string
| values row(1) into dumpfile 'f' 'g'
< ok 0 0
> exec ^
: string
| values row(1) into outfile ''
< ok 0 0
> exec ^
: string
| values row(1) into dumpfile ''
< ok 0 0
> exec ^
: string
| select 1 into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| select 1 into outfile
< error
| 1064 42000 …near 'outfile' at line 1
> exec ^
: string
| select 1 into dumpfile
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| select 1 into @x,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| select 1 from t into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| select 1 from t into outfile
< error
| 1064 42000 …near 'outfile' at line 1
> exec ^
: string
| select 1 from t into @x,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| (select 1) order by 1 into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| (select 1) order by 1 into outfile
< error
| 1064 42000 …near 'outfile' at line 1
> exec ^
: string
| (select 1) union (select 2) into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| (select 1) union (select 2) into @x,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| insert into t values row(1) into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| insert into t select 1 into
< error
| 1064 42000 …near 'into' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields
< error
| 1064 42000 …near 'fields' at line 1
> exec ^
: string
| values row(1) into outfile 'f' columns
< error
| 1064 42000 …near 'columns' at line 1
> exec ^
: string
| values row(1) into outfile 'f' lines
< error
| 1064 42000 …near 'lines' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields lines
< error
| 1064 42000 …near 'lines' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by
< error
| 1064 42000 …near 'by' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields terminated
< error
| 1064 42000 …near 'terminated' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields terminated ','
< error
| 1064 42000 …near '','' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ','
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' columns terminated by ','
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields enclosed by
< error
| 1064 42000 …near 'by' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields enclosed '"'
< error
| 1064 42000 …near ''"'' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields enclosed by '"'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields optionally
< error
| 1064 42000 …near 'optionally' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields optionally enclosed by '"'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields optionally terminated by ','
< error
| 1064 42000 …near 'terminated by ','' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields escaped by
< error
| 1064 42000 …near 'by' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields escaped '\\'
< error
| 1064 42000 …near ''\\'' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields escaped by '\\'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' lines starting by
< error
| 1064 42000 …near 'by' at line 1
> exec ^
: string
| values row(1) into outfile 'f' lines starting 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| values row(1) into outfile 'f' lines starting by 'a'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' lines terminated by
< error
| 1064 42000 …near 'by' at line 1
> exec ^
: string
| values row(1) into outfile 'f' lines terminated 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| values row(1) into outfile 'f' lines terminated by 'a'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ''
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by '' escaped by '' lines terminated by ''
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields enclosed by ''
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields optionally enclosed by ''
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' lines starting by ''
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields escaped by ',' terminated by ';'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields enclosed by '"' terminated by ';'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' lines terminated by 'a' starting by 'b'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ',' fields enclosed by '"'
< error
| 1064 42000 …near 'fields enclosed by '"'' at line 1
> exec ^
: string
| values row(1) into outfile 'f' lines terminated by 'a' fields terminated by ','
< error
| 1064 42000 …near 'fields terminated by ','' at line 1
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ',' terminated by ';'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' character set
< error
| 1064 42000 …near 'set' at line 1
> exec ^
: string
| values row(1) into outfile 'f' charset
< error
| 1064 42000 …near 'charset' at line 1
> exec ^
: string
| values row(1) into outfile 'f' character set utf8mb4
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' charset utf8mb4
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' character utf8mb4
< error
| 1064 42000 …near 'utf8mb4' at line 1
> exec ^
: string
| values row(1) into outfile 'f' charset utf8mb4 fields terminated by ','
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ',' charset utf8mb4
< error
| 1064 42000 …near 'charset utf8mb4' at line 1
> exec ^
: string
| values row(1) into dumpfile 'f' fields terminated by ','
< error
| 1064 42000 …near 'fields terminated by ','' at line 1
> exec ^
: string
| values row(1) into @x fields terminated by ','
< error
| 1064 42000 …near 'fields terminated by ','' at line 1
> exec ^
: string
| values row(1) into 'f'
< error
| 1064 42000 …near ''f'' at line 1
> exec ^
: string
| values row(1) into 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| values row(1) into x
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into x, y
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into @x, y
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into @@autocommit
< error
| 1064 42000 …near '@@autocommit' at line 1
> exec ^
: string
| values row(1) union values row(2) order by 1 for update
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) order by 1 limit 1 for update into @x
< ok 0 0
> exec ^
: string
| select 1 union select 2 order by 1 for update
< ok 0 0
> exec ^
: string
| select 1 union select 2 order by 1 into @x
< ok 0 0
> exec ^
: string
| select 1 union select 2 into @x union select 3
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| select 1 union select 2 limit 1 for update into @x
< ok 0 0
> exec ^
: string
| select 1 from t union select 2 from t into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select 1 from t union select 2 from t order by 1 into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select 1 from t into @x order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| select 1 into @x order by 1
< ok 0 0
> exec ^
: string
| select 1 union select 2 into @x order by 1
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| (select 1 union select 2 for update)
< ok 0 0
> exec ^
: string
| select * from (select 1 union select 2 for update) v
< ok 0 0
> exec ^
: string
| select * from (select 1 union select 2 order by 1 for update) v
< ok 0 0
> exec ^
: string
| select 1 union select 2 for update nowait
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ',' ';'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' lines terminated by 'a' 'b'
< ok 0 0
> exec ^
: string
| values row(1) into dumpfile, x
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into dumpfile, @x
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into x, dumpfile
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) into outfile, x
< error
| 1064 42000 …near ', x' at line 1
> exec ^
: string
| values row(1) into outfile 'f' charset binary
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' character set binary
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' charset default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| values row(1) into outfile 'f' charset 'utf8mb4'
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields optionally enclosed by '"' enclosed by "'"
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields enclosed by '"' optionally enclosed by "'"
< ok 0 0
> exec ^
: string
| values row(1) into outfile 'f' fields terminated by ',' lines terminated by 'a' fields enclosed by '"'
< error
| 1064 42000 …near 'fields enclosed by '"'' at line 1
> exec ^
: string
| select 1 into @x into @y
< error
| 1064 42000 …near 'into @y' at line 1
> exec ^
: string
| select 1 into x for update
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| select 1 for update into x
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| select 1 into outfile 'f' for update
< ok 0 0
> exec ^
: string
| (select 1 into x)
< error
| 1235 0A000 dynamic sql into variable not a question mark not supported
> exec ^
: string
| values row(1) order by 1 into outfile 'f' fields terminated by ','
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) into outfile 'f'
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) order by 1 into dumpfile 'f'
< ok 0 0
> exec ^
: string
| (select 1) limit 1 into outfile 'f' lines starting by 'a'
< ok 0 0
> exec ^
: string
| create table n values row(1) into outfile 'f' fields terminated by ','
< ok 0 0
> exec ^
: string
| insert into t values row(1) into outfile 'f' on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) for update
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) for update for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) order by 1
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) order by 1 order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) limit 1
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) limit 1 limit 1
< error
| 1064 42000 …near 'limit 1' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) order by 1 limit 1 for update into @x
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) into @x
< ok 0 0
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) into @x into @y
< error
| 1064 42000 …near 'into @y' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) for update union values row(4)
< error
| 1064 42000 …near 'union values row(4)' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) order by 1 union values row(4)
< error
| 1064 42000 …near 'union values row(4)' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) limit 1 union values row(4)
< error
| 1064 42000 …near 'union values row(4)' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) for update intersect values row(4)
< error
| 1064 42000 …near 'intersect values row(4)' at line 1
> exec ^
: string
| values row(1) union values row(2) intersect values row(3) intersect values row(4) for update
< ok 0 0
> exec ^
: string
| values row(1) intersect values row(2) union values row(3) for update
< ok 0 0
> exec ^
: string
| values row(1) intersect values row(2) for update union values row(3)
< error
| 1064 42000 …near 'union values row(3)' at line 1
> exec ^
: string
| select 1 union select 2 intersect select 3 for update
< ok 0 0
> exec ^
: string
| select 1 union select 2 intersect select 3 for update for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| select 1 union select 2 intersect select 3 order by 1
< ok 0 0
> exec ^
: string
| select 1 union select 2 intersect select 3 order by 1 order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| select 1 union select 2 intersect select 3 limit 1 limit 1
< error
| 1064 42000 …near 'limit 1' at line 1
> exec ^
: string
| select 1 union select 2 intersect select 3 for update union select 4
< error
| 1064 42000 …near 'union select 4' at line 1
> exec ^
: string
| select 1 union select 2 intersect select 3 order by 1 union select 4
< error
| 1064 42000 …near 'union select 4' at line 1
> exec ^
: string
| select 1 union select 2 intersect select 3 into @x
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| select 1 union select 2 intersect select 3 into @x union select 4
< error
| 5617 HY000 inappropriate INTO
> exec ^
: string
| select 1 union select 2 intersect select 3 order by 1 limit 1 for update into @x
< ok 0 0
> exec ^
: string
| select 1 except select 2 intersect select 3 for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) intersect (select 3) for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) intersect (select 3) for update for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| (select 1) union (select 2) intersect (select 3) order by 1 for update
< ok 0 0
> exec ^
: string
| (select 1) union (select 2) intersect (select 3) into @x
< ok 0 0
> exec ^
: string
| select 1 union (select 2 intersect select 3 for update)
< ok 0 0
> exec ^
: string
| select 1 union (select 2 intersect select 3) for update
< ok 0 0
> exec ^
: string
| select 1 union (select 2 intersect select 3 order by 1) order by 1
< ok 0 0
> exec ^
: string
| select 1 union (select 2 intersect select 3) for update union select 4
< error
| 1064 42000 …near 'union select 4' at line 1
> exec ^
: string
| select 1 union select 2 union select 3 for update for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| select 1 union select 2 union select 3 order by 1 order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
