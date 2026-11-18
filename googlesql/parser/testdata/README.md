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
Rejudged at judge revision 2 (`../../testdata/pins.tsv`), which reads
GoogleSQL's grammar refusals as rejections: 90 rows moved to `reject`.

| Rows | Total | Agree | Only Spanner accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 2,614 | 1,888 | 327 | 399 |
| Near-misses | 6,223 | 5,079 | 1,110 | 34 |

The rows start from `unterminated_statement`, a statement of a script, so
they cover GoogleSQL's scripting and BigQuery statements as well as what
Spanner runs. `select_column_expr`'s `ATTACHED_ALIAS` alternative gets no row:
it refuses every input (`../../testdata/tokens.tsv`).

Most of what only Spanner accepts is GoogleSQL syntax the emulator parses and
then refuses as a statement it does not run (`Statement not supported`):
`MERGE` (80 rows), `LOAD DATA` (74), `TABLE a(…)` / table clauses (56),
`BEGIN … END` (47), `IMPORT` (25), `EXPORT DATA` (23), `SHOW` (23), and
`EXECUTE IMMEDIATE` (21). The pipe rows are different; see below.
Most of what only omni accepts is BigQuery: scripting (`DECLARE`, `IF`,
`CASE`, `BREAK`), and DDL Spanner's DDL parser refuses (`CREATE TABLE a` with
no columns, `CREATE FUNCTION`, `CREATE ROW POLICY`, `GRANT ... TO 'a'`).

## The pipe rows are not a parse verdict

GoogleSQL's pipe syntax is a family the grammar reaches three ways:
`sql_statement_body_no_pipe_suffix → subpipeline_statement →
subpipeline_no_parens` (a pipe suffix alone, `|> AGGREGATE`),
`query → query pipe_and_pipe_operator*` (`SELECT … |> WHERE …`), and
`sql_statement_body_maybe_pipe_suffix pipe_and_pipe_operator*` (`DESCRIBE … |>
…`). The first is a complete statement, not a fragment — the bare suffix is an
alternative of the start rule — so omni parses it (`SubpipelineStmt` in
`googlesql/ast`) rather than refusing it.

A recording cannot judge this family. The emulator answers a statement
containing a pipe operator as one it does not run, without a syntax error, so
`compat.tsv` records `accept` for every one of them — including text the grammar
refuses, such as `|> ALL`, `|> FULL`, `|> WHERE`, or a pipe statement with a
required token dropped (`|> ORDER BY`). Of the 1,165 pipe rows, 407 are valid
per the grammar and now agree; the other 758 are the only-Spanner rows left, and
the grammar, not the recording, is what refuses them. Two of those are the
select-list `expr.*` form over a non-path expression (`|> EXTEND NULL.*`): the
grammar derives it, and omni's dot-star support is path-only, as it was before
this family.

Closing a pipe row therefore means matching the grammar, not the cell: the cell
says `accept` for valid and invalid text alike. A later recording that reads a
real parse verdict for a pipe statement is what would change these rows.

## Suites

The hand-written corpora that held `Parse` to a live emulator are `suite:`
rows, judged by the same recording: the fixtures of the `*_oracle_test.go`
files that compared `Parse` of a statement with the emulator's verdict on it.

| Suite | Rows | Agree | Only Spanner accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| `suite:select` | 95 | 94 | 0 | 1 |
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

The other oracle tests compare something else with the emulator's verdict (an
expression or a type parsed alone, a Spanner docs block with its comments and
line breaks, which a row cannot hold) and replay tapes (`tapes/`): the
harness's answer for each statement they send, recorded with `OMNI_RECORD=1`.

## What the judge needs

- **The harness decides.** `googlesql/spanner-oracle` sends each statement to
  the emulator by its first word: a query runs, DML runs in a transaction it
  rolls back, and DDL is applied to a database of its own
  (`../../spanner-oracle/README.md`). Its verdict is the judge's
  (`spannerVerdict` in `../compat_test.go`).
- **A rejection is the parser's error:** `Syntax error:` from GoogleSQL's
  parser, also where the emulator parses SQL inside DDL (`CREATE PROPERTY
  GRAPH`, a function's body); a message GoogleSQL's grammar writes when an
  action refuses input, as `QUALIFY is not supported` where the emulator
  leaves the feature off (`grammar_refusals.go` in the harness, generated
  from `googlesql.tm` by `../cmd/tm2g4`); or `Error parsing Spanner DDL
  statement:` from the emulator's own DDL parser. Every other answer is an
  acceptance, `Statement not supported` and `Unimplemented` (a query that
  returns a `STRUCT` column) included.
- **Every statement sees an empty schema.** The emulator answers
  `CREATE PROPERTY GRAPH a` by the names the schema holds before it parses the
  rest, so the harness applies DDL to a database it replaces after each DDL
  statement that applies.
- **One statement, no `;`.** The emulator's DDL path refuses a terminator.
