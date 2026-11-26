# CockroachDB parser test data

`compat.tsv` holds a statement per row, with omni's verdict (`cockroachdb/parser`,
delegating through `cockroachdb/analysis` to `github.com/cockroachdb/cockroachdb-parser`
and normalizing failures into `*ParseError`) and the pinned
CockroachDB's (`../../testdata/pins.tsv`, `cockroachdb/cockroach:v25.2.3`).
`TestCompatibility` (`../compat_test.go`) reads it through `enginetest/compat`,
which describes the columns and what passes, and starts no server.

## Rows

| Rule | Rows | From |
|---|---|---|
| Generated alternatives (`Rule/symbols`) | 3,961 | `pkg/sql/parser/sql.y` (`../../testdata/sources.tsv`) and `../../testdata/tokens.tsv` |
| Optional clauses and keyword choices at every site (`Rule/symbols+rule`) | 631 | `(every site)` and `(every alternative)` rows of `../../testdata/tokens.tsv` |
| Generated near-misses (`!drop-N`, `!truncate`, `!insert`, `!swap`) | 11,539 | single-token mutations of the generated alternatives |
| **Total** | **16,131** | |

- **The grammar is fetched, not committed.** CockroachDB's `pkg/sql/parser/sql.y`
  is subject to the CockroachDB Software License (BSL), which is not this repo's
  license. The repo holds only the generated rows and token forms;
  `TestCompatibility` downloads `sql.y` at tag `v25.2.3` once per machine via
  `enginetest.Fetched` and checks its SHA-256 (`../../testdata/sources.tsv`).
- **Token forms** are in `../../testdata/tokens.tsv`, read from CockroachDB
  `v25.2.3`'s scanner (`pkg/sql/scanner/scan.go`), lookahead filter
  (`pkg/sql/parser/lexer.go`), and keyword table (`pkg/sql/lexbase/keywords.go`).
- **330 alternatives get no generated row:** 321 directly contain goyacc's
  built-in `error` recovery token (`SHOW HELP` and unsupported-statement error
  rules), 8 reach rules whose only derivations contain `error`
  (`reindex_stmt`, `alter_unsupported_stmt`, `create_unsupported`,
  `drop_unsupported`, `opt_trusted`, `opt_procedural`), and 1 is
  `const_datetime/DATE` (`date` is not in `GetKeywordID` and lexes as `IDENT`).
  `TestCompatibility` regenerates all 16,131 rows and fails if `compat.tsv`
  holds others.

## Where the verdict comes from

The judge sends `SELECT 1/0;\n` followed by the statement in one PostgreSQL
simple query (`compatJudge`, read by `compat.RawParseVerdict` with `crdbState`).
CockroachDB scans and parses every statement of a simple query before executing
the first, so:

- **SQLSTATE `22012` (division by zero)** means the entire batch parsed and the
  server reached execution of `SELECT 1/0`: **`accept`**. Because `SELECT 1/0`
  fails immediately, nothing of the probed statement ever executes, so DDL,
  DML, transaction control (`BEGIN`/`COMMIT`), `SET`, `USE`, `SHOW`, `DO`,
  `UPSERT`, `IMPORT`, `BACKUP`, and `RESTORE` leave no server state behind.
- **Any other non-operational server error** was raised before `SELECT 1/0`
  executed, so it came from the scanner, goyacc grammar, PL/pgSQL `DO` block
  parser, or a grammar semantic action: **`reject`**. This includes `42601`
  (syntax/lexical errors, help token `??`, removed `RESTORE FROM <uri>` syntax),
  `0A000` (`unimplemented: this syntax`, e.g., `DEFERRABLE` constraints or
  unimplemented routine/attribute clauses), and `22023` (out-of-range type
  modifiers or visibility values checked in grammar actions).
- **Plain `accept` / `reject` (no byte offset in `compat.tsv`).** CockroachDB's
  wire protocol reports parse error locations as a caret diagram in
  `PgError.Detail` while leaving `PgError.Position` at `0`, so the live judge
  emits plain `accept` / `reject` (while `cockroachdb/parser.Parse` extracts
  the byte offset from the error detail into `ParseError.Position` for callers).
- **Operational errors fail the recording.** `nil` error (probe did not fail),
  non-`PgError` transport errors, and operational SQLSTATE classes (`08`
  connection, `53` insufficient resources, `57` operator intervention /
  `statement_timeout`, `58` system error, `XX` internal error) return an error
  rather than a verdict, so they can never be silently recorded as `reject`.
- **Zero bytes** cannot be carried by the PostgreSQL wire protocol's C-string
  query message and fail the judge rather than producing a verdict.

`TestCRDBRawParseVerdict` (`../compat_test.go`) pins every branch of this
verdict classifier offline.

## Counts at `v25.2.3`

| Rows | Total | Agree | Judge only | Omni only | Other |
|---|---|---|---|---|---|
| Alternatives | 4,592 | 4,592 | 0 | 0 | 0 |
| Near-misses | 11,539 | 11,539 | 0 | 0 | 0 |
| **All** | **16,131** | **16,131** | **0** | **0** | **0** |

CockroachDB `v25.2.3` accepts 4,401 of the 4,592 generated alternatives and
rejects 191 in grammar actions or sub-parsers (e.g., `??`, `DO 's'` whose body
is not valid PL/pgSQL, `ALTER INDEX a VISIBILITY 1.5` outside `[0, 1]`,
duplicate options in repeated option lists, and `unimplemented` grammar actions
such as `ALTER TYPE ... ATTRIBUTE` and `CREATE FUNCTION ... TRANSFORM`). It
rejects 10,388 of the 11,539 near-misses and accepts 1,151 whose single-token
mutation still yields valid CockroachDB SQL. `cockroachdb-parser v0.25.2`
agrees with live `cockroachdb/cockroach:v25.2.3` on all 16,131 rows (0
disagreements).

## Commands

From `<repo root>`. Replay offline (no Docker):

```bash
go test ./cockroachdb/parser/ -run '^TestCompatibility$'
```

Judge new rows, or re-record every row after a pin bump in `../../testdata/pins.tsv`:

```bash
OMNI_RECORD=1 go test -timeout 20m ./cockroachdb/parser/ -run '^TestCompatibility$'
```

Verify every committed verdict against a fresh pinned container without writing:

```bash
OMNI_RECORD=verify go test -timeout 20m ./cockroachdb/parser/ -run '^TestCompatibility$'
```

Commit rows a parser change moved:

```bash
go test ./cockroachdb/parser/ -run '^TestCompatibility$' -update-compat
```
