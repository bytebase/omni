# Trino parser test data

## `compat.tsv`: Trino's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/SqlBase.g4`, then its near-misses, with omni's verdict and the
pinned Trino's (`parser` role, judge revision 2 in `../../testdata/pins.tsv`,
at the same `trinodb/trino:482` release as `default`). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no
server. Re-recorded 2026-11-22 at parser pin judge revision 2 and verified
read-only with `OMNI_RECORD=verify` on a fresh container.

| Rows | Total | Agree | Only Trino accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 1,353 | 1,353 | 0 | 0 |
| Near-misses | 3,487 | 3,487 | 0 | 0 |

## Suites

The hand-written corpora that held `Parse` to a live Trino are `suite:` rows,
judged by the same recording. `TestCompatibilitySuitesAreTheCorpora` keeps the
rows equal to the Go corpora they came from (`corpusSuites` in
`../compat_test.go`), which other tests still read.

| Suite | Rows | Agree |
| --- | --- | --- |
| `suite:cte`, `relation`, `select`, `setop`, `window` | 257 | 257 |
| `suite:ddl` | 275 | 275 |
| `suite:dml` | 85 | 85 |
| `suite:dcl-tcl` | 132 | 132 |
| `suite:match-recognize` | 134 | 134 |
| `suite:routines` | 97 | 97 |

All **5,820 / 5,820 (100.00%)** rows agree (`0` only Trino accepts, `0` only
omni accepts).

The other oracle tests compare something else with Trino's verdict (an
expression, a type or a name parsed alone, the lexer, `Split`, SQL omni
generates) and replay tapes (`tapes/`) pinned to the `default` role.

## What the judge needs

- **A statement is run, not prepared.** The judge posts it to
  `/v1/statement`, follows it until it fails or reaches `RUNNING`, and
  cancels it there. The session defaults to the `memory` catalog, so DDL may
  create a table; nothing it leaves changes whether a later statement parses.
- **`SYNTAX_ERROR` is a rejection** (`trinoVerdictOnly`, delegating validation
  to `compat.TrinoVerdict`), read from the error's `errorName`. The parser
  raises it, and so do a few checks after the parse that read only the text,
  which omni's parser makes too: a name with too many parts
  (`ALTER TABLE a.b.c.d RENAME TO e`), an `INTERVAL` qualifier out of order,
  two isolation levels. No generated row reaches those. Every other error, and
  a statement that ran, is an acceptance.
- **Coarse accept/reject verdicts at judge revision 2.** Judge revision 2
  strips `@<byte>` suffixes from `compat.tsv` rejections (matching the other
  engine `compat.tsv` suites) while still validating Trino's error name and
  location bounds through `compat.TrinoVerdict`. Exact `ParseError.Position`
  byte offsets are pinned in `TestErrorPositions` and
  `TestParse_AlignedErrorPositionsAndForms` (`../parse_errors_test.go`) rather
  than degrading omni's diagnostics to match ANTLR's coarser locations.
- **One statement, nothing after it.** The endpoint takes no `;`, so omni's
  verdict (`omniVerdict`) rejects a row with a trailing `;` or an empty text,
  where Trino does.
- **The parser's `AstBuilder` refuses what the grammar derives,** with
  `SYNTAX_ERROR` too: an identifier that starts with a digit or is
  backquoted, a malformed binary or Unicode literal, a control statement other
  than `RETURN` or `BEGIN` as a function's body. Those rows are Trino's
  rejections.

## Error positions (pinned in `TestErrorPositions`)

Under judge revision 1 (`reject@<byte>`), 5,791 / 5,820 rows matched to the
exact byte and 29 rows rejected on both sides at a different byte because
omni's diagnostic position is strictly more informative for callers than
ANTLR's abandoned-prediction or statement-level location. All 29 sharper
positions (plus core SQL syntax, multi-line, multi-byte UTF-8, and lexer error
positions) are permanently pinned in `TestErrorPositions`
(`../parse_errors_test.go`):

| Rows | Cause | Example (omni's byte, Trino's) |
| --- | --- | --- |
| 13 | `IN (…)` (11) and quantified comparison `= ALL(…)` (2) are *alternatives* in `predicate`, told apart only deep inside the group. Trino abandons the whole alternative and reports its first token (`IN`, or the `(` after `ALL`); omni reads into the group and reports the exact token inside it that cannot be taken. | `SELECT 1.5 IN(1.5 zzz)` 18, 11 |
| 9 | A name with too many parts. Trino's check after the parse ("Too many dots in table name") carries **no location**, so the cell was a bare `reject`; omni reports the name's first byte. | `DROP TABLE a.b.c.d` 11, none |
| 3 | A name with too many parts, or `OR REPLACE` with `IF NOT EXISTS`, reported at **byte 0**: Trino's check takes the statement's location. omni reports the name or the `IF NOT EXISTS`. | `CREATE SCHEMA a.b.c` 14, 0 |
| 3 | `<ident>(SELECT …)` where `<ident>` is a non-reserved name (`TABLES`, `STALE`, `ARRAY`), not a comparison quantifier. Trino's prediction stops before the `(`; omni enters the call and reports the `SELECT` where an expression was expected. | `SELECT 1.5 = STALE(SELECT 1.5)` 19, 18 |
| 1 | `SELECT 1.5.*(a)`: omni reads `1.5.*` as `selectAll` and reports the `(` where `AS` is missing; Trino abandons `selectAll` and stops at `*`. | `SELECT 1.5.*(a)` 12, 11 |
