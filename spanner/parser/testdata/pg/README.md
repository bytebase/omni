# Cloud Spanner PostgreSQL-dialect parser test data

## `compat.tsv`: pinned Cloud Spanner emulator verdicts (`DatabaseDialect_POSTGRESQL`)

`compat.tsv` holds a statement for each alternative of each rule in
PostgreSQL 17's grammar (`../../../../pg/testdata/gram.y` and
`../../../../pg/testdata/tokens.tsv`), followed by its near-misses (13,971
generated rows total), plus the `suite:spanner-pg-ddl` corpus (143 hand-curated
Cloud Spanner PostgreSQL DDL extension statements and near-misses), with omni's
verdict (`spanner/parser.ParsePostgres`, delegating to
`pg/parser.ParseWithDialect(..., pgparser.DialectSpannerPG)`) and the pinned
Cloud Spanner emulator's verdict in `DatabaseDialect_POSTGRESQL` mode
(`../../../testdata/pins.tsv`, `default` role:
`gcr.io/cloud-spanner-emulator/emulator:1.5.54`).

`TestPostgresCompatibility` (`../../pg_compat_test.go`) reads the committed
file through `enginetest/compat` and requires **no Docker and no network
access** in offline replay. Recorded live against the pinned Spanner emulator
`1.5.54` (`DatabaseDialect_POSTGRESQL`) and verified end-to-end with
`OMNI_RECORD=verify` with zero differences.

### Current counts

| Rows | Total | Agree | Only Spanner-PG accepts | Only omni/pg accepts | Agreement |
| --- | ---: | ---: | ---: | ---: | ---: |
| One per alternative | 3,591 | 3,591 | 0 | 0 | 100.00% |
| Near-misses | 10,380 | 10,380 | 0 | 0 | 100.00% |
| `suite:spanner-pg-ddl` | 143 | 143 | 0 | 0 | 100.00% |
| **All** | **14,114** | **14,114** | **0** | **0** | **100.00%** |

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
  a DDL statement succeeds, preventing any earlier statement from influencing
  a later verdict.

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

### `DialectSpannerPG` dialect behavior (`emulator:1.5.54`)

Cloud Spanner's PostgreSQL dialect embeds PostgreSQL's Bison parser (based on
**PostgreSQL 16**) with Spanner-specific DDL extensions for both queries/DML
and DDL. `pg/parser.ParseWithDialect(..., pgparser.DialectSpannerPG)` mirrors
this dialect while keeping `pg/parser.Parse` (`DialectDefault`) strictly
PostgreSQL 17:

#### Spanner-PG parser restrictions (0 `Only omni/pg accepts`)

All 107 upstream PG 17 disagreements from the baseline `pg/parser.Parse` run
are closed cleanly at their grammar points in `DialectSpannerPG`, along with
the 7 PG 16 SQL/JSON function-name near-misses:

- **XML expressions disabled**: rejects `IS [NOT] DOCUMENT` (in `a_expr` and
  `b_expr`), `XMLCONCAT`, `XMLELEMENT`, `XMLFOREST`, `XMLPARSE`, `XMLPI`, and
  `XMLROOT`, while retaining `XMLEXISTS`, `XMLSERIALIZE`, and `XMLTABLE`.
- **PostgreSQL 16 SQL/JSON surface**: treats PG 17 SQL/JSON query keywords
  (`JSON(...)`, `JSON_SCALAR`, `JSON_SERIALIZE`, `JSON_QUERY`, `JSON_EXISTS`,
  `JSON_VALUE`, `JSON_TABLE`) as ordinary function identifiers rather than
  PG 17 special-syntax productions, while retaining PG 16 `IS [NOT] JSON`,
  `JSON 'literal'`, `JSON_OBJECT`, `JSON_ARRAY`, `JSON_OBJECTAGG`, and
  `JSON_ARRAYAGG`.
- **PostgreSQL 17 grammar additions gated off**: `ALTER DOMAIN ... ADD [CONSTRAINT] NOT NULL`,
  `COPY ... FORCE [NOT] NULL *`, `COPY ... (opt DEFAULT)`,
  `MERGE ... WHEN NOT MATCHED BY SOURCE|TARGET`, `AT LOCAL`,
  `ALTER COLUMN ... SET EXPRESSION AS`, `SET ACCESS METHOD DEFAULT`, and
  `SET STATISTICS DEFAULT`.
- **Spanner-PG grammar restrictions**: `ALTER TYPE ... SET (...)` requires `=`
  in `operator_def_elem`, `ALTER INDEX` rejects generic `alter_table_cmds`,
  `CREATE RULE` rejects `ON UPDATE`, and `CLUSTER (...)` requires a table name.
- **Non-empty inputs**: rejects 0-statement inputs (`""`, `";"`).

#### Spanner-PG modified grammar behaviors (0 `Only Spanner-PG accepts`)

All 7 Spanner-PG parser acceptances beyond PostgreSQL 17 are matched cleanly at their real grammar points in `DialectSpannerPG` (while remaining strictly rejected in `DialectDefault`):

- **Column-level `ON UPDATE b_expr` constraint** (`ColConstraintElem`, 5 rows): accepts `ON UPDATE b_expr` on column/domain constraints whenever the token after `ON UPDATE` is not in `FIRST(key_action)` (`NO`, `RESTRICT`, `CASCADE`, `SET`), mirroring Spanner-PG's `parser.c` lookahead disambiguation against foreign-key `key_update`.
- **PostgreSQL 16 `createdb_opt_name` identifiers** (1 row): treats the 20 keywords added in PostgreSQL 17 (`TARGET`, `SOURCE`, `MERGE_ACTION`, and the 17 PG 17 SQL/JSON query keywords) as ordinary `IDENT` option names in `CREATE`/`ALTER DATABASE`, matching Spanner-PG's PG 16 lexer.
- **Bare `SET` in `alter_table_cmd`** (1 row): accepts a bare `SET` subcommand when terminated by `,`, `;`, or EOF in `ALTER TABLE`/`VIEW`/`MATERIALIZED VIEW`/`FOREIGN TABLE` (matching Spanner-PG's modified `alter_table_cmd` rule).

#### `suite:spanner-pg-ddl` (143 rows: 91 accept, 52 reject, 100.00% agreement)

Covers Cloud Spanner PostgreSQL DDL extensions supported in `DialectSpannerPG`
(and rejected in `DialectDefault`) together with near-misses and invalid clause
orderings verified live against `emulator:1.5.54`:

- **Table interleaving**: `CREATE TABLE ... [WITH (...)] INTERLEAVE IN [PARENT] parent [ON DELETE {CASCADE | NO ACTION | ...}]`, `ALTER TABLE ... SET INTERLEAVE IN [PARENT] parent [ON DELETE {CASCADE | NO ACTION}]`, `ALTER TABLE ... SET ON DELETE {CASCADE | NO ACTION}`, `ALTER TABLE ... DROP INTERLEAVE`.
- **Row deletion policy (TTL)**: `CREATE TABLE ... TTL INTERVAL ... ON col`, `ALTER TABLE ... ADD|ALTER TTL INTERVAL ... ON col`, `ALTER TABLE ... DROP TTL`.
- **Column attributes**: `HIDDEN` columns and `GENERATED ALWAYS AS (...) VIRTUAL` columns in `CREATE TABLE` and `ALTER TABLE ADD COLUMN`.
- **Change streams**: `CREATE CHANGE STREAM [IF NOT EXISTS] name [FOR ALL | FOR t1, t2(c1)] [WITH (...)]`, `ALTER CHANGE STREAM name SET FOR ... | SET (...) | RESET (...) | DROP FOR ...`, `DROP CHANGE STREAM [IF EXISTS] name`.
- **Views**: `CREATE [OR REPLACE] VIEW name [(cols)] [WITH (...)] SQL SECURITY {INVOKER | DEFINER} AS SELECT ...`.
- **Sequences**: `BIT_REVERSED_POSITIVE`, `SKIP RANGE min max`, `NO SKIP RANGE`, `START COUNTER [WITH] n`, and `RESTART COUNTER [WITH] n` on `CREATE SEQUENCE` and `ALTER SEQUENCE`.
- **Indexes & Search Indexes**: `CREATE [UNIQUE] INDEX ... INTERLEAVE IN parent [WHERE ...]`, `CREATE SEARCH INDEX name ON table (...) [INCLUDE (...)] [PARTITION BY ...] [ORDER BY ...] [INTERLEAVE IN parent] [WHERE ...] [WITH (...)]`, `ALTER SEARCH INDEX name ADD|DROP COLUMN col`, `DROP SEARCH INDEX [IF EXISTS] name`.
- **Locality groups & database options**: `CREATE LOCALITY GROUP [IF NOT EXISTS] name`, `DROP LOCALITY GROUP [IF EXISTS] name`, and `ALTER DATABASE db SET|RESET spanner.*`.

## Recording and verification commands

Offline replay (CI default, zero Docker, zero network):

```bash
go test ./spanner/parser/ -run '^TestPostgresCompatibility$'
```

Live recording against pinned `gcr.io/cloud-spanner-emulator/emulator:1.5.54`:

```bash
OMNI_RECORD=1 go test -v -timeout 30m ./spanner/parser/ -run '^TestPostgresCompatibility$'
```

Read-only live verification against the pinned emulator (re-judges all 14,114
rows, writes nothing, fails on any diff):

```bash
OMNI_RECORD=verify go test -v -timeout 30m ./spanner/parser/ -run '^TestPostgresCompatibility$'
```
