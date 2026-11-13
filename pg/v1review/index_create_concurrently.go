package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// IndexCreateConcurrently reports the top-level statements that build or drop
// an index without CONCURRENTLY: a CREATE INDEX, IndexCreatedNotConcurrently,
// unless its table was created earlier in the script; and a DROP INDEX,
// IndexDroppedNotConcurrently, whatever its table. A statement in a CREATE
// SCHEMA or in a routine body is not judged.
//
// A table counts as created by an earlier CREATE TABLE of the same name,
// compared without schemas, so a table of that name in another schema
// counts too; CREATE TABLE AS and CREATE FOREIGN TABLE do not create one.
//
// A CREATE INDEX finding's Name is the index's name, empty when it has none,
// and Table its table's name without its schema. Statement is the finding's
// statement, since what one statement created changes how the next is
// judged; every finding anchors that whole statement, and Bytebase's V1
// advisor reported the statement text's first non-blank line. Findings are
// in statement order. Ported from Bytebase's V1 rule
// INDEX_CREATE_CONCURRENTLY.
func IndexCreateConcurrently(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	created := make(map[string]bool)
	for i, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if table := relname(n.Relation); table != "" {
				created[table] = true
			}
		case *ast.IndexStmt:
			table := relname(n.Relation)
			if n.Concurrent || created[table] {
				continue
			}
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.IndexCreatedNotConcurrently, Name: n.Idxname, Table: table, Statement: i, Start: start, End: end})
		case *ast.DropStmt:
			if n.RemoveType != int(ast.OBJECT_INDEX) || n.Concurrent {
				continue
			}
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.IndexDroppedNotConcurrently, Statement: i, Start: start, End: end})
		}
	}
	return findings
}
