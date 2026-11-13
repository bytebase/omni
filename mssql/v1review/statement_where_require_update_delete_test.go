package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireUpdateDelete(t *testing.T) {
	missing := func(name, text string) stmtFinding {
		return stmtFinding{Code: v1review.WhereMissing, Name: name, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "update", sql: "UPDATE t SET a = 1", want: []stmtFinding{missing("UPDATE", "UPDATE t SET a = 1")}},
		{name: "update with where", sql: "UPDATE t SET a = 1 WHERE b = 2"},
		{name: "update from without where", sql: "UPDATE t SET a = u.a FROM u", want: []stmtFinding{missing("UPDATE", "UPDATE t SET a = u.a FROM u")}},
		{name: "delete", sql: "DELETE FROM t", want: []stmtFinding{missing("DELETE", "DELETE FROM t")}},
		{name: "delete with where", sql: "DELETE FROM t WHERE a = 1"},
		{name: "in a block not judged", sql: "BEGIN DELETE FROM t; END"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireUpdateDelete)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
