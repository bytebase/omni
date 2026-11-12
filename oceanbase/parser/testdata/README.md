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

What this cannot tell:

- **A row holds one statement.** OceanBase refuses to prepare two before it
  parses either.
- **`?` is a parameter marker,** as `PREPARE` reads it.
- **A refusal on the last reduction can be missed,** when the parser reads the
  tail before it runs that reduction's action. The grammar's actions raise
  almost nothing but 1064, which needs no probe.

### Statements that crash the server

`SET ROLE ALL EXCEPT ALL` and `SET ROLE ALL EXCEPT ALL PRIVILEGES` crash
OceanBase 4.3.5.6 during `PREPARE`: a null dereference in the observer. The
judge starts the observer again, asks the statement again, and, when the
server dies a second time, asks only the probe. Both statements answer the
syntax error at the tail, so they parse and the crash comes after parsing:
their rows say `accept`. Do not send them to a server other tests share.
