# Doris Migration DAG

Derived from `analysis.md`. Each node is a discrete chunk of work that can be sent to `omni-engine-implementing` (which then delegates to `superpowers:brainstorming` for spec -> plan -> implement). Nodes that share the same dependency level can run in parallel.

**Architecture decisions** (locked in `analysis.md`):
- Reference template: `pg/` (structure) + `mysql/` (DDL patterns)
- Grammar origin: Spark-derived (not MySQL), affects expression grammar structure
- Auto-completion and resource change: already on MySQL omni (not in scope)
- SQL review / lint rules: not enabled for Doris in bytebase (not in scope)
- StarRocks: shared parser (both engines use the same code path)
- Test corpora: `doris/parser/testdata/{legacy,official}/`
- Scope: full legacy parity

**Package layout** (target):
```
doris/
  ast/                  # Node interfaces, position types, walker
  parser/               # Lexer + recursive-descent parser + statement splitter
    testdata/{legacy,official}/
  analysis/             # Query span / field lineage + statement classification
```

No `completion/`, `deparse/`, `quality/`, `advisor/`, or `semantic/` packages needed -- auto-completion and resource change already go through MySQL omni, and SQL review/data masking are not enabled for Doris.

---

## Nodes

| ID | Node | Package | Depends On | Parallel With | Tier | Priority | Status |
|----|------|---------|------------|---------------|------|----------|--------|
| **F1** | ast-core | `doris/ast` | -- | F2, C1 | 0 | P0 | **done** (PR #47) |
| **F2** | lexer | `doris/parser` | F1 | C1 | 0 | P0 | **done** (PR #49) |
| **F3** | statement-splitter | `doris/parser` | F2 | F4 | 0 | P0 | **done** (PR #51) |
| **F4** | parser-entry + walker | `doris/parser` | F1, F2 | F3 | 0 | P0 | **done** (PR #52) |
| **C1** | corpus-legacy (lift 51 SQL files) | `doris/parser/testdata/legacy` | -- | F1, F2, F3, F4 | 0 | P0 | **done** (PR #47) |
| **T1.1** | identifiers + qualified names + normalization helpers | `doris/parser` | F4 | T1.8 | 1 | P0 | **done** (PR #54) |
| **T1.2** | data types (INT, FLOAT, DECIMAL, VARCHAR, CHAR, DATE, DATETIME, BOOLEAN, HLL, BITMAP, QUANTILE_STATE, ARRAY, MAP, STRUCT, JSON, VARIANT, AGG_STATE, etc.) | `doris/parser` | T1.1 | T1.8 | 1 | P0 | **done** (PR #57) |
| **T1.3** | expressions (operators, function calls, CASE, CAST, subqueries, IN/BETWEEN/LIKE/REGEXP, window functions, hints) | `doris/parser` | T1.2 | T1.8 | 1 | P0 | not started |
| **T1.4** | SELECT core (list, FROM, WHERE, GROUP BY, HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET) | `doris/parser` | T1.3 | T1.8 | 1 | P0 | not started |
| **T1.5** | joins (INNER/LEFT/RIGHT/FULL OUTER/CROSS with ON/USING, execution hints) | `doris/parser` | T1.4 | T1.6, T1.7, T1.8 | 1 | P0 | not started |
| **T1.6** | CTEs (WITH common table expressions) | `doris/parser` | T1.4 | T1.5, T1.7, T1.8 | 1 | P0 | not started |
| **T1.7** | set operators (UNION [ALL], EXCEPT/MINUS, INTERSECT) | `doris/parser` | T1.4 | T1.5, T1.6, T1.8 | 1 | P0 | not started |
| **T1.8** | statement-classification helper (DDL/DML/SELECT/SelectInfoSchema) | `doris/analysis` | F4 | T1.1--T1.7 | 1 | P0 | **done** (PR #55) |
| **T2.1** | DDL: CREATE TABLE (full -- column declarations, constraints, AGGREGATE/UNIQUE/DUPLICATE KEY, DISTRIBUTED BY, PARTITION BY range/list, ROLLUP, ENGINE, properties) | `doris/parser` | T1.2, T1.3 | T2.3, T2.4, T2.5 | 2 | P0 | not started |
| **T2.2** | DDL: ALTER TABLE (columns, partitions, properties, rollups, distribution, engine, comment) | `doris/parser` | T2.1 | T2.3, T2.4, T2.5 | 2 | P0 | not started |
| **T2.3** | DDL: DATABASE / SCHEMA (CREATE/ALTER/DROP incl. FORCE) | `doris/parser` | T1.1 | T2.1, T2.2, T2.4, T2.5 | 2 | P0 | **done** (PR #58) |
| **T2.4** | DDL: VIEW (CREATE OR REPLACE / ALTER / DROP) | `doris/parser` | T1.4 | T2.1, T2.2, T2.3, T2.5 | 2 | P0 | not started |
| **T2.5** | DDL: INDEX (CREATE/ALTER/DROP with BITMAP, NGRAM_BF, INVERTED, ANN types; BUILD INDEX) | `doris/parser` | T1.1 | T2.1--T2.4 | 2 | P0 | **done** (PR #57) |
| **T3.1** | query span extractor (result-column lineage, table-access, CTE tracking, set-op merging, subquery field resolution) | `doris/analysis` | T1.4, T1.5, T1.6, T1.7 | T3.2 | 3 | P0 | not started |
| **T3.2** | syntax diagnostics (error listing for editor) | `doris/parser` | F4 | T3.1 | 3 | P0 | **done** (PR #56) |
| **T4.1** | DML: INSERT (INTO/OVERWRITE, partitions, labels, column lists, hints) | `doris/parser` | T1.4 | T4.2, T4.3 | 4 | P0 | not started |
| **T4.2** | DML: UPDATE / DELETE (SET, FROM, partitions, USING, WHERE) | `doris/parser` | T1.4, T1.5 | T4.1, T4.3 | 4 | P0 | not started |
| **T4.3** | DML: MERGE INTO (WHEN MATCHED / NOT MATCHED, conditional predicates) | `doris/parser` | T1.4, T1.5 | T4.1, T4.2 | 4 | P0 | not started |
| **M1** | bytebase import switch (replace `bytebase/parser/doris` imports with `omni/doris`) | bytebase repo | all P0 nodes | -- | -- | P0 | not started |
| **T5.1** | DDL: Materialized Views MTMV (CREATE with build mode/refresh triggers/partition/distribution, ALTER, DROP, REFRESH, PAUSE, RESUME, CANCEL) | `doris/parser` | T2.1, T1.4 | T5.2--T5.5 | 5 | P1 | not started |
| **T5.2** | DDL: Catalog (CREATE/ALTER/DROP for external sources -- HIVE, ICEBERG, Delta Lake, Paimon) | `doris/parser` | T1.1 | T5.1, T5.3--T5.5 | 5 | P1 | not started |
| **T5.3** | DDL: Storage infrastructure (STORAGE VAULT, STORAGE POLICY, REPOSITORY, STAGE, FILE) | `doris/parser` | T1.1 | T5.1, T5.2, T5.4, T5.5 | 5 | P1 | not started |
| **T5.4** | DDL: Workload management (WORKLOAD GROUP/POLICY, RESOURCE, COMPUTE GROUP, SQL BLOCK RULE) | `doris/parser` | T1.1 | T5.1--T5.3, T5.5 | 5 | P1 | not started |
| **T5.5** | DDL: Security (ROLE, USER, ROW POLICY, ENCRYPTION KEY, DICTIONARY) | `doris/parser` | T1.1 | T5.1--T5.4 | 5 | P1 | not started |
| **T6.1** | DML: TRUNCATE TABLE, COPY INTO, LOAD (bulk), EXPORT | `doris/parser` | T1.4 | T6.2 | 6 | P1 | not started |
| **T6.2** | DML: ROUTINE LOAD (CREATE/PAUSE/RESUME/STOP/SHOW), SYNC | `doris/parser` | T1.1 | T6.1 | 6 | P1 | not started |
| **T7.1** | DCL: GRANT / REVOKE (table, resource, role privileges) | `doris/parser` | T1.1 | T7.2--T7.5 | 7 | P1 | not started |
| **T7.2** | TCL: BEGIN WITH LABEL, COMMIT, ROLLBACK (WORK/CHAIN/RELEASE) | `doris/parser` | F4 | T7.1, T7.3--T7.5 | 7 | P1 | not started |
| **T7.3** | Utility: SHOW (40+ variants), DESCRIBE / EXPLAIN, USE (database/catalog/cloud cluster), SET / UNSET, HELP | `doris/parser` | T1.1, T1.3 | T7.1, T7.2, T7.4, T7.5 | 7 | P1 | not started |
| **T7.4** | Admin: ADMIN commands (replica, tablet, backend diagnostics), SYSTEM ALTER (backends, frontends, brokers) | `doris/parser` | T1.1 | T7.1--T7.3, T7.5 | 7 | P1 | not started |
| **T7.5** | Utility: BACKUP / RESTORE, KILL, LOCK/UNLOCK, PLUGIN, WARM UP, REFRESH, CLEAN, CANCEL, RECOVER | `doris/parser` | T1.1 | T7.1--T7.4 | 7 | P1 | not started |
| **T8.1** | Job scheduling (CREATE/ALTER/DROP/PAUSE/RESUME JOB, CANCEL TASK) | `doris/parser` | T1.3 | T8.2, T8.3 | 8 | P1 | not started |
| **T8.2** | Stored procedures (CREATE/CALL/DROP PROCEDURE) | `doris/parser` | T1.3 | T8.1, T8.3 | 8 | P1 | not started |
| **T8.3** | Constraint management + ANALYZE (ADD/DROP CONSTRAINT, SHOW CONSTRAINTS, ANALYZE, SHOW ANALYZE) | `doris/parser` | T1.1 | T8.1, T8.2 | 8 | P1 | not started |
| **M2** | Full corpus closure (every file in `testdata/legacy` and `testdata/official` parses cleanly) | `doris/parser` | T8.1, T8.2, T8.3 (last P1) | -- | -- | P1 | not started |

---

## Execution Order

### Wave 0 -- Foundation (3 parallel tracks)

```
F1 (ast-core) --> F2 (lexer) --> F3 (splitter)
                             \-> F4 (parser-entry)
C1 (corpus-legacy)              independent, can start day one
```

### Wave 1 -- Core SQL (linear chain + 1 sibling)

```
F4 --> T1.1 (identifiers) --> T1.2 (types) --> T1.3 (expressions) --> T1.4 (SELECT core)
                                                                       |-> T1.5 (joins)
                                                                       |-> T1.6 (CTEs)
                                                                       \-> T1.7 (set ops)
F4 --> T1.8 (statement-classification helper)       // independent of T1.1-T1.7
```

### Wave 2 -- DDL backbone (mostly parallel after T1.4)

```
T1.2,T1.3    --> T2.1 (CREATE TABLE) --> T2.2 (ALTER TABLE)
T1.4         --> T2.4 (VIEW)
T1.1         --> T2.3 (DB/SCHEMA)
T1.1         --> T2.5 (INDEX)
```

### Wave 3 -- Query span + diagnostics

```
T1.4,T1.5,T1.6,T1.7 --> T3.1 (query span extractor)
F4                   --> T3.2 (syntax diagnostics)
```

### Wave 4 -- P0 DML (parallel with Wave 2/3)

```
T1.4       --> T4.1 (INSERT)
T1.4,T1.5  --> T4.2 (UPDATE/DELETE)
T1.4,T1.5  --> T4.3 (MERGE INTO)
```

### M1 -- Bytebase import switch (P0 done)

After Waves 0-4 complete, switch `backend/plugin/parser/doris` from `bytebase/parser/doris` to `omni/doris`. Also update StarRocks registrations.

### Waves 5-8 -- P1 (legacy parity, post-bytebase-switch)

All P1 nodes can run in parallel within their tier:

```
T5.1, T5.2, T5.3, T5.4, T5.5   // Doris-specific DDL (all parallel)
T6.1, T6.2                      // Remaining DML + data loading (all parallel)
T7.1, T7.2, T7.3, T7.4, T7.5   // DCL + TCL + utility + admin (all parallel)
T8.1, T8.2, T8.3                // Job scheduling + procedures + constraints (all parallel)
```

### M2 -- Full corpus closure

Every file in `testdata/legacy/` and `testdata/official/` parses without error. Migration complete.

---

## Critical Path to Bytebase Migration (P0 only)

```
F1 -> F2 -> F4 -> T1.1 -> T1.2 -> T1.3 -> T1.4 -> T1.5/T1.6/T1.7 -> T3.1 -> M1
                                        \-> T2.1 -> T2.2 -> M1
                                 \-> T1.4 -> T4.1 -> M1
                                          -> T4.2 -> M1
                                          -> T4.3 -> M1
```

Longest dependency chain: **F1 -> F2 -> F4 -> T1.1 -> T1.2 -> T1.3 -> T1.4 -> T1.5/6/7 -> T3.1 -> M1** (10 sequential nodes). Everything else can be parallelized around it.

---

## Notes

1. **Doris grammar is Spark-derived, not MySQL-derived.** The expression grammar follows Spark patterns. This means the expression parser (T1.3) should study `DorisParser.g4`'s `booleanExpression`, `valueExpression`, `primaryExpression` rules carefully -- they differ from MySQL's operator precedence.

2. **StarRocks shares the parser.** Both `Engine_DORIS` and `Engine_STARROCKS` register the same functions. The omni parser must work for both. No StarRocks-specific grammar divergences have been identified in the legacy code.

3. **Auto-completion and resource change are already migrated.** These already route through MySQL omni (`mysql/completion.go` and `mysql/resource_change.go`). The Doris omni parser does not need to provide these features.

4. **No lint rules / advisor needed.** SQL review is not enabled for Doris in bytebase (`return false` for StatementAdvise). No `advisor/` package is needed.

5. **No deparse / LIMIT injection needed.** Unlike Snowflake, bytebase does not perform token-stream rewriting or LIMIT injection for Doris. No `deparse/` package is needed.

6. **F3 (statement splitter) can ship early.** It depends only on the lexer (F2) and unblocks bytebase's `split.go` and `statement_ranges.go` independently of the full parser.

7. **T1.8 (statement classification) goes in `analysis/`** following the `mongo/analysis` pattern. Bytebase's `query_type.go` is the consumer.

8. **T2.1 (CREATE TABLE) is the most complex DDL node** due to Doris-specific extensions: AGGREGATE KEY / UNIQUE KEY / DUPLICATE KEY, DISTRIBUTED BY HASH/RANDOM, PARTITION BY range/list, ROLLUP, ENGINE. This is the defining DDL feature for Doris.

9. **C1 (corpus-legacy) should start in Wave 0.** Lifting the 51 SQL files from `bytebase/parser/doris/examples/` is trivial and independent of all other work.

10. **T7.3 (SHOW 40+ variants) is the largest single P1 node** by statement count. It can be implemented incrementally -- each SHOW variant is independent.

11. **Simpler than Snowflake migration.** Doris has fewer bytebase consumer features (no lint, no deparse, no LIMIT injection), so the P0 critical path has fewer nodes. The total P1 surface is also smaller (DorisParser.g4 is 2,258 lines vs SnowflakeParser.g4's 5,718 lines).

---

## Status Tracking

Update the **Status** column in the node table as work progresses:
- `not started`
- `in progress`
- `done`

When a node moves to `done`, run its associated subset of `testdata/legacy/` and `testdata/official/` to confirm green status before marking it complete.
