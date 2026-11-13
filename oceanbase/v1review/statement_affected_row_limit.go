package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
	"strconv"
)

// StatementAffectedRowLimitProbes returns, in order, the indexes of the
// top-level UPDATE and DELETE statements: the ones the V1 rule
// STATEMENT_AFFECTED_ROW_LIMIT has the database estimate. Bytebase EXPLAINs
// each and passes the results to StatementAffectedRowLimit.
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
// EXPLAIN failed, RowCountFailed when its row count could not be read
// (Value the error for both), RowLimitExceeded when the estimate exceeds
// maxRows (Value the rows). A statement with no probe result was not
// estimated and is not reported. Statement is the index; findings anchor the
// UPDATE or DELETE. Ported from Bytebase's V1 rule
// STATEMENT_AFFECTED_ROW_LIMIT.
func StatementAffectedRowLimit(stmts []ast.Node, probes map[int]v1review.Probe, maxRows int) []v1review.Finding {
	var findings []v1review.Finding
	for _, i := range StatementAffectedRowLimitProbes(stmts) {
		probe, ok := probes[i]
		if !ok {
			continue
		}
		var loc ast.Loc
		switch n := stmts[i].(type) {
		case *ast.UpdateStmt:
			loc = n.Loc
		case *ast.DeleteStmt:
			loc = n.Loc
		}
		if f, ok := rowLimitFinding(probe, maxRows, i); ok {
			f.Start, f.End = anchor(loc)
			findings = append(findings, f)
		}
	}
	return findings
}

// rowLimitFinding judges one statement's row estimate against maxRows.
func rowLimitFinding(probe v1review.Probe, maxRows, statement int) (v1review.Finding, bool) {
	switch {
	case probe.Error != "":
		return v1review.Finding{Code: v1review.DryRunFailed, Value: probe.Error, Statement: statement}, true
	case probe.CountError != "":
		return v1review.Finding{Code: v1review.RowCountFailed, Value: probe.CountError, Statement: statement}, true
	case probe.Rows > int64(maxRows):
		return v1review.Finding{Code: v1review.RowLimitExceeded, Value: strconv.FormatInt(probe.Rows, 10), Statement: statement}, true
	}
	return v1review.Finding{}, false
}
