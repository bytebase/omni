-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestWildcardFieldFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select *, a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select *, a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a, * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a, * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.* from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.* from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select s.t.* from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select s.t.* from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a, t.* from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a, t.* from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select *, t.*, a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select *, t.*, a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select *from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select *from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t . * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t . * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t. * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t. * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t .* from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t .* from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select distinct * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select distinct * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t union all select t.* from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t union all select t.* from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (select * from t limit 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select (select * from t limit 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select exists (select * from t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select exists (select * from t)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t select * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t select * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| with c as (select * from t) select c.* from c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| with c as (select * from t) select c.* from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `*` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `*` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a * b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a * b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.a * 2 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.a * 2 from t
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
| select count(*) = 1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(*) = 1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(*) over () from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(*) over () from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * as x from t
< error
| 1064 42000 …line 1 column 11 near "as x from t" 
> exec ^
: string
| select * x from t
< error
| 1064 42000 …line 1 column 10 near "x from t" 
> exec ^
: string
| select t.* x from t
< error
| 1064 42000 …line 1 column 12 near "x from t" 
> exec ^
: string
| select * = t
< error
| 1064 42000 …line 1 column 10 near "= t" 
> exec ^
: string
| select * = 1 from t
< error
| 1064 42000 …line 1 column 10 near "= 1 from t" 
> exec ^
: string
| select *=1 from t
< error
| 1064 42000 …line 1 column 9 near "=1 from t" 
> exec ^
: string
| select t.* = 1 from t
< error
| 1064 42000 …line 1 column 12 near "= 1 from t" 
> exec ^
: string
| select s.t.* = 1 from t
< error
| 1064 42000 …line 1 column 14 near "= 1 from t" 
> exec ^
: string
| select 1 = * from t
< error
| 1064 42000 …line 1 column 12 near "* from t" 
> exec ^
: string
| select 1 = t.* from t
< error
| 1064 42000 …line 1 column 14 near "* from t" 
> exec ^
: string
| select * <=> 1 from t
< error
| 1064 42000 …line 1 column 12 near "<=> 1 from t" 
> exec ^
: string
| select * + 1 from t
< error
| 1064 42000 …line 1 column 10 near "+ 1 from t" 
> exec ^
: string
| select 1 + * from t
< error
| 1064 42000 …line 1 column 12 near "* from t" 
> exec ^
: string
| select 1 + t.* from t
< error
| 1064 42000 …line 1 column 14 near "* from t" 
> exec ^
: string
| select * * 2 from t
< error
| 1064 42000 …line 1 column 10 near "* 2 from t" 
> exec ^
: string
| select 2 * * from t
< error
| 1064 42000 …line 1 column 12 near "* from t" 
> exec ^
: string
| select * div 2 from t
< error
| 1064 42000 …line 1 column 12 near "div 2 from t" 
> exec ^
: string
| select - * from t
< error
| 1064 42000 …line 1 column 10 near "* from t" 
> exec ^
: string
| select -t.* from t
< error
| 1064 42000 …line 1 column 11 near "* from t" 
> exec ^
: string
| select not * from t
< error
| 1064 42000 …line 1 column 12 near "* from t" 
> exec ^
: string
| select * and 1 from t
< error
| 1064 42000 …line 1 column 12 near "and 1 from t" 
> exec ^
: string
| select 1 and t.* from t
< error
| 1064 42000 …line 1 column 16 near "* from t" 
> exec ^
: string
| select * or 1 from t
< error
| 1064 42000 …line 1 column 11 near "or 1 from t" 
> exec ^
: string
| select * is null from t
< error
| 1064 42000 …line 1 column 11 near "is null from t" 
> exec ^
: string
| select t.* is null from t
< error
| 1064 42000 …line 1 column 13 near "is null from t" 
> exec ^
: string
| select * in (1) from t
< error
| 1064 42000 …line 1 column 11 near "in (1) from t" 
> exec ^
: string
| select * between 1 and 2 from t
< error
| 1064 42000 …line 1 column 16 near "between 1 and 2 from t" 
> exec ^
: string
| select * like 'a' from t
< error
| 1064 42000 …line 1 column 13 near "like 'a' from t" 
> exec ^
: string
| select * regexp 'a' from t
< error
| 1064 42000 …line 1 column 15 near "regexp 'a' from t" 
> exec ^
: string
| select * collate utf8mb4_bin from t
< error
| 1064 42000 …line 1 column 16 near "collate utf8mb4_bin from t" 
> exec ^
: string
| select @a := * from t
< error
| 1064 42000 …line 1 column 14 near "* from t" 
> exec ^
: string
| select (*) from t
< error
| 1064 42000 …line 1 column 9 near "*) from t" 
> exec ^
: string
| select (t.*) from t
< error
| 1064 42000 …line 1 column 11 near "*) from t" 
> exec ^
: string
| select case * when 1 then 2 end from t
< error
| 1064 42000 …line 1 column 13 near "* when 1 then 2 end from t" 
> exec ^
: string
| select case when * then 2 end from t
< error
| 1064 42000 …line 1 column 18 near "* then 2 end from t" 
> exec ^
: string
| select case when 1 then * end from t
< error
| 1064 42000 …line 1 column 25 near "* end from t" 
> exec ^
: string
| select case when 1 then 2 else t.* end from t
< error
| 1064 42000 …line 1 column 34 near "* end from t" 
> exec ^
: string
| select interval * day + d from t
< error
| 1064 42000 …line 1 column 17 near "* day + d from t" 
> exec ^
: string
| select d + interval * day from t
< error
| 1064 42000 …line 1 column 21 near "* day from t" 
> exec ^
: string
| select cast(* as char) from t
< error
| 1064 42000 …line 1 column 13 near "* as char) from t" 
> exec ^
: string
| select cast(t.* as char) from t
< error
| 1064 42000 …line 1 column 15 near "* as char) from t" 
> exec ^
: string
| select convert(*, char) from t
< error
| 1064 42000 …line 1 column 16 near "*, char) from t" 
> exec ^
: string
| select convert(* using utf8mb4) from t
< error
| 1064 42000 …line 1 column 16 near "* using utf8mb4) from t" 
> exec ^
: string
| select count(t.*) from t
< error
| 1064 42000 …line 1 column 16 near "*) from t" 
> exec ^
: string
| select count(distinct *) from t
< error
| 1064 42000 …line 1 column 23 near "*) from t" 
> exec ^
: string
| select count(all *) from t
< error
| 1064 42000 …line 1 column 18 near "*) from t" 
> exec ^
: string
| select 1 from t where *
< error
| 1064 42000 …line 1 column 23 near "*" 
> exec ^
: string
| select 1 from t where t.*
< error
| 1064 42000 …line 1 column 25 near "*" 
> exec ^
: string
| select 1 from t where a = *
< error
| 1064 42000 …line 1 column 27 near "*" 
> exec ^
: string
| select 1 from t where a = t.*
< error
| 1064 42000 …line 1 column 29 near "*" 
> exec ^
: string
| select 1 from t group by *
< error
| 1064 42000 …line 1 column 26 near "*" 
> exec ^
: string
| select 1 from t group by t.*
< error
| 1064 42000 …line 1 column 28 near "*" 
> exec ^
: string
| select 1 from t having *
< error
| 1064 42000 …line 1 column 24 near "*" 
> exec ^
: string
| select 1 from t order by *
< error
| 1064 42000 …line 1 column 26 near "*" 
> exec ^
: string
| select 1 from t order by t.*
< error
| 1064 42000 …line 1 column 28 near "*" 
> exec ^
: string
| do *
< error
| 1064 42000 …line 1 column 4 near "*" 
> exec ^
: string
| do t.*
< error
| 1064 42000 …line 1 column 6 near "*" 
> exec ^
: string
| do * = 1
< error
| 1064 42000 …line 1 column 4 near "* = 1" 
> exec ^
: string
| do 1 = *
< error
| 1064 42000 …line 1 column 8 near "*" 
> exec ^
: string
| do 1, *
< error
| 1064 42000 …line 1 column 7 near "*" 
> exec ^
: string
| set @a = *
< error
| 1064 42000 …line 1 column 10 near "*" 
> exec ^
: string
| set @a = t.*
< error
| 1064 42000 …line 1 column 12 near "*" 
> exec ^
: string
| update t set a = *
< error
| 1064 42000 …line 1 column 18 near "*" 
> exec ^
: string
| update t set a = t.*
< error
| 1064 42000 …line 1 column 20 near "*" 
> exec ^
: string
| insert into t (a) values (*)
< error
| 1064 42000 …line 1 column 27 near "*)" 
> exec ^
: string
| insert into t (a) values (t.*)
< error
| 1064 42000 …line 1 column 29 near "*)" 
> exec ^
: string
| insert into t (a) values (1) on duplicate key update a = *
< error
| 1064 42000 …line 1 column 58 near "*" 
