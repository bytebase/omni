package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertMustSpecifyColumn reports each top-level INSERT that names
// no column list and inserts from a SELECT or a VALUES list, once per
// statement; an INSERT ... SET, and statements nested in a routine body, are
// not judged. Findings anchor nothing (the V1 advisor took the statement's
// first content line). InsertColumnsUnspecified is the only code. Ported
// from Bytebase's V1 rule STATEMENT_INSERT_MUST_SPECIFY_COLUMN.
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
