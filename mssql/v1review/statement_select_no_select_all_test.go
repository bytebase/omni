package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementSelectNoSelectAll(t *testing.T) {
	star := func(text string) statementFinding { return statementFinding{Code: v1review.SelectAll, Text: text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []statementFinding
	}{
		{name: "bare star", sql: "SELECT * FROM t", want: []statementFinding{star("*")}},
		{name: "qualified star", sql: "SELECT a, t.* FROM t", want: []statementFinding{star("t.*")}},
		{name: "no star", sql: "SELECT a FROM t"},
		{name: "count star not judged", sql: "SELECT COUNT(*) FROM t"},
		{name: "subquery", sql: "SELECT a FROM (SELECT * FROM t) x", want: []statementFinding{star("*")}},
		{name: "exists subquery", sql: "SELECT a FROM t WHERE EXISTS (SELECT * FROM u)", want: []statementFinding{star("*")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementSelectNoSelectAll)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
