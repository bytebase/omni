package rules

import (
	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// SelectNoSelectAllRule prohibits SELECT * (bare star) in query lists.
//
// Rule ID:  "snowflake.select.no-select-all"
// Severity: ERROR
//
// This is the production version of the example NoSelectStarRule shipped with
// the T2.6 advisor framework. The differences are:
//   - Rule ID uses the bytebase convention: "snowflake.select.no-select-all"
//   - Severity is ERROR (not WARNING) — SELECT * is a hard requirement violation
//
// Qualified stars (t.*) are not flagged by this rule.
type SelectNoSelectAllRule struct{}

const selectNoSelectAllID = "snowflake.select.no-select-all"

// ID implements Rule.
func (SelectNoSelectAllRule) ID() string { return selectNoSelectAllID }

// Severity implements Rule.
func (SelectNoSelectAllRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (SelectNoSelectAllRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	sel, ok := node.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	var findings []*advisor.Finding
	for _, target := range sel.Targets {
		if !target.Star {
			continue
		}
		starExpr, ok := target.Expr.(*ast.StarExpr)
		if !ok {
			// bare * with nil Expr
			findings = append(findings, &advisor.Finding{
				RuleID:   selectNoSelectAllID,
				Severity: advisor.SeverityError,
				Loc:      target.Loc,
				Message:  "avoid SELECT *: list the required columns explicitly",
			})
			continue
		}
		if starExpr.Qualifier != nil {
			// qualified star (t.*) — skip
			continue
		}
		findings = append(findings, &advisor.Finding{
			RuleID:   selectNoSelectAllID,
			Severity: advisor.SeverityError,
			Loc:      target.Loc,
			Message:  "avoid SELECT *: list the required columns explicitly",
		})
	}
	return findings
}
