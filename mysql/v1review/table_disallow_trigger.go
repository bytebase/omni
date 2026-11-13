package v1review

import (
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// TableDisallowTrigger judges every top-level CREATE TRIGGER; statements
// nested in a routine body are not judged. A trigger with a name is
// TableTriggerCreated, anchoring the whole statement and naming nothing.
// Ported from Bytebase's V1 rule TABLE_DISALLOW_TRIGGER.
func TableDisallowTrigger(stmts []ast.Node) []v1review.Finding {
	var findings []v1review.Finding
	for _, stmt := range stmts {
		if n, ok := stmt.(*ast.CreateTriggerStmt); ok && n.Name != "" {
			start, end := span(n.Loc)
			findings = append(findings, v1review.Finding{Code: v1review.TableTriggerCreated, Start: start, End: end})
		}
	}
	return findings
}
