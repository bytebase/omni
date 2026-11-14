-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestMergeHeaderOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE TOP (1) INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t using t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t [USING] USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t window USING t1 ON window.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO window USING t1 ON window.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t WITH (HOLDLOCK) AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @tv TABLE (a INT); MERGE INTO @tv AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @tv TABLE (a INT); MERGE INTO @tv x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @tv TABLE (a INT); MERGE INTO @tv USING t1 ON a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 AS s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 s (NOLOCK) ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING (SELECT 1 AS a) s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING (SELECT 1 AS a) AS s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING (VALUES (1)) v(a) ON dbo.t.a = v.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 JOIN t2 ON t1.a = t2.a ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE WHEN NOT MATCHED THEN INSERT (a) VALUES (t1.a) WHEN NOT MATCHED BY SOURCE THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE OUTPUT $action OPTION (RECOMPILE);
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t AS USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 102 1 15 1 Incorrect syntax near 'USING'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t AS USING USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 102 1 15 1 Incorrect syntax near 'USING'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t AS x WITH (HOLDLOCK) USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t x t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 102 1 15 1 Incorrect syntax near 't1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING ON dbo.t.a = 1 WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WHEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WHEN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a;
|
< error
| 102 1 15 1 Incorrect syntax near ';'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 a b ON dbo.t.a = a.a WHEN MATCHED THEN DELETE;
|
< error
| 102 1 15 1 Incorrect syntax near 'b'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 AS a AS b ON dbo.t.a = a.a WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING (SELECT 1 AS a) s x ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< error
| 102 1 15 1 Incorrect syntax near 'x'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestMergeHeaderOracle]
< ok 0 -
