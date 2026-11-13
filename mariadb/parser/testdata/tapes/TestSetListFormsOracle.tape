-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSetListFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| set offset = 1
< error
| 1193 HY000 Unknown system variable 'offset'
> exec ^
: string
| set offset := 1
< error
| 1193 HY000 Unknown system variable 'offset'
> exec ^
: string
| set session offset = 1
< error
| 1064 42000 …near 'offset = 1' at line 1
> exec ^
: string
| set global offset = 1
< error
| 1064 42000 …near 'offset = 1' at line 1
> exec ^
: string
| set autocommit = 1, offset = 2
< error
| 1193 HY000 Unknown system variable 'offset'
> exec ^
: string
| alter table t lock = none
< ok 0 0
> exec ^
: string
| alter table t page_compressed = on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| alter table t stats_persistent = 0
< ok 0 0
> exec ^
: string
| alter table where.a add column b int
< ok 0 0
> exec ^
: string
| change master to master_host = 'a'
< ok 0 0
> exec ^
: string
| create index i on where.a (b)
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, autocommit = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, global autocommit = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; set a = binary; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; set a = default; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; set a = on; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare a int; set global autocommit = 1, a = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin declare timestamp int; set timestamp = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set @a = 1, names utf8; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set @a = default; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin set @a = on; end
< error
| 1064 42000 …near 'on; end' at line 1
> exec ^
: string
| create procedure p() begin set autocommit = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set autocommit = all; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set autocommit = binary; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set names utf8; end
< ok 0 0
> exec ^
: string
| create table offset (a int)
< error
| 1064 42000 …near 'offset (a int)' at line 1
> exec ^
: string
| create table select.a (x int)
< ok 0 0
> exec ^
: string
| create table select.from (x int)
< ok 0 0
> exec ^
: string
| create table t (a int) page_compressed = on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| create table t (a int) sequence = on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| create table t (a int) stats_auto_recalc = 1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_persistent = 1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_persistent = on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| create table t (a int) stats_sample_pages = 10
< ok 0 0
> exec ^
: string
| create table t (a int) with system versioning
< ok 0 0
> exec ^
: string
| create table t (vector int)
< error
| 1064 42000 …near 'int)' at line 1
> exec ^
: string
| create table vector (a int)
< error
| 1064 42000 …near 'vector (a int)' at line 1
> exec ^
: string
| create trigger t before insert on a for each row set new.a = 1
< ok 0 0
> exec ^
: string
| create trigger t before insert on a for each row set new.a = 1, new.b = 2
< ok 0 0
> exec ^
: string
| create trigger t before insert on a for each row set new.a = 1, old.b = 2
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| create trigger t before insert on a for each row set new.a = 1, old.b = 2
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| create trigger t before insert on a for each row set new.a = binary
< ok 0 0
> exec ^
: string
| create trigger t before insert on a for each row set new.a = default
< ok 0 0
> exec ^
: string
| create trigger t before insert on a for each row set new.a = on
< ok 0 0
> exec ^
: string
| create trigger t before insert on a for each row set timestamp = 1
< ok 0 0
> exec ^
: string
| delete from select.a
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| delete from select.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| delete from where.a
< error
| 1146 42S02 Table 'where.a' doesn't exist
> exec ^
: string
| delete from where.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| drop table select.a
< ok 0 0
> exec ^
: string
| drop table where.a
< ok 0 0
> exec ^
: string
| fetch c into timestamp
< error
| 1064 42000 …near 'fetch c into timestamp' at line 1
> exec ^
: string
| insert into select.a values (1)
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| insert into select.a values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into where.a values (1)
< error
| 1146 42S02 Table 'where.a' doesn't exist
> exec ^
: string
| insert into where.a values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| reset persist timestamp
< error
| 1064 42000 …near 'persist timestamp' at line 1
> exec ^
: string
| select * from from.select
< error
| 1146 42S02 Table 'from.select' doesn't exist
> exec ^
: string
| select * from from.select
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from select.a
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| select * from select.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t for system_time all
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t for system_time all
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t limit 1 offset 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t limit 1 offset 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t limit 1, 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t limit 1, 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t offset
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select 1 from select.a
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| select 1 from select.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 from where.where
< error
| 1146 42S02 Table 'where.where' doesn't exist
> exec ^
: string
| select 1 from where.where
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1. from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select 1. from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select @@ autocommit
< error
| 1064 42000 …near ' autocommit' at line 1
> exec ^
: string
| select @@`autocommit`
< ok 0 0
> exec ^
: string
| select @@`autocommit`.a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@autocommit . a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@global . autocommit
< ok 0 0
> exec ^
: string
| select @@global.a
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| select @@global.autocommit
< ok 0 0
> exec ^
: string
| select @@local.autocommit
< ok 0 0
> exec ^
: string
| select @@select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| select @@select.a
< error
| 1064 42000 …near 'select.a' at line 1
> exec ^
: string
| select @@session.`autocommit`
< ok 0 0
> exec ^
: string
| select @@session.autocommit
< ok 0 0
> exec ^
: string
| select @a := on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| select `select`.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `select`.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a . select from t
< error
| 1064 42000 …near 'select from t' at line 1
> exec ^
: string
| select a from t where where.b = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a from t where where.b = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a from t where.b
< error
| 1064 42000 …near '.b' at line 1
> exec ^
: string
| select a offset from t
< error
| 1064 42000 …near 'from t' at line 1
> exec ^
: string
| select a.select from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.select from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select a.select.b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.select.b from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select all
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select binary
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select case.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select case.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select count.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select count.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select date.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select date.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select default(a) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select default(a) from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select global.a
< error
| 1109 42S02 Unknown table 'global' in SELECT
> exec ^
: string
| select global.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select if.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select if.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select interval.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select interval.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select local.a
< error
| 1109 42S02 Unknown table 'local' in SELECT
> exec ^
: string
| select local.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select new.a
< error
| 1109 42S02 Unknown table 'new' in SELECT
> exec ^
: string
| select new.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select not.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select not.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select now.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select now.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select null.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select null.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select off
< error
| 1054 42S22 Unknown column 'off' in 'SELECT'
> exec ^
: string
| select off
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select off from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select off from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select offset from t
< error
| 1064 42000 …near 'offset from t' at line 1
> exec ^
: string
| select old.a
< error
| 1109 42S02 Unknown table 'old' in SELECT
> exec ^
: string
| select old.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| select persist.a
< error
| 1109 42S02 Unknown table 'persist' in SELECT
> exec ^
: string
| select persist.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select row
< error
| 1054 42S22 Unknown column 'row' in 'SELECT'
> exec ^
: string
| select row
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select row from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select row from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select select .a from t
< error
| 1064 42000 …near 'select .a from t' at line 1
> exec ^
: string
| select select.
< error
| 1064 42000 …near 'select.' at line 1
> exec ^
: string
| select select. a from t
< error
| 1064 42000 …near 'select. a from t' at line 1
> exec ^
: string
| select select.`a` from t
< error
| 1064 42000 …near 'select.`a` from t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select select.a.b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select select.a.b from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select session.a
< error
| 1109 42S02 Unknown table 'session' in SELECT
> exec ^
: string
| select session.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select show.* from t
< error
| 1064 42000 …near 'show.* from t' at line 1
> exec ^
: string
| select show.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select show.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select t.where from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select t.where from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select table.* from table.a
< error
| 1064 42000 …near 'table.* from table.a' at line 1
> exec ^
: string
| select true.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select true.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select where.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select where.a from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select where.where from where
< error
| 1064 42000 …near 'where' at line 1
> exec ^
: string
| set , autocommit = 1
< error
| 1064 42000 …near ' autocommit = 1' at line 1
> exec ^
: string
| set @@ autocommit = 1
< error
| 1064 42000 …near ' autocommit = 1' at line 1
> exec ^
: string
| set @@ global.autocommit = 1
< error
| 1064 42000 …near ' global.autocommit = 1' at line 1
> exec ^
: string
| set @@'a' = 1
< error
| 1064 42000 …near ''a' = 1' at line 1
> exec ^
: string
| set @@a . b = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'a'
> exec ^
: string
| set @@autocommit = 1
< ok 0 0
> exec ^
: string
| set @@autocommit = 1, names utf8
< ok 0 0
> exec ^
: string
| set @@autocommit = all
< ok 0 0
> exec ^
: string
| set @@autocommit = binary
< ok 0 0
> exec ^
: string
| set @@autocommit = default
< ok 0 0
> exec ^
: string
| set @@autocommit = on
< ok 0 0
> exec ^
: string
| set @@autocommit = row
< ok 0 0
> exec ^
: string
| set @@autocommit = system
< ok 0 0
> exec ^
: string
| set @@global . autocommit = 1
< ok 0 0
> exec ^
: string
| set @@global.autocommit = 1
< ok 0 0
> exec ^
: string
| set @@local.autocommit = 1
< ok 0 0
> exec ^
: string
| set @@persist.autocommit = 1
< error
| 1272 HY000 Variable 'autocommit' is not a variable component (can't be used as XXXX.variable_name)
> exec ^
: string
| set @@persist_only.autocommit = 1
< error
| 1272 HY000 Variable 'autocommit' is not a variable component (can't be used as XXXX.variable_name)
> exec ^
: string
| set @@select.a = 1
< error
| 1064 42000 …near 'select.a = 1' at line 1
> exec ^
: string
| set @@session.autocommit = 1
< ok 0 0
> exec ^
: string
| set @a = 1, autocommit = 1
< ok 0 0
> exec ^
: string
| set @a = 1, character set utf8
< ok 0 0
> exec ^
: string
| set @a = 1, charset utf8
< ok 0 0
> exec ^
: string
| set @a = 1, names utf8
< ok 0 0
> exec ^
: string
| set @a = all
< error
| 1064 42000 …near 'all' at line 1
> exec ^
: string
| set @a = binary
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = binary 'x'
< ok 0 0
> exec ^
: string
| set @a = default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = default(b)
< error
| 1054 42S22 Unknown column 'b' in 'SET'
> exec ^
: string
| set @a = default(b)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = off
< error
| 1054 42S22 Unknown column 'off' in 'SET'
> exec ^
: string
| set @a = off
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| set @a = row
< error
| 1054 42S22 Unknown column 'row' in 'SET'
> exec ^
: string
| set @a = row
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = row(1, 2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| set @a = row(1, 2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set `select` = 1
< error
| 1193 HY000 Unknown system variable 'select'
> exec ^
: string
| set a = 1 , names utf8 collate utf8_bin, b = 2
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set a = 1, char set utf8
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set a.timestamp = 1
< error
| 1272 HY000 Variable 'timestamp' is not a variable component (can't be used as XXXX.variable_name)
> exec ^
: string
| set autocommit
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set autocommit . select = 1
< error
| 1064 42000 …near 'select = 1' at line 1
> exec ^
: string
| set autocommit := 1
< ok 0 0
> exec ^
: string
| set autocommit =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set autocommit = 'off'
< ok 0 0
> exec ^
: string
| set autocommit = 'on'
< ok 0 0
> exec ^
: string
| set autocommit = (on)
< error
| 1064 42000 …near 'on)' at line 1
> exec ^
: string
| set autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1 + 1
< ok 0 0
> exec ^
: string
| set autocommit = 1,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set autocommit = 1, @@autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, @a = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, autocommit = 0
< ok 0 0
> exec ^
: string
| set autocommit = 1, character set utf8
< ok 0 0
> exec ^
: string
| set autocommit = 1, charset = 2
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set autocommit = 1, charset utf8
< ok 0 0
> exec ^
: string
| set autocommit = 1, global autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, global.a = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'global'
> exec ^
: string
| set autocommit = 1, local autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, names = 2
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set autocommit = 1, names utf8
< ok 0 0
> exec ^
: string
| set autocommit = 1, password := 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set autocommit = 1, password = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set autocommit = 1, password.a = 2
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'password'
> exec ^
: string
| set autocommit = 1, persist = 1
< error
| 1193 HY000 Unknown system variable 'persist'
> exec ^
: string
| set autocommit = 1, session autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, session names = 2
< error
| 1193 HY000 Unknown system variable 'names'
> exec ^
: string
| set autocommit = 1, session password = 1
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set autocommit = 1, timestamp = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set autocommit = all
< ok 0 0
> exec ^
: string
| set autocommit = all + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set autocommit = all, autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = binary
< ok 0 0
> exec ^
: string
| set autocommit = binary 'x'
< ok 0 0
> exec ^
: string
| set autocommit = binary 1 + 1
< ok 0 0
> exec ^
: string
| set autocommit = binary, autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = default
< ok 0 0
> exec ^
: string
| set autocommit = default + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set autocommit = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'SET'
> exec ^
: string
| set autocommit = default(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set autocommit = default, autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = off
< ok 0 0
> exec ^
: string
| set autocommit = on
< ok 0 0
> exec ^
: string
| set autocommit = on + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set autocommit = on, autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = row
< ok 0 0
> exec ^
: string
| set autocommit = row + 1
< error
| 1054 42S22 Unknown column 'row' in 'SET'
> exec ^
: string
| set autocommit = row + 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set autocommit = row(1, 2)
< error
| 1241 21000 Operand should contain 1 column(s)
> exec ^
: string
| set autocommit = row(1, 2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set autocommit = row, autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = system
< ok 0 0
> exec ^
: string
| set autocommit = system, autocommit = 1
< ok 0 0
> exec ^
: string
| set binding = 1
< error
| 1193 HY000 Unknown system variable 'binding'
> exec ^
: string
| set char set 'utf8'
< ok 0 0
> exec ^
: string
| set char set binary
< ok 0 0
> exec ^
: string
| set char set default
< ok 0 0
> exec ^
: string
| set char set utf8
< ok 0 0
> exec ^
: string
| set character set 'utf8'
< ok 0 0
> exec ^
: string
| set character set 'utf8mb4'
< ok 0 0
> exec ^
: string
| set character set binary
< ok 0 0
> exec ^
: string
| set character set default
< ok 0 0
> exec ^
: string
| set character set utf8
< ok 0 0
> exec ^
: string
| set character set utf8, autocommit = 1
< ok 0 0
> exec ^
: string
| set charset 'utf8'
< ok 0 0
> exec ^
: string
| set charset binary
< ok 0 0
> exec ^
: string
| set charset default
< ok 0 0
> exec ^
: string
| set charset utf8
< ok 0 0
> exec ^
: string
| set charset utf8 collate utf8_bin
< error
| 1064 42000 …near 'collate utf8_bin' at line 1
> exec ^
: string
| set charset utf8, autocommit = 1
< ok 0 0
> exec ^
: string
| set date = 1
< error
| 1193 HY000 Unknown system variable 'date'
> exec ^
: string
| set default = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set default role a
< error
| 1959 OP000 Invalid role specification `a`
> exec ^
: string
| set default role a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set default role a for b
< error
| 1959 OP000 Invalid role specification `a`
> exec ^
: string
| set default role a for b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set default role none
< ok 0 0
> exec ^
: string
| set from = 1
< error
| 1064 42000 …near 'from = 1' at line 1
> exec ^
: string
| set from.a = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'from'
> exec ^
: string
| set global @@autocommit = 1
< error
| 1064 42000 …near '@@autocommit = 1' at line 1
> exec ^
: string
| set global autocommit = 1
< ok 0 0
> exec ^
: string
| set global autocommit = 1, @a = 1, autocommit = 1
< ok 0 0
> exec ^
: string
| set global autocommit = 1, autocommit = 0
< ok 0 0
> exec ^
: string
| set global autocommit = 1, names utf8
< ok 0 0
> exec ^
: string
| set global autocommit = 1, session autocommit = 0
< ok 0 0
> exec ^
: string
| set global autocommit = 1, session autocommit = 0, autocommit = 1
< ok 0 0
> exec ^
: string
| set global autocommit = all
< ok 0 0
> exec ^
: string
| set global autocommit = binary
< ok 0 0
> exec ^
: string
| set global names utf8
< error
| 1064 42000 …near 'utf8' at line 1
> exec ^
: string
| set global timestamp = 1
< error
| 1228 HY000 Variable 'timestamp' is a SESSION variable and can't be used with SET GLOBAL
> exec ^
: string
| set global timestamp = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set global transaction read only
< ok 0 0
> exec ^
: string
| set global.autocommit = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'global'
> exec ^
: string
| set local autocommit = 1
< ok 0 0
> exec ^
: string
| set local transaction read only
< ok 0 0
> exec ^
: string
| set names 'utf8'
< ok 0 0
> exec ^
: string
| set names @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| set names binary
< ok 0 0
> exec ^
: string
| set names binary collate binary
< ok 0 0
> exec ^
: string
| set names binary collate default
< ok 0 0
> exec ^
: string
| set names default
< ok 0 0
> exec ^
: string
| set names default collate utf8_bin
< error
| 1253 42000 COLLATION 'utf8mb3_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| set names utf8
< ok 0 0
> exec ^
: string
| set names utf8 collate 'utf8_bin'
< ok 0 0
> exec ^
: string
| set names utf8 collate binary
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'utf8mb3'
> exec ^
: string
| set names utf8 collate default
< ok 0 0
> exec ^
: string
| set names utf8 collate utf8_bin
< ok 0 0
> exec ^
: string
| set names utf8 collate utf8_bin, autocommit = 1
< ok 0 0
> exec ^
: string
| set names utf8, autocommit = 1
< ok 0 0
> exec ^
: string
| set names utf8, charset utf8
< ok 0 0
> exec ^
: string
| set names utf8, names latin1
< ok 0 0
> exec ^
: string
| set new = 1
< error
| 1193 HY000 Unknown system variable 'new'
> exec ^
: string
| set new.a.b = 1
< error
| 1064 42000 …near '.b = 1' at line 1
> exec ^
: string
| set old = 1
< ok 0 0
> exec ^
: string
| set password = 'x'
< ok 0 0
> exec ^
: string
| set password = password('x')
< ok 0 0
> exec ^
: string
| set password for u = 'x'
< ok 0 0
> exec ^
: string
| set persist autocommit = 1
< error
| 1064 42000 …near 'autocommit = 1' at line 1
> exec ^
: string
| set persist_only autocommit = 1
< error
| 1064 42000 …near 'autocommit = 1' at line 1
> exec ^
: string
| set resource = 1
< error
| 1193 HY000 Unknown system variable 'resource'
> exec ^
: string
| set resource group rg
< error
| 1064 42000 …near 'group rg' at line 1
> exec ^
: string
| set role a
< error
| 1959 OP000 Invalid role specification `a`
> exec ^
: string
| set role a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set role none
< ok 0 0
> exec ^
: string
| set select = 1
< error
| 1064 42000 …near 'select = 1' at line 1
> exec ^
: string
| set select. a = 1
< error
| 1064 42000 …near 'select. a = 1' at line 1
> exec ^
: string
| set select.a = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'select'
> exec ^
: string
| set session autocommit = 1
< ok 0 0
> exec ^
: string
| set session global = 1
< error
| 1193 HY000 Unknown system variable 'global'
> exec ^
: string
| set session names = 1
< error
| 1193 HY000 Unknown system variable 'names'
> exec ^
: string
| set session password = 1
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set session session = 1
< error
| 1193 HY000 Unknown system variable 'session'
> exec ^
: string
| set session timestamp = 1
< ok 0 0
> exec ^
: string
| set session transaction read only
< ok 0 0
> exec ^
: string
| set statement = 1
< error
| 1193 HY000 Unknown system variable 'statement'
> exec ^
: string
| set status = 1
< error
| 1193 HY000 Unknown system variable 'status'
> exec ^
: string
| set time = 1
< error
| 1193 HY000 Unknown system variable 'time'
> exec ^
: string
| set timestamp = 1
< ok 0 0
> exec ^
: string
| set timestamp.a = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'timestamp'
> exec ^
: string
| set transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set transaction read only
< ok 0 0
> exec ^
: string
| set transaction read write
< ok 0 0
> exec ^
: string
| set user = 1
< error
| 1193 HY000 Unknown system variable 'user'
> exec ^
: string
| show columns from select.a
< ok 0 0
> exec ^
: string
| show columns from where.a
< ok 0 0
> exec ^
: string
| show tables from select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| start slave until master_log_file = 'a', master_log_pos = 1
< ok 0 0
> exec ^
: string
| start slave until sql_after_gtids = 'a'
< ok 0 0
> exec ^
: string
| truncate table where.a
< ok 0 0
> exec ^
: string
| update select.a set b = 1
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| update select.a set b = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update where.a set b = 1
< error
| 1146 42S02 Table 'where.a' doesn't exist
> exec ^
: string
| update where.a set b = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
