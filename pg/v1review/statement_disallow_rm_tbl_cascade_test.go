package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowRmTblCascade(t *testing.T) {
	f := stmtFinding{Code: v1review.DropCascade}
	testStatementRule(t, StatementDisallowRmTblCascade, []statementCase{
		{name: "drop table cascade", sql: "DROP TABLE t, u CASCADE;", want: []stmtFinding{f}},
		{name: "truncate cascade", sql: "TRUNCATE t CASCADE;", want: []stmtFinding{f}},
		{name: "drop table restrict", sql: "DROP TABLE t RESTRICT;"},
		{name: "truncate", sql: "TRUNCATE t;"},
		{name: "drop view cascade not judged", sql: "DROP VIEW v CASCADE;"},
		{name: "drop schema cascade not judged", sql: "DROP SCHEMA s CASCADE;"},
	})
}
