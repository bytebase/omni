// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for Oracle. The contract, and the split between these
// functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"unicode/utf8"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges every table name the statements introduce against
// cfg: the table a CREATE TABLE creates, a CREATE SCHEMA's included, the
// new name an ALTER TABLE ... RENAME TO gives, and the new name of a
// RENAME. RENAME also renames a view, a sequence, or a private synonym,
// which the statement alone cannot tell apart, so every RENAME's new name
// is judged. The name is judged as the parser folds it — an unquoted name
// upper-cased, a quoted one as written — without its schema qualifier,
// with NameMismatch and NameTooLong as the codes. DDL inside PL/SQL runs
// only through EXECUTE IMMEDIATE, whose text is not read. Ported from
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
			case *ast.CreateTableStmt:
				if v.Name != nil {
					naming.check(v.Name.Loc, v.Name.Name)
				}
			case *ast.RenameStmt:
				if v.NewName != nil {
					naming.check(v.NewName.Loc, v.NewName.Name)
				}
			case *ast.AlterTableStmt:
				if v.Actions == nil {
					return true
				}
				for _, item := range v.Actions.Items {
					cmd, ok := item.(*ast.AlterTableCmd)
					if ok && cmd.Action == ast.AT_RENAME && cmd.NewName != "" {
						naming.check(cmd.Loc, cmd.NewName)
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
