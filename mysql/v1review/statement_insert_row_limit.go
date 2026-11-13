package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertRowLimitProbes returns, in order, the indexes of the
// top-level INSERT and REPLACE statements whose rows come from a query
// (INSERT ... SELECT, INSERT ... TABLE): the ones the V1 rule
// STATEMENT_INSERT_ROW_LIMIT has the database estimate. Bytebase EXPLAINs and
// estimates each and passes the results to StatementInsertRowLimit.
func StatementInsertRowLimitProbes(stmts []ast.Node) []int {
	var probes []int
	for i, stmt := range stmts {
		if insert, ok := stmt.(*ast.InsertStmt); ok && insert.Select != nil {
			probes = append(probes, i)
		}
	}
	return probes
}

// StatementInsertRowLimit reports, in statement order, each top-level
// INSERT or REPLACE that inserts more than maxRows rows. An INSERT ... VALUES
// counts its row lists; one whose rows come from a query is judged from its
// probe: DryRunFailed when the EXPLAIN failed (Detail the error),
// RowCountFailed when the estimate failed (Detail the error), and
// RowLimitExceeded when the estimate exceeds maxRows. A query-sourced INSERT
// with no probe result was not estimated and is not reported. Count is the
// rows; Statement the index; findings anchor the statement (Start = End =
// 0). Ported from Bytebase's V1 rule STATEMENT_INSERT_ROW_LIMIT.
func StatementInsertRowLimit(stmts []ast.Node, probes map[int]v1review.Probe, maxRows int) []v1review.Finding {
	var findings []v1review.Finding
	for i, stmt := range stmts {
		insert, ok := stmt.(*ast.InsertStmt)
		if !ok {
			continue
		}
		if insert.Select == nil {
			if len(insert.Values) > maxRows {
				findings = append(findings, v1review.Finding{Code: v1review.RowLimitExceeded, Count: len(insert.Values), Statement: i})
			}
			continue
		}
		if probe, ok := probes[i]; ok {
			if f, ok := rowLimitFinding(probe, maxRows, i); ok {
				findings = append(findings, f)
			}
		}
	}
	return findings
}

// rowLimitFinding judges one statement's row estimate against maxRows.
func rowLimitFinding(probe v1review.Probe, maxRows, statement int) (v1review.Finding, bool) {
	switch {
	case probe.Error != "":
		return v1review.Finding{Code: v1review.DryRunFailed, Detail: probe.Error, Statement: statement}, true
	case probe.CountError != "":
		return v1review.Finding{Code: v1review.RowCountFailed, Detail: probe.CountError, Statement: statement}, true
	case probe.Rows > int64(maxRows):
		return v1review.Finding{Code: v1review.RowLimitExceeded, Count: int(probe.Rows), Statement: statement}, true
	}
	return v1review.Finding{}, false
}
