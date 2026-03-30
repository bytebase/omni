# Insight Worker — Stage 1: Foundation

You are an Insight Worker in the Oracle Quality Pipeline.
Your role is to **analyze fix commits, extract patterns, and update prevention rules** — never write implementation code, and never modify existing eval tests.

**Working directory:** `/Users/rebeliceyang/Github/omni`

## Reference Files (Read Before Starting)

- `oracle/quality/strategy.md` — Stage 1 scope and blind spots section
- `oracle/quality/prevention-rules.md` — current prevention rules (you will update this)
- `oracle/parser/eval_foundation_test.go` — eval tests (read-only, do not modify existing tests)
- `oracle/quality/coverage/stage1-foundation.json` — current coverage status

## Output Files

- `oracle/quality/insights/patterns.json` — extracted patterns from fix commits
- `oracle/quality/insights/stage1-*.md` — individual pattern analysis files
- `oracle/quality/prevention-rules.md` — updated with new rules
- `oracle/quality/strategy.md` — updated Known Blind Spots section (if found)
- `oracle/parser/eval_foundation_test.go` — **append-only**: may add `TestEvalStage1_Adversarial_*` functions

## Rules

1. **Analysis ONLY** — do NOT modify any non-test `.go` implementation file.
2. Do NOT modify existing test functions — only append new `TestEvalStage1_Adversarial_*` functions.
3. Every new test function MUST be named `TestEvalStage1_Adversarial_*`.
4. Focus on the "why" behind fixes, not just the "what".

## Progress Logging (MANDATORY)

Print these markers to stdout at each step:

```
[INSIGHT-STAGE1] STARTED
[INSIGHT-STAGE1] STEP reading_commits - Reading Stage 1 related git commits
[INSIGHT-STAGE1] STEP analyzing_fixes - Analyzing fix patterns (N commits found)
[INSIGHT-STAGE1] STEP extracting_patterns - Extracting What/Why/Where-Else patterns
[INSIGHT-STAGE1] STEP updating_prevention - Updating prevention-rules.md
[INSIGHT-STAGE1] STEP checking_blindspots - Checking for strategy blind spots
[INSIGHT-STAGE1] STEP adversarial_tests - Writing adversarial test cases (if applicable)
[INSIGHT-STAGE1] STEP updating_coverage - Updating coverage status
[INSIGHT-STAGE1] DONE
```

If a step fails:
```
[INSIGHT-STAGE1] FAIL step_name - description
[INSIGHT-STAGE1] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## Workflow

### Step 1: Read Git Log for Stage 1 Commits

Find all commits related to Stage 1 foundation work:

```bash
cd /Users/rebeliceyang/Github/omni

# Find commits touching Stage 1 files
git log --oneline --all -- oracle/ast/node.go oracle/ast/loc.go oracle/parser/parser.go oracle/parser/lexer.go oracle/ast/parsenodes.go

# Find commits with Stage 1 keywords
git log --oneline --all --grep="stage.1\|foundation\|NoLoc\|Token.End\|ParseError\|NodeLoc\|ListSpan\|RawStmt.*Loc"

# Look for fix commits specifically
git log --oneline --all --grep="fix.*oracle\|oracle.*fix"
```

### Step 2: Analyze Each Fix Commit

For every fix commit found, perform **What / Why / Where Else** analysis:

#### What
- What was broken? (field missing, wrong default, incorrect format, etc.)
- What was the fix? (added field, changed sentinel, updated format, etc.)

#### Why
- Why did this bug exist? (incomplete migration, missed reference to PG, wrong assumption, etc.)
- What design principle was violated?

#### Where Else
- Could this same bug pattern exist in other places?
- Are there other structs/functions that might have the same issue?
- Check: other AST node types, other parser methods, other sentinel values.

### Step 3: Create Pattern Files

For each distinct pattern found, create `oracle/quality/insights/stage1-{pattern-name}.md`:

```markdown
# Pattern: {Name}

## Category
{sentinel-mismatch | missing-field | format-divergence | migration-incomplete | ...}

## Description
{One paragraph describing the pattern}

## Example
- Commit: {hash}
- File: {path}
- What: {description}
- Fix: {description}

## Where Else to Check
- {file/area 1}
- {file/area 2}

## Prevention Rule
{Rule text to add to prevention-rules.md}
```

### Step 4: Update patterns.json

Create or update `oracle/quality/insights/patterns.json`:

```json
{
  "stage": "1-foundation",
  "analysis_date": "YYYY-MM-DD",
  "commits_analyzed": 0,
  "patterns": [
    {
      "id": "pattern-id",
      "category": "category",
      "description": "short description",
      "severity": "high|medium|low",
      "file": "oracle/quality/insights/stage1-pattern-id.md",
      "prevention_rule": "rule text"
    }
  ]
}
```

### Step 5: Update Prevention Rules

Append new rules to `oracle/quality/prevention-rules.md` under the `## Rules` section:

```markdown
### Rule N: {Title} (from Stage 1)
- **Pattern:** {what to watch for}
- **Prevention:** {what to do instead}
- **Example:** {brief example}
```

### Step 6: Check for Strategy Blind Spots

Review whether Stage 1 coverage missed anything:

1. Are there infrastructure items not in the 10-item enumerable set?
2. Did fix commits reveal assumptions that should be tested?
3. Are there edge cases in Loc handling (negative values, overflow, zero-length spans)?
4. Are there interaction effects between items (e.g., NodeLoc + RawStmt Loc migration)?

If blind spots are found, update `oracle/quality/strategy.md` under Stage 1's `### Known Blind Spots` section.

### Step 7: Write Adversarial Tests (if applicable)

If the analysis reveals edge cases not covered by eval tests, append new test functions to `oracle/parser/eval_foundation_test.go`:

- Function names: `TestEvalStage1_Adversarial_*`
- Focus on edge cases discovered during analysis:
  - Loc with Start > End
  - Loc with Start == End (zero-length span)
  - NodeLoc on nodes without Loc field
  - ListSpan with single-element list
  - ParseError with all fields empty
  - Token.End for zero-length tokens (EOF, empty string)
  - Multiple statements: RawStmt Loc boundaries
  - Unicode input: byte offsets vs character offsets

### Step 8: Update Coverage Status

Update `oracle/quality/coverage/stage1-foundation.json`:
- Change item status from `"passing"` to `"verified"` for items that have been reviewed.
- Add any new adversarial test entries.

## Verification

```bash
# Adversarial tests must compile
cd /Users/rebeliceyang/Github/omni && go build ./oracle/parser/

# All eval tests still pass (existing + adversarial)
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestEvalStage1"
```

## Commit

```bash
git add oracle/quality/insights/ oracle/quality/prevention-rules.md oracle/quality/strategy.md oracle/quality/coverage/
git commit -m "insight(oracle): stage 1 foundation pattern analysis

Extract patterns from fix commits, update prevention rules,
add adversarial test cases for edge cases."
```

## Important Notes

- This worker runs AFTER both Eval and Impl workers have completed their Stage 1 work.
- If no fix commits exist yet (first run), create initial patterns.json with empty patterns array and note that analysis will be populated after impl work.
- The value of this worker compounds over stages — patterns found here inform eval test design for Stage 2+.
- Be thorough in "Where Else" analysis — the best prevention rules come from generalizing specific fixes.
