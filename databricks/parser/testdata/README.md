# Databricks parser test data

## `compat.tsv`: Apache Spark 4.0.5 SQL compatibility rows

`compat.tsv` holds a statement for each reachable alternative of each rule in
`../../testdata/SqlBaseParser.g4` (paired with `SqlBaseLexer.g4`), followed by
its near-misses (one token dropped, the end cut off, a stray word inserted, or a
keyword swapped), plus hand-curated `suite:` rows covering Spark `AstBuilder`
parse-time checks, keyword-as-identifier edge cases, and reverted
overcorrections. Each row records omni's verdict and Apache Spark 4.0.5's live
parser verdict (`spark.sessionState.sqlParser.parsePlan` with
`spark.sql.scripting.enabled=true`). `TestCompatibility` (`../compat_test.go`)
replays the committed verdicts offline through `enginetest/compat` on every test
run with no container, credentials, or network.

| Rows | Total | Agree | Only Spark accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 2,229 | 2,229 | 0 | 0 |
| Near-misses | 4,209 | 4,209 | 0 | 0 |
| `suite:keywords` | 45 | 45 | 0 | 0 |
| `suite:overcorrections` | 11 | 11 | 0 | 0 |
| `suite:spark-astbuilder` | 232 | 232 | 0 | 0 |
| **All** | **6,726** | **6,726** | **0** | **0** |

- **The judge** is the pinned official `apache/spark:4.0.5` container image
  (`../../testdata/pins.tsv`, role `default`, digest
  `sha256:45eb1b8261a89bc7692be05aa4664676856f1c886cecc7b97b4e10f635b45b3c`,
  source tag `v4.0.5`, judge revision `1`). Each row is judged by calling
  `spark.sessionState.sqlParser.parsePlan(sql)` with
  `spark.sql.scripting.enabled=true`, which runs both Spark's ANTLR 4
  `SqlBaseLexer`/`SqlBaseParser` and Spark's full `AstBuilder` /
  `DataTypeAstBuilder` parse-time validation without executing queries or
  requiring a catalog.
- **The grammar** used to generate the 6,438 baseline rows is
  `../../testdata/SqlBaseParser.g4` with `../../testdata/SqlBaseLexer.g4`
  (`sql/api/src/main/antlr4/org/apache/spark/sql/catalyst/parser/` in
  [`apache/spark`](https://github.com/apache/spark), `v4.0.5`, Apache-2.0 in
  `../../testdata/Spark.LICENSE`).
- **Role `driver` in `../../testdata/pins.tsv` is untouched** and continues to
  pin the Databricks workspace for `databricks/driver` tapes.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to the lexer rules
  of `SqlBaseLexer.g4`. One rule override (`primaryExpression` -> `a`) keeps
  generated rows readable because `SqlBaseParser.g4` lists `CURRENT_DATE` before
  identifiers in `primaryExpression`. `.` is sampled as `a` for the `.*?`
  wildcard productions in `REFRESH`, `ADD`/`LIST`, `SET`/`RESET`, and
  `unsupportedHiveNativeCommands`.
- **22 alternatives have no generated row:** Spark's 8 auxiliary non-statement
  entry points (`singleExpression`, `singleTableIdentifier`,
  `singleMultipartIdentifier`, `singleFunctionIdentifier`, `singleDataType`,
  `singleTableSchema`, `singleRoutineParamList`, `singleQuery`) and 4 unused
  helper rules (`executeImmediateArgument`, `executeImmediateArgumentSeq`,
  `functionIdentifier`, `qualifiedNameList`), none of which are reachable from
  the start rule `compoundOrSingleStatement`.

To judge new rows, or every row after a pin bump (starts the pinned
`apache/spark:4.0.5` container locally via Docker):

```bash
OMNI_RECORD=1 go test ./databricks/parser/ -run '^TestCompatibility$'
```

To verify committed answers read-only against the pinned `apache/spark:4.0.5`
container:

```bash
OMNI_RECORD=verify go test ./databricks/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./databricks/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

- **Live Spark 4.0.5 engine judge (`SourceEngine`), 6,726 rows.** All 6,438
  grammar-generated rows plus 288 hand-curated `suite:` rows
  (`suite:keywords`, `suite:overcorrections`, `suite:spark-astbuilder`) are
  judged by `apache/spark:4.0.5`'s `SparkSession.sessionState.sqlParser.parsePlan`
  with `spark.sql.scripting.enabled=true` and replayed offline in CI.
- **Agreement: 6,726 / 6,726 (100.00%)** — 0 `judge-only` rows and 0 `omni-only`
  rows.

## Databricks SQL dialect deltas beyond Apache Spark 4.0.5

Apache Spark 4.0.5 is the open-source SQL parser foundation of Databricks SQL,
so `compat.tsv` uses `apache/spark:4.0.5` as its reproducible container judge.
`databricks/parser` intentionally accepts the following Databricks SQL
extensions that Spark 4.0.5 rejects (covered directly in `../parser_test.go`,
not hand-edited in `compat.tsv`):

- **`QUALIFY` clause.** Databricks SQL supports `QUALIFY <booleanExpression>`
  after `HAVING`/`WINDOW` in `SELECT` queries (`SELECT ... FROM t QUALIFY
  row_number() OVER (...) = 1`), whereas Apache Spark 4.0.5 rejects `QUALIFY`.
- **Unity Catalog / Hive Metastore `GRANT` and `REVOKE`.** Databricks SQL
  supports `GRANT <privileges> ON <securable> TO <principal>` and `REVOKE
  <privileges> ON <securable> FROM <principal>`. Spark 4.0.5 rejects all `GRANT`
  and `REVOKE` statements (`INVALID_STATEMENT_OR_CLAUSE`), including bare
  `GRANT`, `REVOKE`, `GRANT ROLE`, and `REVOKE zzz`; `databricks/parser`
  accepts only full `GRANT ... ON ... TO ...` / `REVOKE ... ON ... FROM ...`
  statements and continues to reject bare/incomplete `GRANT`/`REVOKE` forms so
  all `compat.tsv` rows agree with Spark 4.0.5.

## Notes on Spark `AstBuilder` checks and parser scope

- **`parsePlan` vs analyzer/execution checks.** `parsePlan` covers all ANTLR
  syntax rules and Spark's `AstBuilder` / `DataTypeAstBuilder` parse-time checks
  (including literal bounds, hex/date/timestamp/interval parsing, collation
  validation on `STRING COLLATE` and `DEFAULT COLLATION`, `TRIM`/`RTRIM`
  collation rejection on expressions, duplicate property/partition/CTE/clause
  checks, `LATERAL` join restrictions, `SQLSTATE` validation, function resource
  types, `DESCRIBE ... AS JSON`, `CREATE VIEW` parameter markers, `IDENTITY`
  column rules, `SQL SECURITY INVOKER` rejection, `CHANGE COLUMN` rename
  rejection, `SET`/`RESET`/`REFRESH` validation, `TABLESAMPLE BUCKET ... ON`
  rejection, `GROUP BY` alias and `ROLLUP(())` restrictions, pipe `AGGREGATE`
  rules, and scripting scope/label/ordering checks). Semantic checks that Spark
  defers to the Catalyst analyzer (such as table/column resolution, type
  checking, or expression `COLLATE` name resolution beyond `TRIM`/`RTRIM`) are
  outside `parsePlan` and live in `databricks/analysis`.
- **Multi-statement scripts.** `Parse()` splits and parses multiple
  semicolon-separated top-level statements and compound `BEGIN ... END` blocks
  in one input, whereas Spark's `parsePlan` accepts a single statement or
  compound block; multi-statement spans are covered directly in
  `../parser_test.go`.
