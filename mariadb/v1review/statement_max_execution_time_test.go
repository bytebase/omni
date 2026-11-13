package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementMaxExecutionTime(t *testing.T) {
	unset := func(stmt int) []stmtFinding {
		return []stmtFinding{{Code: v1review.MaxExecutionTimeUnset, Name: "max_statement_time", Statement: stmt}}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "never set", sql: "SELECT 1; UPDATE t SET a = 1;", want: unset(0)},
		{name: "set first", sql: "SET max_statement_time = 1; SELECT 1;"},
		{name: "numeric string", sql: "SET SESSION MAX_STATEMENT_TIME = '10'; SELECT 1;"},
		{name: "first non-set reported", sql: "SET NAMES utf8mb4; SET max_statement_time = 'x'; SELECT 1; SELECT 2;", want: unset(2)},
		{name: "set after a statement", sql: "SELECT 1; SET max_statement_time = 1; SELECT 2;", want: unset(0)},
		{name: "decimal not counted", sql: "SET max_statement_time = 1.5; SELECT 1;", want: unset(1)},
		{name: "mysql spelling not counted", sql: "SET max_execution_time = 1; SELECT 1;", want: unset(1)},
		{name: "only sets", sql: "SET max_statement_time = 'x';"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementMaxExecutionTime)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
