-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestWithClauseDMLOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) INSERT dbo.t (a) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) INSERT TOP (1) INTO dbo.t (a) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) UPDATE dbo.t SET a = c.a FROM dbo.t JOIN c ON dbo.t.a = c.a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) UPDATE c SET a = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) DELETE FROM dbo.t WHERE a IN (SELECT a FROM c)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) DELETE FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) MERGE INTO dbo.t USING c ON dbo.t.a = c.a WHEN NOT MATCHED THEN INSERT (a) VALUES (c.a);
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c1 AS (SELECT a FROM dbo.t), c2 (a) AS (SELECT a FROM c1) INSERT INTO dbo.t (a) SELECT a FROM c2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ;WITH c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| BEGIN WITH c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c; END
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t)
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t);
|
< error
| 102 1 15 1 Incorrect syntax near ';'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) CREATE TABLE dbo.u (a INT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'CREATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) EXEC sp_who
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'EXEC'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t), INSERT INTO dbo.t (a) SELECT a FROM c
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'INSERT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH INSERT INTO dbo.t (a) VALUES (1)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'INSERT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH SELECT 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH XMLNAMESPACES ('http://x' AS ns) SELECT a FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH XMLNAMESPACES ('http://x' AS ns) INSERT INTO dbo.t (a) SELECT a FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 WHERE EXISTS (WITH c AS (SELECT a FROM dbo.t) DELETE FROM dbo.t)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH XMLNAMESPACES ('http://x' AS ns), c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH XMLNAMESPACES ('http://x' AS ns), INSERT INTO dbo.t (a) SELECT a FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'INSERT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH XMLNAMESPACES ('http://x' AS ns), SELECT a FROM dbo.t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseDMLOracle]
< ok 0 -
