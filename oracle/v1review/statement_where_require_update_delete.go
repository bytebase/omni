package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereRequireUpdateDelete reports each UPDATE and DELETE without a
// WHERE clause among the top-level statements and the statements nested in
// a top-level anonymous PL/SQL block, in walk order. Ones in a procedure,
// function, package or trigger body are not judged, nor are MERGE's UPDATE
// and DELETE clauses. A finding anchors the statement. WhereMissing is the
// only code; Name is "UPDATE" or "DELETE". A statement may be passed bare
// or wrapped in an *ast.RawStmt. Ported from Bytebase's V1 rule
// STATEMENT_WHERE_REQUIRE_UPDATE_DELETE.
func StatementWhereRequireUpdateDelete(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		if block, ok := stmt.(*ast.PLSQLBlock); ok {
			plsqlBlockStatements(block, func(s ast.StmtNode) bool {
				switch s.(type) {
				case *ast.UpdateStmt, *ast.DeleteStmt:
					findings = whereMissing(findings, s)
					return false
				default:
					return true
				}
			})
			continue
		}
		findings = whereMissing(findings, stmt)
	}
	return findings
}

func whereMissing(findings []v1review.Finding, stmt ast.Node) []v1review.Finding {
	switch n := stmt.(type) {
	case *ast.UpdateStmt:
		if n.WhereClause == nil {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "UPDATE", Start: start, End: end})
		}
	case *ast.DeleteStmt:
		if n.WhereClause == nil {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.WhereMissing, Name: "DELETE", Start: start, End: end})
		}
	}
	return findings
}
