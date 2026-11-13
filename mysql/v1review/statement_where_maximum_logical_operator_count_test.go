package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereMaximumLogicalOperatorCount(t *testing.T) {
	exceeded := func(op string, count int, text string) stmtFinding {
		return stmtFinding{Code: v1review.LogicalOperatorsExceeded, Name: op, Count: count, Text: text}
	}
	rule := func(stmts []ast.Node) []v1review.Finding { return StatementWhereMaximumLogicalOperatorCount(stmts, 2) }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "or chain anchors the deepest or", sql: "SELECT * FROM t WHERE a = 1 OR b = 2 OR c = 3;", want: []stmtFinding{exceeded("OR", 3, "a = 1 OR b = 2")}},
		{name: "at the limit", sql: "DELETE FROM t WHERE a = 1 OR b = 2;"},
		{name: "right nested or counts two", sql: "UPDATE t SET a = 1 WHERE a = 1 OR (b = 2 OR c = 3);"},
		{name: "in list", sql: "SELECT * FROM t WHERE a IN (1, 2, 3);", want: []stmtFinding{exceeded("IN", 3, "a IN (1, 2, 3)")}},
		{name: "in list in subquery", sql: "DELETE FROM t WHERE EXISTS (SELECT a FROM s WHERE b IN (1, 2, 3));", want: []stmtFinding{exceeded("IN", 3, "b IN (1, 2, 3)")}},
		{name: "in nested in an in not judged", sql: "DELETE FROM t WHERE a IN (SELECT a FROM s WHERE b IN (1, 2, 3));"},
		{name: "or in subquery not judged", sql: "DELETE FROM t WHERE a IN (SELECT a FROM s WHERE b = 1 OR b = 2 OR b = 3);"},
		{
			name: "union sides and insert select",
			sql:  "INSERT INTO u SELECT a FROM t WHERE a IN (1, 2, 3) UNION SELECT a FROM s WHERE a = 1 OR a = 2 OR a = 3;",
			want: []stmtFinding{exceeded("IN", 3, "a IN (1, 2, 3)"), exceeded("OR", 3, "a = 1 OR a = 2")},
		},
		{name: "having not judged", sql: "SELECT a FROM t GROUP BY a HAVING a IN (1, 2, 3);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, rule)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStatementWhereMaximumLogicalOperatorCountNegative(t *testing.T) {
	rule := func(stmts []ast.Node) []v1review.Finding { return StatementWhereMaximumLogicalOperatorCount(stmts, -1) }
	got := runStatement(t, "SELECT a FROM t WHERE a IN (SELECT b FROM s);", rule)
	want := []stmtFinding{
		{Code: v1review.LogicalOperatorsExceeded, Name: "OR"},
		{Code: v1review.LogicalOperatorsExceeded, Name: "IN", Text: "a IN (SELECT b FROM s)"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
