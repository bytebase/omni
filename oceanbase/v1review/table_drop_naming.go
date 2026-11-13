package v1review

import (
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the names of the tables top-level DROP TABLE
// statements drop against cfg.Pattern, one finding per table that does not
// match, in the order the statement lists them. cfg.MaxLength is not
// judged: the rule has no length. DROP VIEW is not judged, and neither are
// statements nested in a routine body.
//
// The name is judged as written, without its database qualifier. Every
// finding anchors the whole statement. Ported from Bytebase's V1 rule
// TABLE_DROP_NAMING_CONVENTION.
func TableDropNaming(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
	pattern, err := cfg.Compile()
	if err != nil {
		return nil, err
	}
	if pattern == nil {
		return nil, nil
	}
	var findings []v1review.Finding
	for _, stmt := range stmts {
		n, ok := stmt.(*ast.DropTableStmt)
		if !ok {
			continue
		}
		start, end := span(n.Loc)
		for _, table := range n.Tables {
			if table != nil && !pattern.MatchString(table.Name) {
				findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: table.Name, Start: start, End: end})
			}
		}
	}
	return findings, nil
}
