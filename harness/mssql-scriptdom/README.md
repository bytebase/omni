# mssql-scriptdom harness

A small .NET console app that parses T-SQL with
`Microsoft.SqlServer.TransactSql.ScriptDom` and emits a compact JSON shape
covering the fields we diff against omni's mssql parser.

## Build & run

```bash
dotnet build -c Release
echo "SELECT * FROM (VALUES (1,2)) AS t(a,b)" | dotnet bin/Release/net8.0/mssql-scriptdom-harness.dll
```

Output:
```json
{"shape":{"kind":"Script","stmts":[{"kind":"Select","query":{"kind":"QuerySpec","select_list":[{"kind":"SelectStar"}],"from":[{"kind":"InlineDerivedTable","alias":"t","columns":["a","b"]}]}}]}}
```

## Batch / line mode

Set `MSSQL_HARNESS_LINE=1` and feed one base64-encoded SQL per line on
stdin. The harness emits one JSON line per input. Avoids newline-in-SQL
ambiguity and lets callers reuse a single process across many fixtures.

## Consumers

`mssql/parser/scriptdom_harness_test.go` (Go, build tag `scriptdom`) drives
this harness in batch mode and diffs omni's AST against the ScriptDOM shape.

Run with:
```bash
go test -tags scriptdom ./mssql/parser/ -run TestScriptDOMDiff -v
```

## Shape covered

Only the fields we currently diff on:

- **statements**: Select / Insert / other (tagged by class name)
- **select list**: SelectScalar (alias + expr_kind), SelectStar, SelectSetVariable
- **FROM table references**: NamedTableReference, QueryDerivedTable,
  InlineDerivedTable, SchemaObjectFunctionTableReference, OpenJsonTableReference,
  VariableMethodCallTableReference, QualifiedJoin / UnqualifiedJoin
- **alias** + **columns[]** (plain name list) + **with_cols[{name,type}]** (WITH-clause)

Extend `Program.cs::Shape.OfTableRef` / `OfSelectElement` as new scenarios
require additional fields.
