# Omni Engine Capability Guide

This document summarizes the entire journey of building the PostgreSQL engine from scratch to full Bytebase integration. It distills the capability layers each engine must provide, mandatory quality requirements, and hard-won lessons learned — serving as the canonical reference for replicating this work across MySQL, MSSQL, and Oracle.

---

## 1. Capability Layer Model

Each engine's complete capability stack consists of seven layers, where upper layers depend on lower ones. **The layer order is the build order** — skipping layers causes rework.

```
L7  Bytebase Integration   ← Replace ANTLR, wire into product features
L6  Completion              ← SQL autocompletion
L5  Catalog                 ← DDL simulation + Diff + Migration + SDL
L4  PL / Stored Procedures  ← Parsing procedural language inside function/procedure bodies
L3  Error Quality           ← Truncated SQL must not crash; errors must be meaningful
L2  Loc Position Tracking   ← Byte-precise Start + End on every AST node
L1  Parser                  ← Recursive descent, full grammar coverage
```

### L1 Parser — Recursive Descent Parser

**Responsibility**: Convert SQL text into a typed AST.

**Mandatory requirements**:
- Pure Go recursive descent, zero external dependencies
- Cover the target database engine's **complete grammar** (not just a DDL subset)
- Every parse function returns `(Node, error)` (dual return) — never use nil as a substitute for error
- Support multi-statement parsing (`stmtblock → stmtmulti → stmt`)
- Entry point: `Parse(sql string) → ([]Statement, error)`

**Lessons from PG**:
- The PG parser originally returned bare Node; dual return was retrofitted later (post-hoc patches across 78 batches). **New engines must use `(Node, error)` from day one**
- It took 9 audit rounds to cover all grammar edges (ALTER xxx OWNER TO / RENAME TO / SET SCHEMA were missed repeatedly). **At least one systematic audit is required after the initial parser is complete**
- The PG parser has 78 batches and ~29K lines of code. MSSQL is similar in scale. MySQL has 123 batches. This is the normal volume

**Verification methods**:
- PROGRESS.json tracks each batch's status
- Coverage audit against the target database's BNF / official grammar documentation
- pgregress / official test suite compatibility (when available)

### L2 Loc Position Tracking

**Responsibility**: Every AST node carries byte-precise position information via `Loc{Start, End}`.

**Mandatory requirements**:
- **Both Start and End must be set** — Completion needs End to determine which node the cursor is in; Advisor needs Start to report error line numbers
- Use sentinel value `-1` for unknown positions; provide a `NoLoc()` constructor
- End position is the end of the last token (`p.prevEnd()`), not the start of the next token

**Lessons from PG**:
- PG initially only set Start; 8 batches (63–71) were later added specifically to fix End positions. **This caused significant rework**
- MSSQL learned from this and planned a complete 6-phase Loc tracking scheme from the start — zero rework
- Two unfinished nodes (CreateTableSpaceStmt, ObjectWithArgs) remain to this day

**Verification methods**:
- `CheckLocations()` test infrastructure: walks all AST nodes and asserts Loc is not zero-valued
- SCENARIOS-xx-loc.md defines position tracking scenarios per node type

### L3 Error Quality

**Responsibility**: Truncated/incomplete SQL must not cause panics or infinite loops; error messages must be meaningful.

**Mandatory requirements**:
- **Soft-fail**: After every `advance()` + `parseXxx()` combination, if parseXxx returns nil/error, the parser must degrade gracefully rather than continue with a nil node
- **Error position**: Error messages must include line and column numbers
- **Error context**: Use a format similar to PostgreSQL's `"at or near \"xxx\""` so the user knows where the parser got stuck
- **No crashes**: Any input (empty string, binary data, SQL truncated at any position) must never panic

**Lessons from PG**:
- Soft-fail covered 5 phases: expression binary operators → pattern matching → b_expr → JOIN/GROUP BY → DDL constraints
- Truncated SQL is a daily scenario for Completion (user is still typing); the parser must work correctly in this situation
- MSSQL's error alignment was organized into 3 phases: dual return migration → soft-fail → error message enhancement

**Verification methods**:
- SCENARIOS-soft-fail.md defines truncation scenarios
- Fuzz testing: randomly truncate valid SQL and verify no panics
- Compare against the target database's actual error messages (error alignment)

### L4 PL / Stored Procedure Parsing

**Responsibility**: Parse the procedural language inside CREATE FUNCTION / CREATE PROCEDURE bodies.

**Mandatory requirements**:
- Separate parser module (`plpgsql/parser/`, `mysql/parser/compound.go`, etc.)
- Produce independent AST node types (`plpgsql/ast/`)
- Identify SQL statements within the body and delegate to the main parser
- Export `Parse(body string)` for upstream callers

**Lessons from PG**:
- PL/pgSQL has 208 scenarios and is a standalone recursive descent parser
- PL/pgSQL QuerySpan is a separate step in the Bytebase integration — it cannot be merged into the main QuerySpan migration
- MySQL's compound statements (BEGIN...END, IF, LOOP, etc.) are already implemented in `compound.go` within the main parser but not yet independently exported

**Applicability**: PG (PL/pgSQL), MySQL (compound statements), MSSQL (T-SQL batch/proc), Oracle (PL/SQL). Each engine has its own procedural language.

### L5 Catalog — DDL Simulation Engine

This is the most complex layer, comprising four subsystems.

#### L5.1 Catalog Build (DDL Execution Simulation)

**Responsibility**: Maintain an in-memory database schema state; receive DDL statements and update the state.

**Mandatory requirements**:
- `Exec(stmt ast.Node)` interface: accepts an AST node and updates internal state
- Support all DDL object types: schema, table, column, constraint, index, view, function/procedure, trigger, sequence, enum, domain, extension, policy, grant, comment
- Built-in type system: type OIDs, implicit/explicit casts, operators, built-in functions
- Name resolution: unqualified name → schema.name, search_path support

**PG implementation scale**: 183 source files including generated files (pg_proc_generated.go 896KB, pg_operator_generated.go 127KB, pg_type_generated.go 75KB). These generated files are translated from PG source code and contain 193 types, 229 casts, 799 operators, 3,314 functions.

**Lessons from PG**:
- Aggregates and regular functions must be treated differently, otherwise schema load fails (BYT-9076)
- Serial/Identity types implicitly create sequences; this association must be maintained in the catalog
- `GENERATED ALWAYS AS (expr) STORED` has a different internal representation than plain `DEFAULT`

#### L5.2 Diff Engine

**Responsibility**: Compare two catalog states (current vs desired) and produce a list of differences.

**Mandatory requirements**:
- Modularize by object type: `diff_table.go`, `diff_column.go`, `diff_constraint.go`, `diff_index.go`, `diff_view.go`, `diff_function.go`, `diff_trigger.go`, `diff_sequence.go`, `diff_enum.go`, `diff_domain.go`, `diff_schema.go`, `diff_grant.go`, `diff_policy.go`, `diff_extension.go`, `diff_comment.go`, `diff_partition.go`, etc.
- Output a structured list of DiffOps, not raw SQL

**PG scale**: 131 scenarios, 100% complete. Covers add/modify/drop across 18 object types.

#### L5.3 Migration Generation

**Responsibility**: Convert Diff results into an executable sequence of DDL SQL statements.

**Mandatory requirements**:
- Statement ordering must respect dependencies (CREATE schema first, then CREATE table, then CREATE index)
- DROP order must be reversed (DROP index first, then DROP table)
- Handle special cases: serial absorption, identity columns, generated columns, partition tables
- Output formatted SQL text

**PG scale**: 120 scenarios, 98.3% complete.

**Lessons from PG**:
- **DEFERRABLE and EXCLUDE constraints** are easily missed — their DDL syntax differs from regular constraints
- **Column comments** require `COMMENT ON COLUMN table.column`, not `COMMENT ON COLUMN column` (column name qualification is easy to get wrong)
- **Identity columns** have ordering pitfalls: DROP SEQUENCE must come after DROP IDENTITY
- **Extension ordering** must come before types and tables

#### L5.4 SDL (Schema Definition Language)

**Responsibility**: Accept a set of DDL statements as a declarative schema definition and load them into catalog state.

**Mandatory requirements**:
- Name resolution: handle unqualified references
- Dependency extraction: extract dependencies from expressions (CHECK, DEFAULT, view body)
- Topological sort + cycle repair: handle circular FKs and similar situations
- Validation: detect conflicts, duplicate definitions, invalid references

**PG scale**: 82 scenarios, 100% complete. Three phases: validation → dependency resolution → topological sort.

**Lessons from PG**:
- **Circular FKs are the hardest edge case** — A references B and B references A requires cycle repair (create tables without FKs first, then ALTER TABLE ADD FK)
- **Functions referenced by CHECK/view/trigger simultaneously** — sorting must ensure the function comes before all consumers
- Inline FKs (REFERENCES within a column definition) cause LoadSDL failures in circular scenarios

#### L5.5 Oracle Verification (Comparison Against Real Database)

**Responsibility**: Spin up a real database via testcontainer, execute migration DDL, and compare the resulting state against expectations.

**Mandatory requirements**:
- Use Docker testcontainers (PG uses postgres:17, MySQL uses mysql:8.0)
- Test structure: apply desired DDL → build catalog → diff → generate migration → apply to real DB → verify
- Organize by phase: infrastructure → single objects → combinations → regression → real-world scenarios

**PG scale**: 110 scenarios, 90% complete. Discovered 11 migration/SDL bugs that pure unit tests could not catch.

**Lessons from PG**:
- **Dependent view rebuild**: When changing a column type, if a view depends on that column, migration must DROP the view first and then recreate it. Pure diff tests won't catch this — only oracle tests expose it
- **SERIAL and SEQUENCE association**: Sequences implicitly created by `CREATE TABLE ... serial` cannot be independently DROPped
- **Comment types**: Table comments and column comments have different DDL formats
- **Trigger UPDATE OF columns**: The column list in `CREATE TRIGGER ... UPDATE OF col1, col2` is easily lost

### L6 Completion — SQL Autocompletion

**Responsibility**: Given a cursor position and existing catalog state, provide completion candidates.

**Mandatory requirements**:
- Parser-embedded completion mode: record the syntactic context at the cursor position during parsing
- CandidateSet collection: keywords, table names, column names, function names, schema names, type names
- Table reference extraction: extract visible tables from FROM/JOIN clauses in the current query
- Tricky completion fallback: when the parser can't determine context due to truncation, fall back to heuristic completion

**Lessons from PG**:
- Completion has hard dependencies on L2 (Loc End) and L3 (Soft-fail). **If either is incomplete, Completion cannot function**
- Recursive CTE self-references, materialized view context filtering, and similar scenarios are very difficult — these can be marked as partial
- Integration tests verifying real completion behavior across multi-table schemas are essential

### L7 Bytebase Integration

**Responsibility**: Replace ANTLR-based parsing in Bytebase with omni capabilities.

**Integration order** (validated path from PG):

```
 1. Add omni dependency + SQL Split + AST interface adapter
 2. Statement type detection (query type classification)
 3. Catalog Walkthrough (replace ANTLR walk-through)
 4. Advisor migration (in batches, ~10 advisors per batch)
 5. QuerySpan migration
 6. Completion migration
 7. Backup/Restore migration
 8. Diagnose migration
 9. SDL Diff/Migration migration
10. SDL Checker migration
11. Remove old parser dependency
```

**Mandatory requirements**:
- **Build the bridge before burning it** — During the coexistence period, verify omni and ANTLR feature parity one-by-one before replacing
- **Each step is an independent PR** — Must be revertable, with corresponding tests
- **Remove ANTLR only as the last step** — After confirming all features are migrated and tests pass

**Lessons from PG**:
- PG used 21 PRs to complete integration, spanning 5 packages: `parser/pg/`, `schema/pg/`, `advisor/pg/`, `db/pg/`, `taskrun/`
- Nested SELECT in wrapper statements (subqueries within CTEs) was missed during advisor migration
- PL/pgSQL QuerySpan requires separate handling — it cannot be merged into the main QuerySpan PR

---

## 2. Scenario-Driven Development Methodology

All capabilities above are built using the **scenario-driven development** pattern.

### Process

```
1. Define SCENARIOS-{engine}-{capability}.md
   - Organize by phase/section
   - One checkbox per scenario: - [ ] scenario description

2. Create Driver skill (orchestration) and Worker skill (execution)
   - Driver manages progress and dispatches sections
   - Worker executes one section: implement + write tests + run verification + fix + update progress

3. Advance section by section
   - Mark [x] when a section is complete
   - Mark [~] for discovered bugs with explanation

4. Oracle verification (when applicable)
   - Compare against real database behavior
   - Feed discovered differences back into implementation
```

### Naming conventions

```
SCENARIOS-{engine}-{capability}.md    # Scenario definitions
PROGRESS.json                          # Automated progress tracking (parser batches)
PROGRESS_SUMMARY.json                  # Progress summary
```

### Validated scenario volumes

| Capability | Typical Scenario Count | Organization |
|------------|----------------------|--------------|
| Parser | 50–130 batches | By grammar construct |
| Loc | 100–140 | By AST node type |
| Soft-fail | 40–80 | By syntactic context phase |
| Diff | 100–130 | Object type × operation type |
| Migration | 100–120 | Mirrors Diff scenarios |
| SDL | 70–90 | Validation + dependencies + sorting |
| Oracle | 100–120 | Full roundtrip verification |
| Completion | 200–270 | Statement type × clause position |
| PL / Stored procedures | 150–200+ | By statement type |

---

## 3. Current Status by Engine

### PostgreSQL — The Benchmark

| Layer | Status | Scenarios | Notes |
|-------|--------|-----------|-------|
| L1 Parser | ✅ 100% | 78 batches | 9 audit rounds |
| L2 Loc | ✅ 98.5% | 135/137 | 2 residual nodes |
| L3 Error | ✅ 100% | 5 phases | |
| L4 PL/pgSQL | ✅ 100% | 208 | |
| L5.1 Catalog | ✅ ~90% | — | 183 source files |
| L5.2 Diff | ✅ 100% | 131 | |
| L5.3 Migration | ✅ 98.3% | 118/120 | 2 extension scenarios |
| L5.4 SDL | ✅ 100% | 82 | |
| L5.5 Oracle | ✅ 90% | 99/110 | 11 known bugs |
| L6 Completion | ✅ ~95% | 100+ | |
| L7 Integration | ✅ 100% | 21 PRs | ANTLR fully removed |

### MySQL — Missing L5.2+

| Layer | Status | Scenarios | Notes |
|-------|--------|-----------|-------|
| L1 Parser | ✅ 100% | 123 batches | 22 audit rounds |
| L2 Loc | ⚠️ Partial | — | Loc exists but End not systematically tracked |
| L3 Error | ✅ 100% | 4 phases | |
| L4 Stored Procs | ⚠️ Partial | — | compound.go exists but not independently exported |
| L5.1 Catalog | ✅ Extensive | oracle verified | |
| L5.2 Diff | ❌ Missing | — | **Needs to be built** |
| L5.3 Migration | ❌ Missing | — | **Needs to be built** |
| L5.4 SDL | ❌ Missing | — | **Needs to be built** |
| L5.5 Oracle | ✅ Partial | Catalog verified | |
| L6 Completion | ✅ 100% | 9 phases | |
| L7 Integration | ❌ Not started | — | |

### MSSQL — All of L5 Missing

| Layer | Status | Scenarios | Notes |
|-------|--------|-----------|-------|
| L1 Parser | ✅ Complete | — | 70+ Go files |
| L2 Loc | ✅ 100% | 6 phases | Planned from scratch, zero rework |
| L3 Error | ✅ 100% | 3 phases | Dual return + soft-fail + messages |
| L4 Stored Procs | ❌ Missing | — | T-SQL batch/proc |
| L5.1 Catalog | ❌ Missing | — | **Needs to be built from scratch** |
| L5.2 Diff | ❌ Missing | — | |
| L5.3 Migration | ❌ Missing | — | |
| L5.4 SDL | ❌ Missing | — | |
| L5.5 Oracle | ❌ Missing | — | |
| L6 Completion | ✅ 100% | 9 phases | |
| L7 Integration | ❌ Not started | — | |

### Oracle — Earliest Stage

| Layer | Status | Notes |
|-------|--------|-------|
| L1 Parser | 🔧 In progress | Corpus verifier running |
| L2 Loc | 🔧 In progress | Foundation alignment |
| L3–L7 | ❌ Not started | |

---

## 4. Mandatory Engineering Requirements

The following rules are distilled from PG experience and **must be strictly followed**. Violating any of them will cause rework downstream.

### 4.1 Parser Layer

1. **All parse functions return `(Node, error)`** — never return a bare Node
2. **All AST nodes carry `Loc{Start, End}`** — End is set using `p.prevEnd()`
3. **Use `-1` sentinel for unknown positions** — do not use zero values (zero is a valid position)
4. **Multi-statement entry point**: `Parse()` must handle semicolon-separated multiple statements
5. **Zero dependencies**: do not introduce ANTLR, yacc, or any external parser generator

### 4.2 Error Quality

6. **No panics on any input**: empty strings, binary data, SQL truncated at any position must all return errors gracefully
7. **Errors include position**: line number + column number
8. **Nil check after advance**: every `parseXxx()` call after `advance()` must check for nil/error

### 4.3 Catalog Layer

9. **Exec interface**: `Exec(node) error` — accepts only AST nodes, not SQL text
10. **Complete object type coverage**: must not stop at table/column — must cover index, constraint, view, function, trigger, sequence, enum/domain, extension, policy, grant, comment
11. **Type system**: must have a built-in type registry, cast rules, and operator resolution
12. **Name resolution**: search_path / default schema support

### 4.4 Diff + Migration

13. **Diff produces structured DiffOps** — does not generate SQL directly
14. **Migration must topologically sort**: CREATE in dependency order, DROP in reverse dependency order
15. **Special object associations**: serial→sequence, identity columns, generated columns require special handling
16. **Oracle verification is mandatory**: pure unit tests cannot catch dependent view rebuilds, constraint ordering issues, etc.

### 4.5 SDL

17. **Circular dependency handling**: must implement cycle detection + repair (create tables without circular FKs first, then ALTER TABLE ADD FK)
18. **Expression dependency extraction**: objects referenced in CHECK expressions, view bodies, and trigger bodies must be extracted as dependencies

### 4.6 Completion

19. **Depends on Loc.End and Soft-fail**: if either is incomplete, completion is unusable
20. **Parser-embedded completion mode**: do not use external completion algorithms (e.g., C3); collect candidates during the parse process itself

### 4.7 Integration

21. **Build the bridge before burning it**: omni and the old parser coexist; verify and replace features one-by-one; remove the old parser only as the final step
22. **One PR per functional module**: Split, TypeDetection, Walkthrough, Advisor, QuerySpan, Completion, Backup, Diagnose, SDL, Checker

---

## 5. Key Lessons Learned

### Decisions That Proved Right

1. **Recursive descent parser** — Compared to ANTLR: 250–1000x faster parsing, memory usage drops from MB to KB level, agent-development-friendly (fast feedback loops, controllable blast radius)
2. **Scenario-driven development** — SCENARIOS.md + Driver/Worker skill pattern makes systematic progress across hundreds of scenarios feasible
3. **Oracle verification** — Testcontainer comparison against real databases; PG oracle tests uncovered 11 bugs that pure unit tests could not catch
4. **Catalog before everything else** — Diff/Migration/SDL/Advisor/QuerySpan all depend on Catalog; this ordering is non-negotiable
5. **Multiple parser audit rounds** — Covering all grammar edges in one pass is impossible; 9 audit rounds is normal

### Lessons Paid For With Rework

1. **Loc.End retrofitted** — PG retrofitted 8 batches to fix End positions, touching many functions across the parser. MSSQL planned from the start and had zero rework
2. **Dual return retrofitted** — Changing from `Node` to `(Node, error)` is a large-scale refactor. New engines should adopt this from day one
3. **Aggregates are not functions** — Aggregates must be excluded during function sync; otherwise schema load fails. Subtle type system differences can only be caught by oracle tests
4. **Circular FK handling in SDL** — Mutual FK references causing LoadSDL infinite loops were not discovered until the oracle testing phase. Cycle repair should be planned during SDL design
5. **Nested SELECT oversight** — Subqueries within CTEs and wrapper statements (EXPLAIN, PREPARE) were missed during advisor migration. Recursive processing is required

### Scale Reference

| Engine | Parser Lines | Catalog Files | Total Scenarios | Integration PRs |
|--------|-------------|---------------|-----------------|-----------------|
| PG | ~29K | 183 | 901 | 21 |
| MySQL | ~20K (est.) | ~30 (current) | 700+ (current) | — |
| MSSQL | ~25K (est.) | 0 | 985+ (current) | — |

---

## 6. Appendix: PG SCENARIOS File Index

| File | Scenarios | Completion | Coverage |
|------|-----------|------------|---------|
| `SCENARIOS-pg-loc.md` | 137 | 98.5% | AST node position tracking |
| `SCENARIOS-pg-diff.md` | 131 | 100% | Schema diff detection |
| `SCENARIOS-pg-migration.md` | 120 | 98.3% | DDL migration generation |
| `SCENARIOS-pg-oracle.md` | 110 | 90% | Real PG verification |
| `SCENARIOS-pg-sdl.md` | 82 | 100% | Declarative schema loading |
| `SCENARIOS-pg-expr.md` | 61 | 29.5% | Expression deparsing |
| `SCENARIOS-pg-ruleutils.md` | 52 | 0% | ruleutils alignment |
| `SCENARIOS-plpgsql.md` | 208 | 100% | PL/pgSQL parsing |
| `pg/parser/SCENARIOS-soft-fail.md` | ~60 | 100% | Error handling quality |
