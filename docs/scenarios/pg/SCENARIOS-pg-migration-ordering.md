# PG Migration Dependency-Driven Ordering Scenarios

> Goal: Replace hardcoded migration ordering with dependency-driven topological sort using catalog deps
> Verification: LoadSDL before → LoadSDL after → Diff → GenerateMigration → verify op ordering in tests
> Reference sources: pg_dump ordering model, pg/catalog/depend.go (DepEntry), implementation plan at docs/plans/2026-04-01-migration-dep-ordering.md

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Metadata Foundation

### 1.1 MigrationOp Metadata Population

Each MigrationOp must carry ObjOID, ObjType, Phase, and Priority so the sort engine can use them.

- [x] CreateTable op has ObjOID set to the relation OID from `to` catalog
- [x] DropTable op has ObjOID set to the relation OID from `from` catalog
- [x] CreateFunction op has ObjOID set to the UserProc OID from `to` catalog
- [x] DropFunction op has ObjOID set to the UserProc OID from `from` catalog
- [x] CreateView op has ObjOID set to the relation OID from `to` catalog
- [x] DropView op has ObjOID set to the relation OID from `from` catalog
- [x] CreateIndex op has ObjOID set to the index OID from `to` catalog
- [x] DropIndex op has ObjOID set to the index OID from `from` catalog
- [x] CreateTrigger op has ObjOID set to the trigger OID from `to` catalog
- [x] DropTrigger op has ObjOID set to the trigger OID from `from` catalog
- [x] CreateSequence op has ObjOID set to the sequence OID from `to` catalog
- [x] DropSequence op has ObjOID set to the sequence OID from `from` catalog
- [x] CreateSchema op has ObjOID set to the schema OID from `to` catalog
- [x] CreateType (Enum/Domain/Range) op has ObjOID set to the type OID
- [x] CreatePolicy op has ObjOID set to the policy OID from `to` catalog
- [x] AddConstraint op has ObjOID set to the constraint OID from `to` catalog
- [x] AlterColumn op has ObjOID set to owning relation OID from `to` catalog

### 1.2 Phase and Priority Classification

Each op must be classified into PhasePre (DROP), PhaseMain (CREATE/ALTER), or PhasePost (deferred).

- [x] All Drop* ops classified as PhasePre
- [x] All Create* ops classified as PhaseMain
- [x] AlterColumn ops classified as PhaseMain
- [x] AlterFunction ops classified as PhaseMain
- [x] AlterView ops classified as PhaseMain
- [x] AlterSequence ops classified as PhaseMain
- [x] AddConstraint (FK) classified as PhasePost
- [x] AddConstraint (non-FK) classified as PhaseMain
- [x] Comment ops classified as PhaseMain (metadata, high priority)
- [x] Grant/Revoke ops classified as PhaseMain (metadata, high priority)
- [x] CreateSchema op has Priority 0 (lowest = earliest)
- [x] CreateFunction op has Priority 4 (before table's 5 by default)
- [x] CreateTable op has Priority 5
- [x] CreateView op has Priority 8 (after tables)
- [x] Comment/Grant ops have Priority 12 (metadata last)
- [x] Existing tests in TestMigrationOrdering still pass after metadata addition

---

## Phase 2: Sort Engine

### 2.1 Forward Dependency Sorting (CREATE ordering)

When creating objects, dependencies from `to` catalog determine order: depended-on objects first.

- [x] Function referenced by table CHECK → function created before table
- [x] Function referenced by table DEFAULT → function created before table
- [x] Both CHECK function and DEFAULT function → both created before table
- [x] Enum type used as column type → enum created before table
- [x] Domain used as column type → domain created before table
- [x] Sequence in DEFAULT nextval → sequence created before table
- [x] Table INHERITS parent (both new) → parent created before child
- [x] Table PARTITION OF parent (both new) → parent created before child partition
- [x] Partition child dropped → child dropped before parent when both dropped
- [x] View depends on table → table created before view
- [x] View depends on another view (chain of 3) → correct chain order
- [x] Trigger depends on function + table → both created before trigger
- [x] Expression index references function → function created before index
- [x] Policy USING expression references function → function before policy
- [x] Function RETURNS SETOF table → table created before function (dep overrides priority)
- [x] Multiple tables sharing same CHECK function → function before all tables
- [x] No dependencies at all → pure priority ordering (schema < type < table < view)

### 2.2 Reverse Dependency Sorting (DROP ordering)

When dropping objects, dependents from `from` catalog must be dropped first.

- [x] Drop table + dependent view → view dropped before table
- [x] Drop table + dependent trigger → trigger dropped before table
- [x] Drop function + dependent trigger → trigger dropped before function
- [x] Drop table + its indexes → indexes dropped before table
- [x] Drop table with FK referencing another table → FK table can drop independently
- [x] Drop schema + all contained objects → contained objects dropped before schema
- [x] Drop two tables where one has FK to other → FK-referencing table first
- [x] Drop table + dependent policy → policy dropped before table

### 2.3 Dependency Lifting

Catalog deps are recorded at constraint/index/trigger granularity, but migration ops are at table/function level. Deps must be "lifted" to the owning op.

- [x] CHECK constraint → function dep lifted to owning table's CREATE op
- [x] DEFAULT expression → sequence dep lifted to owning table's CREATE op
- [x] Expression index → function dep lifted to index's CREATE op
- [x] Column type → type dep lifted to owning table's CREATE op
- [x] Trigger → function dep mapped correctly (trigger has its own op)
- [x] View query → table dep mapped correctly (view has its own op)
- [x] Constraint FK → target table dep excluded from forward sort (FK deferred)
- [x] Multiple ops sharing same OID (e.g., DROP TABLE + DROP INDEX on same relation) → all participate in ordering
- [x] Column AlterColumn ops ordered via parent table OID relative to dependent views
- [x] Dep referencing OID not in op set → gracefully ignored (no crash)
- [x] Op with zero ObjOID (unpopulated) → excluded from dep graph, ordered by priority only

---

## Phase 3: Heuristic Removal

### 3.1 Replace splitFunctionOps

The dependency graph replaces the string-matching heuristic for function ordering.

- [x] Function referenced by CHECK — ordered correctly by dep graph alone
- [x] Function overload: is_valid(integer) referenced by CHECK, is_valid(text) not — only integer overload forced before table (OID-level precision)
- [x] Function not referenced by any table — placed after tables by priority (late)
- [x] Function RETURNS SETOF table — placed after table by dependency (late)
- [x] No string-matching heuristic used for function ordering (all ordering derived from OID deps)
- [x] Existing test "functions created before views and triggers" still passes

### 3.2 Replace wrapColumnTypeChangesWithViewOps

Column type changes require dependent views to be dropped and recreated. The dependency graph handles ordering; only op injection remains.

- [x] Column type change (int→bigint) with dependent view → DROP VIEW, ALTER COLUMN, CREATE VIEW in correct order
- [x] Column type change with chain of dependent views (v2→v1→table) → both views dropped, column altered, both recreated
- [x] Column type change with no dependent views → ALTER COLUMN only, no extra ops
- [x] Dependent views identified via catalog deps (not string matching on view definition)

---

## Phase 4: Edge Cases and Mixed Operations

### 4.1 Mixed CREATE + DROP Scenarios

Real migrations often have both creates and drops in the same plan.

- [x] Replace table (drop old, create new) with dependent view → DROP VIEW, DROP old table, CREATE new table, CREATE VIEW
- [x] Function signature change with dependent trigger → DROP TRIGGER, DROP old function, CREATE new function, CREATE TRIGGER
- [x] Trigger function identity changed → trigger automatically gets DROP+CREATE injected
- [x] Add new table + drop unrelated table + add view on new table → drops before creates, view after new table
- [x] Enum value addition with dependent table (no table change needed) → ALTER TYPE only
- [x] Drop table + create replacement table with same name → drop first, then create
- [x] Column type change + new function used by CHECK on same table → function created, then view dropped, column altered, view recreated

### 4.2 Cycle Handling and FK Deferral

Circular dependencies must be detected and broken.

- [x] Self-referencing FK → CREATE TABLE, then ADD CONSTRAINT (deferred)
- [x] Circular FK between two tables → both tables created, then both FKs deferred
- [x] Three-way FK cycle (A→B→C→A) → all tables created, all FKs deferred
- [x] CHECK constraint creating cycle with function (function → table type, table → function CHECK) → CHECK deferred to PhasePost if cycle detected
- [x] No circular deps → no ops deferred unnecessarily
- [x] Cycle detection produces clear error for unresolvable cycles (e.g., view A → view B → view A)
- [x] FK deferred ops ordered by name for determinism
- [x] All existing TestMigrationOrdering and TestMigrationRoundtrip tests still pass
