package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireUpdateDelete(t *testing.T) {
	missing := func(name string, start int, text string) stmtFinding {
		return stmtFinding{Code: v1review.WhereMissing, Name: name, Start: start, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "update", sql: "UPDATE t SET a = 1", want: []stmtFinding{missing("UPDATE", 0, "UPDATE t SET a = 1")}},
		{name: "update with where", sql: "UPDATE t SET a = 1 WHERE b = 2"},
		{name: "delete", sql: "DELETE FROM t", want: []stmtFinding{missing("DELETE", 0, "DELETE FROM t")}},
		{name: "delete with where", sql: "DELETE FROM t WHERE a = 1"},
		{name: "merge actions not judged", sql: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET t.a = s.a WHEN NOT MATCHED THEN INSERT (a) VALUES (s.a)"},
		{name: "select not judged", sql: "SELECT a FROM t"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireUpdateDelete)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
