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
| One per alternative | 3,062 | 2,184 | 720 | 158 |
| Near-misses | 8,196 | 5,606 | 152 | 2,438 |

Most of what only omni accepts is a statement cut short or with a word
dropped that omni's parser reads loosely: a bare `SHOW` and `SHOW <anything>`,
`ALTER TABLE` clauses (here behind `ALTER SYSTEM`, the shortest statement
that reaches them), partition and column definitions. Most of what only
StarRocks accepts is a non-reserved keyword as a name (`SELECT ADMIN`),
`USE 'a'`, `COLLATE`, privileges, and `LOAD` and `BACKUP` forms. 3
alternatives get no row: `expressionSingleton`, another entry point,
`hintMap`, which StarRocks reads apart from a hint comment, and `CONCAT`,
which `||` is only under sql_mode `PIPES_AS_CONCAT`.

## Suites

Eight container tests compared `Parse` with a running StarRocks. Their 434
cases are `suite:` rows named after them, judged by the same recording; all
agree: `suite:conformance` (281, `TestStarRocksSyntaxConformance`),
`suite:explain-modifiers` (18), `suite:no-job-statements` (46), and
`suite:pr2a`, `pr2b`, `pr3`, `pr4a` and `pr4b` (89, the parity matrices of
the PRs that added those forms). A statement added to one is a row with `?` in
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
