# Snowflake Migration DAG

Derived from `analysis.md`. Each node is a discrete chunk of work that can be sent to `omni-engine-implementing` (which then delegates to `superpowers:brainstorming` for spec → plan → implement). Nodes that share the same dependency level can run in parallel.

**Architecture decisions** (locked in `analysis.md`):
- Reference template: `pg/` (structure) + `mongo/analysis` (query span shape)
- LIMIT injection: AST-level rewrite via a new `snowflake/deparse/` package
- Test corpora: `snowflake/parser/testdata/{legacy,official}/`
- Scope: full legacy parity (Tier 7 included)

**Package layout** (target):
```
snowflake/
  ast/                  # Node interfaces, position types, walker
  parser/               # Lexer + recursive-descent parser + statement splitter
    testdata/{legacy,official}/
  parsertest/           # Table-driven tests + corpus runner
  analysis/             # Query span / field lineage (mongo/analysis-style)
  advisor/              # 14 lint rules + dispatcher
  deparse/              # AST → SQL string + LIMIT injection rewrite
```

---

## Nodes

| ID | Node | Package | Depends On | Parallel With | Tier | Priority | Status |
|----|------|---------|------------|---------------|------|----------|--------|
| **F1** | ast-core | `snowflake/ast` | — | F2, C1, C2 | 0 | P0 | **done** (PR #14) |
| **F2** | lexer | `snowflake/parser` | F1 | C1, C2 | 0 | P0 | **done** (PR #17) |
| **F3** | statement-splitter | `snowflake/parser` | F2 | F4 | 0 | P0 | **done** (PR #18) |
| **F4** | parser-entry + walker | `snowflake/parser` | F1, F2 | F3 | 0 | P0 | **done** (PR #19) |
| **C1** | corpus-legacy (lift 27 SQL files) | `snowflake/parser/testdata/legacy` | — | F1–F4, C2 | 0 | P0 | **done** |
| **C2** | corpus-official (scrape docs.snowflake.com) | `snowflake/parser/testdata/official` | — | F1–F4, C1 | 0 | P0 | **done** (630 SQL files across 78 commands) |
| **T1.1** | identifiers + qualified names + normalization helpers | `snowflake/parser` | F4 | — | 1 | P0 | **done** (PR #22) |
| **T1.2** | data types (incl. VARIANT/OBJECT/ARRAY/VECTOR/GEO/TIMESTAMP_*) | `snowflake/parser` | T1.1 | — | 1 | P0 | **done** (PR #39) |
| **T1.3** | expressions (operators, function calls, CASE, CAST, JSON path, lambdas, subqueries, IN/BETWEEN/LIKE/RLIKE) | `snowflake/parser` | T1.2 | — | 1 | P0 | **done** (PR #46) |
| **T1.4** | SELECT core (list, FROM, WHERE, GROUP BY incl. CUBE/ROLLUP/GROUPING SETS/ALL, HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET/FETCH/TOP, EXCLUDE) | `snowflake/parser` | T1.3 | — | 1 | P0 | **done** (PR #48, includes T1.6 CTEs) |
| **T1.5** | joins (INNER/LEFT/RIGHT/FULL/CROSS/NATURAL/ASOF/DIRECTED/LATERAL) | `snowflake/parser` | T1.4 | T1.6, T1.7 | 1 | P0 | **done** (PR #50) |
| **T1.6** | CTEs (WITH [RECURSIVE]) | `snowflake/parser` | T1.4 | T1.5, T1.7 | 1 | P0 | **done** (included in T1.4 PR #48) |
| **T1.7** | set operators (UNION/UNION BY NAME/EXCEPT/MINUS/INTERSECT) | `snowflake/parser` | T1.4 | T1.5, T1.6 | 1 | P0 | **done** (PR #53) |
| **T1.8** | statement-classification helper (DDL/DML/SELECT/SHOW/DESCRIBE/Other) | `snowflake/analysis` | F4 | T1.1–T1.7 | 1 | P0 | not started |
| **T2.1** | DDL: DATABASE / SCHEMA (CREATE/ALTER/DROP/UNDROP) | `snowflake/parser` | T1.1 | T2.2, T2.4, T2.5 | 2 | P0 | not started |
| **T2.2** | DDL: CREATE TABLE (full — constraints, CTAS, LIKE, CLUSTER BY, COPY GRANTS, WITH TAGS, MASKING POLICY, COLLATE, IDENTITY/AUTOINCREMENT) | `snowflake/parser` | T1.2, T1.3, T1.4 | T2.1, T2.4 | 2 | P0 | **done** (PR #60) |
| **T2.3** | DDL: ALTER TABLE (full action set) | `snowflake/parser` | T2.2 | T2.4, T2.5 | 2 | P0 | not started |
| **T2.4** | DDL: VIEW + MATERIALIZED VIEW (CREATE/ALTER/DROP) | `snowflake/parser` | T1.4, T1.6 | T2.1, T2.2, T2.3 | 2 | P0 | not started |
| **T2.5** | DDL: DROP / UNDROP (table/schema/db core) | `snowflake/parser` | T1.1 | T2.1–T2.4 | 2 | P0 | not started |
| **T2.6** | advisor dispatcher (generic walker, rule registration) | `snowflake/advisor` | F4 | T2.1–T2.5 | 2 | P0 | not started |
| **T2.7** | 14 lint rules (parallel within node — see breakdown below) | `snowflake/advisor` | T2.2, T2.3, T2.5, T2.6, T1.4 | T3.1, T3.2, T3.4 | 2 | P0 | not started |
| **T3.1** | query span extractor (result-column lineage, table-access, CTE/set-op merging, EXCLUDE, subquery field resolution) | `snowflake/analysis` | T1.4, T1.5, T1.6, T1.7 | T2.7, T3.2, T3.4 | 3 | P0 | not started |
| **T3.2** | deparse-core (AST → SQL string for all P0 statement nodes) | `snowflake/deparse` | T1.4, T2.2, T5.1 | T2.7, T3.1, T3.4 | 3 | P0 | not started |
| **T3.3** | LIMIT injection rewrite (AST-level, mirrors `mysql/deparse/rewrite.go`) | `snowflake/deparse` | T3.2 | T3.4 | 3 | P0 | not started |
| **T3.4** | syntax diagnostics (error listing for editor) | `snowflake/parser` | F4 | T2.7, T3.1–T3.3 | 3 | P0 | not started |
| **T5.1** | DML: INSERT (single + INSERT ALL/FIRST), UPDATE (USING), DELETE (USING), MERGE (matched / not matched / not matched by source) | `snowflake/parser` | T1.4, T1.5 | T2.* | 5 | P0 | not started |
| **M1** | bytebase import switch (replace `bytebase/parser/snowflake` imports with `omni/snowflake`) | bytebase repo | F1–T5.1 (all P0) | — | — | P0 | not started |
| **T4.1** | DDL: STAGE | `snowflake/parser` | T2.2 | T4.2–T4.9 | 4 | P1 | not started |
| **T4.2** | DDL: FILE FORMAT (CSV/JSON/AVRO/ORC/PARQUET) | `snowflake/parser` | T2.2 | T4.1, T4.3–T4.9 | 4 | P1 | not started |
| **T4.3** | DDL: PIPE / STREAM / TASK / ALERT | `snowflake/parser` | T2.2, T1.4 | T4.1, T4.2, T4.4–T4.9 | 4 | P1 | not started |
| **T4.4** | DDL: DYNAMIC TABLE / EXTERNAL TABLE / EVENT TABLE / SEQUENCE | `snowflake/parser` | T2.2, T1.4 | T4.1–T4.3, T4.5–T4.9 | 4 | P1 | not started |
| **T4.5** | DDL: FUNCTION / PROCEDURE (SQL/JS/Python/Java handlers, external functions) | `snowflake/parser` | T1.3 | T4.1–T4.4, T4.6–T4.9 | 4 | P1 | not started |
| **T4.6** | DDL: security (ROLE / USER / MASKING POLICY / ROW ACCESS / SESSION/PASSWORD/NETWORK POLICY / SECURITY INTEGRATION) | `snowflake/parser` | T1.3 | T4.1–T4.5, T4.7–T4.9 | 4 | P1 | not started |
| **T4.7** | DDL: integrations (STORAGE/API/NOTIFICATION INTEGRATION, RESOURCE MONITOR, SECRET, CONNECTION, GIT REPOSITORY) | `snowflake/parser` | T1.3 | T4.1–T4.6, T4.8, T4.9 | 4 | P1 | not started |
| **T4.8** | DDL: replication / sharing (FAILOVER/REPLICATION GROUP, ACCOUNT, MANAGED ACCOUNT, SHARE) | `snowflake/parser` | T1.1 | T4.1–T4.7, T4.9 | 4 | P1 | not started |
| **T4.9** | DDL: newer surface (DATASET, SEMANTIC VIEW/DIMENSION/METRIC, TAG) | `snowflake/parser` | T1.4 | T4.1–T4.8 | 4 | P1 | not started |
| **T5.2** | DML: COPY INTO (load + unload) / PUT / GET / LIST / REMOVE | `snowflake/parser` | T1.4 | T5.3, T5.4 | 5 | P1 | not started |
| **T5.3** | Snowflake-specific clauses: MATCH_RECOGNIZE / PIVOT / UNPIVOT / LATERAL FLATTEN / SPLIT_TO_TABLE / CONNECT BY / START WITH / CHANGES / AT / BEFORE / SAMPLE / TABLESAMPLE / CLONE | `snowflake/parser` | T1.4, T1.5 | T5.2, T5.4 | 5 | P1 | not started |
| **T5.4** | CALL / EXECUTE IMMEDIATE / EXECUTE TASK / EXPLAIN | `snowflake/parser` | T1.3 | T5.2, T5.3 | 5 | P1 | not started |
| **T6.1** | DCL: GRANT / REVOKE (roles + privileges + FROM SHARE + WITH GRANT OPTION) | `snowflake/parser` | T1.1 | T6.2, T6.3 | 6 | P1 | not started |
| **T6.2** | TCL: BEGIN / START TRANSACTION / COMMIT / ROLLBACK / SAVEPOINT | `snowflake/parser` | F4 | T6.1, T6.3 | 6 | P1 | not started |
| **T6.3** | Utility: SHOW (50+) / DESCRIBE / USE / SET / UNSET / COMMENT ON / TRUNCATE / ABORT* | `snowflake/parser` | T1.1 | T6.1, T6.2 | 6 | P1 | not started |
| **T7.1** | Snowflake Scripting bodies (DECLARE, BEGIN…END, `:=`, IF, CASE, FOR, WHILE, REPEAT, EXCEPTION, RETURN, cursors / resultsets) — structural parse only | `snowflake/parser` | T1.3, T4.3, T4.5 | — | 7 | P1 | not started |
| **M2** | Final corpus closure (every file in `testdata/legacy` and `testdata/official` parses cleanly) | `snowflake/parsertest` | T7.1 (last P1) | — | — | P1 | not started |

### T2.7 Sub-nodes (14 lint rules — all parallel once T2.6 lands)

| Sub-ID | Rule |
|--------|------|
| T2.7.a | column_no_null |
| T2.7.b | column_require |
| T2.7.c | column_maximum_varchar_length |
| T2.7.d | table_require_pk |
| T2.7.e | table_no_foreign_key |
| T2.7.f | naming_table |
| T2.7.g | naming_table_no_keyword |
| T2.7.h | naming_identifier_case |
| T2.7.i | naming_identifier_no_keyword |
| T2.7.j | select_no_select_all |
| T2.7.k | where_require_select |
| T2.7.l | where_require_update_delete |
| T2.7.m | migration_compatibility |
| T2.7.n | table_drop_naming_convention |

---

## Execution Order

### Wave 0 — Foundation (3 parallel tracks)

```
F1 (ast-core) ──► F2 (lexer) ──► F3 (splitter)
                             └──► F4 (parser-entry)
C1 (corpus-legacy)            ┐
C2 (corpus-official)          ┴── independent, can start day one
```

### Wave 1 — Core SQL (linear chain + 1 sibling)

```
F4 ──► T1.1 (identifiers) ──► T1.2 (types) ──► T1.3 (expressions) ──► T1.4 (SELECT core)
                                                                        ├──► T1.5 (joins)
                                                                        ├──► T1.6 (CTEs)
                                                                        └──► T1.7 (set ops)
F4 ──► T1.8 (statement-classification helper)            // independent of T1.1–T1.7
```

### Wave 2 — DDL backbone + lint dispatcher (mostly parallel after T1.4)

```
T1.1 ──► T2.1 (DB/SCHEMA)
T1.2,T1.3,T1.4 ──► T2.2 (CREATE TABLE) ──► T2.3 (ALTER TABLE)
T1.4,T1.6      ──► T2.4 (VIEW/MV)
T1.1           ──► T2.5 (DROP/UNDROP)
F4             ──► T2.6 (advisor dispatcher)
T2.2,T2.3,T2.5,T2.6,T1.4 ──► T2.7 (14 lint rules — all parallel)
```

### Wave 3 — Query span + diagnostics + LIMIT injection

```
T1.4,T1.5,T1.6,T1.7 ──► T3.1 (query span extractor)
T1.4,T2.2,T5.1      ──► T3.2 (deparse-core) ──► T3.3 (LIMIT injection)
F4                  ──► T3.4 (syntax diagnostics)
```

### Wave 4 — P0 DML

```
T1.4,T1.5 ──► T5.1 (INSERT/UPDATE/DELETE/MERGE)
```

### M1 — Bytebase import switch (P0 done)

After Waves 0–4 complete, switch `backend/plugin/parser/snowflake`, `backend/plugin/db/snowflake`, and `backend/plugin/advisor/snowflake` from `bytebase/parser/snowflake` to `omni/snowflake`.

### Waves 5–7 — P1 (legacy parity, post-bytebase-switch)

All P1 nodes can run in parallel within their tier — they only depend on Tier 1/2 foundations:

```
T4.1, T4.2, T4.3, T4.4, T4.5, T4.6, T4.7, T4.8, T4.9   // remaining DDL (all parallel)
T5.2, T5.3, T5.4                                       // remaining DML / Snowflake specifics (all parallel)
T6.1, T6.2, T6.3                                       // GRANT/TCL/utility (all parallel)
T7.1                                                   // Snowflake Scripting bodies (depends on T4.3, T4.5)
```

### M2 — Full corpus closure

Every file in `testdata/legacy/` and `testdata/official/` parses without error. Migration complete.

---

## Critical Path to Bytebase Migration (P0 only)

```
F1 → F2 → F4 → T1.1 → T1.2 → T1.3 → T1.4 → T1.5/T1.6/T1.7 → T3.1 → M1
                                          → T2.2 → T2.3 → T2.7 → M1
                                          → T5.1 → M1
                                          → T3.2 → T3.3 → M1
```

Longest dependency chain: **F1 → F2 → F4 → T1.1 → T1.2 → T1.3 → T1.4 → T1.5/6/7 → T3.1 → M1** (10 sequential nodes). Everything else can be parallelized around it.

The dispatcher + 14 lint rules (T2.6 + T2.7.a–n) are the largest parallel-fan-out — once T2.6 lands, all 14 rules can be implemented concurrently.

---

## Notes

1. **F3 (statement splitter) is on the critical path for bytebase's `split.go` and `statement_ranges.go`** but does not block any other parser node — bytebase consumes it directly. Could ship to bytebase as soon as F2 lands, ahead of F4.

2. **C1 + C2 (test corpora) should start in parallel with F1.** Corpus B (official docs scrape) is the longest-running independent task — start it first so the data is ready by the time T1.x nodes need it.

3. **T1.8 (statement classification helper) is unusually placed in `analysis/`** because that's where `mongo/analysis` puts its operation-classifier. Bytebase's `query_type.go` is the consumer.

4. **T3.2 (deparse-core) depends on T5.1** — deparse must cover INSERT/UPDATE/DELETE/MERGE alongside SELECT/DDL because bytebase deparses arbitrary statements. T5.1 therefore needs to land before T3.2 can be considered complete; T3.2 can start partial work (covering Tier 1/2 nodes) earlier.

5. **T2.7 sub-nodes are parallelizable but small** — each lint rule is ~100–250 LOC of porting work. They can be batched into 2–3 brainstorming sessions if running them individually creates too much overhead.

6. **T7.1 (Snowflake Scripting) depends on T4.3 (TASK) and T4.5 (FUNCTION/PROCEDURE)** because the scripting bodies live inside those statements. It is the last P1 node.

7. **No node introduces a `catalog/` or `completion/` package.** Snowflake doesn't need them today (bytebase doesn't consume them) and the analysis explicitly notes the absence. Add them later if a new bytebase feature requires them.

8. **`semantic/` is also intentionally absent.** The legacy parser performs no semantic checks, and bytebase has no semantic-checking call site for Snowflake. Add later if needed.

9. **Brainstorming sessions for each node should start by reading the corresponding `.g4` rule(s) AND the official docs page** for that statement (per the `feedback_parser_testing.md` memory). The two-source rule applies at every level.

---

## Status Tracking

Update the **Status** column in the node table as work progresses:
- `not started`
- `in progress`
- `done`

When a node moves to `done`, run its associated subset of `testdata/legacy/` and `testdata/official/` to confirm green status before marking it complete.
