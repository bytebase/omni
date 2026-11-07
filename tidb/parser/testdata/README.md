# TiDB parser corpus

`pingcap_parser_corpus.tsv` holds every single-statement case from pingcap's own
parser test, with the verdict `omni/tidb/parser` and TiDB v8.5.5 each give it.
It is the measure of how much of the TiDB grammar omni covers, and the gate the
PRs that close that gap work down. `corpus_oracle_test.go` is its only reader;
`docs/plans/2026-11-03-tidb-grammar-gap-inventory.md` groups what is missing.

## Columns

Tab-separated, three fields, one statement per line:

| Field | Values | Means |
|---|---|---|
| 1 | `accept`, `reject` | what `tidb/parser.Parse` does with the statement |
| 2 | `accept`, `reject` | whether TiDB v8.5.5 accepts the syntax |
| 3 | the statement | verbatim, one line |

Lines beginning with `#` are the file's own header and are written back
unchanged by an `-update` run.

Closing a grammar gap flips column 1 from `reject` to `accept`. Column 2 is not
a knob: it records what the engine does, and a disagreement with it is the
finding, not the thing to silence. It may be rewritten only when the pinned
image changes, and then the diff is the list of what the new release altered, or when
the reading of the engine's answer below is corrected, and then the diff is the list
of what the old reading misreported.

Some rows are deliberately awkward and must survive an edit of this file: five
are not valid UTF-8 (pingcap's GBK and big5 cases), sixteen end in a space
(truncated-input cases), and four carry a CR inside a string literal. The
loader splits on `\n` rather than scanning lines so none of that is stripped.

## Provenance

Extracted from `parser_test.go` in `github.com/pingcap/tidb/pkg/parser`, module
version `v0.0.0-20251104011144-749c973e5a70` (commit `749c973e5a70`). Each case
there is a statement with a flag saying whether pingcap's parser should accept
it; only the statement text was taken, and only from cases that fit on one line.
The flag was not: it describes pingcap's parser on master, not the pinned
server, and 109 statements it accepts are grammar TiDB v8.5.5 does not have.
Column 2 asks the server instead.

### Licensing

`github.com/pingcap/tidb` is distributed under the Apache License 2.0
(Copyright 2016 PingCAP, Inc.). The license permits redistribution with
attribution; this file is that attribution. The rows are test fixtures read at
test time and are not compiled into omni, and omni does not depend on pingcap:
nothing here adds a module requirement.

## Where the engine verdict comes from

`PREPARE`, which parses the statement and stops. Executing the corpus cannot
measure it, because statements in the corpus destroy the run that executes
them:

- `shutdown` succeeds and stops the server. Nothing after it gets a verdict and
  no reconnect reaches a process that has exited.
- `SET PASSWORD = 'password';` succeeds and changes the connecting user's
  password. The already-authenticated connection survives, so the run looks
  healthy, but every new connection is refused and the run dies as soon as
  anything replaces it — which `kill CONNECTION_ID()`, later in the corpus,
  then forces.
- `load stats '/tmp/stats.json'` and the `PLAN REPLAYER` file forms make the
  server ask the client for a local file, which the driver refuses with an error
  that is not a server verdict at all.

The reading is the same as everywhere else in this package. Error 1064 is the
parser's rejection, and it is not the only one: TiDB's grammar actions refuse text
their productions read, with the code each chooses, while the parse is still running.
`CREATE VECTOR INDEX i ON nosuch (a)` answers 1149 and `CREATE VECTOR INDEX i ON
nosuch ((a))` answers 1146. Any other server error came after the parse: 1295 for
a statement the prepared-statement protocol does not carry, 1146 for a table that
does not exist. Anything that is not a server error is no verdict and fails the test.

A code does not say which side it is on (1105 is both), so the harness asks. It
prepares the statement again with a syntax error after it as a second statement. The
server reads all of the text before it looks at anything, so a statement refused after
parsing now answers 1064, and a statement an action refused answers its own code again.

One class counts as read even though it is raised while parsing: a name the server looks
up, a system variable (1193), a character set (1115) or a collation (1273, 1253). `Parse`
holds no such list, so those are not grammar it can follow.

Checked against direct execution over 3,897 of the 3,902 statements: the two
agree on every one except the two rows that are themselves `PREPARE`
statements, where nesting answers 1295 rather than judging the inner text. For
those two `PREPARE` is the more accurate syntax verdict — direct execution's
1064 there is about the runtime value of a user variable, not about the
statement. Execution lost five verdicts over the same range; `PREPARE` lost
none.

## Regenerating

From `<repo root>`:

```bash
go test ./tidb/parser/ -run '^TestTiDBCorpusOmniVerdicts$' -count=1 -v -update
```

rewrites column 1 from the parser as it stands now. Read the diff: it should
flip exactly the rows the change was for.

```bash
go test ./tidb/parser/ -run '^TestTiDBCorpusEngineVerdicts$' -count=1 -v -update
```

rewrites both columns from a live server, and is how the file is re-recorded
when `tidbImage` moves to another TiDB release. `tidbCorpusVersion` must move
with it, or the test refuses to run.

Re-extracting from upstream, after a pingcap bump, means taking the statement
of every one-line case in `parser_test.go` again, keeping this file's order, and
re-recording both columns with the second command.

## `parser_y_corpus.tsv` — the completeness corpus

`pingcap_parser_corpus.tsv` is a sample: the statements pingcap happened to
test. It can show a gap but not that none is left. `parser_y_corpus.tsv` is
derived from the grammar itself, so a clean run is evidence that omni covers
every form the grammar has. For each alternative of each rule in TiDB's
`parser.y` at the pinned tag, it holds the shortest statement that reaches that
alternative from the start symbol.

### Columns

Tab-separated, four fields, one statement per line:

| Field | Values | Means |
|---|---|---|
| 1 | `accept`, `reject` | what `tidb/parser.Parse` does with the statement |
| 2 | `accept`, `reject` | whether TiDB v8.5.5 accepts the syntax |
| 3 | `Rule#n` | the grammar rule and the 1-based alternative the row exercises |
| 4 | the statement | the shortest derivation, one line |

The engine verdict is read the same way as the pingcap corpus's, above, by the
same harness. A gap is a row with column 1 `reject` and column 2 `accept`; an
over-acceptance is the reverse. Both are the grammar work list, not a knob:
close a gap in the parser, not in this file.

### Provenance

Generated by `tidb/parser/cmd/grammarcorpus` from `pkg/parser/parser.y` in
`github.com/pingcap/tidb` at the `v8.5.5` tag (2,946 alternatives across 690
rules). The generator reads a vendored copy of `parser.y` (checked in under the
command's `testdata/`, verbatim at the tag) and writes SQL text; it imports
nothing from pingcap or omni. See that command's README for how the statements
are built and which alternatives need a hand-written sentence. The Apache-2.0
attribution above covers this file too.

### Regenerating

From `<repo root>`, `go run ./tidb/parser/cmd/grammarcorpus` fetches `parser.y`
at the tag and rewrites the statements (columns 3 and 4) with placeholder
verdicts. Then, still from `<repo root>`:

```bash
go test ./tidb/parser/ -run '^TestTiDBParserYCorpusOmniVerdicts$' -count=1 -update
go test ./tidb/parser/ -run '^TestTiDBParserYCorpusEngineVerdicts$' -count=1 -update
```

fill columns 1 and 2, exactly as the pingcap corpus's pair does. When the pin
moves, regenerate and re-record: the diff is the list of forms the new release
added or removed.
