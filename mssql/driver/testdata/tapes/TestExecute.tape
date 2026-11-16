-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestExecute default#1
> exec
| CREATE TABLE e_items (id INT PRIMARY KEY, note NVARCHAR(64));
< ok 0 -
> begin
< ok
> exec
| INSERT INTO e_items VALUES (1, N'a'), (2, N'b');
< ok 2 -
> exec
| UPDATE e_items SET note = N'z' WHERE id <= 1;
< ok 1 -
> commit
< ok
> begin
< ok
> exec
| INSERT INTO e_items VALUES (10, N'x');
< ok 1 -
> exec
| INSERT INTO e_items VALUES (1, N'duplicate');
< error
| 2627 1 14 1 Violation of PRIMARY KEY constraint 'PK__e_items__standin1'. Cannot insert duplicate key in object 'dbo.e_items'. The duplicate key value is (1).
> rollback
< ok
> query
| SELECT COUNT(*) FROM e_items WHERE id = 10
< result
< column "" "INT"
< row
: int64 0
> exec
| INSERT INTO e_items VALUES (20, N'kept');
< ok 1 -
> exec
| INSERT INTO e_items VALUES (1, N'duplicate');
< error
| 2627 1 14 1 Violation of PRIMARY KEY constraint 'PK__e_items__standin1'. Cannot insert duplicate key in object 'dbo.e_items'. The duplicate key value is (1).
> query
| SELECT COUNT(*) FROM e_items WHERE id = 20
< result
< column "" "INT"
< row
: int64 1
> exec
| UPDATE e_items SET note = N'r' WHERE id = 2;
< ok 1 -
> exec ^
< ok 1 -
> exec ^
< ok 1 -
> exec
| CREATE DATABASE [omni_created_testexecute]
< ok 0 -
> query
| SELECT SERVERPROPERTY('productversion'), @@VERSION
< result
< column "" "SQL_VARIANT"
< column "" "NVARCHAR"
< row
: string
| 16.0.4305.5
: string
| Microsoft SQL Server 2022 (RTM-CU27-GDR) (KB5126934) - 16.0.4305.5 (X64) 
| 	Oct  2 2026 17:28:11 
| 	Copyright (C) 2022 Microsoft Corporation
| 	Developer Edition (64-bit) on Linux (Ubuntu 22.04.5 LTS) <X64>
> query
| SELECT name, collation_name FROM sys.databases WHERE name NOT IN ('master', 'model', 'msdb', 'tempdb', 'rdscore') ORDER BY name
< result
< column "name" "NVARCHAR"
< column "collation_name" "NVARCHAR"
< row
: string
| omni_created_testexecute
: string
| SQL_Latin1_General_CP1_CI_AS
< row
: string
| omni_TestExecute
: string
| SQL_Latin1_General_CP1_CI_AS
> exec
| DROP DATABASE IF EXISTS [omni_created_testexecute]
< ok 0 -
