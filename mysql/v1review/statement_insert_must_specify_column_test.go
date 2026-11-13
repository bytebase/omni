package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertMustSpecifyColumn(t *testing.T) {
	unspecified := []stmtFinding{{Code: v1review.InsertColumnsUnspecified}}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "insert values", sql: "INSERT INTO t VALUES (1), (2);", want: unspecified},
		{name: "insert select", sql: "INSERT INTO t SELECT * FROM s;", want: unspecified},
		{name: "replace values", sql: "REPLACE INTO t VALUES (1);", want: unspecified},
		{name: "columns given", sql: "INSERT INTO t (a) VALUES (1);"},
		{name: "columns given to select", sql: "INSERT INTO t (a) SELECT a FROM s;"},
		{name: "insert set not judged", sql: "INSERT INTO t SET a = 1;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN INSERT INTO t VALUES (1); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementInsertMustSpecifyColumn)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
