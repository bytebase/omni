-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestTableSampleOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS [s] TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t AS s TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) WHERE s.a = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) JOIN t1 ON s.a = t1.a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) JOIN t1 ON t.a = t1.a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) WITH (NOLOCK) JOIN t1 ON t.a = t1.a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT), t1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) CROSS APPLY (SELECT 1 AS x) ca
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) PIVOT (MAX(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE s SET a = 1 FROM t AS s TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DELETE s FROM t AS s TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 AS s TABLESAMPLE (10 PERCENT) ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) s
|
< error
| 102 1 15 1 Incorrect syntax near 's'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) AS s
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) WITH (NOLOCK) s
|
< error
| 102 1 15 1 Incorrect syntax near 's'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) s2
|
< error
| 102 1 15 1 Incorrect syntax near 's2'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE (10 PERCENT) AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s WITH (NOLOCK) TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WITH (NOLOCK) TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) WITH (NOLOCK) WITH (NOLOCK)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) HOLDLOCK
|
< error
| 1018 1 15 1 Incorrect syntax near 'HOLDLOCK'. If this is intended as a part of a table hint, A WITH keyword and parenthesis are now required. See SQL Server Books Online for proper syntax.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE (10 PERCENT) TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (SELECT 1 AS a) d TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t [TABLESAMPLE] (10 PERCENT)
|
< error
| 102 1 15 1 Incorrect syntax near '10'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS window TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS
|
< error
| 102 1 15 1 Incorrect syntax near 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS WHERE 1 = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WHERE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS JOIN t1 ON t.a = t1.a
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'JOIN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (SELECT 1 AS a) AS WHERE 1 = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WHERE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON('[]') AS WHERE 1 = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WHERE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 AS ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE SYSTEM (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE system (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t TABLESAMPLE SYSTEM (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 ROWS)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 percent)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10.5 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 + 1 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (-10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE ((10) PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @n INT = 10; SELECT * FROM t AS s TABLESAMPLE (@n PERCENT) REPEATABLE (@n)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) repeatable (5)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5 + 1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE SYSTEM (10 PERCENT) REPEATABLE (5) WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s TABLESAMPLE
|
< error
| 102 1 15 1 Incorrect syntax near 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE ()
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'PERCENT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE 10 PERCENT
|
< error
| 102 1 15 1 Incorrect syntax near '10'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 ROW)
|
< error
| 102 1 15 1 Incorrect syntax near 'ROW'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT ROWS)
|
< error
| 102 1 15 1 Incorrect syntax near 'ROWS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE [SYSTEM] (10 PERCENT)
|
< error
| 102 1 15 1 Incorrect syntax near 'SYSTEM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE SYSTEM SYSTEM (10 PERCENT)
|
< error
| 102 1 15 1 Incorrect syntax near 'SYSTEM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE SYSTEM
|
< error
| 102 1 15 1 Incorrect syntax near 'SYSTEM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE BERNOULLI (10 PERCENT)
|
< error
| 102 1 15 1 Incorrect syntax near 'BERNOULLI'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE
|
< error
| 102 1 15 1 Incorrect syntax near 'REPEATABLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE ()
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE 5
|
< error
| 102 1 15 1 Incorrect syntax near '5'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) [REPEATABLE] (5)
|
< error
| 102 1 15 1 Incorrect syntax near 'REPEATABLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5) REPEATABLE (6)
|
< error
| 102 1 15 1 Incorrect syntax near 'REPEATABLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s TABLESAMPLE (10 PERCENT) REPEATABLE (5, 6)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTableSampleOracle]
< ok 0 -
