package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertRowLimitProbes returns, in order, the indexes of the INSERT
// statements, judged through an EXPLAIN ANALYZE to the statement it runs,
// whose rows come from a query rather than a VALUES list: the ones the V1
// rule STATEMENT_INSERT_ROW_LIMIT has the database estimate. Bytebase
// EXPLAINs each and passes the results to StatementInsertRowLimit.
func StatementInsertRowLimitProbes(stmts []ast.Node) []int {
	var probes []int
	for i, stmt := range stmts {
		insert, ok := liveTarget(stmt).(*ast.InsertStmt)
		if !ok {
			continue
		}
		if _, values := valuesRows(insert); !values && insert.SelectStmt != nil {
			probes = append(probes, i)
		}
	}
	return probes
}

// StatementInsertRowLimit reports, in statement order, each INSERT (judged
// through an EXPLAIN ANALYZE) that inserts more than maxRows rows: an INSERT
// ... VALUES by its count of row lists, one whose rows come from a query by
// its probe, as StatementAffectedRowLimit judges one. A query-sourced INSERT
// with no probe result was not estimated and is not reported. Count is the
// rows; Statement the index; findings anchor the statement (Start = End =
// 0). Ported from Bytebase's V1 rule STATEMENT_INSERT_ROW_LIMIT.
func StatementInsertRowLimit(stmts []ast.Node, probes map[int]v1review.Probe, maxRows int) []v1review.Finding {
	var findings []v1review.Finding
	for i, stmt := range stmts {
		insert, ok := liveTarget(stmt).(*ast.InsertStmt)
		if !ok {
			continue
		}
		if rows, values := valuesRows(insert); values {
			if rows > maxRows {
				findings = append(findings, v1review.Finding{Code: v1review.RowLimitExceeded, Count: rows, Statement: i})
			}
			continue
		}
		if probe, ok := probes[i]; ok && insert.SelectStmt != nil {
			if f, ok := rowLimitFinding(probe, maxRows, i); ok {
				findings = append(findings, f)
			}
		}
	}
	return findings
}

// valuesRows returns the row lists of an INSERT ... VALUES, and whether it is
// one.
func valuesRows(insert *ast.InsertStmt) (int, bool) {
	if sel, ok := insert.SelectStmt.(*ast.SelectStmt); ok && sel.ValuesLists != nil {
		return len(sel.ValuesLists.Items), true
	}
	return 0, false
}
