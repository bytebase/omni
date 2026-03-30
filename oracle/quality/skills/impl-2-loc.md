# Impl Worker — Stage 2: Loc Completeness

You are an Impl Worker in the Oracle Quality Pipeline.
Your role is to write **implementation code ONLY** — never modify `*_eval_*_test.go` files.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/prevention-rules.md` — **MUST read before starting any work**
- `oracle/quality/strategy.md` — Stage 2 scope
- `oracle/parser/eval_loc_test.go` — eval tests you must make pass
- `oracle/parser/loc_walker_test.go` — `CheckLocations()` walker (read-only)
- `oracle/ast/parsenodes.go` — AST node types with Loc fields
- `oracle/ast/node.go` — `Loc` struct
- `oracle/parser/parser.go` — parser functions that set Loc values

## Goal

Make **all** Stage 2 eval tests pass:

```bash
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage2"
```

While keeping **all existing tests and Stage 1 eval tests** passing (regression guard):

```bash
cd /Users/rebeliceyang/Github/omni && go test -count=1 ./oracle/...
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"
```

## Rules

1. **Implementation ONLY** — do NOT modify any `*_eval_*_test.go` file.
2. Do NOT break existing tests or Stage 1 eval tests.
3. Read `oracle/quality/prevention-rules.md` before starting.
4. Keep changes minimal and focused — only fix Loc assignments.
5. Do NOT modify `loc_walker_test.go`.

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[IMPL-STAGE2] STARTED
[IMPL-STAGE2] STEP reading_eval - Reading eval test expectations
[IMPL-STAGE2] STEP reading_prevention - Reading prevention rules
[IMPL-STAGE2] STEP analyzing_violations - Running eval tests to identify Loc violations
[IMPL-STAGE2] STEP fixing_select - Fixing Loc for SELECT-related nodes
[IMPL-STAGE2] STEP fixing_dml - Fixing Loc for DML nodes
[IMPL-STAGE2] STEP fixing_ddl - Fixing Loc for DDL nodes
[IMPL-STAGE2] STEP fixing_plsql - Fixing Loc for PL/SQL nodes
[IMPL-STAGE2] STEP fixing_misc - Fixing Loc for remaining nodes
[IMPL-STAGE2] STEP build - Running go build
[IMPL-STAGE2] STEP test_eval_stage2 - Running Stage 2 eval tests
[IMPL-STAGE2] STEP test_eval_stage1 - Running Stage 1 eval tests (regression)
[IMPL-STAGE2] STEP test_existing - Running all existing tests (regression)
[IMPL-STAGE2] DONE
```

If a step fails:
```
[IMPL-STAGE2] FAIL step_name - description
[IMPL-STAGE2] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Implementation Strategy

### Understanding the Problem

Most Loc violations are caused by parser functions that set `Loc.Start` but forget to set `Loc.End`, or set `Loc.End` to 0. The fix is almost always:

```go
// Before (broken):
node.Loc = ast.Loc{Start: startPos}

// After (fixed):
node.Loc = ast.Loc{Start: startPos, End: p.prevEnd()}
```

Where `p.prevEnd()` returns the end position of the most recently consumed token.

### Common Patterns

1. **Missing End**: `Loc.End` is 0 or not set → set it to `p.prevEnd()` at the point where the node production is complete.
2. **Off-by-one**: `Loc.End` is set to `p.pos()` (start of next token) instead of `p.prevEnd()` (end of previous token) → use `p.prevEnd()`.
3. **Compound statements**: `Loc.End` is set at the opening keyword but not updated after the closing keyword → update `Loc.End` after consuming the final token.
4. **Nested nodes**: Parent node's `Loc.End` must encompass all children → set parent's `Loc.End` last.

### Workflow

1. Run Stage 2 eval tests to get the list of violations.
2. For each violation, find the parser function that constructs that node type.
3. Fix the Loc assignment in the parser function.
4. Re-run tests to verify the fix and check for regressions.
5. Repeat until all violations are resolved.

## Verification

After all implementation:

```bash
# Stage 2 eval tests must pass
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
git commit -m "feat(oracle): fix Loc.End for AST node types (stage 2 loc completeness)

Set Loc.End = p.prevEnd() for all Loc-bearing AST nodes so that
sql[node.Loc.Start:node.Loc.End] produces valid source substrings."
```

## Important Notes

- The most common fix is `node.Loc.End = p.prevEnd()` — look for this pattern.
- Some nodes span complex productions (e.g., `CREATE TABLE ... (columns) ...`). Make sure `Loc.End` is set after the **last** token of the production, not an intermediate one.
- Test after each batch of fixes to catch regressions early.
- If a fix for one node type breaks another, it usually means the parser is sharing or reusing position variables. Trace carefully.
