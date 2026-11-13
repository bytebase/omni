package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementAffectedRowLimitProbes returns, in order, the indexes of the
// top-level UPDATE and DELETE statements: the ones the V1 rule
// STATEMENT_AFFECTED_ROW_LIMIT has the database estimate. Bytebase EXPLAINs
// and estimates each and passes the results to StatementAffectedRowLimit.
func StatementAffectedRowLimitProbes(stmts []ast.Node) []int {
	var probes []int
	for i, stmt := range stmts {
		switch stmt.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			probes = append(probes, i)
		}
	}
	return probes
}

// StatementAffectedRowLimit reports, in statement order, each probed
// statement whose estimate failed or exceeds maxRows: DryRunFailed when the
// EXPLAIN failed, RowCountFailed when the estimate failed (Value the error
// for both), RowLimitExceeded when the estimate exceeds maxRows (Value the
// rows). A statement with no probe result was not estimated and is not
// reported. Statement is the index; findings anchor the statement (Start =
// End = 0). Ported from Bytebase's V1 rule STATEMENT_AFFECTED_ROW_LIMIT.
func StatementAffectedRowLimit(stmts []ast.Node, probes map[int]v1review.Probe, maxRows int) []v1review.Finding {
	var findings []v1review.Finding
	for _, i := range StatementAffectedRowLimitProbes(stmts) {
		if probe, ok := probes[i]; ok {
			if f, ok := rowLimitFinding(probe, maxRows, i); ok {
				findings = append(findings, f)
			}
		}
	}
	return findings
}
