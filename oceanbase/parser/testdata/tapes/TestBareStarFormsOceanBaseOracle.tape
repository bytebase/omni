-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestBareStarFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| select * from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select *, a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a, * from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select t.*, * from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select *from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select distinct * from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t union all select * from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select (select * from t limit 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| select exists (select * from t)
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t select * from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select `*` from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a * b from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count(*) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count(all *) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count(*) = 1 from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select t.* from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a, t.* from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select t.* = 1 from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 + t.* from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select -t.* from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select (t.*) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count(t.*) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 from t where t.*
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 from t order by t.*
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t set a = t.*
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t (a) values (t.*)
< error
| 1046 3D000 No database selected
> exec ^
: string
| do *
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec ^
: string
| do 1, *
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
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
| select 1 = * from t
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
| select 1 from t where a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| select 1 from t group by *
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
| do * = 1
< error
| 1064 42000 …near '* = 1' at line 1
> exec ^
: string
| do 1 = *
< error
| 1064 42000 …near '1 = *' at line 1
> exec ^
: string
| set @a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| update t set a = *
< error
| 1064 42000 …near '*' at line 1
> exec ^
: string
| insert into t (a) values (*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| insert into t (a) values (1) on duplicate key update a = *
< error
| 1064 42000 …near '*' at line 1
