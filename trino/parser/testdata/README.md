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
| One per alternative | 1,353 | 1,349 | 1,353 | 0 | 0 | 4 |
| Near-misses | 3,487 | 3,432 | 3,487 | 0 | 0 | 55 |

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
| `suite:match-recognize` | 134 | 132 | 134 |
| `suite:routines` | 97 | 96 | 97 |

Every row agrees on the verdict, and 5,746 agree on the cell; the 74 that do
not reject on both sides at a different byte. *What still differs, and why*
below says where each one stops and why omni keeps its own byte.

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

74 rows reject on both sides at a different byte. omni reports the first token
the grammar cannot take: the token its recursive descent stopped at. Trino's
byte is the token ANTLR's prediction stopped at, which is that token only when
prediction and descent agree. These are the causes, largest first, with the row
that shows each.

| Rows | Cause | Example (omni's byte, Trino's) |
| --- | --- | --- |
| 15 | `IN (…)` and a quantified comparison are *alternatives* in `predicate`, told apart only deep inside the group. Trino abandons the whole alternative and reports its first token (`IN`, or the `(` after `ALL`); omni reads into the group and reports the token that cannot be taken. | `SELECT 1.5 IN(1.5 zzz)` 18, 11 |
| 15 | A row pattern: `MATCH_RECOGNIZE (…)`'s leading clauses, `MEASURES`, and `AFTER MATCH SKIP TO` are all optional, and Trino's prediction for them wanders into the row pattern and stops there ("Expecting: ')', ','"). omni reports the first token of the clause it cannot take. | `SELECT 1.5 FROM a MATCH_RECOGNIZE(zzz PATTERN(a) DEFINE a AS 1.5)` 34, 38 |
| 9 | A name with too many parts. Trino's check after the parse ("Too many dots in table name") carries **no location**, so the cell is a bare `reject`; omni reports the name's first byte. A byte cannot move to nothing without dropping the position a caller reports. | `DROP TABLE a.b.c.d` 11, none |
| 4 | `EXTRACT`'s field and a `U&'…' UESCAPE '…'`'s escape character. Both are checks after the parse, and Trino reports them at the node they refuse (the `EXTRACT`, the literal); omni reports the field or the escape clause. | `SELECT EXTRACT(a FROM 1.5)` 15, 7 |
| 4 | `AS <function definition>`: `functionSpecification`'s `AS DOLLAR_STRING` body. The lexer has no `DOLLAR_STRING` token, so omni stops at the `AS` and Trino at the string. Supporting the form is a lexer token and a body kind, not a byte. | `CREATE FUNCTION a() RETURNS TIMESTAMP AS $$a$$` 38, 41 |
| 3 | `frameBound`'s `UNBOUNDED PRECEDING`: Trino reports the closing `)` of the frame, "Expecting: 'AND'"; omni reports the `AND` where the first `UNBOUNDED` needed `PRECEDING` or `FOLLOWING`. | `SELECT a OVER(RANGE BETWEEN UNBOUNDED AND UNBOUNDED PRECEDING)` 38, 61 |
| 3 | A name with too many parts, or `OR REPLACE` with `IF NOT EXISTS`, reported at **byte 0**: Trino's check takes the statement's location. omni reports the name or the `IF NOT EXISTS`. | `CREATE SCHEMA a.b.c` 14, 0 |
| 2 | A `NOT` that begins a negated predicate (`(NOT)? (BETWEEN\|IN\|LIKE)`) but is followed by none of the three: Trino's prediction walks the negated-predicate path and reports the token after the `NOT`. omni leaves the `NOT` to its caller, which reports the `NOT` — the reading that keeps `a varchar DEFAULT 'x' NOT NULL` a column constraint. | `SELECT 1.5 NOT zzz LIKE 1.5` 11, 15 |
| 2 | A `;` inside a routine body. The endpoint takes no `;`, so Trino's first stray one is its byte; omni parses a script and reports the error of the segment after it. | `CREATE FUNCTION a() RETURNS TIMESTAMP RETURN 1.5; END LOOP` 50, 48 |
| 17 | One row each; see below. | |

The 17 one-off rows:

| Rows | Omni's byte, Trino's | Cause |
| --- | --- | --- |
| `CREATE FUNCTION a(RETURNS TIMESTAMP RETURN 1.5`, `WITH FUNCTION a(RETURNS …`, `WITH FUNCTION a() …, FUNCTION a(RETURNS …`, `DROP FUNCTION a(a zzz TIMESTAMP)` | 18/16/59/16, 36/34/77/22 | `parameterDeclaration`'s `identifier? type`: `RETURNS` is non-reserved in Trino's grammar, so Trino reads it (or `a zzz`) as a parameter and stops at the `RETURN` or `TIMESTAMP` after it; omni's keyword table has `RETURNS` reserved. A keyword-table change, not a byte. |
| `WITH FUNCTION() RETURNS …` | 13, 14 | The function name is missing; Trino stops at the `)`, omni at the `(`. |
| `SELECT 1.5.*(a)` | 12, 11 | Trino refuses `.*` on anything but a relation reference; omni accepts it and stops at the `(`. |
| `SELECT 1.5 FROM TABLE(a(TABLES(SELECT 1.5)))`, `… TABLE(a(COPARTITION(a, a), a, a))`, `… TABLE(a(TABLE(a) ORDER BY()))` | 31/43/42, 30/49/43 | The table-argument readings: Trino's prediction resolves them differently (one stops inside a row pattern's `->`), and one stops at the `TABLES`. |
| `SELECT LISTAGG(1.5,) WITHIN GROUP(ORDER BY 1.5)` | 21, 19 | The `,` separator expects a string; omni's LISTAGG reading and its function-call reading both fail, and the further one is past the `)`. |
| `SELECT JSON_OBJECT(KEY zzz 1.5 VALUE 1.5)`, `… (RETURNING(TIMESTAMP))`, `… (RETURNING < TIMESTAMP >)`, `… (RETURNING < TIMESTAMP, TIMESTAMP >)`, `… (RETURNING ROW(a TIMESTAMP zzz))` | 23/28/29/29/35, 27/39/41/40/45 | `KEY` and `RETURNING` are non-reserved, so the grammar also reads them as a member's key expression. Trino takes that reading when the clause does not parse and stops in the member; omni commits to the clause. |
| `SELECT U&'a' UESCAPE` | 13, 20 | The `UESCAPE` clause is lexed as part of the Unicode literal, so a missing string is reported at the clause rather than at the end of input. |
| `SELECT * FROM (WITH FUNCTION f(x int) RETURNS int RETURN x SELECT f(1) AS c) t` | 15, 29 | `WITH FUNCTION` is not allowed in a subquery, and omni rejects it at the `WITH`; Trino reads the `WITH` as a CTE list and stops later. |
