package v1review

import (
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the name of every table a top-level DROP TABLE
// drops against cfg's pattern, so that only tables named for removal can
// be dropped. DROP VIEW is not judged, and neither are statements nested
// in a routine body or carried by EXPLAIN or PLAN REPLAYER. A name is
// judged without its database qualifier. cfg's MaxLength is not judged:
// the rule has no length. Every finding anchors the whole statement.
// Ported from Bytebase's V1 rule TABLE_DROP_NAMING_CONVENTION.
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
