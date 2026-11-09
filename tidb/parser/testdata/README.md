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

## `parser_y_negatives_corpus.tsv` — grammar-guided negatives (§3.6)

Where `parser_y_corpus.tsv` checks what omni accepts, this file checks what it
rejects. For each reachable alternative's shortest valid sentence, the generator
derives near-misses that break the production in one place: drop a required
token, truncate it, insert a stray token, or swap a keyword for its neighbour in
the keyword set. Each is one step from valid, where a lenient recursive-descent
parser is likeliest to over-accept.

`negatives_corpus_oracle_test.go` is its only reader. The columns are the
completeness corpus's four (omni verdict, engine verdict, label, statement),
read by the same `PREPARE` harness; the label is `Rule#alternative/mutation`
(e.g. `AdminStmt#2/drop-3`, `AlterTableStmt#1/swap`).

The expected verdict is `reject`/`reject`. The gating direction is
over-acceptance: a row omni accepts and the server rejects (`accept`/`reject`)
is a bug, fixed across `mysql`, `mariadb`, and `oceanbase` in the same PR like a
gap. A row omni rejects and the server accepts (`reject`/`accept`) is a §3.3
completeness gap, out of scope here and not silenced by widening `Parse`. Rows
where both accept are near-misses that turned out valid — harmless noise the
measurement sorts out, never a bug.

### Provenance

Generated by `tidb/parser/cmd/grammarcorpus -mode negatives` from the same
vendored `parser.y`, over the same reachable alternatives the completeness
corpus covers. Rows are deduplicated by statement, and a mutation that
reproduces a valid sentence verbatim is dropped, so the file holds only intended
near-misses. The Apache-2.0 attribution above covers this file too.

### Regenerating

From `<repo root>`, `go run ./tidb/parser/cmd/grammarcorpus -mode negatives`
rewrites the statements (columns 3 and 4) with placeholder verdicts. Then:

```bash
go test ./tidb/parser/ -run '^TestTiDBNegativesCorpusOmniVerdicts$' -count=1 -update
go test ./tidb/parser/ -run '^TestTiDBNegativesCorpusEngineVerdicts$' -count=1 -update
```

fill columns 1 and 2, exactly as the completeness corpus's pair does. The engine
test measures the whole corpus live; it needs the pinned TiDB image. When the pin
moves, regenerate and re-record.

## `pingcap_parser_rest_corpus.tsv` — the rest of pingcap's parser tests

The first sample of the design's §3.4 real-SQL corpora. `pingcap_parser_corpus.tsv`
took the single-statement cases of pingcap's parser test one form at a time; this
file takes what that extraction left behind and the forms people combine:

- the `src` of every `[]testCase` / `[]testErrMsgCase` element in
  `parser_test.go` that the first extraction skipped (backtick and multi-line Go
  literals that still yield a one-line statement);
- for every `runNodeRestoreTest` / `...WithFlags` / `...WithFlagsStmtChange` call
  in `ast/*_test.go`, the statement the call parses — `fmt.Sprintf(template,
  sourceSQL)` for each `NodeRestoreTestCase`, so a row is a function inside a
  `SELECT`, a reference option inside a `CREATE TABLE`, a partition bound inside
  an `ALTER`.

`rest_corpus_oracle_test.go` is its only reader. The columns, the `accept`/`reject`
verdicts, the `#` header, and the `PREPARE` engine harness are the pingcap corpus's
(three fields: omni verdict, engine verdict, statement). A gap is `reject`/`accept`;
an over-acceptance is `accept`/`reject` and is a bug.

### Provenance

Extracted from `parser_test.go` and `ast/*_test.go` in
`github.com/pingcap/tidb/pkg/parser` at tag `v8.5.5` (commit
`1fa258b833ff113883beeba40bc130be7ce66610`), by a throwaway tool that reads those
`.go` files with `go/parser` (no build, no network) and emits one statement per
line. Deduplicated against each other and against `pingcap_parser_corpus.tsv`.
Only one-line statements are kept: 30 statements with an embedded newline or tab
are excluded — multi-statement `BEGIN`/`COMMIT` scripts (not single statements),
`LOAD DATA` line terminators and base64 `BINLOG` (newline inside a literal), and
`set`/`enum` members with a tab — as `pingcap_parser_corpus.tsv` excludes them,
because the format is one statement per line. The Apache-2.0 attribution above
covers this file too.

### Regenerating

The engine column is sampled on a normal run and measured whole under `-update`,
so a recorded file is always measured in full. From `<repo root>`:

```bash
go test ./tidb/parser/ -run '^TestTiDBRestCorpusEngineVerdicts$' -count=1 -update
```

rewrites both columns from a live server. `TIDB_CORPUS_N` caps how many rows a
plain (non-`-update`) run measures live, so CI stays bounded as the later §3.4
samples land; the omni-verdict pass always reads every row. Re-extracting after a
pingcap bump means rerunning the extraction above at the new tag and re-recording.

## `tidb_integrationtest_corpus.tsv` — TiDB's integration tests

The second sample of the design's §3.4 real-SQL corpora. Where the two files
above are grammar forms — tested one at a time, then combined the way pingcap's
own parser tests combine them — this one is SQL people wrote to exercise the TiDB
*server*: the statements of TiDB's `tests/integrationtest/t/**/*.test` files, the
closest public stand-in for user SQL. It leans heavily on the DML and query
shapes the other two under-sample.

`integrationtest_corpus_oracle_test.go` is its only reader. The columns, the
`accept`/`reject` verdicts, the `#` header, and the `PREPARE` engine harness are
the pingcap corpus's (three fields: omni verdict, engine verdict, statement). A
gap is `reject`/`accept`; an over-acceptance is `accept`/`reject` and is a bug.
It is far larger than the other samples, so the live engine test measures only a
`TIDB_CORPUS_N` sample on a plain run (every row under `-update`), while the
omni-verdict test reads every row.

### Provenance

Extracted from `tests/integrationtest/t/**/*.test` in `github.com/pingcap/tidb`
at tag `v8.5.5` (commit `1fa258b833ff113883beeba40bc130be7ce66610`, 239 files),
by a throwaway tool (not checked in, like sample 1's). For each file it drops the
mysqltest directive/comment lines (a leading `--` or `#`) and the bare mysqltest
control lines (`connection`, `connect`, `disconnect`, `let`, `echo`, …, which
often lack a trailing `;` and would otherwise merge into the next statement),
runs `tidb/parser.Split` on what remains, and keeps each non-empty segment.
Deduplicated within this corpus and against `pingcap_parser_corpus.tsv` and
`pingcap_parser_rest_corpus.tsv`. Only single-line statements are kept, as the
other corpora are: of 47,322 non-empty segments, 1,236 with an embedded newline
or tab (multi-line `CREATE`/`INSERT`, compound `BEGIN`/`END` blocks) and 243
carrying a mysqltest `$variable` are excluded, 18,515 are duplicates within this
corpus and 82 already appear in the other two, leaving 27,246 rows. The
Apache-2.0 attribution above covers this file too.

### Regenerating

The engine column is sampled on a normal run and measured whole under `-update`.
From `<repo root>`:

```bash
go test ./tidb/parser/ -run '^TestTiDBIntegrationtestCorpusEngineVerdicts$' -count=1 -update
```

rewrites both columns from a live server. Re-extracting after a pingcap bump
means rerunning the extraction above at the new tag and re-recording.

## `bytebase_tidb_corpus.tsv` — Bytebase's own TiDB test files

The third sample of the design's §3.4 real-SQL corpora: the SQL Bytebase already
promises to handle. Where the pingcap corpora are pingcap's grammar cases, these
are the forms Bytebase's users and tests drive through the TiDB path — column
rules, naming checks, DML dry-runs, schema syncs and walk-throughs — so an
over-acceptance here is a form Bytebase ships against.

`bytebase_corpus_oracle_test.go` is its only reader. The columns, the
`accept`/`reject` verdicts, the `#` header, the `PREPARE` engine harness and the
`TIDB_CORPUS_N` sampler are the rest corpus's.

### Provenance

Extracted from the 64 TiDB test files in `github.com/bytebase/bytebase` at commit
`5a6096dfd5`, under `backend/plugin/parser/tidb/test-data` (7 files),
`backend/plugin/advisor/tidb/test` (53 files) and
`backend/plugin/schema/tidb/testdata` (3 of 4 files). A throwaway Go extractor
read the input-side SQL — the values at the top-level case keys `statement`,
`input`, `oldSchema`, `newSchema`, `schema` — split each script with
`tidb/parser.Split`, trimmed, and kept the non-empty single-line segments.
Top-level keys only, so a query-span source column's nested `schema` and a backup
case's generated `result` statements are not picked up; generated-output SQL
(backup/restore rollback, `generate_migration` `expected`,
`get_database_definition`) is Bytebase-synthesized, not promised input, and is
excluded. 527 rows after dropping 194 multi-line statements, 209 internal
duplicates, and 15 already carried by the two pingcap corpora. The method is
recorded in the file's own header and in
`docs/plans/2026-11-08-tidb-bytebase-corpus-gap-inventory.md`.

`github.com/bytebase/bytebase` is this project's own repository; these rows are
test fixtures read at test time, not compiled into omni, and add no module
dependency.

### Regenerating

The engine column is sampled on a normal run and measured whole under `-update`.
From `<repo root>`:

```bash
go test ./tidb/parser/ -run '^TestTiDBBytebaseCorpusEngineVerdicts$' -count=1 -update
```

rewrites both columns from a live server. Re-extracting after a bytebase bump
means rerunning the extraction at the new commit and re-recording.
