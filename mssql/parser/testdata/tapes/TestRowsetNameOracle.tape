-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestRowsetNameOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH db.dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH db..t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH [dbo].[t]
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t AS y
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t y
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t, t1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t JOIN t1 ON 1 = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t y JOIN t1 ON 1 = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH srv.db.dbo.t
|
< error
| 7004 1 16 1 OPENXML does not support retrieving schema from remote tables, as in 'srv.db.dbo.t'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH db..
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo. AS y
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH db.dbo.
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo. t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; DECLARE @tv TABLE (a INT); SELECT * FROM OPENXML(@h, '/r', 1) WITH @tv
|
< error
| 102 1 15 1 Incorrect syntax near '@tv'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t y (c)
|
< error
| 102 1 15 1 Incorrect syntax near 'c'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t WITH (NOLOCK)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t y WITH (NOLOCK)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t (NOLOCK)
|
< error
| 1018 1 15 1 Incorrect syntax near 'NOLOCK'. If this is intended as a part of a table hint, A WITH keyword and parenthesis are now required. See SQL Server Books Online for proper syntax.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH t TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH (SELECT 1 AS a)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'SELECT'.
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.f()
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENXML(@h, '/r', 1) WITH select
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENJSON(N'[]') WITH t
|
< error
| 102 1 15 1 Incorrect syntax near 't'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENJSON(N'[]') WITH dbo.t AS y
|
< error
| 102 1 15 1 Incorrect syntax near 'dbo'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t AS x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t x
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db..t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI')...t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').[db].[dbo].[t]
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI'). db.dbo.t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t JOIN t1 ON 1 = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM t1 JOIN OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t x ON x.a = t1.a
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI')
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI') x
|
< error
| 102 1 15 1 Incorrect syntax near 'x'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').t
|
< error
| 208 1 16 1 Invalid object name 't'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').dbo.t
|
< error
| 208 1 16 1 Invalid object name 'dbo.t'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').x.db.dbo.t
|
< error
| 208 1 16 1 Invalid object name 'x.db.dbo.t'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t.c
|
< error
| 208 1 16 1 Invalid object name 'db.dbo.t.c'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.f()
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t WITH (NOLOCK)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t x WITH (NOLOCK)
|
< error
| 319 1 15 1 Incorrect syntax near the keyword 'with'. If this statement is a common table expression, an xmlnamespaces clause or a change tracking context clause, the previous statement must be terminated with a semicolon.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t x (NOLOCK)
|
< error
| 1018 1 15 1 Incorrect syntax near 'NOLOCK'. If this is intended as a part of a table hint, A WITH keyword and parenthesis are now required. See SQL Server Books Online for proper syntax.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t TABLESAMPLE (10 PERCENT)
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'TABLESAMPLE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI').db.dbo.t FOR SYSTEM_TIME ALL
|
< error
| 102 1 15 1 Incorrect syntax near 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('MSOLEDBSQL').db.dbo.t
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE('a', 'b', 'c').db.dbo.t
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE(@h, 'b').db.dbo.t
|
< error
| 102 1 15 1 Incorrect syntax near '@h'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENDATASOURCE(MSOLEDBSQL, 'b').db.dbo.t
|
< error
| 102 1 15 1 Incorrect syntax near 'MSOLEDBSQL'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENQUERY(srv, 'SELECT 1').db.dbo.t
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @h INT; SELECT * FROM OPENROWSET('SQLNCLI', 'x', 'SELECT 1').db.dbo.t
|
< error
| 102 1 15 1 Incorrect syntax near '.'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestRowsetNameOracle]
< ok 0 -
