# MSSQL Spatial Index And Expression Loc Fixes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Preserve MSSQL spatial index nested option values and make expression `Loc` ranges cover the complete parsed expression.

**Architecture:** Keep the fixes in the MSSQL parser layer, where token source spans are known. Add precise regression tests first, then adjust parser location construction and nested option value consumption without changing AST shapes.

**Tech Stack:** Go, `github.com/bytebase/omni/mssql/parser`, `github.com/bytebase/omni/mssql/ast`, standard `go test`.

### Task 1: Add Spatial Option Regression Tests

**Files:**
- Create: `mssql/parser/loc_precision_test.go`
- Modify: none
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Write failing tests**

Add a parser test that parses:

```sql
CREATE SPATIAL INDEX SI_t ON dbo.t (geom) USING GEOMETRY_GRID
WITH (
  BOUNDING_BOX = (0, 0, 100, 100),
  GRIDS = (LEVEL_1 = LOW, LEVEL_2 = MEDIUM, LEVEL_3 = HIGH, LEVEL_4 = HIGH),
  CELLS_PER_OBJECT = 16
)
```

Assert `CreateSpatialIndexStmt.Options.Items` contains exactly:

```go
[]string{
  "BOUNDING_BOX=(0, 0, 100, 100)",
  "GRIDS=(LEVEL_1 = LOW, LEVEL_2 = MEDIUM, LEVEL_3 = HIGH, LEVEL_4 = HIGH)",
  "CELLS_PER_OBJECT=16",
}
```

**Step 2: Run test to verify failure**

Run:

```bash
go test ./mssql/parser -run TestMSSQLLocPrecision_SpatialIndexNestedOptions -count=1
```

Expected before fix: failure showing `BOUNDING_BOX=` and `GRIDS=`.

### Task 2: Preserve Parenthesized Index Option Values

**Files:**
- Modify: `mssql/parser/alter_objects.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Add a local helper**

Add a helper near `parseAlterIndexOptions`:

```go
func (p *Parser) consumeBalancedValueText() (string, error) {
	start := p.pos()
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
			p.advance()
			if depth == 0 {
				return strings.TrimSpace(p.source[start:p.prevEnd()]), nil
			}
			continue
		}
		p.advance()
	}
	return strings.TrimSpace(p.source[start:p.prevEnd()]), nil
}
```

**Step 2: Use it only for value parentheses**

In the `else if p.cur.Type == '('` branch under `option_name = value`, replace the current “skip nested parens” loop with:

```go
var err error
val, err = p.consumeBalancedValueText()
if err != nil {
	return nil, err
}
```

Keep simple scalar values normalized as today (`FILLFACTOR=80`, `ALLOW_ROW_LOCKS=ON`) to avoid unnecessary golden churn.

**Step 3: Verify**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_SpatialIndexNestedOptions|TestParseAlterIndexOptionsDepth' -count=1
```

Expected after fix: both tests pass.

### Task 3: Add Expression Loc Regression Tests

**Files:**
- Modify: `mssql/parser/loc_precision_test.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Write failing tests**

Add table-driven tests for these target expressions and exact source slices:

```go
[]struct{
	sql string
	want string
}{
	{"SELECT c > 0", "c > 0"},
	{"SELECT c IN (1, 2)", "c IN (1, 2)"},
	{"SELECT c BETWEEN 1 AND 2", "c BETWEEN 1 AND 2"},
	{"SELECT c LIKE 'x%'", "c LIKE 'x%'"},
	{"SELECT c IS NOT NULL", "c IS NOT NULL"},
	{"SELECT a + b * c", "a + b * c"},
	{"SELECT a = 1 AND b = 2", "a = 1 AND b = 2"},
}
```

Unwrap the first `ResTarget` and assert the expression node `Loc` slice equals `want`.

**Step 2: Run test to verify failure**

Run:

```bash
go test ./mssql/parser -run TestMSSQLLocPrecision_ExpressionRanges -count=1
```

Expected before fix: comparison/predicate/binary expressions start at the operator.

### Task 4: Fix Expression Start Ranges

**Files:**
- Modify: `mssql/parser/expr.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Add parser-side Loc helpers**

Add helpers in `expr.go`:

```go
func exprLocStart(expr nodes.ExprNode, fallback int) int {
	if loc, ok := exprLoc(expr); ok && loc.Start >= 0 {
		return loc.Start
	}
	return fallback
}

func exprLoc(expr nodes.ExprNode) (nodes.Loc, bool) {
	switch e := expr.(type) {
	case *nodes.BinaryExpr:
		return e.Loc, true
	case *nodes.UnaryExpr:
		return e.Loc, true
	case *nodes.FuncCallExpr:
		return e.Loc, true
	case *nodes.CaseExpr:
		return e.Loc, true
	case *nodes.BetweenExpr:
		return e.Loc, true
	case *nodes.InExpr:
		return e.Loc, true
	case *nodes.LikeExpr:
		return e.Loc, true
	case *nodes.IsExpr:
		return e.Loc, true
	case *nodes.ExistsExpr:
		return e.Loc, true
	case *nodes.CastExpr:
		return e.Loc, true
	case *nodes.ConvertExpr:
		return e.Loc, true
	case *nodes.TryCastExpr:
		return e.Loc, true
	case *nodes.TryConvertExpr:
		return e.Loc, true
	case *nodes.CoalesceExpr:
		return e.Loc, true
	case *nodes.NullifExpr:
		return e.Loc, true
	case *nodes.IifExpr:
		return e.Loc, true
	case *nodes.ColumnRef:
		return e.Loc, true
	case *nodes.VariableRef:
		return e.Loc, true
	case *nodes.StarExpr:
		return e.Loc, true
	case *nodes.Literal:
		return e.Loc, true
	case *nodes.SubqueryExpr:
		return e.Loc, true
	case *nodes.SubqueryComparisonExpr:
		return e.Loc, true
	case *nodes.CollateExpr:
		return e.Loc, true
	case *nodes.AtTimeZoneExpr:
		return e.Loc, true
	case *nodes.ParenExpr:
		return e.Loc, true
	}
	return nodes.Loc{}, false
}
```

**Step 2: Update infix and predicate nodes**

In `parseOr`, `parseAnd`, `parseAddition`, `parseMultiplication`, and `parseComparison`, construct expression nodes with:

```go
Loc: nodes.Loc{Start: exprLocStart(left, loc), End: p.prevEnd()}
```

Use the same start logic for `IsExpr`, `BetweenExpr`, `InExpr`, `LikeExpr`, `BinaryExpr`, and `SubqueryComparisonExpr`.

**Step 3: Update postfix nodes**

In `parseCollateExpr` and `parseAtTimeZoneExpr`, use the left expression start:

```go
start := exprLocStart(expr, loc)
```

Then set the node `Loc.Start` to `start`.

**Step 4: Verify**

Run:

```bash
go test ./mssql/parser -run TestMSSQLLocPrecision_ExpressionRanges -count=1
```

Expected after fix: exact source slices match all cases.

### Task 5: Add Function And ColumnRef Loc Regression Tests

**Files:**
- Modify: `mssql/parser/loc_precision_test.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Write failing tests**

Add tests for:

```sql
CREATE TABLE t (
  d datetime2 DEFAULT SYSDATETIME(),
  e int DEFAULT dbo.f()
)
```

Assert both default `FuncCallExpr.Loc` slices are `SYSDATETIME()` and `dbo.f()`.

Add tests for:

```sql
SELECT c, t.c, dbo.t.c
```

Assert every `ColumnRef.Loc.End` is non-negative and slices are `c`, `t.c`, and `dbo.t.c`.

**Step 2: Run test to verify failure**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_FunctionCallRanges|TestMSSQLLocPrecision_ColumnRefRanges' -count=1
```

Expected before fix: simple function misses `)`, schema-qualified function has `End=-1`, and column refs have `End=-1`.

### Task 6: Fix Function Call And ColumnRef End Ranges

**Files:**
- Modify: `mssql/parser/name.go`
- Modify: `mssql/parser/expr.go`
- Test: `mssql/parser/loc_precision_test.go`

**Step 1: Fix simple `ColumnRef`**

In `parseIdentExpr`, return:

```go
return &nodes.ColumnRef{
	Column: name,
	Loc:    nodes.Loc{Start: loc, End: p.prevEnd()},
}, nil
```

**Step 2: Fix qualified `ColumnRef` and `StarExpr`**

In `parseQualifiedRef`, set:

```go
Loc: nodes.Loc{Start: loc, End: p.prevEnd()}
```

on the returned `ColumnRef` and qualified `StarExpr`.

**Step 3: Fix simple function call**

In `parseFuncCall`, preserve `nameEnd := p.prevEnd()` before consuming `(`, set `fc.Name.Loc.End = nameEnd`, and after every successful `expect(')')` or closing-paren `advance`, set:

```go
fc.Loc.End = p.prevEnd()
```

After parsing optional `WITHIN GROUP` or `OVER`, set `fc.Loc.End = p.prevEnd()` again so the node covers the full function expression including suffixes.

**Step 4: Fix schema-qualified function call**

In `parseFuncCallWithSchema`, preserve `nameEnd := p.prevEnd()` before consuming `(`, set both `fc.Name.Loc.End` and `fc.Loc.End` using the same rules as simple calls.

**Step 5: Verify**

Run:

```bash
go test ./mssql/parser -run 'TestMSSQLLocPrecision_FunctionCallRanges|TestMSSQLLocPrecision_ColumnRefRanges' -count=1
```

Expected after fix: exact slices match.

### Task 7: Run Focused And Package Verification

**Files:**
- No changes
- Test: full MSSQL parser package

**Step 1: Run all precision regressions**

```bash
go test ./mssql/parser -run TestMSSQLLocPrecision -count=1
```

Expected: pass.

**Step 2: Run existing affected tests**

```bash
go test ./mssql/parser -run 'TestParseAlterIndexOptionsDepth|TestParseCreateIndexStatements|TestParseFunctions|TestParseComparison|TestParseBoolean|TestParseIsNull|TestParseBetween|TestParseIn|TestParseLike' -count=1
```

Expected: pass.

**Step 3: Run package tests**

```bash
go test ./mssql/parser ./mssql -count=1
```

Expected: pass.

### Risk Notes

The main compatibility risk is `ast.String` option serialization. Keep scalar index options exactly normalized as before and only fill previously empty parenthesized values. The main Loc risk is code that accidentally depended on operator-only ranges; that behavior is incorrect for lossless extraction, and focused Loc tests should make the intended contract explicit.
