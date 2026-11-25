package ast

//go:generate go run ./cmd/genwalker

// Visitor is implemented by clients of Walk.
//
// Visit is called for each node during a depth-first walk. If Visit returns
// a non-nil w, Walk recurses into the node's children with w, then calls
// w.Visit(nil) to signal end-of-children (post-order). If Visit returns nil,
// children are not visited.
//
// This is the same shape as pg/ast.Visitor, mysql/ast.Visitor, and
// snowflake/ast.Visitor.
type Visitor interface {
	Visit(node Node) Visitor
}

// Walk traverses an AST in depth-first order. It calls v.Visit(node);
// if that returns a non-nil visitor w, it walks each child node with w,
// then calls w.Visit(nil) to mark post-order.
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}
	w := v.Visit(node)
	if w == nil {
		return
	}
	walkChildren(w, node)
	w.Visit(nil)
}

// Inspect traverses an AST in depth-first order, calling f for each node in
// pre-order. If f returns true, Inspect recurses into the node's children; if
// it returns false, that subtree is pruned.
//
// Unlike the standard library's go/ast.Inspect, f is NOT invoked with a nil
// argument to mark post-order: the underlying inspector translates Walk's
// post-order Visit(nil) into a no-op (it returns nil without calling f). Use
// Walk directly with a Visitor if you need the post-order signal.
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if node != nil && f(node) {
		return f
	}
	return nil
}

// walkNodes visits each entry in nodes by calling Walk on it. Used by
// walk_generated.go to traverse slice-typed child fields like []Node.
func walkNodes(v Visitor, nodes []Node) {
	for _, n := range nodes {
		Walk(v, n)
	}
}
