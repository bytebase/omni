package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementInsertMustSpecifyColumn reports each single-table INSERT without
// a column list, and each INTO clause of a multi-table INSERT ALL or FIRST
// without one, among the top-level statements and the statements nested in
// a top-level anonymous PL/SQL block. INSERTs in a procedure, function,
// package or trigger body are not judged. A single-table finding anchors the
// INSERT and a multi-table one its INTO clause. InsertColumnsUnspecified is
// the only code. A statement may be passed bare or wrapped in an
// *ast.RawStmt. Ported from Bytebase's V1 rule
// STATEMENT_INSERT_MUST_SPECIFY_COLUMN.
func StatementInsertMustSpecifyColumn(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.InsertStmt:
			findings = insertColumnsUnspecified(findings, n)
		case *ast.PLSQLBlock:
			plsqlBlockStatements(n, func(s ast.StmtNode) bool {
				insert, ok := s.(*ast.InsertStmt)
				if !ok {
					return true
				}
				findings = insertColumnsUnspecified(findings, insert)
				return false
			})
		}
	}
	return findings
}

func insertColumnsUnspecified(findings []v1review.Finding, n *ast.InsertStmt) []v1review.Finding {
	if n.InsertType == ast.INSERT_SINGLE && n.Columns == nil {
		start, end := span(n.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.InsertColumnsUnspecified, Start: start, End: end})
	}
	if n.MultiTable == nil {
		return findings
	}
	for _, item := range n.MultiTable.Items {
		if clause, ok := item.(*ast.InsertIntoClause); ok && clause.Columns == nil {
			start, end := span(clause.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.InsertColumnsUnspecified, Start: start, End: end})
		}
	}
	return findings
}
