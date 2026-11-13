package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireSelect(t *testing.T) {
	missing := func(start int, text string) statementFinding {
		return statementFinding{Code: v1review.WhereMissing, Start: start, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []statementFinding
	}{
		{name: "no where", sql: "SELECT a FROM t", want: []statementFinding{missing(0, "SELECT a FROM t")}},
		{name: "where", sql: "SELECT a FROM t WHERE b = 1"},
		{name: "no from not judged", sql: "SELECT 1"},
		{
			name: "derived table",
			sql:  "SELECT a FROM (SELECT b FROM u) x WHERE a = 1",
			want: []statementFinding{missing(15, "SELECT b FROM u")},
		},
		{
			name: "with select anchors the body select",
			sql:  "WITH c AS (SELECT x FROM t WHERE y = 1), d AS (SELECT 1)\n  /* body */ -- next\n  SELECT a FROM c",
			want: []statementFinding{missing(80, "SELECT a FROM c")},
		},
		{
			name: "cte body before the body select",
			sql:  "WITH c AS (SELECT x FROM t) SELECT a FROM c",
			want: []statementFinding{missing(11, "SELECT x FROM t"), missing(28, "SELECT a FROM c")},
		},
		{
			name: "merge source",
			sql:  "MERGE INTO t USING (SELECT id FROM s) AS u ON t.id = u.id WHEN MATCHED THEN DELETE",
			want: []statementFinding{missing(20, "SELECT id FROM s")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, func(stmts []ast.Node) []v1review.Finding { return StatementWhereRequireSelect(stmts, tt.sql) })); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
