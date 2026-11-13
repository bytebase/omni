# Trino parser test data

## `compat.tsv`: Trino's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/SqlBase.g4`, then its near-misses, with omni's verdict and the
pinned Trino's (`../../testdata/pins.tsv`). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no
server. Recorded 2026-11-13 in three minutes; a second recording on a fresh
server wrote the same file.

| Rows | Total | Same cell | Same verdict | Only Trino accepts | Only omni accepts | Both reject, another byte |
| --- | --- | --- | --- | --- | --- | --- |
| One per alternative | 1,353 | 1,195 | 1,202 | 128 | 23 | 7 |
| Near-misses | 3,487 | 2,911 | 3,427 | 55 | 5 | 516 |

## What the judge needs

- **A statement is run, not prepared.** The judge posts it to
  `/v1/statement`, follows it until it fails or reaches `RUNNING`, and
  cancels it there. The session defaults to the `memory` catalog, so DDL may
  create a table; nothing it leaves changes whether a later statement parses.
- **`SYNTAX_ERROR` is a rejection** (`compat.TrinoVerdict`), read from the
  error's `errorName`. The parser raises it, and so do a few checks after the
  parse that read only the text, which omni's parser makes too: a name with
  too many parts (`ALTER TABLE a.b.c.d RENAME TO e`), an `INTERVAL` qualifier
  out of order, two isolation levels. No generated row reaches those. Every
  other error, and a statement that ran, is an acceptance.
- **A rejection carries its byte,** converted from the error's line and
  column: ANTLR's, in code points, lines split at `\n`. A row agrees only
  where omni stops at the same token.
- **One statement, nothing after it.** The endpoint takes no `;`, so omni's
  verdict (`omniVerdict`) rejects a row at its first `;` and an empty text at
  its end, where Trino does.
- **The parser's `AstBuilder` refuses what the grammar derives,** with
  `SYNTAX_ERROR` too: an identifier that starts with a digit or is
  backquoted, a malformed binary or Unicode literal, a control statement other
  than `RETURN` or `BEGIN` as a function's body. Those rows are Trino's
  rejections.
