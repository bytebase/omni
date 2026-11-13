// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for Redshift. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDropNaming judges the tables top-level DROP TABLE statements drop
// against cfg.Pattern: every table the statement lists, in the order
// written. DROP VIEW, DROP MATERIALIZED VIEW, DROP FOREIGN TABLE, TRUNCATE,
// and a DROP TABLE in a routine body are not judged. The name is the last
// part of what was written, however many parts qualify it, as the parser
// stores it: an unquoted name is folded to lower case and truncated to 63
// bytes, a quoted name is kept whole with a doubled quote read as one. An
// empty name is skipped.
//
// cfg.MaxLength is not read: the rule has no length limit. Every finding
// anchors the whole statement. NameMismatch is the only code. Ported from
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
		drop, ok := stmt.(*ast.DropStmt)
		if !ok || drop.RemoveType != int(ast.OBJECT_TABLE) || drop.Objects == nil {
			continue
		}
		start, end := span(drop.Loc)
		for _, item := range drop.Objects.Items {
			name := droppedName(item)
			if name == "" || pattern.MatchString(name) {
				continue
			}
			findings = append(findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
		}
	}
	return findings, nil
}

// droppedName returns the last name part of one DROP object, "" when the
// object holds none.
func droppedName(item ast.Node) string {
	switch n := item.(type) {
	case *ast.List:
		name := ""
		for _, part := range n.Items {
			if s, ok := part.(*ast.String); ok {
				name = s.Str
			}
		}
		return name
	case *ast.RangeVar:
		return n.Relname
	default:
		return ""
	}
}

// span returns a Loc as a finding's byte range, 0,0 when the AST kept none.
func span(loc ast.Loc) (int, int) {
	if loc.Start < 0 || loc.End <= loc.Start {
		return 0, 0
	}
	return loc.Start, loc.End
}
