# tidb-grammar

Builds the parser.y completeness corpus,
`tidb/parser/testdata/parser_y_corpus.tsv`: for every alternative of every rule
in TiDB's grammar, the shortest statement that exercises it.

This is a nested module so `go test ./...` at the repo root does not reach it,
like the other `harness/` modules. It imports nothing from pingcap or omni: it
reads grammar text and writes SQL text. Only the generated TSV is checked in,
never `parser.y` itself.

## Run

```bash
cd harness/tidb-grammar
go run .                      # fetch parser.y at v8.5.5, rewrite the corpus
go run . -tag v8.5.5          # a specific release tag
go run . -grammar parser.y    # read a local file instead of fetching
```

It writes the statements and their `Rule#alternative` labels with placeholder
verdicts. `tidb/parser`'s `TestTiDBParserYCorpus{Omni,Engine}Verdicts` fill the
verdict columns with `-update`, the same way the pingcap corpus's pair does;
see `tidb/parser/testdata/README.md`.

## How a statement is built

A goyacc grammar derives strings from its start symbol. For one alternative `A`
of rule `N`, the shortest statement exercising it is

    shortest prefix reaching N  +  A with each symbol at its shortest  +  shortest suffix

computed in two passes over the grammar, each to a fixpoint:

- **Shortest derivation** of every symbol: a terminal is its own sample text; a
  nonterminal is the cheapest of its alternatives. (`generate.go`)
- **Shortest context** of every nonterminal: the cheapest prefix/suffix of
  terminals around it in a derivation from the start symbol, filling the
  siblings of each occurrence with their shortest derivations. (`generate.go`)

Everything goyacc adds for the Go parser — the prologue and epilogue, the
`%union`, the semantic actions, `%prec`, the precedence declarations — is
dropped: none of it changes which strings the grammar derives. (`grammar.go`)

A terminal's sample text is the literal it lexes from: a keyword or operator is
written as itself, and the literal classes (`identifier`, `stringLit`,
`intLit`, …) and the lexer's `builtinX` function tokens get a fixed sample in
`specialTermText`. Tokens are kept tight to the `(` and `.` around them so a
function call or a dotted name lexes back the way it was written.

## Limits

A generated sentence is grammatically shortest, not idiomatic, and the server
refuses many on a semantic check its grammar actions make (about 146 of them):
those are honest `reject`/`reject` rows, not noise to remove. An alternative
that no derivation from the start symbol reaches, or that never terminates, is
reported as needing a hand-written sentence; at v8.5.5 there are none.
