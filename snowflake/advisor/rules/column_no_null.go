// Package rules provides the production lint rules for the Snowflake advisor.
// Each file in this package implements one Rule from the T2.7 specification.
package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// ColumnNoNullRule requires every non-virtual column in CREATE TABLE to have
// an explicit NOT NULL constraint.
//
// Rule ID:  "snowflake.column.no-null"
// Severity: ERROR
//
// Virtual columns (ColumnDef.VirtualExpr != nil) are exempt because Snowflake
// does not permit NOT NULL on computed columns.
type ColumnNoNullRule struct{}

const columnNoNullID = "snowflake.column.no-null"

// ID implements Rule.
func (ColumnNoNullRule) ID() string { return columnNoNullID }

// Severity implements Rule.
func (ColumnNoNullRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (ColumnNoNullRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	var findings []*advisor.Finding
	for _, col := range stmt.Columns {
		// Virtual columns cannot have NOT NULL — skip them.
		if col.VirtualExpr != nil {
			continue
		}
		if !col.NotNull {
			findings = append(findings, &advisor.Finding{
				RuleID:   columnNoNullID,
				Severity: advisor.SeverityError,
				Loc:      col.Loc,
				Message:  fmt.Sprintf("column %q must have a NOT NULL constraint", col.Name.Name),
			})
		}
	}
	return findings
}
