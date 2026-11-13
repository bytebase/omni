package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIndex judges the names of the plain indexes top-level statements
// create or rename against cfg: the INDEX and KEY definitions of a CREATE
// TABLE, those an ALTER TABLE adds, the new name an ALTER TABLE ... RENAME
// INDEX gives, and the index of a CREATE INDEX that is not UNIQUE, FULLTEXT
// or SPATIAL; a CREATE VECTOR INDEX is judged. UNIQUE, FULLTEXT, SPATIAL
// and VECTOR definitions inside a CREATE TABLE or an ALTER TABLE are not
// judged, nor is a column's own KEY attribute. Statements nested in a
// routine body are not judged.
//
// A renamed index is judged only when schema holds the old name on the
// statement's table, and holds it as not unique: an index the same script
// created is not in schema, so its rename is not judged. The table is looked
// up in the schema with no name, whatever database qualifies it in the
// statement, and the index name without regard to case.
//
// The name is judged as written, and an index with no name as the empty
// name. Its length is counted in bytes. cfg.Format is expanded per index
// with the table's name, without its database qualifier, and the index's
// column names joined with "_". A definition's columns are its key parts
// written as a bare column, a prefix length dropped; when it has none, the
// parenthesized key parts that hold only a column. A CREATE INDEX's columns
// are both kinds together. A renamed index's are the expressions schema
// records for it.
//
// An index whose expanded pattern does not compile yields one
// PatternInvalid finding and is judged no further, not for its length
// either. Otherwise NameMismatch carries the expanded pattern in Pattern.
// Table is the table's name. A CREATE TABLE's finding anchors the
// definition; every other finding anchors the whole statement. Ported from
// Bytebase's V1 rule NAMING_INDEX_IDX.
func NamingIndex(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	return namingIndexes(stmts, cfg, schema, false)
}

// namingIndexes is NamingIndex when unique is false and NamingUniqueKey
// when it is true.
func namingIndexes(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema, unique bool) []v1review.Finding {
	var findings []v1review.Finding
	judge := func(loc ast.Loc, table, name string, columns []string) {
		start, end := span(loc)
		pattern, err := cfg.Expand(map[string]string{
			v1review.TableToken:      table,
			v1review.ColumnListToken: strings.Join(columns, "_"),
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
	definitionType := ast.ConstrIndex
	if unique {
		definitionType = ast.ConstrUnique
	}
	definition := func(loc ast.Loc, table string, constraint *ast.Constraint) {
		if constraint.Type != definitionType {
			return
		}
		columns := constraint.Columns
		if len(columns) == 0 {
			columns = columnRefNames(constraint.IndexColumns)
		}
		judge(loc, table, constraint.Name, columns)
	}
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint != nil {
					definition(constraint.Loc, n.Table.Name, constraint)
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
						definition(n.Loc, table, cmd.Constraint)
					}
				case ast.ATRenameIndex:
					index := schema.TableIndex("", table, cmd.Name, false)
					if index == nil || index.GetUnique() != unique {
						continue
					}
					judge(n.Loc, table, cmd.NewName, index.GetExpressions())
				default:
				}
			}
		case *ast.CreateIndexStmt:
			if n.Unique != unique || (!unique && (n.Fulltext || n.Spatial)) || n.Table == nil {
				continue
			}
			judge(n.Loc, n.Table.Name, n.IndexName, columnRefNames(n.Columns))
		}
	}
	return findings
}

// columnRefNames returns the names of the key parts whose expression is a
// column and nothing else, a parenthesized one included.
func columnRefNames(columns []*ast.IndexColumn) []string {
	var names []string
	for _, column := range columns {
		if ref, ok := column.Expr.(*ast.ColumnRef); ok {
			names = append(names, ref.Column)
		}
	}
	return names
}
