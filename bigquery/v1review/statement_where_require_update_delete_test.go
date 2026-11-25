package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireUpdateDelete(t *testing.T) {
	missing := func(kind, text string) stmtFinding {
		return stmtFinding{Code: v1review.WhereMissing, Name: kind, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "update", sql: "UPDATE `proj.ds.t` SET a = 1;", want: []stmtFinding{missing("UPDATE", "UPDATE `proj.ds.t` SET a = 1")}},
		{name: "delete without FROM", sql: "SELECT 1;\nDELETE ds.t;", want: []stmtFinding{missing("DELETE", "DELETE ds.t")}},
		{name: "where true is a where", sql: "UPDATE ds.t SET a = 1 WHERE TRUE;"},
		{name: "delete with where", sql: "DELETE FROM ds.t WHERE a = 1;"},
		{
			name: "nested in script block",
			sql:  "BEGIN\n  UPDATE ds.t SET a = 1;\n  IF TRUE THEN DELETE FROM ds.t; END IF;\nEND;",
			want: []stmtFinding{missing("UPDATE", "UPDATE ds.t SET a = 1"), missing("DELETE", "DELETE FROM ds.t")},
		},
		{
			name: "merge actions not judged",
			sql:  "MERGE ds.t AS t USING ds.s AS s ON t.id = s.id\nWHEN MATCHED THEN UPDATE SET a = s.a\nWHEN NOT MATCHED BY SOURCE THEN DELETE;",
		},
		{name: "select not judged", sql: "SELECT * FROM ds.t;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireUpdateDelete)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
