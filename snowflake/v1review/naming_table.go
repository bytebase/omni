// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for Snowflake. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges the table names top-level statements introduce
// against cfg: the table a plain CREATE TABLE creates, and the new name an
// ALTER TABLE ... RENAME TO gives. CREATE TABLE ... AS SELECT, LIKE and
// CLONE are not judged, nor are statements inside a Snowflake Scripting
// block, a task's body, an alert's action, or under EXPLAIN. The name is
// judged without its database or schema qualifier, in its source form with
// one surrounding double quote trimmed from each end: an unquoted name is
// not folded, a quoted name keeps its inner quotes doubled ("A""B" is
// judged as A""B), and an IDENTIFIER(...) name is judged as written,
// keyword and parentheses included. Length counts bytes. Every finding
// anchors the statement node, whose range for ALTER TABLE starts at TABLE,
// not ALTER. NameMismatch and NameTooLong are the codes.
// Ported from Bytebase's V1 rule NAMING_TABLE.
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
		switch n := stmt.(type) {
		case *ast.CreateTableStmt:
			if n.AsSelect != nil || n.Like != nil || n.Clone != nil || n.Name == nil {
				continue
			}
			naming.check(n.Loc, n.Name.Name)
		case *ast.AlterTableStmt:
			for _, a := range n.Actions {
				if a == nil || a.Kind != ast.AlterTableRename || a.NewName == nil {
					continue
				}
				naming.check(n.Loc, a.NewName.Name)
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

func (c *tableNaming) check(loc ast.Loc, ident ast.Ident) {
	name := strings.TrimPrefix(strings.TrimSuffix(ident.String(), `"`), `"`)
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
