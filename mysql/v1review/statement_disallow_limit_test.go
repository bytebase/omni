package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowLimit(t *testing.T) {
	limit := func(kind string) stmtFinding { return stmtFinding{Code: v1review.DMLLimit, Name: kind} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "update", sql: "UPDATE t SET a = 1 LIMIT 1;", want: []stmtFinding{limit("UPDATE")}},
		{name: "delete", sql: "DELETE FROM t LIMIT 1;", want: []stmtFinding{limit("DELETE")}},
		{name: "insert select", sql: "INSERT INTO t SELECT a FROM s LIMIT 1;", want: []stmtFinding{limit("INSERT")}},
		{
			name: "union and its sides",
			sql:  "INSERT INTO t (SELECT a FROM s LIMIT 1) UNION (SELECT a FROM u LIMIT 1) LIMIT 2;",
			want: []stmtFinding{limit("INSERT"), limit("INSERT"), limit("INSERT")},
		},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT a FROM (SELECT a FROM s LIMIT 1) x;"},
		{name: "select not judged", sql: "SELECT a FROM s LIMIT 1;"},
		{name: "no limit", sql: "DELETE FROM t WHERE a = 1;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementDisallowLimit)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
