-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSetValueFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set @ 'a' = 1
< error
| 1064 42000 …near ''a' = 1' at line 1
> exec ^
: string
| set @ "a" = 1
< error
| 1064 42000 …near '"a" = 1' at line 1
> exec ^
: string
| select @''
< ok 0 0
> exec ^
: string
| set @'' = 1
< error
| 3061 42000 User variable name '' is illegal
> exec ^
: string
| set @'' = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @"" = 1
< error
| 3061 42000 User variable name '' is illegal
> exec ^
: string
| set @"" = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @'a' = 1
< ok 0 0
> exec ^
: string
| set @"a" = 1
< ok 0 0
> exec ^
: string
| set @`a` = 1
< ok 0 0
> exec ^
: string
| set @'a' := 1
< ok 0 0
> exec ^
: string
| set @'a b' = 1
< ok 0 0
> exec ^
: string
| set @'a'b = 1
< error
| 1064 42000 …near 'b = 1' at line 1
> exec ^
: string
| set @'a'.b = 1
< error
| 1064 42000 …near '.b = 1' at line 1
> exec ^
: string
| set @'a''b' = 1
< ok 0 0
> exec ^
: string
| set @'a\'b' = 1
< ok 0 0
> exec ^
: string
| set @"a""b" = 1
< ok 0 0
> exec ^
: string
| set @'a' 'b' = 1
< error
| 1064 42000 …near ''b' = 1' at line 1
> exec ^
: string
| set @'a
< error
| 1064 42000 …near ''a' at line 1
> exec ^
: string
| set @'a' = 1, @'b' = 2
< ok 0 0
> exec ^
: string
| set @a = @'b'
< ok 0 0
> exec ^
: string
| set @a = 1, @'b' = 2
< ok 0 0
> exec ^
: string
| set global @'a' = 1
< error
| 1064 42000 …near '@'a' = 1' at line 1
> exec ^
: string
| set @'a'
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select @'a'
< ok 0 0
> exec ^
: string
| select @"a"
< ok 0 0
> exec ^
: string
| select @`a`
< ok 0 0
> exec ^
: string
| select @'a', @"b"
< ok 0 0
> exec ^
: string
| select @'a' := 1
< ok 0 0
> exec ^
: string
| select @'a' + 1
< ok 0 0
> exec ^
: string
| select 1 into @'a'
< ok 0 0
> exec ^
: string
| select 1 into @"a"
< ok 0 0
> exec ^
: string
| select 1 into @'a', @'b'
< error
| 1222 21000 The used SELECT statements have a different number of columns
> exec ^
: string
| select 1 into @'a', @'b'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select @'a'.b
< error
| 1064 42000 …near '.b' at line 1
> exec ^
: string
| prepare s from @'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| prepare s from @'a'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| execute s using @'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| execute s using @'a'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select * from t where a = @'b'
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select * from t where a = @'b'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set @'a' = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set @'a' = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin select 1 into @'a'; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin select 1 into @'a'; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create user 'a'@'b'
< ok 0 0
> exec ^
: string
| create user a@b
< ok 0 0
> exec ^
: string
| create user 'a'@b
< ok 0 0
> exec ^
: string
| create user a@'b'
< ok 0 0
> exec ^
: string
| create user `a`@`b`
< ok 0 0
> exec ^
: string
| create user "a"@"b"
< ok 0 0
> exec ^
: string
| create user 'a' @'b'
< ok 0 0
> exec ^
: string
| create user 'a'@ 'b'
< error
| 1064 42000 …near ''b'' at line 1
> exec ^
: string
| drop user 'a'@'b'
< ok 0 0
> exec ^
: string
| grant select on t to 'a'@'b'
< ok 0 0
> exec ^
: string
| grant select on t to a@'b'
< ok 0 0
> exec ^
: string
| create definer = 'a'@'b' view v as select 1
< ok 0 0
> exec ^
: string
| create definer = a@'b' procedure p() select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create definer = a@'b' procedure p() select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set password for 'a'@'b' = 'c'
< error
| 1044 42000 Access denied for user 'omni_test'@'%' to database 'mysql'
> exec ^
: string
| set password for 'a'@'b' = 'c'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set default role none to 'a'@'b'
< ok 0 0
> exec ^
: string
| rename user 'a'@'b' to 'c'@'d'
< ok 0 0
> exec ^
: string
| show grants for 'a'@'b'
< ok 0 0
> exec ^
: string
| select 'a'@'b'
< error
| 1064 42000 …near '@'b'' at line 1
> exec ^
: string
| select 'a' @'b'
< error
| 1064 42000 …near '@'b'' at line 1
> exec ^
: string
| select 1 @'b'
< error
| 1064 42000 …near '@'b'' at line 1
> exec ^
: string
| select 1 as @'b'
< error
| 1064 42000 …near '@'b'' at line 1
> exec ^
: string
| select @'b' as a
< ok 0 0
> exec ^
: string
| select @'a' from t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select @'a' from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 from t where a = @'b' and c = @'d'
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| select 1 from t where a = @'b' and c = @'d'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select @"a" := 1
< ok 0 0
> exec ^
: string
| select (@'a' := 1)
< ok 0 0
> exec ^
: string
| select @'a' = 1
< ok 0 0
> exec ^
: string
| select @'a'.b from t
< error
| 1064 42000 …near '.b from t' at line 1
> exec ^
: string
| select @'a'@'b'
< error
| 1064 42000 …near '@'b'' at line 1
> exec ^
: string
| execute s using @'a', @"b"
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| execute s using @'a', @"b"
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile 'x' into table t (@'a')
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile 'x' into table t (@'a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t set a = @'b'
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| update t set a = @'b'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into t values (@'a')
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t values (@'a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| signal sqlstate '45000' set message_text = @'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| signal sqlstate '45000' set message_text = @'a'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, @'b' := 2
< ok 0 0
> exec ^
: string
| set @"a" = @'b'
< ok 0 0
> exec ^
: string
| /*!50003 set @'a' = 1 */
< ok 0 0
> exec ^
: string
| create definer = a@b view v as select 1
< ok 0 0
> exec ^
: string
| create definer = 'a'@'b' procedure p() select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create definer = 'a'@'b' procedure p() select 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create definer = 'a'@'b' trigger g before insert on t for each row set @a = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create definer = 'a'@'b' trigger g before insert on t for each row set @a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @'a_b' = 1
< ok 0 0
> exec ^
: string
| set @'a.b' = 1
< ok 0 0
> exec ^
: string
| set @'a$b' = 1
< ok 0 0
> exec ^
: string
| set @'a-b' = 1
< ok 0 0
> exec ^
: string
| set @'1a' = 1
< ok 0 0
> exec ^
: string
| set @'é' = 1
< ok 0 0
> exec ^
: string
| set @"a b" = 1
< ok 0 0
> exec ^
: string
| set @`a b` = 1
< ok 0 0
> exec ^
: string
| select @'a b'
< ok 0 0
> exec ^
: string
| select @'a-b'
< ok 0 0
> exec ^
: string
| select @'a_b'
< ok 0 0
> exec ^
: string
| select @'a.b'
< ok 0 0
> exec ^
: string
| select 1 into @'a b'
< ok 0 0
> exec ^
: string
| set @'a\\b' = 1
< ok 0 0
> exec ^
: string
| set @'a%b' = 1
< ok 0 0
> exec ^
: string
| set @'a@b' = 1
< ok 0 0
> exec ^
: string
| set @'a
| b' = 1
< ok 0 0
> exec ^
: string
| select @`a b`
< ok 0 0
> exec ^
: string
| create procedure p() begin set @'a' = 1, @"b" = 2; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set @'a' = 1, @"b" = 2; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = 1, names default
< ok 0 0
> exec ^
: string
| set names default, @a = 1
< ok 0 0
> exec ^
: string
| set @a = 1, session names default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set names default collate utf8mb4_bin
< error
| 1064 42000 …near 'collate utf8mb4_bin' at line 1
> exec ^
: string
| set @a = 1, charset default
< ok 0 0
> exec ^
: string
| set charset default, @a = 1
< ok 0 0
> exec ^
: string
| set @a = 1, session charset default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set charset default collate utf8mb4_bin
< error
| 1064 42000 …near 'collate utf8mb4_bin' at line 1
> exec ^
: string
| set @a = 1, character set default
< ok 0 0
> exec ^
: string
| set character set default, @a = 1
< ok 0 0
> exec ^
: string
| set @a = 1, session character set default
< error
| 1064 42000 …near 'character set default' at line 1
> exec ^
: string
| set character set default collate utf8mb4_bin
< error
| 1064 42000 …near 'collate utf8mb4_bin' at line 1
> exec ^
: string
| set @a = 1, char set default
< ok 0 0
> exec ^
: string
| set char set default, @a = 1
< ok 0 0
> exec ^
: string
| set @a = 1, session char set default
< error
| 1064 42000 …near 'char set default' at line 1
> exec ^
: string
| set char set default collate utf8mb4_bin
< error
| 1064 42000 …near 'collate utf8mb4_bin' at line 1
> exec ^
: string
| set names utf8mb4 collate default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set sql_mode = default(a) + 1
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set sql_mode = default(a) + 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = default(a) + 1
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @a = default(a) + 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set sql_mode = default
< ok 0 0
> exec ^
: string
| set sql_mode = default + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set sql_mode = 1 + default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set sql_mode = default(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set sql_mode = default (a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set sql_mode = (default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| set sql_mode = default default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set sql_mode = default, @b = 1
< ok 0 0
> exec ^
: string
| set sql_mode = `default`
< ok 0 0
> exec ^
: string
| set @@sql_mode = default
< ok 0 0
> exec ^
: string
| set @@sql_mode = default + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set @@sql_mode = 1 + default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @@sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @@sql_mode = default(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @@sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @@sql_mode = default (a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @@sql_mode = (default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| set @@sql_mode = default default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @@sql_mode = default, @b = 1
< ok 0 0
> exec ^
: string
| set @@sql_mode = `default`
< ok 0 0
> exec ^
: string
| set @@session.sql_mode = default
< ok 0 0
> exec ^
: string
| set @@session.sql_mode = default + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set @@session.sql_mode = 1 + default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @@session.sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @@session.sql_mode = default(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @@session.sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @@session.sql_mode = default (a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @@session.sql_mode = (default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| set @@session.sql_mode = default default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @@session.sql_mode = default, @b = 1
< ok 0 0
> exec ^
: string
| set @@session.sql_mode = `default`
< ok 0 0
> exec ^
: string
| set global sql_mode = default
< ok 0 0
> exec ^
: string
| set global sql_mode = default + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set global sql_mode = 1 + default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set global sql_mode = default(a)
< ok 0 0
> exec ^
: string
| set global sql_mode = default (a)
< ok 0 0
> exec ^
: string
| set global sql_mode = (default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| set global sql_mode = default default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set global sql_mode = default, @b = 1
< ok 0 0
> exec ^
: string
| set global sql_mode = `default`
< ok 0 0
> exec ^
: string
| set @a = default + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| set @a = 1 + default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @a = default(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| set @a = default (a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set @a = (default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| set @a = default default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @a = default, @b = 1
< error
| 1064 42000 …near ', @b = 1' at line 1
> exec ^
: string
| set @a = `default`
< error
| 1054 42S22 Unknown column 'default' in 'field list'
> exec ^
: string
| set @a = `default`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set sql_mode := default
< ok 0 0
> exec ^
: string
| set @a := default
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @@sql_mode := default
< ok 0 0
> exec ^
: string
| set @a = 1, sql_mode = default
< ok 0 0
> exec ^
: string
| set @a = 1, global sql_mode = default
< ok 0 0
> exec ^
: string
| set sql_mode = default, @a = 1
< ok 0 0
> exec ^
: string
| set @a = (select default)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| set @a = (select default(a) from t)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| set @a = (select default(a) from t)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into t values (default)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t values (default)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into t (a) values (default)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t (a) values (default)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| update t set a = default
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| update t set a = default
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| insert into t (a) values (1) on duplicate key update a = default
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| insert into t (a) values (1) on duplicate key update a = default
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| install component 'file://c' set a.b = on
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| install component 'file://c' set a.b = on
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
