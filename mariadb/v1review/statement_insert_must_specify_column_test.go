package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertMustSpecifyColumn(t *testing.T) {
	unspecified := stmtFinding{Code: v1review.InsertColumnsUnspecified}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "values", sql: "INSERT INTO t VALUES (1), (2);", want: []stmtFinding{unspecified}},
		{name: "select", sql: "INSERT INTO t SELECT * FROM u;", want: []stmtFinding{unspecified}},
		{name: "replace values", sql: "REPLACE INTO t VALUES (1);", want: []stmtFinding{unspecified}},
		{name: "columns named", sql: "INSERT INTO t (a) VALUES (1); INSERT INTO t (a) SELECT a FROM u;"},
		{name: "set form not judged", sql: "INSERT INTO t SET a = 1;"},
		{name: "insert in routine body not judged", sql: "CREATE PROCEDURE p() BEGIN INSERT INTO t VALUES (1); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementInsertMustSpecifyColumn)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
