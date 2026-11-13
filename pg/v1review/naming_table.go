// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for PostgreSQL. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges the table names top-level statements introduce
// against cfg: the table a CREATE TABLE creates, including a partition or
// typed table, and the new name an ALTER TABLE ... RENAME TO gives. CREATE
// TABLE AS, SELECT INTO, and a statement under EXPLAIN, in a CREATE SCHEMA,
// or in a routine body are not judged, nor are views, materialized views,
// and foreign tables. The name is judged as the parser stores it, without
// its schema qualifier: an unquoted name is folded to lower case and every
// name is truncated to 63 bytes. Its length is in bytes. Every finding
// anchors the whole statement. Ported from Bytebase's V1 rule NAMING_TABLE.
func NamingTable(stmts []ast.Node, cfg v1review.NamingConfig) ([]v1review.Finding, error) {
	pattern, err := cfg.Compile()
	if err != nil {
		return nil, err
	}
	if pattern == nil && cfg.MaxLength <= 0 {
		return nil, nil
	}
	naming := &tableNaming{pattern: pattern, max: cfg.MaxLength}
	for _, stmt := range stmts {
		if raw, ok := stmt.(*ast.RawStmt); ok {
			stmt = raw.Stmt
		}
		switch n := stmt.(type) {
		case *ast.CreateStmt:
			if n.Relation != nil {
				naming.check(n.Loc, n.Relation.Relname)
			}
		case *ast.RenameStmt:
			if n.RenameType == ast.OBJECT_TABLE {
				naming.check(n.Loc, n.Newname)
			}
		}
	}
	return naming.findings, nil
}

// tableNaming applies NamingTable's compiled config.
type tableNaming struct {
	pattern  *regexp.Regexp // nil requires nothing
	max      int            // 0 means no limit
	findings []v1review.Finding
}

func (c *tableNaming) check(loc ast.Loc, name string) {
	start, end := span(loc)
	if c.pattern != nil && !c.pattern.MatchString(name) {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameMismatch, Name: name, Start: start, End: end})
	}
	if c.max > 0 && len(name) > c.max {
		c.findings = append(c.findings, v1review.Finding{Code: v1review.NameTooLong, Name: name, Start: start, End: end})
	}
}

// span returns a Loc as a finding's byte range, 0,0 when the AST kept none.
func span(loc ast.Loc) (int, int) {
	if loc.Start < 0 || loc.End <= loc.Start {
		return 0, 0
	}
	return loc.Start, loc.End
}
