package v1review

import (
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTableNoKeyword reports the table names top-level statements
// introduce that are keywords in the V1 advisor's table (see isKeyword),
// compared without regard to case: the table a CREATE TABLE creates, CTAS
// included, and the new name an ALTER TABLE ... RENAME TO gives. A RENAME
// statement is not judged, nor is a CREATE TABLE nested in CREATE SCHEMA.
// The name is judged as the parser stores it — an unquoted name
// upper-cased, a quoted one as written — without its schema qualifier. A
// CREATE TABLE finding anchors the statement; an ALTER TABLE finding
// anchors its RENAME clause. A statement may be passed bare or wrapped in
// an *ast.RawStmt. Ported from Bytebase's V1 rule NAMING_TABLE_NO_KEYWORD.
func NamingTableNoKeyword(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	check := func(name string, loc ast.Loc) {
		if isKeyword(name) {
			start, end := span(loc)
			findings = append(findings, v1review.Finding{Code: v1review.NameIsKeyword, Name: name, Start: start, End: end})
		}
	}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Name != nil {
				check(n.Name.Name, n.Loc)
			}
		case *ast.AlterTableStmt:
			if n.Actions == nil {
				continue
			}
			for _, item := range n.Actions.Items {
				if cmd, ok := item.(*ast.AlterTableCmd); ok && cmd.Action == ast.AT_RENAME {
					check(cmd.NewName, cmd.Loc)
				}
			}
		}
	}
	return findings
}
