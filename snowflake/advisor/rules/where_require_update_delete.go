package rules

import (
	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// WhereRequireUpdateDeleteRule requires every UPDATE and DELETE statement to
// include a WHERE clause.
//
// Rule ID:  "snowflake.query.where-require-update-delete"
// Severity: ERROR
//
// An UPDATE or DELETE without WHERE modifies or removes all rows in the table,
// which is almost never the intended behaviour in production migrations.
type WhereRequireUpdateDeleteRule struct{}

const whereRequireUpdateDeleteID = "snowflake.query.where-require-update-delete"

// ID implements Rule.
func (WhereRequireUpdateDeleteRule) ID() string { return whereRequireUpdateDeleteID }

// Severity implements Rule.
func (WhereRequireUpdateDeleteRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (WhereRequireUpdateDeleteRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		if n.Where == nil {
			return []*advisor.Finding{{
				RuleID:   whereRequireUpdateDeleteID,
				Severity: advisor.SeverityError,
				Loc:      n.Loc,
				Message:  "UPDATE without a WHERE clause will modify all rows",
			}}
		}
	case *ast.DeleteStmt:
		if n.Where == nil {
			return []*advisor.Finding{{
				RuleID:   whereRequireUpdateDeleteID,
				Severity: advisor.SeverityError,
				Loc:      n.Loc,
				Message:  "DELETE without a WHERE clause will remove all rows",
			}}
		}
	}
	return nil
}
