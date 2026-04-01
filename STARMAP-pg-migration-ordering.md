# Execution Starmap: PG Migration Dependency-Driven Ordering

> Source scenarios: SCENARIOS-pg-migration-ordering.md
> Implementation plan: docs/plans/2026-04-01-migration-dep-ordering.md
> Generated: 2026-04-01

---

## Step 1: File Targets Per Section

### Phase 1: Metadata Foundation

| Section | Files Modified |
|---------|---------------|
| 1.1 MigrationOp Metadata Population | `pg/catalog/migration.go` (add Phase/ObjOID/ObjType/Priority fields to MigrationOp struct), `pg/catalog/migration_table.go`, `pg/catalog/migration_view.go`, `pg/catalog/migration_function.go`, `pg/catalog/migration_constraint.go`, `pg/catalog/migration_column.go`, `pg/catalog/migration_index.go`, `pg/catalog/migration_trigger.go`, `pg/catalog/migration_policy.go`, `pg/catalog/migration_schema.go`, `pg/catalog/migration_extension.go`, `pg/catalog/migration_enum.go`, `pg/catalog/migration_domain.go`, `pg/catalog/migration_range.go`, `pg/catalog/migration_sequence.go`, `pg/catalog/migration_comment.go`, `pg/catalog/migration_grant.go`, `pg/catalog/migration_partition.go` |
| 1.2 Phase and Priority Classification | `pg/catalog/migration.go` (add MigrationPhase type, constants, classifyOps helper; or inline into 1.1's struct literal changes), all `migration_*.go` files (same as 1.1 — the Phase/Priority values are set in the same MigrationOp literals) |

**Key observation**: 1.1 and 1.2 modify the exact same set of files at the exact same MigrationOp literal sites. They are logically separable but physically inseparable — the ObjOID field and the Phase/Priority fields are added to the same struct literals in the same code locations. They MUST be implemented as a single unit.

### Phase 2: Sort Engine

| Section | Files Modified |
|---------|---------------|
| 2.1 Forward Dependency Sorting | `pg/catalog/migration.go` (add `sortMigrationOps`, `topoSortOps`, `depKey` type, OID-to-op index building; reuse `pqHeap` from `sdl.go` or extract shared version), `pg/catalog/migration_ordering_test.go` (new forward-dep test cases) |
| 2.2 Reverse Dependency Sorting | `pg/catalog/migration.go` (extend `topoSortOps` with `reverse` parameter — already designed in 2.1), `pg/catalog/migration_ordering_test.go` (new reverse-dep test cases) |
| 2.3 Dependency Lifting | `pg/catalog/migration.go` (add `liftDepToOp` function, integrate into graph building in `topoSortOps`), `pg/catalog/migration_ordering_test.go` (dep-lifting test cases) |

**Key observation**: 2.1, 2.2, 2.3 all modify `migration.go` and `migration_ordering_test.go`. They build on each other: 2.2 extends 2.1's `topoSortOps`, 2.3 adds lifting logic called by the graph builder from 2.1/2.2. Strictly sequential.

### Phase 3: Heuristic Removal

| Section | Files Modified |
|---------|---------------|
| 3.1 Replace splitFunctionOps | `pg/catalog/migration.go` (delete `splitFunctionOps` function ~55 lines; modify `GenerateMigration` to remove the early/late split call and just append funcOps normally), `pg/catalog/migration_ordering_test.go` (new tests proving dep-graph handles function ordering) |
| 3.2 Replace wrapColumnTypeChangesWithViewOps | `pg/catalog/migration.go` (simplify `wrapColumnTypeChangesWithViewOps` to only inject ops without sorting, or replace with dep-based view detection using `from.deps`), `pg/catalog/migration_ordering_test.go` (column-type-change + view tests) |

**Key observation**: 3.1 and 3.2 both modify `migration.go` (the `GenerateMigration` function body and removing/simplifying different heuristic functions). Sequential.

### Phase 4: Edge Cases and Mixed Operations

| Section | Files Modified |
|---------|---------------|
| 4.1 Mixed CREATE + DROP Scenarios | `pg/catalog/migration_ordering_test.go` (new test cases only — no production code changes expected; these validate the integrated system) |
| 4.2 Cycle Handling and FK Deferral | `pg/catalog/migration.go` (add cycle detection to `topoSortOps`: detect remaining nodes after Kahn's, move deferrable ops to PhasePost, re-run), `pg/catalog/migration_ordering_test.go` (cycle test cases) |

**Key observation**: 4.1 is test-only. 4.2 modifies `migration.go`. But 4.1's tests depend on the full integrated system from Phases 1-3 being complete, so they cannot run until Phase 3 is done. 4.2 adds cycle handling to `topoSortOps` in `migration.go`.

---

## Step 2: Independence Classification

### Shared-file matrix (files modified by each section)

```
                    migration.go  migration_*.go (17 files)  migration_ordering_test.go  sdl.go
1.1+1.2 Metadata    WRITE         WRITE                      —                           —
2.1 Forward Sort    WRITE         —                          WRITE                       READ
2.2 Reverse Sort    WRITE         —                          WRITE                       —
2.3 Dep Lifting     WRITE         —                          WRITE                       —
3.1 splitFunc       WRITE         —                          WRITE                       —
3.2 wrapColType     WRITE         —                          WRITE                       —
4.1 Mixed Tests     —             —                          WRITE                       —
4.2 Cycle Handle    WRITE         —                          WRITE                       —
```

### Pairwise assessment

| Pair | Semantic | Change-surface | Proof | Verdict |
|------|----------|----------------|-------|---------|
| 1.1+1.2 vs 2.1 | 2.1 requires metadata from 1.1+1.2 | Both write `migration.go` | 2.1 tests need ObjOID populated | **SEQUENTIAL: 1 before 2** |
| 2.1 vs 2.2 | 2.2 extends `topoSortOps` from 2.1 | Both write `migration.go` | 2.2 tests need forward sort working | **SEQUENTIAL: 2.1 before 2.2** |
| 2.2 vs 2.3 | 2.3 adds lifting into graph builder from 2.1/2.2 | Both write `migration.go` | 2.3 tests need sort working | **SEQUENTIAL: 2.2 before 2.3** |
| 2.3 vs 3.1 | 3.1 removes heuristic, relies on dep sort from Phase 2 | Both write `migration.go` | 3.1 tests need complete sort engine | **SEQUENTIAL: 2.3 before 3.1** |
| 3.1 vs 3.2 | Independent heuristics, but both modify `GenerateMigration` body | Both write `migration.go` | — | **SEQUENTIAL: 3.1 before 3.2** |
| 3.2 vs 4.1 | 4.1 tests integrated system after heuristic removal | 4.1 only writes test file | 4.1 tests need Phase 3 complete | **SEQUENTIAL: 3.2 before 4.1** |
| 4.1 vs 4.2 | 4.2 adds cycle handling to `migration.go` | 4.1 is test-only, 4.2 writes `migration.go` | 4.1 tests may expose issues 4.2 needs to fix | **SEQUENTIAL: 4.1 before 4.2** |

**Conclusion: The entire starmap is strictly sequential.** Every section shares `migration.go` with its predecessor and/or has a semantic dependency. No parallelism is possible.

---

## Step 3: Execution Shape

```
SERIAL PIPELINE (single worker, 8 steps)

  ┌──────────────────────────────────────────────────────────────┐
  │ Phase 1: Metadata Foundation                                 │
  │   Step 1: [1.1+1.2] MigrationOp Metadata + Phase/Priority   │
  │           → proof: go build + existing tests pass            │
  ├──────────────────────────────────────────────────────────────┤
  │ Phase 2: Sort Engine                                         │
  │   Step 2: [2.1] Forward Dependency Sorting                   │
  │           → proof: go build + new forward-dep tests          │
  │   Step 3: [2.2] Reverse Dependency Sorting                   │
  │           → proof: go build + new reverse-dep tests          │
  │   Step 4: [2.3] Dependency Lifting                           │
  │           → proof: go build + new lifting tests              │
  ├──────────────────────────────────────────────────────────────┤
  │ Phase 3: Heuristic Removal                                   │
  │   Step 5: [3.1] Replace splitFunctionOps                     │
  │           → proof: go build + function ordering tests        │
  │   Step 6: [3.2] Replace wrapColumnTypeChangesWithViewOps     │
  │           → proof: go build + column-type-change tests       │
  ├──────────────────────────────────────────────────────────────┤
  │ Phase 4: Edge Cases                                          │
  │   Step 7: [4.1] Mixed CREATE + DROP Scenarios                │
  │           → proof: new mixed-op tests pass                   │
  │   Step 8: [4.2] Cycle Handling and FK Deferral               │
  │           → proof: go build + cycle tests pass               │
  └──────────────────────────────────────────────────────────────┘
```

---

## Step 4: Proof Checkpoints

### Section-local proofs

| Step | Section | Proof command | What it verifies |
|------|---------|--------------|------------------|
| 1 | 1.1+1.2 | `cd pg/catalog && go build ./...` | New fields compile; no type errors in 17 migration_*.go files |
| 1 | 1.1+1.2 | `cd pg/catalog && go test -run TestMigrationOrdering -count=1` | Existing 6 ordering tests still pass (metadata is additive, no behavior change) |
| 1 | 1.1+1.2 | `cd pg/catalog && go test -run TestMigrationRoundtrip -count=1` | Roundtrip tests still pass (if they exist) |
| 2 | 2.1 | `cd pg/catalog && go test -run TestMigrationOrdering/forward -count=1` | Forward dep scenarios: function before table (CHECK), enum before table, sequence before table, view chain ordering |
| 3 | 2.2 | `cd pg/catalog && go test -run TestMigrationOrdering/reverse -count=1` | Reverse dep scenarios: view dropped before table, trigger before function, index before table |
| 4 | 2.3 | `cd pg/catalog && go test -run TestMigrationOrdering/lifting -count=1` | Dep lifting: CHECK constraint lifted to table, expression index lifted, column type lifted |
| 5 | 3.1 | `cd pg/catalog && go test -run TestMigrationOrdering -count=1` | All ordering tests pass WITHOUT splitFunctionOps; function overload precision test passes |
| 6 | 3.2 | `cd pg/catalog && go test -run TestMigrationOrdering -count=1` | Column type change + view tests pass with simplified wrapColumnTypeChangesWithViewOps |
| 7 | 4.1 | `cd pg/catalog && go test -run TestMigrationOrdering/mixed -count=1` | Mixed create+drop scenarios all pass |
| 8 | 4.2 | `cd pg/catalog && go test -run TestMigrationOrdering/cycle -count=1` | Self-ref FK, circular FK, three-way cycle, unresolvable cycle error |

### Phase-end global proofs

| After Phase | Proof command | What it verifies |
|-------------|--------------|------------------|
| Phase 1 | `cd pg/catalog && go test ./... -count=1` | Full package test suite passes; metadata addition is non-breaking |
| Phase 2 | `cd pg/catalog && go test ./... -count=1` | Sort engine integrated; all existing + new tests pass |
| Phase 3 | `cd pg/catalog && go test ./... -count=1` | Heuristics removed; no regression in full suite |
| Phase 4 | `cd pg/catalog && go test ./... -count=1` | Complete implementation; full regression green |

### Final proof (post-Phase 4)

```
cd pg && go test ./... -count=1
```

Verifies no breakage in any package that imports `pg/catalog`.

---

## Step 5: Execution Contract

```yaml
version: 1
name: pg-migration-dep-ordering
scenario_file: SCENARIOS-pg-migration-ordering.md
plan_file: docs/plans/2026-04-01-migration-dep-ordering.md

execution:
  shape: serial
  reason: >
    All 8 sections share pg/catalog/migration.go as the central orchestration
    file. Each section either adds fields/types to it (Phase 1), adds new
    functions to it (Phase 2), modifies existing functions in it (Phase 3),
    or extends functions in it (Phase 4). No parallelism is safe.

  steps:
    - id: "1.1+1.2"
      name: "MigrationOp Metadata + Phase/Priority Classification"
      sections: ["1.1", "1.2"]
      reason: "1.1 and 1.2 modify the same struct literals in the same files; physically inseparable"
      files:
        - pg/catalog/migration.go           # MigrationOp struct: add Phase, ObjType, ObjOID, Priority fields; add MigrationPhase type+constants
        - pg/catalog/migration_table.go     # generateTableDDL: populate ObjOID=rel.OID, ObjType='r', Phase, Priority in every MigrationOp literal
        - pg/catalog/migration_view.go      # generateViewDDL: populate ObjOID=rel.OID, ObjType='r', Phase, Priority
        - pg/catalog/migration_function.go  # generateFunctionDDL: populate ObjOID=proc.OID, ObjType='f', Phase, Priority
        - pg/catalog/migration_constraint.go # generateConstraintDDL: populate ObjOID=con.OID, ObjType='c', Phase, Priority (FK gets PhasePost)
        - pg/catalog/migration_column.go    # generateColumnDDL: populate ObjOID=rel.OID (parent table), ObjType='r', Phase, Priority
        - pg/catalog/migration_index.go     # generateIndexDDL: populate ObjOID=idx.OID, ObjType='i', Phase, Priority
        - pg/catalog/migration_trigger.go   # generateTriggerDDL: populate ObjOID=trig.OID, ObjType='T' (or trigger-specific), Phase, Priority
        - pg/catalog/migration_policy.go    # generatePolicyDDL: populate ObjOID, Phase, Priority
        - pg/catalog/migration_schema.go    # generateSchemaDDL: populate ObjOID=schema.OID, Phase, Priority=0
        - pg/catalog/migration_extension.go # generateExtensionDDL: populate ObjOID, Phase, Priority=1
        - pg/catalog/migration_enum.go      # generateEnumDDL: populate ObjOID=type.OID, ObjType='t', Phase, Priority=2
        - pg/catalog/migration_domain.go    # generateDomainDDL: populate ObjOID=type.OID, ObjType='t', Phase, Priority=2
        - pg/catalog/migration_range.go     # generateRangeDDL: populate ObjOID=type.OID, ObjType='t', Phase, Priority=2
        - pg/catalog/migration_sequence.go  # generateSequenceDDL: populate ObjOID=seq.OID, ObjType='S', Phase, Priority=3
        - pg/catalog/migration_comment.go   # generateCommentDDL: populate Phase=PhaseMain, Priority=12
        - pg/catalog/migration_grant.go     # generateGrantDDL: populate Phase=PhaseMain, Priority=12
        - pg/catalog/migration_partition.go # generatePartitionDDL: populate ObjOID, Phase, Priority
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"
          - command: "cd pg/catalog && go test -run TestMigrationConstraint -count=1"
            expect: "PASS"
        global:
          - command: "cd pg/catalog && go test ./... -count=1"
            expect: "PASS"

    - id: "2.1"
      name: "Forward Dependency Sorting (CREATE ordering)"
      sections: ["2.1"]
      depends_on: ["1.1+1.2"]
      files:
        - pg/catalog/migration.go            # Add depKey type, topoSortOps func, sortMigrationOps func, wire into GenerateMigration
        - pg/catalog/migration_ordering_test.go  # New subtests: forward dep scenarios (function→CHECK→table, enum→table, sequence→table, view chain, trigger deps)
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"

    - id: "2.2"
      name: "Reverse Dependency Sorting (DROP ordering)"
      sections: ["2.2"]
      depends_on: ["2.1"]
      files:
        - pg/catalog/migration.go            # Extend topoSortOps with reverse=true path
        - pg/catalog/migration_ordering_test.go  # New subtests: reverse dep scenarios (drop view before table, drop trigger before function, drop indexes before table)
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"

    - id: "2.3"
      name: "Dependency Lifting"
      sections: ["2.3"]
      depends_on: ["2.2"]
      files:
        - pg/catalog/migration.go            # Add liftDepToOp func; integrate into graph builder in topoSortOps
        - pg/catalog/migration_ordering_test.go  # New subtests: CHECK→function lifted to table, expression index→function, column type→type, dep not in op set (graceful skip), op with zero ObjOID
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"
        global:
          - command: "cd pg/catalog && go test ./... -count=1"
            expect: "PASS"

    - id: "3.1"
      name: "Replace splitFunctionOps"
      sections: ["3.1"]
      depends_on: ["2.3"]
      files:
        - pg/catalog/migration.go            # Delete splitFunctionOps (~55 lines); modify GenerateMigration to append funcOps without early/late split
        - pg/catalog/migration_ordering_test.go  # New subtests: overload precision (is_valid(integer) vs is_valid(text)), function not referenced by table, RETURNS SETOF table
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"

    - id: "3.2"
      name: "Replace wrapColumnTypeChangesWithViewOps"
      sections: ["3.2"]
      depends_on: ["3.1"]
      files:
        - pg/catalog/migration.go            # Simplify wrapColumnTypeChangesWithViewOps to use from.deps for view detection (not string matching); keep op injection, remove sorting logic
        - pg/catalog/migration_ordering_test.go  # New subtests: column type change + dependent view, chain of dependent views, column type change with no views
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"
        global:
          - command: "cd pg/catalog && go test ./... -count=1"
            expect: "PASS"

    - id: "4.1"
      name: "Mixed CREATE + DROP Scenarios"
      sections: ["4.1"]
      depends_on: ["3.2"]
      files:
        - pg/catalog/migration_ordering_test.go  # New subtests: replace table with dependent view, function signature change with trigger, add+drop+add view, enum value addition, drop+create same name, column type change + new function CHECK
      proof:
        local:
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"

    - id: "4.2"
      name: "Cycle Handling and FK Deferral"
      sections: ["4.2"]
      depends_on: ["4.1"]
      files:
        - pg/catalog/migration.go            # Add cycle detection after Kahn's: identify remaining ops, move deferrable (FK/CHECK) to PhasePost, re-run, error on unresolvable
        - pg/catalog/migration_ordering_test.go  # New subtests: self-ref FK, circular FK (2 tables), three-way FK cycle, CHECK cycle deferral, no false deferrals, unresolvable cycle error, FK deferred ops ordered by name
      proof:
        local:
          - command: "cd pg/catalog && go build ./..."
            expect: "exit 0"
          - command: "cd pg/catalog && go test -run TestMigrationOrdering -count=1"
            expect: "PASS"
        global:
          - command: "cd pg/catalog && go test ./... -count=1"
            expect: "PASS"

  final_proof:
    - command: "cd pg && go test ./... -count=1"
      expect: "PASS"
      reason: "Full pg package tree — catches any breakage in consumers of pg/catalog"

constraints:
  - "No parallelism: all steps share pg/catalog/migration.go"
  - "Phase boundaries are hard gates: Phase N+1 cannot start until Phase N global proof passes"
  - "Within a phase, steps are sequential due to migration.go contention"
  - "Step 1.1+1.2 is merged: the ObjOID and Phase/Priority fields target the same struct literals"
  - "pqHeap from pg/catalog/sdl.go is reused (READ only) — no extraction needed unless desired"
  - "Existing TestMigrationOrdering (6 subtests) must pass at every step"
  - "Existing TestMigrationConstraint tests must pass at every step"

notes:
  - "migration_ordering_test.go is the only test file modified and it accumulates new subtests at each step"
  - "migration.go is the bottleneck: it is modified in 7 of 8 steps (all except 4.1)"
  - "The 17 migration_*.go files are only modified in Step 1 (metadata population); after that they are stable"
  - "sdl.go is READ-only (pqHeap reuse); it is never modified by this starmap"
  - "migration_constraint_test.go contains the opsSQL helper used by ordering tests"
```

---

## Summary

**Shape**: Fully serial, 8 steps, single worker.

**Bottleneck**: `pg/catalog/migration.go` is written by 7 of 8 steps. This file contains the `MigrationOp` struct, `GenerateMigration`, `splitFunctionOps`, `wrapColumnTypeChangesWithViewOps`, and will contain the new `sortMigrationOps`, `topoSortOps`, `liftDepToOp`, and cycle detection logic.

**Why no parallelism**: Every adjacent pair of steps either (a) shares `migration.go` as a write target, or (b) has a semantic dependency where the later step's tests require the earlier step's production code. In most cases both conditions hold simultaneously.

**Critical path**: Steps 1 through 8 in sequence. Estimated touch count: ~18 files in Step 1, then 2 files per step for Steps 2-8.
