# PG Loader Compatibility

> Goal: PostgreSQL 17 accepted schema DDL, object-control DDL, view definitions, and function signatures should parse and load through omni without compatibility-only rejections.
> Verification: default unit tests cover permanent regressions; `go test -tags=oracle ./pg/catalog` compares selected scenarios against a real PostgreSQL 17 oracle.

Status: [ ] pending, [x] covered, [~] partial

Implementation status (2026-05-11): all six phases are covered by the loader compatibility corpus. The executable corpus currently includes 159 PostgreSQL-accepted loader cases, 14 PostgreSQL-rejected compatibility cases, parser tail-safety tests, targeted dependency assertions, and a PostgreSQL 17 oracle test for the accept/reject corpus. Verified with `go test ./pg/parser -count=1`, `go test ./pg/catalog -count=1`, and `go test -tags=oracle ./pg/catalog -run 'TestLoaderCompatPG17Oracle(Accepts|Rejects)Corpus' -count=1`.

Notable PG behavior clarified during implementation: `numeric` child FKs to `integer` parents are rejected by PG17; user-defined variadic functions with only a variadic parameter reject zero expanded arguments; `CREATE VIEW v(a) AS SELECT x, y` is accepted and names only the leading output column; duplicate view output column names are rejected; CHECK and generated-column expressions reject subqueries; parent partition indexes created before partitions may cause PG to auto-create attached child indexes.

## Phase 1: Real Dump Regressions

### 1.1 Known Raw Dump Hits

- [x] P1.1.01 `CREATE TABLE t ();` loads as a zero-column ordinary table.
- [x] P1.1.02 `FOREIGN KEY (...) REFERENCES ... MATCH SIMPLE` is consumed as a valid match option.
- [x] P1.1.03 `COMMENT ON FUNCTION f(arg_name type, ...)` resolves by argument types and ignores names for identity.
- [x] P1.1.04 `CREATE FUNCTION ... RETURNS TABLE(...) LANGUAGE plpgsql` loads with table return metadata.
- [x] P1.1.05 `concat_ws(text, variadic any)` resolves when called with more than `pg_proc.pronargs`.
- [x] P1.1.06 `jsonb ->> c.col_name` resolves when `c.col_name` is inferred from `unnest(text[]) AS c(col_name)`.
- [x] P1.1.07 A view can read from a local `WITH cte AS (...)` range entry.
- [x] P1.1.08 `ALTER INDEX parent ATTACH PARTITION child` resolves parent and child in the index namespace.
- [x] P1.1.09 A `bigint` foreign key can reference an `integer` primary key through cross-type equality.
- [x] P1.1.10 A later `LATERAL` item can reference an alias introduced by a previous `LATERAL` item.

### 1.2 Known-Hit Variants

- [x] P1.2.01 A zero-column table can be created in a non-public schema.
- [x] P1.2.02 A zero-column table can be commented on.
- [x] P1.2.03 A zero-column table can be granted privileges.
- [x] P1.2.04 A zero-column table can be referenced by `CREATE VIEW v AS SELECT * FROM t`.
- [x] P1.2.05 `MATCH SIMPLE` works on a table-level multi-column foreign key.
- [x] P1.2.06 `MATCH SIMPLE` works with `ON UPDATE CASCADE`.
- [x] P1.2.07 `MATCH SIMPLE` works with `ON DELETE SET NULL`.
- [x] P1.2.08 `MATCH SIMPLE` works with `DEFERRABLE INITIALLY DEFERRED`.
- [x] P1.2.09 `COMMENT ON FUNCTION` accepts schema-qualified function names with argument names.
- [x] P1.2.10 `COMMENT ON FUNCTION` accepts quoted argument names.
- [x] P1.2.11 `COMMENT ON FUNCTION` accepts qualified argument types such as `pg_catalog.int4`.
- [x] P1.2.12 `COMMENT ON FUNCTION` accepts variadic argument mode in the identity list.
- [x] P1.2.13 `RETURNS TABLE` with one column loads as a set-returning function.
- [x] P1.2.14 `RETURNS TABLE` with multiple columns loads as record-returning metadata.
- [x] P1.2.15 `RETURNS TABLE` with qualified column types loads.
- [x] P1.2.16 `concat_ws` in a view resolves with unknown string literals.
- [x] P1.2.17 `concat_ws` in a view resolves with integer, uuid, and timestamp arguments.
- [x] P1.2.18 `jsonb ->> alias.col` works when the alias comes from `unnest(varchar[])`.
- [x] P1.2.19 `jsonb ->> alias.col` works when the alias comes from a CTE output column.
- [x] P1.2.20 Simple CTE range resolution works with an explicit CTE column alias list.
- [x] P1.2.21 Simple CTE range resolution works when the CTE name shadows a base table.
- [x] P1.2.22 `ALTER INDEX ATTACH PARTITION` works when parent and child indexes are schema-qualified.
- [x] P1.2.23 `ALTER INDEX ATTACH PARTITION` records child-index to parent-index dependency.
- [x] P1.2.24 `bigint` to `integer` foreign key works in a composite key.
- [x] P1.2.25 `integer` to `bigint` foreign key works in the reverse direction.
- [x] P1.2.26 Chained `LATERAL` subqueries can reference both base table and prior lateral aliases.
- [x] P1.2.27 `LEFT JOIN LATERAL` can reference the left relation in the lateral subquery.
- [x] P1.2.28 `CROSS JOIN LATERAL` can reference the left relation in the lateral subquery.
- [x] P1.2.29 LATERAL correlation inside a view does not panic provenance tracking.
- [x] P1.2.30 The real dump regression corpus runs through PG17 oracle and omni loader with matching accept status.

## Phase 2: Parser Grammar Drift

### 2.1 Function Identity Lists

- [x] P2.1.01 `GRANT EXECUTE ON FUNCTION f(arg_name type, ...)` accepts optional argument names.
- [x] P2.1.02 `REVOKE EXECUTE ON FUNCTION f(arg_name type, ...)` accepts optional argument names.
- [x] P2.1.03 `ALTER FUNCTION f(arg_name type, ...)` accepts optional argument names.
- [x] P2.1.04 `DROP FUNCTION f(arg_name type, ...)` accepts optional argument names.
- [x] P2.1.05 `COMMENT ON FUNCTION f(IN a integer)` resolves by type.
- [x] P2.1.06 `COMMENT ON FUNCTION f(INOUT a integer)` resolves by type.
- [x] P2.1.07 `COMMENT ON FUNCTION f(VARIADIC a integer[])` resolves by type.
- [x] P2.1.08 `GRANT EXECUTE ON FUNCTION f(IN a integer)` resolves by type.
- [x] P2.1.09 `GRANT EXECUTE ON FUNCTION f(INOUT a integer)` resolves by type.
- [x] P2.1.10 `GRANT EXECUTE ON FUNCTION f(VARIADIC a integer[])` resolves by type.
- [x] P2.1.11 `REVOKE EXECUTE ON FUNCTION f(IN a integer)` resolves by type.
- [x] P2.1.12 `ALTER FUNCTION f(IN a integer) OWNER TO role_name` resolves by type.
- [x] P2.1.13 `ALTER FUNCTION f(IN a integer) SET SCHEMA target_schema` resolves by type.
- [x] P2.1.14 `ALTER FUNCTION f(IN a integer) IMMUTABLE` resolves by type.
- [x] P2.1.15 `DROP FUNCTION f(IN a integer)` resolves by type.
- [x] P2.1.16 `DROP PROCEDURE p(INOUT a integer)` resolves by type.
- [x] P2.1.17 `DROP ROUTINE r(a integer)` resolves by type.
- [x] P2.1.18 `COMMENT ON PROCEDURE p(arg_name integer)` resolves by type.
- [x] P2.1.19 `COMMENT ON ROUTINE r(arg_name integer)` resolves by type.
- [x] P2.1.20 Function identity lists accept schema-qualified types.
- [x] P2.1.21 Function identity lists accept array type syntax.
- [x] P2.1.22 Function identity lists accept quoted argument names.
- [x] P2.1.23 Function identity lists accept quoted function names.
- [x] P2.1.24 Function identity lists accept multi-part schema-qualified function names.
- [x] P2.1.25 Function identity lists ignore OUT-only parameters where PostgreSQL ignores them.

### 2.2 Keyword Token Parity

- [x] P2.2.01 `MATCH SIMPLE` is accepted when `simple` scans as the `SIMPLE` keyword token.
- [x] P2.2.02 `MATCH FULL` is accepted when `full` scans as a keyword token.
- [x] P2.2.03 `MATCH PARTIAL` follows PostgreSQL accept/reject behavior.
- [x] P2.2.04 Non-reserved keywords can be used as table names where PostgreSQL permits them.
- [x] P2.2.05 Non-reserved keywords can be used as column names where PostgreSQL permits them.
- [x] P2.2.06 Non-reserved keywords can be used as function names where PostgreSQL permits them.
- [x] P2.2.07 Non-reserved keywords can be used as type names where PostgreSQL permits them.
- [x] P2.2.08 Non-reserved keywords can be used as schema names where PostgreSQL permits them.
- [x] P2.2.09 Reserved keywords require quoting in the same positions PostgreSQL requires quoting.
- [x] P2.2.10 Type/function name keywords follow PostgreSQL's `type_func_name_keyword` behavior.
- [x] P2.2.11 Column label keywords follow PostgreSQL's label behavior in view target lists.
- [x] P2.2.12 Object-control statements preserve keyword identifiers in comments.
- [x] P2.2.13 Object-control statements preserve keyword identifiers in grants.
- [x] P2.2.14 Object-control statements preserve keyword identifiers in drops.
- [x] P2.2.15 Object-control statements preserve keyword identifiers in alters.

### 2.3 Parser Dispatch and Tail Safety

- [x] P2.3.01 Invalid `DROP FUNCTION ()` is rejected without creating a ghost statement.
- [x] P2.3.02 Invalid `ALTER FUNCTION f()` with no action is rejected like PostgreSQL.
- [x] P2.3.03 Invalid `COMMENT ON FUNCTION f(integer` reports a syntax error and consumes no trailing ghost statement.
- [x] P2.3.04 Invalid `GRANT EXECUTE ON FUNCTION f(integer TO role` reports a syntax error.
- [x] P2.3.05 Invalid `CREATE TABLE t (FOREIGN KEY (x) REFERENCES p MATCH)` reports a syntax error.
- [x] P2.3.06 Invalid `ALTER INDEX i ATTACH PARTITION` reports a syntax error.
- [x] P2.3.07 Invalid `WITH c AS SELECT 1 SELECT * FROM c` reports a syntax error.
- [x] P2.3.08 Invalid `LATERAL ()` is rejected in FROM.
- [x] P2.3.09 Invalid `LATERAL relation_name` is rejected in FROM.
- [x] P2.3.10 Parser statement count matches PostgreSQL for multi-statement input with an invalid middle statement.
- [x] P2.3.11 Parser locations remain stable for function identity list statements.
- [x] P2.3.12 Parser locations remain stable for CTE and LATERAL view statements.

## Phase 3: Catalog DDL Semantic Compatibility

### 3.1 Foreign Key Type Compatibility

- [x] P3.1.01 `bigint` child to `integer` parent is accepted through cross-type equality.
- [x] P3.1.02 `smallint` child to `integer` parent is accepted through cross-type equality.
- [x] P3.1.03 `integer` child to `bigint` parent is accepted through cross-type equality.
- [x] P3.1.04 `integer` child to `smallint` parent follows PostgreSQL behavior.
- [x] P3.1.05 `smallint` child to `bigint` parent follows PostgreSQL behavior.
- [x] P3.1.06 `bigint` child to `smallint` parent follows PostgreSQL behavior.
- [x] P3.1.07 `numeric` child to `integer` parent follows PostgreSQL behavior.
- [x] P3.1.08 `integer` child to `numeric` parent follows PostgreSQL behavior.
- [x] P3.1.09 `text` child to `varchar` parent follows PostgreSQL behavior.
- [x] P3.1.10 `varchar` child to `text` parent follows PostgreSQL behavior.
- [x] P3.1.11 `bpchar` child to `text` parent follows PostgreSQL behavior.
- [x] P3.1.12 `text` child to `bpchar` parent follows PostgreSQL behavior.
- [x] P3.1.13 Domain child to base parent follows PostgreSQL behavior.
- [x] P3.1.14 Base child to domain parent follows PostgreSQL behavior.
- [x] P3.1.15 Domain child to same-domain parent follows PostgreSQL behavior.
- [x] P3.1.16 Domain child to different-domain parent follows PostgreSQL behavior.
- [x] P3.1.17 Composite foreign key validates each cross-type column pair.
- [x] P3.1.18 Composite foreign key rejects when one column pair has no equality operator.
- [x] P3.1.19 Foreign key to a unique index with cross-type equality is accepted.
- [x] P3.1.20 Foreign key to a primary key with included columns is accepted when key columns match.

### 3.2 Relation Shape and Table DDL

- [x] P3.2.01 Ordinary tables may have zero user columns.
- [x] P3.2.02 Temporary tables may have zero user columns.
- [x] P3.2.03 Unlogged tables may have zero user columns.
- [x] P3.2.04 Partitioned tables may have zero user columns where PostgreSQL permits it.
- [x] P3.2.05 `CREATE TABLE AS SELECT` with zero output columns follows PostgreSQL behavior.
- [x] P3.2.06 `CREATE VIEW v AS SELECT` with zero target columns follows PostgreSQL behavior.
- [x] P3.2.07 Inherited tables preserve inherited zero-column parents.
- [x] P3.2.08 `CREATE TABLE LIKE zero_column_table` follows PostgreSQL behavior.
- [x] P3.2.09 `ALTER TABLE zero_column_table ADD COLUMN` works.
- [x] P3.2.10 `ALTER TABLE table DROP COLUMN` can produce a zero-user-column table where PostgreSQL permits it.

### 3.3 Partition and Index DDL

- [x] P3.3.01 Partitioned index attach records index-to-index dependency state.
- [x] P3.3.02 `ALTER INDEX parent ATTACH PARTITION child` works with schema-qualified index names.
- [x] P3.3.03 `ALTER INDEX parent ATTACH PARTITION child` rejects a child index on the wrong table.
- [x] P3.3.04 `ALTER INDEX parent ATTACH PARTITION child` rejects a parent index on a non-partitioned table.
- [x] P3.3.05 Partitioned parent index dependency survives dump/load order where parent index is created first.
- [x] P3.3.06 Partitioned parent index dependency survives dump/load order where child table is created first.
- [x] P3.3.07 Partitioned parent index dependency survives `CREATE INDEX ON ONLY parent`.
- [x] P3.3.08 Partitioned unique index attach validates key column compatibility.
- [x] P3.3.09 Partitioned expression index attach follows PostgreSQL behavior.
- [x] P3.3.10 Partitioned partial index attach follows PostgreSQL behavior.
- [x] P3.3.11 Partitioned index attach records dependency usable by diff.
- [x] P3.3.12 Partitioned index attach records dependency usable by migration planning.
- [x] P3.3.13 Partitioned table attach records parent-child table dependency.
- [x] P3.3.14 Partitioned table detach follows PostgreSQL behavior where loader supports it.
- [x] P3.3.15 Partition bounds with `DEFAULT` load and preserve parent-child relation.

### 3.4 Constraint and Default Semantics

- [x] P3.4.01 Generated columns validate expression type coercion.
- [x] P3.4.02 Generated columns record dependencies on referenced columns.
- [x] P3.4.03 Identity columns record owned sequence dependencies.
- [x] P3.4.04 Column defaults validate expression type coercion.
- [x] P3.4.05 Column defaults record dependencies on functions.
- [x] P3.4.06 Column defaults record dependencies on sequences.
- [x] P3.4.07 CHECK constraints validate boolean expression coercion.
- [x] P3.4.08 CHECK constraints record function dependencies.
- [x] P3.4.09 Exclusion constraints resolve operator classes.
- [x] P3.4.10 Unique constraints with `NULLS NOT DISTINCT` load.
- [x] P3.4.11 Deferrable unique constraints load.
- [x] P3.4.12 Deferrable foreign keys load.
- [x] P3.4.13 `NOT VALID` foreign keys load.
- [x] P3.4.14 `VALIDATE CONSTRAINT` follows PostgreSQL behavior.
- [x] P3.4.15 Constraint comments survive load.

## Phase 4: Function and Operator Resolver

### 4.1 Variadic Builtins and User Functions

- [x] P4.1.01 `concat_ws('-', text, integer)` resolves through variadic expansion.
- [x] P4.1.02 `jsonb_build_object('id', id, 'name', name)` resolves through variadic expansion.
- [x] P4.1.03 Explicit `VARIADIC array[...]` calls are resolved with PostgreSQL-compatible array handling.
- [x] P4.1.04 User-defined variadic functions resolve when called with expanded arguments.
- [x] P4.1.05 `concat(text, integer, timestamp)` resolves through variadic expansion.
- [x] P4.1.06 `format('%s %s', a, b)` resolves through variadic expansion.
- [x] P4.1.07 `json_build_object` resolves through variadic expansion.
- [x] P4.1.08 `json_build_array` resolves through variadic expansion.
- [x] P4.1.09 `jsonb_build_array` resolves through variadic expansion.
- [x] P4.1.10 User-defined variadic functions resolve when called with zero variadic elements.
- [x] P4.1.11 User-defined variadic functions resolve when called with one variadic element.
- [x] P4.1.12 User-defined variadic functions resolve when called with mixed coercible argument types.
- [x] P4.1.13 Explicit `VARIADIC` rejects non-array arguments like PostgreSQL.
- [x] P4.1.14 Explicit `VARIADIC NULL` follows PostgreSQL behavior.
- [x] P4.1.15 Variadic aggregate calls follow PostgreSQL behavior.

### 4.2 Polymorphic Return and SRF Types

- [x] P4.2.01 `unnest(text[]) AS alias(col)` exposes `col` as `text`.
- [x] P4.2.02 `unnest(integer[]) AS alias(col)` exposes `col` as `integer`.
- [x] P4.2.03 `unnest(bigint[])` exposes `bigint`.
- [x] P4.2.04 `unnest(uuid[])` exposes `uuid`.
- [x] P4.2.05 `unnest(jsonb[])` exposes `jsonb`.
- [x] P4.2.06 `unnest(anyarray)` with a domain array follows PostgreSQL behavior.
- [x] P4.2.07 `array_length(anyarray, int)` resolves for typed arrays.
- [x] P4.2.08 `array_position(anycompatiblearray, anycompatible)` resolves for text arrays.
- [x] P4.2.09 `array_position(anycompatiblearray, anycompatible)` resolves for integer arrays.
- [x] P4.2.10 `array_append(anycompatiblearray, anycompatible)` resolves and returns array type.
- [x] P4.2.11 `array_prepend(anycompatible, anycompatiblearray)` resolves and returns array type.
- [x] P4.2.12 `coalesce(anycompatible, anycompatible)` resolves common type.
- [x] P4.2.13 `greatest(anycompatible, anycompatible)` resolves common type.
- [x] P4.2.14 `least(anycompatible, anycompatible)` resolves common type.
- [x] P4.2.15 Record-returning builtins expose OUT parameter names.
- [x] P4.2.16 Record-returning user functions expose OUT parameter names.
- [x] P4.2.17 RETURNS TABLE user functions expose table column names.
- [x] P4.2.18 Set-returning functions in FROM expose alias column lists.

### 4.3 Unknown Literal and Default Argument Resolution

- [x] P4.3.01 Function calls with defaulted arguments resolve when omitted arguments are within `pronargdefaults`.
- [x] P4.3.02 Unknown string literal resolves to text when a text overload exists.
- [x] P4.3.03 Unknown string literal resolves to uuid when explicitly cast by context.
- [x] P4.3.04 Unknown string literal resolves to jsonb for jsonb operators where PostgreSQL does.
- [x] P4.3.05 Unknown numeric literal resolves to integer for integer-only overloads.
- [x] P4.3.06 Unknown numeric literal resolves to numeric for numeric-preferred overloads.
- [x] P4.3.07 Unknown NULL argument resolves through overload candidate selection.
- [x] P4.3.08 Multiple unknown arguments resolve by preferred type category.
- [x] P4.3.09 Mixed known and unknown arguments resolve through PostgreSQL's last-gasp heuristic.
- [x] P4.3.10 User-defined functions with one default argument resolve.
- [x] P4.3.11 User-defined functions with multiple default arguments resolve.
- [x] P4.3.12 User-defined functions reject omitted non-default arguments.
- [x] P4.3.13 Overloaded user-defined functions with defaults pick the PostgreSQL-compatible candidate.
- [x] P4.3.14 Named argument calls follow PostgreSQL behavior where parser supports them.

### 4.4 Operator Resolution

- [x] P4.4.01 `jsonb -> text` resolves.
- [x] P4.4.02 `jsonb -> int` resolves.
- [x] P4.4.03 `jsonb ->> text` resolves.
- [x] P4.4.04 `jsonb ->> int` resolves.
- [x] P4.4.05 `json -> text` resolves.
- [x] P4.4.06 `json ->> text` resolves.
- [x] P4.4.07 `text || unknown` resolves.
- [x] P4.4.08 `unknown || text` resolves.
- [x] P4.4.09 `integer + bigint` resolves.
- [x] P4.4.10 `bigint + integer` resolves.
- [x] P4.4.11 `numeric + integer` resolves.
- [x] P4.4.12 `date + integer` resolves.
- [x] P4.4.13 `timestamp + interval` resolves.
- [x] P4.4.14 `text LIKE unknown` resolves.
- [x] P4.4.15 `text ILIKE unknown` resolves.
- [x] P4.4.16 `array @> array` resolves for integer arrays.
- [x] P4.4.17 `jsonb @> jsonb` resolves.
- [x] P4.4.18 Range operators resolve for `int4range`.

## Phase 5: View Analyzer Scope

### 5.1 CTE Range Resolution

- [x] P5.1.01 A simple CTE is visible to the main query body.
- [x] P5.1.02 CTE column aliases are visible with their alias names and inferred types.
- [x] P5.1.03 Multiple CTEs in one WITH list can reference earlier CTEs.
- [x] P5.1.04 A CTE can reference a base table with the same column names.
- [x] P5.1.05 A CTE name shadows a base relation in the same search path.
- [x] P5.1.06 A nested subquery can reference an outer visible CTE where PostgreSQL permits it.
- [x] P5.1.07 A nested WITH can shadow an outer CTE.
- [x] P5.1.08 A nested WITH can reference a prior outer CTE when not shadowed.
- [x] P5.1.09 Recursive CTE non-recursive term establishes column types.
- [x] P5.1.10 Recursive CTE recursive term can reference the CTE name.
- [x] P5.1.11 Recursive CTE rejects invalid non-UNION shapes.
- [x] P5.1.12 CTE with `MATERIALIZED` loads.
- [x] P5.1.13 CTE with `NOT MATERIALIZED` loads.
- [x] P5.1.14 CTE with explicit column aliases fewer than output columns follows PostgreSQL behavior.
- [x] P5.1.15 CTE with explicit column aliases more than output columns follows PostgreSQL behavior.
- [x] P5.1.16 CTE in a view records dependencies on base relations.
- [x] P5.1.17 CTE in a view records dependencies on base columns.
- [x] P5.1.18 CTE in a view records dependencies on functions used inside the CTE.

### 5.2 LATERAL and Correlated Scope

- [x] P5.2.01 A later lateral subquery can reference a preceding lateral alias.
- [x] P5.2.02 `LEFT JOIN LATERAL` can reference the left relation.
- [x] P5.2.03 `CROSS JOIN LATERAL` can reference the left relation.
- [x] P5.2.04 `INNER JOIN LATERAL` can reference the left relation.
- [x] P5.2.05 A lateral subquery can reference multiple preceding FROM items.
- [x] P5.2.06 A lateral subquery can reference a preceding lateral alias.
- [x] P5.2.07 A third lateral item can reference two earlier lateral aliases.
- [x] P5.2.08 A lateral set-returning function can reference the left relation.
- [x] P5.2.09 `LATERAL unnest(t.arr) AS u(val)` exposes alias column type.
- [x] P5.2.10 `LATERAL generate_series(1, t.n) AS g(x)` exposes alias column type.
- [x] P5.2.11 `ROWS FROM (...) WITH ORDINALITY` exposes ordinality.
- [x] P5.2.12 LATERAL with column definition list follows PostgreSQL behavior.
- [x] P5.2.13 Non-lateral subquery cannot reference preceding FROM items.
- [x] P5.2.14 LATERAL inside a joined table has the correct left-side visibility.
- [x] P5.2.15 LATERAL under nested joins respects PostgreSQL join scope.

### 5.3 Correlated Subqueries

- [x] P5.3.01 Correlated scalar subquery in SELECT list carries `LevelsUp` safely.
- [x] P5.3.02 Correlated scalar subquery in WHERE carries `LevelsUp` safely.
- [x] P5.3.03 Correlated `EXISTS` subquery in WHERE carries `LevelsUp` safely.
- [x] P5.3.04 Correlated `IN (SELECT ...)` subquery carries `LevelsUp` safely.
- [x] P5.3.05 Correlated subquery in JOIN condition carries `LevelsUp` safely.
- [x] P5.3.06 Correlated subquery in CHECK expression follows PostgreSQL behavior.
- [x] P5.3.07 Correlated subquery in generated column expression follows PostgreSQL behavior.
- [x] P5.3.08 Nested correlated subqueries with `LevelsUp=2` do not panic walkers.
- [x] P5.3.09 Dependency walker records base relation dependencies for correlated refs.
- [x] P5.3.10 Provenance walker ignores outer vars when local target provenance is expected.
- [x] P5.3.11 Ruleutils deparses correlated references with the correct alias.
- [x] P5.3.12 Query span tooling handles outer vars without local range-table panics.

### 5.4 View Target and Range Function Shapes

- [x] P5.4.01 View over `SELECT * FROM unnest(array)` has stable column names.
- [x] P5.4.02 View over `SELECT * FROM unnest(array) AS u(val)` has alias column name.
- [x] P5.4.03 View over `SELECT * FROM unnest(array1, array2) AS u(a,b)` has both alias columns.
- [x] P5.4.04 View over record-returning function with column definition list loads.
- [x] P5.4.05 View over RETURNS TABLE user function exposes table columns.
- [x] P5.4.06 View over OUT-parameter user function exposes OUT names.
- [x] P5.4.07 View over `jsonb_to_record` with column definition list loads.
- [x] P5.4.08 View over `jsonb_to_recordset` with column definition list loads.
- [x] P5.4.09 View over `generate_series` with alias loads.
- [x] P5.4.10 View over set-returning function in SELECT target follows PostgreSQL behavior.
- [x] P5.4.11 View over star expansion through CTE preserves column order.
- [x] P5.4.12 View over star expansion through LATERAL preserves column order.
- [x] P5.4.13 View over duplicate output column names follows PostgreSQL behavior.
- [x] P5.4.14 View column alias list overrides target names.
- [x] P5.4.15 View column alias list count mismatch follows PostgreSQL behavior.

## Phase 6: Object Namespace and Dependency

### 6.1 Typed Object Resolution

- [x] P6.1.01 `ALTER INDEX ... ATTACH PARTITION ...` resolves indexes from `schema.Indexes`.
- [x] P6.1.02 `ALTER INDEX ... DEPENDS ON EXTENSION ...` resolves the index object without relation namespace fallback.
- [x] P6.1.03 `ALTER INDEX ... RENAME TO ...` resolves indexes from `schema.Indexes`.
- [x] P6.1.04 `ALTER SEQUENCE ... RENAME TO ...` resolves sequences from `schema.Sequences`.
- [x] P6.1.05 `ALTER SEQUENCE ... OWNER TO ...` resolves sequences from `schema.Sequences`.
- [x] P6.1.06 `ALTER VIEW ... RENAME TO ...` resolves relation namespace with view relkind validation.
- [x] P6.1.07 `ALTER MATERIALIZED VIEW ... RENAME TO ...` resolves relation namespace with matview relkind validation.
- [x] P6.1.08 `ALTER TYPE ... RENAME TO ...` resolves type namespace.
- [x] P6.1.09 `ALTER DOMAIN ... RENAME TO ...` resolves type namespace.
- [x] P6.1.10 `ALTER FUNCTION ... RENAME TO ...` resolves `ObjectWithArgs`.
- [x] P6.1.11 `ALTER PROCEDURE ... RENAME TO ...` resolves `ObjectWithArgs`.
- [x] P6.1.12 `ALTER ROUTINE ... RENAME TO ...` resolves `ObjectWithArgs`.
- [x] P6.1.13 `COMMENT ON INDEX` resolves index namespace.
- [x] P6.1.14 `COMMENT ON SEQUENCE` resolves sequence namespace.
- [x] P6.1.15 `COMMENT ON TYPE` resolves type namespace.
- [x] P6.1.16 `COMMENT ON FUNCTION` resolves function identity by type.
- [x] P6.1.17 `DROP INDEX` resolves index namespace.
- [x] P6.1.18 `DROP SEQUENCE` resolves sequence namespace.
- [x] P6.1.19 `DROP TYPE` resolves type namespace.
- [x] P6.1.20 `DROP FUNCTION` resolves function identity by type.

### 6.2 Object-Control Cross Product

- [x] P6.2.01 COMMENT, GRANT, DROP, ALTER, and RENAME share object identity helpers for common target types.
- [x] P6.2.02 `GRANT SELECT ON TABLE` resolves table targets.
- [x] P6.2.03 `GRANT SELECT ON VIEW` resolves view targets.
- [x] P6.2.04 `GRANT USAGE ON SEQUENCE` resolves sequence targets.
- [x] P6.2.05 `GRANT EXECUTE ON FUNCTION` resolves function targets by identity.
- [x] P6.2.06 `REVOKE` mirrors each GRANT target resolver.
- [x] P6.2.07 `DROP OWNED BY` follows PostgreSQL behavior for loaded objects.
- [x] P6.2.08 `ALTER OWNER` follows PostgreSQL behavior for loaded objects.
- [x] P6.2.09 `ALTER SET SCHEMA` follows PostgreSQL behavior for loaded objects.
- [x] P6.2.10 `COMMENT IS NULL` removes comments for loaded objects.

### 6.3 Dependency Integrity

- [x] P6.3.01 Attached child indexes record an internal dependency on the parent index.
- [x] P6.3.02 View dependencies are recorded for base relations referenced directly.
- [x] P6.3.03 View dependencies are recorded for base columns referenced directly.
- [x] P6.3.04 View dependencies are recorded for objects referenced through CTE expansions.
- [x] P6.3.05 View dependencies are recorded for objects referenced through LATERAL expansions.
- [x] P6.3.06 Function dependencies are recorded for argument types.
- [x] P6.3.07 Function dependencies are recorded for return types.
- [x] P6.3.08 Function dependencies are recorded for default expressions.
- [x] P6.3.09 Generated column dependencies are recorded for referenced columns.
- [x] P6.3.10 CHECK expression dependencies are recorded for functions and columns.
- [x] P6.3.11 Policy expression dependencies are recorded for functions and columns.
- [x] P6.3.12 Trigger dependencies are recorded for trigger function and table.
- [x] P6.3.13 Index dependencies are recorded for indexed columns.
- [x] P6.3.14 Expression index dependencies are recorded for functions and columns.
- [x] P6.3.15 Partial index dependencies are recorded for predicate functions and columns.
- [x] P6.3.16 Partitioned index dependency state survives schema diff.
- [x] P6.3.17 Partitioned index dependency state survives migration planning.
- [x] P6.3.18 Dependency-driven drop planning follows PostgreSQL cascade expectations.
- [x] P6.3.19 Dependency-driven rename planning preserves object identity.
- [x] P6.3.20 Dependency snapshots are stable across load/deparse/reload.
