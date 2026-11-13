# Engine compatibility without servers in CI

Status: draft, 2026-11-12. Step 0 is built: `enginetest/compat`, `enginetest/grammar` and TiDB's `compat.tsv` (11,058 rows, recorded on the pinned server in 9 seconds), and `enginetest/tape` with `tidb/catalog` on tapes (131 tapes, 3.1 MB, recorded in 90 seconds; replay takes 1.6 seconds and starts no server). OceanBase has step 1: 20,407 rows, judged on the pinned server in three minutes, 17,728 agreeing (Appendix A). Its step 2 has begun: `oceanbase/parser`, `oceanbase/validate` and 80 tests of `oceanbase/catalog` replay tapes (163 tapes, 6.5 MB) and start no server. The 80 are the tests that always asked OceanBase and every oracle test written against MySQL, moved with the fixes they needed and with what still fails recorded (Appendix A). Catalog's 107 container tests still ask MySQL (D3), and then `review`, the last package under `oceanbase/` that starts a server. A spike measured OceanBase and Snowflake first. Its code was lost. The other engines' grammars were located, not run.

## 1. Summary

**CI starts no database server. The engine's answers are committed to the repo, and tests compare omni with those files.**

```text
        Recording: OMNI_RECORD=1, one engine, a developer machine        Every other run: go test ./...

vendor grammar ─▶ generated statements ──▶ judge ──▶ compat.tsv    ─▶  omni's parser, held to the file
                  DDL × starting schemas ─▶ engine ─▶ simulate.tsv  ─▶  omni's catalog, held to the file
hand-written test ───────────────────────▶ engine ─▶ <Test>.tape   ─▶  the same test, answered from the file
```

| Committed answer | One entry holds | Checks |
| --- | --- | --- |
| `compat.tsv` (§2) | A statement, omni's verdict, and the judge's | Syntax: parser, completion, analysis |
| `simulate.tsv` (§3) | A DDL statement on a starting schema, omni's outcome, and the engine's | DDL behavior: catalog, walk-through |
| Tapes (§4) | Every request a hand-written test sends, and the engine's answer | Everything else: generated SQL, Split, review, validate |

Four terms are used throughout:

| Term | Meaning |
| --- | --- |
| Pin | The exact engine release the answers come from (§5.1) |
| Judge | What decides the expected answer: the pinned engine where there is one, otherwise the vendor's grammar or docs (§2.3) |
| Recording | A run that asks the judge and writes answers: `OMNI_RECORD=1` (§5.2) |
| Replay | Every other run. It reads the files and reaches no server |

**An engine is fully compatible at its pin when all four hold:**

1. **Complete:** its rows cover every grammar alternative and every token form in `tokens.tsv`. They are generated, and `TestCompatibility` regenerates them to prove none was dropped.
2. **Judged:** no row is waiting for a judge.
3. **Agreeing:** on every row, omni's answer equals the judge's.
4. **Replaying:** every tape replays.

Until then, the disagreeing rows are the work list: 2,679 for OceanBase today. The guards in §7 make that count only fall, and once it reaches zero any disagreement fails CI.

**What that does not claim:**
- **Statements no row covers.** Rows sample an infinite language. A statement outside them is unknown until someone adds it as a row, and then it is covered for good.
- **More than the docs, for docs engines.** Where docs are the judge, omni is compatible with the docs, not necessarily with the engine.
- **How omni reads a query.** Agreeing on accept or reject doesn't prove omni built the right tree. `simulate.tsv` proves it for DDL, by comparing what gets stored. Queries are D6.

**Omni's parsers stay hand-written.** A generated parser would add a runtime dependency, lose omni's `Loc`, depth limits and AST shapes, and still miss checks made outside the grammar rules. The grammar generates the test set instead.

**Why now.** OceanBase's container takes 3.5 GiB and about 21 of CI's 30 minutes on a box that runs out of memory, and its packages still test against MySQL. Six engines have no judge at all: cassandra, cosmosdb, elasticsearch, mongo, partiql, snowflake.

## 2. Syntax: `compat.tsv`

### 2.1 The file

`<engine>/parser/testdata/compat.tsv` has one row per statement. It generalizes the corpus TiDB already keeps (`tidb/parser/testdata/parser_y_corpus.tsv`):

```text
# pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 v4.3.5_CE_BP6
# columns: omni, expected, source, rule, statement
accept	accept	engine	stmt_list/BEGIN	BEGIN
reject	accept	engine	alter_table_action/REMOVE_TTL	ALTER TABLE a REMOVE TTL
accept	reject	engine	func_expr/MONTH!drop-3	DO MONTH()
?	?	-	issue:BYT-9120	ALTER TABLE t TRUNCATE PARTITION ALL
```

1. An agreement.
2. A known gap: OceanBase accepts this and omni doesn't.
3. A known gap the other way: omni accepts a broken statement.
4. A bug report just added. It fails the test until a recording judges it, so it is committed judged, never as `?`.

The first three are real spike results. The fourth shows the form of a new row.

| Column | Holds |
| --- | --- |
| `omni` | Omni's verdict when the row was last updated |
| `expected` | The judge's verdict |
| `source` | Which judge: `engine`, `grammar` or `docs` |
| `rule` | The grammar rule the row came from, or a label for a hand-added row |

A verdict is `accept`, `reject`, or `reject@<byte>` where the judge reports a position, as pg does. Two more values exist: `undocumented` (§2.3), and `?` for a cell not filled yet.

The test is one call. The engine supplies only its parser:

```go
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    "testdata/compat.tsv",
		Pins:    "../testdata/pins.tsv",
		Grammar: "../testdata/parser.y",
		Tokens:  "../testdata/tokens.tsv",
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge, // called only in a recording
	})
}

func omniVerdict(sql string) string {
	if _, err := Parse(sql); err != nil {
		return "reject"
	}
	return "accept"
}
```

`compat.Check` runs omni on every row and compares:

| Omni's verdict now | Result |
| --- | --- |
| Equals the `omni` cell | Pass. The row is an agreement or a known gap. |
| Differs from the cell, equals `expected` | Fail: a gap closed. Rerun with `-update-compat` to commit it. |
| Differs from both | Fail: a regression. |
| The row is `?` | Fail: not judged yet. |

Recording both columns is what makes this work. Thousands of known gaps pass quietly, and the one statement a change breaks fails by name:

```text
index_or_heap/ORGANIZATION_HEAP: omni now says reject, engine says accept:
    ALTER TABLE a ORGANIZATION HEAP
```

### 2.2 Where rows come from

**Generated rows** come from the vendor's grammar at the pinned release:

| Kind | Rows | Example | Covers |
| --- | --- | --- | --- |
| Alternative | The shortest statement for each alternative of each rule | `ALTER TABLE a REMOVE TTL` | Every production |
| Near-miss | Each alternative broken once: a token dropped, a word inserted, a keyword swapped, the end cut off | `ALTER TABLE a REMOVE` | Rejections |
| Derivation | N random, deeper expansions per rule, from a fixed seed | `SELECT a FROM t WHERE (b + 1) * 2 > ANY (SELECT c FROM u)` | Combinations, nesting, precedence |
| Token | Every form of every token, in a carrier statement | ``SELECT `a``b` ``, `SELECT 0x1F`, `SELECT /*+ INDEX(t i) */ 1` | Lexing, which grammar rules never see |

- **The judge decides every row.** A near-miss can still be valid: OceanBase's own grammar accepts 900 of its 15,560.
- **Ties follow the grammar's order.** Among equally short choices a statement takes the alternative the grammar lists first, and sits in the first statement that reaches it. So an expression is a plain name, not a keyword that may stand for one, and a hint sits in `SELECT`.
- **Samples resolve.** Where an engine's parser looks a name up while it parses, as TiDB does for charsets and collations, `tokens.tsv` samples a name that exists. An unknown name would answer before a later syntax error, and the row would read as accepted.
- **Rule IDs are stable.** A row is named by its rule and the alternative's symbols, not by an index, so a pin bump's diff shows only what the release changed.
- **Token forms come from a committed `tokens.tsv`,** cited to the vendor's lexer or lexical docs. It is written by hand, because the lexers of MySQL, MariaDB and TiDB are code, not grammar files. Other rows take their sample identifiers and literals from it.

**Hand-added rows** carry a label in `rule` and survive regeneration: `issue:<key>` for a bug report, `suite:<name>` for a corpus. Today's corpora become `suite:` rows: TiDB's, Doris and StarRocks conformance, Trino's differential, pg's fuzz, and the statements that completion and analysis tests check for validity.

**Grammar files.** The generator reads two formats, bison `.y` and ANTLR `.g4`.
- A vendor `.y` or `.g4` is used as published. It is committed where the license allows, as TiDB's `parser.y` already is, and fetched otherwise (D1): `<engine>/testdata/sources.tsv` names the file by the release's commit and its SHA-256, and `enginetest.Fetched` downloads it once per machine and checks it on every use. OceanBase's grammar is read this way.
- Anything else is converted once per pin to a committed `.g4`, each rule citing its source: docs syntax blocks, Doris's CUP, GoogleSQL's Textmapper, ANTLR 2 and 3 (D2).
- A community grammar is never used. Omni's Snowflake parser was ported from one, so that grammar can't judge it.
- Mongo and Elasticsearch have no statement grammar. Their rows cover call shapes (method names, endpoints), not argument documents or request bodies.

### 2.3 Who judges

`expected` comes from the strongest judge the engine has:

1. **`engine`:** the pinned server.
2. **`grammar`:** a recognizer over the vendor's grammar.
3. **`docs`:** a recognizer over the `.g4` converted from the vendor's docs.

Each engine with a server has one function that reads the server's answer as a verdict:

| Engines | Verdict from |
| --- | --- |
| mysql, mariadb, tidb, oceanbase | `PREPARE`, read by TiDB's `prepareVerdict` |
| pg, redshift | SQLSTATE `42601` from the `SELECT 1/0` probe |
| mssql | `SET PARSEONLY ON` |
| oracle | `DBMS_SQL.PARSE` |
| trino, doris, starrocks, googlesql | The engine's error name or message |

An answer the function can't read fails the recording. It is never guessed.

OceanBase needs three things of its own, all measured on the pinned server (`oceanbase/parser/testdata/README.md`):
- **The probe's tail stays inside the statement,** because OceanBase refuses to prepare two statements before it parses either.
- **1064 comes in two kinds.** The parser's says where the error is. A bare "Parse error" comes from later layers, and from `PREPARE` itself for statements it does not carry, so the probe decides it.
- **Two statements crash the server.** The judge starts it again, asks the statement again, and judges one that kills the server twice by its probe alone.

The two recognizers are weaker judges, in known ways:
- **A grammar recognizer misses checks made outside the rules,** in the lexer and in grammar actions. On OceanBase it agreed with the engine on 90% of statements (Appendix A). Each miss found becomes a cited row in `grammar-overrides.tsv`.
- **Docs only say what is valid.** A near-miss judged by docs is `undocumented`, not `reject`. It counts as neither a gap nor a regression, so omni is never pushed to reject forms users write.
- **Snowflake and Redshift can upgrade.** With `SNOWFLAKE_COMPAT_DSN` or `REDSHIFT_COMPAT_DSN` set, a recording asks the engine instead of the docs. The account must be dedicated and scratch-only, because preparing DDL there executes it.

## 3. DDL behavior: `simulate.tsv`

Omni's catalog simulates DDL. `catalog.Exec` either changes the schema as the engine would, or refuses with the engine's error code. `simulate.tsv` holds the engine's side of that. One row, traced:

```text
fixture one_table     CREATE TABLE t (a INT PRIMARY KEY, b INT, KEY idx_b (b))
statement             ALTER TABLE t ADD COLUMN a INT
engine                error 1060 (duplicate column)
omni, catalog.Exec    error 1060                      → the row agrees
```

Rows are generated: every DDL statement the engine accepts in `compat.tsv`, on each fixture, with each object name filled three ways.

```text
# columns: omni, expected, source, fixture, rule, statement
ok:7c21…	ok:7c21…	engine	one_table	alter_table_action/ADD_COLUMN@exists	ALTER TABLE t ADD COLUMN c INT
1060	1060	engine	one_table	alter_table_action/ADD_COLUMN@collides	ALTER TABLE t ADD COLUMN a INT
ok:3f9a…	1146	engine	one_table	alter_table_action/ADD_COLUMN@missing	ALTER TABLE nosuch ADD COLUMN c INT
```

These rows show the format. No simulation has been run.

- **Fixtures** are a few committed starting schemas in `<engine>/catalog/testdata/fixtures/`: an empty database, one table with keys, a foreign-key pair, a partitioned table, a view.
- **Names.** The grammar doesn't say which symbols must name an existing object. `references.tsv` does, in tens of rows per engine. The generator fills each with a name that exists, one that is missing, and one that collides.
- **Outcomes** are an error code, or `ok:<hash>`. The hash is of the database's stored form after the statement: every object's `SHOW CREATE`, in name order. So an `ok` row agrees only when omni stores what the engine stores. The texts are on a tape (§4), so a mismatch prints a diff.
- **Everything else is §2.1:** the same columns plus `fixture`, the same checker, the same guards. The test is `TestSimulation`.

**Walk-through reads the same rows.** For each row, the walk runs on the fixture and the statement:
- On an `ok` row, any finding fails.
- On a refused row, a finding must name the engine's code, or claim only facts true of the fixture.
- Silence passes, since the walk must never be wrong but may be unsure. How often it finds the refusal is printed.

**Cost is unmeasured, and measuring it is the first job.** At about 5 seconds per row on OceanBase, 1,000 statements × 5 fixtures × 3 names would take about 20 hours to record. Rollout starts with one fixture. Docs engines have no `simulate.tsv`, since docs rarely name the error a statement raises.

## 4. Hand-written tests: tapes

### 4.1 What changes in a test

A test that talks to a server keeps its code. It opens its connection differently and iterates maps in sorted order:

```diff
 func TestOracle_ConstraintOfIndexSurvivesOtherNames(t *testing.T) {
-	o := oceanBaseOracle(t)
+	db := testcontainer.Open(t)
 	...
-	for name, script := range constraintOfIndexScripts {
+	for _, name := range slices.Sorted(maps.Keys(constraintOfIndexScripts)) {
```

`testcontainer.Open(t)` returns an ordinary `*sql.DB`. In replay it answers from the test's tape. In a recording it talks to the pinned server and writes the tape. Errors keep the driver's own type, so `errors.As(err, &mysqlErr)` still works.

### 4.2 What a tape holds

One text file per top-level test, `<package>/testdata/tapes/<Test>.tape`. Each request is followed by the engine's answer. Excerpts from `tidb/catalog`:

```text
-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6… v8.5.5

== TestContainerSmoke default#1
> exec
| CREATE TABLE t1 (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `id` int NOT NULL,
|   `name` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
...

== TestContainerCreateInvisiblePrimaryKey default#1
...
< error
| 3522 HY000 A primary key index cannot be invisible
```

Every line starts with a sigil, so SQL and values need no escaping: `|` carries a line of text, `:` a value or an argument (a request keeps its `?` placeholders), `>` a request, `<` an answer. A line-oriented file means a diff shows the line that changed, even inside a `SHOW CREATE`. A file per test means two PRs conflict only when they change the same test.

### 4.3 The rule

**What a test sends depends only on its inputs, never on omni's results.**

Replay holds a test to its tape in both directions. It fails on a request the tape doesn't have next, and on answers the test never asked for.

Two kinds of test send omni's own output, because that output is what they test: round trips (migrations, diff, deparse) and Split. A change to the SQL omni generates therefore fails replay until the engine has seen the new SQL:

```text
no recorded answer for exec "ALTER TABLE `t` ADD COLUMN `c` int NOT NULL" (statement 3 of TestMigration_AddColumn#1).
Record it:
  OMNI_RECORD=1 go test -run '^TestMigration_AddColumn$' ./oceanbase/catalog/
```

This is deliberate, and it is the design's main recurring cost: a PR that changes DDL generation needs the pinned server before it merges.

### 4.4 What goes on tapes

| Component | The tape holds |
| --- | --- |
| Catalog stored form and load | DDL, `SHOW CREATE`, `information_schema` rows |
| Migrations, diff, deparse | Omni's generated SQL, and what the engine stores after it |
| Split | pg's S3 differential: a script run whole, then segment by segment |
| Walk-through on changes | Each statement's outcome, and the database's objects before it |
| Review rules | The query results and stored state a rule predicts |
| Validate | Each routine body's outcome on its fixture |
| Completion | How the engine folds identifiers |
| ScriptDOM | The AST shapes and error positions Microsoft's parser reports |

**Walk-through needs one change to follow the rule.** Today, when a finding doesn't name the server's error code, `Holds` queries the server about the objects the finding claims (`walktest.go:177`). A change to the findings would then change what the test sends. Instead, the engine lists the database's objects before every statement, and `Holds` checks the claim against that list:

```go
type Engine interface {
	// ...Statements, Exec, Report, Violation, Codes as today
	Inventory(t testing.TB, change string) []Object // fixed queries, whatever the walk says
	Same(claimed, stored walk.Key) bool             // the engine's name folding
}
```

The inventory's hash replaces `Digest`. Sweeps keep their fixed seeds, and the tape header records each sweep's size and seed.

### 4.5 Writing a replayable test

- Sort map iteration.
- Use the test's own database. Read no server-wide state outside package setup.
- Take names from `t.Name()`, not from counters or the clock.
- Order every multi-row read.
- Don't skip. Today's skips and DSN fallbacks go. A probe the engine refuses, or one omni fails, is recorded with what it must do instead: `obRefused` and `obGaps` in `oceanbase/catalog/oceanbase_oracle_test.go` each fail an entry that no longer holds.
- Keep a test on one connection. A tape's stream is one, and a helper that asks the pool while its test holds that connection waits forever.
- Read the stored form from the database the change is applied to. A name the engine makes up differs in every database.
- Probe syntax with SQL's `PREPARE ... FROM '<text>'`, not the driver's `Prepare`: a tape holds no driver-level prepared statement.
- Ask a verdict in one request, `PREPARE x FROM ?` with the statement as its argument. Setting a variable, preparing from it and deallocating is three requests, and a sweep of ten thousand statements pays for each: OceanBase's parser tapes went from 13.5 MB to 4.7 MB on this and the two forms below.

A recording catches an answer that changes between runs, but not a request order that does, such as a map iterated while sending. Replay a fresh recording a few times (`-count=10`) before committing it.

**A few tests have nothing to replay:** harness self-tests, timing and hang tests, and the Oracle PL/SQL dictionary test, whose input is Oracle's own source. They call `tape.Live(t)`, which skips them outside a recording. That is the only environment gate. The dictionary test keeps its known-failures manifest as it is today.

The harness is in Appendix C.

## 5. Recording

### 5.1 Pins

`<engine>/testdata/pins.tsv` names each server by image digest, with the source tag it was built from:

```text
# role	image	digest	source tag
default	oceanbase/oceanbase-ce:4.3.5.6-106000012026040916	sha256:…	v4.3.5_CE_BP6
```

- **One release, one grammar.** The grammar comes from the release's source tag, never from a branch.
- **Roles.** `default` judges the rows and records the tapes. An extra row serves tests that need another server, such as pg 18 or MySQL with `lower_case_table_names=1`.
- **Docs engines** pin the conversion date instead of an image.
- **Every committed file names its pin** in a header. A test fails when the header no longer matches `pins.tsv`.

### 5.2 Commands

| Command | Writes | Needs the judge |
| --- | --- | --- |
| `go test ./...` | Nothing | No |
| `go test ./<engine>/<pkg>/ -update-compat` | `omni` cells that now equal `expected` | No |
| `OMNI_RECORD=1 go test ./<engine>/...` | `?` cells, and the tapes of the tests it runs | Yes |

- **One engine at a time.** No machine holds every server.
- **The judge starts itself.** For an engine with a server, the tests start the pinned image by digest, so recording needs Docker. For the others, the judge is a recognizer and needs nothing.
- **A recording must reproduce.** If a test sends what its committed tape holds, under the same pin, the answers must match, or the recording fails. Nothing nondeterministic gets committed, and an answer the engine did not give (a timeout, a driver fault) fails the recording and leaves the committed tape as it was.
- **Missing Docker fails.** A recording never skips.
- **`OMNI_RECORD=1` fails when `CI` is set.** CI never writes answers.
- **Nothing is build-tagged.** The recording code compiles on every PR, so it can't rot unseen.

### 5.3 Where things live

```text
enginetest/                          shared by every engine, knows none of them
  grammar/                           .y and .g4 readers, converters, generator, recognizers
  compat/                            the row checker, prepareVerdict
  tape/                              record and replay
  cmd/compatcheck/                   the CI guards that compare a PR's files with its base

oceanbase/
  testdata/pins.tsv                  the pinned server
  testdata/tokens.tsv                token forms and sample values
  testdata/sources.tsv               vendor files fetched, not committed: name, SHA-256, URL
  parser/testdata/compat.tsv
  catalog/testdata/simulate.tsv
  catalog/testdata/fixtures/         starting schemas
  catalog/testdata/references.tsv
  <package>/testdata/tapes/          one .tape per test
  internal/testcontainer/            Open(t), the error codec, the live connector
```

`enginetest` is new shared code, which D4 must allow. `prepareVerdict` takes plain error codes, so the package imports no driver: each MySQL-family engine converts its driver's error first.

## 6. Everyday work

**Do I need a server?**

| You change | Engine has a server | Engine has none |
| --- | --- | --- |
| Omni's code, and tests send the same requests | No | No |
| Omni's walk-through findings | No | — |
| SQL omni generates, or where `Split` cuts | **Yes:** record the affected tests | — |
| A hand-written test | **Yes:** record that test | — |
| A new row, fixture or reference | **Yes:** judge it | No: the recognizer judges |
| The pin | **Yes:** record the whole engine | Regenerate, or re-convert the docs |

**Fix a gap.** No server.

```bash
go test ./oceanbase/parser/ -run '^TestCompatibility$'
#   alter_table_action/REMOVE_TTL: gap closed; rerun with -update-compat
go test ./oceanbase/parser/ -run '^TestCompatibility$' -update-compat
```

```diff
-reject	accept	engine	alter_table_action/REMOVE_TTL	ALTER TABLE a REMOVE TTL
+accept	accept	engine	alter_table_action/REMOVE_TTL	ALTER TABLE a REMOVE TTL
```

**Report a bug.** Add the statement as a `?` row, and let the engine judge it. Then fix omni as above.

```bash
printf '?\t?\t-\tissue:BYT-9120\tALTER TABLE t TRUNCATE PARTITION ALL\n' >> oceanbase/parser/testdata/compat.tsv
OMNI_RECORD=1 go test ./oceanbase/parser/ -run '^TestCompatibility$'
```

Commit the row as the recording filled it.

**Change generated SQL.** Replay fails and names the test (§4.3). Record it:

```bash
OMNI_RECORD=1 go test ./oceanbase/catalog/ -run '^TestMigration_AddColumn$'
```

**Bump a pin.** In a PR of its own, with no code changes:

```bash
# edit oceanbase/testdata/pins.tsv and the vendored grammar, then:
OMNI_RECORD=1 go test ./oceanbase/...
```

The diff is what the release changed. Fix omni in later PRs.

## 7. Guards

Each guard is a test or a CI step (`go run ./enginetest/cmd/compatcheck origin/main`).

| Risk | Guard |
| --- | --- |
| A change breaks a row omni handled | The row checker fails (§2.1). |
| A regression is committed as the new `omni` | CI fails if an `omni` cell changes to anything but its row's `expected`. |
| `expected` is edited to pass | CI fails if an `expected` cell changes, unless it was `?` or the PR bumps a pin. |
| A `grammar` or `docs` cell is written by hand | CI reruns the recognizer on new and changed rows. |
| Generated rows are added, removed or edited by hand | `TestCompatibility` regenerates them from the pinned grammar and diffs. A recording rewrites them. |
| A regression hides in a pin bump | CI fails a pin-bump PR that changes non-test code. |
| A test sends something unrecorded, or stops sending what was recorded | Replay fails (§4.3). |
| A deleted or renamed test leaves its tape behind | After a passing run of every test (no `-run`, `-skip` or `-list`), the harness fails on a tape no test opened. |
| An answer comes from the wrong server | Headers must match `pins.tsv`. |
| A nondeterministic answer is committed | The recording fails (§5.2). |
| A tape or an `engine` cell is written by hand | Not caught without a server. See D5. |

Progress is visible: `TestCompatibility` prints how many rows of each kind agree, and `docs/engine-test-matrix.md` keeps the numbers per engine.

## 8. Engines

| Engine | Grammar | Judge |
| --- | --- | --- |
| pg | `gram.y` | `postgres:17-alpine` |
| redshift | AWS docs → `.g4` | docs; engine with `REDSHIFT_COMPAT_DSN` |
| mysql | `sql_yacc.yy` @ `mysql-8.0.46` | `mysql:8.0.46` |
| mariadb | `sql_yacc.yy` @ `mariadb-11.8.8` | `mariadb:11.8.8` |
| tidb | `parser.y` @ `v8.5.5` | `pingcap/tidb:v8.5.5` |
| oceanbase | `sql_parser_mysql_mode.y` @ `v4.3.5_CE_BP6` | OceanBase 4.3.5.6 |
| mssql | ScriptDOM `TSql160.g` → `.g4` | SQL Server 2022 CU27 |
| oracle | `oracle/parser/bnf/` docs → `.g4` | Oracle Free 23.26.3 |
| snowflake | Snowflake docs → `.g4` | docs; engine with `SNOWFLAKE_COMPAT_DSN` |
| googlesql | ZetaSQL `googlesql.tm` → `.g4` | Spanner emulator 1.5.54 |
| trino | `SqlBase.g4` @ `482` | `trinodb/trino:482` |
| doris | `DorisParser.g4` + `sql_parser.cup` → `.g4` @ `3.1.4` | Doris 3.1.4 |
| starrocks | `StarRocks.g4` @ `3.4.10` | StarRocks 3.4.10 |
| cassandra | CQL `Parser.g` → `.g4` | grammar |
| partiql | `PartiQLParser.g4` + DynamoDB docs | grammar + docs |
| cosmosdb | `sql.g4` | grammar |
| mongo | mongosh method docs → `.g4` | docs |
| elasticsearch | `elasticsearch-specification` → `.g4` | grammar |

- **TiDB is the starting point.** Its `parser.y` rows were the first `compat.tsv`; its other corpora become `suite:` rows in step 1.
- **Engines judged by grammar or docs** move up to `engine` once an image is pinned: Cassandra, DynamoDB Local, the Cosmos DB emulator, MongoDB and Elasticsearch all ship one.

## 9. Rollout

| Step | What | Done when |
| --- | --- | --- |
| 0 | `enginetest`: generator, row checker, tape harness, `compatcheck` | TiDB's existing corpus runs through it |
| 1 | Per engine: `pins.tsv`, `tokens.tsv`, `compat.tsv`, `TestCompatibility` | Every row judged, baseline committed |
| 2 | Per engine with a server: tapes for every live test, then `simulate.tsv` on one fixture | The engine's tests pass with Docker off |
| 3 | Per engine: remove its server from CI | Its CI starts no container |
| 4 | Per engine: more fixtures, and the disagreements worked down | Fully compatible (§1) |

After step 0 the engines are independent. Order:
1. **OceanBase,** which is the CI pain. Its roughly 160 catalog tests and two review files that run on MySQL today move to OceanBase (D3). The #808 branch has native fixes and 65 cited differences.
2. **Snowflake.**
3. **The other engines with no judge.**
4. **The container engines,** which work today and gain only CI time.

## 10. Decisions

- **D1. Grammar licenses.** MySQL's and MariaDB's grammars are GPL, OceanBase's is Mulan PubL at the pinned tag (its default branch has since moved to Apache-2.0), and omni is MIT. Recommended: don't commit those three. Fetch them by tag, and commit only the generated rows. OceanBase does this (§2.2); the cost is one download per machine, CI's runner included. Permissively licensed grammars are committed, as TiDB's already is.
- **D2. Converted docs.** Confirm that committing `.g4` files derived from vendor docs is acceptable. Oracle's `bnf/` directory already quotes them.
- **D3. MySQL tests in OceanBase and MariaDB.** Recommended: point them at their own engine and record, deleting only tests of MySQL-only behavior.
- **D4. Shared `enginetest`.** Decided: allowed, like `review/walk/walktest`, rather than 18 copies. `AGENTS.md` lists it.
- **D5. A verify lane.** Nothing in §7 catches an engine answer written by hand. Option: a third mode, `OMNI_RECORD=verify`, asks the engine again, fails on any difference, and writes nothing. The self-hosted runner runs it on just the tests and rows a PR changed, so most PRs still start no server. The alternative is to trust review.
- **D6. Queries.** Stored form proves omni reads DDL as the engine does. The same check reaches queries by wrapping each accepted `SELECT` row in `CREATE VIEW`: the engine's stored view text is its own reading of the query. It works where omni's catalog renders views (pg, the MySQL family), and it multiplies the size of `simulate.tsv`. Recommended: decide once step 2 has measured what a row costs.

## Appendix A. Measured results

**OceanBase, recorded.** `oceanbase/parser/testdata/compat.tsv` at `v4.3.5_CE_BP6`: 5,123 alternatives, of which 52 need a token the lexer never produces, and 15,336 near-misses. The pinned server judged all 20,407 rows in three minutes and accepts 4,995 of the 5,071 alternatives.

| Statements | Total | Agree | Only engine accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 5,071 | 3,306 | 1,736 | 29 |
| Near-misses | 15,336 | 14,422 | 527 | 387 |

- **2,679 disagreements,** untriaged. Among the alternatives only OceanBase accepts, the largest groups by leading words are `ALTER SYSTEM` (399), `EXPLAIN` (231, of which 205 are names with a leading dot), `ALTER TABLE` (175), `SELECT` (166) and `DELETE` (131).
- **A carrier can hide its alternatives.** A row is the shortest statement that reaches an alternative, and where omni rejects that carrier, every alternative behind it disagrees for one reason. Fixing the carrier is what shows the rest.
- **Two statements crash the server** (three rows). Their verdict comes from the probe (§2.3).

**OceanBase, catalog's MySQL tests.** Measured 2026-11-12 by pointing them at the pinned OceanBase as written. Nothing of the run is committed.

| Tests | Top-level | Pass | Subtests | Pass | Fail | Skip |
| --- | --- | --- | --- | --- | --- | --- |
| Stored form and migrations (`connectOracle`) | 56 | 20 | 708 | 295 | 357 | 56 |
| `SHOW CREATE` text and scenarios (`startContainer`) | 106 | 15 | 873 | 335 | 537 | 1 |

- **Stored form.** 184 failures are a table that differs from its own readback. Where traced, the cause is OceanBase printing `ROW_FORMAT = DYNAMIC` on every table. 45 are views that differ from theirs. 76 are readbacks omni can't load: 47 foreign keys whose referenced table comes back qualified by its database, 20 triggers, 8 multi-valued indexes printed as `json_query(... RETURNING ...)`, and a routine.
- **Migrations.** 41 generated statements are refused: 14 for combining changes OceanBase takes one at a time (1235, 4179), 7 for dropping a key that is not there (1091), 6 for altering a geometry column's SRID (1235), 5 partition statements as syntax errors, and 9 others. 4 more apply and leave the index visible: `ALTER TABLE ... ADD KEY` does not take `INVISIBLE`, written plainly or as `/*!80000 INVISIBLE */`, where `CREATE TABLE`, `CREATE INDEX` and `ALTER INDEX` do.
- **Text.** The 48 `TestContainer_*` tests mostly compare `SHOW CREATE TABLE` with MySQL 8.0's text, which OceanBase never prints: 301 of their 378 subtests fail. The 27 view tests do the same for view text: 131 of 204 fail. The 25 scenario tests are closest: 163 of 241 pass.
- **The hangs #808 recorded were the tests' own.** A helper asked the pool for a second connection while its test held the only one. The server sat idle each time, and the cross-object ordering test passes 5 of 5 once the pool has a second connection.
- **Not measured:** the keyword sweep builds an index per keyword and did not finish one of its seven positions in five minutes. A tape can't hold a `uint64` yet, which `information_schema`'s unsigned columns return.

**OceanBase, catalog's oracle tests moved.** All 53 now ask OceanBase (2026-11-13; three more read MySQL's `sys` schema and are gone): 80 tapes with the tests that always did, 1.7 MB, 9,393 requests, recorded in 6m10s and replayed in 0.3 s.

- **Six rules changed in the diff,** each measured on the pinned server: `ROW_FORMAT = DYNAMIC` is what an undeclared table stores; a CHECK or foreign key OceanBase named is matched by its definition; an index column's `DESC`, an index's `USING`, a table's `ENGINE` and a view's `ALGORITHM` and `SQL SECURITY` are taken and not stored; `PARTITIONS n` is stored as the partitions it stands for.
- **One changed in the generator:** an invisible index is added and then hidden with `ALTER INDEX`, and one that only changes visibility is altered in place.
- **78 probes are recorded gaps** (`obGaps`), the work list. Stored form (21): `BETWEEN` stored as two comparisons, the default collation of `utf8mb4`, two forms the parser cannot read (`BINARY` before `CHARACTER SET`, and a multi-valued index as OceanBase prints it), seven rewrites of a view's body, and two routine comments `SHOW CREATE FUNCTION` prints unescaped. Migrations (57): 20 `ALTER TABLE`s MySQL needs whole and OceanBase takes only in parts, 10 drops of a foreign key's index that OceanBase never made, 12 changes OceanBase has no statement for (`REMOVE PARTITIONING`, a geometry column's SRID, a generated column's expression), 9 triggers the parser cannot read back once named with their database, 3 multi-valued indexes, and 3 others.
- **25 are recorded refusals** (`obRefused`): forms OceanBase will not take.
- **Events** are read from `information_schema.EVENTS`: OceanBase has no `SHOW CREATE EVENT`.

**OceanBase, the spike.** An earlier run with cruder sample tokens, kept for comparison. TiDB's generator produced 3,969 accepted statements and 15,560 near-misses. All were judged on a real server with `PREPARE`, in about 3.5 minutes.

| Statements | Total | Agree | Only engine accepts | Only omni accepts | Server crash |
| --- | --- | --- | --- | --- | --- |
| One per alternative | 3,969 | 2,622 | 1,287 | 57 | 3 |
| Near-misses | 15,560 | 13,860 | 1,296 | 404 | 0 |

- **3,047 disagreements,** while every existing oracle test passes. They are untriaged. `ALTER SYSTEM` and optimizer hints are the largest families, and some come from the spike's crude sample tokens.
- **The grammar alone** agreed with the engine on 90% of statements. Most of the rest is the lexer (531 malformed hints, which live in comments) and 124 `yyerror` checks in C actions. Token rows and overrides (§2.2, §2.3) exist for these.
- **Replay** took about 1 second per package, against 147 seconds (catalog) and about 300 seconds (review) live.

**Snowflake.** A prototype converted the Syntax blocks of 77 docs pages to a `.g4` and generated 1,557 statements, of which omni accepts 64%. The blocks are near-BNF and map directly:

```text
ALTER WAREHOUSE [ IF EXISTS ] [ <name> ] { SUSPEND | RESUME [ IF SUSPENDED ] }
```

```antlr
// https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse
alterWarehouse : ALTER WAREHOUSE (IF EXISTS)? name? (SUSPEND | RESUME (IF SUSPENDED)?) ;
```

Docs are prose, so some blocks need hand fixes: malformed syntax, undefined fragments. Each fix is made in the `.g4`, with a comment citing its page.

## Appendix B. Rule changes

The **Compatibility** section below replaces the `AGENTS.md` Tests bullets on testing against the engine, CI servers, Docker and environment gates. It lands with the implementation.

> **Compatibility**
>
> - CI starts no database server. Engine answers are committed: syntax verdicts in `<engine>/parser/testdata/compat.tsv`, DDL outcomes in `<engine>/catalog/testdata/simulate.tsv`, and each test's exchanges in `testdata/tapes/<Test>.tape`. Only `OMNI_RECORD=1 go test ./<engine>/...` writes them, from the servers in `<engine>/testdata/pins.tsv`.
> - Grammars are the vendor's `.y` or `.g4` at the pinned release tag, or a committed `.g4` converted from vendor docs. Never a community grammar.
> - Agreement only grows. After a fix, run `-update-compat` and commit the rows it flips. A bug no row covers gets one first, labelled `issue:<key>`: add it with `?` in both verdict cells, record, and commit the filled row.
> - Never edit an `expected` cell or a tape by hand.
> - Record when you add a test, row or fixture, change what a test sends, change SQL omni generates, or bump a pin. Nothing else needs a server.
> - A test reaches the engine only through its engine's `testcontainer` package. What it sends depends only on its inputs, never on omni's results, except where omni's generated SQL is under test.
> - A replayable test sorts map iterations, uses its own database, takes names from `t.Name()`, and orders every multi-row read.
> - Replay never skips. `OMNI_RECORD` is the only environment gate.
> - A pin-bump PR changes the pin and the recorded answers only. Fix omni in a separate PR.

Other files:
- **`AGENTS.md`:** `enginetest` joins `metadata` and `review` as shared code, and the pinned images move to `pins.tsv`.
- **`ci.yml`:** no Docker. The Trino and OceanBase lanes and the live sweep sizes go. `compatcheck` is added.
- **`docs/engine-test-matrix.md`:** becomes the §8 table, with agreement counts per engine.

## Appendix C. The tape harness

`enginetest/tape` is about 1,300 lines. An engine describes its server once:

```go
type Server struct {
	Pins  string                                                // the engine's pins.tsv
	Live  func(t testing.TB, pin enginetest.Pin) driver.Connector // called only in a recording
	Codec Codec                                                  // the driver's errors as tape text and back
}

// Open returns t's connection to the server in pin row role: a replay of its
// tape, or, in a recording, the pinned server with every exchange written down.
func Open(t testing.TB, s *Server, role string) *sql.DB
```

Each engine adds a few small things in `<engine>/internal/testcontainer` (TiDB's is the model):

```go
func Open(t testing.TB) *sql.DB { return tape.Open(t, server, "default") }
func Reset(t testing.TB, fn func(context.Context, *sql.DB) error) // in a recording only, unrecorded
```

- **`Codec`** turns the driver's error into tape text and back (`*mysql.MySQLError`, `*pgconn.PgError`), so `errors.As` works in a replay. It may shorten what repeats: OceanBase's writes `…` for the 122 characters that open every syntax error. A recording compares a committed answer as the codec spells it today, so a new spelling is not a different answer.
- **`Renamer`** is a `Codec` for an engine that makes up names, as OceanBase does for a foreign key or a CHECK written without one (`c_OBFK_1794529609863994`, the time in microseconds). A recording keeps a stand-in for each (`c_OBFK_0000000000000001`) in answers and errors, puts the engine's own name back where the test sends the stand-in, and hands the test the stand-in too. So a test reads the same text in a recording as in a replay, and a second recording writes the same tape.
- **`Live`** starts the pinned image, checks its digest and version, inlines arguments, and sets a 30-second read timeout, so a hung statement fails its test and not the run.
- **`Reset`** puts a shared server into the state a tape starts from, such as dropping what earlier tests left. It runs only in a recording, outside the tape, so what a test sends never depends on what ran before it.

What the harness handles:
- **Every kind of answer:** result sets with column types (a BIGINT UNSIGNED value as the `uint64` the driver returns), several result sets, transactions, errors through the codec, and dropped connections. Driver-level prepared statements are refused, which is why arguments are inlined.
- **Streams.** Each `Open` is one stream, named by test, pin role and open order, with a pool of one connection so session state holds. A recording that sends what the committed tape holds, under the same pin, must get the same answers.
- **`tape.Main(m)`** runs the package's tests. After a passing run of every test, a replay fails on a tape or stream no test opened, and a recording deletes such tapes. A run that failed, selected tests (`-run`, `-skip`) or ran none (`-list`, a fuzz worker) settles nothing: it neither fails on those tapes nor deletes them.
- **Filters.** A stream replays in order, so a `-run` pattern that skips some subtests sharing their parent's connection fails replay. Select whole top-level tests, as `-run '^TestName$'` does. A recording rewrites the whole tape of each test it runs, including one that no longer parses, such as a tape left with merge conflict markers.
- **Compact requests.** A request's text or text argument over 64 KB is written as its SHA-256 and length, since a replay only compares it. A request that repeats the text of the one before it writes `^` for it.
- **Not built yet:** `tape.Live(t)` for tests with nothing to replay, and package setup recorded once in `_setup.tape`. No package on tapes has needed them.
- **Other boundaries.** Trino's HTTP API and the programs tests run (ScriptDOM, the Spanner oracle) record at their own boundary, in the same format.
