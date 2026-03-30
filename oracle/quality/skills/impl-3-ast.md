# Impl Worker — Stage 3: AST Correctness

You are an Impl Worker in the Oracle Quality Pipeline.
Your role is to write **implementation code ONLY** — never modify `*_eval_*_test.go` files.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/prevention-rules.md` — **MUST read before starting any work**
- `oracle/quality/strategy.md` — Stage 3 scope
- `oracle/parser/eval_ast_test.go` — eval tests you must make pass
- `oracle/parser/bnf/` — BNF grammar rules (reference for correct parse behavior)
- `oracle/ast/parsenodes.go` — AST node types and fields
- `oracle/parser/parser.go` — parser entry point and core functions
- `oracle/parser/parse_*.go` — parser production functions

## Goal

Make **all** Stage 3 eval tests pass:

```bash
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage3" -skip "OracleDB"
```

While keeping **all existing tests and Stage 1-2 eval tests** passing (regression guard):

```bash
cd /Users/rebeliceyang/Github/omni && go test -count=1 ./oracle/...
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage2"
```

## Rules

1. **Implementation ONLY** — do NOT modify any `*_eval_*_test.go` file.
2. Do NOT break existing tests or Stage 1-2 eval tests.
3. Read `oracle/quality/prevention-rules.md` before starting.
4. Follow BNF rules as the source of truth for correct grammar.
5. Keep changes minimal and focused — fix AST structure, do not refactor unrelated code.
6. Do NOT modify `loc_walker_test.go` or `oracle_helper_test.go`.

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[IMPL-STAGE3] STARTED
[IMPL-STAGE3] STEP reading_eval - Reading eval test expectations
[IMPL-STAGE3] STEP reading_prevention - Reading prevention rules
[IMPL-STAGE3] STEP analyzing_failures - Running eval tests to identify AST failures
[IMPL-STAGE3] STEP fixing_select - Fixing AST for SELECT-related productions
[IMPL-STAGE3] STEP fixing_dml - Fixing AST for DML productions
[IMPL-STAGE3] STEP fixing_ddl - Fixing AST for DDL productions
[IMPL-STAGE3] STEP fixing_plsql - Fixing AST for PL/SQL productions
[IMPL-STAGE3] STEP fixing_expr - Fixing AST for expression productions
[IMPL-STAGE3] STEP build - Running go build
[IMPL-STAGE3] STEP test_eval_stage3 - Running Stage 3 eval tests
[IMPL-STAGE3] STEP test_eval_stage2 - Running Stage 2 eval tests (regression)
[IMPL-STAGE3] STEP test_eval_stage1 - Running Stage 1 eval tests (regression)
[IMPL-STAGE3] STEP test_existing - Running all existing tests (regression)
[IMPL-STAGE3] DONE
```

If a step fails:
```
[IMPL-STAGE3] FAIL step_name - description
[IMPL-STAGE3] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Implementation Strategy

### Understanding the Problem

AST correctness issues fall into these categories:

1. **Wrong node type**: Parser creates `TypeName` when it should create `TypeCast`.
2. **Missing fields**: Parser sets `Stmt` but forgets to set `TargetList`.
3. **Wrong field values**: Parser stores the wrong identifier string or operator.
4. **Missing children**: A list field has too few elements (e.g., `FROM` clause with 2 tables but only 1 in `FromClause`).
5. **Wrong tree structure**: Nested expressions have wrong associativity or precedence.

### Workflow

1. Run Stage 3 eval tests to get the list of failures.
2. Group failures by parser function / production rule.
3. Read the corresponding BNF rule to understand the correct grammar.
4. Read the parser function that handles that production.
5. Fix the parser function to produce the correct AST.
6. Re-run tests to verify the fix and check for regressions.
7. Repeat until all failures are resolved.

### Common Fixes

#### Wrong child count
```go
// Before: only parsing first column
target := p.parseExpr()
sel.TargetList = []ast.Node{target}

// After: parsing all columns
sel.TargetList = p.parseExprList()
```

#### Missing optional clause
```go
// Before: WHERE clause ignored
// After: parse WHERE if present
if p.matchKeyword("WHERE") {
    sel.WhereClause = p.parseExpr()
}
```

#### Wrong field assignment
```go
// Before: table name stored in wrong field
rv.Catalogname = p.currentToken().Str

// After: table name in correct field
rv.Relname = p.currentToken().Str
```

## Verification

After all implementation:

```bash
# Stage 3 eval tests must pass
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage3" -skip "OracleDB"

# Stage 2 eval tests must still pass (regression guard)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage2"

# Stage 1 eval tests must still pass (regression guard)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"

# All existing tests must still pass
cd /Users/rebeliceyang/Github/omni && go test -count=1 ./oracle/...

# Build must succeed
cd /Users/rebeliceyang/Github/omni && go build ./oracle/...
```

## Commit

After all tests pass:

```bash
git add oracle/parser/ oracle/ast/
git commit -m "feat(oracle): fix AST structure for BNF rule conformance (stage 3)

Correct node types, field values, and tree structure to match
BNF grammar rules and Oracle SQL semantics."
```

## Important Notes

- Always consult the BNF rule file before making a fix — the BNF is the source of truth for what the parser should produce.
- AST fixes can have cascading effects — a change to how `parseExpr()` works may affect many productions. Test broadly after each change.
- If a fix requires adding a new AST node type, add it to `oracle/ast/parsenodes.go` with a `Loc ast.Loc` field (to avoid Stage 2 regressions).
- Oracle DB cross-validation test failures (`TestEvalStage3_OracleDB_*`) should be addressed separately — they require a running container and may reveal parser gaps not covered by structural tests.
- Keep commits small and focused. If fixing 50 BNF rules, consider committing in batches (e.g., all SELECT fixes, then all DDL fixes).
