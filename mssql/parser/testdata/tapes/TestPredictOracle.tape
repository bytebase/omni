-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestPredictOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT d.*, p.Score FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| INSERT INTO dbo.t (a) SELECT p.Score FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float COLLATE Latin1_General_CI_AS NOT NULL, Label nvarchar(10) NULL) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score decimal(10, 2), L nvarchar(max), u dbo.mytype) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t) WITH (Score float)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t d) WITH (Score float) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS [d]) WITH ([Score] float) AS [p]
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = 0x0102, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = N'abc', DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = CAST(@m AS VARBINARY(MAX)) + @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(DATA = dbo.t AS d, MODEL = @m) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| select * from predict(model = @m, data = dbo.t as d) with (score float) as p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM [PREDICT](MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = testdb.dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = #tmp AS d) WITH (s float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| WITH c AS (SELECT 1 AS a) SELECT * FROM PREDICT(MODEL = @m, DATA = c AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = @tv AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.fn(1) AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = (SELECT * FROM dbo.t) AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = (SELECT 1 AS a UNION SELECT 2)) WITH (s float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM dbo.t AS x JOIN PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p ON 1 = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM dbo.t AS x CROSS APPLY PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS p, PREDICT(MODEL = @m, DATA = dbo.t AS e) WITH (s float) AS q
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p PIVOT (SUM(Score) FOR a IN ([1])) AS q
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| UPDATE dbo.t SET a = 1 FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| DELETE dbo.t FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| MERGE dbo.t AS x USING PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p ON 1 = 1 WHEN MATCHED THEN DELETE;
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT predict FROM dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM predict AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM predict(1)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM predict(model)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM dbo.t CROSS APPLY predict(a) AS f
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| PREDICT (MODEL = @m, DATA = dbo.t AS d) WITH (Score float)
|
< error
| 102 1 15 2 Incorrect syntax near 'PREDICT'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| PREDICT (MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near 'PREDICT'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float)
|
< error
| 102 1 15 2 Incorrect syntax near '='.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) AS p
|
< error
| 102 38 15 2 Incorrect syntax near '='.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d)
|
< error
| 102 38 15 2 Incorrect syntax near '='.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH () AS p
|
< error
| 102 1 15 2 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float,) AS p
|
< error
| 102 1 15 2 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (, Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score) AS p
|
< error
| 102 1 15 2 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (select int) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float NOT NULL COLLATE Latin1_General_CI_AS) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float NULL NOT NULL) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'NOT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float PRIMARY KEY) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'PRIMARY'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float DEFAULT 1) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'DEFAULT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float '$.x') AS p
|
< error
| 102 1 15 2 Incorrect syntax near '$.x'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float AS JSON) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float COLLATE [Latin1_General_CI_AS]) AS p
|
< error
| 102 1 15 2 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float SPARSE) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(DATA = dbo.t AS d) WITH (Score float) AS p
|
< error
| 39036 6 16 2 The function PREDICT expects parameter 'MODEL' which was not supplied.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m) WITH (Score float) AS p
|
< error
| 39036 7 16 2 The function PREDICT expects parameter 'DATA' which was not supplied.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< error
| 39035 1 16 2 The function PREDICT has too many arguments supplied.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, DATA = dbo.t AS e) WITH (s float) AS p
|
< error
| 39035 2 16 2 The function PREDICT has too many arguments supplied.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d,) WITH (Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near ')'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m,, DATA = dbo.t AS d) WITH (Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near ','.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL @m, DATA = dbo.t AS d) WITH (Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '@m'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m DATA = dbo.t AS d) WITH (s float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '='.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, FOO = 1) WITH (Score float) AS p
|
< error
| 102 5 16 2 Incorrect syntax near 'FOO'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM predict(model = 1)
|
< error
| 102 38 15 2 Incorrect syntax near '='.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = ONNX) WITH (Score float) AS p
|
< error
| 102 5 16 2 Incorrect syntax near 'RUNTIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = FOO) WITH (Score float) AS p
|
< error
| 102 5 16 2 Incorrect syntax near 'RUNTIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = 'ONNX') WITH (Score float) AS p
|
< error
| 102 5 16 2 Incorrect syntax near 'RUNTIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, 'RUNTIME' = ONNX) WITH (Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '='.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = (SELECT a FROM dbo.t), DATA = dbo.t AS d) WITH (Score float) AS p
|
< error
| 102 3 16 2 Incorrect syntax near '='.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d WITH (NOLOCK)) WITH (Score float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'WITH'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t d (NOLOCK)) WITH (s float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '('.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.fn(1) AS d(x)) WITH (s float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '('.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = (SELECT 1 AS a) AS d(x)) WITH (s float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '('.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t FOR SYSTEM_TIME AS OF '2020-01-01' AS d) WITH (s float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near 'FOR'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d TABLESAMPLE (10 PERCENT)) WITH (s float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'TABLESAMPLE'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS a JOIN dbo.t AS b ON 1 = 1) WITH (s float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'JOIN'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = (dbo.t AS a JOIN dbo.t AS b ON 1 = 1)) WITH (s float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'AS'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = (VALUES (1)) AS d(a)) WITH (s float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'VALUES'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = OPENJSON(@j) AS d) WITH (s float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near 'OPENJSON'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = CONTAINSTABLE(dbo.t, a, 'x') AS d) WITH (s float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'CONTAINSTABLE'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = fn(1) AS d) WITH (s float) AS p
|
< error
| 195 10 15 2 'fn' is not a recognized built-in function name.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = PREDICT(MODEL = @m, DATA = dbo.t AS e) WITH (s float) AS d) WITH (Score float) AS p
|
< error
| 102 1 15 2 Incorrect syntax near '='.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS select) WITH (s float) AS p
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'select'.
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = 1) WITH (s float) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p (x)
|
< error
| 102 13 16 2 Incorrect syntax near 'PREDICT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p WITH (NOLOCK)
|
< error
| 319 1 15 2 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS p (NOLOCK)
|
< error
| 102 13 16 2 Incorrect syntax near 'PREDICT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p FOR SYSTEM_TIME AS OF '2020-01-01'
|
< error
| 102 1 15 2 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM (PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p)
|
< error
| 102 1 15 2 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);
| SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS select
|
< error
| 156 1 15 2 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPredictOracle]
< ok 0 -
