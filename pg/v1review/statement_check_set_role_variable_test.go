package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementCheckSetRoleVariable(t *testing.T) {
	missing := []stmtFinding{{Code: v1review.SetRoleMissing}}
	testStatementRule(t, StatementCheckSetRoleVariable, []statementCase{
		{name: "set role first", sql: "SET ROLE app;\nCREATE TABLE t (a int);"},
		{name: "set role after other sets", sql: "SET statement_timeout = 0;\nSET ROLE app;\nCREATE TABLE t (a int);"},
		{name: "set role after a statement", sql: "CREATE TABLE t (a int);\nSET ROLE app;", want: missing},
		{name: "no set role", sql: "CREATE TABLE t (a int);", want: missing},
		{name: "only other sets", sql: "SET statement_timeout = 0;", want: missing},
		{name: "set local role", sql: "SET LOCAL ROLE app;\nSELECT 1;"},
		{name: "reset role", sql: "RESET ROLE;\nSELECT 1;"},
		{name: "upper-case name", sql: "SET \"ROLE\" = app;\nSELECT 1;"},
		{name: "search_path counts", sql: "SET search_path = s;\nSELECT 1;"},
		{name: "reset all counts", sql: "RESET ALL;\nSELECT 1;"},
		{name: "set session authorization does not count", sql: "SET SESSION AUTHORIZATION u;\nSELECT 1;", want: missing},
		{name: "no statements", sql: "", want: missing},
	})
}
