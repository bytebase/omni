-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestSyncInstance default#1
> exec
| CREATE DATABASE [omni_instance_extra]
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
| omni_instance_extra
: string
| SQL_Latin1_General_CP1_CI_AS
< row
: string
| omni_TestSyncInstance
: string
| SQL_Latin1_General_CP1_CI_AS
