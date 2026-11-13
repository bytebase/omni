package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireSelect(t *testing.T) {
	top := stmtFinding{Code: v1review.WhereMissing}
	nested := func(text string) stmtFinding { return stmtFinding{Code: v1review.WhereMissing, Text: text} }
	testStatementRule(t, StatementWhereRequireSelect, []statementCase{
		{name: "no where", sql: "SELECT a FROM t;", want: []stmtFinding{top}},
		{name: "where", sql: "SELECT a FROM t WHERE a = 1;"},
		{name: "no from", sql: "SELECT 1;"},
		{name: "subquery before outer", sql: "SELECT a FROM (SELECT a FROM t) s;", want: []stmtFinding{nested("SELECT a FROM t"), top}},
		{name: "where subquery", sql: "SELECT a FROM t WHERE a IN (SELECT a FROM u);", want: []stmtFinding{nested("SELECT a FROM u")}},
		{name: "select list subquery", sql: "SELECT (SELECT max(a) FROM u) FROM t WHERE a = 1;", want: []stmtFinding{nested("SELECT max(a) FROM u")}},
		{name: "cte", sql: "WITH x AS (SELECT a FROM t) SELECT a FROM x WHERE a = 1;", want: []stmtFinding{nested("SELECT a FROM t")}},
		{name: "set operation branches", sql: "SELECT a FROM t UNION SELECT a FROM u WHERE a = 1;", want: []stmtFinding{nested("SELECT a FROM t")}},
		{name: "insert select", sql: "INSERT INTO t SELECT a FROM u;", want: []stmtFinding{nested("SELECT a FROM u")}},
		{name: "create view", sql: "CREATE VIEW v AS SELECT a FROM u;", want: []stmtFinding{nested("SELECT a FROM u")}},
		{name: "create materialized view", sql: "CREATE MATERIALIZED VIEW v AS SELECT a FROM u;", want: []stmtFinding{nested("SELECT a FROM u")}},
		{name: "join quals subquery", sql: "SELECT a FROM t JOIN u ON t.a IN (SELECT b FROM v) WHERE a = 1;", want: []stmtFinding{nested("SELECT b FROM v")}},
		{name: "subquery under other expression not judged", sql: "SELECT ARRAY[(SELECT a FROM u)] FROM t WHERE a = 1;"},
		{name: "explain not judged", sql: "EXPLAIN SELECT a FROM t;"},
		{name: "update subquery not judged", sql: "UPDATE t SET a = (SELECT a FROM u);"},
	})
}
