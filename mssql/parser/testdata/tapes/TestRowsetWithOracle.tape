-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestRowsetWithOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT, b INT) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT '$.a') AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT N'$.a') AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT '') AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' AS JSON) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) AS JSON) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' as json) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT '$.a' AS JSON) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE Latin1_General_CI_AS) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE Latin1_General_CI_AS '$.a') AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) COLLATE Latin1_General_CI_AS '$.a' AS JSON) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH ([select] INT) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (window INT) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT, a INT) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a dbo.mytype) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j, '$.x') WITH (a INT) j
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH () AS x
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT,) AS x
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a) AS x
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT NULL) AS x
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'NULL'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT 5) AS x
|
< error
| 102 1 15 1 Incorrect syntax near '5'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT @j) AS x
|
< error
| 102 1 15 1 Incorrect syntax near '@j'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT '$.a' + '') AS x
|
< error
| 102 1 15 1 Incorrect syntax near '+'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT '$.a' '$.b') AS x
|
< error
| 102 1 15 1 Incorrect syntax near '$.b'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) '$.a' COLLATE Latin1_General_CI_AS) AS x
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) AS JSON '$.a') AS x
|
< error
| 102 1 15 1 Incorrect syntax near '$.a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) AS JSON COLLATE Latin1_General_CI_AS) AS x
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT AS) AS x
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT AS XML) AS x
|
< error
| 102 1 15 1 Incorrect syntax near 'XML'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' AS [JSON]) AS x
|
< error
| 102 1 15 1 Incorrect syntax near 'JSON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' JSON) AS x
|
< error
| 102 1 15 1 Incorrect syntax near 'JSON'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE) AS x
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE [Latin1_General_CI_AS]) AS x
|
< error
| 102 1 15 1 Incorrect syntax near 'Latin1_General_CI_AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE 'x') AS x
|
< error
| 102 1 15 1 Incorrect syntax near 'x'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (select INT) AS x
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a.b INT) AS x
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENJSON(@j) WITH (a INT) x (c)
|
< error
| 102 1 15 1 Incorrect syntax near 'c'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT)
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a', b NVARCHAR(10) '@mp:id')
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a NVARCHAR(10) COLLATE Latin1_General_CI_AS '@a') x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r') WITH (a INT) y
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT) AS y
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) y
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT 1)
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a' AS JSON)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a' COLLATE Latin1_General_CI_AS) x
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENQUERY(srv, 'SELECT 1') WITH (a INT) x
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; SELECT * FROM OPENROWSET('SQLNCLI', 'x', 'SELECT 1') WITH (a INT) x
|
< error
| 102 1 15 1 Incorrect syntax near 'WITH'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetWithOracle]
< ok 0 -
