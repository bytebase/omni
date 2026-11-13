package v1review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementMaxExecutionTime reports the first top-level statement that is
// not a SET and runs before any SET assigns max_execution_time (any case,
// any scope) an integer literal or a string literal that parses as a
// decimal integer; any other value, a negative number included, does not
// count as set. Every statement that is not an ast.SetStmt counts, whatever
// it is (SET TRANSACTION among them). There is at most one
// finding: Name is the variable, "max_execution_time", and
// Statement the index of the statement; it anchors nothing (the V1 advice
// had no position). MaxExecutionTimeUnset is the only code. Ported from
// Bytebase's V1 rule STATEMENT_MAX_EXECUTION_TIME.
func StatementMaxExecutionTime(stmts []ast.Node) []v1review.Finding {
	const variable = "max_execution_time"
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
