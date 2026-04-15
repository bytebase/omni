package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

// mustParseJoin parses input, expects exactly one FROM item which must be a
// *ast.JoinClause, and returns it.
func mustParseJoin(t *testing.T, input string) *ast.JoinClause {
	t.Helper()
	stmt := mustParseSelect(t, input)
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	join, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	return join
}

// joinTableName extracts the table name string from the left or right node of
// a JoinClause (assuming it is a *ast.TableRef with a single-part name).
func joinTableName(t *testing.T, n ast.Node) string {
	t.Helper()
	switch v := n.(type) {
	case *ast.TableRef:
		if len(v.Name.Parts) == 0 {
			t.Fatal("TableRef.Name.Parts is empty")
		}
		return v.Name.Parts[len(v.Name.Parts)-1]
	case *ast.JoinClause:
		// For chained joins, return the tag description
		return "(JoinClause)"
	default:
		t.Fatalf("unexpected node type %T", n)
		return ""
	}
}

// ---------------------------------------------------------------------------
// INNER JOIN
// ---------------------------------------------------------------------------

func TestJoinInnerExplicit(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a INNER JOIN b ON a.id = b.id")
	if join.Type != ast.JoinInner {
		t.Errorf("join.Type = %v, want JoinInner", join.Type)
	}
	if joinTableName(t, join.Left) != "a" {
		t.Errorf("left = %q, want %q", joinTableName(t, join.Left), "a")
	}
	if joinTableName(t, join.Right) != "b" {
		t.Errorf("right = %q, want %q", joinTableName(t, join.Right), "b")
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
	if join.Natural {
		t.Error("Natural = true, want false")
	}
}

func TestJoinInnerImplicit(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a JOIN b ON a.id = b.id")
	if join.Type != ast.JoinInner {
		t.Errorf("join.Type = %v, want JoinInner", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// LEFT JOIN
// ---------------------------------------------------------------------------

func TestJoinLeft(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a LEFT JOIN b ON a.id = b.id")
	if join.Type != ast.JoinLeft {
		t.Errorf("join.Type = %v, want JoinLeft", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinLeftOuter(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a LEFT OUTER JOIN b ON a.id = b.id")
	if join.Type != ast.JoinLeft {
		t.Errorf("join.Type = %v, want JoinLeft", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// RIGHT JOIN
// ---------------------------------------------------------------------------

func TestJoinRight(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a RIGHT JOIN b ON a.id = b.id")
	if join.Type != ast.JoinRight {
		t.Errorf("join.Type = %v, want JoinRight", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinRightOuter(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a RIGHT OUTER JOIN b ON a.id = b.id")
	if join.Type != ast.JoinRight {
		t.Errorf("join.Type = %v, want JoinRight", join.Type)
	}
}

// ---------------------------------------------------------------------------
// FULL JOIN
// ---------------------------------------------------------------------------

func TestJoinFullOuter(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a FULL OUTER JOIN b ON a.id = b.id")
	if join.Type != ast.JoinFull {
		t.Errorf("join.Type = %v, want JoinFull", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinFull(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a FULL JOIN b ON a.id = b.id")
	if join.Type != ast.JoinFull {
		t.Errorf("join.Type = %v, want JoinFull", join.Type)
	}
}

// ---------------------------------------------------------------------------
// CROSS JOIN
// ---------------------------------------------------------------------------

func TestJoinCross(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a CROSS JOIN b")
	if join.Type != ast.JoinCross {
		t.Errorf("join.Type = %v, want JoinCross", join.Type)
	}
	if join.On != nil {
		t.Error("On should be nil for CROSS JOIN")
	}
	if len(join.Using) != 0 {
		t.Error("Using should be empty for CROSS JOIN")
	}
}

// ---------------------------------------------------------------------------
// SEMI / ANTI joins
// ---------------------------------------------------------------------------

func TestJoinLeftSemi(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a LEFT SEMI JOIN b ON a.id = b.id")
	if join.Type != ast.JoinLeftSemi {
		t.Errorf("join.Type = %v, want JoinLeftSemi", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinRightSemi(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a RIGHT SEMI JOIN b ON a.id = b.id")
	if join.Type != ast.JoinRightSemi {
		t.Errorf("join.Type = %v, want JoinRightSemi", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinLeftAnti(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a LEFT ANTI JOIN b ON a.id = b.id")
	if join.Type != ast.JoinLeftAnti {
		t.Errorf("join.Type = %v, want JoinLeftAnti", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinRightAnti(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a RIGHT ANTI JOIN b ON a.id = b.id")
	if join.Type != ast.JoinRightAnti {
		t.Errorf("join.Type = %v, want JoinRightAnti", join.Type)
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// NATURAL JOIN
// ---------------------------------------------------------------------------

func TestJoinNatural(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a NATURAL JOIN b")
	if join.Type != ast.JoinInner {
		t.Errorf("join.Type = %v, want JoinInner (natural)", join.Type)
	}
	if !join.Natural {
		t.Error("Natural = false, want true")
	}
	if join.On != nil {
		t.Error("On should be nil for NATURAL JOIN")
	}
}

func TestJoinNaturalLeft(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a NATURAL LEFT JOIN b")
	if join.Type != ast.JoinLeft {
		t.Errorf("join.Type = %v, want JoinLeft", join.Type)
	}
	if !join.Natural {
		t.Error("Natural = false, want true")
	}
}

func TestJoinNaturalRight(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a NATURAL RIGHT JOIN b")
	if join.Type != ast.JoinRight {
		t.Errorf("join.Type = %v, want JoinRight", join.Type)
	}
	if !join.Natural {
		t.Error("Natural = false, want true")
	}
}

// ---------------------------------------------------------------------------
// USING clause
// ---------------------------------------------------------------------------

func TestJoinUsingSingle(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a JOIN b USING (id)")
	if join.Type != ast.JoinInner {
		t.Errorf("join.Type = %v, want JoinInner", join.Type)
	}
	if len(join.Using) != 1 {
		t.Fatalf("Using = %v, want 1 column", join.Using)
	}
	if join.Using[0] != "id" {
		t.Errorf("Using[0] = %q, want %q", join.Using[0], "id")
	}
	if join.On != nil {
		t.Error("On should be nil when USING is present")
	}
}

func TestJoinUsingMultiple(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a JOIN b USING (id, name)")
	if len(join.Using) != 2 {
		t.Fatalf("Using = %v, want 2 columns", join.Using)
	}
	if join.Using[0] != "id" {
		t.Errorf("Using[0] = %q, want %q", join.Using[0], "id")
	}
	if join.Using[1] != "name" {
		t.Errorf("Using[1] = %q, want %q", join.Using[1], "name")
	}
}

func TestJoinLeftUsing(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a LEFT JOIN b USING (id)")
	if join.Type != ast.JoinLeft {
		t.Errorf("join.Type = %v, want JoinLeft", join.Type)
	}
	if len(join.Using) != 1 || join.Using[0] != "id" {
		t.Errorf("Using = %v, want [id]", join.Using)
	}
}

// ---------------------------------------------------------------------------
// Chained joins
// ---------------------------------------------------------------------------

func TestJoinChainedTwoInner(t *testing.T) {
	// SELECT * FROM a JOIN b ON a.id = b.id JOIN c ON b.id = c.id
	// Expected: JoinClause( JoinClause(a, b, ON...), c, ON... )
	stmt := mustParseSelect(t, "SELECT * FROM a JOIN b ON a.id = b.id JOIN c ON b.id = c.id")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	outerJoin, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	if outerJoin.Type != ast.JoinInner {
		t.Errorf("outer join type = %v, want JoinInner", outerJoin.Type)
	}
	if joinTableName(t, outerJoin.Right) != "c" {
		t.Errorf("outer right = %q, want %q", joinTableName(t, outerJoin.Right), "c")
	}
	if outerJoin.On == nil {
		t.Error("outer join On = nil, want non-nil")
	}

	innerJoin, ok := outerJoin.Left.(*ast.JoinClause)
	if !ok {
		t.Fatalf("outer.Left = %T, want *ast.JoinClause", outerJoin.Left)
	}
	if innerJoin.Type != ast.JoinInner {
		t.Errorf("inner join type = %v, want JoinInner", innerJoin.Type)
	}
	if joinTableName(t, innerJoin.Left) != "a" {
		t.Errorf("inner left = %q, want %q", joinTableName(t, innerJoin.Left), "a")
	}
	if joinTableName(t, innerJoin.Right) != "b" {
		t.Errorf("inner right = %q, want %q", joinTableName(t, innerJoin.Right), "b")
	}
}

func TestJoinChainedMixed(t *testing.T) {
	// SELECT * FROM a LEFT JOIN b ON a.id = b.id JOIN c ON b.id = c.id
	stmt := mustParseSelect(t, "SELECT * FROM a LEFT JOIN b ON a.id = b.id JOIN c ON b.id = c.id")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	outerJoin, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	if outerJoin.Type != ast.JoinInner {
		t.Errorf("outer join type = %v, want JoinInner", outerJoin.Type)
	}
	if joinTableName(t, outerJoin.Right) != "c" {
		t.Errorf("outer right = %q, want %q", joinTableName(t, outerJoin.Right), "c")
	}

	innerJoin, ok := outerJoin.Left.(*ast.JoinClause)
	if !ok {
		t.Fatalf("outer.Left = %T, want *ast.JoinClause", outerJoin.Left)
	}
	if innerJoin.Type != ast.JoinLeft {
		t.Errorf("inner join type = %v, want JoinLeft", innerJoin.Type)
	}
}

func TestJoinChainedThree(t *testing.T) {
	// Verify 3-way chain parses without error
	stmt := mustParseSelect(t, "SELECT * FROM a JOIN b ON a.id = b.id LEFT JOIN c ON b.id = c.id CROSS JOIN d")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	outerJoin, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	if outerJoin.Type != ast.JoinCross {
		t.Errorf("outermost join type = %v, want JoinCross", outerJoin.Type)
	}
	if joinTableName(t, outerJoin.Right) != "d" {
		t.Errorf("outermost right = %q, want %q", joinTableName(t, outerJoin.Right), "d")
	}
}

// ---------------------------------------------------------------------------
// Execution hints
// ---------------------------------------------------------------------------

func TestJoinBroadcastHint(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a JOIN [broadcast] b ON a.id = b.id")
	if join.Type != ast.JoinInner {
		t.Errorf("join.Type = %v, want JoinInner", join.Type)
	}
	if len(join.Hints) == 0 {
		t.Error("Hints should be non-empty for [broadcast]")
	}
	if join.Hints[0] != "broadcast" {
		t.Errorf("Hints[0] = %q, want %q", join.Hints[0], "broadcast")
	}
	if join.On == nil {
		t.Error("On = nil, want non-nil")
	}
}

func TestJoinShuffleHint(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a LEFT JOIN [shuffle] b ON a.id = b.id")
	if join.Type != ast.JoinLeft {
		t.Errorf("join.Type = %v, want JoinLeft", join.Type)
	}
	if len(join.Hints) == 0 {
		t.Error("Hints should be non-empty for [shuffle]")
	}
	if join.Hints[0] != "shuffle" {
		t.Errorf("Hints[0] = %q, want %q", join.Hints[0], "shuffle")
	}
}

// ---------------------------------------------------------------------------
// ON condition detail
// ---------------------------------------------------------------------------

func TestJoinOnConditionBinary(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a JOIN b ON a.id = b.id")
	binExpr, ok := join.On.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("On = %T, want *ast.BinaryExpr", join.On)
	}
	if binExpr.Op != ast.BinEq {
		t.Errorf("On.Op = %v, want BinEq", binExpr.Op)
	}
}

// ---------------------------------------------------------------------------
// Node tag and loc
// ---------------------------------------------------------------------------

func TestJoinClauseTag(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM a JOIN b ON a.id = b.id")
	if join.Tag() != ast.T_JoinClause {
		t.Errorf("Tag() = %v, want T_JoinClause", join.Tag())
	}
}

func TestJoinClauseLoc(t *testing.T) {
	input := "SELECT * FROM a JOIN b ON a.id = b.id"
	join := mustParseJoin(t, input)
	if join.Loc.Start < 0 {
		t.Errorf("Loc.Start = %d, want >= 0", join.Loc.Start)
	}
	if join.Loc.End <= join.Loc.Start {
		t.Errorf("Loc.End = %d, want > Start (%d)", join.Loc.End, join.Loc.Start)
	}
}
