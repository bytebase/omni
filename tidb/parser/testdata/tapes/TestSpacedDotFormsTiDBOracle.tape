-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSpacedDotFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| alter table a . select add column x int
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| alter table a . select add column x int
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| call a . select()
< error
| 1305 42000 FUNCTION a.select does not exist
> exec ^
: string
| call a . select()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create database a . select
< error
| 1064 42000 …line 1 column 19 near ". select" 
> exec ^
: string
| create index a . select on t (a)
< error
| 1064 42000 …line 1 column 16 near ". select on t (a)" 
> exec ^
: string
| create procedure a . select() begin end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure a . select() begin end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create sequence a . select
< ok 0 0
> exec ^
: string
| create table a . select (x int)
< ok 0 0
> exec ^
: string
| create table t (a int, b int as (t . select))
< ok 0 0
> exec ^
: string
| create table t (a int, primary key (a), foreign key (a) references b . select (c))
< ok 0 0
> exec ^
: string
| create trigger a . select before insert on t for each row set @a = 1
< error
| 1064 42000 …line 1 column 14 near "trigger a . select before insert on t for each row set @a = 1" 
> exec ^
: string
| create view a . select as select 1
< ok 0 0
> exec ^
: string
| delete from a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| delete from a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop table a . select
< ok 0 0
> exec ^
: string
| drop table a . table
< ok 0 0
> exec ^
: string
| drop view a . select
< ok 0 0
> exec ^
: string
| explain select a . select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain select a . select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| grant all on a . from . select to u
< error
| 1064 42000 …line 1 column 23 near ". select to u" 
> exec ^
: string
| grant select on a . * to u
< ok 0 0
> exec ^
: string
| grant select on a . select to u
< ok 0 0
> exec ^
: string
| insert into a . select values (1)
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| insert into a . select values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a . select) values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a . select) values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile 'x' into table a . select
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile 'x' into table a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| rename table a . select to b . from
< ok 0 0
> exec ^
: string
| select "a." from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select "a." from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 'a.' from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 'a.' from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from a . b . select
< error
| 1064 42000 …line 1 column 21 near ". select" 
> exec ^
: string
| select * from a . from
< error
| 1146 42S02 Table 'a.from' doesn't exist
> exec ^
: string
| select * from a . from
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from a . from where 1
< error
| 1146 42S02 Table 'a.from' doesn't exist
> exec ^
: string
| select * from a . from where 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| select * from a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from a . select . b
< error
| 1064 42000 …line 1 column 26 near ". b" 
> exec ^
: string
| select * from db . of
< error
| 1146 42S02 Table 'db.of' doesn't exist
> exec ^
: string
| select * from db . of
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
| select * from t a join b on a . join = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a join b on a . join = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t a join b on a . on = b . on
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a join b on a . on = b . on
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where a . from = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where a . from = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where a . in (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where a . in (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where a . in in (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where a . in in (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 . from t
< error
| 1064 42000 …line 1 column 10 near ". from t" 
> exec ^
: string
| select @a . select
< error
| 1064 42000 …line 1 column 11 near ". select" 
> exec ^
: string
| select `a.` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `a.` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `a.` select from t
< error
| 1064 42000 …line 1 column 18 near "select from t" 
> exec ^
: string
| select a .
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| select a .	select from t
< error
| 1064 42000 …line 1 column 17 near "select from t" 
> exec ^
: string
| select a .  select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a .  select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . /* c */ b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . /* c */ b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . 1 from t
< error
| 1064 42000 …line 1 column 12 near "1 from t" 
> exec ^
: string
| select a . as b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . as b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . b (1) from t
< error
| 1305 42000 FUNCTION a.b does not exist
> exec ^
: string
| select a . b (1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . b . * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . b . * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . b . c from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . b . c from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . b . select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . b . select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . case from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . case from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . count (1) from t
< error
| 1305 42000 FUNCTION a.count does not exist
> exec ^
: string
| select a . count (1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . current_user from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . current_user from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . default from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . default from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . from from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . from from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . if (1, 2, 3) from t
< error
| 1305 42000 FUNCTION a.if does not exist
> exec ^
: string
| select a . if (1, 2, 3) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . interval from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . interval from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . max (1) from t
< error
| 1305 42000 FUNCTION a.max does not exist
> exec ^
: string
| select a . max (1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . not from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . not from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . now() from t
< error
| 1305 42000 FUNCTION a.now does not exist
> exec ^
: string
| select a . now() from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . null from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . null from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . select (1) from t
< error
| 1305 42000 FUNCTION a.select does not exist
> exec ^
: string
| select a . select (1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . select . b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . select . b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . select(1) from t
< error
| 1305 42000 FUNCTION a.select does not exist
> exec ^
: string
| select a . select(1) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . select, a . from from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . select, a . from from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . true from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a . true from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a . user () from t
< error
| 1305 42000 FUNCTION a.user does not exist
> exec ^
: string
| select a . user () from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a. 
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| select a.  select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.  select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a. * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a. * from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a. /* c */ from t
< error
| 1064 42000 …line 1 column 22 near "from t" 
> exec ^
: string
| select a. as from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a. as from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a. b. c from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a. b. c from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a. select b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a. select b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.b . select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.b . select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a/* . */ from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a/* . */ from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select count(a . select) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count(a . select) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select interval 1 foo + now()
< error
| 1064 42000 …line 1 column 21 near "foo + now()" 
> exec ^
: string
| select lastval(a . select)
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| select lastval(a . select)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select max(a . max) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select max(a . max) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
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
| select now() + interval 1 day
< ok 0 0
> exec ^
: string
| select o . from t
< error
| 1054 42S22 Unknown column 'o.from' in 'field list'
> exec ^
: string
| select o . from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select o .from t
< error
| 1054 42S22 Unknown column 'o.from' in 'field list'
> exec ^
: string
| select o .from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select o. a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select o. a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select o. from from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select o. from from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select o. from t
< error
| 1054 42S22 Unknown column 'o.from' in 'field list'
> exec ^
: string
| select o. from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select o. select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select o. select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select setval(a . select, 1)
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| select setval(a . select, 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t . of from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t . of from t
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
| select t. * , b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t. * , b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t. from from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t. from from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t. of from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t. of from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.a. from t
< error
| 1054 42S22 Unknown column 't.a.from' in 'field list'
> exec ^
: string
| select t.a. from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| truncate table a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| truncate table a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update a . select set b = 1
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| update a . select set b = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a . select = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a . select = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| use a . select
< error
| 1064 42000 …line 1 column 7 near ". select" 
> exec ^
: string
| with a as (select 1) select a . select from a
< error
| 1054 42S22 Unknown column 'a.select' in 'field list'
> exec ^
: string
| with a as (select 1) select a . select from a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a decimal(10,2) default 1. comment 'x')
< error
| 1064 42000 …line 1 column 50 near "comment 'x')" 
> exec ^
: string
| create table t (a decimal(10,2) default 1. not null)
< error
| 1064 42000 …line 1 column 46 near "not null)" 
> exec ^
: string
| insert into t values (1.) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1.) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select .5 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select .5 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1. + 2
< ok 0 0
> exec ^
: string
| select 1. a
< ok 0 0
> exec ^
: string
| select 1. a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1. a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1. and 2
< error
| 1064 42000 …line 1 column 15 near "2" 
> exec ^
: string
| select 1. as from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1. as from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1. from
< ok 0 0
> exec ^
: string
| select 1. from t
< error
| 1064 42000 …line 1 column 16 near "t" 
> exec ^
: string
| select 1. select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1. select from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1. tablesample
< ok 0 0
> exec ^
: string
| select 1. union select 2
< error
| 1064 42000 …line 1 column 22 near "select 2" 
> exec ^
: string
| select 1.0 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1.0 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1.e1 from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1.e1 from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a from t group by 1. having 1
< error
| 1064 42000 …line 1 column 34 near "having 1" 
> exec ^
: string
| select a from t limit 1. offset 2
< error
| 1064 42000 …line 1 column 24 near "1. offset 2" 
> exec ^
: string
| select a from t order by 1. desc
< error
| 1064 42000 …line 1 column 32 near "desc" 
> exec ^
: string
| select a from t where b = 1. and c = 2
< error
| 1064 42000 …line 1 column 32 near "and c = 2" 
> exec ^
: string
| select a from t where b = 1. or c = 2
< error
| 1064 42000 …line 1 column 31 near "or c = 2" 
> exec ^
: string
| select cast(1. as char)
< error
| 1064 42000 …line 1 column 17 near "as char)" 
> exec ^
: string
| select interval 1. foo + now()
< error
| 1064 42000 …line 1 column 22 near "foo + now()" 
> exec ^
: string
| set @a . select = 1
< error
| 1064 42000 …line 1 column 8 near ". select = 1" 
> exec ^
: string
| set a .  select = 1
< ok 0 0
> exec ^
: string
| set a . /* c */ select = 1
< error
| 1064 42000 …line 1 column 22 near "select = 1" 
> exec ^
: string
| set a . from = 1
< ok 0 0
> exec ^
: string
| set a . select = 1
< ok 0 0
> exec ^
: string
| set a. select = 1
< ok 0 0
> exec ^
: string
| set global a . select = 1
< ok 0 0
> exec ^
: string
| set transaction . select = 1
< ok 0 0
> exec ^
: string
| set x = 1, a . select = 2
< ok 0 0
> exec ^
: string
| select 1 + @@a * 2
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@ a
< error
| 1064 42000 …line 1 column 9 near "@@ a" 
> exec ^
: string
| select @@ a.b
< error
| 1064 42000 …line 1 column 9 near "@@ a.b" 
> exec ^
: string
| select @@ global.a
< error
| 1064 42000 …line 1 column 9 near "@@ global.a" 
> exec ^
: string
| select @@/* c */a
< error
| 1064 42000 …line 1 column 9 near "@@/* c */a" 
> exec ^
: string
| select @@`a`
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@`a`.`b`
< error
| 1064 42000 …line 1 column 13 near ".`b`" 
> exec ^
: string
| select @@`a`.b
< error
| 1064 42000 …line 1 column 13 near ".b" 
> exec ^
: string
| select @@`global`.a
< error
| 1064 42000 …line 1 column 18 near ".a" 
> exec ^
: string
| select @@a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@a . b
< error
| 1064 42000 …line 1 column 12 near ". b" 
> exec ^
: string
| select @@a . select
< error
| 1064 42000 …line 1 column 12 near ". select" 
> exec ^
: string
| select @@a .b
< error
| 1064 42000 …line 1 column 12 near ".b" 
> exec ^
: string
| select @@a /* c */ . b
< error
| 1064 42000 …line 1 column 20 near ". b" 
> exec ^
: string
| select @@a as b
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@a b
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@a+1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@a.
< error
| 1193 HY000 Unknown system variable 'a.'
> exec ^
: string
| select @@a. as b
< error
| 1064 42000 …line 1 column 16 near "b" 
> exec ^
: string
| select @@a. b
< error
| 1193 HY000 Unknown system variable 'a.'
> exec ^
: string
| select @@a.*
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| select @@a., 1
< error
| 1193 HY000 Unknown system variable 'a.'
> exec ^
: string
| select @@a.`b`
< error
| 1193 HY000 Unknown system variable 'a.'
> exec ^
: string
| select @@a.b
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| select @@a.b c
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| select @@a.b+1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| select @@a.b.c
< error
| 1193 HY000 Unknown system variable 'a.b.c'
> exec ^
: string
| select @@a.select
< error
| 1193 HY000 Unknown system variable 'a.select'
> exec ^
: string
| select @@a.select.b
< error
| 1193 HY000 Unknown system variable 'a.select.b'
> exec ^
: string
| select @@a=1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@global . a
< error
| 1064 42000 …line 1 column 17 near ". a" 
> exec ^
: string
| select @@global .a
< error
| 1064 42000 …line 1 column 17 near ".a" 
> exec ^
: string
| select @@global.
< error
| 1064 42000 …line 1 column 16 near "@@global." 
> exec ^
: string
| select @@global. a
< error
| 1064 42000 …line 1 column 16 near "@@global. a" 
> exec ^
: string
| select @@global.,1
< error
| 1064 42000 …line 1 column 16 near "@@global.,1" 
> exec ^
: string
| select @@global.`a`
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@global.a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@global.a . b
< error
| 1064 42000 …line 1 column 19 near ". b" 
> exec ^
: string
| select @@global.a .b
< error
| 1064 42000 …line 1 column 19 near ".b" 
> exec ^
: string
| select @@global.a, @@session.b
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@global.select
< error
| 1193 HY000 Unknown system variable 'select'
> exec ^
: string
| select @@local.a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@persist . a
< error
| 1064 42000 …line 1 column 18 near ". a" 
> exec ^
: string
| select @@persist.a
< error
| 1193 HY000 Unknown system variable 'persist.a'
> exec ^
: string
| select @@select
< error
| 1193 HY000 Unknown system variable 'select'
> exec ^
: string
| select @@select.a
< error
| 1193 HY000 Unknown system variable 'select.a'
> exec ^
: string
| select @@session .a
< error
| 1064 42000 …line 1 column 18 near ".a" 
> exec ^
: string
| select @@session.a.b
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set @@ a = 1
< error
| 1064 42000 …line 1 column 6 near "@@ a = 1" 
> exec ^
: string
| set @@a . b = 1
< error
| 1064 42000 …line 1 column 9 near ". b = 1" 
> exec ^
: string
| set @@a .b = 1
< error
| 1064 42000 …line 1 column 9 near ".b = 1" 
> exec ^
: string
| set @@a. b = 1
< error
| 1064 42000 …line 1 column 10 near "b = 1" 
> exec ^
: string
| set @@a.b = 1
< ok 0 0
> exec ^
: string
| set @@global . a = 1
< error
| 1064 42000 …line 1 column 14 near ". a = 1" 
> exec ^
: string
| set @@global .a = 1
< error
| 1064 42000 …line 1 column 14 near ".a = 1" 
> exec ^
: string
| set @@global. a = 1
< error
| 1064 42000 …line 1 column 13 near "@@global. a = 1" 
> exec ^
: string
| set @@global.a . b = 1
< error
| 1064 42000 …line 1 column 16 near ". b = 1" 
> exec ^
: string
| set @@global.a = 1
< ok 0 0
> exec ^
: string
| set @@global.a.b = 1
< ok 0 0
> exec ^
: string
| desc a . select b
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| desc a . select b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| describe a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| describe a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| explain a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show bindings from a.b
< error
| 1064 42000 …line 1 column 18 near "from a.b" 
> exec ^
: string
| show columns from a . b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| show columns from a . b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show columns from a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| show columns from a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show columns from a.b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| show columns from a.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show columns from t from a
< error
| 1146 42S02 Table 'a.t' doesn't exist
> exec ^
: string
| show columns from t from a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show columns from t from a . b
< error
| 1064 42000 …line 1 column 28 near ". b" 
> exec ^
: string
| show columns from t from a.b
< error
| 1064 42000 …line 1 column 27 near ".b" 
> exec ^
: string
| show create database a . b
< error
| 1064 42000 …line 1 column 24 near ". b" 
> exec ^
: string
| show create sequence a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| show create sequence a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show create table a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| show create table a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show create user a . b
< error
| 1064 42000 …line 1 column 20 near ". b" 
> exec ^
: string
| show create view a . select
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| show create view a . select
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show events from a
< ok 0 0
> exec ^
: string
| show events from a . b
< error
| 1064 42000 …line 1 column 20 near ". b" 
> exec ^
: string
| show events from a.b
< error
| 1064 42000 …line 1 column 19 near ".b" 
> exec ^
: string
| show full tables from a
< ok 0 0
> exec ^
: string
| show full tables from a . b
< error
| 1064 42000 …line 1 column 25 near ". b" 
> exec ^
: string
| show full tables in a.b
< error
| 1064 42000 …line 1 column 22 near ".b" 
> exec ^
: string
| show function status from a
< error
| 1064 42000 …line 1 column 25 near "from a" 
> exec ^
: string
| show index from a . b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| show index from a . b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show index from t from a
< error
| 1146 42S02 Table 'a.t' doesn't exist
> exec ^
: string
| show index from t from a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show index from t from a.b
< error
| 1064 42000 …line 1 column 25 near ".b" 
> exec ^
: string
| show keys from t in a
< error
| 1146 42S02 Table 'a.t' doesn't exist
> exec ^
: string
| show keys from t in a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show keys from t in a.b
< error
| 1064 42000 …line 1 column 22 near ".b" 
> exec ^
: string
| show open tables from a
< ok 0 0
> exec ^
: string
| show open tables from a.b
< error
| 1064 42000 …line 1 column 24 near ".b" 
> exec ^
: string
| show placement for table a.b
< error
| 1146 42S02 Table 'a.b' doesn't exist
> exec ^
: string
| show placement for table a.b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show procedure status from a
< error
| 1064 42000 …line 1 column 26 near "from a" 
> exec ^
: string
| show stats_histograms from a.b
< error
| 1064 42000 …line 1 column 26 near "from a.b" 
> exec ^
: string
| show stats_meta from a
< error
| 1064 42000 …line 1 column 20 near "from a" 
> exec ^
: string
| show stats_meta from a.b
< error
| 1064 42000 …line 1 column 20 near "from a.b" 
> exec ^
: string
| show table a . select next_row_id
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| show table a . select next_row_id
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show table a . select regions
< error
| 1146 42S02 Table 'a.select' doesn't exist
> exec ^
: string
| show table a . select regions
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show table status from a
< ok 0 0
> exec ^
: string
| show table status from a . b
< error
| 1064 42000 …line 1 column 26 near ". b" 
> exec ^
: string
| show table status from a.b
< error
| 1064 42000 …line 1 column 25 near ".b" 
> exec ^
: string
| show table status in a like 'x'
< ok 0 0
> exec ^
: string
| show tables from 'a'
< error
| 1064 42000 …line 1 column 20 near "'a'" 
> exec ^
: string
| show tables from `a`
< ok 0 0
> exec ^
: string
| show tables from a
< ok 0 0
> exec ^
: string
| show tables from a . b
< error
| 1064 42000 …line 1 column 20 near ". b" 
> exec ^
: string
| show tables from a . select
< error
| 1064 42000 …line 1 column 20 near ". select" 
> exec ^
: string
| show tables from a .b
< error
| 1064 42000 …line 1 column 20 near ".b" 
> exec ^
: string
| show tables from a as b
< error
| 1064 42000 …line 1 column 21 near "as b" 
> exec ^
: string
| show tables from a b
< error
| 1064 42000 …line 1 column 20 near "b" 
> exec ^
: string
| show tables from a like 'x'
< ok 0 0
> exec ^
: string
| show tables from a where 1
< ok 0 0
> exec ^
: string
| show tables from a. b
< error
| 1064 42000 …line 1 column 19 near ". b" 
> exec ^
: string
| show tables from a.b
< error
| 1064 42000 …line 1 column 19 near ".b" 
> exec ^
: string
| show tables from select
< error
| 1064 42000 …line 1 column 23 near "select" 
> exec ^
: string
| show tables in a . b
< error
| 1064 42000 …line 1 column 18 near ". b" 
> exec ^
: string
| show triggers from a
< ok 0 0
> exec ^
: string
| show triggers from a . b
< error
| 1064 42000 …line 1 column 22 near ". b" 
> exec ^
: string
| show triggers from a.b
< error
| 1064 42000 …line 1 column 21 near ".b" 
