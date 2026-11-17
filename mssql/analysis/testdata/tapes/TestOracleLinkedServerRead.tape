-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestOracleLinkedServerRead default#1
> exec
| CREATE TABLE dbo.base (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT dbo.base VALUES (1, N'alice'), (2, N'bob')
< ok 2 -
> exec
| SELECT a FROM srv.omni_TestOracleLinkedServerRead.dbo.base
< error
| 7202 2 11 1 Could not find server 'srv' in sys.servers. Verify that the correct server name was specified. If necessary, execute the stored procedure sp_addlinkedserver to add the server to sys.servers.
> exec
| SELECT * FROM OPENQUERY(srv, 'SELECT a FROM dbo.base')
< error
| 7202 2 11 1 Could not find server 'srv' in sys.servers. Verify that the correct server name was specified. If necessary, execute the stored procedure sp_addlinkedserver to add the server to sys.servers.
> exec
| EXEC ('SELECT 1') AT srv
< error
| 7202 2 11 1 Could not find server 'srv' in sys.servers. Verify that the correct server name was specified. If necessary, execute the stored procedure sp_addlinkedserver to add the server to sys.servers.
> exec
| EXEC sp_addlinkedserver @server = 'loop', @srvproduct = '', @provider = 'MSOLEDBSQL', @datasrc = '127.0.0.1,1433'
< ok 6 -
> exec
| EXEC sp_addlinkedsrvlogin @rmtsrvname = 'loop', @useself = 'TRUE'
< ok 5 -
> query
| SELECT a, secret FROM dbo.base ORDER BY a
< result
< column "a" "INT"
< column "secret" "NVARCHAR"
< row
: int64 1
: string
| alice
< row
: int64 2
: string
| bob
> query
| SELECT a, secret FROM loop.omni_TestOracleLinkedServerRead.dbo.base ORDER BY a
< result
< column "a" "INT"
< column "secret" "NVARCHAR"
< row
: int64 1
: string
| alice
< row
: int64 2
: string
| bob
> query
| SELECT OBJECT_ID('loop.omni_TestOracleLinkedServerRead.dbo.base'), OBJECT_ID('dbo.base')
< result
< column "" "INT"
< column "" "INT"
< row
: null
: int64 901578250
> query
| SELECT o.bp FROM OPENQUERY(loop, 'SELECT a, a + 1 AS bp FROM omni_TestOracleLinkedServerRead.dbo.base') AS o ORDER BY o.bp
< result
< column "bp" "INT"
< row
: int64 2
< row
: int64 3
> exec
| EXEC ('SELECT 1') AT loop
< error
| 7411 1 16 1 Server 'loop' is not configured for RPC.
> exec
| SELECT a, secret INTO #lt FROM loop.omni_TestOracleLinkedServerRead.dbo.base
< ok 2 -
> query
| SELECT a, secret FROM #lt ORDER BY a
< result
< column "a" "INT"
< column "secret" "NVARCHAR"
< row
: int64 1
: string
| alice
< row
: int64 2
: string
| bob
