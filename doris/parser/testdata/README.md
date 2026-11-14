# Doris parser test data

## `compat.tsv`: Doris's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/DorisParser.g4`, then of `../../testdata/DorisLegacyParser.g4`
(the legacy CUP parser, converted), each followed by its near-misses, with
omni's verdict and the pinned Doris's (`../../testdata/pins.tsv`).
`TestCompatibility` (`../compat_test.go`) reads it through `enginetest/compat`
and starts no server. Recorded 2026-11-14 in 75 seconds; a second recording on
a fresh server wrote the same file.

| Rows | Total | Agree | Only Doris accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| `DorisParser.g4`, one per alternative | 2,354 | 1,910 | 396 | 48 |
| `DorisParser.g4`, near-misses | 5,961 | 4,668 | 67 | 1,226 |
| `DorisLegacyParser.g4`, one per alternative | 1,767 | 1,463 | 251 | 53 |
| `DorisLegacyParser.g4`, near-misses | 5,885 | 4,481 | 55 | 1,349 |

Most of what only omni accepts is a statement cut short or with a word
dropped that omni's parser reads loosely: a bare `SHOW`, `KILL` or `WARM UP`,
`SHOW <anything>`, a `PAUSE JOB zzz`. 24 of `DorisParser.g4`'s alternatives
get no row: its other entry points, `singleStatement` and
`expressionWithEof`, and `selectHint`, which Doris reads apart from a hint
comment. 7 of `DorisLegacyParser.g4`'s need a token its scanner never makes
(`../../testdata/tokens.tsv`).

## Suites

`suite:conformance` holds the 339 cases of the container test that compared
`Parse` with a running Doris (`TestDorisSyntaxConformance`), each judged by the
same recording; all agree. A statement added to it is a row with `?` in both
verdict cells, judged by a recording:

```text
?	?	-	suite:conformance	<statement>
```

## What the judge needs

- **A statement is sent inside a transaction.** The judge sends `BEGIN` and
  then the statement. Doris parses it with both of its parsers and then
  refuses to run anything but `INSERT`, `UPDATE`, `DELETE` and the
  transaction's own statements, so the refusal is an acceptance and nothing
  else runs. Those four run, on a database that holds no table.
- **`LOAD LABEL` and `CREATE ROW POLICY` are sent outside one.** Doris's
  ANTLR parser builds a command for each that hands the text to the CUP
  parser only when it runs, after a transaction has refused it. Outside a
  transaction a load fails on a table that is not there, and a row policy is
  created, which a later one finds made.
- **A rejection is the parser's own exception** (`compat.DorisVerdict`): the
  ANTLR parser's `ParseException`, `(line 1, pos 7)` after its message, for a
  syntax error and for what its AST builder refuses ("Every derived table
  must have its own alias"), and the CUP parser's `Syntax error in line`.
  Every other error comes after the parse and is an acceptance: a lookup or a
  check of the AST builder ("Unknown privilege type a", "Should contain at
  least one column in a table"), and a DML statement that ran. A verdict
  carries no byte: the CUP parser reports none.
- **The whole text is one request.** Doris splits a script into statements
  itself, so omni's verdict (`omniVerdict`) is `Parse` on the whole row.
