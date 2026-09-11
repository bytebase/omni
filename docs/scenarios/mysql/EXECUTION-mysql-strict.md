# Execution Contract: MySQL Parser Strictness Alignment

> Generated: 2026-03-31
> Source: SCENARIOS-mysql-strict.md
> Proof command prefix: `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v`

---

## Step 1: File Targets Per Section

| Section | Parser files modified | Corpus files modified | Test file |
|---------|----------------------|----------------------|-----------|
| 1.1 Main Loop Token Consumption | `mysql/parser/parser.go` | `mysql/quality/corpus/01-select.sql` | `mysql/parser/verify_parse_test.go` (unchanged, just runs) |
| 2.1 SELECT — JOIN Keyword Enforcement | `mysql/parser/select.go` | `mysql/quality/corpus/01-select.sql` | — |
| 2.2 SELECT — Clause Validation | `mysql/parser/select.go` | `mysql/quality/corpus/01-select.sql` | — |
| 2.3 DELETE — FROM Enforcement | `mysql/parser/update_delete.go` | `mysql/quality/corpus/02-dml.sql` | — |
| 2.4 INSERT/REPLACE — Value Source Enforcement | `mysql/parser/insert.go` | `mysql/quality/corpus/02-dml.sql` | — |
| 2.5 UPDATE — SET Enforcement | `mysql/parser/update_delete.go` | `mysql/quality/corpus/02-dml.sql` | — |
| 2.6 DDL — IF NOT EXISTS / IF EXISTS Keyword Chains | `mysql/parser/create_table.go`, `mysql/parser/create_index.go`, `mysql/parser/drop.go`, `mysql/parser/create_database.go` | `mysql/quality/corpus/03-ddl.sql` | — |
| 2.7 DDL — Column Constraint Keyword Chains | `mysql/parser/create_table.go` | `mysql/quality/corpus/03-ddl.sql` | — |
| 2.8 DDL — Structural Delimiter Enforcement | `mysql/parser/create_table.go` | `mysql/quality/corpus/03-ddl.sql` | — |
| 3.1 Parenthesis Balance | `mysql/parser/expr.go` | `mysql/quality/corpus/01-select.sql` | — |
| 3.2 Binary Operator Right-Operand | `mysql/parser/expr.go` | `mysql/quality/corpus/01-select.sql` | — |
| 3.3 IN/BETWEEN/LIKE Expression Completeness | `mysql/parser/expr.go` | `mysql/quality/corpus/01-select.sql` | — |
| 3.4 CASE Expression Completeness | `mysql/parser/expr.go` | `mysql/quality/corpus/01-select.sql` | — |

---

## Step 2: Three-Independence Analysis

### Phase 1 (single section — no pair analysis needed)

Section 1.1 is the only section in Phase 1. It modifies the main parse loop in `parser.go`. No other section touches `parser.go`. Phase 1 must complete before Phase 2 because the trailing-token rejection changes the fundamental contract of the parse loop, and every subsequent section's "still accepted" tests depend on the main loop not introducing false rejections.

### Phase 2: Pairwise Independence

**2.1 vs 2.2** — CONFLICT
- Semantic: both modify SELECT-related parsing
- Change-surface: both modify `mysql/parser/select.go` — 2.1 modifies `matchJoinType()`, 2.2 modifies `parseSelectStmtBase()`
- Proof: both add to `mysql/quality/corpus/01-select.sql`
- Verdict: **Change-surface conflict on select.go.** Although the functions are different, both touch the same file. SEQUENTIAL to avoid merge conflicts.

**2.1+2.2 vs 2.3** — INDEPENDENT
- Semantic: SELECT join/clause vs DELETE FROM enforcement — unrelated
- Change-surface: select.go vs update_delete.go; corpus 01-select.sql vs 02-dml.sql — no overlap
- Proof: independent corpus files
- Verdict: **Fully independent. PARALLEL safe.**

**2.1+2.2 vs 2.4** — INDEPENDENT
- Semantic: SELECT vs INSERT/REPLACE — unrelated
- Change-surface: select.go vs insert.go; corpus 01-select.sql vs 02-dml.sql — no overlap
- Proof: independent corpus files
- Verdict: **Fully independent. PARALLEL safe.**

**2.1+2.2 vs 2.5** — INDEPENDENT
- Semantic: SELECT vs UPDATE — unrelated
- Change-surface: select.go vs update_delete.go; corpus 01-select.sql vs 02-dml.sql — no overlap
- Proof: independent corpus files
- Verdict: **Fully independent. PARALLEL safe.**

**2.3 vs 2.4** — INDEPENDENT
- Semantic: DELETE vs INSERT — unrelated
- Change-surface: update_delete.go vs insert.go; both write to 02-dml.sql — **corpus conflict**
- Proof: shared corpus file 02-dml.sql
- Verdict: **Corpus conflict on 02-dml.sql.** Corpus files are additive-only (append new entries), so this is safe for parallel if we mandate append-only. However, to be conservative: SEQUENTIAL within the DML group.

**2.3 vs 2.5** — CONFLICT
- Semantic: DELETE vs UPDATE — different statements
- Change-surface: **both modify `mysql/parser/update_delete.go`** — 2.3 modifies `parseDeleteStmt()`, 2.5 modifies `parseUpdateStmt()`
- Proof: both write to 02-dml.sql
- Verdict: **Change-surface conflict on update_delete.go.** SEQUENTIAL.

**2.4 vs 2.5** — INDEPENDENT (after 2.3)
- Semantic: INSERT vs UPDATE — unrelated
- Change-surface: insert.go vs update_delete.go; both write to 02-dml.sql — corpus overlap only
- Proof: shared corpus file but additive
- Verdict: **Independent on parser files.** Corpus is additive. PARALLEL safe if 2.3 is already done (since 2.5 shares update_delete.go with 2.3).

**2.6 vs 2.7** — CONFLICT
- Semantic: both DDL strictness
- Change-surface: **both modify `mysql/parser/create_table.go`** — 2.6 modifies `parseCreateTableStmt()` IF NOT EXISTS handling, 2.7 modifies `parseColumnOption()` NOT NULL / PRIMARY KEY handling
- Proof: both write to 03-ddl.sql
- Verdict: **Change-surface conflict on create_table.go.** SEQUENTIAL.

**2.6 vs 2.8** — CONFLICT
- Semantic: both DDL strictness
- Change-surface: **both modify `mysql/parser/create_table.go`** — 2.6 modifies IF NOT EXISTS, 2.8 modifies structural delimiter (closing paren) enforcement
- Proof: both write to 03-ddl.sql
- Verdict: **Change-surface conflict on create_table.go.** SEQUENTIAL.

**2.7 vs 2.8** — CONFLICT
- Semantic: both DDL column/table definition strictness
- Change-surface: **both modify `mysql/parser/create_table.go`**
- Proof: both write to 03-ddl.sql
- Verdict: **Change-surface conflict on create_table.go.** SEQUENTIAL.

**DML group (2.3-2.5) vs DDL group (2.6-2.8)** — INDEPENDENT
- Change-surface: update_delete.go + insert.go vs create_table.go + create_index.go + drop.go + create_database.go — no overlap
- Corpus: 02-dml.sql vs 03-ddl.sql — no overlap
- Verdict: **Fully independent. PARALLEL safe.**

**SELECT group (2.1-2.2) vs DDL group (2.6-2.8)** — INDEPENDENT
- Change-surface: select.go vs create_table.go + create_index.go + drop.go + create_database.go — no overlap
- Corpus: 01-select.sql vs 03-ddl.sql — no overlap
- Verdict: **Fully independent. PARALLEL safe.**

### Phase 3: Pairwise Independence

**3.1 vs 3.2** — CONFLICT
- Change-surface: both modify `mysql/parser/expr.go`
- Verdict: **SEQUENTIAL.**

**3.1 vs 3.3** — CONFLICT
- Change-surface: both modify `mysql/parser/expr.go`
- Verdict: **SEQUENTIAL.**

**3.1 vs 3.4** — CONFLICT
- Change-surface: both modify `mysql/parser/expr.go`
- Verdict: **SEQUENTIAL.**

**3.2 vs 3.3** — CONFLICT
- Change-surface: both modify `mysql/parser/expr.go`
- Verdict: **SEQUENTIAL.**

**3.2 vs 3.4** — CONFLICT
- Change-surface: both modify `mysql/parser/expr.go`
- Verdict: **SEQUENTIAL.**

**3.3 vs 3.4** — CONFLICT
- Change-surface: both modify `mysql/parser/expr.go`
- Verdict: **SEQUENTIAL.**

All Phase 3 sections share `mysql/parser/expr.go`. **Entire phase is SEQUENTIAL.**

---

## Step 3: Execution Shape

### Phase 1: Sequential (single section)

```
1.1
```

### Phase 2: Three parallel lanes, each internally sequential

```
Lane A (SELECT):    2.1 → 2.2
Lane B (DML):       2.3 → 2.5 → 2.4
Lane C (DDL):       2.6 → 2.7 → 2.8
```

Rationale for Lane B ordering: 2.3 and 2.5 both touch `update_delete.go`, so they must be sequential. 2.3 goes first because DELETE FROM enforcement is a prerequisite pattern for UPDATE SET enforcement (same file, simpler change first). 2.4 (INSERT) is independent of 2.5 on parser files but shares the corpus file with 2.3/2.5, and is placed last to avoid interleaving corpus writes.

Rationale for Lane C ordering: 2.6 touches 4 files (create_table.go, create_index.go, drop.go, create_database.go) and must complete first since 2.7 and 2.8 also touch create_table.go. 2.7 before 2.8 because constraint keywords (NOT NULL, PRIMARY KEY) are a prerequisite for structural delimiter enforcement.

### Phase 3: Fully sequential

```
3.1 → 3.2 → 3.3 → 3.4
```

All sections modify `mysql/parser/expr.go`. No parallelism possible.

---

## Step 4: Proof Checkpoints

| Checkpoint | Trigger | Command | Pass criteria |
|-----------|---------|---------|---------------|
| P1-done | After 1.1 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS, 0 new PARSE LENIENT from 1.1 corpus entries |
| P2-A1 | After 2.1 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; JOIN keyword corpus entries pass |
| P2-A2 | After 2.2 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; clause validation corpus entries pass |
| P2-B1 | After 2.3 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; DELETE corpus entries pass |
| P2-B2 | After 2.5 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; UPDATE corpus entries pass |
| P2-B3 | After 2.4 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; INSERT/REPLACE corpus entries pass |
| P2-C1 | After 2.6 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; IF NOT EXISTS / IF EXISTS corpus entries pass |
| P2-C2 | After 2.7 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; column constraint corpus entries pass |
| P2-C3 | After 2.8 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; structural delimiter corpus entries pass |
| P2-gate | All lanes done | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | Full Phase 2 gate: 0 crashes, 0 PARSE VIOLATIONS, PARSE LENIENT count decreased or unchanged |
| P3-1 | After 3.1 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; parenthesis balance corpus entries pass |
| P3-2 | After 3.2 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; binary operator corpus entries pass |
| P3-3 | After 3.3 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; IN/BETWEEN/LIKE corpus entries pass |
| P3-4 | After 3.4 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v` | 0 crashes, 0 PARSE VIOLATIONS; CASE expression corpus entries pass |
| FINAL | After 3.4 | `go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v && go test ./mysql/parser/ -count=1` | 0 crashes, 0 PARSE VIOLATIONS, 0 PARSE LENIENT from all new corpus entries, full test suite green |

---

## Step 5: Execution Contract

```
EXECUTION CONTRACT
==================

Phase 1 — Infrastructure
  depends_on: []
  shape: sequential
  sections:
    - id: "1.1"
      name: "Main Loop Token Consumption"
      modifies:
        - mysql/parser/parser.go
        - mysql/quality/corpus/01-select.sql
      proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
      checkpoint: P1-done

Phase 2 — Statement-Level Strictness
  depends_on: [Phase 1]
  shape: parallel_lanes

  lane_A:
    name: "SELECT strictness"
    sections:
      - id: "2.1"
        name: "SELECT — JOIN Keyword Enforcement"
        modifies:
          - mysql/parser/select.go
          - mysql/quality/corpus/01-select.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-A1

      - id: "2.2"
        name: "SELECT — Clause Validation"
        depends_on: ["2.1"]
        modifies:
          - mysql/parser/select.go
          - mysql/quality/corpus/01-select.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-A2

  lane_B:
    name: "DML strictness"
    sections:
      - id: "2.3"
        name: "DELETE — FROM Enforcement"
        modifies:
          - mysql/parser/update_delete.go
          - mysql/quality/corpus/02-dml.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-B1

      - id: "2.5"
        name: "UPDATE — SET Enforcement"
        depends_on: ["2.3"]
        modifies:
          - mysql/parser/update_delete.go
          - mysql/quality/corpus/02-dml.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-B2

      - id: "2.4"
        name: "INSERT/REPLACE — Value Source Enforcement"
        depends_on: ["2.5"]
        modifies:
          - mysql/parser/insert.go
          - mysql/quality/corpus/02-dml.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-B3

  lane_C:
    name: "DDL strictness"
    sections:
      - id: "2.6"
        name: "DDL — IF NOT EXISTS / IF EXISTS Keyword Chains"
        modifies:
          - mysql/parser/create_table.go
          - mysql/parser/create_index.go
          - mysql/parser/drop.go
          - mysql/parser/create_database.go
          - mysql/quality/corpus/03-ddl.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-C1

      - id: "2.7"
        name: "DDL — Column Constraint Keyword Chains"
        depends_on: ["2.6"]
        modifies:
          - mysql/parser/create_table.go
          - mysql/quality/corpus/03-ddl.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-C2

      - id: "2.8"
        name: "DDL — Structural Delimiter Enforcement"
        depends_on: ["2.7"]
        modifies:
          - mysql/parser/create_table.go
          - mysql/quality/corpus/03-ddl.sql
        proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
        checkpoint: P2-C3

  gate: P2-gate
    trigger: all lanes complete
    proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
    pass: 0 crashes, 0 PARSE VIOLATIONS, PARSE LENIENT unchanged or reduced

Phase 3 — Expression Strictness
  depends_on: [Phase 2]
  shape: sequential
  sections:
    - id: "3.1"
      name: "Parenthesis Balance"
      modifies:
        - mysql/parser/expr.go
        - mysql/quality/corpus/01-select.sql
      proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
      checkpoint: P3-1

    - id: "3.2"
      name: "Binary Operator Right-Operand"
      depends_on: ["3.1"]
      modifies:
        - mysql/parser/expr.go
        - mysql/quality/corpus/01-select.sql
      proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
      checkpoint: P3-2

    - id: "3.3"
      name: "IN/BETWEEN/LIKE Expression Completeness"
      depends_on: ["3.2"]
      modifies:
        - mysql/parser/expr.go
        - mysql/quality/corpus/01-select.sql
      proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
      checkpoint: P3-3

    - id: "3.4"
      name: "CASE Expression Completeness"
      depends_on: ["3.3"]
      modifies:
        - mysql/parser/expr.go
        - mysql/quality/corpus/01-select.sql
      proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v
      checkpoint: P3-4

  gate: FINAL
    trigger: 3.4 complete
    proof: go test ./mysql/parser/ -run TestVerifyCorpus -count=1 -v && go test ./mysql/parser/ -count=1
    pass: 0 crashes, 0 PARSE VIOLATIONS, 0 PARSE LENIENT from new entries, full test suite green

CRITICAL NOTES
==============

1. Phase ordering is strict: Phase 1 → Phase 2 → Phase 3. No phase may
   start until the prior phase gate passes.

2. Within Phase 2, lanes A/B/C run in parallel. Within each lane,
   sections are strictly sequential.

3. Every section MUST add corpus entries BEFORE modifying parser code.
   The pattern is: (a) add @valid: false entries to corpus, (b) verify
   they produce PARSE LENIENT, (c) modify parser, (d) verify they now
   produce parseExpectedFail. This ensures the corpus drives the change.

4. No section may modify a file outside its declared "modifies" list.

5. The verify_parse_test.go file is NOT modified by any section — it is
   infrastructure that already exists and is used as-is.

6. Existing valid SQL in corpus files must NEVER be broken. Every proof
   checkpoint verifies 0 PARSE VIOLATIONS (valid SQL that parser rejects).

7. For section 2.6, the IF NOT EXISTS / IF EXISTS changes span 4 files.
   The current code uses permissive p.match() calls that silently skip
   missing keywords:
     - create_table.go:30   → p.match(kwNOT); p.match(kwEXISTS_KW)
     - create_index.go:31   → p.match(kwNOT); p.match(kwEXISTS_KW)
     - drop.go:24           → p.match(kwEXISTS_KW)
     - drop.go:133          → p.match(kwEXISTS_KW)
     - create_database.go:112 → p.match(kwEXISTS_KW)
   These must be changed to p.expect() with appropriate error handling.

8. For section 2.1, the current matchJoinType() uses p.match(kwJOIN)
   which silently succeeds even when JOIN is absent. These must become
   p.expect(kwJOIN) calls to enforce the keyword chain.

DAG SUMMARY (topological order)
================================

   1.1
    |
    +------ Phase 1 gate ------+
    |            |              |
   2.1          2.3            2.6
    |            |              |
   2.2          2.5            2.7
                 |              |
                2.4            2.8
    |            |              |
    +------ Phase 2 gate ------+
                 |
                3.1
                 |
                3.2
                 |
                3.3
                 |
                3.4
                 |
            FINAL gate
```
