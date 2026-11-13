-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSetValueFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
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
| 3061 42000 User variable name  is illegal
> exec ^
: string
| set @"" = 1
< error
| 3061 42000 User variable name  is illegal
> exec ^
: string
| load data infile 'x' into table t (@'a')
< error
| 1017 HY000 File not exist
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
< error
| 1064 42000 …near '@'a b' = 1' at line 1
> exec ^
: string
| set @'a'b = 1
< error
| 1064 42000 …near '@'a'b = 1' at line 1
> exec ^
: string
| set @'a'.b = 1
< error
| 1064 42000 …near '@'a'.b = 1' at line 1
> exec ^
: string
| set @'a''b' = 1
< ok 0 0
> exec ^
: string
| set @'a\'b' = 1
< error
| 1064 42000 …near '@'a\'b' = 1' at line 1
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
| 1064 42000 …near '@'a' at line 1
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
| 1064 42000 …near '@'a'' at line 1
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
< ok 0 0
> exec ^
: string
| select @'a'.b
< error
| 1064 42000 …near '@'a'.b' at line 1
> exec ^
: string
| select * from t where a = @'b'
< error
| 1046 3D000 No database selected
> exec ^
: string
| create procedure p() begin set @'a' = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin select 1 into @'a'; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| grant select on t to a@'b'
< error
| 1046 3D000 No database selected
> exec ^
: string
| create definer = a@'b' procedure p() select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| set password for 'a'@'b' = 'c'
< ok 0 0
> exec ^
: string
| rename user 'a'@'b' to 'c'@'d'
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
| 1046 3D000 No database selected
> exec ^
: string
| select 1 from t where a = @'b' and c = @'d'
< error
| 1046 3D000 No database selected
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
| 1064 42000 …near '@'a'.b from t' at line 1
> exec ^
: string
| select @'a'@'b'
< error
| 1064 42000 …near '@'b'' at line 1
> exec ^
: string
| call p(@'a')
< error
| 1046 3D000 No database selected
> exec ^
: string
| call p(@'a', @"b")
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t set a = @'b'
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t values (@'a')
< error
| 1046 3D000 No database selected
> exec ^
: string
| signal sqlstate '45000' set message_text = @'a'
< ok 0 0
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
| create definer = 'a'@'b' procedure p() select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create definer = 'a'@'b' trigger g before insert on t for each row set @a = 1
< error
| 1046 3D000 No database selected
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
< error
| 1064 42000 …near '@'a-b' = 1' at line 1
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
< error
| 1064 42000 …near '@"a b" = 1' at line 1
> exec ^
: string
| set @`a b` = 1
< error
| 1064 42000 …near '@`a b` = 1' at line 1
> exec ^
: string
| select @'a b'
< error
| 1064 42000 …near '@'a b'' at line 1
> exec ^
: string
| select @'a-b'
< error
| 1064 42000 …near '@'a-b'' at line 1
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
< error
| 1064 42000 …near '@'a b'' at line 1
> exec ^
: string
| set @'a\\b' = 1
< error
| 1064 42000 …near '@'a\\b' = 1' at line 1
> exec ^
: string
| set @'a%b' = 1
< ok 0 0
> exec ^
: string
| set @'a@b' = 1
< error
| 1064 42000 …near '@'a@b' = 1' at line 1
> exec ^
: string
| set @'a
| b' = 1
< error
| 1064 42000 …near '@'a
| b' = 1' at line 1
> exec ^
: string
| select @`a b`
< error
| 1064 42000 …near '@`a b`' at line 1
> exec ^
: string
| create procedure p() begin set @'a' = 1, @"b" = 2; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
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
< ok 0 0
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
| 1054 42S22 Unknown column 'a' in 'field_list'
> exec ^
: string
| set @a = default(a) + 1
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
> exec ^
: string
| set sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @@sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
> exec ^
: string
| set @@sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @@session.sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
> exec ^
: string
| set @@session.sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set global sql_mode = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
> exec ^
: string
| set global sql_mode = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set @a = default(a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
> exec ^
: string
| set @a = default (a)
< error
| 1054 42S22 Unknown column 'a' in 'field_list'
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
| 1054 42S22 Unknown column 'default' in 'field_list'
> exec ^
: string
| set sql_mode := default
< ok 0 0
> exec ^
: string
| set @a := default
< error
| 1064 42000 …near 'default' at line 1
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
| 1046 3D000 No database selected
> exec ^
: string
| insert into t values (default)
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t (a) values (default)
< error
| 1046 3D000 No database selected
> exec ^
: string
| update t set a = default
< error
| 1046 3D000 No database selected
> exec ^
: string
| insert into t (a) values (1) on duplicate key update a = default
< error
| 1046 3D000 No database selected
> exec ^
: string
| install component 'file://c' set a.b = on
< error
| 1064 42000 …near 'component 'file://c' set a.b = on' at line 1
