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

// Paths, vars, params, subqueries, collection literals, path steps (exprs.go).
var _ ExprNode = (*PathExpr)(nil)
var _ TableExpr = (*PathExpr)(nil)
var _ ExprNode = (*VarRef)(nil)
var _ TableExpr = (*VarRef)(nil)
var _ ExprNode = (*ParamRef)(nil)
var _ ExprNode = (*SubLink)(nil)
var _ TableExpr = (*SubLink)(nil)
var _ ExprNode = (*ListLit)(nil)
var _ ExprNode = (*BagLit)(nil)
var _ ExprNode = (*TupleLit)(nil)
var _ Node = (*TuplePair)(nil)
var _ PathStep = (*DotStep)(nil)
var _ PathStep = (*AllFieldsStep)(nil)
var _ PathStep = (*IndexStep)(nil)
var _ PathStep = (*WildcardStep)(nil)

// Table expression nodes (tableexprs.go).
var _ TableExpr = (*TableRef)(nil)
var _ TableExpr = (*AliasedSource)(nil)
var _ TableExpr = (*JoinExpr)(nil)
var _ TableExpr = (*UnpivotExpr)(nil)

// Type names (types.go).
var _ TypeName = (*TypeRef)(nil)

// Top-level statements (stmts.go).
var _ StmtNode = (*SelectStmt)(nil)
var _ StmtNode = (*SetOpStmt)(nil)
var _ StmtNode = (*ExplainStmt)(nil)
var _ StmtNode = (*InsertStmt)(nil)
var _ StmtNode = (*UpdateStmt)(nil)
var _ StmtNode = (*DeleteStmt)(nil)
var _ StmtNode = (*UpsertStmt)(nil)
var _ StmtNode = (*ReplaceStmt)(nil)
var _ StmtNode = (*RemoveStmt)(nil)
var _ StmtNode = (*CreateTableStmt)(nil)
var _ StmtNode = (*CreateIndexStmt)(nil)
var _ StmtNode = (*DropTableStmt)(nil)
var _ StmtNode = (*DropIndexStmt)(nil)
var _ StmtNode = (*ExecStmt)(nil)

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
		{"PathExpr", &PathExpr{Loc: Loc{Start: 10, End: 20}}},
		{"VarRef", &VarRef{Loc: Loc{Start: 10, End: 20}}},
		{"ParamRef", &ParamRef{Loc: Loc{Start: 10, End: 20}}},
		{"SubLink", &SubLink{Loc: Loc{Start: 10, End: 20}}},
		{"ListLit", &ListLit{Loc: Loc{Start: 10, End: 20}}},
		{"BagLit", &BagLit{Loc: Loc{Start: 10, End: 20}}},
		{"TupleLit", &TupleLit{Loc: Loc{Start: 10, End: 20}}},
		{"TuplePair", &TuplePair{Loc: Loc{Start: 10, End: 20}}},
		{"DotStep", &DotStep{Loc: Loc{Start: 10, End: 20}}},
		{"AllFieldsStep", &AllFieldsStep{Loc: Loc{Start: 10, End: 20}}},
		{"IndexStep", &IndexStep{Loc: Loc{Start: 10, End: 20}}},
		{"WildcardStep", &WildcardStep{Loc: Loc{Start: 10, End: 20}}},
		{"TableRef", &TableRef{Loc: Loc{Start: 10, End: 20}}},
		{"AliasedSource", &AliasedSource{Loc: Loc{Start: 10, End: 20}}},
		{"JoinExpr", &JoinExpr{Loc: Loc{Start: 10, End: 20}}},
		{"UnpivotExpr", &UnpivotExpr{Loc: Loc{Start: 10, End: 20}}},
		{"TypeRef", &TypeRef{Loc: Loc{Start: 10, End: 20}}},
		{"SelectStmt", &SelectStmt{Loc: Loc{Start: 10, End: 20}}},
		{"SetOpStmt", &SetOpStmt{Loc: Loc{Start: 10, End: 20}}},
		{"ExplainStmt", &ExplainStmt{Loc: Loc{Start: 10, End: 20}}},
		{"InsertStmt", &InsertStmt{Loc: Loc{Start: 10, End: 20}}},
		{"UpdateStmt", &UpdateStmt{Loc: Loc{Start: 10, End: 20}}},
		{"DeleteStmt", &DeleteStmt{Loc: Loc{Start: 10, End: 20}}},
		{"UpsertStmt", &UpsertStmt{Loc: Loc{Start: 10, End: 20}}},
		{"ReplaceStmt", &ReplaceStmt{Loc: Loc{Start: 10, End: 20}}},
		{"RemoveStmt", &RemoveStmt{Loc: Loc{Start: 10, End: 20}}},
		{"CreateTableStmt", &CreateTableStmt{Loc: Loc{Start: 10, End: 20}}},
		{"CreateIndexStmt", &CreateIndexStmt{Loc: Loc{Start: 10, End: 20}}},
		{"DropTableStmt", &DropTableStmt{Loc: Loc{Start: 10, End: 20}}},
		{"DropIndexStmt", &DropIndexStmt{Loc: Loc{Start: 10, End: 20}}},
		{"ExecStmt", &ExecStmt{Loc: Loc{Start: 10, End: 20}}},
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
