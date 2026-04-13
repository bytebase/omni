# Doris Migration Analysis

Migration target: full coverage of the legacy ANTLR4 `bytebase/parser/doris` grammar by a hand-written recursive-descent parser under `omni/doris`. Bytebase consumption is used to **prioritize** the order of work; it is never used to **scope down** legacy parser features.

Sources analyzed:

- Legacy parser: `/Users/h3n4l/OpenSource/parser/doris` (`DorisParser.g4` 2,258 lines; `DorisLexer.g4` 14 KB, 200+ keywords)
- Downstream consumer: `/Users/h3n4l/OpenSource/bytebase` (9 import sites in `backend/plugin/parser/doris/`, plus 2 cross-registered from `backend/plugin/parser/mysql/`)

---

## Grammar Coverage

### Origin

The Doris grammar is **derived from Apache Spark SQL** (not MySQL). The header reads: "Copied from Apache Spark and modified for Apache Doris." Despite Doris being MySQL-protocol-compatible, the grammar structure diverges significantly from MySQL's parser.

### Entry rules

| Rule | Parses |
|------|--------|
| `multiStatements` | Complete SQL file (batch of statements) |
| `statement` | Single statement wrapper (includes stored procedures) |
| `statementBase` | Any single SQL statement (21 categories) |
| `query` | SELECT wrapper (CTEs, set operations, ordering) |

### DDL -- CREATE

Covers `CREATE [IF NOT EXISTS]` for:

- Storage objects: `TABLE` (EXTERNAL, TEMPORARY), `VIEW` (OR REPLACE), `DATABASE`/`SCHEMA`, `INDEX` (BITMAP, NGRAM_BF, INVERTED, ANN), `MATERIALIZED VIEW` (MTMV, with BUILD_MODE, REFRESH triggers, partition, distribution)
- Doris-specific: `CATALOG` (multi-catalog for external sources like HIVE, ICEBERG, Delta Lake, Paimon), `STORAGE VAULT`, `STORAGE POLICY`, `DICTIONARY`, `STAGE`, `FILE`, `REPOSITORY`
- Functions: `FUNCTION` (user-defined, alias), `PROCEDURE`
- Security: `ROLE`, `USER`, `ROW POLICY` (RESTRICTIVE/PERMISSIVE), `SQL BLOCK RULE`, `ENCRYPTION KEY`
- Workload management: `WORKLOAD GROUP`, `WORKLOAD POLICY`, `RESOURCE`
- Compute: `COMPUTE GROUP`

### DDL -- ALTER

Supports `ALTER` for: `TABLE` (columns, partitions, properties, rollups, distribution, engine, comment), `DATABASE` (rename, quotas, properties), `VIEW`, `CATALOG`, `ROLE`, `REPOSITORY`, `FUNCTION`, `MATERIALIZED VIEW` (rename, refresh method, replace, set properties), `WORKLOAD GROUP`, `WORKLOAD POLICY`, `SYSTEM` (backends, frontends, brokers, observers, followers), `COLOCATE GROUP`, `RESOURCE`, `USER`, `JOB`

### DDL -- DROP

DROP for all object types above with `[IF EXISTS]`. Additional: `DROP TABLE` supports `TEMPORARY` and `FORCE`. `DROP DATABASE` supports `FORCE`.

### DML

- `INSERT` -- INTO/OVERWRITE TABLE with partitions, labels, column lists, hints
- `UPDATE` -- SET assignments, FROM clauses, WHERE conditions
- `DELETE` -- FROM table with partitions, USING relations, WHERE clauses
- `MERGE INTO` -- WHEN MATCHED/NOT MATCHED clauses with conditional predicates
- `TRUNCATE TABLE` -- with partition specification and FORCE
- `COPY INTO` -- load data from stage (Snowflake-style)
- `LOAD` -- bulk load with data descriptors and properties (MERGE/APPEND/DELETE modes)
- `EXPORT` -- table export to files

### Query/SELECT

- `SELECT` -- `WITH` CTEs, `DISTINCT`/`ALL`, complex joins, subqueries, aggregate functions, window functions, `QUALIFY` clause
- Joins: INNER, LEFT/RIGHT/FULL OUTER, CROSS joins with ON/USING criteria
- Set ops: `UNION [ALL]`, `INTERSECT`, `EXCEPT`/`MINUS`
- `ORDER BY`, `LIMIT`/`OFFSET`

### DCL

- `GRANT` -- table privileges, resource privileges (RESOURCE, CLUSTER, COMPUTE GROUP, STAGE, STORAGE VAULT, WORKLOAD GROUP), roles
- `REVOKE` -- table privileges, resource privileges, roles

### Transaction Control

- `BEGIN` (WITH LABEL), `COMMIT`, `ROLLBACK` -- all with WORK/CHAIN/RELEASE options

### Administrative & System

- **ADMIN**: Show replica distribution/status, rebalance disk, diagnose tablet, compact table, check tablets, set/show config, partition version management, decommission backend
- **SYSTEM ALTER**: Add/drop/modify backends, observers, followers, brokers; set load error hubs
- **SHOW**: 40+ variants (DATABASES, TABLES, COLUMNS, PARTITIONS, CREATE TABLE/VIEW/PROCEDURE, VARIABLES, GRANTS, USERS, ROLES, LOAD, ROUTINE LOAD, BACKUP, RESTORE, PLUGINS, RESOURCES, CATALOGS, STORAGE VAULT, BUILD INDEX, DYNAMIC PARTITION, EVENTS, EXPORT, PROCESSLIST, PRIVILEGES, PROC, FUNCTIONS, CHARSET, COLLATION, ANALYZE, STATS, TABLET, etc.)

### Load/Stream Processing

- `ROUTINE LOAD` -- streaming data load from Kafka/S3/HDFS with CREATE, PAUSE, RESUME, STOP, SHOW
- `SYNC` -- synchronization statement

### Materialized Views

- `CREATE MATERIALIZED VIEW` -- with build mode, refresh triggers (CRON, ON SCHEDULE with EVERY/STARTS/ENDS, ON COMMIT), partition, distribution, rollup properties
- `REFRESH MATERIALIZED VIEW` -- COMPLETE, AUTO, or specific partitions
- `ALTER/DROP/PAUSE/RESUME/CANCEL MATERIALIZED VIEW` -- full lifecycle

### Job Scheduling

- `CREATE JOB` -- STREAMING or SCHEDULE (EVERY/AT), DO statement
- `ALTER/DROP/PAUSE/RESUME JOB`, `CANCEL TASK`

### Statistics & Analysis

- `ANALYZE` -- database, table, profile
- `SHOW ANALYZE` -- job status, queued jobs, tasks, stats (column, table, index)

### Constraints

- `ADD/DROP CONSTRAINT` (PRIMARY KEY, UNIQUE, FOREIGN KEY)
- `SHOW CONSTRAINTS`

### Utility Statements

- `DESCRIBE`/`EXPLAIN` -- table schema, function signatures, query execution plans
- `USE` -- database/catalog/cloud cluster selection (`USE db@cluster`)
- `SET` -- variables, properties, transaction modes, charset/collation, password
- `UNSET` -- variables, storage vault
- `KILL` -- connection, query
- `LOCK`/`UNLOCK` -- table locks (READ/WRITE)
- `INSTALL`/`UNINSTALL PLUGIN`
- `BACKUP`/`RESTORE SNAPSHOT` -- repository-based backup/restore
- `WARM UP` -- cluster/compute group caching
- `REFRESH` -- catalog, database, table, dictionary, LDAP
- `CLEAN` -- profiles, labels, query stats
- `CANCEL` -- load, export, warm-up, backup/restore, build index, ALTER operations
- `RECOVER` -- database, table, partition
- `HELP`

### Doris-specific SQL extensions

These are OLAP-oriented features beyond MySQL compatibility:

- **Key types**: `AGGREGATE KEY`, `UNIQUE KEY`, `DUPLICATE KEY` (table organization models)
- **Distribution**: `DISTRIBUTED BY HASH(columns) BUCKETS n` or `DISTRIBUTED BY RANDOM BUCKETS AUTO`
- **Partitioning**: Range or List partition strategies with `PARTITION BY`
- **Rollups**: Pre-aggregated materialized views via `ROLLUP (name columns)`
- **Table engine**: Explicit `ENGINE = identifier` (doris/ollapdb/memory)
- **Advanced indexing**: BITMAP, NGRAM_BF (N-gram Bloom Filter), INVERTED (full-text), ANN (Approximate Nearest Neighbor)
- **Inverted index components**: ANALYZER, TOKENIZER, TOKEN_FILTER, CHAR_FILTER
- **Multi-catalog**: External data source connectivity (HIVE, ICEBERG, Delta Lake, Paimon)
- **Storage tiering**: Hot/warm/cold via storage policies
- **Workload isolation**: Workload groups and policies for resource management
- **Execution hints**: `/*+ hint */` and `[hint]` for join distribution/skew handling
- **Cloud cluster**: `USE database@cluster` syntax
- **Password policy**: PASSWORD_EXPIRE, PASSWORD_HISTORY, PASSWORD_LOCK_TIME, PASSWORD_REUSE
- **Transaction labels**: `BEGIN WITH LABEL`, `INSERT INTO ... WITH LABEL`

### Data types

The grammar supports standard SQL types plus Doris-specific types. The full type system needs to be extracted from the grammar's `dataType` and `primitiveColType` rules (integers, floats, decimals, strings, dates/times, BOOLEAN, HLL, BITMAP, QUANTILE_STATE, ARRAY, MAP, STRUCT, JSON, VARIANT, AGG_STATE, etc.).

---

## Parse API Surface (legacy)

```go
parser.NewDorisLexer(antlr.CharStream) *DorisLexer
parser.NewDorisParser(antlr.TokenStream) *DorisParser

(*DorisParser).MultiStatements() IMultiStatementsContext   // primary entry
```

Go wrapper functions in `bytebase/bytebase/backend/plugin/parser/doris/`:

| Function | Signature | Purpose |
|----------|-----------|---------|
| `ParseDorisSQL` | `(statement string) ([]*base.ANTLRAST, error)` | Parse multiple statements into ANTLR AST objects |
| `parseSingleDorisSQL` | `(statement string, baseLine int) (*base.ANTLRAST, error)` | Parse single statement with line tracking |
| `parseDorisStatements` | `(statement string) ([]base.ParsedStatement, error)` | Parse statements with text + AST + position info |
| `SplitSQL` | `(statement string) ([]base.Statement, error)` | Split into individual statements (semicolon-delimited) |
| `GetQuerySpan` | `(ctx, statement, database, ...) (*base.QuerySpan, error)` | Extract tables/columns referenced in query |
| `Diagnose` | `(statement string) ([]base.Diagnostic, error)` | Syntax diagnostics for editor |
| `validateQuery` | `(statement string) (bool, bool, error)` | Validate read-only queries for SQL editor |

Errors flow through ANTLR's `ParseErrorListener` with line/column info. Registration model uses `base.Register*Func()` keyed by both `storepb.Engine_DORIS` and `storepb.Engine_STARROCKS` (shared code path).

---

## AST Surface

The legacy parser exposes:

- Listener interface `DorisParserListener` + `BaseDorisParserListener` (Enter/Exit per rule)
- Visitor interface `DorisParserVisitor` + `BaseDorisParserVisitor`
- One Context type per grammar rule (~713 context types)

**Core hierarchy:**

```
MultiStatementsContext
  StatementContext
    StatementBaseContext (21 alias categories)
      SupportedDmlStatementAliasContext
      SupportedCreateStatementAliasContext
      SupportedAlterStatementAliasContext
      SupportedDropStatementAliasContext
      SupportedShowStatementAliasContext
      MaterializedViewStatementAliasContext
      SupportedJobStatementAliasContext
      ConstraintStatementAliasContext
      SupportedLoadStatementAliasContext
      SupportedDescribeStatementAliasContext
      SupportedAdminStatementAliasContext
      SupportedTransactionStatementAliasContext
      SupportedKillStatementAliasContext
      SupportedSetStatementAliasContext
      SupportedUnsetStatementAliasContext
      SupportedRefreshStatementAliasContext
      SupportedCancelStatementAliasContext
      SupportedRecoverStatementAliasContext
      SupportedCleanStatementAliasContext
      SupportedAnalyzeStatementAliasContext
      UnsupportedStatementContext
    CallProcedureContext
    CreateProcedureContext
    DropProcedureContext
```

**Query hierarchy:**

```
QueryContext
  CteContext (Common Table Expressions)
  QueryTermContext
    QueryPrimaryContext
      QuerySpecificationContext (SELECT...FROM...WHERE...GROUP BY...HAVING...QUALIFY)
      SubqueryContext
      InlineTableContext (VALUES clause)
    SetOperationContext (UNION, EXCEPT, INTERSECT)
  QueryOrganizationContext (ORDER BY, LIMIT, OFFSET)
```

**Expression types:** `BooleanExpressionContext`, `PredicateContext`, `ComparisonContext`, `ArithmeticBinaryContext`, `ArithmeticUnaryContext`, `CastExpressionContext`, `FunctionCallExpressionContext`, `CaseExpressionContext`, `WindowFunctionContext`

**Relation types:** `RelationContext`, `RelationPrimaryContext`, `JoinRelationContext`, `JoinCriteriaContext`

**Custom listeners in bytebase:**

| Listener | File | Role |
|----------|------|------|
| `queryValidateListener` | `query.go` | Validate read-only queries for editor |
| `queryTypeListener` | `query_type.go` | Classify statement type (SELECT, DML, DDL, SelectInfoSchema) |
| `accessTableListener` | `query_span_extractor.go` | Collect every table reference |
| CTE listener | `query_span_extractor.go` | Track CTE definitions for exclusion |

---

## Bytebase Import Sites

Import path: `github.com/bytebase/parser/doris`. 9 files in `backend/plugin/parser/doris/`, plus 2 cross-registrations from `backend/plugin/parser/mysql/`.

| File | Parser entry | Purpose |
|------|--------------|---------|
| `backend/plugin/parser/doris/doris.go` | `NewDorisLexer`, `NewDorisParser`, `MultiStatements` | Common parse helper |
| `backend/plugin/parser/doris/split.go` | Lexer + `DorisParserSEMICOLON` | Statement splitting |
| `backend/plugin/parser/doris/statement_ranges.go` | Lexer + token constants | Byte ranges of statements (UTF-16) |
| `backend/plugin/parser/doris/diagnose.go` | Parser + error listener | Syntax diagnostics for editor |
| `backend/plugin/parser/doris/query.go` | Listener walk | Query validation (read-only detection) |
| `backend/plugin/parser/doris/query_type.go` | Listener walk | DDL/DML/SELECT/SelectInfoSchema classification |
| `backend/plugin/parser/doris/query_span.go` | Wrapper | Query span entry point |
| `backend/plugin/parser/doris/query_span_extractor.go` | Listener walk | Table/column access extraction |
| `backend/plugin/parser/doris/normalize.go` | — | Identifier normalization |
| `backend/plugin/parser/mysql/completion.go` | **MySQL omni parser** | Auto-completion (reuses MySQL omni) |
| `backend/plugin/parser/mysql/resource_change.go` | **MySQL omni parser** | Resource change extraction (reuses MySQL omni) |

### Registered functions (8)

| Registration | Key | Handler |
|--------------|-----|---------|
| `RegisterParseStatementsFunc` | `Engine_DORIS`, `Engine_STARROCKS` | `parseDorisStatements` |
| `RegisterSplitterFunc` | `Engine_DORIS`, `Engine_STARROCKS` | `SplitSQL` |
| `RegisterGetQuerySpan` | `Engine_DORIS`, `Engine_STARROCKS` | `GetQuerySpan` |
| `RegisterQueryValidator` | `Engine_DORIS`, `Engine_STARROCKS` | `validateQuery` |
| `RegisterDiagnoseFunc` | `Engine_DORIS`, `Engine_STARROCKS` | `Diagnose` |
| `RegisterStatementRangesFunc` | `Engine_DORIS`, `Engine_STARROCKS` | `GetStatementRanges` |
| `RegisterCompleteFunc` | `Engine_DORIS` | MySQL `Completion` (cross-engine reuse) |
| `RegisterExtractChangedResourcesFunc` | `Engine_DORIS` | MySQL `extractChangedResources` (cross-engine reuse) |

### Engine-specific handling in bytebase

Three Doris-specific code branches exist outside the parser:

1. **Character set clearing** (`plugin/db/starrocks/sync.go`): Doris returns corrupted Unicode in charset/collation fields; bytebase clears them.
2. **Materialized view sync** (`plugin/db/starrocks/sync.go`): Uses Doris-specific `mv_infos()` table-valued function.
3. **Database creation SQL** (`api/v1/rollout_service_task.go`): Simple `CREATE DATABASE` without backticks.

### Feature support matrix in bytebase

| Feature | Doris supported? |
|---------|-----------------|
| Auto-complete | Yes (via MySQL omni parser) |
| Schema sync | Yes |
| Resource change extraction | Yes (via MySQL omni parser) |
| Query validation | Yes |
| Query span / field lineage | Yes |
| Statement splitting | Yes |
| Syntax diagnostics | Yes |
| Data masking | **No** |
| SQL review / lint rules | **No** |
| Query ACLs | **No** |
| Prior backup | **No** |

---

## Feature Dependency Map

| Bytebase feature | Parser APIs used | Data extracted | Critical? |
|------------------|-----------------|----------------|-----------|
| Statement parsing & splitting | Lexer + `MultiStatements`, `DorisParserSEMICOLON`, error listeners | Token boundaries, parse tree, syntax errors | **CRITICAL** -- foundation |
| Statement byte ranges | Lexer + `EOF`/`SEMICOLON` constants | Statement byte positions (UTF-16) | CRITICAL |
| Query type classification | `StatementBaseContext` walk -> DDL/DML/SELECT/SelectInfoSchema | Statement category | HIGH |
| Query validation (editor) | `StatementContext`, `DmlStatementContext` | Read-only/write classification | MEDIUM |
| Query span / field lineage | Walk over `QueryContext`, `QuerySpecificationContext`, `TableNameContext`, `CteContext`, `RelationContext`, etc. | Tables touched, result-column lineage, CTE tracking | **CRITICAL** -- drives field-level access |
| Syntax diagnostics | Parser + error listener | Errors with positions | MEDIUM |
| Auto-completion | **MySQL omni parser** (already migrated) | Completion candidates | N/A (already on omni) |
| Resource change extraction | **MySQL omni parser** (already migrated) | CREATE/DROP/ALTER/DML resource list | N/A (already on omni) |

**Key insight**: Auto-completion and resource change extraction already use the MySQL omni parser for Doris. These do **not** need to be reimplemented for Doris -- they are already migrated.

---

## Gaps and Limitations of the Legacy Parser

1. **TODO: support analyze external table partition spec** (grammar line 990) -- external table analysis lacks partition awareness.
2. **TODO: support other readonly statements** (`query.go:17`) -- query validator doesn't recognize all read-only statements like SHOW TABLES.
3. **FIXME: like should be wildWhere?** (grammar) -- SHOW FRONTEND/BACKEND CONFIG pattern matching.
4. **FIXME: FRONTEND should not contain FROM backendid** (grammar) -- ADMIN SHOW FRONTEND CONFIG syntax.
5. **TODO: need to stay consistent with the legacy** (grammar line 1905) -- legacy compatibility edge case.
6. Query readonly validation incomplete (only covers SELECT, SHOW, EXPLAIN cases).
7. Generated Go files are massive: `doris_parser.go` ~3.4 MB -- clear motivation for rewriting.

---

## Relationship to MySQL

**Divergent from MySQL.** Despite Doris being MySQL-protocol-compatible:

- Grammar is **Spark-derived**, not MySQL-derived
- DorisParser.g4 (2,258 lines) is <50% size of MySQLParser.g4 (5,087 lines)
- Expression syntax follows Spark patterns, not MySQL-specific patterns
- Doris adds OLAP extensions (AGGREGATE KEY, DISTRIBUTED BY, ROLLUP) absent from MySQL

**Overlap areas** (functional MySQL compatibility):
- Basic SELECT/INSERT/UPDATE/DELETE/CREATE TABLE/ALTER TABLE/DROP TABLE
- View/Database management
- User/Role/Grant/Revoke access control
- Transaction control (BEGIN/COMMIT/ROLLBACK)
- Variable/session management (SET/SHOW VARIABLES)

**StarRocks relationship**: Doris and StarRocks share the same parser registration and implementation in bytebase. Both engines are treated identically at the parser level. The Doris omni parser will also serve StarRocks.

---

## Priority Ranking

### Tier 0 -- Foundation (blocks everything)

1. **Lexer** -- 200+ keywords, identifiers (backtick-quoted), numeric/string/binary literals, operators, hints (`/*+ ... */`, `[...]`), semicolons, comments. Must preserve position metadata for `statement_ranges.go`.
2. **AST core types** -- node interfaces, position tracking, walker infrastructure.
3. **Statement splitter** -- uses lexer only; unblocks bytebase's `split.go` and `statement_ranges.go`.
4. **Parser entry** -- `ParseDoris(sql) -> File`; collect syntax errors with positions.

### Tier 1 -- Core SQL surface (unblocks query span + classification)

5. **Identifiers and qualified names** -- `MultipartIdentifier` (1/2/3-part), normalization helpers.
6. **Data types** -- full type grammar (INT, FLOAT, DECIMAL, VARCHAR, CHAR, DATE, DATETIME, BOOLEAN, HLL, BITMAP, QUANTILE_STATE, ARRAY, MAP, STRUCT, JSON, VARIANT, AGG_STATE, etc.).
7. **Expressions** -- operators, function calls, CASE, CAST, subqueries, IN/BETWEEN/LIKE/REGEXP, window functions.
8. **SELECT core** -- SELECT list, FROM (table refs, joins), WHERE, GROUP BY, HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET.
9. **CTEs (`WITH`) and set operators** (UNION/EXCEPT/INTERSECT).
10. **Statement classification helper** -- equivalent of bytebase's `query_type.go`.

### Tier 2 -- DDL backbone (CREATE TABLE is central to Doris)

11. **CREATE TABLE** (full) -- column declarations, constraints (PK/UNIQUE/FK), Doris key types (AGGREGATE KEY, UNIQUE KEY, DUPLICATE KEY), DISTRIBUTED BY, PARTITION BY (range/list), ROLLUP, ENGINE, properties.
12. **ALTER TABLE** (full action set).
13. **CREATE/ALTER/DROP DATABASE / SCHEMA / VIEW**.
14. **CREATE/ALTER/DROP INDEX** (BITMAP, NGRAM_BF, INVERTED, ANN).
15. **DROP** for all above.

### Tier 3 -- Query span extractor parity

16. **Query span extractor** -- reimplementation of the bytebase walker on the new AST. Needs joins, CTE tracking, set-op column merging, subquery field resolution, `TableName` collection.
17. **Syntax diagnostics** -- error listing for editor.

### Tier 4 -- Doris-specific DDL surface

18. **Materialized Views (MTMV)** -- CREATE (with build mode, refresh triggers, partition, distribution), ALTER, DROP, REFRESH, PAUSE, RESUME, CANCEL.
19. **Catalog** -- CREATE/ALTER/DROP CATALOG for external data sources.
20. **Storage & Infrastructure** -- STORAGE VAULT, STORAGE POLICY, REPOSITORY, STAGE, FILE.
21. **Workload management** -- WORKLOAD GROUP, WORKLOAD POLICY, RESOURCE, COMPUTE GROUP, SQL BLOCK RULE.
22. **Security** -- ROW POLICY, ENCRYPTION KEY, DICTIONARY.

### Tier 5 -- DML beyond SELECT

23. **INSERT** (INTO/OVERWRITE, partitions, labels, hints).
24. **UPDATE / DELETE** (FROM clauses, partitions, USING).
25. **MERGE INTO** (MATCHED/NOT MATCHED clauses).
26. **TRUNCATE TABLE**, **COPY INTO**, **LOAD** (bulk), **EXPORT**.
27. **ROUTINE LOAD** -- CREATE, PAUSE, RESUME, STOP, SHOW.

### Tier 6 -- Administrative and utility commands

28. **GRANT / REVOKE** (table, resource, role privileges).
29. **Transaction** (BEGIN WITH LABEL, COMMIT, ROLLBACK).
30. **SHOW** (40+ variants).
31. **DESCRIBE / EXPLAIN**.
32. **USE** (database/catalog/cloud cluster).
33. **SET / UNSET** (variables, properties, transaction modes).
34. **ADMIN** commands (replica, tablet, backend diagnostics).
35. **SYSTEM ALTER** (add/drop/modify backends, frontends, brokers).
36. **BACKUP / RESTORE SNAPSHOT**.
37. **KILL**, **LOCK/UNLOCK**, **INSTALL/UNINSTALL PLUGIN**.
38. **WARM UP**, **REFRESH**, **CLEAN**, **CANCEL**, **RECOVER**, **HELP**.

### Tier 7 -- Job scheduling and procedures

39. **Job scheduling** -- CREATE/ALTER/DROP/PAUSE/RESUME JOB, CANCEL TASK.
40. **Stored procedures** -- CREATE/CALL/DROP PROCEDURE.
41. **Constraint management** -- ADD/DROP CONSTRAINT, SHOW CONSTRAINTS.
42. **SYNC statement**.
43. **ANALYZE / SHOW ANALYZE** (statistics).

---

## Full Coverage Target

Every legacy grammar feature is in scope. Each is tagged P0 (consumed by bytebase today) or P1 (legacy-parity only). P0 must exist before bytebase imports can switch.

| Area | Item | Tier | Tag |
|------|------|------|-----|
| Lexer | Tokens, comments, identifiers, literals, operators, hints, semicolons | 0 | P0 |
| Parser entry | `ParseDoris`, `MultiStatements` | 0 | P0 |
| Statement splitter | Token-driven splitter | 0 | P0 |
| Statement byte ranges | Position metadata | 0 | P0 |
| Identifiers | `MultipartIdentifier`, normalization helpers | 1 | P0 |
| Data types | All Doris types incl. HLL/BITMAP/QUANTILE_STATE/ARRAY/MAP/STRUCT/JSON/VARIANT/AGG_STATE | 1 | P0 |
| Expressions | Arith/cmp/logical, CASE, CAST, function calls, window functions, subqueries | 1 | P0 |
| SELECT core | SELECT list, FROM, WHERE, GROUP BY, HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET | 1 | P0 |
| Joins | INNER/LEFT/RIGHT/FULL OUTER/CROSS with ON/USING | 1 | P0 |
| CTEs | `WITH` common table expressions | 1 | P0 |
| Set ops | UNION [ALL], EXCEPT/MINUS, INTERSECT | 1 | P0 |
| Statement classification | DDL/DML/SELECT/SelectInfoSchema helper | 1 | P0 |
| CREATE TABLE | Full grammar incl. key types, DISTRIBUTED BY, PARTITION BY, ROLLUP, ENGINE, constraints | 2 | P0 |
| ALTER TABLE | Full action set | 2 | P0 |
| CREATE/ALTER/DROP DATABASE/SCHEMA/VIEW | -- | 2 | P0 |
| CREATE/ALTER/DROP INDEX | BITMAP, NGRAM_BF, INVERTED, ANN types | 2 | P0 |
| DROP (all storage objects) | -- | 2 | P0 |
| Query span extractor | Table/column access, CTE tracking, set-op merging | 3 | P0 |
| Syntax diagnostics | Error listing for editor | 3 | P0 |
| INSERT (INTO/OVERWRITE) | Labels, partitions, hints | 5 | P0 |
| UPDATE / DELETE | FROM clauses, partitions | 5 | P0 |
| MERGE INTO | All match clauses | 5 | P0 |
| Materialized Views (MTMV) | CREATE/ALTER/DROP/REFRESH/PAUSE/RESUME/CANCEL | 4 | P1 |
| Catalog DDL | CREATE/ALTER/DROP CATALOG | 4 | P1 |
| Storage infrastructure | STORAGE VAULT/POLICY, REPOSITORY, STAGE, FILE | 4 | P1 |
| Workload management | WORKLOAD GROUP/POLICY, RESOURCE, COMPUTE GROUP, SQL BLOCK RULE | 4 | P1 |
| Security DDL | ROW POLICY, ENCRYPTION KEY, DICTIONARY | 4 | P1 |
| TRUNCATE TABLE, COPY INTO, LOAD, EXPORT | -- | 5 | P1 |
| ROUTINE LOAD | CREATE/PAUSE/RESUME/STOP/SHOW | 5 | P1 |
| GRANT / REVOKE | All privilege types | 6 | P1 |
| Transaction | BEGIN WITH LABEL, COMMIT, ROLLBACK | 6 | P1 |
| SHOW (40+ variants) | -- | 6 | P1 |
| DESCRIBE / EXPLAIN | -- | 6 | P1 |
| USE / SET / UNSET | Database/catalog/cloud cluster, variables, properties | 6 | P1 |
| ADMIN commands | Replica, tablet, backend diagnostics | 6 | P1 |
| SYSTEM ALTER | Add/drop/modify backends, frontends, brokers | 6 | P1 |
| BACKUP / RESTORE SNAPSHOT | -- | 6 | P1 |
| KILL, LOCK/UNLOCK, PLUGIN management | -- | 6 | P1 |
| WARM UP, REFRESH, CLEAN, CANCEL, RECOVER, HELP | -- | 6 | P1 |
| Job scheduling | CREATE/ALTER/DROP/PAUSE/RESUME JOB, CANCEL TASK | 7 | P1 |
| Stored procedures | CREATE/CALL/DROP PROCEDURE | 7 | P1 |
| Constraint management | ADD/DROP CONSTRAINT, SHOW CONSTRAINTS | 7 | P1 |
| SYNC statement | -- | 7 | P1 |
| ANALYZE / SHOW ANALYZE | Statistics | 7 | P1 |

P0 = required before bytebase can switch its imports. P1 = required for full parity with the legacy grammar (the migration target).

---

## Architecture Decisions (preliminary)

### Scope: full legacy parity

The omni Doris parser must cover the entire grammar of `bytebase/parser/doris`. Nothing from the legacy grammar is dropped.

### Relationship to MySQL omni parser

Two bytebase features (auto-completion and resource change extraction) already route Doris through the MySQL omni parser. This means the Doris omni parser does **not** need its own `completion/` or resource-change logic -- these are already handled. The Doris parser must focus on the features that currently use the legacy ANTLR parser: statement parsing, splitting, query span, query validation, query type classification, diagnostics, and statement ranges.

### StarRocks sharing

The legacy parser serves both Doris and StarRocks. The omni Doris parser should also serve StarRocks. Any StarRocks-only divergences (if any) should be handled via engine-type parameters, not separate parsers.

### Reference engine

`mysql/` is the natural reference given Doris's MySQL compatibility, but the grammar is actually Spark-derived. The parser structure should follow `pg/` (most mature AST/parser split), while the DDL-heavy statement categories align with `mysql/` patterns. The Spark-derived expression grammar means Doris expressions are structurally closer to Snowflake than MySQL.

### Target package layout

```
doris/
  ast/                # AST node types
  parser/             # Lexer + recursive descent parser + statement splitter
    testdata/
      legacy/         # Lifted from bytebase/parser/doris/examples (51 SQL files)
      official/       # From Doris documentation
  analysis/           # Query span / field lineage
```

No `completion/`, `deparse/`, `quality/`, or `semantic/` packages needed initially -- auto-completion and resource change already go through MySQL omni, and SQL review/data masking are not currently enabled for Doris.

---

## Test Corpus Strategy (preliminary)

### Corpus A -- Legacy ANTLR4 examples (regression baseline)

Lift every `.sql` file from `/Users/h3n4l/OpenSource/parser/doris/examples/` (51 SQL files covering: account, catalog, cluster, data, DDL, DML, function, job, materialized view, security, session, show, statistics, transaction).

### Corpus B -- Official Apache Doris SQL reference

Extract SQL examples from `https://doris.apache.org/docs/sql-manual/` and its subpages. Cross-reference with the legacy grammar when implementing each rule.
