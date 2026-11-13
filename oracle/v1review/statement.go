package v1review

import "github.com/bytebase/omni/oracle/ast"

// plsqlBlockStatements calls visit for each statement nested in block, in
// walk order, descending into a statement only when visit returns true.
func plsqlBlockStatements(block *ast.PLSQLBlock, visit func(ast.StmtNode) bool) {
	if block == nil {
		return
	}
	visited := make(map[ast.Node]bool)
	ast.Inspect(block, func(node ast.Node) bool {
		if node == nil || visited[node] {
			return false
		}
		visited[node] = true
		if node == block {
			return true
		}
		stmt, ok := node.(ast.StmtNode)
		if !ok {
			return true
		}
		return visit(stmt)
	})
}

// anchor returns a Loc as the byte range a statement rule's finding anchors.
// The V1 statement advisors took a node's line from where it starts alone,
// so a node with a start and no end anchors at its start, and only a node
// with no start anchors nowhere (0,0).
func anchor(loc ast.Loc) (int, int) {
	if loc.Start < 0 {
		return 0, 0
	}
	if loc.End < loc.Start {
		return loc.Start, loc.Start
	}
	return loc.Start, loc.End
}
