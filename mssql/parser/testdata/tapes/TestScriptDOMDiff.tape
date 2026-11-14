-- pin: scriptdom mcr.microsoft.com/dotnet/sdk:8.0.425 sha256:ec9c0a0dc5f60adc2065762050638dd5ed5facfc53716d4580a42d86027c8e80 161.9142.0

== TestScriptDOMDiff scriptdom#1
> query
| SELECT * FROM (SELECT a, b FROM t) AS x(c1, c2)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"QueryDerivedTable","alias":"x","columns":["c1","c2"]}]}}]},"exprs":["*"]}
> query
| SELECT * FROM fn(1, 2) AS t(a, b)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"SchemaObjectFunctionTableReference","alias":"t","columns":["a","b"]}]}}]},"exprs":["*"]}
> query
| SELECT * FROM (VALUES (1, 2), (3, 4)) AS t(a, b)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"InlineDerivedTable","alias":"t","columns":["a","b"]}]}}]},"exprs":["*"]}
> query
| SELECT a = c FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","alias":"a","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| SELECT @v = x FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectSetVariable","variable":"@v","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM OPENJSON(@j) WITH (id int '$.id', name nvarchar(50) '$.name') AS t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"OpenJsonTableReference","alias":"t","with_cols":[{"name":"id","type":"int"},{"name":"name","type":"nvarchar"}]}]}}]},"exprs":["*"]}
> query
| SELECT * FROM (SELECT a FROM t) AS x
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"QueryDerivedTable","alias":"x"}]}}]},"exprs":["*"]}
> query
| SELECT a FROM t
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectScalar","expr_kind":"ColumnReferenceExpression"}],"from":[{"kind":"NamedTableReference"}]}}]},"exprs":["#"]}
> query
| SELECT * FROM CONTAINSTABLE(t, b, 'foo') ct
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"FullTextTableReference","alias":"ct"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM FREETEXTTABLE(t, (a, b), 'foo') ft
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"FullTextTableReference","alias":"ft"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM CONTAINSTABLE(t, b, 'foo', LANGUAGE 'English', 10) ct
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"FullTextTableReference","alias":"ct"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM SEMANTICKEYPHRASETABLE(t, b) s
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"SemanticTableReference","alias":"s"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM SEMANTICSIMILARITYTABLE(t, (a, b), 1) s
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"SemanticTableReference","alias":"s"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(t, a, 1, b, 2) s
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"SemanticTableReference","alias":"s"}]}}]},"exprs":["*"]}
> query
| SELECT * FROM PREDICT(MODEL = @model, DATA = dbo.mytable AS d) WITH (Score FLOAT, Label NVARCHAR(10)) AS p
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"PredictTableReference","alias":"p","with_cols":[{"name":"Score","type":"FLOAT"},{"name":"Label","type":"NVARCHAR"}]}]}}]},"exprs":["*"]}
> query
| SELECT * FROM PREDICT(MODEL = 0x01, DATA = dbo.mytable) WITH (Score FLOAT)
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"PredictTableReference","with_cols":[{"name":"Score","type":"FLOAT"}]}]}}]},"exprs":["*"]}
> query
| SELECT * FROM dbo.t CROSS APPLY PREDICT(MODEL = @m, DATA = dbo.fn(1) AS d) WITH (s float) p
< result
< column "shape" "JSON"
< row
: string
| {"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"UnqualifiedJoin","join_type":"CrossApply","left":{"kind":"NamedTableReference"},"right":{"kind":"PredictTableReference","alias":"p","with_cols":[{"name":"s","type":"float"}]}}]}}]},"exprs":["*"]}
