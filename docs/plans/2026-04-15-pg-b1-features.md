# pg/parser B1 feature cleanup plan

> Status: pre-implementation. Tackles the "high-value + low-cost" subset
> of the `known_failures.json` residue after the follow-up cleanup PR
> (branch `pg-first-sets`).

## Background

The `known_failures.json` audit after PR #71 landed showed 350 remaining
failures in the PG 17 regression corpus. Categorization grouped them
into three buckets:

- **A. Not bugs** (~42): `errors.sql` negative tests + psql meta-commands
- **B1. High-value + low-cost feature gaps** (this plan, ~48 items)
- **B2. Medium-value**
- **B3. High-cost or niche** (JSON / SQL:2016)

B1 targets 8 concrete feature gaps in 7 PG 15/16-era additions plus one
long-standing omission. Each gap is a handful of lines of grammar plus a
test, so the whole batch fits in one PR.

---

## Section 1 — Findings (per feature)

All 8 verified to fail at `HEAD` of `pg-first-sets` (branch pointer
`19d0389`). The PG16 numeric literal gap (`0b101`, `0x1f`, `0o17`,
`1_000_000`) was initially listed in the B1 bucket but verification
showed omni's lexer **already fully supports** them — the 10
`numerology.sql` entries are negative tests (`SELECT 0b` / `0x` / `0o`
with no digits, `SELECT 123abc` trailing-junk) that omni correctly
rejects. They are not bugs and will not be addressed in this plan.

### Finding 1 — GRANT role: comma-separated WITH options + GRANTED BY (25 items)

**Failing SQL**:
- `GRANT r1 TO r2 WITH INHERIT FALSE, ADMIN TRUE`
- `GRANT r1 TO r2 WITH INHERIT TRUE, SET FALSE`
- `GRANT r1 TO r2 GRANTED BY r3`
- `GRANT r1 TO r2 WITH ADMIN OPTION GRANTED BY r3`
- `GRANT r1 TO r2 WITH ADMIN TRUE GRANTED BY r3`

**Root causes (2 separate bugs)**:

**Bug 1a**: `finishGrantRole` at `grant.go:160-170` does not call
`parseOptGrantedBy()`. The helper exists (`grant.go:534-542`) and is
used by `finishRevokeRole`, but nobody wired it into the GRANT path.
Fix: add `grantedBy := p.parseOptGrantedBy()` after the opt-list parse
and set `Grantor: grantedBy` on the returned node.

**Bug 1b**: `parseGrantRoleOptList` at `grant.go:544-557` loops calling
`parseGrantRoleOpt`, but `parseGrantRoleOpt` starts with:

```go
if p.cur.Type != WITH { return nil }
p.advance()
return p.parseGrantRoleOptValue()
```

So the WITH token is consumed at the start, and the loop's second
iteration sees a `,` or end-of-list and returns nil immediately.
Only the first option is ever parsed.

PG's grammar is `grant_role_opt_list: grant_role_opt_list ',' grant_role_opt | grant_role_opt`
— options are comma-separated AFTER a single WITH. Fix: restructure so
WITH is consumed once, then a comma-separated list is parsed:

```go
func (p *Parser) parseGrantRoleOptList() *nodes.List {
    if p.cur.Type != WITH {
        return nil
    }
    p.advance() // consume WITH
    items := []nodes.Node{}
    items = append(items, p.parseGrantRoleOptEntry())
    for p.cur.Type == ',' {
        p.advance()
        items = append(items, p.parseGrantRoleOptEntry())
    }
    return &nodes.List{Items: items}
}
```

where `parseGrantRoleOptEntry` is the renamed `parseGrantRoleOptValue`.

**Regression-sanity tests**: `GRANT r1 TO r2 WITH ADMIN OPTION` and
`GRANT r1 TO r2` (no WITH clause) continue to parse.

### Finding 2 — REVOKE role: ColId OPTION FOR (3+ items)

**Failing SQL**:
- `REVOKE INHERIT OPTION FOR r1 FROM r2`
- `REVOKE SET OPTION FOR r1 FROM r2`
- (working today: `REVOKE ADMIN OPTION FOR r1 FROM r2`)

**Root cause**: `parseRevokeStmt` at `grant.go:53-117` only handles
`GRANT OPTION FOR` and `ADMIN OPTION FOR`. PG's grammar allows
**any ColId**: `REVOKE ColId OPTION FOR privilege_list FROM role_list`.
The ColId becomes a DefElem name with value `false` (the semantic
meaning is "revoke the X-option flag from this grant").

**Fix**: replace the `ADMIN`-specific branch with a lookahead check
that detects `<ColId> OPTION FOR` by peeking two tokens, and generalize
the captured DefElem name. Pseudo-code:

```go
// After the GRANT OPTION FOR branch:
if p.isColId() && p.peekNext().Type == OPTION {
    optName := p.cur.Str
    p.advance() // ColId
    p.expect(OPTION)
    p.expect(FOR)
    roleOptionName = strings.ToLower(optName)
}
```

Then pass `roleOptionName` into `finishRevokeRole` and use it to build
the `DefElem` (replacing the current hardcoded `"admin"`).

**Regression-sanity tests**: `REVOKE ADMIN OPTION FOR r1 FROM r2` and
`REVOKE GRANT OPTION FOR SELECT ON t FROM r` continue to parse.

### Finding 3 — UNIQUE NULLS [NOT] DISTINCT (7 items)

**Failing SQL**:
- `CREATE TABLE t (i int UNIQUE NULLS NOT DISTINCT, t text)`
- `CREATE UNIQUE INDEX i ON t (c) NULLS DISTINCT`
- `CREATE UNIQUE INDEX i ON t (c) NULLS NOT DISTINCT`

**Root cause** (corrected from plan v1 after verifying PG source):
`parseOptUniqueNullTreatment` at `create_table.go:1308-1320` checks for
`NULLS_LA`, but PG's grammar at `gram.y:4019-4021` uses plain `NULLS_P`
in this production:

```
opt_unique_null_treatment:
    NULLS_P DISTINCT        { $$ = true; }
  | NULLS_P NOT DISTINCT    { $$ = false; }
  | /*EMPTY*/               { $$ = true; }
```

PG's `base_yylex` in `parser.c:221-229` only reclassifies `NULLS_P` to
`NULLS_LA` when followed by `FIRST_P`/`LAST_P` — matching omni's
`parser.go:1126-1132`. `NULLS DISTINCT` / `NULLS NOT DISTINCT` keeps
the plain `NULLS_P` token in PG. omni's reclassification is correct;
the probe site is wrong.

**Fix**: one-line change in `parseOptUniqueNullTreatment` — check
`NULLS_P` instead of `NULLS_LA`. No reclassification changes needed.

```go
func (p *Parser) parseOptUniqueNullTreatment() bool {
    if p.cur.Type == NULLS_P { // was: NULLS_LA
        p.advance()
        if p.cur.Type == NOT {
            p.advance()
            p.expect(DISTINCT)
            return true
        }
        p.expect(DISTINCT)
        return false
    }
    return false
}
```

Also update the doc comment above the function — it references
`NULLS_LA` in the grammar shorthand, which is misleading.

**Regression-sanity tests**: `ORDER BY x NULLS FIRST` and
`ORDER BY x NULLS LAST` continue to parse (exercise the unchanged
reclassification path).

### Finding 4 — ALTER POLICY ... RENAME (3 items)

**Failing SQL**:
- `ALTER POLICY p1 ON t RENAME TO p2`

**Root cause**: `parseAlterPolicyStmt` at `grant.go:972` doesn't
handle the RENAME variant. PG's grammar:

```
AlterPolicyStmt: ALTER POLICY name ON qualified_name alter_policy_options
              | ALTER POLICY name ON qualified_name RENAME TO name
```

**Fix**: after consuming `ALTER POLICY name ON qualified_name`, peek
for `RENAME` and branch to a `RenameStmt`-returning path. Keep the
existing alter-options path for the non-rename case.

**Regression-sanity tests**: `ALTER POLICY p1 ON t USING (a > 0)` and
`ALTER POLICY p1 ON t WITH CHECK (a > 0)` continue to parse.

### Finding 5 — ALTER VIEW RENAME COLUMN (3 items)

**Failing SQL**:
- `ALTER VIEW v RENAME COLUMN a TO b`

**Root cause**: omni handles `ALTER VIEW ... RENAME TO` (rename the
view itself) but not `ALTER VIEW ... RENAME COLUMN x TO y`. PG's
grammar allows both under `AlterObjectSchemaStmt`/`RenameStmt` shapes.

**Fix**: in the ALTER VIEW dispatch (find via `grep parseAlterView`),
after RENAME peek for COLUMN; if present, parse `COLUMN name TO name`
and return a `RenameStmt` with the right `RelationType` / `SubName`.

**Regression-sanity tests**: `ALTER VIEW v RENAME TO v2` and
`ALTER VIEW v ALTER COLUMN a SET DEFAULT 1` continue to parse.

### Finding 6 — ALTER DATABASE ... REFRESH COLLATION VERSION (3 items)

**Failing SQL**:
- `ALTER DATABASE postgres REFRESH COLLATION VERSION`

**Root cause**: `parseAlterDatabaseDispatch` at `database.go:140` has
branches for SET/RESET/OWNER TO/RENAME TO/WITH option list but not
REFRESH. Added in PG 15. PG's grammar:

```
AlterDatabaseRefreshCollStmt:
    ALTER DATABASE name REFRESH COLLATION VERSION
```

**Fix**: add a `REFRESH` branch in the dispatch that expects
`COLLATION VERSION` and returns the corresponding AST node.
`T_AlterDatabaseRefreshCollStmt` already exists in `pg/ast/nodetags.go`
but the actual struct is missing in `pg/ast/parsenodes.go` — add it
(single `{DbName string; Loc nodes.Loc}` shape, following the nearby
`AlterDatabaseSetStmt` / `AlterDatabaseStmt` pattern). Walker
regeneration should be mechanical.

**Regression-sanity tests**: `ALTER DATABASE d SET timezone = 'UTC'`
and `ALTER DATABASE d RENAME TO d2` continue to parse.

### Finding 7 — CREATE PUBLICATION FOR TABLES IN SCHEMA ... WHERE (3 items)

**Failing SQL**:
- `CREATE PUBLICATION p FOR TABLES IN SCHEMA s1 WHERE (a = 123)`
- `CREATE PUBLICATION p FOR TABLE t, TABLES IN SCHEMA s1 WHERE (a > 10)`

**Root cause**: omni's `parsePublicationObjSpec` (referenced in `publication.go:185`)
handles both `FOR TABLE t` and `FOR TABLES IN SCHEMA s`, but the row
filter (`WHERE (...)` clause on table specs) is rejected when combined
with `TABLES IN SCHEMA` in a mixed list. Need to inspect the function
to confirm the exact failure mode.

**Fix strategy**: read `parsePublicationObjSpec`, find where it rejects
the combined form, and extend to accept WHERE after `TABLES IN SCHEMA`
per PG's grammar:

```
PublicationObjSpec:
    TABLE relation_expr opt_column_list OptWhereClause
  | TABLES IN_P SCHEMA ColId OptWhereClause
  | ...
```

Note: PG actually rejects `WHERE` on `TABLES IN SCHEMA` at semantic time
— it's a schema-level spec, not a row-level one, so the WHERE clause
doesn't apply. But it's allowed at parse time and errored at
post-parse validation. omni should match the parse-time permissiveness.

**Regression-sanity tests**: `CREATE PUBLICATION p FOR TABLE t` and
`CREATE PUBLICATION p FOR TABLE t WHERE (a > 10)` continue to parse.

### Finding 8 — ALTER TABLE ... ALTER COLUMN ... SET COMPRESSION DEFAULT (1 item)

**Failing SQL**:
- `ALTER TABLE t ALTER COLUMN c SET COMPRESSION default` (compression.sql:63)

**Root cause**: `alter_table.go:1012-1020` parses SET COMPRESSION as
`SET COMPRESSION ColId` only. `DEFAULT` is a reserved keyword in omni
and is therefore not accepted by `parseColId`. The adjacent
`SET STORAGE` branch at lines 1000-1006 already handles this by
special-casing `DEFAULT` before falling through to `parseColId`:

```go
var storageVal string
if p.cur.Type == DEFAULT {
    p.advance()
    storageVal = "default"
} else {
    storageVal, _ = p.parseColId()
}
```

PG's grammar treats both symmetrically — `SET COMPRESSION ColId`
where `default` is admitted because it's an unreserved keyword in
PG's `ColId`. omni has `DEFAULT` in the reserved list, so the
DEFAULT branch has to be inlined (matching the SET STORAGE pattern).

**Fix**: mirror the SET STORAGE DEFAULT handling at the SET
COMPRESSION site. 3-line change.

```go
case COMPRESSION:
    p.advance()
    var compVal string
    if p.cur.Type == DEFAULT {
        p.advance()
        compVal = "default"
    } else {
        compVal, _ = p.parseColId()
    }
    return &nodes.AlterTableCmd{
        Subtype: int(nodes.AT_SetCompression),
        Name:    colname,
        Def:     &nodes.String{Str: compVal},
    }, nil
```

**Regression-sanity tests**: `SET COMPRESSION pglz` and
`SET COMPRESSION lz4` continue to parse (ColId path).

**Note (de-scoped from plan v1)**: The original plan v1 listed
`CAST('42' AS text COLLATE "C")` as Finding 8. On investigation, the
`strings.sql` known_failures entries at lines 1, 8-10, 33-34, 78-91
are NOT CAST COLLATE — they're unicode escape diagnostics, `bytea`
escape-output round-trips, and line-continuation string literals.
No B1-scope CAST COLLATE gap exists in the current corpus. The one
real COLLATE-in-type-literal form (`CAST(x AS text COLLATE "C")`) is
only reachable via `a_expr COLLATE any_name` applied to the CAST
result, which omni already parses. Dropped from this plan.

### Summary — 8 Features

| # | Feature | Count | Fix complexity |
|---|---|---|---|
| 1 | GRANT role WITH opts comma-sep + GRANTED BY | ~20+ | medium (`grant.go` 2 bugs) |
| 2 | REVOKE role ColId OPTION FOR | ~5 | small (`grant.go` generalize) |
| 3 | UNIQUE NULLS [NOT] DISTINCT | 7 | tiny (1 line in `create_table.go`) |
| 4 | ALTER POLICY ... RENAME | 3 | small (dispatch branch) |
| 5 | ALTER VIEW RENAME COLUMN | 3 | small (dispatch branch) |
| 6 | ALTER DATABASE REFRESH COLLATION VERSION | 3 | small (dispatch branch, maybe + AST node) |
| 7 | CREATE PUBLICATION FOR TABLES IN SCHEMA WHERE | 3 | small-medium (`publication.go` WHERE extension) |
| 8 | ALTER TABLE SET COMPRESSION DEFAULT | 1 | tiny (3 lines in `alter_table.go`) |

**Estimated real count**: ~45 known_failures entries covered (minus any
that turn out to be from a different root cause on closer inspection).
The 4-entry CAST COLLATE candidate from plan v1 was dropped after
verification (see Finding 8 note).

---

## Section 2 — Detection method

For each finding:

1. **Grep** omni's parser for the relevant function and PG's `gram.y`
   for the corresponding production.
2. **Read** the surrounding code to understand the current shape and
   identify the specific gap.
3. **Empirically test** with a minimal Go program that calls
   `parser.Parse` on a representative SQL.
4. **Cross-check** against pgregress `known_failures.json` entry list
   — confirm the expected fix would remove the entry.

For findings 1-7 this was straightforward. Finding 8 (SET COMPRESSION
DEFAULT) was surfaced by codex review of plan v1 — compression.sql:63
sits alone in the known_failures list and a quick grep at
`alter_table.go:1012` confirmed the missing DEFAULT branch.

**Regression-sanity test requirement** (from plan v2 review): every
finding's test commit must include at least one subtest covering an
existing, working form that shares the same code path, so the fix
doesn't silently narrow previously-accepted grammar. Per-finding
regression subtests are listed in each Finding section above.

---

## Section 3 — Plan

### Commit breakdown

**Commit 1**: Fix GRANT role statement — add `parseOptGrantedBy()` call
to `finishGrantRole`, restructure `parseGrantRoleOptList` to handle
comma-separated options after a single WITH. Add
`TestGrantRoleOptionsAndGrantedBy` covering 4-6 forms.

**Commit 2**: Generalize REVOKE role `ADMIN OPTION FOR` to
`<ColId> OPTION FOR`. Add subtests for ADMIN/INHERIT/SET option
variants.

**Commit 3**: Change `parseOptUniqueNullTreatment` to probe `NULLS_P`
instead of `NULLS_LA` (match PG's `opt_unique_null_treatment`
production which uses plain `NULLS_P`). Update the doc comment.
Add `TestUniqueNullsNotDistinct` covering CREATE TABLE column
constraint, table-level constraint, and CREATE UNIQUE INDEX.

**Commit 4**: Add `ALTER POLICY ... RENAME TO` to `parseAlterPolicyStmt`.
Add `TestAlterPolicyRename`.

**Commit 5**: Add `ALTER VIEW ... RENAME COLUMN` to the ALTER VIEW
dispatch. Add `TestAlterViewRenameColumn`.

**Commit 6**: Add `ALTER DATABASE ... REFRESH COLLATION VERSION` branch.
Verify `nodes.AlterDatabaseRefreshCollStmt` exists; if not, add it in
the same commit. Add `TestAlterDatabaseRefreshCollation`.

**Commit 7**: Extend `parsePublicationObjSpec` to accept WHERE on
`TABLES IN SCHEMA` entries. Add
`TestCreatePublicationTablesInSchemaWhere`.

**Commit 8**: Add `DEFAULT` handling to the SET COMPRESSION branch in
`alter_table.go` mirroring SET STORAGE. Add
`TestAlterTableSetCompressionDefault`.

**Commit 9**: Run `go test -short ./pg/pgregress/...`; note which
entries in `known_failures.json` are now fixed; remove them.

**Commit 10**: Final codex review on the impl diff.

### Ordering

Commits 3, 4, 5, 6, 7, 8 are independent of each other and of
commits 1/2. Commits 1 and 2 both touch `grant.go` and should be
sequential. Commit 6 is the only risky one — if walker regeneration
for `AlterDatabaseRefreshCollStmt` touches unrelated files, it can be
split into its own PR without blocking the others.

### Test matrix

For each finding, the target test file is:

| Finding | Test file | Subtests |
|---|---|---|
| 1 | `grant_role_opt_list_test.go` (new) | 6+ |
| 2 | `grant_role_opt_list_test.go` (same) | 4 |
| 3 | `unique_nulls_not_distinct_test.go` (new) | 6 |
| 4 | `alter_policy_rename_test.go` (new) | 2 |
| 5 | `alter_view_rename_column_test.go` (new) | 2 |
| 6 | `alter_database_refresh_test.go` (new) | 2 |
| 7 | `create_publication_row_filter_test.go` (new) | 3 |
| 8 | `alter_table_set_compression_test.go` (new) | 3 |

All tests should use AST-shape assertions (not just error-absence),
following the pattern established by `create_function_json_test.go`
(commit `549a1c0`) and the follow-up qualified-type tests.

### known_failures.json cleanup

After all fixes land, re-run `go test -short ./pg/pgregress/...`. Each
"FIXED" report corresponds to a removable entry. Expect ~45 entries
removed, bringing the total from 350 → ~305.

### Risks

1. **AST struct for `AlterDatabaseRefreshCollStmt`** is missing in
   `pg/ast/parsenodes.go` even though `T_AlterDatabaseRefreshCollStmt`
   is present in `pg/ast/nodetags.go`. Adding the struct is
   mechanical; if walker regeneration (`pg/ast/cmd/genwalker`) touches
   unrelated files, split commit 6 into its own PR.

2. **Generalized REVOKE role ColId OPTION FOR** could accidentally
   accept bogus option names like `REVOKE BOGUS OPTION FOR r FROM r`.
   PG accepts this at parse time and errors at semantic time. omni
   should match. Mitigation: no validation at parse level; the test
   covers the known-good names (ADMIN/INHERIT/SET) and rejects
   obviously-bad ones like `REVOKE 123 OPTION FOR ...` via ColId
   rejection at `parseColId`.

3. **PublicationObjSpec WHERE interaction**: PG accepts `WHERE` on
   `TABLES IN SCHEMA` at parse time but errors at validation time.
   If omni's current behavior is semantic rejection rather than parse
   rejection, the "fix" may actually just be removing a parse-time
   check rather than adding grammar. Investigate on implementation.

### Time estimate

| Phase | Time |
|---|---|
| Commit 1+2 (GRANT/REVOKE role) | 60-90 min |
| Commit 3 (NULLS_P probe) | 15 min |
| Commit 4 (ALTER POLICY RENAME) | 30 min |
| Commit 5 (ALTER VIEW RENAME COLUMN) | 30 min |
| Commit 6 (ALTER DATABASE REFRESH) | 30-60 min (AST struct risk) |
| Commit 7 (PUBLICATION row filter) | 30-60 min |
| Commit 8 (SET COMPRESSION DEFAULT) | 15 min |
| Commit 9 (known_failures cleanup) | 15 min |
| Final codex review + iteration | 30-60 min |
| **Total** | **3.5-6 hours** |

### Deferred / out of scope

- **B2 medium-value features**: 60 items across 8 categories
  (array-of-composite access, multi-SET identity chain, foreign table
  column options, GRANT ALL column-level, SUBSTRING SIMILAR, etc.)
- **B3 JSON / SQL:2016**: 80+ items across advanced JSON functions and
  SQL/JSON standard productions
- **PG16 numeric literal negative tests**: 10 `numerology.sql` entries
  that omni already correctly rejects. These should ideally move to a
  separate `intentional_failures.json` or be annotated in the pgregress
  harness as expected-to-fail. Tracked separately.
- **`errors.sql` negative tests**: 31 entries, same situation as above.

---

## Review questions for codex

1. **Finding 1 (GRANT role) root-cause analysis**: is the comma-separated
   loop bug the only issue, or is `parseGrantRoleOpt` also incorrectly
   scoped? Specifically: does it need `parseGrantRoleOptList` consuming
   the WITH, or is there a better shape?
2. **Finding 2 (REVOKE role)**: does PG restrict the ColId to a specific
   set, or accept any ColId at parse time? If any ColId, is there a
   concern about ambiguity with statements like
   `REVOKE role FROM role OPTION FOR ...` where OPTION is a field name?
3. **Finding 3 (NULLS_P)** — RESOLVED in v2. Verified against PG's
   `gram.y:4019-4021` (`opt_unique_null_treatment: NULLS_P DISTINCT`)
   and `parser.c:221-229` (`_LA` only fires on FIRST/LAST). Fix is a
   1-line probe-site change — no reclassification extension needed.
4. **Finding 6 (AlterDatabaseRefreshCollStmt)**: struct is missing in
   `pg/ast/parsenodes.go` but tag exists in `nodetags.go`. Confirm
   walker regeneration is mechanical; if not, split to its own PR.
5. **Finding 8 (CAST COLLATE)** — DROPPED in v2. strings.sql
   known_failures are unicode-escape / bytea-output / line-continuation
   issues, not CAST COLLATE. Replaced with SET COMPRESSION DEFAULT.
6. **Anything I'm missing**: are there other B1-scope items in
   `known_failures.json` I haven't enumerated?
