# MySQL parser test data

## `compat.tsv`: the grammar's rows

`compat.tsv` holds a statement for each alternative of each rule in MySQL's
`sql/sql_yacc.yy` at the pinned release, then its near-misses: the statement
with one token dropped, the end cut off, a stray word inserted, or a keyword
swapped. Each row has omni's verdict and MySQL's. `TestCompatibility`
(`compat_test.go`) reads it through `enginetest/compat`, which describes the
columns and what passes. It starts no server, and it regenerates the rows from
the grammar and fails if the file holds others.

- **The grammar is fetched, not committed.** MySQL is under the GPL v2, so the
  repo holds only the rows. The test downloads the grammar once per machine
  and checks its SHA-256 (`../../testdata/sources.tsv`).
- **Token forms** are in `../../testdata/tokens.tsv`, read off MySQL's keyword
  table (`sql/lex.h`) and lexer (`sql/sql_lex.cc`). A new release's tokens need
  rows there.
- **Names are `utf8mb4`,** a character set, because MySQL looks a character set
  up while it parses (below).
- **A statement ends in `;`.** `sql_statement` lists the terminated form first,
  and ties follow the grammar's order. `PREPARE` takes it.
- **10 alternatives have no row,** because they need a token the lexer never
  produces under the default `sql_mode`: the five `GRAMMAR_SELECTOR_*` the
  server sets before it parses a stored expression, `NOT2_SYM`, `OR_OR_SYM`,
  and `LEX_HOSTNAME`, which the generator cannot write tight against a user
  name.
- **Optimizer hints have no rows.** The lexer hands them to a grammar of their
  own, `sql/sql_hints.yy`, which this file does not read.

To judge new rows, or every row after a pin bump in `../../testdata/pins.tsv`:

```bash
OMNI_RECORD=1 go test ./mysql/parser/ -run '^TestCompatibility$'
```

It starts the pinned image and judges all 12,282 rows in about 15 seconds.
After a parser fix, commit the rows that now agree:

```bash
go test ./mysql/parser/ -run '^TestCompatibility$' -update-compat
```

### How MySQL judges a row

The verdict is the server's answer to `PREPARE omni_compat FROM ?`, which
parses a statement without running it (`compatJudge`), read by
`compat.PrepareVerdict`:

- **1064** is a rejection.
- **Any other error leaves the question open,** because a statement that
  parses can still be refused: a missing table, a statement `PREPARE` does not
  carry (1295). The judge prepares the statement again with `\n; SELECT FROM`
  after it. A 1064 at that tail means the statement itself parsed. The first
  error again means a grammar action refused it before the parser reached the
  tail, which is a rejection: `GET DIAGNOSTICS` outside a routine, `COMMIT` in
  a function, `ALTER EVENT` with no clause.
- **1115, 1193, 1253 and 1273 are acceptances,** without the probe: an unknown
  character set, system variable or collation, looked up while parsing.

What this cannot tell:

- **A lookup hides what follows it.** It ends the parse, so a near-miss broken
  after the name reads as accepted. In the recording 114 rows answered a
  lookup: 71 a system variable (`@@utf8mb4`, `SET utf8mb4 = 1`), 40 a collation
  (`COLLATE utf8mb4`) and 3 a keyword written as a character set; 29 of them
  are near-misses. Writing names as `a` would add 34 near-misses behind
  `CHARSET a`.
- **`?` is a parameter marker,** as `PREPARE` reads it.

## `tapes/`: the server tests' recordings

Every other test here that asks MySQL something reaches it through
`testcontainer.OpenDatabase` and replays `tapes/<Test>.tape`: each request the
test sends, and the pinned server's answer. `enginetest/tape` describes the
format and what fails a replay. After changing what a test sends, record it:

```bash
OMNI_RECORD=1 go test ./mysql/parser/ -run '^TestName$'
```

- **The keyword sweeps walk the lexer's keywords** (`TestKeywordFunctionNamesOracle`,
  `TestSetVariableNameSweepOracle`), so a keyword added to the lexer changes what
  they send: record those two.
- **The operator sweep compares the server with omni's tree** (`TestOperatorSweepMatchesEngine`):
  a parser change that groups an expression otherwise changes what it sends.
  Its size and seed are constants beside it.
- **A verdict is one request,** `PREPARE x FROM ?`, and a tape writes `…near`
  for the 134 characters that open each syntax error.
- **The account reaches only its own database,** so a corpus that runs
  `SHUTDOWN` or `SET GLOBAL` changes nothing for the tests recorded after it.
