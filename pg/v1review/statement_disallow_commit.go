package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementDisallowCommit reports each top-level COMMIT, one Commit finding
// per statement; END is parsed as a COMMIT and counts, while COMMIT
// PREPARED, ROLLBACK and a COMMIT in a routine body do not. Findings anchor
// nothing: Bytebase's V1 advisor placed them on the statement's first line.
// Ported from Bytebase's V1 rule STATEMENT_DISALLOW_COMMIT.
func StatementDisallowCommit(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		txn, ok := unwrapRaw(stmt).(*ast.TransactionStmt)
		if ok && txn.Kind == ast.TRANS_STMT_COMMIT {
			findings = append(findings, v1review.Finding{Code: v1review.Commit})
		}
	}
	return findings
}
