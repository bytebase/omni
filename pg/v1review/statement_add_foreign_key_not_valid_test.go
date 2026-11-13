package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementAddForeignKeyNotValid(t *testing.T) {
	f := stmtFinding{Code: v1review.ForeignKeyValidated}
	testStatementRule(t, StatementAddForeignKeyNotValid, []statementCase{
		{name: "foreign key", sql: "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (id);", want: []stmtFinding{f}},
		{name: "not valid", sql: "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (id) NOT VALID;"},
		{name: "one per command", sql: "ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES u (id), ADD FOREIGN KEY (b) REFERENCES v (id);", want: []stmtFinding{f, f}},
		{name: "check not judged", sql: "ALTER TABLE t ADD CHECK (a > 0);"},
		{name: "reference on added column not judged", sql: "ALTER TABLE t ADD COLUMN a int REFERENCES u (id);"},
		{name: "create table not judged", sql: "CREATE TABLE t (a int REFERENCES u (id));"},
	})
}
