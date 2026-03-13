package ast

import "testing"

func TestWalkSelectStmt(t *testing.T) {
	// SELECT a, b FROM t WHERE x > 0
	tree := &SelectStmt{
		TargetList: &List{Items: []Node{
			&ResTarget{Val: &ColumnRef{Fields: &List{Items: []Node{&String{Str: "a"}}}}},
			&ResTarget{Val: &ColumnRef{Fields: &List{Items: []Node{&String{Str: "b"}}}}},
		}},
		FromClause: &List{Items: []Node{
			&RangeVar{Relname: "t"},
		}},
		WhereClause: &A_Expr{
			Kind:  AEXPR_OP,
			Name:  &List{Items: []Node{&String{Str: ">"}}},
			Lexpr: &ColumnRef{Fields: &List{Items: []Node{&String{Str: "x"}}}},
			Rexpr: &A_Const{Val: &Integer{Ival: 0}},
		},
	}

	var tags []NodeTag
	Inspect(tree, func(n Node) bool {
		if n != nil {
			tags = append(tags, n.Tag())
		}
		return true
	})

	// Should have visited: SelectStmt, List(target), ResTarget, ColumnRef, List, String("a"),
	//   ResTarget, ColumnRef, List, String("b"), List(from), RangeVar,
	//   A_Expr, List(name), String(">"), ColumnRef, List, String("x"), A_Const, Integer
	if len(tags) < 15 {
		t.Errorf("expected at least 15 nodes visited, got %d", len(tags))
	}

	// Check first node is SelectStmt
	if tags[0] != T_SelectStmt {
		t.Errorf("first node should be SelectStmt, got %s", NodeTagName(tags[0]))
	}

	// Verify specific types are present
	tagSet := map[NodeTag]bool{}
	for _, tag := range tags {
		tagSet[tag] = true
	}
	for _, want := range []NodeTag{T_SelectStmt, T_List, T_ResTarget, T_ColumnRef, T_String, T_RangeVar, T_A_Expr, T_A_Const, T_Integer} {
		if !tagSet[want] {
			t.Errorf("expected %s to be visited", NodeTagName(want))
		}
	}
}

func TestWalkNil(t *testing.T) {
	// Should not panic on nil.
	Walk(inspector(func(Node) bool { return true }), nil)
}

func TestInspectPruning(t *testing.T) {
	tree := &SelectStmt{
		TargetList: &List{Items: []Node{
			&ResTarget{Val: &ColumnRef{Fields: &List{Items: []Node{&String{Str: "a"}}}}},
		}},
		WhereClause: &A_Expr{
			Lexpr: &A_Const{Val: &Integer{Ival: 1}},
		},
	}

	// Stop at A_Expr — should not visit its children.
	var visited []NodeTag
	Inspect(tree, func(n Node) bool {
		if n == nil {
			return false
		}
		visited = append(visited, n.Tag())
		if n.Tag() == T_A_Expr {
			return false // prune
		}
		return true
	})

	// A_Const and Integer should NOT be visited (pruned under A_Expr)
	for _, tag := range visited {
		if tag == T_A_Const || tag == T_Integer {
			t.Errorf("should not have visited %s (pruned)", NodeTagName(tag))
		}
	}
	// But A_Expr itself should be visited
	found := false
	for _, tag := range visited {
		if tag == T_A_Expr {
			found = true
		}
	}
	if !found {
		t.Error("A_Expr should have been visited")
	}
}

func TestVisitorPostOrder(t *testing.T) {
	tree := &RangeVar{Relname: "t"}

	type event struct {
		tag  NodeTag
		post bool
	}
	var events []event

	v := &postOrderVisitor{onVisit: func(n Node) {
		if n == nil {
			events = append(events, event{post: true})
		} else {
			events = append(events, event{tag: n.Tag()})
		}
	}}
	Walk(v, tree)

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].tag != T_RangeVar || events[0].post {
		t.Error("first event should be pre-order RangeVar")
	}
	if !events[1].post {
		t.Error("second event should be post-order nil")
	}
}

type postOrderVisitor struct {
	onVisit func(Node)
}

func (v *postOrderVisitor) Visit(node Node) Visitor {
	v.onVisit(node)
	if node != nil {
		return v
	}
	return nil
}
