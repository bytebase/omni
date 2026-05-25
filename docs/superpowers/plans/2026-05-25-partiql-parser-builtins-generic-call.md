# PartiQL parser-builtins-generic-call (15a) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the deferred-feature ParseError in `partiql/parser/parser.go:parseVarRef` (lines 253-258) with a generic `*ast.FuncCall` builder so that any `IDENT(args)` call parses cleanly. This unblocks DynamoDB-native predicates (`attribute_exists`, `attribute_not_exists`, `attribute_type`, `begins_with`, `contains`, plus user-defined function calls) that currently surface as syntax errors in the bytebase SQL editor after the cutover.

**Architecture:** Two-file change in `partiql/parser/`. `parseVarRef` gets a new function-call branch that calls a new `parseFuncCallArgs` helper. The helper parses `(expr (COMMA expr)*)?` matching ANTLR rule `FunctionCallIdent` (PartiQLParser.g4:615). `@foo(x)` is rejected to match the legacy grammar. AST has no change — `*ast.FuncCall` (exprs.go:336) already exists with the right shape. Tests flip the existing `funcall_stub` negative case, add positive cases for the shapes that customers send, and lock in the new negative cases (`@foo(x)`, trailing comma, unclosed paren).

**Tech Stack:** Go, hand-written recursive-descent parser, standard `testing` package, `ast.NodeToString` for AST shape comparison.

**Spec:** `docs/superpowers/specs/2026-05-25-partiql-parser-builtins-generic-call-design.md`

---

### Task 1: Add a single failing positive test for `foo(x)`

**Files:**
- Modify: `partiql/parser/parser_test.go` (add a new `TestParser_FuncCall` test function at the end of the file)

- [ ] **Step 1: Write the failing test**

Append the following to `partiql/parser/parser_test.go`:

```go
// TestParser_FuncCall verifies that generic IDENT(args) function calls
// produce *ast.FuncCall nodes per DAG node 15a (parser-builtins-generic-call).
// Compared against ast.NodeToString output for compact shape assertion.
func TestParser_FuncCall(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single_arg",
			input: "foo(x)",
			want:  "FuncCall{Name:foo Args:[VarRef{Name:x}]}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			if got != tc.want {
				t.Errorf("AST mismatch\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./partiql/parser/ -run TestParser_FuncCall -v`
Expected: FAIL with `unexpected parse error: function call "foo" is deferred to parser-builtins (DAG node 15)`

- [ ] **Step 3: Do not commit yet**

The failing test stays uncommitted; Task 2 lands the implementation + flips the obsolete `funcall_stub` case in the same commit so the whole suite stays green at commit boundaries.

---

### Task 2: Implement `parseFuncCallArgs` and update `parseVarRef`

**Files:**
- Modify: `partiql/parser/parser.go` (replace the function-call ParseError in `parseVarRef` at lines 250-258; add `parseFuncCallArgs` helper)
- Modify: `partiql/parser/parser_test.go` (remove the obsolete `funcall_stub` case at lines 667-671)

- [ ] **Step 1: Modify `parseVarRef` in `partiql/parser/parser.go`**

Replace the existing block:

```go
	// Function call lookahead: <name> ( ... )
	// Applies to BOTH the plain `foo(...)` and `@foo(...)` forms (the
	// latter is unusual but grammar-legal).
	if p.cur.Type == tokPAREN_LEFT {
		return nil, &ParseError{
			Message: fmt.Sprintf("function call %q is deferred to parser-builtins (DAG node 15)", name),
			Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
		}
	}
```

with:

```go
	// Function call: name=symbolPrimitive PAREN_LEFT (expr (COMMA expr)*)? PAREN_RIGHT
	// (PartiQLParser.g4:615 — FunctionCallIdent). Implements DAG node 15a.
	// @-prefix is NOT permitted before a function call: ANTLR's functionCall
	// rule names symbolPrimitive directly, only varRefExpr admits @.
	if p.cur.Type == tokPAREN_LEFT {
		if atPrefixed {
			return nil, &ParseError{
				Message: "@-prefix is not allowed before a function call",
				Loc:     ast.Loc{Start: start, End: p.cur.Loc.End},
			}
		}
		args, endOff, err := p.parseFuncCallArgs()
		if err != nil {
			return nil, err
		}
		return &ast.FuncCall{
			Name: name,
			Args: args,
			Loc:  ast.Loc{Start: start, End: endOff},
		}, nil
	}
```

Note: the unused `caseSensitive` local from `parseSymbolPrimitive` is already discarded for the VarRef path (it sets `CaseSensitive: caseSensitive` on the VarRef return). It is intentionally dropped for FuncCall — see spec §5 (Out of Scope).

- [ ] **Step 2: Add the `parseFuncCallArgs` helper to `partiql/parser/parser.go`**

Add immediately after `parseVarRef` (between `parseVarRef` and `parseType`):

```go
// parseFuncCallArgs consumes the parenthesized argument list of a function
// call: PAREN_LEFT (expr (COMMA expr)*)? PAREN_RIGHT. The opening paren
// must be the current token. Returns the parsed args (nil if empty), the
// 0-based byte End offset of PAREN_RIGHT, and any parse error.
//
// Used by DAG node 15a (parser-builtins-generic-call) and will be reused
// by 15b for the typed builtins whose argument list is also comma-separated
// (DATE_ADD, DATE_DIFF, COALESCE, NULLIF, SIZE, EXISTS, etc.).
func (p *Parser) parseFuncCallArgs() ([]ast.ExprNode, int, error) {
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, 0, err
	}
	var args []ast.ExprNode
	if p.cur.Type != tokPAREN_RIGHT {
		for {
			arg, err := p.parseExprTop()
			if err != nil {
				return nil, 0, err
			}
			args = append(args, arg)
			if p.cur.Type != tokCOMMA {
				break
			}
			p.advance() // consume COMMA
		}
	}
	endOff := p.cur.Loc.End
	if _, err := p.expect(tokPAREN_RIGHT); err != nil {
		return nil, 0, err
	}
	return args, endOff, nil
}
```

- [ ] **Step 3: Remove the obsolete `funcall_stub` case from `parser_test.go`**

Delete lines 667-671 (the test entry):

```go
		{
			name:      "funcall_stub",
			input:     "foo(x)",
			wantErrIn: `function call "foo" is deferred to parser-builtins (DAG node 15)`,
		},
```

Take care to remove the trailing comma on the previous entry only if `funcall_stub` was followed by another stub — check the surrounding context after the edit.

- [ ] **Step 4: Run the parser test suite**

Run: `go test ./partiql/parser/...`
Expected: PASS (TestParser_FuncCall now passes, no other tests broken)

- [ ] **Step 5: Do not commit yet**

Continue to Task 3 to expand positive coverage.

---

### Task 3: Add the remaining positive cases

**Files:**
- Modify: `partiql/parser/parser_test.go` (extend the `TestParser_FuncCall` cases slice from Task 1)

- [ ] **Step 1: Extend `TestParser_FuncCall` with 5 more positive cases**

Edit the `cases` slice in `TestParser_FuncCall` to read in full:

```go
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single_arg",
			input: "foo(x)",
			want:  "FuncCall{Name:foo Args:[VarRef{Name:x}]}",
		},
		{
			name:  "zero_args",
			input: "foo()",
			want:  "FuncCall{Name:foo Args:[]}",
		},
		{
			name:  "multi_arg_literals",
			input: "foo(1, 2, 3)",
			want:  "FuncCall{Name:foo Args:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}]}",
		},
		{
			name:  "nested_calls",
			input: "f(g(x))",
			want:  "FuncCall{Name:f Args:[FuncCall{Name:g Args:[VarRef{Name:x}]}]}",
		},
		{
			name:  "call_with_path_step",
			input: "f(x).bar",
			want:  "PathExpr{Root:FuncCall{Name:f Args:[VarRef{Name:x}]} Steps:[DotStep{Field:bar}]}",
		},
		{
			name:  "quoted_arg_name",
			input: `attribute_exists("a")`,
			want:  "FuncCall{Name:attribute_exists Args:[VarRef{Name:a CaseSensitive:true}]}",
		},
		{
			name:  "dynamodb_begins_with",
			input: `begins_with(addr, '7834')`,
			want:  `FuncCall{Name:begins_with Args:[VarRef{Name:addr} StringLit{Val:"7834"}]}`,
		},
	}
```

- [ ] **Step 2: Run the test, verify all positive cases pass**

Run: `go test ./partiql/parser/ -run TestParser_FuncCall -v`
Expected: PASS for all 7 sub-tests.

The expected strings were derived from `partiql/ast/outfuncs.go`'s case clauses: `FuncCall{Name:%s Args:[...]}`, `VarRef{Name:%s [CaseSensitive:true]}`, `NumberLit{Val:%s}`, `StringLit{Val:%q}` (Go `%q` quotes the value), `PathExpr{Root:... Steps:[...]}`, `DotStep{Field:%s}`. If a sub-test still fails on a rendering mismatch (e.g., because a future `outfuncs.go` edit changed a label), the test report shows `got:` and `want:` — update the `want:` string to match the `got:` value verbatim. This is shape-locking, not semantic verification.

- [ ] **Step 3: Do not commit yet**

Continue to Task 4 for negative cases.

---

### Task 4: Add the negative cases

**Files:**
- Modify: `partiql/parser/parser_test.go` (extend `TestParser_Errors` cases)

- [ ] **Step 1: Add 4 negative cases to `TestParser_Errors`**

Inside the `cases` slice in `TestParser_Errors`, in the `// --- Real syntax errors ---` section (after line 683 in the pre-edit file), add:

```go
		{
			name:      "funcall_at_prefix",
			input:     "@foo(x)",
			wantErrIn: "@-prefix is not allowed before a function call",
		},
		{
			name:      "funcall_trailing_comma",
			input:     "foo(x,)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "funcall_unclosed",
			input:     "foo(x",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "funcall_leading_comma",
			input:     "foo(,)",
			wantErrIn: "unexpected token",
		},
```

- [ ] **Step 2: Run the test, verify negative cases pass**

Run: `go test ./partiql/parser/ -run TestParser_Errors -v`
Expected: PASS for all sub-tests, including the four new ones.

If `funcall_trailing_comma` or `funcall_leading_comma` emit a different message than `"unexpected token"` (because `parseExprTop` may produce a more specific error like `"expected PAREN_RIGHT"` for the trailing case once it tries to parse `)` as an expression), update the `wantErrIn:` string to match the actual error substring. The test framework uses `strings.Contains` so a partial substring is fine.

- [ ] **Step 3: Do not commit yet**

Continue to Task 5 for the broader regression net.

---

### Task 5: Broader regression verification

**Files:** None (verification only)

- [ ] **Step 1: Run the full partiql package test suite**

Run: `go test ./partiql/...`
Expected: PASS for all packages (parser, analysis, completion, ast, catalog).

If any test in `partiql/analysis` or `partiql/completion` now fails because they see `*ast.FuncCall` for the first time in a path that previously errored, inspect the failure and either:
- Update the downstream code to handle `*ast.FuncCall` properly (if the failure is a real bug in those packages), or
- Note the failure as out-of-scope work in the commit message (only if the downstream package's logic was already broken for FuncCall but masked by the parser stub).

- [ ] **Step 2: Verify the AWS corpus summary shift**

Run: `go test ./partiql/parser/ -run TestParser_AWSCorpus -v 2>&1 | grep "AWS corpus:"`
Expected output line: `AWS corpus: N fully parsed, M stubbed, K skipped`

Compare against the pre-change baseline (`74 fully parsed, 14 stubbed, 2 skipped` or whatever the current numbers are). After 15a, `N` should be 3 higher and `M` should be 3 lower (for `function-begins-with-001`, `function-contains-001`, `function-attribute-type-001`).

If the shift is more than 3, that means additional corpus entries also became parseable as a side effect — that's still fine, just note it in the commit message.

If the shift is less than 3, investigate which of the three expected entries did NOT flip and why (likely a parse error elsewhere in the SQL, not in the function call itself).

- [ ] **Step 3: Do not commit yet**

Continue to Task 6 for the commit.

---

### Task 6: Commit the change

**Files:** None (commit only)

- [ ] **Step 1: Review the staged diff**

Run: `git diff partiql/parser/parser.go partiql/parser/parser_test.go`
Confirm: exactly two files changed, no unrelated edits.

- [ ] **Step 2: Stage and commit**

Run:

```bash
git add partiql/parser/parser.go partiql/parser/parser_test.go
git commit -m "$(cat <<'EOF'
feat(partiql): T15a generic IDENT(args) function calls

Replaces the deferred-feature ParseError in parser.go:parseVarRef with
a *ast.FuncCall builder so any IDENT(args) call parses cleanly.

Fixes a post-cutover regression in the bytebase SQL editor where every
DynamoDB-native predicate (attribute_exists / attribute_not_exists /
attribute_type / begins_with / contains, plus any user-defined function
call) surfaced as a syntax error.

@foo(x) is now rejected with a dedicated error to match the legacy
ANTLR grammar (functionCall's `name=symbolPrimitive` rule does not
admit @). The AST shape is unchanged — *ast.FuncCall already existed
with the right fields.

Out of scope: typed keyword builtins (CAST / CASE / SUBSTRING / TRIM /
EXTRACT / COALESCE / NULLIF / DATE_ADD / DATE_DIFF / CHAR_LENGTH /
UPPER / LOWER / SIZE / EXISTS / LIST() / SEXP()) — those are 15b's
keyword-token dispatch and remain stubbed. Aggregates and window
functions are nodes 14 / 13.

Spec: docs/superpowers/specs/2026-05-25-partiql-parser-builtins-generic-call-design.md

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 3: Verify the commit lands**

Run: `git log --oneline -1`
Expected: a single new commit titled `feat(partiql): T15a generic IDENT(args) function calls`.

- [ ] **Step 4: Push (optional, ask user)**

Do NOT push without explicit user confirmation. Ask: "Push to origin/main, or open a PR branch first?"

---

## Out-of-Plan Notes

- **DAG status update**: After this PR lands, mark node 15a status as `done (PR #...)` in `docs/migration/partiql/dag.md` line 41. This is a follow-up doc commit, not part of this plan.

- **Memory updates**: None needed. The work is fully in-tree and committed.
