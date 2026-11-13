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
		{name: "delete", sql: "DELETE FROM t WHERE a = 1 LIMIT 1;", want: []stmtFinding{limit("DELETE")}},
		{name: "insert select", sql: "INSERT INTO t SELECT a FROM u LIMIT 1;", want: []stmtFinding{limit("INSERT")}},
		{
			name: "outer limit before inner ones",
			sql:  "INSERT INTO t (SELECT a FROM u LIMIT 1) UNION (SELECT b FROM v LIMIT 2) LIMIT 3;",
			want: []stmtFinding{limit("INSERT"), limit("INSERT"), limit("INSERT")},
		},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT * FROM (SELECT a FROM v LIMIT 1) x;"},
		{name: "select not judged", sql: "SELECT a FROM t LIMIT 1;"},
		{name: "no limit", sql: "UPDATE t SET a = 1; DELETE FROM t; INSERT INTO t VALUES (1);"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN DELETE FROM t LIMIT 1; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementDisallowLimit)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
