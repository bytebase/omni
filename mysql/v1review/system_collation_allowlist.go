package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// SystemCollationAllowlist judges the collations top-level CREATE and ALTER
// DATABASE and CREATE and ALTER TABLE statements name against the allowed
// list, compared without regard to letter case. Statements nested in a
// routine body are not judged.
//
// Judged are every COLLATE option of a CREATE or ALTER DATABASE, every
// COLLATE table option of a CREATE or ALTER TABLE, and, for each column a
// CREATE TABLE defines or an ADD, MODIFY or CHANGE COLUMN action gives, the
// collation of its type and then each COLLATE among its attributes. A
// column with no type is not judged.
//
// Name is the collation lower-cased. A database option's finding anchors
// the statement, a table option's the option, a CREATE TABLE column's the
// column, and an ALTER TABLE column's its action. Ported from Bytebase's V1
// rule SYSTEM_COLLATION_ALLOWLIST.
func SystemCollationAllowlist(stmts []ast.Node, allowed []string) []v1review.Finding {
	c := &allowlist{code: v1review.CollationNotAllowed, allowed: make(map[string]bool, len(allowed))}
	for _, collation := range allowed {
		c.allowed[strings.ToLower(collation)] = true
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateDatabaseStmt:
			for _, opt := range n.Options {
				if strings.EqualFold(opt.Name, "COLLATE") {
					c.check(n.Loc, opt.Value)
				}
			}
		case *ast.AlterDatabaseStmt:
			for _, opt := range n.Options {
				if strings.EqualFold(opt.Name, "COLLATE") {
					c.check(n.Loc, opt.Value)
				}
			}
		case *ast.CreateTableStmt:
			for _, opt := range n.Options {
				if strings.EqualFold(opt.Name, "COLLATE") {
					c.check(opt.Loc, opt.Value)
				}
			}
			for _, col := range n.Columns {
				if col != nil {
					c.checkColumn(col.Loc, col)
				}
			}
		case *ast.AlterTableStmt:
			for _, cmd := range n.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
					for _, col := range commandColumns(cmd) {
						c.checkColumn(cmd.Loc, col)
					}
				case ast.ATTableOption:
					if cmd.Option != nil && strings.EqualFold(cmd.Option.Name, "COLLATE") {
						c.check(cmd.Option.Loc, cmd.Option.Value)
					}
				}
			}
		}
	}
	return c.findings
}

// checkColumn judges the collation of a column's type, then each COLLATE
// among its attributes, all anchored at loc.
func (c *allowlist) checkColumn(loc ast.Loc, col *ast.ColumnDef) {
	if col == nil || col.TypeName == nil {
		return
	}
	c.check(loc, col.TypeName.Collate)
	for _, attr := range col.Constraints {
		if attr.Type == ast.ColConstrCollate {
			c.check(loc, attr.Name)
		}
	}
}
