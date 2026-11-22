# MariaDB parser test data

## `compat.tsv`: the grammar's rows

`compat.tsv` holds a statement for each alternative of each rule in MariaDB's
`sql_yacc.yy` at the pinned release, then its near-misses: the statement with
one token dropped, the end cut off, a stray word inserted, or a keyword
swapped, followed by 7 `suite:json-operators` rows pinning the `->` and `->>`
operators MariaDB has no grammar for (12,960 generated rows and 7 `suite:`
rows, 12,967 in all: 12,967 / 12,967 (100.00%) agreeing, 0 only MariaDB
accepts, and 0 only omni accepts).
Each row has omni's verdict and MariaDB's. `TestCompatibility`
(`compat_test.go`) reads it through `enginetest/compat`, which describes the
columns and what passes. It starts no server, and it regenerates the generated
rows from the grammar and fails if the file holds others.

- **The grammar is fetched, not committed.** MariaDB Server is under the GPL
  v2, so the repo holds only the rows. The test downloads the grammar once per
  machine and checks its SHA-256 (`../../testdata/sources.tsv`).
- **One file holds two grammars.** `sql_yacc.yy` carries the default grammar
  and the `sql_mode=ORACLE` one between `%ifdef MARIADB` and `%ifdef ORACLE`
  lines, which MariaDB's build resolves before bison runs. The rows are those
  of `MARIADB`: omni's parser does not read Oracle mode.
- **Token forms** are in `../../testdata/tokens.tsv`, read off MariaDB's
  keyword table (`sql/lex.h`) and lexer (`sql/sql_lex.cc`). A new release's
  tokens need rows there.
- **82 alternatives have no row,** because they need a token the lexer never
  produces under the default `sql_mode`: the `_ORACLE_SYM` tokens, `||` under
  `PIPES_AS_CONCAT`, `NOT` under `HIGH_NOT_PRECEDENCE`, and five keywords that
  no entry of `lex.h` spells.
- **A statement ends with `;`.** The grammar lists `statement ';'` before the
  statement alone, and a row takes the first of two equally short forms.

To judge new rows, or every row after a pin bump in `../../testdata/pins.tsv`:

```bash
OMNI_RECORD=1 go test ./mariadb/parser/ -run '^TestCompatibility$'
```

It starts the pinned image and judges all 12,967 rows in about 20 seconds.
After a parser fix, commit the rows that now agree:

```bash
go test ./mariadb/parser/ -run '^TestCompatibility$' -update-compat
```

### How MariaDB judges a row

The verdict is the server's answer to `PREPARE omni_compat_stmt FROM ?`, which
parses a statement without running it (`compatJudge`, pinned under the `parser`
role at judge revision 2 in `../../testdata/pins.tsv` on `mariadb:11.8.8`):

- **Session setup:** The judge selects database `a` (`CREATE DATABASE IF NOT
  EXISTS a; USE a`) so generated statements that qualify `a.a` against an
  unqualified table `a` (such as `DELETE a.a FROM a;` and
  `ALTER TABLE a DROP a.a.a;`) resolve to the same database instead of failing
  on current-database mismatch, and creates a key cache
  (`SET GLOBAL omni_cache.key_buffer_size = 1048576`) for structured system
  variable stand-ins.
- **1064** is a rejection.
- **Parse-time lookup errors are re-asked with a stand-in.** MariaDB's grammar
  actions look names up while the parser runs and abort before reading later
  tokens: a character set (1115), a plugin or loadable function (1128), a
  system variable (1193, 1272, 4081), a collation (1253, 1273), a storage
  engine (1286), and `SHOW ... CODE` in non-debug builds (1289). Whenever
  `PREPARE` answers one of these, `compatJudge` substitutes a known-valid
  MariaDB 11.8.8 stand-in of the same syntactic category (`utf8mb4`,
  `utf8mb4_bin`, `InnoDB`, `autocommit`, `max_statement_time`,
  `omni_cache.key_buffer_size`, `mariadb_schema`, `COMMENT 's'`, or
  `SHOW CREATE ...`) for the offending lookup and re-prepares so MariaDB parses
  past the lookup and rejects any trailing syntax error (e.g. `zzz`) with 1064.
  All 163 lookup-error rows are re-asked with a stand-in (0 blind lookup
  acceptances). `prepareParses`, which the hand-written cases of this package
  ask, reads lookup errors directly.
- **Any other error leaves the question open,** because a statement that
  parses can still be refused: a missing table (1146), an unknown column
  (1054). The judge prepares the statement again with `; SELECT FROM` on a new
  line after it. A 1064 then means the statement itself parsed; the first
  error again means the parser stopped before the tail, and that is a
  rejection.
- **`PREPARE` carries every statement kind but its own.** Measured on the
  pinned server: each of the 59 alternatives of `verb_clause` and `statement`
  gets an answer the rules above read. `PREPARE`, `EXECUTE`,
  `EXECUTE IMMEDIATE` and `DEALLOCATE PREPARE` answer 1295 (not supported in
  the prepared statement protocol), and the probe decides them.

What this cannot tell:

- **A refusal made by an action is a rejection,** though the text follows the
  grammar. The commonest on the pinned server: an identifier as a data type
  (4161, 731 rows, mostly the keywords the grammar lets stand for a type
  name), a name used as a routine variable outside a routine (1327, 315 rows),
  a table named twice (1066, 115 rows). omni accepts some of these, and their
  rows say so.
- **`?` is a parameter marker,** as `PREPARE` reads it.
- **A row holds one statement.**

## `tapes/`: the server tests' recordings

Every other test here that asks MariaDB something reaches it through
`testcontainer.OpenEmptyWith` (`startMariaDB`) and replays
`tapes/<Test>.tape`: each request the test sends, and the pinned server's
answer. `enginetest/tape` describes the format and what fails a replay. After
changing what a test sends, record it:

```bash
OMNI_RECORD=1 go test ./mariadb/parser/ -run '^TestName$'
```

- **A recording starts each test on an empty server.** Before a test's
  connection opens, every database, account and foreign server the image did
  not start with is dropped and `test` is created again, on a connection the
  tape does not see.
- **The keyword sweeps walk the lexer's keywords**
  (`TestReservedWordSweepOracle`, `TestKeywordFunctionNamesOracle`,
  `TestSetVariableNameSweepOracle`), so a keyword added to the lexer changes
  what they send: record those three.
- **`TestOperatorSweepMatchesEngine` sends the grouping `Parse` gave** each
  expression, so a change to how an operator expression is read needs a
  recording. Its size and seed are constants beside it.
- **Tapes are written compactly.** `…near` stands for the 124 characters that
  open each of MariaDB's syntax errors, `^` for a request text repeated from
  the request before, and a `|# sha256:` line for a statement over 64 KB.
- **A verdict is one request,** `PREPARE omni_parses_stmt FROM ?`
  (`prepareParses`), or two when the probe decides it.
