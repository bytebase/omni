package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingAutoIncrementColumn judges the names of AUTO_INCREMENT columns in
// top-level CREATE TABLE and ALTER TABLE statements against cfg: the
// columns a CREATE TABLE defines, and in an ALTER TABLE the columns ADD
// COLUMN adds and the new column of CHANGE COLUMN and MODIFY COLUMN. A
// column without AUTO_INCREMENT is not judged, and neither are statements
// nested in a routine body or carried by EXPLAIN or PLAN REPLAYER. Name is
// the column name and Table the table name, both without a qualifier.
// Length is in bytes. A CREATE TABLE finding anchors the column
// definition; an ALTER TABLE finding anchors the whole statement. Ported
// from Bytebase's V1 rule NAMING_COLUMN_AUTO_INCREMENT.
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
					naming.check(col.Loc, n.Table.Name, col.Name)
				}
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn:
					for _, col := range addedColumns(cmd) {
						if col != nil && col.AutoIncrement {
							naming.check(n.Loc, n.Table.Name, col.Name)
						}
					}
				case ast.ATChangeColumn, ast.ATModifyColumn:
					if cmd.Column != nil && cmd.Column.AutoIncrement {
						naming.check(n.Loc, n.Table.Name, cmd.Column.Name)
					}
				}
			}
		}
	}
	return naming.findings, nil
}
