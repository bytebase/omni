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
| One per alternative | 2,614 | 2,146 | 69 | 399 |
| Near-misses | 6,223 | 5,134 | 1,055 | 34 |

The rows start from `unterminated_statement`, a statement of a script, so
they cover GoogleSQL's scripting and BigQuery statements as well as what
Spanner runs. `select_column_expr`'s `ATTACHED_ALIAS` alternative gets no row:
it refuses every input (`../../testdata/tokens.tsv`).

Most of what only Spanner accepts is GoogleSQL syntax the emulator parses and
then refuses as a statement it does not run (`Statement not supported`):
`MERGE` (80 rows), `LOAD DATA` (74), `TABLE a(…)` / table clauses (56),
`BEGIN … END` (47), `IMPORT` (25), `EXPORT DATA` (23), `SHOW` (23), and
`EXECUTE IMMEDIATE` (21). `SET TRANSACTION` and `SET @@…` are the same, and the
pipe and graph rows are different; see below.
Most of what only omni accepts is BigQuery: scripting (`DECLARE`, `IF`,
`CASE`, `BREAK`), and DDL Spanner's DDL parser refuses (`CREATE TABLE a` with
no columns, `CREATE FUNCTION`, `CREATE ROW POLICY`, `GRANT ... TO 'a'`).

## Most of those families are not parse verdicts either

For `MERGE`, `LOAD DATA`, `BEGIN … END`, `IMPORT`, `MODULE`, `EXPORT DATA`,
`EXPORT MODEL`, `SHOW`, `EXECUTE IMMEDIATE`, `UNDROP` and `CLONE DATA` the
emulator answers `Statement not supported: <Kind>Stmt` from statement dispatch,
BEFORE it parses the body, and the judge reads that as an acceptance. It says so
for `MERGE a USING a ON NULL WHEN MATCHED THEN DELETION` and
`MERGE a USING a ON NULL WHEN MATCHED THEN DELETE` alike, for
`LOAD DATA INTO FROM FILES()` and `SHOW a zzz`, and for `UNDROP CONNECTION`
with no name. So those rows say `accept` whatever the text is, and a near-miss
in them (a dropped or swapped token) cannot be closed by matching the cell: the
cell is not a verdict. `UNDROP` adds a grammar reason of its own —
`undrop_statement`'s action allows only the SCHEMA kind and reports
"UNDROP <kind> is not supported" for the rest, so the `UNDROP CONNECTION …`
rows the generator wrote are text the grammar refuses.

The `TABLE` clause is the exception, and the work here: the emulator parses a
TABLE statement in full before it answers `TABLE statements are not supported`,
so it rejects `TABLE a zzz`, `TABLE a AS x`, `TABLE a PIVOT(…)` and `TABLE a(`
on syntax and accepts the `tvf_argument` shapes. Its rows are a real verdict,
and omni now reads `TABLE a`, `TABLE a(…)` with DESCRIPTOR / MODEL / CONNECTION /
INPUT TABLE / `name => TABLE t` arguments, and `TABLE a WHERE …` — 55 rows. The
same change closed `SHOW` (its three forms, including `SHOW [MATERIALIZED]
VIEWS … FROM … LIKE`, 11), `IMPORT MODULE|PROTO` with `AS`/`INTO` aliases and
`OPTIONS` (9), `MODULE <path> [OPTIONS]` (5), `EXPORT TABLE METADATA FROM` and
`EXPORT DATA AS GRAPH` (3), and the empty `EXCEPTION WHEN ERROR THEN END`
handler that `statement_list`'s empty alternative allows (2): 85 rows in all,
none of them a cell the recording could not judge. One TABLE near-miss is
left: `TABLE a() QUALIFY`, where QUALIFY is the clause word and omni reads no
alias there.

## `SET` is judged on one alternative only

`set_statement` splits by target, and the recording parses exactly one target:
`SET a = NULL` and its near-misses (`SET a = NULL zzz` rejects) are a parse
verdict, and omni agrees on all of them. The other two are not. `SET @@ a = NULL
zzz` and `SET TRANSACTION READ ONLY zzz` are recorded `accept` while their
valid forms are too, so the inserted token changed nothing: the emulator
answers `Statement not supported` for a system-variable or transaction-mode SET
before it parses the assignment. The 10 rows left therefore say `accept` for
text the grammar refuses (`SET TRANSACTION ONLY`, `SET TRANSACTION READ`, `SET
@@ a =`), and matching them would mean accepting what `set_statement` does not
derive.

## The graph rows are a parse verdict

The graph sub-language is the largest family the grammar derives and omni did
not read: `graph_expression`'s endpoint (`IS [NOT] SOURCE|DESTINATION [OF]`) and
label (`IS [NOT] LABELED`) predicates, the braced graph subqueries wherever a
query can appear (`SELECT ARRAY { … }`, `SELECT EXISTS { … }`, `SELECT VALUE {
… }`, `expr IN { … }`, and a CALL's inline subquery), and the operator, pattern
and CALL alternatives inside them. Because a graph query is reached only through
an expression, every one of those rows was a `reject` at the `{`: 208 rows moved
on the one change.

Unlike the pipe and dispatch families, the emulator judges this one: inside the
same rules it rejects `SELECT EXISTS { MATCH }`, `{ MATCH - MATCH }`, `{ - MATCH
- }`, `{ SHORTEST - }`, `{ - [ ] - }` and `{ < - [ ] - }` (the grammar's
`ValidateNoWhitespace` on an edge pattern's delimiters), `{ (a)* }` (a quantifier
on a node pattern), `{(@1 -) }` (a hint opening a path pattern), `{ -@1 - }` (a
hint between two edges) and `{(COST NULL) }`. Text the grammar derives is
accepted and text it refuses is rejected, so the cells are verdicts. The
`COST` filler tail is the one place the recording disagrees with
`googlesql.tm` — `opt_graph_cost/COST_expression` is recorded `reject` — so omni
reads the word as an element named `COST` and not as the tail, which is what the
grammar's own shift preference and the cell agree on.

The rows left in this family are not verdicts: `DYNAMIC LABEL(NULL) DYNAMIC
LABEL(NULL)` and its `drop-10` twin are recorded `accept` with and without the
closing `)`, and `PROPERTIES ALL COLUMNS EXCEPT a))` is recorded `accept` with a
stray `)`, so the emulator answered something other than a syntax error for a
`CREATE PROPERTY GRAPH` whose `DYNAMIC LABEL` support its DDL parser reports
later (4 rows).

`SELECT NULL.*` (3 rows, and the two pipe `|> EXTEND NULL.*` rows) is a real
verdict the grammar derives from `select_column_dot_star`'s
`expression_higher_prec_than_and "." "*"`; omni's dot-star is path-only, as it
was before either family, so they stay. `DROP MODEL a` is the same shape of real
verdict outside every family here: `MODEL` is one alternative of the grammar's
`schema_object_kind`, and the emulator's DDL parser reads it and reports the
missing model later.

## The keyword rows

`common_keyword_as_identifier` is a word the grammar takes as a name, and omni
reserved six of them: `MATCH_RECOGNIZE` (which the legacy grammar lists through
the placeholder token `KW_MATCH_RECOGNIZE_NONRESERVED`, so its reserved split
cannot see the word), `QUALIFY`, `NOTHING`, `SHORTEST` (missing from the older
grammar's list), and `SAFE_CAST` / `REPLACE_FIELDS` (read only as functions).
The pinned grammar lists all six, and the recording agrees: `SELECT
MATCH_RECOGNIZE`, `SELECT NOTHING`, `SELECT QUALIFY`, `SELECT SHORTEST`, `SELECT
SAFE_CAST` and `SELECT REPLACE_FIELDS` are accepted, each also with a trailing
alias. 12 rows closed. `SELECT EXCEPT` is left: `EXCEPT` is a reserved word, so
the grammar cannot derive that statement, and the emulator answered an empty
select list — an analysis error, not a syntax error — which the judge reads as
an acceptance.

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
required token dropped (`|> ORDER BY`). Of the 1,165 pipe rows, 413 are valid
per the grammar and now agree; the other 752 are the only-Spanner rows left, and
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
