# Cosmos DB parser test data

## `compat.tsv`: the grammar's rows

`compat.tsv` holds a statement for each alternative of each rule in `sql.g4`,
then its near-misses: the statement with one token dropped, the end cut off, a
stray word inserted, or a keyword swapped. Each row has omni's verdict and the
judge's. `TestCompatibility` (`../compat_test.go`) reads it through
`enginetest/compat`, which describes the columns and what passes. It starts no
server, regenerates the rows from the grammar, and asks the judge again on
every run.

| Rows | Total | Agree | Only the judge accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 163 | 158 | 5 | 0 |
| Near-misses | 287 | 273 | 3 | 11 |

- **The grammar** is `../../testdata/sql.g4`: Microsoft's grammar of the query
  language, as the .NET SDK publishes it at release `3.64.0`
  (`../../testdata/pins.tsv`). It is MIT-licensed (`sql.g4.LICENSE`) and
  committed byte for byte; its SHA-256 is the pin's digest.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to its lexer rules.

To judge new rows, or every row after a pin bump:

```bash
OMNI_RECORD=1 go test ./cosmosdb/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./cosmosdb/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

Cosmos DB has no server CI can start, so no row is the engine's answer.

- **`grammar`, 386 rows.** A recognizer over `sql.g4`: a statement the grammar
  derives is accepted, any other rejected.
- **`docs`, 64 rows.** `sql.g4` is the SDK's grammar, not the service's parser,
  and Microsoft's query reference takes more. The SELECT page writes
  `<scalar_expression> [ [ AS ] property_alias ]`, where `sql.g4:22` requires
  `AS`. Each row the grammar rejects and derives once `AS` is written is in
  `../../testdata/grammar-overrides.tsv`, accepted and cited to the page.

## What the judge cannot decide

- **What the docs take and no row holds.** The docs' own examples write an
  object's property names bare, `{ name: p.name }`; `sql.g4:199` takes only a
  string there. A hand-added row for one needs an override.
- **Lexing.** The recognizer reads a statement with a generic SQL lexer, not
  with the lexer rules of `sql.g4`. The generated rows write every token apart
  and hold none of these, so a hand-added row is where they would show:
  - A double-quoted string is read as a quoted name; `sql.g4:270` makes it a
    string.
  - `false`, `null`, `true`, `udf` and `undefined` are lower-case only
    (`sql.g4:224`-`244`), and any other case is an identifier. The recognizer
    reads a word in any case.
  - A number takes its sign (`sql.g4:265`), so `1-1` is two numbers to the
    grammar's lexer and a subtraction to the recognizer.
  - `sql.g4` has no comment rule; the recognizer skips `--` and `/* */`.
- **What is checked after parsing.** `SELECT *` needs exactly one alias in
  `FROM`; the grammar derives it with none.
