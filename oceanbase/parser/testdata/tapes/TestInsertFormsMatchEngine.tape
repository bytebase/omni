-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestInsertFormsMatchEngine default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_sweep
< ok 1 0
> exec
| USE omni_sweep
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM ?
: string
| insert into t values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t value ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t value ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert t values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace t value ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (), ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t value (), ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (), (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1), ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) value ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a, b) values (), ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t partition (p0) values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values () on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values( )
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (/* c */)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| explain insert into t values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () value ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t () value ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t ( ) values ( )
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (/* c */) values (/* c */)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () value (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values (), (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t partition (p0) () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into db.t () values ()
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| insert ignore into t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert low_priority into t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace low_priority into t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values () on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| explain insert into t () values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () select 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () (select 1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () with c as (select 1) select * from c
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| replace into t () table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t () values row()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () set a = 1
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec ^
: string
| insert into t (a) set a = 1
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec ^
: string
| replace into t (a) set a = 1
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec ^
: string
| replace into t () set a = 1
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec ^
: string
| insert into t ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t (a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t () on duplicate key update a = 1
< error
| 1064 42000 …near 'on duplicate key update a = 1' at line 1
> exec ^
: string
| insert into t (()) values ()
< error
| 1064 42000 …near ')) values ()' at line 1
> exec ^
: string
| insert into t (,) values ()
< error
| 1064 42000 …near ',) values ()' at line 1
> exec ^
: string
| insert into t (a,) values (1)
< error
| 1064 42000 …near ') values (1)' at line 1
> exec ^
: string
| insert into t (,a) values (1)
< error
| 1064 42000 …near ',a) values (1)' at line 1
> exec ^
: string
| insert into t values
< error
| 1064 42000 …near 'values' at line 1
> exec ^
: string
| insert into t value
< error
| 1064 42000 …near 'value' at line 1
> exec ^
: string
| insert into t values (1) ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| insert into t values () ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| insert into t values default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| insert into t values (),
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| insert into t values ,()
< error
| 1064 42000 …near ',()' at line 1
> exec ^
: string
| replace into t values () on duplicate key update a = 1
< error
| 1064 42000 …near 'on duplicate key update a = 1' at line 1
> exec ^
: string
| insert into t (@a) values (1)
< error
| 1064 42000 …near '@a) values (1)' at line 1
> exec ^
: string
| insert into t (a, @b) values (1, 2)
< error
| 1064 42000 …near '@b) values (1, 2)' at line 1
> exec ^
: string
| insert into t (@a) select 1
< error
| 1064 42000 …near '@a) select 1' at line 1
> exec ^
: string
| replace into t (@a) values (1)
< error
| 1064 42000 …near '@a) values (1)' at line 1
> exec ^
: string
| insert into t (t.*) values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (`@a`) values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (value) values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values) values (1)
< error
| 1064 42000 …near ') values (1)' at line 1
> exec ^
: string
| insert into t (table) values (1)
< error
| 1064 42000 …near ') values (1)' at line 1
> exec ^
: string
| insert into t (select) values (1)
< error
| 1064 42000 …near ') values (1)' at line 1
> exec ^
: string
| insert into t (with) values (1)
< error
| 1064 42000 …near ') values (1)' at line 1
> exec ^
: string
| insert into t values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1), row(2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(), row()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values row()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t value row(1)
< error
| 1064 42000 …near 'row(1)' at line 1
> exec ^
: string
| insert into t value row()
< error
| 1064 42000 …near 'row()' at line 1
> exec ^
: string
| insert into t value row(1), row(2)
< error
| 1064 42000 …near 'row(1), row(2)' at line 1
> exec ^
: string
| insert into t (a) value row(1)
< error
| 1064 42000 …near 'row(1)' at line 1
> exec ^
: string
| insert into t () value row()
< error
| 1064 42000 …near 'row()' at line 1
> exec ^
: string
| replace into t value row(1)
< error
| 1064 42000 …near 'row(1)' at line 1
> exec ^
: string
| insert into t values row(1), (2)
< error
| 1064 42000 …near '(2)' at line 1
> exec ^
: string
| insert into t values (1), row(2)
< error
| 1064 42000 …near 'row(2)' at line 1
> exec ^
: string
| insert into t value (1), row(2)
< error
| 1064 42000 …near 'row(2)' at line 1
> exec ^
: string
| insert into t value row(1), (2)
< error
| 1064 42000 …near 'row(1), (2)' at line 1
> exec ^
: string
| replace into t values (1), row(2)
< error
| 1064 42000 …near 'row(2)' at line 1
> exec ^
: string
| replace into t values row(1), (2)
< error
| 1064 42000 …near '(2)' at line 1
> exec ^
: string
| insert into t values row 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| insert into t values row
< error
| 1064 42000 …near 'row' at line 1
> exec ^
: string
| insert into t values row(default)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(default, 1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row((1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) as new
< error
| 1064 42000 …near 'as new' at line 1
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
| insert into t values row(1), row(2) order by 1 limit 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union values row(2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union all select 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union values row(2) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t select 1 union values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) union values row(1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t ((values row(1)))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) union (values row(2))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) order by 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) values row(1) order by 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values row(1) order by 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values row(1) order by 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values row(1) union values row(2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t () (values row(1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1) union values row(2)
< error
| 1064 42000 …near 'union values row(2)' at line 1
> exec ^
: string
| insert into t values (1) limit 1
< error
| 1064 42000 …near 'limit 1' at line 1
> exec ^
: string
| insert into t values (1) order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| insert into t values () order by 1
< error
| 1064 42000 …near 'order by 1' at line 1
> exec ^
: string
| insert into t () values () limit 1
< error
| 1064 42000 …near 'limit 1' at line 1
> exec ^
: string
| insert into t () values (1) union values (2)
< error
| 1064 42000 …near 'union values (2)' at line 1
> exec ^
: string
| insert into t () value () union values ()
< error
| 1064 42000 …near 'union values ()' at line 1
> exec ^
: string
| insert into t value () union value ()
< error
| 1064 42000 …near 'union value ()' at line 1
> exec ^
: string
| insert into t (values ())
< error
| 1064 42000 …near '())' at line 1
> exec ^
: string
| insert into t () (values ())
< error
| 1064 42000 …near '())' at line 1
> exec ^
: string
| insert into t (select 1) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) (select 1) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t ((select 1)) on duplicate key update a = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) as new
< error
| 1064 42000 …near 'as new' at line 1
> exec ^
: string
| insert into t (table u)
< error
| 1064 42000 …near 'u)' at line 1
> exec ^
: string
| insert into t (table u) on duplicate key update a = 1
< error
| 1064 42000 …near 'u) on duplicate key update a = 1' at line 1
> exec ^
: string
| insert into t table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t table u order by a
< error
| 1064 42000 …near 'table u order by a' at line 1
> exec ^
: string
| insert into t (a) table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t table u as new
< error
| 1064 42000 …near 'table u as new' at line 1
> exec ^
: string
| replace into t table u
< error
| 1064 42000 …near 'table u' at line 1
> exec ^
: string
| insert into t table u into @x
< error
| 1064 42000 …near 'table u into @x' at line 1
> exec ^
: string
| insert into t table u limit 1 into @x
< error
| 1064 42000 …near 'table u limit 1 into @x' at line 1
> exec ^
: string
| insert into t select 1 into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t select 1 union select 2 into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t table u union table v into @x
< error
| 1064 42000 …near 'table u union table v into @x' at line 1
> exec ^
: string
| insert into t values row(1) union select 1 into @x
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (select 1) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| insert into t (values row(1)) into @x
< error
| 1064 42000 …near 'into @x' at line 1
> exec ^
: string
| insert into t values (default)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (default, 1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values ((1))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (row(1,2))
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1), (2,3)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| insert into t values (t.*)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1, *)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| insert into t values (1, t.*)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (t.*, 1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (db.t.*)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t value (*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| replace into t values (t.*)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| insert into t values row(t.*)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| select * from (values row(*)) v
< error
| 1064 42000 …near '*)) v' at line 1
> exec ^
: string
| select * from (values row(t.*)) v
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| values row(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| values row(t.*)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| insert into t values (,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (,1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t value (,1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values (,1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values (,1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) values (,1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1), (,2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (,1), (,2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (,,1)
< error
| 1064 42000 …near ',1)' at line 1
> exec ^
: string
| insert into t values (,1,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (,1,2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (,default)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (,t.*)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (1,,2)
< error
| 1064 42000 …near ',2)' at line 1
> exec ^
: string
| insert into t values (1),(,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values row(,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values row(,1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(,1), row(,2)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values row(1,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| values row(,1)
< error
| 3942 HY000 Each row of a VALUES clause must have at least one column, unless when used as source in an INSERT statement.
> exec ^
: string
| select * from (values row(,1)) v
< error
| 3942 HY000 Each row of a VALUES clause must have at least one column, unless when used as source in an INSERT statement.
> exec ^
: string
| insert high_priority t values ()
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert delayed into t values ()
< error
| 1064 42000 …near 'delayed into t values ()' at line 1
> exec ^
: string
| insert delayed t set a = 1
< error
| 1064 42000 …near 'delayed t set a = 1' at line 1
> exec ^
: string
| insert delayed into t select 1
< error
| 1064 42000 …near 'delayed into t select 1' at line 1
> exec ^
: string
| replace delayed into t values ()
< error
| 1064 42000 …near 'delayed into t values ()' at line 1
> exec ^
: string
| replace delayed t values (1)
< error
| 1064 42000 …near 'delayed t values (1)' at line 1
> exec ^
: string
| replace high_priority into t values ()
< error
| 1064 42000 …near 'high_priority into t values ()' at line 1
> exec ^
: string
| replace high_priority into t values (1)
< error
| 1064 42000 …near 'high_priority into t values (1)' at line 1
> exec ^
: string
| replace ignore into t values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert low_priority ignore into t values (1)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert ignore low_priority into t values (1)
< error
| 1064 42000 …near 'low_priority into t values (1)' at line 1
> exec ^
: string
| insert high_priority delayed into t values (1)
< error
| 1064 42000 …near 'delayed into t values (1)' at line 1
> exec ^
: string
| insert into t values () as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values () as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t () values () as new (a, b)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values () as new on duplicate key update a = new.a
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (), () as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1) as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1) as new (a)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1) as new ()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (1) as new (a,)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| insert into t values (1) new
< error
| 1064 42000 …near 'new' at line 1
> exec ^
: string
| insert into t set a = 1 as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set a = 1 as new on duplicate key update a = new.a
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set a = 1 as new (a)
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t values (1) as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| replace into t set a = 1 as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t select 1 as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t (a) select 1 as new from u as new
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t set a:=1
< error
| 1064 42000 …near ':=1' at line 1
> exec ^
: string
| insert into t set a := default
< error
| 1064 42000 …near ':= default' at line 1
> exec ^
: string
| insert into t set a := 1, b = 2
< error
| 1064 42000 …near ':= 1, b = 2' at line 1
> exec ^
: string
| insert into t set a = 1, b := 2
< error
| 1064 42000 …near ':= 2' at line 1
> exec ^
: string
| insert into t set t.a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t set `@a` := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| replace into t set a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| replace t set a := default
< error
| 1064 42000 …near ':= default' at line 1
> exec ^
: string
| insert into t set a : = 1
< error
| 1064 42000 …near ': = 1' at line 1
> exec ^
: string
| insert into t set a :=
< error
| 1064 42000 …near ':=' at line 1
> exec ^
: string
| insert into t set a := 1 on duplicate key update a := 2
< error
| 1064 42000 …near ':= 1 on duplicate key update a := 2' at line 1
> exec ^
: string
| insert into t set a = 1 on duplicate key update a := 2
< error
| 1064 42000 …near ':= 2' at line 1
> exec ^
: string
| insert into t set a = @b := 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t values(1) on duplicate key update a:=1
< error
| 1064 42000 …near ':=1' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update a := values(a)
< error
| 1064 42000 …near ':= values(a)' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update a := default
< error
| 1064 42000 …near ':= default' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update a = 1, b := 2
< error
| 1064 42000 …near ':= 2' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update t.a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update `@a` := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update a : = 1
< error
| 1064 42000 …near ': = 1' at line 1
> exec ^
: string
| insert into t (a) select 1 on duplicate key update a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t values row(1) on duplicate key update a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t () values () on duplicate key update a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| insert into t values (1), (2) as new (a) on duplicate key update a := new.a
< error
| 1064 42000 …near ':= new.a' at line 1
> exec ^
: string
| INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
< error
| 1064 42000 …near ':=VALUES(a)+VALUES(b)' at line 1
> exec ^
: string
| INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
< error
| 1064 42000 …near ':=VALUES(a)+VALUES(b)' at line 1
> exec ^
: string
| update t set a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t set a:=1
< error
| 1064 42000 …near ':=1' at line 1
> exec ^
: string
| update t set a := default
< error
| 1064 42000 …near ':= default' at line 1
> exec ^
: string
| update t set a := 1, b = 2
< error
| 1064 42000 …near ':= 1, b = 2' at line 1
> exec ^
: string
| update t set t.a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t, u set t.a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t set `@a` := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t set a := (select 1)
< error
| 1064 42000 …near ':= (select 1)' at line 1
> exec ^
: string
| update t set a := 1 where b = 2
< error
| 1064 42000 …near ':= 1 where b = 2' at line 1
> exec ^
: string
| update t set a := 1 order by b limit 1
< error
| 1064 42000 …near ':= 1 order by b limit 1' at line 1
> exec ^
: string
| update t set a : = 1
< error
| 1064 42000 …near ': = 1' at line 1
> exec ^
: string
| update t set a :=
< error
| 1064 42000 …near ':=' at line 1
> exec ^
: string
| update t set a := b := 1
< error
| 1064 42000 …near ':= b := 1' at line 1
> exec ^
: string
| update t set a = b := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t set a := @b := 1
< error
| 1064 42000 …near ':= @b := 1' at line 1
> exec ^
: string
| update t set a = @b := 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| update t set a = 1, b := @c := 2
< error
| 1064 42000 …near ':= @c := 2' at line 1
> exec ^
: string
| explain update t set a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| load data infile 'f' into table t set a = 1
< error
| 1017 HY000 File not exist
> exec ^
: string
| load data infile 'f' into table t set a = default
< error
| 1017 HY000 File not exist
> exec ^
: string
| load data infile 'f' into table t set a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| load data infile 'f' into table t set a:=1
< error
| 1064 42000 …near ':=1' at line 1
> exec ^
: string
| load data infile 'f' into table t set a := default
< error
| 1064 42000 …near ':= default' at line 1
> exec ^
: string
| load data infile 'f' into table t set a = 1, b := 2
< error
| 1064 42000 …near ':= 2' at line 1
> exec ^
: string
| load data infile 'f' into table t set t.a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| load data local infile 'f' into table t set a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| load xml infile 'f' into table t set a := 1
< error
| 1064 42000 …near 'xml infile 'f' into table t set a := 1' at line 1
> exec ^
: string
| load data infile 'f' into table t set a : = 1
< error
| 1064 42000 …near ': = 1' at line 1
> exec ^
: string
| load data infile 'f' into table t (a, @b) set a := @b
< error
| 1064 42000 …near ':= @b' at line 1
> exec ^
: string
| update t set @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| update t set @a := 1
< error
| 1064 42000 …near '@a := 1' at line 1
> exec ^
: string
| update t set @@autocommit = 1
< error
| 1064 42000 …near '@@autocommit = 1' at line 1
> exec ^
: string
| update t set @@autocommit := 1
< error
| 1064 42000 …near '@@autocommit := 1' at line 1
> exec ^
: string
| update t set t.* = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| update t set t.* := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| update t set `@a` = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| insert into t set @a := 1
< error
| 1064 42000 …near '@a := 1' at line 1
> exec ^
: string
| insert into t set @@autocommit = 1
< error
| 1064 42000 …near '@@autocommit = 1' at line 1
> exec ^
: string
| insert into t set @@autocommit := 1
< error
| 1064 42000 …near '@@autocommit := 1' at line 1
> exec ^
: string
| insert into t set t.* = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t set `@a` = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update @a := 1
< error
| 1064 42000 …near '@a := 1' at line 1
> exec ^
: string
| insert into t values (1) on duplicate key update t.* = 1
< error
| 1146 42S02 Table 'omni_sweep.t' doesn't exist
> exec ^
: string
| load data infile 'f' into table t set @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| load data infile 'f' into table t set @a := 1
< error
| 1064 42000 …near '@a := 1' at line 1
> exec ^
: string
| load data infile 'f' into table t set t.* = 1
< error
| 1064 42000 …near '* = 1' at line 1
> exec ^
: string
| load data infile 'f' into table t set `@a` = 1
< error
| 1017 HY000 File not exist
> exec ^
: string
| select a := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| select @a := 1
< ok 0 0
> exec ^
: string
| set @a := 1
< ok 0 0
> exec ^
: string
| set autocommit := 1
< ok 0 0
