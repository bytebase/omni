package rules

import (
	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// WhereRequireSelectRule requires every SELECT that reads from a table (has a
// FROM clause) to also have a WHERE clause.
//
// Rule ID:  "snowflake.query.where-require-select"
// Severity: WARNING
//
// Table-less selects (SELECT 1, SELECT CURRENT_TIMESTAMP()) are exempt because
// they do not read from any table and therefore cannot produce runaway full scans.
type WhereRequireSelectRule struct{}

const whereRequireSelectID = "snowflake.query.where-require-select"

// ID implements Rule.
func (WhereRequireSelectRule) ID() string { return whereRequireSelectID }

// Severity implements Rule.
func (WhereRequireSelectRule) Severity() advisor.Severity { return advisor.SeverityWarning }

// Check implements Rule.
func (WhereRequireSelectRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	sel, ok := node.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	// No FROM clause — exempt (table-less select like SELECT 1).
	if len(sel.From) == 0 {
		return nil
	}
	if sel.Where != nil {
		return nil
	}
	return []*advisor.Finding{{
		RuleID:   whereRequireSelectID,
		Severity: advisor.SeverityWarning,
		Loc:      sel.Loc,
		Message:  "SELECT without a WHERE clause may produce a full table scan",
	}}
}
