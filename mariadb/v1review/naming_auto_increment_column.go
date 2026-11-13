package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingAutoIncrementColumn judges the names of the auto-increment columns
// top-level statements define against cfg: those of a CREATE TABLE, and
// those an ALTER TABLE defines with ADD [COLUMN], in its single and its
// parenthesized form, MODIFY [COLUMN] or CHANGE [COLUMN], whose new name is
// the one judged. A column is auto-increment when its definition says
// AUTO_INCREMENT or SERIAL DEFAULT VALUE; a column of type SERIAL, which
// the engine also makes auto-increment, is not judged. Statements nested in
// a routine body are not judged. The name is judged as parsed, its length
// in bytes, and the finding carries the table's name without its database
// qualifier.
//
// Every finding anchors the whole statement. Findings are in statement
// order, a column's mismatch before its length. Ported from Bytebase's V1
// rule NAMING_COLUMN_AUTO_INCREMENT.
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
				switch cmd.Type {
				case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
					for _, col := range columnDefs(cmd) {
						if col.AutoIncrement {
							naming.check(n.Loc, table, col.Name)
						}
					}
				default:
				}
			}
		}
	}
	return naming.findings, nil
}
