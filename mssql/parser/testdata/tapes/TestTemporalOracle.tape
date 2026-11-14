-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestTemporalOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' TO '2021-01-01'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' AND '2021-01-01'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('2020-01-01', '2021-01-01')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t for system_time all
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR system_time contained in ('a', 'b')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR [SYSTEM_TIME] ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'SYSTEM_TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME [ALL]
|
< error
| 102 1 15 1 Incorrect syntax near 'ALL'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME [CONTAINED] IN ('a', 'b')
|
< error
| 102 1 15 1 Incorrect syntax near 'CONTAINED'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS '2020-01-01'
|
< error
| 102 1 15 1 Incorrect syntax near '2020-01-01'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME OF '2020-01-01'
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'OF'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF
|
< error
| 102 1 15 1 Incorrect syntax near 'OF'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01'
|
< error
| 102 1 15 1 Incorrect syntax near '2020-01-01'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' AND '2021-01-01'
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AND'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' TO '2021-01-01'
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TO'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME CONTAINED ('2020-01-01', '2021-01-01')
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('2020-01-01')
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('a', 'b', 'c')
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('a', 'b',)
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL '2020'
|
< error
| 102 1 15 1 Incorrect syntax near '2020'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF N'2020-01-01'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME AS OF @d
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF @@DBTS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF 1.5
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF 1e3
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF 0x01
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF $5
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF NULL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME FROM @d TO @e
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME FROM 1 TO 2
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME BETWEEN NULL AND @d
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN (@d, @e)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN (1, NULL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF GETDATE()
|
< error
| 102 1 15 1 Incorrect syntax near 'GETDATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF SYSUTCDATETIME()
|
< error
| 102 1 15 1 Incorrect syntax near 'SYSUTCDATETIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF CURRENT_TIMESTAMP
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'CURRENT_TIMESTAMP'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF CAST('2020-01-01' AS DATETIME2)
|
< error
| 102 1 15 1 Incorrect syntax near 'CAST'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF '2020' + '-01-01'
|
< error
| 102 1 15 1 Incorrect syntax near '+'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF -1
|
< error
| 102 1 15 1 Incorrect syntax near '-'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF +1
|
< error
| 102 1 15 1 Incorrect syntax near '+'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME AS OF (@d)
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF DEFAULT
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'DEFAULT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF 'a' COLLATE Latin1_General_CI_AS
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME AS OF @d.x
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL AS s
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL s
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); SELECT * FROM t FOR SYSTEM_TIME FROM @d TO @e AS s
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' AS s TABLESAMPLE (10 PERCENT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL AS s WITH (NOLOCK)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.t FOR SYSTEM_TIME ALL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM srv.db.dbo.t FOR SYSTEM_TIME ALL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL JOIN t1 FOR SYSTEM_TIME AS OF '2020-01-01' ON t.a = t1.a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL CROSS APPLY (SELECT 1 AS x) ca
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL PIVOT (MAX(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL, t1 FOR SYSTEM_TIME ALL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL WHERE 1 = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL FOR XML AUTO
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' FOR XML AUTO
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL ORDER BY 1 FOR JSON PATH
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL UNION SELECT * FROM t1 FOR SYSTEM_TIME ALL
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| WITH c AS (SELECT * FROM t FOR SYSTEM_TIME ALL) SELECT * FROM c
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (SELECT COUNT(*) FROM t FOR SYSTEM_TIME ALL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WHERE a IN (SELECT a FROM t1 FOR SYSTEM_TIME ALL)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE s SET a = 1 FROM t FOR SYSTEM_TIME ALL AS s
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DELETE s FROM t FOR SYSTEM_TIME ALL AS s
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t USING t1 FOR SYSTEM_TIME ALL AS s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR XML AUTO
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR BROWSE
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t AS s FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t s FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t WITH (NOLOCK) FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t FOR SYSTEM_TIME ALL FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (SELECT 1 AS a) d FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM (SELECT 1 AS a) FOR SYSTEM_TIME ALL d
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @tv TABLE (a INT); SELECT * FROM @tv FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM OPENJSON(N'[]') FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM dbo.f() FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE t FOR SYSTEM_TIME ALL SET a = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DELETE FROM t FOR SYSTEM_TIME ALL WHERE a = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| INSERT INTO t FOR SYSTEM_TIME ALL (a) VALUES (1)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| MERGE INTO dbo.t FOR SYSTEM_TIME ALL AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestTemporalOracle]
< ok 0 -
