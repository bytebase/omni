package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertRowLimitProbes returns, in order, the indexes of the
// top-level INSERT and REPLACE statements whose rows come from a query or a
// table: the ones the V1 rule STATEMENT_INSERT_ROW_LIMIT has the database
// estimate. Bytebase EXPLAINs each and passes the results to
// StatementInsertRowLimit.
func StatementInsertRowLimitProbes(stmts []ast.Node) []int {
	var probes []int
	for i, stmt := range stmts {
		if insert, ok := stmt.(*ast.InsertStmt); ok && (insert.Select != nil || insert.TableSource != nil) {
			probes = append(probes, i)
		}
	}
	return probes
}

// StatementInsertRowLimit reports, in statement order, each top-level INSERT
// or REPLACE that inserts more than maxRows rows. A probed INSERT is first
// judged from its probe, as StatementAffectedRowLimit judges one; then every
// INSERT, probed or not, is judged by its count of VALUES row lists, so one
// INSERT can give two findings. A probed INSERT with no probe result gives
// no probe finding. Count is the rows; Statement the index; findings anchor
// the INSERT. Ported from Bytebase's V1 rule STATEMENT_INSERT_ROW_LIMIT.
func StatementInsertRowLimit(stmts []ast.Node, probes map[int]v1review.Probe, maxRows int) []v1review.Finding {
	var findings []v1review.Finding
	for i, stmt := range stmts {
		insert, ok := stmt.(*ast.InsertStmt)
		if !ok {
			continue
		}
		start, end := anchor(insert.Loc)
		if insert.Select != nil || insert.TableSource != nil {
			if probe, ok := probes[i]; ok {
				if f, ok := rowLimitFinding(probe, maxRows, i); ok {
					f.Start, f.End = start, end
					findings = append(findings, f)
				}
			}
		}
		if len(insert.Values) > maxRows {
			findings = append(findings, v1review.Finding{Code: v1review.RowLimitExceeded, Count: len(insert.Values), Statement: i, Start: start, End: end})
		}
	}
	return findings
}
