package v1review

import (
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemFunctionDisallowCreate reports each top-level CREATE FUNCTION,
// CREATE OR ALTER FUNCTION and ALTER FUNCTION, named or not. A statement
// nested in a routine body or an IF or BEGIN block is not judged. Every
// finding anchors the whole statement; FunctionCreated is the only code.
// Ported from Bytebase's V1 rule SYSTEM_FUNCTION_DISALLOW_CREATE.
func SystemFunctionDisallowCreate(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if n, ok := stmt.(*ast.CreateFunctionStmt); ok {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.FunctionCreated, Start: start, End: end})
		}
	}
	return findings
}
