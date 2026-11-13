package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertDisallowOrderByRand(t *testing.T) {
	f := stmtFinding{Code: v1review.InsertOrderByRand}
	testStatementRule(t, StatementInsertDisallowOrderByRand, []statementCase{
		{name: "random", sql: "INSERT INTO t SELECT * FROM u ORDER BY random();", want: []stmtFinding{f}},
		{name: "random_between upper case", sql: "INSERT INTO t SELECT * FROM u ORDER BY a, RANDOM_BETWEEN(1, 2) DESC;", want: []stmtFinding{f}},
		{name: "any name part", sql: "INSERT INTO t SELECT * FROM u ORDER BY random.f();", want: []stmtFinding{f}},
		{name: "set operation branch", sql: "INSERT INTO t (SELECT a FROM u UNION (SELECT a FROM v ORDER BY random()));", want: []stmtFinding{f}},
		{name: "order by other function", sql: "INSERT INTO t SELECT * FROM u ORDER BY lower(a);"},
		{name: "random in select list", sql: "INSERT INTO t SELECT random() FROM u;"},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT * FROM (SELECT * FROM u ORDER BY random()) s;"},
		{name: "values", sql: "INSERT INTO t VALUES (1);"},
		{name: "select not judged", sql: "SELECT * FROM u ORDER BY random();"},
	})
}
