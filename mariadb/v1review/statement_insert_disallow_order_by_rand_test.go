package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementInsertDisallowOrderByRand(t *testing.T) {
	rand := stmtFinding{Code: v1review.InsertOrderByRand}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "order by rand", sql: "INSERT INTO t SELECT a FROM u ORDER BY RAND();", want: []stmtFinding{rand}},
		{name: "each rand item", sql: "INSERT INTO t SELECT a FROM u ORDER BY rand(1), a, RAND();", want: []stmtFinding{rand, rand}},
		{name: "qualified rand", sql: "INSERT INTO t SELECT a FROM u ORDER BY db.RAND();", want: []stmtFinding{rand}},
		{name: "parenthesized select", sql: "INSERT INTO t (SELECT a FROM u ORDER BY RAND() LIMIT 1) ORDER BY RAND() LIMIT 2;", want: []stmtFinding{rand, rand}},
		{name: "set operation", sql: "INSERT INTO t SELECT a FROM u UNION (SELECT b FROM v ORDER BY RAND() LIMIT 1) ORDER BY RAND();", want: []stmtFinding{rand, rand}},
		{name: "rand in expression not judged", sql: "INSERT INTO t SELECT a FROM u ORDER BY RAND() + 1;"},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT * FROM u WHERE a IN (SELECT a FROM v ORDER BY RAND());"},
		{name: "values not judged", sql: "INSERT INTO t VALUES (RAND());"},
		{name: "select not judged", sql: "SELECT a FROM u ORDER BY RAND();"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementInsertDisallowOrderByRand)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
