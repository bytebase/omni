package v1review

import (
	"regexp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingColumn judges the column names top-level CREATE TABLE and ALTER
// TABLE statements introduce against cfg: every column a CREATE TABLE
// defines, and in an ALTER TABLE the columns ADD COLUMN adds, the new
// column of CHANGE COLUMN, and the new name of RENAME COLUMN. The column
// of a MODIFY COLUMN is not judged, and neither are statements nested in a
// routine body or carried by EXPLAIN or PLAN REPLAYER. Name is the column
// name and Table the table name, both without a qualifier. Length is in
// bytes. A CREATE TABLE finding anchors the column definition; an ALTER
// TABLE finding anchors the whole statement. Ported from Bytebase's V1
// rule NAMING_COLUMN.
func NamingColumn(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
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
				if col != nil {
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
				case ast.ATRenameColumn:
					if cmd.NewName != "" {
						naming.check(n.Loc, n.Table.Name, cmd.NewName)
					}
				case ast.ATAddColumn:
					for _, col := range addedColumns(cmd) {
						if col != nil {
							naming.check(n.Loc, n.Table.Name, col.Name)
						}
					}
				case ast.ATChangeColumn:
					if cmd.Column != nil {
						naming.check(n.Loc, n.Table.Name, cmd.Column.Name)
					}
				}
			}
		}
	}
	return naming.findings, nil
}

// addedColumns returns the columns an ADD COLUMN command adds: the
// parenthesized list when it has one, otherwise its single column.
func addedColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	if len(cmd.Columns) > 0 {
		return cmd.Columns
	}
	if cmd.Column != nil {
		return []*ast.ColumnDef{cmd.Column}
	}
	return nil
}

// columnNaming applies a column rule's compiled config.
type columnNaming struct {
	pattern  *regexp.Regexp // nil requires nothing
	max      int            // 0 means no limit
	findings []v1review.Finding
}

func (c *columnNaming) check(loc ast.Loc, table, name string) {
	start, end := span(loc)
	if c.pattern != nil && !c.pattern.MatchString(name) {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Table: table, Start: start, End: end})
	}
	if c.max > 0 && len(name) > c.max {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Table: table, Start: start, End: end})
	}
}
