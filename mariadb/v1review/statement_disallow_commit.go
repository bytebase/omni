package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowCommit reports each top-level COMMIT statement; a COMMIT
// nested in a routine body is not judged. Each finding anchors the COMMIT
// statement. Commit is the only code. Ported from Bytebase's V1 rule
// STATEMENT_DISALLOW_COMMIT.
func StatementDisallowCommit(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		commit, ok := stmt.(*ast.CommitStmt)
		if !ok {
			continue
		}
		start, end := statementAnchor(commit.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.Commit, Start: start, End: end})
	}
	return findings
}
