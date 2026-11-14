-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5

== TestPivotBodyOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (max(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (Sum(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX (col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (dbo.agg(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (db.dbo.agg(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT ([dbo].[agg](col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT ([dbo].[MAX](col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (dbo.window(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col, a) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (STRING_AGG(col, a) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(t.col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(db.dbo.t.col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX([col]) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(window) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (COUNT(*) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '*'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(t.*) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '*'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col + 1) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '+'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(DISTINCT col) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'DISTINCT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(ALL col) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ALL'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(1) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(@v) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '@v'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX((col)) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT ((MAX(col)) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col COLLATE Latin1_General_CI_AS) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COLLATE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (STRING_AGG(col, ',') FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near ','.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX() FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (col FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) OVER () FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'OVER'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(x.db.dbo.t.col) FOR a IN ([1])) p
|
< error
| 117 1 15 1 The column name 'x.db.dbo.t' contains more than the maximum number of prefixes. The maximum is 3.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(db..col) FOR a IN ([1])) p
|
< error
| 1004 1 15 1 Invalid column prefix 'db.': No table name specified
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(.col) FOR a IN ([1])) p
|
< error
| 1004 1 15 1 Invalid column prefix '': No table name specified
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(IDENTITYCOL) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IDENTITYCOL'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX($IDENTITY) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '$IDENTITY'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX($action) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '$action'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT ([MAX](col) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT ("MAX"(col) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (myagg(col) FOR a IN ([1])) p
|
< error
| 195 1 15 1 'myagg' is not a recognized aggregate function.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (GROUPING(col) FOR a IN ([1])) p
|
< error
| 195 1 15 1 'GROUPING' is not a recognized aggregate function.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (GROUPING_ID(col, a) FOR a IN ([1])) p
|
< error
| 195 1 15 1 'GROUPING_ID' is not a recognized aggregate function.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (ABS(col) FOR a IN ([1])) p
|
< error
| 195 1 15 1 'ABS' is not a recognized aggregate function.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (ROW_NUMBER(col) FOR a IN ([1])) p
|
< error
| 195 1 15 1 'ROW_NUMBER' is not a recognized aggregate function.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (LEFT(col) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'LEFT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (COALESCE(col) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'COALESCE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (srv.db.dbo.agg(col) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (db..agg(col) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (.agg(col) FOR a IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (dbo.select(col) FOR a IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'select'.
| 102 1 15 1 Incorrect syntax near 'a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR t.a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR db.dbo.t.a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR [a] IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR window IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR x.db.dbo.t.a IN ([1])) p
|
< error
| 117 1 15 1 The column name 'x.db.dbo.t' contains more than the maximum number of prefixes. The maximum is 3.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR db..a IN ([1])) p
|
< error
| 1004 1 15 1 Invalid column prefix 'db.': No table name specified
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a. IN ([1])) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR (a) IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR $IDENTITY IN ([1])) p
|
< error
| 102 1 15 1 Incorrect syntax near '$IDENTITY'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1], [2])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (x, y)) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ("x")) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1], [1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (window)) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (t.x)) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (db.dbo.t.x)) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]) ) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a) p
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN) p
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN [1]) p
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ()) p
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1],)) p
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1] [2])) p
|
< error
| 102 1 15 1 Incorrect syntax near '2'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (1)) p
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ('x')) p
|
< error
| 102 1 15 1 Incorrect syntax near 'x'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (NULL)) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'NULL'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (@v)) p
|
< error
| 102 1 15 1 Incorrect syntax near '@v'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (select)) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1] AS x)) p
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (y.db.dbo.t.x)) p
|
< error
| 117 1 15 1 The column name 'y.db.dbo.t' contains more than the maximum number of prefixes. The maximum is 3.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN (db..x)) p
|
< error
| 1004 1 15 1 Invalid column prefix 'db.': No table name specified
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]) extra) p
|
< error
| 102 1 15 1 Incorrect syntax near 'extra'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS [p]
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) "p"
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) window
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS window
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS p WHERE p.a = 1
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1]))
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) WHERE 1 = 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'WHERE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) ORDER BY 1
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'ORDER'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS
|
< error
| 102 1 15 1 Incorrect syntax near 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS select
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'select'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) AS 'p'
|
< error
| 102 1 15 1 Incorrect syntax near 'p'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) 'p'
|
< error
| 102 1 15 1 Incorrect syntax near 'p'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (a, col)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (a)) AS window
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT ([v] FOR [c] IN ([a])) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (window FOR c IN (a)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (window)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (t.v FOR c IN (a)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR t.c IN (a)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (t.a, t.col)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (db.dbo.t.a)) u
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (FOR c IN (a)) u
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FOR'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v c IN (a)) u
|
< error
| 102 1 15 1 Incorrect syntax near 'c'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR IN (a)) u
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'IN'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c (a)) u
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN a) u
|
< error
| 102 1 15 1 Incorrect syntax near 'a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN ()) u
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN ('a')) u
|
< error
| 102 1 15 1 Incorrect syntax near 'a'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (1)) u
|
< error
| 102 1 15 1 Incorrect syntax near '1'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (a AS x)) u
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (x.db.dbo.t.a)) u
|
< error
| 117 1 15 1 The column name 'x.db.dbo.t' contains more than the maximum number of prefixes. The maximum is 3.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (db..a)) u
|
< error
| 1004 1 15 1 Invalid column prefix 'db.': No table name specified
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (x.db.dbo.t.v FOR c IN (a)) u
|
< error
| 117 1 15 1 The column name 'x.db.dbo.t' contains more than the maximum number of prefixes. The maximum is 3.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR x.db.dbo.t.c IN (a)) u
|
< error
| 117 1 15 1 The column name 'x.db.dbo.t' contains more than the maximum number of prefixes. The maximum is 3.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (MAX(v) FOR c IN (a)) u
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT ((v) FOR c IN (a)) u
|
< error
| 102 1 15 1 Incorrect syntax near '('.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (@v FOR c IN (a)) u
|
< error
| 102 1 15 1 Incorrect syntax near '@v'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (a))
|
< error
| 102 1 15 1 Incorrect syntax near ')'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t UNPIVOT (v FOR c IN (a)) AS
|
< error
| 102 1 15 1 Incorrect syntax near 'AS'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (APPROX_COUNT_DISTINCT(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (APPROX_PERCENTILE_CONT(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (APPROX_PERCENTILE_DISC(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (AVG(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (CHECKSUM_AGG(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (COUNT(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (COUNT_BIG(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MAX(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (MIN(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (STDEV(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (STDEVP(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (STRING_AGG(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (SUM(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (VAR(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT * FROM t PIVOT (VARP(col) FOR a IN ([1])) p
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestPivotBodyOracle]
< ok 0 -
