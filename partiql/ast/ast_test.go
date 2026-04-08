package ast

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Compile-time interface assertions.
//
// Every node type added to this package gets a `var _ <Interface> = (*Type)(nil)`
// line below. The file fails to compile if a node's interface set drifts.
// Tasks add their assertions to the appropriate section as they grow the AST.
// ---------------------------------------------------------------------------

var _ Node = (*List)(nil)

// Literals — all implement ExprNode.
var _ ExprNode = (*StringLit)(nil)
var _ ExprNode = (*NumberLit)(nil)
var _ ExprNode = (*BoolLit)(nil)
var _ ExprNode = (*NullLit)(nil)
var _ ExprNode = (*MissingLit)(nil)
var _ ExprNode = (*DateLit)(nil)
var _ ExprNode = (*TimeLit)(nil)
var _ ExprNode = (*IonLit)(nil)

// Operators & predicates (exprs.go).
var _ ExprNode = (*BinaryExpr)(nil)
var _ ExprNode = (*UnaryExpr)(nil)
var _ ExprNode = (*InExpr)(nil)
var _ ExprNode = (*BetweenExpr)(nil)
var _ ExprNode = (*LikeExpr)(nil)
var _ ExprNode = (*IsExpr)(nil)

// Special-form expressions (exprs.go).
var _ ExprNode = (*FuncCall)(nil)
var _ ExprNode = (*CaseExpr)(nil)
var _ Node = (*CaseWhen)(nil)
var _ ExprNode = (*CastExpr)(nil)
var _ ExprNode = (*ExtractExpr)(nil)
var _ ExprNode = (*TrimExpr)(nil)
var _ ExprNode = (*SubstringExpr)(nil)
var _ ExprNode = (*CoalesceExpr)(nil)
var _ ExprNode = (*NullIfExpr)(nil)
var _ Node = (*WindowSpec)(nil)

// ---------------------------------------------------------------------------
// TestGetLoc — table-driven Loc round-trip.
//
// One row per node type. Each row constructs the node with Loc{10, 20},
// calls GetLoc(), and asserts the result.
// ---------------------------------------------------------------------------

func TestGetLoc(t *testing.T) {
	cases := []struct {
		name string
		node Node
	}{
		{"List", &List{Loc: Loc{Start: 10, End: 20}}},
		{"StringLit", &StringLit{Loc: Loc{Start: 10, End: 20}}},
		{"NumberLit", &NumberLit{Loc: Loc{Start: 10, End: 20}}},
		{"BoolLit", &BoolLit{Loc: Loc{Start: 10, End: 20}}},
		{"NullLit", &NullLit{Loc: Loc{Start: 10, End: 20}}},
		{"MissingLit", &MissingLit{Loc: Loc{Start: 10, End: 20}}},
		{"DateLit", &DateLit{Loc: Loc{Start: 10, End: 20}}},
		{"TimeLit", &TimeLit{Loc: Loc{Start: 10, End: 20}}},
		{"IonLit", &IonLit{Loc: Loc{Start: 10, End: 20}}},
		{"BinaryExpr", &BinaryExpr{Loc: Loc{Start: 10, End: 20}}},
		{"UnaryExpr", &UnaryExpr{Loc: Loc{Start: 10, End: 20}}},
		{"InExpr", &InExpr{Loc: Loc{Start: 10, End: 20}}},
		{"BetweenExpr", &BetweenExpr{Loc: Loc{Start: 10, End: 20}}},
		{"LikeExpr", &LikeExpr{Loc: Loc{Start: 10, End: 20}}},
		{"IsExpr", &IsExpr{Loc: Loc{Start: 10, End: 20}}},
		{"FuncCall", &FuncCall{Loc: Loc{Start: 10, End: 20}}},
		{"CaseExpr", &CaseExpr{Loc: Loc{Start: 10, End: 20}}},
		{"CaseWhen", &CaseWhen{Loc: Loc{Start: 10, End: 20}}},
		{"CastExpr", &CastExpr{Loc: Loc{Start: 10, End: 20}}},
		{"ExtractExpr", &ExtractExpr{Loc: Loc{Start: 10, End: 20}}},
		{"TrimExpr", &TrimExpr{Loc: Loc{Start: 10, End: 20}}},
		{"SubstringExpr", &SubstringExpr{Loc: Loc{Start: 10, End: 20}}},
		{"CoalesceExpr", &CoalesceExpr{Loc: Loc{Start: 10, End: 20}}},
		{"NullIfExpr", &NullIfExpr{Loc: Loc{Start: 10, End: 20}}},
		{"WindowSpec", &WindowSpec{Loc: Loc{Start: 10, End: 20}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.node.GetLoc()
			if got.Start != 10 || got.End != 20 {
				t.Errorf("GetLoc() = %+v, want {10, 20}", got)
			}
		})
	}
}
