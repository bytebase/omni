package v1review

import (
	"strings"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingForeignKey judges the names of the foreign keys top-level
// statements add against cfg: the table-level FOREIGN KEY constraints of a
// CREATE TABLE, and those an ALTER TABLE adds. A column's own REFERENCES
// clause is not judged. Statements nested in a routine body are not judged.
//
// The name is the CONSTRAINT symbol, else the name written after FOREIGN
// KEY, else empty: a foreign key with no name is judged as the empty name.
// Its length is counted in bytes. cfg.Format is expanded per foreign key
// with the referencing and referenced table names, without their database
// qualifiers, and each side's column names joined with "_"; a referenced
// table or column list that is not written expands to the empty string.
//
// A foreign key whose expanded pattern does not compile yields one
// PatternInvalid finding and is judged no further, not for its length
// either. Otherwise NameMismatch carries the expanded pattern in Pattern.
// Table is the referencing table's name. A CREATE TABLE's finding anchors
// the constraint; an ALTER TABLE's anchors the whole statement. Ported from
// Bytebase's V1 rule NAMING_INDEX_FK.
func NamingForeignKey(stmts []ast.Node, cfg v1review.TemplateConfig) ([]v1review.Finding, error) {
	var findings []v1review.Finding
	check := func(loc ast.Loc, table string, constraint *ast.Constraint) {
		if constraint.Type != ast.ConstrForeignKey {
			return
		}
		name := constraint.Name
		if name == "" {
			name = constraint.IndexName
		}
		referencedTable := ""
		if constraint.RefTable != nil {
			referencedTable = constraint.RefTable.Name
		}
		start, end := span(loc)
		pattern, err := cfg.Expand(map[string]string{
			v1review.ReferencingTableToken:  table,
			v1review.ReferencingColumnToken: strings.Join(constraint.Columns, "_"),
			v1review.ReferencedTableToken:   referencedTable,
			v1review.ReferencedColumnToken:  strings.Join(constraint.RefColumns, "_"),
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
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.Table == nil {
				continue
			}
			for _, constraint := range n.Constraints {
				if constraint != nil {
					check(constraint.Loc, n.Table.Name, constraint)
				}
			}
		case *ast.AlterTableStmt:
			table := ""
			if n.Table != nil {
				table = n.Table.Name
			}
			for _, cmd := range n.Commands {
				if cmd != nil && cmd.Type == ast.ATAddConstraint && cmd.Constraint != nil {
					check(n.Loc, table, cmd.Constraint)
				}
			}
		}
	}
	return findings, nil
}
