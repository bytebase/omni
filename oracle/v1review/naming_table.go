// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for Oracle. The contract, and the split between these
// functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges the table names top-level statements introduce
// against cfg: the table a CREATE TABLE creates, CTAS included, and the new
// name an ALTER TABLE ... RENAME TO gives. A RENAME statement is not
// judged, nor is a CREATE TABLE nested in CREATE SCHEMA. The name is judged
// as the parser stores it — an unquoted name upper-cased, a quoted one as
// written — without its schema qualifier, its length in bytes. A CREATE
// TABLE finding anchors the statement; an ALTER TABLE finding anchors its
// RENAME clause. A statement may be passed bare or wrapped in an
// *ast.RawStmt. Ported from Bytebase's V1 rule NAMING_TABLE.
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
		case *ast.CreateTableStmt:
			if n.Name != nil && n.Name.Name != "" {
				naming.check(n.Loc, n.Name.Name)
			}
		case *ast.AlterTableStmt:
			if n.Actions == nil {
				continue
			}
			for _, item := range n.Actions.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if ok && cmd.Action == ast.AT_RENAME && cmd.NewName != "" {
					naming.check(cmd.Loc, cmd.NewName)
				}
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

// span returns a Loc as a finding's byte range, 0,0 when the AST kept no
// start. V1 anchored by the start alone, so a missing end collapses the
// range onto the start.
func span(loc ast.Loc) (int, int) {
	if loc.Start < 0 {
		return 0, 0
	}
	return loc.Start, max(loc.End, loc.Start)
}
