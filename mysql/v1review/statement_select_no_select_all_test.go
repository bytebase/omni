package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementSelectNoSelectAll(t *testing.T) {
	star := func(text string) stmtFinding { return stmtFinding{Code: v1review.SelectAll, Text: text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "select star", sql: "SELECT * FROM t;", want: []stmtFinding{star("*")}},
		{name: "union sides", sql: "SELECT * FROM t UNION SELECT a FROM s UNION SELECT * FROM u;", want: []stmtFinding{star("*"), star("*")}},
		{name: "subquery and insert select", sql: "INSERT INTO t SELECT a FROM (SELECT * FROM s) x;", want: []stmtFinding{star("*")}},
		{name: "count star not judged", sql: "SELECT COUNT(*) FROM t;"},
		{name: "qualified star not judged", sql: "SELECT t.* FROM t;"},
		{name: "columns", sql: "SELECT a, b FROM t;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementSelectNoSelectAll)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
