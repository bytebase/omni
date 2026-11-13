// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for PostgreSQL. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"unicode/utf8"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges every table name the statements introduce against
// cfg: the table a CREATE TABLE, CREATE TABLE AS, or SELECT INTO creates,
// and the new name an ALTER TABLE ... RENAME TO gives, at top level or
// inside a SQL-standard routine body. The name is judged as the parser
// stores it, without its schema qualifier, so an unquoted name is already
// folded to lower case; views, materialized views, and foreign tables are
// not judged. NameMismatch and NameTooLong are the codes. Ported from
// Bytebase's V1 rule NAMING_TABLE.
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
		if stmt == nil {
			continue
		}
		ast.Inspect(stmt, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.ExplainStmt:
				return explainAnalyzes(v)
			case *ast.CreateStmt:
				if v.Relation != nil {
					naming.check(v.Relation.Loc, v.Relation.Relname)
				}
			case *ast.CreateTableAsStmt:
				if v.Objtype == ast.OBJECT_TABLE && v.Into != nil && v.Into.Rel != nil {
					naming.check(v.Into.Rel.Loc, v.Into.Rel.Relname)
				}
			case *ast.SelectStmt:
				if v.IntoClause != nil && v.IntoClause.Rel != nil {
					naming.check(v.IntoClause.Rel.Loc, v.IntoClause.Rel.Relname)
				}
			case *ast.RenameStmt:
				// The new name carries no Loc of its own; the finding
				// anchors the whole RenameStmt.
				if v.RenameType == ast.OBJECT_TABLE {
					naming.check(v.Loc, v.Newname)
				}
			}
			return true
		})
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
	if c.max > 0 && utf8.RuneCountInString(name) > c.max {
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
