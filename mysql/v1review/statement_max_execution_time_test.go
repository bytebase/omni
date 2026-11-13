package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/v1review"
)

func TestStatementMaxExecutionTime(t *testing.T) {
	unset := func(stmt int) []v1review.Finding {
		return []v1review.Finding{{Code: v1review.MaxExecutionTimeUnset, Name: "max_execution_time", Statement: stmt}}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []v1review.Finding
	}{
		{name: "never set", sql: "SELECT 1; SELECT 2;", want: unset(0)},
		{name: "set first", sql: "SET max_execution_time = 1000; SELECT 1;"},
		{name: "set with scope and case", sql: "SET SESSION MAX_EXECUTION_TIME = 1000; SELECT 1;"},
		{name: "set from integer string", sql: "SET max_execution_time = '1000'; SELECT 1;"},
		{name: "other string does not count", sql: "SET max_execution_time = 'x'; SELECT 1;", want: unset(1)},
		{name: "set too late", sql: "SET foo = 1; SELECT 1; SET max_execution_time = 1000;", want: unset(1)},
		{name: "only sets", sql: "SET foo = 1;"},
		{name: "set transaction is not a set", sql: "SET TRANSACTION ISOLATION LEVEL READ COMMITTED; SET max_execution_time = 1;", want: unset(0)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if diff := cmp.Diff(tt.want, StatementMaxExecutionTime(list.Items)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
