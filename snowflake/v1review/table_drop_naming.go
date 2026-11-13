package v1review

import (
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the table a top-level DROP TABLE drops against
// cfg.Pattern. DROP VIEW, DROP DYNAMIC TABLE, DROP EXTERNAL TABLE and every
// other DROP kind are not judged, nor are TRUNCATE TABLE and statements
// inside a Snowflake Scripting block, a routine or task body, or under
// EXPLAIN. The name is judged without its database or schema qualifier, as
// Snowflake resolves it: an unquoted name upper-cased, a quoted name as its
// content with a doubled quote read as one. An IDENTIFIER(...) name is
// judged as written and upper-cased whole, keyword, parentheses and string
// literal included, and an empty quoted name is judged as the empty string.
//
// cfg.MaxLength is not read: the rule has no length limit. Every finding
// anchors the last part of the dropped name, not the statement.
// NameMismatch is the only code. Ported from Bytebase's V1 rule
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
		drop, ok := stmt.(*ast.DropStmt)
		if !ok || drop.Kind != ast.DropTable || drop.Name == nil {
			continue
		}
		name := drop.Name.Name.Normalize()
		if pattern.MatchString(name) {
			continue
		}
		start, end := span(drop.Name.Name.Loc)
		findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
	}
	return findings, nil
}
