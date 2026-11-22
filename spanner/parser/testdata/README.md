# Cloud Spanner parser test data

## `compat.tsv`: pinned Cloud Spanner emulator verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../../googlesql/testdata/GoogleSQL.g4` (converted from `googlesql.tm` at
GoogleSQL `2026.01.1` by `googlesql/parser/cmd/tm2g4`), followed by its
near-misses, with omni's verdict (`spanner/parser.Parse`, delegating to
`googlesql/parser.ParseWithDialect(..., DialectSpanner)`) and the pinned Cloud
Spanner emulator's verdict (`../../testdata/pins.tsv`, `default` role:
`gcr.io/cloud-spanner-emulator/emulator:1.5.54`).

`TestCompatibility` (`../compat_test.go`) reads the committed file through
`enginetest/compat` and requires **no Docker and no network access** in offline
replay. Recorded against the pinned Spanner emulator `1.5.54` through
`googlesql/spanner-oracle` (~47 seconds for all 8,837 rows) and verified
end-to-end with `OMNI_RECORD=verify` with zero differences.

### Current counts

| Rows | Total | Agree | Only Spanner accepts | Only omni accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 2,614 | 2,576 | 38 | 0 | 98.55% |
| Near-misses | 6,223 | 5,169 | 1,054 | 0 | 83.06% |
| **All** | **8,837** | **7,745** | **1,092** | **0** | **87.64%** |

`select_column_expr`'s `ATTACHED_ALIAS` alternative gets no generated row: it
refuses every input (`../../../googlesql/testdata/tokens.tsv`).

### Comparison with union `googlesql/parser` (`googlesql/parser/testdata/compat.tsv`)

In union mode (`googlesql/parser.Parse`, `DialectDefault`), the 8,837 generated
rows had **7,335 / 8,837** agreeing (`1,093` only Spanner accepts, `409` only
omni accepts; plus 42 `suite:` only-omni rows = 451 total only-omni in
`googlesql/parser/testdata/compat.tsv`).

`DialectSpanner` closes **all 409 generated `only omni accepts` rows (409 → 0)**
and also closes the one non-dispatch `TABLE` near-miss (`TABLE a() QUALIFY`,
where `QUALIFY` is a non-reserved identifier in Spanner: `1,093 → 1,092`),
lifting generated-row agreement from **7,335 / 8,837 (83.00%)** to
**7,745 / 8,837 (87.64%)** with **0 over-acceptances**:

| BigQuery-only family rejected in `DialectSpanner` | Closed generated rows | Notes |
| --- | ---: | --- |
| Procedural scripting (`DECLARE`, `IF`, `CASE`, `WHILE`, `LOOP`, `REPEAT`, `FOR ... IN`, `BREAK`, `LEAVE`, `CONTINUE`, `ITERATE`, `RETURN`, `RAISE`, labeled loops) | 104 | Standalone scripting statements rejected in `DialectSpanner`; `BEGIN ... END`, `a : BEGIN END`, `SET`, `CALL`, and `EXECUTE IMMEDIATE` retained. |
| `CREATE FUNCTION` / `CREATE TABLE FUNCTION` / `CREATE PROCEDURE` | 73 | Rejects `CREATE PROCEDURE`, `CREATE TABLE FUNCTION`, `CREATE AGGREGATE FUNCTION`, scoped functions, `DETERMINISTIC`/`IMMUTABLE`/`STABLE`/`VOLATILE`, `LANGUAGE`, `REMOTE`, `AS 'string'`, unnamed parameters, and missing `RETURNS` or `AS (...)`; preserves Spanner SQL UDFs `CREATE [OR REPLACE] FUNCTION f(x TYPE) RETURNS TYPE [SQL SECURITY INVOKER] AS (expr)`. |
| `CREATE TABLE` without column element list or `PRIMARY KEY` | 63 | Standalone Spanner `CREATE TABLE` requires `(...)` and a table-level or column-level `PRIMARY KEY` (while `\|> CREATE TABLE a` pipe syntax remains accepted). Rejects BigQuery-only `CTAS`, `LIKE`, `CLONE`, `COPY`, `PARTITION BY`, `CLUSTER BY`, `DEFAULT COLLATE`, `WITH CONNECTION`, and table-level `OPTIONS`. |
| `CREATE [UNIQUE\|NULL_FILTERED\|SEARCH\|VECTOR] INDEX` non-identifier keys & `ALL COLUMNS` | 39 | Requires column identifiers for index keys and `STORING(...)`; rejects `(ALL COLUMNS)` and `PARTITION BY`. Preserves `CREATE INDEX a ON a()`, `NULL_FILTERED`, `STORING`, `, INTERLEAVE IN`, `WHERE`, and Spanner `SEARCH`/`VECTOR` indexes. |
| `ANALYZE` with options or targets | 35 | Spanner only supports bare `ANALYZE`. |
| `GRANT` and `REVOKE` | 28 | Requires a target object type (`ON TABLE ...`) and `ROLE` grantees (`TO/FROM ROLE ...`); rejects `ALL [PRIVILEGES]` and string/parameter grantees. |
| `CREATE VIEW` / `MATERIALIZED` / `APPROX` / `RECURSIVE` | 17 | Rejects `MATERIALIZED`, `APPROX`, `RECURSIVE` views, view column lists, and view `OPTIONS`. |
| `DROP` BigQuery-only objects & `DROP SEARCH\|VECTOR INDEX ... ON ...` | 16 | Rejects `DROP DATABASE`, `PROCEDURE`, `MATERIALIZED\|APPROX VIEW`, `SNAPSHOT TABLE`, `EXTERNAL SCHEMA`, `TABLE FUNCTION`, `[ALL] ROW [ACCESS] POLICY/POLICIES`, generic entities, and the `ON table` tail on `DROP SEARCH\|VECTOR INDEX`. |
| `CREATE ROW [ACCESS] POLICY` & `CREATE SNAPSHOT TABLE` & generic entities | 15 | Rejects `CREATE ROW [ACCESS] POLICY`, `CREATE SNAPSHOT TABLE`, and BigQuery generic entities (`PROJECT`, `CAPACITY`, `RESERVATION`, `ASSIGNMENT`, etc.). |
| Query clauses (`QUALIFY`, `IS [NOT] DISTINCT FROM`, `(SELECT ...) AS a`) | 13 | Rejects `QUALIFY` clause, `IS [NOT] DISTINCT FROM`, and parenthesized outer query aliases while preserving `QUALIFY` as a non-reserved identifier (`SELECT QUALIFY`, `SELECT NULL QUALIFY`, `TABLE a() QUALIFY`). |
| `TRUNCATE`, non-`TABLE` `RENAME`, empty DDL `OPTIONS()` | 6 | Rejects `TRUNCATE TABLE`, `RENAME` of non-`TABLE` objects, and empty `OPTIONS()` in Spanner DDL (preserving non-DDL `OPTIONS()` in `MODULE`, `LOAD DATA`, and `SELECT WITH`). |

### Why 1,092 rows remain in `Only Spanner accepts`

Every remaining `Only Spanner accepts` row (38 alternatives, 1,054 near-misses)
is a statement the Spanner emulator `1.5.54` refuses **before parsing** with
`Statement not supported` (or an analysis error after shallow dispatch), which
the oracle harness records as `accept` for valid and broken text alike (see
`../../../googlesql/parser/testdata/README.md` for the full measurement):

- Pipe syntax (`|>`, 750 rows: all 33 pipe-only alternatives are GoogleSQL's
  explicit error-recovery productions such as `|> LOG(SELECT`, `|> GROUP`,
  `|> SET a.`, `|> ELSEIF`, `|> EXPORT DATA AS`, plus 717 malformed pipe
  near-misses)
- Pre-parse `Statement not supported` dispatch families: `MERGE` (80),
  `LOAD DATA` (71), `BEGIN ... END` (42), `EXECUTE IMMEDIATE` (21), `EXPORT`
  (18), `IMPORT` (16), `UNDROP` (13, including 4 `UNDROP CONNECTION`
  alternatives), `SHOW` (12), `CLONE DATA` (10), `SET TRANSACTION` / `SET @@...`
  (10), `MODULE` (7), `ASSERT` (6), `DESCRIBE` (5), `DEFINE` (4), `COMMIT` /
  `ROLLBACK` (6), `EXPLAIN` (2)
- Non-syntax post-dispatch errors: `CREATE PROPERTY GRAPH` duplicate
  `DYNAMIC LABEL` / stray `)` (4) and `SELECT EXCEPT` / shallow query errors (6)

Matching those cells would mean accepting invalid SQL that GoogleSQL's own
grammar refuses.

## Recording and verification commands

Offline replay (CI default, zero Docker, zero network):

```bash
go test ./spanner/parser/ -run '^TestCompatibility$'
```

Live recording against pinned `gcr.io/cloud-spanner-emulator/emulator:1.5.54`
(~47 seconds for all 8,837 rows):

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./spanner/parser/ -run '^TestCompatibility$'
```

Read-only live verification against the pinned emulator (re-judges all 8,837
rows, writes nothing, fails on any diff):

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./spanner/parser/ -run '^TestCompatibility$'
```
