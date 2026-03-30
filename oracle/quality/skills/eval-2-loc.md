# Eval Worker — Stage 2: Loc Completeness

You are an Eval Worker in the Oracle Quality Pipeline.
Your role is to write **tests ONLY** — never implementation code.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/strategy.md` — Stage 2 scope (248 Loc-bearing AST node types)
- `oracle/quality/prevention-rules.md` — **MUST read before starting any work**
- `oracle/ast/parsenodes.go` — all AST node types with Loc fields
- `oracle/ast/node.go` — `Loc` struct, `NoLoc()`, `Node` interface
- `oracle/parser/loc_walker_test.go` — `CheckLocations()` walker you MUST use
- `oracle/parser/parser.go` — `Parse()` entry point
- `oracle/quality/coverage/stage1-foundation.json` — Stage 1 coverage (context)

## Output Files

- **Test file:** `oracle/parser/eval_loc_test.go`
- **Coverage report:** `oracle/quality/coverage/stage2-loc.json`

## Rules

1. **Tests ONLY** — do NOT modify any non-test `.go` file.
2. Every test function MUST be named `TestEvalStage2_*`.
3. Tests should fail clearly with descriptive messages (not just compile errors).
4. Use `CheckLocations()` from `loc_walker_test.go` to verify Loc completeness.
5. For each of the 248 Loc-bearing AST node types, find the **simplest SQL** that produces that node.
6. Parse the SQL, then call `CheckLocations()` and assert zero violations.
7. Do NOT import packages that do not exist yet.
8. Never modify existing test files (`compare_test.go`, `loc_walker_test.go`, etc.).

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[EVAL-STAGE2] STARTED
[EVAL-STAGE2] STEP reading_refs - Reading AST node types and Loc walker
[EVAL-STAGE2] STEP enumerating - Enumerating 248 Loc-bearing node types
[EVAL-STAGE2] STEP writing_tests_select - Writing tests for SELECT-related nodes
[EVAL-STAGE2] STEP writing_tests_dml - Writing tests for DML nodes (INSERT, UPDATE, DELETE, MERGE)
[EVAL-STAGE2] STEP writing_tests_ddl - Writing tests for DDL nodes (CREATE, ALTER, DROP)
[EVAL-STAGE2] STEP writing_tests_plsql - Writing tests for PL/SQL nodes
[EVAL-STAGE2] STEP writing_tests_misc - Writing tests for remaining node types
[EVAL-STAGE2] STEP build_check - Running go build on test file
[EVAL-STAGE2] STEP coverage_report - Generating stage2-loc.json
[EVAL-STAGE2] DONE
```

If a step fails:
```
[EVAL-STAGE2] FAIL step_name - description
[EVAL-STAGE2] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Test Structure

Group tests by SQL category. Each test function exercises one or more node types.

### Category: SELECT

```go
func TestEvalStage2_Select_Basic(t *testing.T) {
    // Exercises: SelectStmt, ResTarget, ColumnRef, RangeVar, etc.
    violations := CheckLocations(t, "SELECT a FROM t")
    if len(violations) > 0 {
        for _, v := range violations {
            t.Errorf("Loc violation: %s", v)
        }
    }
}
```

### Category: DML (INSERT, UPDATE, DELETE, MERGE)

Write the simplest SQL for each DML-related node type.

### Category: DDL (CREATE TABLE, ALTER TABLE, CREATE INDEX, etc.)

Write the simplest SQL for each DDL-related node type.

### Category: PL/SQL (BEGIN..END, DECLARE, IF, LOOP, etc.)

Write the simplest SQL for each PL/SQL-related node type.

### Category: Expressions (CASE, subquery, function call, etc.)

Write the simplest SQL for expression node types.

### Category: Miscellaneous

Cover any remaining node types not in the above categories.

## Test Pattern

Every test follows this pattern:

```go
func TestEvalStage2_Category_NodeType(t *testing.T) {
    sql := "..." // simplest SQL that produces the target node type
    violations := CheckLocations(t, sql)
    if len(violations) > 0 {
        for _, v := range violations {
            t.Errorf("Loc violation: %s", v)
        }
    }
}
```

## How to Find the Simplest SQL for Each Node Type

1. Read `oracle/ast/parsenodes.go` to identify all struct types with a `Loc ast.Loc` field.
2. For each such type, determine which SQL syntax produces it by reading:
   - The parser functions in `oracle/parser/` that construct that node type.
   - The BNF rule files in `oracle/parser/bnf/` for grammar reference.
3. Write the **shortest** SQL that causes the parser to create at least one instance of that node type.
4. If a single SQL statement produces multiple node types, that is fine — one test can cover several types. But every type must appear in at least one test.

## Enumerating Loc-bearing Node Types

To find all 248 Loc-bearing types, run:

```bash
cd /Users/rebeliceyang/Github/omni
grep -c 'Loc  *ast\.Loc\|Loc  *Loc' oracle/ast/parsenodes.go
grep 'type .* struct' oracle/ast/parsenodes.go | head -5
```

Use `reflect` in tests to dynamically verify which types have Loc fields if needed.

## Coverage Report Format

After writing tests, generate `oracle/quality/coverage/stage2-loc.json`:

```json
{
  "stage": "2-loc",
  "total_items": 248,
  "tested_items": 0,
  "items": [
    {"node_type": "SelectStmt", "test": "TestEvalStage2_Select_Basic", "status": "written"},
    {"node_type": "ResTarget", "test": "TestEvalStage2_Select_Basic", "status": "written"},
    {"node_type": "ColumnRef", "test": "TestEvalStage2_Select_Basic", "status": "written"},
    ...
  ]
}
```

The `status` field transitions: `"untested"` → `"written"` → `"passing"` → `"verified"`.

Update `tested_items` to reflect how many node types have at least one test.

## Verification

After writing the test file:

```bash
# Must compile (tests may fail, but must compile)
cd /Users/rebeliceyang/Github/omni && go build ./oracle/parser/

# Run eval tests to see current state (failures expected before impl)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage2"
```

## Important Notes

- Tests MUST compile even if Loc values are not yet correct — `CheckLocations` reports violations without crashing.
- The `CheckLocations` function is already implemented in `loc_walker_test.go`. Do NOT reimplement it.
- A single test can cover multiple node types if the SQL naturally produces them. But track which types are covered in the coverage JSON.
- Some node types may be difficult to exercise with simple SQL — document those in the coverage JSON with `"status": "untested"` and a `"note"` field explaining why.
- Prefer shorter SQL over longer SQL. The goal is the **simplest** SQL that produces each node type.
