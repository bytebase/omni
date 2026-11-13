// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for Snowflake. The contract, and the split between
// these functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"unicode/utf8"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges every table name the statements introduce against
// cfg: the table a CREATE TABLE creates, in every form the parser gives
// the one node (OR REPLACE, OR ALTER, TEMPORARY, TRANSIENT, VOLATILE,
// HYBRID, AS SELECT, LIKE, CLONE), and the new name an ALTER TABLE ...
// RENAME TO gives, at top level or inside a Snowflake Scripting block, a
// task's body, or an alert's action. A statement under EXPLAIN is
// compiled, not run, and is skipped. The name is judged as the parser
// stores it: as written, quotes left out, an unquoted name not folded,
// and without its database or schema qualifier. SWAP WITH introduces no
// name, and a dynamic table is not a table. NameMismatch and NameTooLong
// are the codes. Ported from Bytebase's V1 rule NAMING_TABLE.
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
				// EXPLAIN compiles its statement without running it.
				return false
			case *ast.CreateTableStmt:
				if v.Name != nil {
					naming.check(v.Name.Loc, v.Name.Name.Name)
				}
			case *ast.AlterTableStmt:
				for _, a := range v.Actions {
					if a != nil && a.Kind == ast.AlterTableRename && a.NewName != nil {
						naming.check(a.NewName.Loc, a.NewName.Name.Name)
					}
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
