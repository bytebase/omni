package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireSelect(t *testing.T) {
	missing := func(text string) stmtFinding {
		return stmtFinding{Code: v1review.WhereMissing, Name: "SELECT", Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "no where", sql: "SELECT a FROM t;", want: []stmtFinding{missing("SELECT a FROM t")}},
		{name: "no from not judged", sql: "SELECT 1;"},
		{name: "where", sql: "SELECT a FROM t WHERE a = 1;"},
		{
			name: "cte first, then query, then subquery",
			sql:  "WITH c AS (SELECT a FROM s) SELECT (SELECT b FROM u) FROM c WHERE a = 1;",
			want: []stmtFinding{missing("SELECT a FROM s"), missing("SELECT b FROM u")},
		},
		{name: "union sides", sql: "SELECT a FROM t UNION SELECT a FROM s WHERE a = 1;", want: []stmtFinding{missing("SELECT a FROM t")}},
		{name: "insert select not judged", sql: "INSERT INTO t SELECT a FROM s;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireSelect)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
