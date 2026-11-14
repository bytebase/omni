-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestWithClauseNestedOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c) d
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c) s ON dbo.t.a = s.a WHEN NOT MATCHED THEN INSERT (a) VALUES (s.a);
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (WITH c AS (SELECT a FROM dbo.t) SELECT MAX(a) FROM c) AS m
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @x INT; SET @x = (WITH c AS (SELECT a FROM dbo.t) SELECT MAX(a) FROM c)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM dbo.t WHERE EXISTS (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM dbo.t WHERE a IN (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM dbo.t WHERE a = ANY (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (WITH d AS (SELECT a FROM dbo.t) SELECT a FROM d) SELECT a FROM c
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| INSERT INTO dbo.t (a) WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM dbo.t UNION WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c UNION ALL WITH d AS (SELECT a FROM dbo.t) SELECT a FROM d
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM dbo.t EXCEPT WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE VIEW dbo.v AS WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| ALTER VIEW dbo.v AS WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN (WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE cur CURSOR FOR WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (SELECT a FROM dbo.t) d
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT a FROM dbo.t UNION SELECT a FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| INSERT INTO dbo.t (a) SELECT a FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN (SELECT a FROM dbo.t)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| CREATE FUNCTION dbo.f() RETURNS TABLE AS RETURN SELECT a FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestWithClauseNestedOracle]
< ok 0 -
