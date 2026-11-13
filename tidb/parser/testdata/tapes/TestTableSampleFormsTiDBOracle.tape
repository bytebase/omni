-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTableSampleFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select * from tbl tablesample system (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample system (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample bernoulli (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample bernoulli (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample regions ()
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec ^
: string
| select * from tbl tablesample regions ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample regions (50)
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec ^
: string
| select * from tbl tablesample regions (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample ()
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample system ()
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample system ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample bernoulli ()
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample bernoulli ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50 percent)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50 percent)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50 rows)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50 rows)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50 percent rows)
< error
| 1064 42000 …line 1 column 46 near "rows)" 
> exec ^
: string
| select * from tbl tablesample (50 rows percent)
< error
| 1064 42000 …line 1 column 46 near "percent)" 
> exec ^
: string
| select * from tbl tablesample (percent)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (percent)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (percent percent)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (percent percent)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (rows)
< error
| 1064 42000 …line 1 column 35 near "rows)" 
> exec ^
: string
| select * from tbl tablesample (1+1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (1+1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (a)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (a + b rows)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (a + b rows)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample ((select 1))
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample ((select 1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample ('34')
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample ('34')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample ('34' rows)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample ('34' rows)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50, 50)
< error
| 1064 42000 …line 1 column 34 near ", 50)" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1+1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1+1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable ('a')
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) repeatable ('a')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable ()
< error
| 1064 42000 …line 1 column 48 near ")" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| select * from tbl tablesample () repeatable (1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample () repeatable (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample regions () repeatable (1)
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec ^
: string
| select * from tbl tablesample regions () repeatable (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1) repeatable (2)
< error
| 1064 42000 …line 1 column 60 near "repeatable (2)" 
> exec ^
: string
| select * from tbl tablesample (50) tablesample (50)
< error
| 1064 42000 …line 1 column 46 near "tablesample (50)" 
> exec ^
: string
| select * from tbl tablesample (50) a
< error
| 1064 42000 …line 1 column 36 near "a" 
> exec ^
: string
| select * from tbl tablesample (50) as a
< error
| 1064 42000 …line 1 column 37 near "as a" 
> exec ^
: string
| select * from tbl tablesample
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| select * from tbl tablesample system
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| select * from tbl tablesample (
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| select * from tbl tablesample (50
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| select * from tbl tablesample foo (50)
< error
| 1064 42000 …line 1 column 33 near "foo (50)" 
> exec ^
: string
| select * from tbl tablesample 50
< error
| 1064 42000 …line 1 column 32 near "50" 
> exec ^
: string
| select * from tbl tablesample system 50
< error
| 1064 42000 …line 1 column 39 near "50" 
> exec ^
: string
| select * from db.tbl tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from db.tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl a tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl a tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl as a tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl as a tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl `tablesample` tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl `tablesample` tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl partition (p0) tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl partition (p0) tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl partition (p0) a tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl partition (p0) a tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) partition (p0)
< error
| 1064 42000 …line 1 column 44 near "partition (p0)" 
> exec ^
: string
| select * from tbl use index (i) tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl use index (i) tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl a use index (i) tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl a use index (i) tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) use index (i)
< error
| 1064 42000 …line 1 column 38 near "use index (i)" 
> exec ^
: string
| select * from tbl partition (p0) a use index (i) tablesample (50) repeatable (1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl partition (p0) a use index (i) tablesample (50) repeatable (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (tbl tablesample (50))
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from (tbl tablesample (50))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (tbl tablesample (50), t2)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from (tbl tablesample (50), t2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (select 1) a tablesample (50)
< error
| 1064 42000 …line 1 column 38 near "tablesample (50)" 
> exec ^
: string
| select * from (t1 join t2) tablesample (50)
< error
| 1064 42000 …line 1 column 38 near "tablesample (50)" 
> exec ^
: string
| select * from dual tablesample (50)
< error
| 1064 42000 …line 1 column 30 near "tablesample (50)" 
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$')) a tablesample (50)
< error
| 1064 42000 …line 1 column 25 near "('[]', '$' columns (a int path '$')) a tablesample (50)" 
> exec ^
: string
| select * from t1 tablesample (10), t2 tablesample (20)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 tablesample (10), t2 tablesample (20)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 a tablesample (10) join t2 b tablesample (20) on a.id = b.id
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 a tablesample (10) join t2 b tablesample (20) on a.id = b.id
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 tablesample (10) left join t2 tablesample (20) on t1.id = t2.id
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 tablesample (10) left join t2 tablesample (20) on t1.id = t2.id
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 tablesample (10) straight_join t2 tablesample (20)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 tablesample (10) straight_join t2 tablesample (20)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 tablesample (10) natural join t2 tablesample (20)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 tablesample (10) natural join t2 tablesample (20)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 tablesample (10) cross join t2 tablesample (20)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 tablesample (10) cross join t2 tablesample (20)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 join t2 tablesample (20) using (id)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from t1 join t2 tablesample (20) using (id)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl where id > 20 tablesample (50)
< error
| 1064 42000 …line 1 column 43 near "tablesample (50)" 
> exec ^
: string
| select * from tbl tablesample (50) where id > 20
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) where id > 20
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) group by a
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) group by a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) having a > 1
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) having a > 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) order by a
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) order by a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) limit 1
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) for update
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) union select * from t2 tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) union select * from t2 tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample bernoulli(50) limit 1 into outfile '/tmp/sample.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * from tbl tablesample bernoulli(50) limit 1 into outfile '/tmp/sample.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample bernoulli(50) order by a, b into outfile '/tmp/sample.csv'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * from tbl tablesample bernoulli(50) order by a, b into outfile '/tmp/sample.csv'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (select 1 from tbl tablesample (50))
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select (select 1 from tbl tablesample (50))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| with c as (select * from tbl tablesample (50)) select * from c
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| with c as (select * from tbl tablesample (50)) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| with c as (select 1) select * from c tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| with c as (select 1) select * from c tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t select * from tbl tablesample (50)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t select * from tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t2 as select * from tbl tablesample (50)
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t2 as select * from tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create view v as select * from tbl tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| create view v as select * from tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t select * from tbl tablesample (50)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t select * from tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain select * from tbl tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| explain select * from tbl tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update tbl tablesample (50) set a = 1
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| update tbl tablesample (50) set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update tbl a tablesample (50) set a = 1
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| update tbl a tablesample (50) set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update tbl tablesample (50) join t2 on 1 = 1 set a = 1
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| update tbl tablesample (50) join t2 on 1 = 1 set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from tbl tablesample (50)
< error
| 1064 42000 …line 1 column 27 near "tablesample (50)" 
> exec ^
: string
| delete from tbl a tablesample (50) where 1 = 1
< error
| 1064 42000 …line 1 column 29 near "tablesample (50) where 1 = 1" 
> exec ^
: string
| delete t from tbl t tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| delete t from tbl t tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t using tbl t tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| delete from t using tbl t tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete t from t1 t join t2 tablesample (50)
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| delete t from t1 t join t2 tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| lock tables tbl tablesample (50) read
< error
| 1064 42000 …line 1 column 27 near "tablesample (50) read" 
> exec ^
: string
| lock tables tbl read tablesample (50)
< error
| 1064 42000 …line 1 column 32 near "tablesample (50)" 
> exec ^
: string
| table tbl tablesample (50)
< error
| 1064 42000 …line 1 column 21 near "tablesample (50)" 
> exec ^
: string
| insert into tbl tablesample (50) values (1)
< error
| 1064 42000 …line 1 column 27 near "tablesample (50) values (1)" 
> exec ^
: string
| analyze table tbl tablesample (50)
< error
| 1064 42000 …line 1 column 29 near "tablesample (50)" 
> exec ^
: string
| admin check table tbl tablesample (50)
< error
| 1064 42000 …line 1 column 33 near "tablesample (50)" 
> exec ^
: string
| select * from tablesample
< error
| 1064 42000 …line 1 column 25 near "tablesample" 
> exec ^
: string
| select * from tablesample a
< error
| 1064 42000 …line 1 column 25 near "tablesample a" 
> exec ^
: string
| select * from db.tablesample
< error
| 1146 42S02 Table 'db.tablesample' doesn't exist
> exec ^
: string
| select * from db.tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select tablesample.a from tablesample
< error
| 1064 42000 …line 1 column 37 near "tablesample" 
> exec ^
: string
| select tablesample from t
< error
| 1064 42000 …line 1 column 18 near "tablesample from t" 
> exec ^
: string
| select tablesample.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select db.t.tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select db.t.tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select tablesample, a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample, a from t" 
> exec ^
: string
| select a from t where tablesample = 1
< error
| 1064 42000 …line 1 column 33 near "tablesample = 1" 
> exec ^
: string
| select a from t order by tablesample
< error
| 1064 42000 …line 1 column 36 near "tablesample" 
> exec ^
: string
| select a from t group by tablesample
< error
| 1064 42000 …line 1 column 36 near "tablesample" 
> exec ^
: string
| select tablesample()
< error
| 1064 42000 …line 1 column 18 near "tablesample()" 
> exec ^
: string
| select tablesample(1)
< error
| 1064 42000 …line 1 column 18 near "tablesample(1)" 
> exec ^
: string
| select 1 tablesample
< error
| 1064 42000 …line 1 column 20 near "tablesample" 
> exec ^
: string
| select 1 as tablesample
< error
| 1064 42000 …line 1 column 23 near "tablesample" 
> exec ^
: string
| select * from t tablesample
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| select * from t as tablesample
< error
| 1064 42000 …line 1 column 30 near "tablesample" 
> exec ^
: string
| select * from t tablesample where a = 1
< error
| 1064 42000 …line 1 column 33 near "where a = 1" 
> exec ^
: string
| select * from t tablesample, t2
< error
| 1064 42000 …line 1 column 28 near ", t2" 
> exec ^
: string
| select * from t tablesample join t2
< error
| 1064 42000 …line 1 column 32 near "join t2" 
> exec ^
: string
| select * from t tablesample order by 1
< error
| 1064 42000 …line 1 column 33 near "order by 1" 
> exec ^
: string
| select * from t tablesample limit 1
< error
| 1064 42000 …line 1 column 33 near "limit 1" 
> exec ^
: string
| select * from t tablesample group by a
< error
| 1064 42000 …line 1 column 33 near "group by a" 
> exec ^
: string
| select * from t tablesample use index (i)
< error
| 1064 42000 …line 1 column 31 near "use index (i)" 
> exec ^
: string
| select * from t partition (p0) tablesample
< error
| 1064 42000 …line 1 column 42 near "" 
> exec ^
: string
| select * from (select 1) tablesample
< error
| 1064 42000 …line 1 column 36 near "tablesample" 
> exec ^
: string
| select * from (select 1) as tablesample
< error
| 1064 42000 …line 1 column 39 near "tablesample" 
> exec ^
: string
| select * from t1 join t2 tablesample on 1 = 1
< error
| 1064 42000 …line 1 column 39 near "on 1 = 1" 
> exec ^
: string
| with tablesample as (select 1) select * from tablesample
< error
| 1064 42000 …line 1 column 16 near "tablesample as (select 1) select * from tablesample" 
> exec ^
: string
| update tablesample set a = 1
< error
| 1064 42000 …line 1 column 18 near "tablesample set a = 1" 
> exec ^
: string
| update t tablesample set a = 1
< error
| 1064 42000 …line 1 column 24 near "set a = 1" 
> exec ^
: string
| update t as tablesample set a = 1
< error
| 1064 42000 …line 1 column 23 near "tablesample set a = 1" 
> exec ^
: string
| update t set tablesample = 1
< error
| 1064 42000 …line 1 column 24 near "tablesample = 1" 
> exec ^
: string
| delete from tablesample
< error
| 1064 42000 …line 1 column 23 near "tablesample" 
> exec ^
: string
| delete from db.tablesample
< error
| 1146 42S02 Table 'db.tablesample' doesn't exist
> exec ^
: string
| delete from db.tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete tablesample from t
< error
| 1064 42000 …line 1 column 18 near "tablesample from t" 
> exec ^
: string
| delete tablesample from tablesample
< error
| 1064 42000 …line 1 column 18 near "tablesample from tablesample" 
> exec ^
: string
| delete from t tablesample
< error
| 1064 42000 …line 1 column 25 near "tablesample" 
> exec ^
: string
| delete from t tablesample where a = 1
< error
| 1064 42000 …line 1 column 25 near "tablesample where a = 1" 
> exec ^
: string
| delete from t as tablesample where a = 1
< error
| 1064 42000 …line 1 column 28 near "tablesample where a = 1" 
> exec ^
: string
| delete from t tablesample order by a
< error
| 1064 42000 …line 1 column 25 near "tablesample order by a" 
> exec ^
: string
| delete from t tablesample limit 1
< error
| 1064 42000 …line 1 column 25 near "tablesample limit 1" 
> exec ^
: string
| delete from t tablesample using t
< error
| 1064 42000 …line 1 column 25 near "tablesample using t" 
> exec ^
: string
| delete tablesample.* from t tablesample
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| delete t.* from t tablesample
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| insert into tablesample values (1)
< error
| 1064 42000 …line 1 column 23 near "tablesample values (1)" 
> exec ^
: string
| insert into t (tablesample) values (1)
< error
| 1064 42000 …line 1 column 26 near "tablesample) values (1)" 
> exec ^
: string
| create table tablesample (a int)
< error
| 1064 42000 …line 1 column 24 near "tablesample (a int)" 
> exec ^
: string
| create table t (tablesample int)
< error
| 1064 42000 …line 1 column 27 near "tablesample int)" 
> exec ^
: string
| create table db.tablesample (a int)
< ok 0 0
> exec ^
: string
| create index tablesample on t (a)
< error
| 1064 42000 …line 1 column 24 near "tablesample on t (a)" 
> exec ^
: string
| create index i on t (tablesample)
< error
| 1064 42000 …line 1 column 32 near "tablesample)" 
> exec ^
: string
| create database tablesample
< error
| 1064 42000 …line 1 column 27 near "tablesample" 
> exec ^
: string
| create view tablesample as select 1
< error
| 1064 42000 …line 1 column 23 near "tablesample as select 1" 
> exec ^
: string
| drop table tablesample
< error
| 1064 42000 …line 1 column 22 near "tablesample" 
> exec ^
: string
| alter table t add column tablesample int
< error
| 1064 42000 …line 1 column 36 near "tablesample int" 
> exec ^
: string
| alter table t rename to tablesample
< error
| 1064 42000 …line 1 column 35 near "tablesample" 
> exec ^
: string
| alter table tablesample add column a int
< error
| 1064 42000 …line 1 column 23 near "tablesample add column a int" 
> exec ^
: string
| lock tables tablesample read
< error
| 1064 42000 …line 1 column 23 near "tablesample read" 
> exec ^
: string
| lock tables t tablesample read
< error
| 1064 42000 …line 1 column 25 near "tablesample read" 
> exec ^
: string
| lock tables t as tablesample read
< error
| 1064 42000 …line 1 column 16 near "as tablesample read" 
> exec ^
: string
| truncate table tablesample
< error
| 1064 42000 …line 1 column 26 near "tablesample" 
> exec ^
: string
| select @tablesample
< ok 0 0
> exec ^
: string
| select @@tablesample
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| select 1 from t where a in (select tablesample from t2)
< error
| 1064 42000 …line 1 column 46 near "tablesample from t2)" 
> exec ^
: string
| select * from t tablesample force index (i)
< error
| 1064 42000 …line 1 column 33 near "force index (i)" 
> exec ^
: string
| select * from t as tablesample tablesample (50)
< error
| 1064 42000 …line 1 column 30 near "tablesample tablesample (50)" 
> exec ^
: string
| select * from tablesample tablesample (50)
< error
| 1064 42000 …line 1 column 25 near "tablesample tablesample (50)" 
> exec ^
: string
| select tablesample from tablesample tablesample tablesample (50)
< error
| 1064 42000 …line 1 column 18 near "tablesample from tablesample tablesample tablesample (50)" 
> exec ^
: string
| select * from tablesample.t
< error
| 1146 42S02 Table 'tablesample.t' doesn't exist
> exec ^
: string
| select * from tablesample.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tablesample.tablesample
< error
| 1146 42S02 Table 'tablesample.tablesample' doesn't exist
> exec ^
: string
| select * from tablesample.tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select tablesample.t.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample.t.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (select 1)
< error
| 1064 42000 …line 1 column 37 near "select 1)" 
> exec ^
: string
| select * from tbl tablesample `system` (50)
< error
| 1064 42000 …line 1 column 38 near "`system` (50)" 
> exec ^
: string
| select * from tbl tablesample `bernoulli` (50)
< error
| 1064 42000 …line 1 column 41 near "`bernoulli` (50)" 
> exec ^
: string
| select * from tbl tablesample `regions` ()
< error
| 1064 42000 …line 1 column 39 near "`regions` ()" 
> exec ^
: string
| select * from tbl tablesample (50 `percent`)
< error
| 1064 42000 …line 1 column 43 near "`percent`)" 
> exec ^
: string
| select * from tbl tablesample (50 `rows`)
< error
| 1064 42000 …line 1 column 40 near "`rows`)" 
> exec ^
: string
| select * from tbl tablesample (50) `repeatable` (1)
< error
| 1064 42000 …line 1 column 47 near "`repeatable` (1)" 
> exec ^
: string
| select * from tbl tablesample (@a := 1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (@a := 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample (50) repeatable (@a := 1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl tablesample (50) repeatable (@a := 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl TABLESAMPLE SYSTEM (50 PERCENT) REPEATABLE (1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from tbl TABLESAMPLE SYSTEM (50 PERCENT) REPEATABLE (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from tbl tablesample region ()
< error
| 1064 42000 …line 1 column 36 near "region ()" 
> exec ^
: string
| select * from tbl tablesample (50) (select 1)
< error
| 1064 42000 …line 1 column 36 near "(select 1)" 
> exec ^
: string
| select * from t1 tablesample, t2 tablesample
< error
| 1064 42000 …line 1 column 29 near ", t2 tablesample" 
> exec ^
: string
| select * from tbl a tablesample
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| select 1 from t where tablesample
< error
| 1064 42000 …line 1 column 33 near "tablesample" 
> exec ^
: string
| select count(tablesample) from t
< error
| 1064 42000 …line 1 column 24 near "tablesample) from t" 
> exec ^
: string
| select `tablesample` from `tablesample` `tablesample`
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| select `tablesample` from `tablesample` `tablesample`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t `tablesample` where a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t `tablesample` where a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete `tablesample` from t `tablesample`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete `tablesample` from t `tablesample`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t sample (50)
< error
| 1064 42000 …line 1 column 24 near "(50)" 
> exec ^
: string
| select * from t1 tablesample (10) join t2 tablesample (20) join t3 tablesample (30)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t1 tablesample (10) join t2 tablesample (20) join t3 tablesample (30)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (t1 tablesample (10) join t2) join t3 tablesample (30) on 1 = 1
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from (t1 tablesample (10) join t2) join t3 tablesample (30) on 1 = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t1 join (t2 tablesample (20)) on 1 = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from t1 join (t2 tablesample (20)) on 1 = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select table.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select table.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from table.t
< error
| 1146 42S02 Table 'table.t' doesn't exist
> exec ^
: string
| select * from table.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select table.t.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select table.t.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select select.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select select.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from select.t
< error
| 1146 42S02 Table 'select.t' doesn't exist
> exec ^
: string
| select * from select.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.table from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.table from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db.table
< error
| 1146 42S02 Table 'db.table' doesn't exist
> exec ^
: string
| select * from db.table
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select @@table
< error
| 1193 HY000 Unknown system variable 'table'
> exec ^
: string
| select @@select
< error
| 1193 HY000 Unknown system variable 'select'
> exec ^
: string
| select @@session.table
< error
| 1193 HY000 Unknown system variable 'table'
> exec ^
: string
| select @@global.tablesample
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| select @@session.tablesample
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| select @@tablesample.a
< error
| 1193 HY000 Unknown system variable 'tablesample.a'
> exec ^
: string
| select tablesample .a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample .a from t" 
> exec ^
: string
| select * from tablesample .t
< error
| 1064 42000 …line 1 column 25 near "tablesample .t" 
> exec ^
: string
| select tablesample. a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample. a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t . tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t . tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t. tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t. tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t .tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t .tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db . tablesample
< error
| 1146 42S02 Table 'db.tablesample' doesn't exist
> exec ^
: string
| select * from db . tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db .tablesample
< error
| 1146 42S02 Table 'db.tablesample' doesn't exist
> exec ^
: string
| select * from db .tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table tablesample.t (a int)
< ok 0 0
> exec ^
: string
| drop table tablesample.t
< ok 0 0
> exec ^
: string
| delete from tablesample.t
< error
| 1146 42S02 Table 'tablesample.t' doesn't exist
> exec ^
: string
| delete from tablesample.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update tablesample.t set a = 1
< error
| 1146 42S02 Table 'tablesample.t' doesn't exist
> exec ^
: string
| update tablesample.t set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into tablesample.t values (1)
< error
| 1146 42S02 Table 'tablesample.t' doesn't exist
> exec ^
: string
| insert into tablesample.t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete tablesample.t from tablesample.t
< error
| 1146 42S02 Table 'tablesample.t' doesn't exist
> exec ^
: string
| delete tablesample.t from tablesample.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select tablesample.* from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample.* from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select tablesample . * from t
< error
| 1064 42000 …line 1 column 18 near "tablesample . * from t" 
> exec ^
: string
| select @@ tablesample
< error
| 1064 42000 …line 1 column 9 near "@@ tablesample" 
> exec ^
: string
| select @@ table
< error
| 1064 42000 …line 1 column 9 near "@@ table" 
> exec ^
: string
| select @@table.select
< error
| 1193 HY000 Unknown system variable 'table.select'
> exec ^
: string
| select @@a.table
< error
| 1193 HY000 Unknown system variable 'a.table'
> exec ^
: string
| set @@tablesample = 1
< ok 0 0
> exec ^
: string
| set @@table = 1
< ok 0 0
> exec ^
: string
| set @@session.tablesample = 1
< ok 0 0
> exec ^
: string
| set session tablesample = 1
< error
| 1064 42000 …line 1 column 23 near "tablesample = 1" 
> exec ^
: string
| set tablesample = 1
< error
| 1064 42000 …line 1 column 15 near "tablesample = 1" 
> exec ^
: string
| set global tablesample = 1
< error
| 1064 42000 …line 1 column 22 near "tablesample = 1" 
> exec ^
: string
| select 1. tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1. tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1.tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1.tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1. tablesample
< ok 0 0
> exec ^
: string
| select 1.5 tablesample
< error
| 1064 42000 …line 1 column 22 near "tablesample" 
> exec ^
: string
| select .5 tablesample
< error
| 1064 42000 …line 1 column 21 near "tablesample" 
> exec ^
: string
| select t.a. tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.a. tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a from t  .  tablesample
< error
| 1146 42S02 Table 't.tablesample' doesn't exist
> exec ^
: string
| select a from t  .  tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select tablesample	.a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample	.a from t" 
> exec ^
: string
| select t.	tablesample from t
< error
| 1064 42000 …line 1 column 21 near "tablesample from t" 
> exec ^
: string
| select tablesample	.a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample	.a from t" 
> exec ^
: string
| select t.	tablesample from t
< error
| 1064 42000 …line 1 column 21 near "tablesample from t" 
> exec ^
: string
| select t.
| tablesample from t
< error
| 1064 42000 …line 2 column 12 near "tablesample from t" 
> exec ^
: string
| select tablesample
| .a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample
| .a from t" 
> exec ^
: string
| select t./**/tablesample from t
< error
| 1064 42000 …line 1 column 24 near "tablesample from t" 
> exec ^
: string
| select tablesample/**/.a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample/**/.a from t" 
> exec ^
: string
| select tablesample  .a from t
< error
| 1064 42000 …line 1 column 18 near "tablesample  .a from t" 
> exec ^
: string
| select t  .  tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t  .  tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db  .  tablesample
< error
| 1146 42S02 Table 'db.tablesample' doesn't exist
> exec ^
: string
| select * from db  .  tablesample
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select from.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select from.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from from.t
< error
| 1146 42S02 Table 'from.t' doesn't exist
> exec ^
: string
| select * from from.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select @@from
< error
| 1193 HY000 Unknown system variable 'from'
