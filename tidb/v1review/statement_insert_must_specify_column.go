package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertMustSpecifyColumn reports each top-level INSERT or REPLACE
// that names no columns and takes its rows from VALUES or a query. An empty
// column list () counts as none; the SET form is not judged, nor is an
// INSERT nested in a routine body. A finding anchors the statement.
// InsertColumnsUnspecified is the only code. Ported from Bytebase's V1 rule
// STATEMENT_INSERT_MUST_SPECIFY_COLUMN.
func StatementInsertMustSpecifyColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.InsertStmt)
		if !ok || len(n.Columns) > 0 || (len(n.Values) == 0 && n.Select == nil) {
			continue
		}
		start, end := stmtAnchor(n.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.InsertColumnsUnspecified, Start: start, End: end})
	}
	return findings
}
