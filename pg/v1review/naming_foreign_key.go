package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingForeignKey judges the foreign key constraint names top-level
// statements introduce against cfg: in a CREATE TABLE its table constraints
// and then its columns' constraints, and in an ALTER, command by command,
// an ADD CONSTRAINT's constraint and an ADD COLUMN's column constraints.
// The ALTER form is judged whatever relation it alters. CREATE FOREIGN
// TABLE and a statement in a CREATE SCHEMA or in a routine body are not
// judged.
//
// cfg.Format is expanded per foreign key with the referencing table, the
// referencing columns, the referenced table, and the referenced columns,
// each as the parser stores it: tables without their schema, column lists
// joined with "_", and a referenced column list the statement leaves out
// empty. A column constraint's referencing column is its own column. A
// format that does not compile once expanded gives one PatternInvalid
// finding and that foreign key is judged no further.
//
// The judged name is the constraint's name as the parser stores it: folded
// to lower case unless quoted, truncated to 63 bytes, and empty for a
// constraint the statement does not name, which is judged all the same. Its
// length is in bytes. Table is the relation's name without its schema.
// Every finding anchors the whole statement; Bytebase's V1 advisor
// reported the line of the name's first occurrence in the statement text
// instead, which no node range gives. Ported from Bytebase's V1 rule
// NAMING_INDEX_FK.
func NamingForeignKey(stmts []ast.Node, cfg v1review.TemplateConfig) ([]v1review.Finding, error) {
	naming := &foreignKeyNaming{cfg: cfg}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.TableElts == nil {
				continue
			}
			table := relname(n.Relation)
			for _, item := range n.TableElts.Items {
				if c, ok := item.(*ast.Constraint); ok {
					naming.check(n.Loc, table, c, strings.Join(names(c.FkAttrs), "_"))
				}
			}
			for _, item := range n.TableElts.Items {
				if col, ok := item.(*ast.ColumnDef); ok {
					naming.checkColumn(n.Loc, table, col)
				}
			}
		case *ast.AlterTableStmt:
			if n.Cmds == nil {
				continue
			}
			table := relname(n.Relation)
			for _, item := range n.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch ast.AlterTableType(cmd.Subtype) {
				case ast.AT_AddConstraint:
					if c, ok := cmd.Def.(*ast.Constraint); ok {
						naming.check(n.Loc, table, c, strings.Join(names(c.FkAttrs), "_"))
					}
				case ast.AT_AddColumn:
					if col, ok := cmd.Def.(*ast.ColumnDef); ok {
						naming.checkColumn(n.Loc, table, col)
					}
				}
			}
		}
	}
	return naming.findings, nil
}

// foreignKeyNaming applies NamingForeignKey's config.
type foreignKeyNaming struct {
	cfg      v1review.TemplateConfig
	findings []v1review.Finding
}

// checkColumn judges the foreign key constraints col carries.
func (c *foreignKeyNaming) checkColumn(loc ast.Loc, table string, col *ast.ColumnDef) {
	if col.Constraints == nil {
		return
	}
	for _, item := range col.Constraints.Items {
		if constraint, ok := item.(*ast.Constraint); ok {
			c.check(loc, table, constraint, col.Colname)
		}
	}
}

// check judges constraint when it is a foreign key on table whose
// referencing columns expand to referencing.
func (c *foreignKeyNaming) check(loc ast.Loc, table string, constraint *ast.Constraint, referencing string) {
	if constraint.Contype != ast.CONSTR_FOREIGN {
		return
	}
	start, end := span(loc)
	name := constraint.Conname
	pattern, err := c.cfg.Expand(map[string]string{
		v1review.ReferencingTableToken:  table,
		v1review.ReferencingColumnToken: referencing,
		v1review.ReferencedTableToken:   relname(constraint.Pktable),
		v1review.ReferencedColumnToken:  strings.Join(names(constraint.PkAttrs), "_"),
	})
	if err != nil {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.PatternInvalid, Name: name, Table: table, Detail: err.Error(), Start: start, End: end})
		return
	}
	if !pattern.MatchString(name) {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern.String(), Start: start, End: end})
	}
	if c.cfg.MaxLength > 0 && len(name) > c.cfg.MaxLength {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Table: table, Start: start, End: end})
	}
}

// names returns the names a column name list holds.
func names(list *ast.List) []string {
	if list == nil {
		return nil
	}
	var result []string
	for _, item := range list.Items {
		if s, ok := item.(*ast.String); ok {
			result = append(result, s.Str)
		}
	}
	return result
}
