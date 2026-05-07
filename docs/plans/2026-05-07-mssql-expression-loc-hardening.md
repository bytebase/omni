# MSSQL Expression Loc Hardening Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make MSSQL expression `Loc` ranges reliable enough for lossless source extraction of CHECK constraints, DEFAULT expressions, function calls, and column references.

**Architecture:** First add the MSSQL equivalent of Doris/Snowflake `ast.NodeLoc` helpers so parser code has one canonical way to read child ranges. Then add exact source-slice regression tests for the reported MSSQL failures and fix parser construction sites to use child-node starts and consumed-token ends.

**Tech Stack:** Go, `github.com/bytebase/omni/mssql/ast`, `github.com/bytebase/omni/mssql/parser`, standard `go test`.

### Task 1: Add MSSQL AST Loc Helpers

**Files:**
- Create: `mssql/ast/loc.go`
- Modify: `mssql/ast/node.go`
- Test: `mssql/ast/loc_test.go`

**Step 1: Add failing tests**

Create `mssql/ast/loc_test.go` with focused tests for:

```go
func TestLocMethods(t *testing.T) {
	if (Loc{Start: 1, End: 3}).IsValid() != true { t.Fatal("valid loc") }
	if (Loc{Start: 1, End: -1}).IsValid() != false { t.Fatal("invalid loc") }
	if got := (Loc{Start: 4, End: 8}).Merge(Loc{Start: 1, End: 5}); got != (Loc{Start: 1, End: 8}) {
		t.Fatalf("Merge = %+v", got)
	}
}

func TestNodeLocExpressions(t *testing.T) {
	cases := []struct {
		name string
		node Node
		want Loc
	}{
		{"BinaryExpr", &BinaryExpr{Loc: Loc{Start: 1, End: 2}}, Loc{Start: 1, End: 2}},
		{"FuncCallExpr", &FuncCallExpr{Loc: Loc{Start: 2, End: 5}}, Loc{Start: 2, End: 5}},
		{"ColumnRef", &ColumnRef{Loc: Loc{Start: 3, End: 4}}, Loc{Start: 3, End: 4}},
		{"LikeExpr", &LikeExpr{Loc: Loc{Start: 4, End: 9}}, Loc{Start: 4, End: 9}},
	}
	for _, tt := range cases {
		if got := NodeLoc(tt.node); got != tt.want {
			t.Fatalf("%s NodeLoc = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}
```

**Step 2: Run test to verify RED**

Run:

```bash
go test ./mssql/ast -run 'TestLocMethods|TestNodeLocExpressions' -count=1
```

Expected: compile failure because `Loc.IsValid`, `Loc.Merge`, and `NodeLoc` do not exist.

**Step 3: Implement helpers**

In `mssql/ast/node.go`, add:

```go
func (l Loc) IsValid() bool {
	return l.Start >= 0 && l.End >= 0
}

func (l Loc) Merge(other Loc) Loc {
	if !l.IsValid() { return other }
	if !other.IsValid() { return l }
	out := l
	if other.Start < out.Start { out.Start = other.Start }
	if other.End > out.End { out.End = other.End }
	return out
}
```

Create `mssql/ast/loc.go` mirroring `doris/ast/loc.go` and `snowflake/ast/loc.go`, but scoped at least to all expression nodes used by `mssql/parser/expr.go` and `mssql/parser/name.go`: `BinaryExpr`, `UnaryExpr`, `FuncCallExpr`, `CaseExpr`, `CaseWhen`, `BetweenExpr`, `InExpr`, `LikeExpr`, `IsExpr`, `ExistsExpr`, `FullTextPredicate`, `CastExpr`, `ConvertExpr`, `TryCastExpr`, `TryConvertExpr`, `CoalesceExpr`, `NullifExpr`, `IifExpr`, `ColumnRef`, `VariableRef`, `StarExpr`, `Literal`, `SubqueryExpr`, `SubqueryComparisonExpr`, `CollateExpr`, `AtTimeZoneExpr`, `ParenExpr`, `MethodCallExpr`, `GroupingSetsExpr`, `RollupExpr`, and `CubeExpr`.

Also add:

```go
func SpanNodes(nodes ...Node) Loc {
	out := NoLoc()
	for _, n := range nodes {
		if n == nil { continue }
		out = out.Merge(NodeLoc(n))
	}
	return out
}
```

**Step 4: Verify**

Run:

```bash
go test ./mssql/ast -run 'TestLocMethods|TestNodeLocExpressions' -count=1
```

Expected: pass.

### Task 2: Add Exact Expression Slice Regression Tests

**Files:**
- Modify: `mssql/parser/loc_precision_test.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Write failing SELECT-expression tests**

Add a table-driven test named `TestMSSQLLocPrecision_ExpressionRanges` that parses each SQL, unwraps the first `*ast.ResTarget`, gets `ast.NodeLoc(rt.Val)`, and compares `sql[loc.Start:loc.End]` to `want`.

Use these cases:

```go
{"SELECT c > 0", "c > 0"},
{"SELECT c IN (1, 2)", "c IN (1, 2)"},
{"SELECT c BETWEEN 1 AND 2", "c BETWEEN 1 AND 2"},
{"SELECT c LIKE 'x%'", "c LIKE 'x%'"},
{"SELECT c IS NOT NULL", "c IS NOT NULL"},
{"SELECT a + b * c", "a + b * c"},
{"SELECT a = 1 AND b = 2", "a = 1 AND b = 2"},
{"SELECT name COLLATE Latin1_General_CI_AS LIKE 'A%'", "name COLLATE Latin1_General_CI_AS LIKE 'A%'"},
{"SELECT dt AT TIME ZONE 'UTC'", "dt AT TIME ZONE 'UTC'"},
```

**Step 2: Write failing CHECK-expression tests**

Add `TestMSSQLLocPrecision_CheckConstraintExpressionRanges` with both table-level and column-level constraints:

```sql
CREATE TABLE t (c int, CONSTRAINT ck CHECK (c > 0))
CREATE TABLE t (c int CHECK (c IN (1, 2)))
```

Extract `ConstraintDef.Expr`, call `ast.NodeLoc`, and assert exact slices `c > 0` and `c IN (1, 2)`.

**Step 3: Run tests to verify RED**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_ExpressionRanges|TestMSSQLLocPrecision_CheckConstraintExpressionRanges' -count=1
```

Expected: failures showing operator-only slices like `> 0`, `IN (...)`, `AND b = 2`, `LIKE 'A%'`, and `AT TIME ZONE 'UTC'`.

### Task 3: Fix Infix And Postfix Expression Loc Construction

**Files:**
- Modify: `mssql/parser/expr.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Add small parser helper**

Add near the top of `expr.go`:

```go
func exprStart(expr nodes.ExprNode, fallback int) int {
	loc := nodes.NodeLoc(expr)
	if loc.Start >= 0 {
		return loc.Start
	}
	return fallback
}
```

**Step 2: Update binary operator layers**

In `parseOr`, `parseAnd`, `parseAddition`, and `parseMultiplication`, keep the local operator `loc` for fallback but build nodes with:

```go
Loc: nodes.Loc{Start: exprStart(left, loc), End: nodes.NodeLoc(right).End}
```

If `nodes.NodeLoc(right).End` is invalid, keep `p.prevEnd()` as fallback.

**Step 3: Update predicate/comparison nodes**

In `parseComparison`, set `start := exprStart(left, loc)` for all of:

- `IsExpr`
- `BetweenExpr`
- `InExpr` value-list form
- `InExpr` subquery form
- `LikeExpr`, including `ESCAPE`
- `SubqueryComparisonExpr`
- plain comparison `BinaryExpr`

Each outer predicate node should end at the consumed right-hand side token or closing paren: use `p.prevEnd()` after `expect(')')` or after parsing the right expression.

**Step 4: Update postfix nodes**

In `parseCollateExpr` and `parseAtTimeZoneExpr`, set `Loc.Start` from the left expression, not from the postfix keyword:

```go
Loc: nodes.Loc{Start: exprStart(expr, loc), End: p.prevEnd()}
```

**Step 5: Verify**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_ExpressionRanges|TestMSSQLLocPrecision_CheckConstraintExpressionRanges' -count=1
```

Expected: pass.

### Task 4: Add Function Call And ColumnRef Regression Tests

**Files:**
- Modify: `mssql/parser/loc_precision_test.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Write failing function call tests**

Add `TestMSSQLLocPrecision_FunctionCallRanges` using:

```sql
CREATE TABLE t (
  d datetime2 DEFAULT SYSDATETIME(),
  e int DEFAULT dbo.f(),
  c int DEFAULT COUNT(*)
)
```

Extract each column `DefaultExpr`, assert it is `*ast.FuncCallExpr`, and assert exact slices:

```go
[]string{"SYSDATETIME()", "dbo.f()", "COUNT(*)"}
```

Also include a SELECT target case with suffixes:

```sql
SELECT STRING_AGG(name, ',') WITHIN GROUP (ORDER BY name)
SELECT SUM(x) OVER (PARTITION BY y)
```

The `FuncCallExpr.Loc` should cover the suffix if the suffix is part of the same function expression.

**Step 2: Write failing column reference tests**

Add `TestMSSQLLocPrecision_ColumnRefRanges` using:

```sql
SELECT c, t.c, dbo.t.c, t.*
```

Assert exact slices for `ColumnRef`/`StarExpr`: `c`, `t.c`, `dbo.t.c`, `t.*`.

**Step 3: Run tests to verify RED**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_FunctionCallRanges|TestMSSQLLocPrecision_ColumnRefRanges' -count=1
```

Expected: failures showing `SYSDATETIME(`, `COUNT(`, schema-qualified `End=-1`, and column refs with `End=-1`.

### Task 5: Fix Function Call And Column Reference Loc Ends

**Files:**
- Modify: `mssql/parser/name.go`
- Modify: `mssql/parser/expr.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Fix simple and qualified column refs**

In `parseIdentExpr`, simple `ColumnRef` should use:

```go
Loc: nodes.Loc{Start: loc, End: p.prevEnd()}
```

In `parseQualifiedRef`, set the returned `ColumnRef` and `StarExpr` `Loc.End` to `p.prevEnd()`.

**Step 2: Fix simple function calls**

In `parseFuncCall`, capture `nameEnd := p.prevEnd()` before consuming `(`. Set:

```go
fc.Name.Loc.End = nameEnd
fc.Loc.End = p.prevEnd()
```

after every successful close paren path, including `COUNT(*)`, empty-argument calls, and normal argument calls. After parsing optional `WITHIN GROUP` or `OVER`, update `fc.Loc.End = p.prevEnd()` again so the full function expression is covered.

**Step 3: Fix schema-qualified function calls**

In `parseFuncCallWithSchema`, capture `nameEnd := p.prevEnd()` before consuming `(`. Set `fc.Name.Loc.End = nameEnd` and mirror the same `fc.Loc.End` updates as simple calls.

**Step 4: Verify**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_FunctionCallRanges|TestMSSQLLocPrecision_ColumnRefRanges' -count=1
```

Expected: pass.

### Task 6: Add A Loc Coverage Guard For Expression Nodes

**Files:**
- Modify: `mssql/parser/loc_precision_test.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Add structural guard**

Add a small guard that parses representative expressions and walks them with `ast.Walk`, failing if any expression node under the target expression has `Loc.Start < 0` or `Loc.End < Loc.Start`, except for known synthetic nodes if any are found and documented.

Representative SQL:

```sql
SELECT (a + b) * dbo.f(c) WHERE c IS NOT NULL
SELECT CASE WHEN c BETWEEN 1 AND 2 THEN SYSDATETIME() ELSE dbo.f() END
```

**Step 2: Verify**

Run:

```bash
go test ./mssql/parser -run TestMSSQLLocPrecision_ExpressionNodeLocCompleteness -count=1
```

Expected after previous fixes: pass. If it fails on unrelated AST shapes, document and narrow the guard to expression nodes touched by this plan.

### Task 7: Full Verification

**Files:**
- No changes
- Test: affected packages

**Step 1: Run all precision tests**

```bash
go test ./mssql/parser -run TestMSSQLLocPrecision -count=1
```

Expected: pass.

**Step 2: Run existing expression tests**

```bash
go test ./mssql/parser -run 'TestParseComparison|TestParseBoolean|TestParseIsNull|TestParseBetween|TestParseIn|TestParseLike|TestParseFunctions|TestParseCreateTable|TestLoc_Expressions' -count=1
```

Expected: pass. Note `TestLoc_Expressions` is in package `./mssql`, so if the exact regex does not select it from `./mssql/parser`, run `go test ./mssql -run TestLoc_Expressions -count=1`.

**Step 3: Run package verification**

```bash
go test ./mssql/ast ./mssql/parser ./mssql -count=1
```

Expected: pass.

### Risk Notes

This intentionally changes `Loc.Start` for existing expression nodes from operator-only to full-expression ranges. That is the desired contract for lossless metadata extraction, but any code depending on operator-only spans will observe different offsets. Keep the change scoped to expression nodes and exact source-slice tests; do not attempt to clean up every `End=-1` statement or DDL node in the same pass.
