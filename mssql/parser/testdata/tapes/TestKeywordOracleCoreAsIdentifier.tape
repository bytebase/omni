-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestKeywordOracleCoreAsIdentifier default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_select (select INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_select_q ([select] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 select FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_from (from INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'from'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_from_q ([from] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 from FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_where (where INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'where'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_where_q ([where] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 where FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_insert (insert INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'insert'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_insert_q ([insert] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 insert FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_update (update INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'update'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_update_q ([update] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 update FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_delete (delete INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'delete'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_delete_q ([delete] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 delete FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_create (create INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'create'.
| 343 1 15 1 Unknown object type 'INT' used in a CREATE, DROP, or ALTER statement.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_create_q ([create] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 create FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_drop (drop INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'drop'.
| 343 1 15 1 Unknown object type 'INT' used in a CREATE, DROP, or ALTER statement.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_drop_q ([drop] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 drop FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_alter (alter INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'alter'.
| 343 1 15 1 Unknown object type 'INT' used in a CREATE, DROP, or ALTER statement.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_alter_q ([alter] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 alter FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_order (order INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'order'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_order_q ([order] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 order FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_group (group INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'group'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_group_q ([group] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 group FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_having (having INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'having'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_having_q ([having] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 having FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_join (join INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'join'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_join_q ([join] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 join FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'join'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_on (on INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'on'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_on_q ([on] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 on FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'on'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_set (set INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'set'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_set_q ([set] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 set FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_into (into INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'into'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_into_q ([into] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 into FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_values (values INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'values'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_values_q ([values] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 values FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'values'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_table (table INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'table'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE TABLE dbo.t_core_table_q ([table] INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 table FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'table'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestKeywordOracleCoreAsIdentifier]
< ok 0 -
