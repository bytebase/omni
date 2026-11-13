package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementMaxExecutionTime reports the first top-level statement that is not
// a SET, unless an earlier SET assigned max_statement_time (any case) an
// integer literal or a string literal that parses as a decimal integer. At
// most one finding; it carries no range and sets Statement to the reported
// statement's index (the V1 advisor printed no position).
// MaxExecutionTimeUnset is the only code; Name is "max_statement_time".
// Ported from Bytebase's V1 rule STATEMENT_MAX_EXECUTION_TIME.
func StatementMaxExecutionTime(stmts []ast.Node) []v1review.Finding {
	const variable = "max_statement_time"
	for i, stmt := range stmts {
		set, ok := stmt.(*ast.SetStmt)
		if !ok {
			return []v1review.Finding{{Code: v1review.MaxExecutionTimeUnset, Name: variable, Statement: i}}
		}
		for _, assign := range set.Assignments {
			if assign.Column == nil || assign.Value == nil || strings.ToLower(assign.Column.Column) != variable {
				continue
			}
			switch v := assign.Value.(type) {
			case *ast.IntLit:
				return nil
			case *ast.StringLit:
				if _, err := strconv.Atoi(v.Value); err == nil {
					return nil
				}
			}
		}
	}
	return nil
}
