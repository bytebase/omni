package v1review

import "github.com/bytebase/omni/mysql/ast"

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
