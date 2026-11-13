-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestNextValueColumnNameFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| alter table t add column a int default next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int default next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter a set default (next value for s)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t alter a set default (next value for s)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter a set default next value for s
< error
| 1064 42000 …line 1 column 38 near "next value for s" 
> exec ^
: string
| alter table t alter a set default nextval(s)
< error
| 1064 42000 …line 1 column 41 near "nextval(s)" 
> exec ^
: string
| alter table t modify a int default next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify a int default next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table next (a int)
< ok 0 0
> exec ^
: string
| create table t (a int check (a > next value for s))
< ok 0 0
> exec ^
: string
| create table t (a int default ((next value for s)))
< ok 0 0
> exec ^
: string
| create table t (a int default (next value for s))
< ok 0 0
> exec ^
: string
| create table t (a int default (nextval(s)))
< ok 0 0
> exec ^
: string
| create table t (a int default next value for a.s)
< ok 0 0
> exec ^
: string
| create table t (a int default next value for s on update now())
< ok 0 0
> exec ^
: string
| create table t (a int default next value for s primary key)
< ok 0 0
> exec ^
: string
| create table t (a int default next value for s)
< ok 0 0
> exec ^
: string
| create table t (a int default next value for)
< error
| 1064 42000 …line 1 column 45 near ")" 
> exec ^
: string
| create table t (a int default next value)
< error
| 1064 42000 …line 1 column 41 near ")" 
> exec ^
: string
| create table t (a int default nextval(s))
< ok 0 0
> exec ^
: string
| create table t (a int generated always as (next value for s))
< ok 0 0
> exec ^
: string
| create table t (a int primary key default next value for s, b int)
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update next value for s)
< error
| 1064 42000 …line 1 column 42 near "next value for s)" 
> exec ^
: string
| create table t (a timestamp on update nextval(s))
< error
| 1064 42000 …line 1 column 45 near "nextval(s))" 
> exec ^
: string
| create table t (next int, value int, `for` int)
< ok 0 0
> exec ^
: string
| create table t2 as select next value for s
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t2 as select next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create view v as select next value for s
< ok 0 0
> exec ^
: string
| delete from t where a = next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t where a = next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain select next value for s
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| explain select next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) values (next value for s), (next value for s)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) values (next value for s), (next value for s)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t select next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t select next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a = next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a = next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (next value for s)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (next value for s)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| next value for s
< error
| 1064 42000 …line 1 column 4 near "next value for s" 
> exec ^
: string
| select ((next value for a))
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select ((next value for a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (next value for a)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select (next value for a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t order by next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t order by next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where a = next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where a = next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 + next value for a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select 1 + next value for a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select NEXT VALUE FOR a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select NEXT VALUE FOR a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `next` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `next` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `next` value from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `next` value from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(next value for s) from t
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| select count(next value for s) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select if(1, next value for a, 2)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select if(1, next value for a, 2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select lastval(a)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select lastval(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next
< error
| 1054 42S22 Unknown column 'next' in 'field list'
> exec ^
: string
| select next
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next /* c */ value /* c */ for /* c */ a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next /* c */ value /* c */ for /* c */ a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next VALUE For a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next VALUE For a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next `value` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select next `value` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next a
< error
| 1054 42S22 Unknown column 'next' in 'field list'
> exec ^
: string
| select next a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select next from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| select next value a
< error
| 1064 42000 …line 1 column 19 near "a" 
> exec ^
: string
| select next value for
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| select next value for 'a'
< error
| 1064 42000 …line 1 column 25 near "'a'" 
> exec ^
: string
| select next value for 1
< error
| 1064 42000 …line 1 column 23 near "1" 
> exec ^
: string
| select next value for `a`
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for `a`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for `a`.`b`
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| select next value for `a`.`b`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a + 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a + 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a . b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| select next value for a . b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a . b . c
< error
| 1064 42000 …line 1 column 29 near ". c" 
> exec ^
: string
| select next value for a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| select next value for a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a as b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a as b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a b
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a for update
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a from t
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a, 1
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select next value for a, 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a.b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| select next value for a.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a.b from t
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| select next value for a.b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for a.b.c
< error
| 1064 42000 …line 1 column 26 near ".c" 
> exec ^
: string
| select next value for a.select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| select next value for a.select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value for select
< error
| 1064 42000 …line 1 column 28 near "select" 
> exec ^
: string
| select next value for transaction
< error
| 1146 42S02 Table 'test.transaction' doesn't exist
> exec ^
: string
| select next value for transaction
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next value from t
< error
| 1064 42000 …line 1 column 22 near "from t" 
> exec ^
: string
| select next value(1)
< error
| 1064 42000 …line 1 column 18 near "(1)" 
> exec ^
: string
| select next value, a from t
< error
| 1064 42000 …line 1 column 18 near ", a from t" 
> exec ^
: string
| select next, value from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select next, value from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select next.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select next.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval
< error
| 1054 42S22 Unknown column 'nextval' in 'field list'
> exec ^
: string
| select nextval
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select nextval from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval('a')
< error
| 1064 42000 …line 1 column 18 near "'a')" 
> exec ^
: string
| select nextval()
< error
| 1064 42000 …line 1 column 16 near ")" 
> exec ^
: string
| select nextval(*)
< error
| 1064 42000 …line 1 column 17 near ")" 
> exec ^
: string
| select nextval(1 + 1)
< error
| 1064 42000 …line 1 column 16 near "1 + 1)" 
> exec ^
: string
| select nextval(1)
< error
| 1064 42000 …line 1 column 16 near "1)" 
> exec ^
: string
| select nextval(`a`)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select nextval(`a`)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval(a + 1)
< error
| 1064 42000 …line 1 column 18 near "+ 1)" 
> exec ^
: string
| select nextval(a . select)
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| select nextval(a . select)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval(a)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select nextval(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval(a, b)
< error
| 1064 42000 …line 1 column 17 near ", b)" 
> exec ^
: string
| select nextval(a.*)
< error
| 1064 42000 …line 1 column 18 near "*)" 
> exec ^
: string
| select nextval(a.b)
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| select nextval(a.b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select nextval(a.b.c)
< error
| 1064 42000 …line 1 column 19 near ".c)" 
> exec ^
: string
| select nextval(distinct a)
< error
| 1064 42000 …line 1 column 23 near "distinct a)" 
> exec ^
: string
| select setval(a, 1)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec ^
: string
| select setval(a, 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select value from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select value from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set @a = next value for s
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| set @a = next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set a = next value for s
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| set a = next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a = next value for s
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a = next value for s
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `nextval`(1)
< ok 0 0
> exec ^
: string
| select `nextval`()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'nextval'
> exec ^
: string
| select `nextval`()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `nextval`(a, b)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select `nextval`(a, b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `nextval`('a')
< ok 0 0
> exec ^
: string
| select `nextval`(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| select `nextval`(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `NEXTVAL`(a.b.c)
< error
| 1054 42S22 Unknown column 'a.b.c' in 'field list'
> exec ^
: string
| select `NEXTVAL`(a.b.c)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select d.nextval(1)
< error
| 1305 42000 FUNCTION d.nextval does not exist
> exec ^
: string
| select d.nextval(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `d`.`nextval`(1)
< error
| 1305 42000 FUNCTION d.nextval does not exist
> exec ^
: string
| select `d`.`nextval`(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a timestamp on update now())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now(3))
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now(1 + 1))
< error
| 1064 42000 …line 1 column 45 near "+ 1))" 
> exec ^
: string
| create table t (a timestamp on update (now()))
< error
| 1064 42000 …line 1 column 39 near "(now()))" 
> exec ^
: string
| create table t (a timestamp on update ((now())))
< error
| 1064 42000 …line 1 column 39 near "((now())))" 
> exec ^
: string
| create table t (a timestamp on update (now(3)))
< error
| 1064 42000 …line 1 column 39 near "(now(3)))" 
> exec ^
: string
| create table t (a timestamp on update now)
< error
| 1064 42000 …line 1 column 41 near "now)" 
> exec ^
: string
| create table t (a timestamp on update `now`())
< error
| 1064 42000 …line 1 column 43 near "`now`())" 
> exec ^
: string
| create table t (a timestamp on update current_timestamp)
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update current_timestamp())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update current_timestamp(3))
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update (current_timestamp))
< error
| 1064 42000 …line 1 column 39 near "(current_timestamp))" 
> exec ^
: string
| create table t (a timestamp on update (current_timestamp()))
< error
| 1064 42000 …line 1 column 39 near "(current_timestamp()))" 
> exec ^
: string
| create table t (a timestamp on update localtime)
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update localtime())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update localtime(3))
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update localtimestamp)
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update localtimestamp())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update localtimestamp(3))
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update (localtime))
< error
| 1064 42000 …line 1 column 39 near "(localtime))" 
> exec ^
: string
| create table t (a timestamp on update sysdate())
< error
| 1064 42000 …line 1 column 45 near "sysdate())" 
> exec ^
: string
| create table t (a timestamp on update curdate())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update current_date)
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update current_time)
< error
| 1064 42000 …line 1 column 50 near "current_time)" 
> exec ^
: string
| create table t (a timestamp on update utc_timestamp)
< error
| 1064 42000 …line 1 column 51 near "utc_timestamp)" 
> exec ^
: string
| create table t (a timestamp on update utc_timestamp())
< error
| 1064 42000 …line 1 column 51 near "utc_timestamp())" 
> exec ^
: string
| create table t (a timestamp on update current_user)
< error
| 1064 42000 …line 1 column 50 near "current_user)" 
> exec ^
: string
| create table t (a timestamp on update 1)
< error
| 1064 42000 …line 1 column 39 near "1)" 
> exec ^
: string
| create table t (a timestamp on update 'x')
< error
| 1064 42000 …line 1 column 41 near "'x')" 
> exec ^
: string
| create table t (a timestamp on update null)
< error
| 1064 42000 …line 1 column 42 near "null)" 
> exec ^
: string
| create table t (a timestamp on update true)
< error
| 1064 42000 …line 1 column 42 near "true)" 
> exec ^
: string
| create table t (a timestamp on update default)
< error
| 1064 42000 …line 1 column 45 near "default)" 
> exec ^
: string
| create table t (a timestamp on update b)
< error
| 1064 42000 …line 1 column 39 near "b)" 
> exec ^
: string
| create table t (a timestamp on update (b))
< error
| 1064 42000 …line 1 column 39 near "(b))" 
> exec ^
: string
| create table t (a timestamp on update foo())
< error
| 1064 42000 …line 1 column 41 near "foo())" 
> exec ^
: string
| create table t (a timestamp on update now() + 1)
< error
| 1064 42000 …line 1 column 45 near "+ 1)" 
> exec ^
: string
| create table t (a timestamp on update (now() + 1))
< error
| 1064 42000 …line 1 column 39 near "(now() + 1))" 
> exec ^
: string
| create table t (a timestamp on update 1 + now())
< error
| 1064 42000 …line 1 column 39 near "1 + now())" 
> exec ^
: string
| create table t (a timestamp on update (next value for s))
< error
| 1064 42000 …line 1 column 39 near "(next value for s))" 
> exec ^
: string
| create table t (a timestamp on update ((next value for s)))
< error
| 1064 42000 …line 1 column 39 near "((next value for s)))" 
> exec ^
: string
| create table t (a timestamp on update (nextval(s)))
< error
| 1064 42000 …line 1 column 39 near "(nextval(s)))" 
> exec ^
: string
| create table t (a timestamp on update next value for s + 1)
< error
| 1064 42000 …line 1 column 42 near "next value for s + 1)" 
> exec ^
: string
| create table t (a timestamp on update (next value for s) + 1)
< error
| 1064 42000 …line 1 column 39 near "(next value for s) + 1)" 
> exec ^
: string
| create table t (a timestamp on update if(1, now(), now()))
< error
| 1064 42000 …line 1 column 40 near "if(1, now(), now()))" 
> exec ^
: string
| create table t (a timestamp on update case when 1 then now() end)
< error
| 1064 42000 …line 1 column 42 near "case when 1 then now() end)" 
> exec ^
: string
| create table t (a timestamp on update (select now()))
< error
| 1064 42000 …line 1 column 39 near "(select now()))" 
> exec ^
: string
| create table t (a timestamp default now() on update now())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now() default now())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now() not null)
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update now() comment 'x')
< ok 0 0
> exec ^
: string
| alter table t add column a timestamp on update now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a timestamp on update now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify a timestamp on update (next value for s)
< error
| 1064 42000 …line 1 column 44 near "(next value for s)" 
> exec ^
: string
| alter table t change a b timestamp on update current_timestamp(6)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change a b timestamp on update current_timestamp(6)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a timestamp on update current_date())
< ok 0 0
> exec ^
: string
| create table t (a timestamp on update curdate)
< error
| 1064 42000 …line 1 column 45 near "curdate)" 
> exec ^
: string
| create table t (a timestamp on update curdate(3))
< error
| 1064 42000 …line 1 column 47 near "3))" 
> exec ^
: string
| create table t (a timestamp on update current_date(3))
< error
| 1064 42000 …line 1 column 52 near "3))" 
> exec ^
: string
| create table t (a timestamp on update now('3'))
< error
| 1064 42000 …line 1 column 45 near "'3'))" 
> exec ^
: string
| create table t (a timestamp on update now(-1))
< error
| 1064 42000 …line 1 column 43 near "-1))" 
> exec ^
: string
| create table t (a timestamp on update now(1.5))
< error
| 1064 42000 …line 1 column 45 near "1.5))" 
> exec ^
: string
| create table t (a timestamp on update now(1, 2))
< error
| 1064 42000 …line 1 column 44 near ", 2))" 
> exec ^
: string
| create table t (a timestamp on update now(a))
< error
| 1064 42000 …line 1 column 43 near "a))" 
> exec ^
: string
| create table t (a timestamp on update current_timestamp(a))
< error
| 1064 42000 …line 1 column 57 near "a))" 
> exec ^
: string
| create table t (a timestamp on update localtime(1 + 1))
< error
| 1064 42000 …line 1 column 51 near "+ 1))" 
> exec ^
: string
| create table t (a timestamp on update current_timestamp + 1)
< error
| 1064 42000 …line 1 column 57 near "+ 1)" 
> exec ^
: string
| create table t (a timestamp on update now() + 1 not null)
< error
| 1064 42000 …line 1 column 45 near "+ 1 not null)" 
> exec ^
: string
| create table t (a timestamp on update now() on update now())
< ok 0 0
> exec ^
: string
| create table t (a int references p (x) on update cascade)
< ok 0 0
> exec ^
: string
| create table t (a int references p (x) on update set null)
< ok 0 0
> exec ^
: string
| create table t (a int, foreign key (a) references p (x) on update restrict)
< ok 0 0
> exec ^
: string
| select `nextval`(1), `nextval`()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'nextval'
> exec ^
: string
| select `nextval`(1), `nextval`()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `lastval`(1)
< ok 0 0
> exec ^
: string
| select `next`(1)
< error
| 1305 42000 FUNCTION test.next does not exist
> exec ^
: string
| select `next`(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add (t.a int)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add (t.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column (t.a int, t.b int)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column (t.a int, t.b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a int after a.b.c.d
< error
| 1064 42000 …line 1 column 43 near ".d" 
> exec ^
: string
| alter table t add column a int after b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int after b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a int after d.t.b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int after d.t.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a int after t.b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int after t.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a int first
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int first
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column a.b.c.d int
< error
| 1064 42000 …line 1 column 31 near ".d int" 
> exec ^
: string
| alter table t add column d.t.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column d.t.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column if not exists t.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column if not exists t.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column t.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column t.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add column x.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column x.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add index (t.a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add index (t.a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter column t.a drop default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t alter column t.a drop default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter column t.a set default 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t alter column t.a set default 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter t.a drop default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t alter t.a drop default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter t.a set default 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t alter t.a set default 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t alter t.a set invisible
< error
| 1064 42000 …line 1 column 37 near "invisible" 
> exec ^
: string
| alter table t alter t.a set visible
< error
| 1064 42000 …line 1 column 35 near "visible" 
> exec ^
: string
| alter table t change a b int after t.c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change a b int after t.c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change a t.b int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change a t.b int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change column d.t.a x.b int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change column d.t.a x.b int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change column if exists t.a b int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change column if exists t.a b int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change t.a b int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change t.a b int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t change t.a t.b int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t change t.a t.b int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t drop column if exists t.a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t drop column if exists t.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t drop column t.a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t drop column t.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t drop t.a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t drop t.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify a int after d.t.b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify a int after d.t.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify a int after t.b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify a int after t.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify column d.t.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify column d.t.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify column if exists t.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify column if exists t.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify column x.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify column x.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify t.a int
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify t.a int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t rename column a to t.b
< error
| 1064 42000 …line 1 column 35 near ".b" 
> exec ^
: string
| alter table t rename column t.a to b
< error
| 1064 42000 …line 1 column 30 near ".a to b" 
> exec ^
: string
| create index i on t (d.t.a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index i on t (d.t.a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create index i on t (t.a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index i on t (t.a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table d.t (d.t.a int)
< ok 0 0
> exec ^
: string
| create table d.t (d.x.a int)
< error
| 1103 42000 Incorrect table name 'x'
> exec ^
: string
| create table d.t (d.x.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table d.t (t.a int)
< ok 0 0
> exec ^
: string
| create table d.t (x.t.a int)
< error
| 1102 42000 Incorrect database name 'x'
> exec ^
: string
| create table d.t (x.t.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (.a int)
< error
| 1064 42000 …line 1 column 17 near ".a int)" 
> exec ^
: string
| create table t (`t`.`a` int)
< ok 0 0
> exec ^
: string
| create table t (a . b int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a . b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a . select int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a . select int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a .b int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a .b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int references p ())
< error
| 1064 42000 …line 1 column 37 near "))" 
> exec ^
: string
| create table t (a int references p (a.b.c.d))
< error
| 1064 42000 …line 1 column 42 near ".d))" 
> exec ^
: string
| create table t (a int references p (d.t.x))
< ok 0 0
> exec ^
: string
| create table t (a int references p (x, y.z))
< ok 0 0
> exec ^
: string
| create table t (a int references p (x.y))
< ok 0 0
> exec ^
: string
| create table t (a int, b int, foreign key (t.a) references p (x))
< ok 0 0
> exec ^
: string
| create table t (a int, foreign key (a) references p (x.y))
< ok 0 0
> exec ^
: string
| create table t (a int, index i (t.a))
< ok 0 0
> exec ^
: string
| create table t (a int, key ((t.a)))
< ok 0 0
> exec ^
: string
| create table t (a int, key (a.b.c.d))
< error
| 1064 42000 …line 1 column 34 near ".d))" 
> exec ^
: string
| create table t (a int, key (t.a desc))
< ok 0 0
> exec ^
: string
| create table t (a int, key (t.a(10)))
< ok 0 0
> exec ^
: string
| create table t (a int, key (t.a))
< ok 0 0
> exec ^
: string
| create table t (a int, primary key (d.t.a))
< ok 0 0
> exec ^
: string
| create table t (a int, unique (t.a))
< ok 0 0
> exec ^
: string
| create table t (a. b int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a. b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a.b int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a.b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a.b.c int)
< error
| 1102 42000 Incorrect database name 'a'
> exec ^
: string
| create table t (a.b.c int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a.b.c.d int)
< error
| 1064 42000 …line 1 column 22 near ".d int)" 
> exec ^
: string
| create table t (a.global int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a.global int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a.select int)
< error
| 1103 42000 Incorrect table name 'a'
> exec ^
: string
| create table t (a.select int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (b.a int)
< error
| 1103 42000 Incorrect table name 'b'
> exec ^
: string
| create table t (b.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (d.t.a int)
< error
| 1102 42000 Incorrect database name 'd'
> exec ^
: string
| create table t (d.t.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (global.a int)
< error
| 1103 42000 Incorrect table name 'global'
> exec ^
: string
| create table t (global.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (select.a int)
< error
| 1103 42000 Incorrect table name 'select'
> exec ^
: string
| create table t (select.a int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (t. int)
< error
| 1064 42000 …line 1 column 23 near ")" 
> exec ^
: string
| create table t (t.`a` int)
< ok 0 0
> exec ^
: string
| create table t (t.a int default 1 comment 'x')
< ok 0 0
> exec ^
: string
| create table t (t.a int primary key)
< ok 0 0
> exec ^
: string
| create table t (t.a int)
< ok 0 0
> exec ^
: string
| create table t (t.a int) select 1
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t (t.a int) select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (t.a int, db.t.b int)
< error
| 1102 42000 Incorrect database name 'db'
> exec ^
: string
| create table t (t.a int, db.t.b int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (t.a int, primary key (a))
< ok 0 0
> exec ^
: string
| create table t (t.a int, primary key (t.a))
< ok 0 0
> exec ^
: string
| create table t (t.a)
< error
| 1064 42000 …line 1 column 20 near ")" 
> exec ^
: string
| create table t (t.t.t int)
< error
| 1102 42000 Incorrect database name 't'
> exec ^
: string
| create table t (t.t.t int)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create temporary table t (t.a int)
< ok 0 0
