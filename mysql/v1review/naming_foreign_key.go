package v1review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingForeignKey judges the names of the foreign keys top-level statements
// add against cfg: the table-level FOREIGN KEY constraints of a CREATE
// TABLE, and those an ALTER TABLE adds with ADD [CONSTRAINT] FOREIGN KEY.
// A column's own REFERENCES clause, a foreign key inside ALTER TABLE's
// parenthesized ADD form, and statements nested in a routine body are not
// judged. The name is the CONSTRAINT symbol, else the name written after
// FOREIGN KEY, else empty, and an empty name is judged like any other; its
// length is in bytes.
//
// Each foreign key gets its own pattern: cfg.Format with the referencing
// and referenced table and column tokens replaced by the names as parsed,
// tables without their database qualifier and a column list joined with
// "_". When that pattern does not compile the foreign key gets one
// PatternInvalid finding and is not judged further. A NameMismatch carries
// the expanded pattern.
//
// A CREATE TABLE finding anchors the constraint; an ALTER TABLE finding
// anchors the whole statement. Findings are in statement order, a foreign
// key's mismatch before its length. Ported from Bytebase's V1 rule
// NAMING_INDEX_FK.
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
