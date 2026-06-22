package ast

import (
	"reflect"
	"testing"
)

func TestWalkNil(t *testing.T) {
	var count int
	Inspect(nil, func(n Node) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("expected 0 visits for nil, got %d", count)
	}
}

func TestInspectSelectStmt(t *testing.T) {
	sel := &SelectStmt{
		Elements: []*SelectElement{
			{
				Expr:  &Identifier{Name: "name", Loc: Loc{Start: 7, End: 11}},
				Alias: &Identifier{Name: "n", Loc: Loc{Start: 15, End: 16}},
				Loc:   Loc{Start: 7, End: 16},
			},
		},
		From: &QualifiedName{
			Parts: []*Identifier{
				{Name: "users", Loc: Loc{Start: 22, End: 27}},
			},
			Loc: Loc{Start: 22, End: 27},
		},
		Where: []ExprNode{
			&BinaryExpr{
				Left:  &Identifier{Name: "id", Loc: Loc{Start: 34, End: 36}},
				Op:    "=",
				Right: &IntegerLit{Val: "1", Loc: Loc{Start: 39, End: 40}},
				Loc:   Loc{Start: 34, End: 40},
			},
		},
		Loc: Loc{Start: 0, End: 40},
	}

	var visited []string
	Inspect(sel, func(n Node) bool {
		if n == nil {
			return false
		}
		visited = append(visited, reflect.TypeOf(n).Elem().Name())
		return true
	})

	expected := []string{
		"SelectStmt",
		"SelectElement",
		"Identifier", // name
		"Identifier", // alias n
		"QualifiedName",
		"Identifier", // users
		"BinaryExpr",
		"Identifier", // id
		"IntegerLit", // 1
	}
	if len(visited) != len(expected) {
		t.Fatalf("expected %d visits, got %d: %v", len(expected), len(visited), visited)
	}
	for i, e := range expected {
		if visited[i] != e {
			t.Errorf("visit[%d] = %q, want %q", i, visited[i], e)
		}
	}
}

func TestInspectPruning(t *testing.T) {
	sel := &SelectStmt{
		From: &QualifiedName{
			Parts: []*Identifier{
				{Name: "users", Loc: Loc{Start: 14, End: 19}},
			},
			Loc: Loc{Start: 14, End: 19},
		},
		Loc: Loc{Start: 0, End: 19},
	}

	var visited []string
	Inspect(sel, func(n Node) bool {
		if n == nil {
			return false
		}
		name := reflect.TypeOf(n).Elem().Name()
		visited = append(visited, name)
		if name == "QualifiedName" {
			return false
		}
		return true
	})

	if len(visited) != 2 {
		t.Fatalf("expected 2 visits (SelectStmt, QualifiedName), got %d: %v", len(visited), visited)
	}
}
