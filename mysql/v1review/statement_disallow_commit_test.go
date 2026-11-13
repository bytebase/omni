package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowCommit(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "commit", sql: "COMMIT;", want: []stmtFinding{{Code: v1review.Commit, Text: "COMMIT"}}},
		{name: "commit work and chain", sql: "INSERT INTO t VALUES (1);\nCOMMIT WORK AND CHAIN;", want: []stmtFinding{{Code: v1review.Commit, Text: "COMMIT WORK AND CHAIN"}}},
		{name: "rollback not judged", sql: "ROLLBACK;"},
		{name: "start transaction not judged", sql: "START TRANSACTION;"},
		{name: "commit in routine body not judged", sql: "CREATE PROCEDURE p() BEGIN COMMIT; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementDisallowCommit)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
