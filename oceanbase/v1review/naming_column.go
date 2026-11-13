package v1review

import (
	"regexp"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingColumn judges the column names top-level statements introduce
// against cfg: every column a CREATE TABLE defines, and in an ALTER TABLE
// the columns an ADD [COLUMN] defines, the new name of a RENAME COLUMN, and
// the new name of a CHANGE [COLUMN]. MODIFY is not judged, and neither are
// the columns a CREATE TABLE takes from LIKE or a query. Statements nested
// in a routine body are not judged.
//
// The name is judged as written, its length in bytes, and Table is the
// table's name without its database qualifier. A CREATE TABLE's finding
// anchors the column definition; an ALTER TABLE's anchors the whole
// statement, whichever line the column is on. Ported from Bytebase's V1
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
			table := ""
			if n.Table != nil {
				table = n.Table.Name
			}
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn:
					for _, col := range commandColumns(cmd) {
						naming.check(n.Loc, table, col.Name)
					}
				case ast.ATRenameColumn:
					if cmd.NewName != "" {
						naming.check(n.Loc, table, cmd.NewName)
					}
				case ast.ATChangeColumn:
					if cmd.Column != nil {
						naming.check(n.Loc, table, cmd.Column.Name)
					}
				default:
				}
			}
		}
	}
	return naming.findings, nil
}

// columnNaming applies a column naming rule's compiled config.
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

// commandColumns returns the column definitions an ADD, MODIFY or CHANGE
// command holds: its one column, or the list of a parenthesized ADD.
func commandColumns(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	switch cmd.Type {
	case ast.ATAddColumn:
		if cmd.Column != nil {
			return []*ast.ColumnDef{cmd.Column}
		}
		return cmd.Columns
	case ast.ATModifyColumn, ast.ATChangeColumn:
		if cmd.Column != nil {
			return []*ast.ColumnDef{cmd.Column}
		}
	default:
	}
	return nil
}
