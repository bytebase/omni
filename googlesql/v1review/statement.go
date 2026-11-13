// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for BigQuery, which parses with GoogleSQL. The
// contract, and the split between these functions and Bytebase's V1
// product surface, is github.com/bytebase/omni/v1review.
package v1review

import "github.com/bytebase/omni/googlesql/ast"

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
