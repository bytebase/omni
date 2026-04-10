# MySQL Semantic Layer (L2/L3) — Starmap

> **Status:** Phase 0 (IR design) — **Revision 2** post sub-agent review.
> **Owner:** TBD
> **Final target:** Phase 5 (SDL diff engine).
> **Scope discipline:** omni-only. No bytebase adapter code in this effort.
>
> **Revision 2 (2026-04-10):** incorporates findings from cc + codex sub-agent
> reviews of Revision 1. Key changes:
>   - Phase 1 split into 1a/1b/1c (was a single oversized phase).
>   - § 2 source-of-truth section rewritten to resolve the query-source vs
>     lineage-source contradiction; TiDB removed from references.
>   - Decision D5 reversed: MERGE view substitution moves from analyze-time to
>     consume-time (helper `(*Query).ExpandMergeViews()`).
>   - Decision D12 deferred: scope refactor moves to a new Phase 4.5 cleanup.
>   - § 7 Q7 (WithCheckOption) resolved: View.CheckOption already exists.
>   - Exit criteria tightened: per-query golden IR files, not just "corpus passes".
>   - Phase 1a gets the oracle harness on day one (was Phase 2).
>   - Phase 2 vendors the bytebase corpus snapshot to prevent silent drift.
>   - § 5 risks expanded (lower_case_table_names, 8.0-only, charset coercion,
>     information_schema, ONLY_FULL_GROUP_BY, deparse driver collision).
>   - § 4.1 IR catalog updated to match the actual `mysql/catalog/query.go`
>     skeleton: Q suffix unification, new InListExprQ/BetweenExprQ nodes,
>     RTEFunction reinstated for JSON_TABLE, MySQLBaseType enum, IsAggregate
>     on FuncCallExprQ, WindowClause on Query, ResolvedType minus Nullable.
>   - § 9 review notes appendix populated.

---

## 0. Why this exists

omni's PG side (`pg/catalog`) has been built up to be a real Go in-memory simulation of a PostgreSQL server. It has parser, lexer, DDL state machine, and a 4000-line **semantic analyzer** (`pg/catalog/analyze.go`) translated from PG's `analyze.c`. That analyzer is the load-bearing layer for nine independent DDL pipelines (view materialization, CHECK, DEFAULT, GENERATED, RLS, triggers, indexes, domains, dependencies) and is what makes column-level lineage and SDL diffing tractable on the PG side.

omni's MySQL side has the same end goal — be the in-memory MySQL simulation — but currently only has L0 (parser/lexer) and L1 (DDL state machine via the catalog walkthrough effort). **L2 (semantic analyzer) and L3 (deparse-from-IR) do not exist.** This is the single biggest gap blocking three downstream initiatives:

1. **bytebase MySQL query span migration** (column-level lineage). Today bytebase uses ANTLR listeners; we want to replace them with omni's IR consumed via a thin adapter on bytebase's side.
2. **MySQL SDL** (load / dump / diff schema-as-code). View and expression diffing fundamentally require structural representations, not text.
3. **Future**: advisor over semantic IR, query rewrite, masking policy validation, ACL inference, information_schema simulation.

This plan is **not** "build a query span analyzer". It is "build the MySQL semantic layer that PG already has, sized correctly for SDL as the final consumer, with query span as the first validating consumer along the way."

## 1. Layered architecture (target end state)

```
┌─────────────────────────────────────────────────────────┐
│  Application consumers (outside omni)                    │
│   • bytebase query span                                  │
│   • bytebase SDL drift detection                         │
│   • bytebase advisor / migration plan                    │
└────────────────────────┬────────────────────────────────┘
                         │ public API
┌────────────────────────▼────────────────────────────────┐
│  L4  Schema Diff / SDL                  (mysql/diff,    │
│      ─ catalog state diff                 mysql/sdl)    │
│      ─ uses L3 deparse for normalization                │
│      ─ uses L2 dependency edges for ordering            │
├──────────────────────────────────────────────────────────┤
│  L3  Deparse from IR                    (mysql/deparse) │
│      ─ Query IR → canonical SQL                          │
│      ─ AnalyzedExpr → canonical expression text          │
│      ─ replaces today's "echo original text" path        │
├──────────────────────────────────────────────────────────┤
│  L2  Semantic Analyzer (NEW)            (mysql/catalog) │
│      ─ AnalyzeSelectStmt    → *Query                    │
│      ─ AnalyzeStandaloneExpr → AnalyzedExpr              │
│      ─ produces RangeTableEntryQ / VarExprQ / ...        │
│      ─ stores Relation.AnalyzedQuery on view create      │
│      ─ stores Constraint.CheckAnalyzed, Column.GenAnal   │
├──────────────────────────────────────────────────────────┤
│  L1  Catalog State Machine (EXISTS)     (mysql/catalog) │
│      ─ DDL execution: catalog/exec.go, tablecmds.go,    │
│        viewcmds.go, indexcmds.go, …                     │
│      ─ in-memory Database / Table / View / Routine      │
│      ─ walkthrough oracle test infrastructure            │
├──────────────────────────────────────────────────────────┤
│  L0  Parser + AST + Lexer (EXISTS)      (mysql/parser,  │
│                                          mysql/ast)     │
└──────────────────────────────────────────────────────────┘
```

L0 and L1 are substantially complete. L2/L3/L4 are this plan's scope.

## 2. The MySQL-vs-PG divergence — the most important framing

PG omni was able to translate `analyze.c` because PG's analyzer is a clean isolated tier in PG's source: `raw parse tree → analyze.c → Query → optimizer → executor`. Each layer is a pure function.

**MySQL is structurally different.** MySQL's semantic analysis is in `sql/sql_resolver.cc`, `Query_block::prepare()`, `setup_fields()`, `setup_wild()`, `Item::fix_fields()` — and these functions are deeply intertwined with the optimizer (cost model, JOIN reordering, table elimination, query rewrites). There is no clean tier to translate.

**Therefore the MySQL L2 must be hand-written**, not translated. This is a multi-month commitment with no convenient escape hatch.

### 2.1 Sources of truth (strictly ordered)

For any decision about MySQL behavior the analyzer must produce:

1. **Real MySQL 8.0.34 in a container** — the only **behavioral oracle**. omni already has this infrastructure (`mysql-catalog-oracle-worker`, `mysql-deparse-oracle-driver`). The analyzer effort plugs into it.
2. **MySQL official documentation** — data type rules, conversion, NATURAL/USING semantics, `ONLY_FULL_GROUP_BY`, view algorithm rules, function return types.
3. **MySQL source** (`sql/sql_resolver.cc`, `sql/item_*.cc`, `sql/sql_view.cc`) — read for hard cases, never copy-translated. C++ and not portable.

These are the **only** sources of truth for behavior. Tertiary references — used for shape and patterns, never for behavioral authority:

4. **omni's own assets**:
   - `mysql/deparse/resolver.go` — already implements scope, CTE virtual tables, NATURAL/USING coalescing for column qualification. Roughly 30% of the work.
   - `mysql/completion/resolve.go` — scope walking patterns.
   - `mysql/catalog/viewcmds.go` `deparseViewSelect()` — current view column derivation seam to be replaced.
5. **bytebase's existing ANTLR query span code** — used as a *queries* source for the test corpus and as a behavioral baseline for diff reports. **NOT used as a lineage source-of-truth** — see § 2.2.

**Hard rule for the whole effort:** when MySQL behavior is ambiguous, the answer is "what does MySQL 8.0.34 say in a container", **not** "what does TiDB do" and **not** "what does bytebase ANTLR do". TiDB is intentionally absent from this list — it is a different database with its own divergences from MySQL. Reading TiDB's planner code creates real risk of subtle drift toward TiDB-specific behavior, with no offsetting authority gain.

### 2.2 Query source vs lineage source — they are not the same thing

A common conflation that the Revision 1 draft of this plan made: treating "the test corpus" as a single thing. It is actually two independent inputs:

**Query source** (the input SQL strings):
- Drawn from `bytebase` MySQL query span golden corpus (the largest pool)
- Plus hand-written cases targeting analyzer edge cases
- Plus MySQL documentation examples for canonical syntax forms
- Plus Phase-specific synthetic queries
- Vendored at a pinned hash into omni — see § 3 Phase 2 deliverables
- Updated explicitly via a documented protocol

**Lineage source** (the expected answer for each query — "which columns from which tables does this SELECT read?"):
- Derived from real MySQL 8.0.34 container by querying `information_schema.VIEW_COLUMN_USAGE` / `VIEW_TABLE_USAGE` and / or running `EXPLAIN FORMAT=JSON` and parsing the `used_columns` fields
- Stored as per-query golden IR files committed alongside the query corpus
- The **only** authoritative source

bytebase's existing ANTLR analyzer output is NOT a lineage source. We can run it as a sanity-check and report differences (omni vs bytebase), but where they disagree the container wins. If omni and the container agree but bytebase disagrees, that's a bytebase bug to be filed separately, not a blocker on omni Phase 2.

This split removes the contradiction in Revision 1, where § 2 said "the container is the truth" but § 3 Phase 2 made the bytebase ANTLR output the de facto spec.

## 3. The 6+ phases — final target Phase 5

Each phase has a concrete deliverable, an exit criterion, and a downstream consumer that validates it. **No phase is "done" until something downstream is consuming the output and a verification harness is green.**

Exit criteria use **per-query golden IR files** rather than "the corpus passes". A test fixture is a tuple `(inputSQL, inputCatalogState, expectedQueryIR)`. The expectedQueryIR is hand-curated with container assistance and committed. "Passes" = exact structural match modulo explicitly allowed nondeterminism (e.g. RangeTable order within a single FROM-list).

### Phase 0 — IR design (CURRENT, complete pending Revision 2 sign-off)

**Deliverables:**
- `mysql/catalog/query.go` — type-only declarations, ~900 lines, compiles clean. **DONE.**
- `docs/plans/2026-04-09-mysql-semantic-layer.md` — this document. **DONE (Revision 2).**
- `docs/plans/2026-04-09-mysql-semantic-layer-review-bundle.md` — review onboarding.

**Exit criterion:** review by user + cc + codex; decision log signed off; § 9 review notes populated. **No analyzer code written until exit.**

### Phase 1 — SELECT analyzer, split into 1a / 1b / 1c

The Revision 1 draft tried to do all of SELECT analysis in one phase. Sub-agent review correctly flagged this as a several-month phase masquerading as a milestone. Splitting forces honest exit gates and prevents a half-done analyzer being validated against a half-done corpus.

#### Phase 1a — Foundation: single-relation SELECT

**Deliverables:**
- `mysql/catalog/analyze.go` — `AnalyzeSelectStmt(*ast.SelectStmt) (*Query, error)` entry point
- `mysql/catalog/analyze_targetlist.go` — SELECT list → TargetEntryQ list, including `*` and `t.*` expansion via catalog lookup
- `mysql/catalog/analyze_expr.go` — expression → `AnalyzedExpr`, **column references resolved to VarExprQ (RangeIdx + AttNum), all other expression types stored without type inference**
- `mysql/catalog/scope.go` — scope stack, **copied** (not refactored — see Phase 4.5) from `mysql/deparse/resolver.go` with a `// TODO(scope-unify): see deparse/resolver.go` marker
- `mysql/catalog/analyze_test.go` — golden IR tests against fixture catalog states
- **`mysql/catalog/analyze_oracle_test.go` — oracle harness from day one**, borrowing `mysql-catalog-oracle-worker` container infrastructure. Each fixture is run against a real MySQL 8.0.34 container; expected lineage is derived via `information_schema.VIEW_COLUMN_USAGE` / `VIEW_TABLE_USAGE`.

**Scope:**
- ✅ Single-table SELECT
- ✅ WHERE, GROUP BY, HAVING, ORDER BY, LIMIT, DISTINCT
- ✅ `*` and `t.*` expansion against catalog
- ❌ JOINs (Phase 1b)
- ❌ Subqueries (Phase 1c)
- ❌ Set operations (Phase 1c)
- ❌ Expression type inference (Phase 3)

**Exit criterion:** 30-query golden corpus (single-relation cases drawn from bytebase corpus + hand-written) passes byte-for-byte against committed expected IR files; oracle harness shows zero divergence between omni's lineage and container-derived lineage on the same corpus.

#### Phase 1b — JOINs and FROM-subqueries

**Deliverables:**
- `mysql/catalog/analyze_from.go` — FROM clause → RangeTable + JoinTreeQ
- Update `analyze_expr.go` to resolve column references through multiple RTEs
- USING / NATURAL JOIN column coalescing logic (port from `deparse/resolver.go`)
- FROM-subquery support: `FROM (SELECT ...) AS x`
- Updated tests + 50-query corpus

**Scope:**
- ✅ All JOIN types: INNER, LEFT, RIGHT, CROSS, STRAIGHT_JOIN
- ✅ USING and NATURAL with proper column coalescing
- ✅ FROM-subqueries (non-LATERAL)
- ✅ Multiple FROM items (`FROM a, b, c` cross-product form)
- ❌ Correlated subqueries (Phase 1c)
- ❌ LATERAL subqueries (Phase 1c)

**Exit criterion:** 50-query golden corpus passes; oracle harness clean; corpus includes the known-tricky NATURAL/USING ordering cases from MySQL docs.

#### Phase 1c — Correlated subqueries, CTEs, set operations

**Deliverables:**
- Correlated subquery support (LevelsUp on VarExprQ)
- LATERAL subquery support
- `WITH` and `WITH RECURSIVE` (CommonTableExprQ + RTECTE)
- Scalar / IN / EXISTS / ANY / ALL subquery expressions (SubLinkExprQ)
- `IN (literal_list)` non-subquery form (InListExprQ)
- `BETWEEN ... AND ...` (BetweenExprQ)
- `UNION` / `INTERSECT` / `EXCEPT` (with set-op result-shape unification)
- 100-query corpus

**Scope:**
- ✅ All forms left after 1a/1b
- ❌ Window functions full structure (`WindowClause` populated, `WindowDefQ.FrameClause` raw text only)
- ❌ Expression type inference

**Exit criterion:** 100-query golden corpus passes; oracle harness clean; lineage walker can handle every fixture in the corpus through its full RTE chain.

### Phase 2 — View pipeline + query span end-to-end

**Deliverables:**
- `mysql/catalog/viewcmds.go` updated: `createView` calls `AnalyzeSelectStmt`, stores `*Query` on `View.AnalyzedQuery`
- `View.Columns` derived from `Query.TargetList` instead of from current `deparseViewSelect` text path
- `(*Query).ExpandMergeViews()` helper that performs MERGE-view substitution at consume time. Lineage walkers call this; deparse does not.
- `mysql/catalog/query_span_oracle_test.go` — verification: for a curated set of (catalog state, SELECT) pairs that include views, the analyzed Query's column lineage (computed by an in-test walker) matches container-derived lineage. The walker lives in tests only, mirroring `pg/catalog/query_span_test.go`.
- **Bytebase corpus snapshot vendored**: a snapshot of bytebase's MySQL query span golden corpus is committed under `mysql/catalog/testdata/bytebase_corpus/` at a pinned bytebase hash. A `README.md` documents the update protocol (manual sync, no automation) so silent drift is impossible.

**Scope:**
- View create / replace / drop wired to analyzer
- View-of-view recursion (handled via the RTERelation IsView path)
- Lineage walker covers: direct table refs, JOINs, set ops, subqueries, CTEs, MERGE views (post-expansion), TEMPTABLE views (opaque)

**Exit criterion:** for each fixture in the vendored bytebase corpus, omni's analyzed-IR-based lineage matches container-derived lineage. A separate diff report shows omni vs bytebase ANTLR; non-zero differences are recorded but not blockers (they may be bytebase bugs).

### Phase 3 — Standalone expression analyzer + minimal type inference

**Deliverables:**
- `mysql/catalog/analyze.go` adds `AnalyzeStandaloneExpr(expr ast.ExprNode, table *Table) (AnalyzedExpr, error)`
- `mysql/catalog/types.go` — type inference utilities (function return type table, `*` expansion type propagation, simple coercions like CASE / COALESCE column unification)
- Function return type table for the most common ~50 functions, drawn from a corpus inventory (which functions actually appear in real-world view definitions in the bytebase corpus + MySQL system catalogs). Exact list committed as part of the deliverable.
- DDL hookups:
  - `tablecmds.go` / `altercmds.go`: column DEFAULT and GENERATED expressions analyzed and stored
  - `constraint.go`: CHECK constraint expressions analyzed and stored
  - `indexcmds.go`: functional index expressions analyzed and stored
- Updated walkthrough tests assert structural form is captured

**End-of-phase decision point:** revisit D6 (implicit casts). If CHECK constraint diff or GENERATED diff is observed to false-positive on coercion variations during this phase, materialize implicit casts before moving to Phase 4. The Revision 1 plan deferred this to Phase 5 — too late.

**Exit criterion:** every catalog walkthrough test (the `wt_*_test.go` battery) still passes; new tests assert that `Constraint.CheckAnalyzed`, `Column.DefaultAnalyzed`, `Column.GeneratedAnalyzed`, `Index.ExprAnalyzed` are populated for the relevant test fixtures; view column types in the catalog match what `SHOW COLUMNS FROM view` returns from a real MySQL 8.0.34 container (oracle-verified).

### Phase 4 — Deparse from IR (L3)

**Deliverables:**
- `mysql/deparse/from_query.go` — `DeparseQuery(*catalog.Query) (string, error)` mirroring PG `ruleutils.go::getQueryDef`
- `mysql/deparse/from_expr.go` — `DeparseExpr(catalog.AnalyzedExpr, *DeparseContext) (string, error)` mirroring PG `DeparseExpr`
- `mysql/catalog/show.go` (or wherever `SHOW CREATE VIEW` is implemented) switches to `DeparseQuery(view.AnalyzedQuery)` instead of echoing original text
- `SHOW CREATE TABLE` for tables with CHECK / DEFAULT / GENERATED switches to `DeparseExpr(...)` for those clauses
- Round-trip oracle test: `original SQL → analyze → deparse → re-analyze → assert IR equality`
- **Coordination with `mysql-deparse-driver`**: this phase moves goalposts for the in-flight deparse-driver work. Settle ownership before Phase 1c ends — see § 5 risks.

**Exit criterion:** the existing MySQL deparse oracle (`mysql-deparse-oracle-worker`) passes when fed analyzed-then-deparsed view definitions. This is the load-bearing prerequisite for SDL.

### Phase 4.5 — Scope unification (deferred from Phase 1)

**Deliverables:**
- Merge `mysql/catalog/scope.go` (the Phase 1a copy) and `mysql/deparse/resolver.go`'s scope code into a single shared implementation under `mysql/catalog/scope.go`
- Update `mysql/deparse` to depend on the unified scope
- Update `mysql/completion` if needed

**Why deferred:** sub-agent review found that `mysql/deparse/resolver.go` is **939 lines** (not the ~200 originally estimated) and is intertwined with USING/NATURAL coalescing, CTE virtual tables, set-op CTE hoisting, and ORDER BY ordinal resolution. Doing this refactor in the middle of Phase 1 would entangle two fast-moving fronts. Doing it after Phase 4, when both consumer shapes are stable, makes the unified API design obvious.

**Exit criterion:** no behavior change in `mysql/deparse` (existing deparse tests all pass); the analyzer's scope code now references `mysql/catalog/scope.go` directly.

### Phase 5 — SDL Diff Engine (FINAL TARGET)

**Deliverables:**
- `mysql/diff/` package (new)
  - `diff.go` — top-level diff: two `*catalog.Catalog` snapshots → ordered list of `Change` objects
  - `objects.go` — per-object-kind diff (database, table, view, index, constraint, routine, event, trigger)
  - `expr_diff.go` — structural `AnalyzedExpr` diff (so semantically-equivalent textual variants don't false-positive)
  - `query_diff.go` — structural `Query` diff for view bodies
  - `dependency.go` — dependency graph from analyzer output to compute drop/create order
- `mysql/sdl/` package (new) — paralleling `pg/sdl/`
  - `loadsdl.go` — apply an SDL document to a fresh catalog state via L1 walkthrough
  - `dumpsdl.go` — serialize a catalog state via L3 deparse
  - `apply.go` — given two catalog states, produce ALTER script via diff engine

**Exit criterion:** SDL round-trip test on a substantial schema (50+ objects, including views, generated columns, CHECK constraints, foreign keys, partitioned tables, stored routines): `dump → load → dump` produces byte-identical output, AND `diff(stateA, stateB)` produces a script that when applied via L1 walkthrough yields stateB.

## 4. Phase 0 IR design — Revision 2

This is the centerpiece of the review. § 4.1 mirrors the actual contents of `mysql/catalog/query.go` after Revision 2 fixes.

### 4.1 Type catalog

**Naming convention:** all Query-IR struct types end in `Q` ("Query IR"), with two exceptions: `Query` itself (the namespace anchor) and `AnalyzedExpr` (an interface, not a node). Enums do NOT take the `Q` suffix, with one mechanical exception: `JoinTypeQ` collides with `mysql/ast.JoinType` and is forced to take the suffix.

```go
package catalog

// --- Top-level analyzed query ---

type Query struct {
    CommandType  CmdType
    TargetList   []*TargetEntryQ
    RangeTable   []*RangeTableEntryQ
    JoinTree     *JoinTreeQ
    GroupClause  []*SortGroupClauseQ
    HavingQual   AnalyzedExpr
    SortClause   []*SortGroupClauseQ
    LimitCount   AnalyzedExpr
    LimitOffset  AnalyzedExpr
    Distinct     bool
    SetOp        SetOpType
    AllSetOp     bool
    LArg, RArg   *Query
    CTEList      []*CommonTableExprQ
    IsRecursive  bool
    WindowClause []*WindowDefQ        // NEW (Rev 2): WINDOW w AS (...) declarations
    HasAggs      bool
    LockClause   *LockingClauseQ
    SQLMode      SQLMode              // captured at analyze time
}

type CmdType int
const (
    CmdSelect CmdType = iota
    CmdInsert
    CmdUpdate
    CmdDelete
)

// Removed in Rev 2: Query.DistinctOn (PG-only feature; YAGNI for MySQL).

// --- TargetList ---

type TargetEntryQ struct {
    Expr         AnalyzedExpr
    ResNo        int
    ResName      string
    ResJunk      bool
    ResOrigDB    string
    ResOrigTable string
    ResOrigCol   string
}

// --- Range table ---

type RangeTableEntryQ struct {
    Kind          RTEKind
    DBName        string
    TableName     string
    Alias         string
    ERef          string
    ColNames      []string
    ColTypes      []*ResolvedType
    ColCollations []string
    Subquery      *Query
    JoinType      JoinTypeQ
    JoinUsing     []string
    CTEIndex      int
    CTEName       string
    Lateral       bool
    IsView        bool
    ViewAlgorithm ViewAlgorithm
    FuncExprs     []AnalyzedExpr      // NEW (Rev 2): RTEFunction support for JSON_TABLE
}

type RTEKind int
const (
    RTERelation RTEKind = iota
    RTESubquery
    RTEJoin
    RTECTE
    RTEFunction                       // NEW (Rev 2): MySQL JSON_TABLE
)

type ViewAlgorithm int
const (
    ViewAlgNone       ViewAlgorithm = iota  // NEW (Rev 2): zero value = not a view
    ViewAlgUndefined
    ViewAlgMerge
    ViewAlgTemptable
)

// --- Join tree ---

type JoinTreeQ struct {
    FromList []JoinNode
    Quals    AnalyzedExpr
}

type JoinNode interface{ joinNodeTag() }

type RangeTableRefQ struct{ RTIndex int }
func (*RangeTableRefQ) joinNodeTag() {}

type JoinExprNodeQ struct {
    JoinType    JoinTypeQ
    Left, Right JoinNode
    Quals       AnalyzedExpr
    UsingClause []string
    Natural     bool
    RTIndex     int
}
func (*JoinExprNodeQ) joinNodeTag() {}

type JoinTypeQ int
const (
    JoinInner JoinTypeQ = iota
    JoinLeft
    JoinRight
    JoinCross
    JoinStraight   // STRAIGHT_JOIN — MySQL only
    // No FULL OUTER — MySQL doesn't support it
)

// --- CTE ---

type CommonTableExprQ struct {
    Name        string
    ColumnNames []string
    Query       *Query
    Recursive   bool
}
// Removed in Rev 2: CTERefCount (no consumer; speculative).

// --- SortGroupClause ---

type SortGroupClauseQ struct {
    TargetIdx  int     // 1-based (PG uses 0-based; documented deviation)
    Descending bool
    NullsFirst bool
}

// --- Set op enum ---

type SetOpType int
const (
    SetOpNone SetOpType = iota
    SetOpUnion
    SetOpIntersect
    SetOpExcept
)

// --- Locking clause ---

type LockingClauseQ struct {
    Strength   LockStrength
    Tables     []string         // OF list; must be empty for LockInShareMode
    WaitPolicy LockWaitPolicy   // must be LockWaitDefault for LockInShareMode
}

type LockStrength int
const (
    LockNone        LockStrength = iota
    LockForUpdate
    LockForShare        // FOR SHARE (8.0+, supports OF/NOWAIT/SKIP LOCKED)
    LockInShareMode    // legacy syntax, no modifiers
)

type LockWaitPolicy int
const (
    LockWaitDefault LockWaitPolicy = iota
    LockWaitNowait
    LockWaitSkipLocked
)

// --- AnalyzedExpr interface ---

type AnalyzedExpr interface {
    exprType() *ResolvedType
    exprCollation() string
    analyzedExprTag()
}

// Package-level singleton for boolean-result expressions.
// Used by BoolExprQ, NullTestExprQ, InListExprQ, BetweenExprQ.
var BoolType = &ResolvedType{ BaseType: BaseTypeTinyIntBool }

// --- AnalyzedExpr implementations ---

type VarExprQ struct {
    RangeIdx  int
    AttNum    int          // 1-based
    LevelsUp  int
    Type      *ResolvedType
    Collation string
}

type ConstExprQ struct {
    Type      *ResolvedType
    Collation string
    IsNull    bool
    Value     string
}

type FuncCallExprQ struct {
    Name        string
    Args        []AnalyzedExpr
    IsAggregate bool          // NEW (Rev 2): replaces fragile string-matching
    Distinct    bool
    Over        *WindowDefQ
    ResultType  *ResolvedType
    Collation   string
}

type OpExprQ struct {
    Op         string
    Left       AnalyzedExpr  // nil for prefix
    Right      AnalyzedExpr
    ResultType *ResolvedType
    Collation  string
}

type BoolExprQ struct {
    Op   BoolOpType
    Args []AnalyzedExpr
}
type BoolOpType int
const ( BoolAnd BoolOpType = iota; BoolOr; BoolNot )

type CaseExprQ struct {
    TestExpr   AnalyzedExpr   // nil for searched CASE
    Args       []*CaseWhenQ
    Default    AnalyzedExpr
    ResultType *ResolvedType
    Collation  string
}
type CaseWhenQ struct { Cond, Then AnalyzedExpr }

type CoalesceExprQ struct {
    Args       []AnalyzedExpr
    ResultType *ResolvedType
    Collation  string
}

type NullTestExprQ struct {
    Arg    AnalyzedExpr
    IsNull bool
}

type SubLinkExprQ struct {
    Kind       SubLinkKind
    TestExpr   AnalyzedExpr
    Op         string
    Subquery   *Query
    ResultType *ResolvedType
    Collation  string                // NEW (Rev 2): scalar subqueries returning strings
}
type SubLinkKind int
const (
    SubLinkExists SubLinkKind = iota
    SubLinkScalar
    SubLinkIn
    SubLinkAny
    SubLinkAll
)

// NEW (Rev 2): non-subquery IN.
type InListExprQ struct {
    Arg     AnalyzedExpr
    List    []AnalyzedExpr
    Negated bool   // NOT IN
}

// NEW (Rev 2): BETWEEN as a first-class node (not lowered to AND of comparisons).
type BetweenExprQ struct {
    Arg     AnalyzedExpr
    Lower   AnalyzedExpr
    Upper   AnalyzedExpr
    Negated bool   // NOT BETWEEN
}

type RowExprQ struct {
    Args       []AnalyzedExpr
    ResultType *ResolvedType
}

type CastExprQ struct {
    Arg        AnalyzedExpr
    TargetType *ResolvedType
    Collation  string                 // NEW (Rev 2): CAST(... CHARACTER SET ... COLLATE ...)
}

// --- Window def ---

type WindowDefQ struct {
    Name        string
    PartitionBy []AnalyzedExpr
    OrderBy     []*SortGroupClauseQ
    FrameClause string                // raw text in Phase 1; structured in Phase 4 if needed
}

// --- ResolvedType + MySQLBaseType (formerly ColumnType) ---
//
// Renamed from ColumnType in Rev 1 to avoid shadowing Column.ColumnType (string field).
// Nullable removed in Rev 2 (it's a column property, not a type property).

type ResolvedType struct {
    BaseType   MySQLBaseType  // typed enum, not free-form string
    Unsigned   bool
    Length     int            // VARCHAR(N), CHAR(N), VARBINARY(N), BINARY(N), BIT(N)
    Precision  int            // DECIMAL(P, S)
    Scale      int
    FSP        int            // DATETIME(N), TIME(N), TIMESTAMP(N)
    Charset    string
    Collation  string
    EnumValues []string
    SetValues  []string
}

type MySQLBaseType int
const (
    BaseTypeUnknown MySQLBaseType = iota
    BaseTypeTinyInt
    BaseTypeTinyIntBool         // TINYINT(1) special case (client libraries treat as boolean)
    BaseTypeSmallInt
    BaseTypeMediumInt
    BaseTypeInt
    BaseTypeBigInt
    BaseTypeDecimal
    BaseTypeFloat
    BaseTypeDouble
    BaseTypeBit
    BaseTypeDate
    BaseTypeDateTime
    BaseTypeTimestamp
    BaseTypeTime
    BaseTypeYear
    BaseTypeChar
    BaseTypeVarchar
    BaseTypeBinary
    BaseTypeVarBinary
    BaseTypeTinyText
    BaseTypeText
    BaseTypeMediumText
    BaseTypeLongText
    BaseTypeTinyBlob
    BaseTypeBlob
    BaseTypeMediumBlob
    BaseTypeLongBlob
    BaseTypeEnum
    BaseTypeSet
    BaseTypeJSON
    BaseTypeGeometry
    BaseTypePoint
    BaseTypeLineString
    BaseTypePolygon
    BaseTypeMultiPoint
    BaseTypeMultiLineString
    BaseTypeMultiPolygon
    BaseTypeGeometryCollection
)

// --- SQL mode ---

type SQLMode uint64
const (
    SQLModeAnsiQuotes SQLMode = 1 << iota
    SQLModePipesAsConcat
    SQLModeIgnoreSpace          // NEW (Rev 2): affects function call vs paren expr parsing
    SQLModeHighNotPrecedence    // NEW (Rev 2): affects NOT precedence
    SQLModeNoBackslashEscapes   // NEW (Rev 2): affects string literal parsing
    SQLModeOnlyFullGroupBy
    SQLModeNoUnsignedSubtraction
    SQLModeStrictAllTables
    SQLModeStrictTransTables
    SQLModeNoZeroDate
    SQLModeNoZeroInDate
    SQLModeRealAsFloat
)
```

### 4.2 Decision log

| # | Decision | Rev 2 Status | Rationale / change reason |
|---|---|---|---|
| D1 | Schema identifier: `(DBName, TableName)` instead of OID | **Decided: Approve** | MySQL has no schema namespace; no synthetic OID space worth maintaining. Sub-agent suggested string interning as future optimization; not in Phase 0 |
| D2 | 1-based AttNum (PG convention) | **Decided: Approve** | Algorithm parity with PG lineage walker |
| D3 | Type representation: nullable `*ResolvedType` | **Decided: Approve with sentinel** | Phase 1 leaves nil; Phase 3 fills. `BoolType` package-level singleton means boolean expressions never return nil — eliminates one class of nil-checks |
| D4 | Collation: string ("utf8mb4_0900_ai_ci") | **Decided: Approve** | Right call. MySQL collations are named, not OIDs |
| D5 | View algorithm handling | **REVERSED in Rev 2: consume-time substitution** | Sub-agent reviewer 1 caught this. Analyzer leaves views opaque (`IsView=true`, `Subquery=nil`); consumers call `(*Query).ExpandMergeViews()` if they want lineage transparency. Reasons: (1) deparse must reproduce user's view text, not inlined body; (2) MySQL's actual MERGE is rewrite at query-rewrite time; (3) WITH CHECK OPTION interacts badly with eager substitution |
| D6 | Implicit casts NOT materialized | **Decided: Approve for Phase 1; revisit at end of Phase 3** | Rev 1 said "revisit Phase 5" — too late. CHECK constraint and GENERATED diffing in Phase 3 are likely to surface coercion-fidelity gaps; decide before Phase 4 |
| D7 | SQLMode on Query | **Decided: Approve** | Captured at analyze time so deparse can reproduce session-dependent semantics. Not on every expr (would be bloat) |
| D8 | LATERAL bool on RangeTableEntryQ | **Decided: Approve** | MySQL 8.0.14+ |
| D9 | CmdType enum reservation | **Decided: Approve** | Phase 1 only emits CmdSelect; reserves IR shape for DML in Phase 8+ |
| D10 | WindowDefQ.FrameClause as raw text | **Decided: Approve for Phase 1; revisit Phase 4** | Sub-agent flagged round-trip drift risk; deparse round-trip oracle in Phase 4 will catch it |
| D11 | `query_span.go` not in omni — test-only walker mirrors PG | **Decided: Approve** | Bytebase owns lineage extraction; omni ships only the IR. Do NOT ship a vendored reference under `internal/lineage/` |
| D12 | Scope refactor | **DEFERRED to Phase 4.5** | Rev 1 estimated ~200 lines; actual `mysql/deparse/resolver.go` is 939 lines and entangled with USING/NATURAL/CTE virtual tables. Phase 1a will copy code with `// TODO(scope-unify)` markers |
| D13 | TargetEntryQ provenance fields | **Decided: Approve** | Single-source only; multi-source (COALESCE-of-cols) is the lineage walker's job |
| D14 | STRAIGHT_JOIN as JoinTypeQ variant | **Decided: Approve** | MySQL-only optimizer hint; affects deparse not lineage |
| **D15** *(Rev 2)* | Q-suffix unification | **Decided: Apply uniformly** | All IR structs end in Q except `Query` itself (anchor) and `AnalyzedExpr` (interface). All enums skip Q except `JoinTypeQ` (mechanical collision with `mysql/ast.JoinType`) |
| **D16** *(Rev 2)* | `Nullable` removed from `ResolvedType` | **Decided: Move to column-level** | Type-vs-column conflation would cause SDL diff false positives. `Column.Nullable` already exists in `mysql/catalog/table.go:64` |
| **D17** *(Rev 2)* | `MySQLBaseType` typed enum vs free-form `DataType string` | **Decided: Use typed enum** | ~35 values; safer for switch statements; `BaseTypeTinyIntBool` carries the `TINYINT(1)` special case at the type-identity level |
| **D18** *(Rev 2)* | `IsAggregate` field on `FuncCallExprQ` | **Decided: Add** | Without it, lineage walkers and HasAggs validation fall back to string-matching function names — exactly the fragility this IR removes |

### 4.3 What is explicitly NOT in the IR

- No `Aggref` separate from `FuncCallExprQ` — aggregates are `FuncCallExprQ` with `IsAggregate=true`. Window functions are `FuncCallExprQ` with `Over != nil`.
- No `RelabelExpr` (binary-compatible cast) — implicit casts not materialized (D6).
- No `ArrayExpr` — MySQL doesn't have native array types.
- No `XmlExpr` / `JsonExpr` — JSON handled as `FuncCallExprQ` to `JSON_*` functions (operator forms `->`/`->>` deferred — see "deferred expression nodes" in `query.go`).
- No `RowCompareExpr` — row comparisons handled as `OpExprQ` over `RowExprQ`.
- No `Param` node distinct from `ConstExprQ` — prepared statement params are out of analyzer scope.
- No `Tablesample` — MySQL doesn't support TABLESAMPLE.
- No `LikeExprQ` — Phase 1 lowers `LIKE` without ESCAPE to `OpExprQ{Op: "LIKE"}`. The 3-operand `LIKE x ESCAPE y` form is deferred to Phase 3.
- No `MatchAgainstExprQ` — fulltext predicate, deferred to Phase 4.
- No `WithCheckOption` field on Query/RTE — already on `View.CheckOption` at `mysql/catalog/table.go:89` (see § 7 Q7 resolution).

### 4.4 Surface area mapping: PG → MySQL

| PG type / concept | MySQL equivalent | Notes |
|---|---|---|
| `Query` | `Query` | 1:1, with `SQLMode`, `LockClause`, `WindowClause` MySQL-specific |
| `RangeTblEntry.relid` (OID) | `(DBName, TableName)` | No OID space |
| `Var.varno` | `VarExprQ.RangeIdx` (0-based) | |
| `Var.varattno` | `VarExprQ.AttNum` (1-based) | Matches PG |
| `Var.varlevelsup` | `VarExprQ.LevelsUp` | |
| `Var.vartype` (OID) | `VarExprQ.Type` (`*ResolvedType`) | Nullable in Phase 1 |
| `Const.consttype` | `ConstExprQ.Type` | Nullable in Phase 1 |
| `FuncExpr` + `Aggref` + `WindowFunc` | `FuncCallExprQ` (with `IsAggregate` + `Over`) | Merged |
| `OpExpr` | `OpExprQ` | No operator OID |
| `BoolExpr` | `BoolExprQ` | |
| `CaseExpr` | `CaseExprQ` | |
| `CoalesceExpr` | `CoalesceExprQ` | |
| `NullTest` | `NullTestExprQ` | |
| `SubLink` | `SubLinkExprQ` | |
| `ScalarArrayOpExpr` | `InListExprQ` | Different shape: PG models `x = ANY(array)`; MySQL `x IN (list)` |
| `RowExpr` | `RowExprQ` | |
| `RelabelType`, `CoerceViaIO`, `ArrayCoerceExpr` | (none) | Implicit casts not materialized — D6 |
| `JoinExpr` / `FromExpr` | `JoinExprNodeQ` / `JoinTreeQ` | |
| `CommonTableExpr` | `CommonTableExprQ` | |
| `SortGroupClause` | `SortGroupClauseQ` | Simpler — no OIDs |
| `LockingClause` | `LockingClauseQ` | MySQL syntax different |
| Set ops | `Query.SetOp` + `LArg/RArg` | |
| (PG has no equivalent) | `BetweenExprQ`, `InListExprQ` | PG lowers; MySQL keeps as nodes for round-trip fidelity |
| `RangeFunction` | `RangeTableEntryQ.Kind == RTEFunction` | MySQL only has `JSON_TABLE` |
| `Tablesample`, `WithCheckOption` | (none) | Out of scope — § 4.3 |

### 4.5 Storage hooks on existing catalog types

These are the L1 mutations Phase 2/3 will need:

```go
// mysql/catalog/table.go — additions in Phase 3
type Column struct {
    // existing fields including Nullable...
    DefaultAnalyzed   AnalyzedExpr // Phase 3
    GeneratedAnalyzed AnalyzedExpr // Phase 3
}
type Constraint struct {
    // existing fields...
    CheckAnalyzed AnalyzedExpr // Phase 3
}
type Index struct {
    // existing fields...
    ExprAnalyzed []AnalyzedExpr // Phase 3 (functional index expressions)
}

// mysql/catalog/table.go — addition in Phase 2
type View struct {
    // existing fields including CheckOption (already present at line 89)...
    AnalyzedQuery *Query // Phase 2
}
```

## 5. Risks and mitigations

| Risk | Severity | Mitigation |
|---|---|---|
| MySQL semantic edge cases (NATURAL/USING ordering, ONLY_FULL_GROUP_BY, view algorithm corner cases) consume more time than estimated | High | Phase 1a/1b/1c each have hard exit criteria against real MySQL container; we eat the discovery cost early |
| **Hand-write commitment risk: no TiDB fallback means semantic edge-case discovery cost is fully on the oracle harness** | High | Phase 1a deliverable is the oracle harness on day one. If we discover a class of MySQL semantics we cannot model, we surface it as a Phase 1 finding, not a Phase 5 surprise |
| IR design choice forces Phase 4–5 rewrite | High | Phase 0 review by 3 parties before any consumer code is written; Rev 2 reflects review findings |
| `mysql/deparse/resolver.go` refactor (D12) collides with active deparse work | Medium | **Deferred to Phase 4.5**; Phase 1 copies code with TODO markers |
| **`lower_case_table_names` dev/prod mismatch** (Mac default 2, Linux default 0) breaks identifier resolution | Medium | Capture on `Query.SQLMode` or globally; analyzer's catalog lookup respects the setting. Document the expected default |
| **omni MySQL is 8.0-only** but plan didn't say so | Medium | **Stated explicitly in § 6**: omni MySQL targets 8.0.34 baseline. 5.7 is not supported (CTEs, window functions, CHECK enforcement, INTERSECT/EXCEPT all 8.0+) |
| **Charset/collation propagation through CONCAT** ("illegal mix of collations") | Medium | Phase 3 problem — function return type table must consider collation coercion. Document as Phase 3 deliverable |
| **`information_schema` / `performance_schema` references in views** | Medium | omni catalog must register these as virtual schemas. If not done in Phase 1, views referencing them will fail to analyze. Add to Phase 2 prerequisites |
| **`ONLY_FULL_GROUP_BY` validation is non-trivial functional-dependency analysis** | Low | Defer to Phase 3+; document policy: when `ONLY_FULL_GROUP_BY` is set in `Query.SQLMode`, the analyzer captures it but does not enforce it. Real MySQL would reject — we let it through and surface as a Phase 3+ TODO |
| **bytebase corpus drift**: if bytebase changes its golden corpus, Phase 2's bar silently changes | Low | **Vendored snapshot at pinned hash** with documented manual update protocol (Phase 2 deliverable) |
| **Phase 4 deparse-from-IR collides with `mysql-deparse-driver`** in-flight work | Medium | Settle ownership before Phase 1c ends. Either Phase 4 owns view/CHECK/DEFAULT/GENERATED deparse and the driver owns everything else, or the driver pauses on those object kinds. Decide explicitly |
| MySQL function return type table is too large to maintain | Medium | Phase 3 only requires the ~50 functions used in real-world view definitions; document the policy "add types lazily as test corpus expands" |
| MERGE-vs-TEMPTABLE view substitution semantics wrong | Low | Phase 2 oracle test compares against `EXPLAIN` output of real MySQL; substitution is consume-time so the analyzer's model never has to be lossy |
| SDL diff produces wrong DDL ordering | Medium | Phase 5 has a round-trip exit criterion |
| Bytebase migration stalls because their adapter is independent of this work | Low | Out of scope by user direction; surfaces as a separate timeline issue |

## 6. Out of scope (this plan)

- Bytebase MySQL query span adapter code (lives in bytebase repo, separate effort)
- INSERT / UPDATE / DELETE analysis (deferred until SDL or query span needs it; CmdType reserved)
- Stored routine body analysis (PL/SQL-style — separate effort)
- Information_schema / performance_schema content simulation beyond schema registration
- Trigger body analysis beyond the trigger statement itself
- Event scheduler body analysis
- Optimizer, cost model, execution — omni is not building a query engine
- **MySQL 5.7 compatibility**: omni MySQL targets MySQL 8.0.34 as baseline. CTEs, window functions, CHECK enforcement, INTERSECT/EXCEPT, LATERAL, JSON_TABLE all assume 8.0+ semantics. 5.7-only constructs are out of scope; if encountered, the parser may accept them but the analyzer is not required to model their 5.7-specific behavior.

## 7. Open questions — Rev 2 status

| # | Question | Rev 1 Status | Rev 2 Resolution |
|---|---|---|---|
| Q1 | DBName/TableName vs OID for table identity | Open | **Decided: strings**. Future-optimize via interning if hot path proves it. (D1) |
| Q2 | MERGE substitution: analyze-time vs consume-time | Open | **Decided: consume-time** via `(*Query).ExpandMergeViews()`. (D5 reversed) |
| Q3 | Implicit casts: when to revisit | Open ("Phase 5") | **Revised: revisit at end of Phase 3**. (D6 updated) |
| Q4 | Reference walker location: omni vs bytebase | Open | **Decided: bytebase only**. omni ships test-only walker; no `internal/lineage/` package. (D11 confirmed) |
| Q5 | Scope refactor in Phase 1 — risk acceptable? | Open | **Decided: defer**. Phase 1a copies code with `// TODO(scope-unify)`; full unify in **Phase 4.5**. (D12 deferred) |
| Q6 | Phase 1 expression types nil — does walker need types? | Open | **Decided: no, walker doesn't need types**. Phase 2 lineage walker only reads VarExprQ.RangeIdx + AttNum |
| Q7 | `WithCheckOption` placement | Open | **Closed: already exists** as `View.CheckOption` in `mysql/catalog/table.go:89`. No IR change needed. The original plan author (me) missed this in the existing code |
| Q8 | Risks missing? | Open | **Resolved**: § 5 expanded with 7 new risks (lower_case_table_names, 8.0-only baseline, charset propagation, information_schema, ONLY_FULL_GROUP_BY, deparse driver collision, hand-write commitment) |
| Q9 | Phase 2.5 corpus vendoring | Open | **Decided: yes**. Vendored snapshot at pinned hash is a Phase 2 deliverable, not a separate phase |
| Q10 | Milestone visibility / explicit checkpoints | Open | **Decided: explicit checkpoints at Phase 1c done, Phase 2 done, Phase 4 done**. Each unblocks different downstream work |

All Phase 0 review questions resolved. Phase 1a may begin once Revision 2 is accepted.

## 8. Review process

1. Revision 1 review artifacts: this document + `mysql/catalog/query.go` (Rev 1) + `docs/plans/2026-04-09-mysql-semantic-layer-review-bundle.md`.
2. Reviewers: cc + codex (sub-agent dispatch on 2026-04-10).
3. Findings synthesized into Revision 2 (this document) and `mysql/catalog/query.go` (Rev 2). Both pass `go build` + `go vet`.
4. **Pending: human review of Revision 2 by user (and optionally cc/codex round 2)** before Phase 1a begins.

## 9. Review notes

### Round 1 — Sub-agent review (2026-04-10)

Two parallel sub-agents reviewed Revision 1: a "plan-level" reviewer focused on phase sequencing, exit criteria, and decision log; and a "code-level" reviewer focused on `mysql/catalog/query.go` type-by-type.

**Critical findings (must-fix before Phase 1, all addressed in Rev 2):**

1. **`Nullable` mis-placed on `ResolvedType`** → moved to Column (already exists). Decision D16.
2. **`FuncCallExprQ` lacked `IsAggregate`** → added. Decision D18.
3. **`Query` lacked `WindowClause` slice** → added; `WindowDefQ` doc clarifies its three roles.
4. **Q suffix not applied consistently** → unified rule (D15): all IR structs get Q except `Query` itself; enums skip Q except mechanically-forced `JoinTypeQ`.
5. **`LockInShareMode` doc was factually wrong** ("legacy synonym for FOR SHARE") → rewritten to reflect actual semantics; constraints on Tables/WaitPolicy documented.

**Plan-level findings (addressed in Rev 2):**

1. **Phase 1 oversized** → split into 1a/1b/1c with separate exit criteria.
2. **§ 2 oracle source contradiction** → rewritten to distinguish query source from lineage source.
3. **TiDB consideration** → removed entirely. omni MySQL hand-writes against the MySQL container oracle; TiDB is a different database and using it as a reference creates drift risk with no offsetting authority gain. (User decision, 2026-04-10.)
4. **Exit criteria gameable** → tightened to per-query golden IR file matching, not "corpus passes".
5. **D5 MERGE substitution timing** → reversed to consume-time.
6. **D12 scope refactor underestimated** → deferred to Phase 4.5.
7. **§ 7 Q7 WithCheckOption** → resolved (already exists in code; plan author missed it).
8. **Missing Phase 1 deliverable: oracle harness** → moved to Phase 1a day one.
9. **Missing Phase 2 deliverable: corpus snapshot vendor** → added.
10. **Missing risks** → § 5 expanded.

**Should-fix findings (also addressed in Rev 2):**

- New expression nodes: `InListExprQ`, `BetweenExprQ` added to IR.
- `RTEFunction` reinstated in IR for MySQL 8.0.19+ JSON_TABLE.
- `MySQLBaseType` typed enum replaces free-form `DataType string` (D17).
- `BoolType` package-level singleton defined; `BoolExprQ`/`NullTestExprQ`/`InListExprQ`/`BetweenExprQ` return it from `exprType()`.
- Missing SQLMode flags added: `IgnoreSpace`, `HighNotPrecedence`, `NoBackslashEscapes`.
- `ViewAlgNone` zero value introduced (Rev 1 used `ViewAlgUndefined` as zero, which conflicts with the legitimate "explicitly UNDEFINED" state).
- Per-kind field applicability table added to `RangeTableEntryQ` doc.
- `JoinExprNodeQ ↔ RTEJoin` consistency invariant documented.
- `Query.DistinctOn` and `CommonTableExprQ.CTERefCount` removed (YAGNI).
- `SubLinkExprQ` and `CastExprQ` gained `Collation` fields.
- `WindowDefQ` doc clarifies declaration vs reference vs inline-definition roles.
- `ResolvedType` doc principle stated explicitly: track only modifiers that change value range / storage / sort behavior / client-visible semantics. `INT(11)` vs `INT(5)` and `ZEROFILL` are intentionally NOT tracked because MySQL 8.0.17+ normalizes them away. `TINYINT(1)` is special-cased via its own `BaseType` value (`BaseTypeTinyIntBool`) rather than via a display-width field.

### Round 2 — Pending

Optionally re-dispatch sub-agents to verify Rev 2 fixes are correct and complete.

---

## Appendix A — File map (new files this plan creates)

```
mysql/catalog/
  query.go                   Phase 0  — IR types (DONE Rev 2)
  analyze.go                 Phase 1a — AnalyzeSelectStmt entry
  analyze_targetlist.go      Phase 1a
  analyze_expr.go            Phase 1a
  analyze_oracle_test.go     Phase 1a — oracle harness from day one
  scope.go                   Phase 1a — copied from deparse/resolver.go with TODO marker
  analyze_test.go            Phase 1a — golden IR tests
  analyze_from.go            Phase 1b — JOIN handling
  query_span_oracle_test.go  Phase 2  — view + lineage end-to-end
  testdata/bytebase_corpus/  Phase 2  — vendored snapshot
  types.go                   Phase 3  — type inference utilities
  function_types.go          Phase 3  — function return type table

mysql/deparse/
  from_query.go              Phase 4
  from_expr.go               Phase 4

mysql/diff/                  Phase 5 (new package)
  diff.go
  objects.go
  expr_diff.go
  query_diff.go
  dependency.go

mysql/sdl/                   Phase 5 (new package)
  loadsdl.go
  dumpsdl.go
  apply.go
```

## Appendix B — Reference reading order for reviewers

For cc and codex (without prior conversation context):

1. `pg/catalog/query.go` (PG IR — the model we are adapting)
2. `pg/catalog/analyze.go` first 200 lines (PG analyzer entry shape)
3. `mysql/catalog/catalog.go` + `mysql/catalog/table.go` (current MySQL L1 surface; note `Column.Nullable`, `Column.ColumnType`, `View.CheckOption`)
4. `mysql/catalog/viewcmds.go::createView` (the seam where Phase 2 plugs in)
5. `mysql/deparse/resolver.go` (existing scope handling — Phase 1a source material)
6. `pg/catalog/query_span_test.go` (the lineage walker pattern — Phase 2 will mirror this)
7. This document § 4 (the proposal)
8. This document § 9 (Rev 2 review notes — what changed and why)
