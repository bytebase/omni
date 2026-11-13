package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereNoEqualNull(t *testing.T) {
	eq := func(text string) stmtFinding { return stmtFinding{Code: v1review.EqualNull, Text: text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "equal null", sql: "SELECT * FROM t WHERE a = NULL;", want: []stmtFinding{eq("a = NULL")}},
		{name: "not equal null", sql: "SELECT * FROM t WHERE a = 1 AND b <> NULL;", want: []stmtFinding{eq("b <> NULL")}},
		{name: "null on the left not judged", sql: "SELECT * FROM t WHERE NULL = a;"},
		{name: "is null not judged", sql: "SELECT * FROM t WHERE a IS NULL;"},
		{name: "update where not judged", sql: "UPDATE t SET a = 1 WHERE b = NULL;"},
		{name: "subquery in update judged", sql: "UPDATE t SET a = 1 WHERE b IN (SELECT b FROM s WHERE c = NULL);", want: []stmtFinding{eq("c = NULL")}},
		{
			name: "subquery in a where reported twice",
			sql:  "SELECT * FROM t WHERE a IN (SELECT a FROM s WHERE b = NULL);",
			want: []stmtFinding{eq("b = NULL"), eq("b = NULL")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereNoEqualNull)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
