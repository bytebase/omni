# PG Migration Ordering Test Suite Scenarios

> Goal: Comprehensive, production-quality test suite validating that GenerateMigration produces SQL that executes correctly on a real PostgreSQL instance
> Verification: LoadSDL before → LoadSDL after → Diff → GenerateMigration → assertMigrationValid (phase ordering, no warnings); oracle roundtrip on real PG for Phase 5+
> Reference sources: PostgreSQL pg_depend documentation, pg_dump ordering model, real-world migration patterns

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Object Creation Patterns

Greenfield migrations — creating objects from scratch. Validates forward dependency ordering.

### 1.1 Basic Object Creation

- [ ] Create single table with PK — no dependencies
- [ ] Create table with serial/identity column — sequence created before table
- [ ] Create table with enum column — enum type created before table
- [ ] Create table with domain column — domain created before table
- [ ] Create table with composite type column — composite type created before table
- [ ] Create table with range type column — range type created before table
- [ ] Create table with array of enum column — enum created before table
- [ ] Create table in non-public schema — schema created before table
- [ ] Create sequence with OWNED BY — table created before sequence ownership set

### 1.2 Expression-Dependent Creation

- [ ] Table with CHECK constraint calling user function — function created before table
- [ ] Table with DEFAULT calling user function — function created before table
- [ ] Table with both CHECK function and DEFAULT function — both functions before table
- [ ] Table with generated column referencing user function — function created before table
- [ ] Domain with CHECK constraint calling user function — function created before domain, domain before table
- [ ] Index on expression calling user function — function created before index
- [ ] Partial index with WHERE clause calling user function — function created before index
- [ ] View with SELECT expression calling user function — function created before view
- [ ] RLS policy USING expression calling user function — function created before policy
- [ ] Trigger WHEN clause referencing column — table created before trigger

### 1.3 Multi-Level Dependency Chains

- [ ] Function → CHECK → table → view — 4-level forward chain
- [ ] Function → CHECK → table → view → view — 5-level forward chain
- [ ] Enum → table column → view → materialized view — 4-level type chain
- [ ] Sequence → DEFAULT → table → index → (expression index on same table) — parallel deps on same table
- [ ] Schema → enum → domain(base=enum) → table column(domain) — 4-level type hierarchy
- [ ] Function → trigger → table ← view (two paths converge on table) — diamond pattern

### 1.4 Foreign Key Patterns

- [ ] Simple FK between two new tables — referenced table created first, FK deferred
- [ ] Self-referencing FK (tree structure) — table created, FK deferred
- [ ] Mutual FK between two tables — both tables created, both FKs deferred
- [ ] Three-way FK cycle (A→B→C→A) — all tables created, all FKs deferred
- [ ] FK with ON DELETE CASCADE — same deferral as simple FK
- [ ] FK referencing table in different schema — schema and table before FK
- [ ] Multiple FKs from one table to several tables — all referenced tables before FK

### 1.5 Composite Creation Scenarios

- [ ] E-commerce schema: enum + domain + sequence + function + 3 tables(FK) + 3 indexes + 2 views + trigger + policy — all objects in correct order
- [ ] Multi-tenant schema: 2 schemas + shared function + per-schema tables + cross-schema view — schemas first, then shared objects
- [ ] Audit system: audit table + trigger function + triggers on 3 tables — function before all triggers, tables before triggers
- [ ] Materialized view with index — table before matview, matview before index on matview

---

## Phase 2: Object Drop Patterns

Teardown migrations — removing objects in reverse dependency order.

### 2.1 Simple Drop Ordering

- [ ] Drop table with dependent view — view dropped before table
- [ ] Drop table with dependent trigger — trigger dropped before table
- [ ] Drop table with dependent index — index dropped before table
- [ ] Drop table with dependent RLS policy — policy dropped before table
- [ ] Drop function with dependent trigger — trigger dropped before function
- [ ] Drop function with dependent view (via expression) — view dropped before function
- [ ] Drop enum type used by table column — table/column altered before type dropped
- [ ] Drop sequence used in DEFAULT — default removed before sequence dropped
- [ ] Drop schema with all contained objects — all objects dropped before schema

### 2.2 Cascading Drop Chains

- [ ] Drop table → dependent view → dependent view (chain) — deepest view first
- [ ] Drop function → trigger using it → table trigger is on (if table also dropped) — trigger before both
- [ ] Drop enum → table column using it → view on that table — view, then column, then enum
- [ ] Drop 3 tables with mutual FKs — FKs dropped, then tables in any order
- [ ] Drop table + its CHECK function (both removed) — table before function (table owns the dep)
- [ ] Drop composite type used as table column — table/column altered before type dropped
- [ ] Complete teardown: enum + domain + sequence + function + tables + views + triggers + indexes — full reverse ordering

### 2.3 Selective Drop (Keep Some Objects)

- [ ] Drop one of two tables, other has FK to it — FK constraint dropped, table dropped, other table kept
- [ ] Drop a function used by trigger, keep table — trigger dropped, function dropped, table survives
- [ ] Drop a view in a view chain, keep base view — only dependent views dropped
- [ ] Remove column from table with dependent view — view recreated with remaining columns

---

## Phase 3: Schema Refactoring

Modifications that combine DROP + CREATE + ALTER operations.

### 3.1 Column Changes

- [ ] Add column to table with dependent view — view updated
- [ ] Drop column from table with dependent view — view recreated without column
- [ ] Change column type (int→bigint) with dependent view — DROP VIEW, ALTER COLUMN, CREATE VIEW
- [ ] Change column type with view chain (3 levels) — all views dropped and recreated
- [ ] Change column type with dependent index — index dropped and recreated
- [ ] Change column type with dependent CHECK constraint — constraint updated
- [ ] Add NOT NULL to column with dependent view — view unaffected, column altered
- [ ] Add DEFAULT to column — column altered, no other changes needed

### 3.2 Function Changes

- [ ] Change function body (signature same) — CREATE OR REPLACE, no dep changes
- [ ] Change function signature — DROP old + CREATE new, dependents updated
- [ ] Change function signature with dependent trigger — trigger dropped/recreated around function change
- [ ] Change function signature with dependent CHECK — CHECK recreated (if separate constraint)
- [ ] Change function return type — dependent views need recreation
- [ ] Add function overload — new function created, existing overload untouched
- [ ] Drop function overload — only overload dropped, other signature survives

### 3.3 Table Restructuring

- [ ] Replace table entirely (drop old, create new with same name) — all old deps dropped, new created
- [ ] Replace table with dependent views — views dropped before old table, recreated after new
- [ ] Table split: contacts → people + emails (with FK) — new tables created, old dropped, views migrated
- [ ] Table merge: people + emails → contacts — old tables dropped, new created, views migrated
- [ ] Convert regular table to partitioned — recreate with PARTITION BY, create children
- [ ] Add inheritance (INHERITS) to existing table — parent must exist

### 3.4 View Changes

- [ ] Modify view definition (same columns) — CREATE OR REPLACE VIEW
- [ ] Modify view with column changes — DROP VIEW + CREATE VIEW
- [ ] Modify view that other views depend on — dependent views also recreated
- [ ] Replace view with materialized view — drop view, create matview
- [ ] Add index on materialized view — matview before index

---

## Phase 4: Type System Evolution

Changes to types that cascade through the schema.

### 4.1 Enum Changes

- [ ] Add enum value — ALTER TYPE before any new objects using the value
- [ ] Add enum value used in new partial index WHERE clause — ALTER TYPE before index
- [ ] Add enum value used in new CHECK constraint — ALTER TYPE before constraint
- [ ] Add enum value with dependent view (view definition unchanged) — ALTER TYPE only
- [ ] Create new enum replacing old one — old enum objects migrated, old enum dropped

### 4.2 Domain Changes

- [ ] Add constraint to domain — domain altered
- [ ] Drop constraint from domain — domain altered
- [ ] Change domain base type — requires DROP + CREATE, cascades to columns
- [ ] Domain used by function parameter — function must handle domain change

### 4.3 Composite and Range Type Changes

- [ ] Add field to composite type — ALTER TYPE, dependent functions may need update
- [ ] Create composite type used as function return type — type before function
- [ ] Create composite type used as table column — type before table
- [ ] Composite type A references type B — B created before A
- [ ] Create range type with subtype — subtype/type before range type
- [ ] Range type used in table column — range type before table

---

## Phase 5: Multi-Object Workflows

Real-world migrations involving 5+ coordinated object changes. These should be tested with oracle roundtrip against real PG.

### 5.1 Add Audit Trail to Existing Schema

- [ ] Before: 3 tables (users, orders, products). After: add audit_log table + trigger function + 3 triggers — function before triggers, triggers after their tables
- [ ] Audit trigger function references audit table — audit table before function (for INSERT INTO)

### 5.2 Add Search Infrastructure

- [ ] Before: documents table. After: add tsvector column + trigger function + trigger + GIN index — column added, function created, trigger created, index created in correct order
- [ ] Search function depends on table column — table column before function

### 5.3 Multi-Tenant RLS Setup

- [ ] Before: 3 tables without RLS. After: add tenant_id column + current_tenant() function + ENABLE RLS + policies on all 3 tables — function before policies, tables altered before policies applied

### 5.4 Complete Schema Replacement

- [ ] Before: v1 schema (5 tables, 3 views, 2 functions). After: v2 schema (different tables, views, functions) — all v1 dropped, all v2 created in correct order
- [ ] Schema replacement with some objects surviving — survivors not dropped, new deps correct

### 5.5 Add Partitioning to Existing Table

- [ ] Before: single orders table with indexes + views. After: partitioned orders with 4 range partitions + indexes — views dropped, old table restructured, partitions created, indexes recreated, views recreated

---

## Phase 6: Cycles, Edge Cases, and Boundary Conditions

### 6.1 Cycle Detection and Resolution

Tests that the cycle-breaking mechanism works correctly (FK deferral, warnings). Phase 1.4 tests FK creation ordering; this section tests the cycle detection engine specifically.

- [ ] FK cycle between 2 tables — both tables created, FKs deferred, no warning
- [ ] FK cycle between 3 tables — all tables created, all FKs deferred, no warning
- [ ] FK cycle + shared CHECK function — function before tables, FKs deferred
- [ ] No cycles present — no unnecessary deferrals, no warnings

### 6.2 Boundary Conditions

- [ ] Empty migration (before == after) — zero ops generated
- [ ] Migration with only metadata changes (comments, grants) — only metadata ops
- [ ] Single object migration (add one table) — single CREATE TABLE op
- [ ] Migration with 20+ tables in dependency chain — correct chain ordering
- [ ] Objects with same name in different schemas — no name collision in ordering
- [ ] Function with no dependencies, no dependents — ordered by priority only
- [ ] Table with 10+ indexes, 5+ constraints, 3+ triggers — all created in correct order after table

### 6.3 Diamond and Convergent Dependencies

- [ ] Two tables depending on same enum + same function — function and enum before both tables
- [ ] View joining two tables that both have triggers using same function — function before tables, tables before triggers, tables before view
- [ ] Two views depending on same table, one view depending on the other — base view before dependent view
- [ ] Index and CHECK constraint on same table both referencing same function — function before table (covers both)
