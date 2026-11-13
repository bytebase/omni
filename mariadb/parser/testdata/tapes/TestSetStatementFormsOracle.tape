-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSetStatementFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SET STATEMENT MAX_STATEMENT_TIME = 1 FOR SELECT 1
< ok 0 0
> exec ^
: string
| set /* c */ statement /* c */ max_statement_time = 1 /* c */ for /* c */ select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time := 1 for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time to 1 for select 1
< error
| 1064 42000 …near 'to 1 for select 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1, sql_mode = '' for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1, sql_mode = '', foreign_key_checks = 0 for select 1
< ok 0 0
> exec ^
: string
| set statement for select 1
< error
| 1064 42000 …near 'for select 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1, for select 1
< error
| 1064 42000 …near 'for select 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| set statement , max_statement_time = 1 for select 1
< error
| 1064 42000 …near ' max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement max_statement_time for select 1
< error
| 1064 42000 …near 'for select 1' at line 1
> exec ^
: string
| set statement max_statement_time = for select 1
< error
| 1064 42000 …near 'for select 1' at line 1
> exec ^
: string
| set statement @a = 1 for select 1
< error
| 1064 42000 …near '@a = 1 for select 1' at line 1
> exec ^
: string
| set statement @@max_statement_time = 1 for select 1
< error
| 1064 42000 …near '@@max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement @@session.max_statement_time = 1 for select 1
< error
| 1064 42000 …near '@@session.max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement global max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement session max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement local max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set global statement max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set session statement max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set @a = 1, statement max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set max_statement_time = 1, statement max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement names utf8mb4 for select 1
< error
| 1064 42000 …near 'utf8mb4 for select 1' at line 1
> exec ^
: string
| set statement charset utf8mb4 for select 1
< error
| 1064 42000 …near 'utf8mb4 for select 1' at line 1
> exec ^
: string
| set statement character set utf8mb4 for select 1
< error
| 1064 42000 …near 'character set utf8mb4 for select 1' at line 1
> exec ^
: string
| set statement transaction read only for select 1
< error
| 1064 42000 …near 'read only for select 1' at line 1
> exec ^
: string
| set statement default.key_buffer_size = 1 for select 1
< error
| 1229 HY000 Variable 'key_buffer_size' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set statement default.key_buffer_size = 1 for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement default . key_buffer_size = 1 for select 1
< error
| 1229 HY000 Variable 'key_buffer_size' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set statement default . key_buffer_size = 1 for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement a.key_buffer_size = 1 for select 1
< error
| 1229 HY000 Variable 'key_buffer_size' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set statement a.key_buffer_size = 1 for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement `max_statement_time` = 1 for select 1
< ok 0 0
> exec ^
: string
| set `statement` max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| set statement to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set statement
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement sql_mode = on for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'ON'
> exec ^
: string
| set statement sql_mode = on for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement sql_mode = off for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'off'
> exec ^
: string
| set statement sql_mode = off for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement sql_mode = default for select 1
< ok 0 0
> exec ^
: string
| set statement sql_mode = all for select 1
< ok 0 0
> exec ^
: string
| set statement sql_mode = binary for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'binary'
> exec ^
: string
| set statement sql_mode = binary for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement sql_mode = 'a' for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'a'
> exec ^
: string
| set statement sql_mode = 'a' for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 + 1 for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = (1) for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = (select 1) for select 1
< ok 0 0
> exec ^
: string
| set statement sql_mode = default + 1 for select 1
< error
| 1064 42000 …near '+ 1 for select 1' at line 1
> exec ^
: string
| set statement sql_mode = on + 1 for select 1
< error
| 1064 42000 …near '+ 1 for select 1' at line 1
> exec ^
: string
| set statement max_statement_time = @a for select 1
< error
| 1232 42000 Incorrect argument type to variable 'max_statement_time'
> exec ^
: string
| set statement max_statement_time = @a for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = @@max_statement_time for select 1
< ok 0 0
> exec ^
: string
| set statement sql_mode = system for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'system'
> exec ^
: string
| set statement sql_mode = system for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement sql_mode = binary, max_statement_time = 1 for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'binary'
> exec ^
: string
| set statement sql_mode = binary, max_statement_time = 1 for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement sql_mode = binary 'a' for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'a'
> exec ^
: string
| set statement sql_mode = binary 'a' for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement sql_mode = all, max_statement_time = 1 for select 1
< ok 0 0
> exec ^
: string
| set statement timestamp = 1 for select 1
< ok 0 0
> exec ^
: string
| set statement select = 1 for select 1
< error
| 1064 42000 …near 'select = 1 for select 1' at line 1
> exec ^
: string
| set statement for = 1 for select 1
< error
| 1064 42000 …near 'for = 1 for select 1' at line 1
> exec ^
: string
| set statement offset = 1 for select 1
< error
| 1064 42000 …near 'offset = 1 for select 1' at line 1
> exec ^
: string
| set statement sql_mode = on, max_statement_time = 1 for select 1
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'ON'
> exec ^
: string
| set statement sql_mode = on, max_statement_time = 1 for select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set Statement sql_mode = '' for select 1
< ok 0 0
> exec ^
: string
| l: loop select 1; end loop
< error
| 1064 42000 …near 'l: loop select 1; end loop' at line 1
> exec ^
: string
| set statement sql_mode = (default) for select 1
< error
| 1064 42000 …near ') for select 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for select 1;
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for for select 1
< error
| 1064 42000 …near 'select 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for select * from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for select * from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for select 1 for update
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for select 1 into @a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for (select 1)
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for values (1)
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for with c as (select 1) select * from c
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for insert into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for insert into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for insert into t select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for insert into t select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for replace into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for replace into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for update t set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for update t set a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for delete from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for delete from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for create table x (a int)
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create table x as select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create temporary table x (a int)
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for alter table t add column b int
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for alter table t drop column a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for drop table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for truncate table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for rename table t to x
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create index i on t (a)
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for drop index i on t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create view v as select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for drop view v
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create database d
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for drop database d
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create procedure p() select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create function f() returns int return 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create trigger g before insert on t for each row set @a = 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create event e on schedule every 1 day do select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create sequence s
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for drop procedure p
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set @a = 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set sql_mode = ''
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set global sql_mode = ''
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set names utf8mb4
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set transaction read only
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set password = 'x'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set role none
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set statement sql_mode = '' for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for set statement sql_mode = '' for set statement sql_mode = '' for select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for show tables
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for show create table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for show create table t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for describe t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for explain select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for use test
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for begin
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for begin work
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for begin end
< error
| 1064 42000 …near 'end' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for start transaction
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for commit
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for rollback
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for savepoint a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for release savepoint a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for lock tables t read
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for unlock tables
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for flush tables
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for prepare s from 'select 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| set statement max_statement_time = 1 for prepare s from 'select 1'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for execute s
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| set statement max_statement_time = 1 for execute s
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for deallocate prepare s
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| set statement max_statement_time = 1 for deallocate prepare s
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for call p()
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for do 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for grant select on t to a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for revoke select on t from a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for create user a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for drop user a
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for load data infile 'x' into table t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for load data infile 'x' into table t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for analyze table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for optimize table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for check table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for repair table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for checksum table t
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for kill 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for xa start 'a'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for handler t open
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for signal sqlstate '45000'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for if 1 then select 1; end if
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for while 1 do select 1; end while
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for repeat select 1; until 1 end repeat
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for loop select 1; end loop
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for case 1 when 1 then select 1; end case
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for l: loop select 1; end loop
< error
| 1064 42000 …near 'l: loop select 1; end loop' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for select 1; select 2
< error
| 1064 42000 …near 'select 2' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for 
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for reset master
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for shutdown
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for install plugin a soname 'b'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for get diagnostics @a = number
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for help 'a'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for change master to master_host = 'a'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for start slave
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for stop slave
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for purge binary logs to 'a'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for binlog 'a'
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for return 1
< error
| 1064 42000 …near 'return 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for declare a int
< error
| 1064 42000 …near 'declare a int' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for leave l
< error
| 1064 42000 …near 'leave l' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for iterate l
< error
| 1064 42000 …near 'iterate l' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for open c
< error
| 1064 42000 …near 'open c' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for close c
< error
| 1064 42000 …near 'close c' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for fetch c into a
< error
| 1064 42000 …near 'fetch c into a' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for a := 1
< error
| 1064 42000 …near 'a := 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for a = 1
< error
| 1064 42000 …near 'a = 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for p()
< error
| 1064 42000 …near 'p()' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for @a := 1
< error
| 1064 42000 …near '@a := 1' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for set statement
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for -- c
|  select 1
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for /* c */ if 1 then select 1; end if
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for delete from t where a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for delete from t where a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for update t set a = 1 where a = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| set statement max_statement_time = 1 for update t set a = 1 where a = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set statement max_statement_time = 1 for if 1 then select 1; else select 2; end if
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for case when 1 then select 1; end case
< ok 0 0
> exec ^
: string
| set statement max_statement_time = 1 for if 1 then select 1
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for while 1 do select 1; end
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for repeat select 1; until 1
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for loop select 1; end
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set statement max_statement_time = 1 for loop select 1; end loop l
< error
| 1064 42000 …near 'l' at line 1
> exec ^
: string
| create procedure p() set statement max_statement_time = 1 for select 1
< ok 0 0
> exec ^
: string
| create procedure p() set statement max_statement_time = 1 for
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create procedure p() set statement max_statement_time = 1 for set @a = 1
< ok 0 0
> exec ^
: string
| create procedure p() set statement max_statement_time = 1 for begin select 1; end
< error
| 1064 42000 …near 'select 1; end' at line 1
> exec ^
: string
| create procedure p() set statement max_statement_time = 1 for if 1 then select 1; end if
< ok 0 0
> exec ^
: string
| create procedure p() begin set statement max_statement_time = 1 for select 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set statement max_statement_time = 1 for; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin set statement max_statement_time = 1 for set @a = 1; end
< ok 0 0
> exec ^
: string
| create procedure p() begin set statement max_statement_time = 1 for begin select 1; end; end
< error
| 1064 42000 …near 'select 1; end; end' at line 1
> exec ^
: string
| create procedure p() begin set statement max_statement_time = 1 for if 1 then select 1; end if; end
< ok 0 0
> exec ^
: string
| create function f() returns int begin set statement max_statement_time = 1 for do 1; return 1; end
< ok 0 0
> exec ^
: string
| create function f() returns int begin set statement max_statement_time = 1 for; return 1; end
< error
| 1064 42000 …near '; return 1; end' at line 1
> exec ^
: string
| create function f() returns int begin set statement max_statement_time = 1 for set @a = 1; return 1; end
< ok 0 0
> exec ^
: string
| create function f() returns int begin set statement max_statement_time = 1 for begin select 1; end; return 1; end
< error
| 1064 42000 …near 'select 1; end; return 1; end' at line 1
> exec ^
: string
| create function f() returns int begin set statement max_statement_time = 1 for if 1 then do 1; end if; return 1; end
< ok 0 0
> exec ^
: string
| create trigger g before insert on t for each row set statement max_statement_time = 1 for do 1
< ok 0 0
> exec ^
: string
| create trigger g before insert on t for each row set statement max_statement_time = 1 for
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create trigger g before insert on t for each row set statement max_statement_time = 1 for set @a = 1
< ok 0 0
> exec ^
: string
| create trigger g before insert on t for each row set statement max_statement_time = 1 for if 1 then do 1; end if
< ok 0 0
> exec ^
: string
| create trigger g before insert on t for each row begin set statement max_statement_time = 1 for do 1; end
< ok 0 0
> exec ^
: string
| create trigger g before insert on t for each row begin set statement max_statement_time = 1 for; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create trigger g before insert on t for each row begin set statement max_statement_time = 1 for set @a = 1; end
< ok 0 0
> exec ^
: string
| create trigger g before insert on t for each row begin set statement max_statement_time = 1 for begin select 1; end; end
< error
| 1064 42000 …near 'select 1; end; end' at line 1
> exec ^
: string
| create trigger g before insert on t for each row begin set statement max_statement_time = 1 for if 1 then do 1; end if; end
< ok 0 0
> exec ^
: string
| create event e on schedule every 1 day do set statement max_statement_time = 1 for select 1
< ok 0 0
> exec ^
: string
| create event e on schedule every 1 day do set statement max_statement_time = 1 for
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create event e on schedule every 1 day do set statement max_statement_time = 1 for set @a = 1
< ok 0 0
> exec ^
: string
| create event e on schedule every 1 day do set statement max_statement_time = 1 for begin select 1; end
< error
| 1064 42000 …near 'select 1; end' at line 1
> exec ^
: string
| create event e on schedule every 1 day do set statement max_statement_time = 1 for if 1 then select 1; end if
< ok 0 0
> exec ^
: string
| /*!100200 set statement max_statement_time = 1 for select 1 */
< ok 0 0
> exec ^
: string
| explain set statement max_statement_time = 1 for select 1
< error
| 1064 42000 …near 'set statement max_statement_time = 1 for select 1' at line 1
> exec ^
: string
| prepare s from 'set statement max_statement_time = 1 for select 1'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| prepare s from 'set statement max_statement_time = 1 for select 1'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 for set statement max_statement_time = 1
< error
| 1064 42000 …near 'set statement max_statement_time = 1' at line 1
