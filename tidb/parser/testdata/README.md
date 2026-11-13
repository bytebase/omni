# TiDB parser test data

`compat.tsv` holds a statement per row, with omni's verdict and the pinned
TiDB's (`../../testdata/pins.tsv`). `TestCompatibility` (`../compat_test.go`)
reads it through `enginetest/compat`, which describes the columns and what
passes, and starts no server. `tapes/` holds what the package's other oracle
tests sent to TiDB and its answers (`enginetest/tape`).

## Rows

| Rule | Rows | From |
|---|---|---|
| Generated (`Rule/symbols`, near-misses `!drop-2`, `!swap`) | 11,055 | `../../testdata/parser.y` and `../../testdata/tokens.tsv` |
| `suite:pingcap-parser` | 3,849 | pingcap's parser test, one form at a time |
| `suite:pingcap-parser-rest` | 884 | the rest of pingcap's parser tests, forms combined |
| `suite:tidb-integrationtest` | 27,242 | TiDB's integration tests |
| `suite:bytebase` | 520 | Bytebase's TiDB test files |
| `suite:reject-fuzz-sweep` | 292 | mutants of accepted statements (`../reject_fuzz_test.go`) |
| `suite:oracle-corpus`, `suite:tidb-extensions` | 104 | statements two oracle tests listed |
| `suite:cited`, `suite:reject-fuzz` | 6 | statements a test cites, and minimised fuzz findings |

**Generated rows.** For each alternative of each rule in `parser.y`
(pkg/parser/parser.y at the pinned tag, vendored verbatim), the shortest
statement that reaches it, then its near-misses: the statement with one token
dropped, the end cut off, a stray word inserted, or a keyword swapped. Where
several alternatives are equally short, a statement takes the one the grammar
lists first, and `'binary'` stands for every string, because it is also a
charset and a collation: TiDB looks those names up while parsing, and an
unknown name answers before a later syntax error can. Three alternatives get no
row, because the default SQL mode never lexes their tokens (`||` as
concatenation, high-precedence `NOT`); `TestCompatibility` lists them.
`TestCompatibility` regenerates these rows and fails if the file holds others.

**Hand rows** carry a label in the rule column and survive regeneration. A
statement found elsewhere goes in as `issue:<key>` or `suite:<name>` with `?`
in both verdict columns, and a recording judges it. A suite's statements that
another row already held were left out of it.

**The reject fuzz.** `TestRejectFuzzRows` holds the `suite:reject-fuzz-sweep`
rows to a generator: 300 token and byte mutations, from a fixed seed, of the
statements TiDB accepts among the generated alternatives and the four corpora.
The pool is read from the expected column, so a gap omni closes does not move
it. An over-acceptance is a row omni accepts and TiDB rejects.

## Commands

From `<repo root>`. Judge new rows, or every row after a pin bump:

```bash
OMNI_RECORD=1 go test ./tidb/parser/ -run '^TestCompatibility$'
```

Commit rows a parser fix moved to agreement:

```bash
go test ./tidb/parser/ -run '^TestCompatibility$' -update-compat
```

After a change to the pool or the generator, rewrite the fuzz rows, then judge
the new ones:

```bash
go test ./tidb/parser/ -run '^TestRejectFuzzRows$' -update-compat
```

Record the tape of a test whose requests changed:

```bash
OMNI_RECORD=1 go test ./tidb/parser/ -run '^TestName$'
```

## Where the engine verdict comes from

`PREPARE`, which parses the statement and stops. Executing cannot measure it,
because statements in the corpora destroy the run that executes them:
`shutdown` stops the server, `SET PASSWORD = 'password'` locks out every new
connection, and `load stats '/tmp/stats.json'` and the `PLAN REPLAYER` file
forms ask the client for a local file.

Error 1064 is the parser's rejection, and it is not the only one: TiDB's grammar
actions refuse text their productions read, with the code each chooses, while
the parse is still running. `CREATE VECTOR INDEX i ON nosuch (a)` answers 1149
and `CREATE VECTOR INDEX i ON nosuch ((a))` answers 1146. Any other server error
came after the parse: 1295 for a statement the prepared-statement protocol does
not carry, 1146 for a table that does not exist. Anything that is not a server
error is no verdict and fails the recording.

A code does not say which side it is on (1105 is both), so the judge asks. It
prepares the statement again with a syntax error after it as a second statement.
The server reads all of the text before it looks at anything, so a statement
refused after parsing now answers 1064, and a statement an action refused
answers its own code again (`compat.PrepareVerdict`, held to the server by
`TestSyntaxProbeControls`).

One class counts as read even though it is raised while parsing: a name the
server looks up, a system variable (1193), a character set (1115) or a collation
(1273, 1253). `Parse` holds no such list, so those are not grammar it can follow.

The judge binds each statement to a user variable and prepares from it, so a
row reaches the server byte for byte: five rows are not valid UTF-8 (pingcap's
GBK and big5 cases), 22 end in a space, and four carry a CR inside a string
literal. The file is split on `\n` alone, so none of that is stripped. A taped
test sends `PREPARE omni_stmt FROM ?` with the statement inlined as an escaped
literal: one request per verdict.

## Provenance

`github.com/pingcap/tidb` is distributed under the Apache License 2.0
(Copyright 2016 PingCAP, Inc.). The license permits redistribution with
attribution; this file is that attribution, for `parser.y` and every pingcap
suite. The rows are test fixtures read at test time and are not compiled into
omni, and omni does not depend on pingcap: nothing here adds a module
requirement. Each suite is one statement per line; statements with an embedded
newline were left out.

- **`suite:pingcap-parser`:** the statement of every one-line case in
  `parser_test.go` of `github.com/pingcap/tidb/pkg/parser`, module version
  `v0.0.0-20251104011144-749c973e5a70` (commit `749c973e5a70`). Only the text
  was taken, not pingcap's flag: it describes pingcap's parser on master, and
  109 statements it accepts are grammar TiDB v8.5.5 does not have.
- **`suite:pingcap-parser-rest`:** from `parser_test.go` and `ast/*_test.go` at
  tag `v8.5.5` (commit `1fa258b833ff113883beeba40bc130be7ce66610`): the `src` of
  every test case the first extraction skipped (backtick and multi-line Go
  literals that still yield one line), and for every `runNodeRestoreTest` call
  the statement it parses, `fmt.Sprintf(template, sourceSQL)`, so a row is a
  function inside a `SELECT`, a reference option inside a `CREATE TABLE`, a
  partition bound inside an `ALTER`. Read with `go/parser` by a throwaway tool.
- **`suite:tidb-integrationtest`:** the statements of
  `tests/integrationtest/t/**/*.test` in `github.com/pingcap/tidb` at tag
  `v8.5.5` (239 files), SQL written to exercise the server rather than the
  grammar. A throwaway tool dropped mysqltest directive and control lines, ran
  `tidb/parser.Split` on the rest, and kept each non-empty segment. Of 47,322
  segments, 1,236 multi-line ones and 243 carrying a mysqltest `$variable` were
  excluded, and 18,515 were duplicates.
- **`suite:bytebase`:** the input-side SQL of the 64 TiDB test files of
  `github.com/bytebase/bytebase` at commit `5a6096dfd5`
  (`backend/plugin/parser/tidb/test-data`, `backend/plugin/advisor/tidb/test`,
  `backend/plugin/schema/tidb/testdata`): the values at the top-level keys
  `statement`, `input`, `oldSchema`, `newSchema` and `schema`, split with
  `tidb/parser.Split`. SQL Bytebase generates (backup, `generate_migration`'s
  `expected`) is not promised input and was left out. Bytebase is this
  project's own repository.
