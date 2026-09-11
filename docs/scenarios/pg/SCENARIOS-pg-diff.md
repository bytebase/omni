# PG Catalog Diff Engine Scenarios

> Goal: Implement `func Diff(from, to *Catalog) *SchemaDiff` that compares two catalog states and outputs structured add/drop/modify lists for all PG object types.
> Verification: LoadSQL(schemaA), LoadSQL(schemaB), call Diff, assert expected entries
> Reference sources: omni pg/catalog struct definitions, pgschema internal/diff design, PostgreSQL DDL semantics

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Infrastructure + Core Object Types

### 1.1 Diff Types and Schema Diff

- [x] SchemaDiff struct and all diff entry types compile with no errors
- [x] Diff of two identical empty catalogs returns empty SchemaDiff
- [x] Diff of two identical non-empty catalogs returns empty SchemaDiff
- [x] Schema added — new schema in `to` not in `from`
- [x] Schema dropped — schema in `from` not in `to`
- [x] Schema modified — schema owner changed
- [x] Multiple schemas added/dropped/modified in one diff
- [x] pg_catalog and pg_toast schemas are excluded from diff

### 1.2 Table Diff (Basic)

- [x] Table added — CREATE TABLE in `to` only
- [x] Table dropped — table in `from` not in `to`
- [x] Table unchanged — identical table in both catalogs produces no entry
- [x] Multiple tables across different schemas
- [x] Table persistence changed (permanent ↔ unlogged)
- [x] Table in non-public schema added/dropped
- [x] Table with same name in different schemas treated as distinct objects
- [x] Table renamed detected as drop old + add new (identity is name-based)
- [x] Table ReplicaIdentity changed

### 1.3 Column Diff

- [x] Column added to existing table
- [x] Column dropped from existing table
- [x] Column type changed (e.g., int → bigint) detected via FormatType
- [x] Column nullability changed (NOT NULL added/removed)
- [x] Column default added
- [x] Column default removed
- [x] Column default value changed
- [x] Column identity added (GENERATED ALWAYS AS IDENTITY)
- [x] Column identity removed
- [x] Generated column added (GENERATED ALWAYS AS ... STORED)
- [x] Generated column expression changed
- [x] Generated column removed
- [x] Column collation changed
- [x] Multiple column changes on same table
- [x] Column unchanged — identical column produces no entry

### 1.4 Constraint Diff

- [x] PRIMARY KEY added
- [x] PRIMARY KEY dropped
- [x] UNIQUE constraint added
- [x] UNIQUE constraint dropped
- [x] CHECK constraint added
- [x] CHECK constraint dropped
- [x] CHECK constraint expression changed
- [x] Foreign key added
- [x] Foreign key dropped
- [x] Foreign key actions changed (ON DELETE CASCADE → SET NULL)
- [x] Foreign key match type changed
- [x] Constraint deferrable/deferred changed
- [x] EXCLUDE constraint added/dropped
- [x] Constraint identity by name — same name = same constraint

### 1.5 Index Diff

- [x] Standalone index added (CREATE INDEX)
- [x] Standalone index dropped
- [x] Index columns changed
- [x] Index access method changed (btree → hash)
- [x] Index uniqueness changed
- [x] Partial index WHERE clause changed
- [x] Index INCLUDE columns changed (NKeyColumns)
- [x] Index sort options changed (DESC, NULLS FIRST)
- [x] NULLS NOT DISTINCT changed
- [x] Expression index — expressions changed
- [x] PK/UNIQUE backing index skipped (ConstraintOID != 0)

---

## Phase 2: Functions, Sequences, Types, and Relation Sub-Objects

### 2.1 Sequence Diff

- [x] Standalone sequence added
- [x] Standalone sequence dropped
- [x] Sequence increment changed
- [x] Sequence min/max values changed
- [x] Sequence cycle flag changed
- [x] Sequence start value changed
- [x] Sequence cache value changed
- [x] Sequence type changed (int4 → int8)
- [x] SERIAL/IDENTITY-owned sequence skipped (OwnerRelOID != 0)
- [x] Sequence unchanged produces no entry

### 2.2 Function and Procedure Diff

- [x] Function added
- [x] Function dropped
- [x] Function body changed
- [x] Function volatility changed (VOLATILE → IMMUTABLE)
- [x] Function strictness changed
- [x] Function security changed (DEFINER ↔ INVOKER)
- [x] Function language changed
- [x] Function return type changed
- [x] Function parallel safety changed
- [x] Function leak-proof changed
- [x] Function RETURNS SETOF changed
- [x] Procedure added/dropped (Kind='p')
- [x] Overloaded functions — same name different args are distinct
- [x] Function identity uses schema + name + input arg types (FormatType)
- [x] Function unchanged produces no entry

### 2.3 Enum Type Diff

- [x] Enum type added
- [x] Enum type dropped
- [x] Enum value added (new label appended)
- [x] Enum value added with BEFORE/AFTER positioning
- [x] Enum values reordered detected
- [x] Enum type identity by schema + name
- [x] Enum unchanged produces no entry

### 2.4 Domain Type Diff

- [x] Domain added
- [x] Domain dropped
- [x] Domain base type changed
- [x] Domain NOT NULL changed
- [x] Domain default changed
- [x] Domain constraint added
- [x] Domain constraint dropped
- [x] Domain constraint expression changed
- [x] Domain unchanged produces no entry

### 2.5 Trigger Diff

- [x] Trigger added on table
- [x] Trigger dropped from table
- [x] Trigger timing changed (BEFORE → AFTER)
- [x] Trigger events changed (INSERT → INSERT OR UPDATE)
- [x] Trigger level changed (ROW ↔ STATEMENT)
- [x] Trigger WHEN clause changed
- [x] Trigger function changed
- [x] Trigger enabled state changed
- [x] Trigger UPDATE OF columns changed
- [x] Trigger transition tables changed
- [x] Trigger arguments changed
- [x] Trigger unchanged produces no entry

### 2.6 View and Materialized View Diff

- [x] View added
- [x] View dropped
- [x] View definition changed (detected via GetViewDefinition deparse)
- [x] View CHECK OPTION changed (LOCAL ↔ CASCADED ↔ none)
- [x] Materialized view added
- [x] Materialized view dropped
- [x] Materialized view definition changed
- [x] View/matview treated as relation with appropriate RelKind

### 2.7 Range Type Diff

- [x] Range type added
- [x] Range type dropped
- [x] Range subtype changed
- [x] Range type identity by schema + name
- [x] Range unchanged produces no entry

---

## Phase 3: Metadata Objects, Grants, and Edge Cases

### 3.1 Extension Diff

- [x] Extension added
- [x] Extension dropped
- [x] Extension schema changed
- [x] Extension relocatable flag changed
- [x] Extension unchanged produces no entry

### 3.2 Policy Diff (RLS)

- [x] Policy added on table
- [x] Policy dropped from table
- [x] Policy command type changed
- [x] Policy permissive/restrictive changed
- [x] Policy roles changed
- [x] Policy USING expression changed
- [x] Policy WITH CHECK expression changed
- [x] Table RLS enabled/disabled changed (RowSecurity flag)
- [x] Table FORCE RLS changed (ForceRowSecurity flag)

### 3.3 Comment Diff

- [x] Comment added on table
- [x] Comment dropped from table
- [x] Comment changed on table
- [x] Comment on column added/changed/dropped
- [x] Comment on index
- [x] Comment on function
- [x] Comment on schema
- [x] Comment on type (enum/domain)
- [x] Comment on sequence
- [x] Comment on constraint
- [x] Comment on trigger

### 3.4 Grant Diff

- [x] Grant on table added
- [x] Grant on table revoked
- [x] Grant privilege changed (SELECT → ALL)
- [x] Grant WITH GRANT OPTION changed
- [x] Column-level grant added/revoked
- [x] Grant on function added/revoked
- [x] Grant on schema added/revoked
- [x] Grant on sequence added/revoked
- [x] Grant identity: (objType, objName, grantee, privilege) tuple

### 3.5 Partition and Inheritance Diff

- [x] Partitioned table added (with PARTITION BY)
- [x] Partition child added (PARTITION OF)
- [x] Partition strategy detected in diff (list/range/hash)
- [x] Partition bound values captured in diff
- [x] Table inheritance changed (INHERITS list modified)

### 3.6 Determinism and Edge Cases

- [x] Diff output is deterministic — same input always produces same order
- [x] Diff with many objects of all types simultaneously
- [x] Empty schema (exists but has no objects) vs non-existent schema
- [x] Object with reserved word name (e.g., table named "order")
- [x] FK referencing table in different schema resolves correctly
- [x] Diff is symmetric in structure — Diff(A,B) add = Diff(B,A) drop
