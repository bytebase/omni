package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowOnDelCascade(t *testing.T) {
	f := stmtFinding{Code: v1review.OnDeleteCascade}
	testStatementRule(t, StatementDisallowOnDelCascade, []statementCase{
		{name: "table constraint", sql: "CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES u (id) ON DELETE CASCADE);", want: []stmtFinding{f}},
		{name: "column constraint", sql: "CREATE TABLE t (a int REFERENCES u (id) ON DELETE CASCADE);", want: []stmtFinding{f}},
		{name: "one per statement", sql: "CREATE TABLE t (a int REFERENCES u (id) ON DELETE CASCADE, b int REFERENCES v (id) ON DELETE CASCADE);", want: []stmtFinding{f}},
		{name: "alter add constraint", sql: "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (id) ON DELETE CASCADE;", want: []stmtFinding{f}},
		{name: "alter added column not judged", sql: "ALTER TABLE t ADD COLUMN a int REFERENCES u (id) ON DELETE CASCADE;"},
		{name: "on update cascade not judged", sql: "CREATE TABLE t (a int REFERENCES u (id) ON UPDATE CASCADE);"},
		{name: "on delete set null not judged", sql: "CREATE TABLE t (a int REFERENCES u (id) ON DELETE SET NULL);"},
	})
}
