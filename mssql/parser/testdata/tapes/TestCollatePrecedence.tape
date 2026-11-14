-- pin: default mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04 sha256:cf8f0ab2ebe61feec2eed2629061a07d247f791c2a8ddda74a4ac9f06a291718 16.0.4305.5
-- pin: scriptdom mcr.microsoft.com/dotnet/sdk:8.0.425 sha256:ec9c0a0dc5f60adc2065762050638dd5ed5facfc53716d4580a42d86027c8e80 161.9142.0

== TestCollatePrecedence default#1
> exec
| CREATE TABLE dbo.t (a INT, col INT, partition INT, encryption INT, window INT, bucket INT)
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' COLLATE Latin1_General_CI_AS + 'b'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' + 'b' COLLATE Latin1_General_CI_AS
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 'a' COLLATE Latin1_General_CI_AS + 'b' COLLATE Latin1_General_BIN + 'c'
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x * y COLLATE Latin1_General_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS % 2 FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT -x COLLATE Latin1_General_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT ~x COLLATE Latin1_General_CI_AS & 1 FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT (x + y) COLLATE Latin1_General_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT CAST(x AS VARCHAR(10)) COLLATE Latin1_General_CI_AS + 'b' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT x COLLATE Latin1_General_CI_AS AT TIME ZONE 'UTC' COLLATE Latin1_General_BIN FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 + d AT TIME ZONE 'UTC' COLLATE Latin1_General_CI_AS FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'U' + 'TC' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT 1 + d AT TIME ZONE 'UTC' + 2 FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d * 1 AT TIME ZONE 'UTC' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT d AT TIME ZONE 'UTC' * 2 FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -
> exec
| SET PARSEONLY ON
< ok 0 -
> exec
| SELECT -d AT TIME ZONE 'UTC' FROM t
|
< ok 0 -
> exec
| SET PARSEONLY OFF
| SET QUOTED_IDENTIFIER ON
| SET FIPS_FLAGGER OFF
| USE [omni_TestCollatePrecedence]
< ok 0 -

== TestCollatePrecedence scriptdom#1
> query
| SELECT 'a' COLLATE Latin1_General_CI_AS + 'b'
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}]}}]},"exprs":["(\u002B (collate # Latin1_General_CI_AS) #)"]}
> query
| SELECT 'a' + 'b' COLLATE Latin1_General_CI_AS
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}]}}]},"exprs":["(\u002B # (collate # Latin1_General_CI_AS))"]}
> query
| SELECT 'a' COLLATE Latin1_General_CI_AS + 'b' COLLATE Latin1_General_BIN + 'c'
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}]}}]},"exprs":["(\u002B (\u002B (collate # Latin1_General_CI_AS) (collate # Latin1_General_BIN)) #)"]}
> query
| SELECT x * y COLLATE Latin1_General_CI_AS FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(* # (collate # Latin1_General_CI_AS))"]}
> query
| SELECT x COLLATE Latin1_General_CI_AS % 2 FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(% (collate # Latin1_General_CI_AS) #)"]}
> query
| SELECT -x COLLATE Latin1_General_CI_AS FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"UnaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(u- (collate # Latin1_General_CI_AS))"]}
> query
| SELECT ~x COLLATE Latin1_General_CI_AS & 1 FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(\u0026 (u~ (collate # Latin1_General_CI_AS)) #)"]}
> query
| SELECT (x + y) COLLATE Latin1_General_CI_AS FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ParenthesisExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(collate (paren (\u002B # #)) Latin1_General_CI_AS)"]}
> query
| SELECT CAST(x AS VARCHAR(10)) COLLATE Latin1_General_CI_AS + 'b' FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(\u002B (collate # Latin1_General_CI_AS) #)"]}
> query
| SELECT x COLLATE Latin1_General_CI_AS AT TIME ZONE 'UTC' COLLATE Latin1_General_BIN FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"AtTimeZoneCall"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(atz (collate # Latin1_General_CI_AS) (collate # Latin1_General_BIN))"]}
> query
| SELECT 1 + d AT TIME ZONE 'UTC' COLLATE Latin1_General_CI_AS FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"AtTimeZoneCall"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(atz (\u002B # #) (collate # Latin1_General_CI_AS))"]}
> query
| SELECT d AT TIME ZONE 'U' + 'TC' FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(\u002B (atz # #) #)"]}
> query
| SELECT 1 + d AT TIME ZONE 'UTC' + 2 FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"BinaryExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(\u002B (atz (\u002B # #) #) #)"]}
> query
| SELECT d * 1 AT TIME ZONE 'UTC' FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"AtTimeZoneCall"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(atz (* # #) #)"]}
> query
| SELECT d AT TIME ZONE 'UTC' * 2 FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"AtTimeZoneCall"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(atz # (* # #))"]}
> query
| SELECT -d AT TIME ZONE 'UTC' FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"AtTimeZoneCall"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["(atz (u- #) #)"]}
