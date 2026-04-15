package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// TableNoForeignKeyRule prohibits FOREIGN KEY constraints in CREATE TABLE.
//
// Rule ID:  "snowflake.table.no-foreign-key"
// Severity: ERROR
//
// Both inline (column-level) and table-level FOREIGN KEY constraints are
// flagged. Snowflake supports FK definitions for documentation purposes but
// does not enforce them. Many schemas prefer to avoid FK declarations to keep
// migrations simple.
type TableNoForeignKeyRule struct{}

const tableNoForeignKeyID = "snowflake.table.no-foreign-key"

// ID implements Rule.
func (TableNoForeignKeyRule) ID() string { return tableNoForeignKeyID }

// Severity implements Rule.
func (TableNoForeignKeyRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (TableNoForeignKeyRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	stmt, ok := node.(*ast.CreateTableStmt)
	if !ok {
		return nil
	}
	var findings []*advisor.Finding
	// Check inline column-level FKs.
	for _, col := range stmt.Columns {
		if col.InlineConstraint != nil && col.InlineConstraint.Type == ast.ConstrForeignKey {
			findings = append(findings, &advisor.Finding{
				RuleID:   tableNoForeignKeyID,
				Severity: advisor.SeverityError,
				Loc:      col.Loc,
				Message: fmt.Sprintf(
					"column %q on table %q must not define a FOREIGN KEY constraint",
					col.Name.Name, stmt.Name.String(),
				),
			})
		}
	}
	// Check table-level FK constraints.
	for _, c := range stmt.Constraints {
		if c.Type == ast.ConstrForeignKey {
			findings = append(findings, &advisor.Finding{
				RuleID:   tableNoForeignKeyID,
				Severity: advisor.SeverityError,
				Loc:      c.Loc,
				Message: fmt.Sprintf(
					"table %q must not define a FOREIGN KEY constraint",
					stmt.Name.String(),
				),
			})
		}
	}
	return findings
}
