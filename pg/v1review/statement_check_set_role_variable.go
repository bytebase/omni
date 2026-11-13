package v1review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementCheckSetRoleVariable judges the statements together and reports
// one SetRoleMissing finding when none of the SET statements before the
// first statement that is not a SET sets the role, including when there
// are no statements. A SET of search_path and a RESET ALL count as setting
// the role, as they did in Bytebase's V1 advisor, and so do SET LOCAL and
// RESET ROLE. The finding anchors nothing: the V1 advisor gave it no
// position. Ported from Bytebase's V1 rule
// STATEMENT_CHECK_SET_ROLE_VARIABLE.
func StatementCheckSetRoleVariable(stmts []ast.Node) []v1review.Finding {
	for _, stmt := range stmts {
		set, ok := unwrapRaw(stmt).(*ast.VariableSetStmt)
		if !ok {
			break
		}
		if setsRoleOrSearchPath(set) {
			return nil
		}
	}
	return []v1review.Finding{{Code: v1review.SetRoleMissing}}
}

// setsRoleOrSearchPath reports whether a SET or RESET names the role or
// search_path, in any letter case, or is a RESET ALL.
func setsRoleOrSearchPath(stmt *ast.VariableSetStmt) bool {
	name := strings.ToLower(stmt.Name)
	return name == "role" || name == "search_path" || stmt.Kind == ast.VAR_RESET_ALL
}
