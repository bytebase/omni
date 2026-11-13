# PartiQL parser test data

## `compat.tsv`

`compat.tsv` holds a statement per row, with omni's verdict and the judge's.
`TestCompatibility` (`../compat_test.go`) reads it through `enginetest/compat`,
which describes the columns and what passes. It starts no server, regenerates
the generated rows from the grammar, and asks the judge again on every run.

| Rows | Total | Agree | Only the judge accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative of the grammar | 632 | 251 | 380 | 1 |
| Near-misses: a token dropped, the end cut off, a word inserted, a keyword swapped | 1,217 | 1,200 | 14 | 3 |
| `suite:dynamodb-docs`: the statements DynamoDB's developer guide prints | 61 | 59 | 2 | 0 |

- **The grammar** is `../../testdata/PartiQLParser.g4` with its lexer grammar
  `PartiQLTokens.g4`: the grammar of the PartiQL reference implementation,
  partiql-lang-kotlin, at release `v0.14.9` (`../../testdata/pins.tsv`). Both
  are Apache-2.0 (`PartiQL.LICENSE`, `PartiQL.NOTICE`) and committed byte for
  byte; their SHA-256 is the pin's digest.
- **Why `v0.14.9`.** It is the last release whose grammar holds every form
  DynamoDB documents. The v1 grammar was rewritten along SQL's EBNF and no
  longer derives `INSERT INTO t VALUE {...}`, an `UPDATE`'s `REMOVE`, or
  `RETURNING`.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to the lexer rules.
  Three rows name a rule: `exprTerm` is written `a`, since the grammar lists
  `CURRENT_USER` before a name and every expression row would carry it, and
  `explainOption` and `patternRestrictor` are written with a word the
  reference parser takes.
- **Two alternatives have no row:** `tableName` and `tableConstraintName`,
  which no rule uses.
- **The DynamoDB rows** are `../aws-corpus/index.json`'s examples, each on one
  line; its two syntax skeletons are left out.

To judge new rows, or every row after a pin bump:

```bash
OMNI_RECORD=1 go test ./partiql/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./partiql/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

PartiQL has no server CI can start, so no row is an engine's answer.

- **The grammar judges every row.** It is the reference language: a statement
  it derives is accepted, any other rejected.
- **DynamoDB takes a part of that language,** and its docs judge nothing as
  rejected. A `JOIN`, a `GROUP BY` or a `CREATE TABLE` is accepted here though
  DynamoDB refuses it: omni parses PartiQL, and what DynamoDB leaves out is
  for a review rule to say.
- **DynamoDB's docs can only add.** A statement the developer guide prints and
  the grammar does not derive would be a `docs` row of
  `../../testdata/grammar-overrides.tsv`, accepted and cited to its page. At
  this pin there is none: the grammar derives all 61.
- **31 statements are the reference parser's, not the grammar's.** The grammar
  derives them and the parser built from it refuses them in the visitor that
  reads the parse tree (`PartiQLParserDefault.kt`): an `EXPLAIN` parameter
  that is not `TYPE` or `FORMAT`, `LAG` with no `ORDER BY`, a `DATE` string
  that is no date. Each is a `grammar` row of `grammar-overrides.tsv`, rejected
  and cited to the line.

## The work list

396 rows the judge accepts and omni rejects:

- **271 are one gap:** an expression as a statement (`dql : expr`,
  `PartiQLParser.g4`), which omni's `Parse` refuses and its expression parser
  takes. Most alternatives of the expression rules are carried by one, so
  closing it is what shows the rest. The two DynamoDB rows are this gap:
  `EXISTS(SELECT ...)`, the condition check of a transaction.
- **125 others:** a graph pattern that begins with an edge (48), tokens after
  a statement (19), `EXCLUDE` (12), `VALUES` and value lists (11), `OVERLAY`
  (7), and smaller ones.

## What the judge cannot decide

- **The rest of the visitor.** The 31 overrides are the checks read off
  `PartiQLParserDefault.kt` for the rows this file holds. Running the
  reference parser on every row would find any other; that needs a JVM.
- **`TRIM(a FROM a)`.** The grammar reads `a` as the trim specification or as
  the characters to trim, and the visitor refuses a specification that is not
  `BOTH`, `LEADING` or `TRAILING`. Which reading ANTLR takes decides the
  statement, and a recognizer has no such preference: the row is accepted.
- **Lexing.** The recognizer reads a statement with a generic SQL lexer, not
  with the rules of `PartiQLTokens.g4`. The generated rows write every token
  apart, so these show only in a hand-added row:
  - `$` starts an identifier (`PartiQLTokens.g4:372`); the recognizer reads
    `$a` as a variable.
  - A backtick opens Ion, which runs to the next backtick outside an Ion
    string or comment (`PartiQLTokens.g4:427`-`451`); the recognizer ends it
    at the next backtick.
  - A string has no backslash escape (`PartiQLTokens.g4:360`); the recognizer
    reads one.
