# MySQL Implicit Bugs Phase 1 Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reduce the current MySQL implicit-behavior divergence set by fixing low-risk catalog naming, ALTER behavior, table-option preservation, and SHOW CREATE round-trip gaps first.

**Architecture:** Use the existing `mysql/catalog/scenarios_*_test.go` oracle scenario tests as the source of truth. Fix behavior in catalog ingestion and deparse layers before larger semantic analyzer work. Keep each task scoped to one shared root cause and verify with the narrow failing scenario first, then the full scenario suite.

**Tech Stack:** Go, `go test`, MySQL 8.0 testcontainers, existing `mysql/catalog` catalog/deparse code.

## Phase 0: Status Hygiene

### Task 0.1: Record Current Baseline

**Files:**
- Read: `/tmp/mysql_catalog_scenarios_20260424.json`
- Read: `mysql/catalog/scenarios_bug_queue/*.md`

**Step 1: Run baseline**

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_|TestBugFix_' -count=1 -timeout=30m -json > /tmp/mysql_catalog_scenarios_20260424.json
```

Expected: command exits non-zero while historical scenario failures remain.

**Step 2: Extract failing scenario list**

Run:
```bash
jq -r 'select(.Action=="fail" and .Test!=null and (.Test|test("^TestScenario_.*?/"))) | .Test' /tmp/mysql_catalog_scenarios_20260424.json
```

Expected: list includes AX.3, AX.9, C1.8/C1.10/C1.13, C18.9-C18.13, and other known failures.

## Phase 1: Naming, ALTER, and Round-Trip Metadata

### Task 1.1: Fix ALTER DROP COLUMN Last-Column Rejection

**Files:**
- Modify: `mysql/catalog/altercmds.go`
- Test: `mysql/catalog/scenarios_ax_test.go`

**Step 1: Verify red**

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_AX/AX_3_DropColumn_rejects_last_column' -count=1
```

Expected: FAIL because omni allows dropping the last column.

**Step 2: Implement**

Before removing a column, compute whether the table would have zero columns after the drop. Return MySQL-compatible error 1090 (`ER_CANT_REMOVE_ALL_FIELDS`) before mutating the table.

**Step 3: Verify green**

Run the same test and expect PASS.

### Task 1.2: Ignore Column-Level REFERENCES in CREATE TABLE

**Files:**
- Modify: `mysql/catalog/tablecmds.go`
- Test: `mysql/catalog/scenarios_ax_test.go`

**Step 1: Verify red**

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_AX/AX_9_FK_column_shorthand_silent_ignore' -count=1
```

Expected: FAIL because omni creates an FK for column-level `REFERENCES`.

**Step 2: Implement**

Do not convert column-level `REFERENCES` constraints into catalog FKs. Only table-level `FOREIGN KEY (...) REFERENCES ...` should create FK constraints.

**Step 3: Verify green**

Run the same test and expect PASS.

### Task 1.3: Fix Index and CHECK/FK Name Semantics

**Files:**
- Modify: `mysql/catalog/tablecmds.go`
- Modify: `mysql/catalog/altercmds.go`
- Test: `mysql/catalog/scenarios_c1_test.go`
- Test: `mysql/catalog/scenarios_ps_test.go`

**Step 1: Verify red**

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_C1/(1_8|1_10|1_13)|TestScenario_PS/(PS_2|PS_7|PS_8)' -count=1
```

Expected: FAIL on reserved `PRIMARY`, unique fallback, schema-scoped CHECK, ALTER CHECK max+1, FK collision, and CHECK dup-name cases.

**Step 2: Implement**

Add helpers for:
- Rejecting non-primary indexes named `PRIMARY`.
- Making auto-generated index names skip bare `PRIMARY` and start with `PRIMARY_2`.
- Finding CHECK constraints by name across a database.
- Detecting FK name collisions in a single CREATE TABLE before insertion.
- Generating ALTER ADD CHECK names from max existing `<table>_chk_N`.

**Step 3: Verify green**

Run the same test and expect PASS for fixed subtests.

### Task 1.4: Preserve Table Options and Render SHOW CREATE

**Files:**
- Modify: `mysql/catalog/table.go`
- Modify: `mysql/catalog/tablecmds.go`
- Modify: `mysql/catalog/show.go`
- Test: `mysql/catalog/scenarios_c8_test.go`
- Test: `mysql/catalog/scenarios_c18_test.go`

**Step 1: Verify red**

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_C18/(18_9|18_10|18_11|18_12|18_13)' -count=1
```

Expected: FAIL because explicit table options are dropped or not rendered.

**Step 2: Implement**

Add catalog fields for `COMPRESSION`, `ENCRYPTION`, `STATS_*`, `TABLESPACE`, `MIN_ROWS`, `MAX_ROWS`, `AVG_ROW_LENGTH`, `PACK_KEYS`, `CHECKSUM`, and `DELAY_KEY_WRITE`. Populate them from CREATE TABLE options and emit them only when explicit according to MySQL SHOW CREATE behavior.

**Step 3: Verify green**

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_C8|TestScenario_C18' -count=1
```

Expected: C8 remains green and C18.9-C18.13 turn green.

## Phase 1 Completion Gate

Run:
```bash
go test ./mysql/catalog -run 'TestScenario_AX|TestScenario_C1|TestScenario_C8|TestScenario_C18|TestScenario_PS|TestBugFix_' -count=1 -timeout=20m
```

Expected: the Phase 1 target tests pass except explicitly deferred functional-index cases (`C1.11`, `C1.12`) and partition/date-time items (`PS.5`, `PS.6`) that belong to later phases.
