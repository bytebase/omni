package ast

import (
	"reflect"
	"testing"
)

func TestWalkSelectStmt(t *testing.T) {
	stmt := &SelectStmt{
		TargetList: &List{Items: []Node{
			&ResTarget{Expr: &ColumnRef{Column: "id"}},
			&ResTarget{Expr: &ColumnRef{Column: "name"}},
		}},
		FromClause: &List{Items: []Node{
			&TableRef{Name: &ObjectName{Name: "users"}},
		}},
		WhereClause: &BinaryExpr{
			Op:    "=",
			Left:  &ColumnRef{Column: "id"},
			Right: &NumberLiteral{Val: "1", Ival: 1},
		},
		OrderBy: &List{Items: []Node{
			&SortBy{Expr: &ColumnRef{Column: "name"}},
		}},
		FetchFirst: &FetchFirstClause{
			Count: &NumberLiteral{Val: "10", Ival: 10},
		},
	}

	var visited []string
	Inspect(stmt, func(n Node) bool {
		if n == nil {
			return false
		}
		visited = append(visited, reflect.TypeOf(n).Elem().Name())
		return true
	})

	typeSet := map[string]bool{}
	for _, v := range visited {
		typeSet[v] = true
	}
	for _, want := range []string{"SelectStmt", "List", "ResTarget", "ColumnRef", "TableRef", "ObjectName", "BinaryExpr", "NumberLiteral", "SortBy", "FetchFirstClause"} {
		if !typeSet[want] {
			t.Errorf("expected to visit %s, visited: %v", want, visited)
		}
	}
}

func TestWalkNil(t *testing.T) {
	Walk(inspector(func(n Node) bool { return true }), nil)
}

func TestInspectPruning(t *testing.T) {
	stmt := &SelectStmt{
		WhereClause: &BinaryExpr{
			Op:    "=",
			Left:  &ColumnRef{Column: "id"},
			Right: &NumberLiteral{Val: "1", Ival: 1},
		},
	}

	var visited []string
	Inspect(stmt, func(n Node) bool {
		if n == nil {
			return false
		}
		name := reflect.TypeOf(n).Elem().Name()
		visited = append(visited, name)
		if name == "BinaryExpr" {
			return false
		}
		return true
	})

	typeSet := map[string]bool{}
	for _, v := range visited {
		typeSet[v] = true
	}
	if !typeSet["SelectStmt"] {
		t.Error("expected SelectStmt")
	}
	if !typeSet["BinaryExpr"] {
		t.Error("expected BinaryExpr")
	}
	if typeSet["ColumnRef"] {
		t.Error("ColumnRef should have been pruned")
	}
	if typeSet["NumberLiteral"] {
		t.Error("NumberLiteral should have been pruned")
	}
}
