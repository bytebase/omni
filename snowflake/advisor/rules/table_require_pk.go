package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// TableRequirePKRule requires every persistent CREATE TABLE to define a
// PRIMARY KEY (inline or table-level).
//
// Rule ID:  "snowflake.table.require-pk"
// Severity: ERROR
//
// Exemptions:
//   - TEMPORARY, TRANSIENT, and VOLATILE tables (scratch tables)
//   - CREATE TABLE ... AS SELECT (the schema is derived from the query)
//   - CREATE TABLE ... LIKE source (clone of existing table)
//   - CREATE TABLE ... CLONE source (time-travel clone)
type TableRequirePKRule struct{}

const tableRequirePKID = "snowflake.table.require-pk"

// ID implements Rule.
func (TableRequirePKRule) ID() string { return tableRequirePKID }

// Severity implements Rule.
func (TableRequirePKRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (TableRequirePKRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	// Exempt temporary / transient / volatile tables.
	if stmt.Temporary || stmt.Transient || stmt.Volatile {
		return nil
	}
	// Exempt derived-table forms that don't define columns explicitly.
	if stmt.AsSelect != nil || stmt.Like != nil || stmt.Clone != nil {
		return nil
	}
	// Check for an inline PK on any column.
	for _, col := range stmt.Columns {
		if col.InlineConstraint != nil && col.InlineConstraint.Type == ast.ConstrPrimaryKey {
			return nil
		}
	}
	// Check for a table-level PK constraint.
	for _, c := range stmt.Constraints {
		if c.Type == ast.ConstrPrimaryKey {
			return nil
		}
	}
	return []*advisor.Finding{{
		RuleID:   tableRequirePKID,
		Severity: advisor.SeverityError,
		Loc:      stmt.Loc,
		Message:  fmt.Sprintf("table %q must define a PRIMARY KEY constraint", stmt.Name.String()),
	}}
}
