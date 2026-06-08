package ast

//go:generate go run ./cmd/genwalker

// Visitor defines the interface for AST traversal.
// Visit is called for each node during a depth-first walk.
// If Visit returns a non-nil Visitor, Walk recurses into the node's children
// with the returned Visitor, then calls Visit(nil) to signal post-order.
// If Visit returns nil, children are not visited.
type Visitor interface {
	Visit(node Node) Visitor
}

// Walk traverses an AST in depth-first order. It calls v.Visit(node);
// if that returns a non-nil visitor w, it walks each child node with w,
// then calls w.Visit(nil).
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

// Inspect traverses an AST in depth-first order, calling f for each node.
// If f returns true, Inspect recurses into the node's children.
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

// walkList visits a List node and then walks each of its items.
func walkList(v Visitor, list *List) {
	if list == nil {
		return
	}
	w := v.Visit(list)
	if w == nil {
		return
	}
	for _, item := range list.Items {
		Walk(w, item)
	}
	w.Visit(nil)
}
