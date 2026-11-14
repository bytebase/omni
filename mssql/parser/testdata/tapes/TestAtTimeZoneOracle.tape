-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5
-- pin: scriptdom mcr.microsoft.com/dotnet/sdk:8.0.425 sha256:ec9c0a0dc5f60adc2065762050638dd5ed5facfc53716d4580a42d86027c8e80 161.9142.0

== TestAtTimeZoneOracle default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' AT FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT y FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'y'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT TIME y FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT TIME FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT TIME
|
< error
| 102 1 15 2 Incorrect syntax near 'TIME'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT ZONE 'UTC' FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'ZONE'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x AT TIME ZONE FROM t
|
< error
| 156 1 15 1 Incorrect syntax near the keyword 'FROM'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS AT y FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'y'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' AT y FROM t
|
< error
| 102 1 15 1 Incorrect syntax near 'y'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE x AT y = 1
|
< error
| 4145 1 15 1 An expression of non-boolean type specified in a context where a condition is expected, near 'AT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 FROM t WHERE x AT TIME y = 1
|
< error
| 4145 1 15 1 An expression of non-boolean type specified in a context where a condition is expected, near 'AT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| UPDATE t SET x = y AT TIME 'UTC'
|
< error
| 102 1 15 1 Incorrect syntax near 'AT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| DECLARE @v DATETIMEOFFSET = SYSDATETIMEOFFSET() AT 'UTC'
|
< error
| 102 1 15 1 Incorrect syntax near 'AT'.
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestAtTimeZoneOracle]
< ok 0 -

== TestAtTimeZoneOracle scriptdom#1
> query
| SELECT x AT y FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027y\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"AT","expr_kind":"ColumnReferenceExpression"}]}}]},"exprs":["#"]}
> query
| SELECT x AT TIME y FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Expected WINDOW but encountered TIME instead."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x AT TIME FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027TIME\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"AT","expr_kind":"ColumnReferenceExpression"}]}}]},"exprs":["#"]}
> query
| SELECT x AT TIME
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027TIME\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"AT","expr_kind":"ColumnReferenceExpression"}]}}]},"exprs":["#"]}
> query
| SELECT x AT ZONE 'UTC' FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":13,"message":"Incorrect syntax near \u0027ZONE\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"AT","expr_kind":"ColumnReferenceExpression"}]}}]},"exprs":["#"]}
> query
| SELECT x AT TIME ZONE FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":23,"message":"Incorrect syntax near \u0027FROM\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT x COLLATE Latin1_General_CI_AS AT y FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":42,"message":"Incorrect syntax near \u0027y\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"AT","expr_kind":"ColumnReferenceExpression"}]}}]},"exprs":["(collate # Latin1_General_CI_AS)"]}
> query
| SELECT d AT TIME ZONE 'UTC' AT y FROM t
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":32,"message":"Incorrect syntax near \u0027y\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"AT","expr_kind":"AtTimeZoneCall"}]}}]},"exprs":["(atz # #)"]}
> query
| SELECT 1 FROM t WHERE x AT y = 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":25,"message":"Incorrect syntax near \u0027AT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| SELECT 1 FROM t WHERE x AT TIME y = 1
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":25,"message":"Incorrect syntax near \u0027AT\u0027."}],"shape":{"kind":"Script","stmts":[]},"exprs":null}
> query
| UPDATE t SET x = y AT TIME 'UTC'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":20,"message":"Incorrect syntax near \u0027AT\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"UpdateStatement"}]},"exprs":null}
> query
| DECLARE @v DATETIMEOFFSET = SYSDATETIMEOFFSET() AT 'UTC'
< result
< column "shape" "JSON"
< row
: string
| {"errors":[{"line":1,"column":49,"message":"Incorrect syntax near \u0027AT\u0027."}],"shape":{"kind":"Script","stmts":[{"kind":"DeclareVariableStatement"}]},"exprs":null}
