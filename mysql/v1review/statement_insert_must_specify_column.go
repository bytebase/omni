package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertMustSpecifyColumn reports each top-level INSERT (or
// REPLACE) that names no column list and takes its rows from a query or a
// VALUES list; INSERT ... SET is not judged. A statement yields at most one
// finding, which anchors nothing (the V1 advisor placed it on the
// statement's first line). InsertColumnsUnspecified is the only code.
// Ported from Bytebase's V1 rule STATEMENT_INSERT_MUST_SPECIFY_COLUMN.
func StatementInsertMustSpecifyColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		ins, ok := stmt.(*ast.InsertStmt)
		if !ok || len(ins.Columns) > 0 {
			continue
		}
		if ins.Select != nil || len(ins.Values) > 0 {
			findings = append(findings, v1review.Finding{Code: v1review.InsertColumnsUnspecified})
		}
	}
	return findings
}
