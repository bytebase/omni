# GoogleSQL parser test data

## `compat.tsv`: the Spanner emulator's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/GoogleSQL.g4`, then its near-misses, with omni's verdict and
the pinned Spanner emulator's (`../../testdata/pins.tsv`). `GoogleSQL.g4` is
converted by `../cmd/tm2g4` from `googlesql.tm`, the Textmapper grammar of the
GoogleSQL release the emulator builds against, beside it as published.
`TestCompatibility` (`../compat_test.go`) reads the file through
`enginetest/compat` and starts no server. Recorded 2026-11-14 in about a
minute; a second recording on a fresh emulator wrote the same file.

| Rows | Total | Agree | Only Spanner accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 2,614 | 1,470 | 753 | 391 |
| Near-misses | 6,223 | 4,954 | 1,195 | 74 |

The rows start from `unterminated_statement`, a statement of a script, so
they cover GoogleSQL's scripting and BigQuery statements as well as what
Spanner runs. `select_column_expr`'s `ATTACHED_ALIAS` alternative gets no row:
it refuses every input (`../../testdata/tokens.tsv`).

Most of what only Spanner accepts is GoogleSQL syntax the emulator parses and
then refuses as a statement it does not run (`Statement not supported`): a
statement that starts with a pipe operator (`|> AGGREGATE`, 1,154 rows),
`MERGE`, `LOAD DATA`, `TABLE a(...)`, `IMPORT`, `EXPORT DATA`, `SHOW`,
`EXECUTE IMMEDIATE`, and graph subqueries (`SELECT ARRAY { RETURN NULL }`).
Most of what only omni accepts is BigQuery: scripting (`DECLARE`, `IF`,
`CASE`, `BREAK`), and DDL Spanner's DDL parser refuses (`CREATE TABLE a` with
no columns, `CREATE FUNCTION`, `CREATE ROW POLICY`, `GRANT ... TO 'a'`).

## Suites

The hand-written corpora that held `Parse` to a live emulator are `suite:`
rows, judged by the same recording: the fixtures of the `*_oracle_test.go`
files that compared `Parse` of a statement with the emulator's verdict on it.

| Suite | Rows | Agree | Only Spanner accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| `suite:select` | 95 | 95 | 0 | 0 |
| `suite:pivot-unpivot` | 49 | 49 | 0 | 0 |
| `suite:dml` | 68 | 68 | 0 | 0 |
| `suite:ddl` | 91 | 91 | 0 | 0 |
| `suite:spanner-ddl` | 81 | 81 | 0 | 0 |
| `suite:utility` | 67 | 67 | 0 | 0 |
| `suite:scripting` | 47 | 23 | 3 | 21 |
| `suite:graph-query` | 86 | 85 | 0 | 1 |
| `suite:bq-ddl` | 24 | 4 | 0 | 20 |
| `suite:export-load-clone` | 18 | 11 | 7 | 0 |

The disagreeing suite rows are the forms those tests recorded as divergences:
BigQuery forms Spanner refuses, and `BEGIN SELECT 1 END`, which the emulator
takes and omni refuses. A statement added to a suite is a row with `?` in both
verdict cells, judged by a recording:

```text
?	?	-	suite:select	<statement>
```

## What the judge needs

- **The harness decides.** `googlesql/spanner-oracle` sends each statement to
  the emulator by its first word: a query runs, DML runs in a transaction it
  rolls back, and DDL is applied to a database of its own
  (`../../spanner-oracle/README.md`). Its verdict is the judge's
  (`spannerVerdict` in `../compat_test.go`).
- **A rejection is the parser's error:** `Syntax error:` from GoogleSQL's
  parser, also where the emulator parses SQL inside DDL (`CREATE PROPERTY
  GRAPH`, a function's body), or `Error parsing Spanner DDL statement:` from
  the emulator's own DDL parser. Every other answer is an acceptance,
  `Statement not supported` and `Unimplemented` (a query that returns a
  `STRUCT` column) included.
- **Every statement sees an empty schema.** The emulator answers
  `CREATE PROPERTY GRAPH a` by the names the schema holds before it parses the
  rest, so the harness applies DDL to a database it replaces after each DDL
  statement that applies.
- **One statement, no `;`.** The emulator's DDL path refuses a terminator.
