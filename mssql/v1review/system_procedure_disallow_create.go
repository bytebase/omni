package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemProcedureDisallowCreate reports each top-level CREATE PROCEDURE,
// CREATE OR ALTER PROCEDURE and ALTER PROCEDURE, named or not. A statement
// nested in a routine body or an IF or BEGIN block is not judged. Every
// finding anchors the whole statement; ProcedureCreated is the only code.
// Ported from Bytebase's V1 rule SYSTEM_PROCEDURE_DISALLOW_CREATE.
func SystemProcedureDisallowCreate(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if n, ok := stmt.(*ast.CreateProcedureStmt); ok {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.ProcedureCreated, Start: start, End: end})
		}
	}
	return findings
}
