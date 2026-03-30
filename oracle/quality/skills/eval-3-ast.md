# Eval Worker — Stage 3: AST Correctness

You are an Eval Worker in the Oracle Quality Pipeline.
Your role is to write **tests ONLY** — never implementation code.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/strategy.md` — Stage 3 scope (BNF rules, structural assertions)
- `oracle/quality/prevention-rules.md` — **MUST read before starting any work**
- `oracle/parser/bnf/` — 171 BNF rule files defining Oracle SQL grammar
- `oracle/ast/parsenodes.go` — AST node types and their fields
- `oracle/ast/node.go` — `Node` interface, `Loc` struct
- `oracle/parser/parser.go` — `Parse()` entry point, `ParseAndCheck()` if available
- `oracle/parser/oracle_helper_test.go` — `startOracleDB()` for cross-validation tests
- `oracle/quality/coverage/stage2-loc.json` — Stage 2 coverage (context)

## Output Files

- **Test file:** `oracle/parser/eval_ast_test.go`
- **Coverage report:** `oracle/quality/coverage/stage3-ast.json`

## Rules

1. **Tests ONLY** — do NOT modify any non-test `.go` file.
2. Every test function MUST be named `TestEvalStage3_*`.
3. Oracle DB cross-validation tests MUST be named `TestEvalStage3_OracleDB_*`.
4. Tests should fail clearly with descriptive messages.
5. For each BNF rule in `oracle/parser/bnf/`, write at least 1 SQL + structural assertion.
6. Use `Parse()` (or `ParseAndCheck()` if available) then assert node types and field values.
7. Do NOT import packages that do not exist yet.
8. Never modify existing test files.

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[EVAL-STAGE3] STARTED
[EVAL-STAGE3] STEP reading_refs - Reading BNF rules and AST node types
[EVAL-STAGE3] STEP enumerating_bnf - Enumerating BNF rule files (171 rules)
[EVAL-STAGE3] STEP writing_tests_select - Writing structural tests for SELECT rules
[EVAL-STAGE3] STEP writing_tests_dml - Writing structural tests for DML rules
[EVAL-STAGE3] STEP writing_tests_ddl - Writing structural tests for DDL rules
[EVAL-STAGE3] STEP writing_tests_plsql - Writing structural tests for PL/SQL rules
[EVAL-STAGE3] STEP writing_tests_expr - Writing structural tests for expression rules
[EVAL-STAGE3] STEP writing_tests_oracledb - Writing Oracle DB cross-validation tests
[EVAL-STAGE3] STEP build_check - Running go build on test file
[EVAL-STAGE3] STEP coverage_report - Generating stage3-ast.json
[EVAL-STAGE3] DONE
```

If a step fails:
```
[EVAL-STAGE3] FAIL step_name - description
[EVAL-STAGE3] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Test Categories

### A. Structural Assertion Tests (`TestEvalStage3_*`)

For each BNF rule, write a test that:
1. Parses a SQL statement that exercises the rule.
2. Navigates the parsed AST to find the relevant node.
3. Asserts the node type is correct.
4. Asserts key field values match the SQL semantics.

Example:

```go
func TestEvalStage3_SelectStmt_Basic(t *testing.T) {
    result, err := Parse("SELECT a, b FROM t WHERE x = 1")
    if err != nil {
        t.Fatalf("Parse failed: %v", err)
    }
    if len(result.Items) != 1 {
        t.Fatalf("expected 1 statement, got %d", len(result.Items))
    }
    raw, ok := result.Items[0].(*ast.RawStmt)
    if !ok {
        t.Fatalf("expected *RawStmt, got %T", result.Items[0])
    }
    sel, ok := raw.Stmt.(*ast.SelectStmt)
    if !ok {
        t.Fatalf("expected *SelectStmt, got %T", raw.Stmt)
    }
    // Assert target list has 2 entries
    if len(sel.TargetList) != 2 {
        t.Errorf("expected 2 targets, got %d", len(sel.TargetList))
    }
    // Assert FROM clause present
    if len(sel.FromClause) != 1 {
        t.Errorf("expected 1 FROM item, got %d", len(sel.FromClause))
    }
    // Assert WHERE clause present
    if sel.WhereClause == nil {
        t.Errorf("expected WHERE clause, got nil")
    }
}
```

### B. Oracle DB Cross-Validation Tests (`TestEvalStage3_OracleDB_*`)

These tests verify that SQL accepted by a real Oracle DB is also accepted by our parser.

```go
func TestEvalStage3_OracleDB_SelectVariants(t *testing.T) {
    db := startOracleDB(t) // from oracle_helper_test.go
    defer db.Close()

    sqls := []string{
        "SELECT 1 FROM DUAL",
        "SELECT SYSDATE FROM DUAL",
        "SELECT * FROM ALL_TABLES WHERE ROWNUM <= 1",
        // ... more SQL variants
    }

    for _, sql := range sqls {
        t.Run(sql, func(t *testing.T) {
            // Verify Oracle DB accepts it
            _, dbErr := db.Exec(sql)
            if dbErr != nil {
                t.Skipf("Oracle DB rejected SQL (expected): %v", dbErr)
                return
            }

            // If Oracle accepts it, our parser must also accept it
            _, parseErr := Parse(sql)
            if parseErr != nil {
                t.Errorf("Oracle DB accepts %q but parser rejects: %v", sql, parseErr)
            }
        })
    }
}
```

**Important:** Oracle DB tests require a running Oracle container. Use `startOracleDB(t)` which handles container lifecycle and skips if unavailable.

## BNF Rule Mapping

For each BNF file in `oracle/parser/bnf/`, create at least one test:

```
BNF File                              → Test Function
SELECT.bnf                            → TestEvalStage3_Select_*
INSERT.bnf                            → TestEvalStage3_Insert_*
UPDATE.bnf                            → TestEvalStage3_Update_*
DELETE.bnf                            → TestEvalStage3_Delete_*
MERGE.bnf                             → TestEvalStage3_Merge_*
CREATE-TABLE.bnf                      → TestEvalStage3_CreateTable_*
ALTER-TABLE.bnf                       → TestEvalStage3_AlterTable_*
CREATE-INDEX.bnf                      → TestEvalStage3_CreateIndex_*
CREATE-VIEW.bnf                       → TestEvalStage3_CreateView_*
CREATE-FUNCTION.bnf                   → TestEvalStage3_CreateFunction_*
CREATE-PROCEDURE.bnf                  → TestEvalStage3_CreateProcedure_*
CREATE-PACKAGE.bnf                    → TestEvalStage3_CreatePackage_*
CREATE-TRIGGER.bnf                    → TestEvalStage3_CreateTrigger_*
CREATE-SEQUENCE.bnf                   → TestEvalStage3_CreateSequence_*
... (all 171 rules)
```

## Structural Assertions to Make

For each node type, assert at minimum:

1. **Node type**: the parsed node is the expected Go type.
2. **Field presence**: required fields are non-nil/non-empty.
3. **Field values**: identifiers, literals, and keywords match the SQL input.
4. **Child count**: lists have the expected number of elements.
5. **Clause structure**: optional clauses are present/absent as expected.

## Coverage Report Format

After writing tests, generate `oracle/quality/coverage/stage3-ast.json`:

```json
{
  "stage": "3-ast",
  "total_items": 171,
  "tested_items": 0,
  "items": [
    {"bnf_rule": "SELECT", "test": "TestEvalStage3_Select_Basic", "status": "written"},
    {"bnf_rule": "INSERT", "test": "TestEvalStage3_Insert_Basic", "status": "written"},
    {"bnf_rule": "CREATE-TABLE", "test": "TestEvalStage3_CreateTable_Basic", "status": "written"},
    ...
  ],
  "oracle_db_tests": [
    {"test": "TestEvalStage3_OracleDB_SelectVariants", "sql_count": 10, "status": "written"},
    ...
  ]
}
```

The `status` field transitions: `"untested"` → `"written"` → `"passing"` → `"verified"`.

Update `tested_items` to reflect how many BNF rules have at least one test.

## Verification

After writing the test file:

```bash
# Must compile (tests may fail, but must compile)
cd /Users/rebeliceyang/Github/omni && go build ./oracle/parser/

# Run structural assertion tests (failures expected before impl)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage3" -skip "OracleDB"

# Run Oracle DB tests separately (requires container)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage3_OracleDB" -timeout 300s
```

## Important Notes

- Tests MUST compile even if the AST structure is not yet correct.
- Use type assertions with the `, ok` pattern to handle unexpected types gracefully.
- Oracle DB cross-validation tests should be skipped gracefully when no Oracle container is available.
- Some BNF rules may not be implemented in the parser yet — mark those as `"status": "untested"` with a `"note": "parser not yet implemented"` in the coverage JSON.
- Prefer clear, specific assertions over generic checks. `t.Errorf("SelectStmt.TargetList has %d items, want 2", len(sel.TargetList))` is better than `t.Errorf("wrong result")`.
- Group related BNF rules into single test functions where it makes sense (e.g., all ALTER TABLE variants in one function with subtests).
