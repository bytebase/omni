package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingAutoIncrementColumn judges the names of AUTO_INCREMENT columns in
// top-level statements against cfg: the columns a CREATE TABLE defines, and
// in an ALTER TABLE the columns of ADD COLUMN in both its forms, MODIFY
// COLUMN and CHANGE COLUMN. A column counts when its definition says
// AUTO_INCREMENT or SERIAL DEFAULT VALUE; a column of type SERIAL does not
// count, and statements nested in a routine body are not judged. The name
// is judged as parsed, its length in bytes; Table is the table's name
// without its database qualifier, empty when the ALTER TABLE has none.
//
// Every finding anchors the whole statement, a CREATE TABLE's included.
// Findings are in statement order, a column's mismatch before its length.
// Ported from Bytebase's V1 rule NAMING_COLUMN_AUTO_INCREMENT.
func NamingAutoIncrementColumn(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
	pattern, err := cfg.Compile()
	if err != nil {
		return nil, err
	}
	if pattern == nil && cfg.MaxLength <= 0 {
		return nil, nil
	}
	naming := &columnNaming{pattern: pattern, max: cfg.MaxLength}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, col := range n.Columns {
				if col != nil && col.AutoIncrement {
					naming.check(n.Loc, n.Table.Name, col.Name)
				}
			}
		case *ast.AlterTableStmt:
			table := ""
			if n.Table != nil {
				table = n.Table.Name
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				for _, col := range commandColumns(cmd) {
					if col.AutoIncrement {
						naming.check(n.Loc, table, col.Name)
					}
				}
			}
		}
	}
	return naming.findings, nil
}
