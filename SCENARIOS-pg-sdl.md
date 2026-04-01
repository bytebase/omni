# PG LoadSDL Scenarios

> Goal: Implement `func LoadSDL(sql string) (*Catalog, error)` that loads declarative schema definitions with automatic dependency resolution and SDL validation.
> Verification: Call LoadSDL, verify catalog state matches expected. Compare with LoadSQL on correctly-ordered equivalent DDL.
> Reference sources: omni pg/catalog ProcessUtility, pg_dump dependency ordering, PostgreSQL DDL semantics

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: SDL Validation + Basic Infrastructure

### 1.1 SDL Validation

- [x] Empty string returns empty catalog with no error
- [x] Valid CREATE TABLE is accepted
- [x] Valid CREATE VIEW is accepted
- [x] Valid CREATE FUNCTION is accepted
- [x] Valid CREATE INDEX is accepted
- [x] Valid CREATE SEQUENCE is accepted
- [x] Valid CREATE SCHEMA is accepted
- [x] Valid CREATE TYPE (enum/domain/composite/range) is accepted
- [x] Valid CREATE EXTENSION is accepted
- [x] Valid CREATE TRIGGER is accepted
- [x] Valid CREATE POLICY is accepted
- [x] Valid CREATE MATERIALIZED VIEW is accepted
- [x] Valid CREATE CAST is accepted
- [x] Valid CREATE FOREIGN TABLE is accepted
- [x] Valid COMMENT ON is accepted
- [x] Valid GRANT is accepted
- [x] ALTER SEQUENCE OWNED BY is accepted
- [x] ALTER TABLE ENABLE ROW LEVEL SECURITY is accepted
- [x] ALTER TYPE ADD VALUE is accepted
- [x] INSERT statement is rejected with clear error
- [x] UPDATE statement is rejected with clear error
- [x] DELETE statement is rejected with clear error
- [x] DROP TABLE is rejected with clear error
- [x] ALTER TABLE ADD COLUMN is rejected with clear error
- [x] ALTER TABLE DROP COLUMN is rejected with clear error
- [x] TRUNCATE is rejected with clear error
- [x] DO block is rejected with clear error
- [x] Parse error in SDL returns error

### 1.2 Basic Dependency Resolution — Same Layer

- [x] Two tables with no dependencies load in any order
- [x] Table with FK to another table — FK table defined first in SDL
- [x] Table with FK to another table — FK table defined last in SDL (forward reference)
- [x] Table with column type referencing enum — enum defined after table
- [x] Table with column type referencing domain — domain defined after table
- [x] Table with column type referencing composite type — type defined after table
- [x] View referencing table — table defined after view
- [x] Index on table — table defined after index
- [x] Result catalog matches LoadSQL with correctly-ordered equivalent

### 1.3 Declared Object Collection and Name Resolution

- [x] Unqualified names default to "public" schema
- [x] Schema-qualified names preserved as-is
- [x] Same name in different schemas treated as distinct objects
- [x] Function identity includes argument types (overloads are distinct)
- [x] Built-in function in expression does NOT create dependency (e.g., now())
- [x] Built-in type in column does NOT create dependency (e.g., integer, text)
- [x] Reference to undeclared object is not a dependency (passes through to ProcessUtility for error)

---

## Phase 2: Expression Dependency Extraction

### 2.1 Column DEFAULT and CHECK Expression Dependencies

- [x] DEFAULT nextval('my_seq') creates dependency on sequence
- [x] DEFAULT my_function() creates dependency on function
- [x] CHECK (validate_func(col)) creates dependency on function
- [x] CHECK with subquery referencing table creates dependency on table
- [x] DEFAULT with type cast to user type creates dependency on type
- [x] Multiple defaults referencing different functions resolved correctly
- [x] Column default with built-in function (e.g., now()) does NOT create dependency

### 2.2 View Query Dependencies

- [x] Simple SELECT FROM table creates dependency on table
- [x] JOIN across two tables creates dependencies on both
- [x] Subquery in WHERE referencing another table creates dependency
- [x] Function call in SELECT list creates dependency on function
- [x] CTE (WITH clause) — CTE name NOT treated as external dependency
- [x] CTE body referencing real table creates dependency on that table
- [x] Nested subquery in FROM (subselect) dependencies extracted
- [x] View referencing another view — dependency resolved
- [x] UNION/INTERSECT/EXCEPT across tables — all dependencies extracted
- [x] View with type cast to user type creates dependency

### 2.3 Index, Trigger, and Policy Expression Dependencies

- [x] Expression index with function call creates dependency on function
- [x] Partial index WHERE clause with function creates dependency
- [x] Trigger WHEN clause with function creates dependency
- [x] Policy USING expression with function creates dependency
- [x] Policy WITH CHECK expression with function creates dependency
- [x] Domain CHECK constraint with function creates dependency
- [x] Trigger EXECUTE FUNCTION reference creates dependency on function
- [x] Policy on table creates structural dependency on table
- [x] ALTER SEQUENCE OWNED BY creates dependency on target table

### 2.4 Function and Type Dependencies

- [x] Function parameter type referencing user type creates dependency
- [x] Function RETURNS user_type creates dependency on type
- [x] Function RETURNS SETOF table creates dependency on table
- [x] Function parameter DEFAULT expression with function creates dependency
- [x] Composite type column referencing another user type creates dependency
- [x] Range type with user-defined subtype creates dependency
- [x] Domain based on another domain creates dependency

---

## Phase 3: Topological Sort, Cycles, and Complex Scenarios

### 3.1 Priority Layer Ordering

- [x] Schemas created before tables that reference them
- [x] Extensions created before types/functions they provide
- [x] Types created before tables using them as column types
- [x] Sequences created before tables with DEFAULT nextval referencing them
- [x] Tables created before views referencing them
- [x] Functions created before triggers referencing them
- [x] Tables created before indexes on them
- [x] FK constraints applied after all tables created
- [x] Comments applied after their target objects created
- [x] Grants applied after their target objects created

### 3.2 Topological Sort Within Layers

- [x] View A depends on View B — View B created first
- [x] Chain of 5 views (v5→v4→v3→v2→v1) — all created in correct order
- [x] Two independent views created without error (no false dependency)
- [x] Table with INHERITS parent — parent created first
- [x] Table PARTITION OF parent — parent created first
- [x] Composite type A references composite type B — B created first

### 3.3 Cycle Detection and Repair

- [x] Mutual FK between two tables — both created, FKs applied after (via FK layer)
- [x] Three-way FK cycle (A→B→C→A) — all tables created, FKs deferred
- [x] Self-referencing FK (table references itself) — handled correctly
- [x] Composite type mutual reference via arrays — shell types resolve the cycle
- [x] Unresolvable cycle produces clear error message

### 3.4 Complex Multi-Object SDL

- [x] Full schema with 10+ object types loaded from shuffled SDL
- [x] Schema with function used in CHECK, view, and trigger — single function resolved for all
- [x] Table with SERIAL column — implicit sequence created without explicit CREATE SEQUENCE
- [x] Materialized view with indexes — matview created before its indexes
- [x] Multiple schemas with cross-schema references
- [x] SDL producing identical catalog to LoadSQL with same DDL in correct order
