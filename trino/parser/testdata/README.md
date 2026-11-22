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
| One per alternative | 1,353 | 1,353 | 1,353 | 0 | 0 | 0 |
| Near-misses | 3,487 | 3,470 | 3,487 | 0 | 0 | 17 |

## Suites

The hand-written corpora that held `Parse` to a live Trino are `suite:` rows,
judged by the same recording. `TestCompatibilitySuitesAreTheCorpora` keeps the
rows equal to the Go corpora they came from (`corpusSuites` in
`../compat_test.go`), which other tests still read.

| Suite | Rows | Same cell | Same verdict |
| --- | --- | --- | --- |
| `suite:cte`, `relation`, `select`, `setop`, `window` | 257 | 257 | 257 |
| `suite:ddl` | 275 | 267 | 275 |
| `suite:dml` | 85 | 81 | 85 |
| `suite:dcl-tcl` | 132 | 132 | 132 |
| `suite:match-recognize` | 134 | 134 | 134 |
| `suite:routines` | 97 | 97 | 97 |

Every row agrees on the verdict, and 5,791 agree on the cell (99.50%); the 29
that do not reject on both sides at a different byte. *What still differs, and
why* below says where each one stops and why omni keeps its own byte.

The other oracle tests compare something else with Trino's verdict (an
expression, a type or a name parsed alone, the lexer, `Split`, SQL omni
generates) and replay tapes (`tapes/`).

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

## What still differs, and why

29 rows reject on both sides at a different byte (down from 74 after adding
`DOLLAR_STRING` / `AS $$…$$` function definitions, `UESCAPE` literal parsing,
and non-reserved/further-error alignments). Every remaining difference is a
case where omni's diagnostic position is strictly more informative for callers
than ANTLR's abandoned-prediction or statement-level location:

| Rows | Cause | Example (omni's byte, Trino's) |
| --- | --- | --- |
| 13 | `IN (…)` (11) and quantified comparison `= ALL(…)` (2) are *alternatives* in `predicate`, told apart only deep inside the group. Trino abandons the whole alternative and reports its first token (`IN`, or the `(` after `ALL`); omni reads into the group and reports the exact token inside it that cannot be taken. | `SELECT 1.5 IN(1.5 zzz)` 18, 11 |
| 9 | A name with too many parts. Trino's check after the parse ("Too many dots in table name") carries **no location**, so the cell is a bare `reject`; omni reports the name's first byte. A byte cannot move to nothing without dropping the position a caller reports. | `DROP TABLE a.b.c.d` 11, none |
| 3 | A name with too many parts, or `OR REPLACE` with `IF NOT EXISTS`, reported at **byte 0**: Trino's check takes the statement's location. omni reports the name or the `IF NOT EXISTS`. | `CREATE SCHEMA a.b.c` 14, 0 |
| 3 | `<ident>(SELECT …)` where `<ident>` is a non-reserved name (`TABLES`, `STALE`, `ARRAY`), not a comparison quantifier. Trino's prediction stops before the `(`; omni enters the call and reports the `SELECT` where an expression was expected. | `SELECT 1.5 = STALE(SELECT 1.5)` 19, 18 |
| 1 | `SELECT 1.5.*(a)`: omni reads `1.5.*` as `selectAll` and reports the `(` where `AS` is missing; Trino abandons `selectAll` and stops at `*`. | `SELECT 1.5.*(a)` 12, 11 |
