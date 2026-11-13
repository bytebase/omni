package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowOrderBy(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "update", sql: "UPDATE t SET a = 1 ORDER BY a;", want: []stmtFinding{{Code: v1review.DMLOrderBy, Name: "UPDATE", Text: "UPDATE t SET a = 1 ORDER BY a"}}},
		{name: "delete", sql: "SELECT 1;\nDELETE FROM t ORDER BY a LIMIT 1;", want: []stmtFinding{{Code: v1review.DMLOrderBy, Name: "DELETE", Text: "DELETE FROM t ORDER BY a LIMIT 1"}}},
		{name: "select not judged", sql: "SELECT a FROM t ORDER BY a;"},
		{name: "insert select not judged", sql: "INSERT INTO t SELECT a FROM s ORDER BY a;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementDisallowOrderBy)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
