# MySQL Parser/Validator Split Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Split `mysql/parser` into a pure MySQL-grammar parser + a separate `mysql/validate` package that performs all static semantic checks. Parser becomes MySQL-aligned (accepts everything MySQL's `sql_yacc.yy` accepts); semantic errors (undeclared variable, duplicate declaration, missing RETURN, etc.) move to the validator. Also fix the two reported bugs — string-literal aliases and `SET sys_var` inside routine bodies.

**Architecture:**
- `mysql/parser` — token stream → AST. Only raises `ParseError` for grammar violations. No `procScope`, no `lookupVar`, no label/cursor/condition reachability checks.
- `mysql/validate` — new package. `Validate(*ast.List, Options) []Diagnostic`. Walks AST, rebuilds scope from DECLARE statements, emits structured diagnostics for all checks previously inlined in the parser. Mirrors MySQL's `sp_head::parse` / `sp_pcontext` layer.
- Callers decide: `Parse` alone is enough for catalog/deparse/completion. Advisor/linter/semantic layer call `Validate` on top.

**Tech Stack:** Go, existing `mysql/ast` node types (no schema changes expected), standard testing tools.

**Worktree:** `/Users/rebeliceyang/Github/omni/.worktrees/mysql-fix-2026-04-21`, branch `junyi/mysql-fix-2026-04-21` (based on `origin/main` @ `b22b927`).

---

## Semantic-Check Inventory (what moves out of parser)

Full list of currently inlined semantic errors in `mysql/parser/*.go` (confirmed by grep):

| # | Site | Error text | Moves to |
|---|---|---|---|
| S1 | `scope.go:150` | `duplicate variable declaration: <name>` | validate |
| S2 | `scope.go:163` | `duplicate condition declaration: <name>` | validate |
| S3 | `scope.go:176` | `duplicate cursor declaration: <name>` | validate |
| S4 | `scope.go:194` | `duplicate label: <name>` | validate |
| S5 | `compound.go:494` | `duplicate condition value in handler declaration` | validate |
| S6 | `compound.go:502` | `undeclared condition: <name>` (inside HANDLER FOR list) | validate |
| S7 | `compound.go:1029` | `LEAVE references undeclared label: <name>` | validate |
| S8 | `compound.go:1057` | `ITERATE references undeclared loop label: <name>` | validate |
| S9 | `compound.go:1079` | `RETURN is only allowed inside a function body` | validate |
| S10 | `compound.go:1128` (`requireCursor`) | `undeclared cursor: <name>` (OPEN/FETCH/CLOSE) | validate |
| S11 | `set_show.go:362` | `undeclared variable: <name>` (bare SET target in routine body) | validate **(+Bug 2 fix: sysvar fallback)** |
| S12 | `create_function.go:232` | `no RETURN found in function body` | validate |

Staying in parser (these are grammar-level, not semantic):
- `compound.go:63` — `checkDeclarePhase` (DECLARE ordering) — matches MySQL yacc's grammar ordering, keep as grammar rule.
- `compound.go:84,889,953,1003` — `checkLabelMatch` (begin/end label text matches) — matches yacc-level check in MySQL; keep.

Bugs to fix during the refactor:
- **Bug 1 (grammar, parser-only):** select alias and table alias must accept `TEXT_STRING` as well as `ident`. Fix in Phase 2.
- **Bug 2 (semantic, validator):** SET bare target in routine body that isn't a declared local must fall back to the known-system-variable table before reporting undeclared. Fix in Phase 6 when S11 moves.

---

## Phase Overview

- **Phase 1** — baseline: verify worktree passes existing tests.
- **Phase 2** — fix Bug 1 (pure grammar fix, applies before validator even exists).
- **Phase 3** — scaffold `mysql/validate` package with scope reconstruction + AST walker, no checks yet. All tests pass; nothing wired.
- **Phase 4** — migrate semantic checks one by one, each as a TDD task. Each task: move one check to validator, keep parser side also raising (temporarily) so nothing breaks, add validator test, swap callers.
- **Phase 5** — remove parser-side `procScope` + all S1–S12 checks. Make parser pure grammar.
- **Phase 6** — fix Bug 2 in validator (add sysvar table + fallback).
- **Phase 7** — caller migration & regression tests (completion improvements, catalog unaffected).
- **Phase 8** — docs + cleanup.

Each phase ends with a commit. Phases 4.x ship incrementally so review stays bounded.

---

## Phase 1: Baseline

### Task 1.1: Verify baseline tests

**Files:** none — read-only.

**Step 1: Run the parser tests to confirm green start.**

Run:
```bash
cd /Users/rebeliceyang/Github/omni/.worktrees/mysql-fix-2026-04-21
go test ./mysql/parser/... ./mysql/catalog/... ./mysql/completion/... ./mysql/deparse/... ./mysql/ast/... 2>&1 | tail -50
```
Expected: all packages PASS (no FAIL). Note any baseline failures before starting — must be zero.

**Step 2: Commit plan file.**

```bash
git add docs/plans/2026-04-21-mysql-parser-validator-split.md
git commit -m "docs: plan for MySQL parser/validator split"
```

---

## Phase 2: Bug 1 (string-literal alias)

### Task 2.1: Add `parseIdentOrText` helper

**Files:**
- Modify: `mysql/parser/name.go` (insert after `parseIdent` at line ~548)

**Step 1: Write the helper.**

Insert after the `parseIdent` function:

```go
// parseIdentOrText parses MySQL's `ident_or_text` grammar rule: either an
// identifier (including non-reserved keywords) or a string literal. Used for
// aliases and other positions where MySQL accepts TEXT_STRING_sys equivalently.
//
// Ref: mysql-server sql/sql_yacc.yy — ident_or_text rule
//
//	ident_or_text:
//	    ident
//	  | TEXT_STRING_sys
func (p *Parser) parseIdentOrText() (string, int, error) {
	if p.cur.Type == tokSCONST {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return p.parseIdent()
}
```

**Step 2: Commit.**

```bash
git add mysql/parser/name.go
git commit -m "feat(mysql/parser): add parseIdentOrText for ident_or_text grammar rule"
```

---

### Task 2.2: Use `parseIdentOrText` for SELECT alias

**Files:**
- Modify: `mysql/parser/select.go:771-790`
- Test: `mysql/parser/oracle_corpus_test.go` (add) or new `mysql/parser/alias_test.go`

**Step 1: Write the failing test.**

Create `mysql/parser/alias_test.go`:

```go
package parser

import "testing"

func TestSelectAliasStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT 1 AS 'start'`,
		`SELECT 1 AS "start"`,
		`SELECT 1 'start'`,
		`SELECT a AS 'label', b FROM t`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails.**

```bash
go test ./mysql/parser/ -run TestSelectAliasStringLiteral -v
```
Expected: FAIL with `expected identifier` on the `'start'` cases.

**Step 3: Swap `parseIdent` for `parseIdentOrText` in `parseSelectExpr`.**

In `mysql/parser/select.go`, replace the two call sites inside `parseSelectExpr`:

```go
	// Check for AS alias or implicit alias
	var alias string
	if _, ok := p.match(kwAS); ok {
		aliasLoc := p.pos()
-		name, _, err := p.parseIdent()
+		name, _, err := p.parseIdentOrText()
		if err != nil {
			return nil, err
		}
		alias = name
		if p.completing {
			p.addSelectAliasPosition(aliasLoc)
		}
-	} else if p.isIdentToken() && !p.isSelectTerminator() {
+	} else if (p.isIdentToken() || p.cur.Type == tokSCONST) && !p.isSelectTerminator() {
		// Implicit alias (identifier or string literal without AS)
		aliasLoc := p.pos()
-		alias, _, err = p.parseIdent()
+		alias, _, err = p.parseIdentOrText()
		if err != nil {
			return nil, err
		}
```

**Step 4: Run test to verify it passes.**

```bash
go test ./mysql/parser/ -run TestSelectAliasStringLiteral -v
```
Expected: PASS (all four cases).

**Step 5: Run full parser test suite.**

```bash
go test ./mysql/parser/... 2>&1 | tail -20
```
Expected: all PASS.

**Step 6: Commit.**

```bash
git add mysql/parser/select.go mysql/parser/alias_test.go
git commit -m "fix(mysql/parser): accept string literal as SELECT alias"
```

---

### Task 2.3: Use `parseIdentOrText` for table alias (FROM)

**Files:**
- Modify: `mysql/parser/name.go:795-807` (`parseTableRefWithAlias`)
- Test: extend `mysql/parser/alias_test.go`

**Step 1: Add failing test cases.**

Append to `alias_test.go`:

```go
func TestTableAliasStringLiteral(t *testing.T) {
	cases := []string{
		`SELECT * FROM t AS 'a'`,
		`SELECT * FROM t 'a'`,
		`SELECT * FROM t AS "a" JOIN u AS 'b' ON t.id = b.id`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}
```

**Step 2: Run — expect FAIL.**

```bash
go test ./mysql/parser/ -run TestTableAliasStringLiteral -v
```

**Step 3: Update `parseTableRefWithAlias`.**

In `mysql/parser/name.go`:

```go
	// Optional AS alias
	if _, ok := p.match(kwAS); ok {
-		alias, _, err := p.parseIdentifier()
+		alias, _, err := p.parseIdentOrText()
		if err != nil {
			return nil, err
		}
		ref.Alias = alias
		ref.Loc.End = p.pos()
-	} else if p.isIdentToken() {
+	} else if p.isIdentToken() || p.cur.Type == tokSCONST {
		// Alias without AS keyword — accepts non-reserved keywords and string
		// literals, matching MySQL's opt_table_alias: [AS] ident_or_text
-		alias, _, _ := p.parseIdentifier()
+		alias, _, _ := p.parseIdentOrText()
		ref.Alias = alias
		ref.Loc.End = p.pos()
	}
```

**Step 4: Run test — expect PASS.**

```bash
go test ./mysql/parser/ -run TestTableAliasStringLiteral -v
go test ./mysql/parser/... 2>&1 | tail -10
```

**Step 5: Commit.**

```bash
git add mysql/parser/name.go mysql/parser/alias_test.go
git commit -m "fix(mysql/parser): accept string literal as table alias"
```

---

### Task 2.4: Check JSON_TABLE alias and other alias sites

**Files:**
- Search: `mysql/parser/select.go` for `AS alias` patterns.

**Step 1: Enumerate remaining alias sites.**

```bash
grep -n "parseIdent\b\|parseIdentifier\b" mysql/parser/select.go mysql/parser/update_delete.go mysql/parser/insert.go | grep -v "^[^:]*_test.go"
```

Review each hit. For any that corresponds to an alias position (select item alias, table alias, JSON_TABLE alias, CTE query alias — CTE uses `ident` only per MySQL spec, skip), swap to `parseIdentOrText`.

**Step 2: Add an integration test.**

Append to `alias_test.go`:

```go
func TestJsonTableAliasStringLiteral(t *testing.T) {
	sql := `SELECT * FROM JSON_TABLE('[]', '$' COLUMNS (v INT PATH '$')) AS 'j'`
	if _, err := Parse(sql); err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
}
```

**Step 3: Run + fix until green.**

```bash
go test ./mysql/parser/ -run TestJsonTableAliasStringLiteral -v
```
If FAIL, update the JSON_TABLE alias site in `parseJsonTable`.

**Step 4: Run full parser suite.**

```bash
go test ./mysql/parser/... 2>&1 | tail -10
```
Expected: all PASS.

**Step 5: Commit.**

```bash
git add -u mysql/parser/
git commit -m "fix(mysql/parser): accept string literal alias in remaining positions"
```

---

## Phase 3: Scaffold `mysql/validate` package

### Task 3.1: Create package skeleton

**Files:**
- Create: `mysql/validate/validate.go`
- Create: `mysql/validate/diagnostic.go`
- Create: `mysql/validate/scope.go`
- Create: `mysql/validate/validate_test.go`

**Step 1: Write `diagnostic.go`.**

```go
package validate

// Severity classifies a diagnostic.
type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
)

// Diagnostic is one semantic finding produced by Validate.
type Diagnostic struct {
	Code     string   // stable machine-readable code, e.g. "undeclared_variable"
	Message  string   // human-readable message
	Severity Severity // error or warning
	Position int      // byte offset within the original SQL text
}
```

**Step 2: Write `scope.go`.**

```go
package validate

import nodes "github.com/bytebase/omni/mysql/ast"

// scopeKind classifies a lexical scope.
type scopeKind int

const (
	scopeBlock scopeKind = iota
	scopeHandlerBody
)

// valScope mirrors MySQL's sp_pcontext layout: separate namespaces per kind.
type valScope struct {
	parent     *valScope
	kind       scopeKind
	vars       map[string]*nodes.DeclareVarStmt
	conditions map[string]*nodes.DeclareConditionStmt
	cursors    map[string]*nodes.DeclareCursorStmt
	labels     map[string]labelInfo
	isFunction bool
}

type labelKind int

const (
	labelBlock labelKind = iota
	labelLoop
)

type labelInfo struct {
	kind labelKind
	pos  int
}

func newScope(parent *valScope, kind scopeKind) *valScope {
	s := &valScope{
		parent:     parent,
		kind:       kind,
		vars:       map[string]*nodes.DeclareVarStmt{},
		conditions: map[string]*nodes.DeclareConditionStmt{},
		cursors:    map[string]*nodes.DeclareCursorStmt{},
		labels:     map[string]labelInfo{},
	}
	if parent != nil {
		s.isFunction = parent.isFunction
	}
	return s
}

func (s *valScope) lookupVar(name string) *nodes.DeclareVarStmt {
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.vars[lower(name)]; ok {
			return v
		}
	}
	return nil
}

func (s *valScope) lookupCondition(name string) *nodes.DeclareConditionStmt {
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.conditions[lower(name)]; ok {
			return v
		}
	}
	return nil
}

func (s *valScope) lookupCursor(name string) *nodes.DeclareCursorStmt {
	for cur := s; cur != nil; cur = cur.parent {
		if v, ok := cur.cursors[lower(name)]; ok {
			return v
		}
	}
	return nil
}

// lookupLabel walks the parent chain; labels are blocked at handler-body scopes
// (MySQL's label barrier). loopOnly=true limits matches to loop labels (for ITERATE).
func (s *valScope) lookupLabel(name string, loopOnly bool) (labelInfo, bool) {
	for cur := s; cur != nil; cur = cur.parent {
		if info, ok := cur.labels[lower(name)]; ok {
			if loopOnly && info.kind != labelLoop {
				return labelInfo{}, false
			}
			return info, true
		}
		if cur.kind == scopeHandlerBody {
			return labelInfo{}, false
		}
	}
	return labelInfo{}, false
}

func lower(s string) string {
	// ASCII lowering is sufficient for MySQL identifier case-folding in the
	// contexts the validator cares about. Keep local to avoid importing strings
	// in hot paths if later moved inline.
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
```

**Step 3: Write `validate.go`.**

```go
// Package validate performs static semantic checks on a MySQL AST produced by
// mysql/parser. It is the omni analogue of MySQL's sp_head::parse /
// sp_pcontext phase: grammar errors come from the parser, semantic errors
// (undeclared var/cursor/label, duplicate DECLARE, missing RETURN, etc.) come
// from here.
package validate

import nodes "github.com/bytebase/omni/mysql/ast"

// Options tunes which validators run. Reserved for future strictness toggles.
type Options struct{}

// Validate walks a parsed AST and returns all semantic diagnostics. An empty
// slice means "no issues"; a nil AST returns nil.
func Validate(list *nodes.List, _ Options) []Diagnostic {
	if list == nil {
		return nil
	}
	v := &validator{}
	for _, stmt := range list.Items {
		v.walk(stmt)
	}
	return v.diagnostics
}

type validator struct {
	scope       *valScope
	diagnostics []Diagnostic
}

func (v *validator) push(kind scopeKind) *valScope {
	s := newScope(v.scope, kind)
	v.scope = s
	return s
}

func (v *validator) pop() {
	if v.scope != nil {
		v.scope = v.scope.parent
	}
}

func (v *validator) emit(code, msg string, pos int) {
	v.diagnostics = append(v.diagnostics, Diagnostic{
		Code:     code,
		Message:  msg,
		Severity: SeverityError,
		Position: pos,
	})
}

// walk dispatches on the node type. In this scaffolding commit it only
// descends into routine bodies and compound blocks; each semantic check is
// wired up by a later task.
func (v *validator) walk(n nodes.Node) {
	switch s := n.(type) {
	case nil:
		return
	case *nodes.CreateFunctionStmt:
		v.walkRoutine(s.Body, !s.IsProcedure, s.Params)
	case *nodes.CreateTriggerStmt:
		v.walkRoutine(s.Body, false, nil)
	case *nodes.CreateEventStmt:
		// Event body, if present, is a compound statement. The field is not
		// directly named here; follow the AST definition.
		v.walkEventBody(s)
	case *nodes.BeginEndBlock:
		v.walkBeginEnd(s)
	}
}

func (v *validator) walkRoutine(body nodes.Node, isFunction bool, params []*nodes.FuncParam) {
	if body == nil {
		return
	}
	scope := v.push(scopeBlock)
	scope.isFunction = isFunction
	for _, p := range params {
		scope.vars[lower(p.Name)] = &nodes.DeclareVarStmt{
			Loc:      p.Loc,
			Names:    []string{p.Name},
			TypeName: p.TypeName,
		}
	}
	v.walk(body)
	v.pop()
}

func (v *validator) walkEventBody(_ *nodes.CreateEventStmt) {
	// Fill in once we wire up event-body walking; skeletal for now.
}

func (v *validator) walkBeginEnd(b *nodes.BeginEndBlock) {
	v.push(scopeBlock)
	for _, s := range b.Stmts {
		v.walk(s)
	}
	v.pop()
}
```

**Step 4: Write `validate_test.go` with a smoke test.**

```go
package validate

import (
	"testing"

	parser "github.com/bytebase/omni/mysql/parser"
)

func TestValidateEmpty(t *testing.T) {
	diags := Validate(nil, Options{})
	if diags != nil {
		t.Fatalf("expected nil diagnostics, got %v", diags)
	}
}

func TestValidateCleanProcedure(t *testing.T) {
	list, err := parser.Parse(`CREATE PROCEDURE p() BEGIN DECLARE x INT; SET x = 1; END`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	diags := Validate(list, Options{})
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
```

**Step 5: Build & run the tests.**

```bash
go build ./mysql/validate/...
go test ./mysql/validate/...
```
Expected: PASS (two tests, zero failures).

**Step 6: Commit.**

```bash
git add mysql/validate/
git commit -m "feat(mysql/validate): scaffold package with Diagnostic, scope, walker"
```

---

## Phase 4: Migrate semantic checks

### Strategy per check

For each S# in the inventory:

1. Write a validator test that reproduces the current parser-side rejection — expect it to fail now because the validator doesn't do it yet.
2. Implement the check in `validate.go`/`scope.go`.
3. Validator test PASSes.
4. Parser-side check stays intact for now (dual coverage) — Phase 5 removes parser side.
5. Commit.

Each S# is one task. Tasks 4.1–4.12 below.

### Task 4.1: S1–S4 duplicate DECLARE + labels

**Files:**
- Modify: `mysql/validate/validate.go` (add declare walkers)
- Modify: `mysql/validate/scope.go` (add declare helpers that emit diagnostics)
- Test: `mysql/validate/validate_test.go`

**Step 1: Write the failing tests.**

Append to `validate_test.go`:

```go
func TestValidateDuplicateDeclareVar(t *testing.T) {
	sql := `CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE x INT; SELECT 1; END`
	_, err := parser.Parse(sql)
	if err == nil {
		// Once Phase 5 removes parser-side check this will parse; validate then.
		list, _ := parser.Parse(sql)
		diags := Validate(list, Options{})
		requireCode(t, diags, "duplicate_variable")
	}
}

func TestValidateDuplicateCursor(t *testing.T) { /* analogous */ }
func TestValidateDuplicateCondition(t *testing.T) { /* analogous */ }
func TestValidateDuplicateLabel(t *testing.T) { /* analogous */ }

func requireCode(t *testing.T, diags []Diagnostic, code string) {
	t.Helper()
	for _, d := range diags {
		if d.Code == code {
			return
		}
	}
	t.Fatalf("expected diagnostic code %q, got %v", code, diags)
}
```

Note: while parser-side check is still active, `parser.Parse` returns error and validator never runs. The tests are structured so that once Phase 5 strips parser-side, they start exercising validator. For now, add a companion test that synthesizes AST directly (bypassing parser) so the validator logic is exercised today:

```go
func TestValidateDuplicateVarDirect(t *testing.T) {
	list := &nodes.List{Items: []nodes.Node{
		&nodes.CreateFunctionStmt{
			IsProcedure: true,
			Body: &nodes.BeginEndBlock{Stmts: []nodes.Node{
				&nodes.DeclareVarStmt{Names: []string{"x"}},
				&nodes.DeclareVarStmt{Names: []string{"x"}, Loc: nodes.Loc{Start: 42}},
			}},
		},
	}}
	diags := Validate(list, Options{})
	requireCode(t, diags, "duplicate_variable")
}
```

**Step 2: Run — FAIL.**

```bash
go test ./mysql/validate/ -run TestValidateDuplicateVarDirect -v
```

**Step 3: Implement in validator.**

In `validate.go`, add cases inside `walk`:

```go
	case *nodes.DeclareVarStmt:
		for _, name := range s.Names {
			key := lower(name)
			if _, exists := v.scope.vars[key]; exists {
				v.emit("duplicate_variable", "duplicate variable declaration: "+name, s.Loc.Start)
				continue
			}
			v.scope.vars[key] = s
		}
	case *nodes.DeclareConditionStmt:
		key := lower(s.Name)
		if _, exists := v.scope.conditions[key]; exists {
			v.emit("duplicate_condition", "duplicate condition declaration: "+s.Name, s.Loc.Start)
		} else {
			v.scope.conditions[key] = s
		}
	case *nodes.DeclareCursorStmt:
		key := lower(s.Name)
		if _, exists := v.scope.cursors[key]; exists {
			v.emit("duplicate_cursor", "duplicate cursor declaration: "+s.Name, s.Loc.Start)
		} else {
			v.scope.cursors[key] = s
		}
		// Select subtree not walked: cursor's SELECT isn't a routine-body child.
```

Update `walkBeginEnd` to register labels before walking body:

```go
func (v *validator) walkBeginEnd(b *nodes.BeginEndBlock) {
	if b.Label != "" && v.scope != nil {
		if _, exists := v.scope.labels[lower(b.Label)]; exists {
			v.emit("duplicate_label", "duplicate label: "+b.Label, b.Loc.Start)
		} else {
			v.scope.labels[lower(b.Label)] = labelInfo{kind: labelBlock, pos: b.Loc.Start}
		}
	}
	v.push(scopeBlock)
	for _, s := range b.Stmts {
		v.walk(s)
	}
	v.pop()
}
```

**Step 4: Run — PASS.**

```bash
go test ./mysql/validate/... -v
```

**Step 5: Commit.**

```bash
git add mysql/validate/
git commit -m "feat(mysql/validate): detect duplicate DECLARE/label"
```

---

### Task 4.2: S5–S6 HANDLER condition duplicate + undeclared condition

Write `TestValidateHandlerDuplicateCond` and `TestValidateHandlerUndeclaredCond` constructing AST directly. Implement in `walk` for `*nodes.DeclareHandlerStmt`: iterate `Conditions`, build `handlerCondKey` locally, lookup `scope.conditions` for `HandlerCondName`. Emit `duplicate_handler_condition` and `undeclared_condition`. Open a `scopeHandlerBody` around `s.Stmt` and walk it.

Commit: `feat(mysql/validate): validate HANDLER condition list`.

---

### Task 4.3: S7–S8 LEAVE / ITERATE undeclared label

Write direct-AST tests that build `*nodes.LeaveStmt{Label: "nope"}` / `*nodes.IterateStmt{Label: "nope"}` inside a compound body and assert `undeclared_label` / `undeclared_loop_label`. Register loop labels in `walkWhileStmt`/`walkRepeatStmt`/`walkLoopStmt` walkers (add them). Implement lookup via `scope.lookupLabel`.

Commit: `feat(mysql/validate): validate LEAVE/ITERATE target labels`.

---

### Task 4.4: S10 OPEN/FETCH/CLOSE undeclared cursor

Write direct-AST test for `OpenCursorStmt{Name: "nope"}`. Implement case in `walk`: `if v.scope.lookupCursor(s.Name) == nil { emit("undeclared_cursor", ...) }`. Repeat for `FetchCursorStmt` and `CloseCursorStmt`.

Commit: `feat(mysql/validate): validate cursor reference`.

---

### Task 4.5: S9 RETURN only inside function

Write test: `ReturnStmt` in a procedure body (not function) → `return_outside_function`. Implement: track `v.scope.isFunction`; emit when false.

Commit: `feat(mysql/validate): reject RETURN outside function body`.

---

### Task 4.6: S12 no RETURN in function body

Write test: `CreateFunctionStmt{IsProcedure: false, Body: BeginEndBlock{Stmts: []Node{SelectStmt{...}}}}` → `function_missing_return`. Implement: after walking function body, check `containsReturn(body)` (copy the helper from `mysql/parser/create_function.go` into the validate package). Emit if false.

Commit: `feat(mysql/validate): require RETURN in function body`.

---

### Task 4.7: S11 SET undeclared variable (deferred; revisits in Phase 6)

Write test for an obviously-bogus name: `SET nonexistent_xyz123 = 1` in routine body → `undeclared_variable`. Implement walker case for `*nodes.Assignment` inside routine scope:

```go
case *nodes.Assignment:
    col := s.Column
    if col == nil { return }
    // Skip @user and @@sys variables (parser wraps them with prefix).
    if strings.HasPrefix(col.Column, "@") { return }
    // Skip qualified (schema.var) — those are system vars by MySQL rules.
    if col.Table != "" { return }
    if v.scope == nil { return }
    if v.scope.lookupVar(col.Column) != nil { return }
    // Bug 2 fix wired up in Phase 6 — placeholder for now: keep rejecting.
    v.emit("undeclared_variable", "undeclared variable: "+col.Column, col.Loc.Start)
```

Commit: `feat(mysql/validate): validate SET target in routine body`.

---

## Phase 5: Remove semantic checks from parser

### Task 5.1: Delete parser-side duplicate-DECLARE errors (S1–S4)

**Files:**
- Modify: `mysql/parser/scope.go:140-201`
- Keep `declareVar`/`declareCondition`/`declareCursor`/`declareLabel` signatures but make them no-op (or delete — see below).

**Step 1: Make decision — delete vs neuter.**

Simpler path: delete the `p.procScope != nil` branches and the entire `scope.go` file. But since `pushScope` / `popScope` / parameter registration are still invoked in `parseCreateFunctionStmt` at line 202–213, we have to keep a minimal stub or delete all of that too. Decision: delete. Parser no longer pushes scope.

```bash
rm mysql/parser/scope.go
```

**Step 2: Remove all parser call sites.**

Grep and remove:
```bash
grep -rn "procScope\|pushScope\|popScope\|lookupVar\|lookupCondition\|lookupCursor\|lookupLabel\|declareVar\|declareCondition\|declareCursor\|declareLabel" mysql/parser/*.go | grep -v "_test.go"
```

For each hit, delete the containing `if p.procScope != nil { ... }` block. The remaining site-by-site removals are Tasks 5.2–5.6.

Also remove the `procScope` field from `Parser` in `parser.go`.

**Step 3: Expect build failure until all Tasks 5.x land.** Continue with Task 5.2.

---

### Task 5.2: Remove parser-side HANDLER condition checks (S5–S6)

In `mysql/parser/compound.go:480-513`, strip the `seen` dedup loop and the `lookupCondition` block. Keep only the parse of the condition value list.

Run `go build ./mysql/parser/...`. Iterate until clean.

Commit: `refactor(mysql/parser): strip HANDLER condition semantic checks`.

---

### Task 5.3: Remove parser-side LEAVE/ITERATE checks (S7–S8)

In `compound.go:1017-1067`, drop the `p.procScope != nil` + `lookupLabel` blocks. Keep the raw `parseLabelIdent` call.

Commit: `refactor(mysql/parser): strip LEAVE/ITERATE label checks`.

---

### Task 5.4: Remove parser-side cursor reference check (S10)

In `compound.go:1120-1134` (`requireCursor`), delete the function and its call sites inside `parseOpenCursorStmt`, `parseFetchCursorStmt`, `parseCloseCursorStmt`.

Commit: `refactor(mysql/parser): strip cursor-reference semantic check`.

---

### Task 5.5: Remove parser-side RETURN checks (S9, S12)

Delete the `p.procScope != nil && !p.procScope.isFunction` block in `parseReturnStmt` (`compound.go:1077-1083`). Delete the `containsReturn`-based block in `create_function.go:229-236` and its helper if unused elsewhere.

Commit: `refactor(mysql/parser): strip RETURN reachability checks`.

---

### Task 5.6: Remove parser-side SET undeclared-variable check (S11)

Delete the `p.procScope != nil && col.Table == "" && !systemScope` block in `set_show.go:355-367`. The `systemScope` parameter on `parseSetAssignment` may become unused — leave the signature alone for now (caller still passes it, harmless).

Commit: `refactor(mysql/parser): strip SET undeclared-variable check`.

---

### Task 5.7: Rip out `procScope`, `scopeKind`, all scope plumbing

**Files:**
- Modify: `mysql/parser/parser.go` (drop `procScope` field)
- Modify: `mysql/parser/compound.go:30-43,84,515-520` (drop pushScope/popScope calls)
- Modify: `mysql/parser/create_function.go:199-213` (drop scope push + param registration)
- Modify: `mysql/parser/trigger.go:130,253` (drop pushScope calls)

**Step 1: `grep` every remaining `procScope` or `pushScope` reference.**

```bash
grep -rn "procScope\|pushScope\|popScope\|scopeKind\|scopeBlock\|scopeHandlerBody" mysql/parser/*.go
```
Expected: zero hits after edits.

**Step 2: Build.**

```bash
go build ./mysql/parser/...
```
Expected: clean.

**Step 3: Run parser tests (expect red — semantic tests now pass through).**

```bash
go test ./mysql/parser/... 2>&1 | tee /tmp/parser_test_out.txt
```
Expected FAIL in: `semantic_test.go`, `routine_body_audit_test.go`, any test asserting a semantic-category error text.

**Step 4: Migrate failing tests.**

Every test under `mysql/parser/` that asserted a semantic error (search for `undeclared|duplicate .* declaration|RETURN is only|no RETURN found|duplicate condition value|duplicate label`) belongs to the validator now. For each:
- Remove the assertion from the parser test file (the parse should now succeed).
- Add the equivalent assertion in `mysql/validate/validate_test.go` (parse then validate, check code).

Do this one file at a time, commit per file. Suggested order: `semantic_test.go` first, then `routine_body_audit_test.go`, then any stragglers.

Commit format: `test(mysql/validate): migrate <file>.go assertions from parser`.

**Step 5: Confirm green.**

```bash
go test ./mysql/parser/... ./mysql/validate/... 2>&1 | tail -30
```
Expected: all PASS.

---

## Phase 6: Bug 2 — SET sysvar fallback in validator

### Task 6.1: Collect the MySQL 8.0 system-variable list

**Files:**
- Create: `mysql/validate/sysvars.go`
- Create: `mysql/validate/sysvars_source.md` — provenance note (date dumped, source query)

**Step 1: Dump the list from a MySQL 8.0 container.**

Run against the project's usual 8.0 testcontainer:

```bash
docker run --rm mysql:8.0 mysqld --verbose --help 2>/dev/null \
  | awk '/^---.*---/{p=!p; next} p{print $1}' \
  | grep -E '^[a-z][a-z0-9_]+$' \
  | sort -u > /tmp/sysvars_raw.txt
wc -l /tmp/sysvars_raw.txt
```

Expected: ~400–600 lines.

Alternative if `--help` doesn't cooperate: connect to container and run
```sql
SELECT LOWER(VARIABLE_NAME) FROM performance_schema.global_variables
UNION
SELECT LOWER(VARIABLE_NAME) FROM performance_schema.session_variables
ORDER BY 1;
```

**Step 2: Generate `sysvars.go`.**

```go
// Code generated from MySQL 8.0.<version> — DO NOT EDIT BY HAND.
// See sysvars_source.md for the dump procedure.

package validate

var knownSystemVariables = map[string]struct{}{
	"autocommit":            {},
	"binlog_format":         {},
	"character_set_client":  {},
	// ... full list, lowercased, one per line
	"sql_mode":              {},
	"sql_safe_updates":      {},
	"time_zone":             {},
	"transaction_isolation": {},
	"unique_checks":         {},
	// ...
}

// isSystemVariable reports whether name (case-insensitive) matches a known
// MySQL 8.0 session/global system variable. Used by the SET-assignment
// validator to implement sp_head::find_variable's fallback rule.
func isSystemVariable(name string) bool {
	_, ok := knownSystemVariables[lower(name)]
	return ok
}
```

Populate the map from `/tmp/sysvars_raw.txt`. Tool suggestion: a small `go generate` script, or just paste the sorted list.

**Step 3: Write provenance note.**

`sysvars_source.md`:
```
Generated from MySQL 8.0 on <YYYY-MM-DD> via:
  docker run --rm mysql:8.0 mysqld --verbose --help | awk ...

To refresh, rerun the above, replace sysvars.go contents, and bump the
version comment at the top of the file.
```

**Step 4: Commit.**

```bash
git add mysql/validate/sysvars.go mysql/validate/sysvars_source.md
git commit -m "feat(mysql/validate): embed MySQL 8.0 system variable table"
```

---

### Task 6.2: Wire sysvar fallback into SET validator

**Files:**
- Modify: `mysql/validate/validate.go` (the `*nodes.Assignment` case from Task 4.7)
- Test: `mysql/validate/validate_test.go`

**Step 1: Write failing test.**

```go
func TestValidateSetSystemVariableInRoutine(t *testing.T) {
	cases := []string{
		`CREATE PROCEDURE p() BEGIN SET sql_safe_updates = 0; END`,
		`CREATE PROCEDURE p() BEGIN SET SQL_MODE = ''; END`,
		`CREATE PROCEDURE p() BEGIN SET autocommit = 1; END`,
		`CREATE PROCEDURE p() BEGIN SET foreign_key_checks = 0; SET unique_checks = 0; END`,
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			list, err := parser.Parse(sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			diags := Validate(list, Options{})
			for _, d := range diags {
				if d.Code == "undeclared_variable" {
					t.Fatalf("unexpected undeclared_variable: %v", d)
				}
			}
		})
	}
}

func TestValidateSetUndeclaredStillRejected(t *testing.T) {
	sql := `CREATE PROCEDURE p() BEGIN SET totally_not_a_real_thing = 1; END`
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	diags := Validate(list, Options{})
	requireCode(t, diags, "undeclared_variable")
}
```

**Step 2: Run — expect first group FAIL.**

```bash
go test ./mysql/validate/ -run 'TestValidateSet(System|Undeclared)' -v
```

**Step 3: Update validator case.**

In the `*nodes.Assignment` branch from Task 4.7:

```go
case *nodes.Assignment:
    col := s.Column
    if col == nil { return }
    if strings.HasPrefix(col.Column, "@") { return }
    if col.Table != "" { return }
    if v.scope == nil { return }
    if v.scope.lookupVar(col.Column) != nil { return }
    if isSystemVariable(col.Column) { return } // NEW: sp_head::find_variable fallback
    v.emit("undeclared_variable", "undeclared variable: "+col.Column, col.Loc.Start)
```

**Step 4: Run — PASS.**

```bash
go test ./mysql/validate/... -v 2>&1 | tail -30
```

**Step 5: Commit.**

```bash
git add mysql/validate/validate.go mysql/validate/validate_test.go
git commit -m "fix(mysql/validate): fall back to system variable before undeclared error"
```

---

## Phase 7: Caller integration + completion regression

### Task 7.1: Prove completion improvement

**Files:**
- Test: `mysql/parser/complete_test.go` (extend) or `mysql/completion/completion_test.go`

**Step 1: Add a completion test that previously failed due to Bug 2.**

```go
func TestCompletionAfterSetSysVarInRoutine(t *testing.T) {
	sql := `CREATE PROCEDURE p() BEGIN
  SET sql_safe_updates = 0;
  SELECT |
END`
	// |  marks the cursor. Replace with offset.
	cursor := strings.Index(sql, "|")
	clean := strings.Replace(sql, "|", "", 1)
	cs := parser.Collect(clean, cursor)
	if !cs.HasToken(someFromOrWhereExpectedTokenType) {
		t.Fatalf("expected SELECT-expression candidates after SET, got %v", cs.Tokens)
	}
}
```

Adapt token-type constant and API signature to whatever `complete.go` actually exports.

**Step 2: Run — PASS** (because Phase 5 already removed the early bail).

```bash
go test ./mysql/parser/ -run TestCompletionAfterSetSysVarInRoutine -v
```

**Step 3: Commit.**

```bash
git add mysql/parser/complete_test.go
git commit -m "test(mysql/parser): completion reaches past SET sys_var in routine body"
```

---

### Task 7.2: Catalog / deparse / completion sanity sweep

**Files:** read-only.

**Step 1: Run every downstream MySQL test to confirm no regression.**

```bash
go test ./mysql/... 2>&1 | tee /tmp/mysql_all.txt | tail -40
```

Expected: all PASS. If any regression, investigate — likely an AST consumer that depended on parser rejecting semantically invalid input. Fix by having the caller invoke `validate.Validate` before using the AST.

---

### Task 7.3: Confirm no external-to-mysql regressions

**Step 1:**
```bash
go test ./... 2>&1 | tail -60
```
Expected: whole-repo PASS.

**Step 2: Commit nothing (verification only).**

---

## Phase 8: Docs + cleanup

### Task 8.1: Update SKILL.md / progress docs

**Files:**
- Modify: `mysql/parser/SKILL.md` — note the split: parser is now grammar-only.
- Modify: `mysql/parser/PROGRESS.json` — drop items covering semantic checks that moved.
- Create: `mysql/validate/SKILL.md` — describe package purpose, diagnostic codes, how to add a new rule.

Commit: `docs(mysql): describe parser/validator split`.

---

### Task 8.2: Dead-code sweep

```bash
grep -rn "systemScope\|containsReturn\|isFunction" mysql/parser/ | grep -v "_test.go"
```

Delete any helper now unused. Rerun `go test ./...` to confirm green.

Commit: `chore(mysql/parser): drop unused post-split helpers`.

---

### Task 8.3: Open PR

Push branch and open a PR titled: `feat(mysql): split parser/validator; fix string-literal alias + routine-body SET sysvar`.

PR body:
- Summary: What moved, what's new, the two bugs.
- Semantic-check inventory table (copy from this plan).
- Caller impact: completion unchanged/improved, catalog unchanged, deparse unchanged, advisor/linter now has a proper `Validate` API to call.
- Test plan checklist: `go test ./mysql/...`, targeted completion test from Task 7.1, targeted alias tests from Phase 2.

---

## Rollback plan

Each phase is a clean boundary:
- After Phase 2 the branch is valuable on its own (Bug 1 fixed).
- After Phase 4 the validator exists but parser-side duplicates remain; reverting only Phase 5 restores status quo with no user impact.
- After Phase 6 Bug 2 is fixed; reverting Phase 6 restores the strict check.

Safe rollback granule: commit boundaries. No schema migrations, no external API removed (validator is additive), no config flags to clean up.

---

## Notes for the executor

- **AST immutability assumption.** The validator treats `ast.List` as read-only. It does not mutate nodes. If a future rule needs resolved references, attach them to a parallel map in `Diagnostic` or a returned `Analysis` struct — do **not** stuff them back onto AST nodes.
- **Error codes are stable API.** Once a diagnostic code is shipped, external callers may match on it. Pick names carefully in Tasks 4.x.
- **Tests use direct AST construction in Phase 4** because parser-side checks are still active until Phase 5. After Phase 5, switch to `parser.Parse(sql)`-fed tests for realism.
- **Don't split walkBeginEnd beyond necessity.** Keep scope management in one place.
- **Don't pre-optimize the sysvar lookup.** `map[string]struct{}` is fine at ~600 entries.
