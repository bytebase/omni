# PG Oracle Testing Scenarios

> Goal: Validate the full SDL/Diff/Migration pipeline against real PostgreSQL. For each (before, after) schema pair, verify that omni-generated migration SQL executes successfully on PG and produces the correct schema state.
> Verification: Start PG testcontainer. Apply "before" DDL + migration SQL to schema X. Apply "after" DDL to schema Y. Compare X and Y via system table queries. They must match.
> Reference sources: Real PostgreSQL 16 (testcontainer), pg system catalogs (information_schema, pg_catalog)

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Test Infrastructure + Fully-Loaded Object Validation

### 1.1 PG Oracle Infrastructure

- [x] PG testcontainer starts and connects successfully
- [x] Schema-level isolation works (CREATE SCHEMA, DROP SCHEMA CASCADE between tests)
- [x] SQL execution helper executes multi-statement DDL on PG
- [x] Schema comparison queries return structured results for tables, columns, indexes, constraints, functions, views, triggers, sequences, types, comments, policies
- [x] Comparison function detects identical schemas as equal
- [x] Comparison function detects schema differences correctly
- [x] Migration roundtrip helper: (before, after) → LoadSDL → Diff → GenerateMigration → execute on PG → compare

### 1.2 Fully-Loaded Table Validates on PG

- [x] Create fully-loaded table on PG succeeds (all column types, all constraint types, indexes, comments)
- [x] Drop fully-loaded table on PG succeeds
- [x] Empty schema → fully-loaded table via omni migration succeeds on PG
- [x] Fully-loaded table → empty schema via omni migration succeeds on PG
- [x] Fully-loaded table roundtrip: migration result matches direct creation on PG

### 1.3 Fully-Loaded Function/View/Trigger Validates on PG

- [x] Create fully-loaded function on PG succeeds (all attributes: IMMUTABLE, STRICT, SECURITY DEFINER, LEAKPROOF, PARALLEL SAFE, multi-param with defaults)
- [x] Create fully-loaded procedure on PG succeeds
- [x] Create fully-loaded view on PG succeeds (JOIN, subquery, CTE, WITH CHECK OPTION, comment, trigger)
- [x] Create fully-loaded materialized view on PG succeeds (with indexes)
- [x] Create fully-loaded trigger on PG succeeds (BEFORE UPDATE OF columns, WHEN clause, REFERENCING)
- [x] Empty → all fully-loaded objects via omni migration succeeds and matches direct creation

### 1.4 Fully-Loaded Types/Sequences/Extensions Validates on PG

- [x] Create enum type with multiple values on PG succeeds
- [x] Create domain with NOT NULL, DEFAULT, CHECK constraint on PG succeeds
- [x] Create composite type on PG succeeds
- [x] Create range type on PG succeeds
- [x] Create sequence with all options on PG succeeds
- [x] Create extension (pgcrypto or similar) on PG succeeds
- [x] Empty → all types/sequences via omni migration succeeds and matches

---

## Phase 2: Single-Attribute Change Matrix

### 2.1 Table Column Property Changes

- [x] Change column type: varchar(100) → varchar(200) — migration correct on PG
- [x] Change column type: integer → bigint — migration correct on PG
- [x] Change column type: text → user enum — migration with USING correct on PG
- [x] Add NOT NULL to existing column — migration correct on PG
- [x] Drop NOT NULL from existing column — migration correct on PG
- [x] Add DEFAULT to column — migration correct on PG
- [x] Change DEFAULT value — migration correct on PG
- [x] Drop DEFAULT from column — migration correct on PG
- [x] Add new column with all attributes — migration correct on PG
- [x] Drop existing column — migration correct on PG
- [ ] Add GENERATED ALWAYS AS IDENTITY — migration correct on PG (**BUG: emits SET DEFAULT before ADD IDENTITY**)
- [ ] Drop identity from column — migration correct on PG (**BUG: DROP SEQUENCE before DROP IDENTITY**)
- [x] Change column collation — migration correct on PG
- [ ] Add GENERATED ALWAYS AS (expr) STORED column — migration correct on PG (**BUG: uses DEFAULT instead of GENERATED**)
- [x] Drop generated column — migration correct on PG
- [x] Table UNLOGGED → permanent (persistence change) — migration correct on PG

### 2.2 Table Constraint Changes

- [x] Add PRIMARY KEY — migration correct on PG
- [x] Drop PRIMARY KEY — migration correct on PG
- [x] Add UNIQUE constraint — migration correct on PG
- [x] Drop UNIQUE constraint — migration correct on PG
- [x] Add CHECK constraint (with function call) — migration correct on PG
- [x] Drop CHECK constraint — migration correct on PG
- [x] Change CHECK expression — migration correct on PG
- [x] Add FOREIGN KEY — migration correct on PG
- [x] Drop FOREIGN KEY — migration correct on PG
- [x] Change FK ON DELETE action (CASCADE → SET NULL) — migration correct on PG
- [ ] Add DEFERRABLE constraint — migration correct on PG (**BUG: missing DEFERRABLE clause**)
- [ ] Add EXCLUDE constraint — migration correct on PG (**BUG: missing USING gist**)
- [x] Drop EXCLUDE constraint — migration correct on PG
- [x] Change FK ON UPDATE action — migration correct on PG

### 2.3 Table Attached Object Changes

- [x] Add standalone index — migration correct on PG
- [x] Drop standalone index — migration correct on PG
- [x] Add partial index (WHERE clause) — migration correct on PG
- [x] Add expression index — migration correct on PG
- [x] Add index with INCLUDE columns — migration correct on PG
- [x] Add index with DESC/NULLS FIRST — migration correct on PG
- [x] Add GIN/GiST index — migration correct on PG
- [x] Change index (columns changed) — migration DROP+CREATE correct on PG
- [x] Add trigger with UPDATE OF columns — migration correct on PG
- [x] Drop trigger — migration correct on PG
- [x] Change trigger (WHEN clause changed) — migration DROP+CREATE correct on PG
- [x] Add RLS policy with USING and WITH CHECK — migration correct on PG
- [x] Drop policy — migration correct on PG
- [x] Enable ROW LEVEL SECURITY — migration correct on PG
- [x] Disable ROW LEVEL SECURITY — migration correct on PG
- [x] Add table comment — migration COMMENT ON TABLE correct on PG
- [ ] Add column comment — migration COMMENT ON COLUMN correct on PG (**BUG: missing column name qualification**)
- [x] Change table comment — migration correct on PG
- [x] Drop table comment (set NULL) — migration correct on PG
- [x] ALTER TABLE REPLICA IDENTITY FULL — migration correct on PG

### 2.4 Function/Procedure Changes

- [x] Change function body only — CREATE OR REPLACE FUNCTION correct on PG
- [x] Change function volatility — migration correct on PG
- [x] Change function STRICT/SECURITY/LEAKPROOF — migration correct on PG
- [x] Change function return type — DROP+CREATE correct on PG
- [x] Add function parameter — DROP+CREATE correct on PG
- [x] Change procedure body — CREATE OR REPLACE PROCEDURE correct on PG
- [x] Add function comment — COMMENT ON FUNCTION correct on PG
- [x] Add procedure comment — COMMENT ON PROCEDURE correct on PG

### 2.5 View/MatView Changes

- [x] Change view definition (query changed) — CREATE OR REPLACE VIEW correct on PG
- [x] Add view comment — COMMENT ON VIEW correct on PG
- [x] Change matview definition — DROP+CREATE MATERIALIZED VIEW correct on PG
- [x] Add matview comment — COMMENT ON MATERIALIZED VIEW correct on PG
- [x] Add index on matview — migration correct on PG
- [x] Change view WITH CHECK OPTION — migration correct on PG

### 2.6 Type/Sequence/Extension Changes

- [x] Enum: add value at end — ALTER TYPE ADD VALUE correct on PG
- [x] Enum: add value with BEFORE positioning — correct on PG
- [x] Enum: add value with AFTER positioning — correct on PG
- [x] Domain: change DEFAULT — ALTER DOMAIN correct on PG
- [x] Domain: change NOT NULL — ALTER DOMAIN correct on PG
- [x] Domain: add constraint — ALTER DOMAIN ADD CONSTRAINT correct on PG
- [x] Sequence: change INCREMENT — ALTER SEQUENCE correct on PG (not DROP+CREATE)
- [x] Sequence: change CYCLE — ALTER SEQUENCE correct on PG
- [x] Add new enum type — migration correct on PG
- [x] Drop enum type — migration correct on PG
- [x] Composite type: add/drop/modify column — migration correct on PG
- [x] Range type: change subtype — migration generates warning on PG
- [x] Schema: add new schema — migration correct on PG
- [x] Schema: drop schema — migration correct on PG
- [x] GRANT SELECT on table — migration correct on PG
- [x] REVOKE privilege — migration correct on PG
- [x] Partitioned table (PARTITION BY RANGE) — migration correct on PG
- [x] Partition child (PARTITION OF ... FOR VALUES) — migration correct on PG

---

## Phase 3: Combinations, Interactions, and Known Regressions

### 3.1 Multi-Attribute Simultaneous Changes

- [x] Table: add column + change another column type + drop constraint — all correct in single migration
- [x] Function: change body + change volatility + change parallel — single CREATE OR REPLACE
- [x] Table: add index + add trigger + add policy — all attached objects in one migration
- [x] Multiple tables changed simultaneously — migration correct on PG

### 3.2 Cross-Object Dependency Changes

- [x] Modify function used by CHECK constraint — migration order correct on PG
- [x] Modify function used by trigger — migration order correct on PG
- [x] Modify function used by view — migration order correct on PG
- [x] Add enum type + add column using it — both in single migration, correct order
- [x] Drop table referenced by FK — CASCADE or explicit drop order correct on PG
- [~] Modify table column type used by view — view recreation if needed (**BUG: migration doesn't recreate dependent view**)

### 3.3 Known Regression Tests (from bytebase migration)

- [x] View comment uses COMMENT ON VIEW (not COMMENT ON TABLE) — verified on PG
- [x] Materialized view comment uses COMMENT ON MATERIALIZED VIEW — verified on PG
- [x] Procedure comment uses COMMENT ON PROCEDURE — verified on PG
- [x] Procedure body change uses CREATE OR REPLACE PROCEDURE — verified on PG
- [x] Trigger with UPDATE OF columns preserves column syntax — verified on PG
- [x] Sequence modification uses ALTER SEQUENCE (not DROP+CREATE) — verified on PG
- [x] FK between two tables (forward reference in SDL) — correct on PG
- [~] Mutual FK (A↔B) — both tables created, FKs deferred, correct on PG (**BUG: LoadSDL circular inline FK**)

---

## Phase 4: Real-World Schema Migrations

### 4.1 Complete Schema Lifecycle

- [x] Empty → complex schema (10+ tables, views, functions, triggers, enums) via migration on PG
- [x] Complex schema → modified schema (multiple changes across object types) via migration on PG
- [x] Complex schema → empty (drop everything) via migration on PG
- [x] Migration result matches direct creation: Diff(migrated, direct) is empty on both PG and omni

### 4.2 SDL Forward Reference Stress Test

- [x] SDL with all objects in reverse dependency order — LoadSDL resolves, migration correct on PG
- [~] SDL with circular FK — LoadSDL resolves, migration with deferred FK correct on PG (**BUG: same circular FK**)
- [~] SDL with function used by CHECK, view, and trigger — single function resolved everywhere on PG (**BUG: function ordering before CHECK-dependent table**)
