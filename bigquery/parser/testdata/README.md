# BigQuery parser test data

## `compat.tsv`: live BigQuery dry-run verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../../googlesql/testdata/GoogleSQL.g4` (converted from `googlesql.tm` at
GoogleSQL `2026.01.1` by `googlesql/parser/cmd/tm2g4`), followed by its
near-misses, with omni's verdict (`bigquery/parser.Parse`, delegating to
`googlesql/parser.Parse`) and live BigQuery dry-run's verdict on the pinned
Google Cloud project (`../../testdata/pins.tsv`, `default` role).

`TestCompatibility` (`../compat_test.go`) reads the committed file through
`enginetest/compat` and requires **no Google Cloud credentials and no network
access** in offline replay. Recorded 2026-11-21 against `bytebase-dev` (`US`)
via BigQuery dry-run jobs (`DryRun = true`, `UseLegacySQL = false`), which
execute nothing, cost $0, and create or mutate no datasets or tables, then
verified end-to-end with `OMNI_RECORD=verify` across all 8,837 rows with zero
differences (793 seconds).

### Baseline counts

| Rows | Total | Agree | Only BigQuery accepts | Only omni accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 2,614 | 2,345 | 186 | 83 | 89.71% |
| Near-misses | 6,223 | 6,162 | 43 | 18 | 99.02% |
| **All** | **8,837** | **8,507** | **229** | **101** | **96.27%** |

`select_column_expr`'s `ATTACHED_ALIAS` alternative gets no generated row: it
refuses every input (`../../../googlesql/testdata/tokens.tsv`).

### Comparison with Spanner emulator `1.5.54` (`googlesql/parser/testdata/compat.tsv`)

Spanner emulator `1.5.54` agrees with omni on **7,854 / 9,463 rows (83.00%)**
overall (including 626 hand-added `suite:` rows), or **7,280 / 8,837 (82.38%)**
on the exact same 8,837 generated rows (`2,146 / 2,614` alternatives,
`5,134 / 6,223` near-misses, `1,124` only Spanner accepts, `433` only omni
accepts).

Moving the judge from Spanner emulator `1.5.54` to live BigQuery dry-run on
those exact 8,837 generated rows produces:

| Spanner (`googlesql/parser`) | BigQuery (`bigquery/parser`) | Rows | Meaning |
| --- | --- | ---: | --- |
| `reject` / `reject` (agree) | `reject` / `reject` (agree) | 5,072 | Both engines reject invalid syntax |
| `accept` / `accept` (agree) | `accept` / `accept` (agree) | 1,954 | Both engines accept shared GoogleSQL |
| `reject` / `accept` (Spanner-only) | `reject` / `reject` (agree) | 1,090 | Statements Spanner refused before parsing (pipes, `MERGE`, `LOAD DATA`, `BEGIN ... END`, `SET @@...`, `SHOW`, `EXPORT`, `IMPORT`, `EXECUTE IMMEDIATE`) that BigQuery parses and rejects as syntax errors |
| `accept` / `reject` (omni-only) | `accept` / `accept` (agree) | 391 | BigQuery-dialect statements Spanner refused (scripting `DECLARE`/`IF`/`CASE`/`WHILE`/`LOOP`/`REPEAT`/`FOR`/`BREAK`/`RETURN`/`RAISE`, `QUALIFY`, `IS DISTINCT FROM`, `(SELECT NULL) AS a`, `CREATE TABLE a` without columns, `CREATE FUNCTION`/`PROCEDURE`/`ROW ACCESS POLICY`/`MATERIALIZED VIEW`, `GRANT`/`REVOKE`, `TRUNCATE`, `ANALYZE`) that BigQuery parses |
| `reject` / `reject` (agree) | `reject` / `accept` (BQ-only) | 195 | BigQuery-dialect forms both omni and Spanner rejected (`FROM`-first pipe queries, `CREATE MODEL`, `CREATE EXTERNAL TABLE/SCHEMA`, `ALTER SEARCH/VECTOR INDEX`, `ALTER ROW ACCESS POLICY`, `PRIVILEGE RESTRICTION`, `CONNECTION`, etc.) |
| `accept` / `accept` (agree) | `accept` / `reject` (omni-only) | 59 | Spanner-dialect or disabled-in-BigQuery forms (braced constructors, Spanner `PRIMARY KEY` / `INTERLEAVE` / `NULL_FILTERED` / `LOCALITY GROUP` / `ROW DELETION POLICY`, `FOR UPDATE`, `DESCRIBE ... |>`) |
| `accept` / `reject` (omni-only) | `accept` / `reject` (omni-only) | 42 | Generic entity DDL (`CREATE/ALTER/DROP a ...`), empty `CREATE TABLE a()`, reserved keywords, and strict near-misses rejected by both engines |
| `reject` / `accept` (Spanner-only) | `reject` / `accept` (BQ-only) | 34 | Shared GoogleSQL forms accepted by both engines that omni still rejects (`CREATE MODEL`, `DROP MODEL`, `SELECT NULL.*`, `EXPLAIN`, `CREATE CONSTANT`, `DEFINE TABLE`, `DYNAMIC LABEL` duplicates) |

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
       `googlesql/spanner-oracle/grammar_refusals.go` plus the lookahead
       transformer refusal `EXCEPT must be followed by ALL, DISTINCT, or "("`).
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

## Gap breakdown

### Only BigQuery accepts (229 rows: 186 alternatives, 43 near-misses)

Every one of these 229 rows was verified directly against BigQuery's error
messages: all 229 return `OK` or genuine post-parse semantic/qualification/
unsupported-statement messages (zero missed syntax errors).

| Family | Rows | Notes |
| --- | ---: | --- |
| `FROM`-first queries (`FROM a`, `FROM (SELECT NULL)`, `FROM UNNEST(...)`, `FROM a JOIN a`, `FROM GRAPH_TABLE(...)`, `FROM a PIVOT/UNPIVOT(...)`, `FROM a MATCH_RECOGNIZE(...)`, `WITH ... FROM a`) | 108 | BigQuery natively supports standalone `FROM`-first pipe queries; `googlesql/parser` currently rejects a leading `FROM` (inherited from Spanner). |
| BigQuery ML DDL (`CREATE MODEL`, `DROP MODEL`, `TRANSFORM(...)`, `REMOTE WITH CONNECTION`, `AS (alias AS (...))`, `CREATE TEMP MODEL`) | 37 | Not yet parsed by `googlesql/parser`. |
| External tables and schemas (`CREATE/DROP EXTERNAL TABLE`, `CREATE EXTERNAL SCHEMA`, `WITH PARTITION COLUMNS`, `WITH CONNECTION`) | 15 | Not yet parsed by `googlesql/parser`. |
| Privilege restrictions (`CREATE/ALTER/DROP PRIVILEGE RESTRICTION`) | 14 | Parsed by BigQuery's GoogleSQL frontend (reported as `Statement not supported`). |
| Search / vector index maintenance (`ALTER SEARCH INDEX`, `ALTER VECTOR INDEX` with `REBUILD`, `ADD/DROP/ALTER COLUMN`, `SET OPTIONS`) | 14 | Not yet parsed by `googlesql/parser`. |
| Row access policy alterations (`ALTER [ALL] ROW ACCESS POLICY ... GRANT/REVOKE/FILTER/RENAME`) | 10 | Parsed by BigQuery's GoogleSQL frontend (`Statement not supported`). |
| Connections (`CREATE/DROP CONNECTION`) | 5 | Not yet parsed by `googlesql/parser`. |
| `CREATE TABLE` typeless / trailing-comma column list forms (`CREATE TABLE a(a)`, `CREATE TABLE a(a,)`) | 5 | BigQuery parses untyped column entries in `CREATE TABLE` element lists (e.g., for `AS SELECT` / schema auto-detect). |
| Arbitrary-expression dot-star (`SELECT NULL.*`, `|> EXTEND NULL.*`, `|> AGGREGATE NULL.*`) | 5 | `googlesql/parser` restricts `.*` to path expressions. |
| `DEFINE TABLE` and `EXPLAIN` statements | 4 | Parsed by BigQuery's GoogleSQL frontend (`Statement not supported`). |
| Constants (`CREATE/DROP CONSTANT`) | 3 | Not yet parsed by `googlesql/parser`. |
| Duplicate `DYNAMIC LABEL` in `CREATE PROPERTY GRAPH` | 3 | BigQuery reports a post-parse duplicate-label error even before EOF on the near-miss twin. |
| `CREATE INDEX` table alias (`CREATE INDEX a ON a a(NULL)`) | 2 | BigQuery parses the optional alias before reporting `CREATE INDEX is not supported`. |
| `CREATE TABLE ... COPY/CLONE ... FOR SYSTEM_TIME AS OF` | 2 | Time-travel clause on `COPY`/`CLONE` source. |
| Combined integer + brace statement hint (`@1@{ a = NULL } BEGIN`) | 1 | `googlesql/parser` handles `@1` and `@{...}` separately, not stacked. |
| Adjacent string literal concatenation (`GRANT ALL ON a TO 'a' 'a'`) | 1 | GoogleSQL concatenates adjacent string literals in the lexer/parser. |

### Only omni accepts (101 rows: 83 alternatives, 18 near-misses)

| Family | Rows | BigQuery's answer |
| --- | ---: | --- |
| Braced constructors (`SELECT { }`, `SELECT NEW a { }`, `SELECT STRUCT { }`, etc.) | 38 | `Braced constructors are not supported` (feature disabled in BigQuery). |
| Spanner-specific DDL (`CREATE TABLE a PRIMARY KEY(...)`, `INTERLEAVE IN [PARENT]`, `CREATE NULL_FILTERED INDEX`, `CREATE LOCALITY GROUP`, `ROW DELETION POLICY`) | 20 | Grammar refusals (`PRIMARY KEY must be defined in the table element list...`, `null_filtered is not a supported object type`, `CREATE LOCALITY GROUP is not supported`, `ROW DELETION POLICY clause is not supported`) or syntax errors on `, INTERLEAVE IN`. Shared `googlesql/parser` intentionally accepts Spanner DDL. |
| Generic entity DDL (`CREATE a a`, `ALTER a SET AS 'a'`, `DROP a a`) | 15 | `a is not a supported object type` (grammar refusal on unknown entity kind). |
| BigQuery reserved keywords (`QUALIFY`, `MATCH_RECOGNIZE`, `GRAPH_TABLE`, `LATERAL`, `IGNORE`) | 8 | `Syntax error: Unexpected keyword ...` — BigQuery reserves these keywords in identifier/field positions where `googlesql/parser` permits them (for Spanner compatibility). |
| Pipe `INSERT REPLACE` / `INSERT UPDATE` (`\|> INSERT REPLACE a`) | 4 | `Syntax error: Expected end of input but got identifier "a"`. |
| Empty table element list (`CREATE TABLE a()`, `LOAD DATA INTO a() FROM FILES()`) | 3 | `A table must define at least one column.` |
| `FOR UPDATE` lock mode (`SELECT NULL FOR UPDATE`) | 2 | `FOR UPDATE is not supported` (Spanner feature disabled in BigQuery). |
| Pipe operators on `DESCRIBE` (`DESCRIBE a \|> AGGREGATE`) | 2 | `Pipe operators are not supported on this statement`. |
| Dashed path in `DELETE` (`DELETE a - a`) | 2 | `Syntax error: Unexpected "-"`. |
| Strict syntax / near-miss over-accepts (`a : BEGIN END zzz`, `CREATE INDEX a ON a()`, `\|> EXPORT`, `\|> EXPORT AS`, `SELECT NULL QUALIFY`, `TABLE a WHERE NULL`, `ANALYZE OPTIONS(zzz)`) | 7 | Real `googlesql/parser` over-accepts (mismatched end label, empty index key list, missing `DATA` after `\|> EXPORT`, missing `QUALIFY` predicate, `TABLE ... WHERE`, malformed `OPTIONS(zzz)`). |

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
