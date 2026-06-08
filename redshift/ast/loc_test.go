package ast

import "testing"

func TestNodeLocKnownTypes(t *testing.T) {
	cases := []struct {
		name string
		node Node
		loc  Loc
	}{
		{"SelectStmt", &SelectStmt{Loc: Loc{Start: 0, End: 10}}, Loc{0, 10}},
		{"InsertStmt", &InsertStmt{Loc: Loc{Start: 5, End: 20}}, Loc{5, 20}},
		{"UpdateStmt", &UpdateStmt{Loc: Loc{Start: 3, End: 15}}, Loc{3, 15}},
		{"DeleteStmt", &DeleteStmt{Loc: Loc{Start: 1, End: 8}}, Loc{1, 8}},
		{"MergeStmt", &MergeStmt{Loc: Loc{Start: 0, End: 50}}, Loc{0, 50}},
		{"RangeVar", &RangeVar{Loc: Loc{Start: 10, End: 20}}, Loc{10, 20}},
		{"WithClause", &WithClause{Loc: Loc{Start: 0, End: 30}}, Loc{0, 30}},
		{"ColumnRef", &ColumnRef{Loc: Loc{Start: 7, End: 8}}, Loc{7, 8}},
		{"FuncCall", &FuncCall{Loc: Loc{Start: 2, End: 12}}, Loc{2, 12}},
		{"A_Const", &A_Const{Loc: Loc{Start: 4, End: 5}}, Loc{4, 5}},
		{"A_Expr", &A_Expr{Loc: Loc{Start: 0, End: 9}}, Loc{0, 9}},
		{"RawStmt", &RawStmt{Loc: Loc{Start: 0, End: 25}}, Loc{0, 25}},
		{"TypeCast", &TypeCast{Loc: Loc{Start: 1, End: 10}}, Loc{1, 10}},
		{"SubLink", &SubLink{Loc: Loc{Start: 0, End: 15}}, Loc{0, 15}},
		{"ResTarget", &ResTarget{Loc: Loc{Start: 7, End: 12}}, Loc{7, 12}},
		{"SortBy", &SortBy{Loc: Loc{Start: 30, End: 35}}, Loc{30, 35}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NodeLoc(tc.node)
			if got != tc.loc {
				t.Errorf("NodeLoc(%s) = %+v, want %+v", tc.name, got, tc.loc)
			}
		})
	}
}

func TestNodeLocUnknownType(t *testing.T) {
	got := NodeLoc(&FromExpr{})
	if got != NoLoc() {
		t.Errorf("NodeLoc(FromExpr) = %+v, want NoLoc", got)
	}
}

func TestNodeLocNil(t *testing.T) {
	got := NodeLoc(nil)
	if got != NoLoc() {
		t.Errorf("NodeLoc(nil) = %+v, want NoLoc", got)
	}
}

func TestListSpan(t *testing.T) {
	list := &List{Items: []Node{
		&RangeVar{Relname: "a", Loc: Loc{Start: 10, End: 11}},
		&RangeVar{Relname: "b", Loc: Loc{Start: 13, End: 14}},
	}}
	got := ListSpan(list)
	if got.Start != 10 || got.End != 14 {
		t.Errorf("ListSpan = %+v, want {10, 14}", got)
	}
}

func TestListSpanSingle(t *testing.T) {
	list := &List{Items: []Node{
		&ColumnRef{Loc: Loc{Start: 5, End: 8}},
	}}
	got := ListSpan(list)
	if got.Start != 5 || got.End != 8 {
		t.Errorf("ListSpan = %+v, want {5, 8}", got)
	}
}

func TestListSpanNil(t *testing.T) {
	got := ListSpan(nil)
	if got != NoLoc() {
		t.Errorf("ListSpan(nil) = %+v, want NoLoc", got)
	}
}

func TestListSpanEmpty(t *testing.T) {
	got := ListSpan(&List{})
	if got != NoLoc() {
		t.Errorf("ListSpan(empty) = %+v, want NoLoc", got)
	}
}
