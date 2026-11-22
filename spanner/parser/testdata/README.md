# Cloud Spanner parser test data

## `compat.tsv`: pinned Cloud Spanner emulator verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../../googlesql/testdata/GoogleSQL.g4` (converted from `googlesql.tm` at
GoogleSQL `2026.01.1` by `googlesql/parser/cmd/tm2g4`), followed by its
near-misses, with omni's verdict (`spanner/parser.Parse`, delegating to
`googlesql/parser.ParseWithDialect(..., DialectSpanner)`) and the pinned Cloud
Spanner emulator's verdict (`../../testdata/pins.tsv`, `parser` role, judge
revision 2: `gcr.io/cloud-spanner-emulator/emulator:1.5.54`).

`TestCompatibility` (`../compat_test.go`) reads the committed file through
`enginetest/compat` and requires **no Docker and no network access** in offline
replay. Recorded against the pinned Spanner emulator `1.5.54` through
`googlesql/spanner-oracle` (~39 seconds for all 8,837 rows) and verified
end-to-end with `OMNI_RECORD=verify` with zero differences.

### Current counts

| Rows | Total | Agree | Only Spanner accepts | Only omni accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 2,614 | 2,556 | 0 | 58 | 97.78% |
| Near-misses | 6,223 | 6,205 | 0 | 18 | 99.71% |
| **All** | **8,837** | **8,761** | **0** | **76** | **99.14%** |

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
| `CREATE FUNCTION` / `CREATE TABLE FUNCTION` / `CREATE PROCEDURE` | 73 | Rejects `CREATE PROCEDURE`, `CREATE TABLE FUNCTION`, `CREATE AGGREGATE FUNCTION`, scoped functions, `DETERMINISTIC`/`IMMUTABLE`/`STABLE`/`VOLATILE`, non-`REMOTE` `LANGUAGE`, BigQuery `REMOTE WITH CONNECTION`, `AS 'string'`, unnamed parameters, and missing `RETURNS` or body; preserves Spanner SQL UDFs `CREATE [OR REPLACE] FUNCTION f(x TYPE) RETURNS TYPE [SQL SECURITY INVOKER] AS (expr)` and Spanner remote functions `CREATE [OR REPLACE] FUNCTION f(x TYPE) RETURNS TYPE [NOT DETERMINISTIC] LANGUAGE REMOTE OPTIONS (...)`. |
| `CREATE TABLE` without column element list or `PRIMARY KEY` | 63 | Standalone Spanner `CREATE TABLE` requires `(...)` and a table-level or column-level `PRIMARY KEY` (while `\|> CREATE TABLE a` pipe syntax remains accepted). Rejects BigQuery-only `CTAS`, `LIKE`, `CLONE`, `COPY`, `PARTITION BY`, `CLUSTER BY`, `DEFAULT COLLATE`, `WITH CONNECTION`, and table-level `OPTIONS`. |
| `CREATE [UNIQUE\|NULL_FILTERED\|SEARCH\|VECTOR] INDEX` non-identifier keys & `ALL COLUMNS` | 39 | Requires column identifiers or parenthesized `( expression )` keys on ordinary indexes, column identifiers on `SEARCH`/`VECTOR` indexes and `STORING(...)`; rejects `(ALL COLUMNS)` and `PARTITION BY` on ordinary/vector indexes while preserving `CREATE SEARCH INDEX ... PARTITION BY`, `CREATE INDEX a ON a()`, `NULL_FILTERED`, `STORING`, `, INTERLEAVE IN`, `WHERE`, and Spanner `SEARCH`/`VECTOR` indexes. |
| `ANALYZE` with options or targets | 35 | Spanner only supports bare `ANALYZE`. |
| `GRANT` and `REVOKE` | 28 | Requires a target object type (`ON TABLE ...`) and `ROLE` grantees (`TO/FROM ROLE ...`); rejects `ALL [PRIVILEGES]` and string/parameter grantees. |
| `CREATE VIEW` / `MATERIALIZED` / `APPROX` / `RECURSIVE` | 17 | Rejects `MATERIALIZED`, `APPROX`, `RECURSIVE` views, view column lists, and view `OPTIONS`. |
| `DROP` BigQuery-only objects & `DROP SEARCH\|VECTOR INDEX ... ON ...` | 16 | Rejects `DROP DATABASE`, `PROCEDURE`, `MATERIALIZED\|APPROX VIEW`, `SNAPSHOT TABLE`, `EXTERNAL SCHEMA`, `TABLE FUNCTION`, `[ALL] ROW [ACCESS] POLICY/POLICIES`, generic entities, and the `ON table` tail on `DROP SEARCH\|VECTOR INDEX`. |
| `CREATE ROW [ACCESS] POLICY` & `CREATE SNAPSHOT TABLE` & generic entities | 15 | Rejects `CREATE ROW [ACCESS] POLICY`, `CREATE SNAPSHOT TABLE`, and BigQuery generic entities (`PROJECT`, `CAPACITY`, `RESERVATION`, `ASSIGNMENT`, etc.). |
| Query clauses (`QUALIFY`, `IS [NOT] DISTINCT FROM`, `(SELECT ...) AS a`) | 13 | Rejects `QUALIFY` clause, `IS [NOT] DISTINCT FROM`, and parenthesized outer query aliases while preserving `QUALIFY` as a non-reserved identifier (`SELECT QUALIFY`, `SELECT NULL QUALIFY`, `TABLE a() QUALIFY`). |
| `TRUNCATE`, non-`TABLE` `RENAME`, empty DDL `OPTIONS()` | 6 | Rejects `TRUNCATE TABLE`, `RENAME` of non-`TABLE` objects, and empty `OPTIONS()` in Spanner DDL (preserving non-DDL `OPTIONS()` in `MODULE`, `LOAD DATA`, and `SELECT WITH`). |

### Judge revision 2 (`PinRole: "parser"`)

At judge revision 1 (`default` pin), `compat.tsv` had **7,745 / 8,837 (87.64%)**
agreement with **1,092 `Only Spanner accepts`** and **0 `Only omni accepts`**.
Judge revision 2 eliminates all 1,092 judge false-accepts (**1,092 → 0**):

- **Parse-error fallback (`Statement not supported: <Kind>Stmt`, 1,080 rows):**
  When GoogleSQL's `ParseStatement` succeeds on an unsupported statement kind,
  it reports the AST node name ending in `Statement` or `Assignment`
  (`SubpipelineStatement`, `MergeStatement`, `BeginStatement`,
  `SingleAssignment`, etc. — all 533 remain `accept`). When `ParseStatement`
  fails on an unsupported statement kind, GoogleSQL's fallback path
  (`GetNextStatementKind`) replaces the syntax error with the `ResolvedNodeKind`
  name ending in `Stmt` (`SubpipelineStmt`, `MergeStmt`, `AuxLoadDataStmt`,
  `BeginStmt`, `UndropStmt`, etc.). Judge rev 2 classifies `...Stmt` as `reject`.
- **Lookahead tokenizer and dynamic grammar-action refusals (10 rows):**
  Rejects `EXCEPT must be followed by ALL, DISTINCT, or "("` (7 rows) and
  `DYNAMIC LABEL|PROPERTIES cannot be used more than once in a single element
  table` (3 rows, `MakeSyntaxError` in `googlesql.tm:2862`).
- **`CREATE VIEW` missing `SQL SECURITY` (2 rows):** Re-asks with a
  `SQL SECURITY INVOKER` stand-in before `AS` so broken view bodies
  (`CREATE VIEW a AS NULL`, `CREATE VIEW a AS SELECT`) reject on syntax while
  valid view bodies (`CREATE VIEW a AS SELECT NULL`) stay accepted.

### Remaining 76 `Only omni accepts` rows exposed by judge revision 2

Judge rev 2 revealed 76 rows where the Spanner emulator `1.5.54` genuinely
rejects parsing via the `...Stmt` parse-error fallback while `DialectSpanner`
currently accepts them (tracked for a dedicated `DialectSpanner` parser
follow-up):

- `|> MATCH_RECOGNIZE(...)` pipe operator (47 rows: `MATCH_RECOGNIZE` is
  non-reserved in Spanner emulator `1.5.54`)
- `BEGIN ... END` / `a : BEGIN END` scripting blocks (13 rows: `BeginStmt`;
  SQL transaction `BEGIN` remains accepted as `BeginStatement`)
- `LOAD DATA` with `PARTITIONS(...)`, `TEMP[ORARY] TABLE`, or empty `()` column
  list (7 rows: `AuxLoadDataStmt`)
- Bare `|> INSERT REPLACE` / `|> INSERT UPDATE` without `OR` (4 rows)
- Truncated `|> EXPORT` / `|> EXPORT AS` without `DATA` (2 rows)
- `DESCRIBE ... |> ...` pipe suffix on `DESCRIBE` (2 rows: `DescribeStmt`)
- `|> AGGREGATE GROUP AND ORDER BY ...` (1 row)

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
