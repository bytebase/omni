package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingAutoIncrementColumn judges the names of the AUTO_INCREMENT columns
// top-level statements define against cfg: those a CREATE TABLE defines,
// and in an ALTER TABLE those an ADD [COLUMN], a MODIFY or a CHANGE
// defines. A column is judged only where its own definition says
// AUTO_INCREMENT, so RENAME COLUMN is not judged. Statements nested in a
// routine body are not judged.
//
// The name is judged as written, its length in bytes, and Table is the
// table's name without its database qualifier. Every finding anchors the
// whole statement, whichever line the column is on. Ported from Bytebase's
// V1 rule NAMING_COLUMN_AUTO_INCREMENT.
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
					for _, col := range commandColumns(cmd) {
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
