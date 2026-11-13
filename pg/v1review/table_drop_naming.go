package v1review

import (
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the names of the tables top-level DROP TABLE
// statements drop against cfg.Pattern, one finding per mismatching table in
// the order the statement lists them. Only a name written as table or
// schema.table is judged, without its schema; one written with three parts
// is skipped. DROP VIEW, DROP MATERIALIZED VIEW, DROP FOREIGN TABLE, and a
// statement in a routine body are not judged. The name is judged as the
// parser stores it: an unquoted name is folded to lower case and every name
// is truncated to 63 bytes. cfg.MaxLength is not consulted: the rule has no
// length limit. Every finding anchors the whole statement. Ported from
// Bytebase's V1 rule TABLE_DROP_NAMING_CONVENTION.
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
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		n, ok := stmt.(*ast.DropStmt)
		if !ok || n.RemoveType != int(ast.OBJECT_TABLE) || n.Objects == nil {
			continue
		}
		start, end := span(n.Loc)
		for _, item := range n.Objects.Items {
			list, ok := item.(*ast.List)
			if !ok {
				continue
			}
			parts := names(list)
			if len(parts) != 1 && len(parts) != 2 {
				continue
			}
			if name := parts[len(parts)-1]; !pattern.MatchString(name) {
				findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
			}
		}
	}
	return findings, nil
}
