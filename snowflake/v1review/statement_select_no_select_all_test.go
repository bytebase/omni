package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementSelectNoSelectAll(t *testing.T) {
	star := func(start int, text string) stmtFinding {
		return stmtFinding{Code: v1review.SelectAll, Start: start, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "bare star", sql: "SELECT * FROM t", want: []stmtFinding{star(7, "*")}},
		{name: "qualified star anchors the star", sql: "SELECT a, t.* FROM t", want: []stmtFinding{star(12, "*")}},
		{name: "no star", sql: "SELECT a, b FROM t"},
		{name: "count star not judged", sql: "SELECT COUNT(*) FROM t"},
		{name: "subquery", sql: "SELECT a FROM (SELECT * FROM t) x", want: []stmtFinding{star(22, "*")}},
		{name: "insert select", sql: "INSERT INTO t SELECT * FROM u", want: []stmtFinding{star(21, "*")}},
		{name: "cte body", sql: "WITH c AS (SELECT * FROM t) SELECT a FROM c", want: []stmtFinding{star(18, "*")}},
		{
			name: "anchor order, not walk order",
			sql:  "WITH c AS (SELECT * FROM t) SELECT * FROM (SELECT * FROM c) x",
			want: []stmtFinding{star(18, "*"), star(35, "*"), star(50, "*")},
		},
		{name: "star with exclude", sql: "SELECT * EXCLUDE (a) FROM t", want: []stmtFinding{star(7, "*")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementSelectNoSelectAll)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
