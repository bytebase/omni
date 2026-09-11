# PG Catalog Migration Generator Scenarios

> Goal: Implement `func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan` that produces ordered, structured DDL operations from a SchemaDiff.
> Verification: LoadSQL two schemas, Diff, GenerateMigration, assert ops contain expected DDL. Roundtrip: apply generated SQL to `from` catalog, verify result matches `to`.
> Reference sources: omni pg/catalog diff types, pgschema internal/diff patterns, PostgreSQL DDL semantics

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Infrastructure + Core DDL Generation

### 1.1 MigrationPlan Types and Entry Point

- [x] MigrationPlan, MigrationOp, MigrationOpType types compile
- [x] GenerateMigration on empty SchemaDiff returns empty plan
- [x] GenerateMigration on identical catalogs returns empty plan
- [x] MigrationPlan.SQL() joins all statements with semicolons and newlines
- [x] MigrationPlan.Summary() groups ops by type and counts add/drop/modify
- [x] MigrationPlan.Filter() returns subset matching predicate
- [x] MigrationPlan.HasWarnings() detects ops with non-empty Warning
- [x] MigrationPlan.Warnings() returns only warning ops

### 1.2 Schema DDL

- [x] CREATE SCHEMA for added schema
- [x] DROP SCHEMA CASCADE for dropped schema
- [x] ALTER SCHEMA OWNER TO for modified schema owner
- [x] Schema operations ordered before table operations

### 1.3 Table CREATE and DROP

- [x] CREATE TABLE with columns and inline PK/UNIQUE/CHECK
- [x] CREATE TABLE with column defaults
- [x] CREATE TABLE with NOT NULL columns
- [x] CREATE TABLE with column identity (GENERATED ALWAYS AS IDENTITY)
- [x] CREATE TABLE with generated column (GENERATED ALWAYS AS ... STORED)
- [x] CREATE UNLOGGED TABLE for unlogged tables
- [x] DROP TABLE CASCADE for dropped tables
- [x] All identifiers double-quoted in generated DDL
- [x] Schema-qualified table names in DDL

### 1.4 Column ALTER

- [x] ALTER TABLE ADD COLUMN for new column
- [x] ALTER TABLE ADD COLUMN with NOT NULL and DEFAULT
- [x] ALTER TABLE DROP COLUMN for removed column
- [x] ALTER TABLE ALTER COLUMN TYPE for type change (implicit cast)
- [x] ALTER TABLE ALTER COLUMN TYPE ... USING for type change (no implicit cast)
- [x] USING clause decision uses FindCoercionPathway, not heuristic
- [x] Type change with existing default: DROP DEFAULT → ALTER TYPE → SET DEFAULT
- [x] ALTER TABLE ALTER COLUMN SET NOT NULL
- [x] ALTER TABLE ALTER COLUMN DROP NOT NULL
- [x] ALTER TABLE ALTER COLUMN SET DEFAULT
- [x] ALTER TABLE ALTER COLUMN DROP DEFAULT
- [x] ALTER TABLE ALTER COLUMN ADD GENERATED ... AS IDENTITY
- [x] ALTER TABLE ALTER COLUMN DROP IDENTITY
- [x] Multiple column changes on same table batched into fewer statements

### 1.5 Constraint DDL

- [x] ALTER TABLE ADD CONSTRAINT ... PRIMARY KEY
- [x] ALTER TABLE ADD CONSTRAINT ... UNIQUE
- [x] ALTER TABLE ADD CONSTRAINT ... CHECK (expr)
- [x] ALTER TABLE ADD CONSTRAINT ... FOREIGN KEY REFERENCES
- [x] FK with ON DELETE/UPDATE actions in generated DDL
- [x] FK with DEFERRABLE INITIALLY DEFERRED in generated DDL
- [x] Constraint with NOT VALID generates ADD CONSTRAINT ... NOT VALID
- [x] ALTER TABLE DROP CONSTRAINT for removed constraints
- [x] Modified constraint generated as DROP + ADD (no ALTER CONSTRAINT for structure)
- [x] EXCLUDE constraint ADD/DROP generated correctly

### 1.6 Index DDL

- [x] CREATE INDEX for added standalone index
- [x] CREATE UNIQUE INDEX for unique indexes
- [x] CREATE INDEX ... USING hash/gist/gin for non-btree
- [x] CREATE INDEX with WHERE clause (partial index)
- [x] CREATE INDEX with INCLUDE columns
- [x] CREATE INDEX with expression columns
- [x] CREATE INDEX with DESC/NULLS FIRST options
- [x] DROP INDEX for removed index
- [x] Modified index generated as DROP + CREATE (indexes can't be ALTERed)
- [x] PK/UNIQUE backing indexes NOT generated (managed by constraint DDL)

### 1.7 Partition and Inheritance DDL

- [x] CREATE TABLE ... PARTITION BY RANGE/LIST/HASH for partitioned table
- [x] CREATE TABLE ... PARTITION OF ... FOR VALUES FROM ... TO for range partition
- [x] CREATE TABLE ... PARTITION OF ... FOR VALUES IN (...) for list partition
- [x] CREATE TABLE ... PARTITION OF ... DEFAULT for default partition
- [x] ALTER TABLE REPLICA IDENTITY (DEFAULT/FULL/NOTHING) for replica identity change
- [x] CREATE VIEW ... WITH CHECK OPTION for view check option
- [x] Table inheritance: INHERITS clause in CREATE TABLE

---

## Phase 2: Functions, Sequences, Types, Views, Triggers

### 2.1 Sequence DDL

- [x] CREATE SEQUENCE with all options (INCREMENT, MINVALUE, MAXVALUE, START, CACHE, CYCLE)
- [x] DROP SEQUENCE for removed sequence
- [x] ALTER SEQUENCE for modified properties (increment, min, max, cache, cycle)
- [x] SERIAL-owned sequences not generated (managed by column)

### 2.2 Function and Procedure DDL

- [x] CREATE FUNCTION with full signature (args, return type, language, body)
- [x] CREATE FUNCTION with volatility, strictness, security, parallel, leakproof attributes
- [x] Dollar-quoting for function body ($$...$$)
- [x] DROP FUNCTION with argument types in signature
- [x] CREATE OR REPLACE FUNCTION for modified function (body/attribute changes)
- [x] CREATE PROCEDURE for added procedure
- [x] DROP PROCEDURE with argument types
- [x] Modified procedure generated as DROP + CREATE (no CREATE OR REPLACE PROCEDURE in older PG)
- [x] Overloaded functions produce distinct DDL (different arg signatures)

### 2.3 Enum Type DDL

- [x] CREATE TYPE ... AS ENUM ('val1', 'val2') for added enum
- [x] DROP TYPE for removed enum
- [x] ALTER TYPE ADD VALUE 'new_val' for appended value
- [x] ALTER TYPE ADD VALUE 'new_val' AFTER 'existing' for positioned value
- [x] ALTER TYPE ADD VALUE 'new_val' BEFORE 'existing' for positioned value
- [x] Enum value removal generates Warning (PG limitation)
- [x] ALTER TYPE ADD VALUE cannot run in transaction (Transactional = false)

### 2.4 Domain Type DDL

- [x] CREATE DOMAIN for added domain (base type, default, NOT NULL, constraints)
- [x] DROP DOMAIN for removed domain
- [x] ALTER DOMAIN SET DEFAULT / DROP DEFAULT
- [x] ALTER DOMAIN SET NOT NULL / DROP NOT NULL
- [x] ALTER DOMAIN ADD CONSTRAINT / DROP CONSTRAINT

### 2.5 Trigger DDL

- [x] CREATE TRIGGER with timing, events, level, function
- [x] CREATE TRIGGER with WHEN clause
- [x] CREATE TRIGGER with REFERENCING OLD/NEW TABLE AS
- [x] DROP TRIGGER ON table for removed trigger
- [x] Modified trigger generated as DROP + CREATE

### 2.6 View and Materialized View DDL

- [x] CREATE VIEW AS ... for added view
- [x] DROP VIEW for removed view
- [x] CREATE OR REPLACE VIEW for modified view (definition change)
- [x] CREATE MATERIALIZED VIEW AS ... for added matview
- [x] DROP MATERIALIZED VIEW for removed matview
- [x] Modified matview generated as DROP + CREATE (no CREATE OR REPLACE)
- [x] View dependency chain: DROP dependents before target, recreate after
- [x] Matview indexes generated after matview creation

### 2.7 Range Type DDL

- [x] CREATE TYPE ... AS RANGE for added range type
- [x] DROP TYPE for removed range type
- [x] Range subtype change generates Warning (requires DROP + CREATE, may break dependents)

### 2.8 Extension DDL

- [x] CREATE EXTENSION for added extension
- [x] DROP EXTENSION for removed extension
- [ ] Extension schema change generates appropriate DDL
- [ ] Extension operations ordered before types and tables

---

## Phase 3: Metadata, Ordering, and Roundtrip

### 3.1 Policy and RLS DDL

- [x] CREATE POLICY on table for added policy
- [x] DROP POLICY on table for removed policy
- [x] ALTER POLICY for simple changes (roles, USING, WITH CHECK)
- [x] Complex policy change as DROP + CREATE
- [x] ALTER TABLE ENABLE ROW LEVEL SECURITY
- [x] ALTER TABLE DISABLE ROW LEVEL SECURITY
- [x] ALTER TABLE FORCE ROW LEVEL SECURITY
- [x] ALTER TABLE NO FORCE ROW LEVEL SECURITY

### 3.2 Comment DDL

- [x] COMMENT ON TABLE for added/changed comment
- [x] COMMENT ON TABLE IS NULL for removed comment
- [x] COMMENT ON COLUMN for column comments
- [x] COMMENT ON INDEX, FUNCTION, SCHEMA, TYPE, SEQUENCE, CONSTRAINT, TRIGGER
- [x] Comments generated after object creation

### 3.3 Grant DDL

- [x] GRANT privilege ON object TO role for added grant
- [x] REVOKE privilege ON object FROM role for removed grant
- [x] GRANT WITH GRANT OPTION
- [x] Modified grant generated as REVOKE + GRANT
- [x] Grants generated after object creation

### 3.4 Operation Ordering

- [x] DROP phase before CREATE phase before ALTER phase
- [x] Within DROP: triggers → views → functions → constraints → indexes → columns → tables → sequences → types → extensions → schemas
- [x] Within CREATE: schemas → extensions → types → sequences → tables (no FK) → functions → views → triggers → deferred FKs → indexes → policies
- [x] FK constraints deferred until all tables created
- [x] FK cycle detected — all FKs deferred to ALTER phase
- [x] Types created before tables that use them
- [x] Functions created before views/triggers that reference them

### 3.5 Roundtrip Verification

- [x] Simple roundtrip: apply migration SQL to `from` catalog → Diff with `to` is empty
- [x] Roundtrip with table add + column modify + index drop
- [x] Roundtrip with function + view + trigger changes
- [x] Roundtrip with enum add value + domain alter
- [x] Roundtrip with FK across schemas
- [x] Roundtrip with all object types simultaneously
