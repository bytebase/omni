package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// mustParseSetOp parses input and returns the first statement as *ast.SetOpStmt.
func mustParseSetOp(t *testing.T, input string) *ast.SetOpStmt {
	t.Helper()
	file, errs := Parse(input)
	if len(errs) > 0 {
		t.Fatalf("Parse(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("Parse(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("Parse(%q) got %T, want *ast.SetOpStmt", input, file.Stmts[0])
	}
	return stmt
}

// assertSelectItems is a helper that asserts a node is a *ast.SelectStmt with
// the given number of items and that the first item has the expected expression
// column name or literal value.
func assertIsSelectStmt(t *testing.T, n ast.Node, label string) *ast.SelectStmt {
	t.Helper()
	sel, ok := n.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("%s: got %T, want *ast.SelectStmt", label, n)
	}
	return sel
}

// ---------------------------------------------------------------------------
// UNION
// ---------------------------------------------------------------------------

func TestUnionBasic(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 UNION SELECT a FROM t2")

	if stmt.Op != ast.SetUnion {
		t.Errorf("Op = %v, want SetUnion", stmt.Op)
	}
	if stmt.All {
		t.Error("All = true, want false (DISTINCT is default)")
	}
	assertIsSelectStmt(t, stmt.Left, "Left")
	assertIsSelectStmt(t, stmt.Right, "Right")
}

func TestUnionAll(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 UNION ALL SELECT a FROM t2")

	if stmt.Op != ast.SetUnion {
		t.Errorf("Op = %v, want SetUnion", stmt.Op)
	}
	if !stmt.All {
		t.Error("All = false, want true")
	}
	assertIsSelectStmt(t, stmt.Left, "Left")
	assertIsSelectStmt(t, stmt.Right, "Right")
}

func TestUnionDistinct(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 UNION DISTINCT SELECT a FROM t2")

	if stmt.Op != ast.SetUnion {
		t.Errorf("Op = %v, want SetUnion", stmt.Op)
	}
	if stmt.All {
		t.Error("All = true, want false for DISTINCT")
	}
}

// ---------------------------------------------------------------------------
// INTERSECT
// ---------------------------------------------------------------------------

func TestIntersectBasic(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 INTERSECT SELECT a FROM t2")

	if stmt.Op != ast.SetIntersect {
		t.Errorf("Op = %v, want SetIntersect", stmt.Op)
	}
	if stmt.All {
		t.Error("All = true, want false")
	}
	assertIsSelectStmt(t, stmt.Left, "Left")
	assertIsSelectStmt(t, stmt.Right, "Right")
}

func TestIntersectAll(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 INTERSECT ALL SELECT a FROM t2")

	if stmt.Op != ast.SetIntersect {
		t.Errorf("Op = %v, want SetIntersect", stmt.Op)
	}
	if !stmt.All {
		t.Error("All = false, want true")
	}
}

// ---------------------------------------------------------------------------
// EXCEPT
// ---------------------------------------------------------------------------

func TestExceptBasic(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 EXCEPT SELECT a FROM t2")

	if stmt.Op != ast.SetExcept {
		t.Errorf("Op = %v, want SetExcept", stmt.Op)
	}
	if stmt.All {
		t.Error("All = true, want false")
	}
	assertIsSelectStmt(t, stmt.Left, "Left")
	assertIsSelectStmt(t, stmt.Right, "Right")
}

func TestExceptAll(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 EXCEPT ALL SELECT a FROM t2")

	if stmt.Op != ast.SetExcept {
		t.Errorf("Op = %v, want SetExcept", stmt.Op)
	}
	if !stmt.All {
		t.Error("All = false, want true")
	}
}

// ---------------------------------------------------------------------------
// MINUS (alias for EXCEPT)
// ---------------------------------------------------------------------------

func TestMinusBasic(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 MINUS SELECT a FROM t2")

	if stmt.Op != ast.SetExcept {
		t.Errorf("Op = %v, want SetExcept (MINUS is alias for EXCEPT)", stmt.Op)
	}
	if stmt.All {
		t.Error("All = true, want false")
	}
	assertIsSelectStmt(t, stmt.Left, "Left")
	assertIsSelectStmt(t, stmt.Right, "Right")
}

func TestMinusAll(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT a FROM t1 MINUS ALL SELECT a FROM t2")

	if stmt.Op != ast.SetExcept {
		t.Errorf("Op = %v, want SetExcept", stmt.Op)
	}
	if !stmt.All {
		t.Error("All = false, want true")
	}
}

// ---------------------------------------------------------------------------
// Chained set operations (left-associativity)
// ---------------------------------------------------------------------------

// TestUnionChained verifies that three-way UNION is left-associative:
//
//	SELECT 1 UNION SELECT 2 UNION SELECT 3
//	=> SetOpStmt{ Left=SetOpStmt{Left=S1, Right=S2}, Right=S3 }
func TestUnionChained(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT 1 UNION SELECT 2 UNION SELECT 3")

	if stmt.Op != ast.SetUnion {
		t.Errorf("outer Op = %v, want SetUnion", stmt.Op)
	}

	// Right side of outermost must be a plain SELECT.
	assertIsSelectStmt(t, stmt.Right, "Right")

	// Left side of outermost must be another SetOpStmt.
	inner, ok := stmt.Left.(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("Left = %T, want *ast.SetOpStmt (left-associative)", stmt.Left)
	}
	if inner.Op != ast.SetUnion {
		t.Errorf("inner Op = %v, want SetUnion", inner.Op)
	}
	assertIsSelectStmt(t, inner.Left, "inner.Left")
	assertIsSelectStmt(t, inner.Right, "inner.Right")
}

// ---------------------------------------------------------------------------
// Mixed precedence: INTERSECT binds tighter than UNION
// ---------------------------------------------------------------------------

// TestMixedPrecedenceUnionIntersect verifies:
//
//	SELECT 1 UNION SELECT 2 INTERSECT SELECT 3
//	=> SetOpStmt{ Op=UNION, Left=S1, Right=SetOpStmt{Op=INTERSECT, Left=S2, Right=S3} }
func TestMixedPrecedenceUnionIntersect(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT 1 UNION SELECT 2 INTERSECT SELECT 3")

	if stmt.Op != ast.SetUnion {
		t.Errorf("outer Op = %v, want SetUnion", stmt.Op)
	}

	// Left side of UNION must be a plain SELECT.
	assertIsSelectStmt(t, stmt.Left, "Left")

	// Right side of UNION must be an INTERSECT SetOpStmt.
	inner, ok := stmt.Right.(*ast.SetOpStmt)
	if !ok {
		t.Fatalf("Right = %T, want *ast.SetOpStmt (INTERSECT has higher precedence)", stmt.Right)
	}
	if inner.Op != ast.SetIntersect {
		t.Errorf("inner Op = %v, want SetIntersect", inner.Op)
	}
}

// ---------------------------------------------------------------------------
// AST tag and NodeLoc
// ---------------------------------------------------------------------------

func TestSetOpStmtTag(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT 1 UNION SELECT 2")
	if stmt.Tag() != ast.T_SetOpStmt {
		t.Errorf("Tag() = %v, want T_SetOpStmt", stmt.Tag())
	}
}

func TestSetOpStmtNodeLoc(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT 1 UNION SELECT 2")
	loc := ast.NodeLoc(stmt)
	if loc.Start < 0 || loc.End <= loc.Start {
		t.Errorf("NodeLoc = %+v, expected valid non-empty range", loc)
	}
}

// ---------------------------------------------------------------------------
// Walk visits both sides
// ---------------------------------------------------------------------------

func TestSetOpStmtWalk(t *testing.T) {
	stmt := mustParseSetOp(t, "SELECT 1 UNION SELECT 2")

	count := make(map[ast.NodeTag]int)
	ast.Inspect(stmt, func(n ast.Node) bool {
		if n != nil {
			count[n.Tag()]++
		}
		return true
	})

	// We expect at minimum: SetOpStmt, SelectStmt (left), SelectStmt (right)
	if count[ast.T_SetOpStmt] != 1 {
		t.Errorf("T_SetOpStmt visited %d times, want 1", count[ast.T_SetOpStmt])
	}
	if count[ast.T_SelectStmt] != 2 {
		t.Errorf("T_SelectStmt visited %d times, want 2", count[ast.T_SelectStmt])
	}
}

// ---------------------------------------------------------------------------
// No-op: plain SELECT still returns *ast.SelectStmt (not wrapped)
// ---------------------------------------------------------------------------

func TestPlainSelectUnchanged(t *testing.T) {
	file, errs := Parse("SELECT a FROM t")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatal("no statements")
	}
	if _, ok := file.Stmts[0].(*ast.SelectStmt); !ok {
		t.Fatalf("got %T, want *ast.SelectStmt", file.Stmts[0])
	}
}
