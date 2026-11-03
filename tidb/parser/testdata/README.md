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
image changes, and then the diff is the list of what the new release altered.

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

The reading is the same as everywhere else in this package: error 1064 is the
parser's rejection, any other server error means the text parsed, and anything
that is not a server error is no verdict and fails the test. `PREPARE` adds
1295 to the second group, for statements the prepared-statement protocol does
not carry.

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
