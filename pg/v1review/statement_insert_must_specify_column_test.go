package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertMustSpecifyColumn(t *testing.T) {
	f := stmtFinding{Code: v1review.InsertColumnsUnspecified}
	testStatementRule(t, StatementInsertMustSpecifyColumn, []statementCase{
		{name: "values", sql: "INSERT INTO t VALUES (1);", want: []stmtFinding{f}},
		{name: "select", sql: "INSERT INTO t SELECT * FROM u;", want: []stmtFinding{f}},
		{name: "default values", sql: "INSERT INTO t DEFAULT VALUES;", want: []stmtFinding{f}},
		{name: "columns", sql: "INSERT INTO t (a) VALUES (1);"},
		{name: "with query not judged", sql: "WITH x AS (INSERT INTO t VALUES (1) RETURNING a) SELECT * FROM x;"},
		{name: "each statement", sql: "INSERT INTO t VALUES (1);\nINSERT INTO t (a) VALUES (1);\nINSERT INTO t VALUES (2);", want: []stmtFinding{f, f}},
	})
}
