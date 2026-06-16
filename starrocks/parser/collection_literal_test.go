package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

func firstItemExpr(t *testing.T, input string) ast.Node {
	t.Helper()
	return mustParseSelect(t, input).Items[0].Expr
}

// TestMapLiteralTyped asserts the typed map produces a MapLiteral carrying the
// map<k,v> type and its entries.
func TestMapLiteralTyped(t *testing.T) {
	ml, ok := firstItemExpr(t, "SELECT map<varchar,int>{'x':1,'y':2} FROM t").(*ast.MapLiteral)
	if !ok {
		t.Fatalf("expr = %T, want *ast.MapLiteral", firstItemExpr(t, "SELECT map<varchar,int>{'x':1} FROM t"))
	}
	if ml.MapType == nil {
		t.Error("MapType = nil, want the map<varchar,int> type")
	}
	if len(ml.Entries) != 2 {
		t.Errorf("entries = %d, want 2", len(ml.Entries))
	}
}

// TestMapLiteralUntyped asserts MAP{...} produces a MapLiteral with no type
// (without the feature it silently mis-parses to a bare ColumnRef).
func TestMapLiteralUntyped(t *testing.T) {
	ml, ok := firstItemExpr(t, "SELECT MAP{'x':1} FROM t").(*ast.MapLiteral)
	if !ok {
		t.Fatalf("expr = %T, want *ast.MapLiteral", firstItemExpr(t, "SELECT MAP{'x':1} FROM t"))
	}
	if ml.MapType != nil {
		t.Error("MapType != nil, want nil for the untyped MAP{...} form")
	}
	if len(ml.Entries) != 1 {
		t.Errorf("entries = %d, want 1", len(ml.Entries))
	}
}

func TestMapLiteralEmpty(t *testing.T) {
	ml := firstItemExpr(t, "SELECT map<int,int>{} FROM t").(*ast.MapLiteral)
	if len(ml.Entries) != 0 {
		t.Errorf("entries = %d, want 0", len(ml.Entries))
	}
}

func TestArrayLiteralTyped(t *testing.T) {
	al, ok := firstItemExpr(t, "SELECT array<int>[1,2,3] FROM t").(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr = %T, want *ast.ArrayLiteral", firstItemExpr(t, "SELECT array<int>[1,2,3] FROM t"))
	}
	if al.ElemType == nil {
		t.Error("ElemType = nil, want the array<int> type")
	}
	if len(al.Elements) != 3 {
		t.Errorf("elements = %d, want 3", len(al.Elements))
	}
}

func TestArrayLiteralUntyped(t *testing.T) {
	al, ok := firstItemExpr(t, "SELECT [1,2,3] FROM t").(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr = %T, want *ast.ArrayLiteral", firstItemExpr(t, "SELECT [1,2,3] FROM t"))
	}
	if al.ElemType != nil {
		t.Error("ElemType != nil, want nil for the untyped [...] form")
	}
	if len(al.Elements) != 3 {
		t.Errorf("elements = %d, want 3", len(al.Elements))
	}
}

// Map and array collection literals: map<k,v>{..}, MAP{..}, array<t>[..], [..].

func TestMapLiteralTypedParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT map<varchar,int>{'x':1} FROM t")
}

func TestMapLiteralMultiEntryParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT map<varchar,int>{'x':1,'y':2} FROM t")
}

func TestMapLiteralUntypedParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT MAP{'x':1} FROM t")
}

func TestMapLiteralEmptyParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT map<int,int>{} FROM t")
}

func TestArrayLiteralTypedParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT array<int>[1,2,3] FROM t")
}

func TestArrayLiteralEmptyParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT array<int>[] FROM t")
}

func TestArrayLiteralUntypedParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT [1,2,3] FROM t")
}

func TestMapLiteralNestedParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT map<varchar,array<int>>{'x':[1,2]} FROM t")
}

// Regression: map/array used as ordinary identifiers must still parse.
func TestMapAsIdentifierParses(t *testing.T) {
	_ = mustParseSelect(t, "SELECT map FROM t")
}

// TestMapColumnLessThan guards against the map-literal dispatch over-committing:
// `map` is a non-reserved keyword usable as a column, so `map < 5` is a
// comparison, not a map<...> type. (Single-token peek cannot tell them apart;
// the dispatch must backtrack when no '{' follows the angle brackets.)
func TestMapColumnLessThan(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT cnt FROM t WHERE map < 5")
	be, ok := stmt.Where.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("where = %T, want *ast.BinaryExpr (map < 5)", stmt.Where)
	}
	if _, ok := be.Left.(*ast.ColumnRef); !ok {
		t.Errorf("where.Left = %T, want *ast.ColumnRef (map)", be.Left)
	}
}

// TestArrayColumnLessThan is the same guard for `array` (also non-reserved in
// omni), where `array < 5` must parse as a comparison, not array<...>.
func TestArrayColumnLessThan(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT cnt FROM t WHERE array < 5")
	if _, ok := stmt.Where.(*ast.BinaryExpr); !ok {
		t.Fatalf("where = %T, want *ast.BinaryExpr (array < 5)", stmt.Where)
	}
}
