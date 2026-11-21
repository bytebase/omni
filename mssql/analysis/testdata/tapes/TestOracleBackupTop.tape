-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestOracleBackupTop default#1
> exec
| CREATE TABLE dbo.test (c1 INT, c2 INT)
< ok 0 -
> exec
| INSERT dbo.test VALUES (1, 2), (2, 2), (3, 2), (4, 2)
< ok 4 -
> exec
| SELECT * INTO [omni_TestOracleBackupTop].[dbo].[rollback_test_omni_TestOracleBackupTop] FROM (
|   SELECT TOP (3) [omni_TestOracleBackupTop].[dbo].[test].* FROM test WHERE c2 = 2) AS backup_table;
< ok 3 -
> query
| SELECT COUNT(*) FROM dbo.[rollback_test_omni_TestOracleBackupTop]
< result
< column "" "INT"
< row
: int64 3
