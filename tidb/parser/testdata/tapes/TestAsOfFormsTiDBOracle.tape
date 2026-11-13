-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAsOfFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| select * from t as of
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| select * from t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t a as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as a as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t of
< error
| 1064 42000 …line 1 column 18 near "of" 
> exec ^
: string
| select * from t as `of`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as `of`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t `of`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t `of`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as `of` as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as `of` as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from of
< error
| 1064 42000 …line 1 column 16 near "of" 
> exec ^
: string
| select * from `of`
< error
| 1146 42S02 Table 'test.of' doesn't exist
> exec ^
: string
| select * from `of`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select of from t
< error
| 1064 42000 …line 1 column 9 near "of from t" 
> exec ^
: string
| select `of` from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `of` from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select t.of from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.of from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select of.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select of.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select of .a from t
< error
| 1064 42000 …line 1 column 9 near "of .a from t" 
> exec ^
: string
| select of. a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select of. a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from of.t
< error
| 1146 42S02 Table 'of.t' doesn't exist
> exec ^
: string
| select * from of.t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db.of
< error
| 1146 42S02 Table 'db.of' doesn't exist
> exec ^
: string
| select * from db.of
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 of
< error
| 1064 42000 …line 1 column 11 near "of" 
> exec ^
: string
| select 1 as of
< error
| 1064 42000 …line 1 column 14 near "of" 
> exec ^
: string
| select 1 as `of`
< ok 0 0
> exec ^
: string
| select 1 `of`
< ok 0 0
> exec ^
: string
| create table of (a int)
< error
| 1064 42000 …line 1 column 15 near "of (a int)" 
> exec ^
: string
| create table t (of int)
< error
| 1064 42000 …line 1 column 18 near "of int)" 
> exec ^
: string
| create table t (`of` int)
< ok 0 0
> exec ^
: string
| create table of.t (a int)
< ok 0 0
> exec ^
: string
| select * from t as of.a
< error
| 1064 42000 …line 1 column 22 near ".a" 
> exec ^
: string
| select * from t as of .a
< error
| 1064 42000 …line 1 column 23 near ".a" 
> exec ^
: string
| select * from t as of. a
< error
| 1064 42000 …line 1 column 22 near ". a" 
> exec ^
: string
| select * from t as of.a as of timestamp now()
< error
| 1064 42000 …line 1 column 22 near ".a as of timestamp now()" 
> exec ^
: string
| select * from t as of as of timestamp now()
< error
| 1064 42000 …line 1 column 27 near "of timestamp now()" 
> exec ^
: string
| select * from (select 1) as of
< error
| 1064 42000 …line 1 column 30 near "of" 
> exec ^
: string
| select * from (select 1) as of timestamp now()
< error
| 1064 42000 …line 1 column 30 near "of timestamp now()" 
> exec ^
: string
| select * from (select 1) a as of timestamp now()
< error
| 1064 42000 …line 1 column 32 near "of timestamp now()" 
> exec ^
: string
| select * from (select 1) as a as of timestamp now()
< error
| 1064 42000 …line 1 column 35 near "of timestamp now()" 
> exec ^
: string
| with of as (select 1) select * from of
< error
| 1064 42000 …line 1 column 7 near "of as (select 1) select * from of" 
> exec ^
: string
| with a as of as (select 1) select * from a
< error
| 1064 42000 …line 1 column 12 near "of as (select 1) select * from a" 
> exec ^
: string
| select cast(1 as of)
< error
| 1064 42000 …line 1 column 19 near "of)" 
> exec ^
: string
| create view of as select 1
< error
| 1064 42000 …line 1 column 14 near "of as select 1" 
> exec ^
: string
| create view v as of timestamp now() select 1
< error
| 1064 42000 …line 1 column 19 near "of timestamp now() select 1" 
> exec ^
: string
| select 1 as of timestamp now()
< error
| 1064 42000 …line 1 column 14 near "of timestamp now()" 
> exec ^
: string
| select a as of timestamp now() from t
< error
| 1064 42000 …line 1 column 14 near "of timestamp now() from t" 
> exec ^
: string
| select @of
< ok 0 0
> exec ^
: string
| select @@of
< error
| 1193 HY000 Unknown system variable 'of'
> exec ^
: string
| select * from t AS OF TIMESTAMP now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t AS OF TIMESTAMP now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t As Of TimeStamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t As Of TimeStamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as  of  timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as  of  timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as
| of
| timestamp
| now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as
| of
| timestamp
| now()
| ; SELECT FROM
< error
| 1064 42000 …line 5 column 14 near "FROM" 
> exec ^
: string
| select * from t as	of	timestamp	now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as	of	timestamp	now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as/**/of/**/timestamp/**/now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as/**/of/**/timestamp/**/now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as /* c */ of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as /* c */ of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as -- c
| of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as -- c
| of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 3 column 14 near "FROM" 
> exec ^
: string
| select * from t as # c
| of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as # c
| of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 3 column 14 near "FROM" 
> exec ^
: string
| select * from t as /*! of */ timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as /*! of */ timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t /*! as of */ timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t /*! as of */ timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t /*T! as of timestamp now() */
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t /*T! as of timestamp now() */
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of /*T! timestamp */ now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of /*T! timestamp */ now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as `of` timestamp now()
< error
| 1064 42000 …line 1 column 33 near "timestamp now()" 
> exec ^
: string
| select * from t `as` of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "of timestamp now()" 
> exec ^
: string
| select * from t as of `timestamp` now()
< error
| 1064 42000 …line 1 column 33 near "`timestamp` now()" 
> exec ^
: string
| select * from t 'as' of timestamp now()
< error
| 1064 42000 …line 1 column 20 near "'as' of timestamp now()" 
> exec ^
: string
| select * from t as 'of' timestamp now()
< error
| 1064 42000 …line 1 column 23 near "'of' timestamp now()" 
> exec ^
: string
| select * from t as "of" timestamp now()
< error
| 1064 42000 …line 1 column 23 near ""of" timestamp now()" 
> exec ^
: string
| select * from t of timestamp now()
< error
| 1064 42000 …line 1 column 18 near "of timestamp now()" 
> exec ^
: string
| select * from t as timestamp now()
< error
| 1064 42000 …line 1 column 32 near "now()" 
> exec ^
: string
| select * from t as of now()
< error
| 1064 42000 …line 1 column 25 near "now()" 
> exec ^
: string
| select * from t as of 1
< error
| 1064 42000 …line 1 column 23 near "1" 
> exec ^
: string
| select * from t as of timestamp
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| select * from t as of time now()
< error
| 1064 42000 …line 1 column 26 near "time now()" 
> exec ^
: string
| select * from t as of tso 1
< error
| 1064 42000 …line 1 column 25 near "tso 1" 
> exec ^
: string
| select * from t as of system_time now()
< error
| 1064 42000 …line 1 column 33 near "system_time now()" 
> exec ^
: string
| select * from t as of timestamp_ now()
< error
| 1064 42000 …line 1 column 32 near "timestamp_ now()" 
> exec ^
: string
| select * from t as of ts now()
< error
| 1064 42000 …line 1 column 24 near "ts now()" 
> exec ^
: string
| select * from t as of date now()
< error
| 1064 42000 …line 1 column 26 near "date now()" 
> exec ^
: string
| select * from t as of timestamp()
< error
| 1064 42000 …line 1 column 33 near ")" 
> exec ^
: string
| select * from t as of timestamp(now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp(now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp '2020-01-01'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp '2020-01-01'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp timestamp '2020-01-01 00:00:00'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp timestamp '2020-01-01 00:00:00'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp timestamp
< error
| 1054 42S22 Unknown column 'timestamp' in 'expression'
> exec ^
: string
| select * from t as of timestamp timestamp
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp date '2020-01-01'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp date '2020-01-01'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() - interval 1 second
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() - interval 1 second
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1+1
< error
| 1292 22007 Incorrect time value: '2'
> exec ^
: string
| select * from t as of timestamp 1+1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp t.a
< error
| 1054 42S22 Unknown column 't.a' in 'expression'
> exec ^
: string
| select * from t as of timestamp t.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp null
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec ^
: string
| select * from t as of timestamp null
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp (1)
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp (select 1)
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp ((select 1))
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp ((select 1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp select 1
< error
| 1064 42000 …line 1 column 38 near "select 1" 
> exec ^
: string
| select * from t as of timestamp exists (select 1)
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp exists (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp not 1
< error
| 1292 22007 Incorrect time value: '{0 0 0 0 0 0 0}'
> exec ^
: string
| select * from t as of timestamp not 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 or 2
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp 1 or 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 and 2
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp 1 and 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 xor 2
< error
| 1292 22007 Incorrect time value: '{0 0 0 0 0 0 0}'
> exec ^
: string
| select * from t as of timestamp 1 xor 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 || 2
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp 1 || 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 = 2
< error
| 1292 22007 Incorrect time value: '{0 0 0 0 0 0 0}'
> exec ^
: string
| select * from t as of timestamp 1 = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a is null
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a is null
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a in (1)
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a in (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a between 1 and 2
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a between 1 and 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a like 'x'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a like 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp case when 1 then 2 end
< error
| 1292 22007 Incorrect time value: '2'
> exec ^
: string
| select * from t as of timestamp case when 1 then 2 end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp cast(1 as char)
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp cast(1 as char)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp tidb_bounded_staleness(now() - interval 1 second, now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp tidb_bounded_staleness(now() - interval 1 second, now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp tidb_bounded_staleness()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'tidb_bounded_staleness'
> exec ^
: string
| select * from t as of timestamp tidb_bounded_staleness()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp @a
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec ^
: string
| select * from t as of timestamp @a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp @@a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select * from t as of timestamp -1
< error
| 1292 22007 Incorrect time value: '-1'
> exec ^
: string
| select * from t as of timestamp -1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a := 1
< error
| 1064 42000 …line 1 column 36 near ":= 1" 
> exec ^
: string
| select * from t as of timestamp @a := 1
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp @a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 a
< error
| 1064 42000 …line 1 column 35 near "a" 
> exec ^
: string
| select * from t as of timestamp 1 as a
< error
| 1064 42000 …line 1 column 36 near "as a" 
> exec ^
: string
| select * from t as of timestamp 1,
< error
| 1064 42000 …line 1 column 34 near "" 
> exec ^
: string
| select * from t as of timestamp ,
< error
| 1064 42000 …line 1 column 33 near "," 
> exec ^
: string
| select * from t as of timestamp now() as of timestamp now()
< error
| 1064 42000 …line 1 column 43 near "of timestamp now()" 
> exec ^
: string
| select * from t as of timestamp now() a as of timestamp now()
< error
| 1064 42000 …line 1 column 39 near "a as of timestamp now()" 
> exec ^
: string
| select * from t partition (p0) as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t partition (p0) as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t partition (p0) a as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t partition (p0) a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t partition (p0) as a as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t partition (p0) as a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() partition (p0)
< error
| 1064 42000 …line 1 column 47 near "partition (p0)" 
> exec ^
: string
| select * from t as of timestamp now() a
< error
| 1064 42000 …line 1 column 39 near "a" 
> exec ^
: string
| select * from t as of timestamp now() as a
< error
| 1064 42000 …line 1 column 40 near "as a" 
> exec ^
: string
| select * from t a as of timestamp now() b
< error
| 1064 42000 …line 1 column 41 near "b" 
> exec ^
: string
| select * from t as of timestamp now() use index (i)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() use index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() force index (i)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() force index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() ignore index (i)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() ignore index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() use index (i) ignore index (j)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() use index (i) ignore index (j)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t a as of timestamp now() use index (i)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t a as of timestamp now() use index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t use index (i) as of timestamp now()
< error
| 1064 42000 …line 1 column 35 near "of timestamp now()" 
> exec ^
: string
| select * from t a use index (i) as of timestamp now()
< error
| 1064 42000 …line 1 column 37 near "of timestamp now()" 
> exec ^
: string
| select * from t as of timestamp now() tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t as of timestamp now() tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() tablesample regions() repeatable (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() tablesample regions() repeatable (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() use index (i) tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t as of timestamp now() use index (i) tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t partition (p0) as a as of timestamp now() use index (i) tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t partition (p0) as a as of timestamp now() use index (i) tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t tablesample (50) as of timestamp now()
< error
| 1064 42000 …line 1 column 38 near "of timestamp now()" 
> exec ^
: string
| select * from t use index (i) tablesample (50) as of timestamp now()
< error
| 1064 42000 …line 1 column 52 near "of timestamp now()" 
> exec ^
: string
| select * from db.t as of timestamp now()
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| select * from db.t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db . t as of timestamp now()
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| select * from db . t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db.t a as of timestamp now()
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| select * from db.t a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from `t` as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from `t` as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t.* as of timestamp now()
< error
| 1064 42000 …line 1 column 17 near "* as of timestamp now()" 
> exec ^
: string
| select * from dual as of timestamp now()
< error
| 1064 42000 …line 1 column 24 near "of timestamp now()" 
> exec ^
: string
| select 1 from dual as of timestamp now()
< error
| 1064 42000 …line 1 column 24 near "of timestamp now()" 
> exec ^
: string
| select 1 from dual a as of timestamp now()
< error
| 1064 42000 …line 1 column 20 near "a as of timestamp now()" 
> exec ^
: string
| select * from lateral (select 1) as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "(select 1) as of timestamp now()" 
> exec ^
: string
| select * from lateral (select 1) a as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "(select 1) a as of timestamp now()" 
> exec ^
: string
| select * from json_table('[]', '$[*]' columns (a int path '$')) as of timestamp now()
< error
| 1064 42000 …line 1 column 25 near "('[]', '$[*]' columns (a int path '$')) as of timestamp now()" 
> exec ^
: string
| select * from json_table('[]', '$[*]' columns (a int path '$')) a as of timestamp now()
< error
| 1064 42000 …line 1 column 25 near "('[]', '$[*]' columns (a int path '$')) a as of timestamp now()" 
> exec ^
: string
| select * from (t) as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "of timestamp now()" 
> exec ^
: string
| select * from (t as of timestamp now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from (t as of timestamp now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (t as of timestamp now()) a
< error
| 1064 42000 …line 1 column 41 near "a" 
> exec ^
: string
| select * from (t as of timestamp now(), u)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from (t as of timestamp now(), u)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (t as of timestamp now() join u)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from (t as of timestamp now() join u)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from ((t as of timestamp now()))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from ((t as of timestamp now()))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() where a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() where a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() group by a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() group by a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() having a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() having a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() order by a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() order by a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() window w as ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() window w as ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() for update
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() for share
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| select * from t as of timestamp now() for share
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() lock in share mode
< error
| 1235 42000 function LOCK IN SHARE MODE has only noop implementation in tidb now, use tidb_enable_noop_functions to enable these functions
> exec ^
: string
| select * from t as of timestamp now() lock in share mode
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() union select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() union select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() union select * from u as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() union select * from u as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() except select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() except select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() intersect select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() intersect select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() into outfile 'x'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select * from t as of timestamp now() into outfile 'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now();
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now();
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() select 1
< error
| 1064 42000 …line 1 column 44 near "select 1" 
> exec ^
: string
| (select * from t as of timestamp now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| (select * from t as of timestamp now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select (select 1 from t as of timestamp now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select (select 1 from t as of timestamp now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from (select * from t as of timestamp now()) a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from (select * from t as of timestamp now()) a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t where a in (select a from u as of timestamp now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t where a in (select a from u as of timestamp now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| with c as (select * from t as of timestamp now()) select * from c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| with c as (select * from t as of timestamp now()) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from c as of timestamp now(), d as of timestamp now()
< error
| 1146 42S02 Table 'test.c' doesn't exist
> exec ^
: string
| select * from c as of timestamp now(), d as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now(), u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now(), u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u, t as of timestamp now()
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u, t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() join u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() join u on t.a = u.a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() join u on t.a = u.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() join u as of timestamp now() on t.a = u.a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() join u as of timestamp now() on t.a = u.a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() inner join u using (a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() inner join u using (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() cross join u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() cross join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() left join u on 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() left join u on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() right join u on 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() right join u on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() left outer join u on 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() left outer join u on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() natural join u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() natural join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() natural left join u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() natural left join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp now() straight_join u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() straight_join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 1 natural join u
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| select * from t as of timestamp 1 natural join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a straight_join u
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a straight_join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u join t as of timestamp now()
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u join t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u join t as of timestamp now() on 1
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u join t as of timestamp now() on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u left join t as of timestamp now() on 1
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u left join t as of timestamp now() on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u join t as of timestamp a on 1
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u join t as of timestamp a on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u join t as of timestamp a using (a)
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u join t as of timestamp a using (a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from u natural join t as of timestamp a
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from u natural join t as of timestamp a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t as of timestamp now() set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t as of timestamp now() set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t a as of timestamp now() set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t a as of timestamp now() set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t as a as of timestamp now() set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t as a as of timestamp now() set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t as of timestamp now(), u set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t as of timestamp now(), u set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t as of timestamp now() join u set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t as of timestamp now() join u set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t as of timestamp now() use index (i) set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t as of timestamp now() use index (i) set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t, u as of timestamp now() set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t, u as of timestamp now() set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update u join t as of timestamp now() on 1 set a = 1
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| update u join t as of timestamp now() on 1 set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a = (select 1 from u as of timestamp now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a = (select 1 from u as of timestamp now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t as of timestamp now()
< error
| 1064 42000 …line 1 column 19 near "of timestamp now()" 
> exec ^
: string
| delete from t a as of timestamp now()
< error
| 1064 42000 …line 1 column 21 near "of timestamp now()" 
> exec ^
: string
| delete from t as a as of timestamp now()
< error
| 1064 42000 …line 1 column 24 near "of timestamp now()" 
> exec ^
: string
| delete from t as of timestamp now() where a
< error
| 1064 42000 …line 1 column 19 near "of timestamp now() where a" 
> exec ^
: string
| delete t from t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete t from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete a from t a as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete a from t a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete a from t as a as of timestamp now() where a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete a from t as a as of timestamp now() where a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete t from t as of timestamp now(), u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete t from t as of timestamp now(), u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t using t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t using t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t using t a as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t using t a as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t where a in (select a from u as of timestamp now())
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t where a in (select a from u as of timestamp now())
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t as of timestamp now() values (1)
< error
| 1064 42000 …line 1 column 19 near "of timestamp now() values (1)" 
> exec ^
: string
| insert into t select * from u as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t select * from u as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t as of timestamp now() select 1
< error
| 1064 42000 …line 1 column 19 near "of timestamp now() select 1" 
> exec ^
: string
| replace into t select * from u as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t select * from u as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t as of timestamp now() select 1
< error
| 1064 42000 …line 1 column 20 near "of timestamp now() select 1" 
> exec ^
: string
| create table t2 as select * from t as of timestamp now()
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t2 as select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t2 select * from t as of timestamp now()
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create table t2 select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create view v as select * from t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create view v as select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create view v as select * from t as of timestamp now() union select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create view v as select * from t as of timestamp now() union select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain select * from t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain analyze select * from t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain analyze select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain update t as of timestamp now() set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain update t as of timestamp now() set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace select * from t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace select * from t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| lock tables t as of timestamp now() read
< error
| 1064 42000 …line 1 column 19 near "of timestamp now() read" 
> exec ^
: string
| lock tables t a as of timestamp now() read
< error
| 1064 42000 …line 1 column 15 near "a as of timestamp now() read" 
> exec ^
: string
| lock tables t read as of timestamp now()
< error
| 1064 42000 …line 1 column 24 near "of timestamp now()" 
> exec ^
: string
| lock tables t as of read
< error
| 1064 42000 …line 1 column 19 near "of read" 
> exec ^
: string
| lock table t as of timestamp now() write
< error
| 1064 42000 …line 1 column 18 near "of timestamp now() write" 
> exec ^
: string
| table t as of timestamp now()
< error
| 1064 42000 …line 1 column 13 near "of timestamp now()" 
> exec ^
: string
| truncate t as of timestamp now()
< error
| 1064 42000 …line 1 column 16 near "of timestamp now()" 
> exec ^
: string
| drop table t as of timestamp now()
< error
| 1064 42000 …line 1 column 18 near "of timestamp now()" 
> exec ^
: string
| show columns from t as of timestamp now()
< error
| 1064 42000 …line 1 column 25 near "of timestamp now()" 
> exec ^
: string
| describe t as of timestamp now()
< error
| 1064 42000 …line 1 column 16 near "of timestamp now()" 
> exec ^
: string
| analyze table t as of timestamp now()
< error
| 1064 42000 …line 1 column 21 near "of timestamp now()" 
> exec ^
: string
| do 1 as of timestamp now()
< error
| 1064 42000 …line 1 column 10 near "of timestamp now()" 
> exec ^
: string
| select 1 as of timestamp 1
< error
| 1064 42000 …line 1 column 14 near "of timestamp 1" 
> exec ^
: string
| create index i on t (a) as of timestamp now()
< error
| 1064 42000 …line 1 column 29 near "of timestamp now()" 
> exec ^
: string
| alter table t as of timestamp now() add column b int
< error
| 1064 42000 …line 1 column 19 near "of timestamp now() add column b int" 
> exec ^
: string
| prepare s from 'select 1' as of timestamp now()
< error
| 1064 42000 …line 1 column 31 near "of timestamp now()" 
> exec ^
: string
| create binding for select * from t as of timestamp now() using select * from t as of timestamp now()
< ok 0 0
> exec ^
: string
| start transaction read only as of
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| start transaction as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "of timestamp now()" 
> exec ^
: string
| begin as of timestamp now()
< error
| 1064 42000 …line 1 column 11 near "of timestamp now()" 
> exec ^
: string
| plan replayer dump with stats as of timestamp now() explain select 1
< ok 0 0
> exec ^
: string
| delete of from t as of timestamp now()
< error
| 1064 42000 …line 1 column 9 near "of from t as of timestamp now()" 
> exec ^
: string
| delete t from t as t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete t from t as t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete t from t t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete t from t t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t using t as t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t using t as t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t using t t as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| delete from t using t t as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| delete from t using t as of
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| select * from t for system_time as of timestamp now()
< error
| 1064 42000 …line 1 column 31 near "system_time as of timestamp now()" 
> exec ^
: string
| select * from t for system_time as of now()
< error
| 1064 42000 …line 1 column 31 near "system_time as of now()" 
> exec ^
: string
| select * from t as of timestamp now() for system_time as of timestamp now()
< error
| 1064 42000 …line 1 column 53 near "system_time as of timestamp now()" 
> exec ^
: string
| select * from timestamp as of timestamp now()
< error
| 1146 42S02 Table 'test.timestamp' doesn't exist
> exec ^
: string
| select * from timestamp as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from db.timestamp as of timestamp now()
< error
| 1146 42S02 Table 'db.timestamp' doesn't exist
> exec ^
: string
| select * from db.timestamp as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t timestamp as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t timestamp as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as timestamp as of timestamp now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as timestamp as of timestamp now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select timestamp as of from t
< error
| 1064 42000 …line 1 column 22 near "of from t" 
> exec ^
: string
| select * from t as a as b as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "as b as of timestamp now()" 
> exec ^
: string
| select * from t a b as of timestamp now()
< error
| 1064 42000 …line 1 column 19 near "b as of timestamp now()" 
> exec ^
: string
| select * from t as of timestamp now() as of
< error
| 1064 42000 …line 1 column 43 near "of" 
> exec ^
: string
| select * from t as of timestamp now() use index (i) as of timestamp now()
< error
| 1064 42000 …line 1 column 57 near "of timestamp now()" 
> exec ^
: string
| select * from t as of timestamp now() for update of t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as of timestamp now() for update of t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a left join u on 1
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a left join u on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp left('a', 1)
< error
| 1292 22007 Incorrect datetime value: 'a'
> exec ^
: string
| select * from t as of timestamp left('a', 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp left('a', 1) left join u on 1
< error
| 1292 22007 Incorrect datetime value: 'a'
> exec ^
: string
| select * from t as of timestamp left('a', 1) left join u on 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a use index (i)
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a use index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a ignore index (i)
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a ignore index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a force index (i)
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a force index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t as of timestamp a tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a limit 1
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a from u
< error
| 1064 42000 …line 1 column 38 near "from u" 
> exec ^
: string
| select * from t as of timestamp a b
< error
| 1064 42000 …line 1 column 35 near "b" 
> exec ^
: string
| select * from t as of timestamp (select 1) tablesample (50)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t as of timestamp (select 1) tablesample (50)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a in (1) join u
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a in (1) join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp not a join u
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp not a join u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp timestamp '2020-01-01 00:00:00' tablesample (1)
< error
| 8128 HY000 Invalid TABLESAMPLE: Only supports REGIONS sampling method
> exec ^
: string
| select * from t as of timestamp timestamp '2020-01-01 00:00:00' tablesample (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp 'a' collate utf8mb4_bin
< error
| 1292 22007 Incorrect datetime value: 'a'
> exec ^
: string
| select * from t as of timestamp 'a' collate utf8mb4_bin
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a collate utf8mb4_bin use index (i)
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a collate utf8mb4_bin use index (i)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp interval 1 day + now()
< error
| 1105 HY000 cannot set read timestamp to a future time, readTS: …, currentTS: …
> exec ^
: string
| select * from t as of timestamp interval 1 day + now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a member of ('[]')
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a member of ('[]')
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp binary 'a'
< error
| 1292 22007 Incorrect datetime value: 'a'
> exec ^
: string
| select * from t as of timestamp binary 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp !1
< error
| 1292 22007 Incorrect time value: '{0 0 0 0 0 0 0}'
> exec ^
: string
| select * from t as of timestamp !1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp ~1
< error
| 1105 HY000 cannot convert datum from unsigned bigint to type timestamp
> exec ^
: string
| select * from t as of timestamp ~1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp a -> '$.b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| select * from t as of timestamp a -> '$.b'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as of timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * from u)
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| select * from t as of timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * from u)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (select * from u as of timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * from u))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (select * from u as of timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * from u))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| plan replayer dump with stats as timestamp now() explain select 1
< error
| 1064 42000 …line 1 column 32 near "as timestamp now() explain select 1" 
> exec ^
: string
| plan replayer dump with stats of timestamp now() explain select 1
< error
| 1064 42000 …line 1 column 32 near "of timestamp now() explain select 1" 
> exec ^
: string
| plan replayer dump with stats timestamp now() explain select 1
< error
| 1064 42000 …line 1 column 39 near "timestamp now() explain select 1" 
> exec ^
: string
| plan replayer dump with stats as of now() explain select 1
< error
| 1064 42000 …line 1 column 39 near "now() explain select 1" 
> exec ^
: string
| plan replayer dump with stats as of timestamp explain select 1
< error
| 1064 42000 …line 1 column 53 near "explain select 1" 
> exec ^
: string
| plan replayer dump with stats as of explain select 1
< error
| 1064 42000 …line 1 column 43 near "explain select 1" 
