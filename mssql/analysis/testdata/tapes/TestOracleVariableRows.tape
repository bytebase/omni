-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestOracleVariableRows default#1
> exec
| CREATE TABLE dbo.base (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT dbo.base VALUES (1, N'alice'), (2, N'bob')
< ok 2 -
> query
| DECLARE @v NVARCHAR(20) = (SELECT secret FROM dbo.base WHERE a = 1); SELECT @v AS v
< result
< column "v" "NVARCHAR"
< row
: string
| alice
> query
| DECLARE @v NVARCHAR(20); SET @v = (SELECT secret FROM dbo.base WHERE a = 2); SELECT @v AS v
< result
< column "v" "NVARCHAR"
< row
: string
| bob
> query
| DECLARE @v NVARCHAR(20); SELECT @v = secret FROM dbo.base WHERE a = 1; SELECT @v AS v
< result
< column "v" "NVARCHAR"
< row
: string
| alice
> query
| DECLARE @v INT; SELECT @v = a FROM dbo.base
< result
> exec
| DECLARE @v INT; SELECT @v = a, secret FROM dbo.base
< error
| 141 1 15 1 A SELECT statement that assigns a value to a variable must not be combined with data-retrieval operations.
> query
| DECLARE @v INT = 1
< result
> query
| SELECT @v AS v
< error
| 137 2 15 1 Must declare the scalar variable "@v".
