package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementCreateSpecifySchema(t *testing.T) {
	f := stmtFinding{Code: v1review.SchemaUnspecified, Name: "Table"}
	testStatementRule(t, StatementCreateSpecifySchema, []statementCase{
		{name: "unqualified", sql: "CREATE TABLE t (a int);", want: []stmtFinding{f}},
		{name: "qualified", sql: "CREATE TABLE s.t (a int);"},
		{name: "partition", sql: "CREATE TABLE p PARTITION OF s.t FOR VALUES IN (1);", want: []stmtFinding{f}},
		{name: "create table as not judged", sql: "CREATE TABLE t AS SELECT 1;"},
		{name: "view not judged", sql: "CREATE VIEW v AS SELECT 1;"},
		{name: "index not judged", sql: "CREATE INDEX i ON t (a);"},
		{name: "each statement", sql: "CREATE TABLE t (a int);\nCREATE TABLE s.u (a int);\nCREATE TABLE v (a int);", want: []stmtFinding{f, f}},
	})
}
