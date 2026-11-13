package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// DatabaseDropEmptyDatabase judges each top-level DROP DATABASE against the
// pre-change schema: the dropped database must be the one the schema is of,
// its name compared exactly as written, and must then hold no table.
// Statements nested in a routine body are not judged.
//
// Name is the dropped database as written. Every finding anchors the whole
// statement. Ported from Bytebase's V1 rule DATABASE_DROP_EMPTY_DATABASE.
func DatabaseDropEmptyDatabase(stmts []ast.Node, schema v1review.Schema) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.DropDatabaseStmt)
		if !ok {
			continue
		}
		var code v1review.Code
		switch {
		case schema.DatabaseName() != n.Name:
			code = v1review.DatabaseNotCurrent
		case !schema.HasNoTable():
			code = v1review.DatabaseNotEmpty
		default:
			continue
		}
		start, end := span(n.Loc)
		findings = append(findings, v1review.Finding{Code: code, Name: n.Name, Start: start, End: end})
	}
	return findings
}
