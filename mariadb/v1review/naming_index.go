package v1review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIndex judges the names of the plain indexes top-level statements
// create or rename against cfg:
//
//   - the table-level INDEX and KEY definitions of a CREATE TABLE;
//   - those an ALTER TABLE adds with ADD {INDEX | KEY};
//   - a CREATE INDEX that is not UNIQUE, FULLTEXT or SPATIAL, which leaves
//     CREATE VECTOR INDEX judged, though a VECTOR key of a CREATE TABLE or an
//     ALTER TABLE is not;
//   - the new name of an ALTER TABLE ... RENAME {INDEX | KEY}, when schema
//     holds the old name on that table and the index is not unique. A rename
//     of an index schema does not hold is not judged.
//
// UNIQUE, FULLTEXT, SPATIAL and PRIMARY KEY definitions, an index inside
// ALTER TABLE's parenthesized ADD form, and statements nested in a routine
// body are not judged. The name is the one parsed, and an index written
// without a name is judged with an empty one; its length is in bytes.
//
// Each index gets its own pattern: cfg.Format with the table token replaced
// by the statement's table name without its database qualifier, and the
// column list token by the index's columns joined with "_". The columns are
// the statement's key parts without prefix lengths or directions; a rename
// takes them from the schema's index, which it looks up in the unnamed
// schema by the statement's table name, whatever database qualifies it, with
// the index name compared without case. When the pattern does not compile
// the index gets one PatternInvalid finding and is not judged further. A
// NameMismatch carries the expanded pattern.
//
// A CREATE TABLE finding anchors the index definition; every other finding
// anchors the whole statement. Findings are in statement order, an index's
// mismatch before its length. Ported from Bytebase's V1 rule
// NAMING_INDEX_IDX.
func NamingIndex(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	naming := &indexNaming{
		cfg:        cfg,
		schema:     schema,
		constraint: ast.ConstrIndex,
		created:    func(n *ast.CreateIndexStmt) bool { return !n.Unique && !n.Fulltext && !n.Spatial },
		unique:     false,
	}
	naming.judge(stmts)
	return naming.findings
}

// indexNaming applies the config of NamingIndex or NamingUniqueKey, which
// walk the same statements and differ in the indexes they pick.
type indexNaming struct {
	cfg    v1review.TemplateConfig
	schema v1review.Schema
	// constraint is the table constraint type judged.
	constraint ast.ConstraintType
	// created reports whether a CREATE INDEX is judged.
	created func(*ast.CreateIndexStmt) bool
	// unique is what the schema must say of a renamed index for the rename
	// to be judged.
	unique   bool
	findings []v1review.Finding
}

func (c *indexNaming) judge(stmts []ast.Node) {
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint != nil {
					c.constraintIndex(constraint.Loc, n.Table.Name, constraint)
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
				case ast.ATAddConstraint, ast.ATAddIndex:
					if cmd.Constraint != nil {
						c.constraintIndex(n.Loc, table, cmd.Constraint)
					}
				case ast.ATRenameIndex:
					index := c.schema.TableIndex("", table, cmd.Name, false)
					if index == nil || index.GetUnique() != c.unique {
						continue
					}
					c.check(n.Loc, table, cmd.NewName, index.GetExpressions())
				default:
				}
			}
		case *ast.CreateIndexStmt:
			if !c.created(n) || n.Table == nil {
				continue
			}
			c.check(n.Loc, n.Table.Name, n.IndexName, indexColumnNames(n.Columns))
		}
	}
}

func (c *indexNaming) constraintIndex(loc ast.Loc, table string, constraint *ast.Constraint) {
	if constraint.Type != c.constraint {
		return
	}
	columns := constraint.Columns
	if len(columns) == 0 {
		columns = indexColumnNames(constraint.IndexColumns)
	}
	c.check(loc, table, constraint.Name, columns)
}

func (c *indexNaming) check(loc ast.Loc, table, name string, columns []string) {
	start, end := span(loc)
	pattern, err := c.cfg.Expand(map[string]string{
		v1review.TableToken:      table,
		v1review.ColumnListToken: strings.Join(columns, "_"),
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

// indexColumnNames returns the column each key part names; a key part that
// is not a column reference contributes nothing.
func indexColumnNames(parts []*ast.IndexColumn) []string {
	var names []string
	for _, part := range parts {
		if ref, ok := part.Expr.(*ast.ColumnRef); ok {
			names = append(names, ref.Column)
		}
	}
	return names
}
