package advisor

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// NoSelectStarRule warns when a SELECT list contains a bare * (SELECT *).
//
// This is a minimal example rule included with the framework to prove the
// wiring works end-to-end. It is NOT one of the T2.7 production rules, but
// it exercises the full pipeline: parse -> advisor.Check -> findings with Loc.
//
// Trigger:  a *ast.SelectStmt whose Targets slice contains an entry with
//
//	Star == true AND whose Expr is a *ast.StarExpr with Qualifier == nil
//	(bare * — no table qualifier).
//
// Not flagged: qualified star such as t.* (SelectTarget.Star true with Expr set
//
//	to a *ast.StarExpr carrying a non-nil Qualifier). Qualified star is
//	a different smell and is not flagged by this rule.
//
// Rule ID:  "snowflake.select.no-select-star"
// Severity: WARNING
//
// Implementation note: ast.SelectTarget is not an ast.Node and is therefore
// not visited directly by ast.Walk. The rule matches the parent SelectStmt
// and inspects its Targets slice in-line.
type NoSelectStarRule struct{}

const noSelectStarID = "snowflake.select.no-select-star"

// ID implements Rule.
func (NoSelectStarRule) ID() string { return noSelectStarID }

// Severity implements Rule.
func (NoSelectStarRule) Severity() Severity { return SeverityWarning }

// Check implements Rule.
func (NoSelectStarRule) Check(_ *Context, node ast.Node) []*Finding {
	sel, ok := node.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	var findings []*Finding
	for _, target := range sel.Targets {
		if !target.Star {
			continue
		}
		// target.Expr is always a *ast.StarExpr when Star is true.
		// A bare * has StarExpr.Qualifier == nil.
		// A qualified star (t.*) has StarExpr.Qualifier != nil.
		starExpr, ok := target.Expr.(*ast.StarExpr)
		if !ok {
			continue
		}
		if starExpr.Qualifier != nil {
			continue
		}
		findings = append(findings, &Finding{
			RuleID:   noSelectStarID,
			Severity: SeverityWarning,
			Loc:      target.Loc,
			Message:  "avoid SELECT *: list the required columns explicitly",
		})
	}
	return findings
}
