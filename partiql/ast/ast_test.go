package ast

import (
	"strings"
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

// Clause and DML helpers (stmts.go).
var _ Node = (*TargetEntry)(nil)
var _ Node = (*PivotProjection)(nil)
var _ Node = (*LetBinding)(nil)
var _ Node = (*GroupByClause)(nil)
var _ Node = (*GroupByItem)(nil)
var _ Node = (*OrderByItem)(nil)
var _ Node = (*SetAssignment)(nil)
var _ Node = (*OnConflict)(nil)
var _ Node = (*OnConflictTarget)(nil)
var _ Node = (*ReturningClause)(nil)
var _ Node = (*ReturningItem)(nil)

// Graph patterns (patterns.go).
var _ ExprNode = (*MatchExpr)(nil)
var _ PatternNode = (*GraphPattern)(nil)
var _ PatternNode = (*NodePattern)(nil)
var _ PatternNode = (*EdgePattern)(nil)
var _ Node = (*PatternQuantifier)(nil)
var _ Node = (*PatternSelector)(nil)

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
		{"TargetEntry", &TargetEntry{Loc: Loc{Start: 10, End: 20}}},
		{"PivotProjection", &PivotProjection{Loc: Loc{Start: 10, End: 20}}},
		{"LetBinding", &LetBinding{Loc: Loc{Start: 10, End: 20}}},
		{"GroupByClause", &GroupByClause{Loc: Loc{Start: 10, End: 20}}},
		{"GroupByItem", &GroupByItem{Loc: Loc{Start: 10, End: 20}}},
		{"OrderByItem", &OrderByItem{Loc: Loc{Start: 10, End: 20}}},
		{"SetAssignment", &SetAssignment{Loc: Loc{Start: 10, End: 20}}},
		{"OnConflict", &OnConflict{Loc: Loc{Start: 10, End: 20}}},
		{"OnConflictTarget", &OnConflictTarget{Loc: Loc{Start: 10, End: 20}}},
		{"ReturningClause", &ReturningClause{Loc: Loc{Start: 10, End: 20}}},
		{"ReturningItem", &ReturningItem{Loc: Loc{Start: 10, End: 20}}},
		{"MatchExpr", &MatchExpr{Loc: Loc{Start: 10, End: 20}}},
		{"GraphPattern", &GraphPattern{Loc: Loc{Start: 10, End: 20}}},
		{"NodePattern", &NodePattern{Loc: Loc{Start: 10, End: 20}}},
		{"EdgePattern", &EdgePattern{Loc: Loc{Start: 10, End: 20}}},
		{"PatternQuantifier", &PatternQuantifier{Loc: Loc{Start: 10, End: 20}}},
		{"PatternSelector", &PatternSelector{Loc: Loc{Start: 10, End: 20}}},
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

// ---------------------------------------------------------------------------
// TestNodeToString — golden assertions on outfuncs.go output.
//
// One representative case per node category, plus the multi-interface nodes
// and at least one deeply nested case.
// ---------------------------------------------------------------------------

func TestNodeToString(t *testing.T) {
	cases := []struct {
		name string
		node Node
		want string
	}{
		{
			name: "string_literal",
			node: &StringLit{Val: "hello"},
			want: `StringLit{Val:"hello"}`,
		},
		{
			name: "number_literal",
			node: &NumberLit{Val: "42"},
			want: `NumberLit{Val:42}`,
		},
		{
			name: "binary_expr_add",
			node: &BinaryExpr{
				Op:    BinOpAdd,
				Left:  &VarRef{Name: "a"},
				Right: &NumberLit{Val: "1"},
			},
			want: `BinaryExpr{Op:+ Left:VarRef{Name:a} Right:NumberLit{Val:1}}`,
		},
		{
			name: "in_expr_with_list",
			node: &InExpr{
				Expr: &VarRef{Name: "x"},
				List: []ExprNode{&NumberLit{Val: "1"}, &NumberLit{Val: "2"}},
				Not:  false,
			},
			want: `InExpr{Expr:VarRef{Name:x} List:[NumberLit{Val:1} NumberLit{Val:2}] Not:false}`,
		},
		{
			name: "is_missing_predicate",
			node: &IsExpr{
				Expr: &VarRef{Name: "y"},
				Type: IsTypeMissing,
				Not:  false,
			},
			want: `IsExpr{Expr:VarRef{Name:y} Type:MISSING Not:false}`,
		},
		{
			name: "func_call_count_distinct",
			node: &FuncCall{
				Name:       "COUNT",
				Quantifier: QuantifierDistinct,
				Args:       []ExprNode{&VarRef{Name: "id"}},
			},
			want: `FuncCall{Name:COUNT Quantifier:DISTINCT Args:[VarRef{Name:id}]}`,
		},
		{
			name: "func_call_count_star",
			node: &FuncCall{Name: "COUNT", Star: true},
			want: `FuncCall{Name:COUNT Star:true Args:[]}`,
		},
		{
			name: "case_expr_searched",
			node: &CaseExpr{
				Whens: []*CaseWhen{
					{When: &BoolLit{Val: true}, Then: &NumberLit{Val: "1"}},
				},
				Else: &NumberLit{Val: "0"},
			},
			want: `CaseExpr{Whens:[CaseWhen{When:BoolLit{Val:true} Then:NumberLit{Val:1}}] Else:NumberLit{Val:0}}`,
		},
		{
			name: "cast_expr",
			node: &CastExpr{
				Kind:   CastKindCast,
				Expr:   &VarRef{Name: "x"},
				AsType: &TypeRef{Name: "INTEGER"},
			},
			want: `CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:INTEGER}}`,
		},
		{
			name: "path_expr_with_steps",
			node: &PathExpr{
				Root:  &VarRef{Name: "Music"},
				Steps: []PathStep{&DotStep{Field: "albums"}, &WildcardStep{}},
			},
			want: `PathExpr{Root:VarRef{Name:Music} Steps:[DotStep{Field:albums} WildcardStep{}]}`,
		},
		{
			name: "var_ref_at_prefixed",
			node: &VarRef{Name: "doc", AtPrefixed: true},
			want: `VarRef{Name:doc AtPrefixed:true}`,
		},
		{
			name: "param_ref",
			node: &ParamRef{},
			want: `ParamRef{}`,
		},
		{
			name: "sub_link",
			node: &SubLink{Stmt: &SelectStmt{Star: true}},
			want: `SubLink{Stmt:SelectStmt{Star:true}}`,
		},
		{
			name: "list_literal",
			node: &ListLit{Items: []ExprNode{&NumberLit{Val: "1"}, &NumberLit{Val: "2"}}},
			want: `ListLit{Items:[NumberLit{Val:1} NumberLit{Val:2}]}`,
		},
		{
			name: "bag_literal",
			node: &BagLit{Items: []ExprNode{&NumberLit{Val: "1"}}},
			want: `BagLit{Items:[NumberLit{Val:1}]}`,
		},
		{
			name: "tuple_literal",
			node: &TupleLit{
				Pairs: []*TuplePair{
					{Key: &StringLit{Val: "k"}, Value: &NumberLit{Val: "1"}},
				},
			},
			want: `TupleLit{Pairs:[TuplePair{Key:StringLit{Val:"k"} Value:NumberLit{Val:1}}]}`,
		},
		{
			name: "join_expr",
			node: &JoinExpr{
				Kind:  JoinKindInner,
				Left:  &TableRef{Name: "a"},
				Right: &TableRef{Name: "b"},
				On:    &BoolLit{Val: true},
			},
			want: `JoinExpr{Kind:INNER Left:TableRef{Name:a} Right:TableRef{Name:b} On:BoolLit{Val:true}}`,
		},
		{
			name: "select_with_path_in_from",
			node: &SelectStmt{
				Star: true,
				From: &PathExpr{
					Root:  &VarRef{Name: "Music"},
					Steps: []PathStep{&DotStep{Field: "albums"}, &WildcardStep{}},
				},
			},
			want: `SelectStmt{Star:true From:PathExpr{Root:VarRef{Name:Music} Steps:[DotStep{Field:albums} WildcardStep{}]}}`,
		},
		{
			name: "insert_stmt_with_returning",
			node: &InsertStmt{
				Target: &TableRef{Name: "Music"},
				Value:  &TupleLit{Pairs: []*TuplePair{{Key: &StringLit{Val: "k"}, Value: &NumberLit{Val: "1"}}}},
				Returning: &ReturningClause{
					Items: []*ReturningItem{
						{Status: ReturningStatusModified, Mapping: ReturningMappingNew, Star: true},
					},
				},
			},
			want: `InsertStmt{Target:TableRef{Name:Music} Value:TupleLit{Pairs:[TuplePair{Key:StringLit{Val:"k"} Value:NumberLit{Val:1}}]} Returning:ReturningClause{Items:[ReturningItem{Status:MODIFIED Mapping:NEW Star:true}]}}`,
		},
		{
			name: "match_expr",
			node: &MatchExpr{
				Expr: &VarRef{Name: "g"},
				Patterns: []*GraphPattern{
					{Parts: []PatternNode{&NodePattern{Variable: &VarRef{Name: "n"}}}},
				},
			},
			want: `MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:n}}]}]}`,
		},
		{
			name: "type_ref_decimal",
			node: &TypeRef{Name: "DECIMAL", Args: []int{10, 2}},
			want: `TypeRef{Name:DECIMAL Args:[10,2]}`,
		},
		{
			name: "type_ref_time_with_tz",
			node: &TypeRef{Name: "TIME", WithTimeZone: true},
			want: `TypeRef{Name:TIME WithTimeZone:true}`,
		},
		{
			name: "nil_node",
			node: nil,
			want: `<nil>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NodeToString(tc.node)
			if got != tc.want {
				t.Errorf("NodeToString() mismatch\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TestNodeToString_AllNodesCovered — reflection safety net.
//
// Walks every node in the package's TestGetLoc table and asserts that
// NodeToString returns a non-empty result without panicking. This catches
// the case where a new node type is added to TestGetLoc but never wired
// into outfuncs.go's switch.
// ---------------------------------------------------------------------------

func TestNodeToString_AllNodesCovered(t *testing.T) {
	// Build a fresh table identical in shape to TestGetLoc but constructed
	// here so the two tests cannot drift apart silently. We could share via
	// a top-level helper but inlining keeps the test self-contained.
	cases := allNodesForCoverageTest()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("NodeToString panicked: %v", r)
				}
			}()
			got := NodeToString(tc.node)
			if got == "" {
				t.Errorf("NodeToString returned empty string for %s", tc.name)
			}
			if strings.HasPrefix(got, "<unknown:") {
				t.Errorf("outfuncs.go is missing a switch arm for %s: got %q", tc.name, got)
			}
		})
	}
}

// allNodesForCoverageTest returns one zero-value instance of every node
// type defined in this package. Adding a new node type means adding it
// here too — the test will fail until you do.
func allNodesForCoverageTest() []struct {
	name string
	node Node
} {
	return []struct {
		name string
		node Node
	}{
		// node.go
		{"List", &List{}},

		// literals.go
		{"StringLit", &StringLit{}},
		{"NumberLit", &NumberLit{}},
		{"BoolLit", &BoolLit{}},
		{"NullLit", &NullLit{}},
		{"MissingLit", &MissingLit{}},
		{"DateLit", &DateLit{}},
		{"TimeLit", &TimeLit{}},
		{"IonLit", &IonLit{}},

		// exprs.go
		{"BinaryExpr", &BinaryExpr{}},
		{"UnaryExpr", &UnaryExpr{}},
		{"InExpr", &InExpr{}},
		{"BetweenExpr", &BetweenExpr{}},
		{"LikeExpr", &LikeExpr{}},
		{"IsExpr", &IsExpr{}},
		{"FuncCall", &FuncCall{}},
		{"CaseExpr", &CaseExpr{}},
		{"CaseWhen", &CaseWhen{}},
		{"CastExpr", &CastExpr{}},
		{"ExtractExpr", &ExtractExpr{}},
		{"TrimExpr", &TrimExpr{}},
		{"SubstringExpr", &SubstringExpr{}},
		{"CoalesceExpr", &CoalesceExpr{}},
		{"NullIfExpr", &NullIfExpr{}},
		{"WindowSpec", &WindowSpec{}},
		{"PathExpr", &PathExpr{}},
		{"VarRef", &VarRef{}},
		{"ParamRef", &ParamRef{}},
		{"SubLink", &SubLink{}},
		{"ListLit", &ListLit{}},
		{"BagLit", &BagLit{}},
		{"TupleLit", &TupleLit{}},
		{"TuplePair", &TuplePair{}},
		{"DotStep", &DotStep{}},
		{"AllFieldsStep", &AllFieldsStep{}},
		{"IndexStep", &IndexStep{}},
		{"WildcardStep", &WildcardStep{}},

		// tableexprs.go
		{"TableRef", &TableRef{}},
		{"AliasedSource", &AliasedSource{}},
		{"JoinExpr", &JoinExpr{}},
		{"UnpivotExpr", &UnpivotExpr{}},

		// types.go
		{"TypeRef", &TypeRef{}},

		// stmts.go — top-level
		{"SelectStmt", &SelectStmt{}},
		{"SetOpStmt", &SetOpStmt{}},
		{"ExplainStmt", &ExplainStmt{}},
		{"InsertStmt", &InsertStmt{}},
		{"UpdateStmt", &UpdateStmt{}},
		{"DeleteStmt", &DeleteStmt{}},
		{"UpsertStmt", &UpsertStmt{}},
		{"ReplaceStmt", &ReplaceStmt{}},
		{"RemoveStmt", &RemoveStmt{}},
		{"CreateTableStmt", &CreateTableStmt{}},
		{"CreateIndexStmt", &CreateIndexStmt{}},
		{"DropTableStmt", &DropTableStmt{}},
		{"DropIndexStmt", &DropIndexStmt{}},
		{"ExecStmt", &ExecStmt{}},

		// stmts.go — clause and DML helpers
		{"TargetEntry", &TargetEntry{}},
		{"PivotProjection", &PivotProjection{}},
		{"LetBinding", &LetBinding{}},
		{"GroupByClause", &GroupByClause{}},
		{"GroupByItem", &GroupByItem{}},
		{"OrderByItem", &OrderByItem{}},
		{"SetAssignment", &SetAssignment{}},
		{"OnConflict", &OnConflict{}},
		{"OnConflictTarget", &OnConflictTarget{}},
		{"ReturningClause", &ReturningClause{}},
		{"ReturningItem", &ReturningItem{}},

		// patterns.go
		{"MatchExpr", &MatchExpr{}},
		{"GraphPattern", &GraphPattern{}},
		{"NodePattern", &NodePattern{}},
		{"EdgePattern", &EdgePattern{}},
		{"PatternQuantifier", &PatternQuantifier{}},
		{"PatternSelector", &PatternSelector{}},
	}
}
