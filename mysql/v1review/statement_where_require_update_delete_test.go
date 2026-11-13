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
		{name: "delete", sql: "DELETE FROM t;", want: []stmtFinding{{Code: v1review.WhereMissing, Name: "DELETE"}}},
		{name: "where", sql: "DELETE FROM t WHERE a = 1;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN DELETE FROM t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireUpdateDelete)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
