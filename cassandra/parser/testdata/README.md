# Cassandra parser test data

## `compat.tsv`: the grammar's rows

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/Cql.g4`, then its near-misses: the statement with one token
dropped, the end cut off, a stray word inserted, or a keyword swapped. Each
row has omni's verdict and the grammar's. `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat`, which describes the
columns and what passes. It starts no server, regenerates the rows from the
grammar, and asks the judge again on every run.

| Rows | Total | Agree | Only the grammar accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 904 | 904 | 0 | 0 |
| Near-misses | 2,037 | 2,036 | 1 | 0 |

- **The grammar** is Apache Cassandra's own, `src/antlr/Cql.g`, `Parser.g` and
  `Lexer.g` at `cassandra-5.0.9` (`../../testdata/pins.tsv`), beside it byte
  for byte under the Apache License (`Cassandra.LICENSE`, `Cassandra.NOTICE`);
  their SHA-256 is the pin's digest.
- **`Cql.g4` is converted from them** by `../cmd/antlr3to4`, since they are
  ANTLR 3. Each rule is under the file and line it came from, and what the
  conversion lost is noted beside it: actions, predicates, and alternatives an
  action always refuses (`COMPACT STORAGE` on a view, a bind marker as a
  keyspace, table, index or role name, a quoted user name).
  `TestCqlG4IsTheConversion` fails while `Cql.g4` is anything else. After a pin
  bump replaces the three `.g` files:

  ```bash
  go run ./cassandra/parser/cmd/antlr3to4
  OMNI_RECORD=1 go test ./cassandra/parser/ -run '^TestCompatibility$'
  ```
- **Token forms** are in `../../testdata/tokens.tsv`, cited to `Lexer.g`.
  Three rows pick what rows write: `cident` is a name, since the rule lists
  `""` first and every column would carry it, and `COLUMNFAMILY | TABLE` and
  `COLUMNFAMILIES | TABLES` are written by their current spelling. The judge
  takes either.

After a parser fix, commit the rows that now agree:

```bash
go test ./cassandra/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

No row is a server's answer: a recognizer over `Cql.g4` accepts a statement the
grammar derives and rejects any other. `../../testdata/grammar-overrides.tsv`
corrects it where an action of `Parser.g` refuses, while it parses, a
statement the grammar derives: a column, property, role option or clustering
order given twice, a column set twice, `FROZEN<>` on a native type, a custom
type string naming no class, a JMX name with no domain. 13 rows, each cited to
the `Parser.g` line and the Java it calls.

## The work list

One row disagrees: `DESCRIBE SCHEMA zzz`, which the grammar derives by reading
`SCHEMA` as `K_KEYSPACE`'s second spelling (`Lexer.g:91-92`) while omni reads
it as `K_SCHEMA`, defined first (`Lexer.g:90`). Which token ANTLR 3's lexer
hands the parser decides the row, and a recognizer cannot tell; an engine
judge would.

## What the judge cannot decide

- **Actions the overrides do not cover.** The 13 were read off `Parser.g` and
  its Java for the rows this file holds; a new row that reaches an action
  `Cql.g4` notes ("action may refuse the input") needs reading again. Running
  Cassandra on the rows would settle them: the image exists, and an `engine`
  judge is the upgrade.
- **What fails only when the statement runs,** as a view with no `WHERE` or a
  table with no primary key: Cassandra's parser takes these, so the grammar
  judge accepts them.
- **Lexing.** The recognizer reads a statement with a generic SQL lexer, not
  `Lexer.g`. The generated rows write every token apart, so these show only in
  a hand-added row:
  - `//` starts a comment (`Lexer.g:394`); the recognizer reads `//`.
  - A duration (`1d`), a UUID and a blob are one token each to Cassandra and
    several to the recognizer, which matches them only as `tokens.tsv` spells
    them.
  - Any number stands for any number token: `0x1` where an `INTEGER` is due.
  - `-1` is one `INTEGER` (`Lexer.g:337`); the recognizer reads `-` and `1`.
