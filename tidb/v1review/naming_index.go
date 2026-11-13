package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingIndex judges the names of the non-unique indexes top-level
// statements define or rename against cfg: the table-level INDEX and KEY
// constraints of a CREATE TABLE, those an ALTER TABLE adds, the index a
// CREATE INDEX without UNIQUE creates, FULLTEXT, SPATIAL and VECTOR ones
// included, and the new name of an ALTER TABLE ... RENAME INDEX. A
// FULLTEXT, SPATIAL or VECTOR constraint of a CREATE TABLE or an ALTER
// TABLE is not judged, and neither are statements nested in a routine body
// or carried by EXPLAIN or PLAN REPLAYER.
//
// Name is the name written after INDEX or KEY, or the CONSTRAINT symbol
// when there is none, or empty when there is neither: an unnamed index is
// judged as the empty name. Table is the statement's table without its
// qualifier. The pattern is cfg expanded with that table and the key's
// column names joined with "_"; a key part that is an expression gives no
// name. A key whose pattern does not compile gets one PatternInvalid
// finding and is not judged further. Length is in bytes.
//
// A renamed index is looked up by its old name in the unnamed schema of
// schema, across every table and materialized view and not only the
// statement's table, whatever database the statement names; index names
// compare without case. The rename is judged only when that index exists
// and is not unique, as its new name against the found index's expressions
// and the statement's table. Where two tables hold an index of that name,
// V1 took whichever its map yielded and this takes the first in the
// snapshot's order.
//
// A CREATE TABLE finding anchors the constraint. Every other finding has
// no range: V1 anchors it at the first character of the statement's text
// that is not white space, which a leading comment moves off the
// statement. Ported from Bytebase's V1 rule NAMING_INDEX_IDX.
func NamingIndex(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema) []v1review.Finding {
	return judgeKeys(stmts, cfg, schema, false)
}

// judgeKeys is the walk NamingIndex and NamingUniqueKey share; unique picks
// which of the two kinds of key is judged.
func judgeKeys(stmts []ast.Node, cfg v1review.TemplateConfig, schema v1review.Schema, unique bool) []v1review.Finding {
	kind := ast.ConstrIndex
	if unique {
		kind = ast.ConstrUnique
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint == nil || constraint.Type != kind {
					continue
				}
				start, end := span(constraint.Loc)
				findings = judgeKey(findings, cfg, n.Table.Name, constraint.Name, constraintColumns(constraint), start, end)
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
				case ast.ATAddConstraint, ast.ATAddIndex:
					if constraint := cmd.Constraint; constraint != nil && constraint.Type == kind {
						findings = judgeKey(findings, cfg, n.Table.Name, constraint.Name, constraintColumns(constraint), 0, 0)
					}
				case ast.ATRenameIndex:
					_, index := schema.SchemaIndex("", cmd.Name, false)
					if index == nil || index.GetUnique() != unique {
						continue
					}
					findings = judgeKey(findings, cfg, n.Table.Name, cmd.NewName, index.GetExpressions(), 0, 0)
				default:
				}
			}
		case *ast.CreateIndexStmt:
			if n.Unique != unique || n.Table == nil {
				continue
			}
			findings = judgeKey(findings, cfg, n.Table.Name, n.IndexName, keyPartColumns(n.Columns), 0, 0)
		}
	}
	return findings
}

// judgeKey appends the findings for one key of table.
func judgeKey(findings []v1review.Finding, cfg v1review.TemplateConfig, table, name string, columns []string, start, end int) []v1review.Finding {
	pattern, err := cfg.Expand(map[string]string{
		v1review.TableToken:      table,
		v1review.ColumnListToken: strings.Join(columns, "_"),
	})
	if err != nil {
		return append(findings, v1review.Finding{Code: v1review.PatternInvalid, Name: name, Table: table, Detail: err.Error(), Start: start, End: end})
	}
	if !pattern.MatchString(name) {
		findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Table: table, Pattern: pattern.String(), Start: start, End: end})
	}
	if cfg.MaxLength > 0 && len(name) > cfg.MaxLength {
		findings = append(findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Table: table, Pattern: pattern.String(), Start: start, End: end})
	}
	return findings
}

// constraintColumns returns a key constraint's column names: its Columns,
// or the names of its key parts when it has none.
func constraintColumns(constraint *ast.Constraint) []string {
	if len(constraint.Columns) > 0 {
		return constraint.Columns
	}
	return keyPartColumns(constraint.IndexColumns)
}

// keyPartColumns returns the column names of key parts. A part that is an
// expression gives none.
func keyPartColumns(parts []*ast.IndexColumn) []string {
	var names []string
	for _, part := range parts {
		if part == nil {
			continue
		}
		if ref, ok := part.Expr.(*ast.ColumnRef); ok {
			names = append(names, ref.Column)
		}
	}
	return names
}
