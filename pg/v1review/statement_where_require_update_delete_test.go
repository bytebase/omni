package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireUpdateDelete(t *testing.T) {
	update := stmtFinding{Code: v1review.WhereMissing, Name: "UPDATE"}
	del := stmtFinding{Code: v1review.WhereMissing, Name: "DELETE"}
	testStatementRule(t, StatementWhereRequireUpdateDelete, []statementCase{
		{name: "update", sql: "UPDATE t SET a = 1;", want: []stmtFinding{update}},
		{name: "delete", sql: "DELETE FROM t;", want: []stmtFinding{del}},
		{name: "with where", sql: "UPDATE t SET a = 1 WHERE b = 2;\nDELETE FROM t WHERE b = 2;"},
		{name: "where current of", sql: "DELETE FROM t WHERE CURRENT OF c;"},
		{name: "with query not judged", sql: "WITH x AS (DELETE FROM t RETURNING a) SELECT * FROM x;"},
		{name: "each statement", sql: "UPDATE t SET a = 1;\nDELETE FROM t;", want: []stmtFinding{update, del}},
	})
}
