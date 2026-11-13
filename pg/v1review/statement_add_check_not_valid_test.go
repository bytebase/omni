package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementAddCheckNotValid(t *testing.T) {
	f := stmtFinding{Code: v1review.CheckValidated}
	testStatementRule(t, StatementAddCheckNotValid, []statementCase{
		{name: "check", sql: "ALTER TABLE t ADD CONSTRAINT c CHECK (a > 0);", want: []stmtFinding{f}},
		{name: "not valid", sql: "ALTER TABLE t ADD CONSTRAINT c CHECK (a > 0) NOT VALID;"},
		{name: "one per command", sql: "ALTER TABLE t ADD CHECK (a > 0), ADD CHECK (b > 0) NOT VALID, ADD CHECK (c > 0);", want: []stmtFinding{f, f}},
		{name: "foreign key not judged", sql: "ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES u (id);"},
		{name: "check on added column not judged", sql: "ALTER TABLE t ADD COLUMN a int CHECK (a > 0);"},
		{name: "create table not judged", sql: "CREATE TABLE t (a int, CHECK (a > 0));"},
		{name: "each statement", sql: "ALTER TABLE t ADD CHECK (a > 0);\nALTER TABLE u ADD CHECK (b > 0);", want: []stmtFinding{f, f}},
	})
}
