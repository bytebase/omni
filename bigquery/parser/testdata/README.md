# BigQuery parser test data

## `compat.tsv`: live BigQuery dry-run verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../../googlesql/testdata/GoogleSQL.g4` (converted from `googlesql.tm` at
GoogleSQL `2026.01.1` by `spanner/parser/cmd/tm2g4`), followed by its
near-misses, with omni's verdict (`bigquery/parser.Parse`, BigQuery's own
parser, forked from `spanner/parser`'s `DialectBigQuery` mode) and live BigQuery
dry-run's verdict on the pinned Google Cloud project (`../../testdata/pins.tsv`,
`default` role).

`TestCompatibility` (`../compat_test.go`) reads the committed file through
`enginetest/compat` and requires **no Google Cloud credentials and no network
access** in offline replay. Recorded 2026-11-21 against `bytebase-dev` (`US`)
via BigQuery dry-run jobs (`DryRun = true`, `UseLegacySQL = false`), which
execute nothing, cost $0, and create or mutate no datasets or tables, then
verified end-to-end with `OMNI_RECORD=verify` across all 8,837 rows with zero
differences (793 seconds).

### Current counts

| Rows | Total | Agree | Only BigQuery accepts | Only omni accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 2,737 | 2,737 | 0 | 0 | 100.00% |
| Near-misses | 6,503 | 6,503 | 0 | 0 | 100.00% |
| **All** | **9,240** | **9,240** | **0** | **0** | **100.00%** |

`select_column_expr`'s `ATTACHED_ALIAS` alternative gets no generated row: it
refuses every input (`../../../googlesql/testdata/tokens.tsv`). The optional
clauses that file marks `(every site)` add 403 rows, all agreeing.

### Comparison with Spanner emulator `1.5.54` (`spanner/parser/testdata/compat.tsv`)

The sections below predate the fork: `DialectDefault` (the BigQuery ∪ Spanner
union), `DialectBigQuery` (now `bigquery/parser`) and `DialectSpanner` (now
`spanner/parser`) were modes of one parser.

At judge revision 2, Spanner emulator `1.5.54` alone agreed with
`spanner/parser.Parse` (`DialectDefault`)
on **7,909 / 9,463 rows (83.58%)** overall (including 626 hand-added `suite:`
rows), or **7,335 / 8,837 (83.00%)** on the exact same 8,837 generated rows
(`2,196 / 2,614` alternatives, `5,139 / 6,223` near-misses, `1,093` only Spanner
accepts, `409` only omni accepts). That file now uses the BigQuery ∪ Spanner
union judge (revision 4) and agrees on all 9,463 rows.

Comparing `spanner/parser.Parse` (Spanner emulator `1.5.54`, judge revision 2) with
`bigquery/parser.Parse` (`DialectBigQuery`, live BigQuery dry-run) across all
8,837 generated rows:

| Spanner (`spanner/parser`) | BigQuery (`bigquery/parser`) | Rows | Meaning |
| --- | --- | ---: | --- |
| `reject` / `reject` (agree) | `reject` / `reject` (agree) | 5,096 | Both engines reject invalid syntax |
| `accept` / `accept` (agree) | `accept` / `accept` (agree) | 1,985 | Both engines accept shared GoogleSQL |
| `reject` / `accept` (Spanner-only) | `reject` / `reject` (agree) | 1,093 | Statements Spanner refused before parsing (pipes, `MERGE`, `LOAD DATA`, `BEGIN ... END`, `SET @@...`, `SHOW`, `EXPORT`, `IMPORT`, `EXECUTE IMMEDIATE`, duplicate `DYNAMIC LABEL`) that BigQuery parses and rejects as syntax errors |
| `accept` / `reject` (omni-only) | `accept` / `accept` (agree) | 391 | BigQuery-dialect statements Spanner refused (scripting, `QUALIFY`, `IS DISTINCT FROM`, `(SELECT NULL) AS a`, `CREATE TABLE a` without columns, `CREATE FUNCTION`/`PROCEDURE`/`ROW ACCESS POLICY`/`MATERIALIZED VIEW`, `GRANT`/`REVOKE`, `TRUNCATE`, `ANALYZE`) that BigQuery parses |
| `reject` / `reject` (agree) | `accept` / `accept` (agree) | 195 | BigQuery-dialect forms enabled in `DialectBigQuery` (`FROM`-first pipe queries, `CREATE MODEL` BigQuery ML clauses, `CREATE/DROP EXTERNAL TABLE/SCHEMA`, `ALTER SEARCH/VECTOR INDEX`, `ALTER [ALL] ROW ACCESS POLICY/POLICIES`, `PRIVILEGE RESTRICTION`, `CONNECTION`, `CONSTANT`, untyped `CREATE TABLE` columns, `COPY`/`CLONE` time travel, `CREATE INDEX` table alias, adjacent string concatenation in `GRANT`) |
| `accept` / `accept` (agree) | `reject` / `reject` (agree) | 59 | Spanner-dialect or disabled-in-BigQuery forms rejected in `DialectBigQuery` (braced constructors, `CREATE LOCALITY GROUP`, `FOR UPDATE`, `DESCRIBE ... \|>`, BigQuery-reserved `QUALIFY`/`MATCH_RECOGNIZE`, bare `\|> INSERT REPLACE/UPDATE`, strict `\|> EXPORT`, `CREATE INDEX a ON a()`, `LOAD DATA INTO a()`, `a : BEGIN END zzz`, `TABLE a WHERE NULL`) |
| `accept` / `reject` (omni-only) | `reject` / `reject` (agree) | 18 | Spanner-only DDL (`CREATE TABLE a PRIMARY KEY(...)`, `INTERLEAVE IN [PARENT]`, `CREATE NULL_FILTERED INDEX`, `ROW DELETION POLICY`) and `ANALYZE OPTIONS(zzz)` kept in `DialectDefault` for Spanner/legacy compatibility and rejected in `DialectBigQuery` |

## How the BigQuery dry-run judge works

`compatJudge` (`../compat_test.go`) runs **only** when `OMNI_RECORD=1` or
`OMNI_RECORD=verify` is set:

1. **Credentials & pin check.** Reads `BIGQUERY_CREDENTIALS` (key file path,
   JSON key document, or `"default"` for Application Default Credentials) and
   `BIGQUERY_PROJECT` (no `BIGQUERY_DATASET` is needed). Runs a single
   `SELECT 1` dry-run job and verifies its `gcp:bigquery:<project>:<location>`
   image and SHA-256 digest against the `default` row in `../../testdata/pins.tsv`.
2. **Dry-run execution.** Submits each row with `q.DryRun = true` and
   `q.UseLegacySQL = false` — executing nothing and mutating no datasets or
   tables. Retries with bounded exponential backoff **only** on transient HTTP
   429/5xx (`rateLimitExceeded`, `backendError`, `internalError`) or transport
   timeouts.
3. **Fail-closed error classification (`bigQueryVerdict`).**
   - `err == nil` → `accept`.
   - **HTTP 400 (`invalidQuery` or `invalid`):**
     - `reject` if the message starts with `"Syntax error:"` (covers unexpected
       tokens, unclosed string/comment/identifier literals, illegal input
       characters, missing whitespace, macro-disabled errors) or matches a
       GoogleSQL parser-action refusal in `grammarRefusals`
       (`../grammar_refusals_test.go`, mirrored from
       `spanner/spanner-oracle/grammar_refusals.go` plus the lookahead
       transformer refusal `EXCEPT must be followed by ALL, DISTINCT, or "("`
       and the `MergeDynamicLabelProperties` refusal `DYNAMIC LABEL|PROPERTIES
       cannot be used more than once in a single element table`).
     - Note: only the `"Syntax error:"` prefix is matched (not substrings), so
       syntactically valid outer statements such as `EXECUTE IMMEDIATE 'a'`
       (which BigQuery reports as ``Invalid EXECUTE IMMEDIATE sql string `a`,
       Syntax error: ...``) stay `accept`.
     - All other HTTP 400 `invalidQuery`/`invalid` errors (unqualified table,
       unrecognized name, type not found, `Statement not supported: ...`,
       semantic/analyzer checks) → `accept`.
   - **HTTP 404 (`notFound`):** `accept` only for per-statement SQL object
     lookups (`Not found: Table/Dataset/Routine/Model/Connection/PropertyGraph/...`),
     never for a missing GCP project.
   - **HTTP 403 (`accessDenied`):** `accept` only for per-statement SQL object
     permission checks (`Access Denied: Table/Dataset/Routine/Model/Connection ...`),
     never for project-level IAM, job-creation denial, or billing errors.
   - **All other errors** (HTTP 401, 429, 5xx, billing/quota, unknown 400/403/404
     reasons, empty messages, non-API/transport errors) → **fail closed** (return
     no verdict).
4. **Deterministic BigQuery analyzer HTTP 500 on `DECLARE ... FUNCTION<...>/MAP<...>`.**
   Eight generated rows (`DECLARE a FUNCTION <...>` and `DECLARE a MAP <...>`)
   parse cleanly in BigQuery's parser but deterministically trigger an HTTP 500
   `internalError` inside BigQuery's post-parse script analyzer. Rather than
   weakening fail-closed HTTP 500 handling, `confirmPostParseInternalError`
   accepts such a row **only** when: (a) it matches `DECLARE <var> FUNCTION|MAP <`,
   (b) a live `SELECT 1` dry-run probe immediately succeeds (proving BigQuery is
   healthy), and (c) submitting `sql + " zzz"` returns HTTP 400
   `"Syntax error: Expected end of input..."` (proving BigQuery's parser
   consumed the statement cleanly to EOF).

## Gap closure history (baseline 8,510 / 8,837 → 8,837 / 8,837)

The initial `#1536` baseline measured **8,510 / 8,837 (96.30%)** with **226**
`only BigQuery accepts` rows and **101** `only omni accepts` rows. All 327 gaps
are now closed via shared `spanner/parser` fixes (55 rows that also improved
`spanner/parser/testdata/compat.tsv`) plus `spanner/parser.DialectBigQuery`
(272 dialect-specific rows).

### Closed — previously only BigQuery accepted (226 rows: 185 alternatives, 41 near-misses)

Every one of these 226 rows was verified directly against BigQuery's error
messages: all 226 return `OK` or genuine post-parse semantic/qualification/
unsupported-statement messages (zero missed syntax errors).

| Family | Rows | Notes |
| --- | ---: | --- |
| `FROM`-first queries (`FROM a`, `FROM (SELECT NULL)`, `FROM UNNEST(...)`, `FROM a JOIN a`, `FROM GRAPH_TABLE(...)`, `FROM a PIVOT/UNPIVOT(...)`, `FROM a MATCH_RECOGNIZE(...)`, `WITH ... FROM a`) | 108 | Closed in `DialectBigQuery` (`parseQuery`, `parseTablePrimary` with `LATERAL`, `GRAPH_TABLE`, `MATCH_RECOGNIZE`). |
| BigQuery ML DDL (`CREATE MODEL`, `DROP MODEL`, `TRANSFORM(...)`, `REMOTE WITH CONNECTION`, `AS (alias AS (...))`, `CREATE TEMP MODEL`) | 37 | Shared base `CREATE [OR REPLACE] MODEL [IF NOT EXISTS] ... [REMOTE]` + `DROP MODEL` in `spanner/parser`; BigQuery ML clauses in `DialectBigQuery`. |
| External tables and schemas (`CREATE/DROP EXTERNAL TABLE`, `CREATE EXTERNAL SCHEMA`, `WITH PARTITION COLUMNS`, `WITH CONNECTION`) | 15 | Closed in `DialectBigQuery`. |
| Privilege restrictions (`CREATE/ALTER/DROP PRIVILEGE RESTRICTION`) | 14 | Closed in `DialectBigQuery`. |
| Search / vector index maintenance (`ALTER SEARCH INDEX`, `ALTER VECTOR INDEX` with `REBUILD`, `ADD/DROP/ALTER COLUMN`, `SET OPTIONS`) | 14 | Closed in `DialectBigQuery`. |
| Row access policy alterations (`ALTER [ALL] ROW ACCESS POLICY ... GRANT/REVOKE/FILTER/RENAME`) | 10 | Closed in `DialectBigQuery`. |
| Connections (`CREATE/DROP CONNECTION`) | 5 | Closed in `DialectBigQuery`. |
| `CREATE TABLE` typeless / trailing-comma column list forms (`CREATE TABLE a(a)`, `CREATE TABLE a(a,)`) | 5 | Closed in `DialectBigQuery`. |
| Arbitrary-expression dot-star (`SELECT NULL.*`, `\|> EXTEND NULL.*`, `\|> AGGREGATE NULL.*`) | 5 | Closed in shared `spanner/parser`. |
| `DEFINE TABLE` and `EXPLAIN` statements | 4 | Closed in shared `spanner/parser`. |
| Constants (`CREATE/DROP CONSTANT`) | 3 | Closed in `DialectBigQuery`. |
| `CREATE INDEX` table alias (`CREATE INDEX a ON a a(NULL)`) | 2 | Closed in `DialectBigQuery`. |
| `CREATE TABLE ... COPY/CLONE ... FOR SYSTEM_TIME AS OF` | 2 | Closed in `DialectBigQuery`. |
| Combined integer + brace statement hint (`@1@{ a = NULL } BEGIN`) | 1 | Closed in shared `spanner/parser`. |
| Adjacent string literal concatenation (`GRANT ALL ON a TO 'a' 'a'`) | 1 | Closed in `DialectBigQuery`. |

### Closed — previously only omni accepted (101 rows: 83 alternatives, 18 near-misses)

| Family | Rows | Resolution |
| --- | ---: | --- |
| Braced constructors (`SELECT { }`, `SELECT NEW a { }`, `SELECT STRUCT { }`, etc.) | 38 | Rejected in `DialectBigQuery` (`Braced constructors are not supported`). |
| Spanner-specific DDL (`CREATE TABLE a PRIMARY KEY(...)`, `INTERLEAVE IN [PARENT]`, `CREATE NULL_FILTERED INDEX`, `CREATE LOCALITY GROUP`, `ROW DELETION POLICY`) | 20 | Rejected in `DialectBigQuery`; retained in `DialectDefault` for Spanner. |
| Generic entity DDL (`CREATE a a`, `ALTER a SET AS 'a'`, `DROP a a`) | 15 | Restricted in shared `spanner/parser` to recognized generic entity types (`CAPACITY`, `RESERVATION`, `ASSIGNMENT`, `PROJECT`, `ORGANIZATION`, `FOLDER`, `PLACEMENT`). |
| Reserved keywords (`QUALIFY`, `MATCH_RECOGNIZE`, `GRAPH_TABLE`, `LATERAL`, `IGNORE`) | 8 | `GRAPH_TABLE`, `LATERAL`, and `.IGNORE` rejected in shared `spanner/parser`; `QUALIFY` and `MATCH_RECOGNIZE` reserved in `DialectBigQuery`. |
| Pipe `INSERT REPLACE` / `INSERT UPDATE` (`\|> INSERT REPLACE a`) | 4 | Bare `REPLACE`/`UPDATE` without `OR` rejected in `DialectBigQuery`. |
| Empty table element list (`CREATE TABLE a()`, `LOAD DATA INTO a() FROM FILES()`) | 3 | Empty `CREATE TABLE a()` without Spanner `PRIMARY KEY` rejected in shared `spanner/parser`; all empty element lists rejected in `DialectBigQuery`. |
| `FOR UPDATE` lock mode (`SELECT NULL FOR UPDATE`) | 2 | Rejected in `DialectBigQuery`. |
| Pipe operators on `DESCRIBE` (`DESCRIBE a \|> AGGREGATE`) | 2 | Rejected in `DialectBigQuery`. |
| Whitespace-separated dashed path in `DELETE` (`DELETE a - a`) | 2 | Closed in shared `spanner/parser` (`consumeDashedTail` requires adjacent `-`). |
| Strict syntax / near-miss over-accepts (`a : BEGIN END zzz`, `CREATE INDEX a ON a()`, `\|> EXPORT`, `\|> EXPORT AS`, `SELECT NULL QUALIFY`, `TABLE a WHERE NULL`, `ANALYZE OPTIONS(zzz)`) | 7 | Closed in `DialectBigQuery` without regressing Spanner `compat.tsv`. |

## Recording and verification commands

Offline replay (CI default, zero credentials, zero network):

```bash
go test ./bigquery/parser/ -run '^TestCompatibility$'
```

Live recording against `bytebase-dev` (creates/mutates nothing; ~13 minutes for
all 8,837 rows):

```bash
BIGQUERY_CREDENTIALS=default BIGQUERY_PROJECT=bytebase-dev OMNI_RECORD=1 go test -v -timeout 60m ./bigquery/parser/ -run '^TestCompatibility$'
```

Read-only live verification against `bytebase-dev` (re-judges all 8,837 rows,
writes nothing, fails on any diff):

```bash
BIGQUERY_CREDENTIALS=default BIGQUERY_PROJECT=bytebase-dev OMNI_RECORD=verify go test -v -timeout 60m ./bigquery/parser/ -run '^TestCompatibility$'
```
