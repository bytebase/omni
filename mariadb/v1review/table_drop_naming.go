package v1review

import (
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the name of every table a top-level DROP TABLE
// drops against cfg.Pattern; cfg.MaxLength is not consulted. DROP TEMPORARY
// TABLE is judged like any other; DROP VIEW, DROP SEQUENCE and statements
// nested in a routine body are not. The name is judged as parsed, without
// its database qualifier. Every finding anchors the whole statement, and
// findings are in statement order. Ported from Bytebase's V1 rule
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
