package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementWhereRequireSelect(t *testing.T) {
	missing := func(text string) stmtFinding { return stmtFinding{Code: v1review.WhereMissing, Text: text} }
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "no where", sql: "SELECT a FROM t;", want: []stmtFinding{missing("SELECT a FROM t")}},
		{name: "where", sql: "SELECT a FROM t WHERE b = 1;"},
		{name: "no from", sql: "SELECT 1;"},
		{name: "subquery in where", sql: "SELECT a FROM t WHERE b IN (SELECT c FROM u);", want: []stmtFinding{missing("SELECT c FROM u")}},
		{name: "subquery in from", sql: "SELECT a FROM (SELECT b FROM u) x WHERE a = 1;", want: []stmtFinding{missing("SELECT b FROM u")}},
		{name: "subquery in select list", sql: "SELECT (SELECT b FROM u) FROM t WHERE a = 1;", want: []stmtFinding{missing("SELECT b FROM u")}},
		{name: "cte before owner", sql: "WITH c AS (SELECT a FROM u) SELECT a FROM c;", want: []stmtFinding{missing("SELECT a FROM u"), missing("WITH c AS (SELECT a FROM u) SELECT a FROM c")}},
		{name: "set operation sides", sql: "SELECT a FROM t UNION SELECT b FROM u WHERE c = 1;", want: []stmtFinding{missing("SELECT a FROM t")}},
		{name: "insert select not judged", sql: "INSERT INTO t SELECT a FROM u;"},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN SELECT a FROM t; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementWhereRequireSelect)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
