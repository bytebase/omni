package v1review

import (
	"slices"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowDDL judges the tables every top-level CREATE TABLE, ALTER
// TABLE, DROP TABLE, RENAME TABLE and TRUNCATE names; statements nested in
// a routine body are not judged. A table whose name, without its database,
// equals an entry of disallow exactly (case-sensitively) is
// TableDisallowed, with Name the table name. DROP TABLE and TRUNCATE judge
// each table they name and RENAME TABLE each pair's old name, not the new
// one. Every finding anchors the whole statement. Ported from Bytebase's V1
// rule TABLE_DISALLOW_DDL.
func TableDisallowDDL(stmts []ast.Node, disallow []string) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		var names []string
		var loc ast.Loc
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			loc = n.Loc
			if n.Table != nil {
				names = append(names, n.Table.Name)
			}
		case *ast.AlterTableStmt:
			loc = n.Loc
			if n.Table != nil {
				names = append(names, n.Table.Name)
			}
		case *ast.DropTableStmt:
			loc = n.Loc
			for _, table := range n.Tables {
				names = append(names, table.Name)
			}
		case *ast.RenameTableStmt:
			loc = n.Loc
			for _, pair := range n.Pairs {
				if pair.Old != nil {
					names = append(names, pair.Old.Name)
				}
			}
		case *ast.TruncateStmt:
			loc = n.Loc
			for _, table := range n.Tables {
				names = append(names, table.Name)
			}
		default:
			continue
		}
		start, end := span(loc)
		for _, name := range names {
			if slices.Contains(disallow, name) {
				findings = append(findings, v1review.Finding{Code: v1review.TableDisallowed, Name: name, Start: start, End: end})
			}
		}
	}
	return findings
}
