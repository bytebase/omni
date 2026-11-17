-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestOracleTempTableRows default#1
> exec
| CREATE TABLE dbo.base (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT dbo.base VALUES (1, N'alice'), (2, N'bob')
< ok 2 -
> exec
| CREATE PROCEDURE dbo.p AS SELECT a, secret FROM dbo.base
< ok 0 -
> exec
| SELECT a, secret INTO #into FROM dbo.base
< ok 2 -
> query
| SELECT a, secret FROM #into ORDER BY a
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
> exec
| CREATE TABLE #values (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT #values VALUES ((SELECT MIN(a) FROM dbo.base), (SELECT secret FROM dbo.base WHERE a = 1))
< ok 1 -
> query
| SELECT a, secret FROM #values
< result
< column "a" "INT"
< column "secret" "NVARCHAR"
< row
: int64 1
: string
| alice
> exec
| CREATE TABLE #select (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT #select SELECT a, secret FROM dbo.base
< ok 2 -
> query
| SELECT a, secret FROM #select ORDER BY a
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
> exec
| CREATE TABLE #update (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT #update (a) VALUES (1)
< ok 1 -
> exec
| UPDATE #update SET secret = b.secret FROM dbo.base AS b WHERE b.a = #update.a
< ok 1 -
> query
| SELECT a, secret FROM #update
< result
< column "a" "INT"
< column "secret" "NVARCHAR"
< row
: int64 1
: string
| alice
> exec
| CREATE TABLE #merge (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| MERGE #merge AS tgt USING (SELECT a, secret FROM dbo.base) AS src ON tgt.a = src.a WHEN NOT MATCHED THEN INSERT (a, secret) VALUES (src.a, src.secret);
< ok 2 -
> query
| SELECT a, secret FROM #merge ORDER BY a
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
> exec
| CREATE TABLE #exec (a INT, secret NVARCHAR(20))
< ok 0 -
> exec
| INSERT #exec EXEC dbo.p
< ok 2 -
> query
| SELECT a, secret FROM #exec ORDER BY a
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
| SELECT a, secret INTO #across FROM dbo.base
< result
> query
| SELECT COUNT(*) AS n FROM #across
< result
< column "n" "INT"
< row
: int64 2
> query
| SELECT a, secret FROM #across ORDER BY a
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
