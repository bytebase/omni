# GoogleSQL dialects: BigQuery and Spanner

Date: 2026-11-14. Status: approved 2026-11-14, with every decision of §8 as
recommended.

## 1. Goal

`googlesql/parser.Parse` accepts what BigQuery accepts or what Spanner accepts.
AGENTS.md says `Parse` rejects what the engine rejects, and the union is neither
engine. Bytebase runs both engines on this one parser
(`backend/plugin/parser/{bigquery,spanner}`, through
`backend/plugin/parser/googlesql`), so a Spanner user is never told that
`DECLARE x INT64` or `CREATE TABLE t (a INT64)` is invalid, though Spanner
refuses both.

This design gives the parser a dialect. A strict BigQuery parse accepts what
BigQuery's parser accepts, and a strict Spanner parse what Spanner's parser
accepts, each judged by the engine itself.

## 2. How the two differ

Measured on the pinned Spanner emulator 1.5.54 (the 336 statements of
`docs/migration/googlesql/truth1/bigquery` and about 100 probes) and on
BigQuery by dry runs in `bytebase-dev` (50 probes, §5).

Both engines parse queries and DML with GoogleSQL's parser. A product turns
GoogleSQL features on in `LanguageOptions`; the emulator's list is in
`backend/query/analyzer_options.cc`. `googlesql.tm` checks about two dozen of
those features while it parses (`FEATURE_PIPES`, `FEATURE_QUALIFY`,
`FEATURE_ALLOW_DASHES_IN_TABLE_NAME`, ...). Most other differences are refused
after the parse, by the analyzer ("QUALIFY is not supported", "Statement not
supported: MergeStatement"). GoogleSQL's grammar also holds Spanner's older DDL
(`INTERLEAVE IN PARENT`, `NULL_FILTERED`, a `PRIMARY KEY` after the columns)
behind `FEATURE_SPANNER_LEGACY_DDL`, which BigQuery leaves off: BigQuery
answers "PRIMARY KEY must be defined in the table element list" there.

| Area | Syntax | Measured |
| --- | --- | --- |
| Queries and expressions | Nearly the same | Spanner's parser takes all 52 of BigQuery's documented expressions, and `QUALIFY`, `IS DISTINCT FROM`, `WITH RECURSIVE`, `BY NAME`, `CORRESPONDING`, `FOR SYSTEM_TIME AS OF`, `PIVOT`, `UNPIVOT`, `GROUP BY ALL`, `* EXCEPT`. Its parser refuses an unquoted dashed table name (`FROM my-project.ds.t`), `MATCH_RECOGNIZE`, and a pipe query that starts with `FROM`; BigQuery parses all three. Both parse a table hint (`t@{FORCE_INDEX=i}`) and graph queries. A statement hint (`@{USE_ADDITIONAL_PARALLELISM=TRUE} SELECT 1`) parses on Spanner, and BigQuery refuses it. |
| DML | Nearly the same | 11 of BigQuery's 12 documented DML statements parse on Spanner (`TRUNCATE TABLE` does not), and `MERGE` parses. BigQuery's parser takes Spanner's `INSERT OR UPDATE` and `THEN RETURN` too, and answers that the table does not exist. |
| Scripting | BigQuery only | `DECLARE`, `IF`, `CASE`, `LOOP`, `WHILE`, `REPEAT`, `FOR ... IN`, labels, `RETURN`, `RAISE`: Spanner parses one statement, never a script. 25 of BigQuery's 44 scripting examples fail. |
| DDL | Different grammars | The emulator reads DDL with its own JavaCC parser (`backend/schema/parser/ddl_parser.jjt`), not GoogleSQL's. 64 of BigQuery's 79 DDL examples fail it: `CREATE TABLE t (a INT64)` with no key, `PARTITION BY`, `CLUSTER BY`, table `OPTIONS`, `CREATE OR REPLACE TABLE`, `AS SELECT`, `LIKE`, `COPY`, `CLONE`, temp functions, procedures, models, row access policies, capacity, `RENAME COLUMN`, `SET DATA TYPE`. Spanner's own DDL is not BigQuery's: `PRIMARY KEY (a)` after the columns, `STRING(MAX)`, `INTERLEAVE IN PARENT`, `STORING`, `NULL_FILTERED`, `ROW DELETION POLICY`, `CHANGE STREAM`, roles, `GRANT ... TO ROLE`. On BigQuery the trailing `PRIMARY KEY` and `GRANT ... TO ROLE` are syntax errors, and `CREATE ROLE` and `CREATE CHANGE STREAM` are refused by GoogleSQL's grammar ("ROLE is not a supported object type"). |
| DCL | Different | BigQuery grants IAM roles to principals (``GRANT `roles/...` ON TABLE t TO 'user:...'``); Spanner grants privileges to database roles. |

omni's parser already marks most of this: the `bq_*.go` files, the
`spanner_*.go` files, and `create_table.go`, which names each form it takes
from one dialect.

## 3. Shape: one parser with a dialect

One parser, with a dialect chosen per call. Not a fork engine, as `mariadb` is
of `mysql`:

- GoogleSQL itself is one grammar gated by `LanguageOptions`, and the
  measurements show the query and DML grammar is shared. A fork would copy
  `expr.go`, `select.go`, `from_join.go` and the rest, most of the parser's
  36,000 lines, and every query fix would land twice.
- The split that is real, DDL and scripting, sits at statement dispatch
  (`parseStmt`, `parseCreateStmt`, `parseAlterStmt`, `parseDropStmt`), where a
  dialect check is one line per form.

Within a statement, a gate reads the dialect where `googlesql.tm` reads a
feature, and a dialect-only statement is refused at dispatch:

```go
// Dialect names the engine a parse holds the text to.
type Dialect int

const (
	// Union accepts what either engine's parser accepts. It is Parse's
	// dialect while callers move (§6).
	Union Dialect = iota
	BigQuery
	Spanner
)

// ParseDialect is Parse for one engine.
func ParseDialect(sql string, d Dialect) (*ast.File, error)
```

`ParseBestEffort`, `diagnostics.Analyze` and `parser.Split` take the dialect
too. `Split` needs it for nothing measured so far: a Spanner transaction
(`BEGIN; ...; COMMIT;`) splits the same in both, and `BEGIN ... END` blocks
exist only in BigQuery. `analysis.Dialect` becomes an alias of
`parser.Dialect`, since `analysis` imports `parser` and not the reverse.

## 4. What a strict parse refuses

A dialect parse refuses what that engine's **parser** refuses, and nothing the
engine refuses later. That is the line every other engine's judge draws: Trino
rejects on `SYNTAX_ERROR` only, and SQL Server under `SET PARSEONLY ON`. So
Spanner's dialect accepts `QUALIFY` and `MERGE`, which Spanner's parser takes
and its analyzer refuses. That refusal is semantic, and belongs to a review
rule if Bytebase wants one.

Spanner's dialect refuses:
- the scripting statements, and every statement dispatch has for BigQuery's
  DDL, DCL, `EXPORT` and `LOAD` only;
- the GoogleSQL features the emulator leaves off and `googlesql.tm` checks
  while it parses: dashed names, `MATCH_RECOGNIZE`, a pipe query that starts
  with `FROM`;
- BigQuery's forms inside shared DDL: a table with no primary key, `PARTITION
  BY`, table `OPTIONS`, `OR REPLACE`, `AS SELECT`.

BigQuery's dialect refuses Spanner's DDL (the trailing `PRIMARY KEY`,
`INTERLEAVE`, `NULL_FILTERED`, `CREATE ROLE`, `CREATE CHANGE STREAM`) and
`GRANT ... TO ROLE`. It keeps `INSERT OR UPDATE`, `THEN RETURN` and table
hints, which BigQuery's parser takes (§2).

**The parser's refusals are more than "Syntax error:".** GoogleSQL's grammar
refuses some input in its actions, with a message of its own: "ROLE is not a
supported object type", "PRIMARY KEY must be defined in the table element
list", "RUN statement is not supported." `googlesql.tm` has 216 such calls
(`MakeSyntaxError`, `MakeSqlErrorAt`, `SetOverrideError`). Both judges read
those messages as rejections too: a list generated from `googlesql.tm`, each a
pattern of the call's text, beside the converted grammar. Today the Spanner
harness reads only the prefix, so its `compat.tsv` counts such refusals as
acceptances; step 2 of §7 rejudges it.

`Union` stays exactly today's parser.

## 5. Judges

Each dialect is held to its own `compat.tsv`, with its own judge and rows:

| Dialect | File | Rows | Judge |
| --- | --- | --- | --- |
| Spanner | `googlesql/parser/testdata/compat.tsv` (today's) | `GoogleSQL.g4` from `googlesql.tm`, as now; `SpannerDDL.g4`, converted from the emulator's `ddl_parser.jjt` (Apache-2.0, 1,043 lines) by a new `googlesql/parser/cmd/jjt2g4`; the `suite:` rows | The emulator, as now |
| BigQuery | `googlesql/parser/testdata/compat-bigquery.tsv` | `GoogleSQL.g4`; `suite:` rows from `truth1/bigquery` | Dry runs on BigQuery |

**Spanner's DDL rows.** Today's rows come from `googlesql.tm`, whose DDL is
GoogleSQL's and not what the emulator parses, so Spanner's DDL is covered only
by the 172 `suite:ddl` and `suite:spanner-ddl` rows. The converted
`ddl_parser.jjt` generates rows for every Spanner DDL production, the way Doris
reads its legacy CUP parser through `compat.Suite.More`.

**BigQuery's judge: dry runs.** No official BigQuery emulator exists; the
community ones are a different implementation, which the compatibility spec's
rule against community grammars rules out. A query job with `dryRun: true`
parses and validates the statement on BigQuery, runs nothing and costs nothing.
The judge reads it as Spanner's harness reads the emulator: `Syntax error:`,
or a message of GoogleSQL's grammar (§4), is a rejection, at the line and
column it reports; every other answer is an acceptance, as `Not found:
Dataset`, `Access Denied` on a table, `Table "t" must be qualified with a
dataset`, or "CREATE INDEX is not supported", which BigQuery's front end
writes; a transport or quota error is no verdict.

- **Where:** the dedicated project `bytebase-dev`. A recording reads the
  project from `BIGQUERY_COMPAT_PROJECT` and a token from
  `gcloud auth print-access-token`, so it needs no Google client library: the
  judge is a REST call in `googlesql/internal/testcontainer`, with no nested
  module and no new dependency. CI never records.
- **Account:** the devbox's service account,
  `devbox@bytebase-dev.iam.gserviceaccount.com`, with `roles/bigquery.jobUser`
  on `bytebase-dev` (granted 2026-11-14). A dry run needs no dataset.
- **Pin:** BigQuery has no release to pin. Its `pins.tsv` row names the service
  and the recording date, as a docs engine pins its conversion date; a
  recording that moves a cell bumps the date, so `compatcheck` treats it as a
  pin bump and the diff is BigQuery's drift.
- **Measured:** a dry run parses a whole script (`DECLARE x INT64; SET x = 1;
  SELECT x` answers `SCRIPT`, and `SET x = ;` a syntax error at its column).
  DDL into a dataset that does not exist is parsed first (`CREATE TABLE
  nope_ds.t (a INT64) PRIMARY KEY (a)` answers the grammar's error, `CREATE
  TABLE nope_ds.t (a INT64)` "Not found: Dataset"). A trailing `;` is taken.
  One client asks 3.6 statements a second, so about 9,000 rows take 40
  minutes alone, or five with eight at a time.

Both judges report a line and column, so either file can later hold
`reject@<byte>`, as Trino's does. Not in this design.

## 6. Bytebase

Bytebase already carries the dialect: `googlesql.Config.Dialect` in
`backend/plugin/parser/googlesql/config.go`, set by the `bigquery` and
`spanner` plugins and passed only to `analysis.GetQuerySpan`. The parser calls
(`statement_type.go`, `query_type.go`, `diagnose.go`, `split.go`) take it once
`ParseDialect` exists. Each engine then reports its own syntax errors. That is
a visible change, and its tests decide it: a Spanner statement Bytebase took
before and now refuses is either a gap in omni or a change Bytebase wants.

`Parse` keeps `Union` until Bytebase calls `ParseDialect` everywhere, and then
D2 decides whether it stays.

## 7. Rollout

One PR each; each leaves `go test ./googlesql/...` passing with no server.

1. **`parser.Dialect` and `ParseDialect`,** with `Union` the only behavior:
   the API, `analysis.Dialect` as an alias, the dialect threaded through
   `Parser`. No row moves.
2. **Spanner's query-level gates and scripting.** The grammar's messages
   (§4) join the Spanner judge and `compat.tsv` is rejudged; `omniVerdict`
   parses as `Spanner`, and the rows it closes are committed with
   `-update-compat`.
3. **Spanner's DDL rows.** `jjt2g4`, `SpannerDDL.g4`, the rows judged by the
   emulator; no parser change.
4. **Spanner's DDL gates,** closing what step 3 found.
5. **BigQuery's judge:** the dry-run client, the pin,
   `compat-bigquery.tsv`. `compat.Suite` gains `Role`, the pin role a suite
   judges with, since one `pins.tsv` will hold two judges.
6. **BigQuery's gates,** closing what step 5 found.
7. **Bytebase** passes its dialect, and bumps omni.

Steps 2 to 4 and 5 to 6 are independent and can run side by side. Docs counts
go into one final PR, as for the compatibility spec.

## 8. Decisions

All five were decided as written here on 2026-11-14.

- **D1. One parser with a dialect, not a fork** (§3).
- **D2. What `Parse` means afterwards.** `Union` keeps every caller working
  and breaks AGENTS.md's rule; removing `Parse` in favor of `ParseDialect`
  makes every caller choose. Decided: keep `Union` until step 7, then
  decide with Bytebase's tests in hand.
- **D3. Parser verdicts, not analyzer verdicts** (§4), as every engine's judge
  does. The alternative, refusing what Spanner's analyzer
  refuses, makes the parser track Spanner's feature list and needs a judge
  that reads "not supported" as a rejection.
- **D4. BigQuery's dry runs on `bytebase-dev`.** Done: the devbox's service
  account has `roles/bigquery.jobUser` there.
- **D5. BigQuery's rows.** `GoogleSQL.g4` is GoogleSQL's open-source grammar,
  and BigQuery runs Google's internal build of it, which may be newer. The
  dry run judges every row, so a row the open grammar lacks only goes
  untested, and BigQuery's docs corpus adds `suite:` rows for what its docs
  show.
