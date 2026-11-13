package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertMustSpecifyColumn reports each top-level INSERT that names
// no column list, one InsertColumnsUnspecified finding per statement,
// whatever its source: VALUES, SELECT or DEFAULT VALUES. An INSERT in a
// WITH query or a routine body is not judged. Findings anchor nothing:
// Bytebase's V1 advisor placed them on the statement's first line. Ported
// from Bytebase's V1 rule STATEMENT_INSERT_MUST_SPECIFY_COLUMN.
func StatementInsertMustSpecifyColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		insert, ok := unwrapRaw(stmt).(*ast.InsertStmt)
		if !ok {
			continue
		}
		if insert.Cols == nil || len(insert.Cols.Items) == 0 {
			findings = append(findings, v1review.Finding{Code: v1review.InsertColumnsUnspecified})
		}
	}
	return findings
}
