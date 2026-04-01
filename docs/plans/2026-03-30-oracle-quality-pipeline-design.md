# Oracle Engine Quality Pipeline Design

## Goal

Build a unified quality pipeline for the Oracle database engine, using PG's feature surfaces as the reference template. The pipeline produces strict external evaluation, systematic implementation, and a feedback loop that evolves the evaluation system itself over time.

## Core Principles

1. **Evaluation before implementation** — each feature surface gets eval tests first, implementation second.
2. **Role separation** — eval tests and implementation are written by different workers with opposing goals.
3. **Layered dependency** — 6 stages run serially; each stage's eval tests become permanent regression guards for all subsequent stages.
4. **Three-level feedback** — insights from implementation feed back into (a) implementation code, (b) eval test coverage, and (c) the coverage strategy itself.

## Stages

| Stage | Surface | Depends On | Evaluation Strategy |
|-------|---------|------------|---------------------|
| 1 | Foundation (Loc semantics, Token.End, ParseError) | — | PG behavior reference |
| 2 | Loc completeness (248 nodes) | Stage 1 | Reflection walker + substring verification |
| 3 | AST correctness | Stage 2 | Oracle DB cross-validation + structural assertions |
| 4 | Error quality (soft-fail, message, position) | Stage 3 | Mutation generation + Oracle DB cross-validation |
| 5 | Completion | Stage 4 | PG behavior reference + candidate set assertions |
| 6 | Catalog / Migration | Stage 5 | Oracle DB oracle verification + pairwise combinations |

## Worker Roles

### Eval Worker
- **Goal**: Write maximally strict and comprehensive eval tests.
- **Can read**: AST definitions, BNF grammars, Oracle documentation, PG behavior, coverage strategy.
- **Can write**: `*_eval_test.go`, `coverage.json`.
- **Cannot**: Write implementation code.
- **Enumerable surfaces**: Must cover every item in the enumerable set (BNF rules, AST node types). Gaps tracked in `coverage.json`.
- **Non-enumerable surfaces**: Use generation strategies (mutation, corpus, pairwise). Test data cross-validated against Oracle DB.

### Impl Worker
- **Goal**: Make eval tests pass.
- **Can read**: Eval tests (read-only), existing implementation code, prevention rules.
- **Can write**: Implementation code only.
- **Cannot**: Modify `*_eval_test.go`.
- **Constraint**: Every commit must pass all eval tests from current and all prior stages.

### Insight Worker
- **Goal**: Extract patterns from fix commits, strengthen the evaluation system.
- **Can read**: Fix commit diffs, eval tests, implementation code.
- **Produces**:
  - New patterns → `insights/patterns.json` + `insights/<pattern>.md`
  - Adversarial test cases → appended to eval tests
  - Prevention rules → injected into future Impl Worker instructions
  - Strategy corrections → updates to `strategy.md`

## Stage Internal Flow

```
Step 0: Driver prep
  - Read stage's coverage strategy (strategy.md)
  - Read upstream prevention rules
  - Verify all prior stage eval tests still pass

Step 1: Eval Worker
  - Generate eval tests from strategy + reference sources
  - Produce coverage.json with gap tracking
  - Cross-validate test data against Oracle DB
  - All tests expected to FAIL (no implementation yet)

Step 2: Impl Worker (starmap decomposition, multiple rounds)
  - Receive failing test list + prevention rules
  - Implement in batches, each batch must pass eval tests
  - Regression guard: prior stage tests must pass at every commit

Step 3: Insight Worker
  - Analyze all fix commits from Step 2
  - Extract patterns, generate adversarial tests, update prevention rules
  - If strategy blind spots found → update strategy.md
  - If new test cases generated → Driver returns to Step 2 for new failures

Step 4: Driver verification
  - coverage.json gap_rules is empty
  - All eval tests pass (current + all prior stages)
  - Mark stage done, proceed to next
```

## Coverage Strategy Matrix

| Surface | Enumerable Part | Strategy | Non-enumerable Part | Strategy |
|---------|----------------|----------|--------------------|---------|
| AST correctness | BNF rules (230+) | 1 test per rule minimum | Non-standard user SQL | Oracle DB corpus validation |
| Loc tracking | 248 Loc-bearing nodes | 1 test per node type | — | Fully enumerable |
| Error detection | — | — | All error variations | Mutation from valid SQL (truncate, delete, replace, duplicate) |
| Error messages | BNF rule incomplete variants | 1 error test per rule | Complex nested errors | Mutation generation |
| Completion | Statement types x cursor positions | Enumerate key positions | Complex context | Nested SQL corpus |
| Catalog | Object types x operations | Enumerate basics | Property combinations | Pairwise coverage |

Strategy is maintained in `oracle/quality/strategy.md` and updated by Insight Worker when blind spots are discovered.

## Evaluation Quality Assurance

Three mechanisms ensure eval tests themselves are good enough:

1. **Anchoring to enumerable sources** — coverage.json tracks gaps against known full sets. Driver blocks stage completion if gaps remain.
2. **Oracle DB cross-validation** — test SQL is executed against real Oracle DB. DB-accepted SQL must be parseable; DB-rejected SQL should error. Invalid test data is caught.
3. **Insight accumulation** — every fix commit potentially generates new adversarial test cases. Coverage grows monotonically.

When a fix reveals that the coverage strategy itself has a blind spot (e.g., Oracle-specific idioms not in BNF), the Insight Worker updates `strategy.md`, adding a new coverage dimension. Next Eval Worker run picks up the new dimension.

## File Structure

```
oracle/
  quality/
    strategy.md                    # Coverage strategy matrix (iterable)
    coverage/
      stage1-foundation.json
      stage2-loc.json
      stage3-ast.json
      stage4-error.json
      stage5-completion.json
      stage6-catalog.json
    insights/
      patterns.json                # Pattern index
      <pattern-name>.md            # Detailed pattern records
    prevention-rules.md            # Accumulated prevention rules across stages
    corpus/
      valid-oracle.sql             # Oracle DB validated legal SQL
      idiomatic-oracle.sql         # Oracle-specific idiom corpus
```

## Regression Guard

- Impl Worker runs `go test ./oracle/... -run Eval` (all stages) before every batch completion.
- Stage N changes that break Stage M (M < N) tests must be fixed before proceeding.
- Driver enforces: no batch marked done unless full eval suite passes.

## Feedback Loop Summary

```
Implementation → discover issue → Insight analysis
                                    ├── root cause in implementation → fix code
                                    ├── root cause in test coverage → add tests
                                    └── root cause in strategy blind spot → update strategy → derive new tests → fix code
```

All three feedback paths use the same Driver → Worker pipeline. The only difference is the source of work items.
