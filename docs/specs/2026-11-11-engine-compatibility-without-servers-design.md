# Engine compatibility without servers in CI

Status: draft, 2026-11-13. Step 0 is built: `enginetest/compat`, `enginetest/grammar` and TiDB's `compat.tsv` (11,058 rows, recorded on the pinned server in 9 seconds), and `enginetest/tape` with `tidb/catalog` on tapes (131 tapes, 3.1 MB, recorded in 90 seconds; replay takes 1.6 seconds and starts no server), and `tidb/catalog/testdata/simulate.tsv` on one fixture: 1,724 rows, recorded in 83 seconds at 0.05 seconds a row, 1,151 agreeing (Appendix A). OceanBase has step 1: 20,407 rows, judged on the pinned server in three minutes, 19,732 agreeing once #872 to #884 land (Appendix A). Its step 2 is done: `oceanbase/parser`, `oceanbase/validate`, 162 tests of `oceanbase/catalog` and the four oracle tests of `oceanbase/review` replay tapes (253 tapes, 9.7 MB), nothing under `oceanbase/` starts a server outside a recording, and `oceanbase/catalog/testdata/simulate.tsv` holds one fixture: 2,276 rows, recorded in ten minutes at 0.25 seconds a row, 1,842 agreeing once the catalog fixes of Appendix A land (§3). The 162 are the tests that always asked OceanBase, every oracle test written against MySQL, moved with the fixes they needed and with what still fails recorded (Appendix A), and 77 of the 107 container tests: each asserts OceanBase's answer and, under an `omni gap:` comment, what the catalog still answers as MySQL does. The other 30 container tests and 210 subtests compared MySQL's `SHOW CREATE` text and were deleted (D3). `review`'s walk lists the server's objects (§4.4), and its MySQL oracles were pointed at OceanBase or, where they held what only MySQL has, deleted (Appendix A). MySQL has step 1: 12,282 rows of `sql_yacc.yy` at `mysql-8.0.46`, judged on `mysql:8.0.46` in 13 seconds, 11,159 agreeing (Appendix A). Its step 2 is done: `mysql/parser`, `mysql/validate`, `mysql/catalog` and `mysql/review` replay tapes (269 tapes, 16.8 MB) and start no server, and `mysql/catalog/testdata/simulate.tsv` holds one fixture: 1,648 rows, recorded in 45 seconds, 1,047 agreeing (Appendix A). MariaDB has step 1: 12,960 rows from `sql_yacc.yy` at `mariadb-11.8.8`, judged on the pinned server in 20 seconds, 12,082 agreeing (Appendix A). Its step 2 is done: `mariadb/parser`, `mariadb/validate`, `mariadb/catalog` and `mariadb/review` replay tapes (310 tapes, 16 MB), nothing under `mariadb/` starts a server outside a recording or holds a `t.Skip`, and `mariadb/catalog/testdata/simulate.tsv` holds one fixture: 1,452 rows, recorded in under three minutes at 0.1 seconds a row, 986 agreeing (Appendix A). PartiQL has step 1, judged by its reference grammar with no server: 1,910 rows from `PartiQLParser.g4` at `v0.14.9` and DynamoDB's documented statements, 1,510 agreeing (Appendix A). Snowflake has step 1, judged by its docs: 760 command pages converted to a `.g4`, 58,310 rows, and 9,880 of the 15,653 alternatives agree (Appendix A). A spike measured OceanBase and Snowflake first. Its code was lost. The other engines' grammars were located, not run.

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

Until then, the disagreeing rows are the work list: 675 of `compat.tsv` for OceanBase once #872 to #884 land, and 434 of `simulate.tsv`. The guards in §7 make that count only fall, and once it reaches zero any disagreement fails CI.

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
| Differs from both, and the cell was a gap | Fail: a known gap is now another one, as a rejection at another byte is. `-update-compat` commits it too. |
| Differs from the cell, which equalled `expected` | Fail: a regression. |
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
- **Samples resolve.** Where an engine's parser looks a name up while it parses, as TiDB does for charsets and collations, `tokens.tsv` samples a name that exists. An unknown name would answer before a later syntax error, and the row would read as accepted. Where the grammar gives such a name a rule of its own, as MariaDB's `charset_name`, a `tokens.tsv` row names the rule, and other rows write its text where the rule stands.
- **Rule IDs are stable.** A row is named by its rule and the alternative's symbols, not by an index, so a pin bump's diff shows only what the release changed.
- **Token forms come from a committed `tokens.tsv`,** cited to the vendor's lexer or lexical docs. It is written by hand, because the lexers of MySQL, MariaDB and TiDB are code, not grammar files. Other rows take their sample identifiers and literals from it.

**Hand-added rows** carry a label in `rule` and survive regeneration: `issue:<key>` for a bug report, `suite:<name>` for a corpus. Today's corpora become `suite:` rows: TiDB's, Doris and StarRocks conformance, Trino's differential, pg's fuzz, and the statements that completion and analysis tests check for validity.

**Grammar files.** The generator reads two formats, bison `.y` and ANTLR `.g4`.
- A vendor `.y` or `.g4` is used as published. It is committed where the license allows, as TiDB's `parser.y` already is, and fetched otherwise (D1): `<engine>/testdata/sources.tsv` names the file by the release's commit and its SHA-256, and `enginetest.Fetched` downloads it once per machine and checks it on every use. OceanBase's grammar is read this way, and MariaDB's, which holds its default and its `sql_mode=ORACLE` grammar in one file between `%ifdef` lines: `compat.Suite.Defined` names the one to read.
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
| pg, redshift | The `SELECT 1/0` probe, read by `compat.RawParseVerdict`: the server parses the whole text first, so any error but the probe's own division by zero is the parser's, at the position it reports |
| mssql | `SET PARSEONLY ON` |
| oracle | `DBMS_SQL.PARSE` |
| trino, doris, starrocks, googlesql | The engine's error name or message |

An answer the function can't read fails the recording. It is never guessed.

OceanBase needs three things of its own, all measured on the pinned server (`oceanbase/parser/testdata/README.md`):
- **The probe's tail stays inside the statement,** because OceanBase refuses to prepare two statements before it parses either.
- **1064 comes in two kinds.** The parser's says where the error is. A bare "Parse error" comes from later layers, and from `PREPARE` itself for statements it does not carry, so the probe decides it.
- **Two statements crash the server.** The judge starts it again, asks the statement again, and judges one that kills the server twice by its probe alone.

The two recognizers are weaker judges, in known ways:
- **A grammar recognizer misses checks made outside the rules,** in the lexer and in grammar actions. On OceanBase it agreed with the engine on 90% of statements (Appendix A). Each miss found becomes a cited row in `<engine>/testdata/grammar-overrides.tsv`: the verdict, the source that gives it, the cite, the statement. The source is `grammar` where the grammar's own lexer code or action decides, and `docs` where the vendor's docs take a statement the grammar does not derive; a `docs` row only accepts. `TestCompatibility` fails an override no row holds, and one the recognizer already agrees with.
- **A `.g4`'s lexer rules are read for two things:** the text a token spells, and its keywords. A word a token spells before the last token that matches a class of texts is reserved, as ANTLR gives a text two rules match to the first: it is a name only where a parser rule says so.
- **A grammar judge's pin is its files:** the digest in `pins.tsv` is the SHA-256 of the vendor's files as published (`compat.Suite.Vendor`), so a grammar edited by hand fails.
- **Docs only say what is valid.** A near-miss judged by docs is `undocumented`, not `reject`. It counts as neither a gap nor a regression, so omni is never pushed to reject forms users write.
- **Snowflake and Redshift can upgrade.** With `SNOWFLAKE_COMPAT_DSN` or `REDSHIFT_COMPAT_DSN` set, a recording asks the engine instead of the docs. The account must be dedicated and scratch-only, because preparing DDL there executes it.

## 3. DDL behavior: `simulate.tsv`

Omni's catalog simulates DDL. `catalog.Exec` either changes the schema as the engine would, or refuses with the engine's error code. `simulate.tsv` holds the engine's side of that. One row, traced:

```text
fixture one_table@column   CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))
statement                  ALTER TABLE a ADD a MEDIUMINT
engine                     error 1060 (duplicate column)
omni, catalog.Exec         error 1060                      → the row agrees
```

Rows are generated: every statement of the simulated families that the engine accepts in `compat.tsv`, once, on each fixture. Four rows of `oceanbase/catalog/testdata/simulate.tsv`:

```text
# columns: omni, expected, source, fixture, rule, statement
ok:f3b055810eeb	ok:f3b055810eeb	engine	one_table@table	int_type_i/MEDIUMINT	ALTER TABLE a ADD a MEDIUMINT
1060	1060	engine	one_table@column	int_type_i/MEDIUMINT	ALTER TABLE a ADD a MEDIUMINT
ok	ok:f4692f929084	engine	one_table@table	data_type/NCHAR_VARCHAR_string_length_i_opt_binary	ALTER TABLE a ADD a NCHAR VARCHAR(1)
ok	1072	engine	one_table@table	sort_column_list/sort_column_key	ALTER TABLE a ADD KEY(a)
```

1. An agreement: both take the statement, and omni stores the schema OceanBase stores.
2. An agreement on a refusal: the same error code.
3. A gap: both take it, and omni stores another schema.
4. A gap: omni takes what OceanBase refuses.

- **Fixtures** are committed starting schemas, `<engine>/catalog/testdata/fixtures/<name>.sql`. OceanBase has the first, one table with keys.
- **Names are filled in the fixture, not the statement.** A `compat.tsv` statement writes `a` for every name it holds, so it runs as written and the fixture decides what it finds. `one_table` comes in four namings: `@missing` names nothing `a`, `@table` the table, `@column` a column of it too, `@index` an index of it too. There is no `references.tsv`. The cost: a statement whose names must differ, such as `CREATE TABLE a(a TINYTEXT, a TINYINT)`, is refused on every naming. OceanBase takes 250 of the 569 statements on at least one.
- **Families.** The engine's test lists the leading words it simulates. OceanBase's are the statements on tables, indexes and views, whose whole effect is in the database a row reads back. Databases, sequences, materialized views and server objects are not simulated yet.
- **Outcomes** are an error code, or `ok:<hash>`: the first 12 hex digits of the SHA-256 of every table's and view's `SHOW CREATE`, in name order, with the names the engine makes up counted in order. Omni's cell holds the engine's hash when its catalog stores the same schema, and a bare `ok` when it takes the statement and stores another, or the engine refuses.
- **The same schema** is the engine test's to decide. OceanBase's loads the `SHOW CREATE` texts into a second catalog and takes `Diff` in both directions, the test every oracle of `oceanbase/catalog` uses: the catalog prints MySQL's text, which equalled OceanBase's for 2 of 43 schemas measured.
- **The conversation is on a tape** (§4), `testdata/tapes/TestSimulation.tape`. Per row: the database dropped and built from the fixture, the statement, and after a statement the engine took, each object's `SHOW CREATE`. None of it reads omni. A replay works the engine's outcome out of the tape and fails a row whose `expected` cell is anything else, so neither file is edited without the other.
- **Everything else is §2.1,** with one more move: the same columns plus `fixture`, the same checker (`enginetest/simulate`), `-update-compat`, and the guards in `compatcheck`. A syntax verdict is right or wrong; an outcome can be wrong in many ways, so a fix can take a row from one difference to another, as a parser fix turns omni's 1064 into the catalog's answer. That fails the test until `-update-compat` commits it. Only a row that agreed and no longer does is a regression. The test is `TestSimulation`. `-simulate-show=<text>` prints the fixture, the engine's schema and what omni would change to reach it, for the rows that hold the text.
- **A recording runs every row,** since a tape is written whole: a new row or fixture costs the full run. An answer that differs from the committed one at the same pin fails the recording.

**Walk-through is to read the same rows,** and does not yet. For each row, the walk runs on the fixture and the statement:
- On an `ok` row, any finding fails.
- On a refused row, a finding must name the engine's code, or claim only facts true of the fixture.
- Silence passes, since the walk must never be wrong but may be unsure. How often it finds the refusal is printed.

**Cost, measured:** 0.25 seconds a row on OceanBase, so its 2,276 rows record in ten minutes (Appendix A). The 5 seconds a row this section first guessed would have made 1,000 statements × 5 fixtures × 3 names a 20-hour recording; at the measured rate it is under an hour. Docs engines have no `simulate.tsv`, since docs rarely name the error a statement raises.

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

**Walk-through lists the server's objects to follow the rule.** When a finding doesn't name the server's error code, a `walktest.Holder` queries the server about the objects the finding claims, so a change to the findings changes what the test sends. A `walktest.Lister` lists the database's objects before every statement instead, and `Judge` reads the claim off that list:

```go
type Lister interface {
	Inventory(t testing.TB, change string) []Object // fixed queries, whatever the walk says
	Same(claimed, stored walk.Key) bool             // the engine's name folding
}
```

The list replaces `Digest` too: a statement the server rejects must leave it as it was. OceanBase's walk is a `Lister`, on tapes. The other engines' are `Holder`s and run live until their step 2. A sweep on tape has a fixed size and seed, constants beside its test: its tape holds one stream per change, so another size or seed fails replay.

### 4.5 Writing a replayable test

- Sort map iteration.
- Use the test's own database. Read no server-wide state outside package setup.
- Take names from `t.Name()`, not from counters or the clock.
- Order every multi-row read.
- Don't skip. Today's skips and DSN fallbacks go. A probe the engine refuses, or one omni fails, is recorded with what it must do instead: `obRefused` and `obGaps` in `oceanbase/catalog/oceanbase_oracle_test.go` each fail an entry that no longer holds.
- Keep a test on one connection. A tape's stream is one, and a helper that asks the pool while its test holds that connection waits forever.
- Read the stored form from the database the change is applied to. A name the engine makes up differs in every database.
- Probe syntax with SQL's `PREPARE ... FROM '<text>'`, not the driver's `Prepare`: a tape holds no driver-level prepared statement.
- Keep an id the server makes up out of the tape: a connection's, a table's. Resolve it inside the query (`CONNECTION_ID()`, `IS_USED_LOCK`), or give it a stand-in through the `Renamer`.
- Read what a statement changes at once, not a view that settles later: OceanBase's `information_schema.EVENTS` keeps a dropped database's events for about 20 seconds.
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
  simulate/                          the row checker of simulate.tsv
  tape/                              record and replay
  cmd/compatcheck/                   the CI guards that compare a PR's files with its base

oceanbase/
  testdata/pins.tsv                  the pinned server
  testdata/tokens.tsv                token forms and sample values
  testdata/sources.tsv               vendor files fetched, not committed: name, SHA-256, URL
  testdata/grammar-overrides.tsv     where a grammar judge is wrong, cited (§2.3)
  parser/testdata/compat.tsv
  catalog/testdata/simulate.tsv
  catalog/testdata/fixtures/         starting schemas
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
| A new row or fixture | **Yes:** judge it; for `simulate.tsv`, record every row | No: the recognizer judges |
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
| A regression is committed as the new `omni` | CI fails if an `omni` cell that equalled its row's `expected` changes. |
| `expected` is edited to pass | CI fails if an `expected` cell changes, unless it was `?` or the PR bumps a pin. |
| A `grammar` or `docs` cell is written by hand | `TestCompatibility` asks the recognizer again on every run. |
| A vendored grammar is edited | Its SHA-256 is the pin's digest (§2.3). |
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
| partiql | `PartiQLParser.g4` @ `v0.14.9` of partiql-lang-kotlin | grammar; DynamoDB's docs where they take more |
| cosmosdb | `sql.g4` | grammar |
| mongo | mongosh method docs → `.g4` | docs |
| elasticsearch | `elasticsearch-specification` → `.g4` | grammar |

- **TiDB is the starting point.** Its `parser.y` rows were the first `compat.tsv`; its other corpora became `suite:` rows in step 1.
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
- **D6. Queries.** Stored form proves omni reads DDL as the engine does. The same check reaches queries by wrapping each accepted `SELECT` row in `CREATE VIEW`: the engine's stored view text is its own reading of the query. It works where omni's catalog renders views (pg, the MySQL family), and it multiplies the size of `simulate.tsv`. Recommended: decide once step 2 has measured what a row costs. It has: 0.25 seconds on OceanBase (Appendix A), so a view for each accepted `SELECT` row is minutes of recording, not hours.

## Appendix A. Measured results

**OceanBase, recorded.** `oceanbase/parser/testdata/compat.tsv` at `v4.3.5_CE_BP6`: 5,123 alternatives, of which 52 need a token the lexer never produces, and 15,336 near-misses. The pinned server judged all 20,407 rows in three minutes and accepts 4,995 of the 5,071 alternatives.

| Statements | Total | Agree | Only engine accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 5,070 | 4,665 | 402 | 3 |
| Near-misses | 15,337 | 15,067 | 76 | 194 |

- **675 disagreements once #872, #877, #881, #883 and #884 land** (777 on main), down from the 2,679 first recorded (2026-11-13) and 1,118 before this round. This round closed: mock tablespace, server, logfile group and instance statements (79), outlines, snapshots, limiting and sensitive rules (64), privileges and roles in `GRANT` and `REVOKE` (51), `SHOW` (34), column groups (33), `GET DIAGNOSTICS` (28), transactions, locks, savepoints and XA (24), subpartition actions and name lists (24), `ALGORITHM`, `LOCK`, `ROW_FORMAT` and `INSERT_METHOD` values (21), `ANALYZE` (20), column attributes (17), `RANGE`'s empty and automatic partitioning (15), collection and other types (13), and 20 by other changes. Before it: names with a leading dot or a reserved word after a dot (585), `ALTER SYSTEM` (417), table options (109), `tbl_name`'s clauses (109), tenants and resources (104), tablegroups and sequences (85), `AUDIT` (71), `EXPLAIN` (43), and 38 more.
- **What is left,** by what the statement is:

  | Rows | Omni rejects | Omni accepts | What | Why it is left |
  | --- | --- | --- | --- | --- |
  | 168 | 150 | 18 | `SELECT` | Many small rules, none over 14 rows: function arguments (`func_expr`), window functions, `UNNEST`, `RB_ITERATE`, `JSON_TABLE` columns, natural join types, grouping sets, a flashback on a derived table |
  | 101 | 101 | 0 | Materialized views and their logs | Not read at all; a view's query and columns want a tree, not a matcher |
  | 101 | 0 | 101 | Function calls | A keyword of OceanBase's that names no function (`SELECT gts(a)`), and a keyword function called with no argument (`SELECT date()`): needs the callable keywords measured on the server |
  | 77 | 70 | 7 | External tables and files | `SOURCE()`, `FILES()` and a location as a table, `DELETE FROM` one, `ALTER EXTERNAL TABLE`, `LOAD DATA`'s OceanBase options, `INTO OUTFILE ... PARTITION BY` |
  | 51 | 51 | 0 | `FOR COLUMNS`, `FOR ALL` | Not statements: `method_opt` is the grammar of a `DBMS_STATS` argument, which the parser also reaches as a `stmt`. Left as they are |
  | 32 | 25 | 7 | `ALTER TABLE` | 9 rows OceanBase parses and then refuses with 1064 (a column's `ID n`, `SIZE ('s')` in `ALTER TABLE`, a string as a type): taking them waits for the catalog to refuse them too, or `simulate.tsv` regresses. The rest: `ADD (KEY ...)` in parentheses, `ALTER INDEX` options, `REMOVE TTL`, `JSON (c) STORE AS`, `BEFORE` |
  | 11 | 10 | 1 | `DELETE` and `INSERT` | `DELETE FROM TABLE(a)`, `INSERT OVERWRITE`, a repeated `QUICK` |
  | 3 | 3 | 0 | `SET ROLE ALL EXCEPT ALL` | The statements that crash the server (§2.3) |
  | 2 | 0 | 2 | `CREATE TABLE` | `FOREIGN (a)` without `KEY` |
  | 129 | 68 | 61 | Everything else | No group over ten rows: `FLASHBACK`, `PURGE`, `CACHE INDEX`, users and `SET PASSWORD`, `HANDLER`, databases, `DUMP`, `EXPLAIN`'s forms |
- **A carrier can hide its alternatives.** A row is the shortest statement that reaches an alternative, and where omni rejects that carrier, every alternative behind it disagrees for one reason. Fixing the carrier is what shows the rest.
- **Two statements crash the server** (three rows). Their verdict comes from the probe (§2.3).

**MariaDB, recorded.** `mariadb/parser/testdata/compat.tsv` at `mariadb-11.8.8` (2026-11-13): 3,553 alternatives of the default grammar, of which 82 need a token the lexer never produces under the default `sql_mode` (most are Oracle mode's), and 9,489 near-misses. The pinned server judged all 12,960 rows in 20 seconds, and two recordings on fresh servers left the same file.

| Statements | Total | Agree | Only engine accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 3,471 | 3,025 | 312 | 134 |
| Near-misses | 9,489 | 9,057 | 92 | 340 |

- **878 disagreements,** untriaged. By statement: expressions carried by `DO` (115), `CHANGE MASTER` (57), `ALTER EVENT`, which carries the routine statements (51), `ALTER SERVER` (48), compound statements at the top level (`BEGIN NOT ATOMIC`, 39), `ALTER TABLE` (37), `DELETE` (34), `GET DIAGNOSTICS` (25), `START` (22), `RESIGNAL` (21), `CREATE PACKAGE` (21).
- **`PREPARE` carries every statement kind but its own** (`PREPARE`, `EXECUTE`, `DEALLOCATE`: 1295), which the probe decides (`mariadb/parser/testdata/README.md`).
- **The grammar's actions refuse and look up while the parser runs.** A refusal is a rejection: an identifier as a data type (4161, 731 rows), a routine variable outside a routine (1327, 315), a table named twice (1066, 115). A lookup of a name the server does not have (a character set, collation, system variable, engine or plugin: nine codes) counts as read, as in the package's hand-written cases, and answers before a later syntax error, so `tokens.tsv` gives nine rules a name that exists; 163 rows are still answered that way, most of them `SET GLOBAL a = a` and its near-misses.

**MariaDB, parser and validate moved.** All 117 tests of `mariadb/parser` that asked a server replay (2026-11-13): 87 tapes, 5.5 MB, 53,222 requests, recorded in 45 to 105 seconds and replayed in 25. `mariadb/validate` has one tape of 52 requests. Two recordings of each on fresh servers left the same tapes.

- **The tests moved as written.** They always asked MariaDB, so no expectation changed. Two sweeps walked the keyword map in its own order and now walk it sorted; a verdict is one `PREPARE ... FROM ?` request where it was three.
- **Each test starts on an empty server:** a recording drops what the image did not start with before a test's connection opens (`testcontainer.OpenEmpty`).
- **No skip is left in either package.** Two were dead: a per-case skip no case set, and a skip for a corpus directory that is there.

**MariaDB, catalog moved.** All 215 tests of `mariadb/catalog` that asked a server replay (2026-11-13), a tape each: 3.0 MB, 23,570 requests, recorded in 35 to 65 seconds and replayed in 6. Two recordings on fresh servers left the same tapes.

- **Every connection's clock is stopped** (`SET timestamp`, below the tape): the first two recordings differed in a trigger's `Created` and in the `STARTS` of an event that names none.
- **A test that took a fresh session per case opens a connection per case:** a tape's pool is one connection, and a temporary table or `foreign_key_checks = 0` of one case reached the next.
- **Three map walks sort, and one read is ordered:** a view whose `ORDER BY` ties came back in either order.
- **`information_schema`'s unsigned columns arrive as `uint64`** once arguments are inlined; two helpers read only `int64`.
- **The 78 skips are gone.** 72 were in the view deparse tests: 56 for cases marked partial, which no case is any more, and 16 for a view MariaDB refuses to create, which three cases are (`a->'$.key'` twice, 1064; two columns of one name, 1060). Those three are recorded (`deparseRefusedViews`) and fail when the refusal changes. The other six never fired and now fail.

**MariaDB, review moved.** Its six oracle tests replay (2026-11-13): 6 tapes, 6.6 MB, 13,751 requests, recorded in nine minutes on a loaded box and replayed in seconds. Two recordings on fresh servers left the same tapes.

- **The walk is a `Lister`.** It was a `Holder` and a `Digester`, whose queries followed the walk's findings. Its inventory is three queries before every statement, whatever the walk reports: the user databases, what they hold, and the change's session's temporary tables. Each answers one row, the lines in byte order, since `information_schema` compares names without case and returns tied rows in either order; an index is read by its first column, since a `DISTINCT` took the indexes `x` of tables `b` and `B` for one.
- **It finds 244 of the 296 statements MariaDB rejects in the curated cases and 354 of 547 in the sweep,** with no false finding. The sweep is 40 changes from seed 20261026, constants beside it; `WALK_SWEEP_N`, `WALK_SWEEP_SEED` and `WALK_SWEEP_IMAGE` and its shards are gone.
- **Each case starts on a server that holds no user database,** dropped in a recording below the tape with the global settings a change may leave (`foreign_key_checks`, `system_versioning_alter_history`).
- **The schema sync's reads are ordered by the bytes of their names.**

**MariaDB, `simulate.tsv`.** Recorded 2026-11-13 on the pinned server: the 363 table, index and view statements MariaDB accepts in `compat.tsv`, on the four namings of `one_table` (OceanBase's fixtures, copied). 1,452 rows in 148 seconds, about 0.1 seconds a row; the tape is 0.6 MB, and two recordings on fresh servers left the same two files.

| Fixture | Rows | MariaDB takes | Agree | Only MariaDB takes | Only omni takes | Both take, another schema | Both refuse, another code |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `one_table@missing` | 363 | 6 | 297 | 5 | 10 | 0 | 51 |
| `one_table@table` | 363 | 163 | 196 | 12 | 91 | 9 | 55 |
| `one_table@column` | 363 | 117 | 246 | 12 | 37 | 6 | 62 |
| `one_table@index` | 363 | 112 | 247 | 12 | 36 | 6 | 62 |
| All | 1,452 | 398 | 986 | 41 | 174 | 21 | 230 |

- **The same schema is decided by printing.** MariaDB's catalog has no `Diff`, so the test loads MariaDB's `SHOW CREATE` texts into a second catalog and compares what the two catalogs print; the catalog prints MariaDB's text, which the container tests already hold it to. omni stores MariaDB's schema on 336 of the 398 rows MariaDB takes.
- **466 rows disagree,** untriaged. Most are codes: omni's 1051 where MariaDB answers 1347 or 1146 (52), omni's parser 1064 on 116 rows of 29 statements (`ADD IF NOT EXISTS` among them), and omni taking what MariaDB refuses for a missing key column (1072, 38), a partition change to a table that has none (1505, 24), and an index or a key it cannot build (1911, 1283, 1005).
- **MariaDB makes up no name that changes between runs** (`a_ibfk_1`, `CONSTRAINT_1`), so the test needs no `Renamer`.
- **153 of the 363 statements are refused on all four namings,** as on OceanBase: a partitioned fixture and names that differ within a statement are what reach them.

**OceanBase, catalog's MySQL tests.** Measured 2026-11-12 by pointing them at the pinned OceanBase as written. Nothing of the run is committed.

| Tests | Top-level | Pass | Subtests | Pass | Fail | Skip |
| --- | --- | --- | --- | --- | --- | --- |
| Stored form and migrations (`connectOracle`) | 56 | 20 | 708 | 295 | 357 | 56 |
| `SHOW CREATE` text and scenarios (`startContainer`) | 106 | 15 | 873 | 335 | 537 | 1 |

- **Stored form.** 184 failures are a table that differs from its own readback. Where traced, the cause is OceanBase printing `ROW_FORMAT = DYNAMIC` on every table. 45 are views that differ from theirs. 76 are readbacks omni can't load: 47 foreign keys whose referenced table comes back qualified by its database, 20 triggers, 8 multi-valued indexes printed as `json_query(... RETURNING ...)`, and a routine.
- **Migrations.** 41 generated statements are refused: 14 for combining changes OceanBase takes one at a time (1235, 4179), 7 for dropping a key that is not there (1091), 6 for altering a geometry column's SRID (1235), 5 partition statements as syntax errors, and 9 others. 4 more apply and leave the index visible: `ALTER TABLE ... ADD KEY` does not take `INVISIBLE`, written plainly or as `/*!80000 INVISIBLE */`, where `CREATE TABLE`, `CREATE INDEX` and `ALTER INDEX` do.
- **Text.** The 48 `TestContainer_*` tests mostly compare `SHOW CREATE TABLE` with MySQL 8.0's text, which OceanBase never prints: 301 of their 378 subtests fail. The 27 view tests do the same for view text: 131 of 204 fail. The 25 scenario tests are closest: 163 of 241 pass.
- **The hangs #808 recorded were the tests' own.** A helper asked the pool for a second connection while its test held the only one. The server sat idle each time, and the cross-object ordering test passes 5 of 5 once the pool has a second connection.
- **Not measured then:** the keyword sweep built an index per keyword and did not finish one of its seven positions in five minutes (it takes 11 seconds now that it asks with `PREPARE`), and a tape could not hold a `uint64`, which `information_schema`'s unsigned columns return (it can now).

**OceanBase, catalog's oracle tests moved.** All 53 now ask OceanBase (2026-11-13; three more read MySQL's `sys` schema and are gone): 80 tapes with the tests that always did, 1.7 MB, 9,393 requests, recorded in 6m10s and replayed in 0.3 s.

- **Six rules changed in the diff,** each measured on the pinned server: `ROW_FORMAT = DYNAMIC` is what an undeclared table stores; a CHECK or foreign key OceanBase named is matched by its definition; an index column's `DESC`, an index's `USING`, a table's `ENGINE` and a view's `ALGORITHM` and `SQL SECURITY` are taken and not stored; `PARTITIONS n` is stored as the partitions it stands for. A seventh since: a multi-valued key, `CAST(j->'$.ids' AS UNSIGNED ARRAY)`, is stored as `json_query(... ASIS ...)`, and the two are one index.
- **The generator writes what OceanBase runs.** Its plans were MySQL's, and 42 probes failed on them. Each rule below was measured on the pinned server and is held to it by `TestOracle_OceanBaseAlterTableClauses` and `TestOracle_OceanBaseForeignKeyIndex`:
  - One clause per statement. OceanBase refuses most pairs of `ALTER TABLE` clauses (1235, 4179, 1176) and lets an `AUTO_INCREMENT` column stand without a key, which is what MySQL groups clauses for. The one pair it takes is `DROP PRIMARY KEY, ADD PRIMARY KEY`, and a changed key is that statement: the `DROP` alone is refused on a table a foreign key references (4714).
  - Order does the rest: `AUTO_INCREMENT` leaves one column before another takes it (1075), and a column is dropped or made nullable after the key that holds it is replaced (1235, 1171).
  - Partitioning cannot be removed, so a repartition that lets a column go runs before the key changes, and the column drops after them.
  - A foreign key has no index of its own. While it lives it holds the last index over its columns (5320), so the plan drops it and adds it back around that index's drop.
  - A spatial column becomes `NOT NULL` only from nullable (1138), a generated column's expression changes by dropping and adding the column (1235), and an invisible index is added and then hidden with `ALTER INDEX`.
  - A multi-valued index is built with its table only (1235 afterwards), so a new table carries it in its `CREATE TABLE`.
- **6 probes are recorded gaps** (`obGaps`), of 16: the default collation of `utf8mb4` (OceanBase's is `utf8mb4_general_ci`, the catalog's `utf8mb4_0900_ai_ci`; undecided), `WEIGHT_STRING`'s five-argument readback, which the parser does not read, two routine comments `SHOW CREATE FUNCTION` prints unescaped, and 2 of events. Closed since: `BETWEEN` over constant bounds and nine views, whose bodies the deparser now writes as OceanBase stores them (`WITH RECURSIVE`, each side of a set operation in parentheses, `CASE` for `IF`, `= any` for `IN (subquery)`, `TRIM(BOTH ...)`, a column list as aliases). The parser's 19 are closed: it reads `BINARY` before `CHARACTER SET`, `JSON_QUERY`'s clauses, and a trigger named with its database.
- **10 `omni gap:` comments are left of 59,** and 5 of the 6 `viewNameColumnGaps`. Seven are the parser's: `LINEAR`, `COALESCE PARTITION`, `EXCHANGE PARTITION` without `WITHOUT VALIDATION`, `TABLESPACE=`, a `DEFAULT` list partition and `SHOW CREATE EVENT`. Two are the hidden column of a functional index, which OceanBase names `SYS_NC<column id>$` and the catalog, which keeps no column ids, `!hidden!<index>!0!0`. The view names are an unaliased `(1)`, which keeps its parentheses, and the column order of a `NATURAL` or `USING` join over a join. The 49 closed are the catalog answering as OceanBase does, each rule measured on the pinned server and held to it and to the catalog by `TestOracle_OceanBaseKeyErrors`, `Refusals`, `ColumnTypes`, `IndexTypes`, `ViewAttributes` and `ViewColumnNames`: made-up constraint names, foreign-key and primary-key errors, the statements it refuses, stored column types and character sets, and what it keeps of a view.
- **36 are recorded refusals** (`obRefused`): forms OceanBase will not take. 11 of them are changes it has no statement for (`REMOVE PARTITIONING`, a spatial column's SRID, a generated column's expression under a constraint, a multi-valued index on a table that exists): the plan keeps the statement MySQL would run and warns of it, and the test checks that OceanBase refuses what the plan warns of.
- **Events** are read from `information_schema.EVENTS`: OceanBase has no `SHOW CREATE EVENT`.

**OceanBase, review moved.** Its four oracle tests replay (2026-11-13): 4 tapes, 2.1 MB, 8,140 requests, recorded in 12 minutes and replayed in 0.3 s. Two recordings on fresh servers left the same tapes.

- **The walk** runs 67 curated cases and a sweep of 40 changes, on three connections each: the setup, the catalog, and the change. It finds 223 of the 305 statements OceanBase rejects in the cases and 359 of 617 in the sweep, with no false finding. A walk made to report nothing replays the same 5,788 requests.
- **Its inventory is one query.** Temporary tables come from `oceanbase.__all_table`, by the id of the change's session, which the catalog connection finds through a named lock (`IS_USED_LOCK`). Events come from the scheduler's jobs, by database id. A materialized view's storage table, `__mv_container_<table id>`, gets a stand-in.
- **MySQL's walk oracles are gone:** `mysql:8.0` at `lower_case_table_names` 0 and 1, `mysql:8.4`, and their sweep ran the cases the OceanBase test runs. The pinned server has one tenant, at 1, and its 5 GiB are all assigned, so a tenant at 0 does not fit. Four cases whose setup OceanBase refuses (1193, 1422) were passed over before and are recorded as refused now.
- **Scan rules:** 12 of the 15 cases run whole. OceanBase stops the other three (1235, 1091, 4028), which each case records. Three cases were added: two in the forms OceanBase takes, and the `DROP CONSTRAINT` of a unique key, which it refuses (3940).
- **Column drops:** OceanBase leaves the table its key in all 118 cases: 70 stop at the drop (1235, `drop rowkey column is not supported`), 41 at another error before it, and 7 setups are refused. `ReviewOceanBase` reports nothing on any, since it judges no `DROP COLUMN`.
- **Prior backup** was deleted: Bytebase takes none on OceanBase, and `ReviewOceanBase` answers the request before the rule that test held to MySQL runs.

**OceanBase, `simulate.tsv`.** Recorded 2026-11-13 on the pinned server: the 569 table, index and view statements OceanBase accepts in `compat.tsv`, on the four namings of `one_table`.

| Fixture | Rows | OceanBase takes | Agree | Only OceanBase takes | Only omni takes | Both take, another schema | Both refuse, another code |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `one_table@missing` | 569 | 4 | 471 | 1 | 0 | 0 | 97 |
| `one_table@table` | 569 | 200 | 455 | 22 | 6 | 7 | 79 |
| `one_table@column` | 569 | 122 | 458 | 13 | 6 | 6 | 86 |
| `one_table@index` | 569 | 124 | 458 | 13 | 6 | 6 | 86 |
| All | 2,276 | 450 | 1,842 | 49 | 18 | 19 | 348 |

- **Cost: 0.25 seconds a row.** The 2,276 rows recorded in 9m43s on one connection, the server's start included, and in 10m58s a second time on a busier box. A spike of 300 rows first measured 0.28 seconds: 0.21 to drop the database and build the fixture, 0.05 for the statement, 0.015 to read the schema back. The tape is 0.96 MB, 5,447 requests, and the test replays in 0.3 seconds. Two recordings on fresh servers left the same two files.
- **434 rows disagree** once the catalog fixes below land, down from 1,125 when they were triaged (2026-11-13) and 1,259 at the recording. The table counts that state. 381 of the 434 are statements omni's parser rejects (1064): 98 statements, 49 of them rows OceanBase takes, which are `compat.tsv`'s work list. The catalog's rows were fixed cause by cause, each held to the engine by an oracle test:
  - ALTER TABLE options OceanBase refuses, checked before the other commands (201 rows, #850)
  - the query of a view or CREATE TABLE ... SELECT, resolved before the name (132, #854)
  - partition names and commands (91, #861)
  - a key's columns: missing (1072) or not text for FULLTEXT (1283) (71 at triage, 42 once #887 took the rows both close, #858)
  - a key's parts and block size, before its name (46, #887)
  - a column's definition, before its name (45, #871)
  - a CREATE TABLE with no visible column (44, #864)
  - names looked up in OceanBase's order (30, #880)
  - index options only a vector index or a restore takes (27, #865)
  - TTL columns (21, #873)
  - CHECK expressions (12, #876)
- **53 catalog rows are left,** each for a reason:
  - 19 store another schema: the catalog applies no repartition (`PARTITION BY KEY()`, 15), plus a SERIAL and a GEOMETRYCOLLECTION column and a STORING index.
  - 9 are DROP CHECK, DROP CONSTRAINT and ALTER CHECK of a missing CHECK, where OceanBase answers 3821 or 3940 by the keyword written, which the AST does not keep.
  - 15 are unknown charsets, collations, tablegroups and tablespaces (1115, 1273, 5151, 3510), whose lists in OceanBase are not measured.
  - 7 are partition definitions: duplicate partition names (1517), and RANGE with no partitions (1235 before 1503).
  - 3 are orders: a foreign key to the table itself (1215), and STORING on a missing key column (1072).
- **OceanBase takes 250 of the 569 statements on some naming,** 450 rows, and omni stores its schema on 382. The other 319 statements are refused on all four: a `SELECT a` with no table behind it (1054, 57 statements), forms the release does not support (1235, 42), partition changes to a table that has none (1505, 33), an external table's options (4016, 31). A partitioned fixture and names that differ within a statement are what reach them.
- **A refused statement changed nothing** in the spike, 0 of 237. The file does not hold that: every row builds its fixture again.
- **Text is not the test.** OceanBase's `SHOW CREATE TABLE` equalled the catalog's in 2 of the 43 schemas the spike took; `Diff` against the loaded text was empty in 39.

**TiDB, `simulate.tsv`.** Recorded 2026-11-13 on the pinned server: the 431 table, index and view statements TiDB accepts in `compat.tsv`, on the four namings of `one_table`. Only the generated rows: hand rows (`suite:`, `issue:`) are corpora, left out so the file does not move when they grow. A statement names everything `ident`, as `tokens.tsv` writes an identifier, so the fixtures name `ident`.

| Fixture | Rows | TiDB takes | Agree | Only TiDB takes | Only omni takes | Both take, another schema | Both refuse, another code |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `one_table@missing` | 431 | 3 | 391 | 0 | 24 | 0 | 16 |
| `one_table@table` | 431 | 204 | 243 | 0 | 112 | 4 | 72 |
| `one_table@column` | 431 | 105 | 257 | 0 | 82 | 4 | 88 |
| `one_table@index` | 431 | 111 | 260 | 0 | 82 | 4 | 85 |
| All | 1,724 | 423 | 1,151 | 0 | 300 | 12 | 261 |

- **Cost: 0.05 seconds a row.** The 1,724 rows recorded in 83 seconds, the server's start included, and in 141 seconds a second time on a busier box; a spike of 100 rows took 6. The tape is 0.69 MB, 4,287 requests, and the test replays in 0.1 seconds. Two recordings on fresh servers left the same two files: no answer varies, so no name is made up.
- **The same schema is the same text.** The catalog prints TiDB's `SHOW CREATE`, which the package's differential tests hold it to, so a row agrees when every table's and view's text is equal. Omni stores TiDB's schema on 411 of the 423 rows TiDB takes.
- **573 rows disagree,** untriaged, and none is the parser's: omni reads every statement. 300 omni takes and TiDB refuses, 172 of them with 8200 (62 statements TiDB parses and does not run, such as adding a `UNIQUE` or `AUTO_INCREMENT` column), 30 with 1105 and 21 with 1505. Of the 261 refused with another code, 84 are partition changes to a table that has none, which omni refuses with no code and TiDB with 1505; 45 are omni's 1050 for TiDB's 1051 (a view named as the table, over `SELECT *` with no table); 26 omni's 1060 for TiDB's 8200; 24 omni's 1050 for TiDB's 1113 (`CREATE TABLE ident` with no columns). The 12 other schemas are `ALTER TABLE ident CACHE` and three `PARTITION BY KEY()` forms, on each naming with the table.
- **TiDB takes 234 of the 431 statements on some naming.**

**MySQL, recorded.** `mysql/parser/testdata/compat.tsv` at `mysql-8.0.46` (2026-11-13): 3,173 alternatives, of which 10 need a token the lexer never produces under the default `sql_mode`, and 9,119 near-misses. The pinned server judged all 12,282 rows in 13 seconds and accepts 2,928 of the 3,163 alternatives. Two recordings wrote the same file.

| Statements | Total | Agree | Only engine accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 3,163 | 2,960 | 77 | 126 |
| Near-misses | 9,119 | 8,199 | 322 | 598 |

- **1,123 disagreements,** untriaged. omni rejects 399 statements MySQL takes; 286 are `DO`, and 257 of those are one gap, an alias after a `DO` item (`DO a + a b`), which `select_item_list` allows. The rest are expressions: `MATCH ... AGAINST`, `CAST(... AT TIME ZONE ...)`, `WEIGHT_STRING`'s four arguments, `ST_COLLECT(DISTINCT ...)`, `{ident expr}`.
- **omni accepts 724 statements MySQL refuses:** replication (`CHANGE MASTER` and `CHANGE REPLICATION`, 129; `START`, 45), `ALTER INSTANCE` (56), `GET DIAGNOSTICS` (53), `ALTER EVENT` (49), `ALTER SERVER` (41), `ALTER USER` (26), `RESIGNAL` (19). Many are refused by a grammar action, not by a rule: `GET DIAGNOSTICS` outside a routine, `COMMIT` in a function, `ALTER EVENT` with no clause.
- **Names are `utf8mb4`,** a character set, since MySQL looks one up while it parses. 114 rows still answered a lookup (71 a system variable, 40 a collation, 3 a keyword written as a character set), which reads as accepted; 29 are near-misses whose break the lookup may hide. Written `a`, 34 more near-misses hid behind `CHARSET a`.

**MySQL, moved.** Every test of `mysql/` that asked a server replays (2026-11-13): `mysql/parser` 70 tapes, 6.8 MB; `mysql/validate` one; `mysql/catalog` 182, 4.6 MB, 27,353 requests, recorded in 140 to 200 seconds and replayed in 8; `mysql/review` 7, 5.3 MB, 21,135 requests, recorded in four minutes and replayed in 2. Two recordings on fresh servers left the same tapes, and every package passes with `TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE` pointing nowhere.

- **Every session's clock is pinned** at 2026-01-01 by the live connector, below the tape: a trigger's `Created`, an event's implicit `STARTS` and a table's `CREATE_TIME` are the session's time, and differed between recordings.
- **The review walk is a `walktest.Lister`** on three pinned roles: `default`, `lctn1` (the same image initialized with `lower_case_table_names=1`, fixed at initialization, so a role of its own) and `mysql84` (`mysql:8.4.13`). MySQL lists no session's temporary tables, so `Inventory` asks the change's own session `SHOW CREATE TABLE` for each name the change writes after `TEMPORARY TABLE`, in each database. The walk finds 158 of the 184 statements MySQL rejects in the cases at 0 and on 8.4, 161 of 187 at 1, and 398 of 592 and 389 of 588 in the sweep of 40 changes, with no false finding.
- **151 skips are gone,** and none is left in `mysql/`. Two fired: a view the catalog now prints as MySQL does, whose case asserts it again, and a view MySQL refuses (two columns named `a`, 1060), recorded with the catalog's acceptance as an `omni gap` that fails when either side moves. The rest never fired, a setup or readback the pinned server always gives or a gap since closed, and each is a failure now.
- **Map iteration** sent requests in five tests, and a database name was salted with the process id; both are fixed.

**MySQL, `simulate.tsv`.** Recorded 2026-11-13 on the pinned server: the 412 table, index and view statements MySQL accepts in `compat.tsv`, on the four namings of `one_table`. A statement names everything `utf8mb4`, as `tokens.tsv` writes an identifier, so the fixtures name `utf8mb4`. 1,648 rows in 45 seconds; the tape is 0.7 MB, 4,103 requests, and two recordings left the same two files.

| Fixture | Rows | MySQL takes | Agree | Only MySQL takes | Only omni takes | Both take, another schema | Both refuse, another code |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `one_table@missing` | 412 | 6 | 319 | 4 | 28 | 1 | 60 |
| `one_table@table` | 412 | 191 | 218 | 2 | 105 | 11 | 76 |
| `one_table@column` | 412 | 104 | 253 | 0 | 49 | 8 | 102 |
| `one_table@index` | 412 | 108 | 257 | 0 | 48 | 8 | 99 |
| All | 1,648 | 409 | 1,047 | 6 | 230 | 28 | 337 |

- **601 rows disagree,** untriaged. Of the 230 omni takes, MySQL refuses most as a partition change on a table with none (1505, 81), an unknown column (1054, 68) and a key on a missing column (1072, 44). Where both refuse, omni answers without a code on 87 rows and 1051, 1146 or 1060 on most of the rest.

**OceanBase, the spike.** An earlier run with cruder sample tokens, kept for comparison. TiDB's generator produced 3,969 accepted statements and 15,560 near-misses. All were judged on a real server with `PREPARE`, in about 3.5 minutes.

| Statements | Total | Agree | Only engine accepts | Only omni accepts | Server crash |
| --- | --- | --- | --- | --- | --- |
| One per alternative | 3,969 | 2,622 | 1,287 | 57 | 3 |
| Near-misses | 15,560 | 13,860 | 1,296 | 404 | 0 |

- **3,047 disagreements,** while every existing oracle test passes. They are untriaged. `ALTER SYSTEM` and optimizer hints are the largest families, and some come from the spike's crude sample tokens.
- **The grammar alone** agreed with the engine on 90% of statements. Most of the rest is the lexer (531 malformed hints, which live in comments) and 124 `yyerror` checks in C actions. Token rows and overrides (§2.2, §2.3) exist for these.
- **Replay** took about 1 second per package, against 147 seconds (catalog) and about 300 seconds (review) live.

**TiDB, `tidb/parser`.** Moved 2026-11-13. `compat.tsv` holds 43,952 rows at `v8.5.5`, judged in 98 seconds, and 43,951 agree. The one left, `suite:cited`'s `ALTER DATABASE COLLATE 's' zzz`, is the judge's: TiDB stops at the unknown collation with 1273, which reads as accepted, and never reaches `zzz`, which omni rejects as TiDB does after a known collation.

| Rows | Total | Agree |
| --- | --- | --- |
| Generated: alternatives and near-misses | 11,055 | 11,055 |
| The four corpora: pingcap's parser tests, their rest, TiDB's integration tests, Bytebase's | 32,495 | 32,495 |
| `suite:reject-fuzz-sweep`: 300 mutations from a fixed seed | 292 | 292 |
| Statements two oracle tests listed, cited and fuzz findings | 110 | 109 |

- **The corpora** were 32,560 rows in their own files, sampled live at `TIDB_CORPUS_N`. 65 statements another row already held were left out. The pinned server gave every row the verdict its file recorded.
- **The reject fuzz** draws its seeds from the expected column, so closing a gap no longer moves them. Its minimiser sent requests that depend on omni's verdict, and is gone: a finding is a disagreeing row.
- **99 tests replay tapes,** 7.3 MB, recorded in 55 to 105 seconds and replayed in about 20. Each verdict is one `PREPARE omni_stmt FROM ?`, two with the probe. The codec writes `…` for the 121 characters that open a syntax error and for the TSO numbers in a read-timestamp error, which differ on every run. Six tests iterated maps while sending, and sort now. The operator sweep's size and seed are constants.
- **Deleted:** the live corpus tests and their sampler, and the harness tests of a run that recovers a lost connection or server (`TestTiDBCorpusHarnessSurvivesALostServer`, `TestPrepareSyntaxHoldsOneConnection`, `TestLostHarnessAnswersRatherThanPanics`): a replay has no server to lose, and the judge, which runs only in a recording, fails on a lost answer rather than recovering. `TestSyntaxProbeControls` stays, on a tape.

**TiDB, `tidb/review`.** Its five oracle tests replay (2026-11-13): 5 tapes, 1.6 MB, recorded in 80 to 90 seconds and replayed in 2. Two recordings left the same tapes.

- **The walk is a `Lister`.** One query lists databases, relations, columns, indexes and foreign keys. A local temporary table is in no `information_schema` view, its own session's included, so the inventory asks `SHOW CREATE TABLE` in every database for each name the change's `CREATE TEMPORARY TABLE` statements hold. The `Holder` and its digest are gone.
- **It finds 222 of the 247 statements TiDB rejects in the 60 curated cases, and 390 of 610 in a sweep of 40 changes,** with no false finding. CI swept 500 changes live on four servers; the tape holds 40, as OceanBase's does.
- **Three connections a case,** as OceanBase's: the setup, the catalog, and the change, which sets `foreign_key_checks` itself: a new session could read the global a sweep's base had turned off before the setup turned it on again.
- **Server state** goes through `testcontainer.Reset`: user databases dropped, `test` made again, and the global `foreign_key_checks` and `tidb_enable_check_constraint` on. The prior-backup, scan-rule and column-drop oracles stay as they were, on a tape each.

**Snowflake, converted.** `snowflake/parser/cmd/docsgrammar` read the Syntax sections of the 779 pages the command index lists on 2026-11-13 and wrote `snowflake/testdata/SnowflakeDocs.g4`: 765 pages have syntax, 760 give a rule, with 1,549 statement forms and 151 fragment rules. The blocks are near-BNF and map directly:

```text
ALTER WAREHOUSE [ IF EXISTS ] [ <name> ] { SUSPEND | RESUME [ IF SUSPENDED ] }
```

```antlr
// https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse
alterWarehouse
    : 'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID? ('SUSPEND' | 'RESUME' ('IF' 'SUSPENDED')?)
```

- **Placeholders are read by name** into twelve token classes, sampled in `tokens.tsv`: `<name>` is `ID`, `'<x>'` is `STRING`, `<num>` is `NUMBER`, `<expr>` is `EXPR`, `<query>` is `QUERY`, `<copy_statement>` is `COPY_STATEMENT`, a stage location is `STAGE`, a `file://` path is `FILE`.
- **Docs are prose.** 52 typos, mostly brackets that do not close, are corrected in `snowflake/testdata/docs-fixes.tsv`, each with its page. The converter's own readings include: `SET` or `UNSET` followed only by optional parts sets at least one; a definition with one `NAME = value` per line is any number of them; variants between `--` comments are a choice; a `|` left before `}` adds nothing; a `|` outside brackets chooses what follows the object's name (`ALTER TASK <name> RESUME | SUSPEND`). And 16 `SHOW` pages leave out the `|` after `<schema_name>` in their `IN` list, which the converter puts back.
- **38 forms are not converted,** listed in the `.g4` under their page: 28 write `...` for text they leave out (`SELECT * ILIKE ... REPLACE ...`), 6 name a fragment their page does not define (`{oauth_parameters}`) or are examples (`emp_id`), 2 hold Markdown, 1 has a bracket out of place, and CREATE INTEGRATION writes `<integration_type>` for the keyword each type's own page gives.

The docs judged the 58,310 rows in about a minute, and a replay asks the recognizer again in 8 seconds:

| Statements | Total | Agree | Only docs accept | Undocumented |
| --- | --- | --- | --- | --- |
| One per alternative | 15,653 | 9,880 | 5,773 | 0 |
| Near-misses | 42,657 | 539 | 278 | 41,840 |

- **6,051 disagreements,** all one way, since docs never reject. Omni rejects the shortest documented form of 310 of the 760 commands (106 `ALTER`, 81 `CREATE`, 64 `DROP`, 35 `SHOW`), mostly objects it does not know. By rows, the largest are `CREATE ICEBERG TABLE` and its catalog variants (330), `ALTER ICEBERG TABLE` (94), `ALTER DATABASE` (78), `ALTER POSTGRES INSTANCE` (67) and `ALTER SCHEMA` (66). The prototype's 64% is 61% here.
- **What the docs cannot decide:** an `EXPR` is a constant and a `QUERY` is `SELECT 1`, since the command pages hold no expression or query grammar; a `DATA_TYPE` is `INT`; a class takes any token of its kind, so a value the engine restricts is not checked; a form whose every part is optional (`ALTER STAGE a SET`, `CREATE ICEBERG TABLE a`) is accepted, as written; and clause order is the docs', where Snowflake often takes any.
- **Size.** The file is 9.3 MB. The undocumented near-misses are 6.9 MB of it and decide nothing; generating only the near-misses the docs derive would leave about 2.4 MB.

**PartiQL, recorded.** `partiql/parser/testdata/compat.tsv` (2026-11-13): 632 alternatives and 1,217 near-misses of `PartiQLParser.g4`, the grammar of the reference implementation (partiql-lang-kotlin, Apache-2.0, committed), and the 61 statements DynamoDB's developer guide prints. No server: a recognizer over the grammar judges, in 0.1 seconds.

| Statements | Total | Agree | Only the judge accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 632 | 251 | 380 | 1 |
| Near-misses | 1,217 | 1,200 | 14 | 3 |
| `suite:dynamodb-docs` | 61 | 59 | 2 | 0 |

- **The pin is `v0.14.9`, not the newest release.** The v1 grammar (`v1.6.1`) was rewritten along SQL's EBNF and no longer derives `INSERT INTO t VALUE {...}`, an `UPDATE`'s `REMOVE` or `RETURNING`, which DynamoDB documents. `v0.14.9` is the last release that holds them, and derives all 61 documented statements, so no row needs DynamoDB's docs to override the grammar.
- **400 disagreements, 271 of them one gap:** an expression as a statement (`dql : expr`), which omni's `Parse` refuses. It carries most expression alternatives, and the two DynamoDB rows omni rejects: `EXISTS(SELECT ...)`, a transaction's condition check. The other 125 the judge accepts: a graph pattern that begins with an edge (48), tokens after a statement (19), `EXCLUDE` (12), `VALUES` (11), `OVERLAY` (7).
- **What DynamoDB leaves out is not rejected.** The grammar accepts a `JOIN` or a `CREATE TABLE`; saying DynamoDB refuses one is a review rule's job.
- **The grammar's first one-token expression is `CURRENT_USER`,** which omni does not read, so every expression row carried it and 583 alternatives disagreed. `tokens.tsv` samples `exprTerm` as a name (§2.2).
- **31 statements the grammar derives are refused by the parser built from it,** in the visitor that reads the parse tree: an unknown `EXPLAIN` parameter, `LAG` with no `ORDER BY`, a `DATE` string that is no date. Each is a cited `grammar` row of `grammar-overrides.tsv`. They were read off the source; running the reference parser on every row needs a JVM, and would make it the judge.

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
