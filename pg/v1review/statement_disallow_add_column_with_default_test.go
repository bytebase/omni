package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowAddColumnWithDefault(t *testing.T) {
	f := stmtFinding{Code: v1review.ColumnWithDefault}
	testStatementRule(t, StatementDisallowAddColumnWithDefault, []statementCase{
		{name: "default", sql: "ALTER TABLE t ADD COLUMN a int DEFAULT 0;", want: []stmtFinding{f}},
		{name: "volatile default", sql: "ALTER TABLE t ADD COLUMN a timestamptz DEFAULT now();", want: []stmtFinding{f}},
		{name: "one per statement", sql: "ALTER TABLE t ADD COLUMN a int DEFAULT 0, ADD COLUMN b int DEFAULT 1;", want: []stmtFinding{f}},
		{name: "no default", sql: "ALTER TABLE t ADD COLUMN a int NOT NULL;"},
		{name: "set default not judged", sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT 0;"},
		{name: "create table not judged", sql: "CREATE TABLE t (a int DEFAULT 0);"},
		{name: "each statement", sql: "ALTER TABLE t ADD a int DEFAULT 0;\nALTER TABLE u ADD b int DEFAULT 0;", want: []stmtFinding{f, f}},
	})
}
