package rules

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/advisor"
	"github.com/bytebase/omni/snowflake/ast"
)

// MigrationCompatibilityRule forbids destructive DDL operations that can
// cause data loss or break dependent objects:
//
//   - DROP TABLE
//   - ALTER TABLE ... DROP COLUMN
//
// Rule ID:  "snowflake.migration.compatibility"
// Severity: ERROR
//
// The intent is to protect against accidental destructive migrations in CI/CD
// pipelines. Teams that need to drop objects should use a separate manual
// review process.
type MigrationCompatibilityRule struct{}

const migrationCompatibilityID = "snowflake.migration.compatibility"

// ID implements Rule.
func (MigrationCompatibilityRule) ID() string { return migrationCompatibilityID }

// Severity implements Rule.
func (MigrationCompatibilityRule) Severity() advisor.Severity { return advisor.SeverityError }

// Check implements Rule.
func (MigrationCompatibilityRule) Check(_ *advisor.Context, node ast.Node) []*advisor.Finding {
	switch n := node.(type) {
	case *ast.DropStmt:
		return checkDropStmt(n)
	case *ast.AlterTableStmt:
		return checkAlterTableDrops(n)
	}
	return nil
}

func checkDropStmt(stmt *ast.DropStmt) []*advisor.Finding {
	if stmt.Kind != ast.DropTable {
		return nil
	}
	return []*advisor.Finding{{
		RuleID:   migrationCompatibilityID,
		Severity: advisor.SeverityError,
		Loc:      stmt.Loc,
		Message: fmt.Sprintf(
			"DROP TABLE %q is a destructive operation and is not allowed in migration scripts",
			stmt.Name.String(),
		),
	}}
}

func checkAlterTableDrops(stmt *ast.AlterTableStmt) []*advisor.Finding {
	var findings []*advisor.Finding
	for _, action := range stmt.Actions {
		if action.Kind == ast.AlterTableDropColumn {
			names := make([]string, len(action.DropColumnNames))
			for i, n := range action.DropColumnNames {
				names[i] = n.Name
			}
			findings = append(findings, &advisor.Finding{
				RuleID:   migrationCompatibilityID,
				Severity: advisor.SeverityError,
				Loc:      action.Loc,
				Message: fmt.Sprintf(
					"ALTER TABLE %q DROP COLUMN is a destructive operation and is not allowed in migration scripts",
					stmt.Name.String(),
				),
			})
		}
	}
	return findings
}
