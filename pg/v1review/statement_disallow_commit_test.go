package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowCommit(t *testing.T) {
	f := stmtFinding{Code: v1review.Commit}
	testStatementRule(t, StatementDisallowCommit, []statementCase{
		{name: "commit", sql: "COMMIT;", want: []stmtFinding{f}},
		{name: "end", sql: "END;", want: []stmtFinding{f}},
		{name: "commit and chain", sql: "BEGIN;\nCOMMIT AND CHAIN;", want: []stmtFinding{f}},
		{name: "commit prepared not judged", sql: "COMMIT PREPARED 'x';"},
		{name: "rollback not judged", sql: "ROLLBACK;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() LANGUAGE plpgsql AS $$ BEGIN COMMIT; END $$;"},
	})
}
