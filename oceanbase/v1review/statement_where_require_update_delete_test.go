package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireUpdateDelete(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "update", sql: "UPDATE t SET a = 1;", want: []stmtFinding{{Code: v1review.WhereMissing, Name: "UPDATE"}}},
		{name: "multi-table delete", sql: "DELETE t1 FROM t1 JOIN t2 ON t1.a = t2.a;", want: []stmtFinding{{Code: v1review.WhereMissing, Name: "DELETE"}}},
		{name: "where", sql: "UPDATE t SET a = 1 WHERE b = 2; DELETE FROM t WHERE a = 1;"},
		{name: "select not judged", sql: "SELECT a FROM t;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN DELETE FROM t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireUpdateDelete)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
