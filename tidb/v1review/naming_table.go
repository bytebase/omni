// Package v1review holds the SQL Review V1 judgments ported from
// Bytebase's advisor for TiDB. The contract, and the split between these
// functions and Bytebase's V1 product surface, is
// github.com/bytebase/omni/v1review.
package v1review

import (
	"regexp"
	"unicode/utf8"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// NamingTable judges every table name the statements introduce against
// cfg: the table a CREATE TABLE creates, and the new name a RENAME TABLE
// or an ALTER TABLE ... RENAME gives, at top level or inside a routine
// body. The name is judged as written, without its database qualifier,
// with NameMismatch and NameTooLong as the codes. Ported from Bytebase's
// V1 rule NAMING_TABLE.
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
			case *ast.ExplainStmt, *ast.TraceStmt, *ast.PlanReplayerStmt:
				return runsWrapped(v)
			case *ast.CreateTableStmt:
				if v.Table != nil {
					naming.check(v.Table.Loc, v.Table.Name)
				}
			case *ast.RenameTableStmt:
				for _, pair := range v.Pairs {
					if pair != nil && pair.New != nil {
						naming.check(pair.New.Loc, pair.New.Name)
					}
				}
			case *ast.AlterTableStmt:
				for _, cmd := range v.Commands {
					if cmd == nil || cmd.Type != ast.ATRenameTable {
						continue
					}
					if cmd.NewTable != nil {
						naming.check(cmd.NewTable.Loc, cmd.NewTable.Name)
					} else if cmd.NewName != "" {
						naming.check(cmd.Loc, cmd.NewName)
					}
				}
			}
			return true
		})
	}
	return naming.findings, nil
}

// runsWrapped reports whether one of TiDB's three wrapper statements runs
// the statement it carries, which is when the walk follows it. A plain
// EXPLAIN reports a plan, TRACE PLAN traces the optimizer, and a plain
// PLAN REPLAYER DUMP writes a dump, so DDL under those forms never runs
// on the server and stays unjudged; EXPLAIN ANALYZE, plain TRACE, and
// PLAN REPLAYER DUMP EXPLAIN ANALYZE perform what they carry.
func runsWrapped(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.ExplainStmt:
		return v.Analyze
	case *ast.TraceStmt:
		return !v.TracePlan
	case *ast.PlanReplayerStmt:
		return v.Analyze
	}
	return false
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
