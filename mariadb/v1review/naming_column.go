package v1review

import (
	"regexp"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingColumn judges the column names top-level statements introduce
// against cfg: every column a CREATE TABLE defines, every column an ALTER
// TABLE adds with ADD [COLUMN], in its single and its parenthesized form,
// the new name of CHANGE [COLUMN], and the new name of RENAME COLUMN.
// MODIFY [COLUMN], which keeps the name, and statements nested in a routine
// body are not judged. The name is judged as parsed, its length in bytes,
// and the finding carries the table's name without its database qualifier.
// An empty name is judged in a definition and skipped as RENAME COLUMN's
// new name.
//
// A CREATE TABLE finding anchors the column's definition; an ALTER TABLE
// finding anchors the whole statement. Findings are in statement order, a
// column's mismatch before its length. Ported from Bytebase's V1 rule
// NAMING_COLUMN.
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
				case ast.ATAddColumn, ast.ATChangeColumn:
					for _, col := range columnDefs(cmd) {
						naming.check(n.Loc, table, col.Name)
					}
				case ast.ATRenameColumn:
					if cmd.NewName != "" {
						naming.check(n.Loc, table, cmd.NewName)
					}
				default:
				}
			}
		}
	}
	return naming.findings, nil
}

// columnDefs returns the column definitions an ADD, MODIFY or CHANGE action
// carries: its one column, else the columns of ADD's parenthesized form.
func columnDefs(cmd *ast.AlterTableCmd) []*ast.ColumnDef {
	if cmd.Column != nil {
		return []*ast.ColumnDef{cmd.Column}
	}
	if cmd.Type != ast.ATAddColumn {
		return nil
	}
	var cols []*ast.ColumnDef
	for _, col := range cmd.Columns {
		if col != nil {
			cols = append(cols, col)
		}
	}
	return cols
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
