-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSetListFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| set names utf8 collate @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| alter table t lock = none
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table t page_compressed = on
< error
| 1064 42000 …near 'page_compressed = on' at line 1
> exec ^
: string
| alter table t stats_persistent = 0
< error
| 1046 3D000 No database selected
> exec ^
: string
| alter table where.a add column b int
< error
| 1146 42S02 Table 'where.a' doesn't exist
> exec ^
: string
| create index i on where.a (b)
< error
| 1049 42000 Unknown database
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, autocommit = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, global autocommit = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a = binary; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a = default; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a = on; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set global autocommit = 1, a = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare timestamp int; set timestamp = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set @a = default; end
< error
| 1064 42000 …near 'set @a = default' at line 1
> exec ^
: string
| create procedure p() begin set @a = on; end
< error
| 1064 42000 …near 'set @a = on' at line 1
> exec ^
: string
| create procedure p() begin set autocommit = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set autocommit = all; end
< error
| 1064 42000 …near 'set autocommit = all' at line 1
> exec ^
: string
| create procedure p() begin set autocommit = binary; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set names utf8; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create table offset (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table select.a (x int)
< error
| 1049 42000 Unknown database 'select'
> exec ^
: string
| create table select.from (x int)
< error
| 1049 42000 Unknown database 'select'
> exec ^
: string
| create table t (a int) page_compressed = on
< error
| 1064 42000 …near 'page_compressed = on' at line 1
> exec ^
: string
| create table t (a int) sequence = on
< error
| 1064 42000 …near 'sequence = on' at line 1
> exec ^
: string
| create table t (a int) stats_auto_recalc = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) stats_persistent = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) stats_sample_pages = 10
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table t (a int) with system versioning
< error
| 1064 42000 …near 'versioning' at line 1
> exec ^
: string
| create table t (vector int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create table vector (a int)
< error
| 1046 3D000 No database selected
> exec ^
: string
| create trigger t before insert on a for each row set new.a = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| create trigger t before insert on a for each row set new.a = 1, old.b = 2
< error
| 1046 3D000 No database selected
> exec ^
: string
| create trigger t before insert on a for each row set new.a = binary
< error
| 1046 3D000 No database selected
> exec ^
: string
| create trigger t before insert on a for each row set new.a = default
< error
| 1046 3D000 No database selected
> exec ^
: string
| create trigger t before insert on a for each row set new.a = on
< error
| 1046 3D000 No database selected
> exec ^
: string
| create trigger t before insert on a for each row set timestamp = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| delete from select.a
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| delete from where.a
< error
| 1146 42S02 Table 'where.a' doesn't exist
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
| insert into where.a values (1)
< error
| 1146 42S02 Table 'where.a' doesn't exist
> exec ^
: string
| select * from from.select
< error
| 1146 42S02 Table 'from.select' doesn't exist
> exec ^
: string
| select * from select.a
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| select * from t for system_time all
< error
| 1064 42000 …near 'system_time all' at line 1
> exec ^
: string
| select * from t limit 1 offset 2
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t limit 1, 2
< error
| 1046 3D000 No database selected
> exec ^
: string
| select * from t offset
< error
| 1046 3D000 No database selected
> exec ^
: string
| select 1 from select.a
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| select 1 from where.where
< error
| 1146 42S02 Table 'where.where' doesn't exist
> exec ^
: string
| select 1. from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select @@ autocommit
< error
| 1064 42000 …near '@ autocommit' at line 1
> exec ^
: string
| select @@`autocommit`
< ok 0 0
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
| select @@select.a
< error
| 1064 42000 …near '.a' at line 1
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
| 1046 3D000 No database selected
> exec ^
: string
| select a . select from t
< error
| 1064 42000 …near 'select from t' at line 1
> exec ^
: string
| select a from t where where.b = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a from t where.b
< error
| 1064 42000 …near 'where.b' at line 1
> exec ^
: string
| select a offset from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a.select from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select a.select.b from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select all
< error
| 1064 42000 …near 'all' at line 1
> exec ^
: string
| select binary
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| select case.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select count.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select date.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select default(a) from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select from.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select global.a
< error
| 1109 42S02 Unknown table 'global' in field list
> exec ^
: string
| select if.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select interval.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select local.a
< error
| 1109 42S02 Unknown table 'local' in field list
> exec ^
: string
| select new.a
< error
| 1109 42S02 Unknown table 'new' in field list
> exec ^
: string
| select not.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select now.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select null.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select off
< error
| 1054 42S22 Unknown column 'off' in 'field list'
> exec ^
: string
| select off from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select offset from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select old.a
< error
| 1109 42S02 Unknown table 'old' in field list
> exec ^
: string
| select on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| select persist.a
< error
| 1109 42S02 Unknown table 'persist' in field list
> exec ^
: string
| select row
< error
| 1054 42S22 Unknown column 'row' in 'field list'
> exec ^
: string
| select row from t
< error
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
> exec ^
: string
| select select.a.b from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select session.a
< error
| 1109 42S02 Unknown table 'session' in field list
> exec ^
: string
| select show.* from t
< error
| 1064 42000 …near '* from t' at line 1
> exec ^
: string
| select show.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select t.where from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select table.* from table.a
< error
| 1064 42000 …near '* from table.a' at line 1
> exec ^
: string
| select true.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select where.a from t
< error
| 1046 3D000 No database selected
> exec ^
: string
| select where.where from where
< error
| 1064 42000 …near 'where' at line 1
> exec ^
: string
| set , autocommit = 1
< error
| 1064 42000 …near ', autocommit = 1' at line 1
> exec ^
: string
| set @@ autocommit = 1
< error
| 1064 42000 …near '@ autocommit = 1' at line 1
> exec ^
: string
| set @@ global.autocommit = 1
< error
| 1064 42000 …near '@ global.autocommit = 1' at line 1
> exec ^
: string
| set @@'a' = 1
< error
| 1064 42000 …near '@'a' = 1' at line 1
> exec ^
: string
| set @@a . b = 1
< error
| 1064 42000 …near '. b = 1' at line 1
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
< error
| 1064 42000 …near 'all' at line 1
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
| 1064 42000 …near '.autocommit = 1' at line 1
> exec ^
: string
| set @@persist_only.autocommit = 1
< error
| 1064 42000 …near '.autocommit = 1' at line 1
> exec ^
: string
| set @@select.a = 1
< error
| 1064 42000 …near '.a = 1' at line 1
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
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| set @a = binary 'x'
< ok 0 0
> exec ^
: string
| set @a = default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @a = default(b)
< error
| 1054 42S22 Unknown column 'b' in 'field_list'
> exec ^
: string
| set @a = off
< error
| 1054 42S22 Unknown column 'off' in 'field_list'
> exec ^
: string
| set @a = on
< error
| 1064 42000 …near 'on' at line 1
> exec ^
: string
| set @a = row
< error
| 1054 42S22 Unknown column 'row' in 'field_list'
> exec ^
: string
| set @a = row(1, 2)
< ok 0 0
> exec ^
: string
| set `select` = 1
< ok 0 0
> exec ^
: string
| set a = 1 , names utf8 collate utf8_bin, b = 2
< ok 0 0
> exec ^
: string
| set a = 1, char set utf8
< ok 0 0
> exec ^
: string
| set a.1 = 1
< error
| 1064 42000 …near 'a.1 = 1' at line 1
> exec ^
: string
| set a.timestamp = 1
< error
| 1064 42000 …near 'a.timestamp = 1' at line 1
> exec ^
: string
| set autocommit
< error
| 1064 42000 …near 'autocommit' at line 1
> exec ^
: string
| set autocommit . select = 1
< error
| 1064 42000 …near '. select = 1' at line 1
> exec ^
: string
| set autocommit := 1
< ok 0 0
> exec ^
: string
| set autocommit =
< error
| 1064 42000 …near '=' at line 1
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
| 1064 42000 …near ')' at line 1
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
| 1064 42000 …near ',' at line 1
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
< ok 0 0
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
| 1064 42000 …near 'global.a = 1' at line 1
> exec ^
: string
| set autocommit = 1, local autocommit = 1
< error
| 1064 42000 …near 'autocommit = 1' at line 1
> exec ^
: string
| set autocommit = 1, names = 2
< ok 0 0
> exec ^
: string
| set autocommit = 1, names utf8
< ok 0 0
> exec ^
: string
| set autocommit = 1, password := 2
< error
| 1064 42000 …near 'password := 2' at line 1
> exec ^
: string
| set autocommit = 1, password = 1
< error
| 1064 42000 …near 'password = 1' at line 1
> exec ^
: string
| set autocommit = 1, password.a = 2
< error
| 1064 42000 …near 'password.a = 2' at line 1
> exec ^
: string
| set autocommit = 1, persist = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, session autocommit = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, session names = 2
< ok 0 0
> exec ^
: string
| set autocommit = 1, session password = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, timestamp = 1
< ok 0 0
> exec ^
: string
| set autocommit = 1, transaction = 1
< ok 0 0
> exec ^
: string
| set autocommit = all
< error
| 1064 42000 …near 'all' at line 1
> exec ^
: string
| set autocommit = all + 1
< error
| 1064 42000 …near 'all + 1' at line 1
> exec ^
: string
| set autocommit = all, autocommit = 1
< error
| 1064 42000 …near 'all, autocommit = 1' at line 1
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
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1054 42S22 Unknown column 'row' in 'field_list'
> exec ^
: string
| set autocommit = row(1, 2)
< ok 0 0
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
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set default = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set default role a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| set default role a for b
< error
| 1064 42000 …near 'for b' at line 1
> exec ^
: string
| set default role none
< error
| 1064 42000 …near 'none' at line 1
> exec ^
: string
| set from = 1
< error
| 1064 42000 …near 'from = 1' at line 1
> exec ^
: string
| set from.a = 1
< error
| 1064 42000 …near 'from.a = 1' at line 1
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
< error
| 1064 42000 …near 'all' at line 1
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
< ok 0 0
> exec ^
: string
| set global transaction read only
< ok 0 0
> exec ^
: string
| set global.autocommit = 1
< error
| 1064 42000 …near 'global.autocommit = 1' at line 1
> exec ^
: string
| set local autocommit = 1
< error
| 1064 42000 …near 'autocommit = 1' at line 1
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
< error
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| set names binary collate default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set names default
< ok 0 0
> exec ^
: string
| set names default collate utf8_bin
< ok 0 0
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
| 1064 42000 …near 'binary' at line 1
> exec ^
: string
| set names utf8 collate default
< error
| 1064 42000 …near 'default' at line 1
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
< ok 0 0
> exec ^
: string
| set new.a.b = 1
< error
| 1064 42000 …near 'new.a.b = 1' at line 1
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
| set password = old_password('x')
< error
| 1064 42000 …near 'old_password('x')' at line 1
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
| set resource group rg
< error
| 1064 42000 …near 'resource group rg' at line 1
> exec ^
: string
| set role a
< error
| 3523 HY000 Unknown authorization ID `a`@`%`
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
| 1064 42000 …near 'select.a = 1' at line 1
> exec ^
: string
| set session autocommit = 1
< ok 0 0
> exec ^
: string
| set session global = 1
< ok 0 0
> exec ^
: string
| set session names = 1
< ok 0 0
> exec ^
: string
| set session password = 1
< ok 0 0
> exec ^
: string
| set session session = 1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| set statement autocommit = 1 for select 1
< error
| 1064 42000 …near 'autocommit = 1 for select 1' at line 1
> exec ^
: string
| set status = 1
< ok 0 0
> exec ^
: string
| set time = 1
< ok 0 0
> exec ^
: string
| set timestamp = 1
< ok 0 0
> exec ^
: string
| set timestamp.a = 1
< error
| 1064 42000 …near 'timestamp.a = 1' at line 1
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
< ok 0 0
> exec ^
: string
| show columns from select.a
< error
| 1049 42000 Unknown database
> exec ^
: string
| show columns from where.a
< error
| 1049 42000 Unknown database
> exec ^
: string
| show tables from select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| truncate table where.a
< error
| 1146 42S02 Table 'where.a' doesn't exist
> exec ^
: string
| update select.a set b = 1
< error
| 1146 42S02 Table 'select.a' doesn't exist
> exec ^
: string
| update where.a set b = 1
< error
| 1146 42S02 Table 'where.a' doesn't exist
