-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestAsOfFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| select * from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t as a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as `of` as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| create table t (`of` int)
< ok 0 0
> exec ^
: string
| select * from t as of.a
< error
| 1064 42000 …near '.a' at line 1
> exec ^
: string
| select * from t as of .a
< error
| 1064 42000 …near '.a' at line 1
> exec ^
: string
| select * from t as of. a
< error
| 1064 42000 …near '. a' at line 1
> exec ^
: string
| select * from t as of.a as of timestamp now()
< error
| 1064 42000 …near '.a as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from (select 1) as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from (select 1) a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from (select 1) as a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| with a as of as (select 1) select * from a
< error
| 1064 42000 …near 'of as (select 1) select * from a' at line 1
> exec ^
: string
| create view v as of timestamp now() select 1
< error
| 1064 42000 …near 'of timestamp now() select 1' at line 1
> exec ^
: string
| select 1 as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select a as of timestamp now() from t
< error
| 1064 42000 …near 'timestamp now() from t' at line 1
> exec ^
: string
| select @of
< ok 0 0
> exec ^
: string
| select * from t AS OF TIMESTAMP now()
< error
| 1064 42000 …near 'TIMESTAMP now()' at line 1
> exec ^
: string
| select * from t As Of TimeStamp now()
< error
| 1064 42000 …near 'TimeStamp now()' at line 1
> exec ^
: string
| select * from t as  of  timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t as
| of
| timestamp
| now()
< error
| 1064 42000 …near 'timestamp
| now()' at line 3
> exec ^
: string
| select * from t as	of	timestamp	now()
< error
| 1064 42000 …near 'timestamp	now()' at line 1
> exec ^
: string
| select * from t as/**/of/**/timestamp/**/now()
< error
| 1064 42000 …near 'timestamp/**/now()' at line 1
> exec ^
: string
| select * from t as /* c */ of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t as -- c
| of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 2
> exec ^
: string
| select * from t as # c
| of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 2
> exec ^
: string
| select * from t as /*! of */ timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t /*! as of */ timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as of /*T! timestamp */ now()
< error
| 1064 42000 …near 'now()' at line 1
> exec ^
: string
| select * from t as `of` timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t `as` of timestamp now()
< error
| 1064 42000 …near 'of timestamp now()' at line 1
> exec ^
: string
| select * from t as of `timestamp` now()
< error
| 1064 42000 …near '`timestamp` now()' at line 1
> exec ^
: string
| select * from t 'as' of timestamp now()
< error
| 1064 42000 …near ''as' of timestamp now()' at line 1
> exec ^
: string
| select * from t as 'of' timestamp now()
< error
| 1064 42000 …near ''of' timestamp now()' at line 1
> exec ^
: string
| select * from t as "of" timestamp now()
< error
| 1064 42000 …near '"of" timestamp now()' at line 1
> exec ^
: string
| select * from t of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t as timestamp now()
< error
| 1064 42000 …near 'now()' at line 1
> exec ^
: string
| select * from t as of now()
< error
| 1064 42000 …near 'now()' at line 1
> exec ^
: string
| select * from t as of 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| select * from t as of timestamp
< error
| 1064 42000 …near 'timestamp' at line 1
> exec ^
: string
| select * from t as of time now()
< error
| 1064 42000 …near 'time now()' at line 1
> exec ^
: string
| select * from t as of tso 1
< error
| 1064 42000 …near 'tso 1' at line 1
> exec ^
: string
| select * from t as of system_time now()
< error
| 1064 42000 …near 'system_time now()' at line 1
> exec ^
: string
| select * from t as of timestamp_ now()
< error
| 1064 42000 …near 'timestamp_ now()' at line 1
> exec ^
: string
| select * from t as of ts now()
< error
| 1064 42000 …near 'ts now()' at line 1
> exec ^
: string
| select * from t as of date now()
< error
| 1064 42000 …near 'date now()' at line 1
> exec ^
: string
| select * from t as of timestamp()
< error
| 1064 42000 …near 'timestamp()' at line 1
> exec ^
: string
| select * from t as of timestamp(now())
< error
| 1064 42000 …near 'timestamp(now())' at line 1
> exec ^
: string
| select * from t as of timestamp '2020-01-01'
< error
| 1064 42000 …near 'timestamp '2020-01-01'' at line 1
> exec ^
: string
| select * from t as of timestamp timestamp '2020-01-01 00:00:00'
< error
| 1064 42000 …near 'timestamp timestamp '2020-01-01 00:00:00'' at line 1
> exec ^
: string
| select * from t as of timestamp timestamp
< error
| 1064 42000 …near 'timestamp timestamp' at line 1
> exec ^
: string
| select * from t as of timestamp 'a' 'b'
< error
| 1064 42000 …near 'timestamp 'a' 'b'' at line 1
> exec ^
: string
| select * from t as of timestamp date '2020-01-01'
< error
| 1064 42000 …near 'timestamp date '2020-01-01'' at line 1
> exec ^
: string
| select * from t as of timestamp now() - interval 1 second
< error
| 1064 42000 …near 'timestamp now() - interval 1 second' at line 1
> exec ^
: string
| select * from t as of timestamp 1
< error
| 1064 42000 …near 'timestamp 1' at line 1
> exec ^
: string
| select * from t as of timestamp 1+1
< error
| 1064 42000 …near 'timestamp 1+1' at line 1
> exec ^
: string
| select * from t as of timestamp a
< error
| 1064 42000 …near 'timestamp a' at line 1
> exec ^
: string
| select * from t as of timestamp t.a
< error
| 1064 42000 …near 'timestamp t.a' at line 1
> exec ^
: string
| select * from t as of timestamp null
< error
| 1064 42000 …near 'timestamp null' at line 1
> exec ^
: string
| select * from t as of timestamp (1)
< error
| 1064 42000 …near 'timestamp (1)' at line 1
> exec ^
: string
| select * from t as of timestamp (1, 2)
< error
| 1064 42000 …near 'timestamp (1, 2)' at line 1
> exec ^
: string
| select * from t as of timestamp (select 1)
< error
| 1064 42000 …near 'timestamp (select 1)' at line 1
> exec ^
: string
| select * from t as of timestamp ((select 1))
< error
| 1064 42000 …near 'timestamp ((select 1))' at line 1
> exec ^
: string
| select * from t as of timestamp select 1
< error
| 1064 42000 …near 'timestamp select 1' at line 1
> exec ^
: string
| select * from t as of timestamp exists (select 1)
< error
| 1064 42000 …near 'timestamp exists (select 1)' at line 1
> exec ^
: string
| select * from t as of timestamp not 1
< error
| 1064 42000 …near 'timestamp not 1' at line 1
> exec ^
: string
| select * from t as of timestamp 1 or 2
< error
| 1064 42000 …near 'timestamp 1 or 2' at line 1
> exec ^
: string
| select * from t as of timestamp 1 and 2
< error
| 1064 42000 …near 'timestamp 1 and 2' at line 1
> exec ^
: string
| select * from t as of timestamp 1 xor 2
< error
| 1064 42000 …near 'timestamp 1 xor 2' at line 1
> exec ^
: string
| select * from t as of timestamp 1 || 2
< error
| 1064 42000 …near 'timestamp 1 || 2' at line 1
> exec ^
: string
| select * from t as of timestamp 1 = 2
< error
| 1064 42000 …near 'timestamp 1 = 2' at line 1
> exec ^
: string
| select * from t as of timestamp a is null
< error
| 1064 42000 …near 'timestamp a is null' at line 1
> exec ^
: string
| select * from t as of timestamp a in (1)
< error
| 1064 42000 …near 'timestamp a in (1)' at line 1
> exec ^
: string
| select * from t as of timestamp a between 1 and 2
< error
| 1064 42000 …near 'timestamp a between 1 and 2' at line 1
> exec ^
: string
| select * from t as of timestamp a like 'x'
< error
| 1064 42000 …near 'timestamp a like 'x'' at line 1
> exec ^
: string
| select * from t as of timestamp case when 1 then 2 end
< error
| 1064 42000 …near 'timestamp case when 1 then 2 end' at line 1
> exec ^
: string
| select * from t as of timestamp cast(1 as char)
< error
| 1064 42000 …near 'timestamp cast(1 as char)' at line 1
> exec ^
: string
| select * from t as of timestamp tidb_bounded_staleness(now() - interval 1 second, now())
< error
| 1064 42000 …near 'timestamp tidb_bounded_staleness(now() - interval 1 second, now())' at line 1
> exec ^
: string
| select * from t as of timestamp tidb_bounded_staleness()
< error
| 1064 42000 …near 'timestamp tidb_bounded_staleness()' at line 1
> exec ^
: string
| select * from t as of timestamp @a
< error
| 1064 42000 …near 'timestamp @a' at line 1
> exec ^
: string
| select * from t as of timestamp @@a
< error
| 1064 42000 …near 'timestamp @@a' at line 1
> exec ^
: string
| select * from t as of timestamp ?
< error
| 1064 42000 …near 'timestamp ?' at line 1
> exec ^
: string
| select * from t as of timestamp -1
< error
| 1064 42000 …near 'timestamp -1' at line 1
> exec ^
: string
| select * from t as of timestamp a := 1
< error
| 1064 42000 …near 'timestamp a := 1' at line 1
> exec ^
: string
| select * from t as of timestamp @a := 1
< error
| 1064 42000 …near 'timestamp @a := 1' at line 1
> exec ^
: string
| select * from t as of timestamp 1 a
< error
| 1064 42000 …near 'timestamp 1 a' at line 1
> exec ^
: string
| select * from t as of timestamp 1 as a
< error
| 1064 42000 …near 'timestamp 1 as a' at line 1
> exec ^
: string
| select * from t as of timestamp 1,
< error
| 1064 42000 …near 'timestamp 1,' at line 1
> exec ^
: string
| select * from t as of timestamp ,
< error
| 1064 42000 …near 'timestamp ,' at line 1
> exec ^
: string
| select * from t as of timestamp now() as of timestamp now()
< error
| 1064 42000 …near 'timestamp now() as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() a as of timestamp now()
< error
| 1064 42000 …near 'timestamp now() a as of timestamp now()' at line 1
> exec ^
: string
| select * from t partition (p0) as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t partition (p0) a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t partition (p0) as a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() partition (p0)
< error
| 1064 42000 …near 'timestamp now() partition (p0)' at line 1
> exec ^
: string
| select * from t as of timestamp now() a
< error
| 1064 42000 …near 'timestamp now() a' at line 1
> exec ^
: string
| select * from t as of timestamp now() as a
< error
| 1064 42000 …near 'timestamp now() as a' at line 1
> exec ^
: string
| select * from t a as of timestamp now() b
< error
| 1064 42000 …near 'as of timestamp now() b' at line 1
> exec ^
: string
| select * from t as of timestamp now() use index (i)
< error
| 1064 42000 …near 'timestamp now() use index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp now() force index (i)
< error
| 1064 42000 …near 'timestamp now() force index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp now() ignore index (i)
< error
| 1064 42000 …near 'timestamp now() ignore index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp now() use index (i) ignore index (j)
< error
| 1064 42000 …near 'timestamp now() use index (i) ignore index (j)' at line 1
> exec ^
: string
| select * from t a as of timestamp now() use index (i)
< error
| 1064 42000 …near 'as of timestamp now() use index (i)' at line 1
> exec ^
: string
| select * from t use index (i) as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t a use index (i) as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() tablesample (50)
< error
| 1064 42000 …near 'timestamp now() tablesample (50)' at line 1
> exec ^
: string
| select * from t as of timestamp now() tablesample regions() repeatable (1)
< error
| 1064 42000 …near 'timestamp now() tablesample regions() repeatable (1)' at line 1
> exec ^
: string
| select * from t as of timestamp now() use index (i) tablesample (50)
< error
| 1064 42000 …near 'timestamp now() use index (i) tablesample (50)' at line 1
> exec ^
: string
| select * from t partition (p0) as a as of timestamp now() use index (i) tablesample (50)
< error
| 1064 42000 …near 'as of timestamp now() use index (i) tablesample (50)' at line 1
> exec ^
: string
| select * from t tablesample (50) as of timestamp now()
< error
| 1064 42000 …near '(50) as of timestamp now()' at line 1
> exec ^
: string
| select * from t use index (i) tablesample (50) as of timestamp now()
< error
| 1064 42000 …near 'tablesample (50) as of timestamp now()' at line 1
> exec ^
: string
| select * from db.t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from db . t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from db.t a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from `t` as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t.* as of timestamp now()
< error
| 1064 42000 …near '* as of timestamp now()' at line 1
> exec ^
: string
| select * from dual as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select 1 from dual as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select 1 from dual a as of timestamp now()
< error
| 1064 42000 …near 'a as of timestamp now()' at line 1
> exec ^
: string
| select * from lateral (select 1) as of timestamp now()
< error
| 1064 42000 …near '(select 1) as of timestamp now()' at line 1
> exec ^
: string
| select * from lateral (select 1) a as of timestamp now()
< error
| 1064 42000 …near '(select 1) a as of timestamp now()' at line 1
> exec ^
: string
| select * from json_table('[]', '$[*]' columns (a int path '$')) as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from json_table('[]', '$[*]' columns (a int path '$')) a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from (t) as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from (t as of timestamp now())
< error
| 1064 42000 …near 'timestamp now())' at line 1
> exec ^
: string
| select * from (t as of timestamp now()) a
< error
| 1064 42000 …near 'timestamp now()) a' at line 1
> exec ^
: string
| select * from (t as of timestamp now(), u)
< error
| 1064 42000 …near 'timestamp now(), u)' at line 1
> exec ^
: string
| select * from (t as of timestamp now() join u)
< error
| 1064 42000 …near 'timestamp now() join u)' at line 1
> exec ^
: string
| select * from ((t as of timestamp now()))
< error
| 1064 42000 …near 'timestamp now()))' at line 1
> exec ^
: string
| select * from { oj t as of timestamp now() }
< error
| 1064 42000 …near 'timestamp now() }' at line 1
> exec ^
: string
| select * from { oj t as of timestamp now() left join u on 1 }
< error
| 1064 42000 …near 'timestamp now() left join u on 1 }' at line 1
> exec ^
: string
| select * from t as of timestamp now() where a
< error
| 1064 42000 …near 'timestamp now() where a' at line 1
> exec ^
: string
| select * from t as of timestamp now() group by a
< error
| 1064 42000 …near 'timestamp now() group by a' at line 1
> exec ^
: string
| select * from t as of timestamp now() having a
< error
| 1064 42000 …near 'timestamp now() having a' at line 1
> exec ^
: string
| select * from t as of timestamp now() order by a
< error
| 1064 42000 …near 'timestamp now() order by a' at line 1
> exec ^
: string
| select * from t as of timestamp now() limit 1
< error
| 1064 42000 …near 'timestamp now() limit 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() window w as ()
< error
| 1064 42000 …near 'timestamp now() window w as ()' at line 1
> exec ^
: string
| select * from t as of timestamp now() for update
< error
| 1064 42000 …near 'timestamp now() for update' at line 1
> exec ^
: string
| select * from t as of timestamp now() for share
< error
| 1064 42000 …near 'timestamp now() for share' at line 1
> exec ^
: string
| select * from t as of timestamp now() lock in share mode
< error
| 1064 42000 …near 'timestamp now() lock in share mode' at line 1
> exec ^
: string
| select * from t as of timestamp now() union select 1
< error
| 1064 42000 …near 'timestamp now() union select 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() union select * from u as of timestamp now()
< error
| 1064 42000 …near 'timestamp now() union select * from u as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() except select 1
< error
| 1064 42000 …near 'timestamp now() except select 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() intersect select 1
< error
| 1064 42000 …near 'timestamp now() intersect select 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() into outfile 'x'
< error
| 1064 42000 …near 'timestamp now() into outfile 'x'' at line 1
> exec ^
: string
| select * from t as of timestamp now();
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() select 1
< error
| 1064 42000 …near 'timestamp now() select 1' at line 1
> exec ^
: string
| (select * from t as of timestamp now())
< error
| 1064 42000 …near 'timestamp now())' at line 1
> exec ^
: string
| select (select 1 from t as of timestamp now())
< error
| 1064 42000 …near 'timestamp now())' at line 1
> exec ^
: string
| select * from (select * from t as of timestamp now()) a
< error
| 1064 42000 …near 'timestamp now()) a' at line 1
> exec ^
: string
| select * from t where a in (select a from u as of timestamp now())
< error
| 1064 42000 …near 'timestamp now())' at line 1
> exec ^
: string
| with c as (select * from t as of timestamp now()) select * from c
< error
| 1064 42000 …near 'timestamp now()) select * from c' at line 1
> exec ^
: string
| select * from c as of timestamp now(), d as of timestamp now()
< error
| 1064 42000 …near 'timestamp now(), d as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now(), u
< error
| 1064 42000 …near 'timestamp now(), u' at line 1
> exec ^
: string
| select * from u, t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() join u
< error
| 1064 42000 …near 'timestamp now() join u' at line 1
> exec ^
: string
| select * from t as of timestamp now() join u on t.a = u.a
< error
| 1064 42000 …near 'timestamp now() join u on t.a = u.a' at line 1
> exec ^
: string
| select * from t as of timestamp now() join u as of timestamp now() on t.a = u.a
< error
| 1064 42000 …near 'timestamp now() join u as of timestamp now() on t.a = u.a' at line 1
> exec ^
: string
| select * from t as of timestamp now() inner join u using (a)
< error
| 1064 42000 …near 'timestamp now() inner join u using (a)' at line 1
> exec ^
: string
| select * from t as of timestamp now() cross join u
< error
| 1064 42000 …near 'timestamp now() cross join u' at line 1
> exec ^
: string
| select * from t as of timestamp now() left join u on 1
< error
| 1064 42000 …near 'timestamp now() left join u on 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() right join u on 1
< error
| 1064 42000 …near 'timestamp now() right join u on 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() left outer join u on 1
< error
| 1064 42000 …near 'timestamp now() left outer join u on 1' at line 1
> exec ^
: string
| select * from t as of timestamp now() natural join u
< error
| 1064 42000 …near 'timestamp now() natural join u' at line 1
> exec ^
: string
| select * from t as of timestamp now() natural left join u
< error
| 1064 42000 …near 'timestamp now() natural left join u' at line 1
> exec ^
: string
| select * from t as of timestamp now() straight_join u
< error
| 1064 42000 …near 'timestamp now() straight_join u' at line 1
> exec ^
: string
| select * from t as of timestamp 1 natural join u
< error
| 1064 42000 …near 'timestamp 1 natural join u' at line 1
> exec ^
: string
| select * from t as of timestamp a straight_join u
< error
| 1064 42000 …near 'timestamp a straight_join u' at line 1
> exec ^
: string
| select * from u join t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from u join t as of timestamp now() on 1
< error
| 1064 42000 …near 'timestamp now() on 1' at line 1
> exec ^
: string
| select * from u left join t as of timestamp now() on 1
< error
| 1064 42000 …near 'timestamp now() on 1' at line 1
> exec ^
: string
| select * from u join t as of timestamp a on 1
< error
| 1064 42000 …near 'timestamp a on 1' at line 1
> exec ^
: string
| select * from u join t as of timestamp a using (a)
< error
| 1064 42000 …near 'timestamp a using (a)' at line 1
> exec ^
: string
| select * from u natural join t as of timestamp a
< error
| 1064 42000 …near 'timestamp a' at line 1
> exec ^
: string
| update t as of timestamp now() set a = 1
< error
| 1064 42000 …near 'timestamp now() set a = 1' at line 1
> exec ^
: string
| update t a as of timestamp now() set a = 1
< error
| 1064 42000 …near 'as of timestamp now() set a = 1' at line 1
> exec ^
: string
| update t as a as of timestamp now() set a = 1
< error
| 1064 42000 …near 'as of timestamp now() set a = 1' at line 1
> exec ^
: string
| update t as of timestamp now(), u set a = 1
< error
| 1064 42000 …near 'timestamp now(), u set a = 1' at line 1
> exec ^
: string
| update t as of timestamp now() join u set a = 1
< error
| 1064 42000 …near 'timestamp now() join u set a = 1' at line 1
> exec ^
: string
| update t as of timestamp now() use index (i) set a = 1
< error
| 1064 42000 …near 'timestamp now() use index (i) set a = 1' at line 1
> exec ^
: string
| update t, u as of timestamp now() set a = 1
< error
| 1064 42000 …near 'timestamp now() set a = 1' at line 1
> exec ^
: string
| update u join t as of timestamp now() on 1 set a = 1
< error
| 1064 42000 …near 'timestamp now() on 1 set a = 1' at line 1
> exec ^
: string
| update t set a = (select 1 from u as of timestamp now())
< error
| 1064 42000 …near 'timestamp now())' at line 1
> exec ^
: string
| delete from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| delete from t a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| delete from t as a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| delete from t as of timestamp now() where a
< error
| 1064 42000 …near 'timestamp now() where a' at line 1
> exec ^
: string
| delete a from t a as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| delete a from t as a as of timestamp now() where a = 1
< error
| 1064 42000 …near 'as of timestamp now() where a = 1' at line 1
> exec ^
: string
| delete from t where a in (select a from u as of timestamp now())
< error
| 1064 42000 …near 'timestamp now())' at line 1
> exec ^
: string
| insert into t as of timestamp now() values (1)
< error
| 1064 42000 …near 'as of timestamp now() values (1)' at line 1
> exec ^
: string
| insert into t select * from u as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| insert into t as of timestamp now() select 1
< error
| 1064 42000 …near 'as of timestamp now() select 1' at line 1
> exec ^
: string
| replace into t select * from u as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| replace into t as of timestamp now() select 1
< error
| 1064 42000 …near 'as of timestamp now() select 1' at line 1
> exec ^
: string
| create table t2 as select * from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| create table t2 select * from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| create view v as select * from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| create view v as select * from t as of timestamp now() union select 1
< error
| 1064 42000 …near 'timestamp now() union select 1' at line 1
> exec ^
: string
| explain select * from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| explain analyze select * from t as of timestamp now()
< error
| 1064 42000 …near 'analyze select * from t as of timestamp now()' at line 1
> exec ^
: string
| desc select * from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| explain update t as of timestamp now() set a = 1
< error
| 1064 42000 …near 'timestamp now() set a = 1' at line 1
> exec ^
: string
| trace select * from t as of timestamp now()
< error
| 1064 42000 …near 'trace select * from t as of timestamp now()' at line 1
> exec ^
: string
| lock tables t as of timestamp now() read
< error
| 1064 42000 …near 'timestamp now() read' at line 1
> exec ^
: string
| lock tables t a as of timestamp now() read
< error
| 1064 42000 …near 'as of timestamp now() read' at line 1
> exec ^
: string
| lock tables t read as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| lock table t as of timestamp now() write
< error
| 1064 42000 …near 'timestamp now() write' at line 1
> exec ^
: string
| table t as of timestamp now()
< error
| 1064 42000 …near 'table t as of timestamp now()' at line 1
> exec ^
: string
| truncate t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| drop table t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| show columns from t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| describe t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| analyze table t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| do 1 as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select 1 as of timestamp 1
< error
| 1064 42000 …near 'timestamp 1' at line 1
> exec ^
: string
| create index i on t (a) as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| alter table t as of timestamp now() add column b int
< error
| 1064 42000 …near 'as of timestamp now() add column b int' at line 1
> exec ^
: string
| prepare s from 'select 1' as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| create binding for select * from t as of timestamp now() using select * from t as of timestamp now()
< error
| 1064 42000 …near 'binding for select * from t as of timestamp now() using select * from t as of...' at line 1
> exec ^
: string
| start transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| start transaction read only as of
< error
| 1064 42000 …near 'as of' at line 1
> exec ^
: string
| start transaction as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set session transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| begin as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| plan replayer dump with stats as of timestamp now() explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats as of timestamp now() explain select 1' at line 1
> exec ^
: string
| delete of from t as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| delete t from t as t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| delete t from t t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| delete from t using t as t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| delete from t using t t as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t for system_time as of now()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t for system_time as of now()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t as of timestamp now() for system_time as of timestamp now()
< error
| 1064 42000 …near 'timestamp now() for system_time as of timestamp now()' at line 1
> exec ^
: string
| select * from timestamp as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from db.timestamp as of timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| select * from t timestamp as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t as timestamp as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| select * from t as a as b as of timestamp now()
< error
| 1064 42000 …near 'as b as of timestamp now()' at line 1
> exec ^
: string
| select * from t a b as of timestamp now()
< error
| 1064 42000 …near 'b as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() as of
< error
| 1064 42000 …near 'timestamp now() as of' at line 1
> exec ^
: string
| select * from t as of timestamp now() use index (i) as of timestamp now()
< error
| 1064 42000 …near 'timestamp now() use index (i) as of timestamp now()' at line 1
> exec ^
: string
| select * from t as of timestamp now() for update of t
< error
| 1064 42000 …near 'timestamp now() for update of t' at line 1
> exec ^
: string
| select * from t as of timestamp a left join u on 1
< error
| 1064 42000 …near 'timestamp a left join u on 1' at line 1
> exec ^
: string
| select * from t as of timestamp left('a', 1)
< error
| 1064 42000 …near 'timestamp left('a', 1)' at line 1
> exec ^
: string
| select * from t as of timestamp left('a', 1) left join u on 1
< error
| 1064 42000 …near 'timestamp left('a', 1) left join u on 1' at line 1
> exec ^
: string
| select * from t as of timestamp a use index (i)
< error
| 1064 42000 …near 'timestamp a use index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp a ignore index (i)
< error
| 1064 42000 …near 'timestamp a ignore index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp a force index (i)
< error
| 1064 42000 …near 'timestamp a force index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp a tablesample (50)
< error
| 1064 42000 …near 'timestamp a tablesample (50)' at line 1
> exec ^
: string
| select * from t as of timestamp a limit 1
< error
| 1064 42000 …near 'timestamp a limit 1' at line 1
> exec ^
: string
| select * from t as of timestamp a from u
< error
| 1064 42000 …near 'timestamp a from u' at line 1
> exec ^
: string
| select * from t as of timestamp a b
< error
| 1064 42000 …near 'timestamp a b' at line 1
> exec ^
: string
| select * from t as of timestamp (select 1) tablesample (50)
< error
| 1064 42000 …near 'timestamp (select 1) tablesample (50)' at line 1
> exec ^
: string
| select * from t as of timestamp a in (1) join u
< error
| 1064 42000 …near 'timestamp a in (1) join u' at line 1
> exec ^
: string
| select * from t as of timestamp not a join u
< error
| 1064 42000 …near 'timestamp not a join u' at line 1
> exec ^
: string
| select * from t as of timestamp timestamp '2020-01-01 00:00:00' tablesample (1)
< error
| 1064 42000 …near 'timestamp timestamp '2020-01-01 00:00:00' tablesample (1)' at line 1
> exec ^
: string
| select * from t as of timestamp 'a' collate utf8mb4_bin
< error
| 1064 42000 …near 'timestamp 'a' collate utf8mb4_bin' at line 1
> exec ^
: string
| select * from t as of timestamp a collate utf8mb4_bin use index (i)
< error
| 1064 42000 …near 'timestamp a collate utf8mb4_bin use index (i)' at line 1
> exec ^
: string
| select * from t as of timestamp interval 1 day + now()
< error
| 1064 42000 …near 'timestamp interval 1 day + now()' at line 1
> exec ^
: string
| select * from t as of timestamp a member of ('[]')
< error
| 1064 42000 …near 'timestamp a member of ('[]')' at line 1
> exec ^
: string
| select * from t as of timestamp binary 'a'
< error
| 1064 42000 …near 'timestamp binary 'a'' at line 1
> exec ^
: string
| select * from t as of timestamp !1
< error
| 1064 42000 …near 'timestamp !1' at line 1
> exec ^
: string
| select * from t as of timestamp ~1
< error
| 1064 42000 …near 'timestamp ~1' at line 1
> exec ^
: string
| select * from t as of timestamp a -> '$.b'
< error
| 1064 42000 …near 'timestamp a -> '$.b'' at line 1
> exec ^
: string
| select * from t as of timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * from u)
< error
| 1064 42000 …near 'timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * fr...' at line 1
> exec ^
: string
| insert into t (select * from u as of timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * from u))
< error
| 1064 42000 …near 'timestamp (date_sub(now(), interval _utf8mb4'60' minute)) not in (select * fr...' at line 1
> exec ^
: string
| plan replayer dump with stats as timestamp now() explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats as timestamp now() explain select 1' at line 1
> exec ^
: string
| plan replayer dump with stats of timestamp now() explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats of timestamp now() explain select 1' at line 1
> exec ^
: string
| plan replayer dump with stats timestamp now() explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats timestamp now() explain select 1' at line 1
> exec ^
: string
| plan replayer dump with stats as of now() explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats as of now() explain select 1' at line 1
> exec ^
: string
| plan replayer dump with stats as of timestamp explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats as of timestamp explain select 1' at line 1
> exec ^
: string
| plan replayer dump with stats as of explain select 1
< error
| 1064 42000 …near 'plan replayer dump with stats as of explain select 1' at line 1
