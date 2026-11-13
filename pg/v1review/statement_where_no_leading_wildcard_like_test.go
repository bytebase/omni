package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereNoLeadingWildcardLike(t *testing.T) {
	f := stmtFinding{Code: v1review.LeadingWildcardLike}
	testStatementRule(t, StatementWhereNoLeadingWildcardLike, []statementCase{
		{name: "percent", sql: "SELECT * FROM t WHERE a LIKE '%x';", want: []stmtFinding{f}},
		{name: "underscore", sql: "SELECT * FROM t WHERE a LIKE '_x';", want: []stmtFinding{f}},
		{name: "ilike", sql: "DELETE FROM t WHERE a ILIKE '%x';", want: []stmtFinding{f}},
		{name: "not like", sql: "UPDATE t SET b = 1 WHERE a NOT LIKE '%x';", want: []stmtFinding{f}},
		{name: "one per statement", sql: "SELECT * FROM t WHERE a LIKE '%x' OR b LIKE '%y';", want: []stmtFinding{f}},
		{name: "select list", sql: "SELECT a LIKE '%x' FROM t;", want: []stmtFinding{f}},
		{name: "subquery", sql: "SELECT * FROM t WHERE a IN (SELECT a FROM u WHERE b LIKE '%x');", want: []stmtFinding{f}},
		{name: "trailing wildcard", sql: "SELECT * FROM t WHERE a LIKE 'x%';"},
		{name: "similar to not judged", sql: "SELECT * FROM t WHERE a SIMILAR TO '%x';"},
		{name: "non-constant pattern not judged", sql: "SELECT * FROM t WHERE a LIKE '%' || b;"},
		{name: "each statement", sql: "SELECT * FROM t WHERE a LIKE '%x';\nSELECT 1;\nSELECT * FROM t WHERE a LIKE '_x';", want: []stmtFinding{f, f}},
	})
}
