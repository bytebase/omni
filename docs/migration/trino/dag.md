# Trino Migration DAG

Seeded into the v2 SQLite store (`omni-v2-migration.db`, migration id 3,
`dag_version = v2-2026-06-02`, `oracle_tier = docker`) on 2026-06-02. The store is
the source of truth the orchestrator schedules from; this file is the readable
companion. Built from [`analysis.md`](analysis.md).

**18 nodes** (14 P0 · 4 P1), **30 deps**, **7 divergences**. Single root
(`ast-core`); 9 merge layers; acyclic. P0 = required for the bytebase import
switch; P1 = legacy-parity / advanced, deferrable but in scope.

## Nodes

| Node | Kind | P0 | Depends on (merged) | Writes |
|---|---|---|---|---|
| `trino/ast-core` | grammar | ✅ | — | `trino/ast/*.go` |
| `trino/lexer` | grammar | ✅ | ast-core | `trino/parser/{lexer,keywords,tokens}.go` |
| `trino/parser-foundation` | grammar | ✅ | ast-core, lexer | `trino/parser/{parser,errors,split,diagnose,identifiers}.go` |
| `trino/types` | grammar | ✅ | parser-foundation | `trino/parser/datatypes.go` |
| `trino/expressions` | grammar | ✅ | types | `trino/parser/{expr,function,predicate}.go` |
| `trino/expr-json` | grammar | ✅ | expressions | `trino/parser/json.go` |
| `trino/parser-select` | grammar | ✅ | expressions | `trino/parser/{select,relation,ctes,window,setops}.go` |
| `trino/parser-match-recognize` | grammar | — | parser-select | `trino/parser/{match_recognize,row_pattern}.go` |
| `trino/parser-ddl` | grammar | ✅ | parser-select, types | `trino/parser/{create_table,alter_table,drop,schema_view,catalog_ddl,comment_analyze}.go` |
| `trino/parser-dml` | grammar | ✅ | parser-select | `trino/parser/{insert,update_delete,merge}.go` |
| `trino/parser-routines` | grammar | — | expressions, parser-select | `trino/parser/{function_def,routine_body}.go` |
| `trino/parser-dcl-tcl` | grammar | — | parser-foundation, expressions | `trino/parser/{grant_revoke,transaction,prepared}.go` |
| `trino/parser-utility` | grammar | — | parser-foundation, expressions | `trino/parser/{show,session,explain_call}.go` |
| `trino/catalog` | feature | ✅ | ast-core | `trino/catalog/*.go` |
| `trino/analysis` | feature | ✅ | parser-foundation, parser-select | `trino/analysis/*.go` |
| `trino/completion` | feature | ✅ | parser-select, catalog, analysis | `trino/completion/*.go` |
| `trino/deparse` | feature | ✅ | ast-core, types | `trino/deparse/*.go` |
| `trino/bytebase-switch` | feature | ✅ | analysis, completion, deparse, parser-ddl, parser-dml | `bytebase:backend/plugin/{parser,schema}/trino/**` |

## Merge layers (parallelization)

The orchestrator drains one merge-layer per fire (deps must be *merged*). Within a
layer, nodes with non-overlapping `writes` globs run concurrently.

| Layer | Nodes |
|---|---|
| 0 | ast-core |
| 1 | catalog, lexer |
| 2 | parser-foundation |
| 3 | types |
| 4 | deparse, expressions |
| 5 | expr-json, parser-select, parser-dcl-tcl, parser-utility |
| 6 | analysis, parser-ddl, parser-dml, parser-match-recognize, parser-routines |
| 7 | completion |
| 8 | bytebase-switch (INTEGRATE — human-guided, not auto-dispatched) |

Note: all `trino/parser/*` nodes write distinct files, so they are co-schedulable
within a layer; they serialize only via the dep edges above.

## Mapping to the consumed contract (P0 rationale)

| bytebase P0 API | satisfied by |
|---|---|
| `SplitSQL` | lexer, parser-foundation |
| `validateQuery` / query-type | parser-foundation, parser-select, analysis |
| `GetQuerySpan` (masking + lineage) | parser-select, expressions, analysis |
| `Completion` (c3) | parser-select, catalog, completion |
| `Diagnose` | parser-foundation (+ all statement parsers for clean coverage) |
| `GetDatabaseDefinition` | deparse |

Deferred (P1) nodes — match-recognize, routines, dcl-tcl, utility — are rarer
statement/expression forms; the import switch (`bytebase-switch`) depends only on
the P0 consumed-feature path, matching the doris precedent.

## Divergence ledger (seeded, oracle-adjudicated)

| Node | Case | Status | Disposition |
|---|---|---|---|
| parser-select | NEAREST JOIN | confirmed | docs-crawl error — exclude |
| lexer | `COLON_` literal `_:` | confirmed | implement `:` |
| parser-select | JSON_TABLE | flagged | P1 extension (docs-ahead) |
| parser-select | WITH SESSION | flagged | P1 extension (docs-ahead) |
| parser-dml | DML `@branch` | flagged | P1 extension (docs-ahead) |
| parser-select | set-op CORRESPONDING | flagged | P1 extension (docs-ahead) |
| parser-select | GROUP BY ALL/AUTO | proposed | pin semantics via oracle |

## Next

Hand back to `omni-engine-v2-migration` to run EXECUTE waves. First fire's ready-set
is `trino/ast-core`; draining proceeds layer-by-layer. Every grammar node validates
against the live Trino 481 oracle via `trino/internal/trinooracle`.
