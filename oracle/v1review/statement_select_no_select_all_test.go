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
		{name: "bare star", sql: "SELECT * FROM t", want: []stmtFinding{star("*")}},
		{name: "qualified star not judged", sql: "SELECT t.* FROM t"},
		{name: "count star not judged", sql: "SELECT COUNT(*) FROM t"},
		{name: "subquery", sql: "SELECT a FROM (SELECT * FROM t)", want: []stmtFinding{star("*")}},
		{name: "insert select", sql: "INSERT INTO t SELECT * FROM u", want: []stmtFinding{star("*")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementSelectNoSelectAll)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
