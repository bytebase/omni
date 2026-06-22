package ast

// Visitor defines the interface for AST traversal.
type Visitor interface {
	Visit(node Node) Visitor
}

// Walk traverses an AST in depth-first order. It calls v.Visit(node); if that
// returns a non-nil Visitor w, Walk recurses into the children of node with w,
// then calls w.Visit(nil).
func Walk(v Visitor, node Node) {
	if node == nil {
		return
	}
	if v = v.Visit(node); v == nil {
		return
	}
	walkChildren(v, node)
	v.Visit(nil)
}

// Inspect traverses an AST in depth-first order, calling f for each node.
// If f returns true, Inspect recurses into the children of node.
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if node == nil || !f(node) {
		return nil
	}
	return f
}
