# StarRocks parser test data

## `compat.tsv`: StarRocks's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/StarRocks.g4`, then its near-misses, with omni's verdict and
the pinned StarRocks's (`../../testdata/pins.tsv`). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no
server. Recorded 2026-11-14 in 45 seconds; a second recording on a fresh
server wrote the same file.

| Rows | Total | Agree | Only StarRocks accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 3,062 | 3,062 | 0 | 0 |
| Near-misses | 8,196 | 8,196 | 0 | 0 |
| An optional clause at every site, with near-misses | 29 | 29 | 0 | 0 |

Every row agrees: the 11,287 generated rows and the 914 `suite:` rows below,
12,201 in all. #1467 closed the last 203, 202 generated rows and one `suite:`
row. The generated ones were mostly statements omni did not read at all:
`PIPE`, resource groups, `CATALOG`, data cache rules, `TRANSLATE`, the SQL and
backend blacklists, plan advisors, `REFRESH EXTERNAL TABLE`, `EXECUTE AS` and
`RESTORE SNAPSHOT`; the rest were forms StarRocks refuses and omni accepted, a
reserved word read as a name (`SELECT FLOOR`), `SELECT COUNT(DISTINCT)`,
`GROUP BY ROLLUP()`, `SHOW … WHERE`, and `CREATE CATALOG a` with no
`PROPERTIES`. 3
alternatives get no row: `expressionSingleton`, another entry point,
`hintMap`, which StarRocks reads apart from a hint comment, and `CONCAT`,
which `||` is only under sql_mode `PIPES_AS_CONCAT`.

The optional clauses `../../testdata/tokens.tsv` names with `(every site)`
get a row at each alternative that can write them. `StarRocks.g4` gives most
statements a rule of its own, so most of its `(IF NOT EXISTS)?` groups have
one site and their row already; the rows are for `ifNotExists`, `orReplace`
and the partition and `ALTER USER` clauses written at more than one. Their 29
rows agreed when first recorded, 2026-11-23.

## Suites

Eight container tests compared `Parse` with a running StarRocks. Their 434
cases are `suite:` rows named after them, judged by the same recording; all 434
agree: `suite:conformance` (281, `TestStarRocksSyntaxConformance`),
`suite:explain-modifiers` (18), `suite:no-job-statements` (46), and
`suite:pr2a`, `pr2b`, `pr3`, `pr4a` and `pr4b` (89, the parity matrices of
the PRs that added those forms). The 15 `suite:characters` rows, the
characters no rule of `StarRocksLex.g4` matches, are a group of their own; the
one that disagreed, `SELECT 1　AS x` with an ideographic space, agrees now that
omni reads the character as one StarRocks refuses (#1467). The other 465 rows
are 7 `suite:comments`, 5 `suite:copy-into`, 396 `suite:keywords` and 57
`suite:doris-residue`; all agree. A statement added
to one is a row with `?` in
both verdict cells, judged by a recording:

```text
?	?	-	suite:conformance	<statement>
```

## What the judge needs

- **A statement is sent behind a probe,** as one request:
  `SELECT omni_compat_probe();` and then the statement. StarRocks parses a
  request whole, building every statement's AST, before it runs the first,
  and stops at the first that fails. The probe parses and fails when it runs,
  so its error is an acceptance and nothing of the statement runs
  (`compat.StarRocksVerdict`).
- **Every other error is the parse's, and a rejection:** a syntax error
  (`Getting syntax error at line 2, column 10`), and what the AST builder
  refuses (`not support.`) or fails on (a `ClassCastException`'s message, as
  for `SELECT FLOOR`). A verdict carries no byte: the AST builder's errors
  report none.
- **The whole text is one request.** StarRocks splits a script into
  statements itself, so omni's verdict (`omniVerdict`) is `Parse` on the whole
  row.
