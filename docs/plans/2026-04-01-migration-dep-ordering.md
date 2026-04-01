# Migration Dependency-Driven Ordering Plan

## Problem Statement

`GenerateMigration` uses hardcoded call ordering + two heuristic hacks
(`splitFunctionOps`, `wrapColumnTypeChangesWithViewOps`) to determine operation
order. This is fragile and has known bugs:

1. **DROP ordering wrong**: `DROP TABLE` before `DROP VIEW` (masked by CASCADE)
2. **Same-type CREATE ordering missing**: table A INHERITS table B, both new — no
   guarantee B is created first
3. **`splitFunctionOps` is fragile**: uses `strings.Contains(checkExpr, funcName+"(")`
   which fails on overloads, substring matches, and schema-qualified names
4. **ALTER cascading incomplete**: changing a function signature doesn't trigger
   drop/recreate of dependent triggers

## Design Principles

1. **Catalog deps are the single source of truth** — both `from` and `to` catalogs
   already have `deps []DepEntry` populated during ProcessUtility. No string matching.
2. **Three operation phases**: DROP (reverse dep order) → CREATE (forward dep order) → METADATA
3. **Global topological sort per phase** — same Kahn + priority heap algorithm as SDL
4. **FK and CHECK constraint splitting** — pg_dump strategy: defer to POST_DATA when
   they create cycles

## Key Concepts

### Object Identity

Objects are matched between `from` and `to` catalogs by identity string:

| Object Type | Identity Format | Example |
|-------------|----------------|---------|
| Schema | `schema_name` | `public` |
| Table/View | `schema.relname` | `public.users` |
| Sequence | `schema.seqname` | `public.users_id_seq` |
| Function | `schema.name(arg_types)` | `public.is_valid(integer)` |
| Enum | `schema.typename` | `public.mood` |
| Domain | `schema.domainname` | `public.email` |
| Index | `schema.indexname` | `public.users_pkey` |
| Trigger | `schema.table.triggername` | `public.users.audit_trig` |
| Policy | `schema.table.policyname` | `public.users.user_isolation` |
| Constraint | `schema.table.conname` | `public.users.users_pkey` |

**Function identity includes argument types** — this distinguishes overloads.
`splitFunctionOps` currently uses bare `Name` for CHECK expression matching, which
can't distinguish `is_valid(integer)` from `is_valid(text)`. The OID-based deps
don't have this problem: the CHECK expression was analyzed and resolved to a
specific function OID.

### Dependency Sources

| Phase | Source Catalog | Direction | Meaning |
|-------|---------------|-----------|---------|
| DROP | `from` | Reverse | If A depends on B, drop A before B |
| CREATE | `to` | Forward | If A depends on B, create B before A |
| ALTER | Both | Both | May need auxiliary DROP/CREATE |

### Operation Categories

Each `MigrationOp` falls into one of three phases:

| Phase | Op Types | Ordering |
|-------|----------|----------|
| **PRE** (DROP) | DropView, DropTrigger, DropPolicy, DropIndex, DropConstraint, DropColumn, DropFunction, DropTable, DropSequence, DropType, DropDomain, DropEnum, DropRange, DropExtension, DropSchema, Revoke | Reverse dependency order |
| **MAIN** (CREATE + ALTER) | CreateSchema, CreateExtension, CreateType/Enum/Domain/Range, CreateSequence, CreateFunction, CreateTable, AddColumn, AlterColumn, AddConstraint (non-FK), CreateIndex, AlterFunction, AlterView, CreateView, CreateTrigger, CreatePolicy, Comment, Grant | Forward dependency order |
| **POST** (Deferred) | AddConstraint (FK), AddConstraint (CHECK with function dep cycle) | After all MAIN ops |

## Implementation Plan

### Step 1: Add dependency metadata to MigrationOp

```go
type MigrationOp struct {
    Type          MigrationOpType
    SchemaName    string
    ObjectName    string     // identity string
    SQL           string
    Warning       string
    Transactional bool
    ParentObject  string

    // New fields for dependency-driven sorting
    Phase         MigrationPhase  // Pre, Main, Post
    ObjType       byte            // 'r', 'f', 'i', 'c', 't', 'S', etc.
    ObjOID        uint32          // OID in source catalog (from for DROP, to for CREATE)
    Priority      int             // tie-breaker within phase (same as stmtPriority concept)
}

type MigrationPhase int
const (
    PhasePre  MigrationPhase = iota  // DROP operations
    PhaseMain                         // CREATE + ALTER operations
    PhasePost                         // Deferred FK/CHECK
)
```

Each `generateXxxDDL` function populates `ObjOID` from the catalog object it
generates the op for. This is the key that links ops to the dependency graph.

### Step 2: Build dependency graph from catalog deps

```go
// buildMigrationGraph constructs the dependency graph for sorting migration ops.
//
// For each op, find its dependencies in the appropriate catalog:
// - DROP ops: look up from.deps where RefOID == op.ObjOID (who depends on me?)
// - CREATE ops: look up to.deps where ObjOID == op.ObjOID (who do I depend on?)
//
// Then match those OIDs to other ops in the same plan.
func buildMigrationGraph(from, to *Catalog, ops []MigrationOp) (adj [][]int, inDegree []int)
```

The matching between OIDs and ops works via two index maps:
- `fromOIDToOp`: maps (objType, OID) → index of DROP op
- `toOIDToOp`: maps (objType, OID) → index of CREATE/ALTER op

### Step 3: Replace GenerateMigration sorting

```go
func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
    // Phase 1: Generate all ops (unordered)
    var allOps []MigrationOp
    allOps = append(allOps, generateSchemaDDL(from, to, diff)...)
    allOps = append(allOps, generateExtensionDDL(from, to, diff)...)
    // ... all generateXxxDDL calls ...
    allOps = append(allOps, generateFunctionDDL(from, to, diff)...)  // no split
    allOps = append(allOps, generateTableDDL(from, to, diff)...)
    allOps = append(allOps, generateColumnDDL(from, to, diff)...)
    allOps = append(allOps, generateConstraintDDL(from, to, diff)...)
    allOps = append(allOps, generateViewDDL(from, to, diff)...)
    allOps = append(allOps, generateIndexDDL(from, to, diff)...)
    allOps = append(allOps, generateTriggerDDL(from, to, diff)...)
    allOps = append(allOps, generatePolicyDDL(from, to, diff)...)
    allOps = append(allOps, generateCommentDDL(from, to, diff)...)
    allOps = append(allOps, generateGrantDDL(from, to, diff)...)

    // Phase 2: Classify ops into phases
    classifyOps(allOps)

    // Phase 3: Build dependency graph and sort
    sorted := sortMigrationOps(from, to, allOps)

    return &MigrationPlan{Ops: sorted}
}
```

### Step 4: Implement sortMigrationOps

```go
func sortMigrationOps(from, to *Catalog, ops []MigrationOp) []MigrationOp {
    // Separate into three phases
    var preOps, mainOps, postOps []MigrationOp
    for _, op := range ops {
        switch op.Phase {
        case PhasePre:
            preOps = append(preOps, op)
        case PhaseMain:
            mainOps = append(mainOps, op)
        case PhasePost:
            postOps = append(postOps, op)
        }
    }

    // Sort PRE phase: reverse dependency order using from.deps
    sortedPre := topoSortOps(from, preOps, true /*reverse*/)

    // Sort MAIN phase: forward dependency order using to.deps
    sortedMain := topoSortOps(to, mainOps, false /*forward*/)

    // POST phase: FK/deferred constraints, sorted by name for determinism
    sortedPost := sortByName(postOps)

    // Concatenate: DROP first, then CREATE/ALTER, then deferred
    return append(append(sortedPre, sortedMain...), sortedPost...)
}
```

### Step 5: Implement topoSortOps (reuse SDL pattern)

```go
// topoSortOps sorts migration ops using Kahn's algorithm with priority tie-break.
// If reverse=true, edges are reversed (for DROP ordering).
//
// Uses catalog.deps to build the dependency graph:
// - forward (CREATE): if dep says A depends on B, B must come before A
// - reverse (DROP): if dep says A depends on B, A must come before B (drop dependent first)
func topoSortOps(c *Catalog, ops []MigrationOp, reverse bool) []MigrationOp {
    // Build OID → op index map
    oidToIdx := map[depKey]int{}  // (objType, objOID) → op index
    for i, op := range ops {
        if op.ObjOID != 0 {
            oidToIdx[depKey{op.ObjType, op.ObjOID}] = i
        }
    }

    // Build adjacency from catalog deps
    n := len(ops)
    adj := make([][]int, n)
    inDegree := make([]int, n)

    for _, d := range c.deps {
        if d.DepType == DepInternal {
            continue // internal deps are inseparable, handled implicitly
        }

        // Find the op for the dependent object and the referenced object
        depIdx, depOK := oidToIdx[depKey{d.ObjType, d.ObjOID}]
        refIdx, refOK := oidToIdx[depKey{d.RefType, d.RefOID}]

        if !depOK || !refOK {
            continue // one side not in our op set
        }
        if depIdx == refIdx {
            continue // self-reference
        }

        if reverse {
            // DROP: dependent (depIdx) must come BEFORE referenced (refIdx)
            adj[depIdx] = append(adj[depIdx], refIdx)
            inDegree[refIdx]++
        } else {
            // CREATE: referenced (refIdx) must come BEFORE dependent (depIdx)
            adj[refIdx] = append(adj[refIdx], depIdx)
            inDegree[depIdx]++
        }
    }

    // Kahn's algorithm with min-heap (priority, index) — same as SDL
    // ... (reuse pqHeap from sdl.go, or extract to shared utility)
}
```

### Step 6: Delete heuristic hacks

- Remove `splitFunctionOps` — dependency graph handles function→table ordering
- Remove `wrapColumnTypeChangesWithViewOps` — dependency graph handles view→column ordering
- Remove per-type internal sorting in each `generateXxxDDL` — global sort handles it

### Step 7: Handle ALTER edge cases

ALTER operations that change object identity (e.g., function signature change)
are already handled as DROP + CREATE pairs by `generateFunctionDDL`. The dependency
graph naturally orders them:

1. DROP trigger (depends on old function via `from.deps`)
2. DROP old function
3. CREATE new function
4. CREATE trigger (depends on new function via `to.deps`)

For column type changes that break dependent views:
- The diff generates ALTER COLUMN TYPE ops
- Dependent views need to be dropped and recreated
- Instead of `wrapColumnTypeChangesWithViewOps`, the dependency graph handles this:
  - View depends on table column (via `to.deps`)
  - If column type changed, the view's DiffModify generates DROP VIEW + CREATE VIEW
  - DROP VIEW is in PhasePre, ALTER COLUMN is in PhaseMain, CREATE VIEW is in PhaseMain
  - Dependency order: DROP VIEW → ALTER COLUMN → CREATE VIEW

**Caveat**: This requires that column type changes are detected as DiffModify on
the view (since the view definition hasn't changed, only its dependency's type).
Currently `wrapColumnTypeChangesWithViewOps` explicitly handles this. With the
new approach, the diff engine needs to detect "implicit view modification" when a
dependency's column type changes. This is the trickiest part.

**Proposed solution**: Keep a simplified version of `wrapColumnTypeChangesWithViewOps`
as a **pre-processing step** that injects additional DROP VIEW + CREATE VIEW ops
into the allOps list when column type changes are detected. The dependency-driven
sort then orders them correctly. This separates "what ops are needed" (pre-processing)
from "what order" (dependency sort).

## Migration from Current to New Architecture

To minimize risk, implement incrementally:

### Phase A: Add ObjOID to MigrationOp (non-breaking)
- Add `ObjOID`, `ObjType`, `Phase`, `Priority` fields to `MigrationOp`
- Each `generateXxxDDL` populates them
- No change to sorting logic yet
- All existing tests pass

### Phase B: Implement dependency-driven sort (parallel path)
- Add `sortMigrationOps` and `topoSortOps`
- Wire into `GenerateMigration` behind a flag or as the new default
- Keep old path as fallback for comparison
- Add tests for all problematic scenarios

### Phase C: Remove heuristic hacks
- Delete `splitFunctionOps`
- Simplify `wrapColumnTypeChangesWithViewOps` to only inject ops, not sort
- Remove per-type internal sorting
- Comprehensive regression testing

## Priority Mapping (tie-breaker for topo sort)

| Priority | Op Types | Rationale |
|----------|----------|-----------|
| 0 | CreateSchema, DropSchema | Schema is the container |
| 1 | CreateExtension, DropExtension | Extensions provide types/functions |
| 2 | CreateType (Enum/Domain/Range/Composite), DropType | Types before columns |
| 3 | CreateSequence, DropSequence | Sequences before DEFAULT |
| 4 | CreateFunction, DropFunction, AlterFunction | Functions before tables (CHECK/DEFAULT) |
| 5 | CreateTable, DropTable | Tables before views/indexes |
| 6 | AddColumn, DropColumn, AlterColumn | Columns within tables |
| 7 | AddConstraint (non-FK), DropConstraint | Constraints after columns |
| 8 | CreateView, DropView, AlterView | Views after tables |
| 9 | CreateIndex, DropIndex | Indexes after tables |
| 10 | CreateTrigger, DropTrigger | Triggers after functions + tables |
| 11 | CreatePolicy, DropPolicy, AlterPolicy | Policies after tables |
| 12 | Comment, Grant, Revoke | Metadata last |
| 99 | AddConstraint (FK) | POST phase |

Note: Functions at priority 4 (before tables at 5) is the **default** tie-breaker.
When there's no dependency between a function and a table, functions naturally
come first. When a function depends on a table (RETURNS SETOF table), the
dependency edge overrides priority and places the function after the table.
This eliminates the need for `splitFunctionOps` entirely.

## Function Name vs Identity

The current `splitFunctionOps` matches by bare function name:
```go
funcName := fn.To.Name  // e.g. "is_valid"
strings.Contains(c.CheckExpr, funcName+"(")
```

This has problems:
1. **Overload ambiguity**: `is_valid(integer)` and `is_valid(text)` both match
2. **Substring false positive**: function `is` matches `is_valid(`
3. **Schema mismatch**: unqualified name in expression won't match qualified name

The new approach uses **OID-based deps** from the catalog:
- When LoadSDL/LoadSQL processes `CREATE TABLE ... CHECK(is_valid(qty))`,
  the analyzer resolves `is_valid` to a specific function OID
- `recordDependencyOnSingleRelExpr` records `constraint OID → function OID`
- The migration graph uses these OID links — no string matching needed
- Overloads are correctly distinguished because each has a unique OID

## Review Issues and Resolutions

### Issue 1 (High): Constraint deps don't map to table ops

Catalog deps record `constraint OID → function OID`, but migration ops operate
at the table level (CREATE TABLE includes inline constraints). The OID-to-op
mapping would fail because the constraint OID doesn't match any op's ObjOID.

**Resolution**: Dep lifting. When building the migration graph, for each dep
involving a constraint OID (`ObjType='c'`), look up the owning relation via
`Constraint.RelOID` and map the dep to the table's op instead. Same for index
deps (`ObjType='i'`) — lift to owning relation via `Index.RelOID`.

```go
// liftDepToOp maps a dep entry's ObjType/ObjOID to the migration op that
// "owns" it. Constraints and indexes are lifted to their parent table.
func liftDepToOp(c *Catalog, d DepEntry, oidToIdx map[depKey][]int) []int {
    // Direct match first
    if idxs, ok := oidToIdx[depKey{d.ObjType, d.ObjOID}]; ok {
        return idxs
    }
    // Lift constraint → owning table
    if d.ObjType == 'c' {
        if con, ok := c.constraints[d.ObjOID]; ok {
            return oidToIdx[depKey{'r', con.RelOID}]
        }
    }
    // Lift index → owning table
    if d.ObjType == 'i' {
        if idx, ok := c.indexByOID[d.ObjOID]; ok {
            return oidToIdx[depKey{'r', idx.RelOID}]
        }
    }
    return nil
}
```

### Issue 2 (High): Multiple ops per OID

A DiffModify on a function with signature change produces both a DROP op
(from OID) and a CREATE op (to OID). A table may have column ops, constraint
ops, and the table create op all referencing the same relation OID.

**Resolution**: Change `oidToIdx` from `map[depKey]int` to `map[depKey][]int`.
When a dep edge is resolved, it creates edges to ALL ops matching that OID.
This is correct because if two ops share an OID (e.g., DROP TABLE and DROP INDEX
on the same relation), they all participate in the dependency ordering.

### Issue 3 (Medium): Trigger cascade on function signature change

When function `f(int)` changes to `f(int, text)`, triggers referencing the old
function need to be dropped and recreated. The diff engine currently does NOT
detect this — it only compares triggers by name/timing/events.

**Resolution (Phase B)**: When generating trigger ops, check if the trigger's
function identity has changed between `from` and `to`. If the function was
DiffDrop+DiffAdd (signature change), any trigger referencing the old function
identity should generate a DiffModify (DROP + CREATE). This check goes in
`diffTriggers` or `generateTriggerDDL`.

### Issue 4 (Medium): Column ops have no own OID

ALTER COLUMN ops don't have a dedicated OID — they operate on a relation's column.

**Resolution**: Key column ops by `(ObjType='r', ObjOID=relOID)`. All column
ops on the same table share the table's OID. This means column ops are ordered
relative to other object types (table comes before its dependent views) but
column ops on the same table are ordered by the generateColumnDDL internal logic
(which knows the correct ADD before ALTER before DROP order within a table).

### Issue 5 (Medium): Cycle detection

The topoSortOps pseudocode has no cycle handling.

**Resolution**: After Kahn's algorithm, if `len(sorted) < len(ops)`, there's a
cycle. Handle by:
1. Identify the ops involved in the cycle
2. If any are AddConstraint (CHECK or FK), move them to PhasePost (deferred)
3. Re-run sort without the deferred ops
4. If cycle persists with no deferrable ops, return error

### Issue 6 (Low): Partition DDL omitted

**Resolution**: Add `generatePartitionDDL` to the allOps collection and
`migration_partition.go` to the files-modified list.

### Issue 7: wrapColumnTypeChangesWithViewOps replacement

The plan keeps a simplified version as a pre-processing step.

**Refined approach**: Use `from.deps` to find views that depend on a table being
modified (not string matching on view definition). When a column type change is
detected, check `from.deps` for any view OIDs that depend on that relation OID.
Inject synthetic DROP VIEW (PhasePre) + CREATE VIEW (PhaseMain) ops for those
views. The topological sort then handles the ordering.

## Testing Strategy

### New test scenarios to add:

1. **Cross-type DROP ordering**: drop table + drop dependent view → view dropped first
2. **Same-type CREATE ordering**: INHERITS parent + child both new → parent first
3. **Function overload with CHECK**: two overloads, only one referenced by CHECK
4. **Function signature change + trigger**: trigger dropped before function, recreated after
5. **Column type change + view**: view dropped, column altered, view recreated
6. **Domain used in table column**: domain created before table
7. **Sequence in DEFAULT + table**: sequence before table
8. **Chain: function → CHECK → table → view**: correct 4-step ordering
9. **Everything at once**: schema + enum + function + table + view + trigger + index + FK

### Regression: all existing `TestMigrationOrdering` tests must pass unchanged.

## Files Modified

| File | Changes |
|------|---------|
| `migration.go` | Rewrite `GenerateMigration`, add `sortMigrationOps`, delete `splitFunctionOps`, simplify `wrapColumnTypeChangesWithViewOps` |
| `migration_table.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_view.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_function.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_constraint.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_column.go` | Add ObjOID/Phase to ops |
| `migration_index.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_trigger.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_policy.go` | Add ObjOID/Phase to ops, remove internal sort |
| `migration_schema.go` | Add ObjOID/Phase to ops |
| `migration_extension.go` | Add ObjOID/Phase to ops |
| `migration_enum.go` | Add ObjOID/Phase to ops |
| `migration_domain.go` | Add ObjOID/Phase to ops |
| `migration_range.go` | Add ObjOID/Phase to ops |
| `migration_sequence.go` | Add ObjOID/Phase to ops |
| `migration_comment.go` | Add Phase to ops |
| `migration_grant.go` | Add Phase to ops |
| `migration_ordering_test.go` | Add new dependency ordering tests |
| `sdl.go` | Extract `pqHeap` to shared utility (optional) |
