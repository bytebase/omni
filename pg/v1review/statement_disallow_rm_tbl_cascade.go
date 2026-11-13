package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowRmTblCascade reports each top-level DROP TABLE or
// TRUNCATE with CASCADE, one DropCascade finding per statement. DROP of any
// other object kind is not judged. Findings anchor nothing: Bytebase's V1
// advisor placed them on the statement's first line. Ported from
// Bytebase's V1 rule STATEMENT_DISALLOW_RM_TBL_CASCADE.
func StatementDisallowRmTblCascade(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var cascade bool
		switch n := unwrapRaw(stmt).(type) {
		case *ast.DropStmt:
			cascade = n.RemoveType == int(ast.OBJECT_TABLE) && ast.DropBehavior(n.Behavior) == ast.DROP_CASCADE
		case *ast.TruncateStmt:
			cascade = n.Behavior == ast.DROP_CASCADE
		}
		if cascade {
			findings = append(findings, v1review.Finding{Code: v1review.DropCascade})
		}
	}
	return findings
}
