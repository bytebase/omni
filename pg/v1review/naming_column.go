package v1review

import (
	"regexp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingColumn judges the column names top-level statements introduce
// against cfg: each column a CREATE TABLE defines, each column an ALTER ...
// ADD COLUMN adds, and the new name a RENAME COLUMN gives. The ALTER and
// RENAME forms are judged whatever relation they alter, so a view's or a
// foreign table's renamed column and an attribute ALTER TYPE adds are
// judged; CREATE FOREIGN TABLE, CREATE TABLE AS, ALTER TYPE ... RENAME
// ATTRIBUTE, and a statement in a CREATE SCHEMA or in a routine body are
// not. The name is judged as the parser stores it: an unquoted name is
// folded to lower case and every name is truncated to 63 bytes. Its length
// is in bytes. Table is the relation's name without its schema. Findings
// come in statement order, a CREATE TABLE's and an ALTER's in the order
// they write their columns. Every finding anchors the whole statement;
// Bytebase's V1 advisor reported the line of the name's first occurrence in
// the statement text instead, which no node range gives. Ported from
// Bytebase's V1 rule NAMING_COLUMN.
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
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.TableElts == nil {
				continue
			}
			for _, item := range n.TableElts.Items {
				if col, ok := item.(*ast.ColumnDef); ok {
					naming.check(n.Loc, relname(n.Relation), col.Colname)
				}
			}
		case *ast.AlterTableStmt:
			if n.Cmds == nil {
				continue
			}
			for _, item := range n.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok || ast.AlterTableType(cmd.Subtype) != ast.AT_AddColumn {
					continue
				}
				if col, ok := cmd.Def.(*ast.ColumnDef); ok {
					naming.check(n.Loc, relname(n.Relation), col.Colname)
				}
			}
		case *ast.RenameStmt:
			if n.RenameType == ast.OBJECT_COLUMN {
				naming.check(n.Loc, relname(n.Relation), n.Newname)
			}
		}
	}
	return naming.findings, nil
}

// columnNaming applies NamingColumn's compiled config.
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

// relname returns a relation's name without its schema, "" when the
// statement has none.
func relname(rv *ast.RangeVar) string {
	if rv == nil {
		return ""
	}
	return rv.Relname
}
