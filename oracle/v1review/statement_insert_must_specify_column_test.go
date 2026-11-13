package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertMustSpecifyColumn(t *testing.T) {
	missing := func(text string) stmtFinding {
		return stmtFinding{Code: v1review.InsertColumnsUnspecified, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "values", sql: "INSERT INTO t VALUES (1)", want: []stmtFinding{missing("INSERT INTO t VALUES (1)")}},
		{name: "columns", sql: "INSERT INTO t (a) VALUES (1)"},
		{name: "select", sql: "INSERT INTO t SELECT a FROM u", want: []stmtFinding{missing("INSERT INTO t SELECT a FROM u")}},
		{
			name: "insert all anchors each into clause",
			sql:  "INSERT ALL INTO t VALUES (1) INTO u (a) VALUES (2) SELECT 1 FROM dual",
			want: []stmtFinding{missing("INTO t VALUES (1)")},
		},
		{
			name: "anonymous block",
			sql:  "BEGIN\n  INSERT INTO t VALUES (1);\nEND;",
			want: []stmtFinding{missing("INSERT INTO t VALUES (1)")},
		},
		{name: "procedure body not judged", sql: "CREATE PROCEDURE p AS BEGIN INSERT INTO t VALUES (1); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementInsertMustSpecifyColumn)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
