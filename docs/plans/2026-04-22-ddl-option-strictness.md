# DDL Option-Value Strictness — Plan

Date: 2026-04-22
Scope: close the remaining `isAnyKeywordIdent` over-acceptance in omni's mssql parser
Depends on: `docs/plans/2026-04-22-ddl-option-taxonomy.md` (the 426-site catalog)

## Why this exists

After closing the identifier-vs-reserved axis (commit `3b86999`), `isAnyKeywordIdent` still has **~420 remaining call sites** that the earlier fix deliberately left alone. The taxonomy pass shows these fall into **six categories**, of which only three are actually strictness gaps.

**Verified live bug example**:
```
ALTER DATABASE foo SET RECOVERY SELECT
  SqlScriptDOM → REJECT
  omni         → accept    ❌
```

`mssql-strict` (the existing starmap) has already closed the **option NAME** axis (17 categories passing `TestKeywordOracleOptionPositions`). What remains is the **option VALUE** axis plus a few unrelated chores.

## Category taxonomy (from the 426-site audit)

| Cat | What it is | Count | Action |
|---|---|---|---|
| **A** | Statement / object names | ~95 | Replace with `isIdentLike` (low risk) |
| **B** | Option **VALUE** enums (right side of `=`) | ~95 | Per-option `optionSet` — **main remaining work** |
| **C** | Schema / filegroup / column identifier references | ~150 | Replace with `isIdentLike` (medium risk, needs per-case check) |
| **D** | Subcommand dispatch (`ALTER DATABASE name ADD/MODIFY/REMOVE`) | ~10 | Dispatch `optionSet` per statement |
| **E** | Dynamic property bags (service broker message props, XML attrs) | ~20 | **Leave alone** — grammar genuinely accepts any keyword |
| **F** | Identifier-like misses (similar to A, slightly different context) | ~40 | Same treatment as A |

**Effective untouched work = A (95) + B (95) + C (150) + D (10) + F (40) ≈ 390 sites.** E is preserved as-is (20).

## What `mssql-strict` already did vs. didn't

### Did (NAME side — 17 section headers)
`SCENARIOS-mssql-strict.md` closed the option-NAME whitelist at every `WITH (NAME = value)` entry:
- SET predicate (25 names), ALTER DATABASE SET, index WITH, table WITH, proc WITH, view WITH, query hints, FOR XML/JSON, cursor, backup/restore, bulk insert, fulltext, service broker NAMES, availability group, endpoint protocol
- All verified via `TestKeywordOracleOptionPositions` → 0 mismatch

### Didn't (what's still in `isAnyKeywordIdent` territory)

1. **Option VALUES** — `SET RECOVERY {FULL|SIMPLE|BULK_LOGGED}`, `DATA_COMPRESSION = {NONE|ROW|PAGE|COLUMNSTORE}`, `AUTOMATED_BACKUP_PREFERENCE = {PRIMARY|SECONDARY_ONLY|...}`. Category B.
2. **Statement-level object names** — `CREATE SERVICE svc`, `ALTER ROLE role`, `DROP ASSEMBLY a` etc. Category A.
3. **Subcommand dispatch** — `ALTER DATABASE {ADD|MODIFY|REMOVE}`, `ALTER AVAILABILITY GROUP {MODIFY|FORCE_FAILOVER_ALLOW_DATA_LOSS|...}`. Category D.

## Clustering by file (top 10, 300/426 sites covered)

| File | Sites | Cat split | Priority | SqlScriptDOM source |
|---|---|---|---|---|
| service_broker.go | 46 | 16 A / 15 B / 3 D / 12 E | High | ServiceBrokerOptionsHelper; ~12 MUST NOT change |
| fulltext.go | 38 | 8 A / 10 B / 20 C | High | FulltextIndexOptionHelper etc. |
| external.go | 32 | 14 A / 6 B / 12 C | Medium | ExternalDataSourceOptionsHelper, ExternalFileFormatOptionsHelper |
| endpoint.go | 32 | 3 A / 0 B / 29 C | Low | (covered by strict Phase 6.4 for names; C sites are filegroup/cert references) |
| utility.go | 28 | 4 A / 12 B / 12 C | Medium | scattered — DBCC / SET options |
| security_audit.go | 22 | — | Medium | AuditOptionKindHelper, AuditActionSpecOptionsHelper |
| create_table.go | 21 | mostly C + A | Low | already covered by Phase 3.2 for NAMES |
| resource_governor.go | 20 | mostly B | Medium | ResourcePoolOptionsHelper, WorkloadGroupOptionsHelper |
| alter_objects.go | 18 | mixed | Medium | partially covered by Phase 2.2 |
| event.go | 17 | mostly B+D | Medium | EventSessionOptionsHelper |

Lower-volume files (< 10 each) total ~120 sites across 30 files.

## Proposed phase breakdown

### Phase X.A — Statement name capture (~95 A sites)

**Mechanical**: replace `isAnyKeywordIdent` with `isIdentLike` at every "stmt.Name = p.cur.Str" capture.

Example diff shape:
```go
-	if p.isAnyKeywordIdent() {
+	if p.isIdentLike() {
 		stmt.Name = p.cur.Str
 		p.advance()
 	}
```

- Risk: inputs like `CREATE SERVICE FROM` were previously accepted, will now reject. **SqlScriptDOM rejects them too**, so this is pure alignment.
- Rollout: one commit per file to keep diffs reviewable.
- Verification: each file gets 1-2 reject-alignment fixtures in `TestScriptDOMRejectAlignment`.

### Phase X.B — Option VALUE enums (~95 B sites)

**Requires work per site**: find the SqlScriptDOM helper class and extract the enum.

Pattern:
```go
// Before
if p.isAnyKeywordIdent() {
    val := strings.ToUpper(p.cur.Str)
    p.advance()
    switch val {
    case "FULL", "SIMPLE", "BULK_LOGGED": /* ok */
    default: /* silently kept */
    }
}

// After
var recoveryValues = newOptionSet().withIdents("FULL", "SIMPLE", "BULK_LOGGED")
val, err := p.expectOption(recoveryValues)
if err != nil { return err }
```

- Risk: each option needs a verified enum list. Harness oracle scripted per enum:
  ```bash
  for v in FULL SIMPLE BULK_LOGGED SELECT FROM BOGUS; do
      echo "ALTER DATABASE foo SET RECOVERY $v" | harness...
  done
  ```
- Rollout: one commit per SqlScriptDOM helper class (~10-15 commits).
- Verification: 3-5 reject fixtures per enum (bad enum name + reserved keyword + bogus ident).

### Phase X.C — Schema/filegroup/column references (~150 C sites)

**Usually an identifier position**: replace `isAnyKeywordIdent` with `isIdentLike`.

Example: filegroup names in `CREATE TABLE ... ON [PRIMARY]`. `PRIMARY` is CoreKeyword; must be bracketed per SqlScriptDOM.

- Risk: higher — some sites may legitimately accept keywords (e.g. `ON DEFAULT`).
- Rollout: per-file, with harness validation of DEFAULT / PRIMARY / unknown keyword acceptance.

### Phase X.D — Subcommand dispatch (~10 D sites)

**Small, targeted**: at each `ALTER X` site, define the valid subcommand set.

```go
var alterDatabaseSubcommands = newOptionSet().withKeywords(kwADD, kwMODIFY, kwREMOVE, kwSET, kwCOLLATE, ...)
```

- Risk: low, but need to ensure the dispatch lookup still works (optionSet returns matched token, caller switches on it).
- Rollout: one commit per statement.

### Phase X.E — DO NOT TOUCH (~20 E sites)

Dynamic property-bag sites:
- service broker message type properties (`SENT BY { INITIATOR | TARGET | ANY }` — ANY is a keyword, accepted intentionally)
- XML attribute name storage
- RESTORE FILELISTONLY metadata

Mark each with a comment explaining why it accepts arbitrary keywords.

### Phase X.F — Identifier-like misses (~40 F sites)

Similar to A, but in less common positions (e.g. role names, schema names in nested DDL). Same mechanical `isAnyKeywordIdent → isIdentLike` replacement.

## Implementation order recommendation

1. **Phase X.A first (95 A sites)** — lowest risk, highest clarity. Mostly mechanical, likely 1 day of work.
2. **Phase X.F alongside A** — same treatment, merge into same batches (total ~135 sites).
3. **Phase X.D next (10 D sites)** — small, discrete, catches several user-visible bugs.
4. **Phase X.B by SqlScriptDOM helper class** — the real work. Break into ~10-15 commits, one per helper. Each needs harness fixtures.
5. **Phase X.C last** — review per-site; some may stay as-is.
6. **Phase X.E** — defensive comments only, no code change.

## Estimated effort

- Phase X.A + X.F (mechanical ident replacement): ~0.5 day; ~15 commits
- Phase X.D (subcommand dispatch): ~0.5 day; ~10 commits
- Phase X.B (option VALUE enums): 2-4 days; ~15 commits with per-enum harness fixtures
- Phase X.C (schema/filegroup refs): 1 day; ~10 commits
- Phase X.E: 0.25 day (pure doc comments)

**Total**: ~5 working days if done linearly. Parallelizable into 3 independent tracks (A+F, B, C+D).

## What to verify after each phase

1. `go build ./mssql/...`
2. `go test ./mssql/... -count=1 -short` — full mssql suite
3. `go test -tags scriptdom ./mssql/parser/ -run TestScriptDOM` — harness
4. `TestKeywordOracleOptionPositions` — existing oracle strictness fence (must stay 0 mismatch)
5. New reject fixtures added for the phase-specific enum

## Starmap integration

This plan is a natural extension of `SCENARIOS-mssql-strict.md`. Recommend:

1. Add a **Phase 7 block** to `SCENARIOS-mssql-strict.md` titled "Option VALUE strictness" with sections 7.1 through 7.N matching the SqlScriptDOM helper classes enumerated above.
2. Reuse `mssql-strict-worker` / `mssql-strict-driver` skills — the scenario shape matches (one helper class per section).
3. Track in `PROGRESS.json` alongside the existing option-NAME work.

## POC candidate for this worktree

Before handing off to the mssql-strict starmap, pick **one** well-bounded cluster as a proof-of-concept here. Best candidate:

### `SET RECOVERY` + a handful of other simple ALTER DATABASE SET value enums

- **Size**: small (~6 value-enum positions in `alter_objects.go`)
- **Visibility**: the bug example (`SET RECOVERY SELECT`) is the one I already showed the user
- **SqlScriptDOM source**: `DatabaseOptionKindHelper.cs` lists the enum values per option
- **Verification**: harness + oracle test already set up

If POC succeeds, we have empirical evidence that the `optionSet` approach works at this scale and the estimated effort is realistic. Then confidently hand the rest to mssql-strict.

## Out of scope

- Data type keywords (INT, VARCHAR, etc.) — already handled by lexer/datatype parser, orthogonal axis.
- String literal option values (file paths, connection strings) — already correct, use `tokSCONST`.
- User-defined option names (in `OPENROWSET(BULK ...)` bulk options) — semantically strings, leave alone.
