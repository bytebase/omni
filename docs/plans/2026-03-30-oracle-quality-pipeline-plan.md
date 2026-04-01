# Oracle Quality Pipeline Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build the infrastructure for a 6-stage quality pipeline that evaluates and implements Oracle engine feature surfaces with strict external evaluation, role-separated workers, and a self-evolving feedback loop.

**Architecture:** A quality driver script orchestrates three worker roles (Eval, Impl, Insight) per stage. Each stage's eval tests become permanent regression guards. Coverage is tracked in JSON, insights accumulate in markdown, and the coverage strategy itself is iterable.

**Tech Stack:** Bash (driver), Go (tests/implementation), testcontainers (Oracle DB), Claude CLI (worker dispatch)

**Design doc:** `docs/plans/2026-03-30-oracle-quality-pipeline-design.md`

---

## Task 1: Create oracle/quality/ directory structure and strategy.md

**Files:**
- Create: `oracle/quality/strategy.md`
- Create: `oracle/quality/prevention-rules.md`
- Create: `oracle/quality/insights/patterns.json`
- Create: `oracle/quality/corpus/.gitkeep`

**Step 1: Create directory structure**

```bash
mkdir -p oracle/quality/coverage
mkdir -p oracle/quality/insights
mkdir -p oracle/quality/corpus
```

**Step 2: Write strategy.md**

This is the coverage strategy matrix that Eval Workers read and Insight Workers update. It defines what each stage must cover and how.

```markdown
# Oracle Quality Pipeline — Coverage Strategy

This file is read by Eval Workers to determine what tests to generate,
and updated by Insight Workers when blind spots are discovered.

## Stage 1: Foundation

### Scope
- Loc sentinel: all unknown positions use `Loc{Start: -1, End: -1}` (not End=0)
- Token.End: lexer tracks end byte offset for every token
- ParseError: has Severity, Code, Message, Position fields
- Parser.source: stores input SQL for tokenText extraction
- RawStmt: uses `Loc` instead of `StmtLocation`/`StmtLen`
- Utility functions: `NoLoc()`, `NodeLoc()`, `ListSpan()`

### Evaluation Strategy
- **PG behavior reference (C)**: each feature must match PG's equivalent behavior
- Eval tests verify structural properties (field existence, sentinel values, function signatures)

### Enumerable Set
- 7 infrastructure items listed above — all must have eval tests

### Known Blind Spots
(initially empty — updated by Insight Worker)

---

## Stage 2: Loc Completeness

### Scope
- Every Loc-bearing AST node (248 types) must have valid Start/End after parsing
- `sql[node.Loc.Start:node.Loc.End]` must produce a reasonable substring

### Evaluation Strategy
- **Reflection walker (automated)**: walk parsed AST, check every Loc field
- For each node type, at least one SQL statement that exercises it

### Enumerable Set
- 248 Loc-bearing struct types in `oracle/ast/parsenodes.go`
- Coverage: each type must appear in at least one eval test SQL

### Generation Strategy
- For each AST node type, find the simplest SQL that produces that node
- Parse → walk → verify Loc.Start >= 0 && Loc.End > Loc.Start

### Known Blind Spots
(initially empty)

---

## Stage 3: AST Correctness

### Scope
- Parser produces correct AST structure for all supported SQL
- Field values match SQL semantics (correct identifiers, operators, clause types)

### Evaluation Strategy
- **Oracle DB cross-validation (B)**: SQL executed on real Oracle DB; if DB accepts, parser must parse successfully
- **Structural assertions (A)**: parse result fields checked against expected values

### Enumerable Set
- BNF production rules (230+ files in oracle/parser/bnf/)
- Each rule: at least 1 test case with structural assertion

### Generation Strategy (non-enumerable)
- Oracle documentation example SQL → validated against Oracle DB
- Real-world Oracle SQL corpus → validated against Oracle DB

### Known Blind Spots
(initially empty)

---

## Stage 4: Error Quality

### Scope
- Invalid SQL detected (not silently accepted)
- Parser does not panic on any input
- Error position points to the correct token
- Error message is descriptive

### Evaluation Strategy
- **Mutation generation**: from valid SQL, systematically produce invalid variants
  - Truncation: cut at each token boundary
  - Deletion: remove one required constituent
  - Replacement: swap keyword with invalid token
  - Duplication: repeat a clause
- **Oracle DB cross-validation (B)**: if Oracle DB rejects SQL, parser should also reject

### Enumerable Set
- None (error space is infinite)

### Generation Strategy
- Take all Layer 1/2 valid SQL from Stage 3
- Apply 4 mutation types to each
- Verify: no panic, returns error, position >= 0

### Known Blind Spots
(initially empty)

---

## Stage 5: Completion

### Scope
- Parser suggests valid syntax continuations at cursor position
- Candidates include relevant keywords, identifiers, clause types

### Evaluation Strategy
- **PG behavior reference (C)**: completion infrastructure mirrors PG's pattern
- **Candidate set assertions (A)**: at cursor position X, candidates must include/exclude specific items

### Enumerable Set
- Statement types (30+) x key cursor positions (after keyword, after table, in WHERE, etc.)

### Generation Strategy (non-enumerable)
- Nested queries, CTE context, PL/SQL block context

### Known Blind Spots
(initially empty)

---

## Stage 6: Catalog / Migration

### Scope
- In-memory schema model from parsed DDL
- Migration DDL generation between schema states
- Round-trip: apply migration to source schema → matches target schema

### Evaluation Strategy
- **Oracle DB oracle verification (B)**: DDL round-trip on real Oracle container
- **Pairwise combinations**: object types x operations x properties

### Enumerable Set
- Object types: table, view, index, trigger, sequence, function, procedure, package, type, synonym
- Operations: create, alter, drop, rename, comment

### Generation Strategy (non-enumerable)
- Pairwise: cover all 2-way dimension combinations
- Property combinations generated systematically, not exhaustively

### Known Blind Spots
(initially empty)
```

**Step 3: Write initial prevention-rules.md**

```markdown
# Prevention Rules

Accumulated rules from Insight Worker analysis across all stages.
Impl Workers MUST read this file before starting work.

## Rules

(No rules yet — this file is populated by Insight Workers as stages complete.)
```

**Step 4: Write initial patterns.json**

```json
{
  "version": 1,
  "patterns": []
}
```

**Step 5: Create corpus .gitkeep**

```bash
touch oracle/quality/corpus/.gitkeep
```

**Step 6: Commit**

```bash
git add oracle/quality/
git commit -m "feat(oracle): create quality pipeline directory structure and coverage strategy"
```

---

## Task 2: Create the quality driver script

**Files:**
- Create: `scripts/quality-driver.sh`

**Step 1: Write quality-driver.sh**

This driver orchestrates the eval → impl → insight flow for each stage.
It extends the existing `driver.sh` pattern but adds:
- Three worker phases per stage (eval, impl, insight)
- Regression guard (run prior stage tests before marking done)
- Coverage gap checking
- Insight feedback loop (re-enter impl if insight generates new tests)

```bash
#!/bin/bash
# Oracle Quality Pipeline Driver
#
# Usage: ./quality-driver.sh <stage> [--eval-only] [--impl-only] [--insight-only]
#
# Stages: 1-foundation, 2-loc, 3-ast, 4-error, 5-completion, 6-catalog
#
# Each stage runs three phases:
#   1. Eval: generate eval tests (expected to fail)
#   2. Impl: implement code to pass eval tests
#   3. Insight: extract patterns, generate adversarial tests
#
# If insight generates new tests, the driver re-enters impl phase.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OMNI_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
QUALITY_DIR="$OMNI_DIR/oracle/quality"
LOG_DIR="$QUALITY_DIR/logs"

# Parse arguments
STAGE=""
EVAL_ONLY=false
IMPL_ONLY=false
INSIGHT_ONLY=false
MAX_IMPL_ITERATIONS=50
MAX_INSIGHT_ROUNDS=3

while [[ $# -gt 0 ]]; do
    case "$1" in
        --eval-only) EVAL_ONLY=true; shift ;;
        --impl-only) IMPL_ONLY=true; shift ;;
        --insight-only) INSIGHT_ONLY=true; shift ;;
        --max-impl-iterations) MAX_IMPL_ITERATIONS="$2"; shift 2 ;;
        --max-insight-rounds) MAX_INSIGHT_ROUNDS="$2"; shift 2 ;;
        *) STAGE="$1"; shift ;;
    esac
done

if [ -z "$STAGE" ]; then
    echo "Usage: $0 <stage> [--eval-only] [--impl-only] [--insight-only]"
    echo "Stages: 1-foundation, 2-loc, 3-ast, 4-error, 5-completion, 6-catalog"
    exit 1
fi

# Map stage to number for ordering
stage_number() {
    case "$1" in
        1-foundation) echo 1 ;;
        2-loc) echo 2 ;;
        3-ast) echo 3 ;;
        4-error) echo 4 ;;
        5-completion) echo 5 ;;
        6-catalog) echo 6 ;;
        *) echo "Unknown stage: $1" >&2; exit 1 ;;
    esac
}

STAGE_NUM=$(stage_number "$STAGE")
SKILL_DIR="$QUALITY_DIR/skills"
STAGE_COVERAGE="$QUALITY_DIR/coverage/stage${STAGE_NUM}-${STAGE#*-}.json"

mkdir -p "$LOG_DIR" "$QUALITY_DIR/coverage"

echo "=== Oracle Quality Pipeline: Stage $STAGE ==="
echo "Working directory: $OMNI_DIR"
echo ""

# --- Regression guard: verify all prior stages pass ---
run_regression_guard() {
    if [ "$STAGE_NUM" -le 1 ]; then
        return 0
    fi
    echo "  [Regression Guard] Running eval tests for stages 1..$(($STAGE_NUM - 1))..."
    local test_pattern="TestEvalStage[1-$(($STAGE_NUM - 1))]"
    cd "$OMNI_DIR" && go test ./oracle/... -run "$test_pattern" -count=1 2>&1 | tail -5
    if [ ${PIPESTATUS[0]} -ne 0 ]; then
        echo "  REGRESSION DETECTED: prior stage eval tests failing. Fix before proceeding."
        exit 1
    fi
    echo "  [Regression Guard] All prior stages pass."
}

# --- Phase 1: Eval ---
run_eval() {
    echo ""
    echo "========================================="
    echo "  Phase: EVAL (Stage $STAGE)"
    echo "========================================="

    local timestamp=$(date '+%Y%m%d_%H%M%S')
    local log_file="$LOG_DIR/eval_${STAGE}_${timestamp}.log"
    local skill_file="$SKILL_DIR/eval-${STAGE}.md"

    if [ ! -f "$skill_file" ]; then
        echo "Error: $skill_file not found"
        exit 1
    fi

    cd "$OMNI_DIR" && claude -p "$(cat "$skill_file")" \
        --dangerously-skip-permissions \
        2>&1 | tee "$log_file"

    # Verify eval tests were created and fail as expected
    echo ""
    echo "  Verifying eval tests exist and fail..."
    local test_pattern="TestEvalStage${STAGE_NUM}"
    cd "$OMNI_DIR" && go test ./oracle/... -run "$test_pattern" -count=1 2>&1 | tail -5 || true
    echo "  Eval phase complete."
}

# --- Phase 2: Impl ---
run_impl() {
    echo ""
    echo "========================================="
    echo "  Phase: IMPL (Stage $STAGE)"
    echo "========================================="

    local skill_file="$SKILL_DIR/impl-${STAGE}.md"

    if [ ! -f "$skill_file" ]; then
        echo "Error: $skill_file not found"
        exit 1
    fi

    local iteration=0
    while [ $iteration -lt $MAX_IMPL_ITERATIONS ]; do
        iteration=$((iteration + 1))
        local timestamp=$(date '+%Y%m%d_%H%M%S')
        local log_file="$LOG_DIR/impl_${STAGE}_${iteration}_${timestamp}.log"

        # Check if all eval tests pass
        local test_pattern="TestEvalStage${STAGE_NUM}"
        local test_exit=0
        cd "$OMNI_DIR" && go test ./oracle/... -run "$test_pattern" -count=1 > /dev/null 2>&1 || test_exit=$?

        if [ "$test_exit" -eq 0 ]; then
            echo "  [Impl $iteration] All eval tests pass."
            # Also run regression guard
            run_regression_guard
            echo "  Impl phase complete."
            return 0
        fi

        echo "  [Impl $iteration/$MAX_IMPL_ITERATIONS] Eval tests still failing, dispatching worker..."

        cd "$OMNI_DIR" && claude -p "$(cat "$skill_file")" \
            --dangerously-skip-permissions \
            2>&1 | tee "$log_file"

        sleep 2
    done

    echo "  Reached max impl iterations."
    return 1
}

# --- Phase 3: Insight ---
run_insight() {
    echo ""
    echo "========================================="
    echo "  Phase: INSIGHT (Stage $STAGE)"
    echo "========================================="

    local timestamp=$(date '+%Y%m%d_%H%M%S')
    local log_file="$LOG_DIR/insight_${STAGE}_${timestamp}.log"
    local skill_file="$SKILL_DIR/insight-${STAGE}.md"

    if [ ! -f "$skill_file" ]; then
        echo "Error: $skill_file not found"
        exit 1
    fi

    cd "$OMNI_DIR" && claude -p "$(cat "$skill_file")" \
        --dangerously-skip-permissions \
        2>&1 | tee "$log_file"

    # Check if insight generated new test cases by re-running eval tests
    local test_pattern="TestEvalStage${STAGE_NUM}"
    local test_exit=0
    cd "$OMNI_DIR" && go test ./oracle/... -run "$test_pattern" -count=1 > /dev/null 2>&1 || test_exit=$?

    if [ "$test_exit" -ne 0 ]; then
        echo "  Insight generated new failing tests. Re-entering impl phase."
        return 1  # signal to re-enter impl
    fi

    echo "  Insight phase complete. No new failures."
    return 0
}

# --- Main flow ---

run_regression_guard

if [ "$EVAL_ONLY" = true ]; then
    run_eval
    exit 0
fi

if [ "$IMPL_ONLY" = true ]; then
    run_impl
    exit $?
fi

if [ "$INSIGHT_ONLY" = true ]; then
    run_insight
    exit $?
fi

# Full pipeline: eval → impl → insight (with re-entry loop)
run_eval
run_impl || { echo "Impl phase did not converge."; exit 1; }

insight_round=0
while [ $insight_round -lt $MAX_INSIGHT_ROUNDS ]; do
    insight_round=$((insight_round + 1))

    insight_exit=0
    run_insight || insight_exit=$?

    if [ "$insight_exit" -eq 0 ]; then
        break  # no new failures
    fi

    echo ""
    echo "  [Insight round $insight_round/$MAX_INSIGHT_ROUNDS] Re-entering impl..."
    run_impl || { echo "Impl phase did not converge after insight round $insight_round."; exit 1; }
done

# Final verification
echo ""
echo "========================================="
echo "  VERIFICATION (Stage $STAGE)"
echo "========================================="

cd "$OMNI_DIR" && go build ./oracle/... 2>&1
cd "$OMNI_DIR" && go test ./oracle/... -run "TestEval" -count=1 -v 2>&1 | tail -20

echo ""
echo "========================================="
echo "  STAGE COMPLETE: $STAGE"
echo "========================================="
```

**Step 2: Make executable**

```bash
chmod +x scripts/quality-driver.sh
```

**Step 3: Commit**

```bash
git add scripts/quality-driver.sh
git commit -m "feat(oracle): add quality pipeline driver script with eval/impl/insight phases"
```

---

## Task 3: Create worker skill templates

**Files:**
- Create: `oracle/quality/skills/eval-1-foundation.md`
- Create: `oracle/quality/skills/impl-1-foundation.md`
- Create: `oracle/quality/skills/insight-1-foundation.md`

These are the concrete skill files for Stage 1. They serve as both functional skills AND templates for subsequent stages.

**Step 1: Write eval-1-foundation.md**

```markdown
# Eval Worker — Stage 1: Foundation

You are writing **evaluation tests** for the Oracle parser's foundation infrastructure.
Your goal is to write strict, comprehensive tests that verify the parser has the correct
infrastructure in place. You do NOT write implementation code.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Your Role

- You write tests. That's it.
- Your tests should be maximally strict — test every edge case you can think of.
- You CANNOT write or modify any non-test files.
- All test functions MUST be named `TestEvalStage1_*` so the driver can filter them.
- Output test files to: `oracle/parser/eval_foundation_test.go`

## What You're Testing

Read `oracle/quality/strategy.md` Stage 1 section for the full scope. Summary:

### 1. Loc Sentinel Value
- `ast.NoLoc()` must return `Loc{Start: -1, End: -1}`
- All "unknown" Loc fields must use -1, never 0

### 2. Token.End Field
- Lexer must track end byte offset for every token
- `token.End` must equal `token.Loc.Start + len(token.Str)` for identifiers/keywords
- `token.End` must be > `token.Loc.Start` for all tokens

### 3. ParseError Enhancement
- ParseError must have fields: Severity, Code, Message, Position
- Default Severity: "ERROR"
- Default Code: "42601" (syntax error)
- `Error()` method must format as: `"ERROR: <message> (SQLSTATE <code>)"`

### 4. Parser.source Field
- Parser must store the input SQL string
- Needed for tokenText extraction in error messages

### 5. RawStmt Loc
- RawStmt must use `Loc` field (not StmtLocation/StmtLen)
- `Loc.Start` = byte offset of statement start
- `Loc.End` = byte offset past statement end (exclusive)

### 6. Utility Functions
- `NodeLoc(node Node) Loc` — extract Loc from any AST node
- `ListSpan(list) Loc` — compute byte range covering all nodes in a list

## Reference

Read the PG equivalents to understand expected behavior:
- `pg/ast/node.go` — PG's Loc struct and Node interface
- `pg/parser/parser.go` — PG's ParseError struct
- `pg/parser/loc_test.go` — PG's Loc validation patterns

## Output

1. Write `oracle/parser/eval_foundation_test.go` with all tests
2. Write `oracle/quality/coverage/stage1-foundation.json`:
   ```json
   {
     "stage": 1,
     "surface": "foundation",
     "items": [
       {"id": "no_loc", "description": "NoLoc() returns {-1,-1}", "tested": true},
       {"id": "token_end", "description": "Token.End field exists and is correct", "tested": true},
       ...
     ],
     "total": 7,
     "tested": 7,
     "gaps": []
   }
   ```
3. Run `go build ./oracle/...` — if it doesn't compile, that's expected (features not implemented yet). Write tests that check for the existence of fields/functions using build tags or reflection.
4. Commit: `test(oracle): eval stage 1 — foundation infrastructure tests`

## Important

- Do NOT look at implementation code to decide what to test. Test what SHOULD exist per the strategy.
- If something doesn't compile because the feature doesn't exist yet, use a build-tag-separated test file or write tests that verify via reflection.
- Prefer tests that will fail clearly (not just "doesn't compile") so the Impl Worker gets actionable feedback.
```

**Step 2: Write impl-1-foundation.md**

```markdown
# Impl Worker — Stage 1: Foundation

You are implementing Oracle parser foundation infrastructure to make eval tests pass.
You do NOT modify test files.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Your Role

- You write implementation code. That's it.
- You CANNOT modify any `*_eval_*_test.go` files.
- Your goal: make `go test ./oracle/... -run TestEvalStage1 -count=1` pass.

## What to Implement

Run the eval tests first to see what's failing:
```bash
go test ./oracle/... -run TestEvalStage1 -count=1 -v 2>&1 | head -100
```

Then fix failures one by one. The expected changes are:

### 1. Loc Sentinel (`oracle/ast/node.go`)
- Add `NoLoc() Loc` function returning `Loc{Start: -1, End: -1}`
- Audit existing code: replace `Loc{}` or `Loc{0, 0}` with `NoLoc()` where meaning is "unknown"

### 2. Token.End (`oracle/parser/lexer.go`)
- Add `End int` field to Token struct
- Set `End` in the lexer after each token is scanned

### 3. ParseError Enhancement (`oracle/parser/parser.go`)
- Add `Severity string` and `Code string` fields to ParseError
- Default Severity to "ERROR", Code to "42601"
- Update `Error()` method to format as PG does

### 4. Parser.source (`oracle/parser/parser.go`)
- Store input SQL in parser struct
- Add helper methods: `tokenText(tok Token) string`

### 5. RawStmt Loc (`oracle/ast/parsenodes.go`)
- Change RawStmt from `StmtLocation int` / `StmtLen int` to `Loc Loc`
- Update parser to set `Loc` on every RawStmt
- Fix all code that reads StmtLocation/StmtLen

### 6. Utility Functions (`oracle/ast/node.go`)
- Add `NodeLoc(node Node) Loc`
- Add `ListSpan(items []Node) Loc`

## Prevention Rules

Read `oracle/quality/prevention-rules.md` before starting.

## Process

1. Read all failing eval tests to understand what's expected
2. Implement changes in small batches
3. After each change: `go build ./oracle/... && go test ./oracle/... -run TestEvalStage1 -count=1`
4. Also run existing parser tests: `go test ./oracle/parser/... -count=1` — do NOT break them
5. Commit after each logical group of changes

## Progress Logging

Print to stdout:
```
[STAGE1] STARTED
[STAGE1] STEP <name> — description
[STAGE1] DONE
```
```

**Step 3: Write insight-1-foundation.md**

```markdown
# Insight Worker — Stage 1: Foundation

You analyze the implementation work that just completed for Stage 1 and extract
lessons that will prevent similar issues in future stages.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Your Role

- You do NOT write implementation code or modify eval tests directly.
- You analyze git diffs, identify patterns, and produce structured outputs.
- You may add NEW adversarial test cases to eval test files (append only, never modify existing tests).

## Process

### Step 1: Collect Fix Commits

Find all commits related to Stage 1 implementation:
```bash
git log --oneline --since="$(git log -1 --format=%ci -- oracle/quality/skills/impl-1-foundation.md)" -- oracle/
```

### Step 2: Analyze Each Fix

For each commit that contains a "fix" or required multiple attempts:

1. **What was wrong?** — describe the bug
2. **Why did the Impl Worker get it wrong initially?** — root cause
3. **Is this a pattern?** — could the same mistake happen in Stages 2-6?

### Step 3: Write Patterns

For each identified pattern, create `oracle/quality/insights/<pattern-name>.md`:

```markdown
---
pattern: <pattern-name>
severity: high|medium|low
discovered_in: stage1
occurrences: N
---

## Root Cause
Why this happened.

## Affected Commits
- <hash> — description

## Where Else (scan results)
- [ ] Stage 2: ...
- [ ] Stage 3: ...

## Prevention Rule
Rule text for future workers.
```

Update `oracle/quality/insights/patterns.json` with the new pattern.

### Step 4: Update Prevention Rules

Append new rules to `oracle/quality/prevention-rules.md`.

### Step 5: Check Strategy Blind Spots

Did any fix reveal a type of problem not covered by `oracle/quality/strategy.md`?
If yes, add it to the "Known Blind Spots" section of the relevant stage.

### Step 6: Generate Adversarial Tests (if applicable)

If patterns suggest edge cases not covered by existing eval tests,
append new `TestEvalStage1_Adversarial_*` functions to the eval test file.

### Step 7: Commit

```bash
git add oracle/quality/insights/ oracle/quality/prevention-rules.md oracle/quality/strategy.md
git commit -m "insight(oracle): stage 1 pattern analysis and prevention rules"
```

## Output

Print a summary to stdout:
```
[INSIGHT] Stage 1 Analysis
[INSIGHT] Commits analyzed: N
[INSIGHT] Patterns found: N
[INSIGHT] New prevention rules: N
[INSIGHT] Strategy updates: N
[INSIGHT] Adversarial tests added: N
[INSIGHT] DONE
```
```

**Step 4: Commit**

```bash
git add oracle/quality/skills/
git commit -m "feat(oracle): add eval/impl/insight worker skill files for stage 1"
```

---

## Task 4: Create Oracle DB testcontainer helper

**Files:**
- Create: `oracle/parser/oracle_helper_test.go`

This provides the Oracle DB connection infrastructure for cross-validation in Stages 3-6.
It follows the same pattern as `pg/catalog/oracle_helper_test.go` and `mysql/catalog/oracle_test.go`.

**Step 1: Write oracle_helper_test.go**

```go
package parser

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	_ "github.com/sijms/go-ora/v2"
	"github.com/testcontainers/testcontainers-go"
	tcoracledb "github.com/testcontainers/testcontainers-go/modules/oracledb"
)

// oracleDB wraps a real Oracle DB container connection for cross-validation testing.
type oracleDB struct {
	db  *sql.DB
	ctx context.Context
}

var (
	oracleDBOnce    sync.Once
	oracleDBInst    *oracleDB
	oracleDBCleanup func()
)

// startOracleDB starts a shared Oracle Free container. The container is reused
// across all tests via sync.Once.
func startOracleDB(t *testing.T) *oracleDB {
	t.Helper()
	oracleDBOnce.Do(func() {
		ctx := context.Background()
		container, err := tcoracledb.Run(ctx, "gvenzl/oracle-free:23-slim-faststart",
			tcoracledb.WithInitScripts(), // no init scripts
			testcontainers.WithEnv(map[string]string{
				"ORACLE_PASSWORD": "test",
			}),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to start Oracle container: %v", err))
		}

		connStr, err := container.ConnectionString(ctx)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			panic(fmt.Sprintf("failed to get connection string: %v", err))
		}

		db, err := sql.Open("oracle", connStr)
		if err != nil {
			_ = testcontainers.TerminateContainer(container)
			panic(fmt.Sprintf("failed to open database: %v", err))
		}

		if err := db.PingContext(ctx); err != nil {
			_ = testcontainers.TerminateContainer(container)
			panic(fmt.Sprintf("failed to ping Oracle: %v", err))
		}

		oracleDBInst = &oracleDB{db: db, ctx: ctx}
		oracleDBCleanup = func() {
			db.Close()
			_ = testcontainers.TerminateContainer(container)
		}
	})
	if oracleDBInst == nil {
		t.Fatal("Oracle DB container failed to start")
	}
	t.Cleanup(func() {
		// Cleanup happens when all tests finish, not per-test
	})
	return oracleDBInst
}

// canParseOnOracle checks if Oracle DB accepts the given SQL by preparing it.
// Returns nil if Oracle accepts, error if Oracle rejects.
func (o *oracleDB) canParseOnOracle(sql string) error {
	// Use EXPLAIN PLAN to validate SQL without executing it
	_, err := o.db.ExecContext(o.ctx, fmt.Sprintf("EXPLAIN PLAN FOR %s", sql))
	return err
}

// execSQL executes SQL on the Oracle container.
func (o *oracleDB) execSQL(t *testing.T, sql string) {
	t.Helper()
	_, err := o.db.ExecContext(o.ctx, sql)
	if err != nil {
		t.Fatalf("execSQL(%q): %v", sql, err)
	}
}
```

**Step 2: Add go-ora dependency**

```bash
cd /Users/rebeliceyang/Github/omni && go get github.com/sijms/go-ora/v2
```

Note: The exact Oracle testcontainer module and driver may need adjustment based on what's available. The `gvenzl/oracle-free:23-slim-faststart` image is ~1.5GB but starts in ~10 seconds. The `github.com/testcontainers/testcontainers-go/modules/oracledb` module may need to be verified — if unavailable, use a generic container setup. The key pattern is the same as the PG helper.

**Step 3: Verify it compiles**

```bash
go build ./oracle/...
```

**Step 4: Commit**

```bash
git add oracle/parser/oracle_helper_test.go go.mod go.sum
git commit -m "feat(oracle): add Oracle DB testcontainer helper for cross-validation"
```

---

## Task 5: Create reflection-based Loc walker for Oracle

**Files:**
- Create: `oracle/parser/loc_walker_test.go`

This is needed for Stage 2 (Loc completeness) evaluation. It mirrors PG's `CheckLocations` in `pg/parser/loc_test.go`.

**Step 1: Write loc_walker_test.go**

```go
package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// LocViolation records a single location validation failure.
type LocViolation struct {
	Path    string // e.g. "Items[0](SelectStmt).TargetList.Items[0](ResTarget)"
	NodeTag string // node type name
	Start   int
	End     int
	Reason  string
}

func (v LocViolation) String() string {
	return fmt.Sprintf("%s [%s]: Start=%d End=%d — %s", v.Path, v.NodeTag, v.Start, v.End, v.Reason)
}

// CheckLocations parses sql via Parse(), recursively walks the AST using
// reflection, and returns all Loc violations where Start >= 0 but End <= Start.
func CheckLocations(t *testing.T, sql string) []LocViolation {
	t.Helper()

	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("CheckLocations Parse(%q): %v", sql, err)
	}

	var violations []LocViolation
	if result != nil {
		for i, item := range result.Items {
			path := fmt.Sprintf("Items[%d]", i)
			walkNodeLocs(reflect.ValueOf(item), path, &violations)
		}
	}
	return violations
}

// walkNodeLocs recursively walks a reflected AST value, checking every Loc field.
func walkNodeLocs(v reflect.Value, path string, violations *[]LocViolation) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		typeName := v.Type().Name()

		locField := v.FieldByName("Loc")
		if locField.IsValid() && locField.Type() == reflect.TypeOf(ast.Loc{}) {
			loc := locField.Interface().(ast.Loc)
			if loc.Start >= 0 && loc.End <= loc.Start {
				*violations = append(*violations, LocViolation{
					Path:    path,
					NodeTag: typeName,
					Start:   loc.Start,
					End:     loc.End,
					Reason:  "Start >= 0 but End <= Start",
				})
			}
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			if field.Name == "Loc" {
				continue
			}
			childPath := path
			if typeName != "" {
				childPath = fmt.Sprintf("%s(%s).%s", path, typeName, field.Name)
			} else {
				childPath = fmt.Sprintf("%s.%s", path, field.Name)
			}
			walkNodeLocs(v.Field(i), childPath, violations)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			walkNodeLocs(v.Index(i), elemPath, violations)
		}
	}
}

// ListLocBearingTypes returns all struct type names in the AST that have a Loc field.
// Used by eval tests to verify coverage.
func ListLocBearingTypes() []string {
	// This uses the oracle/ast package's outfuncs.go registrations or
	// direct reflection on known node types.
	// For now, we scan parsenodes.go types via a known list.
	// The eval worker will generate the actual comprehensive list.
	return nil // placeholder — eval worker populates via reflection
}
```

**Step 2: Verify it compiles**

```bash
go build ./oracle/...
go test ./oracle/parser/ -run "^$" -count=1  # compile check only
```

**Step 3: Commit**

```bash
git add oracle/parser/loc_walker_test.go
git commit -m "feat(oracle): add reflection-based Loc walker for AST validation"
```

---

## Task 6: Create coverage.json schema and initial stage files

**Files:**
- Create: `oracle/quality/coverage/stage1-foundation.json`

**Step 1: Write initial coverage file for Stage 1**

This file is updated by the Eval Worker after generating tests, and checked by the Driver for gaps.

```json
{
  "stage": 1,
  "surface": "foundation",
  "status": "pending",
  "items": [
    {"id": "no_loc", "description": "NoLoc() returns Loc{-1,-1}", "tested": false},
    {"id": "loc_sentinel_consistency", "description": "All unknown Loc fields use -1 not 0", "tested": false},
    {"id": "token_end", "description": "Token struct has End field, lexer sets it", "tested": false},
    {"id": "parse_error_severity", "description": "ParseError has Severity field", "tested": false},
    {"id": "parse_error_code", "description": "ParseError has Code field", "tested": false},
    {"id": "parse_error_format", "description": "ParseError.Error() matches PG format", "tested": false},
    {"id": "parser_source", "description": "Parser stores source SQL string", "tested": false},
    {"id": "raw_stmt_loc", "description": "RawStmt uses Loc instead of StmtLocation/StmtLen", "tested": false},
    {"id": "node_loc", "description": "NodeLoc() extracts Loc from any AST node", "tested": false},
    {"id": "list_span", "description": "ListSpan() computes byte range of node list", "tested": false}
  ],
  "total": 10,
  "tested": 0,
  "gaps": ["no_loc", "loc_sentinel_consistency", "token_end", "parse_error_severity", "parse_error_code", "parse_error_format", "parser_source", "raw_stmt_loc", "node_loc", "list_span"]
}
```

**Step 2: Commit**

```bash
git add oracle/quality/coverage/
git commit -m "feat(oracle): add initial coverage tracking for stage 1 foundation"
```

---

## Task 7: Wire up the driver — verify end-to-end flow

**Files:**
- Modify: `scripts/quality-driver.sh` (fix paths if needed)

**Step 1: Verify directory structure is complete**

```bash
ls -la oracle/quality/
ls -la oracle/quality/skills/
ls -la oracle/quality/coverage/
ls -la oracle/quality/insights/
```

Expected:
```
oracle/quality/
  strategy.md
  prevention-rules.md
  skills/
    eval-1-foundation.md
    impl-1-foundation.md
    insight-1-foundation.md
  coverage/
    stage1-foundation.json
  insights/
    patterns.json
  corpus/
    .gitkeep
```

**Step 2: Dry run the driver (eval-only)**

```bash
./scripts/quality-driver.sh 1-foundation --eval-only
```

This should:
1. Pass regression guard (no prior stages)
2. Invoke `claude -p` with the eval skill
3. Eval Worker creates `oracle/parser/eval_foundation_test.go`
4. Driver verifies the tests exist

**Step 3: Check that eval tests fail as expected**

```bash
go test ./oracle/... -run TestEvalStage1 -count=1 -v 2>&1 | head -30
```

Expected: compilation errors or test failures (features not implemented yet).

**Step 4: Dry run impl phase**

```bash
./scripts/quality-driver.sh 1-foundation --impl-only
```

This should dispatch the Impl Worker to make eval tests pass.

**Step 5: Dry run insight phase**

```bash
./scripts/quality-driver.sh 1-foundation --insight-only
```

**Step 6: If everything works, run full pipeline**

```bash
./scripts/quality-driver.sh 1-foundation
```

---

## Task 8: Create skill templates for remaining stages (2-6)

**Files:**
- Create: `oracle/quality/skills/eval-2-loc.md`
- Create: `oracle/quality/skills/impl-2-loc.md`
- Create: `oracle/quality/skills/insight-2-loc.md`
- (repeat for stages 3-6)

Each stage's skills follow the same structure as Stage 1 but with stage-specific content.

**Step 1: Create Stage 2 (Loc) skills**

The eval skill for Stage 2 is unique because it uses the reflection walker:

- Eval Worker: for each of the 248 Loc-bearing node types, find/write a SQL that produces that node, then verify `CheckLocations` reports no violations.
- Impl Worker: fix Loc.End assignments across parser files.
- Insight Worker: same pattern analysis as Stage 1.

**Step 2: Create Stage 3 (AST) skills**

- Eval Worker: for each BNF rule, write a SQL + structural assertion. Cross-validate SQL against Oracle DB.
- Impl Worker: fix parser to produce correct AST.
- Insight Worker: pattern analysis.

**Step 3: Create Stage 4 (Error) skills**

- Eval Worker: take valid SQL from Stage 3, apply mutations, verify parser doesn't panic and returns error.
- Impl Worker: fix error handling.
- Insight Worker: pattern analysis.

**Step 4: Create Stage 5-6 skills**

These depend on Stages 1-4 being complete. Write templates now, fill in details later based on accumulated prevention rules and insights.

**Step 5: Commit**

```bash
git add oracle/quality/skills/
git commit -m "feat(oracle): add worker skill files for stages 2-6"
```

---

## Summary of Deliverables

| Task | Deliverable | Purpose |
|------|------------|---------|
| 1 | `oracle/quality/` directory + strategy.md | Coverage strategy and file structure |
| 2 | `scripts/quality-driver.sh` | Pipeline orchestration |
| 3 | Stage 1 skill files (eval/impl/insight) | First concrete worker instructions |
| 4 | `oracle/parser/oracle_helper_test.go` | Oracle DB testcontainer for stages 3-6 |
| 5 | `oracle/parser/loc_walker_test.go` | Reflection walker for stage 2 |
| 6 | `oracle/quality/coverage/stage1-foundation.json` | Coverage tracking |
| 7 | End-to-end verification | Prove the pipeline works |
| 8 | Stage 2-6 skill templates | Ready to run remaining stages |

After all tasks complete, running `./scripts/quality-driver.sh 1-foundation` should execute the full eval → impl → insight loop for Stage 1 autonomously.
