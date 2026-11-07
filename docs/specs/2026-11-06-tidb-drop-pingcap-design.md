# Dropping pingcap/tidb from Bytebase

Date: 2026-11-06. Status: draft for review. §8 lists the decisions it needs.

## 1. Goal

No Bytebase package imports `github.com/pingcap/tidb/...`. Every TiDB path
parses with `omni/tidb` and nothing else. The root module
`github.com/pingcap/tidb` leaves Bytebase's `go.mod`, and so do both `replace`
directives to `github.com/bytebase/tidb`.

`github.com/pingcap/tidb/pkg/parser` itself stays as an indirect requirement,
because a dependency of a dependency imports it:
`github.com/go-mysql-org/go-mysql/client` imports `pkg/parser/charset`, and
Bytebase links that through `gh-ost`. Removing it means changing go-mysql or
Bytebase's gh-ost fork, which is outside this design (§8, D5).

"Full compatibility" with pingcap's parser means three things here, and only
these:

1. **Acceptance parity with the TiDB server.** `tidb/parser.Parse` accepts what
   the pinned TiDB release accepts and rejects what it rejects. The authority is
   the server, not pingcap's parser:
   - pingcap's parser on master accepts grammar no release has (109 corpus rows
     today);
   - the server refuses some statements during parsing in ways pingcap's
     grammar file does not show.
2. **Capability parity for Bytebase's consumers.** Everything Bytebase does
   with a pingcap AST today (§4) has an omni equivalent with the same
   observable result, or a deliberate, listed change.
3. **Text parity where text is compared.** Where Bytebase compares SQL text
   with what TiDB itself reports, omni writes what TiDB writes: column
   defaults, generated expressions, view bodies, type names. omni already
   writes `SHOW CREATE TABLE` and view definitions in TiDB's format (#701,
   #704).

### Non-goals

- **A pingcap-shaped AST, or a converter from omni's AST to pingcap's.**
  Consumers move to omni's AST.
- **pingcap's error text.** omni's `ParseError` carries a byte offset, and
  Bytebase converts it to a position as it already does for MySQL.
- **Optimizer hint structure.** TiDB reads `/*+ ... */` with a separate hint
  parser and turns a bad hint into a warning, never a syntax error, so hints do
  not affect acceptance. No Bytebase consumer reads hints.
- **TiDB releases other than the pinned one** (§3.6).

## 2. Approach, and what was rejected

Keep omni's hand-written parser and measure it against the server, as the
other engines do. Move Bytebase's consumers onto omni one at a time, following
what Bytebase's MySQL adapter already does on `omni/mysql`. Delete the pingcap
fallback last.

Rejected:

- **Vendor pingcap's goyacc parser into omni.** It brings pingcap's AST, its
  nodes without `Loc`, its `types` and `charset` packages, and its error model.
  It breaks the engine contract in `AGENTS.md` and moves the dependency rather
  than removing it.
- **An omni-to-pingcap AST converter in Bytebase.** It keeps the dependency.
- **Port `parser.y` production by production into omni's code.** omni's parser
  is hand-written to its own contract. `parser.y` tells it what to cover
  (§2.1), not how to write it.
- **Start from the published docs.** `tidb/parser` began from the MySQL
  reference manual's syntax descriptions (the 140 files in `tidb/parser/bnf/`).
  That is how it stood at 65.5% against TiDB on 2026-11-03. The docs leave out
  words TiDB reads and ignores, and some TiDB-only options, and they describe a
  "stable" version that moves.

### 2.1 Authority

Two sources, highest first. The lower never overrules the higher.

1. **The TiDB v8.5.5 server** decides what is accepted and what a statement
   means. It is asked in three ways:
   - `PREPARE` for whether a statement parses (§3.1);
   - comparing results for the meaning of an expression (§3.5);
   - running DDL and reading `SHOW CREATE TABLE` and `information_schema` for
     the meaning of a definition.
2. **`pkg/parser/parser.y` at the `v8.5.5` tag** decides what to build. It
   lists every statement form, the order and nesting of clauses, and operator
   precedence. Read it at the release tag, never master: master has grammar no
   release has. Three things it does not settle, and the server answers
   instead:
   - **The code attached to its grammar rules also refuses input.** About 146
     of these actions call `AppendError` or return 1. What each one refuses is
     read from the server.
   - **Lexing is elsewhere.** Keywords, `/*T! */` comments, hex and bit
     literals, names that start with a digit, and `_charset` introducers live
     in `lexer.go` and `misc.go`. Hints have their own grammar in
     `hintparser.y`.
   - **Grammar does not fix meaning.** A tree can follow the productions and
     still group or attach wrongly.

omni's TiDB parser already cites `parser.y` at v8.5.5 403 times, e.g.
`// Ref: AdminStmt in TiDB v8.5.5 parser/parser.y:11204-11501`. This section
makes that the rule for the rest of the work.

## 3. Workstream A: grammar (omni)

### 3.1 Where it stands

`tidb/parser/testdata/pingcap_parser_corpus.tsv` holds every one-line case in
pingcap's `parser_test.go`. Its engine column comes from TiDB v8.5.5 through
`PREPARE`; `docs/plans/2026-11-03-tidb-grammar-gap-inventory.md` describes the
method.

| | 2026-11-03 | 2026-11-06 |
|---|---:|---:|
| both accept | 2090 | 3060 |
| gap: TiDB accepts, omni rejects | 1099 | 56 |
| over-acceptance: TiDB rejects, omni accepts | 122 | 24 |
| both reject | 591 | 762 |
| omni coverage of what TiDB accepts | 65.5% | 98.2% |

The engine column changed between the two dates (3,189 accepted, now 3,116).
That happened when the harness learned to tell a refusal during parsing from
an error that comes after it (#670, #698).

### 3.2 Remaining corpus rows

Gaps (56):

- ODBC escapes `{ts '...'}`, `{d ...}`, `{t ...}`, `{ident expr}` — 9
- `_charset` introducer before a hex or bit literal (`_utf8 0xD0B1`) — 4
- `FETCH FIRST|NEXT [n] ROW[S] ONLY` — 6
- Window null treatment and `NTH_VALUE ... FROM FIRST|LAST` — 5
- Column type words TiDB keeps from MySQL 5 (`ASCII`, `BYTE`, `LONG`,
  `STORAGE DISK|DEFAULT`, repeated `ZEROFILL` on `YEAR`) — 12
- Modifiers TiDB reads and ignores (`SELECT LOW_PRIORITY|DELAYED|SQL_CACHE`,
  `UPDATE|DELETE HIGH_PRIORITY|DELAYED`) — 7
- `LOCK TABLES ... WRITE LOCAL` — 3
- `SET @ = 1`, `SET @' ' = 1`, `SET SESSION_STATES` — 3
- One family each: `DROP TABLES` (2 rows), `UNION DISTINCTROW`,
  `USE INDEX (PRIMARY)`, `DELETE ... USE INDEX`, nested `JOIN ... ON ... ON`,
  `SCHEMA()` — 7

Over-acceptances (24):

- `RIGHT JOIN ... LEFT JOIN t3` with no join condition — 2
- Multi-table `UPDATE`/`DELETE` with `ORDER BY` or `LIMIT`, and
  `DELETE FROM t1.*` — 5
- `->`/`->>` on a non-column left side or with a non-string path — 4
- Parenthesised column `DEFAULT (expr)` that TiDB refuses — 5
- `0X11` and `x'0xaa'` literals — 2
- `SECONDARY_ENGINE_ATTRIBUTE` on an index — 3
- `VALUES LESS THAN (DEFAULT)`, `NTH_VALUE(val)`,
  `MOD(10, 2 BETWEEN 0 AND 5)` — 3

Each family is a small PR of the kind landing now. Each one:

- re-baselines with `TestTiDBCorpusOmniVerdicts -update`;
- checks whether `mysql`, `mariadb`, and `oceanbase` share the gap, and fixes
  them in the same PR if they do.

### 3.3 A corpus from parser.y: what "complete" means

The pingcap corpus is a sample: the statements pingcap happened to test. It
can show a gap but not that none is left. `parser.y` at v8.5.5 lists every
form, so a corpus built from it can.

**How it is built.** A generator reads the grammar rules of `parser.y` at the
tag: about 700 rules with about 3,000 alternatives. For each alternative it
writes the shortest statement that uses it:

- **Shortest derivation.** Expand every other rule along its shortest path to
  text, so recursive rules end.
- **Fixed sample text per token.** A keyword becomes its word, a name `t`, a
  number `1`, a string `'s'`, a hex literal `0x1`.
- **Hand-written sentences where generation fails.** Some alternatives only
  make sense in context, or are reached only through precedence. Each gets a
  hand-written sentence, kept in the generator beside the rule it covers. A
  few dozen are expected.

The generator is a root-module command, `tidb/parser/cmd/grammarcorpus`, like
each engine's `genwalker`: `./...` compiles it but never runs it. It reads a
vendored copy of `parser.y` (`testdata/parser.y`, verbatim at the tag,
Apache-2.0 header intact) and writes the corpus, so regeneration is
deterministic and offline. The vendored grammar and the TSV are the only
checked-in inputs; nothing imports or compiles the grammar, tests read the TSV
and never `parser.y`, and omni's own code imports nothing from pingcap.

**What it records.** `tidb/parser/testdata/parser_y_corpus.tsv` has four
columns: the omni and engine verdicts, the rule and the alternative's index
(e.g. `AlterTableSpec#41`), and the statement, which stays last so the loader
keeps reading it whole. The same `PREPARE` harness fills the engine column.

**What "complete" means.** Every row agrees: omni accepts each alternative the
server accepts, and refuses each one the server refuses (the grammar actions
of §2.1). A grammar PR names the rows it flips, and a PR that leaves a rule
without a row is incomplete.

**When the pin moves**, the new tag's `parser.y` regenerates the file. The
diff is the list of forms the release added or removed.

**Status (2026-11-07).** Landed: `tidb/parser/cmd/grammarcorpus` and the first
measurement. parser.y at v8.5.5 has 690 rules and 2,946 alternatives; omni
parses 2,654 of the 2,807 the server accepts (94.5%) with zero
over-acceptances. The 153 gaps, grouped by rule, are
`docs/plans/2026-11-07-tidb-parser-y-gap-inventory.md`, and are the grammar
work list that replaces §3.2.

### 3.4 Samples of real SQL

The `parser.y` corpus proves coverage of each form alone, not of forms
combined the way people write them. The one-line pingcap corpus stays, and
three more samples join it. Each becomes its own TSV in the same three-column
format, measured by the same `PREPARE` harness:

1. **The rest of pingcap's parser tests.** These are the raw-string and
   multi-line cases in `parser_test.go` that the first extraction skipped, plus
   the restore cases in `ast/*_test.go`: several hundred rows between them.
2. **TiDB's integration tests.** These are `tests/integrationtest/t/**/*.test`
   in the TiDB repository: 240 files and about 50,000 statements before
   deduplication. People wrote them to exercise the server, not the grammar,
   which makes them the closest public stand-in for user SQL.
   - Split them with `tidb/parser.Split`, deduplicate, and drop the `--`
     directives.
   - A statement that names a missing table still gets a verdict, because
     `PREPARE` parses before it resolves names.
3. **Bytebase's own TiDB test files.** These are the 64 files under
   `backend/plugin/parser/tidb/test-data`, `backend/plugin/advisor/tidb/test`,
   and `backend/plugin/schema/tidb/testdata`. They are what Bytebase already
   promises to handle.

Each source lands with its first measurement, and its gaps join §3.2's list.
At these sizes a live run of every row on every test run is expensive:

- the live engine test may run a sample, controlled by `TIDB_CORPUS_N` like the
  other sweeps;
- the omni-verdict test always reads every row.

### 3.5 Tree shape

Bytebase's consumers read trees, not verdicts. A tree that groups or attaches
wrongly still parses, and then gives wrong metadata, wrong lineage, or a wrong
stored view. Acceptance corpora cannot see that, so each grammar PR also checks
the tree against the server:

- **Expressions:** add the new form to `operatorTreeCases`.
  `TestOperatorTreesMatchEngine` evaluates the expression and omni's fully
  parenthesised reading of it on TiDB over every combination of 0, 1, 2 and
  NULL, and requires `NOT (expr <=> shape)` to match no row.
- **Definitions:** run the DDL on TiDB, and compare `SHOW CREATE TABLE` and
  `information_schema` with `tidb/catalog`, as the catalog's container tests
  already do.
- **Stored text:** compare view bodies and defaults with what TiDB stores.

### 3.6 Target version

Stay on v8.5.5. It is the latest LTS, and every TiDB oracle test pins it.

TiDB's grammar has only grown across releases, so accepting what v8.5 accepts
also serves users on v7.5 and v8.1. The exception is a statement only a newer
release accepts: omni passes it, and the older server rejects it when the
statement runs.

When the next LTS ships, move the pin and rerun the engine verdicts. The diff is
the list of grammar that release added.

## 4. Workstream B: Bytebase consumers

Only the parse dispatcher is omni-first today. Every other consumer still works
on a pingcap AST. It gets one in one of two ways:

- **Directly**, by calling `ParseTiDB` or `parser.New()`.
- **Through the bridge.** `GetTiDBAST` calls `OmniAST.AsPingCapAST`, which
  parses the text again with pingcap (`parser/tidb/omni.go`).

Bytebase's MySQL adapter runs every one of these capabilities on `omni/mysql`,
so each row below has a working precedent. Paths are under
`backend/plugin/`.

| Consumer | pingcap use today | Move to | Change to expect |
|---|---|---|---|
| `parser/tidb/diagnose.go` | `ParseTiDB` only | the dispatcher's omni parse and `convertOmniParseError` | editor diagnostics show omni's error text |
| `parser/tidb/statement_type.go` | about 15 statement types, through the bridge | a switch on `omni/tidb/ast`, as MySQL's does | none |
| `parser/tidb/query.go` (`validateQuery`, `hasSelectInto`, `ExtractMySQLTableList`) | `ParseTiDB`, `ast.Visitor` | `omni/tidb/ast` and `ast.Inspect`, as MySQL's does | pingcap puts `INTO` on a set-operation arm; check omni's shape |
| `parser/tidb/resource_change.go` | pingcap first, omni only for DML and `BATCH` | omni only, with DDL coverage ported from MySQL's | DDL changes are no longer lost when pingcap fails |
| `db/tidb/stats.go` (`countDMLTargets`) | `ParseTiDB`, `DeleteStmt` and `UpdateStmt` table refs | `omni/tidb/ast` | none |
| `advisor/tidb` prior-backup check, `getTiDBNodes` | `DDLNode` assertion, through the bridge | an omni DDL classifier, as MySQL's `isOmniDDL` | none |
| `parser/tidb/limit.go` | `ParseTiDB`, `ValueExpr`, `format.Restore` | add or replace `LIMIT n` in place by `Loc`, as MySQL's `limit.go` does (§8, D2) | the query keeps the user's text; `limit_test.go` goldens change |
| `parser/tidb/query_span_extractor.go`, `query_type.go`, `util.go` | about 40 expression node types; names unaliased columns from `Text()`; DML/DDL from pingcap's interfaces | port Bytebase's MySQL omni span extractor to `omni/tidb/ast`; name unaliased columns from the `Loc` slice | the largest consumer port; the recursive-CTE split must handle omni's set-operation shape |
| `schema/tidb/get_database_metadata.go` | `types.FieldType` (`CompactStr`, flen, charset), `Text()` for defaults | `tidb/catalog.Exec`, then a TiDB `catalogToProto`, as MySQL's does | type strings follow the catalog; `UNSIGNED` and `ZEROFILL`, dropped today, come back (§8, D3) |
| `schema/tidb/walk_through.go` | the bridge; `types.Datum` default checks; `format.Restore` for defaults | `tidb/catalog` (`LoadMetadata`, `Exec`), as MySQL's walk does (§8, D1) | goldens with pingcap error text such as `[types:1292]` change |
| `schema/tidb/generate_migration_test.go` | `parser.New()` checks that migrations parse | `omni/tidb/parser.Parse` | none |

Dead code goes in the first PR that touches each file:

- `TypeString` and `ParseTiDBForSyntaxCheck` in `parser/tidb/tidb.go`;
- `SplitTiDBMultiSQL` in `parser/tokenizer`.

omni needs nothing new for this beyond the grammar. Every capability above
exists in `tidb/catalog`, `tidb/deparse`, `tidb/parser.Tokenize`, or
`tidb/review` today. If a port finds a missing piece, omni gets it with a test
against the engine before Bytebase uses it.

## 5. Workstream C: delete the fallback

### 5.1 Gate

All three must hold:

1. **No consumer needs pingcap.** Every row in §4 is done, so nothing calls
   `GetTiDBAST` or `ParseTiDB`.
2. **The corpora are clean.** The `parser.y` corpus (§3.3) and every sample
   (§3.4) have no gap and no over-acceptance, or each remaining row is waived
   in its TSV with a reason. The `parser.y` corpus is the proof of
   completeness; the samples catch forms that only fail in combination.
3. **Production agrees, where it can see.** `tidb_dispatcher_omni_fallback_total`
   shows no `unknown` reason and stays under 0.1% of statements for 30 days.
   That is the threshold Bytebase's `metrics.go` already states. The counter
   sees only deployments that report metrics, which excludes most self-hosted
   installs, so it supports gate 2 and does not replace it.

### 5.2 Deletion

One Bytebase PR:

1. Make the dispatcher omni-only. When omni rejects a statement, that is the
   syntax error.
2. Delete pingcap's side of the bridge:
   - the `AST` type, `PingCapASTProvider`, `GetTiDBAST`, and
     `OmniAST.AsPingCapAST`;
   - `ParseTiDB`, `newTiDBParser`, `convertParserError`, and
     `applyTiDBLineTracking`;
   - `metrics.go`;
   - `tokenizer.SetLineForMySQLCreateTableStmt`;
   - the tests that pin the bridge (`omni_test.go`, `tidb_test.go`).
3. Drop both `replace` lines and run `go mod tidy`:
   - the root `github.com/pingcap/tidb` leaves `go.mod`;
   - `pkg/parser` stays as an indirect requirement through gh-ost and
     go-mysql, now resolved to the upstream version.

## 6. Risks

- **Users see errors they did not see before.** Two kinds of statement are
  involved:
  - **Correct new errors.** Bytebase's pingcap fork accepts grammar v8.5.5
    rejects, and those statements become syntax errors. The server would have
    rejected them anyway.
  - **Wrong new errors.** Statements only a TiDB release newer than the pin
    accepts become syntax errors too. Moving the pin promptly (§3.6) keeps
    these few.
- **Visible output changes:**
  - error messages in diagnostics and walk-through advice;
  - type strings in parsed schema metadata, which feed schema comparison;
  - the text of queries rewritten for a result limit.

  Each is listed in §4. Each should land in its own PR with the golden diff
  reviewed, not folded into a larger change.
- **The query span port is large.** It is the one consumer with no small slice.
  To limit risk, run the MySQL extractor's test fixtures and TiDB's
  `query_span.yaml` against both implementations, and switch only when they
  agree.
- **The production counter is weak evidence** (§5.1). The corpora carry the
  gate.

## 7. Order of work

Bytebase pins omni by commit, so each Bytebase step below starts by bumping
that pin.

1. **omni: close §3.2's 80 rows.** Small PRs that can run in parallel, each
   with §3.5's tree checks.
2. **omni: build the `parser.y` corpus (§3.3).** One PR with the generator and
   its first measurement. Its gaps, grouped by rule, become the grammar work
   list, replacing §3.2 once it is closed.
3. **omni: add §3.4's three samples.** One PR each, each listing its gaps;
   then close those gaps.
4. **Bytebase: the small consumers.** `diagnose.go`, `statement_type.go`,
   `query.go`, `resource_change.go`, `db/tidb/stats.go`, the prior-backup
   check, `generate_migration_test.go`, and dead code. These are independent of
   each other.
5. **Bytebase: `limit.go`**, after D2.
6. **Bytebase: the query span extractor.**
7. **Bytebase: `get_database_metadata.go`, then `walk_through.go`.** The second
   reuses the first's `catalogToProto`.
8. **Check §5.1's gate, then make §5.2's deletion.**

Steps 1–3 (omni) and steps 4–7 (Bytebase) can run at the same time. Steps 4–7
need only grammar the dispatcher already parses with omni. A consumer that
meets a statement omni cannot parse returns no result, as the bridge does
today.

## 8. Decisions needed

- **D1. Walk-through.** Port MySQL's catalog walk to `tidb/catalog`
  (recommended), or send Bytebase's V1 walk to omni's `tidb/review` walk, which
  the V2 review path already runs?
- **D2. Limit rewriting.** Add or replace `LIMIT n` in the user's text by
  `Loc`, as MySQL does (recommended, goldens change)? Or reprint with
  `tidb/deparse.RestoreSelectWith`, which keeps today's goldens but makes
  `Restore` match pingcap's character for character?
- **D3. Metadata type strings.** Follow the catalog's TiDB `SHOW CREATE TABLE`
  form (recommended; closer to what sync reports), or reproduce pingcap's
  `CompactStr`, including dropping `UNSIGNED` and `ZEROFILL`?
- **D4. Gate.** Do the corpora alone gate the deletion, or also the production
  counter's 30 days?
- **D5. The `pkg/parser` that remains through gh-ost and go-mysql.** Accept it
  as indirect (recommended for this project), or patch Bytebase's gh-ost fork
  so it no longer pulls in go-mysql's `client`?
