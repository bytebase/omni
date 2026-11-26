# OceanBase parser test data

## `compat.tsv`: the grammar's rows

`compat.tsv` holds a statement for each alternative of each rule in OceanBase's
`sql_parser_mysql_mode.y` at the pinned release, then its near-misses: the
statement with one token dropped, the end cut off, a stray word inserted, or a
keyword swapped. Each row has omni's verdict and OceanBase's.
`TestCompatibility` (`compat_test.go`) reads it through `enginetest/compat`,
which describes the columns and what passes. It starts no server, and it
regenerates the rows from the grammar and fails if the file holds others.

- **The grammar is fetched, not committed.** OceanBase CE 4.3.5 is under the
  Mulan PubL v2, so the repo holds only the rows. The test downloads the
  grammar once per machine and checks its SHA-256 (`../../testdata/sources.tsv`).
- **Token forms** are in `../../testdata/tokens.tsv`, read off OceanBase's lexer
  and keyword table. A new release's tokens need rows there.
- **Optional clauses are tried at every site.** `tokens.tsv` marks
  `opt_if_not_exists`, `opt_if_exists` and the other optional clauses the
  grammar names at many sites (`opt_equal_mark`, `opt_replace`, `opt_temporary`
  and the rest) `(every site)`, so each alternative naming one gets a row
  labelled `<alternative>+<rule>` with the clause written out, and near-misses
  dropping each of its words.
- **Keyword choices are tried in every form at every site.** A rule marked
  `(every alternative)` (`table_or_tables`, `collation_name`, `opt_limit`
  and the rest) gets a row per form, labelled `<alternative>+<rule>/<form>`.
- **`suite:` rows are hand-added** for forms no generated row writes, such as
  a decimal string length or `COLLATE BINARY`.
- **52 alternatives have no row,** because they need a token the lexer never
  produces: dead grammar, or `||` and `NOT` under a non-default `sql_mode`.

To judge new rows, or every row after a pin bump in `../../testdata/pins.tsv`:

```bash
OMNI_RECORD=1 go test ./oceanbase/parser/ -run '^TestCompatibility$'
```

It boots the pinned image, which takes about 3.5 GiB, and judges all 20,407
rows in about three minutes. After a parser fix, commit the rows that now agree:

```bash
go test ./oceanbase/parser/ -run '^TestCompatibility$' -update-compat
```

### How OceanBase judges a row

The verdict is the server's answer to `PREPARE`, which parses a statement
without running it (`compatJudge`):

- **1064 with the parser's message** ("You have an error in your SQL syntax ...
  near ...") is a rejection.
- **Any other error leaves the question open,** because a statement that
  parses can still be refused: a missing table, an unsupported statement type.
  The judge prepares the statement again with ` @@@ (((` on a new line after
  it. The parser's syntax error at that tail means the statement itself
  parsed. The first error again means the parser never reached the tail.
- **1064 with only "Parse error"** is such an error. It comes from the layers
  after the parser, and from `PREPARE` itself, which does not carry `PREPARE`,
  `EXECUTE` or `DEALLOCATE` though the server takes each when sent directly.
- **1149 ("You have an error in your SQL syntax", with no `near` clause)** can
  come either from a post-parse resolver check (`SELECT RANK() OVER(ROWS 1
  FOLLOWING)`, whose probe fails with 1064 `at line 2` near `@@@ (((`) or from
  a parse-time syntax error (`DO a zzz`, whose probe fails with 1064 `at line
  1` near `a zzz\n @@@ (((`). At judge revision 2 (`PinRole: "parser"` in
  `../../testdata/pins.tsv`), when the first error is 1149, a probe 1064 must
  land on the probe tail's own line (`at line 2` for a single-line statement)
  to count as parsed; an earlier-line 1064 rejects the row.

What this cannot tell:

- **A row holds one statement.** OceanBase refuses to prepare two before it
  parses either.
- **`?` is a parameter marker,** as `PREPARE` reads it.
- **A refusal on the last reduction can be missed,** when the parser reads the
  tail before it runs that reduction's action. The grammar's actions raise
  almost nothing but 1064, which needs no probe.

### Agreement (`20,407 / 20,407`, `100.00%`)

At judge revision 2 (`PinRole: "parser"` in `../../testdata/pins.tsv`), all
`20,407` rows of `compat.tsv` agree (`5,071 / 5,071` alternatives and
`15,336 / 15,336` near-misses, `0` disagreements):

- `stmt/pl_expr_stmt!insert` (`DO a zzz`) is judged `reject` by judge revision
  2 because its `1149` probe fails with `1064` on line 1 rather than on the
  probe tail's line.
- `stmt/pl_expr_stmt` and `pl_expr_stmt/DO_expr` (`DO a`) agree as `reject`
  because `parseDoStmt` matches OceanBase CE `v4.3.5_CE_BP6`'s PL parser
  (`src/pl/parser/pl_parser_mysql_mode.y`, `do_sp_stmt`), which rejects any
  top-level `T_COLUMN_REF` item in `DO expr_list` with `1064` while accepting
  non-column expressions (`DO 1`, `DO 1, 2`, `DO SLEEP(0)`, `DO @x := 1`,
  `DO @@version`).

### Statements that crash the server

`SET ROLE ALL EXCEPT ALL` and `SET ROLE ALL EXCEPT ALL PRIVILEGES` crash
OceanBase 4.3.5.6 during `PREPARE`: a null dereference in the observer. The
judge starts the observer again, asks the statement again, and, when the
server dies a second time, asks only the probe. Both statements answer the
syntax error at the tail, so they parse and the crash comes after parsing:
their rows say `accept`. Do not send them to a server other tests share.

## `tapes/`: the server tests' recordings

Every other test here that asks OceanBase something reaches it through
`testcontainer.Open` and replays `tapes/<Test>.tape`: each request the test
sends, and the pinned server's answer. `enginetest/tape` describes the format
and what fails a replay. After changing what a test sends, record it:

```bash
OMNI_RECORD=1 go test ./oceanbase/parser/ -run '^TestName$'
```

- **The keyword sweeps walk the lexer's keywords** (`TestReservedWordSweepOceanBaseOracle`,
  `TestKeywordFunctionNamesOceanBaseOracle`, `TestSetVariableNameSweepOceanBaseOracle`),
  so a keyword added to the lexer changes what they send: record those three.
- **Tapes are written compactly.** `…near` stands for the 122 characters that
  open each of OceanBase's syntax errors, `^` for a request text repeated from
  the request before, and a `|# sha256:` line for a statement over 64 KB.
- **A verdict is one request,** `PREPARE x FROM ?` (`prepareAcceptsOceanBase`):
  the three sweeps alone ask 36,500 of them.
