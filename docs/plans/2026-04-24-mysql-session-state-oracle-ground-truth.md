# MySQL Session-State and Oracle Ground-Truth Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Resolve the final pending MySQL implicit-behavior scenarios by modeling the session state that changes DDL semantics and by converting invalid oracle assumptions into explicit MySQL-aligned rejection tests.

**Architecture:** Treat `Catalog` as one MySQL session: `SET` statements mutate catalog-local session state, and DDL reads that state when MySQL behavior depends on it. For unsupported or non-MySQL syntax, the test oracle must prove rejection and omni should reject rather than carrying a skipped "no ground truth" scenario.

**Tech Stack:** Go, `mysql/catalog`, `mysql/parser`, MySQL 8.0 oracle container via testcontainers.

## Context

The progress tracker in `mysql/catalog/SCENARIOS-mysql-implicit-behavior.md` now has only three pending rows:

- `3.8` second `TIMESTAMP` legacy implicit zero default, `LOW`.
- `6.10` `LIST DEFAULT` partition catch-all, `MED`.
- `16.12` `TIMESTAMP` promotion inherits fsp, `MED`.

`3.8` and `16.12` are blocked by the same missing catalog capability: `SET SESSION explicit_defaults_for_timestamp=0` has no effect on omni. `6.10` is different: the current MySQL oracle rejects `VALUES IN (DEFAULT)`, and official MySQL partitioning docs describe LIST partition value lists as integer lists, not a DEFAULT catch-all partition. Treating this as an implementation gap would add non-MySQL behavior.

Useful references:

- MySQL data type defaults: https://dev.mysql.com/doc/mysql/8.0/en/data-type-defaults.html
- MySQL LIST partitioning: https://dev.mysql.com/doc/refman/en/partitioning-list.html

## Design

Add an explicit, small session-state model to `Catalog` instead of scattering individual booleans forever. Keep current fields working initially, but introduce a struct so future settings do not keep expanding the top-level catalog shape.

```go
type SessionState struct {
    ForeignKeyChecks              bool
    GenerateInvisiblePrimaryKey   bool
    ShowGIPK                      bool
    CharsetClient                 string
    CollationConnection           string
    ExplicitDefaultsForTimestamp  bool
    SQLMode                       string
}
```

Default `ExplicitDefaultsForTimestamp` to `true`, matching MySQL 8.0 default behavior. `SET SESSION explicit_defaults_for_timestamp=0`, `SET @@session.explicit_defaults_for_timestamp=0`, and bare `SET explicit_defaults_for_timestamp=0` should all update this field. `SET ... = DEFAULT` should restore `true`. Store `SQLMode` as raw text for now; do not implement broad SQL mode semantics in this task.

The timestamp DDL path should call a helper after column attributes are parsed and before table constraints/validation finish:

```go
func (c *Catalog) applyTimestampSessionDefaults(cols []*Column, defs []*nodes.ColumnDef) {
    if c.session.ExplicitDefaultsForTimestamp {
        return
    }
    // Legacy MySQL behavior:
    // - TIMESTAMP columns without explicit NULL become NOT NULL.
    // - The first TIMESTAMP without explicit DEFAULT/ON UPDATE gets
    //   DEFAULT CURRENT_TIMESTAMP[(fsp)] ON UPDATE CURRENT_TIMESTAMP[(fsp)].
    // - Later TIMESTAMP columns without explicit DEFAULT get zero timestamp default.
}
```

The helper needs small predicate helpers: `hasExplicitNull`, `hasExplicitDefault`, `hasExplicitOnUpdate`, and `timestampCurrentExpr(fsp int)`. Do not infer from the final `Column` alone because the final nullable/default state loses whether the user explicitly wrote `NULL`.

For `6.10`, replace the skipped test with an oracle rejection test. If the parser accepts `DEFAULT` as a `DefaultExpr` inside partition values, add validation that rejects `DEFAULT` for `LIST` / `LIST COLUMNS` partition definitions. This keeps parser permissiveness acceptable only if catalog strictness matches MySQL.

## Task 1: Lock Oracle Ground Truth for `6.10`

**Files:**
- Modify: `mysql/catalog/scenarios_c6_test.go`
- Modify: `mysql/catalog/tablecmds.go`
- Modify: `mysql/catalog/SCENARIOS-mysql-implicit-behavior.md`

**Step 1: Write the failing test**

Change `TestScenario_C6/6_10_list_default_partition` so it no longer skips on oracle rejection. It should assert MySQL rejects the DDL and omni rejects the same DDL.

```go
ddl := `CREATE TABLE t (c INT) PARTITION BY LIST(c)
    (PARTITION p0 VALUES IN (1,2), PARTITION pd VALUES IN (DEFAULT))`

if _, err := mc.db.ExecContext(mc.ctx, ddl); err == nil {
    t.Fatalf("oracle: expected MySQL to reject LIST DEFAULT partition syntax")
}

results, err := c.Exec(ddl, nil)
omniErr := firstExecErr(results, err)
if omniErr == nil {
    t.Fatalf("omni: expected rejection for LIST DEFAULT partition syntax")
}
```

Use an existing local error helper if one exists; otherwise add a small C6-local helper.

**Step 2: Run test to verify it fails**

Run:

```bash
go test ./mysql/catalog -run 'TestScenario_C6/6_10' -count=1 -timeout=20m
```

Expected: FAIL if omni currently accepts `DEFAULT` in LIST partition values.

**Step 3: Implement minimal rejection**

In `validatePartitionClause`, add a recursive check for `*nodes.DefaultExpr` inside `pc.Partitions[*].Values` when `pc.Type == nodes.PartitionList`.

```go
if pc.Type == nodes.PartitionList {
    for _, pd := range pc.Partitions {
        if containsPartitionDefaultExpr(pd.Values) {
            return &Error{Code: 1064, SQLState: "42000", Message: "DEFAULT is not valid in LIST partition values"}
        }
    }
}
```

**Step 4: Update tracker**

Mark `6.10` verified, with a note like:

```markdown
| 6.10 | LIST DEFAULT partition catch-all | verified | MED | oracle rejection verified; MySQL has no DEFAULT catch-all LIST partition |
```

Also update the section body to remove the false MySQL support claim.

**Step 5: Verify and commit**

Run:

```bash
go test ./mysql/catalog -run 'TestScenario_C6/6_10' -count=1 -timeout=20m
go test ./mysql/catalog -run 'TestScenario_C6$' -count=1 -timeout=30m
git diff --check
git add mysql/catalog/scenarios_c6_test.go mysql/catalog/tablecmds.go mysql/catalog/SCENARIOS-mysql-implicit-behavior.md
git commit -m "reject mysql list default partition syntax"
```

## Task 2: Add Catalog Session-State Tests

**Files:**
- Create or modify: `mysql/catalog/session_state_test.go`
- Modify: `mysql/catalog/catalog.go`
- Modify: `mysql/catalog/exec.go`

**Step 1: Write tests for `SET explicit_defaults_for_timestamp`**

Add focused catalog-only tests:

```go
func TestCatalogSessionStateExplicitDefaultsForTimestamp(t *testing.T) {
    c := New()
    if !c.session.ExplicitDefaultsForTimestamp {
        t.Fatal("default explicit_defaults_for_timestamp should be true")
    }
    mustExec(t, c, "SET SESSION explicit_defaults_for_timestamp=0")
    if c.session.ExplicitDefaultsForTimestamp {
        t.Fatal("expected explicit_defaults_for_timestamp=false")
    }
    mustExec(t, c, "SET @@session.explicit_defaults_for_timestamp=DEFAULT")
    if !c.session.ExplicitDefaultsForTimestamp {
        t.Fatal("DEFAULT should restore true")
    }
}
```

Also cover `ON`, `OFF`, `1`, `0`, and bare `SET explicit_defaults_for_timestamp=0`.

**Step 2: Run test to verify it fails**

Run:

```bash
go test ./mysql/catalog -run 'TestCatalogSessionStateExplicitDefaultsForTimestamp' -count=1
```

Expected: FAIL because no session field exists.

**Step 3: Introduce `SessionState`**

Add the struct in `catalog.go`, initialize it in `New()`, and either migrate existing top-level fields in one shot or keep compatibility accessors while moving logic to `c.session`.

```go
func defaultSessionState() SessionState {
    return SessionState{
        ForeignKeyChecks:             true,
        CharsetClient:                "utf8mb4",
        CollationConnection:          "utf8mb4_0900_ai_ci",
        ExplicitDefaultsForTimestamp: true,
        SQLMode:                      "DEFAULT",
    }
}
```

Normalize SET variable names by stripping `@@`, optional `session.`, `global.`, and by lowercasing. Rejecting unsupported scopes is not necessary for this task; keep existing permissive SET behavior.

**Step 4: Update `execSet`**

Handle:

```go
case "explicit_defaults_for_timestamp":
    if isDefaultSetValue(asgn.Value) {
        c.session.ExplicitDefaultsForTimestamp = true
    } else if v, ok := parseSetBool(asgn.Value); ok {
        c.session.ExplicitDefaultsForTimestamp = v
    }
case "sql_mode":
    c.session.SQLMode = nodeToSQLValue(asgn.Value)
```

Existing `foreign_key_checks`, GIPK, charset, and collation handling should keep passing.

**Step 5: Verify and commit**

Run:

```bash
go test ./mysql/catalog -run 'TestCatalogSessionState|TestWalkThrough_2_4|TestWalkThrough_13_1' -count=1
git diff --check
git add mysql/catalog/catalog.go mysql/catalog/exec.go mysql/catalog/session_state_test.go
git commit -m "model mysql catalog session state"
```

## Task 3: Implement Legacy TIMESTAMP Promotion

**Files:**
- Modify: `mysql/catalog/tablecmds.go`
- Modify: `mysql/catalog/scenarios_c3_test.go`
- Modify: `mysql/catalog/scenarios_c16_test.go`
- Modify: `mysql/catalog/SCENARIOS-mysql-implicit-behavior.md`

**Step 1: Tighten scenario tests**

Update `3.8` to run the session SET statements on both sides:

```go
setLegacyTimestampMode(t, mc)
mustExec(t, c, "SET SESSION explicit_defaults_for_timestamp=0; SET SESSION sql_mode='';")
```

Assert omni exactly matches the oracle:

- `ts1.Nullable == false`
- `ts1.Default` is `CURRENT_TIMESTAMP` or equivalent
- `ts1.OnUpdate` is `CURRENT_TIMESTAMP` or equivalent
- `ts2.Nullable == false`
- `ts2.Default` is zero timestamp
- `ts2.OnUpdate == ""`

Unskip `16.12` and assert `TIMESTAMP(3) NOT NULL` under legacy mode gets `CURRENT_TIMESTAMP(3)` for both default and on-update.

**Step 2: Run tests to verify they fail**

Run:

```bash
go test ./mysql/catalog -run 'TestScenario_C3/3_8|TestScenario_C16/16_12' -count=1 -timeout=20m
```

Expected: FAIL until session-driven promotion is implemented.

**Step 3: Implement timestamp promotion helper**

After all columns have been created, call a helper before constraints and partition validation rely on final column metadata.

Rules for this task:

- Only applies when `!c.session.ExplicitDefaultsForTimestamp`.
- Only applies to `DataType == "timestamp"`.
- Explicit `NULL` keeps nullable behavior and blocks legacy NOT NULL promotion for that column.
- Explicit `DEFAULT` must not be overwritten.
- Explicit `ON UPDATE` must not be overwritten.
- First eligible timestamp gets `DEFAULT CURRENT_TIMESTAMP[(fsp)]` and `ON UPDATE CURRENT_TIMESTAMP[(fsp)]`.
- Later eligible timestamp columns get zero timestamp default when they are NOT NULL and no explicit default exists.

Use `colDef.TypeName.Length` as fsp; default `0`.

**Step 4: Update tracker and docs**

Mark `3.8` and `16.12` verified. Update C3/C16 bug queues if they still mention active session-state limitations.

**Step 5: Verify and commit**

Run:

```bash
go test ./mysql/catalog -run 'TestScenario_C3/3_8|TestScenario_C16/16_12' -count=1 -timeout=20m
go test ./mysql/catalog -run 'TestScenario_C(3|16)$' -count=1 -timeout=30m
git diff --check
git add mysql/catalog/tablecmds.go mysql/catalog/scenarios_c3_test.go mysql/catalog/scenarios_c16_test.go mysql/catalog/SCENARIOS-mysql-implicit-behavior.md mysql/catalog/scenarios_bug_queue/c3.md mysql/catalog/scenarios_bug_queue/c16.md
git commit -m "honor legacy timestamp session defaults"
```

## Task 4: Final Reconciliation

**Files:**
- Modify: `mysql/catalog/SCENARIOS-mysql-implicit-behavior.md`
- Modify as needed: `mysql/catalog/scenarios_bug_queue/*.md`

**Step 1: Run all scenario tests**

Run:

```bash
go test ./mysql/catalog -run 'TestScenario_(C|AX|PS)' -count=1 -timeout=40m
```

Expected: PASS.

**Step 2: Recount tracker**

Run:

```bash
awk -F'|' '/^\| [A-Z0-9]+\.[0-9]+ \|/ {gsub(/^ +| +$/,"",$2); gsub(/^ +| +$/,"",$4); gsub(/^ +| +$/,"",$5); status[$4]++; if($4!="verified"){remain++; print $2 " | " $4 " | " $5}; total++} END {print "total", total; for (s in status) print s, status[s]; print "remaining", remain}' mysql/catalog/SCENARIOS-mysql-implicit-behavior.md
```

Expected: `remaining 0`.

**Step 3: Commit final reconciliation**

Run:

```bash
git diff --check
git add mysql/catalog/SCENARIOS-mysql-implicit-behavior.md mysql/catalog/scenarios_bug_queue
git commit -m "complete mysql implicit scenario reconciliation"
```

## Open Decisions

`SQLMode` is stored but not interpreted beyond allowing timestamp tests to mirror the oracle setup. Do not expand into full SQL mode behavior in this task; that is a separate compatibility surface.

If future MySQL versions add a real LIST DEFAULT catch-all partition, `Task 1` will fail on the oracle side. At that point, treat it as a new version-gated scenario, not as current MySQL 8.0 behavior.
