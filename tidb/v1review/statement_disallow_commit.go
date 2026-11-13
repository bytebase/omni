package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowCommit reports each top-level COMMIT statement; a COMMIT
// nested in a routine body is not judged. Findings anchor nothing (the V1
// advisor placed them on the statement's first line). Commit is the only
// code. Ported from Bytebase's V1 rule STATEMENT_DISALLOW_COMMIT.
func StatementDisallowCommit(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if _, ok := stmt.(*ast.CommitStmt); ok {
			findings = append(findings, v1review.Finding{Code: v1review.Commit})
		}
	}
	return findings
}
