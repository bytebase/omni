package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementSelectNoSelectAll(t *testing.T) {
	top := stmtFinding{Code: v1review.SelectAll}
	nested := func(text string) stmtFinding { return stmtFinding{Code: v1review.SelectAll, Text: text} }
	testStatementRule(t, StatementSelectNoSelectAll, []statementCase{
		{name: "star", sql: "SELECT * FROM t;", want: []stmtFinding{top}},
		{name: "qualified star", sql: "SELECT a, t.* FROM t;", want: []stmtFinding{top}},
		{name: "one per select", sql: "SELECT *, * FROM t;", want: []stmtFinding{top}},
		{name: "columns", sql: "SELECT a FROM t;"},
		{name: "count star not judged", sql: "SELECT count(*) FROM t;"},
		{name: "subquery before outer", sql: "SELECT * FROM (SELECT * FROM t) s;", want: []stmtFinding{nested("SELECT * FROM t"), top}},
		{name: "where subquery", sql: "SELECT a FROM t WHERE a IN (SELECT * FROM u);", want: []stmtFinding{nested("SELECT * FROM u")}},
		{name: "cte", sql: "WITH x AS (SELECT * FROM t) SELECT a FROM x;", want: []stmtFinding{nested("SELECT * FROM t")}},
		{name: "set operation branches", sql: "SELECT * FROM t UNION SELECT a FROM u UNION ALL SELECT * FROM v;", want: []stmtFinding{nested("SELECT * FROM t"), nested("SELECT * FROM v")}},
		{name: "insert select", sql: "INSERT INTO t SELECT * FROM u;", want: []stmtFinding{nested("SELECT * FROM u")}},
		{name: "create view", sql: "CREATE VIEW v AS SELECT * FROM u;", want: []stmtFinding{nested("SELECT * FROM u")}},
		{name: "create table as", sql: "CREATE TABLE t AS SELECT * FROM u;", want: []stmtFinding{nested("SELECT * FROM u")}},
		{name: "case expression subquery", sql: "SELECT CASE WHEN (SELECT * FROM u) IS NULL THEN 1 END FROM t;", want: []stmtFinding{nested("SELECT * FROM u")}},
		{name: "subquery under other expression not judged", sql: "SELECT ARRAY[(SELECT * FROM u)] FROM t;"},
		{name: "array subquery", sql: "SELECT ARRAY(SELECT * FROM u) FROM t;", want: []stmtFinding{nested("SELECT * FROM u")}},
		{name: "explain not judged", sql: "EXPLAIN SELECT * FROM t;"},
		{name: "delete subquery not judged", sql: "DELETE FROM t WHERE a IN (SELECT * FROM u);"},
	})
}
