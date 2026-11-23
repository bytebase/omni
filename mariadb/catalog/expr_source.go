package catalog

import (
	"maps"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// exprSource is the AST a stored expression (a generated column's, a CHECK's,
// a partitioning function's) was declared with, and the column renames since,
// so a rename rewrites the stored text from the AST rather than parsing the
// text again. The AST is never changed: a snapshot of the table may share it.
type exprSource struct {
	node nodes.ExprNode
	// renames maps the lowered name a column reference was declared with to
	// the name the column has now.
	renames map[string]string
}

func newExprSource(node nodes.ExprNode) *exprSource {
	if node == nil {
		return nil
	}
	return &exprSource{node: node}
}

// rename follows a rename of the column oldName to newName. It returns the
// source after it and the expression's text, written as the declaration's
// was (nodeToSQL); ok is false when the expression reads no such column, or
// when no AST was kept.
func (s *exprSource) rename(oldName, newName string) (next *exprSource, text string, ok bool) {
	if s == nil || s.node == nil {
		return s, "", false
	}
	renames := maps.Clone(s.renames)
	if renames == nil {
		renames = map[string]string{}
	}
	nodes.Inspect(s.node, func(n nodes.Node) bool {
		ref, isRef := n.(*nodes.ColumnRef)
		if !isRef || ref.Star || ref.Column == "" {
			return true
		}
		declared := toLower(ref.Column)
		current := ref.Column
		if now, moved := renames[declared]; moved {
			current = now
		}
		if toLower(current) == toLower(oldName) {
			renames[declared] = newName
			ok = true
		}
		return true
	})
	if !ok {
		return s, "", false
	}
	next = &exprSource{node: s.node, renames: renames}
	return next, next.text(), true
}

// text writes the expression with its renames, setting each reference to its
// column's name for the length of the write and back after it.
func (s *exprSource) text() string {
	type saved struct {
		ref  *nodes.ColumnRef
		name string
	}
	var restore []saved
	nodes.Inspect(s.node, func(n nodes.Node) bool {
		if ref, isRef := n.(*nodes.ColumnRef); isRef && !ref.Star {
			if now, moved := s.renames[toLower(ref.Column)]; moved {
				restore = append(restore, saved{ref, ref.Column})
				ref.Column = now
			}
		}
		return true
	})
	out := nodeToSQL(s.node)
	for _, r := range restore {
		r.ref.Column = r.name
	}
	return out
}
