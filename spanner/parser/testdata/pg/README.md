# Cloud Spanner PostgreSQL-dialect parser test data

## `compat.tsv`: pinned Cloud Spanner emulator verdicts (`DatabaseDialect_POSTGRESQL`)

`compat.tsv` holds a statement for each alternative of each rule in
PostgreSQL 17's grammar (`../../../../pg/testdata/gram.y` and
`../../../../pg/testdata/tokens.tsv`), followed by its near-misses (13,700
generated rows total; excluding the hand-added `suite:` rows that live in
`pg/parser/testdata/compat.tsv`), with omni's verdict
(`spanner/parser.ParsePostgres`, delegating to `pg/parser.Parse`) and the
pinned Cloud Spanner emulator's verdict in `DatabaseDialect_POSTGRESQL` mode
(`../../../testdata/pins.tsv`, `default` role:
`gcr.io/cloud-spanner-emulator/emulator:1.5.54`).

`TestPostgresCompatibility` (`../../pg_compat_test.go`) reads the committed
file through `enginetest/compat` and requires **no Docker and no network
access** in offline replay. Recorded live against the pinned Spanner emulator
`1.5.54` (`DatabaseDialect_POSTGRESQL`, ~66 seconds for all 13,700 rows) and
verified end-to-end with `OMNI_RECORD=verify` with zero differences.

### Current counts

| Rows | Total | Agree | Only Spanner-PG accepts | Only omni/pg accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 3,394 | 3,305 | 0 | 89 | 97.38% |
| Near-misses | 10,306 | 10,274 | 14 | 18 | 99.69% |
| **All** | **13,700** | **13,579** | **14** | **107** | **99.12%** |

13 `gram.y` alternatives get no generated row (same as `pg/parser`): the five
internal `MODE_*` entry points (`MODE_TYPE_NAME`, `MODE_PLPGSQL_EXPR`,
`MODE_PLPGSQL_ASSIGN1..3`) and the 8 rules reachable only from them.

### How the live judge works

`testcontainer.JudgePostgres` (`../../../internal/testcontainer/pg_judge.go`)
creates a fresh empty `DatabaseDialect_POSTGRESQL` database on the pinned
Spanner emulator `1.5.54` and routes each statement to the matching Spanner API:

- **Queries** (`SELECT`, `WITH`, `VALUES`, `TABLE`, `(`, utility statements,
  empty inputs) run in a single-use read-only snapshot (`Single().Query`).
  CTE-led DML that returns `DML statements may not be performed in single-use
  transactions` automatically falls back to the DML path.
- **DML** (`INSERT`, `UPDATE`, `DELETE`, `MERGE`) runs inside a
  `ReadWriteTransaction` that always aborts before commit so no state is
  mutated.
- **DDL** (`CREATE`, `ALTER`, `DROP`, `GRANT`, `REVOKE`, `RENAME`, `TRUNCATE`,
  `ANALYZE`) runs via `UpdateDatabaseDdl` against a dedicated empty
  `DatabaseDialect_POSTGRESQL` database that is dropped and recreated whenever
  a DDL statement succeeds (7 resets across all 13,700 rows), preventing any
  earlier statement from influencing a later verdict.

`spannerPGVerdict` (`../../pg_compat_test.go`) is fail-closed:

- **Reject** on raw-parser, lexer, single-statement cardinality, and
  PostgreSQL `gram.y` parser-action errors (`[ERROR] syntax error at ...`,
  `; failed to parse the DDL statements.`, `SQL queries are limited to single
  statements`, `[ERROR] XML expressions are not supported`, and PostgreSQL
  `scan.l` / `gram.y` action errors such as `CREATE ASSERTION`, `OVERLAPS`
  arity, invalid window frame bounds, `WITH TIES` without `ORDER BY`,
  `LIMIT #,#`, `UNIQUE (SELECT)`, `SET CATALOG`, unrecognized role/column/row
  security/partitioning options).
- **Accept** on successful execution and post-parse semantic, catalog, runtime,
  or unsupported-statement errors (`relation "..." does not exist`,
  `column "..." does not exist`, `Statement is not supported.`,
  `Utility SQL statements are not supported`,
  `statement type "MergeStmt" is not supported`,
  `Primary key must be defined...`, `Table not found`, `OutOfRange`, and
  post-parse DDL `GOOGLESQL_RET_CHECK` quirks).
- **Error (no verdict)** on any non-gRPC error, transport/lifecycle status
  (`Unavailable`, `DeadlineExceeded`, `Canceled`, `Aborted`, etc.),
  resource-level `NotFound`/`AlreadyExists`, generic `Internal`, or unexpected
  gRPC status code.

### Spanner PostgreSQL-dialect findings (`emulator:1.5.54`)

Cloud Spanner's PostgreSQL dialect embeds PostgreSQL's Bison parser (based on
**PostgreSQL 16**) for both queries/DML and DDL. Unsupported PostgreSQL
statements (`BEGIN`, `COMMIT`, `LISTEN`, `NOTIFY`, `VACUUM`, `EXPLAIN`, `COPY`,
`CREATE TRIGGER`, `CREATE EXTENSION`, `CREATE POLICY`, `CREATE RULE`,
`CREATE TYPE`, `CREATE DOMAIN`, `CREATE MATERIALIZED VIEW`, `MERGE`, `ANALYZE`,
`TRUNCATE`, etc.) are **parsed in full** before Spanner rejects them post-parse
with `Utility SQL statements are not supported`, `Statement is not supported.`,
or `statement type "MergeStmt" is not supported` — while malformed near-misses
of those same statements fail at parse time with `[ERROR] syntax error ...`.

#### `Only omni/pg accepts` (107 rows: 89 alternatives, 18 near-misses)

Every row omni/pg accepts and Spanner-PG rejects falls into five clean groups:

| Family | Rows | Notes |
| --- | ---: | --- |
| XML expressions disabled in Spanner's parser (`IS [NOT] DOCUMENT`, `XMLCONCAT`, `XMLELEMENT`, `XMLATTRIBUTES`, `XMLFOREST`, `XMLPARSE`, `XMLPI`, `XMLROOT`) | 45 | Spanner's PG parser aborts during parsing with `[ERROR] XML expressions are not supported` (even before trailing invalid tokens). |
| PostgreSQL 17 SQL/JSON query functions & clauses (`JSON_QUERY(... PASSING / WRAPPER / ON EMPTY / ON ERROR / KEEP\|OMIT QUOTES)`, `JSON_EXISTS(... ON ERROR)`) | 30 | Spanner `1.5.54` is based on PostgreSQL 16, before PG 17 added SQL/JSON query function grammar. |
| Other PostgreSQL 17 grammar additions (`ALTER DOMAIN ... ADD [CONSTRAINT] NOT NULL [...]` (8), `COPY ... FORCE [NOT] NULL *` / `(opt DEFAULT)` (3), `MERGE ... WHEN NOT MATCHED BY SOURCE\|TARGET` (2), `AT LOCAL` (2), `ALTER COLUMN ... SET EXPRESSION AS` (1), `SET ACCESS METHOD DEFAULT` (1), `SET STATISTICS DEFAULT` (1)) | 18 | Added in PostgreSQL 17; rejected with `syntax error` by Spanner's PG 16-based parser. |
| Spanner-PG grammar restrictions (`ALTER TYPE ... SET(...)` (5), `ALTER INDEX [IF EXISTS] ... DROP ...` (2), `CREATE RULE ... ON UPDATE` (1), `CLUSTER(...)` (1)) | 9 | Rejected at parse time in Spanner's modified `gram.y`. |
| Empty and bare-semicolon inputs (`""` (4), `";"` (1)) | 5 | PostgreSQL's raw parser accepts 0-statement inputs; Spanner requires exactly one statement (`SQL queries are limited to single statements`). |

#### `Only Spanner-PG accepts` (14 rows: 0 alternatives, 14 near-misses)

Every row Spanner-PG accepts and omni/pg rejects is a PG 17 near-miss:

| Family | Rows | Notes |
| --- | ---: | --- |
| PG 17 SQL/JSON keywords parsed as ordinary functions in PG 16 (`SELECT JSON_EXISTS(a)`, `SELECT FROM JSON_EXISTS(a)`, `SELECT JSON()`, `SELECT JSON_SCALAR()`, `SELECT JSON_SERIALIZE()`, `SELECT JSON_TABLE(a)`, `SELECT JSON_SCALAR(a, a)`) | 7 | Not special grammar keywords in Spanner's PG 16-based parser, so they parse as ordinary function calls and fail post-parse during name resolution. |
| Permissive `ON UPDATE <action>` in `key_action` (`CREATE DOMAIN a a REFERENCES a ON UPDATE ACTION\|RETURN\|CASCADED\|NULL\|SETOF NULL`) | 5 | Parsed by Spanner's modified `key_action` rule, then rejected post-parse (`Statement is not supported.`). |
| Spanner-PG `ALTER TABLE` / `ALTER DATABASE` grammar extensions (`ALTER TABLE a SET`, `ALTER DATABASE a TARGET 1.5`) | 2 | Parsed by Spanner's modified DDL grammar before failing post-parse (`GOOGLESQL_RET_CHECK` in `pg_to_spanner_ddl_translator.cc` / `Statement is not supported.`). |

#### Spanner-PG-specific DDL extensions (beyond upstream `pg/testdata/gram.y`)

Cloud Spanner's PostgreSQL dialect also adds Spanner-specific DDL syntax that
is not part of upstream PostgreSQL 17's `gram.y` (and therefore not among the
13,700 generated rows here), including:

- Table interleaving: `CREATE TABLE child (...) INTERLEAVE IN PARENT parent [ON DELETE {CASCADE | NO ACTION}]`
- Row deletion policy (TTL): `CREATE TABLE t (..., ts timestamptz) TTL INTERVAL '30 days' ON ts`
- Change streams: `CREATE CHANGE STREAM cs FOR ALL`
- Spanner views: `CREATE VIEW v SQL SECURITY INVOKER AS SELECT ...`
- Spanner sequences: `CREATE SEQUENCE s BIT_REVERSED_POSITIVE` / `ALTER DATABASE db SET spanner.default_sequence_kind = 'bit_reversed_positive'`

Because `ParsePostgres` currently delegates directly to `pg/parser.Parse`
(keeping `pg/parser` strictly aligned with upstream PostgreSQL 17), this
baseline records the live Spanner-PG `1.5.54` compatibility over all 13,700
upstream PostgreSQL 17 grammar rows without modifying `pg/parser`.

## Recording and verification commands

Offline replay (CI default, zero Docker, zero network):

```bash
go test ./spanner/parser/ -run '^TestPostgresCompatibility$'
```

Live recording against pinned `gcr.io/cloud-spanner-emulator/emulator:1.5.54`
(~66 seconds for all 13,700 rows):

```bash
OMNI_RECORD=1 go test -v -timeout 30m ./spanner/parser/ -run '^TestPostgresCompatibility$'
```

Read-only live verification against the pinned emulator (re-judges all 13,700
rows, writes nothing, fails on any diff):

```bash
OMNI_RECORD=verify go test -v -timeout 30m ./spanner/parser/ -run '^TestPostgresCompatibility$'
```
