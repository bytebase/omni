# grammarcorpus

Builds the parser.y completeness corpus,
`tidb/parser/testdata/parser_y_corpus.tsv`: for every alternative of every rule
in TiDB's grammar, the shortest statement that exercises it.

It is a plain root-module command, like each engine's `genwalker`: `go test
./...` compiles it but never runs it. It imports nothing from pingcap or omni:
it reads grammar text and writes SQL text.

## When to run it

Only when moving to a new TiDB release. Day to day nobody runs it — the
checked-in TSV is the artifact the oracle tests read. You run it when you bump
the pinned grammar: refresh `testdata/parser.y`, rerun, and the TSV diff is the
set of forms the release added or removed.

## The vendored grammar

`testdata/parser.y` is a verbatim copy of `pkg/parser/parser.y` from
`github.com/pingcap/tidb` at tag `v8.5.5` (commit
`1fa258b833ff113883beeba40bc130be7ce66610`), Apache-2.0, license header intact.
It is vendored, not fetched, so regeneration is deterministic and offline and
does not depend on pingcap keeping that raw path alive. It is reference input
for this tool only: nothing imports or compiles it.

## Run

From the repo root:

```bash
go run ./tidb/parser/cmd/grammarcorpus                   # read the vendored parser.y, rewrite the corpus
```

To move to a new release, refresh the vendored copy, then rerun:

```bash
curl -fsSL https://raw.githubusercontent.com/pingcap/tidb/<tag>/pkg/parser/parser.y \
  -o tidb/parser/cmd/grammarcorpus/testdata/parser.y
go run ./tidb/parser/cmd/grammarcorpus -tag <tag>
```

`-grammar <file>` reads any other local file; `-grammar ''` fetches at `-tag`.

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
refuses many on a semantic check its grammar actions make (about 121 of them):
those are honest `reject`/`reject` rows, not noise to remove. An alternative
that no derivation from the start symbol reaches, that never terminates, or that
the lexer only produces under a non-default SQL mode, is reported as needing a
hand-written sentence; at v8.5.5 that is one, `SimpleExpr`'s `||` concat, which
needs `PIPES_AS_CONCAT` (in the default mode `||` lexes as OR).
