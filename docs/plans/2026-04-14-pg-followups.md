# pg/parser follow-up cleanup plan

> Status: pre-implementation. Findings + plan for review BEFORE any code
> changes. Targets the three follow-up items left over after the
> backtrack-fix PR (commits dbb6471..efb0ff1 on `pg-first-sets`).

## Background

After completing the FIRST-set consolidation refactor (`pg-first-sets`
branch, commits `0a390de..efb0ff1`) and its codex review, three
follow-up items were documented but not fixed:

1. **3-component qualified type names** (`db.schema.mytype`) — fail in all
   type contexts. Documented in `pg/parser/CLAUDE.md` "Known limitations".
2. **3 `isColId() || isTypeFunctionName()` call sites** — flagged in commit
   `fc43201` as Phase 3 audit candidates because the `isAExprStart` oracle
   showed TypeFuncNameKeyword tokens are not real expression starters.
3. **Ryuk reaper flake** — testcontainers-go race condition when tests
   re-run rapidly in the same shell session. Documented in CLAUDE.md
   "Known operational quirk".

This document captures investigation results for each, decides whether
each is a real bug, and proposes per-item fixes. The user explicitly
asked for all three to be addressed, but two of them turn out to be
non-bugs (with different reasons), so the actual implementation work is
narrower than expected.

---

## Section 1 — Findings

### Finding 1: `parseGenericType` 1-dot limitation (REAL BUG, but downstream catalog also needs a fix)

**SCOPE EXPANSION (per codex review of commit f7731c2)**: fixing
`parseGenericType` alone is insufficient. The catalog's
`typeNameParts` helper silently drops qualification for any `Names`
list with more than 2 components — so a parser that produces
`Names=["db","schema","mytype"]` would feed the catalog a 3-component
list, and the catalog would silently treat it as bare `"mytype"` with
no schema. This is a **silent semantic bug that the parser fix would
expose**. The plan now includes a paired catalog fix.



**Symptom**: any type position that uses a 3-or-more-component qualified
name fails to parse:

| SQL | Result |
|---|---|
| `SELECT 1::pg_catalog.int4` (2-component) | ✅ OK |
| `CREATE TABLE t (c pg_catalog.int4)` | ✅ OK |
| `SELECT 1::db.schema.mytype` (3-component) | ❌ FAIL |
| `CREATE TABLE t (c db.schema.mytype)` | ❌ FAIL |
| `SELECT CAST(NULL AS db.schema.mytype)` | ❌ FAIL |
| `ALTER TABLE t ALTER COLUMN c TYPE db.schema.mytype` | ❌ FAIL |
| `CREATE TABLE t (c a.b.c.d)` (4-component) | ❌ FAIL |
| `CREATE FUNCTION f(x db.schema.mytype) RETURNS int AS '...'` | ❌ FAIL |
| `CREATE FUNCTION f() RETURNS db.schema.mytype AS '...'` | ❌ FAIL |

Every type position that goes through `parseGenericType` is affected.
Verified empirically against the current HEAD of `pg-first-sets`.

**Comparison with PG**: PG's grammar uses a recursive `attrs` production:

```bison
attrs:  '.' attr_name                   { $$ = list_make1(makeString($2)); }
      | attrs '.' attr_name             { $$ = lappend($1, makeString($3)); }
      ;

GenericType:
        type_function_name opt_type_modifiers
      | type_function_name attrs opt_type_modifiers
      ;
```

Reference: `postgres/src/backend/parser/gram.y:7007-7011, 14385-14398`.

Because `attrs` is recursive, PG accepts arbitrarily many `.attr_name`
continuations. `pg_catalog.int4`, `db.schema.mytype`, and even
`a.b.c.d.e` are all syntactically valid.

**omni's current code** (`pg/parser/type.go:192-232`):

```go
func (p *Parser) parseGenericType() (*nodes.TypeName, error) {
    name, err := p.parseTypeFunctionName()
    if err != nil {
        return nil, err
    }

    if p.cur.Type == '.' {     // ← consumes ONE dot
        p.advance()
        attr, err := p.parseAttrName()
        if err != nil {
            return nil, err
        }
        typmods, err := p.parseOptTypeModifiers()
        if err != nil {
            return nil, err
        }
        return &nodes.TypeName{
            Names: &nodes.List{Items: []nodes.Node{
                &nodes.String{Str: name},
                &nodes.String{Str: attr},
            }},
            Typmods: typmods,
            Loc:     nodes.NoLoc(),
        }, nil
    }
    // ... unqualified case ...
}
```

The `if p.cur.Type == '.'` consumes exactly one dot. PG's `attrs` is
recursive (`for p.cur.Type == '.'`).

**Crucial detail**: omni already has a recursive `parseAttrs` helper at
`pg/parser/name.go:457`:

```go
func (p *Parser) parseAttrs() (*nodes.List, error) {
    if _, err := p.expect('.'); err != nil {
        return nil, err
    }
    name, err := p.parseAttrName()
    // ...
    items := []nodes.Node{&nodes.String{Str: name}}
    for p.cur.Type == '.' {       // ← LOOP
        p.advance()
        name, err = p.parseAttrName()
        // ...
        items = append(items, &nodes.String{Str: name})
    }
    return &nodes.List{Items: items}, nil
}
```

It's already used by 3 other call sites (`extension.go`, `fdw.go`,
`type.go:parseFuncType`). The fix for `parseGenericType` is to delegate
to `parseAttrs` instead of hand-rolling a single-dot case.

**Fix**:

```go
func (p *Parser) parseGenericType() (*nodes.TypeName, error) {
    name, err := p.parseTypeFunctionName()
    if err != nil {
        return nil, err
    }

    nameItems := []nodes.Node{&nodes.String{Str: name}}
    if p.cur.Type == '.' {
        attrs, err := p.parseAttrs()
        if err != nil {
            return nil, err
        }
        nameItems = append(nameItems, attrs.Items...)
    }

    typmods, err := p.parseOptTypeModifiers()
    if err != nil {
        return nil, err
    }
    return &nodes.TypeName{
        Names:   &nodes.List{Items: nameItems},
        Typmods: typmods,
        Loc:     nodes.NoLoc(),
    }, nil
}
```

~15 lines vs the current ~36. Matches PG's grammar exactly. No new
abstraction — reuses the existing `parseAttrs`.

**Paired catalog fix** (added per codex review):

`pg/catalog/nodeutil.go:189-206` `typeNameParts` currently has:

```go
func typeNameParts(tn *nodes.TypeName) (schema, name string) {
    if tn.Names == nil {
        return "", ""
    }
    items := tn.Names.Items
    switch len(items) {
    case 1:
        return "", stringVal(items[0])
    case 2:
        return stringVal(items[0]), stringVal(items[1])
    default:
        // Last element is the type name.
        if len(items) > 0 {
            return "", stringVal(items[len(items)-1])  // ← drops schema AND catalog
        }
        return "", ""
    }
}
```

For `len > 2`, it returns `("", lastItem)` — silently dropping all
qualification beyond the last component. Today this path is unreachable
(parser produces at most 2-component names), but the parser fix above
would produce `Names=["db","schema","mytype"]` for `db.schema.mytype`,
and `typeNameParts` would silently turn it into bare `"mytype"`.

PG's actual semantics: in `catalog.schema.type`, the catalog must match
the current database name. PG accepts the syntax and validates the
catalog at name-resolution time. omni doesn't track the current
database, so we can't perform that validation, but we can still extract
the schema and name correctly:

```go
case 1:
    return "", stringVal(items[0])
case 2:
    return stringVal(items[0]), stringVal(items[1])
case 3:
    // catalog.schema.name — drop the catalog prefix. PG would validate
    // it matches the current database, but we don't track that here.
    return stringVal(items[1]), stringVal(items[2])
default:
    // 4+ components: PG also rejects these at name resolution time
    // ("improper qualified name"). Fall back to returning the last
    // component as the name with empty schema, so the downstream
    // resolver error at least mentions an identifier the user can
    // recognize. A better error message would require changing
    // typeNameParts's signature to return a bool/error, which is
    // out of scope for this fix.
    if len(items) > 0 {
        return "", stringVal(items[len(items)-1])
    }
    return "", ""
}
```

**Important caveat (per codex review of commit 08090e6)**: the catalog
fix is **NOT behavior-neutral** for one currently-reachable path —
`%TYPE` references like `RETURNS schema.tab.col%TYPE`. omni's
`parseFuncType` already produces `TypeName{PctType: true, Names:
[schema, tab, col]}` for such inputs. Today, `typeNameParts` returns
`("", "col")` and resolution fails as `type "col" does not exist`.
After the fix, it returns `("tab", "col")` and resolution fails as
`type "tab.col" does not exist` (or similar).

**This is not a regression — both error paths are wrong.** omni's
catalog has **no PctType handler at all** (verified by `grep -rn
"PctType" pg/catalog/` — zero matches). `%TYPE` semantics are
"resolve to the type of column X in table Y", which requires actual
column-type lookup, not a `(schema, name)` split. Fixing PctType
resolution properly is **out of scope** for this PR; documenting
the existing brokenness is enough.

The commit message for the catalog fix should explicitly call out:

1. Behavior change is intentional for non-PctType 3-component types
   (the bug being fixed)
2. Behavior also changes incidentally for PctType inputs, but on a
   path that was never working — error message changes from
   `type "col" does not exist` to `type "tab.col" does not exist`,
   which is more informative even though both are wrong
3. Full PctType resolution is out of scope; tracked separately.

**Call site audit for `typeNameParts`**: 7 sites. 5 use only the name
(`_, name := ...`), 2 use both (`schema, name := ...` at
`nodeutil.go:145` and `nodeutil.go:350`). For the name-only sites, the
fix is transparent. For the schema-using sites, the fix changes the
returned `schema` from `""` to `items[len-2]` for 3-component names —
which is the correct behavior matching PG.

**Test coverage**: in addition to parser tests for 6 type positions at
2/3/4-component depths, add a catalog test that asserts
`typeNameParts` returns the expected `(schema, name)` pair for 3-component
input, and that 4+component returns an empty/invalid result rather than
silently truncating.

**Test matrix expansion (per codex item 6 — broader parseGenericType
call paths)**: codex's audit found additional broken positions I
missed. These all hit `parseGenericType` via different call chains and
all fail at HEAD:

- `RETURNS TABLE (x db.schema.mytype)` (parseTableFuncColumn → parseFuncType → parseTypename → parseSimpleTypename → parseGenericType)
- `CREATE OPERATOR === (..., LEFTARG = db.schema.mytype, ...)` (parseDefArg → parseFuncType → ...)
- `CREATE SEQUENCE s AS db.schema.mytype` (sequence AS → parseSimpleTypename → parseGenericType)
- `XMLSERIALIZE(... AS db.schema.mytype)` (parseSimpleTypename → parseGenericType)
- `json_serialize(... RETURNING db.schema.mytype)` (parseTypename → ...)

These should all be in the test matrix because they exercise different
call chains to the same buggy function.

**Cost**: ~1 hour of work (parser fix + catalog fix + 11-position test
matrix). 1-2 commits.

### Finding 2: 3 `isColId() || isTypeFunctionName()` call sites are DEAD CODE (NOT a bug, just style)

**Sites** (from commit `fc43201` follow-up notes):
- `pg/parser/expr.go:1475` in `parseAExprPrimary` default branch
- `pg/parser/create_table.go:901` in `parseOptQualifiedName`
- `pg/parser/create_index.go:273` in `isIndexElemOpclassStart`

**Original concern** (from Phase 2 audit): the `isAExprStart` oracle
test against PG 17 showed that 22 TypeFuncNameKeyword tokens (INNER,
LEFT, JOIN, CROSS, NATURAL, etc.) are NOT in PG's expression FIRST set,
but the `isColId() || isTypeFunctionName()` predicate accepts them.
This looked like a divergence from PG's strict grammar.

**Investigation result**: at all 3 call sites, the OR-with-`isTypeFunctionName`
is **dead code at runtime**, because the immediately-following parse
function uses `parseColId` (which rejects TypeFuncNameKeyword) anyway.

#### Trace 1: `expr.go:1475 → parseColumnRefOrFuncCall`

```go
// expr.go:1475
if p.isColId() || p.isTypeFunctionName() {
    return p.parseColumnRefOrFuncCall()
}

// expr.go:2420 parseColumnRefOrFuncCall:
func (p *Parser) parseColumnRefOrFuncCall() (nodes.Node, error) {
    name, err := p.parseColId()  // ← rejects TypeFuncNameKeyword
    // ...
}
```

For a TypeFuncNameKeyword token like `INNER`:
- Predicate returns true (via `isTypeFunctionName`)
- Enters `parseColumnRefOrFuncCall`
- `parseColId` errors out because `INNER` is `TypeFuncNameKeyword`, not
  `ColId` (UnreservedKeyword + ColNameKeyword + IDENT)

**Empirical verification**: `SELECT inner` (and 8 other TypeFuncNameKeyword
single-token expressions) all fail at HEAD of `pg-first-sets`. omni
already rejects them — the `|| isTypeFunctionName()` is just unreachable
acceptance.

#### Trace 2: `create_table.go:901 → parseOptQualifiedName → parseAnyName → parseColId`

```go
// create_table.go:900 parseOptQualifiedName:
func (p *Parser) parseOptQualifiedName() *nodes.List {
    if p.isColId() || p.isTypeFunctionName() {
        name, _ := p.parseAnyName()
        return name
    }
    return nil
}

// name.go:164 parseAnyName:
func (p *Parser) parseAnyName() (*nodes.List, error) {
    // ...
    id, err := p.parseColId()  // ← rejects TypeFuncNameKeyword
    // ...
}
```

PG's `any_name: ColId | ColId attrs` — uses `ColId`, not
`type_function_name`. omni's `parseAnyName` correctly enforces this.
Same dead-code-OR situation.

#### Trace 3: `create_index.go:273 → isIndexElemOpclassStart → parseOptQualifiedName`

```go
// create_index.go:266 isIndexElemOpclassStart:
func (p *Parser) isIndexElemOpclassStart() bool {
    switch p.cur.Type {
    case ASC, DESC, NULLS_LA, ',', ')', 0,
        WITH, INCLUDE, WHERE, TABLESPACE, USING, DO:
        return false
    }
    return p.isColId() || p.isTypeFunctionName()
}

// create_index.go:243 parseIndexElemOpclass:
if !p.isIndexElemOpclassStart() {
    return nil, nil
}
opclass := p.parseOptQualifiedName()  // → parseAnyName → parseColId
```

The downstream `parseColId` enforces ColId. The OR-with-`isTypeFunctionName`
in the predicate causes `parseIndexElemOpclass` to **try** parsing for
TypeFuncNameKeyword tokens, then fail — instead of the predicate rejecting
upfront and the outer parser trying the next clause.

This means a TypeFuncNameKeyword token at the index-element position
produces a slightly different error message under the current code vs
after the cleanup. But neither produces correct parse output, and **no
valid SQL changes behavior** because the `isIndexElemOpclassStart` negative
filter already covers the only TypeFuncNameKeyword tokens that are valid
index clause terminators.

#### Conclusion

All 3 sites are dead-code ORs. **Removing them is a pure cleanup, not a
behavior change for valid SQL.** The cleanup makes the predicate match
PG's grammar (ColId, not type_function_name) and removes confusion for
future readers. Error messages for invalid SQL with TypeFuncNameKeyword
tokens at these positions may shift slightly.

**Fix**: 3 one-line edits, removing `|| p.isTypeFunctionName()` from each
of the 3 sites. Plus a small test that `SELECT inner`, `SELECT join`,
`SELECT left` continue to fail (sanity that the cleanup doesn't
accidentally accept anything new).

**Cost**: ~15 minutes. 1 commit.

### Finding 3: Ryuk reaper flake (rare, low impact)

**Symptom**: when `pg/parser` tests are re-run rapidly in the same shell
session, the testcontainers-go Ryuk reaper occasionally terminates the
previous test process's container before the new test process's
`startFirstSetOracle` can ping it. The test fails with
`connection refused`.

**Reproduction**: 1 confirmed occurrence during the original Phase 1
smoke test in this conversation. 5 follow-up rapid reruns all passed.
The frequency cannot be validated from source alone — it depends on
local Docker/Ryuk timing — so any "X% per run" estimate is a guess. CI
should be unaffected because per-job process isolation gives Ryuk a
clean lifetime per process. (Codex correctly flagged my earlier "rare"
claim as not validatable from source; this paragraph reflects what we
actually know.)

**Investigation summary**:

- omni's mysql/catalog uses **per-test fresh containers** (no `sync.Once`),
  so it doesn't hit the Ryuk reuse race. Each test creates and destroys
  its own container.
- omni's pg/catalog and pg/parser both use `sync.Once` for efficiency —
  containers are shared across tests in the same process, started once.
- Neither package has a `TestMain` for explicit shutdown.
- testcontainers-go's Ryuk feature can be disabled via
  `TESTCONTAINERS_RYUK_DISABLED=true`. Trade-off: with Ryuk disabled,
  containers leak indefinitely (until manual cleanup or system reboot)
  on test process exit, because there's no other shutdown hook.

**Three options**:

#### Option A: Status quo (no code change, document only)

Keep CLAUDE.md's "Known operational quirk" note. Suggest local dev set
`TESTCONTAINERS_RYUK_DISABLED=true` in shell rc if the flake bothers
them. CI is unaffected.

- Cost: 0 lines
- Local cost: ~30 seconds of retry per month (estimated)
- Container leak risk: 0

#### Option B: Default-disable Ryuk via `init()` in `first_set_oracle_test.go`

```go
func init() {
    if os.Getenv("TESTCONTAINERS_RYUK_DISABLED") == "" {
        os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
    }
}
```

- Cost: 1 line + import
- Local cost: 0 flakes
- Container leak risk: ~50MB per test run, indefinite. Each
  `go test ./pg/parser/...` invocation leaks one postgres:17-alpine
  container until reboot or `docker container prune`.
- Trade-off: most dev environments have plenty of memory and reboot
  occasionally, so the leak is usually invisible. But the leak IS real
  and could surprise someone running tests in a long-lived CI runner.

#### Option C: Add `TestMain` for explicit shutdown

```go
// pg/parser/main_test.go (new file)
package parser

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    code := m.Run()
    shutdownFirstSetOracle() // explicit cleanup
    os.Exit(code)
}
```

Plus a `shutdownFirstSetOracle()` helper that terminates the container
if it was started.

- Cost: ~15 lines (TestMain + shutdown helper + exposing the container
  reference from the sync.Once closure)
- Local cost: 0 flakes (because explicit shutdown happens before the
  next test process starts; Ryuk has nothing to race against)
- Container leak risk: 0 (explicit shutdown path, even on test panic if
  we use `m.Run()` deferred properly)
- Risk: novel pattern in pg/parser. omni has no `TestMain` files
  anywhere in this package; introducing one could surprise reviewers.
  Need to verify Go's test framework handles file-level `TestMain` in
  one `_test.go` file when other `_test.go` files in the same package
  exist (it does — `TestMain` is per-package, not per-file).

#### Recommendation

**Option A** for the immediate plan. The flake is rare enough that the
fix isn't urgent, and both code options have non-trivial downsides:
B leaks containers, C introduces a new pattern.

If the user wants a code fix, **Option C** is the right one — it solves
both the flake AND the implicit container leak that's already happening
on every test process exit. But this should be a separate PR,
not bundled with the type-name fix, because it changes a different
subsystem (test infrastructure, not the parser proper).

This is the only one of the 3 follow-ups where I'd recommend NOT
shipping a fix in this PR.

---

## Section 2 — Detection method

### Finding 1 detection

1. **Empirical test** — wrote a Go program iterating 10 SQL strings
   covering 2/3/4-component qualified names in 6 type positions. Ran
   against HEAD of `pg-first-sets`. Result: 2-component works, all
   3+component fail across all positions.
2. **Source trace** — read `parseGenericType` (`type.go:192-232`),
   confirmed the `if p.cur.Type == '.'` consumes exactly one dot.
3. **PG comparison** — read `gram.y:7003-7011` for `attrs` (recursive)
   and `gram.y:14385-14398` for `GenericType` (uses `attrs`). Confirmed
   PG's grammar accepts arbitrary depth.
4. **Discovered existing helper** — `grep "parseAttrs"` found
   `name.go:457` already implements the recursive form. 3 other call
   sites use it (`extension.go`, `fdw.go`, `type.go:parseFuncType`).
   Conclusion: the fix is to make `parseGenericType` use it.

### Finding 2 detection

1. **Source trace per call site** — read each of the 3 sites and the
   immediately-following parse function.
2. **Identified the dead-code pattern** — at each site, the predicate
   is `isColId() || isTypeFunctionName()` but the next call uses
   `parseColId` (which only accepts ColId-category tokens). So
   TypeFuncNameKeyword tokens enter the parse function and are
   immediately rejected.
3. **Empirical verification** — wrote a Go program that tries `SELECT
   inner`, `SELECT left`, `SELECT join`, etc. (9 TypeFuncNameKeyword
   single-token expressions). All 9 fail at HEAD. Confirmed omni already
   rejects them.
4. **PG grammar check** — `parseColumnRefOrFuncCall` mirrors
   `a_expr → c_expr → columnref → ColId`. `parseAnyName` mirrors PG's
   `any_name: ColId | ColId attrs`. Both use ColId. The OR-with-
   `isTypeFunctionName` in the predicates is inconsistent with both
   PG and the downstream parse functions.

### Finding 3 detection

1. **Reproduction attempt** — 5 rapid `go test` reruns. None failed.
   Flake is intermittent.
2. **Sibling pattern check** — read `mysql/catalog/container_test.go`.
   It uses per-test fresh containers (no sync.Once), so it doesn't hit
   the Ryuk reuse race.
3. **TestMain check** — grep'd `pg/parser/*_test.go` and
   `pg/catalog/*_test.go`. No `TestMain` exists. Adding one would be
   novel for the package.
4. **Trade-off analysis** — 3 options enumerated, each with cost, leak
   risk, and code surface.

---

## Section 3 — Plan

### Goal

Fix finding 1 (real bug) and finding 2 (style cleanup with PG-grammar
alignment). Document but don't fix finding 3 (Option A: doc-only).

### Commit breakdown

**Commit 1 (CATALOG)**: Fix `pg/catalog/nodeutil.go:typeNameParts` to
handle 3-component names by treating `[catalog, schema, name]` as
`(schema, name)`, and explicitly returning empty for 4+component
(rather than silently truncating). Add `TestTypeNamePartsThreeComponent`
covering 1/2/3/4-component inputs. **Land this BEFORE commit 2** so the
parser fix doesn't expose the silent-truncation bug.

**Commit 2 (PARSER)**: Fix `parseGenericType` to use `parseAttrs` for
arbitrary qualified type depth. Add `TestParseQualifiedTypeMultiComponent`
covering all 11 positions (6 originally listed + 5 found by codex
review): CAST, TYPECAST, CREATE TABLE column, ALTER TABLE, CREATE
FUNCTION param, CREATE FUNCTION return, RETURNS TABLE column, CREATE
OPERATOR LEFTARG, CREATE SEQUENCE AS, XMLSERIALIZE, json_serialize
RETURNING — at 2-component, 3-component, and 4-component depths where
applicable.

**Commit 3 (CLEANUP)**: Remove dead-code `|| p.isTypeFunctionName()`
from the 3 call sites at `expr.go:1475`, `create_table.go:901`,
`create_index.go:273`. Add a small
`TestTypeFuncNameKeywordRejectedAsExpression` covering `SELECT inner`,
`SELECT left`, etc. as a regression sanity (these already fail today;
the test locks in that they continue to fail after the cleanup).

**Commit 4 (DOCS)**: Update `pg/parser/CLAUDE.md` to remove the
"Three-component qualified type names" entry from "Known limitations"
(it will be fixed by commits 1+2). Leave the Ryuk reaper note as-is
(Option A for finding 3).

**Commit ordering matters**: catalog fix MUST land before parser fix.
Otherwise the parser starts producing 3-component ASTs while the catalog
still silently drops them, creating a temporary window of silent
semantic bug.

**Caveat on "additive"**: the catalog fix is NOT purely additive — see
the "Important caveat" section above for the `%TYPE` interaction.
However, both behaviors on the `%TYPE` path are wrong (omni has no
PctType handler), so the change is "different broken" rather than
"newly broken". The catalog-before-parser ordering is still the right
choice.

### Test matrix

```go
// pg/parser/qualified_type_multi_test.go (new file)
func TestParseQualifiedTypeMultiComponent(t *testing.T) {
    cases := []struct {
        name      string
        sql       string
        wantNames []string  // expected ArgType.Names list
    }{
        // 2-component (regression: must still work)
        {
            name:      "2-component CAST",
            sql:       `SELECT CAST(NULL AS pg_catalog.int4)`,
            wantNames: []string{"pg_catalog", "int4"},
        },
        // 3-component — the 6 positions originally listed
        {
            name:      "3-component CAST",
            sql:       `SELECT CAST(NULL AS db.schema.mytype)`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component TYPECAST",
            sql:       `SELECT 1::db.schema.mytype`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component CREATE TABLE column",
            sql:       `CREATE TABLE t (c db.schema.mytype)`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component ALTER TABLE",
            sql:       `ALTER TABLE t ALTER COLUMN c TYPE db.schema.mytype`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component CREATE FUNCTION param",
            sql:       `CREATE FUNCTION f(x db.schema.mytype) RETURNS int AS 'select 1' LANGUAGE sql`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component CREATE FUNCTION return",
            sql:       `CREATE FUNCTION f() RETURNS db.schema.mytype AS 'select 1' LANGUAGE sql`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        // 3-component — additional positions found by codex review
        {
            name:      "3-component RETURNS TABLE column",
            sql:       `CREATE FUNCTION f() RETURNS TABLE (x db.schema.mytype) AS 'select 1' LANGUAGE sql`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component CREATE OPERATOR LEFTARG",
            sql:       `CREATE OPERATOR === (FUNCTION = int4eq, LEFTARG = db.schema.mytype, RIGHTARG = int4)`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component CREATE SEQUENCE AS",
            sql:       `CREATE SEQUENCE s AS db.schema.mytype`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component XMLSERIALIZE",
            sql:       `SELECT XMLSERIALIZE(CONTENT '<a/>' AS db.schema.mytype)`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        {
            name:      "3-component json_serialize RETURNING",
            sql:       `SELECT JSON_SERIALIZE('{}' RETURNING db.schema.mytype)`,
            wantNames: []string{"db", "schema", "mytype"},
        },
        // 4-component (PG accepts at parse time but errors at semantic
        // resolution because the catalog component must match the current
        // database AND the chain can't go deeper). omni's parser should
        // accept the syntax; the catalog fix below should reject the
        // semantic.
        {
            name:      "4-component parser accepts",
            sql:       `CREATE TABLE t (c a.b.c.d)`,
            wantNames: []string{"a", "b", "c", "d"},
        },
    }
    // ... walk to TypeName and assert Names ...
}

// pg/catalog/nodeutil_test.go addition (new test in existing file
// or new file if nodeutil_test.go doesn't exist)
func TestTypeNamePartsThreeComponent(t *testing.T) {
    cases := []struct {
        name       string
        items      []string  // raw Names list strings
        wantSchema string
        wantName   string
    }{
        {
            name:       "1-component bare",
            items:      []string{"int4"},
            wantSchema: "",
            wantName:   "int4",
        },
        {
            name:       "2-component schema-qualified",
            items:      []string{"pg_catalog", "int4"},
            wantSchema: "pg_catalog",
            wantName:   "int4",
        },
        {
            name:       "3-component catalog.schema.name",
            items:      []string{"db", "schema", "mytype"},
            wantSchema: "schema",
            wantName:   "mytype",
        },
        {
            // 4+ components: PG also rejects at name resolution. omni
            // returns the last component with empty schema as a soft
            // fallback so the downstream resolver error mentions a
            // recognizable identifier. A typed error would require
            // changing typeNameParts's signature, which is out of scope.
            name:       "4-component soft fallback",
            items:      []string{"a", "b", "c", "d"},
            wantSchema: "",
            wantName:   "d",
        },
    }
    // ... build TypeName, call typeNameParts, assert ...
}

// pg/parser/typefuncname_keyword_rejected_test.go (new file)
func TestTypeFuncNameKeywordRejectedAsExpression(t *testing.T) {
    // After removing the dead-code OR at the 3 call sites, these MUST
    // still fail (they fail today; the test locks the behavior in).
    sqls := []string{
        `SELECT inner`,
        `SELECT left`,
        `SELECT right`,
        `SELECT join`,
        `SELECT cross`,
        `SELECT outer`,
        `SELECT full`,
        `SELECT natural`,
        `SELECT verbose`,
    }
    for _, sql := range sqls {
        t.Run(sql, func(t *testing.T) {
            _, err := Parse(sql)
            if err == nil {
                t.Errorf("expected parse error for %q (TypeFuncNameKeyword as expression atom), got nil", sql)
            }
        })
    }
}
```

### Risks

1. **Finding 1 fix might surface other 3+component-related bugs.** The
   resulting AST has `Names = [a, b, c, ...]`. Downstream code (catalog
   resolution, deparser, etc.) may or may not handle 3+component name
   lists. The fix changes the parser to PRODUCE such ASTs; if downstream
   crashes on them, the surface area is wider than expected. Mitigation:
   run the full omni test suite after commit 1, not just `pg/parser`.

2. **Finding 2 cleanup might change error messages for invalid SQL.**
   Specifically, `CREATE INDEX ... ON t (col INNER)` where INNER is in
   an unusual position. Today: tries to parse INNER as opclass, fails
   inside `parseColId`. After: `isIndexElemOpclassStart` returns false,
   parser tries the next clause (probably WHERE / WITH / ASC), also
   fails but with a different message. **No valid SQL is affected** —
   only error messages for already-broken SQL. Acceptable.

3. **Test interaction with the existing first_sets discipline.** If
   any FIRST-set predicate (like `isAExprStart`) used the old `isColId()
   || isTypeFunctionName()` shape internally, removing it would change
   predicate behavior. Mitigation: `isAExprStart` is composed
   differently (uses `isColId()` alone, per the Phase 2 codex review).
   Verified by reading `first_sets.go:170-179`. No interaction.

4. **Code reviewer may push back on dropping the OR.** The OR has been
   in the codebase for a long time; someone might have added it for a
   reason that's now lost. Mitigation: the commit message and trace
   evidence in this plan doc make the dead-code analysis clear. If
   review surfaces a new use case I missed, fall back to keeping the
   OR but adding an explanatory comment.

### Time estimate (revised after codex review)

| Phase | Time |
|---|---|
| Commit 1 (catalog typeNameParts fix + 4-case unit test) | 20-30 min |
| Commit 2 (parseGenericType fix + 11-position multi-component tests) | 45-60 min |
| Commit 3 (3-site dead-code cleanup + sanity test) | 15-30 min |
| Commit 4 (CLAUDE.md update) | 10 min |
| Full pg/parser + pg/catalog + omni-wide smoke | 30 min |
| Codex review + iteration | 30-60 min |
| **Total** | **~3 hours** |

Increased from "~2 hours" because the codex review revealed (a) the
catalog dependency means there's a 4th commit, (b) the test matrix is
~2x larger (11 positions instead of 6), (c) a wider smoke is needed
since the catalog change touches code paths beyond pg/parser.

### Deferred / out of scope

- **Finding 3 (Ryuk flake)**: Option A (no code change). Rationale: the
  flake is rare, the workaround is documented, and both code options
  have downsides (B leaks containers, C introduces a new TestMain
  pattern). If a future PR wants to address it, Option C is the
  recommended approach.
- **mysql / mssql / oracle parser audits** for analogous `parseGenericType`
  limitations. Each may have its own grammar shape; out of scope for
  this PR.
- **Catalog / deparser support for 3+component type names**. Risk #1
  above mentions this. If commit 1's full-suite run surfaces a downstream
  issue, triage and either fix or revert commit 1 with a documented
  follow-up.

---

## Review questions for codex

1. **Finding 1 fix correctness**: does the proposed `parseAttrs`-based
   rewrite of `parseGenericType` produce the correct AST shape for both
   2-component and 3+component names? Specifically, is
   `Names = lcons(name, attrs.Items)` the right list-construction order?
2. **Finding 2 dead-code claim**: trace each of the 3 call sites and
   verify the claim that `parseColId` downstream rejects
   TypeFuncNameKeyword. Look for any ColId-bypassing code path I missed.
3. **Finding 2 risk**: are there valid SQL forms with TypeFuncNameKeyword
   tokens at expression / qualified-name / opclass positions that
   currently parse via the dead-code OR? If yes, my "no valid SQL is
   affected" claim is wrong.
4. **Finding 3 trade-off**: is Option A really the right choice, or should
   we ship Option C (TestMain) as part of this PR?
5. **Test matrix completeness**: are there 3+component qualified type
   positions I'm not testing? Spot-check by grepping `parseTypename` /
   `parseFuncType` / `parseSimpleTypename` call sites.
6. **Downstream integration risk**: is there evidence in the catalog /
   deparser code that 3+component type ASTs are or are not handled?
7. **Anything I'm missing**: are there other follow-up items from the
   `pg-first-sets` branch's commit history that I haven't enumerated?
