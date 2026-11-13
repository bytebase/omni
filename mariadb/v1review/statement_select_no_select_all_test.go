package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementSelectNoSelectAll(t *testing.T) {
	star := stmtFinding{Code: v1review.SelectAll, Text: "*"}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "select star", sql: "SELECT * FROM t;", want: []stmtFinding{star}},
		{name: "star among columns", sql: "SELECT a, b FROM t;\nSELECT *, a FROM t;", want: []stmtFinding{star}},
		{name: "subquery and set operation", sql: "SELECT a FROM (SELECT * FROM u) x UNION SELECT * FROM v;", want: []stmtFinding{star, star}},
		{name: "insert select", sql: "INSERT INTO t SELECT * FROM u;", want: []stmtFinding{star}},
		{name: "routine body", sql: "CREATE PROCEDURE p() BEGIN SELECT * FROM t; END;", want: []stmtFinding{star}},
		{name: "explain", sql: "EXPLAIN SELECT * FROM t;", want: []stmtFinding{star}},
		{name: "qualified star not judged", sql: "SELECT t.* FROM t;"},
		{name: "count star not judged", sql: "SELECT COUNT(*) FROM t;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementSelectNoSelectAll)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
