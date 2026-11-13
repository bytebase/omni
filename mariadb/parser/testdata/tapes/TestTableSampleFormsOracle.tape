-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestTableSampleFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| select * from tbl tablesample system (50)
< error
| 1064 42000 …near 'system (50)' at line 1
> exec ^
: string
| select * from tbl tablesample bernoulli (50)
< error
| 1064 42000 …near 'bernoulli (50)' at line 1
> exec ^
: string
| select * from tbl tablesample regions ()
< error
| 1064 42000 …near 'regions ()' at line 1
> exec ^
: string
| select * from tbl tablesample regions (50)
< error
| 1064 42000 …near 'regions (50)' at line 1
> exec ^
: string
| select * from tbl tablesample ()
< error
| 1064 42000 …near '()' at line 1
> exec ^
: string
| select * from tbl tablesample system ()
< error
| 1064 42000 …near 'system ()' at line 1
> exec ^
: string
| select * from tbl tablesample bernoulli ()
< error
| 1064 42000 …near 'bernoulli ()' at line 1
> exec ^
: string
| select * from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| select * from tbl tablesample (50 percent)
< error
| 1064 42000 …near '(50 percent)' at line 1
> exec ^
: string
| select * from tbl tablesample (50 rows)
< error
| 1064 42000 …near '(50 rows)' at line 1
> exec ^
: string
| select * from tbl tablesample (50 percent rows)
< error
| 1064 42000 …near '(50 percent rows)' at line 1
> exec ^
: string
| select * from tbl tablesample (50 rows percent)
< error
| 1064 42000 …near '(50 rows percent)' at line 1
> exec ^
: string
| select * from tbl tablesample (percent)
< error
| 1064 42000 …near '(percent)' at line 1
> exec ^
: string
| select * from tbl tablesample (percent percent)
< error
| 1064 42000 …near '(percent percent)' at line 1
> exec ^
: string
| select * from tbl tablesample (rows)
< error
| 1064 42000 …near '(rows)' at line 1
> exec ^
: string
| select * from tbl tablesample (1+1)
< error
| 1064 42000 …near '(1+1)' at line 1
> exec ^
: string
| select * from tbl tablesample (a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| select * from tbl tablesample (a + b rows)
< error
| 1064 42000 …near '(a + b rows)' at line 1
> exec ^
: string
| select * from tbl tablesample ((select 1))
< error
| 1064 42000 …near '((select 1))' at line 1
> exec ^
: string
| select * from tbl tablesample ('34')
< error
| 1064 42000 …near '('34')' at line 1
> exec ^
: string
| select * from tbl tablesample ('34' rows)
< error
| 1064 42000 …near '('34' rows)' at line 1
> exec ^
: string
| select * from tbl tablesample (50, 50)
< error
| 1064 42000 …near '(50, 50)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1)
< error
| 1064 42000 …near '(50) repeatable (1)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1+1)
< error
| 1064 42000 …near '(50) repeatable (1+1)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable ('a')
< error
| 1064 42000 …near '(50) repeatable ('a')' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable ()
< error
| 1064 42000 …near '(50) repeatable ()' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable
< error
| 1064 42000 …near '(50) repeatable' at line 1
> exec ^
: string
| select * from tbl tablesample () repeatable (1)
< error
| 1064 42000 …near '() repeatable (1)' at line 1
> exec ^
: string
| select * from tbl tablesample regions () repeatable (1)
< error
| 1064 42000 …near 'regions () repeatable (1)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable (1) repeatable (2)
< error
| 1064 42000 …near '(50) repeatable (1) repeatable (2)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) tablesample (50)
< error
| 1064 42000 …near '(50) tablesample (50)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) a
< error
| 1064 42000 …near '(50) a' at line 1
> exec ^
: string
| select * from tbl tablesample (50) as a
< error
| 1064 42000 …near '(50) as a' at line 1
> exec ^
: string
| select * from tbl tablesample
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec ^
: string
| select * from tbl tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from tbl tablesample system
< error
| 1064 42000 …near 'system' at line 1
> exec ^
: string
| select * from tbl tablesample (
< error
| 1064 42000 …near '(' at line 1
> exec ^
: string
| select * from tbl tablesample (50
< error
| 1064 42000 …near '(50' at line 1
> exec ^
: string
| select * from tbl tablesample foo (50)
< error
| 1064 42000 …near 'foo (50)' at line 1
> exec ^
: string
| select * from tbl tablesample 50
< error
| 1064 42000 …near '50' at line 1
> exec ^
: string
| select * from tbl tablesample system 50
< error
| 1064 42000 …near 'system 50' at line 1
> exec ^
: string
| select * from db.tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| select * from tbl a tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl as a tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl `tablesample` tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl partition (p0) tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| select * from tbl partition (p0) a tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) partition (p0)
< error
| 1064 42000 …near '(50) partition (p0)' at line 1
> exec ^
: string
| select * from tbl use index (i) tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl a use index (i) tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) use index (i)
< error
| 1064 42000 …near '(50) use index (i)' at line 1
> exec ^
: string
| select * from tbl partition (p0) a use index (i) tablesample (50) repeatable (1)
< error
| 1064 42000 …near 'tablesample (50) repeatable (1)' at line 1
> exec ^
: string
| select * from tbl as of timestamp now() tablesample (50)
< error
| 1064 42000 …near 'timestamp now() tablesample (50)' at line 1
> exec ^
: string
| select * from (tbl tablesample (50))
< error
| 1064 42000 …near '(50))' at line 1
> exec ^
: string
| select * from (tbl tablesample (50), t2)
< error
| 1064 42000 …near '(50), t2)' at line 1
> exec ^
: string
| select * from (select 1) a tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from (t1 join t2) tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from dual tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from json_table('[]', '$' columns (a int path '$')) a tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from t1 tablesample (10), t2 tablesample (20)
< error
| 1064 42000 …near '(10), t2 tablesample (20)' at line 1
> exec ^
: string
| select * from t1 a tablesample (10) join t2 b tablesample (20) on a.id = b.id
< error
| 1064 42000 …near 'tablesample (10) join t2 b tablesample (20) on a.id = b.id' at line 1
> exec ^
: string
| select * from t1 tablesample (10) left join t2 tablesample (20) on t1.id = t2.id
< error
| 1064 42000 …near '(10) left join t2 tablesample (20) on t1.id = t2.id' at line 1
> exec ^
: string
| select * from t1 tablesample (10) straight_join t2 tablesample (20)
< error
| 1064 42000 …near '(10) straight_join t2 tablesample (20)' at line 1
> exec ^
: string
| select * from t1 tablesample (10) natural join t2 tablesample (20)
< error
| 1064 42000 …near '(10) natural join t2 tablesample (20)' at line 1
> exec ^
: string
| select * from t1 tablesample (10) cross join t2 tablesample (20)
< error
| 1064 42000 …near '(10) cross join t2 tablesample (20)' at line 1
> exec ^
: string
| select * from t1 join t2 tablesample (20) using (id)
< error
| 1064 42000 …near '(20) using (id)' at line 1
> exec ^
: string
| select * from tbl where id > 20 tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) where id > 20
< error
| 1064 42000 …near '(50) where id > 20' at line 1
> exec ^
: string
| select * from tbl tablesample (50) group by a
< error
| 1064 42000 …near '(50) group by a' at line 1
> exec ^
: string
| select * from tbl tablesample (50) having a > 1
< error
| 1064 42000 …near '(50) having a > 1' at line 1
> exec ^
: string
| select * from tbl tablesample (50) order by a
< error
| 1064 42000 …near '(50) order by a' at line 1
> exec ^
: string
| select * from tbl tablesample (50) limit 1
< error
| 1064 42000 …near '(50) limit 1' at line 1
> exec ^
: string
| select * from tbl tablesample (50) for update
< error
| 1064 42000 …near '(50) for update' at line 1
> exec ^
: string
| select * from tbl tablesample (50) union select * from t2 tablesample (50)
< error
| 1064 42000 …near '(50) union select * from t2 tablesample (50)' at line 1
> exec ^
: string
| select * from tbl tablesample bernoulli(50) limit 1 into outfile '/tmp/sample.csv'
< error
| 1064 42000 …near 'bernoulli(50) limit 1 into outfile '/tmp/sample.csv'' at line 1
> exec ^
: string
| select * from tbl tablesample bernoulli(50) order by a, b into outfile '/tmp/sample.csv'
< error
| 1064 42000 …near 'bernoulli(50) order by a, b into outfile '/tmp/sample.csv'' at line 1
> exec ^
: string
| select (select 1 from tbl tablesample (50))
< error
| 1064 42000 …near '(50))' at line 1
> exec ^
: string
| with c as (select * from tbl tablesample (50)) select * from c
< error
| 1064 42000 …near '(50)) select * from c' at line 1
> exec ^
: string
| with c as (select 1) select * from c tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| insert into t select * from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| create table t2 as select * from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| create view v as select * from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| replace into t select * from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| explain select * from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| update tbl tablesample (50) set a = 1
< error
| 1064 42000 …near '(50) set a = 1' at line 1
> exec ^
: string
| update tbl a tablesample (50) set a = 1
< error
| 1064 42000 …near 'tablesample (50) set a = 1' at line 1
> exec ^
: string
| update tbl tablesample (50) join t2 on 1 = 1 set a = 1
< error
| 1064 42000 …near '(50) join t2 on 1 = 1 set a = 1' at line 1
> exec ^
: string
| delete from tbl tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| delete from tbl a tablesample (50) where 1 = 1
< error
| 1064 42000 …near 'tablesample (50) where 1 = 1' at line 1
> exec ^
: string
| delete t from tbl t tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| delete from t using tbl t tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| delete t from t1 t join t2 tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| lock tables tbl tablesample (50) read
< error
| 1064 42000 …near '(50) read' at line 1
> exec ^
: string
| lock tables tbl read tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| table tbl tablesample (50)
< error
| 1064 42000 …near 'table tbl tablesample (50)' at line 1
> exec ^
: string
| insert into tbl tablesample (50) values (1)
< error
| 1064 42000 …near 'tablesample (50) values (1)' at line 1
> exec ^
: string
| analyze table tbl tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| admin check table tbl tablesample (50)
< error
| 1064 42000 …near 'admin check table tbl tablesample (50)' at line 1
> exec ^
: string
| select * from tablesample
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| select * from tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from tablesample a
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| select * from tablesample a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample.a from tablesample
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| select tablesample.a from tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample, a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample, a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t where tablesample = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t where tablesample = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t order by tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t order by tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t group by tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t group by tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample()
< error
| 1305 42000 FUNCTION test.tablesample does not exist
> exec ^
: string
| select tablesample()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample(1)
< error
| 1305 42000 FUNCTION test.tablesample does not exist
> exec ^
: string
| select tablesample(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 tablesample
< ok 0 0
> exec ^
: string
| select 1 as tablesample
< ok 0 0
> exec ^
: string
| select * from t tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample where a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample where a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample, t2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample, t2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample join t2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample join t2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample order by 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample order by 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample limit 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample group by a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample group by a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample use index (i)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample use index (i)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t partition (p0) tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t partition (p0) tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from (select 1) tablesample
< ok 0 0
> exec ^
: string
| select * from (select 1) as tablesample
< ok 0 0
> exec ^
: string
| select * from t1 join t2 tablesample on 1 = 1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| select * from t1 join t2 tablesample on 1 = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| with tablesample as (select 1) select * from tablesample
< ok 0 0
> exec ^
: string
| update tablesample set a = 1
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| update tablesample set a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t tablesample set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t tablesample set a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t as tablesample set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t as tablesample set a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t set tablesample = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set tablesample = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from tablesample
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| delete from tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete tablesample from t
< error
| 1109 42S02 Unknown table 'tablesample' in MULTI DELETE
> exec ^
: string
| delete tablesample from t
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'tablesample' in MULTI DELETE
> exec ^
: string
| delete tablesample from tablesample
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| delete tablesample from tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from t tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from t tablesample where a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t tablesample where a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from t as tablesample where a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t as tablesample where a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from t tablesample order by a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t tablesample order by a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from t tablesample limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t tablesample limit 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete tablesample.* from t tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete tablesample.* from t tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete t.* from t tablesample
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| delete t.* from t tablesample
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| insert into tablesample values (1)
< error
| 1146 42S02 Table 'test.tablesample' doesn't exist
> exec ^
: string
| insert into tablesample values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into t (tablesample) values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (tablesample) values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table tablesample (a int)
< ok 0 0
> exec ^
: string
| create table t (tablesample int)
< ok 0 0
> exec ^
: string
| create table db.tablesample (a int)
< ok 0 0
> exec ^
: string
| create index tablesample on t (a)
< ok 0 0
> exec ^
: string
| create index i on t (tablesample)
< ok 0 0
> exec ^
: string
| create database tablesample
< ok 0 0
> exec ^
: string
| create view tablesample as select 1
< ok 0 0
> exec ^
: string
| drop table tablesample
< ok 0 0
> exec ^
: string
| alter table t add column tablesample int
< ok 0 0
> exec ^
: string
| alter table t rename to tablesample
< ok 0 0
> exec ^
: string
| alter table tablesample add column a int
< ok 0 0
> exec ^
: string
| lock tables tablesample read
< ok 0 0
> exec ^
: string
| lock tables t tablesample read
< ok 0 0
> exec ^
: string
| lock tables t as tablesample read
< ok 0 0
> exec ^
: string
| truncate table tablesample
< ok 0 0
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
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1 from t where a in (select tablesample from t2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t tablesample force index (i)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t tablesample force index (i)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as tablesample tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
> exec ^
: string
| select * from tablesample tablesample (50)
< error
| 1064 42000 …near '(50)' at line 1
> exec ^
: string
| select tablesample from tablesample tablesample tablesample (50)
< error
| 1064 42000 …near 'tablesample (50)' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from tbl tablesample (select 1)
< error
| 1064 42000 …near '(select 1)' at line 1
> exec ^
: string
| select * from tbl tablesample `system` (50)
< error
| 1064 42000 …near '`system` (50)' at line 1
> exec ^
: string
| select * from tbl tablesample `bernoulli` (50)
< error
| 1064 42000 …near '`bernoulli` (50)' at line 1
> exec ^
: string
| select * from tbl tablesample `regions` ()
< error
| 1064 42000 …near '`regions` ()' at line 1
> exec ^
: string
| select * from tbl tablesample (50 `percent`)
< error
| 1064 42000 …near '(50 `percent`)' at line 1
> exec ^
: string
| select * from tbl tablesample (50 `rows`)
< error
| 1064 42000 …near '(50 `rows`)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) `repeatable` (1)
< error
| 1064 42000 …near '(50) `repeatable` (1)' at line 1
> exec ^
: string
| select * from tbl tablesample (@a := 1)
< error
| 1064 42000 …near '(@a := 1)' at line 1
> exec ^
: string
| select * from tbl tablesample (50) repeatable (@a := 1)
< error
| 1064 42000 …near '(50) repeatable (@a := 1)' at line 1
> exec ^
: string
| select * from tbl TABLESAMPLE SYSTEM (50 PERCENT) REPEATABLE (1)
< error
| 1064 42000 …near 'SYSTEM (50 PERCENT) REPEATABLE (1)' at line 1
> exec ^
: string
| select * from tbl tablesample region ()
< error
| 1064 42000 …near 'region ()' at line 1
> exec ^
: string
| select * from tbl tablesample (50) (select 1)
< error
| 1064 42000 …near '(50) (select 1)' at line 1
> exec ^
: string
| select * from t1 tablesample, t2 tablesample
< error
| 1066 42000 Not unique table/alias: 'tablesample'
> exec ^
: string
| select * from t1 tablesample, t2 tablesample
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'tablesample'
> exec ^
: string
| select * from tbl a tablesample
< error
| 1064 42000 …near 'tablesample' at line 1
> exec ^
: string
| select 1 from t where tablesample
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1 from t where tablesample
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count(tablesample) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(tablesample) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t1 tablesample (10) join t2 tablesample (20) join t3 tablesample (30)
< error
| 1064 42000 …near '(10) join t2 tablesample (20) join t3 tablesample (30)' at line 1
> exec ^
: string
| select * from (t1 tablesample (10) join t2) join t3 tablesample (30) on 1 = 1
< error
| 1064 42000 …near '(10) join t2) join t3 tablesample (30) on 1 = 1' at line 1
> exec ^
: string
| select * from t1 join (t2 tablesample (20)) on 1 = 1
< error
| 1064 42000 …near '(20)) on 1 = 1' at line 1
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
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select tablesample .a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample .a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from tablesample .t
< error
| 1146 42S02 Table 'tablesample.t' doesn't exist
> exec ^
: string
| select * from tablesample .t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample . * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample . * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @@tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set @@session.tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set session tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
> exec ^
: string
| set global tablesample = 1
< error
| 1193 HY000 Unknown system variable 'tablesample'
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1. tablesample
< ok 0 0
> exec ^
: string
| select 1.5 tablesample
< ok 0 0
> exec ^
: string
| select .5 tablesample
< ok 0 0
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample	.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample	.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.	tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.	tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample	.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample	.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.	tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.	tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.
| tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.
| tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 3
> exec ^
: string
| select tablesample
| .a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample
| .a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 3
> exec ^
: string
| select t./**/tablesample from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t./**/tablesample from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample/**/.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample/**/.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select tablesample  .a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select tablesample  .a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
