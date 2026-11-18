package analysis

import (
	"testing"

	"github.com/bytebase/omni/mssql"
)

// TestStatementTypeOfExecute holds a procedure called by name without EXEC as
// a batch's first statement to the type of EXEC: it is the same statement,
// with ExecStmt.Implied saying the call has no EXEC in the text.
func TestStatementTypeOfExecute(t *testing.T) {
	for _, sql := range []string{
		"EXEC my_proc 1",
		"EXECUTE my_proc",
		"my_proc 1",
		"my_proc",
		"dbo.my_proc @p = 1",
	} {
		stmts, err := mssql.Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		if len(stmts) != 1 {
			t.Fatalf("Parse(%q): %d statements, want 1", sql, len(stmts))
		}
		if got := StatementTypeOf(stmts[0].AST); got != Execute {
			t.Errorf("StatementTypeOf(%q) = %v, want Execute", sql, got)
		}
	}
}
