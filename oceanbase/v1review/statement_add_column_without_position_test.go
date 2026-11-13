package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementAddColumnWithoutPosition(t *testing.T) {
	positioned := func(keyword string) stmtFinding { return stmtFinding{Code: v1review.ColumnPositioned, Name: keyword} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "first", sql: "ALTER TABLE t ADD COLUMN a INT FIRST;", want: []stmtFinding{positioned("FIRST")}},
		{name: "after", sql: "ALTER TABLE t ADD a INT AFTER b;", want: []stmtFinding{positioned("AFTER")}},
		{name: "each command", sql: "ALTER TABLE t ADD a INT FIRST, ADD b INT, ADD c INT AFTER a;", want: []stmtFinding{positioned("FIRST"), positioned("AFTER")}},
		{name: "no position", sql: "ALTER TABLE t ADD COLUMN a INT;"},
		{name: "modify not judged", sql: "ALTER TABLE t MODIFY COLUMN a INT FIRST;"},
		{name: "change not judged", sql: "ALTER TABLE t CHANGE a b INT AFTER c;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE t ADD a INT FIRST; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementAddColumnWithoutPosition)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
