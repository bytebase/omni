# MariaDB parser test data

## `compat.tsv`: the grammar's rows

`compat.tsv` holds a statement for each alternative of each rule in MariaDB's
`sql_yacc.yy` at the pinned release, then its near-misses: the statement with
one token dropped, the end cut off, a stray word inserted, or a keyword
swapped. Each row has omni's verdict and MariaDB's. `TestCompatibility`
(`compat_test.go`) reads it through `enginetest/compat`, which describes the
columns and what passes. It starts no server, and it regenerates the rows from
the grammar and fails if the file holds others.

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

It starts the pinned image and judges all 12,960 rows in about 20 seconds.
After a parser fix, commit the rows that now agree:

```bash
go test ./mariadb/parser/ -run '^TestCompatibility$' -update-compat
```

### How MariaDB judges a row

The verdict is the server's answer to `PREPARE omni_compat_stmt FROM ?`, which
parses a statement without running it (`compatJudge`):

- **1064** is a rejection.
- **A name the server does not have counts as read.** MariaDB's grammar
  actions look names up while the parser runs: a character set (1115), a
  plugin or loadable function (1128), a system variable (1193, 1272, 4081), a
  collation (1253, 1273), a storage engine (1286), and a form the build
  leaves out (1289). They are about what the server holds, and omni keeps no
  such list. `prepareParses`, which the hand-written cases of this package
  ask, reads them the same way.
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

- **A name the server does not have answers before a later syntax error.**
  `SET GLOBAL a = a zzz` answers 1193, not 1064, and reads as accepted. Where
  the grammar gives such a name a rule of its own, `tokens.tsv` writes a name
  that exists (`charset_name` is `utf8mb4`), and the rows that pass through
  the rule are judged on their syntax. 163 rows are still answered by such an
  error: the alternatives that are the lookup itself (`SET NAMES a`), and the
  `SET GLOBAL`, `SET STATEMENT` and `@@DEFAULT.` forms, whose variable is any
  identifier. A near-miss among them may be wrong. The recording logs the
  list.
- **A refusal made by an action is a rejection,** though the text follows the
  grammar. The commonest on the pinned server: an identifier as a data type
  (4161, 731 rows, mostly the keywords the grammar lets stand for a type
  name), a name used as a routine variable outside a routine (1327, 315 rows),
  a table named twice (1066, 115 rows). omni accepts some of these, and their
  rows say so.
- **`?` is a parameter marker,** as `PREPARE` reads it.
- **A row holds one statement.**
