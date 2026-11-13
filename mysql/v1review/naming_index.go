package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIndex judges the names top-level statements give ordinary indexes
// against cfg: the table-level INDEX and KEY definitions of a CREATE TABLE,
// those an ALTER TABLE adds with ADD {INDEX | KEY}, the index a CREATE INDEX
// creates, and the new name of an ALTER TABLE ... RENAME {INDEX | KEY}.
// UNIQUE, FULLTEXT and SPATIAL definitions, an index inside ALTER TABLE's
// parenthesized ADD form, and statements nested in a routine body are not
// judged. The name is judged as parsed, an unnamed index as the empty name,
// its length in bytes.
//
// A rename is judged only when schema holds the old index on the table of
// that name, read without the database qualifier the statement wrote, and
// holds it as not unique: an index the same script created is not found, and
// a FULLTEXT or SPATIAL index is, though creating one is not judged.
//
// Each index gets its own pattern: cfg.Format with the table token replaced
// by the table name as parsed, without its database qualifier, and the
// column list token by the index's columns joined with "_". For a rename
// those are the expressions schema holds. For a definition they are its key
// parts that are plain columns, without prefix length or direction; when a
// table-level definition has none, and always for CREATE INDEX, a
// parenthesized key part that is one column counts too, without its table
// qualifier. When the pattern does not compile the index gets one
// PatternInvalid finding and is not judged further. A NameMismatch carries
// the expanded pattern.
//
// A CREATE TABLE finding anchors the index definition; every other finding
// anchors the whole statement. Findings are in statement order, an index's
// mismatch before its length. Ported from Bytebase's V1 rule
// NAMING_INDEX_IDX.
func NamingIndex(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	return namingKeys(stmts, cfg, schema, false)
}

// namingKeys judges ordinary indexes, or unique keys when unique is set.
func namingKeys(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema, unique bool) []v1review.Finding {
	var findings []v1review.Finding
	check := func(loc ast.Loc, table, name string, columns []string) {
		start, end := span(loc)
		pattern, err := cfg.Expand(map[string]string{
			v1review.ColumnListToken: strings.Join(columns, "_"),
			v1review.TableToken:      table,
		})
		if err != nil {
			findings = append(findings, v1review.Finding{Code: v1review.PatternInvalid, Name: name, Table: table, Detail: err.Error(), Start: start, End: end})
			return
		}
		if !pattern.MatchString(name) {
			findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern.String(), Start: start, End: end})
		}
		if cfg.MaxLength > 0 && len(name) > cfg.MaxLength {
			findings = append(findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Table: table, Start: start, End: end})
		}
	}
	constraintType := ast.ConstrIndex
	if unique {
		constraintType = ast.ConstrUnique
	}
	checkConstraint := func(loc ast.Loc, table string, constraint *ast.Constraint) {
		if constraint.Type != constraintType {
			return
		}
		columns := constraint.Columns
		if len(columns) == 0 {
			columns = keyPartColumns(constraint.IndexColumns)
		}
		check(loc, table, constraint.Name, columns)
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint != nil {
					checkConstraint(constraint.Loc, n.Table.Name, constraint)
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
						checkConstraint(n.Loc, table, cmd.Constraint)
					}
				case ast.ATRenameIndex:
					index := schema.TableIndex("", table, cmd.Name, false)
					if index == nil || index.GetUnique() != unique {
						continue
					}
					check(n.Loc, table, cmd.NewName, index.GetExpressions())
				default:
				}
			}
		case *ast.CreateIndexStmt:
			judged := n.Unique
			if !unique {
				judged = !n.Unique && !n.Fulltext && !n.Spatial
			}
			if !judged || n.Table == nil {
				continue
			}
			check(n.Loc, n.Table.Name, n.IndexName, keyPartColumns(n.Columns))
		}
	}
	return findings
}

// keyPartColumns returns the columns of the key parts that are one column,
// a parenthesized one included.
func keyPartColumns(parts []*ast.IndexColumn) []string {
	var columns []string
	for _, part := range parts {
		if ref, ok := part.Expr.(*ast.ColumnRef); ok {
			columns = append(columns, ref.Column)
		}
	}
	return columns
}
