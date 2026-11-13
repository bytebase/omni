package v1review

import (
	"strings"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingForeignKey judges the names of the foreign keys top-level CREATE
// TABLE and ALTER TABLE statements define against cfg: the table-level
// FOREIGN KEY constraints of a CREATE TABLE, and those an ALTER TABLE adds.
// A REFERENCES written on a column is not judged, and neither are
// statements nested in a routine body or carried by EXPLAIN or PLAN
// REPLAYER.
//
// Name is the CONSTRAINT symbol, or the name written after FOREIGN KEY when
// there is no symbol, or empty when there is neither: an unnamed foreign
// key is judged as the empty name. Table is the table name without its
// qualifier. The pattern is cfg expanded with that table as the
// referencing table, the referenced table without its qualifier, and each
// side's column names joined with "_". A key whose pattern does not
// compile gets one PatternInvalid finding and is not judged further.
// Length is in bytes.
//
// A CREATE TABLE finding anchors the constraint. An ALTER TABLE finding
// has no range: V1 anchors it at the first character of the statement's
// text that is not white space, which a leading comment moves off the
// statement. Ported from Bytebase's V1 rule NAMING_INDEX_FK.
func NamingForeignKey(stmts []ast.Node, cfg v1review.TemplateConfig) ([]v1review.Finding, error) {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint == nil {
					continue
				}
				start, end := span(constraint.Loc)
				findings = judgeForeignKey(findings, cfg, n.Table.Name, constraint, start, end)
			}
		case *ast.AlterTableStmt:
			if n.Table == nil {
				continue
			}
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATAddConstraint {
					findings = judgeForeignKey(findings, cfg, n.Table.Name, cmd.Constraint, 0, 0)
				}
			}
		}
	}
	return findings, nil
}

// judgeForeignKey appends the findings for one constraint of table, which
// are none unless it is a foreign key.
func judgeForeignKey(findings []v1review.Finding, cfg v1review.TemplateConfig, table string, constraint *ast.Constraint, start, end int) []v1review.Finding {
	if constraint == nil || constraint.Type != ast.ConstrForeignKey {
		return findings
	}
	name := constraint.Name
	if name == "" {
		name = constraint.IndexName
	}
	referencedTable := ""
	if constraint.RefTable != nil {
		referencedTable = constraint.RefTable.Name
	}
	pattern, err := cfg.Expand(map[string]string{
		v1review.ReferencingTableToken:  table,
		v1review.ReferencingColumnToken: strings.Join(constraint.Columns, "_"),
		v1review.ReferencedTableToken:   referencedTable,
		v1review.ReferencedColumnToken:  strings.Join(constraint.RefColumns, "_"),
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
